# Error Handling PR Template

This template is for PRs that involve error handling improvements or changes.

## Description

<!-- Describe the error handling changes you've made -->

## Error Types Added/Modified

<!-- List the error types that have been added or modified -->

- [ ] Invalid Path Errors
- [ ] Authentication Errors
- [ ] Connection Errors
- [ ] Resource Not Found Errors
- [ ] Schema Validation Errors
- [ ] Corrupted Data Errors
- [ ] Configuration Errors
- [ ] Rate Limiting Errors
- [ ] Server Errors
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
- [ ] Integration tests for error handling
- [ ] Error recovery tests (if applicable)
- [ ] Manual testing with the test_error_handling.sh script

## Documentation Updates

<!-- List the documentation files you've updated -->

- [ ] Updated ERROR_HANDLING.md
- [ ] Updated relevant component documentation
- [ ] Added examples to CLI documentation
- [ ] Updated error message references

## Checklist

- [ ] Error messages are clear and actionable
- [ ] Error handling is consistent with existing patterns
- [ ] Errors include context about what went wrong
- [ ] Errors suggest how to fix the issue when possible
- [ ] Error handling doesn't expose sensitive information
- [ ] Tests verify both the happy path and error scenarios
- [ ] Documentation is updated to reflect the changes
