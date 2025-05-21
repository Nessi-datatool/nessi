package cli

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/nessi-dev/nessi/pkg/logger"
	"github.com/spf13/cobra"
)

// ErrorHandler handles errors in the CLI
type ErrorHandler struct {
	Interactive bool
	Logger      logger.Logger
	Resolver    *common.InteractiveErrorResolver
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(interactive bool) *ErrorHandler {
	return &ErrorHandler{
		Interactive: interactive,
		Logger:      logger.DefaultLogger,
		Resolver:    common.NewInteractiveErrorResolver(),
	}
}

// HandleError handles an error in the CLI
func (h *ErrorHandler) HandleError(err error) error {
	if err == nil {
		return nil
	}

	// Log the error
	h.Logger.Error(err, "Command failed")

	// Check if it's a NessiError
	var nessiErr *common.NessiError
	if errors.As(err, &nessiErr) {
		// Print the error with color
		printColoredError(nessiErr)

		// Return the error without interactive resolution if not enabled
		if !h.Interactive {
			return fmt.Errorf("error: %w", err)
		}

		// Check if it's a resolvable error
		resolvableErr, ok := err.(common.ResolvableError)
		if !ok {
			return fmt.Errorf("error: %w", err)
		}

		// Try to resolve the error interactively
		return h.Resolver.ResolveError(resolvableErr)
	}

	// For non-NessiError, just print the error
	color.Red("Error: %s", err.Error())
	return fmt.Errorf("error: %w", err)
}

// printColoredError prints a NessiError with color
func printColoredError(err *common.NessiError) {
	// Print error code and type
	color.New(color.FgRed, color.Bold).Printf("[%s] %s: ", err.Code, common.GetErrorDescription(err.Code))

	// Print error message
	color.New(color.FgRed).Println(err.Message)

	// Print details if available
	if err.Details != "" {
		color.New(color.FgYellow).Printf("Details: %s\n", err.Details)
	}

	// Print suggestion if available
	if err.Suggestion != "" {
		color.New(color.FgGreen).Printf("Suggestion: %s\n", err.Suggestion)
	}
}

// AddErrorHandlingFlags adds error handling flags to a command
func AddErrorHandlingFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("interactive", true, "Enable interactive error resolution")
}

// GetErrorHandler creates an error handler from command flags
func GetErrorHandler(cmd *cobra.Command) *ErrorHandler {
	interactive, _ := cmd.Flags().GetBool("interactive")
	return NewErrorHandler(interactive)
}

// WrapCobraCommand wraps a cobra command with error handling
func WrapCobraCommand(runE func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Get error handler
		handler := GetErrorHandler(cmd)

		// Run the command
		err := runE(cmd, args)

		// Handle the error
		return handler.HandleError(err)
	}
}

// ConvertToResolvableError converts a standard error to a resolvable error if possible
func ConvertToResolvableError(err error) error {
	if err == nil {
		return nil
	}

	// Check for path not found errors
	if os.IsNotExist(err) {
		// Extract path from error message
		path := extractPathFromError(err.Error())
		if path != "" {
			return common.NewResolvablePathError(path)
		}
	}

	// Check for Delta table errors
	if strings.Contains(err.Error(), "not a Delta table") || strings.Contains(err.Error(), "_delta_log") {
		// Extract path from error message
		path := extractPathFromError(err.Error())
		if path != "" {
			return common.NewResolvableDeltaTableError(path)
		}
	}

	// Check for configuration errors
	if strings.Contains(err.Error(), "configuration") || strings.Contains(err.Error(), "config") {
		// Try to extract key and value from error message
		key, value := extractConfigFromError(err.Error())
		if key != "" {
			return common.NewResolvableConfigError(key, value)
		}
	}

	return err
}

// extractPathFromError attempts to extract a path from an error message
func extractPathFromError(errMsg string) string {
	// Common patterns for path errors
	patterns := []string{
		"path '(.+)' does not exist",
		"'(.+)' is not a Delta Lake table",
		"no such file or directory: (.+)",
		"cannot access '(.+)'",
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(errMsg)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// extractConfigFromError attempts to extract a configuration key and value from an error message
func extractConfigFromError(errMsg string) (string, string) {
	// Common patterns for configuration errors
	patterns := []string{
		"invalid configuration value for '(.+)': '(.*)'",
		"missing required configuration: (.+)",
		"(.+) is required",
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(errMsg)
		if len(matches) > 1 {
			key := matches[1]
			value := ""
			if len(matches) > 2 {
				value = matches[2]
			}
			return key, value
		}
	}

	return "", ""
}
