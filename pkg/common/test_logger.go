package common

import (
	"fmt"
	"os"
	"time"
)

// TestLogger is a simplified logger for testing purposes
type TestLogger struct {
	Level string
}

// NewTestLogger creates a new test logger
func NewTestLogger(level string) *TestLogger {
	return &TestLogger{
		Level: level,
	}
}

// Log logs a message with the specified level
func (l *TestLogger) Log(level, msg string) {
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, msg)
}

// Debug logs a debug message
func (l *TestLogger) Debug(msg string) {
	l.Log("DEBUG", msg)
}

// Info logs an info message
func (l *TestLogger) Info(msg string) {
	l.Log("INFO", msg)
}

// Warn logs a warning message
func (l *TestLogger) Warn(msg string) {
	l.Log("WARN", msg)
}

// Error logs an error message
func (l *TestLogger) Error(err error, msg string) {
	if err != nil {
		l.Log("ERROR", fmt.Sprintf("%s: %s", msg, err))
	} else {
		l.Log("ERROR", msg)
	}
}

// Fatal logs a fatal message
func (l *TestLogger) Fatal(msg string) {
	l.Log("FATAL", msg)
	os.Exit(1)
}

// ErrorCode logs an error with a specific error code
func (l *TestLogger) ErrorCode(code ErrorCode, msg string, details ...string) {
	l.Log("ERROR", fmt.Sprintf("[%s] %s", code, msg))
}
