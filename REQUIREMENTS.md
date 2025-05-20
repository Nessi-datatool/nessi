# Nessi Project Requirements

This document outlines the features and requirements of the Nessi open-source project.

## Project Overview
Nessi is an open-source data quality and Delta Lake management tool with a Go-first architecture. Python components are optional and can be enabled by users as needed. The goal is to become the go-to tool for Delta Lake data quality and management.

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

### CLI Monitoring
- ✅ Table health scores and trends via CLI commands
- ✅ Quality metric tracking over time with exportable data
- ✅ Recent validation failures reporting
- ✅ Table health and ingestion metrics via CLI

### Monitoring Capabilities
- ✅ Core monitoring functionality via CLI commands
- ✅ Metrics export for external visualization
- ✅ System health status commands

### Alerting
- ✅ Email notifications for critical issues
- ✅ Threshold-based triggers
- ✅ Daily/Weekly summary reports
- ✅ Customizable alert templates

### Monitoring History
- ✅ 30-day metrics retention
- ✅ Trend visualization and comparison
- ✅ Performance impact analysis

### CLI Reports
- ✅ Quality score summaries via CLI
- ✅ Table and column level summaries
- ✅ Trend data for external visualization
- ✅ Distribution data export
- ✅ Pattern frequency analysis
- ✅ Failure rate tracking
- ✅ Export to JSON/CSV for automation and external reporting
- ✅ Machine-readable output for integration with other tools

### Security
- ✅ Secure local file operations
- ✅ Environment variable-based authentication for cloud services
- ✅ API key management for external services
- ✅ Basic access controls for CLI commands

### Containerization
- ✅ 100% Docker-native for easy installation
- ✅ No JVM setup or Spark tuning required
- ✅ Single container deployment for CLI tools

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
- ✅ ASCII schema trees and CLI heatmaps for improved visualization
- ✅ Terraform-style CLI flags for automation-friendly interface
- ✅ Local-only data handling for privacy and security
- ✅ GitHub usage badge and telemetry opt-in (with privacy controls)
- ✅ Quickstart demo with embedded examples

### dbt Integration
- ✅ Opt-in dbt plugin for data quality and profiling
- ✅ Execute Nessi.dev data quality rules against dbt models
- ✅ Enhanced model selection with dbt syntax (e.g., tag:daily,+downstream)
- ✅ Generate data profiles for Delta tables associated with dbt models
- ✅ Automatic mapping between dbt models and Delta tables
- ✅ Export results in various formats (JSON, CSV, table)
- ✅ Generate artifacts for dbt Docs integration
- ✅ Alerting capabilities via email

### Workflow Orchestration Integration
- ✅ Apache Airflow operators and sensors for Nessi.dev tasks
- ✅ Prefect tasks and flows integration
- ✅ Dagster ops and resources integration
- ✅ Kubernetes operators for Nessi.dev jobs
- ✅ Workflow status monitoring and callbacks
- ✅ Pipeline-aware data quality checks
### Collaboration & Governance
- ✅ Shared rule libraries and configurations

## Architecture Requirements

### Primary Go Components
1. ✅ Data Quality Engine
2. ✅ Monitoring Service
3. ✅ Report Generator
4. ✅ Delta Lake Connector
5. ✅ Security Manager
6. ✅ CLI Interface
7. ✅ Email Alerting System

### Extension System
- ✅ Plugin architecture for enabling/disabling features
- ✅ Feature flags for toggling functionality
- ✅ Dynamic configuration loading
- ✅ Custom rule extensions
- ✅ Integration points for external systems

### Community
- ✅ Documentation and tutorials
- ✅ Example projects and use cases
- ✅ GitHub repository with issue tracking
- ✅ Slack/Discord community support channels

### Testing & Quality Assurance
- ✅ Comprehensive unit test coverage
- ✅ Integration tests for key components
- ✅ DBT plugin test coverage (AlertManager, selection system, validator functions)
- ✅ Profiler test coverage
- ✅ Artifacts generation test coverage
- ✅ Test strategy documentation
- ✅ In-memory testing for Freshness SLA functionality
- ✅ Mock connectors for reliable test execution