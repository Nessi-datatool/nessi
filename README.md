# Nessi

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Python](https://img.shields.io/badge/Python-3.8%2B-blue)](https://www.python.org/downloads/)
[![PyPI](https://img.shields.io/pypi/v/nessi.svg)](https://pypi.org/project/nessi/)
[![Tests](https://github.com/Nessi-datatool/nessi/actions/workflows/tests.yml/badge.svg)](https://github.com/Nessi-datatool/nessi/actions/workflows/tests.yml)
[![Codecov](https://codecov.io/gh/Nessi-datatool/nessi/branch/main/graph/badge.svg)](https://codecov.io/gh/Nessi-datatool/nessi)

Nessi is a Python-based data processing and analysis tool built with PySpark and Delta Lake. It provides a comprehensive suite of features for data generation, table scanning, and analysis.

## Features

- **Sample Data Generation**: Create realistic test data with various data types and distributions
- **Multiple Format Support**: Work with Delta Lake, Parquet, and CSV formats
- **Table Scanning**: Profile and analyze tables with detailed statistics
- **Delta Lake Integration**: Full support for Delta Lake features and optimizations
- **Comprehensive Testing**: Robust test suite with high coverage

## Prerequisites

- Python 3.8 or higher
- Java 8 or higher
- Git

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Nessi-datatool/nessi.git
   cd nessi
   ```

2. Create and activate a virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

## Configuration

Configure Spark with Delta Lake support:

```python
from pyspark.sql import SparkSession

spark = SparkSession.builder \
    .appName("Nessi") \
    .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
    .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
    .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0,io.delta:delta-storage:2.4.0") \
    .getOrCreate()
```

Optional environment variables:
```bash
export SPARK_HOME=/path/to/spark
export PYSPARK_PYTHON=python
export PYSPARK_DRIVER_PYTHON=python
```

## Usage

### Generate Sample Data

```python
from src.sample_data_generator import SampleDataGenerator

generator = SampleDataGenerator()
df = generator.generate_sample_data(rows=1000)
df.show()
```

### Scan Tables

```python
from src.scanner.table_scanner import TableScanner

scanner = TableScanner()
result = scanner.scan_delta_table("path/to/delta/table")
print(result)
```

### Create Test Tables

```python
from src.create_test_delta_table import create_test_delta_table

create_test_delta_table(
    path="path/to/table",
    rows=1000,
    schema={"id": "int", "name": "string", "value": "double"}
)
```

## Testing

Run the test suite:
```bash
python -m pytest tests/ -v
```

Generate coverage report:
```bash
python -m pytest tests/ -v --cov=src --cov-report=html
```

## Project Structure

```
nessi/
├── src/                    # Source code
│   ├── scanner/           # Table scanning functionality
│   ├── sample_data_generator.py
│   └── create_test_delta_table.py
├── tests/                 # Test files
├── data/                  # Data files
├── venv/                  # Virtual environment
├── requirements.txt       # Python dependencies
├── setup.py              # Package configuration
└── README.md             # Project documentation
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## Code of Conduct

Please review our [Code of Conduct](CODE_OF_CONDUCT.md) before participating in the project.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE.txt](LICENSE.txt) file for details.

## Documentation

For detailed API documentation, see [API.md](API.md).

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes. 