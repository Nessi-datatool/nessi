# Nessi Plugin API Reference

This document provides detailed information about Nessi's Plugin API, including interfaces, methods, and examples for developing custom plugins.

## Overview

Nessi's plugin system is designed to be extensible and flexible, allowing developers to add new functionality without modifying the core codebase. Plugins are implemented as Go packages that adhere to specific interfaces.

## Core Interfaces

### PluginInterface

All plugins must implement the `PluginInterface` interface:

```go
type PluginInterface interface {
	Initialize() error
	GetMetadata() PluginMetadata
	Execute(args []string, data interface{}) (interface{}, error)
}
```

#### Methods

- **Initialize()**: Sets up the plugin and initializes any required resources
- **GetMetadata()**: Returns metadata about the plugin (name, version, description, etc.)
- **Execute()**: Executes the plugin's functionality with the provided arguments and data

### PluginMetadata

The `PluginMetadata` struct contains information about the plugin:

```go
type PluginMetadata struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	Author       string   `json:"author"`
	Capabilities []string `json:"capabilities"`
}
```

## Specialized Interfaces

Depending on the plugin type, you may also implement one or more specialized interfaces:

### CommandExecutor

For plugins that add new commands to the Nessi CLI:

```go
type CommandExecutor interface {
	GetCommands() []*cobra.Command
	ExecuteCommand(cmd string, args []string) (interface{}, error)
}
```

### QualityRuleProvider

For plugins that add custom data quality validation rules:

```go
type QualityRuleProvider interface {
	GetRules() []QualityRule
	ValidateData(data interface{}, rule QualityRule) (RuleResult, error)
}
```

### ReportExtension

For plugins that extend report templates with custom sections:

```go
type ReportExtension interface {
	GetReportSections() []ReportSection
	GenerateContent(data interface{}, section ReportSection) (interface{}, error)
}
```

### FormatHandler

For plugins that add support for additional data formats:

```go
type FormatHandler interface {
	GetSupportedFormats() []string
	Read(path string, options map[string]interface{}) (interface{}, error)
	Write(data interface{}, path string, options map[string]interface{}) error
}
```

## Viral Growth Plugin Interfaces

Nessi provides specialized interfaces for plugins that enhance viral growth and community engagement:

### SocialSharingProvider

For plugins that add social sharing capabilities to reports:

```go
type SocialSharingProvider interface {
	EnhanceReport(data interface{}) (interface{}, error)
	GenerateShareLinks(data interface{}) (interface{}, error)
	GenerateQRCode(data interface{}) (interface{}, error)
}
```

### CommunityEngagementProvider

For plugins that facilitate community contributions and engagement:

```go
type CommunityEngagementProvider interface {
	CollectFeedback(data interface{}) (interface{}, error)
	SuggestContributions(data interface{}) (interface{}, error)
	EnhanceReportWithCommunity(data interface{}) (interface{}, error)
}
```

### BadgeProvider

For plugins that generate embeddable badges:

```go
type BadgeProvider interface {
	GenerateBadge(data interface{}) (interface{}, error)
	EnhanceReportWithBadge(data interface{}) (interface{}, error)
	GetBadgeMarkdown(data interface{}) (interface{}, error)
	GetBadgeHTML(data interface{}) (interface{}, error)
}
```

## Plugin Development Guide

### Basic Plugin Structure

Here's a basic structure for a Nessi plugin:

```go
package main

import (
	"fmt"
	"github.com/nessi-dev/nessi/pkg/plugins"
)

// MyPlugin implements the PluginInterface
type MyPlugin struct {
	Name        string
	Version     string
	Description string
	Author      string
}

// Initialize sets up the plugin
func (p *MyPlugin) Initialize() error {
	p.Name = "my_plugin"
	p.Version = "1.0.0"
	p.Description = "My custom Nessi plugin"
	p.Author = "Your Name"
	return nil
}

// GetMetadata returns plugin metadata
func (p *MyPlugin) GetMetadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        p.Name,
		Version:     p.Version,
		Description: p.Description,
		Author:      p.Author,
		Capabilities: []string{
			"custom_capability",
		},
	}
}

// Execute runs the plugin functionality
func (p *MyPlugin) Execute(args []string, data interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	switch args[0] {
	case "my_command":
		return p.myCommand(data)
	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}

// myCommand implements a custom command
func (p *MyPlugin) myCommand(data interface{}) (interface{}, error) {
	// Implement your command logic here
	return map[string]interface{}{
		"result": "Command executed successfully",
	}, nil
}

// Plugin is the exported symbol that Nessi will look for
var Plugin MyPlugin
```

### Building and Packaging

To build your plugin:

```bash
go build -buildmode=plugin -o my_plugin.so main.go
```

Your plugin should be packaged with a README.md file that explains its functionality, usage, and any configuration options.

## Viral Growth Plugin Examples

### Social Sharing Plugin Example

Here's a simplified example of a social sharing plugin:

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/nessi-dev/nessi/pkg/plugins"
)

// SocialSharingPlugin implements the PluginInterface
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
	Hashtags    string `json:"hashtags"`
}

// Initialize sets up the plugin
func (p *SocialSharingPlugin) Initialize() error {
	p.Name = "social_sharing_plugin"
	p.Version = "1.0.0"
	p.Description = "Enhances Nessi reports with social sharing capabilities"
	p.Author = "Your Name"
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
			"social_sharing",
			"report_enhancement",
		},
	}
}

// Execute runs the plugin functionality
func (p *SocialSharingPlugin) Execute(args []string, data interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	switch args[0] {
	case "enhance_report":
		return p.enhanceReport(data)
	case "generate_share_links":
		return p.generateShareLinks(data)
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

	// Generate HTML for social sharing enhancements
	htmlContent := generateSocialSharingHTML(shareData)

	return map[string]interface{}{
		"html_content": htmlContent,
		"share_data":   shareData,
	}, nil
}

// generateShareLinks creates sharing links for various platforms
func (p *SocialSharingPlugin) generateShareLinks(data interface{}) (interface{}, error) {
	// Implementation details omitted for brevity
	return map[string]interface{}{
		"twitter_link":  "https://twitter.com/intent/tweet?text=...",
		"linkedin_link": "https://www.linkedin.com/sharing/share-offsite/?url=...",
		"email_link":    "mailto:?subject=...&body=...",
	}, nil
}

// Helper function to generate HTML for social sharing features
func generateSocialSharingHTML(data SocialShareData) string {
	// Implementation details omitted for brevity
	return "<div class=\"social-share\">...</div>"
}

// Plugin is the exported symbol that Nessi will look for
var Plugin SocialSharingPlugin
```

## Best Practices

### Error Handling

Plugins should handle errors gracefully and return meaningful error messages. Use the `fmt.Errorf` function to create descriptive error messages.

### Documentation

Each plugin should include comprehensive documentation, including:

- A README.md file with usage instructions
- Code comments explaining the plugin's functionality
- Examples of how to use the plugin

### Testing

Plugins should include tests to ensure they work correctly. Use Go's testing package to write unit tests for your plugin.

### Versioning

Use semantic versioning (MAJOR.MINOR.PATCH) for your plugin versions. Increment the:

- MAJOR version when you make incompatible API changes
- MINOR version when you add functionality in a backward-compatible manner
- PATCH version when you make backward-compatible bug fixes

## Plugin Distribution

Plugins can be distributed in several ways:

1. **GitHub Repository**: Host your plugin on GitHub and allow users to install it using `nessi plugins install github.com/username/plugin`
2. **Pre-built Binaries**: Provide pre-built binaries for different platforms
3. **Package Managers**: Distribute your plugin through package managers like Homebrew or apt

## Community Guidelines

When developing plugins for Nessi, please follow these community guidelines:

1. **Be Respectful**: Respect the Nessi community and its users
2. **Follow Go Best Practices**: Write clean, idiomatic Go code
3. **Provide Documentation**: Include comprehensive documentation with your plugin
4. **Respond to Issues**: Address issues and pull requests in a timely manner
5. **Share Your Work**: Share your plugin with the Nessi community

## Conclusion

Nessi's plugin system provides a powerful way to extend its functionality. By developing plugins, you can contribute to the Nessi ecosystem and help improve the tool for everyone.

For more information, see the [Plugins Documentation](PLUGINS.md) or join the Nessi community channels.
