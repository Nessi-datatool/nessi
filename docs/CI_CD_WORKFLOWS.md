# CI/CD Workflows Documentation

This document describes the Continuous Integration (CI) and Continuous Deployment (CD) workflows implemented for the Nessi project.

## Overview

Nessi uses GitHub Actions for automating testing, building, and releasing the software. The workflows are designed to ensure code quality, security, and reliable releases. The project employs a comprehensive set of automated workflows to maintain high standards of code quality, security, and documentation.

## CI Workflow

The CI workflow runs on every pull request to the `main` branch and on every push to the `main` branch. It performs a series of checks to ensure code quality and functionality.

### Workflow Jobs

#### 1. Test

The test job runs unit and integration tests to ensure the code functions correctly.

- **Unit Tests**: Runs all tests with the `-short` flag to exclude long-running tests
- **Integration Tests**: Runs tests with the `integration` tag
- **Coverage Check**: Ensures test coverage meets the minimum threshold (currently 70%)
- **Coverage Reports**: Generates both text and HTML coverage reports

#### 2. Security Scan

The security scan job runs various security tools to identify potential vulnerabilities.

- **Gosec**: Scans for security issues in the Go code
- **Nancy**: Checks dependencies for known vulnerabilities
- **Govulncheck**: Checks for vulnerabilities in the Go code and dependencies
- **License Compliance**: Checks licenses of all dependencies

#### 3. Databricks Integration Test

This job specifically tests the Databricks and Delta Lake integration.

- **Integration Tests**: Tests with the `integration` tag for Databricks and Delta Lake
- **Mock Tests**: Tests using mock implementations
- **Error Handling Tests**: Tests specifically for error handling scenarios

#### 4. Lint

The lint job checks code style and formatting.

- **Golangci-lint**: Runs a comprehensive set of linters
- **Gofmt**: Ensures consistent code formatting

#### 5. Benchmark

The benchmark job runs performance benchmarks.

- **Go Benchmarks**: Runs all benchmark tests and measures performance
- **Benchmark Reports**: Generates reports with benchmark results

#### 6. Build

The build job compiles the code for multiple platforms.

- **Cross-Platform Builds**: Builds for Linux, macOS, and Windows
- **Multi-Architecture**: Builds for amd64 and arm64 architectures
- **Version Information**: Embeds version, build date, and git information

## CD Workflow (Release Pipeline)

The CD workflow runs when a tag with the pattern `v*` is pushed to the repository. It builds and releases the software.

### Workflow Jobs

#### 1. Build Release Binaries

This job builds the release binaries for multiple platforms.

- **Cross-Platform Builds**: Builds for Linux, macOS, and Windows
- **Multi-Architecture**: Builds for amd64 and arm64 architectures
- **Archives**: Creates `.tar.gz` archives for Linux and macOS, and `.zip` archives for Windows

#### 2. Docker

This job builds and pushes Docker images to GitHub Container Registry.

- **Multi-Architecture Images**: Builds for amd64 and arm64 architectures
- **Tags**: Tags images with semantic versioning (e.g., `v1.2.3`, `v1.2`, `v1`, `latest`)
- **Metadata**: Includes build information in image metadata

#### 3. Release

This job creates a GitHub release with all the built artifacts.

- **Changelog Generation**: Automatically generates a changelog from commit messages
- **Release Notes**: Creates comprehensive release notes
- **Asset Upload**: Uploads all built binaries as release assets

## How to Use

### Triggering CI

The CI workflow is automatically triggered on:
- Pull requests to the `main` branch
- Pushes to the `main` branch

No manual action is required.

### Creating a Release

To create a new release:

1. Ensure all tests pass on the `main` branch
2. Create and push a tag with the pattern `v*` (e.g., `v1.2.3`)
   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```
3. The release workflow will automatically build and publish the release

### Release Versioning

Nessi follows semantic versioning (SemVer):
- **Major version** (`v1.0.0`): Incompatible API changes
- **Minor version** (`v0.1.0`): Backwards-compatible functionality
- **Patch version** (`v0.0.1`): Backwards-compatible bug fixes

Pre-release versions can be tagged with suffixes:
- Alpha: `v1.0.0-alpha.1`
- Beta: `v1.0.0-beta.1`
- Release Candidate: `v1.0.0-rc.1`

## Configuration

### Environment Variables

The workflows use the following environment variables:

- `GO_VERSION`: The Go version to use (currently 1.21)
- `COVERAGE_THRESHOLD`: The minimum test coverage percentage (currently 70%)

### Secrets

The workflows use the following GitHub secrets:

- `GITHUB_TOKEN`: Automatically provided by GitHub, used for authentication

## Troubleshooting

### Common Issues

#### Failed Tests

If tests fail in the CI workflow:
1. Check the test logs in the GitHub Actions interface
2. Run the tests locally to reproduce the issue
3. Fix the failing tests and push the changes

#### Failed Builds

If builds fail in the CI workflow:
1. Check the build logs in the GitHub Actions interface
2. Ensure the code compiles locally
3. Fix any build issues and push the changes

#### Failed Release

If the release workflow fails:
1. Check the workflow logs in the GitHub Actions interface
2. Ensure the tag follows the correct pattern (`v*`)
3. Ensure all CI checks pass on the `main` branch
4. Fix any issues and create a new tag

## Contributing

When contributing to the Nessi project, please ensure:

1. All tests pass locally before creating a pull request
2. Code follows the project's style guidelines
3. New features include appropriate tests
4. Documentation is updated as needed

The CI workflow will automatically check your contribution for issues.

## Additional Workflows

In addition to the main CI and CD workflows, Nessi employs several specialized GitHub Actions workflows to automate various aspects of project maintenance:

### Dependency Updates

The dependency updates workflow automatically keeps project dependencies up-to-date.

- **Schedule**: Runs every Monday at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Actions**:
  - Updates Go dependencies to their latest versions
  - Creates a pull request with the changes
  - Labels the PR with "dependencies" and "automated"

### Security Scan

The security scan workflow performs comprehensive security analysis of the codebase.

- **Schedule**: Runs every Sunday at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Tools**:
  - Gosec for Go security scanning
  - Nancy for dependency vulnerability scanning
  - Govulncheck for Go vulnerability scanning
- **Actions**:
  - Uploads SARIF reports to GitHub Security tab
  - Creates an issue if high or critical vulnerabilities are found
  - Uploads all security reports as artifacts

### CodeQL Analysis

The CodeQL analysis workflow provides deep security and code quality analysis.

- **Triggers**: Runs on pushes to main, pull requests to main, and every Wednesday at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Features**:
  - Analyzes Go code for security vulnerabilities and quality issues
  - Uses GitHub's CodeQL engine for advanced static analysis
  - Results appear in the GitHub Security tab

### Stale Issue Management

The stale issue management workflow helps keep the issue tracker clean and organized.

- **Schedule**: Runs every day at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Actions**:
  - Marks issues as stale after 60 days of inactivity
  - Closes stale issues after an additional 14 days of inactivity
  - Exempts issues with specific labels (no-stale, security, bug, enhancement, documentation)
  - Applies similar rules to pull requests

### Documentation Updates

The documentation updates workflow automatically keeps documentation in sync with code changes.

- **Triggers**: Runs on pushes to main that modify Go files, docs, or README.md
- **Manual Trigger**: Available via workflow_dispatch
- **Actions**:
  - Generates API documentation using gomarkdoc
  - Updates CLI documentation using the built-in doc command
  - Creates a pull request with the changes
  - Labels the PR with "documentation" and "automated"

### Performance Benchmarks

The performance benchmarks workflow tracks performance over time and detects regressions.

- **Triggers**: Runs on pushes to main and every Tuesday at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Actions**:
  - Runs benchmarks on the current commit
  - Runs benchmarks on the previous commit
  - Compares results using benchstat
  - Creates an issue if significant performance regressions are detected
  - Uploads benchmark results as artifacts

### Code Quality Analysis

The code quality analysis workflow provides comprehensive code quality metrics and analysis.

- **Triggers**: Runs on pushes to main, pull requests to main, and every Thursday at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Components**:
  - **SonarCloud Analysis**: Provides detailed code quality metrics, technical debt tracking, and security analysis
  - **GolangCI-Lint with Reviewdog**: Provides inline code review comments on pull requests
  - **Go Formatting Check**: Ensures all code follows Go formatting standards
  - **Cyclomatic Complexity Analysis**: Identifies overly complex functions that need refactoring
- **Actions**:
  - Analyzes code quality and provides detailed metrics
  - Comments on pull requests with code quality issues
  - Fails the workflow if critical issues are found
  - Uploads code quality reports as artifacts

### Docker Image Testing

The Docker image testing workflow ensures that the Docker image works correctly and meets all requirements.

- **Triggers**: Runs on pushes to main and pull requests that modify Docker-related files
- **Manual Trigger**: Available via workflow_dispatch
- **Tests**:
  - Builds the Docker image and verifies it works correctly
  - Tests basic functionality by running commands in the container
  - Tests Docker Compose configuration
  - Verifies multi-architecture build capability
  - Scans for security vulnerabilities using Trivy
- **Actions**:
  - Builds the Docker image
  - Runs various tests to verify functionality
  - Scans for vulnerabilities
  - Fails the workflow if critical issues are found

### API Testing

The API testing workflow ensures that the API endpoints work correctly and maintain backward compatibility.

- **Triggers**: Runs on pushes to main and pull requests that modify API-related files, and every Friday at midnight
- **Manual Trigger**: Available via workflow_dispatch
- **Components**:
  - **Integration Tests**: Tests API functionality with real dependencies
  - **Contract Tests**: Verifies API backward compatibility
  - **Performance Tests**: Measures API performance under load
- **Actions**:
  - Starts a test environment with PostgreSQL and MinIO
  - Builds and runs the API server
  - Runs REST, gRPC, and GraphQL API tests
  - Runs contract tests to verify backward compatibility
  - Runs performance tests using k6
  - Generates and uploads test reports

## Future Improvements

Planned improvements to the CI/CD workflows:

1. End-to-end testing with real-world scenarios
2. Deployment to package managers (Homebrew, apt, etc.)
3. Automated changelog generation with categorized changes
4. Integration with code coverage services
5. Automated release notes generation
