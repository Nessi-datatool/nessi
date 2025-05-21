# Error Handling in Nessi

## Overview

This document describes the comprehensive error handling implementation in the Nessi project. Proper error handling is critical for providing a good user experience, especially in a CLI-only application where clear error messages help users diagnose and fix issues quickly.

The Nessi error handling system consists of several key components:

1. **Error Codes System** - Standardized error codes with structured error types
2. **Interactive Error Resolution** - Guided resolution for common errors
3. **Error Suggestions** - Contextual suggestions for fixing errors
4. **Error Telemetry** - Anonymous collection of error statistics
5. **Retry Mechanism** - Automatic retry for transient failures

## Delta Lake Error Handling

### Error Types

1. **Invalid Path Errors**
   - Occurs when a specified path doesn't exist or isn't accessible
   - Example: `Error: Invalid path: /nonexistent/path does not exist or is not accessible`

2. **Not a Delta Table Errors**
   - Occurs when a directory is not a valid Delta Lake table
   - Example: `Error: Not a Delta table: /path/to/regular/directory is not a Delta Lake table`

3. **Schema Validation Errors**
   - Occurs when data doesn't match the expected schema
   - Example: `Error: Schema validation failed: Field 'age' expected type INT but got STRING`

4. **Corrupted Table Errors**
   - Occurs when a Delta table's transaction log is corrupted
   - Example: `Error: Corrupted Delta table: Transaction log in /path/to/corrupted/table is corrupted`

### Implementation Details

The Delta Lake error handling is implemented in the following files:

- `pkg/datalake/delta.go`: Main implementation of Delta Lake operations
- `pkg/datalake/delta_error_handling_test.go`: Tests for Delta Lake error handling
- `pkg/datalake/delta_error_handling_integration_test.go`: Integration tests for Delta Lake error handling

## Databricks Integration Error Handling

### Error Types

1. **Authentication Errors**
   - Occurs when the Databricks token is invalid or expired
   - Example: `Error: Authentication failed: Invalid or expired Databricks token`

2. **Connection Errors**
   - Occurs when the Databricks host is unreachable
   - Example: `Error: Connection failed: Could not connect to Databricks host`

3. **Resource Not Found Errors**
   - Occurs when a requested resource (catalog, schema, table) doesn't exist
   - Example: `Error: Resource not found: Catalog 'nonexistent_catalog' does not exist`

4. **Workspace ID Validation Errors**
   - Occurs when the workspace ID is missing in a multi-workspace environment
   - Example: `Error: Configuration error: Workspace ID is required for multi-workspace environments`

5. **Rate Limiting Errors**
   - Occurs when too many requests are made to the Databricks API
   - Example: `Error: Rate limit exceeded: Too many requests to Databricks API, please try again later`

6. **Server Errors**
   - Occurs when the Databricks server returns an error
   - Example: `Error: Server error: Databricks API returned an error (500)`

### Implementation Details

The Databricks integration error handling is implemented in the following files:

- `pkg/catalog/databricks/databricks_catalog.go`: Main implementation of Databricks catalog operations
- `pkg/catalog/databricks/databricks_test_utils.go`: Mock client for testing error scenarios
- `pkg/catalog/databricks/enhanced_integration_test.go`: Integration tests for Databricks error handling
- `pkg/catalog/databricks/databricks_delta_integration_test.go`: Integration tests for Databricks and Delta Lake

## Testing Error Handling

A comprehensive test script is available to verify the error handling implementation:

```bash
# Run the error handling test script
./scripts/test_error_handling.sh
```

This script tests various error scenarios for both Delta Lake and Databricks integration, including:

- Invalid paths
- Non-Delta tables
- Corrupted Delta tables
- Authentication failures
- Connection issues
- Missing workspace IDs
- Resource not found errors
- Schema validation errors

## Best Practices

### For Users

1. **Check error messages carefully** - They provide specific guidance on what went wrong
2. **Verify environment variables** - Ensure all required variables are set correctly
3. **Validate resources exist** - Confirm catalogs, schemas, and tables exist before operations
4. **Use verbose logging** - Set `NESSI_LOG_LEVEL=debug` for detailed error information

### For Developers

1. **Provide specific error messages** - Include details about what went wrong and how to fix it
2. **Implement proper error types** - Use structured errors with appropriate context
3. **Add comprehensive tests** - Test all error scenarios to ensure proper handling
4. **Document error messages** - Ensure all possible error messages are documented
5. **Consider user experience** - Format error messages for readability in CLI output

## Error Codes System

All errors in Nessi have a unique error code in the format `NXXX`, where:

- `N` is the prefix for Nessi errors
- `X` is a digit representing the error category
- The first digit indicates the general category:
  - `1XX` - Path and file errors
  - `2XX` - Configuration errors
  - `3XX` - Authentication errors
  - `4XX` - Connection errors
  - `5XX` - Delta Lake errors
  - `6XX` - Databricks errors
  - `7XX` - Schema errors
  - `8XX` - Validation errors
  - `9XX` - Internal errors
  - `VXX` - Viral growth feature errors

Examples:
- `N101` - Invalid path error
- `N201` - Invalid configuration error
- `N301` - Authentication failed error
- `V101` - Share command error
- `V102` - Badge command error
- `V103` - Community command error

### NessiError Structure

All errors in Nessi are structured using the `NessiError` type, which includes:

- `Code` - The error code (e.g., `N101`)
- `Message` - A user-friendly error message
- `Details` - Additional details about the error (optional)
- `Suggestions` - List of suggestions to fix the error (optional)

## Interactive Error Resolution

For common errors, Nessi provides an interactive resolution mechanism that guides users through fixing the issue. This feature is enabled by default and can be controlled with the `--interactive` flag.

Example:
```
nessi schema show --path /nonexistent/path
❌ Error: Invalid path: /nonexistent/path does not exist or is not accessible

Would you like to create this directory? [y/N]: y
Creating directory /nonexistent/path...
✅ Directory created successfully!
```

### Resolvable Errors

The following error types support interactive resolution:

1. **Path Errors** - Offers to create missing directories
2. **Configuration Errors** - Helps set correct configuration values
3. **Delta Table Errors** - Assists with table creation or repair

## Error Suggestions System

Even when interactive resolution is not possible, Nessi provides contextual suggestions to help users fix errors. These suggestions include:

- Description of the problem
- Suggested solution
- Link to relevant documentation

Example:
```
❌ Error: Authentication failed: Invalid Databricks token

🔍 Suggested solutions:

📝 Authentication failed due to invalid credentials
🔧 Solution: Check your credentials and ensure they are correctly configured
📚 Documentation: https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#authentication-errors

📝 Token may have expired
🔧 Solution: Try regenerating your authentication token
```

## Error Telemetry

Nessi includes an error telemetry system that anonymously collects error statistics to help improve the product. This feature:

- Is enabled by default but can be disabled
- Collects only error codes and counts, no personal data
- Stores data locally, not sent to any server

### Telemetry Commands

```bash
# View telemetry status
nessi telemetry status

# Enable error telemetry
nessi telemetry enable-error

# Disable error telemetry
nessi telemetry disable-error

# View error statistics
nessi telemetry report-errors

# Export error statistics to a file
nessi telemetry export-errors /path/to/export.json
```

## Retry Mechanism

For transient errors (like network issues or rate limiting), Nessi implements an automatic retry mechanism with:

- Exponential backoff - Increasing delay between retries
- Jitter - Random variation in retry timing to prevent thundering herd problems
- Configurable policies - Customizable retry counts and delays

## Testing Error Handling

Nessi provides tools to test the error handling system:

### Test Error Command

```bash
# List all available error codes
nessi test-error --list

# Generate a specific error
nessi test-error N101

# Test with custom details and suggestion
nessi test-error N201 --details "Custom details" --suggestion "Try this solution"

# Test resolvable errors
nessi test-error N101 --resolvable

# Test error telemetry
nessi test-error N301 --telemetry

# Export telemetry data
nessi test-error N401 --telemetry --export /path/to/export.json
```

### Test Script

A comprehensive test script is available to verify all aspects of error handling:

```bash
# Run the error handling test script
./scripts/test_error_handling.sh
```
