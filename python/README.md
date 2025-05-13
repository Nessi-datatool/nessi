# Nessi Monitoring Client

A Python client library for interacting with the Nessi monitoring system.

## Installation

```bash
pip install nessi-client
```

Or install from source:

```bash
git clone https://github.com/nessi-dev/nessi-dev.git
cd nessi-dev/python
pip install -e .
```

## Features

- **Metrics Management**: Send and retrieve monitoring metrics
- **Alerts Handling**: Access and manage monitoring alerts
- **Data Quality Rules**: Create, update, and manage data quality rules
- **Data Profiling**: Access data profiles and statistics
- **Data Validation**: Validate data against quality rules
- **Authentication**: Secure access with API keys or username/password
- **Multi-Format Support**: Work with Delta Lake, Parquet, and CSV files
- **Automatic Schema Inference**: Detect schemas from various data formats

## Quick Start

```python
from nessi_client import NessiClient, Metric
from datetime import datetime

# Initialize the client
client = NessiClient(
    base_url="http://localhost:8080",
    api_key="your-api-key"  # or use username/password
)

# Send a metric
metric = Metric(
    name="cpu_usage",
    value=45.2,
    timestamp=datetime.now(),
    tags={"host": "server-01", "environment": "production"}
)
client.send_metric(metric)

# Get recent metrics
metrics = client.get_metrics(
    name="cpu_usage",
    limit=10
)
for m in metrics:
    print(f"{m.timestamp}: {m.name} = {m.value}")

# Get alerts
alerts = client.get_alerts(
    severity="high",
    triggered=True
)
for alert in alerts:
    print(f"ALERT: {alert.name} - {alert.description}")

# Validate data
data = {
    "id": "12345",
    "name": "Sample Product",
    "price": 99.99,
    "quantity": 10
}
results = client.validate_data(data)
for result in results:
    status = "PASSED" if result.passed else "FAILED"
    print(f"{status}: {result.rule_id} - {result.message}")
```

## Authentication

The client supports two authentication methods:

### API Key (Recommended)

```python
client = NessiClient(
    base_url="http://localhost:8080",
    api_key="your-api-key"
)
```

### Username and Password

```python
client = NessiClient(
    base_url="http://localhost:8080",
    username="your-username",
    password="your-password"
)
```

## Advanced Usage

### Working with Data Profiles

```python
# Get profiles for a specific dataset
profiles = client.get_profiles(
    dataset_name="customer_data",
    limit=5
)

# Get a specific profile
profile = client.get_profile("profile-id-123")

# Access profile statistics
for column, stats in profile.column_stats.items():
    print(f"Column: {column}")
    print(f"  Type: {stats.get('type')}")
    print(f"  Null %: {stats.get('null_percentage')}")
    print(f"  Unique %: {stats.get('unique_percentage')}")
    
    if stats.get('type') == 'numeric':
        print(f"  Min: {stats.get('min')}")
        print(f"  Max: {stats.get('max')}")
        print(f"  Mean: {stats.get('mean')}")
        print(f"  Median: {stats.get('median')}")
```

### Managing Data Quality Rules

```python
from nessi_client import Rule

# Create a new rule
rule = Rule(
    id="price_range",
    name="Price Range Check",
    description="Validates that product prices are within an acceptable range",
    severity="error",
    field="price",
    rule_type="range",
    config={
        "min": 0.01,
        "max": 9999.99
    },
    tags=["product", "validation"]
)
client.create_rule(rule)

# Update an existing rule
updates = {
    "config": {
        "min": 0.01,
        "max": 5000.00  # Updated maximum price
    },
    "severity": "warning"
}
client.update_rule("price_range", updates)

# Get all rules
rules = client.get_rules()
for rule in rules:
    print(f"{rule.id}: {rule.name} ({rule.severity})")
```

## Error Handling

The client raises exceptions for API errors:

```python
from nessi_client import NessiClient
import requests

try:
    client = NessiClient(base_url="http://localhost:8080")
    metrics = client.get_metrics()
except requests.RequestException as e:
    print(f"API Error: {e}")
```

## IQR-based Outlier Detection

The client supports accessing profiles with IQR-based outlier detection, which is more robust for skewed distributions compared to the z-score method.

```python
# Get profiles with IQR-based outlier detection
profiles = client.get_profiles(
    dataset_name="sales_data",
    limit=1
)

# Check the outlier detection method
for profile in profiles:
    print(f"Profile: {profile.name}")
    print(f"Detection method: {profile.detection_method}")
    
    # Access outliers
    if profile.outliers:
        for column, values in profile.outliers.items():
            print(f"Outliers in {column}: {len(values)} found")
```

## Sudden Change Detection

The client can access metrics with sudden change detection capabilities:

```python
# Get metrics with change detection information
metrics = client.get_metrics(
    name="daily_sales",
    limit=30  # Get enough data points to analyze patterns
)

# The metadata contains change detection information
for metric in metrics:
    if "change_detection" in metric.metadata:
        change_info = metric.metadata["change_detection"]
        print(f"Metric: {metric.name} = {metric.value}")
        print(f"Change type: {change_info.get('pattern', 'normal')}")
        print(f"Magnitude: {change_info.get('magnitude')}")
        print(f"Direction: {change_info.get('direction')}")
```

## Multi-Format Support

The client provides comprehensive support for working with multiple data formats, including Delta Lake, Parquet, and CSV files.

### Format Detection

```python
from nessi_client import NessiClient

client = NessiClient()

# Detect the format of a file or directory
format_result = client.detect_format("/path/to/data")
print(f"Detected format: {format_result.format} with {format_result.confidence} confidence")

# Access the schema
if format_result.schema:
    print(f"Schema has {len(format_result.schema.fields)} fields:")
    for field in format_result.schema.fields:
        print(f"  - {field.name}: {field.data_type} (nullable: {field.nullable})")
```

### Reading Data

```python
# Read data from a file or directory
data_batch = client.read_data("/path/to/data")
print(f"Read {len(data_batch.data)} rows with {len(data_batch.schema.fields)} columns")

# Access the data (list of dictionaries)
for i, record in enumerate(data_batch.data[:3]):
    print(f"Record {i}: {record}")
```

### Schema Inference

```python
# Infer schema from a list of records
records = [
    {"id": 1, "name": "Alice", "score": 95.5},
    {"id": 2, "name": "Bob", "score": 87.0}
]
schema = client.infer_schema(records)

# Infer schema from a file
schema = client.infer_schema("/path/to/data.csv")

# Infer schema from a file-like object
with open("/path/to/data.csv", "rb") as f:
    schema = client.infer_schema(f)
```

### Validating Files

```python
# Validate a file against quality rules
results = client.validate_file("/path/to/data.csv")
for result in results:
    status = "PASSED" if result.passed else "FAILED"
    print(f"{status}: {result.rule_id} - {result.message}")

# Validate with specific rules
results = client.validate_file("/path/to/data.csv", rules=["rule1", "rule2"])
```

### Profiling Files

```python
# Create a profile for a file
profile = client.profile_file("/path/to/data.csv")
print(f"Profile created: {profile.id} - {profile.name}")
print(f"Row count: {profile.row_count}")
print(f"Column count: {profile.column_count}")

# Access column statistics
for column, stats in profile.column_stats.items():
    print(f"Column: {column}")
    print(f"  Type: {stats.get('type')}")
    print(f"  Null %: {stats.get('null_percentage')}")
```

### Custom Format Configuration

```python
from nessi_client import NessiClient
from nessi_client.format_models import FormatConfig

# Create a custom format configuration
format_config = FormatConfig(
    date_formats=["%Y-%m-%d", "%Y/%m/%d", "%d-%m-%Y"],
    csv_delimiter=",",
    csv_has_header=True,
    max_rows_for_inference=500,
    min_confidence_threshold=0.8
)

# Initialize client with custom format configuration
client = NessiClient(format_config=format_config)
```

## License

MIT
