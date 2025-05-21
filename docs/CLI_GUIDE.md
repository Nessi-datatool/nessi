# Nessi CLI Guide

## Overview

Nessi is a powerful CLI-only tool for data quality management and Delta Lake operations. This guide provides detailed information about using the Nessi command-line interface for all your data quality needs.

## Command Structure

The Nessi CLI follows a consistent command structure:

```
nessi <command> [subcommand] [options] [arguments]
```

Where:
- `<command>` is the main command category (e.g., `report`, `scan`, `validate`)
- `[subcommand]` is an optional subcommand for more specific operations
- `[options]` are flags that modify the command behavior
- `[arguments]` are the required inputs for the command

## Core Commands

### Report Generation

Generate comprehensive reports in various formats (HTML, PDF, JSON, CSV):

```bash
nessi report <type> <data_path> [options]
```

Where `<type>` can be:
- `quality` - Generate a data quality report
- `schema` - Generate a schema validation report
- `freshness` - Generate a data freshness report
- `performance` - Generate a performance optimization report

Options:
- `--format, -f` - Report format (html, pdf, json, csv). Default: html
- `--output, -o` - Output path for the report
- `--templates, -t` - Directory containing report templates

Examples:

```bash
# Generate an HTML quality report
nessi report quality s3://my-bucket/my-table --output=quality-report.html

# Generate a PDF schema report
nessi report schema s3://my-bucket/my-table --format=pdf --output=schema-report.pdf

# Generate a JSON freshness report
nessi report freshness s3://my-bucket/my-table --format=json --output=freshness-report.json
```

### Data Scanning

Scan Delta Lake tables for quality issues:

```bash
nessi scan <table_path> [options]
```

Options:
- `--rules, -r` - Path to quality rules YAML file
- `--output, -o` - Output path for scan results
- `--format, -f` - Output format (json, yaml). Default: json
- `--verbose, -v` - Enable verbose output

Example:

```bash
# Scan a table with custom rules
nessi scan s3://my-bucket/my-table --rules=rules.yaml --output=scan-results.json
```

### Schema Operations

View and validate table schemas:

```bash
nessi schema <table_path> [options]
```

Options:
- `--history` - Show schema evolution history
- `--version, -v` - View schema at a specific version
- `--output, -o` - Output path for schema information
- `--format, -f` - Output format (json, yaml). Default: json

Examples:

```bash
# View current schema
nessi schema s3://my-bucket/my-table

# View schema history
nessi schema s3://my-bucket/my-table --history

# View schema at version 3
nessi schema s3://my-bucket/my-table --version=3
```

### Time Travel

Access historical versions of Delta Lake tables:

```bash
nessi timetravel <table_path> [options]
```

Options:
- `--version, -v` - Table version to travel to
- `--timestamp, -t` - Timestamp to travel to (format: YYYY-MM-DDTHH:MM:SS)
- `--output, -o` - Output path for time travel results

Examples:

```bash
# Time travel to version 5
nessi timetravel s3://my-bucket/my-table --version=5

# Time travel to a specific timestamp
nessi timetravel s3://my-bucket/my-table --timestamp="2023-01-15T14:30:00"
```

## Report Types

Nessi generates comprehensive reports in multiple formats:

### HTML Reports

HTML reports provide interactive visualizations and detailed information about your data:

- **Quality Reports**: Data quality metrics, issue summaries, and recommendations
- **Schema Reports**: Schema structure, field details, and validation results
- **Freshness Reports**: Data freshness metrics, update frequency, and staleness alerts
- **Performance Reports**: Optimization recommendations, file statistics, and partition analysis

### PDF Reports

PDF reports offer the same information as HTML reports in a printable format, perfect for sharing with stakeholders or archiving.

### JSON/CSV Reports

Machine-readable formats for integration with other tools and systems.

## Configuration

Nessi can be configured using a YAML configuration file or environment variables.

### Configuration File

Create a `nessi.yaml` file in your working directory or specify a path with `--config`:

```yaml
# nessi.yaml example
security:
  credentials:
    enabled: true
    file: "credentials.yaml"

reporting:
  enabled: true
  output_path: "./reports"
  default_format: "html"
  templates_path: "./templates"

monitoring:
  metrics_enabled: true
  metrics_path: "./metrics"
  metrics_retention: 30d

quality:
  default_rules: "rules.yaml"
  threshold: 0.95
```

### Environment Variables

All configuration options can also be set using environment variables:

```bash
export NESSI_REPORTING_ENABLED=true
export NESSI_REPORTING_OUTPUT_PATH="./reports"
export NESSI_REPORTING_DEFAULT_FORMAT="html"
```

## Databricks Integration

Nessi integrates with Databricks for seamless Delta Lake operations:

```bash
# Set Databricks credentials
export DATABRICKS_HOST="https://your-workspace.cloud.databricks.com"
export DATABRICKS_TOKEN="your-personal-access-token"

# List catalogs in Databricks
nessi catalog list --type databricks

# List tables in a specific catalog and schema
nessi catalog tables --type databricks --database main.default
```

## Batch Processing

Nessi supports batch operations for processing multiple tables:

```bash
# Scan multiple tables
nessi batch scan --tables-file=tables.txt --rules=rules.yaml --output-dir=scan-results

# Generate reports for multiple tables
nessi batch report quality --tables-file=tables.txt --format=html --output-dir=reports
```

## Logging

Control logging verbosity with the `--log-level` flag:

```bash
nessi scan s3://my-bucket/my-table --log-level=debug
```

Available log levels: error, warn, info, debug, trace

## Exit Codes

Nessi returns the following exit codes:

- `0`: Success
- `1`: General error
- `2`: Configuration error
- `3`: Input validation error
- `4`: Data quality error (when quality checks fail)
- `5`: Permission error

## Further Resources

- [Report Templates Guide](REPORT_TEMPLATES.md)
- [Quality Rules Reference](QUALITY_RULES.md)
- [Delta Lake Operations](DELTA_OPERATIONS.md)
- [Databricks Integration](DATABRICKS.md)
- [Plugin System Commands](#plugin-system-commands)
- [Viral Growth Commands](#viral-growth-commands)

## Plugin System Commands

Manage and use Nessi plugins:

```bash
nessi plugins [subcommand] [options]
```

Subcommands:
- `list` - List installed plugins
- `install` - Install a plugin
- `remove` - Remove a plugin
- `info` - Display information about a plugin

Options:
- `--path, -p` - Path to plugin (for install)
- `--name, -n` - Name of plugin (for remove, info)

Examples:

```bash
# List installed plugins
nessi plugins list

# Install a plugin from a local path
nessi plugins install --path=/path/to/plugin

# Install a plugin from a GitHub repository
nessi plugins install github.com/username/nessi-plugin

# Remove a plugin
nessi plugins remove --name=plugin_name

# Get information about a plugin
nessi plugins info --name=plugin_name
```

## Viral Growth Commands

Enhance Nessi's visibility and foster community engagement:

```bash
nessi viral [subcommand] [options]
```

Subcommands:
- `share` - Generate shareable reports
- `badge` - Create 'Powered by Nessi' badges
- `community` - Engage with the Nessi community

### Share Command

Generate shareable reports with social sharing features:

```bash
nessi viral share [table_name] [options]
```

Options:
- `--output, -o` - Output path for the shareable report
- `--title, -t` - Title for the shareable report
- `--description, -d` - Description for the shareable report
- `--hashtags` - Comma-separated hashtags for social media sharing

Examples:

```bash
# Generate a shareable report
nessi viral share my_table --output=report.html

# Generate a report with custom title and description
nessi viral share my_table --title="Quality Report" --description="Monthly quality check"

# Specify hashtags for social media sharing
nessi viral share my_table --hashtags="dataquality,datalake,opensource"
```

### Badge Command

Generate 'Powered by Nessi' badges for embedding in projects:

```bash
nessi viral badge [options]
```

Options:
- `--format, -f` - Badge format (markdown, html, rst)
- `--label, -l` - Text on the left side of the badge
- `--message, -m` - Text on the right side of the badge
- `--color, -c` - Color of the right side (hex or named color)
- `--style, -s` - Badge style (flat, flat-square, plastic, etc.)
- `--quality-score, -q` - Quality score to include in the badge (0-100)

Examples:

```bash
# Generate a badge in markdown format
nessi viral badge --format=markdown

# Generate a badge with custom text
nessi viral badge --label="verified by" --message="nessi" --color=blue

# Generate a quality score badge
nessi viral badge --quality-score=95
```

### Community Command

Engage with the Nessi community:

```bash
nessi viral community [subcommand] [options]
```

Subcommands:
- `feedback` - Provide feedback on Nessi
- `contribute` - Get suggestions for contributing to Nessi

Feedback options:
- `--type, -t` - Type of feedback (general, bug, feature)
- `--text` - Feedback text
- `--name, -n` - Your name (optional)
- `--email, -e` - Your email (optional)

Contribute options:
- `--experience, -e` - Your experience level (beginner, intermediate, advanced)
- `--github-profile, -g` - Your GitHub profile (optional)

Examples:

```bash
# Submit feedback
nessi viral community feedback --type=feature --text="It would be great to have..."

# Get contribution suggestions for beginners
nessi viral community contribute --experience=beginner
```
