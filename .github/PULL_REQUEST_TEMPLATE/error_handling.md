# Error Handling PR Template

This template is for PRs that involve error handling improvements or changes.

## Description

<!-- Describe the error handling changes you've made -->

## Error Components Modified

<!-- Check the error handling components that have been modified -->

- [ ] Error Codes System (structured error types, error codes)
- [ ] Interactive Error Resolution (resolvable errors)
- [ ] Error Suggestions (contextual suggestions)
- [ ] Error Telemetry (error statistics collection)
- [ ] Retry Mechanism (automatic retry for transient errors)

## Error Types Added/Modified

<!-- List the error types that have been added or modified -->

- [ ] Path Errors (N1XX)
- [ ] Configuration Errors (N2XX)
- [ ] Authentication Errors (N3XX)
- [ ] Connection Errors (N4XX)
- [ ] Delta Lake Errors (N5XX)
- [ ] Databricks Errors (N6XX)
- [ ] Schema Errors (N7XX)
- [ ] Validation Errors (N8XX)
- [ ] Internal Errors (N9XX)
- [ ] Other (please specify): 

## Error Messages

<!-- List the new or modified error messages -->

```
# Example:
# Error: Authentication failed: Invalid or expired Databricks token
```

## Tests Added

<!-- Describe the tests you've added to verify the error handling -->

- [ ] Unit tests for each error scenario
- [ ] Tests for error suggestions
- [ ] Tests for error telemetry
- [ ] Tests for interactive resolution
- [ ] Integration tests for error handling
- [ ] Error recovery tests (if applicable)
- [ ] Testing with the test-error command
- [ ] Manual testing with the test_error_handling.sh script

## Documentation Updates

<!-- List the documentation files you've updated -->

- [ ] Updated ERROR_HANDLING.md
- [ ] Updated CONTRIBUTING_ERROR_HANDLING.md
- [ ] Updated relevant component documentation
- [ ] Added examples to CLI documentation
- [ ] Updated error message references
- [ ] Added error code documentation

## Checklist

- [ ] Error messages are clear and actionable
- [ ] Error handling uses the NessiError type with appropriate error codes
- [ ] Error suggestions are registered for new error codes
- [ ] Interactive resolution is implemented for resolvable errors (if applicable)
- [ ] Error telemetry is properly integrated
- [ ] Retry mechanism is used for transient errors (if applicable)
- [ ] Errors include context about what went wrong
- [ ] Error handling doesn't expose sensitive information
- [ ] Tests verify both the happy path and error scenarios
- [ ] Documentation is updated to reflect the changes
