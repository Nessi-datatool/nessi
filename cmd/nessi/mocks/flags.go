package mocks

import (
	"flag"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	// Ensure we only initialize once
	initOnce sync.Once

	// Original command line arguments
	originalArgs []string
)

// InitMockFlags initializes mock flags for testing
func InitMockFlags() {
	initOnce.Do(func() {
		// Save original arguments
		originalArgs = os.Args

		// Set minimal arguments to avoid parsing issues
		os.Args = []string{"nessi-test"}

		// Reset the flag package
		flag.CommandLine = flag.NewFlagSet("nessi-test", flag.ExitOnError)
	})
}

// ResetFlags resets the flags for testing
func ResetFlags() {
	// Reset the flag package
	flag.CommandLine = flag.NewFlagSet("nessi-test", flag.ExitOnError)
}

// RestoreFlags restores the original flags after testing
func RestoreFlags() {
	// Restore original arguments
	os.Args = originalArgs
}

// SetupTestCommandFlags sets up a command's flags for testing
func SetupTestCommandFlags(cmd *cobra.Command) {
	// Reset the command's flags
	cmd.ResetFlags()
	cmd.ResetCommands()
}
