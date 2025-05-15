# Nessi.dev Test Data

This directory contains test data files used by the Nessi.dev test suite.

## Files

### test_data.parquet

This is a Parquet file used for testing the Delta Lake and Parquet file handling capabilities of Nessi.dev. It contains sample data that is used by various tests to verify the correct functioning of data loading, schema inference, and data quality validation features.

### Test Rules

The `test_rules.yaml` file contains sample quality rules used for testing the quality rules functionality. These rules are loaded by the `TestYAMLLoader` function in the quality rules tests.

### rules/test_rules.yaml

This YAML file contains sample validation rules used for testing the rule engine in Nessi.dev. It includes examples of regex, enum, length, and date format validation rules that are used by the test suite to verify the correct functioning of the rule validation system.

### test_rules.yaml

This YAML file in the root test_data directory contains a simplified set of test rules used by the quality rules tests. It includes basic validation rules for testing the YAML loader functionality without requiring the full rule set.

### rca/sample_anomaly.json

This JSON file contains a sample anomaly record used for testing the Root Cause Analysis (RCA) feature. It includes metadata about a detected anomaly such as timestamp, metric, threshold, and severity.

### rca/schema_changes.json

This JSON file contains sample schema change records used for testing the RCA feature's ability to detect and analyze schema changes as potential root causes for anomalies.

### rca/data_quality_issues.json

This JSON file contains sample data quality issues used for testing the RCA feature's ability to correlate anomalies with data quality rule failures in both the target table and upstream data sources.

## Freshness SLA Testing

The Freshness SLA testing now uses an in-memory approach with the following components:

### BetterMockDeltaConnector

A mock implementation of the Delta connector interface that doesn't rely on file operations and can be configured with predefined responses. This connector allows tests to simulate Delta Lake tables with configurable last modified times without requiring actual table files.

### TestSLAManager

A simplified version of the SLAManager specifically designed for testing. It works with the BetterMockDeltaConnector to provide a complete in-memory testing environment for the Freshness SLA functionality.

This in-memory testing approach avoids potential deadlocks and timeouts caused by file I/O operations during testing, making the tests more reliable and faster to execute.

### rca/dashboard_test_data.json

This JSON file contains sample RCA results used for testing the RCA dashboard integration. It includes multiple analysis results with various root causes, confidence scores, and affected tables for testing the dashboard's visualization capabilities.

### rca/insights_test_data.json

This JSON file contains sample RCA insights data used for testing the dashboard's insights feature. It includes aggregated statistics about root cause distributions, affected tables, and common root causes for testing the charts and visualizations.

### freshness/sla_test_data.json

This JSON file contains sample SLA configurations used for testing the data freshness monitoring feature. It includes various tables with different expected update frequencies, warning and critical thresholds, and metadata for testing the SLA management functionality.

### freshness/status_test_data.json

This JSON file contains sample freshness status data used for testing the freshness dashboard. It includes current status information for tables with different freshness levels (up-to-date, warning, critical) to test the dashboard's visualization and alerting capabilities.

## Test Data Structure

The test data is organized into the following directories:

- `dq/`: Contains data quality test data
- `metrics/`: Contains metrics test data
- `rca/`: Contains root cause analysis test data
- `freshness/`: Contains data freshness and SLA test data

## Usage

Test files in this directory are referenced by the test suite and should not be modified or deleted unless you are also updating the corresponding tests.

When adding new test data files, please document them in this README.md file with a brief description of their purpose and the tests that use them.

## Running Tests with This Data

The test data in this directory is used by various test scripts. To run tests that use this data, use our unified test script:

```bash
# Run all tests including those that use test data
./scripts/run_tests.sh --all

# Run specific tests that use test data
./scripts/run_tests.sh --package="./pkg/delta/..."

# Run only fast tests (some may use test data)
./scripts/run_tests.sh --fast
```

For more information about our testing strategy and available options, see:
- `./scripts/run_tests.sh --help`
- `./docs/test_strategy.md`
