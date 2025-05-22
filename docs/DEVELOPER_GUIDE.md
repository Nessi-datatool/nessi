# Nessi Developer Guide

This guide provides comprehensive information for developers who want to extend, customize, or contribute to the Nessi project.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Core Components](#core-components)
- [Building from Source](#building-from-source)
- [Testing](#testing)
- [Creating Plugins](#creating-plugins)
- [API Development](#api-development)
- [Contributing Guidelines](#contributing-guidelines)
- [Code Style](#code-style)
- [Documentation](#documentation)
- [Release Process](#release-process)

## Development Environment Setup

### Prerequisites

- **Go**: Version 1.21 or later
- **Python**: Version 3.8 or later (for Python API and testing)
- **Docker**: Latest version (for containerized development and testing)
- **Git**: Latest version
- **Make**: Latest version

### Setting Up Your Environment

1. Clone the repository:

   ```bash
   git clone https://github.com/nessi-dev/nessi.git
   cd nessi
   ```

2. Install development dependencies:

   ```bash
   # Go dependencies
   go mod download

   # Python dependencies (for testing and Python API)
   pip install -r requirements-dev.txt
   ```

3. Set up pre-commit hooks:

   ```bash
   pre-commit install
   ```

4. Configure your IDE:

   - **VS Code**: Use the provided `.vscode/settings.json` for consistent settings
   - **GoLand/IntelliJ**: Import the provided code style settings

## Project Structure

The Nessi project is organized as follows:

```
nessi/
├── cmd/                    # Command-line applications
│   ├── nessi/              # Main CLI application
│   └── ...                 # Other executables
├── pkg/                    # Library packages
│   ├── datalake/           # Data lake operations
│   ├── quality/            # Data quality functionality
│   ├── schema/             # Schema management
│   ├── metrics/            # Metrics collection and analysis
│   ├── report/             # Report generation
│   └── ...                 # Other packages
├── internal/               # Internal packages (not for external use)
│   ├── cli/                # CLI implementation
│   ├── api/                # API implementation
│   └── ...                 # Other internal packages
├── api/                    # API definitions
│   ├── rest/               # REST API
│   ├── grpc/               # gRPC API
│   └── graphql/            # GraphQL API
├── web/                    # Web UI
├── scripts/                # Build and utility scripts
├── docs/                   # Documentation
├── examples/               # Example code and configurations
└── test/                   # Test data and integration tests
```

## Core Components

Nessi is built around several core components:

### Delta Lake Handler

The Delta Lake handler (`pkg/datalake/delta.go`) provides functionality for working with Delta Lake tables, including:

- Detecting Delta Lake tables
- Reading data from Delta Lake tables
- Writing data to Delta Lake tables
- Extracting metadata
- Time travel capabilities

### Quality Engine

The quality engine (`pkg/quality/`) provides data quality validation functionality:

- Rule-based validation
- Quality scoring
- Custom rule definitions
- Quality metrics calculation

### Schema Manager

The schema manager (`pkg/schema/`) handles schema operations:

- Schema validation
- Schema evolution
- Schema comparison
- Schema registry integration

### Metrics Collector

The metrics collector (`pkg/metrics/`) handles metrics collection and analysis:

- Quality metrics
- Performance metrics
- Freshness metrics
- Trend analysis

### Report Generator

The report generator (`pkg/report/`) creates reports in various formats:

- HTML reports with visualizations
- PDF reports
- JSON/CSV exports
- Custom templates

## Building from Source

### Building the CLI

```bash
# Build for your current platform
go build -o nessi ./cmd/nessi

# Build for a specific platform
GOOS=linux GOARCH=amd64 go build -o nessi-linux-amd64 ./cmd/nessi
GOOS=darwin GOARCH=amd64 go build -o nessi-darwin-amd64 ./cmd/nessi
GOOS=windows GOARCH=amd64 go build -o nessi-windows-amd64.exe ./cmd/nessi
```

### Building with Make

```bash
# Build for your current platform
make build

# Build for all supported platforms
make build-all

# Build and install
make install
```

### Building Docker Image

```bash
# Build Docker image
make docker-build

# Run Docker image
docker run --rm nessi/nessi:latest --help
```

## Testing

Nessi uses a comprehensive testing approach:

### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run tests for a specific package
go test ./pkg/quality/...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run integration tests
make integration-test

# Run specific integration test
go test -tags=integration ./test/integration/delta_test.go
```

### End-to-End Tests

```bash
# Run end-to-end tests
make e2e-test
```

### Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Run benchmarks for a specific package
go test -bench=. ./pkg/datalake/...
```

## Creating Plugins

Nessi supports plugins for extending functionality. Plugins can be written in Go and are loaded dynamically at runtime.

### Plugin Types

- **Format Handlers**: Add support for new data formats
- **Quality Rules**: Add custom quality validation rules
- **Report Templates**: Add custom report templates
- **Metrics Collectors**: Add custom metrics collection
- **Command Extensions**: Add new CLI commands

### Plugin Structure

A basic plugin structure:

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/plugin"
)

// Plugin is the main plugin struct
type Plugin struct{}

// Init initializes the plugin
func (p *Plugin) Init() error {
    // Initialization code
    return nil
}

// Name returns the plugin name
func (p *Plugin) Name() string {
    return "my-plugin"
}

// Version returns the plugin version
func (p *Plugin) Version() string {
    return "1.0.0"
}

// Register registers the plugin components
func (p *Plugin) Register(registry plugin.Registry) error {
    // Register components
    return nil
}

// Plugin entry point
var Plugin plugin.Plugin = &Plugin{}
```

### Building a Plugin

```bash
go build -buildmode=plugin -o my-plugin.so my-plugin.go
```

### Installing a Plugin

```bash
mkdir -p ~/.nessi/plugins
cp my-plugin.so ~/.nessi/plugins/
```

### Using a Plugin

```bash
nessi --plugin my-plugin command
```

## API Development

Nessi provides multiple APIs for integration:

### REST API

The REST API is defined in OpenAPI format in `api/rest/openapi.yaml`. To extend the REST API:

1. Update the OpenAPI specification
2. Implement the new endpoints in `internal/api/rest/`
3. Add tests for the new endpoints

### gRPC API

The gRPC API is defined in Protocol Buffers format in `api/grpc/nessi.proto`. To extend the gRPC API:

1. Update the Protocol Buffers definition
2. Generate the gRPC code: `make generate-grpc`
3. Implement the new services in `internal/api/grpc/`
4. Add tests for the new services

### GraphQL API

The GraphQL API is defined in GraphQL schema format in `api/graphql/schema.graphql`. To extend the GraphQL API:

1. Update the GraphQL schema
2. Generate the GraphQL code: `make generate-graphql`
3. Implement the new resolvers in `internal/api/graphql/`
4. Add tests for the new resolvers

## Contributing Guidelines

We welcome contributions to Nessi! Please follow these guidelines:

1. **Fork the repository** and create a branch for your feature or bug fix
2. **Write tests** for your changes
3. **Ensure all tests pass**: `make test`
4. **Update documentation** as needed
5. **Follow the code style** guidelines
6. **Submit a pull request** with a clear description of your changes

For more details, see [CONTRIBUTING.md](../CONTRIBUTING.md).

## Code Style

Nessi follows these code style guidelines:

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `golint` and `go vet` to check for issues
- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines

### Python Code Style

- Follow [PEP 8](https://www.python.org/dev/peps/pep-0008/)
- Use `black` for formatting
- Use `isort` for import sorting
- Use `flake8` for linting

### Commit Messages

- Use the imperative mood ("Add feature" not "Added feature")
- First line should be 50 characters or less
- Reference issues and pull requests where appropriate
- Use the body to explain what and why, not how

## Documentation

Documentation is a crucial part of Nessi. When contributing, please:

1. **Update the relevant documentation** for your changes
2. **Add examples** where appropriate
3. **Use clear, concise language**
4. **Follow the documentation style guide**

### Documentation Structure

- `docs/`: Main documentation directory
  - `cli/`: CLI documentation
  - `api/`: API documentation
  - `*.md`: General documentation files

### Building Documentation

```bash
# Generate CLI documentation
make docs-cli

# Generate API documentation
make docs-api

# Generate all documentation
make docs
```

## Release Process

Nessi follows a structured release process:

1. **Version Bump**: Update version in `version.go` and `CHANGELOG.md`
2. **Release Branch**: Create a release branch (`release-vX.Y.Z`)
3. **Testing**: Run all tests on the release branch
4. **Documentation**: Ensure documentation is up to date
5. **Release Commit**: Commit the version bump and changelog updates
6. **Tag**: Tag the release commit (`git tag vX.Y.Z`)
7. **Build**: Build the release artifacts
8. **Publish**: Publish the release artifacts
9. **Announce**: Announce the release

### Release Checklist

- [ ] Update version in `version.go`
- [ ] Update `CHANGELOG.md`
- [ ] Run all tests
- [ ] Build for all platforms
- [ ] Update documentation
- [ ] Create GitHub release
- [ ] Publish Docker image
- [ ] Publish to package repositories

## Advanced Topics

### Performance Optimization

Nessi is designed for performance when working with large datasets:

- **Parallelization**: Use goroutines for parallel processing
- **Memory Management**: Minimize memory usage with streaming operations
- **Caching**: Cache results where appropriate
- **Profiling**: Use Go's profiling tools to identify bottlenecks

### Security Considerations

When developing for Nessi, consider these security aspects:

- **Input Validation**: Validate all user input
- **Authentication**: Implement proper authentication for APIs
- **Authorization**: Check permissions for all operations
- **Secrets Management**: Never hardcode secrets
- **Dependency Security**: Regularly update dependencies

### Error Handling

Nessi uses a standardized error handling approach:

- **Error Codes**: Use error codes for categorization
- **Error Messages**: Provide clear, actionable error messages
- **Error Details**: Include relevant details for debugging
- **Error Recovery**: Implement recovery mechanisms where appropriate

For more information on error handling, see [ERROR_HANDLING.md](ERROR_HANDLING.md).

## Community and Support

- **GitHub Issues**: Report bugs and request features
- **GitHub Discussions**: Ask questions and discuss ideas
- **Slack Channel**: Join our community on Slack
- **Mailing List**: Subscribe to our mailing list for announcements

For more information, visit [nessi.dev](https://nessi.dev).
