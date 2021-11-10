package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// RESTControllerConfiguration defines a struct with required environment variables for rest controller
type RESTControllerConfiguration struct {
	UserID             string        `envconfig:"USER_ID" required:"true"`
	UserIDPrefix       string        `envconfig:"USER_ID_PREFIX" default:""`
	HTTPRequestTimeout time.Duration `envconfig:"HTTP_REQUEST_TIMEOUT" default:"30s"`
}

// LoadEnvVars reads all env vars required for the server package
func (c *RESTControllerConfiguration) LoadFromEnvVars() error {
	if err := envconfig.Process("", c); err != nil {
		return err
	}

	return nil
}
