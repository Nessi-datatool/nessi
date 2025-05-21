# Implementing Error Handling in Nessi

This guide provides practical examples and best practices for implementing error handling in the Nessi codebase. It is intended for developers who are contributing to Nessi and want to ensure their code follows the project's error handling standards.

## Table of Contents

- [Basic Error Handling](#basic-error-handling)
- [Using Error Codes](#using-error-codes)
- [Creating Structured Errors](#creating-structured-errors)
- [Adding Error Suggestions](#adding-error-suggestions)
- [Implementing Resolvable Errors](#implementing-resolvable-errors)
- [Integrating with Error Telemetry](#integrating-with-error-telemetry)
- [Using the Retry Mechanism](#using-the-retry-mechanism)
- [Testing Error Handling](#testing-error-handling)

## Basic Error Handling

All errors in Nessi should use the structured `NessiError` type from the `common` package. Here's a basic example:

```go
package mypackage

import (
	"fmt"
	"github.com/nessi-dev/nessi/pkg/common"
)

func ReadFile(path string) ([]byte, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create a structured error with error code
		return nil, common.NewError(common.ErrInvalidPath, fmt.Sprintf("Path %s does not exist", path))
	}
	
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		// Wrap the original error with context
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return data, nil
}
```

## Using Error Codes

Nessi uses standardized error codes in the format `NXXX` where the first digit indicates the error category. Always use the predefined error codes from the `common` package:

```go
// Path error (N1XX)
err := common.NewError(common.ErrInvalidPath, "Invalid path")

// Configuration error (N2XX)
err := common.NewError(common.ErrInvalidConfig, "Invalid configuration")

// Authentication error (N3XX)
err := common.NewError(common.ErrAuthenticationFailed, "Authentication failed")

// Connection error (N4XX)
err := common.NewError(common.ErrConnectionFailed, "Connection failed")

// Delta Lake error (N5XX)
err := common.NewError(common.ErrInvalidDeltaTable, "Not a valid Delta table")
```

## Creating Structured Errors

The `NessiError` type includes fields for error code, message, details, and suggestions. Here's how to create a fully structured error:

```go
func ValidateConfig(config *Config) error {
	if config.Host == "" {
		// Create a new error
		err := common.NewError(common.ErrInvalidConfig, "Invalid configuration: host is required")
		
		// Add details
		err.Details = "The host field in the configuration is empty"
		
		// Add suggestions
		err.Suggestions = []string{
			"Set the host field in your configuration file",
			"Use the --host flag to specify the host",
		}
		
		return err
	}
	
	return nil
}
```

## Adding Error Suggestions

In addition to including suggestions in the error itself, you can register global suggestions for error codes. This is useful for providing consistent suggestions for common errors:

```go
func init() {
	// Register suggestions for an error code
	common.RegisterErrorSuggestion(common.ErrorSuggestion{
		ErrorCode:   common.ErrInvalidConfig,
		Description: "The configuration is invalid or missing required fields",
		Solution:    "Check your configuration file for syntax errors or missing required fields",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#configuration-errors",
	})
	
	// You can register multiple suggestions for the same error code
	common.RegisterErrorSuggestion(common.ErrorSuggestion{
		ErrorCode:   common.ErrInvalidConfig,
		Description: "Environment variables might override configuration",
		Solution:    "Check if any environment variables are overriding your configuration",
		DocumentURL: "",
	})
}
```

## Implementing Resolvable Errors

For errors that can be resolved interactively, use the resolvable error types:

```go
func ValidatePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create a resolvable path error
		return common.NewResolvablePathError(path)
	}
	return nil
}

func ValidateConfig(configKey, configValue string) error {
	if !isValidConfigValue(configKey, configValue) {
		// Create a resolvable configuration error
		return common.NewResolvableConfigError(configKey, configValue)
	}
	return nil
}
```

Resolvable errors will be automatically handled by the CLI if interactive mode is enabled.

## Integrating with Error Telemetry

The error telemetry system automatically records errors when they're handled by the CLI. However, you can also record errors manually:

```go
func ProcessData(data []byte) error {
	// Get error telemetry
	telemetryConfig := common.DefaultErrorTelemetryConfig()
	telemetry := common.NewErrorTelemetry(telemetryConfig)
	
	// Process data
	err := processDataInternal(data)
	if err != nil {
		// Record the error
		telemetry.RecordError(err)
		return err
	}
	
	return nil
}
```

## Using the Retry Mechanism

For operations that might fail transiently (like network requests), use the retry mechanism:

```go
func FetchData(url string) ([]byte, error) {
	// Create retry policy
	policy := common.RetryPolicy{
		MaxAttempts:     3,
		InitialBackoff:  time.Second,
		MaxBackoff:      30 * time.Second,
		BackoffFactor:   2.0,
		Jitter:          true,
	}
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	
	// Use retry with exponential backoff
	var data []byte
	err := common.RetryWithBackoff(ctx, policy, func() error {
		var fetchErr error
		data, fetchErr = httpGet(url)
		if fetchErr != nil {
			// Check if the error is retryable
			if isRetryableError(fetchErr) {
				return fetchErr // Will be retried
			}
			return common.NonRetryableError{Err: fetchErr} // Won't be retried
		}
		return nil
	})
	
	return data, err
}
```

## Testing Error Handling

Always write tests for error scenarios. Here are some examples:

```go
func TestReadFileNonExistentPath(t *testing.T) {
	// Test with non-existent path
	_, err := ReadFile("/nonexistent/path")
	
	// Check if it's a NessiError
	var nessiErr *common.NessiError
	assert.True(t, errors.As(err, &nessiErr))
	
	// Check error code
	assert.Equal(t, common.ErrInvalidPath, nessiErr.Code)
	
	// Check message content
	assert.Contains(t, nessiErr.Message, "/nonexistent/path")
}

func TestErrorSuggestions(t *testing.T) {
	// Create an error
	err := common.NewError(common.ErrInvalidConfig, "test error")
	
	// Get suggestions
	suggestions := common.GetSuggestionsForError(err)
	
	// Check that suggestions are provided
	assert.NotEmpty(t, suggestions)
}

func TestResolvableError(t *testing.T) {
	// Create a resolvable error
	err := common.NewResolvablePathError("/test/path")
	
	// Check that it implements the ResolvableError interface
	_, ok := err.(common.ResolvableError)
	assert.True(t, ok)
}
```

## Conclusion

By following these guidelines, you'll ensure that your code provides a consistent and user-friendly error handling experience. Remember to:

1. Always use structured errors with appropriate error codes
2. Provide clear, actionable error messages
3. Add detailed suggestions for fixing errors
4. Implement interactive resolution for common errors
5. Write tests for error scenarios

For more information, see the [Error Handling Documentation](ERROR_HANDLING.md) and [Contributing Guidelines](CONTRIBUTING_ERROR_HANDLING.md).
