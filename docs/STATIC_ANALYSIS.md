# Static Analysis Tools for Nessi

This document outlines the static analysis tools used in the Nessi project to ensure code quality and catch potential issues early in the development process.

## Golangci-lint

We use [golangci-lint](https://golangci-lint.run/) as our primary static analysis tool for Go code. It combines multiple linters into a single tool, making it easy to enforce coding standards and catch common mistakes.

### Installation

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Usage

```bash
# Run linters on the entire project
golangci-lint run ./...

# Run linters on specific packages or files
golangci-lint run ./pkg/... ./cmd/...
golangci-lint run ./examples/plugins/badge_plugin/main.go
```

### Enabled Linters

The following linters are enabled in our configuration:

#### Default Linters
- **errcheck**: Checks for unchecked errors
- **gosimple**: Suggests code simplifications
- **govet**: Reports suspicious constructs
- **ineffassign**: Detects ineffectual assignments
- **staticcheck**: Comprehensive static analyzer
- **typecheck**: Type-checks Go code
- **unused**: Checks for unused constants, variables, functions, and types
- **gofmt**: Checks if code was gofmt-ed

#### Additional Linters
- **gocyclo**: Detects complex functions
- **dupl**: Detects code duplication
- **goconst**: Finds repeated strings that could be constants
- **goimports**: Checks import formatting
- **gosec**: Performs security checks
- **misspell**: Finds commonly misspelled words
- **revive**: Drop-in replacement for golint
- **unconvert**: Removes unnecessary type conversions
- **unparam**: Finds unused function parameters
- **whitespace**: Checks for trailing whitespace

#### Viral Growth Feature Specific Linters
- **exportloopref**: Checks for pointers to enclosing loop variables
- **prealloc**: Finds slice declarations that could use make()
- **bodyclose**: Checks for unclosed response bodies

## Running Linters in CI

We have configured GitHub Actions to run linters automatically on every pull request. This ensures that code quality is maintained throughout the development process.

## Pre-commit Hooks

We recommend setting up pre-commit hooks to run linters before committing code. This helps catch issues early and ensures that only high-quality code is committed to the repository.

```bash
# Install pre-commit
pip install pre-commit

# Set up pre-commit hooks
pre-commit install
```

## Fixing Linting Issues

Many linting issues can be fixed automatically using the `--fix` flag:

```bash
golangci-lint run --fix ./...
```

For more complex issues, manual intervention may be required. Refer to the specific linter's documentation for guidance on fixing issues.

## Ignoring Linting Issues

In rare cases, it may be necessary to ignore linting issues. This should be done sparingly and with clear justification. Use the following comment to ignore a specific linting issue:

```go
// nolint: lintername
```

For example:

```go
// nolint: errcheck
f.Close()
```

## Conclusion

Static analysis tools are an essential part of our quality assurance process. They help us maintain high code quality, catch potential issues early, and ensure consistency throughout the codebase.
