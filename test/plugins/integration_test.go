package plugins

import (
	"testing"

	"github.com/nessi-dev/nessi/examples/plugins/badge_plugin"
	"github.com/nessi-dev/nessi/examples/plugins/community_engagement_plugin"
	"github.com/nessi-dev/nessi/examples/plugins/social_sharing_plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPluginInteroperability tests that all viral growth plugins can work together
func TestPluginInteroperability(t *testing.T) {
	// Initialize all plugins
	socialPlugin := &social_sharing_plugin.SocialSharingPlugin{}
	badgePlugin := &badge_plugin.BadgePlugin{}
	communityPlugin := &community_engagement_plugin.CommunityEngagementPlugin{}

	require.NoError(t, socialPlugin.Initialize())
	require.NoError(t, badgePlugin.Initialize())
	require.NoError(t, communityPlugin.Initialize())

	// Create a basic HTML report
	baseReport := "<html><body><div id=\"report-content\">Test Quality Report</div></body></html>"

	// Enhance with social sharing
	socialResult, err := socialPlugin.Execute([]string{"enhance_report"}, map[string]interface{}{
		"report_html": baseReport,
		"title":       "Test Quality Report",
		"description": "A test report for integration testing",
		"report_url":  "https://example.com/reports/test-report",
		"hashtags":    "nessi,dataquality,testing",
	})
	require.NoError(t, err)

	// Extract enhanced HTML
	socialResultMap, ok := socialResult.(map[string]interface{})
	require.True(t, ok)
	socialEnhancedHTML, ok := socialResultMap["html_content"].(string)
	require.True(t, ok)
	assert.Contains(t, socialEnhancedHTML, "social-share")

	// Enhance with badge
	badgeResult, err := badgePlugin.Execute([]string{"enhance_report_with_badge"}, map[string]interface{}{
		"report_html": socialEnhancedHTML,
		"label":       "powered by",
		"message":     "nessi",
		"color":       "blue",
	})
	require.NoError(t, err)

	// Extract enhanced HTML
	badgeResultMap, ok := badgeResult.(map[string]interface{})
	require.True(t, ok)
	badgeEnhancedHTML, ok := badgeResultMap["enhanced_html"].(string)
	require.True(t, ok)
	assert.Contains(t, badgeEnhancedHTML, "social-share")
	assert.Contains(t, badgeEnhancedHTML, "shields.io")

	// Enhance with community
	communityResult, err := communityPlugin.Execute([]string{"enhance_report_with_community"}, map[string]interface{}{
		"report_html": badgeEnhancedHTML,
	})
	require.NoError(t, err)

	// Extract final enhanced HTML
	communityResultMap, ok := communityResult.(map[string]interface{})
	require.True(t, ok)
	finalHTML, ok := communityResultMap["enhanced_html"].(string)
	require.True(t, ok)

	// Verify all enhancements are present
	assert.Contains(t, finalHTML, "social-share", "Social sharing elements should be present")
	assert.Contains(t, finalHTML, "shields.io", "Badge should be present")
	assert.Contains(t, finalHTML, "community", "Community section should be present")
	assert.Contains(t, finalHTML, "Test Quality Report", "Original content should be preserved")
}

// TestPluginLoadOrder tests that plugins can be loaded in any order
func TestPluginLoadOrder(t *testing.T) {
	// Test different loading orders
	loadOrders := [][]string{
		{"social", "badge", "community"},
		{"social", "community", "badge"},
		{"badge", "social", "community"},
		{"badge", "community", "social"},
		{"community", "social", "badge"},
		{"community", "badge", "social"},
	}

	for _, order := range loadOrders {
		t.Run("Load order: "+order[0]+","+order[1]+","+order[2], func(t *testing.T) {
			// Create plugin instances
			var socialPlugin *social_sharing_plugin.SocialSharingPlugin
			var badgePlugin *badge_plugin.BadgePlugin
			var communityPlugin *community_engagement_plugin.CommunityEngagementPlugin

			// Initialize plugins in the specified order
			for _, pluginType := range order {
				switch pluginType {
				case "social":
					socialPlugin = &social_sharing_plugin.SocialSharingPlugin{}
					require.NoError(t, socialPlugin.Initialize())
				case "badge":
					badgePlugin = &badge_plugin.BadgePlugin{}
					require.NoError(t, badgePlugin.Initialize())
				case "community":
					communityPlugin = &community_engagement_plugin.CommunityEngagementPlugin{}
					require.NoError(t, communityPlugin.Initialize())
				}
			}

			// Verify all plugins are initialized correctly
			assert.Equal(t, "social_sharing_plugin", socialPlugin.GetMetadata().Name)
			assert.Equal(t, "badge_plugin", badgePlugin.GetMetadata().Name)
			assert.Equal(t, "community_engagement_plugin", communityPlugin.GetMetadata().Name)

			// Verify each plugin can execute its primary function
			_, err := socialPlugin.Execute([]string{"generate_share_links"}, map[string]interface{}{
				"title":      "Test Report",
				"report_url": "https://example.com/reports/test-report",
			})
			assert.NoError(t, err)

			_, err = badgePlugin.Execute([]string{"generate_badge"}, map[string]interface{}{
				"label":   "powered by",
				"message": "nessi",
			})
			assert.NoError(t, err)

			_, err = communityPlugin.Execute([]string{"suggest_contributions"}, map[string]interface{}{
				"experience": "beginner",
			})
			assert.NoError(t, err)
		})
	}
}

// TestPluginResourceIsolation tests that plugins don't interfere with each other's resources
func TestPluginResourceIsolation(t *testing.T) {
	// Initialize all plugins
	socialPlugin1 := &social_sharing_plugin.SocialSharingPlugin{}
	socialPlugin2 := &social_sharing_plugin.SocialSharingPlugin{}
	badgePlugin := &badge_plugin.BadgePlugin{}
	communityPlugin := &community_engagement_plugin.CommunityEngagementPlugin{}

	require.NoError(t, socialPlugin1.Initialize())
	require.NoError(t, socialPlugin2.Initialize())
	require.NoError(t, badgePlugin.Initialize())
	require.NoError(t, communityPlugin.Initialize())

	// Verify that multiple instances of the same plugin type don't interfere
	result1, err := socialPlugin1.Execute([]string{"generate_share_links"}, map[string]interface{}{
		"title":      "Report 1",
		"report_url": "https://example.com/reports/report-1",
	})
	require.NoError(t, err)

	result2, err := socialPlugin2.Execute([]string{"generate_share_links"}, map[string]interface{}{
		"title":      "Report 2",
		"report_url": "https://example.com/reports/report-2",
	})
	require.NoError(t, err)

	// Verify results are different
	resultMap1, ok := result1.(map[string]interface{})
	require.True(t, ok)
	resultMap2, ok := result2.(map[string]interface{})
	require.True(t, ok)

	twitterLink1, ok := resultMap1["twitter_link"].(string)
	require.True(t, ok)
	twitterLink2, ok := resultMap2["twitter_link"].(string)
	require.True(t, ok)

	assert.NotEqual(t, twitterLink1, twitterLink2, "Different plugin instances should produce different results")
	assert.Contains(t, twitterLink1, "Report 1", "First plugin should use first report title")
	assert.Contains(t, twitterLink2, "Report 2", "Second plugin should use second report title")

	// Verify that different plugin types don't interfere
	_, err = badgePlugin.Execute([]string{"generate_badge"}, map[string]interface{}{
		"label":   "powered by",
		"message": "nessi",
	})
	require.NoError(t, err)

	// Verify social plugin still works after badge plugin execution
	result3, err := socialPlugin1.Execute([]string{"generate_share_links"}, map[string]interface{}{
		"title":      "Report 1",
		"report_url": "https://example.com/reports/report-1",
	})
	require.NoError(t, err)

	resultMap3, ok := result3.(map[string]interface{})
	require.True(t, ok)
	twitterLink3, ok := resultMap3["twitter_link"].(string)
	require.True(t, ok)

	assert.Contains(t, twitterLink3, "Report 1", "Social plugin should still work after badge plugin execution")
}
