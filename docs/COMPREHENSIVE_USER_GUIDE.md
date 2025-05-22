# Nessi Comprehensive User Guide

## Table of Contents

- [Introduction](#introduction)
- [Core Concepts](#core-concepts)
- [Installation](#installation)
- [Configuration](#configuration)
- [Command Line Interface](#command-line-interface)
- [Data Lake Operations](#data-lake-operations)
- [Schema Management](#schema-management)
- [Data Quality](#data-quality)
- [Metrics and Monitoring](#metrics-and-monitoring)
- [Reporting](#reporting)
- [Integration](#integration)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)
- [Glossary](#glossary)
- [FAQ](#faq)

## Introduction

Nessi is a powerful data quality and observability tool designed for data engineers and analysts working with data lakes. It provides comprehensive capabilities for monitoring, validating, and reporting on data quality across your data ecosystem.

### About This Guide

This comprehensive user guide provides detailed information on using Nessi effectively. Whether you're new to Nessi or an experienced user, this guide will help you understand and utilize all of Nessi's features.

### Editions

Nessi is available in two editions:

#### Core Features (Community Edition)

The Community Edition is open-source and includes:

- Basic Delta Lake table operations
- Schema validation and evolution tracking
- Data quality rule definition and validation
- Basic metrics collection and analysis
- HTML/PDF report generation
- Command-line interface
- Local file system support

#### Premium Features (Pro Edition)

The Pro Edition includes all Community Edition features plus:

- Databricks integration
- AWS S3 storage support
- Workflow orchestration
- Advanced metrics and analytics
- Team collaboration features
- Enhanced support options

For more details on editions and licensing, see [LICENSE_MANAGEMENT.md](LICENSE_MANAGEMENT.md).

## Core Concepts

### Data Lake

A data lake is a centralized repository that allows you to store structured and unstructured data at any scale. Nessi works with data lakes to provide quality monitoring and observability.

### Delta Lake

Delta Lake is an open-source storage layer that brings reliability to data lakes. Nessi has built-in support for Delta Lake tables, allowing you to perform quality checks and schema validation on Delta Lake data.

### Data Quality

Data quality refers to the condition of data based on factors such as accuracy, completeness, consistency, and reliability. Nessi helps you define, measure, and monitor data quality across your data lake.

### Schema

A schema defines the structure of your data, including field names, data types, and constraints. Nessi provides tools for schema validation, evolution tracking, and comparison.

### Metrics

Metrics are quantitative measurements that help you understand the state and quality of your data. Nessi collects various metrics, including quality metrics, performance metrics, and freshness metrics.

### Reports

Reports provide a structured way to view and share information about your data quality and metrics. Nessi generates comprehensive reports in HTML and PDF formats.

## Installation

For detailed installation instructions, please refer to [INSTALLATION.md](INSTALLATION.md).

### Quick Installation

#### Using Installation Script (Linux/macOS)

```bash
curl -sSL https://nessi.dev/install.sh | bash
```

#### Using Pre-built Binaries

1. Download the appropriate binary for your platform from [https://nessi.dev/downloads](https://nessi.dev/downloads)
2. Extract the archive
3. Move the binary to a directory in your PATH

#### Using Docker

```bash
docker pull nessi/nessi:latest
docker run --rm nessi/nessi:latest --help
```

#### Building from Source

```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
go build -o nessi ./cmd/nessi
```

### Verifying Installation

To verify that Nessi is installed correctly:

```bash
nessi --version
```

## Configuration

Nessi can be configured using a configuration file, environment variables, or command-line flags.

### Configuration File

By default, Nessi looks for a configuration file at `~/.nessi/config.yaml`. You can specify a different configuration file using the `--config` flag.

Example configuration file:

```yaml
# General configuration
general:
  log_level: info
  output_format: text

# Storage configuration
storage:
  type: local
  path: /path/to/data

# Delta Lake configuration
delta:
  enable_time_travel: true
  
# Quality configuration
quality:
  rules_path: /path/to/rules
  
# Reporting configuration
reporting:
  output_dir: /path/to/reports
  template: default
```

### Environment Variables

Nessi supports configuration through environment variables. Environment variables take precedence over configuration file settings.

Examples:

```bash
# Set log level
export NESSI_LOG_LEVEL=debug

# Set storage type
export NESSI_STORAGE_TYPE=s3

# Set AWS credentials (for S3 storage)
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
```

### Command-line Flags

Command-line flags take precedence over both environment variables and configuration file settings.

Examples:

```bash
# Set log level
nessi --log-level debug

# Set output format
nessi --output-format json

# Set storage path
nessi --storage-path /path/to/data
```

## Command Line Interface

Nessi provides a comprehensive command-line interface (CLI) for interacting with your data lake and performing various operations.

For detailed CLI documentation, please refer to the [CLI Reference](cli/index.md).

### Basic Usage

```bash
nessi [global options] command [command options] [arguments...]
```

### Global Options

- `--config FILE`: Load configuration from FILE
- `--log-level LEVEL`: Set log level (debug, info, warn, error)
- `--output-format FORMAT`: Set output format (text, json, csv)
- `--help, -h`: Show help
- `--version, -v`: Print version information

### Commands

Nessi provides several command categories:

- `tables`: Delta Lake table operations
- `quality`: Data quality checking and reporting
- `scan`: Data scanning and analysis
- `schema`: Schema management and validation
- `metrics`: Metrics collection and analysis
- `report`: Report generation and management
- `integration`: External system integration
- `config`: Configuration management

For detailed information on each command, see the [CLI Reference](cli/REFERENCE.md).

## Data Lake Operations

Nessi provides various commands for working with data lakes, particularly Delta Lake tables.

### Listing Tables

```bash
nessi tables list --path /path/to/data
```

### Viewing Table Details

```bash
nessi tables describe --path /path/to/table
```

### Viewing Table Schema

```bash
nessi tables schema --path /path/to/table
```

### Viewing Table History

```bash
nessi tables history --path /path/to/table
```

### Time Travel

Nessi supports time travel for Delta Lake tables, allowing you to view the table as it existed at a specific version or timestamp.

```bash
# View table at a specific version
nessi tables show --path /path/to/table --version 3

# View table at a specific timestamp
nessi tables show --path /path/to/table --timestamp "2023-01-01T12:00:00Z"
```

## Schema Management

Nessi provides tools for managing and validating schemas.

### Validating Schema

```bash
nessi schema validate --path /path/to/table --schema-file /path/to/schema.json
```

### Comparing Schemas

```bash
nessi schema compare --source /path/to/table1 --target /path/to/table2
```

### Tracking Schema Evolution

```bash
nessi schema history --path /path/to/table
```

### Visualizing Schema

```bash
nessi schema tree --path /path/to/table
```

## Data Quality

Nessi provides comprehensive data quality validation capabilities.

### Defining Quality Rules

Quality rules can be defined in YAML format:

```yaml
# quality_rules.yaml
rules:
  - name: not_null
    description: Check if column values are not null
    columns:
      - id
      - name
    
  - name: unique
    description: Check if column values are unique
    columns:
      - id
    
  - name: range
    description: Check if column values are within range
    columns:
      - age
    parameters:
      min: 0
      max: 120
```

### Running Quality Checks

```bash
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml
```

### Viewing Quality Results

```bash
nessi quality results --path /path/to/table
```

### Quality Metrics

Nessi calculates various quality metrics, including:

- Completeness: Percentage of non-null values
- Accuracy: Percentage of values meeting accuracy rules
- Consistency: Percentage of values meeting consistency rules
- Uniqueness: Percentage of unique values
- Timeliness: Freshness of data

To view quality metrics:

```bash
nessi quality metrics --path /path/to/table
```

## Metrics and Monitoring

Nessi collects and analyzes various metrics to help you monitor your data lake.

### Collecting Metrics

```bash
nessi metrics collect --path /path/to/table
```

### Viewing Metrics

```bash
nessi metrics show --path /path/to/table
```

### Metrics Types

Nessi collects several types of metrics:

- **Quality Metrics**: Measurements of data quality
- **Performance Metrics**: Measurements of query and processing performance
- **Freshness Metrics**: Measurements of data freshness and timeliness
- **Volume Metrics**: Measurements of data volume and growth

### Visualizing Metrics

```bash
nessi metrics visualize --path /path/to/table --type quality
```

## Reporting

Nessi generates comprehensive reports on data quality and metrics.

### Generating Reports

```bash
nessi report generate --path /path/to/table --output /path/to/report.html
```

### Report Types

Nessi supports several report types:

- **Quality Report**: Detailed report on data quality
- **Schema Report**: Report on schema structure and evolution
- **Freshness Report**: Report on data freshness
- **Performance Report**: Report on performance metrics

To specify a report type:

```bash
nessi report generate --path /path/to/table --type quality --output /path/to/report.html
```

### Report Formats

Nessi supports multiple report formats:

- **HTML**: Interactive HTML reports with visualizations
- **PDF**: Printable PDF reports
- **JSON**: Machine-readable JSON reports
- **CSV**: Tabular CSV reports

To specify a report format:

```bash
nessi report generate --path /path/to/table --format pdf --output /path/to/report.pdf
```

## Integration

Nessi provides integration with various external systems.

### Databricks Integration (Pro Edition)

```bash
# Set Databricks credentials
export DATABRICKS_HOST=your_databricks_host
export DATABRICKS_TOKEN=your_databricks_token

# List Databricks tables
nessi integration databricks list-tables --catalog your_catalog --schema your_schema

# Run quality checks on Databricks table
nessi quality check --integration databricks --catalog your_catalog --schema your_schema --table your_table --rules /path/to/quality_rules.yaml
```

### AWS S3 Integration (Pro Edition)

```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key

# List S3 tables
nessi integration s3 list-tables --bucket your_bucket --prefix your_prefix

# Run quality checks on S3 table
nessi quality check --integration s3 --bucket your_bucket --key your_key --rules /path/to/quality_rules.yaml
```

## Troubleshooting

### Common Issues

#### Installation Issues

- **Issue**: `nessi: command not found`
  - **Solution**: Ensure the Nessi binary is in your PATH

- **Issue**: Permission denied when running Nessi
  - **Solution**: Ensure the Nessi binary has execute permissions (`chmod +x nessi`)

#### Configuration Issues

- **Issue**: Configuration file not found
  - **Solution**: Ensure the configuration file exists at the expected location or specify the path using `--config`

- **Issue**: Invalid configuration
  - **Solution**: Check the configuration file syntax and ensure all required fields are present

#### Delta Lake Issues

- **Issue**: Unable to read Delta Lake table
  - **Solution**: Ensure the path points to a valid Delta Lake table (contains `_delta_log` directory)

- **Issue**: Schema validation errors
  - **Solution**: Check the table schema and ensure it matches the expected schema

#### Performance Issues

- **Issue**: Slow performance with large tables
  - **Solution**: Consider using sampling (`--sample-size`) for large tables

### Error Codes

Nessi uses standardized error codes to help diagnose issues:

- **N1XX**: Input/output errors
- **N2XX**: Configuration errors
- **N3XX**: Delta Lake errors
- **N4XX**: Schema errors
- **N5XX**: Quality errors
- **N6XX**: Metrics errors
- **N7XX**: Report errors
- **N8XX**: Integration errors
- **N9XX**: Internal errors

For detailed error information, use the `--log-level debug` flag.

### Logs

Nessi logs information to help diagnose issues. By default, logs are written to stderr. You can control the log level using the `--log-level` flag.

To save logs to a file:

```bash
nessi --log-level debug command > nessi.log 2>&1
```

### Getting Help

If you encounter issues that you cannot resolve, please:

1. Check the [FAQ](#faq) section
2. Search for similar issues on the [GitHub Issues](https://github.com/nessi-dev/nessi/issues) page
3. Open a new issue if your problem is not already reported

## Best Practices

### Data Quality Management

- Define clear quality rules for each table
- Run quality checks regularly, ideally as part of your data pipeline
- Monitor quality metrics over time to identify trends
- Set up alerts for significant quality degradation

### Schema Management

- Document your schemas thoroughly
- Track schema changes over time
- Validate schema changes before applying them
- Use schema comparison to ensure consistency across environments

### Performance Optimization

- Use sampling for large tables when appropriate
- Schedule resource-intensive operations during off-peak hours
- Monitor performance metrics to identify bottlenecks
- Consider partitioning large tables for better performance

### Reporting

- Generate regular reports for key stakeholders
- Customize reports to focus on relevant metrics
- Use HTML reports for interactive exploration
- Use PDF reports for sharing and archiving

### Security

- Secure your configuration file, especially if it contains sensitive information
- Use environment variables for sensitive configuration when possible
- Follow the principle of least privilege when setting up integrations
- Regularly update Nessi to get security fixes

## Glossary

- **Delta Lake**: An open-source storage layer that brings reliability to data lakes
- **Data Quality**: The condition of data based on factors such as accuracy, completeness, consistency, and reliability
- **Schema**: The structure of data, including field names, data types, and constraints
- **Metrics**: Quantitative measurements that help understand the state and quality of data
- **Time Travel**: The ability to access previous versions of data
- **Quality Rule**: A definition of what constitutes quality for a specific aspect of data
- **Report**: A structured presentation of data quality and metrics information

## FAQ

### General Questions

**Q: Is Nessi open source?**

A: Yes, the Community Edition of Nessi is open source under the Apache 2.0 license. The Pro Edition includes additional features and is available under a commercial license.

**Q: What data formats does Nessi support?**

A: Nessi primarily supports Delta Lake tables. Support for other formats may be available in future releases.

**Q: Can Nessi work with my existing data lake?**

A: Yes, Nessi can work with existing Delta Lake tables without any modifications to the data.

### Installation Questions

**Q: Does Nessi require any dependencies?**

A: Nessi is distributed as a single binary with no external dependencies for the core functionality. Some integrations may require additional configuration.

**Q: Is there a Docker image available?**

A: Yes, Nessi provides official Docker images that you can use to run Nessi in a containerized environment.

### Usage Questions

**Q: How do I define custom quality rules?**

A: Custom quality rules can be defined in YAML format. See the [Data Quality](#data-quality) section for details.

**Q: Can Nessi automatically fix quality issues?**

A: Currently, Nessi focuses on identifying and reporting quality issues. Automatic remediation is not yet supported.

**Q: How do I integrate Nessi with my data pipeline?**

A: Nessi can be integrated into data pipelines using its command-line interface. You can call Nessi commands from your pipeline scripts or orchestration tools.

### Technical Questions

**Q: How does Nessi handle large tables?**

A: Nessi can use sampling to efficiently process large tables. You can control the sample size using the `--sample-size` flag.

**Q: Does Nessi support incremental processing?**

A: Yes, Nessi can process only the new data since the last run when working with Delta Lake tables that have a timestamp column.

**Q: How secure is Nessi?**

A: Nessi is designed with security in mind. It does not require special privileges to run and can be configured to use secure connections for integrations.

### Pro Edition Questions

**Q: How do I upgrade from Community Edition to Pro Edition?**

A: You can upgrade by obtaining a Pro Edition license and activating it using the `nessi config license activate` command.

**Q: What additional features are included in the Pro Edition?**

A: The Pro Edition includes Databricks integration, AWS S3 support, workflow orchestration, and other premium features. See [LICENSE_MANAGEMENT.md](LICENSE_MANAGEMENT.md) for details.

**Q: Is there a trial available for the Pro Edition?**

A: Yes, you can start a 1-month free trial of the Pro Edition using the `nessi config license trial` command.

---

For more information, visit [nessi.dev](https://nessi.dev) or contact support@nessi.dev.
