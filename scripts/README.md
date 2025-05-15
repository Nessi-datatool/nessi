# Nessi.dev Scripts

This directory contains utility scripts for the Nessi.dev project.

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

## Test Data

Test data files have been moved to the `/test_data` directory for better organization.
