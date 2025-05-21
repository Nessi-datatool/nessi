package plugins

import (
	"encoding/json"
	"testing"

	"github.com/nessi-dev/nessi/examples/plugins/social_sharing_plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocialSharingPlugin_Initialize(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}

	// Initialize the plugin
	err := plugin.Initialize()
	require.NoError(t, err, "Plugin initialization should not fail")

	// Verify metadata
	metadata := plugin.GetMetadata()
	assert.Equal(t, "social_sharing_plugin", metadata.Name, "Plugin name should match")
	assert.NotEmpty(t, metadata.Version, "Plugin version should not be empty")
	assert.NotEmpty(t, metadata.Description, "Plugin description should not be empty")
	assert.NotEmpty(t, metadata.Author, "Plugin author should not be empty")
	assert.Contains(t, metadata.Capabilities, "social_sharing", "Plugin should have social_sharing capability")
}

func TestSocialSharingPlugin_EnhanceReport(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"title":       "Test Report",
		"description": "A test report for unit testing",
		"report_url":  "https://example.com/reports/test-report",
		"hashtags":    "nessi,dataquality,testing",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the enhance_report command
	result, err := plugin.Execute([]string{"enhance_report"}, data)
	require.NoError(t, err, "Enhance report should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for HTML content
	htmlContent, ok := resultMap["html_content"].(string)
	require.True(t, ok, "Result should contain HTML content")
	assert.Contains(t, htmlContent, "<div class=\"social-share\">", "HTML should contain social share div")

	// Check for sharing data
	shareData, ok := resultMap["share_data"].(map[string]interface{})
	require.True(t, ok, "Result should contain share data")
	assert.Equal(t, "Test Report", shareData["title"], "Title should match input")
}

func TestSocialSharingPlugin_GenerateShareLinks(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"title":       "Test Report",
		"description": "A test report for unit testing",
		"report_url":  "https://example.com/reports/test-report",
		"hashtags":    "nessi,dataquality,testing",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the generate_share_links command
	result, err := plugin.Execute([]string{"generate_share_links"}, data)
	require.NoError(t, err, "Generate share links should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for sharing links
	assert.Contains(t, resultMap, "twitter_link", "Result should contain Twitter link")
	assert.Contains(t, resultMap, "linkedin_link", "Result should contain LinkedIn link")
	assert.Contains(t, resultMap, "email_link", "Result should contain email link")

	// Verify link content
	twitterLink, ok := resultMap["twitter_link"].(string)
	require.True(t, ok, "Twitter link should be a string")
	assert.Contains(t, twitterLink, "twitter.com", "Twitter link should contain twitter.com")
	assert.Contains(t, twitterLink, "Test Report", "Twitter link should contain the title")
	assert.Contains(t, twitterLink, "nessi", "Twitter link should contain hashtags")
}

func TestSocialSharingPlugin_GenerateQRCode(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"report_url": "https://example.com/reports/test-report",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the generate_qr_code command
	result, err := plugin.Execute([]string{"generate_qr_code"}, data)
	require.NoError(t, err, "Generate QR code should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for QR code data
	qrCodeData, ok := resultMap["qr_code_data"].(string)
	require.True(t, ok, "Result should contain QR code data")
	assert.NotEmpty(t, qrCodeData, "QR code data should not be empty")

	// Check for QR code HTML
	qrCodeHTML, ok := resultMap["qr_code_html"].(string)
	require.True(t, ok, "Result should contain QR code HTML")
	assert.Contains(t, qrCodeHTML, "<img", "QR code HTML should contain img tag")
	assert.Contains(t, qrCodeHTML, "data:image/png;base64,", "QR code HTML should contain base64 image data")
}

func TestSocialSharingPlugin_InvalidCommand(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute an invalid command
	_, err := plugin.Execute([]string{"invalid_command"}, nil)
	assert.Error(t, err, "Invalid command should return an error")
	assert.Contains(t, err.Error(), "unknown command", "Error should mention unknown command")
}

func TestSocialSharingPlugin_EmptyCommand(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute with empty command list
	_, err := plugin.Execute([]string{}, nil)
	assert.Error(t, err, "Empty command should return an error")
	assert.Contains(t, err.Error(), "no command specified", "Error should mention no command specified")
}

func TestSocialSharingPlugin_InvalidData(t *testing.T) {
	// Create a new plugin instance
	plugin := &social_sharing_plugin.SocialSharingPlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute with invalid data
	_, err := plugin.Execute([]string{"enhance_report"}, "not a map")
	assert.Error(t, err, "Invalid data should return an error")
}
