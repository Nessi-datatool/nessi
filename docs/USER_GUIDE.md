# Nessi User Guide

## Navigation

- [Documentation Home](README.md)
- [Installation & Quickstart](../QUICKSTART.md)
- [Configuration](CONFIGURATION.md)
- [CLI Reference](cli/README.md)
- [Reporting](REPORTING.md)
- [Quality Rules](QUALITY_RULES.md)
- [Integrations](integration_guide.md)
- [Developer Guide](developer_experience.md)
- [FAQ](faq.md)

---

## Features

- Data Quality Checks
- System Metrics
- CLI Interface
- Plugin System
- Viral Growth Features
  - Social Sharing
  - Community Engagement
  - 'Powered by Nessi' Badges
- Python API
- Reporting System
- Alerting System
- Metrics Export
- Profiling
- Freshness Monitoring
- Workflow Integrations
- Email Alerting
- Comprehensive Reporting

---

## Table of Contents

1. Introduction
2. Installation & Quickstart
3. Core Concepts
4. Security Features
5. Command-Line Interface
6. Python API
7. Data Quality Features
8. Reporting System
9. Monitoring & Metrics
10. Integrations
11. Configuration
12. Best Practices
13. Troubleshooting
14. FAQ

---

## Introduction

Nessi is a comprehensive monitoring system designed for data quality and system monitoring. It provides tools for profiling data, validating data against rules, and monitoring system metrics.

### Key Features

- **Delta Lake Support**: Schema evolution tracking, version control, and time travel
- **Data Quality**: Profiling, validation rules, and quality scoring
- **Monitoring**: Real-time metrics, comprehensive reporting, and alerting
- **Integration**: dbt, Airflow, and other workflow tools

## Installation & Quickstart

### Prerequisites

- Python 3.8+
- Go 1.18+
- Docker (optional)

### Installation

```bash
# Install from source
git clone https://github.com/nessi-dev/nessi.git
cd nessi
make install

# Or use Docker
docker pull nessi-dev/nessi:latest
```

For detailed installation instructions, see the [QUICKSTART.md](../QUICKSTART.md) guide.

## Core Concepts

### Data Quality

Nessi's data quality features include:

- **Profiling**: Statistical analysis of data
- **Validation**: Rule-based checks against data
- **Monitoring**: Tracking quality metrics over time
- **Alerting**: Notifications for quality issues

### Monitoring

Nessi provides monitoring for:

- **System Metrics**: CPU, memory, disk usage
- **Data Metrics**: Volume, freshness, quality
- **Process Metrics**: Job duration, success rates

## Security Features

Nessi includes several security features:

- File-level access controls
- Encrypted configuration storage
- Secure credential management
- Audit logging

## Reporting System

Nessi provides a comprehensive reporting system with multiple output formats:

- **HTML Reports**: Static reports with quality score cards and visualizations for viewing in browsers
- **PDF Reports**: Portable reports for sharing and documentation
- **JSON/CSV**: Machine-readable formats for integration with other tools

Generate reports using the CLI:

```bash
# Generate a quality report for a table
nessi report --table path/to/table --format html --output report.html

# Generate a freshness report
nessi freshness report --format pdf --output freshness.pdf

# Export metrics as JSON for integration
nessi metrics export --format json --output metrics.json
```

Reports include detailed information about:

- Data quality scores and trends
- Schema evolution and changes
- Freshness monitoring results
- Validation rule results
- System performance metrics

## Command-Line Interface

Nessi provides a comprehensive CLI:

```bash
# Get help
nessi --help

# Run data quality checks
nessi validate --table path/to/table

# Generate a profile
nessi profile --table path/to/table

# View schema as ASCII tree
nessi schema-tree --table path/to/table
```

### Schema Visualization

The `schema-tree` command displays the schema of a Delta Lake table as an ASCII tree, making it easy to understand the structure of your data:

```bash
# View table schema as a tree
nessi schema-tree --table path/to/table
```

Example output:

```
├── id: integer (not null)
├── name: string
├── email: string
├── signup_date: timestamp
└── profile: struct
    ├── age: integer
    ├── country: string
    └── preferences: array
        └── string
```

### Data Visualization

Nessi provides comprehensive visualization commands to help you understand your data quality metrics:

```bash
# Generate a heatmap visualization of data quality metrics
nessi visualize heatmap --table path/to/table --metric completeness

# Generate a bar chart visualization of data quality metrics
nessi visualize barchart --table path/to/table

# Show all available metrics in the bar chart
nessi visualize barchart --table path/to/table --all-metrics

# Generate a scatter plot comparing two metrics
nessi visualize scatter --table path/to/table --x-metric completeness --y-metric accuracy

# Other visualization options
nessi visualize --help
```

#### Visualization Types

1. **Heatmap**: Displays data quality metrics across different dimensions (e.g., time periods and data partitions) using color coding to indicate quality levels.

2. **Bar Chart**: Shows comparative metrics for different quality dimensions in an easy-to-read bar format, with color coding to highlight good, medium, and poor quality areas.

3. **Scatter Plot**: Allows comparison of two different metrics to identify correlations and patterns in your data quality measurements.

## Python API

Nessi can be used programmatically:

```python
from nessi import Validator, Profiler

# Create a validator
validator = Validator()
validator.add_rule("completeness", "column_name", threshold=0.95)
results = validator.validate("path/to/table")

# Create a profiler
profiler = Profiler()
profile = profiler.profile("path/to/table")
```

## Data Quality Features

### Profiling

```bash
nessi profile --table path/to/table --output profile.json
```

### Validation

```bash
nessi validate --table path/to/table --rules rules.yaml
```

### Custom Rules

Create custom rules in YAML:

```yaml
rules:
  - name: completeness_check
    type: completeness
    column: user_id
    threshold: 0.99
```

## Monitoring

### Metrics Collection

Nessi collects metrics for:

- Table health
- Data quality
- System resources
- Process performance

### Alerting

Configure alerts via the CLI:

```bash
# Create an alert configuration
nessi alerts create --name data_quality_alert --condition "quality_score < 0.9" --output email --recipients team@example.com --severity critical

# List configured alerts
nessi alerts list

# Test an alert configuration
nessi alerts test --name data_quality_alert
```

## Integrations

### dbt Integration

```bash
# Install dbt plugin
pip install nessi-dbt

# Run validation on dbt models
dbt run --nessi-validate
```

### Airflow Integration

```python
# Import Nessi Airflow operators
from nessi_airflow.operators import NessiQualityOperator

# Add a Nessi quality check to your DAG
quality_check = NessiQualityOperator(
    task_id='quality_check',
    table='my_table',
    rules='rules.yaml',
    dag=dag
)
```

## Plugin System

Nessi provides a flexible plugin system that allows you to extend its functionality. Plugins can add new commands, quality rules, report extensions, format handlers, and more.

```bash
# List available plugins
nessi plugins list

# Install a plugin
nessi plugins install plugin_name

# Use a plugin
nessi report generate --table my_table --plugin report_extensions_plugin
```

For more information on developing plugins, see the [Plugin API Documentation](PLUGIN_API.md).

## Viral Growth Features

Nessi includes several features designed to increase its visibility, drive adoption, and foster community engagement.

### Social Sharing

The social sharing plugin enhances Nessi reports with advanced sharing capabilities:

```bash
# Generate a shareable report
nessi viral share my_table --output report.html
```

Shareable reports include:

- Social media sharing buttons for Twitter, LinkedIn, and more
- QR codes for easy access from mobile devices
- Embed codes for including reports in websites and documentation
- Pre-formatted messages optimized for different platforms

### 'Powered by Nessi' Badges

The badge plugin generates embeddable badges that can be added to your project's README, documentation, or website:

```bash
# Generate a badge in markdown format
nessi viral badge --format markdown

# Generate a badge with custom text
nessi viral badge --label "verified by" --message "nessi" --color blue

# Generate a quality score badge
nessi viral badge --quality-score 95
```

Badges can be generated in multiple formats (Markdown, HTML, RST) and customized with different styles, colors, and text.

### Community Engagement

The community engagement plugin facilitates feedback collection and contribution suggestions:

```bash
# Submit feedback
nessi viral community feedback --type feature --text "It would be great to have..."

# Get contribution suggestions
nessi viral community contribute --experience beginner
```

These features help foster community contributions and increase Nessi's adoption.

For more information on leveraging these features, see the [Viral Growth Guide](VIRAL_GROWTH.md).

```python
from airflow import DAG
from nessi.airflow import NessiValidateOperator

validate_task = NessiValidateOperator(
    task_id="validate_data",
    table_path="path/to/table",
    rules="rules.yaml",
    dag=dag
)
```

## Configuration

Nessi can be configured via:

- Config file (`~/.nessi/config.yaml`)
- Environment variables
- CLI flags

Example config:

```yaml
monitoring:
  retention_days: 30
  metrics_interval: 60

security:
  jwt_secret: "your-secret-key"
  enable_tls: true
```

## Best Practices

- Run validation after each data load
- Set up regular profiling jobs
- Configure alerts for critical metrics
- Use version control for rule definitions

## Error Handling & Troubleshooting

### Error Handling System

Nessi includes a robust error handling system that provides clear, actionable error messages and automatic recovery from transient failures. Key features include:

- **Standardized Error Codes**: All errors have unique codes (e.g., N101, N201) that help identify the issue
- **Interactive Resolution**: Many common errors can be resolved interactively
- **Helpful Suggestions**: Error messages include suggestions for resolving the issue
- **Error Telemetry**: Optional collection of error statistics to improve the product
- **Automatic Retry**: Transient errors are automatically retried with exponential backoff

For detailed information, see the [ERROR_HANDLING.md](ERROR_HANDLING.md) documentation.

### Testing Error Handling

You can test the error handling system using the `test-error` command:

```bash
# List all available error codes
nessi test-error --list

# Generate a specific error
nessi test-error N101

# Test with custom details and suggestion
nessi test-error N201 --details "Custom details" --suggestion "Try this solution"
```

### Error Telemetry

Nessi includes an error telemetry system that anonymously collects error statistics to help improve the product:

```bash
# View telemetry status
nessi telemetry status

# Enable error telemetry
nessi telemetry enable-error

# Disable error telemetry
nessi telemetry disable-error

# View error statistics
nessi telemetry report-errors
```

### Common Issues

- **Connection errors**: Check network and credentials
- **Validation failures**: Examine data quality issues
- **Performance problems**: Review resource allocation
- **Path errors**: Verify that file paths exist and are accessible
- **Delta table errors**: Ensure tables are valid Delta Lake tables

### Logs

Logs are stored in `~/.nessi/logs/` by default.

## FAQ

### How do I contribute to Nessi?

See our [CONTRIBUTING.md](../CONTRIBUTING.md) guide.

### How do I report bugs?

Open an issue on our [GitHub repository](https://github.com/nessi-dev/nessi/issues).

### How do I get help?

Join our community on [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions).
