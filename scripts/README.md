# Nessi.dev Scripts

This directory contains utility scripts for the Nessi.dev project. These scripts help with development, testing, and verification of various features.

## Testing Scripts

### run_tests.sh

A comprehensive test runner with multiple modes and options:

```bash
./scripts/run_tests.sh [options]
```

Options:
- `--fast` - Run only fast tests (< 5 seconds)
- `--short` - Run only specific short test files
- `--all` - Run all tests including long-running ones
- `--package PKG` - Specify a package to test (e.g., ./pkg/dbt/...)
- `--no-cache` - Disable test caching (adds -count=1)
- `--files FILES` - Run specific test files (comma-separated)
- `--help` - Show help message

Examples:
```bash
# Run fast tests only
./scripts/run_tests.sh --fast

# Run all tests for a specific package
./scripts/run_tests.sh --all --package="./pkg/freshness"

# Run specific test files
./scripts/run_tests.sh --short --files="./pkg/monitoring/export_short_test.go,./pkg/monitoring/dashboard/alerts_short_test.go"
```

### debug_long_tests.sh

A script specifically for debugging long-running tests with verbose output, timeouts, and race detection:

```bash
./scripts/debug_long_tests.sh
```

This script:
- Sets environment variables for debugging
- Creates log files for test output
- Runs specific long-running packages with verbose output
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

#### Test Modes

- **Standard Tests**: Run tests that complete in under 30 seconds (default)
- **Fast Tests**: Run only tests that complete in under 5 seconds, skipping known long-running packages
- **Short Tests**: Run only specific short test files with minimal timeout settings
- **All Tests**: Run all tests including long-running ones

#### Options

- `--package PKG`: Specify a package to test (e.g., `./pkg/dbt/...`)
- `--no-cache`: Disable test caching (adds `-count=1`)
- `--files FILES`: Run specific test files (comma-separated)

#### Usage Examples

```bash
# Run standard tests (default)
./scripts/run_tests.sh

# Run only fast tests
./scripts/run_tests.sh --fast

# Run specific short test files
./scripts/run_tests.sh --short

# Run specific short test files with custom file list
./scripts/run_tests.sh --short --files="./pkg/dbt/alert_test.go,./pkg/dbt/config_test.go"

# Run tests for a specific package
./scripts/run_tests.sh --package="./pkg/dbt/..."

# Run tests without caching
./scripts/run_tests.sh --no-cache

# Run all tests including long-running ones
./scripts/run_tests.sh --all

# Show help
./scripts/run_tests.sh --help
```

This script replaces all previous test scripts, including:
- `run_fast_tests.sh`
- `run_only_fast_tests.sh`
- `run_short_tests.sh`
- `fast_tests.sh`
- `fix_tests.sh`

All test-related functionality has been consolidated into this single script for easier maintenance and usage.

### debug_long_tests.sh

A specialized script for debugging long-running tests in the monitoring packages. This script:

- Runs tests with verbose output
- Enables race detection
- Saves detailed logs to the `logs` directory
- Uses the unified test runner with appropriate options

Usage:
```bash
./scripts/debug_long_tests.sh
```

This script is particularly useful when investigating timeouts or race conditions in the long-running tests.

## Test Data

Test data files have been moved to the `/test_data` directory for better organization.
