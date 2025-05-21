package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/rs/zerolog"
)

// Log levels
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Logger is a wrapper around zerolog.Logger with additional methods
type Logger struct {
	zerolog.Logger
	ctx context.Context
}

// contextKey is the type used for context values
type contextKey string

// loggerKey is the key used to store the logger in the context
const loggerKey = contextKey("logger")

// DefaultLogger is the default logger instance
var DefaultLogger = NewLogger(InfoLevel)

// init initializes the default logger
func init() {
	// Set default time format to ISO8601
	zerolog.TimeFieldFormat = time.RFC3339

	// Get log level from environment variable
	logLevel := os.Getenv("NESSI_LOG_LEVEL")
	if logLevel == "" {
		logLevel = InfoLevel
	}

	// Initialize default logger
	DefaultLogger = NewLogger(logLevel)
}

// NewLogger creates a new logger with the specified log level
func NewLogger(level string) Logger {
	// Set log level
	var logLevel zerolog.Level
	switch strings.ToLower(level) {
	case DebugLevel:
		logLevel = zerolog.DebugLevel
	case InfoLevel:
		logLevel = zerolog.InfoLevel
	case WarnLevel:
		logLevel = zerolog.WarnLevel
	case ErrorLevel:
		logLevel = zerolog.ErrorLevel
	case FatalLevel:
		logLevel = zerolog.FatalLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Create writer
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		},
	}

	// Create logger
	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Logger()

	return Logger{
		Logger: logger,
		ctx:    context.Background(),
	}
}

// WithOutput returns a new logger with the specified output writer
func (l Logger) WithOutput(w io.Writer) Logger {
	return Logger{
		Logger: l.Output(w),
		ctx:    l.ctx,
	}
}

// WithContext returns a new logger with the specified context
func (l Logger) WithContext(ctx context.Context) Logger {
	return Logger{
		Logger: l.Logger,
		ctx:    ctx,
	}
}

// FromContext returns a logger from the context or the default logger if none is found
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return DefaultLogger
	}

	logger, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return DefaultLogger
	}

	return logger
}

// WithLogger returns a new context with the logger attached
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Error logs an error with additional context
func (l Logger) Error(err error, msg string) {
	// Get caller information
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Extract filename from path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	// Check if it's a NessiError
	var nessiErr *common.NessiError
	if err != nil && errors.As(err, &nessiErr) {
		// Log with error code and details
		l.Logger.Error().
			Str("error_code", string(nessiErr.Code)).
			Str("error_type", common.GetErrorDescription(nessiErr.Code)).
			Str("error", nessiErr.Error()).
			Str("file", file).
			Int("line", line).
			Msg(msg)
	} else if err != nil {
		// Log regular error
		l.Logger.Error().
			Err(err).
			Str("file", file).
			Int("line", line).
			Msg(msg)
	} else {
		// Log error message without error
		l.Logger.Error().
			Str("file", file).
			Int("line", line).
			Msg(msg)
	}
}

// ErrorCode logs an error with a specific error code
func (l Logger) ErrorCode(code common.ErrorCode, msg string, details ...string) {
	// Get caller information
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Extract filename from path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	// Create error event
	event := l.Logger.Error().
		Str("error_code", string(code)).
		Str("error_type", common.GetErrorDescription(code)).
		Str("file", file).
		Int("line", line)

	// Add details if provided
	if len(details) > 0 {
		event = event.Str("details", strings.Join(details, "; "))
	}

	// Log message
	event.Msg(msg)
}

// Debug logs a debug message
func (l Logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Debug()

	// Add fields if provided
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = addField(event, k, v)
		}
	}

	event.Msg(msg)
}

// Info logs an info message
func (l Logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Info()

	// Add fields if provided
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = addField(event, k, v)
		}
	}

	event.Msg(msg)
}

// Warn logs a warning message
func (l Logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Warn()

	// Add fields if provided
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = addField(event, k, v)
		}
	}

	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l Logger) Fatal(msg string, fields ...map[string]interface{}) {
	event := l.Logger.Fatal()

	// Add fields if provided
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = addField(event, k, v)
		}
	}

	event.Msg(msg)
}

// addField adds a field to the event based on its type
func addField(event *zerolog.Event, key string, value interface{}) *zerolog.Event {
	switch v := value.(type) {
	case string:
		return event.Str(key, v)
	case int:
		return event.Int(key, v)
	case int64:
		return event.Int64(key, v)
	case float64:
		return event.Float64(key, v)
	case bool:
		return event.Bool(key, v)
	case time.Time:
		return event.Time(key, v)
	case []string:
		return event.Strs(key, v)
	case []int:
		return event.Ints(key, v)
	case []int64:
		return event.Ints64(key, v)
	case error:
		return event.AnErr(key, v)
	default:
		return event.Interface(key, v)
	}
}
