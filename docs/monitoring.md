# Nessi Monitoring Guide

This guide explains how to use Nessi's monitoring capabilities with Grafana and Prometheus.

## Overview

Nessi provides comprehensive monitoring through:
- Grafana dashboards for visualization
- Prometheus for metrics collection
- Custom metrics for table statistics

## Architecture

```
┌──────────┐    ┌──────────┐    ┌──────────┐
│  Nessi   │───▶│Prometheus│───▶│ Grafana  │
└──────────┘    └──────────┘    └──────────┘
     │               │               │
     └───────────────┴───────────────┘
         Docker Network (nessi-network)
```

## Metrics Collection

### Available Metrics

1. **Table Metrics**
   - Row count
   - File size
   - Last modified time
   - Schema version

2. **Column Statistics**
   - Count
   - Mean
   - Standard deviation
   - Min/Max values
   - Null count

3. **Performance Metrics**
   - Scan duration
   - Memory usage
   - CPU utilization

### Prometheus Configuration

The `prometheus.yml` configuration:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'nessi'
    static_configs:
      - targets: ['backend:8000']
    metrics_path: '/metrics'

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'grafana'
    static_configs:
      - targets: ['grafana:3000']
```

## Grafana Dashboards

### Default Dashboard

The default dashboard includes:
1. Table Metrics Panel
   - Shows basic table statistics
   - Real-time updates
   - Configurable refresh rate

2. Column Statistics Panel
   - Displays column-level metrics
   - Statistical visualizations
   - Trend analysis

### Creating Custom Dashboards

Using the Python API:

```python
from src.scanner.grafana_dashboard import GrafanaDashboard

# Initialize dashboard
dashboard = GrafanaDashboard(
    grafana_url="http://localhost:3000",
    password="admin"
)

# Create new dashboard
response = dashboard.create_dashboard("Custom Dashboard")
dashboard_uid = response["dashboard"]["uid"]

# Add custom panels
dashboard.add_table_metrics_panel(dashboard_uid, scan_results)
dashboard.add_column_stats_panel(dashboard_uid, column_stats)
```

Using the Grafana UI:
1. Log in to Grafana
2. Click "Create" → "Dashboard"
3. Add panels using the metrics from Prometheus
4. Save and configure refresh settings

### Dashboard Export/Import

Export dashboard:
```python
dashboard.export_dashboard(dashboard_uid, "dashboard.json")
```

Import dashboard:
1. In Grafana UI: Dashboards → Import
2. Upload JSON file or paste content
3. Configure data source mappings

## Alerting

### Setting Up Alerts

1. In Grafana:
   - Open dashboard panel
   - Edit panel
   - Add alert rule
   - Configure conditions and notifications

2. Using the API:
```python
# Alert configuration example (future feature)
dashboard.add_alert_rule(
    dashboard_uid,
    panel_id,
    condition="row_count > 1000000",
    notification_channel="email"
)
```

### Alert Channels

Supported notification channels:
- Email
- Slack
- Webhook
- PagerDuty

## Best Practices

1. **Dashboard Organization**
   - Use meaningful names
   - Group related panels
   - Add descriptions
   - Use appropriate visualizations

2. **Performance**
   - Set appropriate scrape intervals
   - Use efficient queries
   - Monitor resource usage

3. **Security**
   - Change default passwords
   - Use secure connections
   - Implement access control

## Troubleshooting

### Common Issues

1. **Metrics Not Showing**
   - Check Prometheus targets
   - Verify scrape configuration
   - Check network connectivity

2. **Dashboard Problems**
   - Verify data source connection
   - Check panel queries
   - Validate permissions

3. **Alert Issues**
   - Test notification channels
   - Check alert conditions
   - Verify alert history

### Debugging

1. Check Prometheus targets:
```bash
curl http://localhost:9090/api/v1/targets
```

2. Verify metrics collection:
```bash
curl http://localhost:9090/api/v1/query?query=nessi_table_metrics
```

3. Check Grafana logs:
```bash
docker-compose logs grafana
```

## Advanced Configuration

### Custom Metrics

Adding custom metrics:
1. Define metric in Prometheus format
2. Implement collection logic
3. Add to scrape configuration
4. Create visualization panel

### High Availability

For production environments:
1. Use Prometheus federation
2. Set up Grafana clustering
3. Implement backup solutions
4. Monitor the monitoring system

## Next Steps

- Read the [Configuration Guide](configuration.md)
- Explore [API Documentation](../API.md)
- Learn about [Contributing](../CONTRIBUTING.md)
- Check [Security Guidelines](../SECURITY.md) 