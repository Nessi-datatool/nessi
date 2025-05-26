# End-to-End Testing in Nessi

This document outlines the end-to-end (E2E) testing strategy for the Nessi CLI application, mapping each feature to its corresponding test case.

## Overview

E2E tests validate the entire application workflow from start to finish, ensuring that all components work together correctly. These tests interact with the Nessi CLI just as a user would, verifying that commands produce the expected outputs and handle errors appropriately.

## Test Structure

The E2E tests are organized into three main categories:

1. **User Journeys**: Tests that simulate common user workflows
2. **Integrations**: Tests that verify integration with external systems and formats
3. **Stress Tests**: Tests that verify error handling and edge cases

## Feature Coverage

### Core Features (Community Edition)

| Feature | Test File | Test Function | Status |
|---------|-----------|---------------|--------|
| Config Management | `user-journeys/new_user_onboarding_test.go` | `TestNewUserOnboarding` | ✅ |
| Delta Lake Table Operations | `integrations/delta_lake_test.go` | `TestDeltaLakeIntegration` | ✅ |
| Data Quality Checks | `user-journeys/new_user_onboarding_test.go` | `TestNewUserOnboarding` | ✅ |
| Enhanced Report Generation | `integrations/enhanced_report_test.go` | `TestEnhancedReportGeneration` | ✅ |
| CLI Visualization | `integrations/cli_visualization_test.go` | `TestCLIVisualization` | ✅ |
| Schema Management | `integrations/delta_lake_test.go` | `TestDeltaLakeIntegration` | ✅ |
| Comprehensive Error Handling | `stress-tests/error_handling_comprehensive_test.go` | `TestErrorHandlingComprehensive` | ✅ |
| Performance Benchmarking | `stress-tests/performance_benchmark_test.go` | `TestPerformanceBenchmark` | ✅ |
| License Management | `user-journeys/new_user_onboarding_test.go` | `TestNewUserOnboarding` | ✅ |

### Premium Features (Pro Edition)

| Feature | Test File | Test Function | Status |
|---------|-----------|---------------|--------|
| Time Travel | `integrations/delta_lake_test.go` | `TestDeltaLakeIntegration` | ✅ |
| Databricks Integration | `integrations/databricks_client_test.go` | `TestDatabricksClient` | ✅ |
| Databricks Error Handling | `integrations/databricks_error_handling_test.go` | `TestDatabricksErrorHandling` | ✅ |
| Cloud Storage Integration | `integrations/cloud_storage_test.go` | `TestCloudStorageIntegration` | ✅ |
| Workflow Orchestration | `integrations/workflow_test.go` | `TestWorkflowOrchestration` | ✅ |

## Running the Tests

The E2E tests can be run using the provided scripts:

```bash
# Run all tests with detailed output
./e2e-tests/run_detailed_tests.sh

# Run all tests with minimal output
./e2e-tests/run_tests.sh

# Enable tests (remove skip directives)
./e2e-tests/run_detailed_tests.sh --enable
```

## Test Requirements

The E2E tests require:

1. A compiled Nessi binary in the PATH or in the project's `bin` directory
2. Sample test data (included in the repository)
3. Go testing tools

## Adding New Tests

When adding new features to Nessi, follow these guidelines to ensure proper E2E test coverage:

1. **Identify the Feature Category**: Determine if the feature is part of a user journey, an integration, or requires stress testing
2. **Create a Test File**: Add a new test file in the appropriate directory
3. **Implement Test Functions**: Create test functions that verify all aspects of the feature
4. **Update This Document**: Add the new feature and test to the feature coverage table

## Test Data

The E2E tests use sample data located in the `e2e-tests/testdata` directory:

- `delta_tables/`: Sample Delta Lake tables
- `schemas/`: Sample schema files
- `configs/`: Sample configuration files

## Troubleshooting

If the E2E tests fail, check the following:

1. **Binary Availability**: Ensure the Nessi binary is available in the PATH or in the project's `bin` directory
2. **Test Data**: Verify that the test data is properly set up
3. **Environment Variables**: Check if any required environment variables are missing
4. **Permissions**: Ensure the test scripts have execution permissions

## Recent Improvements

The following improvements have been made to the E2E testing framework:

1. **Enhanced Report Generation Tests**: Added tests for generating reports in different formats (HTML, JSON, Markdown, CSV) with various templates (quality, schema, freshness, performance).
2. **CLI Visualization Tests**: Added tests for schema tree visualization, quality heatmap visualization, ASCII bar chart, and progress visualization.
3. **Comprehensive Error Handling Tests**: Added tests for various error scenarios including non-existent tables, invalid command flags, missing required arguments, invalid file formats, permission errors, and more.
4. **Performance Benchmark Tests**: Added tests to measure performance of CLI operations under different data volumes and conditions.
5. **Databricks Integration Tests**: Added tests for the Databricks client functionality, including authentication, resource access, and catalog operations.

## Future Improvements

The following improvements are planned for the E2E testing framework:

1. **Parallel Test Execution**: Enable parallel execution of tests to reduce testing time
2. **Containerized Testing**: Add support for running tests in containers to ensure consistent environments
3. **Continuous Integration**: Integrate E2E tests into CI/CD pipelines
4. **Test Coverage Reporting**: Add tools to measure and report E2E test coverage
5. **Automated Test Generation**: Explore tools for generating E2E tests based on code changes
