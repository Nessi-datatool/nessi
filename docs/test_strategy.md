# Nessi.dev Test Strategy

This document outlines the testing strategy for Nessi.dev, focusing on improving test performance and developer experience.

## Test Categories

Tests are categorized into three tiers:

1. **Fast Tests**: Run in < 5 seconds, test core functionality only
2. **Standard Tests**: Run in < 30 seconds, test integration between components
3. **Long Tests**: May take minutes to run, test end-to-end functionality

## Component Test Coverage

### DBT Plugin

The dbt plugin has comprehensive test coverage for the following components:

1. **AlertManager**: Tests for creating alert managers, sending alerts, and formatting failed rules
2. **Selection System**: Tests for applying modifiers, upstream/downstream selection, and parent/child model relationships
3. **Validator Functions**: Tests for finding dbt project paths, validating models, and handling edge cases
4. **Profiler**: Tests for setting profile types, detecting failures, and handling edge cases
5. **Artifacts**: Tests for generating different types of artifacts

## Running Tests

We provide several scripts to run different test categories:

- `run_fast_tests.sh`: Runs only the fast tests (< 5 seconds)
- `run_standard_tests.sh`: Runs fast and standard tests (< 30 seconds)
- `run_all_tests.sh`: Runs all tests including long-running ones

## CI/CD Integration

For CI/CD pipelines:
- PR validation: Run fast and standard tests
- Nightly builds: Run all tests including long-running ones
- Release validation: Run all tests with additional validation

## Test Optimization Techniques

1. **Isolation**: Tests should be isolated and not depend on external resources
2. **Mocking**: Use mocks for external dependencies
3. **Timeouts**: All tests should have appropriate timeouts
4. **Resource Cleanup**: Ensure proper cleanup of resources
5. **Parallel Execution**: Tests should be designed to run in parallel when possible

## Implementation Details

The test optimization includes:
- Reduced sleep times in cleanup routines
- Proper resource management with `t.Cleanup()`
- Test timeouts to prevent hanging
- Environment variables to control test behavior
- Build tags to separate long and short tests

## Environment Variables

- `NESSI_SKIP_LONG_TESTS`: Set to "true" to skip long-running tests
- `NESSI_TEST_TIMEOUT_MS`: Set timeout for tests in milliseconds
- `NESSI_TEST_DEBUG`: Set to "true" for verbose debug output

## Build Tags

- `skiplong`: Skip long-running tests
- `integration`: Only run integration tests
- `e2e`: Only run end-to-end tests
