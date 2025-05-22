# Basic Commands

This document provides detailed information about Nessi's basic commands.

## `help`

Displays help information about available commands.

### Usage

```bash
nessi help [command]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `command` | (Optional) The command to get help for |

### Examples

```bash
# Show general help
nessi help

# Show help for a specific command
nessi help tables list

# Show help for a command category
nessi help quality
```

## `version`

Displays version information for Nessi.

### Usage

```bash
nessi version
```

### Output

The `version` command displays:
- Version number (e.g., v1.2.3)
- Git commit hash
- Build date
- Go version used for building
- Operating system and architecture

### Example

```bash
$ nessi version
Nessi CLI v1.2.3
Commit: a1b2c3d4e5f6
Built: 2025-05-22T12:00:00Z
Go version: go1.21.0
OS/Arch: darwin/amd64
```

## `info`

Displays system information and current configuration.

### Usage

```bash
nessi info [--verbose]
```

### Options

| Option | Description |
|--------|-------------|
| `--verbose` | Show detailed configuration information |

### Output

The `info` command displays:
- System information (OS, architecture)
- Configuration file location
- Storage configuration
- Logging configuration
- License information (if applicable)

### Example

```bash
$ nessi info
System Information:
  OS: darwin
  Architecture: amd64
  
Configuration:
  File: /Users/username/.nessi/config.yaml
  
Storage:
  Type: local
  Path: /Users/username/data
  
Logging:
  Level: info
  File: /Users/username/logs/nessi.log
  
License:
  Status: Active
  Type: Community Edition
```

With `--verbose` flag:

```bash
$ nessi info --verbose
System Information:
  OS: darwin
  Architecture: amd64
  CPU Cores: 8
  Memory: 16GB
  
Configuration:
  File: /Users/username/.nessi/config.yaml
  Environment Variables:
    NESSI_STORAGE_TYPE: local
    NESSI_LOGGING_LEVEL: info
  
Storage:
  Type: local
  Path: /Users/username/data
  Disk Space:
    Total: 500GB
    Used: 200GB
    Available: 300GB
  
Logging:
  Level: info
  File: /Users/username/logs/nessi.log
  Format: json
  Max Size: 100MB
  Max Age: 30 days
  
License:
  Status: Active
  Type: Community Edition
  Features:
    - Basic Table Operations
    - Quality Checks
    - Schema Validation
```

## `completion`

Generates shell completion scripts for various shells.

### Usage

```bash
nessi completion [shell]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `shell` | The shell to generate completion for (bash, zsh, fish, powershell) |

### Examples

```bash
# Generate Bash completion
nessi completion bash > /etc/bash_completion.d/nessi

# Generate Zsh completion
nessi completion zsh > "${fpath[1]}/_nessi"

# Generate Fish completion
nessi completion fish > ~/.config/fish/completions/nessi.fish

# Generate PowerShell completion
nessi completion powershell > nessi.ps1
```

### Installation Instructions

#### Bash

```bash
# Option 1: Source the completion script in your .bashrc
echo 'source <(nessi completion bash)' >> ~/.bashrc

# Option 2: Add the completion script to the completion directory
nessi completion bash > /etc/bash_completion.d/nessi
```

#### Zsh

```bash
# Option 1: Source the completion script in your .zshrc
echo 'source <(nessi completion zsh)' >> ~/.zshrc

# Option 2: Add the completion script to the completion directory
nessi completion zsh > "${fpath[1]}/_nessi"
```

#### Fish

```bash
nessi completion fish > ~/.config/fish/completions/nessi.fish
```

#### PowerShell

```powershell
nessi completion powershell > nessi.ps1
. ./nessi.ps1
```

## `setup`

Runs an interactive setup wizard to configure Nessi.

### Usage

```bash
nessi setup [--non-interactive]
```

### Options

| Option | Description |
|--------|-------------|
| `--non-interactive` | Run in non-interactive mode using default values |

### Interactive Setup

The setup wizard guides you through configuring:
- Storage settings
- Logging settings
- Integration settings (if applicable)
- License activation (if applicable)

### Example

```bash
$ nessi setup
Welcome to the Nessi setup wizard!

Storage Configuration:
1. Local
2. S3
3. Azure Blob Storage
4. Google Cloud Storage
Select storage type [1]: 1

Enter local storage path [/Users/username/data]: 

Logging Configuration:
1. Console only
2. File only
3. Console and file
Select logging option [3]: 3

Enter log level (debug, info, warn, error) [info]: 

Enter log file path [/Users/username/logs/nessi.log]: 

Setup complete! Configuration saved to /Users/username/.nessi/config.yaml
```

## `license`

Manages Nessi license (available in both Community and Pro editions).

### Usage

```bash
nessi license [command]
```

### Commands

| Command | Description |
|---------|-------------|
| `status` | Shows current license status |
| `activate` | Activates a license key |
| `deactivate` | Deactivates the current license |

### Examples

```bash
# Check license status
nessi license status

# Activate a license
nessi license activate XXXX-XXXX-XXXX-XXXX

# Deactivate a license
nessi license deactivate
```

### License Status Output

```bash
$ nessi license status
License Status: Active
Edition: Pro
Expiration: 2026-05-22
Features:
  - Basic Table Operations
  - Quality Checks
  - Schema Validation
  - Cloud Storage (S3, Azure, GCP)
  - Databricks Integration
  - dbt Integration
  - Data Catalog
  - Workflow Orchestration
```

## Common Options for All Basic Commands

All basic commands support the following global options:

| Option | Description |
|--------|-------------|
| `--config` | Path to configuration file |
| `--verbose` | Enable verbose output |
| `--quiet` | Suppress all output except errors |
| `--output` | Output format (text, json, csv) |
| `--help, -h` | Show help |

## Error Handling

Basic commands use the following error codes:

- `N200`: Configuration file not found
- `N201`: Invalid configuration format
- `N202`: Permission denied when reading configuration
- `N300`: License activation failed
- `N301`: License validation failed
- `N302`: License expired

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
