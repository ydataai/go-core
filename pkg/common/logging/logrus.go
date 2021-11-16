package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type logrusProxy struct {
	logger *logrus.Logger
}

// newLogrusLogger construct and initializes the logrus logger.
func newLogrusLogger(config LoggerConfiguration) logrusProxy {
	log := logrus.StandardLogger()
	// formatter
	if strings.ToUpper(config.Output) == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetReportCaller(true)
		log.SetFormatter(&Formatter{
			TimestampFormat: "[2006-01-02 15:04:05]",
			CallerFirst:     config.CallerFirst,
			TrimMessages:    config.TrimMessages,
			HideKeys:        config.HideKeys,
		})
	}
	// level
	loglevel, err := logrus.ParseLevel(config.Level)
	if err != nil {
		log.Errorf("An error occurred while parsing LogLevel: %v", err)
		log.Infof("The debug level will be set.")
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(loglevel)
	}
	return logrusProxy{logger: log}
}

func (l logrusProxy) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l logrusProxy) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l logrusProxy) Print(args ...interface{}) {
	l.logger.Print(args...)
}

func (l logrusProxy) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l logrusProxy) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l logrusProxy) Warning(args ...interface{}) {
	l.logger.Warning(args...)
}

func (l logrusProxy) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l logrusProxy) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l logrusProxy) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l logrusProxy) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l logrusProxy) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l logrusProxy) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l logrusProxy) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l logrusProxy) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l logrusProxy) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}

func (l logrusProxy) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l logrusProxy) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l logrusProxy) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l logrusProxy) Traceln(args ...interface{}) {
	l.logger.Traceln(args...)
}

func (l logrusProxy) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l logrusProxy) Println(args ...interface{}) {
	l.logger.Println(args...)
}

func (l logrusProxy) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l logrusProxy) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

func (l logrusProxy) Warningln(args ...interface{}) {
	l.logger.Warningln(args...)
}

func (l logrusProxy) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

func (l logrusProxy) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}

func (l logrusProxy) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}
