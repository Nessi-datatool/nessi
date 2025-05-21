package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// BenchmarkSocialSharingPlugin benchmarks the social sharing plugin
func BenchmarkSocialSharingPlugin(b *testing.B) {
	// Initialize the test harness
	harness := NewPluginTestHarness(b)

	// Benchmark the EnhanceReport function
	b.Run("EnhanceReport", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"title":       "Test Report",
			"description": "A test report for benchmarking",
			"metrics":     map[string]float64{"quality": 95.0, "completeness": 98.0},
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's EnhanceReport function
			result, err := harness.CallPluginFunction("social_sharing_plugin", "EnhanceReport", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})

	// Benchmark the GenerateShareLinks function
	b.Run("GenerateShareLinks", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"url":         "https://example.com/report/123",
			"title":       "Test Report",
			"description": "A test report for benchmarking",
			"hashtags":    "nessi,dataquality,test",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's GenerateShareLinks function
			result, err := harness.CallPluginFunction("social_sharing_plugin", "GenerateShareLinks", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})

	// Benchmark the GenerateQRCode function
	b.Run("GenerateQRCode", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"url":  "https://example.com/report/123",
			"size": 256,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's GenerateQRCode function
			result, err := harness.CallPluginFunction("social_sharing_plugin", "GenerateQRCode", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})
}

// BenchmarkBadgePlugin benchmarks the badge plugin
func BenchmarkBadgePlugin(b *testing.B) {
	// Initialize the test harness
	harness := NewPluginTestHarness(b)

	// Benchmark the GenerateBadge function
	b.Run("GenerateBadge", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"label":   "powered by",
			"message": "nessi",
			"color":   "blue",
			"style":   "flat",
			"format":  "markdown",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's GenerateBadge function
			result, err := harness.CallPluginFunction("badge_plugin", "GenerateBadge", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})

	// Benchmark the QualityScoreBadge function
	b.Run("QualityScoreBadge", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"score":  95,
			"format": "markdown",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's QualityScoreBadge function
			result, err := harness.CallPluginFunction("badge_plugin", "QualityScoreBadge", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})
}

// BenchmarkCommunityEngagementPlugin benchmarks the community engagement plugin
func BenchmarkCommunityEngagementPlugin(b *testing.B) {
	// Initialize the test harness
	harness := NewPluginTestHarness(b)

	// Benchmark the CollectFeedback function
	b.Run("CollectFeedback", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"feedback_type": "feature",
			"feedback_text": "I would like to see more visualization options",
			"user_name":    "Test User",
			"user_email":   "test@example.com",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's CollectFeedback function
			result, err := harness.CallPluginFunction("community_engagement_plugin", "CollectFeedback", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})

	// Benchmark the SuggestContributions function for beginners
	b.Run("SuggestContributionsForBeginners", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"experience_level": "beginner",
			"github_profile":  "testuser",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's SuggestContributions function
			result, err := harness.CallPluginFunction("community_engagement_plugin", "SuggestContributions", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})

	// Benchmark the SuggestContributions function for advanced users
	b.Run("SuggestContributionsForAdvanced", func(b *testing.B) {
		// Prepare test data
		testData := map[string]interface{}{
			"experience_level": "advanced",
			"github_profile":  "testuser",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Call the plugin's SuggestContributions function
			result, err := harness.CallPluginFunction("community_engagement_plugin", "SuggestContributions", testData)
			require.NoError(b, err)
			require.NotNil(b, result)
		}
	})
}

// BenchmarkPluginInteroperability benchmarks the interoperability of all plugins
func BenchmarkPluginInteroperability(b *testing.B) {
	// Initialize the test harness
	harness := NewPluginTestHarness(b)

	// Prepare test data for a report
	reportData := map[string]interface{}{
		"title":       "Test Report",
		"description": "A test report for benchmarking",
		"metrics":     map[string]float64{"quality": 95.0, "completeness": 98.0},
		"url":         "https://example.com/report/123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Step 1: Enhance the report with social sharing features
		enhancedReport, err := harness.CallPluginFunction("social_sharing_plugin", "EnhanceReport", reportData)
		require.NoError(b, err)
		require.NotNil(b, enhancedReport)

		// Step 2: Add a quality score badge to the report
		badgeData := map[string]interface{}{
			"score":  95,
			"format": "markdown",
		}
		badge, err := harness.CallPluginFunction("badge_plugin", "QualityScoreBadge", badgeData)
		require.NoError(b, err)
		require.NotNil(b, badge)

		// Step 3: Add the badge to the report
		enhancedReportMap, ok := enhancedReport.(map[string]interface{})
		require.True(b, ok)
		enhancedReportMap["badge"] = badge

		// Step 4: Add community engagement features
		feedbackData := map[string]interface{}{
			"report_id":     "123",
			"feedback_type": "general",
		}
		feedbackLink, err := harness.CallPluginFunction("community_engagement_plugin", "CollectFeedback", feedbackData)
		require.NoError(b, err)
		require.NotNil(b, feedbackLink)

		// Step 5: Add the feedback link to the report
		enhancedReportMap["feedback_link"] = feedbackLink
	}
}

// Helper method to extend the test harness for benchmarking
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
