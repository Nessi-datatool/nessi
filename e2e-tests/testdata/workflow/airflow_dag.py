from airflow import DAG
from airflow.operators.bash import BashOperator
from datetime import datetime, timedelta

default_args = {
    'owner': 'nessi',
    'depends_on_past': False,
    'start_date': datetime(2023, 1, 1),
    'email': ['alerts@nessi.dev'],
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

dag = DAG(
    'nessi_data_quality',
    default_args=default_args,
    description='Nessi data quality checks for Delta Lake tables',
    schedule_interval=timedelta(days=1),
    catchup=False,
)

scan_users_table = BashOperator(
    task_id='scan_users_table',
    bash_command='nessi scan users --output-format=json --output-file=/tmp/users_scan_results.json',
    dag=dag,
)

validate_users_schema = BashOperator(
    task_id='validate_users_schema',
    bash_command='nessi validate schema users --schema-file=/path/to/schemas/user_schema_v1.json --output-format=json --output-file=/tmp/users_schema_validation.json',
    dag=dag,
)

apply_validation_rules = BashOperator(
    task_id='apply_validation_rules',
    bash_command='nessi validate rules users --rules-file=/path/to/configs/validation_rules.json --output-format=json --output-file=/tmp/users_rule_validation.json',
    dag=dag,
)

generate_report = BashOperator(
    task_id='generate_report',
    bash_command='nessi report generate --input-files=/tmp/users_scan_results.json,/tmp/users_schema_validation.json,/tmp/users_rule_validation.json --output-format=html --output-file=/tmp/users_quality_report.html',
    dag=dag,
)

scan_users_table >> validate_users_schema >> apply_validation_rules >> generate_report
