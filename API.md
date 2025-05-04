# NESSI API Documentation

## Overview

NESSI provides a RESTful API for data quality scanning and monitoring. The API is designed to run in Docker containers and uses FastAPI for the web framework.

## Deployment

### Docker Deployment

```bash
# Build and run the API container
docker-compose build nessi-api
docker-compose up -d nessi-api
```

The API will be available at `http://localhost:8000`.

## API Endpoints

### Scanner Endpoints

#### POST /api/v1/scanner/delta/scan
Scan a Delta table for data quality issues.

**Request Body:**
```json
{
    "table_path": "path/to/delta/table",
    "options": {
        "validate_schema": true,
        "check_constraints": true
    }
}
```

**Response:**
```json
{
    "status": "success",
    "data": {
        "table_name": "example_table",
        "row_count": 1000,
        "schema": {
            "fields": [...]
        },
        "quality_metrics": {
            "completeness": 0.99,
            "accuracy": 0.98,
            "consistency": 0.97
        }
    }
}
```

### Monitoring Endpoints

#### GET /api/v1/monitoring/metrics
Get current system metrics.

**Response:**
```json
{
    "status": "success",
    "data": {
        "cpu_usage": 45.2,
        "memory_usage": 60.1,
        "disk_usage": 30.5,
        "network_io": {
            "bytes_in": 1024,
            "bytes_out": 2048
        }
    }
}
```

## Authentication

All API endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <token>
```

## Error Handling

The API uses standard HTTP status codes and returns error messages in the following format:

```json
{
    "status": "error",
    "error": {
        "code": "ERROR_CODE",
        "message": "Error description"
    }
}
```

## Rate Limiting

API requests are rate-limited to prevent abuse. The current limits are:
- 100 requests per minute per IP
- 1000 requests per hour per API key

## WebSocket Endpoints

### /ws/v1/monitoring/stream
Stream real-time monitoring data.

**Message Format:**
```json
{
    "type": "metric_update",
    "data": {
        "metric_name": "cpu_usage",
        "value": 45.2,
        "timestamp": "2024-01-01T00:00:00Z"
    }
}
```

## Docker Health Checks

The API container includes health check endpoints:

### GET /health
Check API health status.

**Response:**
```json
{
    "status": "healthy",
    "version": "1.0.0",
    "timestamp": "2024-01-01T00:00:00Z"
}
```

## Development

To run the API in development mode:

```bash
docker-compose -f docker-compose.dev.yml up -d nessi-api
```

The development API will be available at `http://localhost:8001` with additional debugging features enabled.

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

### `generate_pdf_report`

```python
def generate_pdf_report(
    self,
    scan_result: Dict[str, Any],
    output_path: Optional[str] = None
) -> str
```

Generates a PDF report from scan results.

**Parameters:**
- `scan_result`: Results from scan_table
- `output_path`: Path to save the PDF report (optional)

**Returns:**
- Path to the generated PDF report

**Example:**
```python
pdf_path = scanner.generate_pdf_report(result)
```

The PDF report includes:
- Data quality metrics with visual indicators
- Performance statistics
- Table scan results
- Schema information
- Timestamp and copyright information

## Report Generation

### `ReportGenerator`

```