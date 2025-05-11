# Nessi.dev Project Requirements

## Project Overview
Nessi.dev is a data quality and Delta Lake management tool with a Go-first architecture. Python components are optional and can be enabled by users as needed.

## Core Features

### Delta Lake Support
- Native Delta Lake support with schema evolution, version control, and time travel
- Partition intelligence and Z-ordering awareness
- Metadata inspection and optimization hints

### Data Quality Intelligence
- Completeness, accuracy, and consistency scoring
- Built-in profiling for every column (min, max, unique %, null %)
- Anomaly and outlier detection
- Pattern recognition (e.g., regex pattern mismatches, column formats)
- Row-level rule validation with custom logic

### Monitoring
- Live metrics for table health and ingestion
- Grafana dashboards with panel templates
- Prometheus integration for alerts and historical trends
- Supports threshold-based alerting and reporting

### Report Generation
- HTML and PDF reports with summary cards, histograms, and partition heatmaps
- Export to JSON/CSV for automation

### Security
- SSL/TLS encrypted interfaces
- Role-based access control (RBAC) for APIs and dashboards
- API rate limiting and optional IP allowlisting

### Containerization
- 100% Docker-native for easy installation
- No JVM setup or Spark tuning required
- Multi-service orchestration via Docker Compose

### Data Format Support
- Delta Lake, Parquet, and CSV support
- Auto schema inference and fallback logic
- Handles nested schemas and edge cases

### Developer Features
- CLI-first with Python API support
- Embeddable in Airflow, GitHub Actions, or custom scripts
- Detailed documentation and examples

## Architecture Requirements

### Primary Go Components (Required)
1. API Gateway
2. Data Quality Engine
3. Monitoring Service
4. Report Generator
5. Basic Delta Lake Connector
6. Security Manager
7. CLI Interface

### Optional Python Extensions
1. Advanced Delta Lake Features
2. ML-based Anomaly Detection
3. Advanced Statistical Analysis

### Extension System
- Plugin architecture for enabling/disabling features
- Feature flags for toggling functionality
- Dynamic configuration loading

## Pro Features (Future)
- Data contracts and schema diff automation
- SLA compliance monitoring
- Git integration for schema diffs and versioning audits
- BI tool connectors (e.g., dbt, Tableau lineage sync)