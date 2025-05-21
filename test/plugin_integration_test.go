package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSocialSharingPlugin tests the social sharing plugin functionality
func TestSocialSharingPlugin(t *testing.T) {
	// Create a test harness
	harness := NewPluginTestHarness(t)

	// Test enhance_report command
	t.Run("EnhanceReport", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "social_sharing_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"title":       "Test Report",
			"description": "A test report for unit testing",
			"report_url":  "https://example.com/reports/test-report",
			"hashtags":    "nessi,dataquality,testing",
		}

		// Execute the enhance_report command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "enhance_report", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result["html_content"], "social-share")
		shareData, ok := result["share_data"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "Test Report", shareData["title"])
	})

	// Test generate_share_links command
	t.Run("GenerateShareLinks", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "social_sharing_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"title":       "Test Report",
			"description": "A test report for unit testing",
			"report_url":  "https://example.com/reports/test-report",
			"hashtags":    "nessi,dataquality,testing",
		}

		// Execute the generate_share_links command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "generate_share_links", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "twitter_link")
		assert.Contains(t, result, "linkedin_link")
		assert.Contains(t, result, "email_link")
		assert.Contains(t, result["twitter_link"], "Test Report")
	})

	// Test generate_qr_code command
	t.Run("GenerateQRCode", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "social_sharing_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"report_url": "https://example.com/reports/test-report",
		}

		// Execute the generate_qr_code command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "generate_qr_code", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "qr_code_data")
		assert.Contains(t, result, "qr_code_html")
		assert.Contains(t, result["qr_code_html"], "<img")
	})
}

// TestBadgePlugin tests the badge plugin functionality
func TestBadgePlugin(t *testing.T) {
	// Create a test harness
	harness := NewPluginTestHarness(t)

	// Test generate_badge command
	t.Run("GenerateBadge", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "badge_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"label":   "powered by",
			"message": "nessi",
			"color":   "blue",
			"style":   "flat",
		}

		// Execute the generate_badge command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "generate_badge", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "badge_url")
		assert.Contains(t, result, "badge_markdown")
		assert.Contains(t, result, "badge_html")
		assert.Contains(t, result["badge_url"], "shields.io")
		assert.Contains(t, result["badge_markdown"], "![")
		assert.Contains(t, result["badge_html"], "<img")
	})

	// Test quality score badge
	t.Run("QualityScoreBadge", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "badge_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"quality_score": 95,
			"style":         "flat",
		}

		// Execute the generate_badge command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "generate_badge", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "badge_url")
		assert.Contains(t, result, "badge_markdown")
		assert.Contains(t, result, "badge_html")
	})
}

// TestCommunityEngagementPlugin tests the community engagement plugin functionality
func TestCommunityEngagementPlugin(t *testing.T) {
	// Create a test harness
	harness := NewPluginTestHarness(t)

	// Test collect_feedback command
	t.Run("CollectFeedback", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "community_engagement_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"type":  "feature",
			"text":  "It would be great to have a dashboard for quality metrics",
			"name":  "Test User",
			"email": "test@example.com",
		}

		// Execute the collect_feedback command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "collect_feedback", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "issue_url")
		assert.Contains(t, result, "feedback_details")
		assert.Contains(t, result["issue_url"], "github.com")
	})

	// Test suggest_contributions command for beginners
	t.Run("SuggestContributionsForBeginners", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "community_engagement_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"experience":     "beginner",
			"github_profile": "testuser",
		}

		// Execute the suggest_contributions command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "suggest_contributions", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "suggestions")
		suggestions, ok := result["suggestions"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, suggestions)
	})

	// Test suggest_contributions command for advanced users
	t.Run("SuggestContributionsForAdvanced", func(t *testing.T) {
		// Create a mock plugin
		plugin, err := harness.LoadPlugin(t, "community_engagement_plugin")
		require.NoError(t, err)

		// Create test data
		testData := map[string]interface{}{
			"experience":     "advanced",
			"github_profile": "testuser",
		}

		// Execute the suggest_contributions command
		mockPlugin := plugin.(*MockPlugin)
		result, err := mockPlugin.Execute(t, "suggest_contributions", testData)
		require.NoError(t, err)

		// Verify result
		assert.Contains(t, result, "suggestions")
		suggestions, ok := result["suggestions"].([]interface{})
		require.True(t, ok)
		assert.NotEmpty(t, suggestions)
		assert.Len(t, suggestions, 2) // Should have both beginner and advanced suggestions
	})
}

// TestPluginInteroperability tests that all viral growth plugins can work together
func TestPluginInteroperability(t *testing.T) {
	// Create a test harness
	harness := NewPluginTestHarness(t)

	// Create mock plugins
	socialPlugin, err := harness.LoadPlugin(t, "social_sharing_plugin")
	require.NoError(t, err)
	badgePlugin, err := harness.LoadPlugin(t, "badge_plugin")
	require.NoError(t, err)
	communityPlugin, err := harness.LoadPlugin(t, "community_engagement_plugin")
	require.NoError(t, err)

	// Create a basic HTML report
	baseReport := "<html><body><div id=\"report-content\">Test Quality Report</div></body></html>"

	// Enhance with social sharing
	mockSocialPlugin := socialPlugin.(*MockPlugin)
	socialResult, err := mockSocialPlugin.Execute(t, "enhance_report", map[string]interface{}{
		"report_html": baseReport,
		"title":       "Test Quality Report",
		"description": "A test report for integration testing",
		"report_url":  "https://example.com/reports/test-report",
		"hashtags":    "nessi,dataquality,testing",
	})
	require.NoError(t, err)

	// Extract enhanced HTML
	socialEnhancedHTML, ok := socialResult["html_content"].(string)
	require.True(t, ok)
	assert.Contains(t, socialEnhancedHTML, "social-share")

	// Enhance with badge
	mockBadgePlugin := badgePlugin.(*MockPlugin)
	badgeResult, err := mockBadgePlugin.Execute(t, "generate_badge", map[string]interface{}{
		"label":   "powered by",
		"message": "nessi",
		"color":   "blue",
	})
	require.NoError(t, err)

	// Extract badge HTML
	badgeHTML, ok := badgeResult["badge_html"].(string)
	require.True(t, ok)
	assert.Contains(t, badgeHTML, "<img")

	// Verify we can get contribution suggestions
	mockCommunityPlugin := communityPlugin.(*MockPlugin)
	communityResult, err := mockCommunityPlugin.Execute(t, "suggest_contributions", map[string]interface{}{
		"experience": "beginner",
	})
	require.NoError(t, err)

	// Verify suggestions are present
	suggestions, ok := communityResult["suggestions"].([]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, suggestions)
}
