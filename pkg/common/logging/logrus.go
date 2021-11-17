package logging

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	logCalldepth = 2 // specific to this wrapper implementation
	callerInfo   = "callerInfo"
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
			TimestampFormat:       "[2006-01-02 15:04:05]",
			CallerFirst:           config.CallerFirst,
			TrimMessages:          config.TrimMessages,
			HideKeys:              config.HideKeys,
			CustomCallerFormatter: formatter,
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

func fileInfo(callDepth int) string {
	// Inspect runtime call stack
	pc := make([]uintptr, callDepth)
	runtime.Callers(callDepth, pc)
	f := runtime.FuncForPC(pc[callDepth-1])
	file, line := f.FileLine(pc[callDepth-1])

	// Truncate abs file path
	if slash := strings.LastIndex(file, "/"); slash >= 0 {
		file = file[slash+1:]
	}

	return fmt.Sprintf("%s:%d", file, line-1)
}

func formatter(entry *logrus.Entry) string {
	if info, ok := entry.Data[callerInfo]; ok {
		return fmt.Sprintf(" [%v]", info)
	}
	return ""
}

func (l logrusProxy) Trace(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Trace(args...)
}

func (l logrusProxy) Debug(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Debug(args...)
}

func (l logrusProxy) Print(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Print(args...)
}

func (l logrusProxy) Info(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Info(args...)
}

func (l logrusProxy) Warn(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Warn(args...)
}

func (l logrusProxy) Warning(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Warning(args...)
}

func (l logrusProxy) Error(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Error(args...)
}

func (l logrusProxy) Panic(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Panic(args...)
}

func (l logrusProxy) Fatal(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Fatal(args...)
}

func (l logrusProxy) Tracef(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Tracef(format, args...)
}

func (l logrusProxy) Debugf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Debugf(format, args...)
}

func (l logrusProxy) Printf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Printf(format, args...)
}

func (l logrusProxy) Infof(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Infof(format, args...)
}

func (l logrusProxy) Warnf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Warnf(format, args...)
}

func (l logrusProxy) Warningf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Warningf(format, args...)
}

func (l logrusProxy) Errorf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Errorf(format, args...)
}

func (l logrusProxy) Panicf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Panicf(format, args...)
}

func (l logrusProxy) Fatalf(format string, args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Fatalf(format, args...)
}

func (l logrusProxy) Traceln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Traceln(args...)
}

func (l logrusProxy) Debugln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Debugln(args...)
}

func (l logrusProxy) Println(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Println(args...)
}

func (l logrusProxy) Infoln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Infoln(args...)
}

func (l logrusProxy) Warnln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Warnln(args...)
}

func (l logrusProxy) Warningln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Warningln(args...)
}

func (l logrusProxy) Errorln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Errorln(args...)
}

func (l logrusProxy) Panicln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Panicln(args...)
}

func (l logrusProxy) Fatalln(args ...interface{}) {
	f := fileInfo(logCalldepth)
	l.logger.WithField(callerInfo, f).Fatalln(args...)
}
