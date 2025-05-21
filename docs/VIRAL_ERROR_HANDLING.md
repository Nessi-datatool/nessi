# Viral Growth Features Error Handling

This document describes the error handling implementation for the viral growth features in the Nessi project.

## Overview

The viral growth features in Nessi include robust error handling to provide clear, actionable error messages when operations fail. This enhances the user experience by helping users diagnose and fix issues quickly.

## Error Codes

Viral growth feature errors use the `VXX` format:

- `V101` - Share command errors
- `V102` - Badge command errors
- `V103` - Community command errors

## Share Command Errors

### Error Types

1. **Invalid Path Errors**
   - Occurs when a specified output path doesn't exist or isn't accessible
   - Error Code: `V101`
   - Example: `Error: output path directory does not exist: /nonexistent/path [Error Code: V101]`

2. **Missing Table Name**
   - Occurs when no table name is provided
   - Error Code: `V101`
   - Example: `Error: share command requires a table name [Error Code: V101]`

3. **Invalid Table**
   - Occurs when the specified table doesn't exist or isn't accessible
   - Error Code: `V101`
   - Example: `Error: table 'nonexistent_table' does not exist or is not accessible [Error Code: V101]`

### Handling Share Errors

When a share command error occurs, the system provides suggestions for resolving the issue:

```
❌ Sharing Error: output path directory does not exist: /nonexistent/path [Error Code: V101]

Details: The specified output directory does not exist

Suggestions:
1. Check if the table exists and is accessible
2. Verify that the social sharing plugin is installed correctly
3. Try using a different output format or path
4. Run 'nessi viral share --help' for usage information
```

## Badge Command Errors

### Error Types

1. **Invalid Format**
   - Occurs when an unsupported badge format is specified
   - Error Code: `V102`
   - Example: `Error: Invalid badge format: invalid [Error Code: V102]`

2. **Invalid Color**
   - Occurs when an unrecognized color is specified
   - Error Code: `V102`
   - Example: `Warning: Color 'invalid-color' may not be recognized. Using default color instead. [Error Code: V102]`

3. **Invalid Quality Score**
   - Occurs when a quality score outside the valid range (0-100) is specified
   - Error Code: `V102`
   - Example: `Error: Invalid quality score: 101. Score must be between 0 and 100 [Error Code: V102]`

### Handling Badge Errors

When a badge command error occurs, the system provides suggestions for resolving the issue:

```
❌ Badge Generation Error: Invalid badge format: invalid [Error Code: V102]

Details: Supported formats are: markdown, html, rst

Suggestions:
1. Check if the badge plugin is installed correctly
2. Try using a different badge format (markdown, html)
3. Verify that the color value is valid (hex or named color)
4. Run 'nessi viral badge --help' for usage information
```

## Community Command Errors

### Error Types

1. **Missing Feedback Text**
   - Occurs when no feedback text is provided
   - Error Code: `V103`
   - Example: `Error: Feedback text cannot be empty [Error Code: V103]`

2. **Invalid Experience Level**
   - Occurs when an invalid experience level is specified
   - Error Code: `V103`
   - Example: `Error: invalid experience level: invalid [Error Code: V103]`

3. **Network Errors**
   - Occurs when there are issues connecting to the community services
   - Error Code: `V103`
   - Example: `Error: Failed to connect to community service [Error Code: V103]`

### Handling Community Errors

When a community command error occurs, the system provides suggestions for resolving the issue:

```
❌ Community Engagement Error: Feedback text cannot be empty [Error Code: V103]

Details: The feedback text is required to submit feedback

Suggestions:
1. Check your internet connection
2. Verify that the community engagement plugin is installed correctly
3. Try again later if the issue persists
4. Run 'nessi viral community --help' for usage information
```

## Implementation Details

The viral growth features error handling is implemented in the following files:

- `cmd/nessi/cli/viral_error_handler.go`: Main implementation of viral error handling
- `cmd/nessi/cli/viral_cmd.go`: Integration with viral commands
- `pkg/common/resolvable_errors.go`: Definition of resolvable errors
- `pkg/errors/error_codes.go`: Definition of error codes

## Testing Error Handling

The viral growth features error handling can be tested using the test-error command:

```bash
# Test share command error
nessi test-error V101 --details "Custom details" --suggestion "Try this solution"

# Test badge command error
nessi test-error V102 --details "Custom details" --suggestion "Try this solution"

# Test community command error
nessi test-error V103 --details "Custom details" --suggestion "Try this solution"
```

## Best Practices

When implementing new viral growth features, follow these best practices for error handling:

1. **Use the appropriate error code** - Use the correct error code for the feature (V101, V102, V103)
2. **Provide clear error messages** - Make error messages clear and actionable
3. **Include details** - Provide additional details about the error
4. **Offer suggestions** - Suggest ways to resolve the issue
5. **Use the HandleViralError function** - Use the HandleViralError function to handle errors consistently

## Conclusion

The viral growth features error handling system in Nessi enhances the user experience by providing clear, actionable error messages when operations fail. By following the patterns and best practices described in this document, you can ensure that your viral growth features provide a consistent and user-friendly experience.
