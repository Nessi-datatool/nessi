# Nessi

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Python Version](https://img.shields.io/badge/python-3.8%2B-blue)](https://www.python.org/downloads/)
[![PyPI Version](https://img.shields.io/pypi/v/nessi.svg)](https://pypi.org/project/nessi/)
[![Tests](https://github.com/Nessi-datatool/nessi/actions/workflows/tests.yml/badge.svg)](https://github.com/Nessi-datatool/nessi/actions/workflows/tests.yml)
[![Code Coverage](https://codecov.io/gh/Nessi-datatool/nessi/branch/main/graph/badge.svg)](https://codecov.io/gh/Nessi-datatool/nessi)

A Python-based data processing and analysis tool built with PySpark and Delta Lake.

## Features

- Sample data generation with various data patterns
- Support for multiple file formats (Delta, Parquet, CSV)
- Table scanning and profiling capabilities
- Delta Lake integration for versioned data storage

## Prerequisites

- Python 3.8 or higher
- Java 8 or higher (required for PySpark)
- Git

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Nessi-datatool/nessi.git
cd nessi
```

### 2. Create and Activate Virtual Environment

```bash
# Create virtual environment
python -m venv venv

# Activate virtual environment
# On macOS/Linux:
source venv/bin/activate
# On Windows:
# venv\Scripts\activate
```

### 3. Install Dependencies

```bash
# Install required packages
pip install -r requirements.txt

# Install the package in development mode
pip install -e .
```

### 4. Configure Spark

The project uses PySpark with Delta Lake. The default configuration includes:

```python
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("Nessi") \
    .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0,io.delta:delta-storage:2.4.0") \
    .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
    .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
    .getOrCreate()
```

### 5. Set Environment Variables (Optional)

You can customize the behavior of the application using environment variables:

```bash
# Set test data directory
export TEST_DATA_DIR=/path/to/test/data

# Set Spark memory configuration
export SPARK_DRIVER_MEMORY=2g
export SPARK_EXECUTOR_MEMORY=2g
```

## Usage

### Sample Data Generation

```python
from src.sample_data_generator import generate_sample_data, save_as_delta

# Generate sample data with 1000 rows
data = generate_sample_data(num_rows=1000)

# Save as Delta table
save_as_delta(data, "path/to/output")

# Save as Parquet
save_as_parquet(data, "path/to/output.parquet")

# Save as CSV
save_as_csv(data, "path/to/output.csv")
```

### Table Scanning

```python
from src.scanner.table_scanner import TableScanner

# Create scanner instance
scanner = TableScanner()

# Scan a Delta table
delta_result = scanner.scan_table("path/to/delta/table")

# Scan a Parquet table
parquet_result = scanner.scan_table("path/to/parquet/table")

# Scan a CSV table
csv_result = scanner.scan_table("path/to/csv/table")

# Generate a report
report = scanner.generate_report(delta_result)
```

### Creating Test Tables

```python
from src.create_test_delta_table import create_test_table

# Create a test Delta table with default settings
table_path = create_test_table()

# Create a test Delta table with custom schema
custom_schema = {
    "id": "integer",
    "name": "string",
    "value": "double"
}
table_path = create_test_table(schema=custom_schema)

# Create a partitioned table
table_path = create_test_table(partition_by=["date"])
```

## Testing

Run the test suite:

```bash
python -m pytest tests/ -v
```

For test coverage report:

```bash
python -m pytest tests/ --cov=src --cov-report=html
```

## Project Structure

```
nessi/
├── data/              # Sample and test data
├── src/               # Source code
│   ├── scanner/       # Table scanning functionality
│   └── ...           # Other modules
├── tests/             # Test files
├── requirements.txt   # Python dependencies
├── setup.py          # Package setup
├── pyproject.toml    # Modern Python packaging config
├── API.md            # API documentation
├── CHANGELOG.md      # Version history
├── CONTRIBUTING.md   # Contribution guidelines
├── CODE_OF_CONDUCT.md # Community guidelines
├── NOTICE            # Attribution notices
└── README.md         # Project documentation
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## Code of Conduct

Please review our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE.txt](LICENSE.txt) file for details.

## Documentation

For detailed API documentation, see [API.md](API.md).

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and changes. 