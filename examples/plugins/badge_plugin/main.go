package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/nessi-dev/nessi/pkg/plugins"
)

// BadgePlugin implements the PluginInterface for generating embeddable badges
type BadgePlugin struct {
	Name        string
	Version     string
	Description string
	Author      string
}

// BadgeData contains data for badge generation
type BadgeData struct {
	Label        string `json:"label"`
	Message      string `json:"message"`
	Color        string `json:"color"`
	Style        string `json:"style"`
	Logo         string `json:"logo"`
	LogoWidth    int    `json:"logo_width"`
	LinkURL      string `json:"link_url"`
	AltText      string `json:"alt_text"`
	EmbedType    string `json:"embed_type"`
	QualityScore int    `json:"quality_score"`
}

// Initialize sets up the plugin
func (p *BadgePlugin) Initialize() error {
	p.Name = "badge_plugin"
	p.Version = "1.0.0"
	p.Description = "Generates embeddable 'Powered by Nessi' badges for reports and projects"
	p.Author = "Nessi Community"
	return nil
}

// GetMetadata returns plugin metadata
func (p *BadgePlugin) GetMetadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Capabilities: []string{
			"badge_generation",
			"report_enhancement",
		},
	}
}

// Execute runs the plugin functionality
func (p *BadgePlugin) Execute(args []string, data interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified for badge plugin")
	}

	switch args[0] {
	case "generate_badge":
		return p.generateBadge(data)
	case "enhance_report":
		return p.enhanceReportWithBadge(data)
	case "get_badge_markdown":
		return p.getBadgeMarkdown(data)
	case "get_badge_html":
		return p.getBadgeHTML(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}

// generateBadge creates a badge based on the provided data
func (p *BadgePlugin) generateBadge(data interface{}) (interface{}, error) {
	// Parse input data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input data: %v", err)
	}

	var badgeData BadgeData
	if err := json.Unmarshal(dataBytes, &badgeData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
	}

	// Set defaults if not provided
	if badgeData.Label == "" {
		badgeData.Label = "powered by"
	}

	if badgeData.Message == "" {
		badgeData.Message = "nessi"
	}

	if badgeData.Color == "" {
		badgeData.Color = "3498db"
	}

	if badgeData.Style == "" {
		badgeData.Style = "flat"
	}

	if badgeData.Logo == "" {
		// Base64 encoded small Nessi logo
		badgeData.Logo = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmZmZmYiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIj48cGF0aCBkPSJNMjAgMTRINFYxMGE2IDYgMCAwIDEgMTItMHY0eiIvPjxwYXRoIGQ9Ik0xMiAxNnY0Ii8+PHBhdGggZD0iTTggMTh2MiIvPjxwYXRoIGQ9Ik0xNiAxOHYyIi8+PC9zdmc+"
	}

	if badgeData.LogoWidth == 0 {
		badgeData.LogoWidth = 16
	}

	if badgeData.LinkURL == "" {
		badgeData.LinkURL = "https://github.com/nessi-dev/nessi"
	}

	if badgeData.AltText == "" {
		badgeData.AltText = "Powered by Nessi"
	}

	if badgeData.EmbedType == "" {
		badgeData.EmbedType = "markdown"
	}

	// Generate badge URL using shields.io
	badgeURL := fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s?style=%s&logo=%s&logoWidth=%d",
		url.QueryEscape(badgeData.Label),
		url.QueryEscape(badgeData.Message),
		badgeData.Color,
		badgeData.Style,
		base64.StdEncoding.EncodeToString([]byte(badgeData.Logo)),
		badgeData.LogoWidth)

	// Generate badge code based on embed type
	var badgeCode string
	switch badgeData.EmbedType {
	case "markdown":
		badgeCode = fmt.Sprintf("[![%s](%s)](%s)", badgeData.AltText, badgeURL, badgeData.LinkURL)
	case "html":
		badgeCode = fmt.Sprintf("<a href=\"%s\"><img src=\"%s\" alt=\"%s\"></a>", badgeData.LinkURL, badgeURL, badgeData.AltText)
	case "rst":
		badgeCode = fmt.Sprintf(".. image:: %s\n   :target: %s\n   :alt: %s", badgeURL, badgeData.LinkURL, badgeData.AltText)
	default:
		badgeCode = fmt.Sprintf("[![%s](%s)](%s)", badgeData.AltText, badgeURL, badgeData.LinkURL)
	}

	return map[string]interface{}{
		"badge_url":  badgeURL,
		"badge_code": badgeCode,
		"embed_type": badgeData.EmbedType,
		"link_url":   badgeData.LinkURL,
	}, nil
}

// enhanceReportWithBadge adds badge elements to reports
func (p *BadgePlugin) enhanceReportWithBadge(data interface{}) (interface{}, error) {
	// Parse input data if provided
	var badgeData BadgeData
	if data != nil {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input data: %v", err)
		}

		if err := json.Unmarshal(dataBytes, &badgeData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal input data: %v", err)
		}
	}

	// Set quality score badge if provided
	qualityBadge := ""
	if badgeData.QualityScore > 0 {
		qualityColor := "red"
		if badgeData.QualityScore >= 90 {
			qualityColor = "brightgreen"
		} else if badgeData.QualityScore >= 70 {
			qualityColor = "green"
		} else if badgeData.QualityScore >= 50 {
			qualityColor = "yellow"
		} else if badgeData.QualityScore >= 30 {
			qualityColor = "orange"
		}

		qualityBadge = fmt.Sprintf("https://img.shields.io/badge/quality-%d%%25-%s",
			badgeData.QualityScore, qualityColor)
	}

	// Generate HTML for badge section
	htmlContent := `
<div class="nessi-badge-section">
  <h3>Powered by Nessi</h3>
  <p>This report was generated using Nessi, an open-source Delta Lake quality tool.</p>
  
  <div class="badge-container">
    <div class="badge-display">
      <img src="https://img.shields.io/badge/powered%20by-nessi-3498db?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IiNmZmZmZmYiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIj48cGF0aCBkPSJNMjAgMTRINFYxMGE2IDYgMCAwIDEgMTItMHY0eiIvPjxwYXRoIGQ9Ik0xMiAxNnY0Ii8+PHBhdGggZD0iTTggMTh2MiIvPjxwYXRoIGQ9Ik0xNiAxOHYyIi8+PC9zdmc+" alt="Powered by Nessi">
`

	// Add quality badge if available
	if qualityBadge != "" {
		htmlContent += fmt.Sprintf("      <img src=\"%s\" alt=\"Quality Score: %d percent\" style=\"margin-left: 10px;\">\n",
			qualityBadge, badgeData.QualityScore)
	}

	htmlContent += `    </div>
    
    <div class="badge-embed">
      <p>Add this badge to your project:</p>
      <div class="code-tabs">
        <div class="tab active" data-tab="markdown">Markdown</div>
        <div class="tab" data-tab="html">HTML</div>
        <div class="tab" data-tab="rst">RST</div>
      </div>
      
      <div class="tab-content active" id="markdown">
        <pre><code>[![Powered by Nessi](https://img.shields.io/badge/powered%20by-nessi-3498db?style=flat)](https://github.com/nessi-dev/nessi)</code></pre>
        <button class="copy-button" onclick="copyToClipboard('markdown')">Copy</button>
      </div>
      
      <div class="tab-content" id="html">
        <pre><code>&lt;a href="https://github.com/nessi-dev/nessi"&gt;&lt;img src="https://img.shields.io/badge/powered%20by-nessi-3498db?style=flat" alt="Powered by Nessi"&gt;&lt;/a&gt;</code></pre>
        <button class="copy-button" onclick="copyToClipboard('html')">Copy</button>
      </div>
      
      <div class="tab-content" id="rst">
        <pre><code>.. image:: https://img.shields.io/badge/powered%20by-nessi-3498db?style=flat
   :target: https://github.com/nessi-dev/nessi
   :alt: Powered by Nessi</code></pre>
        <button class="copy-button" onclick="copyToClipboard('rst')">Copy</button>
      </div>
    </div>
  </div>
  
  <p class="badge-note">Using Nessi in your projects? Let us know by <a href="https://github.com/nessi-dev/nessi/discussions/new?category=show-and-tell" target="_blank">sharing your story</a>!</p>
</div>

<style>
.nessi-badge-section {
  background-color: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  margin: 30px 0;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}
.nessi-badge-section h3 {
  margin-top: 0;
  color: #2c3e50;
  font-size: 20px;
  margin-bottom: 15px;
}
.badge-container {
  background: white;
  border-radius: 8px;
  padding: 15px;
  margin-top: 20px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.05);
}
.badge-display {
  text-align: center;
  margin-bottom: 20px;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 5px;
}
.badge-embed p {
  margin-bottom: 10px;
  font-weight: 600;
}
.code-tabs {
  display: flex;
  border-bottom: 1px solid #ddd;
  margin-bottom: 15px;
}
.tab {
  padding: 8px 15px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
}
.tab.active {
  border-bottom: 2px solid #3498db;
  color: #3498db;
}
.tab-content {
  display: none;
  position: relative;
}
.tab-content.active {
  display: block;
}
.tab-content pre {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 5px;
  overflow-x: auto;
  margin: 0;
}
.tab-content code {
  font-family: monospace;
  font-size: 12px;
}
.copy-button {
  position: absolute;
  top: 5px;
  right: 5px;
  background: #3498db;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 3px 8px;
  font-size: 12px;
  cursor: pointer;
}
.badge-note {
  margin-top: 15px;
  font-size: 14px;
  color: #7f8c8d;
  text-align: center;
}
</style>

<script>
function copyToClipboard(tabId) {
  const codeElement = document.querySelector('#' + tabId + ' code');
  const textArea = document.createElement('textarea');
  textArea.value = codeElement.textContent;
  document.body.appendChild(textArea);
  textArea.select();
  document.execCommand('copy');
  document.body.removeChild(textArea);
  
  const button = document.querySelector('#' + tabId + ' .copy-button');
  const originalText = button.textContent;
  button.textContent = 'Copied!';
  setTimeout(function() {
    button.textContent = originalText;
  }, 2000);
}

document.querySelectorAll('.tab').forEach(tab => {
  tab.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
    
    tab.classList.add('active');
    const tabId = tab.getAttribute('data-tab');
    document.getElementById(tabId).classList.add('active');
  });
});
</script>
`

	return map[string]interface{}{
		"html_content": htmlContent,
		"template_vars": map[string]interface{}{
			"badge_enabled":  true,
			"badge_url":      "https://img.shields.io/badge/powered%20by-nessi-3498db?style=flat",
			"badge_markdown": "[![Powered by Nessi](https://img.shields.io/badge/powered%20by-nessi-3498db?style=flat)](https://github.com/nessi-dev/nessi)",
			"badge_html":     "<a href=\"https://github.com/nessi-dev/nessi\"><img src=\"https://img.shields.io/badge/powered%20by-nessi-3498db?style=flat\" alt=\"Powered by Nessi\"></a>",
		},
	}, nil
}

// getBadgeMarkdown returns markdown code for embedding a badge
func (p *BadgePlugin) getBadgeMarkdown(data interface{}) (interface{}, error) {
	// Generate badge data
	result, err := p.generateBadge(data)
	if err != nil {
		return nil, err
	}

	// Extract badge code
	badgeData, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	// Force markdown format
	badgeURL, _ := badgeData["badge_url"].(string)
	linkURL, _ := badgeData["link_url"].(string)

	return map[string]interface{}{
		"markdown": fmt.Sprintf("[![Powered by Nessi](%s)](%s)", badgeURL, linkURL),
	}, nil
}

// getBadgeHTML returns HTML code for embedding a badge
func (p *BadgePlugin) getBadgeHTML(data interface{}) (interface{}, error) {
	// Generate badge data
	result, err := p.generateBadge(data)
	if err != nil {
		return nil, err
	}

	// Extract badge code
	badgeData, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	// Force HTML format
	badgeURL, _ := badgeData["badge_url"].(string)
	linkURL, _ := badgeData["link_url"].(string)

	return map[string]interface{}{
		"html": fmt.Sprintf("<a href=\"%s\"><img src=\"%s\" alt=\"Powered by Nessi\"></a>", linkURL, badgeURL),
	}, nil
}

// Plugin is the exported symbol that Nessi will look for
var Plugin BadgePlugin
