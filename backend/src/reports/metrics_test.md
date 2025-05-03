# Metrics Test Documentation

## Overview

The `metrics_test.py` script generates test metrics for the nessi.dev monitoring system. It simulates various operational scenarios and provides real-time metrics for demonstration and testing purposes.

## Metrics

### Table Scan Metrics
- `nessi_table_scan_total_rows`: Total number of rows scanned
- `nessi_scan_efficiency`: Efficiency of table scans (0-100)
- `nessi_partition_utilization`: Partition utilization score (0-100)
- `nessi_index_usage`: Index usage efficiency (0-100)

### Data Quality Metrics
- `nessi_data_quality_score`: Overall data quality score (0-100)
- `nessi_data_health`: Data health assessment (0-100)
- `nessi_anomaly_count`: Number of anomalies detected

### Performance Metrics
- `nessi_processing_time_seconds`: Processing time in seconds
- `nessi_performance_status`: System performance status (0-100)
- `nessi_resource_utilization`: Resource utilization score (0-100)

### Custom Report Metrics
- `nessi_custom_report_count`: Number of custom reports generated

## Scenarios

The script simulates five different operational scenarios:

### 1. Normal Operation
- Optimal performance metrics
- High data quality scores
- Efficient resource utilization
- Low anomaly count

### 2. High Load
- Increased processing time
- Moderate performance impact
- Higher resource utilization
- Slightly reduced data quality

### 3. Data Quality Issues
- Reduced data quality scores
- Increased anomaly detection
- Moderate performance impact
- Resource utilization issues

### 4. Performance Degradation
- Significant processing delays
- Resource constraints
- System bottlenecks
- Reduced efficiency scores

### 5. Custom Analysis
- Specialized metrics
- Advanced analysis
- Custom reporting
- High efficiency scores

## Usage

1. Start the metrics server:
```bash
python src/reports/metrics_test.py
```

2. Access metrics at:
```
http://localhost:8000/metrics
```

## Configuration

### Metric Ranges

#### Normal Operation
- Table scan rows: 1,000,000 - 2,000,000
- Data quality score: 90-100
- Processing time: 1.0-2.0 seconds
- Scan efficiency: 85-95
- Partition utilization: 80-90
- Index usage: 85-95
- Data health: 90-100
- Performance status: 85-95
- Resource utilization: 80-90

#### High Load
- Table scan rows: 4,000,000 - 5,000,000
- Data quality score: 85-95
- Processing time: 3.0-5.0 seconds
- Scan efficiency: 70-80
- Partition utilization: 60-70
- Index usage: 65-75
- Data health: 80-90
- Performance status: 70-80
- Resource utilization: 60-70

#### Data Quality Issues
- Table scan rows: 1,000,000 - 2,000,000
- Data quality score: 70-80
- Processing time: 2.0-3.0 seconds
- Scan efficiency: 60-70
- Partition utilization: 50-60
- Index usage: 55-65
- Data health: 60-70
- Performance status: 65-75
- Resource utilization: 50-60

#### Performance Degradation
- Table scan rows: 3,000,000 - 4,000,000
- Data quality score: 80-90
- Processing time: 4.0-5.0 seconds
- Scan efficiency: 50-60
- Partition utilization: 40-50
- Index usage: 45-55
- Data health: 70-80
- Performance status: 50-60
- Resource utilization: 40-50

#### Custom Analysis
- Table scan rows: 2,000,000 - 3,000,000
- Data quality score: 95-100
- Processing time: 0.5-1.5 seconds
- Scan efficiency: 90-100
- Partition utilization: 85-95
- Index usage: 90-100
- Data health: 95-100
- Performance status: 90-100
- Resource utilization: 85-95

## Integration

The metrics test script is designed to work with:
- Prometheus for metrics collection
- Grafana for visualization
- Docker Compose for containerization

## Troubleshooting

### Common Issues

1. **Metrics Not Updating**
   - Check if the script is running
   - Verify the update interval
   - Check for errors in the console

2. **Incorrect Metric Values**
   - Verify scenario selection
   - Check random number generation
   - Validate metric ranges

3. **Connection Issues**
   - Check port availability
   - Verify network configuration
   - Check firewall settings

## Development

### Adding New Metrics

1. Define the metric:
```python
new_metric = Gauge('nessi_new_metric', 'Description of new metric')
```

2. Update the metric in scenarios:
```python
if scenario == "normal_operation":
    new_metric.set(random.randint(min_value, max_value))
```

### Modifying Scenarios

1. Add a new scenario:
```python
elif scenario == "new_scenario":
    # Set metrics for new scenario
    pass
```

2. Update scenario selection:
```python
scenarios = ["normal_operation", "high_load", "data_quality_issues", 
             "performance_degradation", "custom_analysis", "new_scenario"]
```

## Best Practices

1. **Metric Naming**
   - Use consistent prefixes
   - Be descriptive
   - Follow Prometheus naming conventions

2. **Value Ranges**
   - Use appropriate ranges for each scenario
   - Ensure realistic values
   - Maintain consistency across scenarios

3. **Update Frequency**
   - Consider system load
   - Balance between real-time updates and performance
   - Use appropriate sleep intervals

## Examples

### Example 1: Normal Operation Metrics
```python
# Metrics during normal operation (0-12 minutes)
{
    'nessi_table_scan_total_rows': 50000,
    'nessi_data_quality_score': 95.0,
    'nessi_processing_time_seconds': 2.5,
    'nessi_scan_efficiency': 0.85,
    'nessi_partition_utilization': 0.65,
    'nessi_index_usage': 0.75,
    'nessi_data_health': 0.90,
    'nessi_performance_status': 0.95,
    'nessi_resource_utilization': 0.60
}
```

### Example 2: High Load Metrics
```python
# Metrics during high load (12-24 minutes)
{
    'nessi_table_scan_total_rows': 150000,
    'nessi_data_quality_score': 85.0,
    'nessi_processing_time_seconds': 4.5,
    'nessi_scan_efficiency': 0.65,
    'nessi_partition_utilization': 0.85,
    'nessi_index_usage': 0.55,
    'nessi_data_health': 0.80,
    'nessi_performance_status': 0.75,
    'nessi_resource_utilization': 0.85
}
```

### Example 3: Data Quality Issues
```python
# Metrics during data quality issues (24-36 minutes)
{
    'nessi_table_scan_total_rows': 75000,
    'nessi_data_quality_score': 65.0,
    'nessi_processing_time_seconds': 3.5,
    'nessi_scan_efficiency': 0.55,
    'nessi_partition_utilization': 0.75,
    'nessi_index_usage': 0.65,
    'nessi_data_health': 0.60,
    'nessi_performance_status': 0.70,
    'nessi_resource_utilization': 0.75,
    'nessi_anomaly_count': 5
}
```

### Example 4: Custom Analysis
```python
# Metrics during custom analysis (48-60 minutes)
{
    'nessi_table_scan_total_rows': 100000,
    'nessi_data_quality_score': 92.0,
    'nessi_processing_time_seconds': 3.0,
    'nessi_scan_efficiency': 0.80,
    'nessi_partition_utilization': 0.70,
    'nessi_index_usage': 0.80,
    'nessi_data_health': 0.88,
    'nessi_performance_status': 0.90,
    'nessi_resource_utilization': 0.70,
    'nessi_custom_report_count': 8
}
```

### Example 5: Adding a New Metric
```python
# Define a new custom metric
custom_metric = Gauge('nessi_custom_metric', 'Custom metric description')

# Set metric value based on scenario
if scenario == 'normal_operation':
    custom_metric.set(0.9)
elif scenario == 'high_load':
    custom_metric.set(0.7)
elif scenario == 'data_quality_issues':
    custom_metric.set(0.5)
elif scenario == 'performance_degradation':
    custom_metric.set(0.6)
elif scenario == 'custom_analysis':
    custom_metric.set(0.8)
```

### Example 6: Custom Alert Rule
```yaml
groups:
  - name: nessi-custom-alerts
    rules:
      - alert: LowCustomMetric
        expr: nessi_custom_metric < 0.7
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Custom metric below threshold"
          description: "Custom metric has been below 0.7 for 5 minutes"
```

### Example 7: Grafana Panel Configuration
```json
{
  "title": "Custom Metric",
  "type": "gauge",
  "datasource": "Prometheus",
  "targets": [
    {
      "expr": "nessi_custom_metric",
      "legendFormat": "Custom"
    }
  ],
  "fieldConfig": {
    "defaults": {
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {"color": "red", "value": null},
          {"color": "yellow", "value": 0.7},
          {"color": "green", "value": 0.9}
        ]
      },
      "unit": "percentunit"
    }
  }
}
```

### Example 8: Docker Compose Service
```yaml
services:
  metrics-test:
    build:
      context: .
      dockerfile: Dockerfile.test
    ports:
      - "8000:8000"
    networks:
      - monitoring
    environment:
      - PYTHONUNBUFFERED=1
      - SCENARIO=normal_operation
      - METRICS_INTERVAL=5s
    volumes:
      - ./src/reports:/app/src/reports
    command: python src/reports/metrics_test.py
```

### Example 9: Prometheus Scrape Config
```yaml
scrape_configs:
  - job_name: 'nessi-metrics'
    scrape_interval: 5s
    static_configs:
      - targets: ['metrics-test:8000']
    metrics_path: '/metrics'
    metric_relabel_configs:
      - source_labels: [__name__]
        regex: 'nessi_(.*)'
        target_label: 'metric_type'
        replacement: '$1'
```

### Example 10: Metric Relabeling
```yaml
metric_relabel_configs:
  - source_labels: [__name__]
    regex: 'nessi_(.*)'
    target_label: 'metric_type'
    replacement: '$1'
  - source_labels: [instance]
    regex: '(.*):\d+'
    target_label: 'instance'
    replacement: '$1'
``` 