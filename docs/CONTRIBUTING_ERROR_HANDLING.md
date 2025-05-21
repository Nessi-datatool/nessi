# Contributing to Error Handling in Nessi

This guide provides best practices and standards for implementing error handling in the Nessi project. Following these guidelines ensures a consistent and user-friendly error handling experience across the codebase.

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

### Go Implementation

1. **Use structured errors**: Define error types with appropriate fields.

```go
type PathError struct {
    Path string
    Err  error
}

func (e *PathError) Error() string {
    return fmt.Sprintf("Error: Invalid path: %s does not exist or is not accessible", e.Path)
}
```

2. **Wrap errors with context**: Use `fmt.Errorf` with `%w` to wrap errors with context.

```go
if err != nil {
    return fmt.Errorf("failed to read Delta table: %w", err)
}
```

3. **Error handling in functions**: Return meaningful errors from functions.

```go
func ReadDeltaTable(path string) ([]Row, error) {
    if !fileExists(path) {
        return nil, &PathError{Path: path, Err: errors.New("path does not exist")}
    }
    
    if !isDeltaTable(path) {
        return nil, fmt.Errorf("Error: Not a Delta table: %s is not a Delta Lake table", path)
    }
    
    // Implementation...
}
```

## Testing Error Handling

1. **Test error scenarios**: Write tests for error scenarios, not just the happy path.

```go
func TestReadDeltaTableNonExistentPath(t *testing.T) {
    _, err := ReadDeltaTable("/nonexistent/path")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "Invalid path")
    assert.Contains(t, err.Error(), "/nonexistent/path")
}
```

2. **Use the test script**: Use the `test_error_handling.sh` script to verify error handling.

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
- [ ] Error handling follows the standard format and categories
- [ ] Tests cover error scenarios
- [ ] Documentation is updated
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
