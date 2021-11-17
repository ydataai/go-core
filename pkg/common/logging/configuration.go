package logging

import (
	"github.com/kelseyhightower/envconfig"
)

// LoggerConfiguration defines a struct with required environment variables for logger
type LoggerConfiguration struct {
	Level        string `envconfig:"LOG_LEVEL" default:"debug"`
	Output       string `envconfig:"LOG_OUTPUT" default:""`
	CallerFirst  bool   `envconfig:"LOG_CALLER_FIRST" default:"false"`
	TrimMessages bool   `envconfig:"LOG_TRIM_MESSAGES" default:"true"`
	HideKeys     bool   `envconfig:"LOG_HIDE_KEYS" default:"false"`
}

// LoadFromEnvVars from the Logger
func (lc *LoggerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", lc)
}
