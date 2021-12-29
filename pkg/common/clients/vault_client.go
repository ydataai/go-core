package clients

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
)

const path = "datasources"

// VaultClient defines the Vault client struct, holding all the required dependencies
type VaultClient struct {
	configuration VaultClientConfiguration
	logger        logging.Logger
	client        *api.Client
	secret        *api.Secret
}

// NewVaultClient returns an initialized struct with the required dependencies injected
func NewVaultClient(configuration VaultClientConfiguration, logger logging.Logger) (*VaultClient, error) {
	config := &api.Config{Address: configuration.VaultURL}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	vc := &VaultClient{
		configuration: configuration,
		logger:        logger,
		client:        client,
	}

	if err = vc.login(); err != nil {
		return nil, err
	}

	go vc.renew()

	return vc, nil
}

// login the k8s service account
func (vc *VaultClient) login() error {
	vc.logger.Info("performing vault k8s login.")
	// reads jwt from service account
	jwt, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return fmt.Errorf("unable to read file containing service account token: %v 😱", err)
	}
	params := map[string]interface{}{
		"jwt":  string(jwt),
		"role": "datasource-controller-role", // the name of the role in Vault that was created with this app's Kubernetes service account bound to it
	}
	// perform login
	resp, err := vc.client.Logical().Write("auth/kubernetes/login", params)
	if err != nil {
		return fmt.Errorf("unable to log in with Kubernetes auth: %v 😱", err)
	}
	if resp == nil || resp.Auth == nil || resp.Auth.ClientToken == "" {
		return errors.New("login response did not return client token 😱")
	}
	// client update with the access token
	vc.logger.Info("login: client logged in successfully 🔑")
	token := strings.TrimSuffix(resp.Auth.ClientToken, "\n")
	vc.client.SetToken(token)
	// stores login response secret
	vc.secret = resp
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

// StoreCredentials receives the name of the DataSource and the respective map of credentials and attempts to store them
// on the Vault server.
func (vc *VaultClient) StoreCredentials(datasourceName string, credentials map[string]string) error {
	vc.logger.Info("Sending credentials to Vault ☄️")

	_, err := vc.client.Logical().Write(fmt.Sprintf("%s/data/%s", path, datasourceName), map[string]interface{}{
		"data": credentials,
	})
	if err != nil {
		vc.logger.Errorf("Unable to store credentials in Vault 😱. Err: %v ", err)
		return err
	}

	vc.logger.Info("Credentials safely secured in Vault 🔑")
	return nil
}

// GetCredentials receives the name of the DataSource and attemps to retrieve the map of credentials present
// on the Vault server.
func (vc *VaultClient) GetCredentials(datasourceName string) (*config.Credentials, error) {
	vc.logger.Info("Fetching credentials from Vault ☄️")

	secret, err := vc.client.Logical().Read(fmt.Sprintf("%s/data/%s", path, datasourceName))
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

// DeleteCredentials receives the name of the DataSource and attempts to delete the existing credentials on Vault.
// Is performs a soft delete, per docs > https://www.vaultproject.io/docs/commands/kv/delete
func (vc *VaultClient) DeleteCredentials(datasourceName string) error {
	vc.logger.Info("Deleting credentials from Vault ☄️")

	_, err := vc.client.Logical().Delete(fmt.Sprintf("%s/data/%s", path, datasourceName))
	if err != nil {
		vc.logger.Errorf("Unable to delete credentials from Vault 😱. Err: %v", err)
		return err
	}

	vc.logger.Info("Credentials deleted from Vault ☠️")
	return nil
}

// CheckIfEngineExists attempts to call the /tune API endpoint on the Secrets Engine. Should it fail, it might be an
// indication that the Secrets Engine is not created, which it's useful to know whether or not to call CreateEngine
func (vc *VaultClient) CheckIfEngineExists() bool {
	vc.logger.Info("Checking if datasource vault engine exists☄️")

	epath := "sys/mounts/" + path + "/tune"

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
