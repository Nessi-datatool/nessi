# Contributing to Nessi

Thank you for your interest in contributing to Nessi! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to foster an inclusive and respectful community.

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report:

1. Check the [issue tracker](https://github.com/nessi-dev/nessi/issues) to see if the issue has already been reported
2. Update your copy of Nessi to the latest version to see if the issue persists

When submitting a bug report, please include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior and what actually happened
- Nessi version, OS, and relevant environment details
- Any relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

1. Use a clear and descriptive title
2. Provide a detailed description of the suggested enhancement
3. Explain why this enhancement would be useful to Nessi users
4. Include examples of how the feature would work

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:

- `good-first-issue`: Issues suitable for newcomers
- `help-wanted`: Issues that need assistance
- `documentation`: Improvements or additions to documentation

### Pull Requests

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Make your changes
4. Run tests to ensure your changes don't break existing functionality
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.18 or higher
- Python 3.8 or higher (for Python integrations)
- Git

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/nessi-dev/nessi.git
   cd nessi
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

### Coding Standards

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Add comments to explain non-obvious code behavior
- Write tests for new functionality

## Testing

We aim for high test coverage. Please include tests with your contributions:

- Unit tests for individual functions and methods
- Integration tests for feature workflows
- End-to-end tests for critical user journeys

Run tests with:

```bash
# Run all tests
make test

# Run specific tests
go test ./pkg/specific-package

# Run tests with coverage
make test-coverage
```

## Documentation

Documentation is crucial for Nessi. Please update or add documentation for your changes:

- Update README.md if necessary
- Add or update documentation in the docs/ directory
- Include code comments for public APIs
- Update examples if relevant

## Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Reference issues and pull requests in the description
- Keep the first line under 72 characters
- Consider using the following format:
  ```
  category: Brief description

  Longer detailed description if necessary.

  Fixes #123
  ```

## Release Process

1. Maintainers will periodically create releases from the main branch
2. Version numbers follow [Semantic Versioning](https://semver.org/)
3. Release notes will be generated based on merged pull requests

## Community

- Join our [Discussions](https://github.com/nessi-dev/nessi/discussions) to ask questions
- Follow our [Twitter](https://twitter.com/nessidev) for updates
- Participate in community calls (schedule in the README)

## License

By contributing to Nessi, you agree that your contributions will be licensed under the project's [Apache License 2.0](LICENSE).
