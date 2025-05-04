# Apache Airflow Integration

This guide shows how to integrate Nessi.dev with Apache Airflow for data quality monitoring and pipeline orchestration.

## Example DAG

Create a file `dags/nessi_quality_check.py`:

```python
from datetime import datetime, timedelta
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from nessi import NessiClient

default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

def check_data_quality(**context):
    """Run data quality checks using Nessi.dev."""
    client = NessiClient(api_key='your-api-key')
    
    # Run quality check
    result = client.quality.check(
        table_path='s3://bucket/table',
        rules={
            'completeness': 0.95,
            'accuracy': 0.98,
            'consistency': 0.97
        }
    )
    
    # Check quality score
    if result.score < 0.9:
        raise ValueError(f"Quality score below threshold: {result.score}")
    
    # Get detailed report
    report = client.quality.report(
        table_path='s3://bucket/table',
        format='html'
    )
    
    # Save report
    with open('/tmp/quality_report.html', 'w') as f:
        f.write(report)
    
    return result.score

def monitor_metrics(**context):
    """Monitor system metrics using Nessi.dev."""
    client = NessiClient(api_key='your-api-key')
    
    # Get current metrics
    metrics = client.metrics.get()
    
    # Check for alerts
    alerts = client.alerts.get()
    if alerts:
        raise ValueError(f"Active alerts found: {alerts}")
    
    return metrics

with DAG(
    'nessi_quality_check',
    default_args=default_args,
    description='Data quality monitoring with Nessi.dev',
    schedule_interval='@daily',
    start_date=datetime(2024, 1, 1),
    catchup=False,
) as dag:
    
    check_quality = PythonOperator(
        task_id='check_quality',
        python_callable=check_data_quality,
        provide_context=True,
    )
    
    monitor = PythonOperator(
        task_id='monitor_metrics',
        python_callable=monitor_metrics,
        provide_context=True,
    )
    
    upload_report = BashOperator(
        task_id='upload_report',
        bash_command='aws s3 cp /tmp/quality_report.html s3://reports/quality_report.html',
    )
    
    check_quality >> monitor >> upload_report
```

## Configuration

1. Install the Nessi.dev Python package:
```bash
pip install nessi
```

2. Configure Airflow connections:
```bash
airflow connections add nessi_api \
    --conn-type http \
    --conn-host https://api.nessi.dev \
    --conn-extra '{"api_key": "your-api-key"}'
```

## Advanced Usage

### Custom Operators

Create custom Airflow operators for Nessi.dev:

```python
from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults

class NessiQualityCheckOperator(BaseOperator):
    """Custom operator for Nessi.dev quality checks."""
    
    @apply_defaults
    def __init__(
        self,
        table_path: str,
        rules: dict,
        *args,
        **kwargs
    ):
        super().__init__(*args, **kwargs)
        self.table_path = table_path
        self.rules = rules
    
    def execute(self, context):
        client = NessiClient(api_key='your-api-key')
        result = client.quality.check(
            table_path=self.table_path,
            rules=self.rules
        )
        return result.score

# Use in DAG
check_quality = NessiQualityCheckOperator(
    task_id='check_quality',
    table_path='s3://bucket/table',
    rules={
        'completeness': 0.95,
        'accuracy': 0.98
    }
)
```

### Dynamic Task Generation

Generate tasks dynamically based on table configurations:

```python
def create_quality_tasks(dag):
    """Create quality check tasks for multiple tables."""
    tables = [
        {'path': 's3://bucket/table1', 'rules': {'completeness': 0.95}},
        {'path': 's3://bucket/table2', 'rules': {'completeness': 0.98}},
    ]
    
    tasks = []
    for table in tables:
        task = NessiQualityCheckOperator(
            task_id=f'check_quality_{table["path"].split("/")[-1]}',
            table_path=table['path'],
            rules=table['rules'],
            dag=dag
        )
        tasks.append(task)
    
    return tasks

# Use in DAG
quality_tasks = create_quality_tasks(dag)
```

### Error Handling and Notifications

Add error handling and notifications:

```python
from airflow.operators.slack_operator import SlackAPIPostOperator

def handle_quality_failure(context):
    """Handle quality check failures."""
    task_instance = context['task_instance']
    error = task_instance.error
    
    slack_alert = SlackAPIPostOperator(
        task_id='slack_alert',
        token='your-slack-token',
        channel='#alerts',
        text=f'Quality check failed: {error}',
        dag=dag
    )
    
    return slack_alert.execute(context)

# Add to DAG
check_quality.on_failure_callback = handle_quality_failure
```

### Metrics Visualization

Create a custom view for quality metrics:

```python
from airflow.www.views import AirflowViewMixin
from flask import Blueprint, render_template

class NessiMetricsView(AirflowViewMixin):
    """Custom view for Nessi.dev metrics."""
    
    @app.route('/nessi/metrics')
    def metrics(self):
        client = NessiClient(api_key='your-api-key')
        metrics = client.metrics.get()
        return render_template(
            'nessi_metrics.html',
            metrics=metrics
        )

# Register view
admin_view = NessiMetricsView(category="Nessi", name="Metrics")
``` 