package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolvableError represents an error that can be resolved interactively
type ResolvableError struct {
	*NessiError
	ResolutionSteps []string
}

// ResolveError implements the error resolution interface
func (e *ResolvableError) ResolveError() error {
	// This is a base implementation that doesn't do anything
	// Specific error types will override this method
	return nil
}

// NewResolvableError creates a new resolvable error
func NewResolvableError(code ErrorCode, message string, details string, resolutionSteps []string) *ResolvableError {
	nessiErr := NewError(code, message)
	if details != "" {
		nessiErr.Details = details
	}

	return &ResolvableError{
		NessiError:      nessiErr,
		ResolutionSteps: resolutionSteps,
	}
}

// ResolvablePathError represents a path error that can be resolved
type ResolvablePathError struct {
	*ResolvableError
	Path string
}

// NewResolvablePathError creates a new resolvable path error
func NewResolvablePathError(path string) *ResolvablePathError {
	message := fmt.Sprintf("Path '%s' does not exist or is not accessible", path)
	resolvableErr := NewResolvableError(ErrInvalidPath, message, "", []string{
		fmt.Sprintf("Create the directory: mkdir -p %s", path),
		"Check permissions and try again",
	})

	return &ResolvablePathError{
		ResolvableError: resolvableErr,
		Path:            path,
	}
}

// ResolveError implements the error resolution interface for path errors
func (e *ResolvablePathError) ResolveError() error {
	// Create the directory if it doesn't exist
	fmt.Printf("Creating directory: %s\n", e.Path)
	err := os.MkdirAll(e.Path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ResolvableDeltaTableError represents a Delta table error that can be resolved
type ResolvableDeltaTableError struct {
	*ResolvableError
	Path string
}

// NewResolvableDeltaTableError creates a new resolvable Delta table error
func NewResolvableDeltaTableError(path string) *ResolvableDeltaTableError {
	message := fmt.Sprintf("'%s' is not a Delta Lake table", path)
	resolvableErr := NewResolvableError(ErrInvalidDeltaTable, message, "", []string{
		fmt.Sprintf("Create a Delta table at: %s", path),
		"Specify a different path that contains a valid Delta table",
	})

	return &ResolvableDeltaTableError{
		ResolvableError: resolvableErr,
		Path:            path,
	}
}

// ResolveError implements the error resolution interface for Delta table errors
func (e *ResolvableDeltaTableError) ResolveError() error {
	// Create the directory if it doesn't exist
	fmt.Printf("Creating directory structure for Delta table: %s\n", e.Path)
	err := os.MkdirAll(e.Path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create _delta_log directory
	deltaLogPath := filepath.Join(e.Path, "_delta_log")
	err = os.MkdirAll(deltaLogPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Create an empty transaction log file
	txnLogPath := filepath.Join(deltaLogPath, "00000000000000000000.json")
	err = os.WriteFile(txnLogPath, []byte("{}"), 0644)
	if err != nil {
		return fmt.Errorf("failed to create transaction log: %w", err)
	}

	return nil
}

// ResolvableConfigError represents a configuration error that can be resolved
type ResolvableConfigError struct {
	*ResolvableError
	Key   string
	Value string
}

// NewResolvableConfigError creates a new resolvable configuration error
func NewResolvableConfigError(key string, value string) *ResolvableConfigError {
	message := fmt.Sprintf("Invalid configuration value for '%s': '%s'", key, value)
	resolvableErr := NewResolvableError(ErrInvalidConfig, message, "", []string{
		fmt.Sprintf("Set a valid value for '%s'", key),
		"Check the configuration documentation",
	})

	return &ResolvableConfigError{
		ResolvableError: resolvableErr,
		Key:             key,
		Value:           value,
	}
}

// ResolveError implements the error resolution interface for configuration errors
func (e *ResolvableConfigError) ResolveError() error {
	// This is a placeholder implementation
	// In a real implementation, this would prompt the user for a new value
	// and update the configuration
	fmt.Printf("Would update configuration key '%s' with a new value\n", e.Key)

	return nil
}

// GetAllErrorCodes returns all available error codes
func GetAllErrorCodes() []ErrorCode {
	// This is a simplified implementation that returns a subset of error codes
	return []ErrorCode{
		ErrUnknown,
		ErrInvalidPath,
		ErrInternalError,
		ErrInvalidDeltaTable,
		ErrInvalidConfig,
		ErrNotDeltaTable,
		ErrAuthFailed,
		ErrConnectionFailed,
		ErrResourceNotFound,
		ErrSchemaValidation,
	}
}
