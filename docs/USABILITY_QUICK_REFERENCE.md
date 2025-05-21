# Nessi Usability Features - Quick Reference

This quick reference guide provides an overview of Nessi's usability features that make it more error-free and easier to use.

## Setup Wizard

```bash
# Run the interactive setup wizard
nessi setup

# Run in non-interactive mode with default values
nessi setup --non-interactive --config-dir ~/.nessi
```

## Dry Run Mode

```bash
# Show what would happen without making changes
nessi schema create --path /path/to/table --dry-run

# Preview repair operations
nessi repair --table /path/to/table --dry-run

# Preview any command's actions
nessi <command> --dry-run
```

## Input Validation

Nessi automatically validates and normalizes input:

### Path Validation

- Expands home directory (`~`)
- Converts relative paths to absolute paths
- Verifies that paths exist (configurable)
- Verifies that paths are directories or files as required

### Configuration Validation

```bash
# Validate a configuration file
nessi config validate --file config.yaml
```

### Output Format Validation

Supported formats:
- HTML
- PDF
- JSON
- CSV
- Text

## Common Command Line Flags

```bash
# Enable interactive mode
nessi <command> --interactive

# Preview operations without making changes
nessi <command> --dry-run

# Specify configuration file
nessi <command> --config /path/to/config.yaml

# Enable verbose output
nessi <command> --verbose
```

## Error Handling

```bash
# List all available error codes
nessi test-error --list

# Generate a specific error to test handling
nessi test-error N101  # Path error

# Test interactive resolution
nessi test-error N201 --resolvable

# Test error telemetry
nessi test-error N301 --telemetry
```

## Configuration

Usability features can be configured in `~/.nessi/config.yaml`:

```yaml
# Error handling configuration
error_handling:
  # Whether to enable interactive error resolution
  interactive_resolution: true

  # Error telemetry configuration
  telemetry:
    enabled: true
    max_errors: 100

  # Automatic retry configuration
  retry:
    enabled: true
    max_attempts: 3
    initial_backoff: 1
    max_backoff: 30
    backoff_factor: 2.0
    jitter: true

# Validation configuration
validation:
  # Whether to validate paths before operations
  validate_paths: true
  # Whether to validate configuration values
  validate_config: true
  # Whether to normalize paths (convert relative to absolute, expand ~)
  normalize_paths: true
```

## Troubleshooting

If you encounter issues:

1. Run with verbose logging:
   ```bash
   NESSI_LOG_LEVEL=debug nessi <command>
   ```

2. Check error codes reference:
   ```bash
   nessi test-error --list
   ```

3. Validate your configuration:
   ```bash
   nessi config validate --file ~/.nessi/config.yaml
   ```

4. Run the setup wizard to reset your configuration:
   ```bash
   nessi setup
   ```

For more detailed information, see the [Usability Features Documentation](USABILITY_FEATURES.md) and [Error Handling Documentation](ERROR_HANDLING.md).
