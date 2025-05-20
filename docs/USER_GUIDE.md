# Nessi User Guide

---

## Features

- Data Quality Checks
- System Metrics
- CLI Interface
- Python API
- Reporting System
- Alerting System
- Metrics Export
- Profiling
- Freshness Monitoring
- Workflow Integrations
- Email Alerting
- Monitoring Dashboards

---

## Table of Contents

1. Introduction
2. Installation & Quickstart
3. Core Concepts
4. Security Features
5. Command-Line Interface
6. Python API
7. Data Quality Features
8. Reporting System
9. Monitoring & Metrics
10. Integrations
11. Configuration
12. Best Practices
13. Troubleshooting
14. FAQ

---

## Introduction

Nessi is a comprehensive monitoring system designed for data quality and system monitoring. It provides tools for profiling data, validating data against rules, and monitoring system metrics.

### Key Features

- **Delta Lake Support**: Schema evolution tracking, version control, and time travel
- **Data Quality**: Profiling, validation rules, and quality scoring
- **Monitoring**: Real-time metrics, comprehensive reporting, and alerting
- **Integration**: dbt, Airflow, and other workflow tools

## Installation & Quickstart

### Prerequisites

- Python 3.8+
- Go 1.18+
- Docker (optional)

### Installation

```bash
# Install from source
git clone https://github.com/nessi-dev/nessi.git
cd nessi
make install

# Or use Docker
docker pull nessi-dev/nessi:latest
```

For detailed installation instructions, see the [QUICKSTART.md](../QUICKSTART.md) guide.

## Core Concepts

### Data Quality

Nessi's data quality features include:

- **Profiling**: Statistical analysis of data
- **Validation**: Rule-based checks against data
- **Monitoring**: Tracking quality metrics over time
- **Alerting**: Notifications for quality issues

### Monitoring

Nessi provides monitoring for:

- **System Metrics**: CPU, memory, disk usage
- **Data Metrics**: Volume, freshness, quality
- **Process Metrics**: Job duration, success rates

## Security Features

Nessi includes several security features:

- JWT-based authentication
- TLS encryption
- API key management
- Security headers

## Reporting System

Nessi provides a comprehensive reporting system with multiple output formats:

- **HTML Reports**: Interactive reports with quality score cards and visualizations
- **PDF Reports**: Static reports for sharing and documentation
- **JSON/CSV**: Machine-readable formats for integration with other tools

Generate reports using the CLI:

```bash
# Generate a quality report for a table
nessi report --table path/to/table --format html --output report.html

# Generate a freshness report
nessi freshness report --format pdf --output freshness.pdf

# Export metrics as JSON for integration
nessi metrics export --format json --output metrics.json
```

Reports include detailed information about:

- Data quality scores and trends
- Schema evolution and changes
- Freshness monitoring results
- Validation rule results
- System performance metrics

## Command-Line Interface

Nessi provides a comprehensive CLI:

```bash
# Get help
nessi --help

# Run data quality checks
nessi validate --table path/to/table

# Generate a profile
nessi profile --table path/to/table

# View schema as ASCII tree
nessi schema tree --table path/to/table
```

## Python API

Nessi can be used programmatically:

```python
from nessi import Validator, Profiler

# Create a validator
validator = Validator()
validator.add_rule("completeness", "column_name", threshold=0.95)
results = validator.validate("path/to/table")

# Create a profiler
profiler = Profiler()
profile = profiler.profile("path/to/table")
```

## Data Quality Features

### Profiling

```bash
nessi profile --table path/to/table --output profile.json
```

### Validation

```bash
nessi validate --table path/to/table --rules rules.yaml
```

### Custom Rules

Create custom rules in YAML:

```yaml
rules:
  - name: completeness_check
    type: completeness
    column: user_id
    threshold: 0.99
```

## Monitoring

### Metrics Collection

Nessi collects metrics for:

- Table health
- Data quality
- System resources
- Process performance

### Alerting

Configure email alerts:

```yaml
alerts:
  - name: data_quality_alert
    condition: "quality_score < 0.9"
    recipients: ["team@example.com"]
    severity: critical
```

## Integrations

### dbt Integration

```bash
# Install dbt plugin
pip install nessi-dbt

# Run validation on dbt models
dbt run --nessi-validate
```

### Airflow Integration

```python
from airflow import DAG
from nessi.airflow import NessiValidateOperator

validate_task = NessiValidateOperator(
    task_id="validate_data",
    table_path="path/to/table",
    rules="rules.yaml",
    dag=dag
)
```

## Configuration

Nessi can be configured via:

- Config file (`~/.nessi/config.yaml`)
- Environment variables
- CLI flags

Example config:

```yaml
monitoring:
  retention_days: 30
  metrics_interval: 60

security:
  jwt_secret: "your-secret-key"
  enable_tls: true
```

## Best Practices

- Run validation after each data load
- Set up regular profiling jobs
- Configure alerts for critical metrics
- Use version control for rule definitions

## Troubleshooting

### Common Issues

- **Connection errors**: Check network and credentials
- **Validation failures**: Examine data quality issues
- **Performance problems**: Review resource allocation

### Logs

Logs are stored in `~/.nessi/logs/` by default.

## FAQ

### How do I contribute to Nessi?

See our [CONTRIBUTING.md](../CONTRIBUTING.md) guide.

### How do I report bugs?

Open an issue on our [GitHub repository](https://github.com/nessi-dev/nessi/issues).

### How do I get help?

Join our community on [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions).
