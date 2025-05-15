# Nessi.dev Test Data

This directory contains test data files used by the Nessi.dev test suite.

## Files

### test_data.parquet

This is a Parquet file used for testing the Delta Lake and Parquet file handling capabilities of Nessi.dev. It contains sample data that is used by various tests to verify the correct functioning of data loading, schema inference, and data quality validation features.

### rules/test_rules.yaml

This YAML file contains sample validation rules used for testing the rule engine in Nessi.dev. It includes examples of regex and enum validation rules that are used by the test suite to verify the correct functioning of the rule validation system.

### rca/sample_anomaly.json

This JSON file contains a sample anomaly record used for testing the Root Cause Analysis (RCA) feature. It includes metadata about a detected anomaly such as timestamp, metric, threshold, and severity.

### rca/schema_changes.json

This JSON file contains sample schema change records used for testing the RCA feature's ability to detect and analyze schema changes as potential root causes for anomalies.

### rca/data_quality_issues.json

This JSON file contains sample data quality issues used for testing the RCA feature's ability to correlate anomalies with data quality rule failures in both the target table and upstream data sources.

### rca/dashboard_test_data.json

This JSON file contains sample RCA results used for testing the RCA dashboard integration. It includes multiple analysis results with various root causes, confidence scores, and affected tables for testing the dashboard's visualization capabilities.

### rca/insights_test_data.json

This JSON file contains sample RCA insights data used for testing the dashboard's insights feature. It includes aggregated statistics about root cause distributions, affected tables, and common root causes for testing the charts and visualizations.

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
