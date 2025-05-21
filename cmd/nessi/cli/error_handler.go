package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/nessi-dev/nessi/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ErrorHandler handles errors in the CLI
type ErrorHandler struct {
	Interactive bool
	Logger      logger.Logger
	Resolver    *common.InteractiveErrorResolver
	Telemetry   *common.ErrorTelemetry
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(interactive bool) *ErrorHandler {
	// Create error telemetry
	telemetryConfig := common.DefaultErrorTelemetryConfig()
	telemetry := common.NewErrorTelemetry(telemetryConfig)
	
	return &ErrorHandler{
		Interactive: interactive,
		Logger:      logger.DefaultLogger,
		Resolver:    common.NewInteractiveErrorResolver(),
		Telemetry:   telemetry,
	}
}

// HandleError handles an error in the CLI
func HandleError(err error) {
	// Check if it's a NessiError
	var nessiErr *common.NessiError
	if errors.As(err, &nessiErr) {
		// Get suggestions for this error
		suggestions := common.GetSuggestionsForError(nessiErr)

		// Record error in telemetry if enabled
		recordErrorInTelemetry(nessiErr)

		// Print error with code
		error := color.New(color.FgRed, color.Bold)
		error.Printf("\n❌ Error: %s\n", nessiErr.Message)

		// Print error code if available
		if nessiErr.Code != "" {
			fmt.Printf("Code: %s\n", nessiErr.Code)
			fmt.Printf("For more information: nessi test-error %s --info\n", nessiErr.Code)
		}

		// Print details if available
		if nessiErr.Details != "" {
			fmt.Printf("Details: %s\n", nessiErr.Details)
		}

		// Print suggestions if available
		if len(suggestions) > 0 {
			fmt.Println("\nSuggestions:")
			for i, suggestion := range suggestions {
				fmt.Printf("%d. %s\n", i+1, suggestion.Description)
				fmt.Printf("   Solution: %s\n", suggestion.Solution)
				if suggestion.DocumentURL != "" {
					fmt.Printf("   Documentation: %s\n", suggestion.DocumentURL)
				}
			}
		}

		// Check if it's a resolvable error
		var resolvableErr common.ResolvableError
		if errors.As(err, &resolvableErr) {
			// Check if interactive mode is enabled
			if isInteractiveModeEnabled() {
				fmt.Println("\nThis error can be resolved interactively.")
				fmt.Print("Do you want to resolve it now? [Y/n]: ")

				// Read user input
				var input string
				fmt.Scanln(&input)
				input = strings.ToLower(strings.TrimSpace(input))

				// If user wants to resolve the error
				if input == "" || input == "y" || input == "yes" {
					// Resolve the error
					if err := resolvableErr.Resolve(); err != nil {
						fmt.Printf("Failed to resolve the error: %s\n", err)
					} else {
						success := color.New(color.FgGreen, color.Bold)
						success.Println("\n✅ Error resolved successfully!")
						fmt.Println("You can now retry the command.")
					}
				}
			} else {
				fmt.Println("\nThis error can be resolved interactively. Run with --interactive flag to enable interactive resolution.")
			}
		}

		// Check if automatic retry is enabled for this error
		if shouldRetryError(nessiErr) {
			fmt.Println("\nAutomatically retrying operation...")
			// Note: The actual retry logic would be implemented in the command execution
		}
	} else {
		// Print generic error
		error := color.New(color.FgRed, color.Bold)
		error.Printf("\n❌ Error: %s\n", err)
	}
}

// printColoredError prints a NessiError with color
func printColoredError(err *common.NessiError) {
	// Create color printers
	error := color.New(color.FgRed, color.Bold)
	detail := color.New(color.FgYellow)

	// Print error message
	error.Printf("❌ Error: %s\n", err.Message)

	// Print error code if available
	if err.Code != "" {
		fmt.Printf("Code: %s\n", err.Code)
		fmt.Printf("For more information: nessi test-error %s --info\n", err.Code)
	}

	// Print details if available
	if err.Details != "" {
		detail.Printf("Details: %s\n", err.Details)
	}

	// Print suggestion if available
	if err.Suggestion != "" {
		color.New(color.FgGreen).Printf("Suggestion: %s\n", err.Suggestion)
	}
}

// isInteractiveModeEnabled checks if interactive mode is enabled
func isInteractiveModeEnabled() bool {
	// Check if interactive mode is enabled in config
	if viper.IsSet("error_handling.interactive_resolution") {
		return viper.GetBool("error_handling.interactive_resolution")
	}

	// Default to true if not set
	return true
}

// recordErrorInTelemetry records an error in telemetry if enabled
func recordErrorInTelemetry(err *common.NessiError) {
	// Check if telemetry is enabled in config
	if viper.IsSet("error_handling.telemetry.enabled") && viper.GetBool("error_handling.telemetry.enabled") {
		// Get telemetry instance
		telemetryConfig := common.DefaultErrorTelemetryConfig()
		telemetry := common.NewErrorTelemetry(telemetryConfig)

		// Record the error
		telemetry.RecordError(err)
	}
}

// shouldRetryError checks if an error should be automatically retried
func shouldRetryError(err *common.NessiError) bool {
	// Check if retry is enabled in config
	if !viper.IsSet("error_handling.retry.enabled") || !viper.GetBool("error_handling.retry.enabled") {
		return false
	}

	// Check if the error is retryable
	// For now, we'll consider connection errors (N4XX) as retryable
	return strings.HasPrefix(string(err.Code), "N4")
}

// AddErrorHandlingFlags adds error handling flags to a command
func AddErrorHandlingFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("interactive", true, "Enable interactive error resolution")
}

// WrapAllCommands wraps all commands with error handling
func WrapAllCommands(cmd *cobra.Command) {
	// Add common flags to all commands
	cmd.PersistentFlags().Bool("interactive", false, "Enable interactive mode")
	cmd.PersistentFlags().Bool("dry-run", false, "Show what would be done without making changes")

	// Wrap the command's RunE function with error handling
	if cmd.RunE != nil {
		originalRunE := cmd.RunE
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			// Check if this is a dry run
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if dryRun {
				// Set dry run mode in context
				cmd.SetContext(context.WithValue(cmd.Context(), "dry-run", true))
			}

			// Check if interactive mode is enabled via flag
			interactive, _ := cmd.Flags().GetBool("interactive")
			if interactive {
				// Set interactive mode in viper for this session
				viper.Set("error_handling.interactive_resolution", true)
			}

			// Execute the command
			err := originalRunE(cmd, args)
			if err != nil {
				// Handle the error
				HandleError(err)
				return nil // Return nil to prevent Cobra from printing the error
			}

			// If this was a dry run, print a message
			if dryRun {
				fmt.Println("\nThis was a dry run. No changes were made.")
				fmt.Println("To execute these actions, run the command without the --dry-run flag.")
			}

			return nil
		}
	} else if cmd.Run != nil {
		originalRun := cmd.Run
		cmd.Run = func(cmd *cobra.Command, args []string) {
			// Check if this is a dry run
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if dryRun {
				// Set dry run mode in context
				cmd.SetContext(context.WithValue(cmd.Context(), "dry-run", true))
			}

			// Check if interactive mode is enabled via flag
			interactive, _ := cmd.Flags().GetBool("interactive")
			if interactive {
				// Set interactive mode in viper for this session
				viper.Set("error_handling.interactive_resolution", true)
			}

			// Execute the command
			originalRun(cmd, args)

			// If this was a dry run, print a message
			if dryRun {
				fmt.Println("\nThis was a dry run. No changes were made.")
				fmt.Println("To execute these actions, run the command without the --dry-run flag.")
			}
		}
	}

	// Recursively wrap all subcommands
	for _, subCmd := range cmd.Commands() {
		WrapAllCommands(subCmd)
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
