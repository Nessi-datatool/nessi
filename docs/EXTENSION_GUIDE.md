# Nessi Extension Guide

This guide provides comprehensive information on extending and customizing Nessi's functionality through plugins, custom rules, report templates, and integrations.

## Table of Contents

- [Overview](#overview)
- [Plugin System](#plugin-system)
- [Custom Quality Rules](#custom-quality-rules)
- [Custom Report Templates](#custom-report-templates)
- [Custom Integrations](#custom-integrations)
- [Custom Commands](#custom-commands)
- [Custom Metrics](#custom-metrics)
- [Extending the API](#extending-the-api)
- [Packaging and Distribution](#packaging-and-distribution)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

Nessi is designed to be extensible, allowing you to customize and extend its functionality to meet your specific needs. This guide covers the various extension points available in Nessi and provides examples for each.

## Plugin System

Nessi's plugin system allows you to extend its functionality without modifying the core codebase. Plugins are loaded dynamically at runtime and can add new commands, quality rules, report templates, and integrations.

### Plugin Structure

A basic plugin structure:

```
my-plugin/
├── plugin.yaml         # Plugin metadata
├── plugin.go           # Plugin entry point
├── rules/              # Custom quality rules
│   └── my_rule.go
├── templates/          # Custom report templates
│   └── my_template.html
├── commands/           # Custom commands
│   └── my_command.go
└── README.md           # Plugin documentation
```

### Plugin Metadata

The `plugin.yaml` file contains metadata about your plugin:

```yaml
name: my-plugin
version: 1.0.0
description: My custom Nessi plugin
author: Your Name
email: your.email@example.com
website: https://example.com
license: Apache-2.0
requires:
  nessi: ">=1.0.0"
provides:
  rules:
    - my_rule
  templates:
    - my_template
  commands:
    - my_command
```

### Plugin Entry Point

The `plugin.go` file is the entry point for your plugin:

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/plugin"
)

// Plugin is the main plugin struct
type Plugin struct{}

// Init initializes the plugin
func (p *Plugin) Init() error {
    // Initialization code
    return nil
}

// Name returns the plugin name
func (p *Plugin) Name() string {
    return "my-plugin"
}

// Version returns the plugin version
func (p *Plugin) Version() string {
    return "1.0.0"
}

// Register registers the plugin components
func (p *Plugin) Register(registry plugin.Registry) error {
    // Register rules
    registry.RegisterRule("my_rule", &MyRule{})
    
    // Register templates
    registry.RegisterTemplate("my_template", "templates/my_template.html")
    
    // Register commands
    registry.RegisterCommand("my_command", &MyCommand{})
    
    return nil
}

// Plugin entry point
var Plugin plugin.Plugin = &Plugin{}
```

### Building a Plugin

```bash
go build -buildmode=plugin -o my-plugin.so my-plugin.go
```

### Installing a Plugin

```bash
mkdir -p ~/.nessi/plugins
cp my-plugin.so ~/.nessi/plugins/
```

### Using a Plugin

```bash
nessi --plugin my-plugin my-command
```

## Custom Quality Rules

Custom quality rules allow you to define your own data quality validations.

### Rule Interface

Custom rules must implement the `Rule` interface:

```go
type Rule interface {
    Name() string
    Description() string
    Validate(data []byte, parameters map[string]interface{}) (Result, error)
}
```

### Example Rule Implementation

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/quality"
)

// EmailRule validates email addresses
type EmailRule struct{}

// Name returns the rule name
func (r *EmailRule) Name() string {
    return "email_format"
}

// Description returns the rule description
func (r *EmailRule) Description() string {
    return "Validates that values match email format"
}

// Validate implements the validation logic
func (r *EmailRule) Validate(data []byte, parameters map[string]interface{}) (quality.Result, error) {
    // Parse parameters
    columns, ok := parameters["columns"].([]string)
    if !ok {
        return quality.Result{}, fmt.Errorf("columns parameter is required")
    }
    
    // Parse data
    records, err := parseData(data)
    if err != nil {
        return quality.Result{}, err
    }
    
    // Validate each column
    results := make(map[string]quality.ColumnResult)
    for _, column := range columns {
        validCount := 0
        invalidValues := []string{}
        
        for _, record := range records {
            value := record[column]
            if isValidEmail(value) {
                validCount++
            } else {
                invalidValues = append(invalidValues, value)
            }
        }
        
        // Calculate quality score
        score := float64(validCount) / float64(len(records))
        
        // Create column result
        results[column] = quality.ColumnResult{
            Valid:         score >= 0.95,
            Score:         score,
            InvalidValues: invalidValues[:min(10, len(invalidValues))],
            Message:       fmt.Sprintf("%.2f%% of values are valid emails", score*100),
        }
    }
    
    // Create overall result
    return quality.Result{
        Valid:        allColumnsValid(results),
        ColumnResults: results,
    }, nil
}

// Helper function to validate email format
func isValidEmail(email string) bool {
    // Implement email validation logic
    return regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
}
```

### Using Custom Rules

Custom rules can be used in quality rule definitions:

```yaml
# quality_rules.yaml
rules:
  - name: email_format
    description: Check if email columns have valid format
    columns:
      - email
      - contact_email
```

## Custom Report Templates

Custom report templates allow you to create your own report formats.

### Template Structure

Report templates are HTML files with Go template syntax:

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
    <style>
        /* Custom CSS styles */
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
        }
        .header {
            background-color: #2c3e50;
            color: white;
            padding: 10px;
            border-radius: 5px;
        }
        .summary {
            margin: 20px 0;
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 5px;
        }
        .metrics {
            display: flex;
            flex-wrap: wrap;
        }
        .metric {
            flex: 1;
            min-width: 200px;
            margin: 10px;
            padding: 15px;
            background-color: #e9ecef;
            border-radius: 5px;
        }
        .good {
            color: #28a745;
        }
        .warning {
            color: #ffc107;
        }
        .error {
            color: #dc3545;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{ .Title }}</h1>
        <p>Generated on {{ .GeneratedAt }}</p>
    </div>
    
    <div class="summary">
        <h2>Summary</h2>
        <p>Table: {{ .TableName }}</p>
        <p>Path: {{ .TablePath }}</p>
        <p>Overall Quality Score: 
            <span class="{{ if ge .OverallScore 0.9 }}good{{ else if ge .OverallScore 0.7 }}warning{{ else }}error{{ end }}">
                {{ printf "%.2f%%" (mul .OverallScore 100) }}
            </span>
        </p>
    </div>
    
    <div class="metrics">
        {{ range .Metrics }}
        <div class="metric">
            <h3>{{ .Name }}</h3>
            <p>Score: 
                <span class="{{ if ge .Score 0.9 }}good{{ else if ge .Score 0.7 }}warning{{ else }}error{{ end }}">
                    {{ printf "%.2f%%" (mul .Score 100) }}
                </span>
            </p>
            <p>{{ .Description }}</p>
        </div>
        {{ end }}
    </div>
    
    <div class="details">
        <h2>Quality Details</h2>
        {{ range .RuleResults }}
        <div class="rule">
            <h3>{{ .Name }}</h3>
            <p>{{ .Description }}</p>
            <p>Status: 
                <span class="{{ if .Valid }}good{{ else }}error{{ end }}">
                    {{ if .Valid }}Passed{{ else }}Failed{{ end }}
                </span>
            </p>
            
            {{ if not .Valid }}
            <div class="failures">
                <h4>Failures</h4>
                <ul>
                    {{ range .Failures }}
                    <li>{{ . }}</li>
                    {{ end }}
                </ul>
            </div>
            {{ end }}
        </div>
        {{ end }}
    </div>
</body>
</html>
```

### Template Functions

Nessi provides several template functions for use in report templates:

- `mul`: Multiply two numbers
- `div`: Divide two numbers
- `add`: Add two numbers
- `sub`: Subtract two numbers
- `formatDate`: Format a date string
- `formatNumber`: Format a number
- `formatPercent`: Format a percentage
- `formatBytes`: Format a byte size
- `formatDuration`: Format a duration

### Registering Templates

Register your template in your plugin:

```go
func (p *Plugin) Register(registry plugin.Registry) error {
    registry.RegisterTemplate("my_template", "templates/my_template.html")
    return nil
}
```

### Using Custom Templates

Use your custom template when generating reports:

```bash
nessi report generate --path /path/to/table --template my_template --output /path/to/report.html
```

## Custom Integrations

Custom integrations allow you to connect Nessi to external systems.

### Integration Interface

Custom integrations must implement the `Integration` interface:

```go
type Integration interface {
    Name() string
    Description() string
    Connect(parameters map[string]interface{}) error
    Disconnect() error
    ListResources() ([]Resource, error)
    GetResource(id string) (Resource, error)
}
```

### Example Integration Implementation

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/integration"
)

// MyIntegration is a custom integration
type MyIntegration struct {
    client *MyClient
    connected bool
}

// Name returns the integration name
func (i *MyIntegration) Name() string {
    return "my_integration"
}

// Description returns the integration description
func (i *MyIntegration) Description() string {
    return "Integration with My System"
}

// Connect establishes a connection to the external system
func (i *MyIntegration) Connect(parameters map[string]interface{}) error {
    // Extract parameters
    host, ok := parameters["host"].(string)
    if !ok {
        return fmt.Errorf("host parameter is required")
    }
    
    apiKey, ok := parameters["api_key"].(string)
    if !ok {
        return fmt.Errorf("api_key parameter is required")
    }
    
    // Create client
    client, err := NewMyClient(host, apiKey)
    if err != nil {
        return err
    }
    
    // Test connection
    if err := client.Ping(); err != nil {
        return err
    }
    
    // Store client
    i.client = client
    i.connected = true
    
    return nil
}

// Disconnect closes the connection
func (i *MyIntegration) Disconnect() error {
    if !i.connected {
        return nil
    }
    
    if err := i.client.Close(); err != nil {
        return err
    }
    
    i.connected = false
    return nil
}

// ListResources lists available resources
func (i *MyIntegration) ListResources() ([]integration.Resource, error) {
    if !i.connected {
        return nil, fmt.Errorf("not connected")
    }
    
    // Get resources from external system
    myResources, err := i.client.ListResources()
    if err != nil {
        return nil, err
    }
    
    // Convert to Nessi resources
    resources := make([]integration.Resource, len(myResources))
    for j, res := range myResources {
        resources[j] = integration.Resource{
            ID:   res.ID,
            Name: res.Name,
            Type: res.Type,
            Metadata: map[string]interface{}{
                "created_at": res.CreatedAt,
                "owner":      res.Owner,
                "size":       res.Size,
            },
        }
    }
    
    return resources, nil
}

// GetResource gets a specific resource
func (i *MyIntegration) GetResource(id string) (integration.Resource, error) {
    if !i.connected {
        return integration.Resource{}, fmt.Errorf("not connected")
    }
    
    // Get resource from external system
    myResource, err := i.client.GetResource(id)
    if err != nil {
        return integration.Resource{}, err
    }
    
    // Convert to Nessi resource
    resource := integration.Resource{
        ID:   myResource.ID,
        Name: myResource.Name,
        Type: myResource.Type,
        Metadata: map[string]interface{}{
            "created_at": myResource.CreatedAt,
            "owner":      myResource.Owner,
            "size":       myResource.Size,
        },
        Data: myResource.Data,
    }
    
    return resource, nil
}
```

### Registering Integrations

Register your integration in your plugin:

```go
func (p *Plugin) Register(registry plugin.Registry) error {
    registry.RegisterIntegration("my_integration", &MyIntegration{})
    return nil
}
```

### Using Custom Integrations

Use your custom integration in Nessi commands:

```bash
nessi integration my_integration list-resources --host example.com --api-key your-api-key
```

## Custom Commands

Custom commands allow you to add new functionality to the Nessi CLI.

### Command Interface

Custom commands must implement the `Command` interface:

```go
type Command interface {
    Name() string
    Description() string
    Execute(args []string, flags map[string]interface{}) error
}
```

### Example Command Implementation

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/cli"
)

// MyCommand is a custom command
type MyCommand struct{}

// Name returns the command name
func (c *MyCommand) Name() string {
    return "my-command"
}

// Description returns the command description
func (c *MyCommand) Description() string {
    return "My custom command"
}

// Execute implements the command logic
func (c *MyCommand) Execute(args []string, flags map[string]interface{}) error {
    // Parse flags
    verbose, _ := flags["verbose"].(bool)
    output, _ := flags["output"].(string)
    
    // Command logic
    fmt.Println("Executing my custom command")
    
    if verbose {
        fmt.Println("Verbose mode enabled")
    }
    
    // Process arguments
    for i, arg := range args {
        fmt.Printf("Argument %d: %s\n", i, arg)
    }
    
    // Generate output
    result := map[string]interface{}{
        "status": "success",
        "message": "Command executed successfully",
        "data": args,
    }
    
    // Format output
    if output == "json" {
        jsonResult, err := json.Marshal(result)
        if err != nil {
            return err
        }
        fmt.Println(string(jsonResult))
    } else {
        fmt.Println("Status: success")
        fmt.Println("Message: Command executed successfully")
        fmt.Println("Data:", args)
    }
    
    return nil
}

// GetFlags returns the command flags
func (c *MyCommand) GetFlags() []cli.Flag {
    return []cli.Flag{
        {
            Name:        "verbose",
            Shorthand:   "v",
            Type:        "bool",
            Default:     false,
            Description: "Enable verbose output",
        },
        {
            Name:        "output",
            Shorthand:   "o",
            Type:        "string",
            Default:     "text",
            Description: "Output format (text, json)",
        },
    }
}
```

### Registering Commands

Register your command in your plugin:

```go
func (p *Plugin) Register(registry plugin.Registry) error {
    registry.RegisterCommand("my-command", &MyCommand{})
    return nil
}
```

### Using Custom Commands

Use your custom command:

```bash
nessi my-command --verbose --output json arg1 arg2
```

## Custom Metrics

Custom metrics allow you to define your own metrics for data quality and performance.

### Metric Interface

Custom metrics must implement the `Metric` interface:

```go
type Metric interface {
    Name() string
    Description() string
    Calculate(data []byte, parameters map[string]interface{}) (float64, error)
    Format(value float64) string
}
```

### Example Metric Implementation

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/metrics"
)

// CardinalityMetric calculates the cardinality of a column
type CardinalityMetric struct{}

// Name returns the metric name
func (m *CardinalityMetric) Name() string {
    return "cardinality"
}

// Description returns the metric description
func (m *CardinalityMetric) Description() string {
    return "Calculates the cardinality (number of unique values) of a column"
}

// Calculate implements the metric calculation logic
func (m *CardinalityMetric) Calculate(data []byte, parameters map[string]interface{}) (float64, error) {
    // Parse parameters
    column, ok := parameters["column"].(string)
    if !ok {
        return 0, fmt.Errorf("column parameter is required")
    }
    
    // Parse data
    records, err := parseData(data)
    if err != nil {
        return 0, err
    }
    
    // Calculate cardinality
    uniqueValues := make(map[string]struct{})
    for _, record := range records {
        value := record[column]
        uniqueValues[value] = struct{}{}
    }
    
    // Return cardinality
    return float64(len(uniqueValues)), nil
}

// Format formats the metric value
func (m *CardinalityMetric) Format(value float64) string {
    return fmt.Sprintf("%.0f unique values", value)
}
```

### Registering Metrics

Register your metric in your plugin:

```go
func (p *Plugin) Register(registry plugin.Registry) error {
    registry.RegisterMetric("cardinality", &CardinalityMetric{})
    return nil
}
```

### Using Custom Metrics

Use your custom metric in quality checks:

```yaml
# quality_rules.yaml
rules:
  - name: high_cardinality
    description: Check if column has high cardinality
    columns:
      - id
    parameters:
      metric: cardinality
      min_value: 1000
```

## Extending the API

Nessi provides extension points for the API (Pro Edition).

### REST API Extensions

To extend the REST API:

1. Define your API endpoints in OpenAPI format
2. Implement the endpoint handlers
3. Register your endpoints with the API server

Example endpoint definition:

```yaml
# my_api.yaml
paths:
  /api/v1/my-extension:
    get:
      summary: My custom endpoint
      description: Returns custom data
      parameters:
        - name: param1
          in: query
          required: true
          schema:
            type: string
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                  data:
                    type: object
```

Example endpoint implementation:

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/api"
)

// MyEndpoint is a custom API endpoint
type MyEndpoint struct{}

// Handle implements the endpoint handler
func (e *MyEndpoint) Handle(req api.Request) (api.Response, error) {
    // Get parameters
    param1 := req.Query.Get("param1")
    
    // Process request
    data := map[string]interface{}{
        "param1": param1,
        "timestamp": time.Now().Unix(),
    }
    
    // Return response
    return api.Response{
        StatusCode: 200,
        Body: map[string]interface{}{
            "status": "success",
            "data": data,
        },
    }, nil
}
```

Register your endpoint in your plugin:

```go
func (p *Plugin) Register(registry plugin.Registry) error {
    registry.RegisterAPIEndpoint("/api/v1/my-extension", &MyEndpoint{})
    return nil
}
```

## Packaging and Distribution

### Plugin Packaging

Package your plugin for distribution:

```bash
# Create a distribution directory
mkdir -p dist/my-plugin

# Copy plugin files
cp plugin.yaml dist/my-plugin/
cp my-plugin.so dist/my-plugin/
cp -r rules/ dist/my-plugin/
cp -r templates/ dist/my-plugin/
cp README.md dist/my-plugin/

# Create a zip archive
cd dist
zip -r my-plugin.zip my-plugin/
```

### Plugin Distribution

Distribute your plugin:

1. Host the plugin zip file on a web server or GitHub
2. Create a plugin manifest file:

```yaml
# plugin-manifest.yaml
plugins:
  - name: my-plugin
    version: 1.0.0
    description: My custom Nessi plugin
    author: Your Name
    url: https://example.com/plugins/my-plugin.zip
    sha256: 1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
```

3. Share the plugin manifest with users

### Plugin Installation

Users can install your plugin:

```bash
# Install from a local file
nessi plugin install --file my-plugin.zip

# Install from a URL
nessi plugin install --url https://example.com/plugins/my-plugin.zip

# Install from a manifest
nessi plugin install --manifest plugin-manifest.yaml --plugin my-plugin
```

## Best Practices

### Plugin Development

- Follow Go best practices for code organization and style
- Use meaningful names for your plugin components
- Provide comprehensive documentation
- Include examples of how to use your plugin
- Write tests for your plugin components
- Handle errors gracefully
- Respect Nessi's plugin API contracts

### Performance Considerations

- Optimize your code for performance
- Use efficient algorithms and data structures
- Minimize memory allocations
- Use parallel processing where appropriate
- Profile your code to identify bottlenecks
- Consider resource usage in your plugin

### Security Considerations

- Validate all user input
- Handle sensitive data securely
- Use secure connections for external systems
- Follow the principle of least privilege
- Document security considerations for your plugin

## Troubleshooting

### Common Issues

#### Plugin Not Loading

**Issue**: Nessi doesn't load your plugin

**Solution**:
1. Verify the plugin is in the correct location (`~/.nessi/plugins/`)
2. Check that the plugin is built with the correct Go version
3. Ensure the plugin implements the required interfaces
4. Check for errors in the plugin initialization

#### Plugin Compatibility

**Issue**: Plugin is not compatible with your Nessi version

**Solution**:
1. Check the plugin's required Nessi version in `plugin.yaml`
2. Update the plugin to be compatible with your Nessi version
3. Use a compatible version of Nessi

#### Plugin Performance

**Issue**: Plugin causes performance issues

**Solution**:
1. Profile the plugin to identify bottlenecks
2. Optimize the plugin code
3. Consider using sampling or batching for large datasets
4. Use more efficient algorithms or data structures

### Debugging Plugins

Enable plugin debugging:

```bash
# Enable plugin debugging
export NESSI_PLUGIN_DEBUG=1
nessi --plugin my-plugin my-command
```

Use the plugin development mode:

```bash
# Run Nessi in plugin development mode
nessi --plugin-dev-mode --plugin-path /path/to/plugin/directory my-command
```

---

For more information on extending Nessi, please visit [nessi.dev/extensions](https://nessi.dev/extensions) or contact support@nessi.dev.
