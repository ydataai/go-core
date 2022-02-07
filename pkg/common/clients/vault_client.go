package clients

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
)

// VaultClient defines the Vault client struct, holding all the required dependencies
type VaultClient struct {
	configuration VaultClientConfiguration
	role          string
	logger        logging.Logger
	client        *api.Client
	secret        *api.Secret
}

// NewVaultClient returns an initialized struct with the required dependencies injected
func NewVaultClient(role string, configuration VaultClientConfiguration, logger logging.Logger) (*VaultClient, error) {
	config := &api.Config{Address: configuration.VaultURL}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	vc := &VaultClient{
		configuration: configuration,
		role:          role,
		logger:        logger,
		client:        client,
	}

	if err = vc.login(); err != nil {
		return nil, err
	}

	return vc, nil
}

// login the k8s service account
func (vc *VaultClient) login() error {
	// Make sure it's development mode
	if vc.configuration.VaultDevMode {
		vc.logger.Info("[development mode] performing vault login.")
		if vc.configuration.VaultToken == "" {
			return errors.New("unable to get token as env var ('VAULT_TOKEN') 😱")
		}
		vc.client.SetToken(vc.configuration.VaultToken)
		return nil
	}

	// Production mode
	vc.logger.Info("performing vault k8s login.")
	// reads jwt from service account
	jwt, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return fmt.Errorf("unable to read file containing service account token: %v 😱", err)
	}
	params := map[string]interface{}{
		"jwt":  string(jwt),
		"role": vc.role, // the name of the role in Vault that was created with this app's Kubernetes service account bound to it
	}
	// perform login
	secret, err := vc.client.Logical().Write("auth/kubernetes/login", params)
	if err != nil {
		return fmt.Errorf("unable to log in with Kubernetes auth: %v 😱", err)
	}
	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		return errors.New("login response did not return client token 😱")
	}
	// client update with the access token
	vc.logger.Info("login: client logged in successfully 🔑")
	token := strings.TrimSuffix(secret.Auth.ClientToken, "\n")
	vc.client.SetToken(token)
	// stores login response secret
	vc.secret = secret

	// do the token renewal cycle
	go vc.renew()

	return nil
}

// renew the token according to secret.Auth.LeaseDuration automatically
func (vc *VaultClient) renew() {
	vc.logger.Info("stating vault token auto renew ...")
	// schedule the token renew operation
	for range time.Tick(time.Second * time.Duration(vc.secret.Auth.LeaseDuration-(vc.secret.Auth.LeaseDuration/10))) {
		// perform renew
		resp, err := vc.client.Auth().Token().Renew(vc.secret.Auth.ClientToken, vc.secret.Auth.LeaseDuration)
		if err != nil {
			vc.logger.Errorf("unable to renew the access token %v 😱", err)
		}
		// client update with the renewed token
		if resp != nil && resp.Auth != nil && resp.Auth.ClientToken != "" {
			token := strings.TrimSuffix(resp.Auth.ClientToken, "\n")
			vc.logger.Info("renew: client token renewed successfully 🔑")
			vc.client.SetToken(token)
		} else {
			// new login to deal with system token expiration
			vc.login()
		}
	}
}

// StoreCredentials receives the path and the respective map of credentials and attempts to store them
// on the Vault server.
func (vc *VaultClient) StoreCredentials(path string, credentials map[string]string) error {
	vc.logger.Info("Sending credentials to Vault ☄️")

	_, err := vc.client.Logical().Write(path, map[string]interface{}{
		"data": credentials,
	})
	if err != nil {
		vc.logger.Errorf("Unable to store credentials in Vault 😱. Err: %v ", err)
		return err
	}

	vc.logger.Info("Credentials safely secured in Vault 🔑")
	return nil
}

// GetCredentials receives the path and attemps to retrieve the map of credentials present
// on the Vault server.
func (vc *VaultClient) GetCredentials(path string) (*config.Credentials, error) {
	vc.logger.Info("Fetching credentials from Vault ☄️")

	secret, err := vc.client.Logical().Read(path)
	if err != nil {
		vc.logger.Errorf("Unable to fetch credentials from Vault 😱. Err: %v", err)
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	vc.logger.Info("Credentials safely retrieved from Vault 🔑")

	secretsMap, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		vc.logger.Errorf("Unable to decipher received credentials from Vault 😱. Err: %v", err)
		return nil, err
	}

	vc.logger.Info("Processing credentials map 🔎")
	credentials := config.Credentials{}

	for key, value := range secretsMap {
		credentials[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", value)
	}

	return &credentials, nil
}

// DeleteCredentials receives the path and attempts to delete the existing credentials on Vault.
// Is performs a soft delete, per docs > https://www.vaultproject.io/docs/commands/kv/delete
func (vc *VaultClient) DeleteCredentials(path string) error {
	vc.logger.Info("Deleting credentials from Vault ☄️")

	_, err := vc.client.Logical().Delete(path)
	if err != nil {
		vc.logger.Errorf("Unable to delete credentials from Vault 😱. Err: %v", err)
		return err
	}

	vc.logger.Info("Credentials deleted from Vault ☠️")
	return nil
}

// CheckIfEngineExists attempts to call the /tune API endpoint on the Secrets Engine. Should it fail, it might be an
// indication that the Secrets Engine is not created, which it's useful to know whether or not to call CreateEngine
func (vc *VaultClient) CheckIfEngineExists(path string) bool {
	vc.logger.Info("Checking if vault engine exists☄️")

	epath := fmt.Sprintf("sys/mounts/%s/tune", path)

	if _, err := vc.client.Logical().Read(epath); err != nil {
		switch err.(type) {
		case *api.ResponseError:
			vc.logger.Infof("%v Secrets engine seems to be non existing 🤔. Err: %v", epath, err)
			return false
		default:
			vc.logger.Errorf("An error occurred fetching %v Secrets Engine 😵. Err: %v", epath, err)
			return false
		}
	}
	return true
}

// List ...
func (vc *VaultClient) List(path string) (interface{}, error) {
	vc.logger.Infof("[Vault] Listing the path: '%s' ☄️", path)

	secret, err := vc.client.Logical().List(path)
	if err != nil {
		vc.logger.Errorf("[Vault] Unable to list the path: '%s' 😱. Err: %v", path, err)
		return nil, err
	}

	if secret == nil {
		vc.logger.Infof("[Vault] ❌ No data found in path: '%s'", path)
		return nil, nil
	}

	vc.logger.Infof("[Vault] Listed the path: '%s' ☄️", path)
	return secret.Data["keys"], nil
}

// Get ...
func (vc *VaultClient) Get(path string) (map[string]interface{}, error) {
	vc.logger.Infof("[Vault] Getting the '%s' ☄️", path)

	secret, err := vc.client.Logical().Read(path)
	if err != nil {
		vc.logger.Errorf("[Vault] Unable to get '%s' 😱. Err: %v", path, err)
		return nil, err
	}

	if secret == nil {
		vc.logger.Infof("[Vault] ❌ No data found: %s", path)
		return nil, nil
	}

	vc.logger.Infof("[Vault] Got the '%s' ☄️", path)
	return secret.Data, nil
}

// Delete ...
func (vc *VaultClient) Delete(path string) error {
	vc.logger.Infof("[Vault] Deleting the path: '%s'", path)

	secret, err := vc.client.Logical().Read(path)
	if err != nil {
		return fmt.Errorf("[Vault] Unable to delete the path: '%s' 😱. Err: %v", path, err)
	}

	if secret == nil {
		vc.logger.Infof("[Vault] ❌ No data found in path: '%s'", path)
		return nil
	}

	_, err = vc.client.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("[Vault] Unable to delete the path: '%s' 😱. Err: %v", path, err)
	}

	vc.logger.Infof("[Vault] Deleted the path: '%s' ☄️", path)
	return nil
}

// Put ...
func (vc *VaultClient) Put(path string, data map[string]interface{}) error {
	vc.logger.Infof("[Vault] Creating the '%s' ☄️", path)

	_, err := vc.client.Logical().Write(path, data)
	if err != nil {
		return fmt.Errorf("[Vault] Unable to create '%s' 😱. Err: %v", path, err)
	}

	vc.logger.Infof("[Vault] Created the '%s' ☄️", path)
	return nil
}

// Patch ...
func (vc *VaultClient) Patch(path string, data map[string]interface{}) error {
	vc.logger.Infof("[Vault] Patch the '%s' ☄️", path)
	// try to patch the path
	_, err := vc.client.Logical().JSONMergePatch(context.Background(), path, data)
	if err == nil {
		return nil
	}
	// If it's a 405, that probably means the server is running a pre-1.9
	// Vault version that doesn't support the HTTP PATCH method.
	if re, ok := err.(*api.ResponseError); ok && re.StatusCode != 405 {
		return fmt.Errorf("[Vault] Unable to add the path: '%s' 😱. Err: %v", path, err)
	}
	// get data to update it in memory
	existingData, err := vc.Get(path)
	if err != nil {
		return fmt.Errorf("[Vault] Unable to get the path: '%s' 😱. Err: %v", path, err)
	}
	// if exists data, then update
	if existingData != nil {
		// if it exists, then update
		for key, value := range data {
			existingData[key] = value
		}
		if err := vc.Put(path, existingData); err != nil {
			return err
		}
		return nil
	}
	// it doesn't exists, create
	if err := vc.Put(path, data); err != nil {
		return err
	}
	return nil
}
