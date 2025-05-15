# Data Freshness Monitoring

This document describes the freshness monitoring feature of Nessi.dev.

## Overview

Freshness monitoring allows users to define Service Level Agreements (SLAs) for data freshness and monitor compliance with these SLAs over time. This feature is essential for ensuring that data is up-to-date and reliable for downstream consumers.

The freshness monitoring system integrates with Delta Lake to track the last update time of tables and compare it against defined SLAs. It provides visibility into data freshness through a dashboard, CLI commands, and a Python API.

## Overview

The freshness monitoring feature helps you:

- Define expected update frequencies for tables
- Monitor adherence to these SLAs
- Get alerts when tables are not updated as expected
- Visualize freshness status in a dashboard
- Export freshness reports for compliance and auditing

## Configuration

### CLI Configuration

You can configure SLAs for tables using the Nessi CLI:

```bash
# Define an SLA for a table
nessi sla --table customer_orders --define "frequency=1h,warning=150,critical=200"

# Check the current freshness status of a table
nessi freshness --table customer_orders

# List all SLA configurations
nessi sla --list

# Delete an SLA configuration
nessi sla --table customer_orders --delete
```

### SLA Configuration Parameters

When defining an SLA, you can specify the following parameters:

- `frequency`: Expected update frequency (e.g., 15m, 1h, 1d, 1w, 1mo)
- `warning`: Warning threshold as a percentage of the expected frequency (e.g., 150 for 1.5x)
- `critical`: Critical threshold as a percentage of the expected frequency (e.g., 200 for 2x)
- `description`: Optional description of the SLA
- `tags`: Optional comma-separated list of tags for categorization
- `grafana`: Optional URL to a Grafana dashboard for the table

Example:
```bash
nessi sla --table customer_orders --define "frequency=1h,warning=150,critical=200,description='Customer orders table with hourly updates',tags=sales,core,production,grafana=https://grafana.example.com/d/abc123/customer-orders"
```

## Dashboard

The freshness dashboard provides a visual interface for monitoring the freshness of your tables. It includes:

- Current freshness status for all tables
- SLA configuration management
- Historical trends and compliance statistics
- Export functionality for reporting

To access the dashboard, navigate to:
```
http://localhost:8080/freshness
```

### Dashboard Tabs

The dashboard has three main tabs:

1. **Status**: Shows the current freshness status of all tables, including:
   - Last update time
   - Time since last update
   - Expected frequency
   - Next expected update
   - Status (Up-to-date, Warning, Critical)

2. **SLA Configurations**: Allows you to view, add, edit, and delete SLA configurations.

3. **Trends**: Shows historical freshness data and compliance statistics.

## API

The freshness feature provides the following API endpoints:

### Status API

```
GET /api/freshness/status
GET /api/freshness/status?table=<table_name>
```

Returns the current freshness status of all tables or a specific table.

### SLA API

```
GET /api/freshness/sla
GET /api/freshness/sla?table=<table_name>
POST /api/freshness/sla
DELETE /api/freshness/sla?table=<table_name>
```

Manages SLA configurations.

### Trends API

```
GET /api/freshness/trends
GET /api/freshness/trends?table=<table_name>
```

Returns historical freshness data and compliance statistics.

### Export API

```
GET /api/freshness/export?format=csv
GET /api/freshness/export?format=json
```

Exports freshness data in CSV or JSON format.

## Integration with Monitoring

The freshness feature integrates with the Nessi.dev monitoring system to:

- Record metrics about table freshness
- Trigger alerts when SLAs are violated
- Provide links to Grafana dashboards for detailed monitoring

## Configuration File

SLA configurations are stored in a JSON file at:

```
$NESSI_CONFIG_DIR/freshness/sla_config.json
```

This file is automatically created and managed by the Nessi.dev system, but can also be edited manually if needed.

## Best Practices

1. **Set Realistic SLAs**: Define SLAs based on actual business requirements and technical capabilities.

2. **Use Tags**: Tag your tables for easier filtering and categorization.

3. **Link to Grafana**: Provide Grafana dashboard URLs for detailed monitoring.

4. **Regular Review**: Periodically review and adjust SLAs based on changing business needs.

5. **Monitor Trends**: Use the trends tab to identify patterns and potential issues.

## Troubleshooting

### Common Issues

1. **Table Not Found**: Ensure the table path is correct and the table exists in your Delta Lake.

2. **Incorrect Freshness Status**: Check if the table's transaction log is accessible and contains valid timestamps.

3. **Dashboard Not Showing Data**: Verify that the freshness manager is properly initialized and SLAs are defined.

### Logs

Freshness-related logs can be found in the Nessi.dev log files with the component name `freshness-manager`.

## Examples

### Example 1: Hourly Table

```bash
# Define an SLA for a table that should be updated hourly
nessi sla --table hourly_metrics --define "frequency=1h,warning=150,critical=200"

# Check the freshness status
nessi freshness --table hourly_metrics
```

### Example 2: Daily Table with Tags

```bash
# Define an SLA for a table that should be updated daily, with tags
nessi sla --table daily_sales --define "frequency=1d,warning=125,critical=150,tags=sales,finance,daily"

# Check the freshness status
nessi freshness --table daily_sales
```

### Example 3: Weekly Table with Grafana Link

```bash
# Define an SLA for a table that should be updated weekly, with a Grafana link
nessi sla --table weekly_reports --define "frequency=1w,warning=120,critical=150,grafana=https://grafana.example.com/d/abc123/weekly-reports"

# Check the freshness status
nessi freshness --table weekly_reports
```
