# Developer Experience Features

## Overview

Nessi.dev provides a comprehensive set of developer experience features designed to make data quality management and Delta Lake operations more efficient and user-friendly. This document outlines the implementation details of these features.

## CLI-First Approach

### One-Line Validation Commands

The one-line validation command provides a quick way to validate Delta tables:

```bash
# Basic validation with default rules
nessi validate /path/to/table

# Output in JSON format
nessi validate /path/to/table --format json

# Set validation threshold
nessi validate /path/to/table --threshold 0.9

# Save results to a file
nessi validate /path/to/table --output-file results.json
```

Implementation details:
- Performs basic validation checks on Delta tables
- Verifies Delta format, schema, data files, and transaction logs
- Configurable validation threshold
- Multiple output formats (text, JSON)
- Detailed validation results

### Batch Processing Scripts

The batch processing feature allows running data quality checks on multiple tables:

```bash
# Run batch job with a config file
nessi batch config.json

# Run batch job with specific output directory
nessi batch config.json --output ./results

# Run parallel processing
nessi batch config.json --parallel 4

# Filter batch jobs by tags
nessi batch --tags production,daily
```

Batch configuration example (JSON):
```json
{
  "name": "Daily Quality Check",
  "description": "Run daily quality checks on production tables",
  "tables": [
    "/path/to/table1",
    "/path/to/table2",
    "/path/to/table3"
  ],
  "rules": "production_rules.json",
  "output": "daily_results",
  "format": "json",
  "tags": ["production", "daily"],
  "params": {
    "threshold": "0.9",
    "verbose": "true"
  }
}
```

Implementation details:
- Support for JSON and YAML configuration files
- Configurable output directory
- Parallel processing capabilities
- Tagging system for organizing batch jobs
- Comprehensive result summaries

### Configuration Management

The configuration management system allows viewing and modifying Nessi.dev settings:

```bash
# View all configuration settings
nessi config

# View a specific configuration setting
nessi config server.port

# Set a configuration setting
nessi config server.port --set 8080

# Reset configuration to defaults
nessi config --reset

# Export configuration to a file
nessi config export config.json

# Import configuration from a file
nessi config import config.json
```

Implementation details:
- Hierarchical configuration structure
- Support for environment variables
- Default configuration generation
- Configuration import/export (JSON, YAML)
- User-specific configuration overrides

## Integration with Other Tools

### Metrics Export

Nessi provides flexible metrics export capabilities:

- Exports metrics in multiple formats (JSON, CSV)
- Supports custom metrics collection
- Enables integration with external monitoring tools

### Reporting System

The reporting system provides comprehensive visualization capabilities:

- HTML reports with interactive elements
- PDF reports for documentation and sharing
- Customizable report templates
- Scheduled report generations
- Shareable insights with team members
- Historical trend visualization

## Security Features

Building on the comprehensive security implementation, the developer experience includes:

- Secure configuration management
- Authentication for CLI commands
- Role-based access control
- Audit logging for configuration changes

## Testing

Comprehensive tests ensure the reliability of the developer experience features:

- **Command Tests**: Verify CLI command functionality
- **Configuration Tests**: Validate configuration management
- **Batch Processing Tests**: Ensure reliable batch operations
- **Integration Tests**: Verify seamless integration with other components

## Webhook Integration

Nessi.dev provides comprehensive webhook integration for event notifications:

```bash
# List all registered webhooks
nessi webhook list

# Create a new webhook
nessi webhook create --name "Quality Alert" --url "https://example.com/webhook" --events "data.validation.failed,alert.triggered"

# Get details of a webhook
nessi webhook get webhook-id

# Update a webhook
nessi webhook update webhook-id --name "Updated Webhook" --url "https://example.com/new-webhook"

# Enable/disable a webhook
nessi webhook enable webhook-id
nessi webhook disable webhook-id

# List supported event types
nessi webhook events

# Test a webhook
nessi webhook test webhook-id data.validation.failed
```

Implementation details:
- Event-based notifications for various system events
- Configurable webhook endpoints with custom headers
- Flexible event filtering with wildcard support
- Automatic retry mechanism for failed deliveries
- JSON-based webhook payloads with standardized format

See [webhook_integration.md](./webhook_integration.md) for more details.

## Plugin System

Nessi.dev includes a flexible plugin system for custom extensions:

```bash
# List all installed plugins
nessi plugin list

# Install a plugin
nessi plugin install /path/to/plugin.so

# Get plugin information
nessi plugin info plugin-name

# Enable/disable a plugin
nessi plugin enable plugin-name
nessi plugin disable plugin-name

# List supported plugin types
nessi plugin types

# Execute a plugin function
nessi plugin exec plugin-name function-name arg1 arg2
```

Implementation details:
- Multiple plugin types (validation, alert, metric, storage, export, UI)
- Dynamic loading of Go plugins as shared libraries
- Plugin metadata and versioning support
- Thread-safe plugin management

See [plugin_system.md](./plugin_system.md) for more details.

## Integration Examples

Nessi.dev can be easily integrated with various external systems:

### Workflow Orchestration Integrations

Nessi.dev provides dedicated integration packages for popular workflow orchestration tools and container orchestration platforms:

- **Apache Airflow**: Custom operators and sensors for Airflow 2.0+
- **Prefect**: Custom tasks and flows for Prefect 2.0+
- **Dagster**: Custom ops, resources, and jobs for Dagster 1.0+

These integrations allow you to incorporate data quality checks, profiling, validation, and lineage tracking into your data pipelines. See the [Workflow Orchestration Integration Guide](./workflow_orchestration.md) for detailed documentation.

#### Airflow Integration Example

```python
from airflow import DAG
from airflow.utils.dates import days_ago
from nessi_airflow.operators.data_quality_operator import NessiDataQualityOperator
from nessi_airflow.sensors.data_quality_sensor import NessiDataQualitySensor

dag = DAG('nessi_data_quality', start_date=days_ago(1))

# Run a data quality check
quality_check = NessiDataQualityOperator(
    task_id='run_quality_check',
    table_name='my_table',
    rules=[{'name': 'not_null_check', 'rule_type': 'not_null', 'column': 'id'}],
    quality_threshold=0.9,
    conn_id='nessi_default',
    dag=dag,
)

# Wait for the quality check to complete
wait_for_quality = NessiDataQualitySensor(
    task_id='wait_for_quality',
    check_id="{{ task_instance.xcom_pull(task_ids='run_quality_check')['check_id'] }}",
    quality_threshold=0.9,
    conn_id='nessi_default',
    dag=dag,
)

# Define task dependencies
quality_check >> wait_for_quality
```

#### Prefect Integration Example

```python
from prefect import flow
from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results

@flow(name="Data Quality Flow")
def data_quality_flow(table_name: str):
    # Run quality check
    quality_response = run_quality_check(
        table_name=table_name,
        rules=[{'name': 'not_null_check', 'rule_type': 'not_null', 'column': 'id'}],
        profile=True,
        api_host="http://your-nessi-instance",
        api_key="your-api-key",
    )
    
    # Wait for quality results
    quality_results = wait_for_quality_results(
        check_id=quality_response["check_id"],
        quality_threshold=0.9,
        fail_on_rule_failure=True,
        api_host="http://your-nessi-instance",
        api_key="your-api-key",
    )
    
    return quality_results

if __name__ == "__main__":
    data_quality_flow("my_table")
```

#### Dagster Integration Example

```python
from dagster import job
from nessi_dagster.resources import nessi_resource
from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results
```

#### Kubernetes Integration Example

```python
from nessi_k8s.client import NessiK8sClient

# Create client
client = NessiK8sClient(
    namespace="nessi",
    image="nessi/nessi:latest",
    api_host="http://your-nessi-instance",
    api_key="your-api-key",
)

# Run a data quality check
result = client.run_quality_check(
    table_name="my_table",
    rules=[
        {
            "name": "not_null_check",
            "description": "Check that id column is not null",
            "rule_type": "not_null",
            "column": "id",
        },
    ],
    profile=True,
    wait_for_completion=True,
)

print(f"Job status: {result}")
```

```yaml
# Example Kubernetes manifest for running a Nessi job
apiVersion: batch/v1
kind: Job
metadata:
  name: nessi-quality-check
  namespace: nessi
spec:
  template:
    spec:
      containers:
      - name: nessi
        image: nessi/nessi:latest
        command: ["nessi", "quality", "check", "my_table", "--format", "json"]
        env:
        - name: NESSI_API_HOST
          value: "http://nessi-api:8080"
        - name: NESSI_API_KEY
          valueFrom:
            secretKeyRef:
              name: nessi-secrets
              key: api-key
      restartPolicy: Never
  backoffLimit: 3
```

from dagster import job
from nessi_dagster.resources import nessi_resource
from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results

@job(
    resource_defs={
        "nessi": nessi_resource.configured({
            "api_host": "http://your-nessi-instance",
            "api_key": "your-api-key",
            "api_secret": "",
            "timeout": 300,
        }),
    },
)
def data_quality_job():
    # Run quality check
    quality_response = run_quality_check(
        table_name="my_table",
        rules=[{'name': 'not_null_check', 'rule_type': 'not_null', 'column': 'id'}],
        profile=True,
    )
    
    # Wait for quality results
    quality_results = wait_for_quality_results(
        check_id=quality_response["check_id"],
        quality_threshold=0.9,
        fail_on_rule_failure=True,
    )
    
    return quality_results

if __name__ == "__main__":
    result = data_quality_job.execute_in_process()
```

### GitHub Actions Integration

```yaml
name: Nessi Data Quality Checks

on:
  schedule:
    - cron: '0 0 * * *'  # Run daily at midnight

jobs:
  data-quality-checks:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        
      - name: Install Nessi
        run: |
          # Install Nessi
          ...
      
      - name: Run data quality checks
        run: |
          nessi validate /path/to/table --format json > results.json
```

### Custom Script Integration

```python
import subprocess
import json

def run_nessi_validation(table_path):
    """Run Nessi validation on a Delta table and return the results."""
    cmd = ['nessi', 'validate', table_path, '--format', 'json']
    result = subprocess.run(cmd, capture_output=True, text=True, check=True)
    return json.loads(result.stdout)

results = run_nessi_validation('/path/to/delta/table')
if results.get('passed', False):
    print("Validation passed!")
else:
    print("Validation failed!")
```

See [integration_guide.md](./integration_guide.md) for more detailed examples and best practices.

## Future Enhancements

Planned enhancements for developer experience include:

1. Interactive command-line interface with auto-completion
2. Scheduled batch jobs with cron-like syntax
3. Enhanced reporting capabilities
4. Additional workflow orchestration integrations (e.g., Luigi, Argo Workflows)
5. Expanded Kubernetes integration with Helm charts and operators

## Completed Enhancements

Recently completed enhancements:

1. **Workflow Orchestration Integrations**: Added comprehensive integrations for Apache Airflow, Prefect, Dagster, and Kubernetes
2. **Plugin System**: Implemented a flexible plugin system for custom extensions
3. **Webhook Integration**: Added comprehensive webhook support for event notifications
4. **Cloud Integration**: Implemented seamless integration with AWS, Azure, and GCP
5. **Intelligent Alerting**: Added advanced anomaly detection and pattern recognition
6. **Kubernetes Integration**: Added native support for running Nessi operations as Kubernetes jobs
7. **Observability Features**: Implemented metrics collection and structured logging for better monitoring
