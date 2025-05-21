# Nessi Plugin System

This document provides an overview of Nessi's plugin system, including installation, usage, and development guidelines.

## Overview

Nessi's plugin system allows users to extend its functionality with custom features, integrations, and enhancements. Plugins can add new commands, data quality rules, report extensions, format handlers, and more.

## Installation

### Installing Plugins

Plugins can be installed using the `nessi plugins install` command:

```bash
# Install a plugin from a local path
nessi plugins install /path/to/plugin

# Install a plugin from a GitHub repository
nessi plugins install github.com/username/nessi-plugin
```

### Listing Installed Plugins

To see all installed plugins:

```bash
nessi plugins list
```

### Removing Plugins

To remove a plugin:

```bash
nessi plugins remove plugin_name
```

## Using Plugins

Plugins can be used in various ways depending on their type:

### Using with Reports

Many plugins enhance Nessi's reporting capabilities. To use a plugin with a report:

```bash
nessi report generate --table my_table --plugin plugin_name
```

You can also pass arguments to the plugin:

```bash
nessi report generate --table my_table --plugin plugin_name --plugin-args "key1=value1,key2=value2"
```

### Using Plugin Commands

Some plugins add new commands to Nessi. These commands can be accessed directly:

```bash
nessi plugin_command [args]
```

## Viral Growth Plugins

Nessi includes several plugins designed to increase its visibility and adoption through viral growth mechanisms:

### Social Sharing Plugin

Enhances reports with advanced social sharing capabilities:

```bash
# Generate a report with social sharing features
nessi viral share my_table --title "My Quality Report" --hashtags "dataquality,datalake"
```

### Badge Plugin

Generates embeddable "Powered by Nessi" badges for projects:

```bash
# Generate a badge in markdown format
nessi viral badge --format markdown

# Generate a badge with a quality score
nessi viral badge --quality-score 95
```

### Community Engagement Plugin

Facilitates community contributions and engagement:

```bash
# Get contribution suggestions
nessi viral community contribute --experience beginner

# Submit feedback
nessi viral community feedback --type feature --text "It would be great to have..."
```

## Built-in Plugins

Nessi comes with several built-in plugins:

### Quality Rules Plugin

Adds custom data validation rules:

```bash
# Run quality checks with custom rules
nessi quality check --table my_table --plugin quality_rules_plugin
```

### Report Extensions Plugin

Extends reporting capabilities with custom sections:

```bash
# Generate a report with extended sections
nessi report generate --table my_table --plugin report_extensions_plugin
```

### Format Handler Plugin

Adds support for additional data formats:

```bash
# Scan a JSONL file
nessi scan --file data.jsonl --plugin format_handler_plugin
```

## Plugin Development

To develop your own plugins, see the [Plugin API Documentation](PLUGIN_API.md) for detailed information on interfaces, methods, and examples.

### Plugin Types

Nessi supports several types of plugins:

1. **Command Executor**: Adds new commands to the Nessi CLI
2. **Quality Rule Provider**: Adds custom data quality validation rules
3. **Report Extension**: Extends report templates with custom sections
4. **Format Handler**: Adds support for additional data formats
5. **Integration**: Connects Nessi with external systems

### Example Plugin Structure

```
plugin_name/
├── main.go         # Main plugin implementation
├── README.md       # Plugin documentation
└── [other files]   # Additional plugin files
```

## Community Contributions

We encourage the community to develop and share plugins that extend Nessi's functionality. Some ideas for plugins:

- Integration with additional data sources
- Custom visualization formats
- Specialized quality rules for specific domains
- CI/CD integrations
- Additional report formats

To share your plugin with the community, create a GitHub repository and announce it in the Nessi community channels.

## Troubleshooting

### Common Issues

- **Plugin not found**: Ensure the plugin is installed correctly and the path is valid
- **Incompatible plugin**: Check that the plugin is compatible with your version of Nessi
- **Plugin errors**: Check the plugin logs for detailed error messages

### Getting Help

If you encounter issues with plugins, you can:

- Check the plugin documentation
- Open an issue on the plugin's GitHub repository
- Ask for help in the Nessi community channels
