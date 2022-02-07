package clients

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Authenticator is an interface to identify which way to authenticate
type Authenticator interface {
	Authenticate(vc *VaultClient) error
}

// K8sAuthenticator defines a struct for authenticating with Kubernetes.
type K8sAuthenticator struct{}

// NewK8sAuthenticator defines a new K8sAuthenticator struct.
func NewK8sAuthenticator() Authenticator {
	return &K8sAuthenticator{}
}

// Authenticate is used to authenticate using Kubernetes.
func (a *K8sAuthenticator) Authenticate(vc *VaultClient) error {
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
	go vc.renew(a)

	return nil
}

// LocalAuthenticator is used to configure the development mode
type LocalAuthenticator struct {
	token string
}

// NewLocalAuthenticator defines a new LocalAuthenticator
func NewLocalAuthenticator(token string) Authenticator {
	return &LocalAuthenticator{token: token}
}

// Authenticate is used to authenticate local (development mode).
func (a *LocalAuthenticator) Authenticate(vc *VaultClient) error {
	vc.client.SetToken(a.token)
	return nil
}
