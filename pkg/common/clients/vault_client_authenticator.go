package clients

// Authenticator is an interface to identify which way to authenticate
type Authenticator interface {
	Authenticate(vc *VaultClient) error
}
