# Grafana Dashboard Documentation

## Overview

The nessi.dev Reports Dashboard provides a comprehensive view of system metrics and performance indicators. It's designed to help monitor and analyze various aspects of the nessi.dev system in real-time.

## Dashboard Structure

### 1. Table Scan Reports Section

#### Total Rows Scanned Panel
- **Type**: Time Series
- **Metric**: `nessi_table_scan_total_rows`
- **Description**: Shows the total number of rows scanned over time
- **Visualization**: Line graph with smooth interpolation
- **Thresholds**: 
  - Green: Normal operation
  - Red: High load conditions

#### Data Quality Score Panel
- **Type**: Gauge
- **Metric**: `nessi_data_quality_score`
- **Description**: Displays current data quality score
- **Thresholds**:
  - Green: 90-100%
  - Yellow: 70-90%
  - Red: <70%

### 2. Data Quality Reports Section

#### Data Quality Trends Panel
- **Type**: Time Series
- **Metric**: `nessi_data_quality_score`
- **Description**: Shows data quality trends over time
- **Visualization**: Line graph with smooth interpolation
- **Unit**: Percentage

#### Processing Time Panel
- **Type**: Gauge
- **Metric**: `nessi_processing_time_seconds`
- **Description**: Displays current processing time
- **Unit**: Seconds
- **Thresholds**:
  - Green: <2s
  - Yellow: 2-4s
  - Red: >4s

### 3. Performance Reports Section

#### Processing Time Trends Panel
- **Type**: Time Series
- **Metric**: `nessi_processing_time_seconds`
- **Description**: Shows processing time trends over time
- **Visualization**: Line graph with smooth interpolation
- **Unit**: Seconds

#### Custom Reports Generated Panel
- **Type**: Gauge
- **Metric**: `nessi_custom_report_count`
- **Description**: Displays number of custom reports generated
- **Unit**: Count

### 4. Custom Reports Section

#### Custom Report Trends Panel
- **Type**: Time Series
- **Metric**: `nessi_custom_report_count`
- **Description**: Shows custom report generation trends
- **Visualization**: Line graph with smooth interpolation

#### Data Quality Score Panel
- **Type**: Gauge
- **Metric**: `nessi_data_quality_score`
- **Description**: Displays current data quality score
- **Thresholds**:
  - Green: 90-100%
  - Yellow: 70-90%
  - Red: <70%

## Dashboard Configuration

### General Settings
- **Refresh Interval**: 5 seconds
- **Time Range**: Last 1 hour
- **Theme**: Dark
- **Timezone**: System default

### Panel Configuration

#### Time Series Panels
- **Line Width**: 2
- **Fill Opacity**: 20%
- **Point Size**: 5
- **Line Interpolation**: Smooth
- **Show Points**: Never
- **Legend**: Bottom
- **Tooltip**: Single

#### Gauge Panels
- **Orientation**: Auto
- **Thresholds**: 3 levels
- **Show Threshold Labels**: False
- **Show Threshold Markers**: True

## Integration

### Prometheus Data Source
- **Name**: Prometheus
- **Type**: Prometheus
- **URL**: http://prometheus:9090
- **Access**: Server
- **Scrape Interval**: 5s

### Alerting
- **Alert Rules**: Configured in Prometheus
- **Notification Channels**: Email
- **Default Recipient**: admin@example.com

## Customization

### Adding New Panels

1. Click "Add Panel" in the dashboard
2. Select panel type
3. Configure data source and query
4. Set visualization options
5. Add to appropriate section

### Modifying Existing Panels

1. Click panel title
2. Select "Edit"
3. Modify configuration
4. Save changes

### Dashboard Variables

Currently no variables are configured. To add:
1. Go to Dashboard Settings
2. Select "Variables"
3. Add new variable
4. Configure options

## Best Practices

### Panel Layout
- Group related metrics
- Use consistent sizing
- Maintain clear hierarchy
- Leave space for expansion

### Visualization
- Use appropriate chart types
- Maintain consistent color schemes
- Include clear labels
- Add helpful tooltips

### Performance
- Limit data points
- Use appropriate refresh rates
- Optimize queries
- Monitor dashboard load time

## Troubleshooting

### Common Issues

1. **Panels Not Loading**
   - Check Prometheus connection
   - Verify metric names
   - Check time range
   - Validate permissions

2. **Incorrect Data**
   - Verify metric collection
   - Check scrape configuration
   - Validate time ranges
   - Review query syntax

3. **Performance Issues**
   - Reduce refresh rate
   - Optimize queries
   - Limit time range
   - Check system resources

### Logs
- Grafana logs: `docker-compose logs grafana`
- Prometheus logs: `docker-compose logs prometheus`
- Metrics logs: `docker-compose logs metrics-test`

## Maintenance

### Regular Tasks
1. Monitor dashboard performance
2. Review alert configurations
3. Update panel configurations
4. Check data source connections

### Backup
1. Export dashboard JSON
2. Save configuration files
3. Document customizations
4. Version control changes

## Examples

### Example 1: Table Scan Panel Configuration
```json
{
  "title": "Total Rows Scanned",
  "type": "timeseries",
  "datasource": "Prometheus",
  "targets": [
    {
      "expr": "nessi_table_scan_total_rows",
      "legendFormat": "Rows"
    }
  ],
  "fieldConfig": {
    "defaults": {
      "custom": {
        "drawStyle": "line",
        "lineInterpolation": "smooth",
        "barAlignment": 0,
        "lineWidth": 2,
        "fillOpacity": 20,
        "gradientMode": "none",
        "spanNulls": false,
        "showPoints": "never",
        "pointSize": 5,
        "stacking": {
          "mode": "none",
          "group": "A"
        }
      }
    }
  }
}
```

### Example 2: Data Quality Gauge Configuration
```json
{
  "title": "Data Quality Score",
  "type": "gauge",
  "datasource": "Prometheus",
  "targets": [
    {
      "expr": "nessi_data_quality_score",
      "legendFormat": "Quality"
    }
  ],
  "fieldConfig": {
    "defaults": {
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {"color": "red", "value": null},
          {"color": "yellow", "value": 70},
          {"color": "green", "value": 90}
        ]
      },
      "unit": "percent"
    }
  }
}
```

### Example 3: Processing Time Panel
```json
{
  "title": "Processing Time Trends",
  "type": "timeseries",
  "datasource": "Prometheus",
  "targets": [
    {
      "expr": "nessi_processing_time_seconds",
      "legendFormat": "Time"
    }
  ],
  "fieldConfig": {
    "defaults": {
      "custom": {
        "drawStyle": "line",
        "lineInterpolation": "smooth",
        "barAlignment": 0,
        "lineWidth": 2,
        "fillOpacity": 20,
        "gradientMode": "none",
        "spanNulls": false,
        "showPoints": "never",
        "pointSize": 5
      },
      "unit": "s"
    }
  }
}
```

### Example 4: Custom Report Count Panel
```json
{
  "title": "Custom Reports Generated",
  "type": "gauge",
  "datasource": "Prometheus",
  "targets": [
    {
      "expr": "nessi_custom_report_count",
      "legendFormat": "Reports"
    }
  ],
  "fieldConfig": {
    "defaults": {
      "thresholds": {
        "mode": "absolute",
        "steps": [
          {"color": "red", "value": null},
          {"color": "yellow", "value": 5},
          {"color": "green", "value": 10}
        ]
      },
      "unit": "short"
    }
  }
}
```

### Example 5: Dashboard Layout
```json
{
  "panels": [
    {
      "title": "Table Scan Reports",
      "type": "row",
      "gridPos": {"h": 1, "w": 24, "x": 0, "y": 0}
    },
    {
      "title": "Total Rows Scanned",
      "type": "timeseries",
      "gridPos": {"h": 8, "w": 12, "x": 0, "y": 1}
    },
    {
      "title": "Data Quality Score",
      "type": "gauge",
      "gridPos": {"h": 8, "w": 12, "x": 12, "y": 1}
    }
  ]
}
```

### Example 6: Dashboard Variables
```json
{
  "templating": {
    "list": [
      {
        "name": "instance",
        "type": "query",
        "datasource": "Prometheus",
        "query": "label_values(nessi_table_scan_total_rows, instance)"
      },
      {
        "name": "metric_type",
        "type": "query",
        "datasource": "Prometheus",
        "query": "label_values(nessi_table_scan_total_rows, metric_type)"
      }
    ]
  }
}
```

### Example 7: Alert Configuration
```json
{
  "alert": {
    "name": "High Processing Time",
    "message": "Processing time exceeds threshold",
    "conditions": [
      {
        "evaluator": {
          "params": [4],
          "type": "gt"
        },
        "operator": {
          "type": "and"
        },
        "query": {
          "params": ["A", "5m", "now"]
        },
        "reducer": {
          "params": [],
          "type": "avg"
        },
        "type": "query"
      }
    ],
    "frequency": "1m",
    "handler": 1,
    "noDataState": "no_data",
    "notifications": []
  }
}
```

### Example 8: Dashboard Annotations
```json
{
  "annotations": {
    "list": [
      {
        "name": "Anomalies",
        "datasource": "Prometheus",
        "query": "nessi_anomaly_count > 0",
        "enable": true,
        "iconColor": "red",
        "showIn": 0
      }
    ]
  }
}
```

### Example 9: Time Range Configuration
```json
{
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ]
  }
}
```

### Example 10: Dashboard Links
```json
{
  "links": [
    {
      "title": "Prometheus",
      "url": "http://localhost:9090",
      "icon": "external link",
      "targetBlank": true
    },
    {
      "title": "Metrics",
      "url": "http://localhost:8000/metrics",
      "icon": "external link",
      "targetBlank": true
    }
  ]
} 