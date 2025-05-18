# Grafana and Prometheus Dashboarding in Nessi OSS

Advanced dashboarding and monitoring features (Grafana and Prometheus integration) are only available in the LakeDiff Enterprise version of Nessi. The open-source version provides a basic built-in dashboard for core monitoring functionality.

For advanced visualization, alerting, and analytics, please contact the LakeDiff team for enterprise licensing and support.


This document explains how to use the Grafana dashboards included with Nessi for monitoring alerts and metrics.

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
- Table sizes and record counts
- Error rates
- Operation latencies
- System resource usage

### Nessi Alerts Dashboard

This dashboard focuses on the alerting system:
- Alerts by severity (critical, warning, info)
- Alerts by status (active, acknowledged, resolved, silenced)
- Active alerts by type
- Alert rule evaluation rates
- Top alert sources

Key panels:
- **Alerts by Severity** - Time series showing alert counts by severity level
- **Active Critical Alerts** - Gauge showing the number of active critical alerts
- **Alert Rule Evaluation Rate** - Shows how frequently alert rules are being evaluated

### Nessi Data Quality Dashboard

This dashboard focuses on data quality metrics:
- Overall data quality score
- Data quality by table
- Completeness and accuracy metrics
- Rule validation failures
- Anomalies detected

Key panels:
- **Overall Data Quality Score** - Gauge showing the average quality score
- **Data Quality Score by Table** - Time series showing quality scores by table
- **Rule Validation Failures by Table** - Shows which tables have the most rule failures

### Nessi Performance Dashboard

This dashboard focuses on system performance:
- Operation latencies
- API request rates
- Error rates
- Resource usage (CPU, memory)
- Table sizes and record counts

Key panels:
- **Operation Latency** - 95th percentile latency for different operations
- **API Request Rate** - Requests per second by endpoint
- **Memory Usage** - Memory consumption over time

### Nessi Anomaly Detection Dashboard

This dashboard focuses on anomaly detection:
- Anomalies by table and type
- Top columns with anomalies
- Anomaly detection performance
- Specific anomaly type counts

Key panels:
- **Total Anomalies by Table** - Time series showing anomaly counts by table
- **Anomaly Distribution by Type** - Pie chart showing distribution of anomaly types
- **Anomaly Detection Duration** - Performance metrics for anomaly detection algorithms

## Using Dashboards with Alerts

The Grafana dashboards integrate with Nessi's alerting system to provide visual monitoring of alerts and metrics. Here's how to use them effectively:

### Monitoring Alerts

1. **View Active Alerts**: The Alerts Dashboard shows all active alerts, their severity, and status.

2. **Filter Alerts**: Use the dashboard filters to focus on specific alert types, severities, or sources.

3. **Track Alert Trends**: The time series panels show how alert counts change over time, helping identify patterns.

### Alert Actions from Grafana

While you can view alerts in Grafana, actions like acknowledging or resolving alerts should be performed through:

1. **Nessi Web Dashboard**: Use the Nessi alerts dashboard at `/alerts` for managing alerts.

2. **Nessi API**: Use the alerts API endpoints for programmatic control.

### Setting Up Grafana Alerting

You can also configure Grafana to send notifications based on the metrics displayed in the dashboards:

1. Hover over a panel and click the edit icon
2. Go to the "Alert" tab
3. Configure alert conditions based on the panel's metrics
4. Set up notification channels (email, Slack, etc.)

## Customizing Dashboards

You can customize the provided dashboards or create new ones:

1. Click the settings icon on any dashboard
2. Select "Save As" to create a copy
3. Modify panels, add new ones, or adjust existing queries
4. Save your customized dashboard

## Prometheus Integration

The dashboards use the following Prometheus metrics from Nessi:

- `nessi_alerts_total` - Alert counts by various labels
- `nessi_data_quality_score` - Data quality scores
- `nessi_operation_latency_seconds` - Operation latencies
- `nessi_api_requests_total` - API request counts
- `nessi_errors_total` - Error counts
- `nessi_table_size_bytes` - Table sizes
- `nessi_record_count` - Record counts
- `nessi_rule_validation_failures_total` - Rule validation failures
- `nessi_anomalies_detected_total` - Anomaly counts
- `nessi_alert_rule_evaluations_total` - Alert rule evaluation counts
- `nessi_alert_rule_evaluation_errors_total` - Alert rule evaluation errors
- `nessi_anomaly_detection_duration_seconds` - Anomaly detection duration
- `nessi_anomaly_detection_runs_total` - Anomaly detection run counts
- `nessi_anomaly_detection_failures_total` - Anomaly detection failures

## Troubleshooting

### Dashboard Shows No Data

1. Check if Prometheus is running:
   ```
   docker-compose ps
   ```

2. Verify Nessi is exposing metrics:
   ```
   curl http://localhost:9090/metrics
   ```

3. Check Prometheus target status in Grafana:
   - Go to Configuration > Data Sources > Prometheus
   - Click "Explore" and run a simple query like `up`

### Missing Dashboards

If dashboards are not automatically provisioned:

1. Import them manually:
   - Go to Dashboards > Import
   - Upload the JSON files from `config/grafana/dashboards/`

## Further Customization

For advanced customization:

1. **Add New Metrics**: Extend Nessi to expose additional Prometheus metrics

2. **Create New Dashboards**: Design custom dashboards for specific monitoring needs

3. **Set Up Additional Alert Channels**: Configure Grafana to send alerts to additional notification channels

4. **Use Variables**: Add dashboard variables to make dashboards more dynamic and reusable
