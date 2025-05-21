# Databricks Integration

Nessi provides integration with Databricks, allowing you to work with Databricks Unity Catalog and Delta Lake tables directly from the command line.

## Overview

The Databricks integration in Nessi enables you to:

- Browse and search Databricks Unity Catalog
- View and update table metadata
- Work with Delta Lake tables
- Perform data quality checks on Databricks tables

## Configuration

To use the Databricks integration, you need to set the following environment variables:

```bash
export DATABRICKS_HOST="https://your-workspace.cloud.databricks.com"
export DATABRICKS_TOKEN="your-personal-access-token"
export DATABRICKS_WORKSPACE_ID="your-workspace-id"  # Required for multi-workspace environments
export DATABRICKS_DEFAULT_SCHEMA="default"          # Optional, defaults to "default"
export DATABRICKS_DEFAULT_CATALOG="hive_metastore"  # Optional, defaults to "hive_metastore"
```

### Obtaining a Databricks Token

1. Log in to your Databricks workspace
2. Click on your user profile icon in the top-right corner
3. Select "User Settings"
4. Go to the "Access Tokens" tab
5. Click "Generate New Token"
6. Provide a name and expiration for your token
7. Copy the generated token (you won't be able to see it again)

## CLI Commands

### Catalog Commands

```bash
# List catalogs in Databricks
nessi catalog list --type databricks

# List schemas in a catalog
nessi catalog schemas --type databricks --catalog main

# List tables in a specific catalog and schema
nessi catalog tables --type databricks --database main.default

# Get details of a specific table
nessi catalog describe --type databricks --database main.default --table customers
```

### Delta Lake Commands

```bash
# Work with Delta Lake tables
nessi tables info /path/to/delta/table

# Run quality checks on a Delta Lake table
nessi quality check /path/to/delta/table

# Generate quality metrics for a Delta Lake table
nessi quality metrics /path/to/delta/table
```

## Examples

### Example 1: Exploring Databricks Unity Catalog

```bash
# List all catalogs
nessi catalog list --type databricks

# List all schemas in the 'main' catalog
nessi catalog schemas --type databricks --catalog main

# List all tables in the 'main.default' schema
nessi catalog tables --type databricks --database main.default

# Get details of the 'customers' table
nessi catalog describe --type databricks --database main.default --table customers
```

### Example 2: Running Quality Checks on Delta Tables

```bash
# Run quality checks on a Delta Lake table
nessi quality check --catalog databricks --database main.default --table customers

# Generate a quality report in HTML format
nessi report --catalog databricks --database main.default --table customers --format html --output ./reports/quality_report.html

# Generate a quality report in PDF format
nessi report --catalog databricks --database main.default --table customers --format pdf --output ./reports/quality_report.pdf

# Export quality metrics as JSON for further processing
nessi report --catalog databricks --database main.default --table customers --format json --output ./reports/quality_metrics.json
```

## Error Handling

Nessi provides robust error handling for Databricks integration:

### Authentication Errors

```bash
# Error when token is invalid or expired
nessi catalog list --type databricks
# Output: Error: Authentication failed: Invalid or expired Databricks token
```

### Connection Errors

```bash
# Error when Databricks host is unreachable
nessi catalog list --type databricks
# Output: Error: Connection failed: Could not connect to Databricks host
```

### Resource Not Found Errors

```bash
# Error when catalog doesn't exist
nessi catalog schemas --type databricks --catalog nonexistent_catalog
# Output: Error: Resource not found: Catalog 'nonexistent_catalog' does not exist

# Error when table doesn't exist
nessi catalog describe --type databricks --database main.default --table nonexistent_table
# Output: Error: Resource not found: Table 'main.default.nonexistent_table' does not exist
```

### Workspace ID Validation

```bash
# Error when workspace ID is not provided in a multi-workspace environment
# DATABRICKS_WORKSPACE_ID not set
nessi catalog list --type databricks
# Output: Error: Configuration error: Workspace ID is required for multi-workspace environments
```

### Rate Limiting Errors

```bash
# Error when too many requests are made
nessi catalog list --type databricks
# Output: Error: Rate limit exceeded: Too many requests to Databricks API, please try again later
```

### Server Errors

```bash
# Error when Databricks server returns an error
nessi catalog list --type databricks
# Output: Error: Server error: Databricks API returned an error (500)
```

### Error Handling Best Practices

1. **Always set required environment variables** before running commands
2. **Validate resource existence** before performing operations
3. **Check authentication** regularly as tokens expire
4. **Implement retry logic** for transient errors
5. **Monitor rate limits** to avoid throttling

### Debugging Databricks Integration Issues

```bash
# Enable verbose logging for detailed error information
NESSI_LOG_LEVEL=debug nessi catalog list --type databricks

# Test connection to Databricks
nessi test connection --type databricks
```

## Limitations

- The current implementation supports read-only operations for Unity Catalog metadata
- Table updates are limited to metadata changes
- Some advanced Databricks features like table ACLs are not yet supported

## Troubleshooting

### Common Issues

1. **Authentication Errors**:
   - Ensure your Databricks token is valid and has not expired
   - Check that you have the correct permissions in Databricks

2. **Connection Issues**:
   - Verify that your DATABRICKS_HOST is correct
   - Ensure your network allows connections to Databricks

3. **Missing Tables**:
   - Confirm that you have access to the specified catalog and schema
   - Check that the table exists in the specified location
