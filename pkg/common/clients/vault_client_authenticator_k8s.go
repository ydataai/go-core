package clients

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// K8sAuthenticator defines a struct for authenticating with Kubernetes.
type K8sAuthenticator struct{}

// NewK8sAuthenticator defines a new K8sAuthenticator struct.
func NewK8sAuthenticator() Authenticator {
	return &K8sAuthenticator{}
}

// Authenticate is used to authenticate using Kubernetes.
func (a *K8sAuthenticator) Authenticate(vc *VaultClient) error {
	if err := a.login(vc); err != nil {
		return err
	}

	// do the token renewal cycle
	go a.renew(vc)

	return nil
}

func (a *K8sAuthenticator) login(vc *VaultClient) error {
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

	return nil
}

// renew the token according to secret.Auth.LeaseDuration automatically
func (a *K8sAuthenticator) renew(vc *VaultClient) {
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
			a.login(vc)
		}
	}
}
