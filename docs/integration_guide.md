# Nessi Integration Guide

This guide explains how to integrate Nessi with external systems such as Apache Airflow, GitHub Actions, and custom scripts. Nessi is designed to be easily embeddable in various environments, allowing you to incorporate data quality checks into your existing workflows.

## Table of Contents

- [Workflow Orchestration Integrations](#workflow-orchestration-integrations)
  - [Airflow Integration](#airflow-integration)
  - [Prefect Integration](#prefect-integration)
  - [Dagster Integration](#dagster-integration)
  - [Kubernetes Integration](#kubernetes-integration)
- [GitHub Actions Integration](#github-actions-integration)
- [Custom Script Integration](#custom-script-integration)
  - [Python Scripts](#python-scripts)
  - [Shell Scripts](#shell-scripts)
- [Viral Growth Features Integration](#viral-growth-features-integration)
  - [Social Sharing Integration](#social-sharing-integration)
  - [Badge Integration](#badge-integration)
  - [Community Engagement Integration](#community-engagement-integration)
- [Integration Best Practices](#integration-best-practices)
- [Troubleshooting](#troubleshooting)

## Workflow Orchestration Integrations

Nessi.dev provides dedicated integrations for popular workflow orchestration tools and container orchestration platforms, allowing you to incorporate data quality checks, profiling, validation, and lineage tracking into your data pipelines. These integrations are more powerful and easier to use than the basic script-based approaches.

For detailed documentation on these integrations, see the [Workflow Orchestration Integration Guide](./workflow_orchestration.md).

## Airflow Integration

Apache Airflow is a popular platform for programmatically authoring, scheduling, and monitoring workflows. You can integrate Nessi into your Airflow DAGs to run data quality checks as part of your data pipelines.

### Native Airflow Integration

Nessi.dev provides a native Airflow integration package with custom operators and sensors:

```bash
pip install nessi-airflow
```

```python
from nessi_airflow.operators.data_quality_operator import NessiDataQualityOperator
from nessi_airflow.operators.profile_operator import NessiProfileOperator
from nessi_airflow.operators.validation_operator import NessiValidationOperator
from nessi_airflow.operators.lineage_operator import NessiLineageOperator
from nessi_airflow.sensors.data_quality_sensor import NessiDataQualitySensor

# Create a data quality check task
quality_check = NessiDataQualityOperator(
    task_id='run_quality_check',
    table_name='my_table',
    rules=[{'name': 'not_null_check', 'rule_type': 'not_null', 'column': 'id'}],
    quality_threshold=0.9,
    conn_id='nessi_default',
    dag=dag,
)

# Create a sensor to wait for the quality check to complete
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

For more details, see the [Workflow Orchestration Integration Guide](./workflow_orchestration.md#apache-airflow-integration).

### Script-Based Airflow Integration

### Example Airflow DAG

```python
from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.bash import BashOperator
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago
import json
import subprocess
import os

# Define default arguments for the DAG
default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

# Create the DAG
dag = DAG(
    'nessi_data_quality_checks',
    default_args=default_args,
    description='Run Nessi data quality checks on Delta tables',
    schedule_interval=timedelta(days=1),
    start_date=days_ago(1),
    tags=['nessi', 'data_quality'],
)

# Function to run Nessi validation and parse results
def run_nessi_validation(table_path, **kwargs):
    """Run Nessi validation on a Delta table and return the results."""
    try:
        # Run Nessi validation command
        cmd = ['nessi', 'validate', table_path, '--format', 'json']
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        
        # Parse the JSON output
        validation_results = json.loads(result.stdout)
        
        # Check if validation passed
        if not validation_results.get('passed', False):
            # Raise an exception if validation failed
            raise Exception(f"Data quality validation failed for {table_path}")
        
        return validation_results
    
    except Exception as e:
        print(f"Error: {e}")
        raise

# Create tasks for each Delta table
for table_path in ['/path/to/delta/table1', '/path/to/delta/table2']:
    table_name = os.path.basename(table_path)
    
    # Task to run Nessi validation
    validate_task = PythonOperator(
        task_id=f'validate_{table_name}',
        python_callable=run_nessi_validation,
        op_kwargs={'table_path': table_path},
        dag=dag,
    )
```

For a complete example, see [airflow_integration.py](../examples/integrations/airflow_integration.py).

### Airflow Integration Tips

1. **Use the native integration package**: The `nessi-airflow` package provides a more robust and maintainable integration than script-based approaches.
2. **Install Nessi in your Airflow environment**: Make sure Nessi is installed and accessible in your Airflow worker environment.
3. **Use XComs for passing results**: Use Airflow's XComs to pass validation results between tasks.
4. **Handle failures appropriately**: Configure your DAG to handle validation failures according to your requirements.
5. **Use Airflow variables for configuration**: Store Nessi configuration in Airflow variables for flexibility.


## Prefect Integration

Prefect is a modern workflow orchestration tool designed for data engineers and scientists. Nessi.dev provides a native Prefect integration package with tasks and flows for data quality checks, profiling, validation, and lineage tracking.

### Installation

```bash
pip install nessi-prefect
```

### Example Prefect Flow

```python
from prefect import flow
from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results

@flow(name="Data Quality Flow")
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
        ],
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

For more details, see the [Workflow Orchestration Integration Guide](./workflow_orchestration.md#prefect-integration).

### Prefect Integration Tips

1. **Use the native integration package**: The `nessi-prefect` package provides a more robust and maintainable integration than script-based approaches.
2. **Leverage Prefect Secret blocks**: Store API keys and other sensitive information in Prefect Secret blocks.
3. **Use Prefect deployments**: Deploy your flows to run on a schedule or in response to events.
4. **Configure appropriate retries**: Set up retries for tasks that might fail due to transient issues.
5. **Use Prefect notifications**: Configure notifications for flow run failures.

## Dagster Integration

Dagster is a data orchestrator for machine learning, analytics, and ETL. Nessi.dev provides a native Dagster integration package with resources, ops, and jobs for data quality checks, profiling, validation, and lineage tracking.

### Installation

```bash
pip install nessi-dagster
```

### Example Dagster Job

```python
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
        rules=[
            {
                "name": "not_null_check",
                "description": "Check that id column is not null",
                "rule_type": "not_null",
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

## Kubernetes Integration

Kubernetes is a powerful container orchestration platform. Nessi.dev provides a dedicated Kubernetes integration that allows you to run data quality operations as Kubernetes jobs.

### Installation

```bash
pip install nessi-k8s
```

### Native Kubernetes Integration

Nessi.dev provides a native Kubernetes integration package with operators for running data quality checks, profiling, and validation as Kubernetes jobs:

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

### Kubernetes Manifest Example

You can also create Kubernetes manifests directly to run Nessi operations:

```yaml
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

### Command Line Interface

The Kubernetes integration also provides a command-line interface:

```bash
# Run a data quality check
python -m nessi_k8s.cli quality-check --table my_table --rules rules.json --namespace nessi

# Create a default configuration file
python -m nessi_k8s.cli create-config --output nessi-k8s-config.yaml
```

For more details, see the [Workflow Orchestration Integration Guide](./workflow_orchestration.md#kubernetes-integration).

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

For more details, see the [Workflow Orchestration Integration Guide](./workflow_orchestration.md#dagster-integration).

### Dagster Integration Tips

1. **Use the native integration package**: The `nessi-dagster` package provides a more robust and maintainable integration than script-based approaches.
2. **Leverage Dagster resources**: Use the Nessi resource to configure the API connection once and reuse it across ops.
3. **Use Dagster schedules**: Schedule your jobs to run on a regular basis.
4. **Configure appropriate retries**: Set up retries for ops that might fail due to transient issues.


## GitHub Actions Integration

GitHub Actions allows you to automate your software workflows directly in your GitHub repository. You can use Nessi with GitHub Actions to run data quality checks as part of your CI/CD pipeline.

### Example GitHub Actions Workflow

```yaml
name: Nessi Data Quality Checks

on:
  schedule:
    # Run daily at midnight
    - cron: '0 0 * * *'
  workflow_dispatch:
    # Allow manual triggering
  push:
    branches: [ main ]
    paths:
      # Only run when data files change
      - 'data/**'

jobs:
  data-quality-checks:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
          
      - name: Install Nessi
        run: |
          # Clone Nessi repository
          git clone git@github.com:nessi-dev/nessi.git
          cd nessi
          
          # Build and install Nessi
          make build
          sudo cp bin/nessi /usr/local/bin/
      
      - name: Run data quality checks
        run: |
          # Run checks on all Delta tables in the data directory
          for table in data/*/; do
            echo "Running checks on $table"
            nessi validate "$table" --format json > "${table//\//_}_results.json"
          done
      
      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: data-quality-results
          path: |
            *_results.json
            report.html
```

For a complete example, see [github_actions_workflow.yml](../examples/integrations/github_actions_workflow.yml).

### GitHub Actions Integration Tips

1. **Cache Nessi installation**: Use GitHub Actions caching to speed up your workflow.
2. **Use artifacts for results**: Upload validation results as artifacts for later analysis.
3. **Configure appropriate triggers**: Set up your workflow to run on schedule, on push, or manually.
4. **Use GitHub Secrets for sensitive information**: Store API keys and other sensitive information in GitHub Secrets.
5. **Set up notifications**: Configure GitHub Actions to send notifications when data quality checks fail.

## Custom Script Integration

You can integrate Nessi with your own custom scripts to run data quality checks as part of your workflows.

### Python Scripts

Here's a simple example of integrating Nessi with a Python script:

```python
import subprocess
import json
import sys

def run_nessi_validation(table_path):
    """Run Nessi validation on a Delta table and return the results."""
    try:
        # Run Nessi validation command
        cmd = ['nessi', 'validate', table_path, '--format', 'json']
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        
        # Parse the JSON output
        validation_results = json.loads(result.stdout)
        
        return validation_results
    
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

# Run validation on a Delta table
results = run_nessi_validation('/path/to/delta/table')

# Check if validation passed
if results.get('passed', False):
    print("Validation passed!")
else:
    print("Validation failed!")
    sys.exit(1)
```

For a complete example, see [custom_script_integration.py](../examples/integrations/custom_script_integration.py).

### Shell Scripts

Here's a simple example of integrating Nessi with a shell script:

```bash
#!/bin/bash

# Run Nessi validation on a Delta table
echo "Running validation on /path/to/delta/table"
nessi validate /path/to/delta/table --format json > results.json

# Check if validation passed
passed=$(jq -r '.passed' results.json)
if [ "$passed" = "true" ]; then
    echo "Validation passed!"
    exit 0
else
    echo "Validation failed!"
    exit 1
fi
```

For a complete example, see [shell_script_integration.sh](../examples/integrations/shell_script_integration.sh).

## Integration Best Practices



When integrating Nessi with external systems, follow these best practices:

1. **Use JSON output format**: Nessi's JSON output format is easy to parse and process programmatically.
2. **Handle errors gracefully**: Make sure your integration handles errors and exceptions properly.
3. **Use appropriate exit codes**: Return non-zero exit codes when validation fails to signal failure to the calling system.
4. **Configure timeouts**: Set appropriate timeouts for Nessi commands to prevent hanging processes.
5. **Use webhook integration**: Use Nessi's webhook feature to notify external systems of validation results.
6. **Log validation results**: Store validation results for historical analysis and trend detection.
7. **Use configuration files**: Use Nessi configuration files to define validation rules and thresholds.
8. **Automate report generation**: Generate reports automatically after validation for easy sharing and analysis.
9. **Implement observability**: Add metrics collection and structured logging for better monitoring and debugging.
10. **Use resource constraints**: Set appropriate resource requests and limits for Kubernetes jobs.

## Viral Growth Features Integration

Nessi provides several viral growth features that can be integrated into your existing workflows to increase visibility, drive adoption, and foster community engagement. These features are implemented as plugins and can be accessed through the `nessi viral` command group.

### Social Sharing Integration

The social sharing plugin enhances Nessi reports with advanced sharing capabilities, making it easy to share quality reports with stakeholders and the broader data community.

#### Integrating with CI/CD Pipelines

You can automatically generate and share reports as part of your CI/CD pipeline:

```yaml
# Example GitHub Actions workflow
name: Data Quality Check and Share

on:
  schedule:
    - cron: '0 8 * * 1' # Weekly on Monday at 8am

jobs:
  quality_check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install Nessi
        run: pip install nessi
      - name: Run Nessi quality check and generate shareable report
        run: |
          nessi quality check --table my_table
          nessi viral share my_table --output report.html
      - name: Upload report
        uses: actions/upload-artifact@v2
        with:
          name: quality-report
          path: report.html
```

#### Integrating with Slack

Automatically share reports with your team via Slack:

```bash
#!/bin/bash

# Generate shareable report
nessi viral share my_table --output report.html

# Share link to Slack
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"New data quality report available: https://example.com/reports/report.html"}' \
  $SLACK_WEBHOOK_URL
```

### Badge Integration

The badge plugin generates embeddable "Powered by Nessi" badges that can be added to your project's README, documentation, or website.

#### Integrating with GitHub READMEs

Add a badge to your project's README.md file:

```bash
# Generate a badge in markdown format
badge=$(nessi viral badge --format markdown)

# Add the badge to the README.md file
sed -i "1s/^/$badge\n\n/" README.md
```

#### Integrating with Documentation Sites

Add badges to your documentation site:

```python
import subprocess
import re

# Generate a badge in HTML format
result = subprocess.run(["nessi", "viral", "badge", "--format", "html"], capture_output=True, text=True)
badge_html = result.stdout.strip()

# Add the badge to your documentation site
with open('docs/index.html', 'r') as file:
    content = file.read()

# Insert the badge after the opening body tag
content = re.sub(r'<body>', f'<body>\n{badge_html}', content)

with open('docs/index.html', 'w') as file:
    file.write(content)
```

### Community Engagement Integration

The community engagement plugin facilitates feedback collection and contribution suggestions, helping to foster community contributions.

#### Integrating with Issue Tracking Systems

Automatically create GitHub issues from user feedback:

```python
import subprocess
import json
import requests

# Collect feedback using Nessi
result = subprocess.run(
    ["nessi", "viral", "community", "feedback", "--type", "feature", "--text", "It would be great to have..."],
    capture_output=True, text=True
)

# Extract the GitHub issue URL from the output
output_lines = result.stdout.strip().split('\n')
for line in output_lines:
    if line.startswith('https://github.com/'):
        issue_url = line.strip()
        break

# Create the issue using the GitHub API
if issue_url:
    # Extract the repo and issue details from the URL
    parts = issue_url.replace('https://github.com/', '').split('/')
    repo_owner = parts[0]
    repo_name = parts[1]
    
    # Create the issue
    github_token = 'your_github_token'
    headers = {
        'Authorization': f'token {github_token}',
        'Accept': 'application/vnd.github.v3+json'
    }
    
    issue_data = {
        'title': 'Feature Request from Nessi',
        'body': 'It would be great to have...'
    }
    
    response = requests.post(
        f'https://api.github.com/repos/{repo_owner}/{repo_name}/issues',
        headers=headers,
        json=issue_data
    )
```

For more information on Nessi's viral growth features, see the [Viral Growth Guide](./VIRAL_GROWTH.md) and the [Plugin API Documentation](./PLUGIN_API.md).

## Troubleshooting

If you encounter issues when integrating Nessi with external systems, try the following:

1. **Check Nessi installation**: Make sure Nessi is installed and accessible in your environment.
2. **Verify permissions**: Ensure that Nessi has the necessary permissions to access your Delta tables.
3. **Check command syntax**: Verify that you're using the correct Nessi command syntax.
4. **Inspect error messages**: Look at the error messages returned by Nessi for clues about the issue.
5. **Enable debug logging**: Run Nessi with debug logging enabled to get more information.
6. **Check for version compatibility**: Ensure that you're using a compatible version of Nessi.
7. **Test commands manually**: Try running Nessi commands manually to verify they work as expected.

For more help, refer to the [Nessi documentation](https://github.com/nessi-dev/nessi) or open an issue on GitHub.
