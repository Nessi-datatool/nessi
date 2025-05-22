# Nessi Migration Guide

This guide provides comprehensive instructions for migrating between different versions of Nessi, including upgrading from previous versions and migrating data, configurations, and integrations.

## Table of Contents

- [Overview](#overview)
- [Version Compatibility](#version-compatibility)
- [Upgrade Paths](#upgrade-paths)
- [Backup Procedures](#backup-procedures)
- [Configuration Migration](#configuration-migration)
- [Data Migration](#data-migration)
- [Rule Migration](#rule-migration)
- [Report Template Migration](#report-template-migration)
- [Integration Migration](#integration-migration)
- [License Migration](#license-migration)
- [Post-Migration Verification](#post-migration-verification)
- [Rollback Procedures](#rollback-procedures)
- [Version-Specific Migration Notes](#version-specific-migration-notes)
- [Troubleshooting](#troubleshooting)

## Overview

Migrating to a new version of Nessi involves several steps, including backing up your existing data and configurations, installing the new version, and migrating your configurations and data. This guide provides detailed instructions for each step of the migration process.

## Version Compatibility

### Compatibility Matrix

| From Version | To Version | Compatibility | Migration Complexity |
|--------------|------------|---------------|----------------------|
| 0.x.x        | 1.0.0      | Partial       | High                 |
| 1.0.x        | 1.1.x      | Full          | Low                  |
| 1.1.x        | 1.2.x      | Full          | Low                  |
| 1.x.x        | 2.0.0      | Partial       | Medium               |

### Compatibility Notes

- **Full Compatibility**: Direct upgrade supported with minimal changes
- **Partial Compatibility**: Upgrade requires configuration changes or data migration
- **Not Compatible**: Requires a multi-step migration process

## Upgrade Paths

### Recommended Upgrade Paths

For the smoothest upgrade experience, follow these recommended paths:

- **0.x.x → 1.0.0**: Upgrade directly, but expect configuration changes
- **1.0.x → 1.1.x → 1.2.x**: Incremental upgrades for best results
- **1.x.x → 2.0.0**: Direct upgrade supported, but review breaking changes

### Skipping Versions

If you need to skip multiple versions, we recommend:

1. Backup your data and configurations
2. Review the release notes for all intermediate versions
3. Follow the migration notes for each major version change
4. Consider a test migration in a non-production environment first

## Backup Procedures

Before migrating, always back up your Nessi data and configurations:

### Configuration Backup

```bash
# Create a backup directory
mkdir -p ~/nessi-backup/config

# Copy configuration files
cp -r ~/.nessi/config.yaml ~/nessi-backup/config/
cp -r ~/.nessi/rules/ ~/nessi-backup/rules/
cp -r ~/.nessi/templates/ ~/nessi-backup/templates/
```

### Data Backup

```bash
# Create a backup directory for data
mkdir -p ~/nessi-backup/data

# Backup metrics data
cp -r ~/.nessi/metrics/ ~/nessi-backup/data/metrics/

# Backup reports
cp -r ~/.nessi/reports/ ~/nessi-backup/data/reports/
```

### License Backup

```bash
# Backup license information
cp ~/.nessi/license.json ~/nessi-backup/license.json
```

## Configuration Migration

### Automatic Configuration Migration

Nessi provides a tool to automatically migrate your configuration:

```bash
# Migrate configuration from previous version
nessi config migrate --from-version 1.0.0 --to-version 1.1.0
```

### Manual Configuration Migration

If automatic migration fails or you prefer manual control:

1. Install the new version of Nessi
2. Create a new default configuration:
   ```bash
   nessi config init --output ~/new-config.yaml
   ```
3. Compare with your old configuration:
   ```bash
   diff ~/.nessi/config.yaml ~/new-config.yaml
   ```
4. Merge your customizations into the new configuration
5. Replace the default configuration:
   ```bash
   cp ~/new-config.yaml ~/.nessi/config.yaml
   ```

### Configuration Format Changes

#### Version 1.0.0 to 1.1.0

```yaml
# Old format (1.0.0)
storage:
  type: local
  path: /path/to/data

# New format (1.1.0)
storage:
  provider: local
  options:
    path: /path/to/data
```

#### Version 1.1.0 to 1.2.0

```yaml
# Old format (1.1.0)
quality:
  rules_path: /path/to/rules

# New format (1.2.0)
quality:
  rules:
    path: /path/to/rules
    auto_reload: true
```

## Data Migration

### Metrics Data Migration

For metrics data format changes:

```bash
# Migrate metrics data
nessi metrics migrate --from-version 1.0.0 --to-version 1.1.0 --metrics-dir ~/.nessi/metrics
```

### Report Data Migration

For report data format changes:

```bash
# Migrate report data
nessi report migrate --from-version 1.0.0 --to-version 1.1.0 --reports-dir ~/.nessi/reports
```

### Delta Table Compatibility

Nessi maintains backward compatibility with Delta Lake tables. However, new features may require table optimization:

```bash
# Optimize Delta tables for new version
nessi tables optimize --path /path/to/table --compatibility-version 1.1.0
```

## Rule Migration

### Quality Rule Format Changes

Quality rule formats may change between versions:

#### Version 1.0.0 to 1.1.0

```yaml
# Old format (1.0.0)
rules:
  - name: not_null
    columns: [id, name]
    
# New format (1.1.0)
rules:
  - name: not_null
    columns:
      - name: id
        severity: error
      - name: name
        severity: warning
```

### Rule Migration Tool

Use the rule migration tool to update your rules:

```bash
# Migrate quality rules
nessi quality migrate-rules --from-version 1.0.0 --to-version 1.1.0 --rules-dir ~/.nessi/rules
```

## Report Template Migration

### Template Format Changes

Report templates may change between versions:

```bash
# Migrate report templates
nessi report migrate-templates --from-version 1.0.0 --to-version 1.1.0 --templates-dir ~/.nessi/templates
```

### Custom Template Migration

For custom templates:

1. Review the template changes in the new version
2. Update your custom templates to match the new format
3. Test your templates with the new version

## Integration Migration

### Databricks Integration Changes (Pro Edition)

Databricks integration settings may change between versions:

```yaml
# Old format (1.0.0)
integrations:
  databricks:
    host: your-databricks-instance.cloud.databricks.com
    token: dapi_xxxxxxxxxxxxxxxxxxxxxxxx

# New format (1.1.0)
integrations:
  databricks:
    connection:
      host: your-databricks-instance.cloud.databricks.com
      token: dapi_xxxxxxxxxxxxxxxxxxxxxxxx
    options:
      timeout_seconds: 300
      retry_count: 3
```

### AWS Integration Changes (Pro Edition)

AWS integration settings may change between versions:

```yaml
# Old format (1.0.0)
integrations:
  aws:
    region: us-west-2
    access_key_id: AKIAXXXXXXXXXXXXXXXX
    secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# New format (1.1.0)
integrations:
  aws:
    connection:
      region: us-west-2
      credentials:
        access_key_id: AKIAXXXXXXXXXXXXXXXX
        secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    options:
      s3:
        max_concurrent_requests: 20
```

## License Migration

### License Migration Process

When upgrading Nessi with an active license:

1. Backup your license file:
   ```bash
   cp ~/.nessi/license.json ~/nessi-backup/license.json
   ```

2. Install the new version of Nessi

3. Restore your license:
   ```bash
   cp ~/nessi-backup/license.json ~/.nessi/license.json
   ```

4. Verify your license status:
   ```bash
   nessi config license status
   ```

### License Format Changes

If the license format has changed:

```bash
# Migrate license to new format
nessi config license migrate --from-version 1.0.0 --to-version 1.1.0
```

## Post-Migration Verification

After migration, verify that everything is working correctly:

### Configuration Verification

```bash
# Verify configuration
nessi config validate
```

### Functionality Verification

```bash
# Verify basic functionality
nessi --version
nessi tables list --path /path/to/data
nessi quality check --path /path/to/table --rules /path/to/rules.yaml
```

### Integration Verification (Pro Edition)

```bash
# Verify Databricks integration
nessi integration databricks test-connection

# Verify AWS integration
nessi integration aws test-connection
```

### License Verification (Pro Edition)

```bash
# Verify license status
nessi config license status
```

## Rollback Procedures

If you encounter issues after migration, you can roll back to the previous version:

### Configuration Rollback

```bash
# Restore configuration from backup
cp -r ~/nessi-backup/config/* ~/.nessi/
```

### Binary Rollback

```bash
# Remove new version
rm /usr/local/bin/nessi

# Restore previous version
cp ~/nessi-backup/bin/nessi /usr/local/bin/nessi
```

### Complete Rollback

For a complete rollback:

1. Uninstall the new version
2. Reinstall the previous version
3. Restore all backups

## Version-Specific Migration Notes

### Migrating from 0.x to 1.0

#### Breaking Changes

- Configuration format has changed significantly
- Quality rule format has been updated
- Report templates use a new structure
- Command-line interface has been standardized

#### Migration Steps

1. Backup all data and configurations
2. Install version 1.0.0
3. Use the migration tools to update configurations and rules
4. Update custom scripts that use the CLI
5. Verify all functionality

### Migrating from 1.0 to 1.1

#### Changes

- Storage configuration format has been updated
- New metrics collection capabilities
- Enhanced report templates

#### Migration Steps

1. Backup all data and configurations
2. Install version 1.1.0
3. Use the migration tools to update configurations
4. Verify all functionality

### Migrating from 1.x to 2.0

#### Breaking Changes

- Quality rule engine has been completely redesigned
- New plugin architecture
- Configuration structure has been reorganized

#### Migration Steps

1. Backup all data and configurations
2. Install version 2.0.0
3. Use the migration tools to update configurations and rules
4. Update custom plugins to use the new architecture
5. Verify all functionality

## Troubleshooting

### Common Migration Issues

#### Configuration Not Recognized

**Issue**: New version doesn't recognize your configuration

**Solution**:
1. Verify the configuration format matches the new version
2. Run the configuration migration tool
3. Check for syntax errors in the configuration file

#### Rules Not Working

**Issue**: Quality rules don't work after migration

**Solution**:
1. Verify the rule format matches the new version
2. Run the rule migration tool
3. Check for syntax errors in the rule definitions

#### Integration Failures

**Issue**: Integrations fail after migration

**Solution**:
1. Verify the integration configuration matches the new version
2. Update authentication credentials if necessary
3. Check network connectivity
4. Verify that the integration is supported in the new version

#### License Issues

**Issue**: License not recognized after migration

**Solution**:
1. Verify the license file is in the correct location
2. Run the license migration tool
3. Contact support if the issue persists

### Getting Help

If you encounter issues during migration:

1. Check the [troubleshooting guide](TROUBLESHOOTING.md)
2. Search for similar issues on the [GitHub Issues](https://github.com/nessi-dev/nessi/issues) page
3. Contact support at support@nessi.dev (Pro Edition)
4. Open a new issue if your problem is not already reported

---

For more information on Nessi migration, please visit [nessi.dev/migration](https://nessi.dev/migration) or contact support@nessi.dev.
