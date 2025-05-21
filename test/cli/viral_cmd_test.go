package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/cmd/nessi/cli"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupViralCommand creates a test command with the viral subcommands
func setupViralCommand() *cobra.Command {
	// Create a root command
	rootCmd := &cobra.Command{
		Use:   "nessi",
		Short: "Nessi CLI",
	}

	// Add the viral command
	rootCmd.AddCommand(cli.GetViralCommand())

	return rootCmd
}

// captureOutput captures stdout during command execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestViralShareCommand(t *testing.T) {
	// Create a temporary directory for test outputs
	tmpDir, err := os.MkdirTemp("", "nessi-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set up the command
	rootCmd := setupViralCommand()

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Basic share command",
			args:     []string{"viral", "share", "test_table", "--output", filepath.Join(tmpDir, "report.html")},
			expected: "Report generated successfully",
		},
		{
			name:     "Share with custom title",
			args:     []string{"viral", "share", "test_table", "--title", "Custom Report", "--output", filepath.Join(tmpDir, "titled_report.html")},
			expected: "Report generated successfully",
		},
		{
			name:     "Share with hashtags",
			args:     []string{"viral", "share", "test_table", "--hashtags", "dataquality,nessi", "--output", filepath.Join(tmpDir, "hashtag_report.html")},
			expected: "Report generated successfully",
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

			// Check output
			assert.Contains(t, output, tc.expected)

			// Check if file was created
			outputFile := tc.args[len(tc.args)-1]
			_, err := os.Stat(outputFile)
			assert.NoError(t, err, "Output file should exist")
		})
	}
}

func TestViralBadgeCommand(t *testing.T) {
	// Create a temporary directory for test outputs
	tmpDir, err := os.MkdirTemp("", "nessi-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Set up the command
	rootCmd := setupViralCommand()

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Markdown badge",
			args:     []string{"viral", "badge", "--format", "markdown", "--output", filepath.Join(tmpDir, "badge.md")},
			expected: "Badge generated successfully",
		},
		{
			name:     "HTML badge",
			args:     []string{"viral", "badge", "--format", "html", "--output", filepath.Join(tmpDir, "badge.html")},
			expected: "Badge generated successfully",
		},
		{
			name:     "Custom badge",
			args:     []string{"viral", "badge", "--label", "verified by", "--message", "nessi", "--color", "blue", "--output", filepath.Join(tmpDir, "custom_badge.md")},
			expected: "Badge generated successfully",
		},
		{
			name:     "Quality score badge",
			args:     []string{"viral", "badge", "--quality-score", "95", "--output", filepath.Join(tmpDir, "score_badge.md")},
			expected: "Badge generated successfully",
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

			// Check output
			assert.Contains(t, output, tc.expected)

			// Check if file was created
			outputFile := tc.args[len(tc.args)-1]
			_, err := os.Stat(outputFile)
			assert.NoError(t, err, "Output file should exist")
		})
	}
}

func TestViralCommunityCommand(t *testing.T) {
	// Set up the command
	rootCmd := setupViralCommand()

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Feedback command",
			args:     []string{"viral", "community", "feedback", "--type", "feature", "--text", "It would be great to have..."},
			expected: "Feedback submitted successfully",
		},
		{
			name:     "Contribute command for beginners",
			args:     []string{"viral", "community", "contribute", "--experience", "beginner"},
			expected: "Contribution suggestions for beginners",
		},
		{
			name:     "Contribute command for advanced",
			args:     []string{"viral", "community", "contribute", "--experience", "advanced"},
			expected: "Contribution suggestions for advanced",
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

			// Check output
			assert.Contains(t, output, tc.expected)

			// For feedback command, check if it contains a GitHub issue URL
			if tc.name == "Feedback command" {
				assert.Contains(t, output, "github.com")
			}
		})
	}
}

func TestViralCommandHelp(t *testing.T) {
	// Set up the command
	rootCmd := setupViralCommand()

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name: "Viral command help",
			args: []string{"viral", "--help"},
			expected: []string{
				"share",
				"badge",
				"community",
			},
		},
		{
			name: "Share command help",
			args: []string{"viral", "share", "--help"},
			expected: []string{
				"--output",
				"--title",
				"--description",
				"--hashtags",
			},
		},
		{
			name: "Badge command help",
			args: []string{"viral", "badge", "--help"},
			expected: []string{
				"--format",
				"--label",
				"--message",
				"--color",
				"--style",
				"--quality-score",
			},
		},
		{
			name: "Community command help",
			args: []string{"viral", "community", "--help"},
			expected: []string{
				"feedback",
				"contribute",
			},
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

			// Check that output contains all expected strings
			for _, expected := range tc.expected {
				assert.Contains(t, output, expected, "Help output should contain %s", expected)
			}
		})
	}
}

func TestViralCommandErrors(t *testing.T) {
	// Set up the command
	rootCmd := setupViralCommand()

	// Test cases
	tests := []struct {
		name              string
		args              []string
		errorContains     string
		errorCodeContains string // New field for error codes
	}{
		{
			name:              "Share with missing table",
			args:              []string{"viral", "share"},
			errorContains:     "requires a table name",
			errorCodeContains: "", // No specific error code expected
		},
		{
			name:              "Badge with invalid format",
			args:              []string{"viral", "badge", "--format", "invalid"},
			errorContains:     "Invalid badge format",
			errorCodeContains: "V102", // Badge error code
		},
		{
			name:              "Community feedback with missing text",
			args:              []string{"viral", "community", "feedback", "--type", "feature"},
			errorContains:     "Feedback text cannot be empty",
			errorCodeContains: "V103", // Community error code
		},
		{
			name:              "Community contribute with invalid experience",
			args:              []string{"viral", "community", "contribute", "--experience", "invalid"},
			errorContains:     "invalid experience level",
			errorCodeContains: "", // No specific error code expected
		},
		// New test cases for error handling
		{
			name:              "Badge with invalid color",
			args:              []string{"viral", "badge", "--color", "invalid-color"},
			errorContains:     "may not be recognized", // This is a warning, not an error
			errorCodeContains: "",
		},
		{
			name:              "Share with invalid output path",
			args:              []string{"viral", "share", "test_table", "--output", "/invalid/path/report.html"},
			errorContains:     "output path",
			errorCodeContains: "V101", // Share error code
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set command args
			rootCmd.SetArgs(tc.args)

			// Execute command and expect error
			output := captureOutput(func() {
				err := rootCmd.Execute()
				// For the invalid color test, we expect a warning but not an error
				if tc.name == "Badge with invalid color" {
					assert.NoError(t, err, "Command should succeed with a warning")
				} else {
					assert.Error(t, err, "Command should fail")
				}
			})

			// Check error message
			assert.Contains(t, output, tc.errorContains, "Error message should contain %s", tc.errorContains)

			// Check error code if specified
			if tc.errorCodeContains != "" {
				assert.Contains(t, output, tc.errorCodeContains, "Error code should contain %s", tc.errorCodeContains)
			}
		})
	}
}
