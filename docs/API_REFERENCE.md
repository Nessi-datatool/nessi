# Nessi API Reference

This document provides a comprehensive reference for the Nessi API, allowing you to integrate Nessi's data quality, schema management, and metrics capabilities into your applications.

## Table of Contents

- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [Core Concepts](#core-concepts)
- [Python Client](#python-client)
- [Go Client](#go-client)
- [REST API](#rest-api)
- [API Reference](#api-reference)
  - [Tables API](#tables-api)
  - [Quality API](#quality-api)
  - [Scan API](#scan-api)
  - [Schema API](#schema-api)
  - [Metrics API](#metrics-api)
  - [Report API](#report-api)
- [Error Handling](#error-handling)
- [Examples](#examples)

## Getting Started

### Installation

#### Python Client

```bash
pip install nessi-client
```

#### Go Client

```bash
go get github.com/nessi-dev/nessi-go-client
```

### Basic Usage

#### Python Example

```python
from nessi import Client

# Initialize client
client = Client(api_key="your-api-key")

# List tables
tables = client.tables.list()

# Run quality check
result = client.quality.check("/path/to/table")

# Generate report
report = client.report.generate("/path/to/table", type="quality", format="html")
```

#### Go Example

```go
package main

import (
    "fmt"
    "github.com/nessi-dev/nessi-go-client"
)

func main() {
    // Initialize client
    client := nessi.NewClient("your-api-key")

    // List tables
    tables, err := client.Tables.List()
    if err != nil {
        panic(err)
    }

    // Run quality check
    result, err := client.Quality.Check("/path/to/table")
    if err != nil {
        panic(err)
    }

    // Generate report
    report, err := client.Report.Generate("/path/to/table", "quality", "html")
    if err != nil {
        panic(err)
    }
}
```

## Authentication

Nessi API supports the following authentication methods:

- **API Key**: Pass your API key in the `X-API-Key` header
- **Bearer Token**: Use a JWT token with the `Authorization: Bearer <token>` header
- **Basic Auth**: Use username and password with Basic Authentication

### Obtaining an API Key

```bash
# Generate an API key using the CLI
nessi auth generate-key --name "My API Key" --expiry 365d
```

## Core Concepts

The Nessi API is organized around the following core concepts:

- **Tables**: Delta Lake tables and other data sources
- **Quality**: Data quality checks and rules
- **Scan**: Comprehensive data scanning and analysis
- **Schema**: Schema management and validation
- **Metrics**: Data quality metrics collection and analysis
- **Report**: Report generation and management

## Python Client

The Python client provides a Pythonic interface to the Nessi API.

### Installation

```bash
pip install nessi-client
```

### Client Initialization

```python
from nessi import Client

# Initialize with API key
client = Client(api_key="your-api-key")

# Initialize with bearer token
client = Client(token="your-bearer-token")

# Initialize with basic auth
client = Client(username="your-username", password="your-password")

# Initialize with custom endpoint
client = Client(api_key="your-api-key", endpoint="https://your-nessi-server/api")
```

## Go Client

The Go client provides a Go interface to the Nessi API.

### Installation

```bash
go get github.com/nessi-dev/nessi-go-client
```

### Client Initialization

```go
import "github.com/nessi-dev/nessi-go-client"

// Initialize with API key
client := nessi.NewClient("your-api-key")

// Initialize with bearer token
client := nessi.NewClientWithToken("your-bearer-token")

// Initialize with basic auth
client := nessi.NewClientWithBasicAuth("your-username", "your-password")

// Initialize with custom endpoint
client := nessi.NewClientWithOptions(nessi.ClientOptions{
    APIKey:   "your-api-key",
    Endpoint: "https://your-nessi-server/api",
})
```

## REST API

The Nessi REST API provides a standard HTTP interface for integrating with Nessi.

### Base URL

```
https://api.nessi.dev/v1
```

### Authentication

```
X-API-Key: your-api-key
```

or

```
Authorization: Bearer your-bearer-token
```

### Content Type

```
Content-Type: application/json
```

## API Reference

### Tables API

#### List Tables

Lists all available tables.

**Python**

```python
tables = client.tables.list()
```

**Go**

```go
tables, err := client.Tables.List()
```

**REST**

```
GET /tables
```

Response:

```json
{
  "tables": [
    {
      "name": "customers",
      "path": "/data/customers",
      "format": "delta",
      "size": 1288490188,
      "modified": "2025-05-01T12:34:56Z"
    },
    {
      "name": "orders",
      "path": "/data/orders",
      "format": "delta",
      "size": 3758096384,
      "modified": "2025-05-02T10:11:12Z"
    }
  ]
}
```

#### Describe Table

Shows detailed information about a table.

**Python**

```python
table = client.tables.describe("/path/to/table")
```

**Go**

```go
table, err := client.Tables.Describe("/path/to/table")
```

**REST**

```
GET /tables/{table_path}
```

Response:

```json
{
  "table": {
    "name": "customers",
    "path": "/data/customers",
    "format": "delta",
    "size_bytes": 1288490188,
    "row_count": 1000000,
    "created_at": "2025-04-15T08:30:45Z",
    "modified_at": "2025-05-01T12:34:56Z",
    "current_version": 5,
    "schema": [
      {
        "name": "id",
        "type": "string",
        "nullable": false
      },
      {
        "name": "name",
        "type": "string",
        "nullable": false
      }
    ]
  }
}
```

### Quality API

#### Run Quality Check

Runs quality checks on a table.

**Python**

```python
result = client.quality.check("/path/to/table", rules="rules.json")
```

**Go**

```go
result, err := client.Quality.Check("/path/to/table", "rules.json")
```

**REST**

```
POST /quality/check
```

Request:

```json
{
  "source": "/path/to/table",
  "rules": "rules.json"
}
```

Response:

```json
{
  "result": {
    "source": "/path/to/table",
    "timestamp": "2025-05-22T22:12:09Z",
    "rules_checked": 4,
    "rules_passed": 3,
    "rules_failed": 1,
    "overall_score": 0.95,
    "details": [
      {
        "rule": "email_validation",
        "type": "pattern",
        "column": "email",
        "threshold": 0.99,
        "actual": 0.998,
        "status": "PASS",
        "severity": "high"
      }
    ]
  }
}
```

### Scan API

#### Run Scan

Runs a data scan on a table.

**Python**

```python
result = client.scan.run("/path/to/table", type="all")
```

**Go**

```go
result, err := client.Scan.Run("/path/to/table", "all")
```

**REST**

```
POST /scan/run
```

Request:

```json
{
  "source": "/path/to/table",
  "type": "all"
}
```

Response:

```json
{
  "result": {
    "source": "/path/to/table",
    "timestamp": "2025-05-22T22:12:09Z",
    "format": "delta",
    "rows": 1000000,
    "scan_duration": 45.2,
    "profile_summary": {
      "columns": 6,
      "data_size": 1288490188,
      "null_values_ratio": 0.0256,
      "distinct_values_ratio": 0.82
    }
  }
}
```

### Schema API

#### Show Schema

Shows the schema of a table.

**Python**

```python
schema = client.schema.show("/path/to/table")
```

**Go**

```go
schema, err := client.Schema.Show("/path/to/table")
```

**REST**

```
GET /schema/{table_path}
```

Response:

```json
{
  "schema": {
    "name": "customers",
    "fields": [
      {
        "name": "id",
        "type": "string",
        "nullable": false,
        "metadata": {
          "description": "Unique customer identifier"
        }
      },
      {
        "name": "name",
        "type": "string",
        "nullable": false,
        "metadata": {
          "description": "Customer's full name"
        }
      }
    ],
    "metadata": {
      "created": "2025-04-15T08:30:45Z",
      "modified": "2025-05-01T12:34:56Z",
      "primaryKey": ["id"],
      "partitionColumns": ["signup_date"]
    }
  }
}
```

### Metrics API

#### Collect Metrics

Collects metrics from a table.

**Python**

```python
metrics = client.metrics.collect("/path/to/table", type="quality")
```

**Go**

```go
metrics, err := client.Metrics.Collect("/path/to/table", "quality")
```

**REST**

```
POST /metrics/collect
```

Request:

```json
{
  "source": "/path/to/table",
  "type": "quality"
}
```

Response:

```json
{
  "metrics": {
    "source": "/path/to/table",
    "timestamp": "2025-05-22T22:12:09Z",
    "quality": {
      "completeness": 0.9744,
      "accuracy": 0.9950,
      "consistency": 0.9900,
      "uniqueness": 0.9950,
      "timeliness": 0.9850,
      "overall_score": 0.9879
    }
  }
}
```

### Report API

#### Generate Report

Generates a report from scan results, quality checks, or metrics.

**Python**

```python
report = client.report.generate("/path/to/table", type="quality", format="html")
```

**Go**

```go
report, err := client.Report.Generate("/path/to/table", "quality", "html")
```

**REST**

```
POST /report/generate
```

Request:

```json
{
  "source": "/path/to/table",
  "type": "quality",
  "format": "html"
}
```

Response:

```json
{
  "report": {
    "source": "/path/to/table",
    "type": "quality",
    "format": "html",
    "timestamp": "2025-05-22T22:12:09Z",
    "url": "https://api.nessi.dev/reports/report_12345.html",
    "content": "<!DOCTYPE html>..."
  }
}
```

## Error Handling

The API uses standard HTTP status codes to indicate success or failure:

- `200 OK`: The request was successful
- `400 Bad Request`: The request was invalid
- `401 Unauthorized`: Authentication failed
- `403 Forbidden`: The authenticated user doesn't have permission
- `404 Not Found`: The requested resource was not found
- `500 Internal Server Error`: An error occurred on the server

Error responses include a JSON object with details:

```json
{
  "error": {
    "code": "N500",
    "message": "Delta Lake table not found",
    "details": "The table '/path/to/nonexistent' does not exist or is not a valid Delta Lake table",
    "suggestion": "Check the table path and ensure it is a valid Delta Lake table"
  }
}
```

### Error Codes

The API uses the same error codes as the CLI:

- `N500-N599`: Table Operations
- `N600-N699`: Scan Operations
- `N700-N799`: Schema Operations
- `N800-N899`: Quality Operations
- `N900-N999`: Metrics Operations
- `N1000-N1099`: Report Operations

## Examples

### Complete Python Example

```python
from nessi import Client

# Initialize client
client = Client(api_key="your-api-key")

# List tables
tables = client.tables.list()
print(f"Found {len(tables)} tables")

# Check table schema
schema = client.schema.show("/data/customers")
print(f"Table has {len(schema.fields)} columns")

# Run quality check
result = client.quality.check("/data/customers")
print(f"Quality score: {result.overall_score:.2f}")

# Collect metrics
metrics = client.metrics.collect("/data/customers")
print(f"Completeness: {metrics.quality.completeness:.2f}")

# Generate report
report = client.report.generate("/data/customers", type="quality", format="html")
print(f"Report generated: {report.url}")

# Handle errors
try:
    result = client.quality.check("/nonexistent/table")
except Exception as e:
    print(f"Error: {e}")
```

### Complete Go Example

```go
package main

import (
    "fmt"
    "github.com/nessi-dev/nessi-go-client"
)

func main() {
    // Initialize client
    client := nessi.NewClient("your-api-key")

    // List tables
    tables, err := client.Tables.List()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d tables\n", len(tables))

    // Check table schema
    schema, err := client.Schema.Show("/data/customers")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Table has %d columns\n", len(schema.Fields))

    // Run quality check
    result, err := client.Quality.Check("/data/customers")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Quality score: %.2f\n", result.OverallScore)

    // Collect metrics
    metrics, err := client.Metrics.Collect("/data/customers")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Completeness: %.2f\n", metrics.Quality.Completeness)

    // Generate report
    report, err := client.Report.Generate("/data/customers", "quality", "html")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Report generated: %s\n", report.URL)

    // Handle errors
    _, err = client.Quality.Check("/nonexistent/table")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
    }
}
```

## Feature Availability

| Feature | Community Edition | Pro Edition |
|---------|-------------------|-------------|
| Basic Table Operations | ✓ | ✓ |
| Basic Schema Operations | ✓ | ✓ |
| Basic Quality Checks | ✓ | ✓ |
| Basic Data Scanning | ✓ | ✓ |
| Basic Metrics Collection | ✓ | ✓ |
| Simple Reports | ✓ | ✓ |
| Advanced Schema Operations | | ✓ |
| Advanced Quality Checks | | ✓ |
| Comprehensive Metrics | | ✓ |
| Enhanced Reports | | ✓ |
| Grafana Integration | | ✓ |
| Workflow Orchestration | | ✓ |
| Databricks Integration | | ✓ |
| AWS S3 Storage | | ✓ |

For more information on licensing, see the [License Management documentation](LICENSE_MANAGEMENT.md).
