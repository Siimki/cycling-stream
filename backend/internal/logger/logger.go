package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init(env string) {
	Log = logrus.New()

	// Set log level based on environment
	switch env {
	case "production":
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "development", "dev":
		Log.SetLevel(logrus.DebugLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Set output to stdout
	Log.SetOutput(os.Stdout)
}

// Helper functions for common log operations
func Info(args ...interface{}) {
	Log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// WithField creates a new entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields creates a new entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithError creates a new entry with an error field
func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}

