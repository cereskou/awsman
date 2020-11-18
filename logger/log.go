package logger

import (
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log Level
const (
	_          = iota
	TRACE uint = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

//初期化
func init() {
	//Log format
	format := &prefixed.TextFormatter{
		DisableColors: true,
		// TimestampFormat: "2006-01-02 15:04:05.000",
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(format)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

}

//SetOutput -
func SetOutput(target string, logfile string) {
	var output io.Writer

	switch strings.ToLower(target) {
	case "screen":
		output = os.Stdout
		logrus.SetOutput(output)
	case "file":
		output = &lumberjack.Logger{
			Filename:  logfile,
			MaxSize:   150,
			MaxAge:    28,
			LocalTime: true,
			Compress:  true,
		}
		logrus.SetOutput(output)
	default:
		output = &lumberjack.Logger{
			Filename:  logfile,
			MaxSize:   150,
			MaxAge:    28,
			LocalTime: true,
			Compress:  true,
		}
		mw := io.MultiWriter(os.Stdout, output)
		logrus.SetOutput(mw)
	}
}

//SetLevel -
func SetLevel(level string) {
	lvl, _ := logrus.ParseLevel(level)
	logrus.SetLevel(lvl)
}

//Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Trace  logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	logrus.Trace(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	logrus.Tracef(format, args...)
}

// Warn  logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}
