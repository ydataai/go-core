package logger

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Configuration defines a struct with required environment variables for logger
type Configuration struct {
	Level        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
	Output       string `envconfig:"LOG_OUTPUT" default:""`
	CallerFirst  bool   `envconfig:"LOG_CALLER_FIRST" default:"false"`
	TrimMessages bool   `envconfig:"LOG_TRIM_MESSAGES" default:"true"`
	HideKeys     bool   `envconfig:"LOG_HIDE_KEYS" default:"false"`
}

// LoadFromEnvVars from the Logger
func (c *Configuration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}

// NewConfiguration for logger
func NewConfiguration(level, output string, cf, trim, hk bool) *Configuration {
	return &Configuration{
		Level:        level,
		Output:       output,
		CallerFirst:  cf,
		TrimMessages: trim,
		HideKeys:     hk,
	}
}

// InitLogger initializes the logger
func (c *Configuration) InitLogger() *logrus.Logger {
	log := logrus.StandardLogger()

	if strings.ToUpper(c.Output) == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetReportCaller(true)
		log.SetFormatter(&Formatter{
			TimestampFormat: "[2006-01-02 15:04:05]",
			CallerFirst:     c.CallerFirst,
			TrimMessages:    c.TrimMessages,
			HideKeys:        c.HideKeys,
		})
	}

	loglevel, err := logrus.ParseLevel(c.Level)
	if err != nil {
		log.Errorf("An error occurred while parsing LogLevel: %v", err)
		log.Infof("The debug level will be set.")
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(loglevel)
	}

	return log
}
