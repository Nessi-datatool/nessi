# Nessi.dev Scripts

This directory contains utility scripts for the Nessi.dev project. These scripts help with development, testing, and verification of various features.

## Testing Scripts

### run_tests.sh

A comprehensive test runner with multiple modes and options:

```bash
./scripts/run_tests.sh [options]
```

Options:
- `--package PKG` - Specify a package to test (e.g., ./pkg/dbt/...)
- `--no-cache` - Disable test caching (adds -count=1)
- `--files FILES` - Run specific test files (comma-separated)
- `--help` - Show help message

Examples:
```bash
# Run all tests
./scripts/run_tests.sh

# Run tests for a specific package
./scripts/run_tests.sh --package="./pkg/freshness"

# Run specific test files
./scripts/run_tests.sh --files="./pkg/monitoring/export_test.go,./pkg/monitoring/dashboard/alerts_test.go"
```

### debug_tests.sh

A script for debugging tests with verbose output, timeouts, and race detection:

```bash
./scripts/debug_tests.sh
```

This script:
- Sets environment variables for debugging
- Creates log files for test output
- Runs packages with verbose output
- Performs race detection on critical packages
- Saves all logs to the `./logs` directory

## Development Scripts

### run-dev.sh

A script for setting up and running the development environment:

```bash
./scripts/run-dev.sh
```

This script:
- Installs air for hot reloading (if not already installed)
- Creates necessary directories for data, logs, and certificates
- Generates self-signed certificates if they don't exist
- Starts the monitoring stack with docker-compose (if available)
- Sets development environment variables
- Runs the application with hot reloading

### generate_certs.sh

Generates self-signed SSL certificates for development:

```bash
./scripts/generate_certs.sh
```

## Security Verification Scripts

The security verification scripts have been organized into separate directories to avoid package conflicts:

### verify_security_pkg/verify_security.go

A simple script that verifies the basic SSL certificate generation functionality. It creates a temporary directory, generates SSL certificates, and verifies that the certificate files were created successfully.

### verify_security_features_pkg/verify_security_features.go

A comprehensive script that verifies all security features including:
- Authentication manager creation
- Certificate manager creation
- User management (creation, listing)
- Authentication (username/password, token validation)
- API key management (generation, authentication)
- SSL certificate generation

### test_security_pkg/test_security.go

A test script that runs a series of tests on the security package and reports detailed results. It tests:
- Auth manager creation
- Cert manager creation
- User creation
- Authentication
- Token validation
- API key generation
- API key authentication
- Admin user creation
- Role middleware
- SSL certificate generation

## Testing Scripts

### run_tests.sh

A comprehensive test runner script that consolidates all testing functionality. It supports multiple test modes, package selection, and other options.

#### Test Execution

The test runner executes tests with appropriate timeout settings and provides detailed output for debugging purposes.

#### Options

- `--package PKG`: Specify a package to test (e.g., `./pkg/dbt/...`)
- `--no-cache`: Disable test caching (adds `-count=1`)
- `--files FILES`: Run specific test files (comma-separated)

#### Usage Examples

```bash
# Run all tests
./scripts/run_tests.sh

# Run specific test files
./scripts/run_tests.sh --files="./pkg/dbt/alert_test.go,./pkg/dbt/config_test.go"

# Run tests for a specific package
./scripts/run_tests.sh --package="./pkg/dbt/..."

# Run tests without caching
./scripts/run_tests.sh --no-cache

# Show help
./scripts/run_tests.sh --help
```

This script provides a unified interface for running all tests in the project.

All test-related functionality has been consolidated into this single script for easier maintenance and usage.

### debug_tests.sh

A specialized script for debugging tests in the monitoring packages. This script:

- Runs tests with verbose output
- Enables race detection
- Saves detailed logs to the `logs` directory
- Uses the unified test runner with appropriate options

Usage:
```bash
./scripts/debug_tests.sh
```

This script is particularly useful when investigating timeouts or race conditions in tests.

## Test Data

Test data files have been moved to the `/test_data` directory for better organization.
