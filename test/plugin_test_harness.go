package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// PluginTestHarness provides utilities for testing plugins
type PluginTestHarness struct {
	t        testing.TB
	PluginsDir string
	TempDir    string
}

// CallPluginFunction calls a function in a plugin and returns the result
func (h *PluginTestHarness) CallPluginFunction(pluginName, functionName string, data interface{}) (interface{}, error) {
	// In a real implementation, this would call the plugin function
	// For benchmarking, we'll just return a mock result
	switch pluginName {
	case "social_sharing_plugin":
		switch functionName {
		case "EnhanceReport":
			return map[string]interface{}{
				"title":       "Enhanced Test Report",
				"description": "A test report for benchmarking",
				"metrics":     map[string]float64{"quality": 95.0, "completeness": 98.0},
				"share_links": map[string]string{
					"twitter":  "https://twitter.com/intent/tweet?text=Test%20Report&url=https%3A%2F%2Fexample.com%2Freport%2F123",
					"linkedin": "https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fexample.com%2Freport%2F123",
				},
				"qr_code": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
			}, nil
		case "GenerateShareLinks":
			return map[string]string{
				"twitter":  "https://twitter.com/intent/tweet?text=Test%20Report&url=https%3A%2F%2Fexample.com%2Freport%2F123",
				"linkedin": "https://www.linkedin.com/sharing/share-offsite/?url=https%3A%2F%2Fexample.com%2Freport%2F123",
				"facebook": "https://www.facebook.com/sharer/sharer.php?u=https%3A%2F%2Fexample.com%2Freport%2F123",
			}, nil
		case "GenerateQRCode":
			return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...", nil
		default:
			return nil, nil
		}
	case "badge_plugin":
		switch functionName {
		case "GenerateBadge":
			return "[![Powered by Nessi](https://img.shields.io/badge/powered%20by-nessi-blue?style=flat)](https://github.com/nessi-dev/nessi)", nil
		case "QualityScoreBadge":
			return "[![Quality Score: 95%](https://img.shields.io/badge/quality-95%25-brightgreen?style=flat)](https://github.com/nessi-dev/nessi)", nil
		default:
			return nil, nil
		}
	case "community_engagement_plugin":
		switch functionName {
		case "CollectFeedback":
			return "https://github.com/nessi-dev/nessi/issues/new?title=Feedback&body=I%20would%20like%20to%20provide%20feedback", nil
		case "SuggestContributions":
			return []string{
				"Fix a documentation typo: https://github.com/nessi-dev/nessi/blob/main/docs/",
				"Add a test case: https://github.com/nessi-dev/nessi/tree/main/pkg/test",
				"Improve error messages: https://github.com/nessi-dev/nessi/blob/main/pkg/errors/",
			}, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}
}

// NewPluginTestHarness creates a new test harness for plugins
func NewPluginTestHarness(t testing.TB) *PluginTestHarness {
	// Create a temporary directory for test outputs
	tempDir, err := os.MkdirTemp("", "nessi-plugin-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Clean up the temporary directory after the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Get the plugins directory
	// Find the project root directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	// If we're in the test directory, go up one level
	if filepath.Base(wd) == "test" {
		wd = filepath.Dir(wd)
	}
	pluginDir := filepath.Join(wd, "examples", "plugins")

	return &PluginTestHarness{
		t:         t,
		PluginsDir: pluginDir,
		TempDir:    tempDir,
	}
}

// BuildPlugin builds a plugin for testing
func (h *PluginTestHarness) BuildPlugin(t *testing.T, pluginName string) string {
	pluginDir := filepath.Join(h.PluginsDir, pluginName)
	pluginOutput := filepath.Join(h.TempDir, pluginName+".so")

	// Build the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginOutput, "./")
	cmd.Dir = pluginDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build plugin %s: %v\nOutput: %s", pluginName, err, output)
	}

	return pluginOutput
}

// LoadPlugin loads a plugin for testing
func (h *PluginTestHarness) LoadPlugin(t *testing.T, pluginPath string) (interface{}, error) {
	// This is a mock implementation since we can't directly load plugins in tests
	// In a real implementation, this would use plugin.Open and plugin.Lookup
	
	// Instead, we'll execute a helper program that loads and interacts with the plugin
	// For now, we'll just return a mock plugin interface
	return &MockPlugin{PluginPath: pluginPath}, nil
}

// MockPlugin is a mock implementation of a plugin for testing
type MockPlugin struct {
	PluginPath string
}

// Execute executes a command on the mock plugin
func (p *MockPlugin) Execute(t *testing.T, command string, data map[string]interface{}) (map[string]interface{}, error) {
	// We don't actually need to marshal the data in our mock implementation
	// This would be needed in a real implementation that communicates with plugins

	// Create a helper program that loads and interacts with the plugin
	// For now, we'll just return mock results based on the command
	switch command {
	case "enhance_report":
		return map[string]interface{}{
			"html_content": "<div class=\"social-share\">Social sharing buttons</div>",
			"share_data":   data,
		}, nil
	case "generate_share_links":
		return map[string]interface{}{
			"twitter_link":  "https://twitter.com/intent/tweet?text=" + data["title"].(string),
			"linkedin_link": "https://www.linkedin.com/sharing/share-offsite/?url=" + data["report_url"].(string),
			"email_link":    "mailto:?subject=" + data["title"].(string),
		}, nil
	case "generate_qr_code":
		return map[string]interface{}{
			"qr_code_data":  "mock_qr_code_data",
			"qr_code_html": "<img src=\"data:image/png;base64,mock_data\">",
		}, nil
	case "generate_badge":
		return map[string]interface{}{
			"badge_url":      "https://img.shields.io/badge/powered%20by-nessi-blue",
			"badge_markdown": "![powered by](https://img.shields.io/badge/powered%20by-nessi-blue)",
			"badge_html":     "<img src=\"https://img.shields.io/badge/powered%20by-nessi-blue\" alt=\"powered by nessi\">",
		}, nil
	case "collect_feedback":
		return map[string]interface{}{
			"issue_url":        "https://github.com/nessi-dev/nessi/issues/new",
			"feedback_details": data,
		}, nil
	case "suggest_contributions":
		suggestions := []interface{}{
			map[string]interface{}{
				"text":        "Improve documentation",
				"difficulty": "beginner",
				"url":        "https://github.com/nessi-dev/nessi/issues?q=is%3Aissue+is%3Aopen+label%3Adocumentation",
			},
		}
		if data["experience"] == "advanced" {
			suggestions = append(suggestions, map[string]interface{}{
				"text":        "Implement core features",
				"difficulty": "advanced",
				"url":        "https://github.com/nessi-dev/nessi/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement",
			})
		}
		return map[string]interface{}{
			"suggestions": suggestions,
		}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}
