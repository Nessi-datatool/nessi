package plugins

import (
	"fmt"
	"net/url"
	"strings"
)

// MockSocialSharingPlugin is a mock implementation of the social sharing plugin for testing
type MockSocialSharingPlugin struct {
	initialized bool
}

// Initialize initializes the plugin
func (p *MockSocialSharingPlugin) Initialize() error {
	p.initialized = true
	return nil
}

// GetMetadata returns the plugin metadata
func (p *MockSocialSharingPlugin) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":        "social_sharing_plugin",
		"version":     "1.0.0",
		"description": "Plugin for social sharing functionality",
		"author":      "Nessi Team",
		"capabilities": []string{
			"social_sharing",
		},
	}
}

// Execute executes a command
func (p *MockSocialSharingPlugin) Execute(command []string, data interface{}) (interface{}, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	switch command[0] {
	case "generate_share_links":
		return p.generateShareLinks(data)
	case "enhance_report":
		return p.enhanceReport(data)
	case "generate_qr_code":
		return p.generateQRCode(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", command[0])
	}
}

// generateShareLinks generates share links for social media
func (p *MockSocialSharingPlugin) generateShareLinks(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	title, ok := dataMap["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title is required")
	}

	reportURL, ok := dataMap["report_url"].(string)
	if !ok {
		return nil, fmt.Errorf("report_url is required")
	}

	description := "Check out this Nessi report!"
	if desc, ok := dataMap["description"].(string); ok {
		description = desc
	}

	hashtags := "nessi,dataquality"
	if tags, ok := dataMap["hashtags"].(string); ok {
		hashtags = tags
	}

	// Generate Twitter link with the title directly visible in the URL (not encoded)
	twitterText := url.QueryEscape(title + ": " + description)
	twitterURL := url.QueryEscape(reportURL)
	twitterHashtags := url.QueryEscape(hashtags)
	twitterLink := fmt.Sprintf("https://twitter.com/intent/tweet?text=%s&url=%s&hashtags=%s",
		twitterText, twitterURL, twitterHashtags)

	// Generate LinkedIn link
	linkedinTitle := url.QueryEscape(title)
	linkedinSummary := url.QueryEscape(description)
	linkedinURL := url.QueryEscape(reportURL)
	linkedinLink := fmt.Sprintf("https://www.linkedin.com/shareArticle?mini=true&url=%s&title=%s&summary=%s",
		linkedinURL, linkedinTitle, linkedinSummary)

	// Generate Facebook link
	facebookURL := url.QueryEscape(reportURL)
	facebookLink := fmt.Sprintf("https://www.facebook.com/sharer/sharer.php?u=%s", facebookURL)

	// Generate email link
	emailSubject := url.QueryEscape("Nessi Report: " + title)
	emailBody := url.QueryEscape(description + "\n\n" + reportURL)
	emailLink := fmt.Sprintf("mailto:?subject=%s&body=%s", emailSubject, emailBody)

	return map[string]interface{}{
		"twitter_link":  twitterLink,
		"linkedin_link": linkedinLink,
		"facebook_link": facebookLink,
		"email_link":    emailLink,
	}, nil
}

// enhanceReport enhances a report with social sharing buttons
func (p *MockSocialSharingPlugin) enhanceReport(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	reportHTML, ok := dataMap["report_html"].(string)
	if !ok {
		return nil, fmt.Errorf("report_html is required")
	}

	// Generate share links
	shareLinks, err := p.generateShareLinks(data)
	if err != nil {
		return nil, err
	}

	shareLinksMap := shareLinks.(map[string]interface{})
	twitterLink := shareLinksMap["twitter_link"].(string)
	linkedinLink := shareLinksMap["linkedin_link"].(string)
	facebookLink := shareLinksMap["facebook_link"].(string)

	// Create social sharing buttons HTML
	socialHTML := fmt.Sprintf(`
	<div class="social-share">
		<h3>Share this report</h3>
		<a href="%s" target="_blank" class="social-button twitter">Share on Twitter</a>
		<a href="%s" target="_blank" class="social-button linkedin">Share on LinkedIn</a>
		<a href="%s" target="_blank" class="social-button facebook">Share on Facebook</a>
	</div>
	`, twitterLink, linkedinLink, facebookLink)

	// Insert social sharing buttons before the end of the body
	enhancedHTML := strings.Replace(reportHTML, "</body>", socialHTML+"</body>", 1)

	return map[string]interface{}{
		"html_content": enhancedHTML,
		"share_data": map[string]interface{}{
			"title":       dataMap["title"],
			"description": dataMap["description"],
			"report_url":  dataMap["report_url"],
			"hashtags":    dataMap["hashtags"],
		},
	}, nil
}

// generateQRCode generates a QR code for a URL
func (p *MockSocialSharingPlugin) generateQRCode(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}

	reportURL, ok := dataMap["report_url"].(string)
	if !ok {
		return nil, fmt.Errorf("report_url is required")
	}

	// Mock QR code data (base64 encoded)
	qrCodeData := "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAAA1BMVEX///+nxBvIAAAAH0lEQVRoge3BAQ0AAADCIPunNsc3YAAAAAAAAAAAADwDTbgAAUiIBVoAAAAASUVORK5CYII="

	// Generate HTML for the QR code
	qrCodeHTML := fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="QR Code for %s">`,
		qrCodeData, reportURL)

	return map[string]interface{}{
		"qr_code_data": qrCodeData,
		"qr_code_html": qrCodeHTML,
		"report_url":   reportURL,
	}, nil
}
