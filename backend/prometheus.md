# Prometheus Configuration Documentation

## Overview

This document describes the Prometheus configuration for the nessi.dev monitoring system. Prometheus is used to collect, store, and query time series data from various metrics endpoints.

## Configuration File

The main configuration file is located at `prometheus.yml` and contains the following sections:

### Global Configuration
```yaml
global:
  scrape_interval: 5s
  evaluation_interval: 5s
```

- **scrape_interval**: How frequently to scrape targets
- **evaluation_interval**: How frequently to evaluate rules

### Scrape Configurations

#### nessi.dev Metrics Job
```yaml
scrape_configs:
  - job_name: 'nessi-metrics'
    static_configs:
      - targets: ['metrics-test:8000']
    metrics_path: '/metrics'
    scheme: 'http'
```

- **job_name**: Identifier for the metrics collection job
- **targets**: List of endpoints to scrape
- **metrics_path**: Path to metrics endpoint
- **scheme**: Protocol to use (http/https)

### Relabeling Configuration

#### Instance Relabeling
```yaml
relabel_configs:
  - source_labels: ['__address__']
    target_label: 'instance'
    regex: '(.*):.*'
    replacement: '$1'
```

- Extracts instance name from address
- Used for better metric organization

#### Metric Relabeling
```yaml
metric_relabel_configs:
  - source_labels: ['__name__']
    regex: 'nessi_(.*)'
    target_label: 'metric_type'
    replacement: '$1'
```

- Categorizes metrics by type
- Improves metric organization

## Metrics Collection

### Table Scan Metrics
- `nessi_table_scan_total_rows`
- `nessi_scan_efficiency`
- `nessi_partition_utilization`
- `nessi_index_usage`

### Data Quality Metrics
- `nessi_data_quality_score`
- `nessi_data_health`
- `nessi_anomaly_count`

### Performance Metrics
- `nessi_processing_time_seconds`
- `nessi_performance_status`
- `nessi_resource_utilization`

### Custom Report Metrics
- `nessi_custom_report_count`

## Alerting Rules

### Table Scan Alerts
```yaml
groups:
  - name: table_scan_alerts
    rules:
      - alert: HighRowCount
        expr: nessi_table_scan_total_rows > 4000000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: High row count detected
          description: Table scan shows more than 4M rows
```

### Data Quality Alerts
```yaml
      - alert: LowDataQuality
        expr: nessi_data_quality_score < 70
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: Low data quality detected
          description: Data quality score below 70%
```

### Performance Alerts
```yaml
      - alert: HighProcessingTime
        expr: nessi_processing_time_seconds > 4
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: High processing time
          description: Processing time exceeds 4 seconds
```

## Integration

### Grafana Integration
- Data source configuration
- Dashboard queries
- Alert visualization

### Metrics Test Integration
- Endpoint configuration
- Scrape interval
- Metric validation

## Best Practices

### Configuration
1. Use appropriate scrape intervals
2. Implement proper relabeling
3. Configure alert rules
4. Set up proper retention

### Monitoring
1. Monitor scrape success
2. Track metric cardinality
3. Watch for failed scrapes
4. Monitor alert states

### Performance
1. Optimize scrape intervals
2. Manage metric cardinality
3. Configure proper retention
4. Monitor resource usage

## Troubleshooting

### Common Issues

1. **Scrape Failures**
   - Check endpoint availability
   - Verify network connectivity
   - Check firewall settings
   - Validate credentials

2. **High Cardinality**
   - Review metric labels
   - Optimize relabeling
   - Check for label explosions
   - Monitor storage usage

3. **Alert Issues**
   - Verify rule syntax
   - Check evaluation interval
   - Validate alert conditions
   - Test notification channels

### Logs
- Prometheus logs: `docker-compose logs prometheus`
- Metrics logs: `docker-compose logs metrics-test`

## Maintenance

### Regular Tasks
1. Monitor scrape success
2. Review alert rules
3. Check storage usage
4. Validate configurations

### Backup
1. Export configurations
2. Backup alert rules
3. Document changes
4. Version control files

## Security

### Access Control
1. Configure authentication
2. Set up authorization
3. Manage API access
4. Secure endpoints

### Network Security
1. Use HTTPS where possible
2. Configure firewalls
3. Monitor access
4. Log security events 