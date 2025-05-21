# Nessi Usability Features

This document describes the usability features implemented in Nessi to make it more error-free and easier to use.

## Table of Contents

- [Input Validation](#input-validation)
- [Guided Setup Wizard](#guided-setup-wizard)
- [Dry Run Mode](#dry-run-mode)
- [Configuration System](#configuration-system)
- [Command Line Flags](#command-line-flags)

## Input Validation

Nessi includes a comprehensive input validation system that validates and normalizes user input before executing commands. This helps prevent errors by catching invalid input early and providing clear error messages.

### Path Validation

The path validation system ensures that file paths are valid and accessible before operations are performed:

```go
// Validate a path
validator := common.NewPathValidator()
path, err := validator.ValidatePath(inputPath)
if err != nil {
    // Handle invalid path error
}
```

Features include:

- Expansion of home directory (`~`)
- Conversion of relative paths to absolute paths
- Verification that paths exist (configurable)
- Verification that paths are directories or files as required
- Normalization of path separators

### Configuration Validation

The configuration validation system ensures that configuration values are valid:

```go
// Validate a configuration value
validator := common.NewConfigValidator()
validator.Required = true
err := validator.ValidateConfig("host", configValue)
if err != nil {
    // Handle invalid configuration error
}
```

Features include:

- Validation of required fields
- Type checking
- Range validation
- Format validation

### Output Format Validation

The output format validation system ensures that output formats are valid:

```go
// Validate an output format
format, err := common.ValidateOutputFormat(inputFormat)
if err != nil {
    // Handle invalid format error
}
```

Supported formats include:

- HTML
- PDF
- JSON
- CSV
- Text

## Guided Setup Wizard

The guided setup wizard helps new users configure Nessi for optimal use. It walks users through the initial setup process, creating configuration files and validating the environment.

```bash
# Run the setup wizard
nessi setup
```

The wizard covers:

1. **Configuration Directory**: Sets up the Nessi configuration directory
2. **Error Handling Configuration**: Configures interactive error resolution, telemetry, and retries
3. **Environment Check**: Verifies that required dependencies are installed
4. **Configuration File Creation**: Creates a configuration file with user-specified settings

## Dry Run Mode

Dry run mode allows users to see what would happen without making any changes. This is useful for understanding the effects of a command before executing it.

```bash
# Run a command in dry run mode
nessi repair --table mytable --dry-run
```

Dry run mode is implemented using the `DryRunManager` class:

```go
// Create a dry run manager
dryRun, _ := cmd.Flags().GetBool("dry-run")
manager := common.NewDryRunManager(dryRun)

// Add actions that would be performed
manager.AddAction("Create directory: %s", path)
manager.AddAction("Write file: %s", filePath)

// Check if we should execute
if manager.ShouldExecute() {
    // Perform the actual operations
    os.MkdirAll(path, 0755)
    os.WriteFile(filePath, data, 0644)
}

// Print actions at the end
manager.PrintActions()
```

The `DryRunManager` class provides:

- Tracking of actions that would be performed
- Conditional execution based on dry run mode
- Pretty printing of actions

## Configuration System

Nessi uses a flexible configuration system that supports multiple sources of configuration with clear precedence:

1. Command line flags (highest precedence)
2. Environment variables
3. Configuration file
4. Default values (lowest precedence)

The configuration file is located at `~/.nessi/config.yaml` by default, but can be specified with the `--config` flag.

See the [example configuration file](../config/config.yaml.example) for all available options.

### Configuration Structure

The configuration file is organized into sections:

- `error_handling`: Error handling configuration
  - `interactive_resolution`: Whether to enable interactive error resolution
  - `telemetry`: Error telemetry configuration
  - `retry`: Automatic retry configuration
  - `suggestions`: Error suggestion configuration
- `logging`: Logging configuration
- `validation`: Validation configuration
- `delta`: Delta Lake configuration
- `aliases`: Command aliases
- `preferences`: User preferences

## Command Line Flags

Nessi provides several command line flags to control its behavior:

### Global Flags

These flags are available for all commands:

- `--interactive`: Enable interactive mode
- `--dry-run`: Show what would be done without making changes
- `--config`: Specify the configuration file to use
- `--verbose`: Enable verbose output
- `--quiet`: Suppress all output except errors

### Command-Specific Flags

Each command may have additional flags specific to its functionality. Use `nessi <command> --help` to see the available flags for a command.

## Integration with Error Handling

All of these usability features are integrated with Nessi's error handling system to provide a consistent and user-friendly experience. When errors occur:

1. The error is validated and normalized
2. Clear error messages are displayed with error codes
3. Suggestions for resolving the error are provided
4. Interactive resolution is offered if enabled
5. Automatic retry is attempted if enabled

See the [Error Handling Documentation](ERROR_HANDLING.md) for more information on Nessi's error handling system.
