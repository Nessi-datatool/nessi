# Nessi Monitoring System User Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Core Concepts](#core-concepts)
4. [Security Features](#security-features)
5. [Dashboard](#dashboard)
6. [Command-Line Interface](#command-line-interface)
7. [Python API](#python-api)
8. [Data Quality Features](#data-quality-features)
9. [Monitoring & Alerts](#monitoring--alerts)
10. [Configuration](#configuration)
11. [Best Practices](#best-practices)
12. [Troubleshooting](#troubleshooting)
13. [FAQ](#faq)

## Introduction

Nessi is a comprehensive monitoring system designed for data quality and system monitoring. It provides tools for profiling data, detecting anomalies, validating data against rules, and monitoring system metrics.

### Key Features

- **Data Quality Monitoring**: Profile your data, detect outliers, and validate against rules
- **System Monitoring**: Track system metrics like CPU, memory, and disk usage
- **Alerting**: Get notified when metrics exceed thresholds or data quality issues arise
- **Security**: Role-based access control, authentication, and SSL/TLS support
- **Multiple Interfaces**: Web dashboard, command-line tools, and Python API
- **Extensibility**: Plug in custom data sources, rules, and alert handlers

## Installation

### Prerequisites

- Go 1.16 or later
- Python 3.7 or later (for Python API)
- PostgreSQL or SQLite (for persistent storage)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/nessi-dev/nessi-dev.git
cd nessi-dev

# Build the binaries
go build -o nessi cmd/nessi/main.go

# Start the server
./nessi server start
```

The server will be available at http://localhost:8080 by default.

### Docker Installation

```bash
# Pull the Docker image
docker pull nessi-dev/nessi:latest

# Run the container
docker run -p 8080:8080 -v nessi-data:/data nessi-dev/nessi:latest
```

## Core Concepts

### Metrics

Metrics are the fundamental unit of monitoring in Nessi. A metric consists of:

- **Name**: Identifier for the metric (e.g., `cpu_usage`, `row_count`)
- **Value**: Numeric or string value
- **Timestamp**: When the metric was recorded
- **Tags**: Key-value pairs for categorization (e.g., `{"host": "server-01"}`)
- **Metadata**: Additional information about the metric

### Alerts

Alerts are triggered when metrics exceed defined thresholds or when data quality rules are violated. An alert includes:

- **Name**: Identifier for the alert
- **Description**: Human-readable description
- **Severity**: Importance level (e.g., `info`, `warning`, `error`, `critical`)
- **Triggered**: Whether the alert is active
- **Related Metric**: The metric that triggered the alert
- **Threshold**: The value that was exceeded

### Data Profiles

Data profiles provide statistical information about datasets, including:

- **Column Types**: Data types detected for each column
- **Basic Statistics**: Min, max, mean, median, etc. for numeric columns
- **Categorical Statistics**: Frequency counts for categorical columns
- **Null Values**: Percentage of null values per column
- **Outliers**: Detected outliers using Z-score or IQR methods

### Rules

Rules define expectations for data quality. A rule includes:

- **Name**: Identifier for the rule
- **Description**: Human-readable description
- **Severity**: Importance level when violated
- **Field**: The data field to validate
- **Rule Type**: The type of validation (e.g., `range`, `regex`, `enum`)
- **Configuration**: Type-specific configuration

## Security Features

### Authentication

Nessi supports multiple authentication methods:

#### Username and Password

```bash
# Login via CLI
nessi login --username admin --password your_password

# Login via API
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}'
```

#### API Keys

API keys provide programmatic access without username/password:

```bash
# Generate an API key
nessi apikey create --description "My API key"

# Use the API key
curl -X GET http://localhost:8080/api/metrics \
  -H "X-API-Key: your_api_key"
```

### Role-Based Access Control

Nessi implements role-based access control with three predefined roles:

- **Admin**: Full access to all features
- **User**: Can view and create metrics, alerts, and rules
- **Viewer**: Read-only access to metrics and alerts

```bash
# Create a new user with a specific role
nessi user create --username john --password secret --role user

# List users
nessi user list

# Update a user's role
nessi user update --username john --role admin
```

### SSL/TLS Support

Nessi supports HTTPS for secure communication:

```bash
# Generate self-signed certificates
nessi ssl generate

# Start the server with SSL
nessi server start --ssl
```

Configure SSL in the configuration file:

```yaml
ssl:
  enabled: true
  cert_file: /path/to/cert.pem
  key_file: /path/to/key.pem
  require_https: true  # Redirect HTTP to HTTPS
```

## Dashboard

The web dashboard provides a user-friendly interface for monitoring and management.

### Features

- **Metrics Visualization**: Charts and graphs for metrics
- **Alert Management**: View and manage alerts
- **Data Quality Reports**: View data profiles and validation results
- **User Management**: Manage users and roles (admin only)
- **Configuration**: Configure system settings (admin only)

### Access

The dashboard is available at `http://localhost:8080` by default. Log in with your username and password.

## Command-Line Interface

The CLI provides command-line access to Nessi features.

### Basic Commands

```bash
# Start the server
nessi server start

# Stop the server
nessi server stop

# View metrics
nessi metrics list

# Send a metric
nessi metrics send --name cpu_usage --value 45.2 --tags host=server-01

# View alerts
nessi alerts list

# View data profiles
nessi profiles list

# Validate data against rules
nessi validate --file data.csv
```

### Advanced Commands

```bash
# Export metrics to CSV
nessi metrics export --output metrics.csv

# Import rules from YAML
nessi rules import --file rules.yaml

# Generate a report
nessi report generate --output report.html
```

## Python API

The Python API provides programmatic access to Nessi features.

### Installation

```bash
# Install from PyPI
pip install nessi-client

# Or install from source
cd nessi-dev/python
pip install -e .
```

### Basic Usage

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
```

### Authentication

```python
# API key authentication (recommended)
client = NessiClient(
    base_url="http://localhost:8080",
    api_key="your-api-key"
)

# Username and password authentication
client = NessiClient(
    base_url="http://localhost:8080",
    username="admin",
    password="your-password"
)
```

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
    
    if stats.get('type') == 'numeric':
        print(f"  Min: {stats.get('min')}")
        print(f"  Max: {stats.get('max')}")
        print(f"  Mean: {stats.get('mean')}")
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
```

## Data Quality Features

### Data Profiling

Data profiling analyzes datasets to provide statistical information:

```bash
# Profile a CSV file
nessi profile --file data.csv

# Profile with custom options
nessi profile --file data.csv --sample-size 1000 --outlier-method iqr
```

Python API:

```python
# Get profiles
profiles = client.get_profiles(dataset_name="customer_data")

# Access outliers
for profile in profiles:
    if profile.outliers:
        for column, values in profile.outliers.items():
            print(f"Outliers in {column}: {len(values)} found")
```

### Outlier Detection

Nessi supports two methods for outlier detection:

- **Z-Score**: Identifies values that deviate from the mean by a certain number of standard deviations
- **IQR (Interquartile Range)**: More robust for skewed distributions, identifies values outside the interquartile range

```bash
# Profile with Z-score outlier detection
nessi profile --file data.csv --outlier-method zscore --zscore-threshold 3.0

# Profile with IQR outlier detection
nessi profile --file data.csv --outlier-method iqr --iqr-multiplier 1.5
```

### Sudden Change Detection

Detects abrupt changes in metrics over time:

```bash
# Enable change detection
nessi metrics monitor --name daily_sales --change-detection
```

Python API:

```python
# Get metrics with change detection information
metrics = client.get_metrics(
    name="daily_sales",
    limit=30
)

# The metadata contains change detection information
for metric in metrics:
    if "change_detection" in metric.metadata:
        change_info = metric.metadata["change_detection"]
        print(f"Change type: {change_info.get('pattern')}")
        print(f"Magnitude: {change_info.get('magnitude')}")
```

### Rule Validation

Validate data against quality rules:

```bash
# Validate a CSV file against all rules
nessi validate --file data.csv

# Validate against specific rules
nessi validate --file data.csv --rules rule1,rule2
```

Python API:

```python
# Validate data
data = [
    {"id": "12345", "name": "Product A", "price": 99.99},
    {"id": "67890", "name": "Product B", "price": -10.00}  # Invalid price
]
results = client.validate_data(data)

for result in results:
    status = "PASSED" if result.passed else "FAILED"
    print(f"{status}: {result.rule_id} - {result.message}")
```

## Monitoring & Alerts

### System Metrics

Monitor system metrics like CPU, memory, and disk usage:

```bash
# Start the system monitor
nessi monitor system start

# View system metrics
nessi metrics list --name cpu_usage
```

### Custom Metrics

Send custom metrics for application monitoring:

```bash
# Send a custom metric
nessi metrics send --name api_response_time --value 120 --tags endpoint=/users
```

Python API:

```python
# Send a custom metric
metric = Metric(
    name="api_response_time",
    value=120,
    timestamp=datetime.now(),
    tags={"endpoint": "/users", "method": "GET"}
)
client.send_metric(metric)
```

### Alert Configuration

Configure alerts based on metric thresholds:

```yaml
# alerts.yaml
alerts:
  - name: high_cpu_usage
    description: CPU usage is too high
    severity: warning
    metric: cpu_usage
    condition: ">= 80"
    duration: 5m  # Must exceed threshold for 5 minutes
```

```bash
# Import alert configuration
nessi alerts import --file alerts.yaml
```

### Alert Notifications

Configure notification channels for alerts:

```yaml
# notifications.yaml
channels:
  - name: email
    type: email
    config:
      recipients: ["admin@example.com"]
  
  - name: slack
    type: slack
    config:
      webhook_url: "https://hooks.slack.com/services/..."
```

```bash
# Import notification configuration
nessi notifications import --file notifications.yaml
```

## Configuration

### Configuration File

The main configuration file is located at `~/.nessi/config.yaml` by default:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  data_dir: ~/.nessi/data

database:
  type: sqlite  # or postgresql
  path: ~/.nessi/nessi.db  # for sqlite
  # connection_string: "host=localhost user=nessi password=secret dbname=nessi" # for postgresql

security:
  auth_enabled: true
  default_admin_username: admin
  default_admin_password: admin  # Change this!
  token_expiry: 24h
  
ssl:
  enabled: false
  cert_file: ~/.nessi/ssl/cert.pem
  key_file: ~/.nessi/ssl/key.pem
  require_https: false

profiling:
  sample_size: 1000
  outlier_method: zscore
  zscore_threshold: 3.0
  iqr_multiplier: 1.5
```

### Environment Variables

Configuration can also be provided via environment variables:

```bash
# Server configuration
export NESSI_SERVER_HOST=0.0.0.0
export NESSI_SERVER_PORT=8080

# Database configuration
export NESSI_DB_TYPE=postgresql
export NESSI_DB_CONNECTION_STRING="host=localhost user=nessi password=secret dbname=nessi"

# Security configuration
export NESSI_SECURITY_AUTH_ENABLED=true
export NESSI_SECURITY_DEFAULT_ADMIN_PASSWORD=secure_password
```

## Best Practices

### Security

- **Change default passwords**: Always change the default admin password
- **Use API keys**: For programmatic access, use API keys instead of username/password
- **Enable SSL**: In production, always enable SSL and require HTTPS
- **Principle of least privilege**: Assign the minimum necessary role to each user

### Performance

- **Sample large datasets**: Use sampling for large datasets to improve profiling performance
- **Index frequently queried fields**: If using PostgreSQL, create indexes for frequently queried fields
- **Limit metric retention**: Configure retention policies to limit the storage of old metrics

### Monitoring

- **Start with basic metrics**: Begin with CPU, memory, and disk usage metrics
- **Add application metrics**: Monitor application-specific metrics like response times
- **Configure meaningful alerts**: Set thresholds that indicate actual problems, not normal variations
- **Use tags**: Tag metrics with relevant information for better filtering and analysis

## Troubleshooting

### Common Issues

#### Server Won't Start

- Check if another process is using the same port
- Verify that the data directory is writable
- Check the logs for detailed error messages

```bash
# Check for processes using port 8080
lsof -i :8080

# View logs
cat ~/.nessi/logs/nessi.log
```

#### Authentication Issues

- Verify that you're using the correct username and password
- Check if your token has expired
- For API keys, ensure the key is valid and not revoked

```bash
# Reset the admin password
nessi user reset-password --username admin
```

#### SSL/TLS Issues

- Verify that the certificate and key files exist and are readable
- Check if the certificate is valid and not expired
- For self-signed certificates, disable SSL verification in clients

```bash
# Generate new self-signed certificates
nessi ssl generate
```

### Logs

Logs are stored in `~/.nessi/logs/` by default:

```bash
# View server logs
cat ~/.nessi/logs/nessi.log

# View access logs
cat ~/.nessi/logs/access.log
```

## FAQ

### General

**Q: What makes Nessi different from other monitoring systems?**  
A: Nessi combines system monitoring with data quality monitoring, providing a unified view of both infrastructure and data health.

**Q: Is Nessi suitable for production use?**  
A: Yes, Nessi is designed for production use with features like authentication, SSL support, and scalable storage options.

### Data Quality

**Q: How does Nessi detect outliers?**  
A: Nessi supports two methods for outlier detection: Z-score and IQR (Interquartile Range). Z-score is better for normally distributed data, while IQR is more robust for skewed distributions.

**Q: Can I create custom validation rules?**  
A: Yes, Nessi supports custom validation rules through a flexible rule configuration system.

### Security

**Q: How secure is the authentication system?**  
A: Nessi uses industry-standard security practices, including password hashing with bcrypt, JWT tokens with expiration, and role-based access control.

**Q: Can I integrate with existing authentication systems?**  
A: Nessi can be extended to support external authentication providers like LDAP or OAuth.

### Technical

**Q: What databases does Nessi support?**  
A: Nessi supports SQLite for simple deployments and PostgreSQL for production use.

**Q: How can I back up Nessi data?**  
A: For SQLite, back up the database file. For PostgreSQL, use standard PostgreSQL backup procedures.

**Q: Is there an API rate limit?**  
A: By default, there is no rate limit, but you can configure one in the server settings for production deployments.
