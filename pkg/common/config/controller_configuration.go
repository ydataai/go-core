package config

import "github.com/kelseyhightower/envconfig"

// ControllerConfiguration defines required and common values for the controllers
type ControllerConfiguration struct {
	Environment string `envconfig:"ENVIRONMENT" required:"true"`
}

// LoadFromEnvVars for the controller configuration
func (c *ControllerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
