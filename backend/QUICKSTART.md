# Quick Start Guide

## Prerequisites

Before you begin, ensure you have:
- Docker installed
- Docker Compose installed
- Git installed (for cloning the repository)

## Important: License Information

The system includes a 14-day trial period. After 14 days, you will need a license key to continue using the system:
1. The trial period starts from first activation
2. After 14 days, you must obtain a license key
3. Without a valid license key, the system will continue to operate but may have limited functionality
4. Contact support at support@nessi.com to obtain a license key

For detailed license information, see [LICENSE.md](../LICENSE.md).

## Installation Steps

1. **Clone the Repository**
```bash
git clone <repository-url>
cd backend
```

2. **Start the Monitoring System**
```bash
docker-compose -f docker-compose.test.yml up
```

3. **Access the Services**
- Open Grafana Dashboard: http://localhost:3000
  - Username: admin
  - Password: admin
- View Prometheus: http://localhost:9090
- Check Metrics: http://localhost:8000/metrics

## Basic Usage

### Viewing Metrics

1. **Access the Dashboard**
   - Open http://localhost:3000
   - Log in with admin/admin
   - Navigate to "Nessi Reports Dashboard"

2. **Monitor Key Metrics**
   - Table Scan Reports: View row counts and efficiency
   - Data Quality Reports: Check quality scores and trends
   - Performance Reports: Monitor processing times
   - Custom Reports: View generated report counts

### Understanding Scenarios

The system cycles through different scenarios every minute:

1. **Normal Operation (0-12 minutes)**
   - Optimal performance
   - High data quality
   - Green indicators

2. **High Load (12-24 minutes)**
   - Increased processing time
   - Yellow indicators
   - Higher resource usage

3. **Data Quality Issues (24-36 minutes)**
   - Reduced quality scores
   - Red indicators
   - Increased anomalies

4. **Performance Degradation (36-48 minutes)**
   - Slow processing
   - Resource constraints
   - Red indicators

5. **Custom Analysis (48-60 minutes)**
   - Specialized metrics
   - Advanced analysis
   - Mixed indicators

## Common Tasks

### 1. Check System Health
```bash
# View all service logs
docker-compose logs

# Check specific service
docker-compose logs prometheus
```

### 2. Restart Services
```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart grafana
```

### 3. Stop the System
```bash
# Stop all services
docker-compose down
```

## Troubleshooting

### Quick Fixes

1. **Dashboard Not Loading**
   - Check if Grafana is running: `docker-compose ps`
   - Verify credentials: admin/admin
   - Check browser console for errors

2. **Metrics Not Showing**
   - Verify Prometheus is running
   - Check metrics endpoint: http://localhost:8000/metrics
   - Review service logs

3. **High Resource Usage**
   - Adjust scrape intervals in prometheus.yml
   - Reduce dashboard refresh rate
   - Check system resources

## Next Steps

1. **Explore the Dashboard**
   - Review different panels
   - Check metric trends
   - Monitor alerts

2. **Customize Configuration**
   - Adjust alert thresholds
   - Modify scrape intervals
   - Add new metrics

3. **Learn More**
   - Read detailed documentation
   - Review example scenarios
   - Explore customization options

## Support

For additional help:
- Check the full documentation
- Review troubleshooting guides
- Contact support team

## Examples

### Example 1: Starting the System
```bash
# Clone the repository
git clone https://github.com/your-org/nessi.git
cd nessi/backend

# Start the services
docker-compose -f docker-compose.test.yml up
```

### Example 2: Accessing the Dashboard
```bash
# Open Grafana in your browser
open http://localhost:3000

# Login credentials
Username: admin
Password: admin
```

### Example 3: Viewing Metrics
```bash
# View raw metrics
curl http://localhost:8000/metrics

# Example output:
nessi_table_scan_total_rows 50000
nessi_data_quality_score 95.0
nessi_processing_time_seconds 2.5
```

### Example 4: Prometheus Queries
```promql
# Basic queries
nessi_table_scan_total_rows
nessi_data_quality_score
nessi_processing_time_seconds

# Aggregated queries
avg(nessi_data_quality_score)
max(nessi_processing_time_seconds)
```

### Example 5: Grafana Dashboard
```json
{
  "title": "Nessi Dashboard",
  "panels": [
    {
      "title": "Table Scan Reports",
      "type": "row"
    },
    {
      "title": "Total Rows Scanned",
      "type": "timeseries"
    }
  ]
}
```

### Example 6: Alert Configuration
```yaml
groups:
  - name: nessi-alerts
    rules:
      - alert: HighRowCount
        expr: nessi_table_scan_total_rows > 1000000
        for: 5m
        labels:
          severity: warning
```

### Example 7: Docker Compose
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
```

### Example 8: Environment Variables
```bash
# Set environment variables
export SCENARIO=normal_operation
export METRICS_INTERVAL=5s
export LOG_LEVEL=INFO
```

### Example 9: Custom Metrics
```python
# Define a new metric
custom_metric = Gauge('nessi_custom_metric', 'Custom metric description')

# Set metric value
custom_metric.set(0.9)
```

### Example 10: Troubleshooting
```bash
# View logs
docker-compose -f docker-compose.test.yml logs

# Check service status
docker-compose -f docker-compose.test.yml ps

# Restart services
docker-compose -f docker-compose.test.yml restart
``` 