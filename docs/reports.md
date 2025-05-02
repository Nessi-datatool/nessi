# Nessi Report Generation

## Overview

Nessi provides comprehensive report generation capabilities for data analysis and monitoring. Reports can be generated in multiple formats and include various types of analysis.

## Report Types

### 1. Table Scan Reports

Table scan reports provide detailed information about your data tables, including:

- **File Information**
  - Format (Parquet, CSV, Delta)
  - Path
  - Size
  - Last Modified Date

- **Schema Analysis**
  - Column names
  - Data types
  - Nullability

- **Statistical Analysis**
  - Row count
  - Column-level statistics:
    - Count
    - Null count
    - Min/Max values
    - Mean
    - Standard deviation

Example usage:
```python
from src.scanner.table_scanner import TableScanner

scanner = TableScanner()
result = scanner.scan_table("path/to/table", format="delta")
report = scanner.generate_report(result, output_path="reports/table_analysis.md")
```

### 2. Grafana Dashboard Reports

Nessi can generate interactive Grafana dashboards for real-time monitoring:

- **Table Metrics Panel**
  - Row count trends
  - File size monitoring
  - Schema changes

- **Column Statistics Panel**
  - Distribution analysis
  - Null value tracking
  - Statistical metrics

Example usage:
```python
from src.monitoring.grafana_dashboard import GrafanaDashboard

dashboard = GrafanaDashboard(
    grafana_url="http://localhost:3000",
    api_key="your-api-key"
)

# Create dashboard
dashboard_response = dashboard.create_dashboard("Nessi Table Metrics")

# Add panels
dashboard.add_table_metrics_panel(
    dashboard_uid=dashboard_response['uid'],
    scan_results=scan_results
)

# Export dashboard
dashboard.export_dashboard(
    dashboard_uid=dashboard_response['uid'],
    output_path="dashboards/nessi_metrics.json"
)
```

### 3. Data Quality Reports

Nessi generates data quality reports that include:

- **Completeness Analysis**
  - Missing value detection
  - Null value statistics
  - Required field validation

- **Consistency Checks**
  - Data type validation
  - Value range verification
  - Pattern matching

- **Uniqueness Analysis**
  - Duplicate detection
  - Primary key validation
  - Unique constraint checking

### 4. Performance Reports

Performance reports provide insights into:

- **Processing Metrics**
  - Execution time
  - Memory usage
  - CPU utilization

- **Resource Utilization**
  - Disk I/O
  - Network traffic
  - Cache hit rates

### 5. Custom Reports

Nessi supports custom report generation through:

- **Template System**
  - Jinja2 template support
  - Custom formatting
  - Dynamic content

- **Export Formats**
  - Markdown
  - HTML
  - JSON
  - CSV

## Report Generation Process

1. **Data Collection**
   - Table scanning
   - Statistics calculation
   - Metric aggregation

2. **Analysis**
   - Data profiling
   - Quality assessment
   - Performance evaluation

3. **Report Creation**
   - Template selection
   - Data formatting
   - Output generation

4. **Export**
   - File generation
   - Dashboard creation
   - API integration

## Best Practices

1. **Report Configuration**
   - Set appropriate refresh rates
   - Configure alert thresholds
   - Define report schedules

2. **Performance Optimization**
   - Use efficient queries
   - Implement caching
   - Optimize data collection

3. **Storage Management**
   - Implement retention policies
   - Archive old reports
   - Clean up temporary files

4. **Security**
   - Secure report access
   - Encrypt sensitive data
   - Implement access controls

## Troubleshooting

Common issues and solutions:

1. **Report Generation Fails**
   - Check file permissions
   - Verify data access
   - Review error logs

2. **Performance Issues**
   - Optimize queries
   - Increase resources
   - Implement caching

3. **Formatting Problems**
   - Validate templates
   - Check data types
   - Review export settings

## Next Steps

- Explore the [API Documentation](../API.md)
- Learn about [Monitoring](monitoring.md)
- Read the [Configuration Guide](configuration.md) 