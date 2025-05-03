# Nessi Monitoring System

A comprehensive monitoring system for Nessi that provides real-time insights into data processing, quality, and performance metrics.

## Features

- **Real-time Metrics Monitoring**
  - Table scan statistics
  - Data quality scores
  - Processing performance
  - Custom report generation
  - Detailed analysis metrics

- **Detailed Analysis**
  - Scan efficiency tracking
  - Partition utilization monitoring
  - Index usage analysis
  - Data health assessment
  - Performance status tracking
  - Anomaly detection
  - Resource utilization monitoring

- **Visualization**
  - Interactive Grafana dashboard
  - Real-time metrics display
  - Trend analysis
  - Threshold-based alerts
  - Customizable panels

## Components

### 1. Metrics Generator (`metrics_test.py`)
Generates test metrics for demonstration purposes:
- Table scan metrics
- Data quality scores
- Processing time measurements
- Custom report counts
- Detailed analysis metrics

### 2. Grafana Dashboard (`grafana_dashboard.json`)
Interactive dashboard with:
- Table Scan Reports section
- Data Quality Reports section
- Performance Reports section
- Custom Reports section
- Real-time metrics visualization
- Trend analysis panels

### 3. Prometheus Configuration (`prometheus.yml`)
Metrics collection configuration:
- 5-second scrape interval
- Metric relabeling for better organization
- Instance labeling
- Metric type categorization

### 4. Docker Compose Setup (`docker-compose.test.yml`)
Containerized environment with:
- Metrics test service
- Prometheus server
- Grafana dashboard
- Network configuration
- Volume mounts

## Getting Started

### Prerequisites
- Docker
- Docker Compose
- Python 3.8+

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd backend
```

2. Start the monitoring system:
```bash
docker-compose -f docker-compose.test.yml up
```

3. Access the services:
- Grafana Dashboard: http://localhost:3000
- Prometheus: http://localhost:9090
- Metrics: http://localhost:8000/metrics

### Default Credentials
- Grafana:
  - Username: admin
  - Password: admin

## Metrics Overview

### Table Scan Metrics
- `nessi_table_scan_total_rows`: Total rows scanned
- `nessi_scan_efficiency`: Scan efficiency score (0-100)
- `nessi_partition_utilization`: Partition utilization score (0-100)
- `nessi_index_usage`: Index usage efficiency score (0-100)

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

The system simulates different operational scenarios:

1. **Normal Operation**
   - Optimal performance
   - High data quality
   - Efficient resource utilization

2. **High Load**
   - Increased processing time
   - Moderate performance impact
   - Higher resource utilization

3. **Data Quality Issues**
   - Reduced data quality scores
   - Increased anomaly detection
   - Performance impact

4. **Performance Degradation**
   - Significant processing delays
   - Resource constraints
   - System bottlenecks

5. **Custom Analysis**
   - Specialized metrics
   - Advanced analysis
   - Custom reporting

## Alerting

The system includes threshold-based alerts for:
- High row counts
- Low data quality scores
- Excessive processing time
- Resource utilization issues
- Anomaly detection

## Customization

### Dashboard Customization
- Modify panel layouts
- Adjust thresholds
- Add new metrics
- Customize visualizations

### Metric Configuration
- Adjust scrape intervals
- Modify metric labels
- Add new metrics
- Configure alert rules

## Troubleshooting

### Common Issues

1. **Metrics Not Showing**
   - Check Prometheus connection
   - Verify metrics endpoint
   - Check scrape configuration

2. **Dashboard Not Loading**
   - Verify Grafana credentials
   - Check dashboard JSON
   - Validate data source

3. **High Resource Usage**
   - Adjust scrape intervals
   - Optimize metric collection
   - Review alert rules

### Logs
- Metrics test service: `docker-compose logs metrics-test`
- Prometheus: `docker-compose logs prometheus`
- Grafana: `docker-compose logs grafana`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This software is provided under a 14-day trial period. After the trial period, a license is required to continue using the software.

## Support

For support, please contact [support contact information] 