# Nessi Features Documentation

## Overview
Nessi is a powerful data processing and analysis tool built on Apache Spark, designed for efficient table scanning, data generation, and analysis. It provides a robust set of features for working with various data formats and performing comprehensive data analysis.

## Core Features

### 1. Table Scanning
The `TableScanner` class provides comprehensive table analysis capabilities:

#### Supported Formats
- **Parquet**: Native support for Parquet file format
- **CSV**: Full support for CSV files with header and schema inference

#### Scanning Capabilities
- **Schema Analysis**
  - Automatic schema detection
  - Detailed type information for each column
  - Support for complex data types

- **Statistical Analysis**
  - Row count calculation
  - Column-level statistics including:
    - Count
    - Mean
    - Standard deviation
    - Minimum value
    - Maximum value

- **Special Character Handling**
  - Automatic handling of special characters in column names
  - Conversion of '@' to '_at_'
  - Conversion of '.' to '_dot_'

### 2. Sample Data Generation
The `SampleDataGenerator` class provides flexible data generation capabilities:

#### Data Generation Features
- **Structured Data Creation**
  - Support for multiple data types:
    - Integer
    - String
    - Timestamp
  - Customizable schema definition
  - Automatic timestamp generation

- **Data Export Options**
  - **Parquet Export**
    - Efficient columnar storage
    - Automatic directory creation
    - Overwrite mode support
  
  - **CSV Export**
    - Header inclusion
    - Proper quote handling
    - Escape character support
    - Single file output (coalesced)

### 3. Spark Configuration
The `spark_config` module provides optimized Spark session management:

#### Configuration Features
- **Resource Management**
  - Memory configuration for driver and executor
  - Local execution mode support
  - Dynamic allocation control

- **Performance Optimization**
  - Shuffle partition configuration
  - Parallelism settings
  - UI control for reduced overhead

- **Session Management**
  - Automatic session cleanup
  - Custom application naming
  - Local host binding

### 4. Grafana Dashboard Integration
The `GrafanaDashboard` class provides comprehensive monitoring capabilities:

#### Dashboard Management
- **Dashboard Creation**
  - Customizable dashboard titles
  - Folder organization
  - Automatic refresh configuration
  - Tag-based organization

- **Metrics Visualization**
  - **Table Metrics Panel**
    - Row count visualization
    - File information display
    - Real-time updates
  
  - **Column Statistics Panel**
    - Statistical metrics visualization
    - Threshold-based coloring
    - Multi-column support
    - Time-series data display

- **Dashboard Export**
  - JSON format export
  - Version control integration
  - Backup capabilities
  - Template sharing

#### Integration Features
- **Prometheus Integration**
  - Native Prometheus data source support
  - Custom metric queries
  - Alert configuration
  - Threshold monitoring

- **API Management**
  - Secure API key authentication
  - RESTful API integration
  - Dashboard programmatic management
  - Panel customization

## Technical Specifications

### System Requirements
- Python 3.11 or higher
- Java 17
- Apache Spark 3.5.0
- Delta Lake 3.0.0
- Grafana 8.0 or higher
- Prometheus 2.0 or higher

### Dependencies
- numpy<2.0.0
- pandas==2.1.4
- pyarrow==14.0.1
- pyspark==3.5.0
- delta-spark==3.0.0
- requests>=2.31.0
- Additional utility packages for testing and development

### Performance Characteristics
- Optimized for local development
- Memory-efficient configuration
- Automatic resource cleanup
- Efficient data processing
- Real-time monitoring capabilities

## Usage Examples

### Table Scanning
```python
from backend.src.scanner.table_scanner import TableScanner

# Initialize scanner
scanner = TableScanner()

# Scan a Parquet table
result = scanner.scan_table("path/to/table.parquet", "parquet")

# Access results
print(f"Row count: {result['row_count']}")
print(f"Schema: {result['schema']}")
print(f"Statistics: {result['column_stats']}")
```

### Sample Data Generation
```python
from backend.src.scanner.sample_data_generator import SampleDataGenerator

# Initialize generator
generator = SampleDataGenerator()

# Generate sample data
df = generator.generate_sample_data()

# Save as Parquet
generator.save_as_parquet(df, "output/path")

# Save as CSV
generator.save_as_csv(df, "output/path")
```

### Grafana Dashboard Integration
```python
from backend.src.scanner.grafana_dashboard import GrafanaDashboard

# Initialize dashboard manager
dashboard = GrafanaDashboard(
    grafana_url="http://localhost:3000",
    api_key="your-api-key"
)

# Create a new dashboard
dashboard_response = dashboard.create_dashboard("Nessi Table Metrics")

# Add table metrics panel
dashboard.add_table_metrics_panel(
    dashboard_uid=dashboard_response['uid'],
    scan_results=scan_results
)

# Add column statistics panel
dashboard.add_column_stats_panel(
    dashboard_uid=dashboard_response['uid'],
    column_stats=column_stats
)

# Export dashboard
dashboard.export_dashboard(
    dashboard_uid=dashboard_response['uid'],
    output_path="dashboards/nessi_metrics.json"
)
```

## Best Practices

### Performance Optimization
1. Use Parquet format for better performance
2. Monitor memory usage with large datasets
3. Clean up Spark sessions when done
4. Configure appropriate refresh rates for dashboards
5. Use efficient Prometheus queries

### Error Handling
1. Always check for file existence before scanning
2. Handle special characters in column names
3. Validate format types before processing
4. Implement proper API error handling
5. Monitor dashboard creation status

### Data Management
1. Use appropriate data types for columns
2. Consider data size when choosing format
3. Implement proper cleanup procedures
4. Regular dashboard backups
5. Version control for dashboard configurations

## Future Enhancements
- Delta Lake table support
- Additional file format support
- Enhanced statistical analysis
- Distributed processing capabilities
- Advanced data validation features
- Custom dashboard templates
- Advanced alerting rules
- Multi-datasource support
- Dashboard sharing capabilities 