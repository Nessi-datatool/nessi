# Nessi.dev Test Strategy

This document outlines the testing strategy for Nessi.dev, focusing on improving test performance and developer experience.

## Test Categories

Tests are categorized based on their scope:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test interaction between components
3. **End-to-End Tests**: Test complete workflows

## Component Test Coverage

### DBT Plugin

The dbt plugin has comprehensive test coverage for the following components:

1. **AlertManager**: Tests for creating alert managers, sending alerts, and formatting failed rules
2. **Selection System**: Tests for applying modifiers, upstream/downstream selection, and parent/child model relationships
3. **Validator Functions**: Tests for finding dbt project paths, validating models, and handling edge cases
4. **Profiler**: Tests for setting profile types, detecting failures, and handling edge cases
5. **Artifacts**: Tests for generating different types of artifacts

## Running Tests

We provide a unified test script (`scripts/run_tests.sh`) that supports various options:

```bash
# Run all tests
./scripts/run_tests.sh

# Run tests for a specific package
./scripts/run_tests.sh --package="./pkg/dbt/..."

# Run specific test files
./scripts/run_tests.sh --files="path/to/file1_test.go,path/to/file2_test.go"

# Run tests without caching
./scripts/run_tests.sh --no-cache
```

All test scripts are located in the `scripts/` directory, and test data files are in the `test_data/` directory.

## CI/CD Integration

For CI/CD pipelines:
- PR validation: Run unit and integration tests
- Nightly builds: Run all tests including end-to-end tests
- Release validation: Run all tests with additional validation

## Test Optimization Techniques

1. **Isolation**: Tests should be isolated and not depend on external resources
2. **Mocking**: Use mocks for external dependencies
3. **Timeouts**: All tests should have appropriate timeouts
4. **Resource Cleanup**: Ensure proper cleanup of resources
5. **Parallel Execution**: Tests should be designed to run in parallel when possible
   - Use Go's built-in `t.Parallel()` for parallel test execution
   - Do not mix `t.Parallel()` with `testutil.RunInParallel(t)` in the same test
   - Be aware of potential race conditions when running tests in parallel

## Implementation Details

The test optimization includes:
- Standardized test approach without fast/long test distinctions
- Reduced sleep times in cleanup routines
- Proper resource management with `t.Cleanup()`
- Test timeouts to prevent hanging
- Environment variables to control test behavior
- Standardized test helpers for consistent execution
- Improved parallel test execution with clear warnings about potential conflicts

## Environment Variables

- `NESSI_TEST_TIMEOUT_MS`: Set timeout for tests in milliseconds
- `NESSI_TEST_DEBUG`: Set to "true" for verbose debug output

## Build Tags

- `integration`: Only run integration tests
- `e2e`: Only run end-to-end tests
