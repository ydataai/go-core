package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// RESTControllerConfiguration defines a struct with required environment variables for rest controller
type RESTControllerConfiguration struct {
	UserID             string        `envconfig:"USER_ID" default:""`
	UserIDPrefix       string        `envconfig:"USER_ID_PREFIX" default:""`
	HTTPRequestTimeout time.Duration `envconfig:"HTTP_REQUEST_TIMEOUT" default:"30s"`
}

// LoadFromEnvVars reads all env vars required
func (c *RESTControllerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
