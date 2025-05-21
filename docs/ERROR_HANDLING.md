# Error Handling in Nessi

## Overview

This document describes the comprehensive error handling implementation in the Nessi project, focusing on the Delta Lake and Databricks components. Proper error handling is critical for providing a good user experience, especially in a CLI-only application where clear error messages help users diagnose and fix issues quickly.

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

## Future Improvements

1. **Error codes** - Add numeric error codes for easier reference
2. **Interactive error resolution** - Guide users through fixing common errors
3. **Extended logging** - Provide more detailed logs for debugging complex issues
4. **Retry mechanisms** - Automatically retry operations for transient errors
5. **Error telemetry** - Collect anonymous error statistics to improve the product
