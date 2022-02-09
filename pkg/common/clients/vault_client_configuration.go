package clients

import "github.com/kelseyhightower/envconfig"

// VaultClientConfiguration is a struct that holds all the environment variables required to the Vault client
type VaultClientConfiguration struct {
	VaultURL  string `envconfig:"VAULT_SERVER_URL" required:"true"`
	VaultRole string `envconfig:"VAULT_ROLE" required:"true"`
}

// LoadFromEnvVars parses the required configuration variables. Throws an error if the validations aren't met
func (vlc *VaultClientConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", vlc)
}
