package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestViralShareCommand tests the 'nessi viral share' command functionality
func TestViralShareCommand(t *testing.T) {
	// Create a temporary directory for test outputs
	tmpDir, err := os.MkdirTemp("", "nessi-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

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
			// Run the command
			cmd := exec.Command("nessi", tc.args...)
			output, err := cmd.CombinedOutput()
			
			// For now, we'll just check if the command runs without error
			// In a real implementation, we would verify the output and check the generated file
			if err != nil {
				t.Logf("Command output: %s", output)
				t.Fatalf("Command failed: %v", err)
			}
			
			// Check if the output contains the expected string
			assert.Contains(t, string(output), tc.expected)
			
			// Check if the file was created
			outputFile := tc.args[len(tc.args)-1]
			_, err = os.Stat(outputFile)
			assert.NoError(t, err, "Output file should exist")
			
			// Check file content (basic check)
			content, err := os.ReadFile(outputFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), "<html>", "Output should be HTML")
			
			// If this is the hashtag test, check for hashtags in the content
			if strings.Contains(tc.name, "hashtags") {
				assert.Contains(t, string(content), "dataquality", "Hashtags should be in the output")
				assert.Contains(t, string(content), "nessi", "Hashtags should be in the output")
			}
		})
	}
}

// TestViralBadgeCommand tests the 'nessi viral badge' command functionality
func TestViralBadgeCommand(t *testing.T) {
	// Create a temporary directory for test outputs
	tmpDir, err := os.MkdirTemp("", "nessi-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
		format   string
	}{
		{
			name:     "Markdown badge",
			args:     []string{"viral", "badge", "--format", "markdown", "--output", filepath.Join(tmpDir, "badge.md")},
			expected: "Badge generated successfully",
			format:   "markdown",
		},
		{
			name:     "HTML badge",
			args:     []string{"viral", "badge", "--format", "html", "--output", filepath.Join(tmpDir, "badge.html")},
			expected: "Badge generated successfully",
			format:   "html",
		},
		{
			name:     "Custom badge",
			args:     []string{"viral", "badge", "--label", "verified by", "--message", "nessi", "--color", "blue", "--output", filepath.Join(tmpDir, "custom_badge.md")},
			expected: "Badge generated successfully",
			format:   "markdown",
		},
		{
			name:     "Quality score badge",
			args:     []string{"viral", "badge", "--quality-score", "95", "--output", filepath.Join(tmpDir, "score_badge.md")},
			expected: "Badge generated successfully",
			format:   "markdown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run the command
			cmd := exec.Command("nessi", tc.args...)
			output, err := cmd.CombinedOutput()
			
			// For now, we'll just check if the command runs without error
			if err != nil {
				t.Logf("Command output: %s", output)
				t.Fatalf("Command failed: %v", err)
			}
			
			// Check if the output contains the expected string
			assert.Contains(t, string(output), tc.expected)
			
			// Check if the file was created
			outputFile := tc.args[len(tc.args)-1]
			_, err = os.Stat(outputFile)
			assert.NoError(t, err, "Output file should exist")
			
			// Check file content based on format
			content, err := os.ReadFile(outputFile)
			assert.NoError(t, err)
			
			switch tc.format {
			case "markdown":
				assert.Contains(t, string(content), "![")
				assert.Contains(t, string(content), "](")
			case "html":
				assert.Contains(t, string(content), "<img")
				assert.Contains(t, string(content), "src=")
			}
			
			// Check for custom content if applicable
			if strings.Contains(tc.name, "Custom") {
				assert.Contains(t, string(content), "verified by")
				assert.Contains(t, string(content), "nessi")
				assert.Contains(t, string(content), "blue")
			}
			
			// Check for quality score if applicable
			if strings.Contains(tc.name, "Quality") {
				assert.Contains(t, string(content), "95")
			}
		})
	}
}

// TestViralCommunityCommand tests the 'nessi viral community' command functionality
func TestViralCommunityCommand(t *testing.T) {
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
			// Run the command
			cmd := exec.Command("nessi", tc.args...)
			output, err := cmd.CombinedOutput()
			
			// For now, we'll just check if the command runs without error
			if err != nil {
				t.Logf("Command output: %s", output)
				t.Fatalf("Command failed: %v", err)
			}
			
			// Check if the output contains the expected string
			assert.Contains(t, string(output), tc.expected)
			
			// For feedback command, check if it contains a GitHub issue URL
			if strings.Contains(tc.name, "Feedback") {
				assert.Contains(t, string(output), "github.com")
			}
			
			// For contribute commands, check for appropriate content
			if strings.Contains(tc.name, "Contribute") {
				if strings.Contains(tc.name, "beginners") {
					assert.Contains(t, string(output), "beginner")
					assert.Contains(t, string(output), "documentation")
				} else if strings.Contains(tc.name, "advanced") {
					assert.Contains(t, string(output), "advanced")
					assert.Contains(t, string(output), "core")
				}
			}
		})
	}
}

// TestPluginIntegration tests that the viral growth plugins can be loaded and used together
func TestPluginIntegration(t *testing.T) {
	// Get list of plugins
	cmd := exec.Command("nessi", "plugins", "list", "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", output)
		t.Fatalf("Command failed: %v", err)
	}
	
	// Parse the JSON output
	var plugins []map[string]interface{}
	err = json.Unmarshal(output, &plugins)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	// Check if the viral growth plugins are in the list
	foundSocialSharing := false
	foundBadge := false
	foundCommunity := false
	
	for _, plugin := range plugins {
		name, ok := plugin["name"].(string)
		if !ok {
			continue
		}
		
		if strings.Contains(name, "social_sharing") {
			foundSocialSharing = true
		} else if strings.Contains(name, "badge") {
			foundBadge = true
		} else if strings.Contains(name, "community") {
			foundCommunity = true
		}
	}
	
	assert.True(t, foundSocialSharing, "Social sharing plugin should be installed")
	assert.True(t, foundBadge, "Badge plugin should be installed")
	assert.True(t, foundCommunity, "Community engagement plugin should be installed")
	
	// Test that the viral command group is available
	cmd = exec.Command("nessi", "help", "viral")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", output)
		t.Fatalf("Command failed: %v", err)
	}
	
	// Check that the output contains information about all subcommands
	assert.Contains(t, string(output), "share")
	assert.Contains(t, string(output), "badge")
	assert.Contains(t, string(output), "community")
}

// TestErrorHandling tests that the viral growth commands handle errors gracefully
func TestErrorHandling(t *testing.T) {
	// Test cases for commands that should fail
	tests := []struct {
		name          string
		args          []string
		errorContains string
	}{
		{
			name:          "Share with missing table",
			args:          []string{"viral", "share"},
			errorContains: "table name is required",
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
			// Run the command
			cmd := exec.Command("nessi", tc.args...)
			output, err := cmd.CombinedOutput()
			
			// The command should fail
			assert.Error(t, err, "Command should fail")
			
			// Check if the output contains the expected error message
			assert.Contains(t, string(output), tc.errorContains)
		})
	}
}
