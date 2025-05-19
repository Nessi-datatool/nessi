# Nessi Plugin System

Nessi supports a flexible plugin system for extending core functionality. Plugins can add new data quality rules, integrate with external systems, or provide custom reporting and analytics.

## Types of Plugins
- **Python Extensions:** For ML, advanced analytics, and custom rules (see `docs/python_extensions.md`)
- **Go Plugins:** For performance-critical extensions or native integrations

## Building a Plugin
Nessi supports the following types of plugins:

- **Validation Plugins**: Custom validation rules for data quality checks
- **Alert Plugins**: Custom alerting mechanisms for notifications
- **Metric Plugins**: Custom metrics collection for monitoring
- **Storage Plugins**: Custom storage backends for data persistence
- **Export Plugins**: Custom export formats for data sharing
- **UI Plugins**: Custom UI components for visualization

## Plugin Management

Nessi provides a set of commands to manage plugins:

```bash
# List all installed plugins
nessi plugin list

# Install a plugin from a shared library file
nessi plugin install /path/to/plugin.so

# Uninstall a plugin
nessi plugin uninstall plugin-name

# Get detailed information about a plugin
nessi plugin info plugin-name

# Enable a plugin
nessi plugin enable plugin-name

# Disable a plugin
nessi plugin disable plugin-name

# List supported plugin types
nessi plugin types

# Execute a function in a plugin
nessi plugin exec plugin-name function-name [args...]
```

## Creating a Plugin

Plugins for Nessi are written in Go and compiled as shared libraries (`.so` files). Here's how to create a simple validation plugin:

### 1. Create a new Go module

```bash
mkdir my-validation-plugin
cd my-validation-plugin
go mod init github.com/yourusername/my-validation-plugin
```

### 2. Create the plugin code

Create a file named `plugin.go` with the following content:

```go
package main

import (
	"fmt"
)

// PluginInfo contains metadata about this plugin
var PluginInfo = &PluginInfo{
	Name:        "my-validation-plugin",
	Version:     "1.0.0",
	Description: "A custom validation plugin for Nessi",
	Author:      "Your Name",
	Type:        "validation",
	Enabled:     true,
}

// PluginInfo represents the structure for plugin metadata
type PluginInfo struct {
	Name        string
	Version     string
	Description string
	Author      string
	Type        string
	Enabled     bool
}

// Validate is the main function for validation plugins
func Validate(args ...interface{}) ([]interface{}, error) {
	if len(args) < 1 {
		return []interface{}{false, "No data provided"}, nil
	}

	// Implement your validation logic here
	// For example, check if a string is not empty
	data, ok := args[0].(string)
	if !ok {
		return []interface{}{false, "Data must be a string"}, nil
	}

	if data == "" {
		return []interface{}{false, "Data cannot be empty"}, nil
	}

	return []interface{}{true, "Validation passed"}, nil
}

// GetOptions returns the available options for this plugin
func GetOptions() map[string]interface{} {
	return map[string]interface{}{
		"option1": "default value",
		"option2": 42,
	}
}

// main is required for Go plugins
func main() {}
```

### 3. Build the plugin

```bash
go build -buildmode=plugin -o my-validation-plugin.so plugin.go
```

### 4. Install the plugin

```bash
nessi plugin install /path/to/my-validation-plugin.so
```

## Plugin Interfaces

Depending on the type of plugin you're creating, you'll need to implement different interfaces:

### Validation Plugin

```go
// Validate validates data against custom rules
func Validate(args ...interface{}) ([]interface{}, error)
```

### Alert Plugin

```go
// SendAlert sends an alert to a custom destination
func SendAlert(alert interface{}, options map[string]interface{}) error
```

### Metric Plugin

```go
// CollectMetrics collects custom metrics
func CollectMetrics(options map[string]interface{}) (map[string]interface{}, error)
```

### Storage Plugin

```go
// Store stores data in a custom backend
func Store(key string, data interface{}, options map[string]interface{}) error

// Retrieve retrieves data from a custom backend
func Retrieve(key string, options map[string]interface{}) (interface{}, error)
```

### Export Plugin

```go
// Export exports data to a custom format
func Export(data interface{}, options map[string]interface{}) ([]byte, error)
```

### UI Plugin

```go
// GetComponent returns a UI component
func GetComponent(name string, options map[string]interface{}) (interface{}, error)
```

## Example Plugins

Nessi includes several example plugins to help you get started:

- **Email Validator**: A validation plugin that validates email addresses
- **Slack Alerter**: An alert plugin that sends alerts to Slack
- **Prometheus Metrics**: A metric plugin that exports metrics to Prometheus
- **S3 Storage**: A storage plugin that stores data in Amazon S3
- **Excel Export**: An export plugin that exports data to Excel format
- **Chart Component**: A UI plugin that provides custom chart components

You can find these examples in the `pkg/plugin/example` directory.

## Best Practices

When creating plugins for Nessi, follow these best practices:

1. **Provide clear documentation**: Document your plugin's purpose, usage, and options
2. **Handle errors gracefully**: Return meaningful error messages
3. **Validate inputs**: Check that inputs are of the expected type and format
4. **Provide default options**: Set sensible defaults for your plugin's options
5. **Version your plugins**: Use semantic versioning for your plugins
6. **Test thoroughly**: Write tests for your plugin's functionality
7. **Keep it simple**: Focus on a single, well-defined functionality
8. **Be thread-safe**: Ensure your plugin can be used concurrently

## Troubleshooting

If you encounter issues with plugins, try the following:

1. **Check plugin compatibility**: Ensure the plugin is compatible with your version of Nessi
2. **Verify plugin permissions**: Make sure the plugin file has the correct permissions
3. **Check for dependencies**: Ensure all required dependencies are installed
4. **Look for error messages**: Check the error output for clues about the issue
5. **Enable debug logging**: Run Nessi with debug logging enabled to get more information

## Contributing Plugins

If you've created a useful plugin for Nessi, consider contributing it to the community. You can submit your plugin to the Nessi Plugin Registry by following these steps:

1. **Create a GitHub repository** for your plugin
2. **Add documentation** explaining how to use your plugin
3. **Include tests** to verify your plugin's functionality
4. **Submit a pull request** to the Nessi Plugin Registry repository

## Conclusion

The Nessi Plugin System provides a powerful way to extend Nessi's functionality with custom plugins. By creating plugins, you can add new validation rules, alerting mechanisms, metrics collection, storage backends, export formats, and UI components to Nessi.
