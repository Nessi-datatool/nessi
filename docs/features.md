# Nessi Features Documentation

## Delta Lake Superpowers

### Native Delta Lake Support
- **Schema Evolution**: Automatically handle schema changes with `evolve_schema`
  ```python
  delta_manager.evolve_schema(path, {
      "column_name": {
          "type": "new_type",
          "nullable": True,
          "comment": "Updated column description"
      }
  })
  ```
- **Version Control**: Track and access historical versions
  ```python
  # Get version history
  history = delta_manager.get_history(path)
  
  # Time travel to specific version
  df = delta_manager.time_travel(path, version=2)
  # Or to specific timestamp
  df = delta_manager.time_travel(path, timestamp="2024-01-01")
  ```
- **Time Travel**: Access data at any point in time
  ```python
  # Access data as of specific version
  df = delta_manager.time_travel(path, version=5)
  ```

### Partition Intelligence
- **Automatic Partition Detection**
  ```python
  # Get partition candidates
  candidates = delta_manager._get_partition_candidates(path)
  ```
- **Z-Ordering Optimization**
  ```python
  # Optimize table with Z-ordering
  delta_manager.optimize_table(path, z_order_by=["column1", "column2"])
  ```

### Metadata and Optimization
- **Detailed Metadata Inspection**
  ```python
  # Get comprehensive table metadata
  metadata = delta_manager.get_metadata(path)
  ```
- **Optimization Hints**
  ```python
  # Get optimization recommendations
  analysis = delta_manager.analyze_table(path)
  ```

## Deep Data Quality Intelligence

### Quality Scoring
- **Completeness Score**: Measure data completeness
  ```python
  analyzer = QualityAnalyzer()
  results = analyzer.analyze_dataframe(df)
  completeness = results["completeness"]
  ```
- **Consistency Score**: Evaluate data consistency
  ```python
  consistency = results["consistency"]
  ```
- **Overall Quality Score**: Combined quality metric
  ```python
  quality = results["overall_quality"]
  ```

### Column Profiling
- **Basic Statistics**
  ```python
  # Get column statistics
  stats = analyzer._analyze_column(df["column"])
  print(f"Min: {stats['min']}, Max: {stats['max']}")
  print(f"Unique %: {stats['unique_percentage']}")
  print(f"Null %: {stats['null_percentage']}")
  ```

### Anomaly Detection
- **Outlier Detection**
  ```python
  # Get detected anomalies
  anomalies = results["anomalies"]
  ```
- **Pattern Recognition**
  ```python
  # Get detected patterns
  patterns = results["patterns"]
  ```

### Rule Validation
- **Custom Rule Validation**
  ```python
  rules = [
      {
          "name": "age_check",
          "condition": "age >= 0 and age <= 120",
          "description": "Age must be between 0 and 120"
      }
  ]
  validation_results = analyzer.validate_rules(df, rules)
  ```

## Real-Time Monitoring

### Metrics Collection
- **Table Health Metrics**
  ```python
  monitor = Monitor()
  monitor.start_scan()
  monitor.end_scan()
  ```
- **Ingestion Metrics**
  ```python
  monitor.start_validation()
  monitor.end_validation()
  ```

### Grafana Integration
- **Pre-configured Dashboards**
  - Table health dashboard
  - Data quality dashboard
  - Performance metrics dashboard

### Alerting
- **Threshold-based Alerts**
  ```yaml
  # alertmanager.yml
  routes:
    - match:
        severity: critical
      receiver: email
  ```
- **Multiple Notification Channels**
  - Email
  - Slack
  - Webhook

## Interactive Report Generation

### Report Types
- **HTML Reports**
  ```python
  generator = ReportGenerator()
  html_report = generator.generate_report(data, format="html")
  ```
- **PDF Reports**
  ```python
  pdf_report = generator.generate_report(data, format="pdf")
  ```
- **JSON/CSV Export**
  ```python
  json_report = generator.generate_report(data, format="json")
  ```

### Visualizations
- **Summary Cards**
  - Overall quality score
  - Completeness metrics
  - Consistency metrics
- **Distribution Charts**
  - Value distributions
  - Pattern distributions
- **Heatmaps**
  - Anomaly distribution
  - Partition distribution

## Enterprise-Grade Security

### Authentication & Authorization
- **JWT-based Authentication**
  ```python
  token = security.create_token(user_id, roles)
  ```
- **Role-Based Access Control**
  ```python
  @security.check_role(["admin"])
  def admin_function():
      pass
  ```

### Rate Limiting
- **API Rate Limiting**
  ```python
  @security.rate_limit_middleware(limit=100, window=60)
  def rate_limited_function():
      pass
  ```

### Network Security
- **SSL/TLS Encryption**
  ```python
  @security.ssl_middleware()
  def secure_endpoint():
      pass
  ```
- **IP Whitelisting**
  ```python
  @security.ip_whitelist_middleware()
  def restricted_endpoint():
      pass
  ```

## Containerized Simplicity

### Docker Deployment
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### Service Configuration
- **Nessi Service**
  - Port: 8443
  - Volumes: /app/data, /app/logs, /app/config
- **Redis Service**
  - Port: 6379
  - Persistent storage
- **Prometheus Service**
  - Port: 9090
  - Metrics collection
- **Grafana Service**
  - Port: 3000
  - Dashboard visualization
- **Alertmanager Service**
  - Port: 9093
  - Alert handling
- **Nginx Service**
  - Port: 443
  - SSL termination

## Data Format Versatility

### Supported Formats
- **Delta Lake**
  ```python
  df = spark.read.format("delta").load(path)
  ```
- **Parquet**
  ```python
  df = spark.read.parquet(path)
  ```
- **CSV**
  ```python
  df = spark.read.csv(path)
  ```

### Schema Handling
- **Auto Schema Inference**
  ```python
  df = spark.read.format("delta").load(path)
  schema = df.schema
  ```
- **Nested Schema Support**
  ```python
  # Handle complex nested structures
  df = spark.read.format("delta").load(path)
  ```

## Developer Focused

### CLI Interface
```bash
# Run data quality analysis
nessi analyze --format delta --path /data/table

# Generate report
nessi report --format html --output report.html

# Monitor table
nessi monitor --table my_table
```

### Python API
```python
from nessi import DeltaManager, QualityAnalyzer, ReportGenerator

# Initialize components
delta_manager = DeltaManager(spark)
analyzer = QualityAnalyzer()
generator = ReportGenerator()

# Use features
df = delta_manager.time_travel(path, version=2)
results = analyzer.analyze_dataframe(df)
report = generator.generate_report(results)
```

### Integration Examples
- **Airflow Integration**
  ```python
  from airflow import DAG
  from nessi.operators import NessiOperator
  
  with DAG('nessi_pipeline') as dag:
      analyze = NessiOperator(
          task_id='analyze',
          command='analyze',
          params={'path': '/data/table'}
      )
  ```
- **GitHub Actions**
  ```yaml
  - name: Run Nessi Analysis
    uses: nessi/action@v1
    with:
      command: analyze
      path: data/table
  ```

## Configuration

### Security Configuration
```json
{
  "jwt_secret": "your-secret",
  "token_expiry_hours": 24,
  "rate_limit": {
    "requests": 100,
    "window": 60
  },
  "ip_whitelist": ["192.168.1.0/24"]
}
```

### Monitoring Configuration
```json
{
  "metrics_port": 9090,
  "alert_log_path": "logs/alerts.log",
  "grafana": {
    "host": "localhost",
    "port": 3000
  },
  "prometheus": {
    "scrape_interval": "15s",
    "evaluation_interval": "15s"
  }
}
```

## Best Practices

### Data Quality
1. Run regular quality checks
2. Set up automated alerts
3. Monitor data completeness
4. Track consistency metrics

### Security
1. Use strong JWT secrets
2. Implement proper role hierarchy
3. Monitor rate limits
4. Regular security audits

### Performance
1. Use appropriate partitioning
2. Implement Z-ordering for frequent queries
3. Regular table optimization
4. Monitor query performance

### Maintenance
1. Regular backups
2. Monitor disk usage
3. Update dependencies
4. Review logs regularly 