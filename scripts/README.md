# Nessi.dev Scripts

This directory contains utility scripts for the Nessi.dev project.

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
