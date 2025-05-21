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
		name          string
		args          []string
		errorContains string
	}{
		{
			name:          "Share with missing table",
			args:          []string{"viral", "share"},
			errorContains: "requires a table name",
		},
		{
			name:          "Badge with invalid format",
			args:          []string{"viral", "badge", "--format", "invalid"},
			errorContains: "invalid format",
		},
		{
			name:          "Community feedback with missing text",
			args:          []string{"viral", "community", "feedback", "--type", "feature"},
			errorContains: "feedback text is required",
		},
		{
			name:          "Community contribute with invalid experience",
			args:          []string{"viral", "community", "contribute", "--experience", "invalid"},
			errorContains: "invalid experience level",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set command args
			rootCmd.SetArgs(tc.args)

			// Execute command and expect error
			output := captureOutput(func() {
				err := rootCmd.Execute()
				assert.Error(t, err, "Command should fail")
			})

			// Check error message
			assert.Contains(t, output, tc.errorContains, "Error message should contain %s", tc.errorContains)
		})
	}
}
