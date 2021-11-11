package logging

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// LoggerConfiguration defines a struct with required environment variables for logger
type LoggerConfiguration struct {
	Level        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
	Output       string `envconfig:"LOG_OUTPUT" default:""`
	CallerFirst  bool   `envconfig:"LOG_CALLER_FIRST" default:"false"`
	TrimMessages bool   `envconfig:"LOG_TRIM_MESSAGES" default:"true"`
	HideKeys     bool   `envconfig:"LOG_HIDE_KEYS" default:"false"`
}

// LoadFromEnvVars from the Logger
func (lc *LoggerConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", lc)
}

// InitLogger initializes the logger
func (lc *LoggerConfiguration) InitLogger() *logrus.Logger {
	log := logrus.StandardLogger()

	if strings.ToUpper(lc.Output) == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetReportCaller(true)
		log.SetFormatter(&Formatter{
			TimestampFormat: "[2006-01-02 15:04:05]",
			CallerFirst:     lc.CallerFirst,
			TrimMessages:    lc.TrimMessages,
			HideKeys:        lc.HideKeys,
		})
	}

	loglevel, err := logrus.ParseLevel(lc.Level)
	if err != nil {
		log.Errorf("An error occurred while parsing LogLevel: %v", err)
		log.Infof("The debug level will be set.")
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(loglevel)
	}

	return log
}
