# Nessi Integration Guide

This guide explains how to integrate Nessi with external systems such as Apache Airflow, GitHub Actions, and custom scripts. Nessi is designed to be easily embeddable in various environments, allowing you to incorporate data quality checks into your existing workflows.

## Table of Contents

- [Airflow Integration](#airflow-integration)
- [GitHub Actions Integration](#github-actions-integration)
- [Custom Script Integration](#custom-script-integration)
  - [Python Scripts](#python-scripts)
  - [Shell Scripts](#shell-scripts)
- [Integration Best Practices](#integration-best-practices)
- [Troubleshooting](#troubleshooting)

## Airflow Integration

Apache Airflow is a popular platform for programmatically authoring, scheduling, and monitoring workflows. You can integrate Nessi into your Airflow DAGs to run data quality checks as part of your data pipelines.

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

1. **Install Nessi in your Airflow environment**: Make sure Nessi is installed and accessible in your Airflow worker environment.
2. **Use XComs for passing results**: Use Airflow's XComs to pass validation results between tasks.
3. **Handle failures appropriately**: Configure your DAG to handle validation failures according to your requirements.
4. **Use Airflow variables for configuration**: Store Nessi configuration in Airflow variables for flexibility.
5. **Set up alerting**: Configure Airflow to send alerts when data quality checks fail.

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
          git clone https://github.com/nessi-dev/nessi-dev.git
          cd nessi-dev
          
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

## Troubleshooting

If you encounter issues when integrating Nessi with external systems, try the following:

1. **Check Nessi installation**: Make sure Nessi is installed and accessible in your environment.
2. **Verify permissions**: Ensure that Nessi has the necessary permissions to access your Delta tables.
3. **Check command syntax**: Verify that you're using the correct Nessi command syntax.
4. **Inspect error messages**: Look at the error messages returned by Nessi for clues about the issue.
5. **Enable debug logging**: Run Nessi with debug logging enabled to get more information.
6. **Check for version compatibility**: Ensure that you're using a compatible version of Nessi.
7. **Test commands manually**: Try running Nessi commands manually to verify they work as expected.

For more help, refer to the [Nessi documentation](https://github.com/nessi-dev/nessi-dev) or open an issue on GitHub.
