# Contributing to Error Handling in Nessi

This guide provides best practices and standards for implementing error handling in the Nessi project. Following these guidelines ensures a consistent and user-friendly error handling experience across the codebase.

## Error Handling Components

The Nessi error handling system consists of several components that work together:

1. **Error Codes System** - Standardized error codes with structured error types
2. **Interactive Error Resolution** - Guided resolution for common errors
3. **Error Suggestions** - Contextual suggestions for fixing errors
4. **Error Telemetry** - Anonymous collection of error statistics
5. **Retry Mechanism** - Automatic retry for transient failures

## Error Handling Principles

1. **Be specific and actionable**: Error messages should clearly describe what went wrong and suggest how to fix it.
2. **Include context**: Include relevant context in error messages (e.g., file paths, resource names).
3. **Categorize errors**: Use appropriate error types to categorize different error scenarios.
4. **Fail early**: Validate inputs and preconditions early to provide better error messages.
5. **Don't expose sensitive information**: Ensure error messages don't contain sensitive data like credentials.

## Error Message Format

Follow this format for error messages:

```
Error: <Category>: <Specific message with context> [<Optional suggestion>]
```

Examples:
```
Error: Invalid path: /path/to/file does not exist or is not accessible
Error: Authentication failed: Invalid or expired Databricks token
Error: Schema validation failed: Field 'age' expected type INT but got STRING
```

## Error Categories

Use these standard error categories:

1. **Invalid Path**: For file or directory path issues
2. **Authentication**: For credential and authentication issues
3. **Connection**: For network and connectivity issues
4. **Resource Not Found**: For missing resources (tables, files, etc.)
5. **Schema Validation**: For data schema issues
6. **Corrupted Data**: For corrupted files or data
7. **Configuration**: For configuration issues
8. **Rate Limit**: For API rate limiting issues
9. **Server Error**: For remote server errors
10. **Permission Denied**: For access control issues

## Implementation Guidelines

### Using the Error Codes System

1. **Use NessiError**: Always use the structured NessiError type for errors.

```go
import "github.com/nessi-dev/nessi/pkg/common"

// Create a new error with error code
err := common.NewError(common.ErrInvalidPath, fmt.Sprintf("Path %s does not exist", path))

// Add details if needed
err.Details = "Additional context about the error"

// Add suggestions
err.Suggestions = []string{"Check if the path exists", "Verify you have the correct permissions"}
```

2. **Use helper functions**: Use the provided helper functions for common error types.

```go
// For path errors
err := common.NewPathError(path)

// For configuration errors
err := common.NewConfigError("database.host", "invalid hostname")

// For authentication errors
err := common.NewAuthenticationError("Invalid token")
```

3. **Register error suggestions**: Add suggestions for error codes.

```go
common.RegisterErrorSuggestion(common.ErrorSuggestion{
    ErrorCode:   common.ErrInvalidPath,
    Description: "The specified path does not exist",
    Solution:    "Check that the path exists and you have permissions",
    DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#path-errors",
})
```

4. **Create resolvable errors**: For errors that can be resolved interactively.

```go
// Create a resolvable path error
err := common.NewResolvablePathError(path)

// Create a resolvable configuration error
err := common.NewResolvableConfigError("database.host", "invalid-host")
```

## Testing Error Handling

1. **Test error scenarios**: Write tests for error scenarios, not just the happy path.

```go
func TestReadDeltaTableNonExistentPath(t *testing.T) {
    _, err := ReadDeltaTable("/nonexistent/path")
    
    // Check if it's a NessiError
    var nessiErr *common.NessiError
    assert.True(t, errors.As(err, &nessiErr))
    
    // Check error code
    assert.Equal(t, common.ErrInvalidPath, nessiErr.Code)
    
    // Check message content
    assert.Contains(t, nessiErr.Message, "/nonexistent/path")
}
```

2. **Test error suggestions**: Test that appropriate suggestions are provided.

```go
func TestErrorSuggestions(t *testing.T) {
    err := common.NewError(common.ErrInvalidPath, "test error")
    suggestions := common.GetSuggestionsForError(err)
    assert.NotEmpty(t, suggestions)
}
```

3. **Test error telemetry**: Verify telemetry recording works.

```go
func TestErrorTelemetry(t *testing.T) {
    // Create telemetry with test config
    config := common.ErrorTelemetryConfig{
        Enabled: true,
        StoragePath: "test_telemetry.json",
    }
    telemetry := common.NewErrorTelemetry(config)
    
    // Record an error
    err := common.NewError(common.ErrInvalidPath, "test error")
    telemetry.RecordError(err)
    
    // Check stats
    stats := telemetry.GetErrorStats()
    assert.Equal(t, 1, stats["total_errors"])
}
```

4. **Use the test command**: Use the `test-error` command to test error handling.

```bash
# Test a specific error code
nessi test-error N101

# Test with telemetry
nessi test-error N201 --telemetry
```

5. **Use the test script**: Use the `test_error_handling.sh` script to verify error handling.

```bash
./scripts/test_error_handling.sh
```

## Documenting Errors

1. **Update ERROR_HANDLING.md**: Document new error types and messages.
2. **Update component documentation**: Include error handling in component documentation.
3. **Add examples**: Provide examples of error scenarios and how to resolve them.

## Review Checklist

Before submitting a PR with error handling changes, ensure:

- [ ] Error messages are clear, specific, and actionable
- [ ] Error handling uses the NessiError type with appropriate error codes
- [ ] Error suggestions are registered for new error codes
- [ ] Interactive resolution is implemented for resolvable errors
- [ ] Error telemetry is properly integrated
- [ ] Tests cover error scenarios, suggestions, and telemetry
- [ ] Documentation is updated in ERROR_HANDLING.md
- [ ] No sensitive information is exposed in error messages
- [ ] Errors are properly propagated and wrapped with context

## Examples

### Delta Lake Error Handling

```go
func IsDeltaTable(path string) (bool, error) {
    // Check if path exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false, &PathError{Path: path, Err: err}
    }
    
    // Check if _delta_log directory exists
    deltaLogPath := filepath.Join(path, "_delta_log")
    if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
        return false, fmt.Errorf("Error: Not a Delta table: %s does not contain _delta_log directory", path)
    }
    
    return true, nil
}
```

### Databricks Error Handling

```go
func (c *DatabricksClient) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
    if workspaceID == "" {
        return nil, fmt.Errorf("Error: Configuration error: Workspace ID is required for multi-workspace environments")
    }
    
    resp, err := c.httpClient.Get(fmt.Sprintf("%s/api/2.0/unity-catalog/catalogs?workspaceId=%s", c.host, workspaceID))
    if err != nil {
        return nil, fmt.Errorf("Error: Connection failed: Could not connect to Databricks host: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 401 || resp.StatusCode == 403 {
        return nil, fmt.Errorf("Error: Authentication failed: Invalid or expired Databricks token")
    }
    
    if resp.StatusCode == 404 {
        return nil, fmt.Errorf("Error: Resource not found: Catalogs endpoint not found")
    }
    
    if resp.StatusCode == 429 {
        return nil, fmt.Errorf("Error: Rate limit exceeded: Too many requests to Databricks API, please try again later")
    }
    
    if resp.StatusCode >= 500 {
        return nil, fmt.Errorf("Error: Server error: Databricks API returned an error (%d)", resp.StatusCode)
    }
    
    // Process response...
}
```
