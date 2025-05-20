# Audit Logging System

## Overview

Nessi.dev implements a comprehensive audit logging system that tracks key operations and changes throughout the application. This system provides accountability, helps with troubleshooting, and supports compliance requirements by maintaining a detailed record of user actions and system events.

## Core Components

### Action Types

The audit logging system tracks the following types of actions:

- **Create**: Creation of resources
- **Read**: Reading or viewing resources
- **Update**: Modification of existing resources
- **Delete**: Deletion of resources
- **Execute**: Execution of operations
- **Login**: User authentication
- **Logout**: User session termination
- **ConfigChange**: Changes to system configuration

### Resource Types

The system logs actions performed on the following resource types:

- **Table**: Delta tables and their metadata
- **Rule**: Validation rules
- **User**: User accounts
- **Webhook**: Webhook configurations
- **Plugin**: Plugin management
- **Config**: System configuration
- **Report**: Generated reports

### Audit Entry Structure

Each audit entry contains the following information:

- **ID**: Unique identifier for the audit entry
- **Timestamp**: When the action occurred
- **UserID**: The user who performed the action
- **ActionType**: The type of action performed
- **ResourceType**: The type of resource affected
- **ResourceID**: The specific resource identifier
- **Details**: Additional context about the action
- **Success**: Whether the action succeeded
- **ErrorMessage**: Error details if the action failed
- **IPAddress**: The IP address of the user
- **UserAgent**: The user's browser or client information

## Usage

### CLI Commands

```bash
# View recent audit logs
nessi audit logs

# View audit logs for a specific user
nessi audit logs --user user1

# View audit logs for a specific resource
nessi audit logs --resource-type table --resource-id table1

# View audit logs for a specific action type
nessi audit logs --action-type update

# View audit logs within a date range
nessi audit logs --from 2025-01-01 --to 2025-01-31

# Export audit logs to a file
nessi audit logs --export audit_logs.json
```

### API Endpoints

The audit logging system exposes the following API endpoints:

- `GET /api/v1/audit/logs` - Query audit logs with filtering options
- `GET /api/v1/audit/logs/recent` - Get recent audit logs
- `GET /api/v1/audit/logs/users/{id}` - Get audit logs for a specific user
- `GET /api/v1/audit/logs/resources/{type}/{id}` - Get audit logs for a specific resource

## Configuration

The audit logging system can be configured in the `config/config.yaml` file:

```yaml
security:
  audit:
    enabled: true
    log_path: "/var/log/nessi/audit.log"
    retention_days: 90
    log_level: info
```

## Implementation Details

### Storage

The audit logging system supports multiple storage backends:

- **File**: Logs are stored in a file on disk
- **Database**: Logs are stored in a database for easier querying
- **Syslog**: Logs are sent to the system's syslog service

### Performance Considerations

The audit logging system is designed to have minimal impact on application performance:

- Asynchronous logging for non-critical operations
- Batched writes to reduce I/O overhead
- Configurable log levels to control verbosity

### Security

The audit logs themselves are protected:

- Access to audit logs requires administrative privileges
- Logs are tamper-evident through sequential IDs and timestamps
- Sensitive information is redacted from log entries

## Integration with Other Systems

### RBAC Integration

The audit logging system integrates with the RBAC system to log permission checks and changes to roles and permissions.

### Webhook Integration

Audit events can trigger webhooks to notify external systems about important actions.

### Monitoring Integration

The audit logging system provides metrics about system usage that can be exported to various formats (JSON, CSV) for monitoring and analysis.

## Best Practices

1. **Regular Review**: Regularly review audit logs to identify unusual patterns or security issues
2. **Retention Policy**: Establish a retention policy that balances storage requirements with compliance needs
3. **Backup Strategy**: Include audit logs in your backup strategy to ensure they are preserved
4. **Monitoring**: Set up alerts for critical audit events that may indicate security issues

## Examples

### Tracking Configuration Changes

```bash
# View all configuration changes
nessi audit logs --action-type config_change

# View configuration changes by a specific user
nessi audit logs --action-type config_change --user admin1
```

### Investigating Security Incidents

```bash
# View failed login attempts
nessi audit logs --action-type login --success false

# View all actions by a specific IP address
nessi audit logs --ip-address 192.168.1.100
```

## Troubleshooting

### Common Issues

1. **Missing Logs**: Ensure the audit logging system is enabled and configured correctly
2. **Performance Impact**: If logging is affecting performance, consider adjusting the log level or using asynchronous logging
3. **Disk Space**: Monitor disk space usage if file-based logging is used with a long retention period

### Debugging

Use the system logs to troubleshoot issues with the audit logging system:

```bash
nessi system logs --component audit
```
