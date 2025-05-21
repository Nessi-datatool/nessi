package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/cmd/nessi/cli"
	"github.com/stretchr/testify/assert"
)

// TestProgressIndicators tests the progress indicators in the viral commands
func TestProgressIndicators(t *testing.T) {
	// Set up the command
	rootCmd := setupViralCommand()

	// Test cases
	tests := []struct {
		name             string
		args             []string
		progressContains string
		successContains  string
	}{
		{
			name:             "Share command progress",
			args:             []string{"viral", "share", "test_table"},
			progressContains: "Generating shareable report",
			successContains:  "Successfully generated",
		},
		{
			name:             "Badge command progress",
			args:             []string{"viral", "badge"},
			progressContains: "Generating 'Powered by Nessi' badge",
			successContains:  "Badge generated successfully",
		},
		{
			name:             "Community feedback progress",
			args:             []string{"viral", "community", "feedback", "--text", "This is a test feedback"},
			progressContains: "Submitting feedback",
			successContains:  "Feedback submitted successfully",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set command args
			rootCmd.SetArgs(tc.args)

			// Capture output and execute command
			output := captureOutput(func() {
				err := rootCmd.Execute()
				assert.NoError(t, err)
			})

			// Check for progress indicator
			assert.Contains(t, output, tc.progressContains, "Output should contain progress indicator: %s", tc.progressContains)

			// Check for success message
			assert.Contains(t, output, tc.successContains, "Output should contain success message: %s", tc.successContains)

			// Check for spinner characters (this is a basic check since the actual spinner is hard to test)
			assert.True(t, strings.Contains(output, "⠋") || strings.Contains(output, "⠙") ||
				strings.Contains(output, "⠹") || strings.Contains(output, "⠸") ||
				strings.Contains(output, "⠼") || strings.Contains(output, "⠴") ||
				strings.Contains(output, "⠦") || strings.Contains(output, "⠧") ||
				strings.Contains(output, "⠇") || strings.Contains(output, "⠏"),
				"Output should contain spinner characters")
		})
	}
}

// TestShowProgress tests the ShowProgress function directly
func TestShowProgress(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		message  string
		duration time.Duration
	}{
		{
			name:     "Short progress",
			message:  "Short progress test",
			duration: 100 * time.Millisecond,
		},
		{
			name:     "Medium progress",
			message:  "Medium progress test",
			duration: 200 * time.Millisecond,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output
			output := captureOutput(func() {
				cli.ShowProgress(tc.message, tc.duration)
			})

			// Check that output contains the message
			assert.Contains(t, output, tc.message, "Output should contain the progress message")
		})
	}
}

// TestShowSuccess tests the ShowSuccess function directly
func TestShowSuccess(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Simple success",
			message: "Operation completed successfully",
		},
		{
			name:    "Detailed success",
			message: "All files processed successfully (10/10)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output
			output := captureOutput(func() {
				cli.ShowSuccess(tc.message)
			})

			// Check that output contains the message
			assert.Contains(t, output, tc.message, "Output should contain the success message")

			// Check for checkmark symbol
			assert.Contains(t, output, "✅", "Output should contain checkmark symbol")
		})
	}
}

// TestShowWarning tests the ShowWarning function directly
func TestShowWarning(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Simple warning",
			message: "Operation completed with warnings",
		},
		{
			name:    "Detailed warning",
			message: "Some files could not be processed (8/10)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output
			output := captureOutput(func() {
				cli.ShowWarning(tc.message)
			})

			// Check that output contains the message
			assert.Contains(t, output, tc.message, "Output should contain the warning message")

			// Check for warning symbol
			assert.Contains(t, output, "⚠️", "Output should contain warning symbol")
		})
	}
}
