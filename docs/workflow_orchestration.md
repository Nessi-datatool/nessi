# Workflow Orchestration Integration

Nessi.dev provides deep integration with popular workflow orchestration tools, allowing you to incorporate data quality checks, profiling, validation, and lineage tracking into your data pipelines.

## Supported Orchestration Tools

- **Apache Airflow**: Custom operators and sensors for Airflow 2.0+
- **Prefect**: Custom tasks and flows for Prefect 2.0+
- **Dagster**: Custom ops, resources, and jobs for Dagster 1.0+
- **Kubernetes**: Native operators for running Nessi operations as Kubernetes jobs

## Common Features

All workflow orchestration integrations provide the following capabilities:

- **Data Quality Checks**: Run quality checks on your data and wait for results
- **Data Profiling**: Generate profiles of your data to understand its characteristics
- **Data Validation**: Validate your data against custom rules
- **Lineage Tracking**: Capture and visualize data lineage information
- **Failure Handling**: Configure whether to fail workflows based on quality scores or rule failures
- **Authentication**: Secure communication with Nessi.dev API using API keys
- **Observability**: Metrics collection and structured logging for monitoring and debugging

## Apache Airflow Integration

### Installation

```bash
pip install nessi-airflow
```

### Configuration

1. Set up Airflow connection for Nessi.dev:

```bash
airflow connections add 'nessi_default' \
    --conn-type 'http' \
    --conn-host 'http://your-nessi-instance' \
    --conn-login 'your-api-key' \
    --conn-password 'your-api-secret'
```

2. Import and use the operators and sensors in your DAGs:

```python
from nessi_airflow.operators.data_quality_operator import NessiDataQualityOperator
from nessi_airflow.operators.profile_operator import NessiProfileOperator
from nessi_airflow.operators.validation_operator import NessiValidationOperator
from nessi_airflow.operators.lineage_operator import NessiLineageOperator
from nessi_airflow.sensors.data_quality_sensor import NessiDataQualitySensor
from nessi_airflow.sensors.validation_sensor import NessiValidationSensor
from nessi_airflow.sensors.profile_sensor import NessiProfileSensor
```

### Example DAG

```python
from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.python import PythonOperator

from nessi_airflow.operators.data_quality_operator import NessiDataQualityOperator
from nessi_airflow.sensors.data_quality_sensor import NessiDataQualitySensor

default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'email_on_failure': False,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

with DAG(
    'nessi_data_quality_example',
    default_args=default_args,
    description='A simple Nessi.dev data quality example',
    schedule_interval=timedelta(days=1),
    start_date=datetime(2022, 1, 1),
    catchup=False,
    tags=['nessi', 'data-quality'],
) as dag:

    # Run a data quality check
    run_quality_check = NessiDataQualityOperator(
        task_id='run_quality_check',
        table_name='my_table',
        rules=[
            {
                'name': 'not_null_check',
                'description': 'Check that id column is not null',
                'rule_type': 'not_null',
                'column': 'id',
            },
            {
                'name': 'unique_check',
                'description': 'Check that id column is unique',
                'rule_type': 'unique',
                'column': 'id',
            },
        ],
        profile=True,
        quality_threshold=0.9,
        conn_id='nessi_default',
    )

    # Wait for the quality check to complete
    wait_for_quality = NessiDataQualitySensor(
        task_id='wait_for_quality',
        check_id="{{ task_instance.xcom_pull(task_ids='run_quality_check')['check_id'] }}",
        quality_threshold=0.9,
        fail_on_rule_failure=True,
        conn_id='nessi_default',
    )

    # Define task dependencies
    run_quality_check >> wait_for_quality
```

## Prefect Integration

### Installation

```bash
pip install nessi-prefect
```

### Configuration

1. Set up Prefect Secret block for Nessi.dev (optional):

```python
from prefect.blocks.system import Secret

Secret(value={"api_key": "your-api-key", "api_secret": "your-api-secret"}).save(name="nessi-credentials")
```

2. Import and use the tasks in your flows:

```python
from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results
from nessi_prefect.tasks.profile import run_profile, wait_for_profile
from nessi_prefect.tasks.validation import run_validation, wait_for_validation_results
from nessi_prefect.tasks.lineage import get_lineage, visualize_lineage
```

### Example Flow

```python
from prefect import flow, task
from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results

@flow(name="Nessi Data Quality Flow")
def data_quality_flow(table_name: str):
    # Run quality check
    quality_response = run_quality_check(
        table_name=table_name,
        rules=[
            {
                "name": "not_null_check",
                "description": "Check that id column is not null",
                "rule_type": "not_null",
                "column": "id",
            },
            {
                "name": "unique_check",
                "description": "Check that id column is unique",
                "rule_type": "unique",
                "column": "id",
            },
        ],
        profile=True,
        api_host="http://your-nessi-instance",
        api_key="your-api-key",
        # Alternatively, use a Secret block:
        # secret_block_name="nessi-credentials",
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

### Using Pre-built Flows

Nessi.dev provides pre-built flows for common use cases:

```python
from nessi_prefect.flows.data_quality_flow import data_quality_flow
from nessi_prefect.flows.data_pipeline_flow import data_pipeline_flow

# Run a simple data quality flow
result = data_quality_flow(
    table_name="my_table",
    rules=[...],
    quality_threshold=0.9,
    api_host="http://your-nessi-instance",
    api_key="your-api-key",
)

# Run a complete data pipeline flow
result = data_pipeline_flow(
    input_table="source_table",
    output_table="target_table",
    quality_rules=[...],
    validation_rules=[...],
    api_host="http://your-nessi-instance",
    api_key="your-api-key",
)
```

## Dagster Integration

### Installation

```bash
pip install nessi-dagster
```

### Configuration

1. Configure the Nessi resource in your Dagster job:

```python
from nessi_dagster.resources import nessi_resource

@job(
    resource_defs={
        "nessi": nessi_resource.configured({
            "api_host": "http://your-nessi-instance",
            "api_key": "your-api-key",
            "api_secret": "your-api-secret",
            "timeout": 300,
        }),
    },
)
def my_job():
    # Job definition
```

2. Import and use the ops in your jobs:

```python
from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results
from nessi_dagster.ops.profile_ops import run_profile, wait_for_profile
from nessi_dagster.ops.validation_ops import run_validation, wait_for_validation_results
from nessi_dagster.ops.lineage_ops import get_lineage, visualize_lineage
```

### Example Job

```python
from dagster import job, op
from nessi_dagster.resources import nessi_resource
from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results

@job(
    resource_defs={
        "nessi": nessi_resource.configured({
            "api_host": "http://your-nessi-instance",
            "api_key": "your-api-key",
            "api_secret": "your-api-secret",
            "timeout": 300,
        }),
    },
)
def data_quality_job():
    # Run quality check
    quality_response = run_quality_check(
        table_name="my_table",
        rules=[
            {
                "name": "not_null_check",
                "description": "Check that id column is not null",
                "rule_type": "not_null",
                "column": "id",
            },
            {
                "name": "unique_check",
                "description": "Check that id column is unique",
                "rule_type": "unique",
                "column": "id",
            },
        ],
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

### Using Pre-built Jobs

Nessi.dev provides pre-built jobs for common use cases:

```python
from nessi_dagster.jobs.data_quality_job import data_quality_job
from nessi_dagster.jobs.data_pipeline_job import data_pipeline_job

# Configure and execute the jobs
data_quality_job.execute_in_process(
    run_config={
        "ops": {
            "run_quality_check": {
                "config": {
                    "table_name": "my_table",
                    "rules": [...],
                }
            },
            "wait_for_quality_results": {
                "config": {
                    "quality_threshold": 0.9,
                }
            }
        }
    }
)
```

## Best Practices

1. **Error Handling**: Always configure appropriate error handling for your workflows. Use the `fail_on_rule_failure` and `quality_threshold` parameters to control when workflows should fail.

2. **Timeouts**: Set appropriate timeouts for operations, especially when dealing with large datasets.

3. **Authentication**: Use API keys for authentication and store them securely using your orchestration tool's secret management capabilities.

4. **Monitoring**: Monitor the execution of your workflows and set up alerts for failures.

5. **Lineage**: Capture lineage information to understand the relationships between your data assets.

## Testing

All workflow orchestration integrations come with comprehensive test suites to ensure reliability and correctness. These tests cover the client communication, operators/tasks/ops, and complete workflows.

### Running Tests

```bash
# For Airflow integration tests
cd integrations/airflow
python -m unittest discover tests

# For Prefect integration tests
cd integrations/prefect
python -m unittest discover tests

# For Dagster integration tests
cd integrations/dagster
python -m unittest discover tests

# For Kubernetes integration tests
cd integrations/kubernetes
python -m unittest discover tests
```

### Test Coverage

The tests cover:

1. **API Communication**: All API calls are correctly formatted and responses are properly handled
2. **Error Handling**: Errors are properly caught and reported
3. **Workflow Logic**: The workflow components (operators, tasks, ops) work as expected
4. **Integration**: The components work together in complete workflows

The tests use mocking extensively to avoid actual API calls during testing, making them fast and reliable.

## Kubernetes Integration

### Installation

```bash
pip install nessi-k8s
```

### Configuration

The Kubernetes integration can be configured in multiple ways:

1. **Configuration file** - YAML or JSON file with all settings
2. **Environment variables** - Prefixed with `NESSI_K8S_`
3. **Direct parameters** - Passed to the client constructor

#### Example Configuration File

```yaml
# nessi-k8s-config.yaml
namespace: nessi
image: nessi/nessi:latest
service_account_name: nessi-service-account
api_host: http://nessi-api:8080
api_key: your-api-key
api_secret: your-api-secret
job_ttl_seconds_after_finished: 3600
job_backoff_limit: 3
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

### Example Usage

```python
from nessi_k8s.client import NessiK8sClient

# Create client
client = NessiK8sClient(
    config_file="nessi-k8s-config.yaml",
    # Or provide configuration directly
    # namespace="nessi",
    # image="nessi/nessi:latest",
    # api_host="http://nessi-api:8080",
    # api_key="your-api-key",
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
        {
            "name": "unique_check",
            "description": "Check that id column is unique",
            "rule_type": "unique",
            "column": "id",
        },
    ],
    profile=True,
    wait_for_completion=True,
)

print(f"Job status: {result}")
```

### Command Line Interface

The package includes a command-line interface for running Nessi operations in Kubernetes:

```bash
# Run a data quality check
python -m nessi_k8s.cli quality-check --table my_table --rules rules.json --namespace nessi

# Create a default configuration file
python -m nessi_k8s.cli create-config --output nessi-k8s-config.yaml
```

## Troubleshooting

### Common Issues

1. **Authentication Failures**: Ensure your API key and secret are correct and properly configured.

2. **Timeouts**: Increase the timeout value for operations on large datasets.

3. **Missing Tables**: Verify that the table names you're using exist in your Nessi.dev instance.

4. **Rule Failures**: Check the rule definitions and ensure they match the expected data format.

5. **Kubernetes Permissions**: Ensure the service account has the necessary permissions to create and manage jobs.

6. **Resource Constraints**: Check that the resource requests and limits are appropriate for your workload.

### Getting Help

If you encounter issues with the workflow orchestration integrations, please:

1. Check the documentation for your specific orchestration tool
2. Review the logs for detailed error messages
3. Contact support at support@nessi.dev
