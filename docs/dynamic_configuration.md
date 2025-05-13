# Dynamic Configuration System

## Overview

Nessi.dev implements a dynamic configuration system that allows for loading and reloading configuration at runtime. This system enables configuration changes without requiring application restarts, providing flexibility and improved operational efficiency.

## Core Components

### Configuration Formats

The system supports the following configuration formats:

- **JSON**: Standard JSON format
- **YAML**: YAML format for more human-readable configuration

### Dynamic Config Manager

The Dynamic Config Manager is responsible for:

- Loading configuration from files
- Monitoring configuration files for changes
- Notifying components when configuration changes
- Providing thread-safe access to configuration values

## Usage

### CLI Commands

```bash
# View current configuration
nessi config view

# Reload configuration from disk
nessi config reload

# Set a configuration value
nessi config set server.port 8080

# Get a configuration value
nessi config get server.port

# Watch configuration changes
nessi config watch

# Export configuration to a file
nessi config export config.yaml

# Import configuration from a file
nessi config import config.yaml
```

### API Endpoints

The dynamic configuration system exposes the following API endpoints:

- `GET /api/v1/config` - Get the current configuration
- `POST /api/v1/config/reload` - Reload configuration from disk
- `PUT /api/v1/config/{path}` - Update a configuration value
- `GET /api/v1/config/{path}` - Get a specific configuration value
- `POST /api/v1/config/import` - Import configuration from a file
- `GET /api/v1/config/export` - Export configuration to a file

## Configuration

The dynamic configuration system itself can be configured in the `config/config.yaml` file:

```yaml
config:
  watch_enabled: true
  watch_interval: 30s
  backup_enabled: true
  backup_directory: "/var/lib/nessi/config_backups"
```

## Implementation Details

### Thread Safety

The dynamic configuration system is designed to be thread-safe, allowing multiple components to access configuration values concurrently without race conditions.

### File Watching

The system can monitor configuration files for changes and automatically reload when changes are detected. This is useful for environments where configuration is managed by external systems.

### Validation

Configuration values are validated when loaded to ensure they meet the required format and constraints. This prevents invalid configuration from being applied.

### Hierarchical Access

Configuration values can be accessed using dot notation (e.g., `server.port`), allowing for hierarchical organization of configuration.

### Change Notifications

Components can register for notifications when specific configuration values change, allowing them to adapt to new settings dynamically.

## Integration with Other Systems

### Feature Flags Integration

The dynamic configuration system integrates with the feature flags system, allowing feature flags to be updated through the configuration system.

### Audit Logging Integration

Configuration changes are logged in the audit system, providing accountability and traceability for configuration changes.

### RBAC Integration

Access to configuration values and the ability to make changes is controlled by the RBAC system, ensuring that only authorized users can modify configuration.

## Best Practices

1. **Validation**: Always validate configuration values before applying them
2. **Defaults**: Provide sensible default values for all configuration options
3. **Documentation**: Document all configuration options and their impact
4. **Versioning**: Consider versioning your configuration schema to handle upgrades
5. **Security**: Protect sensitive configuration values (e.g., API keys, passwords)

## Examples

### Loading Configuration from Different Sources

```go
// Load configuration from a file
config, err := dynamicConfig.LoadFromFile("/path/to/config.yaml")

// Load configuration from environment variables
config, err := dynamicConfig.LoadFromEnv()

// Load configuration from multiple sources with priority
config, err := dynamicConfig.LoadFromSources([]Source{
    {Type: FileSource, Path: "/path/to/config.yaml"},
    {Type: EnvSource, Prefix: "NESSI_"},
})
```

### Watching for Configuration Changes

```go
// Start watching a configuration file for changes
watcher, err := dynamicConfig.WatchFile("/path/to/config.yaml", func(newConfig *DynamicConfig) {
    // Handle configuration change
    applyNewConfiguration(newConfig)
})

// Stop watching when done
watcher.Stop()
```

## Troubleshooting

### Common Issues

1. **File Permissions**: Ensure the application has permission to read and write configuration files
2. **Format Errors**: Verify that configuration files are correctly formatted (valid JSON or YAML)
3. **Conflicting Changes**: Be cautious when multiple sources can modify the same configuration values

### Debugging

Use the system logs to troubleshoot issues with the dynamic configuration system:

```bash
nessi system logs --component config
```

### Configuration Backup and Recovery

The system automatically creates backups of configuration before applying changes, allowing for recovery if needed:

```bash
# List configuration backups
nessi config backups list

# Restore from a backup
nessi config backups restore backup-20250101-120000.yaml
```
