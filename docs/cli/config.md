# Configuration Commands

This document provides detailed information about Nessi's configuration commands for managing settings, profiles, and preferences.

## `config set`

Sets a configuration value.

### Usage

```bash
nessi config set <key> <value> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `key` | Configuration key (e.g., `log.level`, `output.format`) |
| `value` | Configuration value |

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--global` | Set configuration globally for all projects |
| `--project` | Set configuration for the current project only |
| `--temp` | Set configuration temporarily for this session only |

### Examples

```bash
# Set log level
nessi config set log.level debug

# Set default output format
nessi config set output.format json

# Set configuration for a specific profile
nessi config set metrics.store.path /data/metrics --profile production

# Set global configuration
nessi config set delta.default_path /data/delta --global

# Set project-specific configuration
nessi config set quality.rules_path /project/rules --project
```

### Output

```
Configuration set successfully:
  Key: log.level
  Value: debug
  Scope: project
  Profile: default
  File: /Users/username/.nessi/config.yaml
```

## `config get`

Gets a configuration value.

### Usage

```bash
nessi config get <key> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `key` | Configuration key (e.g., `log.level`, `output.format`) |

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--output` | Output format (text, json) |
| `--show-origin` | Show where the configuration value comes from |
| `--show-all` | Show all matching configuration values from different scopes |

### Examples

```bash
# Get log level
nessi config get log.level

# Get configuration for a specific profile
nessi config get metrics.store.path --profile production

# Show configuration origin
nessi config get delta.default_path --show-origin

# Show all matching configurations
nessi config get quality.rules_path --show-all
```

### Output

#### Text Format (Default)

```
log.level = debug
```

#### With Origin

```
log.level = debug (project: /Users/username/.nessi/config.yaml)
```

#### Show All

```
log.level:
  global: info (/Users/username/.nessi/config.yaml)
  project: debug (./.nessi.yaml)
  environment: debug (NESSI_LOG_LEVEL)
  effective: debug
```

## `config list`

Lists all configuration values.

### Usage

```bash
nessi config list [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--output` | Output format (text, json, yaml) |
| `--show-origin` | Show where each configuration value comes from |
| `--filter` | Filter configuration keys by pattern |
| `--scope` | Filter by scope (global, project, env, all) |

### Examples

```bash
# List all configuration
nessi config list

# List configuration for a specific profile
nessi config list --profile production

# Filter configuration
nessi config list --filter "log.*"

# Show configuration origin
nessi config list --show-origin

# List only project configuration
nessi config list --scope project
```

### Output

#### Text Format (Default)

```
Configuration:

log.level = debug
log.file = /var/log/nessi.log
output.format = json
delta.default_path = /data/delta
quality.rules_path = /project/rules
metrics.store.path = /data/metrics
report.template_path = /templates
```

#### With Origin

```
Configuration:

log.level = debug (project: ./.nessi.yaml)
log.file = /var/log/nessi.log (global: /Users/username/.nessi/config.yaml)
output.format = json (environment: NESSI_DEFAULT_FORMAT)
delta.default_path = /data/delta (global: /Users/username/.nessi/config.yaml)
quality.rules_path = /project/rules (project: ./.nessi.yaml)
metrics.store.path = /data/metrics (profile: production)
report.template_path = /templates (global: /Users/username/.nessi/config.yaml)
```

## `config unset`

Unsets a configuration value.

### Usage

```bash
nessi config unset <key> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `key` | Configuration key to unset |

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--global` | Unset from global configuration |
| `--project` | Unset from project configuration |
| `--all` | Unset from all scopes |

### Examples

```bash
# Unset a configuration value
nessi config unset log.level

# Unset from global configuration
nessi config unset delta.default_path --global

# Unset from a specific profile
nessi config unset metrics.store.path --profile production

# Unset from all scopes
nessi config unset quality.rules_path --all
```

### Output

```
Configuration unset successfully:
  Key: log.level
  Scope: project
  Profile: default
  File: /Users/username/.nessi/config.yaml
```

## `config init`

Initializes a new configuration file.

### Usage

```bash
nessi config init [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--global` | Initialize global configuration |
| `--project` | Initialize project configuration |
| `--template` | Template to use (default, minimal, full) |
| `--force` | Overwrite existing configuration file |
| `--interactive` | Use interactive mode to set values |

### Examples

```bash
# Initialize project configuration
nessi config init --project

# Initialize global configuration
nessi config init --global

# Initialize with a specific template
nessi config init --project --template full

# Initialize interactively
nessi config init --project --interactive

# Initialize a specific profile
nessi config init --profile production
```

### Output

```
Configuration initialized successfully:
  Scope: project
  Profile: default
  Template: default
  File: ./.nessi.yaml

The configuration file has been created with default settings.
You can modify it using 'nessi config set' or by editing the file directly.
```

## `config profile`

Manages configuration profiles.

### Usage

```bash
nessi config profile [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `list` | List all profiles |
| `create` | Create a new profile |
| `delete` | Delete a profile |
| `rename` | Rename a profile |
| `copy` | Copy a profile |
| `set-default` | Set the default profile |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--output` | Output format (text, json) |
| `--show-details` | Show detailed profile information |

### Options for `create`

| Option | Description |
|--------|-------------|
| `--name` | Profile name |
| `--base` | Base profile to copy from |
| `--description` | Profile description |
| `--interactive` | Use interactive mode to set values |

### Options for `delete`, `rename`, `copy`, and `set-default`

| Option | Description |
|--------|-------------|
| `--name` | Profile name |
| `--new-name` | New profile name (for rename and copy) |
| `--force` | Force operation without confirmation |

### Examples

```bash
# List all profiles
nessi config profile list

# Create a new profile
nessi config profile create --name production --base default

# Create a profile interactively
nessi config profile create --name staging --interactive

# Delete a profile
nessi config profile delete --name staging

# Rename a profile
nessi config profile rename --name production --new-name prod

# Copy a profile
nessi config profile copy --name prod --new-name prod-backup

# Set the default profile
nessi config profile set-default --name prod
```

### Output for `list`

```
Configuration Profiles:

NAME         DEFAULT  DESCRIPTION                       LAST MODIFIED
default      *        Default configuration             2025-05-01 12:34:56
production            Production environment settings   2025-05-10 09:23:45
development           Development environment settings  2025-05-15 14:56:23
testing               Testing environment settings      2025-05-20 10:11:12
```

### Output for `create`

```
Profile created successfully:
  Name: production
  Base: default
  File: /Users/username/.nessi/profiles/production.yaml

You can now set configuration values for this profile using:
  nessi config set <key> <value> --profile production
```

## `config import`

Imports configuration from a file.

### Usage

```bash
nessi config import <file-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `file-path` | Path to the configuration file to import |

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--global` | Import to global configuration |
| `--project` | Import to project configuration |
| `--merge` | Merge with existing configuration |
| `--overwrite` | Overwrite existing configuration |
| `--format` | Input file format (yaml, json, toml) |

### Examples

```bash
# Import configuration
nessi config import /path/to/config.yaml

# Import to a specific profile
nessi config import /path/to/config.yaml --profile production

# Import and merge with existing configuration
nessi config import /path/to/config.yaml --merge

# Import to global configuration
nessi config import /path/to/config.yaml --global

# Import from JSON format
nessi config import /path/to/config.json --format json
```

### Output

```
Configuration imported successfully:
  Source: /path/to/config.yaml
  Destination: /Users/username/.nessi/config.yaml
  Profile: default
  Scope: project
  Mode: merge
  Keys Imported: 15
  Keys Updated: 5
  Keys Added: 10
```

## `config export`

Exports configuration to a file.

### Usage

```bash
nessi config export <file-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `file-path` | Path to export the configuration to |

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--global` | Export global configuration |
| `--project` | Export project configuration |
| `--all` | Export all configuration (global, project, profiles) |
| `--format` | Output file format (yaml, json, toml) |
| `--filter` | Filter configuration keys by pattern |
| `--include-secrets` | Include secrets in the export (not recommended) |

### Examples

```bash
# Export configuration
nessi config export /path/to/config.yaml

# Export a specific profile
nessi config export /path/to/config.yaml --profile production

# Export global configuration
nessi config export /path/to/config.yaml --global

# Export all configuration
nessi config export /path/to/config.yaml --all

# Export to JSON format
nessi config export /path/to/config.json --format json

# Export filtered configuration
nessi config export /path/to/config.yaml --filter "log.*"
```

### Output

```
Configuration exported successfully:
  Destination: /path/to/config.yaml
  Profile: default
  Scope: project
  Format: yaml
  Keys Exported: 15
```

## `config validate`

Validates configuration.

### Usage

```bash
nessi config validate [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--global` | Validate global configuration |
| `--project` | Validate project configuration |
| `--all` | Validate all configuration (global, project, profiles) |
| `--fix` | Attempt to fix validation issues |
| `--output` | Output format (text, json) |

### Examples

```bash
# Validate current configuration
nessi config validate

# Validate a specific profile
nessi config validate --profile production

# Validate global configuration
nessi config validate --global

# Validate all configuration
nessi config validate --all

# Validate and fix issues
nessi config validate --fix
```

### Output

```
Configuration Validation Results:

Status: VALID
Profile: default
Scope: project
File: /Users/username/.nessi/config.yaml

No validation issues found.
```

### Output with Issues

```
Configuration Validation Results:

Status: INVALID
Profile: default
Scope: project
File: /Users/username/.nessi/config.yaml

Issues:
  - log.level: Invalid value 'trace' (valid values: debug, info, warn, error)
  - metrics.store.path: Path does not exist: /nonexistent/path
  - delta.default_path: Path is not writable: /readonly/path

Run 'nessi config validate --fix' to attempt to fix these issues automatically.
```

## `config env`

Manages environment variables for configuration.

### Usage

```bash
nessi config env [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `list` | List environment variables used by Nessi |
| `export` | Generate shell commands to export environment variables |
| `unset` | Generate shell commands to unset environment variables |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--output` | Output format (text, json) |
| `--show-values` | Show current values of environment variables |
| `--filter` | Filter by variable name pattern |

### Options for `export` and `unset`

| Option | Description |
|--------|-------------|
| `--profile` | Configuration profile name |
| `--shell` | Shell type (bash, zsh, fish, cmd, powershell) |
| `--file` | Output file path |
| `--include-secrets` | Include secrets in the export (not recommended) |

### Examples

```bash
# List environment variables
nessi config env list

# List with current values
nessi config env list --show-values

# Generate export commands for bash
nessi config env export --shell bash

# Generate export commands for a specific profile
nessi config env export --profile production --shell bash

# Generate unset commands
nessi config env unset --shell bash

# Save export commands to a file
nessi config env export --shell bash --file .env
```

### Output for `list`

```
Nessi Environment Variables:

VARIABLE                DESCRIPTION                                  OVERRIDES
NESSI_CONFIG            Path to configuration file                   --config
NESSI_LOG_LEVEL         Log level (debug, info, warn, error)         log.level
NESSI_LICENSE_PATH      Path to license file                         license.path
NESSI_NO_COLOR          Disable colored output if set to true        output.color
NESSI_DEFAULT_FORMAT    Default output format                        output.format
NESSI_METRICS_STORE     Path or URL to metrics store                 metrics.store.path
NESSI_SCHEMA_REGISTRY   URL to schema registry                       schema.registry.url
NESSI_DATABRICKS_TOKEN  Databricks access token                      integration.databricks.token
NESSI_DATABRICKS_HOST   Databricks host URL                          integration.databricks.host
NESSI_AWS_REGION        AWS region for S3 access                     integration.aws.region
```

### Output for `export` (Bash)

```
# Nessi environment variables for profile: default
export NESSI_CONFIG="/Users/username/.nessi/config.yaml"
export NESSI_LOG_LEVEL="debug"
export NESSI_DEFAULT_FORMAT="json"
export NESSI_METRICS_STORE="/data/metrics"
export NESSI_SCHEMA_REGISTRY="http://schema-registry:8081"
export NESSI_NO_COLOR="false"
```

## Configuration File Format

Nessi uses YAML as the default format for configuration files. The configuration is organized hierarchically with sections and subsections.

### Example Configuration File

```yaml
# Global settings
log:
  level: info
  file: /var/log/nessi.log
  format: json

# Output settings
output:
  format: text
  color: true
  pager: true

# Delta Lake settings
delta:
  default_path: /data/delta
  time_travel_enabled: true

# Quality settings
quality:
  rules_path: /path/to/rules
  default_threshold: 0.95

# Metrics settings
metrics:
  store:
    path: /path/to/metrics
    type: file
  collect_on_scan: true

# Report settings
report:
  template_path: /path/to/templates
  default_format: html

# Integration settings
integration:
  grafana:
    url: http://grafana:3000
    api_key: your-api-key
  databricks:
    host: https://your-workspace.cloud.databricks.com
    token: your-token
  aws:
    region: us-west-2
    profile: nessi
```

### Configuration Hierarchy

Nessi uses a hierarchical approach to configuration, with the following precedence (from highest to lowest):

1. Command-line arguments
2. Environment variables
3. Project-specific configuration (`.nessi.yaml` in the current directory)
4. User profile configuration (in `~/.nessi/profiles/<profile>.yaml`)
5. Global user configuration (in `~/.nessi/config.yaml`)
6. Default values

This allows for flexible configuration management across different environments and projects.

## Error Handling

Configuration commands use the following error codes:

- `N1200`: Configuration file not found
- `N1201`: Invalid configuration format
- `N1202`: Configuration validation failed
- `N1203`: Profile not found
- `N1204`: Permission denied for configuration file
- `N1205`: Configuration import/export failed
- `N1206`: Configuration key not found
- `N1207`: Invalid configuration value

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
