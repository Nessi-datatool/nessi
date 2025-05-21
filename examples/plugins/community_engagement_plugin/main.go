package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nessi-dev/nessi/pkg/plugins"
)

// CommunityEngagementPlugin implements the PluginInterface for community engagement features
type CommunityEngagementPlugin struct {
	Name        string
	Version     string
	Description string
	Author      string
}

// CommunityData contains data for community engagement features
type CommunityData struct {
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	ProjectName   string `json:"project_name"`
	FeedbackType  string `json:"feedback_type"`
	FeedbackText  string `json:"feedback_text"`
	Contribution  bool   `json:"contribution"`
	GithubProfile string `json:"github_profile"`
}

// Initialize sets up the plugin
func (p *CommunityEngagementPlugin) Initialize() error {
	p.Name = "community_engagement_plugin"
	p.Version = "1.0.0"
	p.Description = "Enhances Nessi with community engagement features"
	p.Author = "Nessi Community"
	return nil
}

// GetMetadata returns plugin metadata
func (p *CommunityEngagementPlugin) GetMetadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Capabilities: []string{
			"community_engagement",
			"feedback_collection",
			"contribution_guidance",
		},
	}
}

// Execute runs the plugin functionality
func (p *CommunityEngagementPlugin) Execute(args []string, data interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified for community engagement plugin")
	}

	switch args[0] {
	case "collect_feedback":
		return p.collectFeedback(data)
	case "suggest_contributions":
		return p.suggestContributions(data)
	case "enhance_report":
		return p.enhanceReportWithCommunity(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}

// collectFeedback processes user feedback
func (p *CommunityEngagementPlugin) collectFeedback(data interface{}) (interface{}, error) {
	// Parse input data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %v", err)
	}

	var communityData CommunityData
	if err := json.Unmarshal(dataBytes, &communityData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
	}

	// Validate required fields
	if communityData.FeedbackText == "" {
		return nil, fmt.Errorf("feedback text is required")
	}

	// Generate feedback ID and timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// In a real implementation, this would save the feedback to a database or file
	return map[string]interface{}{
		"status":      "success",
		"message":     "Thank you for your feedback!",
		"feedback_id": fmt.Sprintf("fb-%d", time.Now().Unix()),
		"timestamp":   timestamp,
		"github_issue_url": fmt.Sprintf("https://github.com/nessi-dev/nessi/issues/new?title=%s&body=%s",
			"Feedback: "+communityData.FeedbackType,
			"Feedback from Nessi user:\n\n"+communityData.FeedbackText),
	}, nil
}

// suggestContributions suggests ways to contribute to Nessi
func (p *CommunityEngagementPlugin) suggestContributions(data interface{}) (interface{}, error) {
	// Parse input data if provided
	var communityData CommunityData
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input data: %v", err)
		}

		if err := json.Unmarshal(dataBytes, &communityData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
		}
	}

	// Generate contribution suggestions
	suggestions := []map[string]interface{}{
		{
			"title":       "Add a New Quality Rule",
			"description": "Contribute a new data quality rule to help validate Delta Lake tables",
			"difficulty":  "Easy",
			"link":        "https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#adding-quality-rules",
		},
		{
			"title":       "Improve Documentation",
			"description": "Help improve Nessi's documentation with examples, tutorials, or clarifications",
			"difficulty":  "Easy",
			"link":        "https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#improving-documentation",
		},
		{
			"title":       "Add Format Support",
			"description": "Extend Nessi to support additional data formats beyond Delta Lake",
			"difficulty":  "Medium",
			"link":        "https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#adding-format-support",
		},
		{
			"title":       "Create a New Plugin",
			"description": "Develop a new plugin to extend Nessi's functionality",
			"difficulty":  "Medium",
			"link":        "https://github.com/nessi-dev/nessi/blob/main/docs/PLUGIN_API.md",
		},
		{
			"title":       "Performance Optimization",
			"description": "Help optimize Nessi's performance for large Delta Lake tables",
			"difficulty":  "Hard",
			"link":        "https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md#performance-optimization",
		},
	}

	// Add first-time contributor suggestions
	firstTimeContributions := []map[string]interface{}{
		{
			"title":       "Fix a Documentation Typo",
			"description": "Start with a simple documentation fix to learn the contribution process",
			"difficulty":  "Very Easy",
			"link":        "https://github.com/nessi-dev/nessi/blob/main/docs/",
		},
		{
			"title":       "Add Test Case",
			"description": "Add a test case for an existing feature",
			"difficulty":  "Easy",
			"link":        "https://github.com/nessi-dev/nessi/tree/main/pkg/test",
		},
	}

	return map[string]interface{}{
		"suggestions":            suggestions,
		"first_time_suggestions": firstTimeContributions,
		"contribution_guide":     "https://github.com/nessi-dev/nessi/blob/main/docs/CONTRIBUTING.md",
		"community_slack":        "https://join.slack.com/t/nessi-community/shared_invite/...",
		"github_repo":            "https://github.com/nessi-dev/nessi",
	}, nil
}

// enhanceReportWithCommunity adds community engagement elements to reports
func (p *CommunityEngagementPlugin) enhanceReportWithCommunity(data interface{}) (interface{}, error) {
	// Generate HTML for community engagement section
	htmlContent := `
<div class="community-engagement">
  <h3>Join the Nessi Community</h3>
  <p>Nessi is an open-source project that thrives on community contributions. Here's how you can get involved:</p>
  
  <div class="community-cards">
    <div class="community-card">
      <div class="card-icon"><i class="fas fa-code-branch"></i></div>
      <h4>Contribute</h4>
      <p>Help improve Nessi by contributing code, documentation, or ideas.</p>
      <a href="https://github.com/nessi-dev/nessi" target="_blank" class="card-link">GitHub Repository</a>
    </div>
    
    <div class="community-card">
      <div class="card-icon"><i class="fas fa-comments"></i></div>
      <h4>Discuss</h4>
      <p>Join our community discussions on Slack or GitHub Discussions.</p>
      <a href="https://join.slack.com/t/nessi-community/shared_invite/..." target="_blank" class="card-link">Join Slack</a>
    </div>
    
    <div class="community-card">
      <div class="card-icon"><i class="fas fa-star"></i></div>
      <h4>Star Us</h4>
      <p>Show your support by starring Nessi on GitHub.</p>
      <a href="https://github.com/nessi-dev/nessi" target="_blank" class="card-link">Star on GitHub</a>
    </div>
  </div>
  
  <div class="feedback-form">
    <h4>Quick Feedback</h4>
    <p>Help us improve Nessi by sharing your thoughts:</p>
    <div class="feedback-buttons">
      <button class="feedback-button" onclick="window.open('https://github.com/nessi-dev/nessi/issues/new?title=Feature%20Request&labels=enhancement', '_blank')">
        <i class="fas fa-lightbulb"></i> Suggest Feature
      </button>
      <button class="feedback-button" onclick="window.open('https://github.com/nessi-dev/nessi/issues/new?title=Bug%20Report&labels=bug', '_blank')">
        <i class="fas fa-bug"></i> Report Bug
      </button>
      <button class="feedback-button" onclick="window.open('https://forms.gle/nessi-feedback', '_blank')">
        <i class="fas fa-comment"></i> General Feedback
      </button>
    </div>
  </div>
</div>

<style>
.community-engagement {
  background-color: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  margin: 30px 0;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}
.community-engagement h3 {
  margin-top: 0;
  color: #2c3e50;
  font-size: 20px;
  margin-bottom: 15px;
}
.community-cards {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin: 20px 0;
}
.community-card {
  flex: 1;
  min-width: 200px;
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.05);
  text-align: center;
  transition: transform 0.3s ease;
}
.community-card:hover {
  transform: translateY(-5px);
}
.card-icon {
  font-size: 24px;
  color: #3498db;
  margin-bottom: 10px;
}
.community-card h4 {
  margin: 10px 0;
  color: #2c3e50;
}
.community-card p {
  color: #7f8c8d;
  font-size: 14px;
  margin-bottom: 15px;
}
.card-link {
  display: inline-block;
  color: #3498db;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
}
.card-link:hover {
  text-decoration: underline;
}
.feedback-form {
  background: white;
  border-radius: 8px;
  padding: 15px;
  margin-top: 20px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.05);
}
.feedback-form h4 {
  margin-top: 0;
  color: #2c3e50;
}
.feedback-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 15px;
}
.feedback-button {
  padding: 8px 15px;
  border: none;
  border-radius: 5px;
  background: #3498db;
  color: white;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  transition: background 0.3s ease;
}
.feedback-button:hover {
  background: #2980b9;
}
@media (max-width: 768px) {
  .community-cards {
    flex-direction: column;
  }
  .feedback-buttons {
    flex-direction: column;
  }
  .feedback-button {
    width: 100%;
  }
}
</style>
`

	return map[string]interface{}{
		"html_content": htmlContent,
		"template_vars": map[string]interface{}{
			"community_engagement_enabled": true,
			"github_repo":                  "https://github.com/nessi-dev/nessi",
			"slack_invite":                 "https://join.slack.com/t/nessi-community/shared_invite/...",
			"feedback_form":                "https://forms.gle/nessi-feedback",
		},
	}, nil
}

// Plugin is the exported symbol that Nessi will look for
var Plugin CommunityEngagementPlugin
