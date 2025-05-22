package plugins

import (
	"fmt"
	"net/url"
	"strings"
)

// MockCommunityEngagementPlugin is a mock implementation of the community engagement plugin for testing
type MockCommunityEngagementPlugin struct {
	initialized bool
}

// Initialize initializes the plugin
func (p *MockCommunityEngagementPlugin) Initialize() error {
	p.initialized = true
	return nil
}

// GetMetadata returns the plugin metadata
func (p *MockCommunityEngagementPlugin) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":        "community_engagement_plugin",
		"version":     "1.0.0",
		"description": "Plugin for community engagement functionality",
		"author":      "Nessi Team",
		"capabilities": []string{
			"community_engagement",
		},
	}
}

// Execute executes a command
func (p *MockCommunityEngagementPlugin) Execute(command []string, data interface{}) (interface{}, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	switch command[0] {
	case "suggest_contributions":
		return p.suggestContributions(data)
	case "enhance_report_with_community":
		return p.enhanceReportWithCommunity(data)
	case "collect_feedback":
		return p.collectFeedback(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", command[0])
	}
}

// suggestContributions suggests ways to contribute to the Nessi project
func (p *MockCommunityEngagementPlugin) suggestContributions(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	experience := "intermediate"
	if exp, ok := dataMap["experience"].(string); ok {
		experience = exp
	}

	var suggestionItems []map[string]interface{}
	switch experience {
	case "beginner":
		suggestionItems = []map[string]interface{}{
			{"text": "Improve documentation", "difficulty": "easy"},
			{"text": "Add examples", "difficulty": "easy"},
			{"text": "Report bugs", "difficulty": "easy"},
			{"text": "Translate documentation", "difficulty": "easy"},
		}
	case "intermediate":
		suggestionItems = []map[string]interface{}{
			{"text": "Fix open issues", "difficulty": "medium"},
			{"text": "Add new features", "difficulty": "medium"},
			{"text": "Improve test coverage", "difficulty": "medium"},
			{"text": "Create tutorials", "difficulty": "medium"},
		}
	case "advanced":
		suggestionItems = []map[string]interface{}{
			{"text": "Implement new data quality checks in core functionality", "difficulty": "hard"},
			{"text": "Optimize performance of core operations", "difficulty": "hard"},
			{"text": "Add support for new data formats", "difficulty": "hard"},
			{"text": "Enhance visualization capabilities", "difficulty": "hard"},
		}
	default:
		suggestionItems = []map[string]interface{}{
			{"text": "Contribute to documentation", "difficulty": "easy"},
			{"text": "Report bugs", "difficulty": "easy"},
			{"text": "Suggest new features", "difficulty": "medium"},
			{"text": "Share Nessi with others", "difficulty": "easy"},
		}
	}

	return map[string]interface{}{
		"suggestions": suggestionItems,
	}, nil
}

// collectFeedback collects feedback from users
func (p *MockCommunityEngagementPlugin) collectFeedback(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	// Extract feedback data
	feedbackType, _ := dataMap["type"].(string)
	feedbackText, _ := dataMap["text"].(string)
	name, _ := dataMap["name"].(string)
	email, _ := dataMap["email"].(string)

	// Generate a mock GitHub issue URL
	issueURL := fmt.Sprintf("https://github.com/nessi-dev/nessi/issues/new?title=%s&body=%s",
		url.QueryEscape(fmt.Sprintf("%s feedback from %s", feedbackType, name)),
		url.QueryEscape(feedbackText))

	return map[string]interface{}{
		"issue_url": issueURL,
		"feedback_details": map[string]interface{}{
			"type":  feedbackType,
			"text":  feedbackText,
			"name":  name,
			"email": email,
		},
	}, nil
}

// enhanceReportWithCommunity enhances a report with community engagement information
func (p *MockCommunityEngagementPlugin) enhanceReportWithCommunity(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	reportHTML, ok := dataMap["report_html"].(string)
	if !ok {
		return nil, fmt.Errorf("report_html is required")
	}

	// Get contribution suggestions
	result, err := p.suggestContributions(map[string]interface{}{
		"experience": "intermediate",
	})
	if err != nil {
		return nil, err
	}

	resultMap := result.(map[string]interface{})
	suggestionItems := resultMap["suggestions"].([]map[string]interface{})

	// Create community section HTML
	var suggestionsHTML strings.Builder
	for _, item := range suggestionItems {
		suggestionsHTML.WriteString(fmt.Sprintf("<li>%s</li>", item["text"]))
	}

	communityHTML := fmt.Sprintf(`
	<div class="community-section">
		<h3>Join the Nessi Community</h3>
		<p>Nessi is an open-source project that thrives on community contributions. Here are some ways you can contribute:</p>
		<ul>
			%s
		</ul>
		<p>Visit our <a href="https://github.com/nessi-dev/nessi" target="_blank">GitHub repository</a> to get started!</p>
		<div class="feedback-section">
			<h4>Share Your Feedback</h4>
			<p>We value your feedback! Let us know what you think about Nessi:</p>
			<a href="https://github.com/nessi-dev/nessi/issues/new?template=feedback.md" class="feedback-button">Submit Feedback</a>
		</div>
	</div>
	`, suggestionsHTML.String())

	// Insert community section before the end of the body
	enhancedHTML := strings.Replace(reportHTML, "</body>", communityHTML+"</body>", 1)

	return map[string]interface{}{
		"enhanced_html": enhancedHTML,
	}, nil
}
