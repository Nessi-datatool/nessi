# Nessi Quick Start Guide

This guide will help you get started with Nessi quickly, covering installation, basic configuration, and your first commands.

## Installation

### Option 1: Using the Installation Script (Recommended)

The easiest way to install Nessi is using our installation script:

```bash
curl -sSL https://nessi.dev/install.sh | bash
```

### Option 2: Using Pre-built Binaries

1. Download the appropriate binary for your platform from the [releases page](https://github.com/nessi-dev/nessi/releases).
2. Extract and move the binary to a directory in your PATH.

For detailed installation instructions, see the [Installation Guide](INSTALLATION.md).

## Verify Installation

Verify that Nessi is installed correctly:

```bash
nessi version
```

You should see output showing the version, build date, and other information.

## Basic Configuration

Initialize your configuration:

```bash
# Initialize with interactive prompts
nessi config init --interactive

# Or use the default configuration
nessi config init
```

## Your First Commands

### 1. List Available Tables

```bash
# List all Delta Lake tables
nessi tables list
```

### 2. View Table Schema

```bash
# Show schema of a table
nessi schema show /path/to/table
```

Example output:
```
Schema: /path/to/table
Format: delta
Version: 5 (latest)

COLUMN       TYPE        NULLABLE  DESCRIPTION
id           string      false     Unique identifier
name         string      false     Full name
email        string      true      Email address
age          integer     true      Age in years
signup_date  date        true      Date when signed up
last_login   timestamp   true      Timestamp of last login
```

### 3. Run Quality Checks

```bash
# Run quality checks on a table
nessi quality check /path/to/table
```

### 4. Generate a Report

```bash
# Generate a quality report
nessi report generate /path/to/table --type quality --output html --file quality_report.html
```

## Core Workflows

### Data Quality Workflow

```bash
# 1. List tables
nessi tables list

# 2. Check table schema
nessi schema show /path/to/table

# 3. Run quality checks
nessi quality check /path/to/table

# 4. Generate a quality report
nessi report generate /path/to/table --type quality --output html --file quality_report.html
```

### Data Scanning Workflow

```bash
# 1. Run a comprehensive scan
nessi scan run /path/to/table --type all

# 2. View scan results
nessi scan export /results/latest_scan.json /reports/scan_report.html --format html

# 3. Compare with previous scan
nessi scan compare /results/latest_scan.json /results/previous_scan.json
```

### Metrics Collection Workflow

```bash
# 1. Collect metrics
nessi metrics collect /path/to/table --store

# 2. View metrics
nessi metrics show latest

# 3. Analyze trends
nessi metrics trend /path/to/table --from 2025-04-01
```

## Using the Enhanced Report Templates

Nessi includes enhanced report templates with modern styling and interactive elements:

```bash
# Generate an enhanced quality report
nessi report generate /path/to/table --type quality --template enhanced --output html --file enhanced_report.html
```

## Command Line Help

Get help for any command:

```bash
# General help
nessi --help

# Command-specific help
nessi tables --help
nessi quality check --help
```

## Next Steps

After getting familiar with the basics, explore these features:

1. **Schema Management**: Validate and evolve schemas
   ```bash
   nessi schema validate /path/to/table /path/to/schema.json
   nessi schema evolve /path/to/table --add-column "new_column:string:true:Description"
   ```

2. **Metrics Visualization**: Create visual representations of your metrics
   ```bash
   nessi metrics visualize /path/to/table --type radar --metric quality
   ```

3. **Integration with External Systems**: Connect with Databricks, AWS, or Grafana
   ```bash
   nessi integration databricks connect --host your-workspace.cloud.databricks.com --token your-token
   ```

4. **Scheduled Operations**: Set up regular scans and reports
   ```bash
   nessi scan schedule create --source /path/to/table --type all --schedule "0 0 * * *"
   ```

## Common Configurations

### Setting Default Output Format

```bash
nessi config set output.format json
```

### Configuring Log Level

```bash
nessi config set log.level debug
```

### Setting Default Delta Path

```bash
nessi config set delta.default_path /data/delta
```

## Troubleshooting

If you encounter issues:

1. Check your configuration:
   ```bash
   nessi config list
   ```

2. Enable debug logging:
   ```bash
   nessi --log-level debug <command>
   ```

3. Verify installation:
   ```bash
   nessi info
   ```

For more detailed troubleshooting, see the [Installation Guide](INSTALLATION.md#troubleshooting).

## Additional Resources

- [User Guide](USER_GUIDE.md) - Comprehensive documentation
- [CLI Reference](cli/REFERENCE.md) - Quick reference for all commands
- [Installation Guide](INSTALLATION.md) - Detailed installation instructions
- [FAQ](FREQUENTLY_ASKED_QUESTIONS.md) - Answers to common questions
- [Website](https://nessi.dev) - Official website with tutorials and examples
