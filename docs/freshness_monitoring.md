# Data Freshness Monitoring

This document describes the freshness monitoring feature of Nessi.dev.

## Overview

Freshness monitoring allows users to define Service Level Agreements (SLAs) for data freshness and monitor compliance with these SLAs over time. This feature is essential for ensuring that data is up-to-date and reliable for downstream consumers.

The freshness monitoring system integrates with Delta Lake to track the last update time of tables and compare it against defined SLAs. It provides visibility into data freshness through CLI commands and a Python API, with comprehensive reporting capabilities.

## Overview

The freshness monitoring feature helps you:

- Define expected update frequencies for tables
- Monitor adherence to these SLAs
- Get alerts when tables are not updated as expected
- Generate HTML and PDF reports for freshness status
- Export freshness data in JSON/CSV formats for compliance and auditing

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
- `report_url`: Optional URL to a generated report for the table

Example:
```bash
nessi sla --table customer_orders --define "frequency=1h,warning=150,critical=200,description='Customer orders table with hourly updates',tags=sales,core,production,grafana=https://grafana.example.com/d/abc123/customer-orders"
```

## Reporting

Nessi provides comprehensive reporting capabilities for monitoring the freshness of your tables. Reports include:

- Current freshness status for all tables
- SLA configuration details
- Historical trends and compliance statistics
- Detailed metrics for analysis

Generate reports using the CLI:
```bash
# Generate a freshness report for all tables
nessi freshness report --format html --output freshness-report.html

# Generate a freshness report for a specific table
nessi freshness report --table customer_orders --format pdf --output customer-orders-freshness.pdf

# Export freshness data as JSON for integration
nessi freshness export --format json --output freshness-data.json
```

### Report Contents

Freshness reports include the following information:

1. **Status Summary**: Shows the current freshness status of all tables, including:
   - Last update time
   - Time since last update
   - Expected frequency
   - Next expected update
   - Status (Up-to-date, Warning, Critical)

2. **SLA Configurations**: Allows you to view, add, edit, and delete SLA configurations.

3. **Trends**: Shows historical freshness data and compliance statistics.

## CLI Commands

The freshness feature provides the following CLI commands:

### Status Commands

```bash
# Get freshness status for all tables
nessi freshness status

# Get freshness status for a specific table
nessi freshness status --table <table_name>

# Get detailed freshness status with metrics
nessi freshness status --detailed
```

These commands return the current freshness status of all tables or a specific table.

### SLA Management Commands

```bash
# List all SLA configurations
nessi freshness sla list

# Get SLA configuration for a specific table
nessi freshness sla get --table <table_name>

# Create or update an SLA configuration
nessi freshness sla set --table <table_name> --frequency <frequency> --warning <threshold> --critical <threshold>

# Delete an SLA configuration
nessi freshness sla delete --table <table_name>
```

These commands allow you to manage SLA configurations directly from the command line.

### Trends Commands

```bash
# Get historical freshness trends for all tables
nessi freshness trends

# Get historical freshness trends for a specific table
nessi freshness trends --table <table_name>

# Get trends for a specific time period
nessi freshness trends --from "2025-05-01" --to "2025-05-20"
```

These commands provide historical freshness data and compliance statistics for analysis and reporting.

### Export Commands

```bash
# Export freshness data as JSON for integration
nessi freshness export --format json --output freshness-data.json

# Export freshness data as CSV for compliance and auditing
nessi freshness export --format csv --output freshness-data.csv
```

These commands allow you to export freshness data in JSON or CSV format.

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

- **Define Realistic SLAs**: Set frequency expectations that match your actual data pipeline schedules
- **Start with Warning Thresholds**: Begin with higher warning thresholds and adjust as you understand your data patterns
- **Tag Tables by Importance**: Use tags to categorize tables by importance and criticality
- **Monitor Compliance Trends**: Use the reporting system to track compliance over time
- **Generate Regular Reports**: Schedule regular report generation for documentation and auditing.

## Troubleshooting

### Common Issues

1. **SLA Not Being Monitored**: Verify that the table path is correct and accessible.
2. **Incorrect Freshness Status**: Check if the Delta Lake transaction log is accessible and up-to-date.
3. **Reports Not Showing Data**: Verify that the freshness manager is properly initialized and SLAs are defined.

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

### Example 3: Weekly Table with Report URL

```bash
# Define an SLA for a table that should be updated weekly, with a report URL
nessi sla --table weekly_reports --define "frequency=1w,warning=120,critical=150,report_url=https://reports.example.com/weekly-reports.html"

# Check the freshness status
nessi freshness --table weekly_reports
```
