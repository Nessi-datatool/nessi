# Getting Started with Nessi

This guide will help you quickly get up and running with Nessi, the open-source CLI-only data quality and Delta Lake management tool.

## Installation

### Prerequisites

- Go 1.18 or higher
- Python 3.8 or higher (for Python integrations)
- Access to Delta Lake tables

### Option 1: Install Using Go

```bash
go install github.com/nessi-dev/nessi/cmd/nessi@latest
```

### Option 2: Download Prebuilt Binary

1. Download the latest release from the [GitHub Releases page](https://github.com/nessi-dev/nessi/releases)
2. Make the binary executable:
   ```bash
   chmod +x nessi
   ```
3. Move to your PATH:
   ```bash
   mv nessi /usr/local/bin/
   ```

### Option 3: Docker Installation

```bash
docker pull nessi/nessi:latest
```

## Quick Start

### Basic Commands

Check your installation:

```bash
nessi --version
```

Get help:

```bash
nessi --help
```

### Working with Delta Tables

#### Scan a Delta Table

```bash
# Local table
nessi scan /path/to/delta/table

# S3 table (requires Pro license or trial)
nessi scan s3://bucket/path/to/table
```

#### Generate a Data Quality Report

```bash
nessi report /path/to/delta/table --format html --output ./reports
```

#### Check Schema Evolution

```bash
nessi schema history /path/to/delta/table
```

#### Time Travel

```bash
# View table at a specific version
nessi view /path/to/delta/table --version 5

# View table at a specific timestamp
nessi view /path/to/delta/table --timestamp "2025-01-01T00:00:00Z"
```

## Configuration

Create a `nessi.yaml` configuration file in your home directory at `~/.nessi/config.yaml`:

```yaml
storage:
  type: local  # or s3, azure, gcp with Pro license
  base_path: /path/to/delta/tables

quality:
  default_rules: path/to/rules.yaml
  threshold: 0.95

reporting:
  output_path: ./reports
  default_format: html
```

## Using Premium Features

Start a free trial to access premium features:

```bash
nessi license start-trial
```

Check your license status:

```bash
nessi license status
```

## Common Workflows

### Data Quality Validation Workflow

1. Define quality rules in a YAML file:
   ```yaml
   rules:
     - name: "not_null_check"
       column: "user_id"
       type: "not_null"
     - name: "range_check"
       column: "age"
       type: "range"
       min: 0
       max: 120
   ```

2. Run validation:
   ```bash
   nessi validate /path/to/delta/table --rules ./rules.yaml
   ```

3. Generate a report:
   ```bash
   nessi report /path/to/delta/table --rules ./rules.yaml --format html
   ```

### Delta Lake Management Workflow

1. Check table health:
   ```bash
   nessi health /path/to/delta/table
   ```

2. Optimize table:
   ```bash
   nessi optimize /path/to/delta/table
   ```

3. View transaction history:
   ```bash
   nessi transactions /path/to/delta/table
   ```

## Next Steps

- Read the [full documentation](https://github.com/nessi-dev/nessi/docs)
- Try the [example projects](https://github.com/nessi-dev/nessi/examples)
- Join the [community](https://github.com/nessi-dev/nessi/discussions)

## Getting Help

If you encounter any issues:

1. Check the [FAQ](https://github.com/nessi-dev/nessi/docs/FAQ.md)
2. Search [existing issues](https://github.com/nessi-dev/nessi/issues)
3. Ask in the [discussions forum](https://github.com/nessi-dev/nessi/discussions)
