package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Logger
var Logger = logrus.New()

// init initializes the logger
func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}

// SetLevel sets the logging level
func SetLevel(level string) error {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	Logger.SetLevel(l)
	return nil
}

// Error logs an error message
func Error(msg string, err error) {
	entry := Logger.WithFields(logrus.Fields{
		"error": err.Error(),
	})
	entry.Error(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn(msg)
}

// Info logs an info message
func Info(msg string) {
	Logger.Info(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug(msg)
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	return Logger
}
