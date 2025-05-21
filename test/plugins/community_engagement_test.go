package plugins

import (
	"encoding/json"
	"testing"

	"github.com/nessi-dev/nessi/examples/plugins/community_engagement_plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunityEngagementPlugin_Initialize(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}

	// Initialize the plugin
	err := plugin.Initialize()
	require.NoError(t, err, "Plugin initialization should not fail")

	// Verify metadata
	metadata := plugin.GetMetadata()
	assert.Equal(t, "community_engagement_plugin", metadata.Name, "Plugin name should match")
	assert.NotEmpty(t, metadata.Version, "Plugin version should not be empty")
	assert.NotEmpty(t, metadata.Description, "Plugin description should not be empty")
	assert.NotEmpty(t, metadata.Author, "Plugin author should not be empty")
	assert.Contains(t, metadata.Capabilities, "community_engagement", "Plugin should have community_engagement capability")
}

func TestCommunityEngagementPlugin_CollectFeedback(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"type":  "feature",
		"text":  "It would be great to have a dashboard for quality metrics",
		"name":  "Test User",
		"email": "test@example.com",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the collect_feedback command
	result, err := plugin.Execute([]string{"collect_feedback"}, data)
	require.NoError(t, err, "Collect feedback should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for issue URL
	issueURL, ok := resultMap["issue_url"].(string)
	require.True(t, ok, "Result should contain issue URL")
	assert.Contains(t, issueURL, "github.com", "Issue URL should point to GitHub")

	// Check for feedback details
	feedbackDetails, ok := resultMap["feedback_details"].(map[string]interface{})
	require.True(t, ok, "Result should contain feedback details")
	assert.Equal(t, "feature", feedbackDetails["type"], "Feedback type should match input")
	assert.Equal(t, "It would be great to have a dashboard for quality metrics", feedbackDetails["text"], "Feedback text should match input")
	assert.Equal(t, "Test User", feedbackDetails["name"], "Feedback name should match input")
}

func TestCommunityEngagementPlugin_SuggestContributions(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data for beginners
	beginnerData := map[string]interface{}{
		"experience":     "beginner",
		"github_profile": "testuser",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(beginnerData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the suggest_contributions command for beginners
	result, err := plugin.Execute([]string{"suggest_contributions"}, data)
	require.NoError(t, err, "Suggest contributions for beginners should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for suggestions
	suggestions, ok := resultMap["suggestions"].([]interface{})
	require.True(t, ok, "Result should contain suggestions array")
	assert.NotEmpty(t, suggestions, "Suggestions should not be empty")

	// Check for beginner-appropriate suggestions
	suggestionText := ""
	for _, suggestion := range suggestions {
		if suggestionMap, ok := suggestion.(map[string]interface{}); ok {
			if text, ok := suggestionMap["text"].(string); ok {
				suggestionText += text
			}
		}
	}
	assert.Contains(t, suggestionText, "documentation", "Beginner suggestions should mention documentation")

	// Create test data for advanced users
	advancedData := map[string]interface{}{
		"experience":     "advanced",
		"github_profile": "testuser",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err = json.Marshal(advancedData)
	require.NoError(t, err)

	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the suggest_contributions command for advanced users
	result, err = plugin.Execute([]string{"suggest_contributions"}, data)
	require.NoError(t, err, "Suggest contributions for advanced users should not fail")

	// Verify result
	resultMap, ok = result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for suggestions
	suggestions, ok = resultMap["suggestions"].([]interface{})
	require.True(t, ok, "Result should contain suggestions array")
	assert.NotEmpty(t, suggestions, "Suggestions should not be empty")

	// Check for advanced-appropriate suggestions
	suggestionText = ""
	for _, suggestion := range suggestions {
		if suggestionMap, ok := suggestion.(map[string]interface{}); ok {
			if text, ok := suggestionMap["text"].(string); ok {
				suggestionText += text
			}
		}
	}
	assert.Contains(t, suggestionText, "core", "Advanced suggestions should mention core functionality")
}

func TestCommunityEngagementPlugin_EnhanceReportWithCommunity(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}
	require.NoError(t, plugin.Initialize())

	// Create test data
	testData := map[string]interface{}{
		"report_html": "<html><body><div id=\"report-content\">Test Report</div></body></html>",
	}

	// Convert to JSON and back to simulate how it would be passed in real usage
	jsonData, err := json.Marshal(testData)
	require.NoError(t, err)

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	require.NoError(t, err)

	// Execute the enhance_report_with_community command
	result, err := plugin.Execute([]string{"enhance_report_with_community"}, data)
	require.NoError(t, err, "Enhance report with community should not fail")

	// Verify result
	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "Result should be a map")

	// Check for enhanced HTML
	enhancedHTML, ok := resultMap["enhanced_html"].(string)
	require.True(t, ok, "Result should contain enhanced HTML")
	assert.Contains(t, enhancedHTML, "community", "Enhanced HTML should contain community section")
	assert.Contains(t, enhancedHTML, "contribute", "Enhanced HTML should mention contributions")
	assert.Contains(t, enhancedHTML, "feedback", "Enhanced HTML should mention feedback")
	assert.Contains(t, enhancedHTML, "Test Report", "Enhanced HTML should contain original report content")
}

func TestCommunityEngagementPlugin_InvalidCommand(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute an invalid command
	_, err := plugin.Execute([]string{"invalid_command"}, nil)
	assert.Error(t, err, "Invalid command should return an error")
	assert.Contains(t, err.Error(), "unknown command", "Error should mention unknown command")
}

func TestCommunityEngagementPlugin_EmptyCommand(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute with empty command list
	_, err := plugin.Execute([]string{}, nil)
	assert.Error(t, err, "Empty command should return an error")
	assert.Contains(t, err.Error(), "no command specified", "Error should mention no command specified")
}

func TestCommunityEngagementPlugin_InvalidData(t *testing.T) {
	// Create a new plugin instance
	plugin := &community_engagement_plugin.CommunityEngagementPlugin{}
	require.NoError(t, plugin.Initialize())

	// Execute with invalid data
	_, err := plugin.Execute([]string{"collect_feedback"}, "not a map")
	assert.Error(t, err, "Invalid data should return an error")
}
