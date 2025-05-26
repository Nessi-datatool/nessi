# Nessi End-to-End Tests

This directory contains end-to-end tests for the Nessi CLI application. These tests verify that the entire application works correctly from a user's perspective by simulating real user interactions with the CLI.

## Test Structure

The tests are organized into the following categories:

- **User Journeys**: Tests that simulate complete user workflows, such as onboarding a new user.
- **Integrations**: Tests that verify integration with external systems like Delta Lake and Databricks.
- **Stress Tests**: Tests that verify error handling and performance under various conditions.

## Key Test Files

### Integration Tests
- `enhanced_report_test.go`: Tests report generation in different formats (HTML, JSON, Markdown, CSV) with various templates.
- `cli_visualization_test.go`: Tests CLI visualization features like schema trees, quality heatmaps, and ASCII charts.
- `databricks_client_test.go`: Tests Databricks client functionality including authentication and catalog operations.
- `databricks_error_handling_test.go`: Tests error handling for Databricks integration.

### Stress Tests
- `error_handling_comprehensive_test.go`: Tests various error scenarios including non-existent tables, invalid flags, and permission errors.
- `performance_benchmark_test.go`: Tests performance of CLI operations under different data volumes and conditions.

## Running the Tests

To run all E2E tests:

```bash
cd /path/to/nessi
ENABLE_E2E_TESTS=true go test -v ./e2e-tests/...
```

To run a specific test category:

```bash
ENABLE_E2E_TESTS=true go test -v ./e2e-tests/scenarios/user-journeys/...
ENABLE_E2E_TESTS=true go test -v ./e2e-tests/scenarios/integrations/...
ENABLE_E2E_TESTS=true go test -v ./e2e-tests/scenarios/stress-tests/...
```

To run a specific test file:

```bash
ENABLE_E2E_TESTS=true go test -v ./e2e-tests/scenarios/integrations/enhanced_report_test.go
ENABLE_E2E_TESTS=true go test -v ./e2e-tests/scenarios/stress-tests/performance_benchmark_test.go
```

## Test Requirements

- The Nessi binary must be built and available in your PATH or in the project's bin directory.
- Sample data files must be available in the test_data directory.
- For integration tests, appropriate credentials and connection details must be configured.
- Environment variable `ENABLE_E2E_TESTS=true` must be set to run the tests.
- For email distribution tests, set `ENABLE_EMAIL_TESTS=true`.

## Test Data

The tests use sample data files located in the `test_data` directory. These files include:
- Sample Delta Lake tables
- Schema definition files
- Configuration templates

## Adding New Tests

When adding new E2E tests, follow these guidelines:
1. Place the test in the appropriate category directory
2. Use the `testing` package for Go tests
3. Implement proper setup and teardown to ensure tests are isolated
4. Document any special requirements or setup needed for the test
5. Use the `testutil` package for common testing utilities
6. Add environment variable checks to skip tests when not explicitly enabled
7. Update the documentation in `docs/E2E_TESTING.md` to reflect new tests

## Mock CLI Implementation

The E2E tests use a mock CLI implementation located in the `mock` directory. This mock CLI simulates the behavior of the actual Nessi CLI without requiring the full application to be built. The mock CLI supports:

- Delta Lake table operations
- Schema management
- Quality checks
- Report generation
- Error handling
- Databricks integration

When adding new features to the Nessi CLI, make sure to update the mock CLI implementation to support testing those features.
