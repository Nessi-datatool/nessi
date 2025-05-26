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

- Go 1.21 or higher
- Python 3.8 or higher (for Python integrations)
- Git
- Docker and Docker Compose (for containerized development)

### Local Development

Nessi provides multiple ways to set up your development environment:

#### Standard Setup

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

#### Docker-based Development (Recommended)

We provide a comprehensive Docker-based development environment that includes hot-reloading and all necessary dependencies:

1. Start the development environment:
   ```bash
   docker-compose --profile dev up
   ```

2. Access the Nessi development server at http://localhost:8081

3. Make changes to the code and see them automatically reflected in the running application

The development environment includes:
- Hot-reloading with Air
- PostgreSQL database for testing
- MinIO for S3-compatible storage testing
- Grafana for visualization testing

For more details, see [DEVELOPMENT_ENVIRONMENT.md](docs/DEVELOPMENT_ENVIRONMENT.md).

### Coding Standards

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Add comments to explain non-obvious code behavior
- Write tests for new functionality

## Testing

We aim for high test coverage (minimum 70%). Please include tests with your contributions:

- Unit tests for individual functions and methods
- Integration tests for feature workflows
- End-to-end tests for critical user journeys
- Error handling tests for failure scenarios

### Running Tests Locally

```bash
# Run all tests
make test

# Run only unit tests (faster)
go test -short ./...

# Run specific tests
go test ./pkg/specific-package

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run integration tests
go test -tags=integration ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Using the Docker Testing Environment

For integration tests that require external services:

```bash
# Start the test environment
docker-compose --profile test up -d

# Run integration tests
go test -tags=integration ./...

# Stop the test environment
docker-compose --profile test down
```

### CI/CD Pipeline

All pull requests are automatically tested by our CI pipeline, which runs:

1. Unit and integration tests
2. Code coverage analysis (must meet 70% threshold)
3. Security scanning with multiple tools
4. Linting and code style checks
5. Performance benchmarks

The CI pipeline results will be visible in your pull request. All checks must pass before a PR can be merged.

For more details on our CI/CD workflows, see [CI_CD_WORKFLOWS.md](docs/CI_CD_WORKFLOWS.md).

## Documentation

Documentation is crucial for Nessi. Please update or add documentation for your changes:

- Update README.md if necessary
- Add or update documentation in the docs/ directory
- Include code comments for public APIs
- Update examples if relevant

We have an automated documentation workflow that generates API documentation from code comments. To ensure your code is properly documented:

1. Add comprehensive comments to all exported functions, types, and methods
2. Follow the [Go documentation conventions](https://golang.org/doc/comment)
3. Include examples where appropriate

The documentation workflow will automatically create a PR to update the generated documentation when changes are merged to the main branch.

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

Nessi uses an automated release process through GitHub Actions:

1. Version numbers follow [Semantic Versioning](https://semver.org/)
2. To create a new release, maintainers tag the main branch with a version tag (e.g., `v1.2.3`)
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```
3. The CD workflow automatically:
   - Builds binaries for multiple platforms (Linux, macOS, Windows) and architectures (amd64, arm64)
   - Creates Docker images and pushes them to GitHub Container Registry
   - Generates a changelog from commit messages
   - Creates a GitHub release with all artifacts

Pre-release versions can be created using suffixes:
- Alpha: `v1.2.3-alpha.1`
- Beta: `v1.2.3-beta.1`
- Release Candidate: `v1.2.3-rc.1`

For more details on the release process, see [CI_CD_WORKFLOWS.md](docs/CI_CD_WORKFLOWS.md).

## Community

- Join our [Discussions](https://github.com/nessi-dev/nessi/discussions) to ask questions
- Follow our [Twitter](https://twitter.com/nessidev) for updates
- Participate in community calls (schedule in the README)

## License

By contributing to Nessi, you agree that your contributions will be licensed under the project's [Apache License 2.0](LICENSE).
