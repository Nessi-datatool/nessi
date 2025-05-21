# Nessi Error Codes Reference

This document provides a quick reference for all error codes used in Nessi. Each error code follows the format `NXXX` where the first digit indicates the error category.

## Error Code Categories

| Category | Code Range | Description |
|----------|------------|-------------|
| Path Errors | N1XX | Errors related to file paths and directories |
| Configuration Errors | N2XX | Errors related to configuration settings |
| Authentication Errors | N3XX | Errors related to authentication and credentials |
| Connection Errors | N4XX | Errors related to network connections |
| Delta Lake Errors | N5XX | Errors related to Delta Lake operations |
| Databricks Errors | N6XX | Errors related to Databricks integration |
| Schema Errors | N7XX | Errors related to schema validation and handling |
| Validation Errors | N8XX | Errors related to data validation |
| Internal Errors | N9XX | Internal errors and unexpected issues |

## Detailed Error Codes

### Path Errors (N1XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N101 | ErrInvalidPath | Path does not exist or is not accessible | - Path does not exist<br>- Incorrect path separators<br>- Insufficient permissions |
| N102 | ErrPermissionDenied | Permission denied for path | - Insufficient permissions<br>- File locked by another process |
| N103 | ErrNotADirectory | Path is not a directory | - Path points to a file, not a directory |
| N104 | ErrNotAFile | Path is not a file | - Path points to a directory, not a file |
| N105 | ErrPathTooLong | Path exceeds maximum length | - Path is too long for the operating system |

### Configuration Errors (N2XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N201 | ErrInvalidConfig | Invalid configuration value | - Syntax error in config file<br>- Invalid value type |
| N202 | ErrMissingConfig | Missing required configuration | - Required field not specified<br>- Config file not found |
| N203 | ErrInvalidFormat | Invalid configuration format | - JSON/YAML syntax error<br>- Incorrect format specification |
| N204 | ErrConfigConflict | Conflicting configuration values | - Mutually exclusive options specified |
| N205 | ErrUnsupportedConfig | Unsupported configuration | - Configuration option not supported in current context |

### Authentication Errors (N3XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N301 | ErrAuthenticationFailed | Authentication failed | - Invalid credentials<br>- Account locked or disabled |
| N302 | ErrInvalidToken | Invalid token | - Token malformed<br>- Token not recognized |
| N303 | ErrExpiredToken | Expired token | - Token has expired<br>- Token revoked |
| N304 | ErrInsufficientPermissions | Insufficient permissions | - Account lacks required permissions |
| N305 | ErrMissingCredentials | Missing credentials | - Required credentials not provided |

### Connection Errors (N4XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N401 | ErrConnectionFailed | Connection failed | - Server unavailable<br>- Network issues |
| N402 | ErrTimeout | Connection timeout | - Server overloaded<br>- Network latency |
| N403 | ErrNetworkError | Network error | - Network connectivity issues<br>- DNS resolution failure |
| N404 | ErrResourceNotFound | Resource not found | - URL or endpoint does not exist |
| N405 | ErrRateLimitExceeded | Rate limit exceeded | - Too many requests in a short period |

### Delta Lake Errors (N5XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N501 | ErrInvalidDeltaTable | Not a valid Delta Lake table | - Directory is not a Delta table<br>- Missing _delta_log directory |
| N502 | ErrCorruptedDeltaTable | Corrupted Delta Lake table | - Transaction log corrupted<br>- Incomplete write operation |
| N503 | ErrSchemaMismatch | Schema mismatch | - Schema evolution incompatible<br>- Column types don't match |
| N504 | ErrInvalidPartition | Invalid partition | - Partition column does not exist<br>- Partition value invalid |
| N505 | ErrTableLocked | Table locked | - Concurrent write operation<br>- Failed transaction not cleaned up |

### Databricks Errors (N6XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N601 | ErrDatabricksAPIError | Databricks API error | - API endpoint error<br>- Invalid API request |
| N602 | ErrWorkspaceNotFound | Workspace not found | - Workspace ID does not exist<br>- Workspace access revoked |
| N603 | ErrCatalogNotFound | Catalog not found | - Catalog does not exist<br>- Insufficient permissions |
| N604 | ErrSchemaNotFound | Schema not found | - Schema does not exist<br>- Schema access revoked |
| N605 | ErrTableNotFound | Table not found | - Table does not exist<br>- Table access revoked |

### Schema Errors (N7XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N701 | ErrInvalidSchema | Invalid schema | - Schema definition invalid<br>- Unsupported data type |
| N702 | ErrSchemaValidationFailed | Schema validation failed | - Data does not match schema<br>- Required fields missing |
| N703 | ErrIncompatibleSchema | Incompatible schema | - Schema evolution not possible<br>- Breaking changes |
| N704 | ErrSchemaInferenceError | Schema inference error | - Cannot infer schema from data<br>- Ambiguous types |
| N705 | ErrSchemaEvolutionError | Schema evolution error | - Incompatible schema changes<br>- Column type changes |

### Validation Errors (N8XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N801 | ErrValidationFailed | Validation failed | - Data does not meet validation criteria |
| N802 | ErrQualityCheckFailed | Quality check failed | - Data quality below threshold |
| N803 | ErrConstraintViolation | Constraint violation | - Unique constraint violated<br>- Foreign key constraint violated |
| N804 | ErrDataTypeMismatch | Data type mismatch | - Value cannot be converted to expected type |
| N805 | ErrInvalidFormat | Invalid data format | - Date/time format invalid<br>- Number format invalid |

### Internal Errors (N9XX)

| Code | Name | Description | Common Causes |
|------|------|-------------|---------------|
| N901 | ErrInternalError | Internal error | - Unexpected internal failure |
| N902 | ErrUnexpectedError | Unexpected error | - Unhandled exception |
| N903 | ErrNotImplemented | Feature not implemented | - Feature requested but not yet implemented |
| N904 | ErrDependencyError | Dependency error | - Required dependency missing or incompatible |
| N905 | ErrSystemResourceError | System resource error | - Out of memory<br>- Disk space exhausted |

## Troubleshooting

When you encounter an error code, you can:

1. **Get more information**: Use `nessi test-error <code> --list` to see details about the error code
2. **View suggestions**: Error messages include suggestions for resolving the issue
3. **Check telemetry**: Use `nessi telemetry report-errors` to see common errors
4. **Consult documentation**: See the [Error Handling Documentation](ERROR_HANDLING.md) for detailed guidance

## Testing Error Codes

You can test how different error codes are handled using the `test-error` command:

```bash
# List all available error codes
nessi test-error --list

# Generate a specific error
nessi test-error N101

# Test with custom details and suggestion
nessi test-error N201 --details "Custom details" --suggestion "Try this solution"
```

For more information on error handling in Nessi, see the [Error Handling Documentation](ERROR_HANDLING.md).
