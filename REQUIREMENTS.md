# Nessi.dev Project Requirements

## Project Overview
Nessi.dev is a data quality and Delta Lake management tool with a Go-first architecture. Python components are optional and can be enabled by users as needed. The goal is to become the go-to tool for Delta Lake data quality and management.

## Core Features

### Full Schema Management
- ✅ Schema evolution tracking and history
- ✅ Schema validation on write operations
- ✅ Field-level metadata exploration
- ✅ Partition intelligence and Z-ordering awareness
- ✅ Metadata inspection and optimization hints

### Version Control
- ✅ Transaction history with commit details
- ✅ Data rollback capabilities (up to 30 days)
- ✅ Version comparison with change summaries

### Time Travel
- ✅ Point-in-time querying (up to 30 days history)
- ✅ Timestamp and version-based access
- ✅ Historical state reconstruction

### Automated Profiling
- ✅ Statistical summaries (min, max, mean, median)
- ✅ Distribution analysis with histograms
- ✅ Null percentage and unique value ratios
- ✅ Type inference and consistency checks
- ✅ Completeness, accuracy, and consistency scoring

### Anomaly Detection
- ✅ Outlier detection using z-score and IQR methods
- ✅ Sudden change detection between runs
- ✅ Constant changes detection between runs
- ✅ Sudden spikes/dips detection between runs
- ✅ Oscillations detection between runs
- ✅ Historical trend deviation alerts

### Pattern Recognition
- ✅ Regex pattern extraction and matching
- ✅ Common format validation (emails, dates, IDs)
- ✅ Data format consistency monitoring

### Rule Validation
- ✅ Predefined rules library (null checks, range validation, length, regex, enum)
- ✅ Custom rule creation via YAML/SQL/Python
- ✅ Rule execution history and impact analysis
- ✅ Row-level rule validation with custom logic

### Real-time Dashboard
- ✅ Table health scores and trends
- ✅ Quality metric tracking over time
- ✅ Recent validation failures
- ✅ Live metrics for table health and ingestion

### Grafana Integration
- ✅ 5 pre-built dashboard templates
- ✅ Custom visualization options
- ✅ Shareable insights with team members
- ✅ Prometheus integration for alerts and historical trends

### Advanced Alerting
- ✅ Email notifications for critical issues
- ✅ Threshold-based triggers
- ✅ Daily/Weekly summary reports
- ✅ Integration with collaboration tools (Slack, Teams)
- ✅ Configurable alert routing and escalation policies
- ✅ Intelligent alerting based on historical trends and anomalies
  - ✅ Outlier detection using statistical analysis
  - ✅ Trend deviation detection for gradual changes
  - ✅ Seasonal pattern recognition (hourly, daily, weekly)
  - ✅ Automatic rule creation and updating based on historical data
  - ✅ Self-tuning sensitivity configuration

### Monitoring History
- ✅ 30-day metrics retention
- ✅ Trend visualization and comparison
- ✅ Performance impact analysis
- ✅ Historical data analysis for intelligent alerting

### Interactive Reports
- ✅ Quality score cards with drill-down
- ✅ Table and column level summaries
- ✅ Trend charts and comparison views
- ✅ Distribution histograms
- ✅ Pattern frequency charts
- ✅ Failure rate trends
- ✅ HTML and PDF reports with summary cards, histograms, and partition heatmaps
- ✅ Export to JSON/CSV for automation
- ✅ Shareable links with expiration

### Security
- ✅ SSL/TLS encrypted interfaces
- ✅ User authentication with JWT
- ✅ API key management
- ✅ Role-based access control (RBAC) for APIs and dashboards
- ✅ API rate limiting and optional IP allowlisting
- ✅ Audit logging for key operations

### Containerization
- 100% Docker-native for easy installation
- No JVM setup or Spark tuning required
- Multi-service orchestration via Docker Compose

### Enhanced Multi-format Support
- Delta Lake (primary)
- Parquet files
- CSV with automatic schema inference
- Auto schema inference and fallback logic
- Handles nested schemas and edge cases
- Configuration options for schema inference (e.g., date formats)
- Basic data quality checks across different formats

### Developer Experience
- CLI-first with one-line validation commands
- Batch processing scripts
- Configuration management
- Python API for programmatic access to all features
- Webhook integration
- Custom extension support
- Embeddable in Airflow, GitHub Actions, or custom scripts
- Interactive examples
- API reference
- Best practices guides

### Collaboration & Governance
- Role-Based Access Control (RBAC) for the CLI/API
- Audit Logging for key operations (configuration changes, rule executions)
- Team-based permissions and access controls
- Shared rule libraries and configurations

## Architecture Requirements

### Primary Go Components (Required)
1. ✅ API Gateway
2. ✅ Data Quality Engine
3. ✅ Monitoring Service
4. ✅ Report Generator
5. ✅ Delta Lake Connector
6. ✅ Security Manager
7. ✅ CLI Interface
8. ✅ Alerting System
   - ✅ Standard alerting with multiple notification channels
   - ✅ Intelligent alerting with pattern recognition

### Optional Python Extensions
1. Advanced Delta Lake Features
2. ML-based Anomaly Detection
3. Advanced Statistical Analysis
4. Custom Rule Execution Engine
5. Advanced Visualization Components

### Extension System
- Plugin architecture for enabling/disabling features
- Feature flags for toggling functionality
- Dynamic configuration loading
- Custom rule extensions
- Integration points for external systems