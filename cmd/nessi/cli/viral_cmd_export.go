package cli

import (
	"os"

	"github.com/spf13/cobra"
)

// GetViralCommand returns the viral command for testing
func GetViralCommand() *cobra.Command {
	// Set environment variable to indicate test mode
	os.Setenv("GO_TESTING", "1")
	return viralCmd
}
