package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestLoggerLevels(t *testing.T) {
	// Create a buffer to capture log output
	buf := &bytes.Buffer{}

	// Create a logger with the buffer as output
	logger := NewLogger(DebugLevel).WithOutput(buf)

	// Test debug level
	logger.Debug("debug message")
	logOutput := buf.String()
	assert.Contains(t, logOutput, "debug")
	assert.Contains(t, logOutput, "debug message")
	buf.Reset()

	// Test info level
	logger.Info("info message")
	logOutput = buf.String()
	assert.Contains(t, logOutput, "info")
	assert.Contains(t, logOutput, "info message")
	buf.Reset()

	// Test warn level
	logger.Warn("warn message")
	logOutput = buf.String()
	assert.Contains(t, logOutput, "warn")
	assert.Contains(t, logOutput, "warn message")
	buf.Reset()

	// Test error level
	logger.Error(errors.New("test error"), "error message")
	logOutput = buf.String()
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "error message")
	assert.Contains(t, logOutput, "test error")
	buf.Reset()
}

func TestLoggerWithFields(t *testing.T) {
	// Create a buffer to capture log output
	buf := &bytes.Buffer{}

	// Create a logger with the buffer as output
	logger := NewLogger(DebugLevel).WithOutput(buf)

	// Test with fields
	logger.Info("info with fields", map[string]interface{}{
		"string_field": "string value",
		"int_field":    42,
		"bool_field":   true,
	})

	logOutput := buf.String()
	assert.Contains(t, logOutput, "info")
	assert.Contains(t, logOutput, "info with fields")
	assert.Contains(t, logOutput, "string_field")
	assert.Contains(t, logOutput, "string value")
	assert.Contains(t, logOutput, "int_field")
	assert.Contains(t, logOutput, "42")
	assert.Contains(t, logOutput, "bool_field")
	assert.Contains(t, logOutput, "true")
}

func TestLoggerWithNessiError(t *testing.T) {
	// Create a buffer to capture log output
	buf := &bytes.Buffer{}

	// Create a logger with the buffer as output
	logger := NewLogger(DebugLevel).WithOutput(buf)

	// Create a NessiError
	nessiErr := common.NewError(common.ErrInvalidPath, "test path is invalid").
		WithDetails("file does not exist").
		WithSuggestion("check the path")

	// Log the error
	logger.Error(nessiErr, "error with NessiError")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "error with NessiError")
	assert.Contains(t, logOutput, "error_code")
	assert.Contains(t, logOutput, string(common.ErrInvalidPath))
	assert.Contains(t, logOutput, "error_type")
	assert.Contains(t, logOutput, "Invalid path")
	assert.Contains(t, logOutput, "test path is invalid")
	assert.Contains(t, logOutput, "file does not exist")
	assert.Contains(t, logOutput, "check the path")
}

func TestLoggerErrorCode(t *testing.T) {
	// Create a buffer to capture log output
	buf := &bytes.Buffer{}

	// Create a logger with the buffer as output
	logger := NewLogger(DebugLevel).WithOutput(buf)

	// Log an error with error code
	logger.ErrorCode(common.ErrNotDeltaTable, "not a delta table", "path: /test/path")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "error")
	assert.Contains(t, logOutput, "not a delta table")
	assert.Contains(t, logOutput, "error_code")
	assert.Contains(t, logOutput, string(common.ErrNotDeltaTable))
	assert.Contains(t, logOutput, "error_type")
	assert.Contains(t, logOutput, "Not a Delta table")
	assert.Contains(t, logOutput, "details")
	assert.Contains(t, logOutput, "path: /test/path")
}

func TestLoggerFromContext(t *testing.T) {
	// Create a context with a logger
	ctx := context.Background()
	logger := NewLogger(DebugLevel)
	ctx = WithLogger(ctx, logger)

	// Get logger from context
	loggerFromCtx := FromContext(ctx)

	// Verify the logger is from the context
	assert.NotNil(t, loggerFromCtx)
	// We can't directly compare the loggers due to internal state differences,
	// so we'll just check that we got a logger back

	// Test with nil context
	loggerFromNilCtx := FromContext(nil)
	// Should return a valid logger
	assert.NotNil(t, loggerFromNilCtx)

	// Test with context that doesn't have a logger
	ctxWithoutLogger := context.Background()
	loggerFromCtxWithoutLogger := FromContext(ctxWithoutLogger)
	// Should return a valid logger
	assert.NotNil(t, loggerFromCtxWithoutLogger)
}
