# Root Cause Analysis (RCA) for Anomaly Diagnosis

## Overview

The Root Cause Analysis (RCA) feature in Nessi.dev helps users quickly diagnose the underlying causes of data quality anomalies. When anomalies are detected, understanding why they occurred is crucial for effective remediation. RCA capabilities allow users to pinpoint the source of data quality issues, saving time and effort in debugging.

## Key Features

### Comprehensive Analysis Types

- **Schema Change Detection**: Identifies recent schema changes that might have impacted data quality
- **Data Quality Rule Correlation**: Links anomalies to recent rule validation failures
- **Lineage-based Analysis**: Traces issues to upstream data sources and processing steps
- **System Performance Correlation**: Connects anomalies with system resource constraints or performance issues

### Multiple Output Formats

- **Structured JSON Reports**: Machine-readable format for integration with other tools
- **Human-readable Text**: Clear, concise summaries for command-line usage
- **HTML Dashboard Integration**: Visual representation of RCA results in the Nessi.dev dashboard

### Actionable Insights

- **Recommended Actions**: Specific steps to address the identified root causes
- **Affected Tables**: List of tables impacted by the anomaly
- **Related Anomalies**: Identification of connected anomalies for holistic problem-solving

### Integration Capabilities

- **Grafana Dashboard Links**: Direct links to relevant Grafana dashboards for further investigation
- **Alerting Integration**: Option to trigger alerts based on RCA findings
- **Feature Flag Control**: Enable/disable RCA functionality as needed

## Usage

### CLI Command

```bash
# Basic usage
nessi rca <anomaly_id>

# Output as JSON
nessi rca <anomaly_id> --format=json

# Save output to a file
nessi rca <anomaly_id> --format=html --output=report.html

# Send alert for critical findings
nessi rca <anomaly_id> --alert
```

### Configuration

RCA can be configured in the Nessi.dev configuration file:

```yaml
rca:
  enable_rca: true
  max_history_days: 7
  detail_level: "standard"  # basic, standard, detailed
  include_lineage: true
  include_schema_change: true
  alert_threshold: 2  # 0=info, 1=warning, 2=error, 3=critical
```

## Implementation Details

The RCA feature is implemented as a dedicated service within the Nessi.dev architecture:

- **RCA Service**: Core component that performs the analysis
- **Analyzer**: Main logic for identifying potential root causes
- **CLI Integration**: Command-line interface for triggering analysis
- **Dashboard Integration**: Visual representation of results in the web UI

### Dependencies

The RCA service integrates with several other Nessi.dev components:

- **Monitoring Service**: For anomaly data and historical metrics
- **Delta Lake Connector**: For table metadata and schema history
- **Data Quality Engine**: For rule validation history
- **Alerting System**: For sending notifications about critical findings

## Example Output

### Text Format

```
Root Cause Analysis for Anomaly anom-20250514-001
Analysis Time: 2025-05-14 15:30:45

PRIMARY ROOT CAUSE:
  Type: schema_change
  Confidence: 80.0%
  Description: Column type change detected prior to anomaly
  Timestamp: 2025-05-14 13:30:45

OTHER POTENTIAL CAUSES:
  1. Upstream data quality rule failures detected (60.0% confidence)
  2. High system load during data processing (40.0% confidence)

AFFECTED TABLES:
  - sales/transactions
  - reporting/customer_metrics
  - analytics/customer_segmentation

RELATED ANOMALIES:
  - anom-20250514-002
  - anom-20250514-007

RECOMMENDED ACTIONS:
  1. Review recent schema changes to table sales/transactions
  2. Verify data type compatibility for column customer_id
  3. Check for missing data transformation steps in ETL pipeline

Grafana Dashboard: https://grafana.example.com/d/nessi-anomaly/anomaly-details?var-anomaly_id=anom-20250514-001
```

## Benefits

- **Faster Resolution**: Quickly identify and address the root causes of data quality issues
- **Reduced Debugging Time**: Eliminate manual investigation of complex data pipelines
- **Proactive Problem Prevention**: Identify patterns that lead to anomalies
- **Improved Data Quality**: Address systemic issues rather than symptoms
- **Enhanced Visibility**: Provide stakeholders with clear explanations of data quality issues
