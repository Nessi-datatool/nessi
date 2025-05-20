"""
Example DAG demonstrating Nessi.dev integration with Apache Airflow.

This DAG shows how to use Nessi.dev operators and sensors to perform data quality checks,
profiling, validation, and lineage tracking as part of a data pipeline.
"""

from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.dummy import DummyOperator
from airflow.operators.python import PythonOperator
from airflow.utils.dates import days_ago

from nessi_airflow.operators.data_quality_operator import NessiDataQualityOperator
from nessi_airflow.operators.profile_operator import NessiProfileOperator
from nessi_airflow.operators.validation_operator import NessiValidationOperator
from nessi_airflow.operators.lineage_operator import NessiLineageOperator
from nessi_airflow.sensors.data_quality_sensor import NessiDataQualitySensor


# Default arguments for the DAG
default_args = {
    'owner': 'nessi',
    'depends_on_past': False,
    'email': ['data-quality@example.com'],
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

# Define the DAG
dag = DAG(
    'nessi_data_quality_example',
    default_args=default_args,
    description='Example DAG for Nessi.dev data quality integration',
    schedule_interval=timedelta(days=1),
    start_date=days_ago(1),
    tags=['nessi', 'data_quality', 'example'],
)

# Define tasks

# Start task
start = DummyOperator(
    task_id='start',
    dag=dag,
)

# Data ingestion task (simulated)
def simulate_data_ingestion(**kwargs):
    """Simulate data ingestion."""
    print("Simulating data ingestion...")
    # In a real scenario, this would be a task that loads data into a table
    return {'table_name': 'example_sales_data'}

ingest_data = PythonOperator(
    task_id='ingest_data',
    python_callable=simulate_data_ingestion,
    dag=dag,
)

# Profile the data
profile_data = NessiProfileOperator(
    task_id='profile_data',
    table_name="{{ ti.xcom_pull(task_ids='ingest_data')['table_name'] }}",
    sample_size=10000,  # Sample 10,000 rows for profiling
    timeout=300,  # 5 minutes timeout
    dag=dag,
)

# Run data quality checks
quality_check = NessiDataQualityOperator(
    task_id='quality_check',
    table_name="{{ ti.xcom_pull(task_ids='ingest_data')['table_name'] }}",
    rules=[
        {
            "name": "sales_amount_positive",
            "description": "Sales amount should be positive",
            "rule_type": "range",
            "column": "sales_amount",
            "min_value": 0,
        },
        {
            "name": "transaction_date_valid",
            "description": "Transaction date should be valid",
            "rule_type": "not_null",
            "column": "transaction_date",
        },
        {
            "name": "customer_id_unique",
            "description": "Customer ID should be unique",
            "rule_type": "unique",
            "column": "customer_id",
        },
    ],
    quality_threshold=0.9,  # 90% quality score threshold
    fail_on_rule_failure=True,
    dag=dag,
)

# Run data validation
validate_data = NessiValidationOperator(
    task_id='validate_data',
    table_name="{{ ti.xcom_pull(task_ids='ingest_data')['table_name'] }}",
    rules=[
        {
            "name": "region_in_allowed_values",
            "description": "Region should be in allowed values",
            "rule_type": "enum",
            "column": "region",
            "allowed_values": ["North", "South", "East", "West"],
        },
        {
            "name": "product_code_pattern",
            "description": "Product code should match pattern",
            "rule_type": "regex",
            "column": "product_code",
            "pattern": "^[A-Z]{2}-\\d{4}$",
        },
    ],
    dag=dag,
)

# Capture lineage
capture_lineage = NessiLineageOperator(
    task_id='capture_lineage',
    table_name="{{ ti.xcom_pull(task_ids='ingest_data')['table_name'] }}",
    push_airflow_lineage=True,
    format='json',
    dag=dag,
)

# Data transformation task (simulated)
def simulate_data_transformation(**kwargs):
    """Simulate data transformation."""
    print("Simulating data transformation...")
    # Pull quality results from XCom
    ti = kwargs['ti']
    quality_results = ti.xcom_pull(task_ids='quality_check')
    
    # Log quality score
    quality_score = quality_results.get('quality_score', 0)
    print(f"Quality score: {quality_score}")
    
    # In a real scenario, this would transform the data based on quality results
    return {'transformed_table': 'example_sales_data_transformed'}

transform_data = PythonOperator(
    task_id='transform_data',
    python_callable=simulate_data_transformation,
    provide_context=True,
    dag=dag,
)

# Run quality check on transformed data
quality_check_transformed = NessiDataQualityOperator(
    task_id='quality_check_transformed',
    table_name="{{ ti.xcom_pull(task_ids='transform_data')['transformed_table'] }}",
    rules=[
        {
            "name": "sales_amount_positive",
            "description": "Sales amount should be positive",
            "rule_type": "range",
            "column": "sales_amount",
            "min_value": 0,
        },
        {
            "name": "all_dates_valid",
            "description": "All dates should be valid",
            "rule_type": "not_null",
            "column": "transaction_date",
        },
    ],
    quality_threshold=0.95,  # Higher threshold for transformed data
    dag=dag,
)

# Data export task (simulated)
def simulate_data_export(**kwargs):
    """Simulate data export."""
    print("Simulating data export...")
    # In a real scenario, this would export the data to a destination
    return {'export_id': '12345'}

export_data = PythonOperator(
    task_id='export_data',
    python_callable=simulate_data_export,
    dag=dag,
)

# End task
end = DummyOperator(
    task_id='end',
    dag=dag,
)

# Define task dependencies
start >> ingest_data >> [profile_data, quality_check, validate_data]
[profile_data, quality_check, validate_data] >> capture_lineage >> transform_data
transform_data >> quality_check_transformed >> export_data >> end

# Example of using a sensor (commented out as it requires an existing check_id)
"""
wait_for_quality_check = NessiDataQualitySensor(
    task_id='wait_for_quality_check',
    check_id='existing-check-id',  # This would be a pre-existing check ID
    quality_threshold=0.9,
    fail_on_rule_failure=True,
    poke_interval=60,  # Check every minute
    timeout=3600,  # Timeout after 1 hour
    mode='poke',
    dag=dag,
)
"""
