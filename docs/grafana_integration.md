# Grafana and Prometheus Dashboarding in Nessi

Nessi provides built-in dashboard functionality for core monitoring capabilities, with additional support for Grafana and Prometheus integration for more advanced visualization needs.

## Overview

Nessi includes 5 pre-built Grafana dashboard templates that provide comprehensive monitoring capabilities:

1. **Nessi Dashboard** - General overview of system metrics
2. **Nessi Alerts Dashboard** - Detailed view of all alerts and their statuses
3. **Nessi Data Quality Dashboard** - Data quality metrics and trends
4. **Nessi Performance Dashboard** - System performance metrics
5. **Nessi Anomaly Detection Dashboard** - Anomaly detection metrics and trends

## Setup

The Grafana dashboards are automatically provisioned when you run Nessi with the included Docker Compose configuration. The dashboards connect to Prometheus, which collects metrics from the Nessi application.

### Prerequisites

- Docker and Docker Compose installed
- Nessi running with the default configuration

### Accessing Grafana

1. Start Nessi using Docker Compose:
   ```
   docker-compose up -d
   ```

2. Access Grafana at:
   ```
   http://localhost:3000
   ```

3. Log in with the default credentials:
   - Username: `admin`
   - Password: `admin`

4. Navigate to "Dashboards" > "Nessi" to see all available dashboards

## Dashboard Descriptions

### Nessi Dashboard

The main dashboard provides an overview of key metrics:

- System health status
- Active alerts
- Data quality scores
- Resource utilization
- Recent validation results

### Nessi Alerts Dashboard

The alerts dashboard shows:

- Active alerts by severity
- Alert history and trends
- Alert resolution times
- Top alert sources
- Alert categories

### Nessi Data Quality Dashboard

The data quality dashboard displays:

- Quality scores by table/dataset
- Trend analysis of quality metrics
- Validation failure rates
- Data completeness metrics
- Accuracy and consistency scores

### Nessi Performance Dashboard

The performance dashboard monitors:

- CPU, memory, and disk usage
- Query performance metrics
- Response times
- Throughput statistics
- Resource bottlenecks

### Nessi Anomaly Detection Dashboard

The anomaly detection dashboard highlights:

- Detected anomalies by type
- Anomaly severity distribution
- Historical anomaly patterns
- False positive rates
- Detection sensitivity metrics

## Custom Dashboards

You can create custom dashboards by:

1. Duplicating an existing dashboard
2. Adding new panels with custom metrics
3. Configuring alerts for specific thresholds
4. Saving your custom dashboard

## Prometheus Integration

Nessi exposes metrics in Prometheus format at:

```
http://localhost:8080/metrics
```

You can configure Prometheus to scrape these metrics by adding the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'nessi'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:8080']
```

## Troubleshooting

Common issues:

- **No data in dashboards**: Check that Prometheus is correctly scraping metrics
- **Missing metrics**: Verify that the metric collection is enabled in Nessi configuration
- **Dashboard errors**: Ensure Grafana can connect to Prometheus data source
