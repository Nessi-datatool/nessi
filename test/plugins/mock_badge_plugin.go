package plugins

import (
	"fmt"
	"net/url"
	"strings"
)

// MockBadgePlugin is a mock implementation of the badge plugin for testing
type MockBadgePlugin struct {
	initialized bool
}

// Initialize initializes the plugin
func (p *MockBadgePlugin) Initialize() error {
	p.initialized = true
	return nil
}

// GetMetadata returns the plugin metadata
func (p *MockBadgePlugin) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":        "badge_plugin",
		"version":     "1.0.0",
		"description": "Plugin for generating badges",
		"author":      "Nessi Team",
		"capabilities": []string{
			"badge_generation",
		},
	}
}

// Execute executes a command
func (p *MockBadgePlugin) Execute(command []string, data interface{}) (interface{}, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	switch command[0] {
	case "generate_badge":
		return p.generateBadge(data)
	case "get_badge_markdown":
		return p.getBadgeMarkdown(data)
	case "get_badge_html":
		return p.getBadgeHTML(data)
	case "enhance_report_with_badge":
		return p.enhanceReportWithBadge(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", command[0])
	}
}

// generateBadge generates a badge
func (p *MockBadgePlugin) generateBadge(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	label := "powered by"
	message := "nessi"
	color := "blue"
	style := "flat"
	qualityScore := 0

	if val, ok := dataMap["label"].(string); ok {
		label = val
	}
	if val, ok := dataMap["message"].(string); ok {
		message = val
	}
	if val, ok := dataMap["color"].(string); ok {
		color = val
	}
	if val, ok := dataMap["style"].(string); ok {
		style = val
	}
	if val, ok := dataMap["quality_score"].(float64); ok {
		qualityScore = int(val)
	}

	// Generate badge URL
	badgeURL := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s?style=%s",
		url.PathEscape(label),
		url.PathEscape(message),
		color,
		style)

	if qualityScore > 0 {
		badgeURL = fmt.Sprintf("https://img.shields.io/badge/quality-%d%%25-%s?style=%s",
			qualityScore,
			p.getColorForScore(qualityScore),
			style)
	}

	// Generate markdown
	badgeMarkdown := fmt.Sprintf("![%s](%s)", label, badgeURL)

	// Generate HTML
	badgeHTML := fmt.Sprintf("<img src=\"%s\" alt=\"%s %s\">", badgeURL, label, message)

	return map[string]interface{}{
		"badge_url":      badgeURL,
		"badge_markdown": badgeMarkdown,
		"badge_html":     badgeHTML,
	}, nil
}

// getBadgeMarkdown gets the badge markdown
func (p *MockBadgePlugin) getBadgeMarkdown(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	label := "powered by"
	message := "nessi"
	color := "blue"
	style := "flat"

	if val, ok := dataMap["label"].(string); ok {
		label = val
	}
	if val, ok := dataMap["message"].(string); ok {
		message = val
	}
	if val, ok := dataMap["color"].(string); ok {
		color = val
	}
	if val, ok := dataMap["style"].(string); ok {
		style = val
	}

	// Generate badge URL
	badgeURL := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s?style=%s",
		url.PathEscape(label),
		url.PathEscape(message),
		color,
		style)

	// Generate markdown
	badgeMarkdown := fmt.Sprintf("![%s](%s)", label, badgeURL)

	return map[string]interface{}{
		"markdown": badgeMarkdown,
	}, nil
}

// getBadgeHTML gets the badge HTML
func (p *MockBadgePlugin) getBadgeHTML(data interface{}) (interface{}, error) {
	result, err := p.generateBadge(data)
	if err != nil {
		return nil, err
	}

	resultMap := result.(map[string]interface{})
	return map[string]interface{}{
		"html": resultMap["badge_html"],
	}, nil
}

// enhanceReportWithBadge enhances a report with a badge
func (p *MockBadgePlugin) enhanceReportWithBadge(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	reportHTML, ok := dataMap["report_html"].(string)
	if !ok {
		return nil, fmt.Errorf("report_html is required")
	}

	result, err := p.generateBadge(data)
	if err != nil {
		return nil, err
	}

	resultMap := result.(map[string]interface{})
	badgeHTML := resultMap["badge_html"].(string)

	// Insert badge at the end of the body
	enhancedHTML := strings.Replace(reportHTML, "</body>", badgeHTML+"</body>", 1)

	return map[string]interface{}{
		"enhanced_html": enhancedHTML,
	}, nil
}

// getColorForScore returns a color based on the quality score
func (p *MockBadgePlugin) getColorForScore(score int) string {
	if score >= 90 {
		return "brightgreen"
	} else if score >= 80 {
		return "green"
	} else if score >= 70 {
		return "yellowgreen"
	} else if score >= 60 {
		return "yellow"
	} else if score >= 50 {
		return "orange"
	} else {
		return "red"
	}
}
