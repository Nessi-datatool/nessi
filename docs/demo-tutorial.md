# nessi.dev Demo Tutorial

This tutorial provides a comprehensive walkthrough of nessi.dev's features using the demo script. We'll explain each component and show you how to customize the demo for your needs.

## Overview

The demo showcases three main features:
1. Sample data generation
2. Table scanning and analysis
3. Metrics visualization with Grafana

## Prerequisites

- Completed [Installation Guide](installation.md)
- Running Docker containers (check with `docker-compose ps`)
- Basic Python knowledge
- Basic understanding of Apache Spark

## Step-by-Step Tutorial

### 1. Understanding the Demo Script

The demo script is located at `backend/src/demo.py`. Let's break down its components:

```python
"""
Demo script to showcase nessi.dev features.
Generates sample data, scans tables, and creates a Grafana dashboard.
"""

from pathlib import Path
import time
import os
import requests
from src.scanner.sample_data_generator import SampleDataGenerator
from src.scanner.table_scanner import TableScanner
from src.scanner.grafana_dashboard import GrafanaDashboard
```

### 2. Sample Data Generation

```python
# Create demo directory
demo_dir = Path("data/demo")
demo_dir.mkdir(parents=True, exist_ok=True)

# Generate sample data
print("Generating sample data...")
generator = SampleDataGenerator()
df = generator.generate_sample_data()

# Save in different formats
parquet_path = str(demo_dir / "sample.parquet")
csv_path = str(demo_dir / "sample.csv")

print(f"Saving data to {parquet_path} and {csv_path}...")
generator.save_as_parquet(df, parquet_path)
generator.save_as_csv(df, csv_path)
```

This section:
- Creates a directory for demo data
- Generates sample data with predefined schema
- Saves data in both Parquet and CSV formats

### 3. Table Scanning

```python
print("\nScanning tables...")
scanner = TableScanner()
parquet_results = scanner.scan_table(parquet_path, "parquet")
csv_results = scanner.scan_table(csv_path, "csv")

print("\nParquet table statistics:")
print(f"- Row count: {parquet_results['row_count']}")
print(f"- Columns: {', '.join(field['name'] for field in parquet_results['schema'])}")
print("\nColumn statistics:")
for col, stats in parquet_results['column_stats'].items():
    print(f"- {col}:")
    for stat_name, stat_value in stats.items():
        print(f"  - {stat_name}: {stat_value}")
```

This section:
- Scans both Parquet and CSV files
- Extracts schema information
- Calculates column statistics
- Displays results

### 4. Grafana Dashboard Creation

```python
print("\nSetting up Grafana dashboard...")
grafana_url = os.getenv("GRAFANA_URL", "http://localhost:3000")
grafana_admin_password = os.getenv("GF_SECURITY_ADMIN_PASSWORD", "admin")

# Wait for Grafana to be ready
print("Waiting for Grafana to be ready...")
if not wait_for_grafana(grafana_url):
    print("Warning: Grafana is not responding. Skipping dashboard creation.")
    return

try:
    # Create dashboard
    print("Creating Grafana dashboard...")
    dashboard = GrafanaDashboard(grafana_url, grafana_admin_password)
    dashboard_response = dashboard.create_dashboard("nessi.dev Demo Dashboard")
    dashboard_uid = dashboard_response["dashboard"]["uid"]

    # Add panels
    print("Adding dashboard panels...")
    dashboard.add_table_metrics_panel(dashboard_uid, parquet_results)
    dashboard.add_column_stats_panel(dashboard_uid, parquet_results["column_stats"])

    # Export dashboard
    dashboard_path = str(demo_dir / "dashboard.json")
    print(f"\nExporting dashboard to {dashboard_path}...")
    dashboard.export_dashboard(dashboard_uid, dashboard_path)
```

This section:
- Connects to Grafana
- Creates a new dashboard
- Adds visualization panels
- Exports dashboard configuration

## Customizing the Demo

### 1. Modifying Sample Data

To generate different sample data:

```python
# Change number of rows
df = generator.generate_sample_data(num_rows=1000)

# Add custom columns
from pyspark.sql.functions import lit
df = df.withColumn("custom_column", lit("value"))
```

### 2. Customizing Scans

To scan different file types or locations:

```python
# Scan a specific directory
results = scanner.scan_table("/path/to/your/data.parquet", "parquet")

# Get specific statistics
print(f"Number of rows: {results['row_count']}")
print(f"Schema: {results['schema']}")
```

### 3. Customizing Dashboards

To create custom dashboard panels:

```python
# Create dashboard with custom title
dashboard_response = dashboard.create_dashboard("Custom Dashboard")

# Add specific panels
dashboard.add_table_metrics_panel(dashboard_uid, scan_results)
```

## Running the Demo

1. Start the services:
```bash
docker-compose up -d
```

2. Run the demo script:
```bash
docker-compose exec backend python src/demo.py
```

3. Access the dashboard:
- Open http://localhost:3000 in your browser
- Log in with username: admin, password: admin
- Navigate to "nessi.dev Demo Dashboard"

## Expected Output

The demo will create:
1. Sample data files:
   - `data/demo/sample.parquet`
   - `data/demo/sample.csv`
2. Grafana dashboard with:
   - Table metrics panel
   - Column statistics panel
3. Dashboard export file:
   - `data/demo/dashboard.json`

## Troubleshooting

### Common Issues

1. **Data Generation Fails**
   - Check disk space
   - Verify permissions on data directory

2. **Scanning Errors**
   - Ensure files exist
   - Check file permissions
   - Verify file formats

3. **Grafana Issues**
   - Check if Grafana is running
   - Verify credentials
   - Check network connectivity

## Next Steps

- Explore the [Features Documentation](features.md)
- Learn about [Monitoring](monitoring.md)
- Read the [Configuration Guide](configuration.md)
- Check the [API Documentation](../API.md) 