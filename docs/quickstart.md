# Quick Start Guide

## Getting Started with Nessi

This guide will help you get up and running with Nessi quickly. We'll cover the basic operations and most common use cases.

## Basic Operations

### 1. Initialize Delta Lake Table

```python
from nessi import DeltaManager

# Initialize Delta Manager
delta_manager = DeltaManager(spark)

# Create a new Delta table
delta_manager.create_table(
    path="/data/my_table",
    schema="id INT, name STRING, age INT",
    partition_by=["age"]
)
```

### 2. Data Quality Analysis

```python
from nessi import QualityAnalyzer

# Initialize Quality Analyzer
analyzer = QualityAnalyzer()

# Analyze DataFrame
results = analyzer.analyze_dataframe(df)

# Print quality scores
print(f"Overall Quality: {results['overall_quality']}")
print(f"Completeness: {results['completeness']}")
print(f"Consistency: {results['consistency']}")
```

### 3. Generate Report

```python
from nessi import ReportGenerator

# Initialize Report Generator
generator = ReportGenerator()

# Generate HTML report
report = generator.generate_report(results, format="html")

# Save report
generator.save_report(report, "quality_report.html")

# Display report
generator.display_report(report)
```

## Common Use Cases

### 1. Time Travel Analysis

```python
# Access historical data
df = delta_manager.time_travel(
    path="/data/my_table",
    version=2  # or timestamp="2024-01-01"
)

# Compare with current data
current_df = spark.read.format("delta").load("/data/my_table")
```

### 2. Schema Evolution

```python
# Evolve schema
delta_manager.evolve_schema("/data/my_table", {
    "new_column": {
        "type": "STRING",
        "nullable": True,
        "comment": "New column description"
    }
})
```

### 3. Quality Monitoring

```python
from nessi import Monitor

# Initialize Monitor
monitor = Monitor()

# Start monitoring
monitor.start_scan()

# Record quality metrics
monitor.update_quality_score(0.95)
monitor.record_error("Data validation failed")

# End monitoring
monitor.end_scan()
```

## CLI Commands

### Basic Operations

```bash
# Analyze data quality
nessi analyze --format delta --path /data/my_table

# Generate report
nessi report --format html --output report.html

# Monitor table
nessi monitor --table my_table
```

### Advanced Operations

```bash
# Time travel
nessi time-travel --table my_table --version 2

# Schema evolution
nessi evolve-schema --table my_table --schema schema.json

# Quality check
nessi quality-check --table my_table --rules rules.json
```

## Integration Examples

### 1. Airflow Integration

```python
from airflow import DAG
from nessi.operators import NessiOperator

with DAG('nessi_pipeline') as dag:
    analyze = NessiOperator(
        task_id='analyze',
        command='analyze',
        params={'path': '/data/table'}
    )
```

### 2. GitHub Actions

```yaml
- name: Run Nessi Analysis
  uses: nessi/action@v1
  with:
    command: analyze
    path: data/table
```

## Security Configuration

### 1. API Authentication

```python
from nessi import SecurityMiddleware

# Initialize security middleware
security = SecurityMiddleware()

# Create token
token = security.create_token(
    user_id="user123",
    roles=["admin", "analyst"]
)

# Verify token
payload = security.verify_token(token)
```

### 2. Rate Limiting

```python
@security.rate_limit_middleware(limit=100, window=60)
def rate_limited_endpoint():
    pass
```

## Monitoring Setup

### 1. Grafana Dashboard

1. Access Grafana at http://localhost:3000
2. Import dashboards from `grafana/dashboards`
3. Configure data sources

### 2. Alert Configuration

```yaml
# alertmanager.yml
routes:
  - match:
      severity: critical
    receiver: email
receivers:
  - name: email
    email_configs:
      - to: admin@example.com
```

## Best Practices

### 1. Data Quality

- Run regular quality checks
- Set up automated alerts
- Monitor data completeness
- Track consistency metrics

### 2. Performance

- Use appropriate partitioning
- Implement Z-ordering for frequent queries
- Regular table optimization
- Monitor query performance

### 3. Security

- Use strong JWT secrets
- Implement proper role hierarchy
- Monitor rate limits
- Regular security audits

## Next Steps

- [Configuration Guide](configuration.md)
- [API Reference](api/python.md)
- [Best Practices](best_practices/quality.md)
- [Troubleshooting Guide](troubleshooting/common_issues.md) 