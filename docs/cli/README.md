# Nessi CLI Documentation

## Overview

The Nessi Command Line Interface (CLI) provides a powerful set of tools for working with Delta Lake tables, performing data quality checks, and managing your data lake. This documentation covers all available commands, their options, and usage examples.

## Installation

### Binary Installation

Download the appropriate binary for your platform from the [releases page](https://github.com/nessi-dev/nessi/releases).

```bash
# Linux/macOS
chmod +x nessi
sudo mv nessi /usr/local/bin/

# Verify installation
nessi version
```

### Docker Installation

```bash
docker pull ghcr.io/nessi-dev/nessi:latest

# Run Nessi commands
docker run --rm -v $(pwd):/data ghcr.io/nessi-dev/nessi:latest <command>
```

### Build from Source

```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
go build -o nessi ./cmd/nessi
```

## Configuration

Nessi can be configured using a YAML configuration file, environment variables, or command-line flags.

### Configuration File

Default location: `~/.nessi/config.yaml`

```yaml
# Example configuration file
storage:
  type: local
  path: /path/to/data
logging:
  level: info
  file: /path/to/logs/nessi.log
```

Specify a custom configuration file:

```bash
nessi --config /path/to/config.yaml <command>
```

### Environment Variables

All configuration options can be set using environment variables with the `NESSI_` prefix:

```bash
export NESSI_STORAGE_TYPE=local
export NESSI_STORAGE_PATH=/path/to/data
export NESSI_LOGGING_LEVEL=info
```

## Command Structure

Nessi commands follow a consistent structure:

```
nessi [global options] command [command options] [arguments...]
```

## Global Options

| Option | Description |
|--------|-------------|
| `--config` | Path to configuration file |
| `--verbose` | Enable verbose output |
| `--quiet` | Suppress all output except errors |
| `--output` | Output format (text, json, csv) |
| `--help, -h` | Show help |
| `--version, -v` | Show version |

## Core Commands (Community Edition)

### Basic Commands

| Command | Description |
|---------|-------------|
| `help` | Shows a list of commands or help for one command |
| `version` | Shows the Nessi version information |
| `info` | Shows system information and configuration |
| `completion` | Generates shell completion scripts |

### Table Commands

| Command | Description |
|---------|-------------|
| `tables list` | Lists all available Delta Lake tables |
| `tables describe` | Shows detailed information about a table |
| `tables schema` | Shows the schema of a table |
| `tables history` | Shows the history of a table |
| `tables profile` | Generates a profile of a table's data |

### Quality Commands

| Command | Description |
|---------|-------------|
| `quality check` | Runs quality checks on a table |
| `quality rules list` | Lists all available quality rules |
| `quality rules create` | Creates a new quality rule |
| `quality rules delete` | Deletes a quality rule |
| `quality report` | Generates a quality report for a table |

### Schema Commands

| Command | Description |
|---------|-------------|
| `schema validate` | Validates a table's schema |
| `schema compare` | Compares schemas between tables or versions |
| `schema tree` | Displays a schema as a tree structure |

## Premium Commands (Pro Edition)

### Cloud Storage Commands

| Command | Description |
|---------|-------------|
| `cloud s3 list` | Lists tables in S3 storage |
| `cloud azure list` | Lists tables in Azure storage |
| `cloud gcp list` | Lists tables in GCP storage |

### Integration Commands

| Command | Description |
|---------|-------------|
| `databricks connect` | Connects to a Databricks workspace |
| `databricks list` | Lists tables in a Databricks workspace |
| `databricks import` | Imports a table from Databricks |
| `dbt connect` | Connects to a dbt project |
| `dbt run` | Runs dbt models and validates results |

### Catalog Commands

| Command | Description |
|---------|-------------|
| `catalog list` | Lists all tables in the catalog |
| `catalog register` | Registers a table in the catalog |
| `catalog search` | Searches for tables in the catalog |
| `catalog tag` | Manages tags for tables |

### Workflow Commands

| Command | Description |
|---------|-------------|
| `workflow list` | Lists all workflows |
| `workflow create` | Creates a new workflow |
| `workflow run` | Runs a workflow |
| `workflow schedule` | Schedules a workflow |

## Command Documentation

For detailed documentation on each command, see:

- [Basic Commands](basic.md)
- [Table Commands](tables.md)
- [Quality Commands](quality.md)
- [Schema Commands](schema.md)
- [Cloud Storage Commands](cloud.md) (Pro Edition)
- [Integration Commands](integrations.md) (Pro Edition)
- [Catalog Commands](catalog.md) (Pro Edition)
- [Workflow Commands](workflow.md) (Pro Edition)

## Examples

### List all tables

```bash
nessi tables list
```

### Describe a table

```bash
nessi tables describe my_table
```

### Generate a table profile

```bash
nessi tables profile my_table --sample 50
```

### Run quality checks

```bash
nessi quality check my_table --rules completeness,uniqueness
```

### Generate a quality report

```bash
nessi quality report my_table --output json > quality_report.json
```

### Connect to Databricks (Pro Edition)

```bash
nessi databricks connect --host https://dbc-123456789.cloud.databricks.com --token <token>
```

## Shell Completion

Nessi provides shell completion for Bash, Zsh, Fish, and PowerShell.

```bash
# Bash
nessi completion bash > /etc/bash_completion.d/nessi

# Zsh
nessi completion zsh > "${fpath[1]}/_nessi"

# Fish
nessi completion fish > ~/.config/fish/completions/nessi.fish

# PowerShell
nessi completion powershell > nessi.ps1
```

## License Management

Some commands are only available in the Pro Edition. To check your license status:

```bash
nessi license status
```

To activate a license:

```bash
nessi license activate <license-key>
```

For more information on Nessi's licensing, see the [License Management documentation](../LICENSE_MANAGEMENT.md).

## Troubleshooting

### Common Issues

#### Command not found

Ensure Nessi is properly installed and in your PATH:

```bash
which nessi
```

#### Permission denied

Ensure the Nessi binary has execute permissions:

```bash
chmod +x /path/to/nessi
```

#### Configuration errors

Check your configuration file for errors:

```bash
nessi info --verbose
```

### Error Codes

Nessi uses structured error codes to help diagnose issues:

- `N1XX`: Path errors
- `N2XX`: Configuration errors
- `N3XX`: Authentication errors
- `N4XX`: Connection errors
- `N5XX`: Delta Lake errors
- `N6XX`: Databricks errors
- `N7XX`: Schema errors
- `N8XX`: Validation errors
- `N9XX`: Internal errors

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).

## Getting Help

If you need additional help:

- Run `nessi help <command>` for detailed command help
- Visit the [Nessi documentation](https://nessi.dev/docs)
- Join the [Nessi community](https://github.com/nessi-dev/nessi/discussions)
- For Pro Edition support, contact [support@nessi.dev](mailto:support@nessi.dev)
