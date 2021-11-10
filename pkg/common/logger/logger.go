package logger

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	Level        string `envconfig:"LOG_LEVEL" default:"DEBUG"`
	Type         string `envconfig:"LOG_TYPE" default:""`
	CallerFirst  bool   `envconfig:"LOG_CALLER_FIRST" default:"false"`
	TrimMessages bool   `envconfig:"LOG_TRIM_MESSAGES" default:"true"`
	HideKeys     bool   `envconfig:"LOG_HIDE_KEYS" default:"false"`
}

func NewLogger() *logrus.Logger {
	config := LoggerConfig{}
	if err := envconfig.Process("", &config); err != nil {
		panic(err.Error())
	}

	loglevel, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Errorf("Error: %v", err)
		panic(err.Error())
	}

	log := logrus.StandardLogger()

	if strings.ToUpper(config.Type) == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetLevel(loglevel)
		log.SetReportCaller(true)
		log.SetFormatter(&Formatter{
			TimestampFormat: "[2006-01-02 15:04:05]",
			CallerFirst:     config.CallerFirst,
			TrimMessages:    config.TrimMessages,
			HideKeys:        config.HideKeys,
		})
	}

	return log
}
