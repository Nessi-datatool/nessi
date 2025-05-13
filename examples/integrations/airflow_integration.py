"""
Nessi Airflow Integration Example

This example demonstrates how to integrate Nessi with Apache Airflow
to run data quality checks as part of your data pipelines.
"""

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

# Define the paths to your Delta tables
delta_tables = [
    '/path/to/delta/table1',
    '/path/to/delta/table2',
    '/path/to/delta/table3',
]

# Function to run Nessi validation and parse results
def run_nessi_validation(table_path, **kwargs):
    """Run Nessi validation on a Delta table and return the results."""
    try:
        # Run Nessi validation command
        cmd = ['nessi', 'validate', table_path, '--format', 'json']
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        
        # Parse the JSON output
        validation_results = json.loads(result.stdout)
        
        # Log the results
        print(f"Validation results for {table_path}:")
        print(json.dumps(validation_results, indent=2))
        
        # Check if validation passed
        if not validation_results.get('passed', False):
            # Get the failed rules
            failed_rules = [rule for rule in validation_results.get('rules', []) 
                           if not rule.get('passed', False)]
            
            # Log the failed rules
            print(f"Failed rules for {table_path}:")
            for rule in failed_rules:
                print(f"- {rule.get('name')}: {rule.get('message')}")
            
            # Raise an exception if validation failed
            raise Exception(f"Data quality validation failed for {table_path}")
        
        return validation_results
    
    except subprocess.CalledProcessError as e:
        print(f"Error running Nessi validation: {e}")
        print(f"Stderr: {e.stderr}")
        raise
    except json.JSONDecodeError as e:
        print(f"Error parsing Nessi output: {e}")
        print(f"Output: {result.stdout}")
        raise
    except Exception as e:
        print(f"Unexpected error: {e}")
        raise

# Function to send validation results to a webhook
def send_results_to_webhook(table_path, **kwargs):
    """Send validation results to a webhook."""
    ti = kwargs['ti']
    validation_results = ti.xcom_pull(task_ids=f'validate_{os.path.basename(table_path)}')
    
    if not validation_results:
        print(f"No validation results found for {table_path}")
        return
    
    # Here you would send the results to your webhook
    # This is a placeholder for your webhook integration code
    print(f"Sending validation results for {table_path} to webhook")
    
    # Example of using Nessi's webhook command
    try:
        webhook_cmd = [
            'nessi', 'webhook', 'trigger',
            '--event', 'validation.complete',
            '--payload', json.dumps({
                'table': table_path,
                'results': validation_results,
                'timestamp': datetime.now().isoformat()
            })
        ]
        subprocess.run(webhook_cmd, check=True)
        print(f"Webhook triggered successfully for {table_path}")
    except subprocess.CalledProcessError as e:
        print(f"Error triggering webhook: {e}")
        # Don't fail the task if webhook fails
        pass

# Create tasks for each Delta table
for table_path in delta_tables:
    table_name = os.path.basename(table_path)
    
    # Task to run Nessi validation
    validate_task = PythonOperator(
        task_id=f'validate_{table_name}',
        python_callable=run_nessi_validation,
        op_kwargs={'table_path': table_path},
        dag=dag,
    )
    
    # Task to send results to webhook
    webhook_task = PythonOperator(
        task_id=f'webhook_{table_name}',
        python_callable=send_results_to_webhook,
        op_kwargs={'table_path': table_path},
        provide_context=True,
        dag=dag,
    )
    
    # Set task dependencies
    validate_task >> webhook_task

# Task to generate a summary report
generate_report_task = BashOperator(
    task_id='generate_summary_report',
    bash_command='nessi report generate --output-dir /path/to/reports --format html',
    dag=dag,
)

# Set final task dependencies
for table_path in delta_tables:
    table_name = os.path.basename(table_path)
    dag.get_task(f'webhook_{table_name}') >> generate_report_task
