package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// InteractiveErrorResolver provides methods for interactively resolving errors
type InteractiveErrorResolver struct {
	Reader *bufio.Reader
	Writer *bufio.Writer
}

// NewInteractiveErrorResolver creates a new interactive error resolver
func NewInteractiveErrorResolver() *InteractiveErrorResolver {
	return &InteractiveErrorResolver{
		Reader: bufio.NewReader(os.Stdin),
		Writer: bufio.NewWriter(os.Stdout),
	}
}

// ResolvableError is an error that can be resolved interactively
type ResolvableError interface {
	// GetResolutionOptions returns a list of options for resolving the error
	GetResolutionOptions() []string
	// ResolveWithOption attempts to resolve the error with the selected option
	ResolveWithOption(option int) error
	// Error returns the error message
	Error() string
}

// ResolveError attempts to resolve an error interactively
func (r *InteractiveErrorResolver) ResolveError(err error) error {
	// Check if the error is resolvable
	resolvableErr, ok := err.(ResolvableError)
	if !ok {
		// Not a resolvable error
		return err
	}

	// Get resolution options
	options := resolvableErr.GetResolutionOptions()
	if len(options) == 0 {
		// No resolution options available
		return err
	}

	// Display error and resolution options
	fmt.Fprintf(r.Writer, "Error: %s\n\n", resolvableErr.Error())
	fmt.Fprintln(r.Writer, "Resolution options:")
	for i, option := range options {
		fmt.Fprintf(r.Writer, "%d. %s\n", i+1, option)
	}
	fmt.Fprintln(r.Writer, "0. Do nothing (continue with error)")
	r.Writer.Flush()

	// Get user selection
	fmt.Fprint(r.Writer, "\nSelect an option: ")
	r.Writer.Flush()

	selection, err := r.Reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read selection: %w", err)
	}

	// Parse selection
	selection = strings.TrimSpace(selection)
	var option int
	if _, err := fmt.Sscanf(selection, "%d", &option); err != nil {
		return fmt.Errorf("invalid selection: %s", selection)
	}

	// Handle selection
	if option == 0 {
		// User chose to do nothing
		return err
	}

	if option < 1 || option > len(options) {
		return fmt.Errorf("invalid selection: %d", option)
	}

	// Attempt to resolve the error
	return resolvableErr.ResolveWithOption(option - 1)
}

// ConfirmAction asks the user to confirm an action
func (r *InteractiveErrorResolver) ConfirmAction(prompt string) (bool, error) {
	fmt.Fprintf(r.Writer, "%s (y/n): ", prompt)
	r.Writer.Flush()

	response, err := r.Reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

// ResolvablePathError is a path error that can be resolved interactively
type ResolvablePathError struct {
	*NessiError
	Path string
}

// NewResolvablePathError creates a new resolvable path error
func NewResolvablePathError(path string) *ResolvablePathError {
	return &ResolvablePathError{
		NessiError: NewError(ErrInvalidPath, fmt.Sprintf("Path '%s' does not exist or is not accessible", path)).
			WithSuggestion("Check that the path exists and you have permission to access it"),
		Path: path,
	}
}

// GetResolutionOptions returns a list of options for resolving the path error
func (e *ResolvablePathError) GetResolutionOptions() []string {
	return []string{
		"Create the directory",
		"Specify a different path",
	}
}

// ResolveWithOption attempts to resolve the path error with the selected option
func (e *ResolvablePathError) ResolveWithOption(option int) error {
	switch option {
	case 0: // Create the directory
		return os.MkdirAll(e.Path, 0755)
	case 1: // Specify a different path
		resolver := NewInteractiveErrorResolver()
		fmt.Fprintf(resolver.Writer, "Enter a new path: ")
		resolver.Writer.Flush()

		newPath, err := resolver.Reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read new path: %w", err)
		}

		newPath = strings.TrimSpace(newPath)
		if newPath == "" {
			return fmt.Errorf("path cannot be empty")
		}

		// Check if the new path exists
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			// Ask if the user wants to create the new path
			confirm, err := resolver.ConfirmAction(fmt.Sprintf("Path '%s' does not exist. Create it?", newPath))
			if err != nil {
				return err
			}

			if confirm {
				if err := os.MkdirAll(newPath, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			} else {
				return fmt.Errorf("path does not exist: %s", newPath)
			}
		}

		// Update the path
		e.Path = newPath
		return nil
	default:
		return fmt.Errorf("invalid option: %d", option)
	}
}

// ResolvableDeltaTableError is a Delta table error that can be resolved interactively
type ResolvableDeltaTableError struct {
	*NessiError
	Path string
}

// NewResolvableDeltaTableError creates a new resolvable Delta table error
func NewResolvableDeltaTableError(path string) *ResolvableDeltaTableError {
	return &ResolvableDeltaTableError{
		NessiError: NewError(ErrNotDeltaTable, fmt.Sprintf("'%s' is not a Delta Lake table", path)).
			WithSuggestion("Ensure the path points to a valid Delta Lake table with a _delta_log directory"),
		Path: path,
	}
}

// GetResolutionOptions returns a list of options for resolving the Delta table error
func (e *ResolvableDeltaTableError) GetResolutionOptions() []string {
	return []string{
		"Initialize as a Delta table",
		"Specify a different path",
	}
}

// ResolveWithOption attempts to resolve the Delta table error with the selected option
func (e *ResolvableDeltaTableError) ResolveWithOption(option int) error {
	switch option {
	case 0: // Initialize as a Delta table
		// Create _delta_log directory
		deltaLogPath := fmt.Sprintf("%s/_delta_log", e.Path)
		if err := os.MkdirAll(deltaLogPath, 0755); err != nil {
			return fmt.Errorf("failed to create _delta_log directory: %w", err)
		}

		// Create an empty transaction log file
		transactionLogPath := fmt.Sprintf("%s/00000000000000000000.json", deltaLogPath)
		file, err := os.Create(transactionLogPath)
		if err != nil {
			return fmt.Errorf("failed to create transaction log file: %w", err)
		}
		defer file.Close()

		// Write an empty Delta table transaction log
		_, err = file.WriteString("{\"commitInfo\":{\"timestamp\":0,\"operation\":\"CREATE TABLE\",\"operationParameters\":{},\"isBlindAppend\":true}}\n")
		if err != nil {
			return fmt.Errorf("failed to write transaction log: %w", err)
		}

		return nil
	case 1: // Specify a different path
		resolver := NewInteractiveErrorResolver()
		fmt.Fprintf(resolver.Writer, "Enter a path to a Delta table: ")
		resolver.Writer.Flush()

		newPath, err := resolver.Reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read new path: %w", err)
		}

		newPath = strings.TrimSpace(newPath)
		if newPath == "" {
			return fmt.Errorf("path cannot be empty")
		}

		// Check if the new path is a Delta table
		deltaLogPath := fmt.Sprintf("%s/_delta_log", newPath)
		if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
			// Ask if the user wants to initialize the new path as a Delta table
			confirm, err := resolver.ConfirmAction(fmt.Sprintf("Path '%s' is not a Delta table. Initialize it as one?", newPath))
			if err != nil {
				return err
			}

			if confirm {
				// Create _delta_log directory
				if err := os.MkdirAll(deltaLogPath, 0755); err != nil {
					return fmt.Errorf("failed to create _delta_log directory: %w", err)
				}

				// Create an empty transaction log file
				transactionLogPath := fmt.Sprintf("%s/00000000000000000000.json", deltaLogPath)
				file, err := os.Create(transactionLogPath)
				if err != nil {
					return fmt.Errorf("failed to create transaction log file: %w", err)
				}
				defer file.Close()

				// Write an empty Delta table transaction log
				_, err = file.WriteString("{\"commitInfo\":{\"timestamp\":0,\"operation\":\"CREATE TABLE\",\"operationParameters\":{},\"isBlindAppend\":true}}\n")
				if err != nil {
					return fmt.Errorf("failed to write transaction log: %w", err)
				}
			} else {
				return fmt.Errorf("path is not a Delta table: %s", newPath)
			}
		}

		// Update the path
		e.Path = newPath
		return nil
	default:
		return fmt.Errorf("invalid option: %d", option)
	}
}

// ResolvableConfigError is a configuration error that can be resolved interactively
type ResolvableConfigError struct {
	*NessiError
	ConfigKey   string
	ConfigValue string
}

// NewResolvableConfigError creates a new resolvable configuration error
func NewResolvableConfigError(key, value string) *ResolvableConfigError {
	return &ResolvableConfigError{
		NessiError: NewError(ErrInvalidConfig, fmt.Sprintf("Invalid configuration value for '%s': '%s'", key, value)).
			WithSuggestion(fmt.Sprintf("Provide a valid value for '%s'", key)),
		ConfigKey:   key,
		ConfigValue: value,
	}
}

// GetResolutionOptions returns a list of options for resolving the configuration error
func (e *ResolvableConfigError) GetResolutionOptions() []string {
	return []string{
		"Provide a new value",
		"Use default value",
	}
}

// ResolveWithOption attempts to resolve the configuration error with the selected option
func (e *ResolvableConfigError) ResolveWithOption(option int) error {
	switch option {
	case 0: // Provide a new value
		resolver := NewInteractiveErrorResolver()
		fmt.Fprintf(resolver.Writer, "Enter a new value for '%s': ", e.ConfigKey)
		resolver.Writer.Flush()

		newValue, err := resolver.Reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read new value: %w", err)
		}

		newValue = strings.TrimSpace(newValue)
		if newValue == "" {
			return fmt.Errorf("value cannot be empty")
		}

		// Update the value
		e.ConfigValue = newValue
		return nil
	case 1: // Use default value
		// Get default value based on config key
		defaultValue := getDefaultConfigValue(e.ConfigKey)
		if defaultValue == "" {
			return fmt.Errorf("no default value available for '%s'", e.ConfigKey)
		}

		// Update the value
		e.ConfigValue = defaultValue
		return nil
	default:
		return fmt.Errorf("invalid option: %d", option)
	}
}

// getDefaultConfigValue returns the default value for a configuration key
func getDefaultConfigValue(key string) string {
	// Define default values for common configuration keys
	defaultValues := map[string]string{
		"DATABRICKS_HOST":            "https://dbc-12345-abcd.cloud.databricks.com",
		"DATABRICKS_TOKEN":           "",
		"DATABRICKS_WORKSPACE_ID":    "0",
		"DATABRICKS_DEFAULT_SCHEMA":  "default",
		"DATABRICKS_DEFAULT_CATALOG": "hive_metastore",
		"AWS_REGION":                 "us-west-2",
		"GCP_PROJECT_ID":             "",
		"AZURE_TENANT_ID":            "",
		"NESSI_LOG_LEVEL":            "info",
	}

	return defaultValues[key]
}
