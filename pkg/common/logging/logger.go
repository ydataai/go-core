package logging

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

// NewLogger creates a Logger proxy instance.
func NewLogger(config LoggerConfiguration) Logger {
	return newLogrusLogger(config)
}
