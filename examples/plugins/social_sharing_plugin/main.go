package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/nessi-dev/nessi/pkg/plugins"
)

// SocialSharingPlugin implements the PluginInterface for enhancing report sharing capabilities
type SocialSharingPlugin struct {
	Name        string
	Version     string
	Description string
	Author      string
}

// SocialShareData contains data for social sharing features
type SocialShareData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ReportURL   string `json:"report_url"`
	QRCodeURL   string `json:"qr_code_url"`
	Hashtags    string `json:"hashtags"`
	AnalyticsID string `json:"analytics_id"`
}

// Initialize sets up the plugin
func (p *SocialSharingPlugin) Initialize() error {
	p.Name = "social_sharing_plugin"
	p.Version = "1.0.0"
	p.Description = "Enhances Nessi reports with advanced social sharing capabilities"
	p.Author = "Nessi Community"
	return nil
}

// GetMetadata returns plugin metadata
func (p *SocialSharingPlugin) GetMetadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Capabilities: []string{
			"report_enhancement",
			"social_sharing",
		},
	}
}

// Execute runs the plugin functionality
func (p *SocialSharingPlugin) Execute(args []string, data interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified for social sharing plugin")
	}

	switch args[0] {
	case "enhance_report":
		return p.enhanceReport(data)
	case "generate_share_links":
		return p.generateShareLinks(data)
	case "generate_qr_code":
		return p.generateQRCode(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}

// enhanceReport adds social sharing elements to the report template
func (p *SocialSharingPlugin) enhanceReport(data interface{}) (interface{}, error) {
	// Parse input data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %v", err)
	}

	var shareData SocialShareData
	if err := json.Unmarshal(dataBytes, &shareData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
	}

	// Generate QR code URL if not provided
	if shareData.QRCodeURL == "" && shareData.ReportURL != "" {
		shareData.QRCodeURL = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=%s", url.QueryEscape(shareData.ReportURL))
	}

	// Set default hashtags if not provided
	if shareData.Hashtags == "" {
		shareData.Hashtags = "nessi,dataquality,datalake,opensource"
	}

	// Generate HTML for social sharing enhancements
	htmlContent := generateSocialSharingHTML(shareData)

	return map[string]interface{}{
		"html_content":  htmlContent,
		"share_data":    shareData,
		"template_vars": generateTemplateVars(shareData),
	}, nil
}

// generateShareLinks creates optimized sharing links for various platforms
func (p *SocialSharingPlugin) generateShareLinks(data interface{}) (interface{}, error) {
	// Parse input data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %v", err)
	}

	var shareData SocialShareData
	if err := json.Unmarshal(dataBytes, &shareData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
	}

	// Create sharing links
	twitterText := fmt.Sprintf("Check out this data quality report: %s #%s",
		shareData.Title, strings.ReplaceAll(shareData.Hashtags, ",", " #"))

	linkedinText := fmt.Sprintf("I wanted to share this data quality report: %s", shareData.Title)

	emailSubject := fmt.Sprintf("Data Quality Report: %s", shareData.Title)
	emailBody := fmt.Sprintf("I thought you might be interested in this data quality report:\n\n%s\n\nView the report here: %s",
		shareData.Description, shareData.ReportURL)

	// Add analytics tracking if provided
	reportURL := shareData.ReportURL
	if shareData.AnalyticsID != "" {
		// Add UTM parameters for tracking
		reportURL = addUTMParams(reportURL, "social_share", "plugin")
	}

	return map[string]interface{}{
		"twitter_link":  fmt.Sprintf("https://twitter.com/intent/tweet?text=%s&url=%s", url.QueryEscape(twitterText), url.QueryEscape(reportURL)),
		"linkedin_link": fmt.Sprintf("https://www.linkedin.com/sharing/share-offsite/?url=%s&summary=%s", url.QueryEscape(reportURL), url.QueryEscape(linkedinText)),
		"email_link":    fmt.Sprintf("mailto:?subject=%s&body=%s", url.QueryEscape(emailSubject), url.QueryEscape(emailBody)),
		"whatsapp_link": fmt.Sprintf("https://wa.me/?text=%s", url.QueryEscape(fmt.Sprintf("%s: %s", shareData.Title, reportURL))),
		"slack_text":    fmt.Sprintf("Check out this data quality report: %s\n%s", shareData.Title, reportURL),
	}, nil
}

// generateQRCode creates a QR code for the report URL
func (p *SocialSharingPlugin) generateQRCode(data interface{}) (interface{}, error) {
	// Parse input data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %v", err)
	}

	var shareData SocialShareData
	if err := json.Unmarshal(dataBytes, &shareData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
	}

	if shareData.ReportURL == "" {
		return nil, fmt.Errorf("report URL is required to generate QR code")
	}

	// Generate QR code URL using an external service
	qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=%s",
		url.QueryEscape(shareData.ReportURL))

	return map[string]interface{}{
		"qr_code_url": qrCodeURL,
		"report_url":  shareData.ReportURL,
	}, nil
}

// Helper function to generate HTML for social sharing features
func generateSocialSharingHTML(data SocialShareData) string {
	html := `
<div class="social-share-enhanced">
  <h3>Share This Report</h3>
  <div class="share-container">
    <div class="qr-code">
      <img src="%s" alt="QR Code for Report" />
      <p>Scan to view report</p>
    </div>
    <div class="share-buttons">
      <a href="https://twitter.com/intent/tweet?text=%s&url=%s&hashtags=%s" target="_blank" class="share-button share-twitter">
        <i class="fab fa-twitter"></i> Twitter
      </a>
      <a href="https://www.linkedin.com/sharing/share-offsite/?url=%s" target="_blank" class="share-button share-linkedin">
        <i class="fab fa-linkedin"></i> LinkedIn
      </a>
      <a href="mailto:?subject=%s&body=%s" class="share-button share-email">
        <i class="fas fa-envelope"></i> Email
      </a>
      <a href="https://wa.me/?text=%s" target="_blank" class="share-button share-whatsapp">
        <i class="fab fa-whatsapp"></i> WhatsApp
      </a>
    </div>
  </div>
  <div class="embed-code">
    <p>Embed this report:</p>
    <textarea readonly><iframe src="%s" width="800" height="600" frameborder="0"></iframe></textarea>
  </div>
</div>
<style>
.social-share-enhanced {
  background-color: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  margin: 30px 0;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}
.social-share-enhanced h3 {
  margin-top: 0;
  color: #2c3e50;
  font-size: 20px;
  margin-bottom: 15px;
}
.share-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 20px;
}
.qr-code {
  text-align: center;
}
.qr-code img {
  border: 1px solid #ddd;
  padding: 5px;
  background: white;
  border-radius: 5px;
}
.qr-code p {
  margin-top: 5px;
  font-size: 14px;
  color: #7f8c8d;
}
.share-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}
.share-whatsapp {
  background-color: #25D366;
}
.embed-code textarea {
  width: 100%%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  height: 60px;
  font-family: monospace;
  font-size: 12px;
  resize: none;
}
</style>
`

	twitterText := url.QueryEscape(fmt.Sprintf("Check out this data quality report: %s", data.Title))
	emailSubject := url.QueryEscape(fmt.Sprintf("Data Quality Report: %s", data.Title))
	emailBody := url.QueryEscape(fmt.Sprintf("I thought you might be interested in this data quality report:\n\n%s\n\nView the report here: %s",
		data.Description, data.ReportURL))
	whatsappText := url.QueryEscape(fmt.Sprintf("%s: %s", data.Title, data.ReportURL))

	return fmt.Sprintf(html,
		data.QRCodeURL,
		twitterText, url.QueryEscape(data.ReportURL), data.Hashtags,
		url.QueryEscape(data.ReportURL),
		emailSubject, emailBody,
		whatsappText,
		data.ReportURL)
}

// Helper function to generate template variables
func generateTemplateVars(data SocialShareData) map[string]interface{} {
	return map[string]interface{}{
		"social_sharing_enabled": true,
		"qr_code_url":            data.QRCodeURL,
		"twitter_text":           fmt.Sprintf("Check out this data quality report: %s #%s", data.Title, strings.ReplaceAll(data.Hashtags, ",", " #")),
		"linkedin_text":          fmt.Sprintf("I wanted to share this data quality report: %s", data.Title),
		"email_subject":          fmt.Sprintf("Data Quality Report: %s", data.Title),
		"email_body":             fmt.Sprintf("I thought you might be interested in this data quality report:\n\n%s\n\nView the report here: %s", data.Description, data.ReportURL),
		"whatsapp_text":          fmt.Sprintf("%s: %s", data.Title, data.ReportURL),
		"embed_code":             fmt.Sprintf("<iframe src=\"%s\" width=\"800\" height=\"600\" frameborder=\"0\"></iframe>", data.ReportURL),
	}
}

// Helper function to add UTM parameters to URLs for tracking
func addUTMParams(baseURL string, source string, medium string) string {
	if !strings.Contains(baseURL, "?") {
		baseURL += "?"
	} else {
		baseURL += "&"
	}

	return fmt.Sprintf("%sutm_source=%s&utm_medium=%s&utm_campaign=nessi_report",
		baseURL, url.QueryEscape(source), url.QueryEscape(medium))
}

// Plugin is the exported symbol that Nessi will look for
var Plugin SocialSharingPlugin
