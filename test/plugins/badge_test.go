package plugins

import (
	"encoding/json"
	"testing"

	"github.com/nessi-dev/nessi/examples/plugins/badge_plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadgePlugin_Initialize(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}

	// Initialize the plugin
	err := plugin.Initialize()
	require.NoError(t, err, "Plugin initialization should not fail")

	// Verify metadata
	metadata := plugin.GetMetadata()
	assert.Equal(t, "badge_plugin", metadata.Name, "Plugin name should match")
	assert.NotEmpty(t, metadata.Version, "Plugin version should not be empty")
	assert.NotEmpty(t, metadata.Description, "Plugin description should not be empty")
	assert.NotEmpty(t, metadata.Author, "Plugin author should not be empty")
	assert.Contains(t, metadata.Capabilities, "badge_generation", "Plugin should have badge_generation capability")
}

func TestBadgePlugin_GenerateBadge(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"label":   "powered by",
		"message": "nessi",
		"color":   "blue",
		"style":   "flat",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the generate_badge command
	result, err := plugin.Execute([]string{"generate_badge"}, data)
	require.NoError(t, err, "Generate badge should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for badge URL
	badgeURL, ok := resultMap["badge_url"].(string)
	require.True(t, ok, "Result should contain badge URL")
	assert.Contains(t, badgeURL, "shields.io", "Badge URL should use shields.io")
	assert.Contains(t, badgeURL, "powered%20by", "Badge URL should contain the label")
	assert.Contains(t, badgeURL, "nessi", "Badge URL should contain the message")
	assert.Contains(t, badgeURL, "blue", "Badge URL should contain the color")

	// Check for badge markdown
	badgeMarkdown, ok := resultMap["badge_markdown"].(string)
	require.True(t, ok, "Result should contain badge markdown")
	assert.Contains(t, badgeMarkdown, "![powered by](")
	assert.Contains(t, badgeMarkdown, badgeURL)

	// Check for badge HTML
	badgeHTML, ok := resultMap["badge_html"].(string)
	require.True(t, ok, "Result should contain badge HTML")
	assert.Contains(t, badgeHTML, "<img")
	assert.Contains(t, badgeHTML, "src=")
	assert.Contains(t, badgeHTML, badgeURL)
}

func TestBadgePlugin_GenerateQualityScoreBadge(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"quality_score": 95,
		"style":         "flat",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the generate_badge command
	result, err := plugin.Execute([]string{"generate_badge"}, data)
	require.NoError(t, err, "Generate quality score badge should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for badge URL
	badgeURL, ok := resultMap["badge_url"].(string)
	require.True(t, ok, "Result should contain badge URL")
	assert.Contains(t, badgeURL, "quality", "Badge URL should contain 'quality'")
	assert.Contains(t, badgeURL, "95", "Badge URL should contain the quality score")

	// Check badge color based on score
	assert.Contains(t, badgeURL, "brightgreen", "High quality score should have brightgreen color")
}

func TestBadgePlugin_GetBadgeMarkdown(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"label":   "powered by",
		"message": "nessi",
		"color":   "blue",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the get_badge_markdown command
	result, err := plugin.Execute([]string{"get_badge_markdown"}, data)
	require.NoError(t, err, "Get badge markdown should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for markdown content
	markdown, ok := resultMap["markdown"].(string)
	require.True(t, ok, "Result should contain markdown")
	assert.Contains(t, markdown, "![powered by](")
	assert.Contains(t, markdown, "shields.io")
	assert.Contains(t, markdown, "powered%20by-nessi-blue")
}

func TestBadgePlugin_GetBadgeHTML(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"label":   "powered by",
		"message": "nessi",
		"color":   "blue",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the get_badge_html command
	result, err := plugin.Execute([]string{"get_badge_html"}, data)
	require.NoError(t, err, "Get badge HTML should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for HTML content
	html, ok := resultMap["html"].(string)
	require.True(t, ok, "Result should contain HTML")
	assert.Contains(t, html, "<img")
	assert.Contains(t, html, "src=")
	assert.Contains(t, html, "shields.io")
	assert.Contains(t, html, "alt=\"powered by nessi\"")
}

func TestBadgePlugin_EnhanceReportWithBadge(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"report_html": "<html><body><div id=\"report-content\">Test Report</div></body></html>",
		"label":      "powered by",
		"message":    "nessi",
		"color":      "blue",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the enhance_report_with_badge command
	result, err := plugin.Execute([]string{"enhance_report_with_badge"}, data)
	require.NoError(t, err, "Enhance report with badge should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for enhanced HTML
	enhancedHTML, ok := resultMap["enhanced_html"].(string)
	require.True(t, ok, "Result should contain enhanced HTML")
	assert.Contains(t, enhancedHTML, "<img")
	assert.Contains(t, enhancedHTML, "shields.io")
	assert.Contains(t, enhancedHTML, "Test Report")
}

func TestBadgePlugin_InvalidCommand(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute an invalid command
	_, err := plugin.Execute([]string{"invalid_command"}, nil)
	assert.Error(t, err, "Invalid command should return an error")
	assert.Contains(t, err.Error(), "unknown command", "Error should mention unknown command")
}

func TestBadgePlugin_EmptyCommand(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute with empty command list
	_, err := plugin.Execute([]string{}, nil)
	assert.Error(t, err, "Empty command should return an error")
	assert.Contains(t, err.Error(), "no command specified", "Error should mention no command specified")
}

func TestBadgePlugin_InvalidData(t *testing.T) {
	// Create a new plugin instance
	plugin := &badge_plugin.BadgePlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute with invalid data
	_, err := plugin.Execute([]string{"generate_badge"}, "not a map")
	assert.Error(t, err, "Invalid data should return an error")
}
