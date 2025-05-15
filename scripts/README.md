# Nessi.dev Scripts

This directory contains utility scripts for the Nessi.dev project.

## Testing Scripts

### run_tests.sh

A unified test runner script that supports different test modes:

- **Fast Tests**: Run only tests that complete in under 5 seconds
- **Standard Tests**: Run tests that complete in under 30 seconds (default)
- **All Tests**: Run all tests including long-running ones

Usage:
```bash
# Run standard tests (default)
./scripts/run_tests.sh

# Run only fast tests
./scripts/run_tests.sh --fast

# Run all tests including long-running ones
./scripts/run_tests.sh --all

# Show help
./scripts/run_tests.sh --help
```

This script replaces the older `run_fast_tests.sh`, `run_only_fast_tests.sh`, and `fix_tests.sh` scripts.

## Test Data

Test data files have been moved to the `/test_data` directory for better organization.
