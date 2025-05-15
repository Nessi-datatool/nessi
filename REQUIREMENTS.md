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
  - ✅ Root Cause Analysis (RCA) for anomaly diagnosis
    - ✅ Schema change detection and impact analysis
    - ✅ Interactive RCA dashboard with visualization of root causes, confidence scoring, and insights
    - ✅ Data quality rule failure correlation
    - ✅ Lineage-based upstream issue detection
    - ✅ System performance correlation
    - ✅ Actionable recommendations based on findings
    - ✅ Integration with Grafana dashboards

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
- ✅ 100% Docker-native for easy installation
- ✅ No JVM setup or Spark tuning required
- ✅ Multi-service orchestration via Docker Compose

### Enhanced Multi-format Support
- ✅ Delta Lake (primary)
- ✅ Parquet files
- ✅ CSV with automatic schema inference
- ✅ Auto schema inference and fallback logic
- ✅ Handles nested schemas and edge cases
- ✅ Configuration options for schema inference (e.g., date formats)
- ✅ Basic data quality checks across different formats

### Cloud Integration
- ✅ AWS S3 integration for Delta Lake tables
- ✅ Azure Blob Storage integration with latest SDK (v1.6.1+)
- ✅ GCP Cloud Storage integration
- ✅ Multi-cloud Delta Lake management
- ✅ Cloud-native authentication methods
- ✅ Cross-cloud data quality validation

### Data Catalog Integration
- ✅ AWS Glue Data Catalog integration
- ✅ Azure Purview integration
- ✅ Google Cloud Data Catalog integration
- ✅ Data asset discovery and connection
- ✅ Metadata utilization (descriptions, tags, lineage)
- ✅ Quality metrics and lineage publishing back to catalogs

### Data Lineage Visualization and Comparison
- ✅ Interactive lineage graph visualization
- ✅ Upstream and downstream dependency tracking
- ✅ Version comparison of lineage changes
- ✅ Impact analysis for schema changes
- ✅ Time-based lineage evolution view
- ✅ Export lineage to various formats

### Developer Experience
- ✅ CLI-first with one-line validation commands
- ✅ Batch processing scripts
- ✅ Configuration management
- ✅ Python API for programmatic access to all features
- ✅ Webhook integration
- ✅ Custom extension support
- ✅ Embeddable in Airflow, GitHub Actions, or custom scripts
- ✅ Interactive examples
- ✅ API reference
- ✅ Best practices guides

### dbt Integration
- ✅ Opt-in dbt plugin for data quality and profiling
- ✅ Execute Nessi.dev data quality rules against dbt models
- ✅ Enhanced model selection with dbt syntax (e.g., tag:daily,+downstream)
- ✅ Generate data profiles for Delta tables associated with dbt models
- ✅ Lineage-aware validation for root cause analysis
- ✅ Automatic mapping between dbt models and Delta tables
- ✅ Export results in various formats (JSON, CSV, table)
- ✅ Generate artifacts for dbt Docs integration
- ✅ Alerting capabilities via Slack and email

### Workflow Orchestration Integration
- ✅ Apache Airflow operators and sensors for Nessi.dev tasks
- ✅ Prefect tasks and flows integration
- ✅ Dagster ops and resources integration
- ✅ Kubernetes operators for Nessi.dev jobs
- ✅ Workflow status monitoring and callbacks
- ✅ Pipeline-aware data quality checks

### Collaboration & Governance
- ✅ Role-Based Access Control (RBAC) for the CLI/API
- ✅ Audit Logging for key operations (configuration changes, rule executions)
- ✅ Team-based permissions and access controls
- ✅ Shared rule libraries and configurations

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
9. ✅ Root Cause Analysis (RCA) Service
   - ✅ Anomaly diagnosis engine
   - ✅ Interactive dashboard integration
   - ✅ Historical insights and trend analysis
   - ✅ CLI accessibility via `nessi rca <anomaly_id>`
   - ✅ Feature flag control with `--enable-rca`
   - ✅ Structured JSON reports and HTML dashboard integration
   - ✅ Grafana dashboard linking for contextual metrics

### Optional Python Extensions
1. ✅ Advanced Delta Lake Features
2. ✅ ML-based Anomaly Detection
3. ✅ Advanced Statistical Analysis
4. ✅ Custom Rule Execution Engine
5. ✅ Advanced Visualization Components

### Extension System
- ✅ Plugin architecture for enabling/disabling features
- ✅ Feature flags for toggling functionality
- ✅ Dynamic configuration loading
- ✅ Custom rule extensions
- ✅ Integration points for external systems

### Testing & Quality Assurance
- ✅ Comprehensive unit test coverage
- ✅ Integration tests for key components
- ✅ DBT plugin test coverage (AlertManager, selection system, validator functions)
- ✅ Profiler test coverage
- ✅ Artifacts generation test coverage
- ✅ Test strategy documentation
- ✅ In-memory testing for Freshness SLA functionality
- ✅ Mock connectors for reliable test execution