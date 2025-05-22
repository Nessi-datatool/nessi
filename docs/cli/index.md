# Nessi CLI Documentation

This document provides an overview of the Nessi Command Line Interface (CLI) and links to detailed documentation for each command group.

## Overview

Nessi CLI is a powerful tool for data quality management, schema validation, and metrics collection. It provides a comprehensive set of commands for working with Delta Lake tables, validating data quality, collecting metrics, and generating reports.

## Installation

```bash
# Install from binary release
curl -sSL https://nessi.dev/install.sh | bash

# Or install from source
git clone https://github.com/nessi-dev/nessi.git
cd nessi
go install ./cmd/nessi
```

## Basic Usage

```bash
# Get help
nessi --help

# Get version information
nessi version

# Get information about your environment
nessi info

# Generate shell completion
nessi completion bash > ~/.bash_completion.d/nessi
```

## Command Groups

Nessi CLI is organized into the following command groups:

| Command Group | Description |
|---------------|-------------|
| [tables](tables.md) | Commands for working with Delta Lake tables |
| [quality](quality.md) | Commands for data quality checks and reporting |
| [scan](scan.md) | Commands for scanning and analyzing data |
| [schema](schema.md) | Commands for schema management and validation |
| [metrics](metrics.md) | Commands for collecting and analyzing metrics |
| [report](report.md) | Commands for generating reports |
| [integration](integration.md) | Commands for external system integration |
| [config](config.md) | Commands for configuration management |

## Core Features (Community Edition)

The following features are available in the Community Edition of Nessi:

- **Delta Lake Table Management**: List, describe, and export Delta Lake tables
- **Basic Schema Operations**: View, validate, and compare schemas
- **Data Quality Checks**: Run basic quality checks against predefined rules
- **Data Scanning**: Scan tables and files to analyze their structure and content
- **Basic Metrics Collection**: Collect basic quality and performance metrics
- **Simple Reports**: Generate basic text and JSON reports

## Premium Features (Pro Edition)

The following features require a Pro license:

- **Advanced Schema Operations**: Schema evolution, registry integration, and inference
- **Advanced Quality Checks**: Custom rules, advanced validation, and quality scoring
- **Comprehensive Metrics**: Trend analysis, alerts, and visualizations
- **Enhanced Reports**: Interactive HTML and PDF reports with visualizations
- **Grafana Integration**: Create and update Grafana dashboards from Nessi metrics
- **Workflow Orchestration**: Schedule and automate Nessi commands
- **Databricks Integration**: Seamless integration with Databricks
- **AWS S3 Storage**: Direct access to data stored in AWS S3

For more information on licensing, see the [License Management documentation](../LICENSE_MANAGEMENT.md).

## Global Options

The following options are available for all Nessi commands:

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help for a command |
| `--verbose`, `-v` | Enable verbose output |
| `--quiet`, `-q` | Suppress all output except errors |
| `--json` | Output in JSON format |
| `--config` | Path to configuration file |
| `--log-level` | Set log level (debug, info, warn, error) |
| `--no-color` | Disable colored output |
| `--license` | Path to license file |

## Environment Variables

Nessi CLI supports the following environment variables:

| Variable | Description |
|----------|-------------|
| `NESSI_CONFIG` | Path to configuration file |
| `NESSI_LOG_LEVEL` | Log level (debug, info, warn, error) |
| `NESSI_LICENSE_PATH` | Path to license file |
| `NESSI_NO_COLOR` | Disable colored output if set to true |
| `NESSI_DEFAULT_FORMAT` | Default output format |
| `NESSI_METRICS_STORE` | Path or URL to metrics store |
| `NESSI_SCHEMA_REGISTRY` | URL to schema registry |
| `NESSI_DATABRICKS_TOKEN` | Databricks access token |
| `NESSI_DATABRICKS_HOST` | Databricks host URL |
| `NESSI_AWS_REGION` | AWS region for S3 access |

## Configuration File

Nessi CLI can be configured using a YAML configuration file. By default, it looks for a file named `.nessi.yaml` in the current directory or the user's home directory.

Example configuration file:

```yaml
# Global configuration
log_level: info
output_format: text
color: true

# Delta Lake configuration
delta:
  default_path: /data/delta
  time_travel_enabled: true

# Quality configuration
quality:
  rules_path: /path/to/rules
  default_threshold: 0.95

# Metrics configuration
metrics:
  store_path: /path/to/metrics
  collect_on_scan: true

# Report configuration
report:
  template_path: /path/to/templates
  default_format: html

# Integration configuration
integration:
  grafana:
    url: http://grafana:3000
    api_key: your-api-key
  databricks:
    host: https://your-workspace.cloud.databricks.com
    token: your-token
```

## Examples

Here are some common usage examples:

```bash
# List all Delta tables
nessi tables list

# Show schema of a table
nessi schema show /data/customers

# Run quality checks
nessi quality check /data/customers

# Scan a table and generate a report
nessi scan run /data/customers --output html --file report.html

# Collect metrics
nessi metrics collect /data/customers --store

# Analyze metrics trends
nessi metrics trend /data/customers --from 2025-04-01 --to 2025-05-22

# Compare two versions of a table
nessi tables diff /data/customers /data/customers --version1 3 --version2 5

# Validate data against a schema
nessi schema validate /data/customers /schemas/customers_schema.json
```

## Error Handling

Nessi CLI uses a consistent error handling approach with error codes and detailed error messages. For more information, see the [Error Handling documentation](../ERROR_HANDLING.md).

## Getting Help

If you need help with Nessi CLI, you can:

- Use the `--help` flag with any command to see usage information
- Visit the [Nessi website](https://nessi.dev) for documentation and tutorials
- Join the [Nessi community forum](https://community.nessi.dev) for support and discussions
- Report issues on [GitHub](https://github.com/nessi-dev/nessi/issues)

## License

Nessi is available in Community (free) and Pro (paid) editions. For more information on licensing, see the [License Management documentation](../LICENSE_MANAGEMENT.md).
