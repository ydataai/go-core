package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/ydataai/go-core/pkg/common"
)

// ControllerConfiguration defines required and common values for the controllers
type ControllerConfiguration struct {
	Environment common.Environment `envconfig:"ENVIRONMENT" required:"true"`
}

// LoadFromEnvVars for the controller configuration
func (c *ControllerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
