# Nessi Quick Start Guide

This guide will help you get started with Nessi quickly. We'll cover the basic operations and show you how to perform common tasks.

## Prerequisites

Ensure you have completed the [Installation Guide](installation.md) and have Nessi running in Docker.

## Basic Operations

### 1. Generate Sample Data

```python
from src.scanner.sample_data_generator import SampleDataGenerator

# Initialize generator
generator = SampleDataGenerator()

# Generate sample data
df = generator.generate_sample_data()

# Save in different formats
generator.save_as_parquet(df, "data/sample.parquet")
generator.save_as_csv(df, "data/sample.csv")
```

### 2. Scan Tables

```python
from src.scanner.table_scanner import TableScanner

# Initialize scanner
scanner = TableScanner()

# Scan Parquet file
parquet_results = scanner.scan_table("data/sample.parquet", "parquet")

# Print results
print(f"Row count: {parquet_results['row_count']}")
print("Columns:", ", ".join(field['name'] for field in parquet_results['schema']))

# Print column statistics
for col, stats in parquet_results['column_stats'].items():
    print(f"\n{col} statistics:")
    for stat_name, value in stats.items():
        print(f"  {stat_name}: {value}")
```

### 3. Monitor with Grafana

```python
from src.scanner.grafana_dashboard import GrafanaDashboard

# Initialize dashboard manager
dashboard = GrafanaDashboard("http://localhost:3000", "admin")

# Create dashboard
response = dashboard.create_dashboard("My First Dashboard")
dashboard_uid = response["dashboard"]["uid"]

# Add panels
dashboard.add_table_metrics_panel(dashboard_uid, parquet_results)
dashboard.add_column_stats_panel(dashboard_uid, parquet_results["column_stats"])
```

## Running the Demo

The easiest way to see everything in action is to run the demo script:

```bash
docker-compose exec backend python src/demo.py
```

This will:
1. Generate sample data
2. Save it in both Parquet and CSV formats
3. Scan the tables
4. Create a Grafana dashboard with metrics

## Common Tasks

### Working with Different File Formats

Nessi supports multiple file formats:

```python
# CSV files
csv_results = scanner.scan_table("data/sample.csv", "csv")

# Parquet files
parquet_results = scanner.scan_table("data/sample.parquet", "parquet")
```

### Customizing Sample Data

You can customize the sample data generation:

```python
# Generate data with specific size
df = generator.generate_sample_data(num_rows=1000)

# Add custom columns
df = df.withColumn("custom_col", ...)
```

### Monitoring and Metrics

Access your metrics:
1. Open Grafana at http://localhost:3000
2. Log in with username: admin, password: admin
3. Navigate to your dashboard
4. View real-time metrics and statistics

## Best Practices

1. **Resource Management**
   - Close Spark sessions when done
   - Monitor memory usage
   - Clean up temporary files

2. **Error Handling**
   - Always check return values
   - Handle exceptions appropriately
   - Validate input data

3. **Performance**
   - Use appropriate file formats (Parquet for large datasets)
   - Monitor Spark execution
   - Configure memory settings appropriately

## Next Steps

- Read the [Features Documentation](features.md) for detailed feature information
- Check the [API Documentation](../API.md) for complete API reference
- Follow the [Demo Tutorial](demo-tutorial.md) for a comprehensive walkthrough
- Learn about [Monitoring](monitoring.md) for advanced Grafana usage 