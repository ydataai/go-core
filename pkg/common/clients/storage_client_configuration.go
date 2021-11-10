package clients

import (
	"github.com/kelseyhightower/envconfig"
)

// StorageClientConfiguration is a struct that holds all the environment variables required to the Storage Client
type StorageClientConfiguration struct {
	BasePath string `envconfig:"STORAGE_PATH" required:"true"`
}

// LoadEnvVars parses the required configuration variables. Throws an error if the validations aren't met
func (c *StorageClientConfiguration) LoadFromEnvVars() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}

	return nil
}
