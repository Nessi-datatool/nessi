# Nessi.dev Test Data

This directory contains test data files used by the Nessi.dev test suite.

## Files

### test_data.parquet

This is a Parquet file used for testing the Delta Lake and Parquet file handling capabilities of Nessi.dev. It contains sample data that is used by various tests to verify the correct functioning of data loading, schema inference, and data quality validation features.

### rules/test_rules.yaml

This YAML file contains sample validation rules used for testing the rule engine in Nessi.dev. It includes examples of regex and enum validation rules that are used by the test suite to verify the correct functioning of the rule validation system.

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
