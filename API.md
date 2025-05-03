# nessi.dev API Documentation

## Overview

nessi.dev is a Python-based data processing and analysis tool built with PySpark and Delta Lake. This document provides detailed API documentation for all major components.

## Table of Contents

1. [Sample Data Generation](#sample-data-generation)
2. [Table Scanning](#table-scanning)
3. [Report Generation](#report-generation)
4. [Data Format Support](#data-format-support)

## Sample Data Generation

### `generate_sample_data`

```python
def generate_sample_data(
    num_rows: int = 1000,
    patterns: Optional[Dict[str, str]] = None,
    seed: Optional[int] = None
) -> DataFrame
```

Generates sample data with various patterns.

**Parameters:**
- `num_rows`: Number of rows to generate (default: 1000)
- `patterns`: Dictionary of column patterns (default: None)
- `seed`: Random seed for reproducibility (default: None)

**Returns:**
- PySpark DataFrame with generated data

**Example:**
```python
from src.sample_data_generator import generate_sample_data

# Generate data with custom patterns
patterns = {
    "id": "increment",
    "name": "random_string",
    "age": "random_int(18, 65)",
    "timestamp": "current_timestamp"
}
data = generate_sample_data(num_rows=1000, patterns=patterns)
```

### `save_as_delta`

```python
def save_as_delta(
    df: DataFrame,
    path: str,
    mode: str = "overwrite",
    partition_by: Optional[List[str]] = None
) -> None
```

Saves DataFrame as a Delta table.

**Parameters:**
- `df`: PySpark DataFrame to save
- `path`: Output path for Delta table
- `mode`: Write mode (default: "overwrite")
- `partition_by`: List of columns to partition by (default: None)

**Example:**
```python
from src.sample_data_generator import save_as_delta

save_as_delta(data, "path/to/delta/table", partition_by=["date"])
```

## Table Scanning

### `TableScanner`

```python
class TableScanner:
    def scan_table(
        self,
        path: str,
        format: str = "delta",
        options: Optional[Dict[str, str]] = None
    ) -> Dict[str, Any]
```

Scans a table and returns statistics and metadata.

**Parameters:**
- `path`: Path to the table
- `format`: Table format ("delta", "parquet", or "csv")
- `options`: Additional options for reading the table

**Returns:**
- Dictionary containing table statistics and metadata

**Example:**
```python
from src.scanner.table_scanner import TableScanner

scanner = TableScanner()
result = scanner.scan_table("path/to/table", format="delta")
```

### `generate_report`

```python
def generate_report(
    self,
    scan_result: Dict[str, Any],
    template: str = "default",
    output_path: Optional[str] = None
) -> str
```

Generates a report from scan results.

**Parameters:**
- `scan_result`: Results from scan_table
- `template`: Report template to use
- `output_path`: Path to save the report (optional)

**Returns:**
- Generated report as string

**Example:**
```python
report = scanner.generate_report(result, template="detailed")
```

## Report Generation

### `ReportGenerator`

```python
class ReportGenerator:
    def __init__(
        self,
        template_dir: str = "templates",
        default_template: str = "default.html"
    )
```

Generates reports from scan results using Jinja2 templates.

**Parameters:**
- `template_dir`: Directory containing report templates
- `default_template`: Default template file name

**Methods:**
- `render(template_name: str, data: Dict[str, Any]) -> str`
- `save_report(content: str, output_path: str) -> None`

**Example:**
```python
from src.scanner.report_generator import ReportGenerator

generator = ReportGenerator()
report = generator.render("detailed.html", scan_result)
generator.save_report(report, "output/report.html")
```

## Data Format Support

### Delta Lake

```python
def read_delta_table(
    path: str,
    version: Optional[int] = None,
    timestamp: Optional[str] = None
) -> DataFrame
```

Reads a Delta table with optional version or timestamp.

**Parameters:**
- `path`: Path to Delta table
- `version`: Specific version to read
- `timestamp`: Timestamp to read at

**Example:**
```python
from src.scanner.delta_scanner import read_delta_table

# Read latest version
df = read_delta_table("path/to/delta/table")

# Read specific version
df = read_delta_table("path/to/delta/table", version=1)
```

### Parquet

```python
def read_parquet_table(
    path: str,
    options: Optional[Dict[str, str]] = None
) -> DataFrame
```

Reads a Parquet table.

**Parameters:**
- `path`: Path to Parquet table
- `options`: Additional read options

**Example:**
```python
from src.scanner.parquet_scanner import read_parquet_table

df = read_parquet_table("path/to/parquet/table")
```

### CSV

```python
def read_csv_table(
    path: str,
    options: Optional[Dict[str, str]] = None
) -> DataFrame
```

Reads a CSV table.

**Parameters:**
- `path`: Path to CSV file
- `options`: Additional read options

**Example:**
```python
from src.scanner.csv_scanner import read_csv_table

df = read_csv_table("path/to/csv/table")
```

## Error Handling

All functions raise appropriate exceptions:

- `FileNotFoundError`: When input file/directory doesn't exist
- `ValueError`: For invalid parameters
- `RuntimeError`: For processing errors

## Configuration

Configuration can be set through environment variables:

- `TEST_DATA_DIR`: Directory for test data
- `SPARK_DRIVER_MEMORY`: Spark driver memory
- `SPARK_EXECUTOR_MEMORY`: Spark executor memory

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to the API documentation. 