# Pull Request

## Description

<!-- Describe the changes you've made -->

## Type of Change

<!-- Mark the appropriate option with an [x] -->

- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Test update
- [ ] Error handling improvement
- [ ] Performance optimization
- [ ] Other (please describe):

## Related Issues

<!-- Link to any related issues here -->

## Error Handling Considerations

<!-- If your PR involves error handling, please answer these questions -->

- Does this PR add or modify error messages? Yes/No
- If yes, have you used the NessiError type with appropriate error codes? Yes/No
- Have you registered error suggestions for any new error codes? Yes/No
- Have you implemented interactive resolution for resolvable errors (if applicable)? Yes/No
- Have you integrated with the error telemetry system? Yes/No
- Have you added tests for error scenarios, suggestions, and telemetry? Yes/No
- Have you updated ERROR_HANDLING.md and CONTRIBUTING_ERROR_HANDLING.md? Yes/No

<!-- For significant error handling changes, consider using the specialized error_handling.md template -->

## Testing

<!-- Describe the tests you've run to verify your changes -->

## Checklist

- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published in downstream modules
- [ ] Error handling is comprehensive and follows project standards
