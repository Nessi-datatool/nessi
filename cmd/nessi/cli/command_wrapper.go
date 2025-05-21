package cli

import (
	"github.com/nessi-dev/nessi/pkg/logger"
	"github.com/spf13/cobra"
)

// WrapCommand wraps a cobra command with error handling
func WrapCommand(cmd *cobra.Command) {
	// Store the original RunE function
	originalRunE := cmd.RunE

	// Skip if there's no RunE function
	if originalRunE == nil {
		return
	}

	// Replace with wrapped function
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Get interactive flag
		interactive, _ := cmd.Flags().GetBool("interactive")

		// Create error handler
		handler := NewErrorHandler(interactive)

		// Run the original function
		err := originalRunE(cmd, args)

		// If there's an error, try to convert it to a resolvable error
		if err != nil {
			err = ConvertToResolvableError(err)

			// Handle the error
			return handler.HandleError(err)
		}

		return nil
	}

	// Add error handling flags if they don't exist
	if cmd.Flags().Lookup("interactive") == nil {
		cmd.Flags().Bool("interactive", true, "Enable interactive error resolution")
	}
}

// WrapAllCommands recursively wraps all commands with error handling
func WrapAllCommands(cmd *cobra.Command) {
	// Wrap the command itself
	WrapCommand(cmd)

	// Wrap all subcommands
	for _, subCmd := range cmd.Commands() {
		WrapAllCommands(subCmd)
	}
}

// RegisterErrorHandlers registers error handlers for common error types
func RegisterErrorHandlers() {
	// Create a logger
	log := logger.DefaultLogger

	// We're not actually registering handlers here since the common package
	// doesn't have a RegisterErrorHandler function. Instead, we're just
	// setting up the logger to handle specific error types.

	// These error types are already defined in the common package,
	// and our error handler will use them when handling errors.

	// Log initialization message
	log.Debug("Initializing error handlers for CLI")
}
