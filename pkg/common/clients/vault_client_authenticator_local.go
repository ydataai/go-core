package clients

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
