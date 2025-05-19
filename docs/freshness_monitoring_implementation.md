# Freshness Monitoring Implementation

## Overview

The Freshness Monitoring system in Nessi.dev provides a comprehensive solution for tracking and monitoring the freshness of Delta Lake tables. It allows users to define Service Level Agreements (SLAs) for data freshness and monitor compliance with these SLAs over time.

## Core Components

### SLA Configuration

The `SLAConfig` structure defines the freshness expectations for a table:

```go
type SLAConfig struct {
    TableName         string        // Name of the table
    TablePath         string        // Path to the table
    ExpectedFrequency time.Duration // Expected update frequency
    WarningThreshold  float64       // Warning threshold as percentage of expected frequency
    CriticalThreshold float64       // Critical threshold as percentage of expected frequency
    Metadata          map[string]string // Additional metadata
    Enabled           bool          // Whether this SLA is enabled
}
```

Key validation rules:
- Table name and path cannot be empty
- Expected frequency must be greater than zero
- Warning threshold must be at least 100% (equal to expected frequency)
- Critical threshold must be greater than warning threshold

### SLA Manager

The `SLAManager` is the central component responsible for:
- Managing SLA configurations (add, update, delete)
- Checking freshness status of tables
- Tracking freshness history
- Generating compliance reports and trends

### Freshness Status

The `FreshnessStatus` structure represents the current freshness state of a table:

```go
type FreshnessStatus struct {
    TableName      string    // Name of the table
    TablePath      string    // Path to the table
    LastUpdateTime time.Time // Last update time of the table
    TimeSinceUpdate time.Duration // Time since last update
    ExpectedFrequency time.Duration // Expected update frequency
    Status         SLALevel  // Current status (info, warning, critical)
}
```

### History Tracking

The system maintains a history of freshness checks, allowing for trend analysis and compliance reporting:

```go
type FreshnessHistoryEntry struct {
    TableName         string
    TablePath         string
    Timestamp         time.Time
    LastUpdateTime    time.Time
    TimeSinceUpdate   time.Duration
    ExpectedFrequency time.Duration
    Status            string
}
```

### Trends and Compliance

The system provides trend analysis and compliance reporting:

```go
type FreshnessTrends struct {
    TableName string
    TablePath string
    History   []FreshnessHistoryEntry
    Compliance FreshnessCompliance
}

type FreshnessCompliance struct {
    InfoCount     int
    WarningCount  int
    CriticalCount int
    TotalCount    int
    ComplianceRate float64
}
```

## Integration with Other Nessi.dev Components

### Delta Lake Integration

The freshness monitoring system integrates with the Delta Lake connector to retrieve table metadata, including the last modified time. This integration supports:

- Local Delta Lake tables
- Cloud-hosted Delta Lake tables (via the CloudDeltaConnector)
- Time travel operations for historical analysis

### Output Formats

The freshness monitoring system provides multiple output formats for reporting:

- JSON output for programmatic consumption and integration with other tools
- Table format for terminal viewing with color-coded status indicators
- CSV export for further analysis and reporting

### CLI Integration

The system provides CLI commands for managing SLAs and checking freshness status:

```
nessi freshness sla list
nessi freshness sla add --table <table_name> --path <table_path> --frequency <frequency>
nessi freshness sla delete --table <table_name>
nessi freshness status --table <table_name>
nessi freshness status --all
```

### Security Integration

The freshness monitoring system integrates with the Nessi.dev security system:

- CLI commands are protected by API key authentication
- Role-based access control for SLA management
- Audit logging for SLA changes and status checks

### Plugin System Integration

The freshness monitoring system supports extension through the plugin system:

- Custom freshness check logic via ValidationPlugin
- Custom alerting for freshness violations via AlertPlugin
- Custom storage backends for SLA configurations and history via StoragePlugin

## Configuration Options

### Duration Formats

The system supports various duration formats:

- Standard Go duration formats (e.g., "1h", "30m")
- Custom formats for longer durations:
  - "1d" for 1 day
  - "1w" for 1 week
  - "1mo" for 1 month
- Named frequencies:
  - "hourly" for 1 hour
  - "daily" for 1 day
  - "weekly" for 1 week
  - "monthly" for 1 month

### Threshold Configuration

- Warning threshold: Percentage of expected frequency (must be at least 100%)
- Critical threshold: Percentage of expected frequency (must be greater than warning threshold)

## Python API Integration

The freshness monitoring system is accessible through the Nessi.dev Python API:

```python
from nessi import Client

# Initialize client
client = Client(host="localhost:8080", api_key="your-api-key")

# List SLAs
slas = client.freshness.list_slas()

# Add SLA
client.freshness.add_sla(
    table_name="my_table",
    table_path="/path/to/my_table",
    expected_frequency="1d",
    warning_threshold=150,
    critical_threshold=200
)

# Check freshness
status = client.freshness.check_freshness("my_table")

# Get trends
trends = client.freshness.get_trends("my_table")
```

## Workflow Integration

The freshness monitoring system can be integrated into workflow orchestration tools:

### Airflow Integration

```python
from airflow import DAG
from airflow.operators.python_operator import PythonOperator
from datetime import datetime, timedelta
from nessi import Client

def check_freshness():
    client = Client(host="localhost:8080", api_key="your-api-key")
    statuses = client.freshness.check_all_freshness()
    for status in statuses:
        if status.status == "critical":
            raise Exception(f"Table {status.table_name} is not fresh!")

dag = DAG(
    'freshness_check',
    default_args={
        'owner': 'airflow',
        'depends_on_past': False,
        'start_date': datetime(2023, 1, 1),
        'email_on_failure': True,
        'email_on_retry': False,
        'retries': 1,
        'retry_delay': timedelta(minutes=5),
    },
    schedule_interval=timedelta(days=1),
)

check_task = PythonOperator(
    task_id='check_freshness',
    python_callable=check_freshness,
    dag=dag,
)
```

### GitHub Actions Integration

```yaml
name: Data Freshness Check

on:
  schedule:
    - cron: '0 8 * * *'  # Run daily at 8 AM
  workflow_dispatch:  # Allow manual triggering

jobs:
  check-freshness:
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@v2

      - name: Set up Python
        uses: actions/setup-python@v2
        with:
          python-version: '3.9'

      - name: Install dependencies
        run: |
          python -m pip install --upgrade pip
          pip install nessi-client

      - name: Check data freshness
        run: |
          python -c "
          from nessi import Client
          client = Client(
              host='${{ secrets.NESSI_API_HOST }}',
              api_key='${{ secrets.NESSI_API_KEY }}'
          )
          statuses = client.freshness.check_all_freshness()
          critical_tables = [s.table_name for s in statuses if s.status == 'critical']
          if critical_tables:
              print(f'Critical freshness issues found in tables: {critical_tables}')
              exit(1)
          "
```

## Best Practices

1. **Setting Appropriate Thresholds**:
   - Warning threshold should be set to detect early signs of staleness
   - Critical threshold should indicate a serious SLA violation

2. **Monitoring Frequency**:
   - Set monitoring frequency based on the expected update frequency
   - For hourly tables, check at least every 15-30 minutes
   - For daily tables, check a few times per day

3. **Trend Analysis**:
   - Regularly review freshness trends to identify patterns
   - Look for recurring issues at specific times or days

4. **Integration with Alerting**:
   - Configure alerts for critical freshness violations
   - Use different notification channels based on severity

5. **Documentation**:
   - Document SLA configurations and their business justification
   - Keep a record of known issues and their resolutions
