package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nessi-dev/nessi/pkg/errorcode"
)

// PathValidator validates and normalizes file paths
type PathValidator struct {
	AllowNonExistent bool // Whether to allow paths that don't exist
	RequireDirectory bool // Whether to require the path to be a directory
	RequireFile      bool // Whether to require the path to be a file
}

// NewPathValidator creates a new path validator with default settings
func NewPathValidator() *PathValidator {
	return &PathValidator{
		AllowNonExistent: false,
		RequireDirectory: false,
		RequireFile:      false,
	}
}

// ValidatePath validates and normalizes a file path
// Returns the normalized path if valid, or an error if invalid
func (v *PathValidator) ValidatePath(path string) (string, error) {
	// Check if path is empty
	if path == "" {
		return "", NewError(errorcode.ErrInvalidPath, "Path cannot be empty")
	}

	// Expand home directory if path starts with ~
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		path = filepath.Join(cwd, path)
	}

	// Clean the path
	path = filepath.Clean(path)

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) && !v.AllowNonExistent {
			return "", NewError(errorcode.ErrInvalidPath, fmt.Sprintf("Path %s does not exist", path))
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to access path: %w", err)
		}
	} else {
		// Path exists, check if it's a directory or file as required
		if v.RequireDirectory && !info.IsDir() {
			return "", NewError(errorcode.ErrInvalidArgument, fmt.Sprintf("%s is not a directory", path))
		}
		if v.RequireFile && info.IsDir() {
			return "", NewError(errorcode.ErrInvalidArgument, fmt.Sprintf("%s is not a file", path))
		}
	}

	return path, nil
}

// ValidateTablePath validates and normalizes a Delta Lake table path
func ValidateTablePath(path string) (string, error) {
	// Create validator for Delta tables
	validator := NewPathValidator()
	validator.RequireDirectory = true

	// Validate the base path
	normalizedPath, err := validator.ValidatePath(path)
	if err != nil {
		return "", err
	}

	// Check if it's a Delta table (has _delta_log directory)
	deltaLogPath := filepath.Join(normalizedPath, "_delta_log")
	if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
		return "", NewError(errorcode.ErrNotDeltaTable, fmt.Sprintf("%s is not a Delta Lake table (missing _delta_log directory)", path))
	} else if err != nil {
		return "", fmt.Errorf("failed to access _delta_log directory: %w", err)
	}

	return normalizedPath, nil
}

// ConfigValidator validates configuration values
type ConfigValidator struct {
	Required bool // Whether the configuration is required
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		Required: false,
	}
}

// ValidateConfig validates a configuration value
func (v *ConfigValidator) ValidateConfig(key string, value string) error {
	// Check if value is empty
	if value == "" && v.Required {
		return NewError(errorcode.ErrConfigNotFound, fmt.Sprintf("Configuration %s is required", key))
	}

	return nil
}

// ValidateOutputFormat validates and normalizes an output format
func ValidateOutputFormat(format string) (string, error) {
	// Normalize to lowercase
	format = strings.ToLower(format)

	// Check if format is valid
	validFormats := map[string]bool{
		"html": true,
		"pdf":  true,
		"json": true,
		"csv":  true,
		"text": true,
	}

	if !validFormats[format] {
		return "", NewError(errorcode.ErrInvalidArgument, fmt.Sprintf("Invalid output format: %s. Valid formats are: html, pdf, json, csv, text", format))
	}

	return format, nil
}
