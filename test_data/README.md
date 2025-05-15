# Nessi.dev Test Data

This directory contains test data files used by the Nessi.dev test suite.

## Files

### test_data.parquet

This is a Parquet file used for testing the Delta Lake and Parquet file handling capabilities of Nessi.dev. It contains sample data that is used by various tests to verify the correct functioning of data loading, schema inference, and data quality validation features.

## Usage

Test files in this directory are referenced by the test suite and should not be modified or deleted unless you are also updating the corresponding tests.

When adding new test data files, please document them in this README.md file with a brief description of their purpose and the tests that use them.
