package testutil

import (
	"bytes"

	"github.com/nessi-dev/nessi/cmd/nessi/mocks"
	"github.com/spf13/cobra"
)

// ExecuteCommand executes a command for testing and returns its output
func ExecuteCommand(cmd *cobra.Command, args ...string) (string, error) {
	// Reset flags to avoid conflicts
	mocks.ResetFlags()
	// Set up command flags for testing
	mocks.SetupTestCommandFlags(cmd)

	// Execute command using our mock implementation
	return mocks.ExecuteCommand(cmd, args...)
}

// SetupTestCommand sets up a command for testing with output and error buffers
func SetupTestCommand(cmd *cobra.Command) (*bytes.Buffer, *bytes.Buffer) {
	// Reset flags to avoid conflicts
	mocks.ResetFlags()
	// Set up command flags for testing
	mocks.SetupTestCommandFlags(cmd)

	// Set up command for testing
	return mocks.SetupTestCommand(cmd)
}
