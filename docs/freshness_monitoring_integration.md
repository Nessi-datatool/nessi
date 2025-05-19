# Freshness Monitoring Integration Guide

This guide explains how to integrate the Freshness Monitoring system with other components of Nessi.dev.

## Integration with Delta Lake

The Freshness Monitoring system integrates with Delta Lake to track table freshness:

```go
// Example of integrating with Delta Lake connector
func (m *SLAManager) CheckFreshness(tableName string) (*FreshnessStatus, error) {
    // Get SLA configuration
    config, ok := m.GetSLA(tableName)
    if !ok {
        return nil, fmt.Errorf("no SLA configuration found for table %s", tableName)
    }
    
    // Get table metadata from Delta Lake
    metadata, err := m.deltaConnector.GetTableMetadata(config.TablePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read table metadata: %w", err)
    }
    
    // Check freshness status
    status := m.calculateFreshnessStatus(config, metadata.LastModified)
    
    // Update history
    m.updateHistory(status)
    
    return status, nil
}
```

### Cloud Integration

For cloud-hosted Delta Lake tables, use the CloudDeltaConnector:

```go
// Initialize with appropriate cloud provider
cloudConnector, err := NewCloudDeltaConnector(CloudProviderAWS, awsConfig)
if err != nil {
    log.Fatalf("Failed to create cloud connector: %v", err)
}

// Create SLA manager with cloud connector
manager, err := NewSLAManager(configPath, cloudConnector)
if err != nil {
    log.Fatalf("Failed to create SLA manager: %v", err)
}
```

## Integration with CLI

### CLI Commands

The following CLI commands are available for freshness monitoring:

- `nessi freshness list-sla` - List all SLA configurations
- `nessi freshness create-sla` - Create a new SLA configuration
- `nessi freshness get-sla <tableName>` - Get SLA configuration for a table
- `nessi freshness update-sla <tableName>` - Update SLA configuration for a table
- `nessi freshness delete-sla <tableName>` - Delete SLA configuration for a table
- `nessi freshness status` - Get freshness status for all tables
- `nessi freshness status <tableName>` - Get freshness status for a table
- `nessi freshness trends` - Get freshness trends for all tables
- `nessi freshness trends <tableName>` - Get freshness trends for a table

### Output Formats

The CLI commands support various output formats:

- JSON output for programmatic consumption
- Table format for terminal viewing
- CSV export for further analysis

## Integration with Security System

### Authentication

All freshness monitoring CLI commands are protected by API key authentication:

```go
// Example of securing freshness CLI commands
func executeFreshnessCommand(cmd *cobra.Command, args []string) error {
    // Verify API key from environment or config file
    apiKey := viper.GetString("api_key")
    if apiKey == "" {
        return fmt.Errorf("API key is required for freshness monitoring commands")
    }
    
    // Authenticate with API key
    user, err := authManager.ValidateAPIKey(apiKey)
    if err != nil {
        return fmt.Errorf("authentication failed: %w", err)
    }
    
    // Continue with command execution
    // ...
}
```

### Role-Based Access Control

Different operations require different roles:

- Viewing freshness status: `viewer` role
- Managing SLA configurations: `admin` role

```go
// Example of RBAC for SLA management in CLI commands
func executeSLACreateCommand(cmd *cobra.Command, args []string) error {
    // Get user from API key authentication
    user, err := authManager.ValidateAPIKey(apiKey)
    if err != nil {
        return fmt.Errorf("authentication failed: %w", err)
    }
    
    // Check if user has admin role
    if !authManager.CheckUserRole(user, "admin") {
        return fmt.Errorf("insufficient permissions: admin role required")
    }
    
    // Process command
    // ...
    return nil
}
```

### Audit Logging

All SLA management operations are logged for audit purposes:

```go
// Example of audit logging for SLA management in CLI commands
func executeSLACreateCommand(cmd *cobra.Command, args []string) error {
    // Get user from API key authentication
    user, err := authManager.ValidateAPIKey(apiKey)
    if err != nil {
        return fmt.Errorf("authentication failed: %w", err)
    }
    
    // Log the operation
    auditLogger.LogUserAction(user.Username, "create_sla", map[string]interface{}{
        "table_name": slaConfig.TableName,
        "table_path": slaConfig.TablePath,
    })
    
    // Process command
    // ...
    return nil
}
```

## Integration with Plugin System

### ValidationPlugin

Custom validation logic for freshness checks:

```go
// Example of a ValidationPlugin for freshness checks
type FreshnessValidationPlugin struct {
    // Plugin implementation
}

func (p *FreshnessValidationPlugin) Validate(data interface{}) (bool, error) {
    // Custom validation logic
    freshnessStatus, ok := data.(*FreshnessStatus)
    if !ok {
        return false, fmt.Errorf("invalid data type")
    }
    
    // Custom validation logic
    // ...
    
    return isValid, nil
}
```

### AlertPlugin

Custom alerting for freshness violations:

```go
// Example of an AlertPlugin for freshness violations
type FreshnessAlertPlugin struct {
    // Plugin implementation
}

func (p *FreshnessAlertPlugin) SendAlert(data interface{}) error {
    // Custom alerting logic
    freshnessStatus, ok := data.(*FreshnessStatus)
    if !ok {
        return fmt.Errorf("invalid data type")
    }
    
    // Send alert
    // ...
    
    return nil
}
```

### StoragePlugin

Custom storage backends for SLA configurations and history:

```go
// Example of a StoragePlugin for SLA configurations
type SLAStoragePlugin struct {
    // Plugin implementation
}

func (p *SLAStoragePlugin) Store(key string, data interface{}) error {
    // Custom storage logic
    slaConfig, ok := data.(*SLAConfig)
    if !ok {
        return fmt.Errorf("invalid data type")
    }
    
    // Store data
    // ...
    
    return nil
}
```

## Integration with Python API

### Python Client

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

### Format Handler Integration

For multi-format support:

```python
from nessi import Client, FormatHandler

# Initialize client
client = Client(host="localhost:8080", api_key="your-api-key")

# Create format handler
format_handler = FormatHandler()

# Detect format
format_info = format_handler.detect_format("/path/to/my_table")

# Add SLA with format-specific configuration
client.freshness.add_sla(
    table_name="my_table",
    table_path="/path/to/my_table",
    expected_frequency="1d",
    warning_threshold=150,
    critical_threshold=200,
    metadata={"format": format_info.format_type}
)
```

## Integration with Workflow Orchestration

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

## Integration with Intelligent Alerting

```go
// Example of integrating with intelligent alerting
func (m *SLAManager) CheckAllFreshness() ([]*FreshnessStatus, error) {
    // Check freshness for all tables
    statuses, err := m.checkAllFreshness()
    if err != nil {
        return nil, err
    }
    
    // Send to intelligent alerting system
    for _, status := range statuses {
        if status.Status == SLALevelCritical {
            m.intelligentAlertManager.SendAlert(AlertTypeFreshness, status)
        }
    }
    
    return statuses, nil
}
```

## Integration with Cloud Providers

### AWS Integration

```go
// Example of AWS S3 integration
awsConfig := &CloudProviderConfig{
    Region:          "us-west-2",
    Bucket:          "my-delta-lake-bucket",
    AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
    SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
}

// Create AWS Delta connector
awsConnector, err := NewCloudDeltaConnector(CloudProviderAWS, awsConfig)
if err != nil {
    log.Fatalf("Failed to create AWS connector: %v", err)
}

// Create SLA manager with AWS connector
manager, err := NewSLAManager(configPath, awsConnector)
if err != nil {
    log.Fatalf("Failed to create SLA manager: %v", err)
}
```

### Azure Integration

```go
// Example of Azure Blob Storage integration
azureConfig := &CloudProviderConfig{
    AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT"),
    AccountKey:    os.Getenv("AZURE_STORAGE_KEY"),
    ContainerName: "my-delta-lake-container",
}

// Create Azure Delta connector
azureConnector, err := NewCloudDeltaConnector(CloudProviderAzure, azureConfig)
if err != nil {
    log.Fatalf("Failed to create Azure connector: %v", err)
}

// Create SLA manager with Azure connector
manager, err := NewSLAManager(configPath, azureConnector)
if err != nil {
    log.Fatalf("Failed to create SLA manager: %v", err)
}
```

### GCP Integration

```go
// Example of GCP Cloud Storage integration
gcpConfig := &CloudProviderConfig{
    ProjectID:      os.Getenv("GCP_PROJECT_ID"),
    CredentialsFile: os.Getenv("GCP_CREDENTIALS_FILE"),
    Bucket:         "my-delta-lake-bucket",
}

// Create GCP Delta connector
gcpConnector, err := NewCloudDeltaConnector(CloudProviderGCP, gcpConfig)
if err != nil {
    log.Fatalf("Failed to create GCP connector: %v", err)
}

// Create SLA manager with GCP connector
manager, err := NewSLAManager(configPath, gcpConnector)
if err != nil {
    log.Fatalf("Failed to create SLA manager: %v", err)
}
```
