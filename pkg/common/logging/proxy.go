package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Logger defines the common logging operations.
type Logger interface {
	// Trace logs a message at level Trace on the standard logger.
	Trace(args ...interface{})
	// Debug logs a message at level Debug on the standard logger.
	Debug(args ...interface{})
	// Print logs a message at level Info on the standard logger.
	Print(args ...interface{})
	// Info logs a message at level Info on the standard logger.
	Info(args ...interface{})
	// Warn logs a message at level Warn on the standard logger.
	Warn(args ...interface{})
	// Warning logs a message at level Warn on the standard logger.
	Warning(args ...interface{})
	// Error logs a message at level Error on the standard logger.
	Error(args ...interface{})
	// Panic logs a message at level Panic on the standard logger.
	Panic(args ...interface{})
	// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
	Fatal(args ...interface{})
	// Tracef logs a message at level Trace on the standard logger.
	Tracef(format string, args ...interface{})
	// Debugf logs a message at level Debug on the standard logger.
	Debugf(format string, args ...interface{})
	// Printf logs a message at level Info on the standard logger.
	Printf(format string, args ...interface{})
	// Infof logs a message at level Info on the standard logger.
	Infof(format string, args ...interface{})
	// Warnf logs a message at level Warn on the standard logger.
	Warnf(format string, args ...interface{})
	// Warningf logs a message at level Warn on the standard logger.
	Warningf(format string, args ...interface{})
	// Errorf logs a message at level Error on the standard logger.
	Errorf(format string, args ...interface{})
	// Panicf logs a message at level Panic on the standard logger.
	Panicf(format string, args ...interface{})
	// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
	Fatalf(format string, args ...interface{})
	// Traceln logs a message at level Trace on the standard logger.
	Traceln(args ...interface{})
	// Debugln logs a message at level Debug on the standard logger.
	Debugln(args ...interface{})
	// Println logs a message at level Info on the standard logger.
	Println(args ...interface{})
	// Infoln logs a message at level Info on the standard logger.
	Infoln(args ...interface{})
	// Warnln logs a message at level Warn on the standard logger.
	Warnln(args ...interface{})
	// Warningln logs a message at level Warn on the standard logger.
	Warningln(args ...interface{})
	// Errorln logs a message at level Error on the standard logger.
	Errorln(args ...interface{})
	// Panicln logs a message at level Panic on the standard logger.
	Panicln(args ...interface{})
	// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
	Fatalln(args ...interface{})
}

type loggerProxy struct {
	logger *logrus.Logger
}

// NewLogger creates a Logger proxy instance.
func NewLogger() Logger {
	config := LoggerConfiguration{}
	if err := config.LoadFromEnvVars(); err != nil {
		fmt.Printf("An error occurred while logger initialization. Err: %v", err)
		fmt.Printf("Initializing standard logger.")
		return &loggerProxy{logger: logrus.StandardLogger()}
	}
	return &loggerProxy{logger: config.InitLogger()}
}

func (l loggerProxy) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l loggerProxy) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l loggerProxy) Print(args ...interface{}) {
	l.logger.Print(args...)
}

func (l loggerProxy) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l loggerProxy) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l loggerProxy) Warning(args ...interface{}) {
	l.logger.Warning(args...)
}

func (l loggerProxy) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l loggerProxy) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l loggerProxy) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l loggerProxy) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l loggerProxy) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l loggerProxy) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l loggerProxy) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l loggerProxy) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l loggerProxy) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
}

func (l loggerProxy) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l loggerProxy) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l loggerProxy) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l loggerProxy) Traceln(args ...interface{}) {
	l.logger.Traceln(args...)
}

func (l loggerProxy) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l loggerProxy) Println(args ...interface{}) {
	l.logger.Println(args...)
}

func (l loggerProxy) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l loggerProxy) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

func (l loggerProxy) Warningln(args ...interface{}) {
	l.logger.Warningln(args...)
}

func (l loggerProxy) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

func (l loggerProxy) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}

func (l loggerProxy) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}
