# Nessi Examples

This document provides practical examples of using Nessi for various data quality and observability tasks.

## Table of Contents

- [Basic Usage Examples](#basic-usage-examples)
- [Delta Lake Table Operations](#delta-lake-table-operations)
- [Schema Management Examples](#schema-management-examples)
- [Data Quality Examples](#data-quality-examples)
- [Metrics Collection Examples](#metrics-collection-examples)
- [Report Generation Examples](#report-generation-examples)
- [Integration Examples](#integration-examples)
- [Advanced Usage Examples](#advanced-usage-examples)
- [Workflow Examples](#workflow-examples)
- [Scripting Examples](#scripting-examples)

## Basic Usage Examples

### Getting Help

```bash
# Get general help
nessi --help

# Get help for a specific command
nessi tables --help

# Get help for a specific subcommand
nessi tables describe --help
```

### Version Information

```bash
# Get version information
nessi --version
```

### Configuration

```bash
# Initialize default configuration
nessi config init

# View current configuration
nessi config show

# Validate configuration
nessi config validate
```

## Delta Lake Table Operations

### Listing Tables

```bash
# List all tables in a directory
nessi tables list --path /path/to/data

# List tables with a specific pattern
nessi tables list --path /path/to/data --pattern "*_daily"

# List tables with detailed information
nessi tables list --path /path/to/data --detailed
```

### Describing Tables

```bash
# Get basic table information
nessi tables describe --path /path/to/table

# Get detailed table information
nessi tables describe --path /path/to/table --detailed

# Get table information in JSON format
nessi tables describe --path /path/to/table --output-format json
```

### Viewing Table Schema

```bash
# View table schema
nessi tables schema --path /path/to/table

# View schema in JSON format
nessi tables schema --path /path/to/table --output-format json

# View schema with column statistics
nessi tables schema --path /path/to/table --with-stats
```

### Viewing Table History

```bash
# View table history
nessi tables history --path /path/to/table

# View detailed history
nessi tables history --path /path/to/table --detailed

# Limit history entries
nessi tables history --path /path/to/table --limit 10
```

### Time Travel

```bash
# View table at a specific version
nessi tables show --path /path/to/table --version 3

# View table at a specific timestamp
nessi tables show --path /path/to/table --timestamp "2023-01-01T12:00:00Z"

# Compare table versions
nessi tables diff --path /path/to/table --from-version 2 --to-version 3
```

## Schema Management Examples

### Schema Validation

```bash
# Validate table schema against a schema file
nessi schema validate --path /path/to/table --schema-file /path/to/schema.json

# Validate with detailed output
nessi schema validate --path /path/to/table --schema-file /path/to/schema.json --detailed

# Validate with relaxed mode (ignore additional fields)
nessi schema validate --path /path/to/table --schema-file /path/to/schema.json --relaxed
```

### Schema Comparison

```bash
# Compare schemas of two tables
nessi schema compare --source /path/to/table1 --target /path/to/table2

# Compare with detailed output
nessi schema compare --source /path/to/table1 --target /path/to/table2 --detailed

# Compare and ignore specific fields
nessi schema compare --source /path/to/table1 --target /path/to/table2 --ignore-fields created_at,updated_at
```

### Schema Evolution

```bash
# Track schema evolution
nessi schema history --path /path/to/table

# View schema at a specific version
nessi schema show --path /path/to/table --version 3

# Compare schema versions
nessi schema diff --path /path/to/table --from-version 2 --to-version 3
```

### Schema Visualization

```bash
# Visualize schema as a tree
nessi schema tree --path /path/to/table

# Visualize schema with field types
nessi schema tree --path /path/to/table --with-types

# Export schema visualization to file
nessi schema tree --path /path/to/table --output /path/to/schema-tree.txt
```

## Data Quality Examples

### Defining Quality Rules

Create a file named `quality_rules.yaml`:

```yaml
# quality_rules.yaml
rules:
  - name: not_null
    description: Check if column values are not null
    columns:
      - id
      - name
    
  - name: unique
    description: Check if column values are unique
    columns:
      - id
    
  - name: range
    description: Check if column values are within range
    columns:
      - age
    parameters:
      min: 0
      max: 120
      
  - name: regex
    description: Check if column values match a pattern
    columns:
      - email
    parameters:
      pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
      
  - name: referential_integrity
    description: Check if column values exist in reference table
    columns:
      - customer_id
    parameters:
      reference_table: /path/to/customers
      reference_column: id
```

### Running Quality Checks

```bash
# Run quality checks with defined rules
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml

# Run checks with detailed output
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --detailed

# Run checks with specific rule categories
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --categories completeness,accuracy

# Run checks with sampling
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --sample-size 1000
```

### Viewing Quality Results

```bash
# View latest quality results
nessi quality results --path /path/to/table

# View results in JSON format
nessi quality results --path /path/to/table --output-format json

# View historical results
nessi quality history --path /path/to/table

# Compare quality results
nessi quality compare --path /path/to/table --from-date "2023-01-01" --to-date "2023-01-31"
```

### Quality Metrics

```bash
# View quality metrics
nessi quality metrics --path /path/to/table

# View specific metrics
nessi quality metrics --path /path/to/table --metrics completeness,accuracy

# View metrics over time
nessi quality metrics --path /path/to/table --period "last-30-days"

# Export metrics to CSV
nessi quality metrics --path /path/to/table --output-format csv --output /path/to/metrics.csv
```

## Metrics Collection Examples

### Collecting Metrics

```bash
# Collect all metrics
nessi metrics collect --path /path/to/table

# Collect specific metrics
nessi metrics collect --path /path/to/table --metrics volume,freshness,performance

# Collect metrics with tags
nessi metrics collect --path /path/to/table --tags environment=production,dataset=customers
```

### Viewing Metrics

```bash
# View all metrics
nessi metrics show --path /path/to/table

# View specific metrics
nessi metrics show --path /path/to/table --metrics volume,freshness

# View metrics for a specific period
nessi metrics show --path /path/to/table --from "2023-01-01" --to "2023-01-31"

# View metrics with filtering
nessi metrics show --path /path/to/table --filter "value > 0.9"
```

### Metrics Analysis

```bash
# Analyze metrics trends
nessi metrics analyze --path /path/to/table --period "last-30-days"

# Detect anomalies in metrics
nessi metrics anomalies --path /path/to/table --sensitivity high

# Forecast metrics
nessi metrics forecast --path /path/to/table --days 7

# Compare metrics between tables
nessi metrics compare --source /path/to/table1 --target /path/to/table2
```

### Metrics Visualization

```bash
# Visualize quality metrics
nessi metrics visualize --path /path/to/table --type quality

# Visualize freshness metrics
nessi metrics visualize --path /path/to/table --type freshness

# Visualize metrics as a heatmap
nessi metrics visualize --path /path/to/table --type heatmap

# Export visualization to file
nessi metrics visualize --path /path/to/table --type quality --output /path/to/visualization.html
```

## Report Generation Examples

### Generating Reports

```bash
# Generate a quality report
nessi report generate --path /path/to/table --type quality --output /path/to/report.html

# Generate a schema report
nessi report generate --path /path/to/table --type schema --output /path/to/report.html

# Generate a freshness report
nessi report generate --path /path/to/table --type freshness --output /path/to/report.html

# Generate a performance report
nessi report generate --path /path/to/table --type performance --output /path/to/report.html
```

### Report Formats

```bash
# Generate HTML report
nessi report generate --path /path/to/table --format html --output /path/to/report.html

# Generate PDF report
nessi report generate --path /path/to/table --format pdf --output /path/to/report.pdf

# Generate JSON report
nessi report generate --path /path/to/table --format json --output /path/to/report.json

# Generate CSV report
nessi report generate --path /path/to/table --format csv --output /path/to/report.csv
```

### Custom Reports

```bash
# Generate report with custom template
nessi report generate --path /path/to/table --template /path/to/custom_template.html --output /path/to/report.html

# Generate report with specific sections
nessi report generate --path /path/to/table --sections summary,quality,schema --output /path/to/report.html

# Generate report with custom title
nessi report generate --path /path/to/table --title "Monthly Quality Report" --output /path/to/report.html
```

### Report Management

```bash
# List generated reports
nessi report list

# View report details
nessi report show --report-id 12345

# Delete a report
nessi report delete --report-id 12345

# Archive reports
nessi report archive --older-than 30d
```

## Integration Examples

### Databricks Integration (Pro Edition)

```bash
# Set Databricks credentials
export DATABRICKS_HOST=your-databricks-instance.cloud.databricks.com
export DATABRICKS_TOKEN=dapi_xxxxxxxxxxxxxxxxxxxxxxxx

# List Databricks catalogs
nessi integration databricks list-catalogs

# List schemas in a catalog
nessi integration databricks list-schemas --catalog your_catalog

# List tables in a schema
nessi integration databricks list-tables --catalog your_catalog --schema your_schema

# Run quality checks on a Databricks table
nessi quality check --integration databricks --catalog your_catalog --schema your_schema --table your_table --rules /path/to/quality_rules.yaml
```

### AWS S3 Integration (Pro Edition)

```bash
# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=us-west-2

# List S3 buckets
nessi integration aws list-buckets

# List objects in a bucket
nessi integration aws list-objects --bucket your_bucket --prefix your_prefix

# Run quality checks on a Delta Lake table in S3
nessi quality check --integration aws --bucket your_bucket --key your_key --rules /path/to/quality_rules.yaml
```

## Advanced Usage Examples

### Sampling

```bash
# Sample by percentage
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --sample-percentage 10

# Sample by row count
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --sample-size 1000

# Stratified sampling
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --sample-strategy stratified --strata-column category
```

### Incremental Processing

```bash
# Process incrementally by timestamp
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --incremental --timestamp-column updated_at

# Process incrementally with lookback
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --incremental --timestamp-column updated_at --lookback 1d
```

### Partitioning

```bash
# Process specific partition
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --partition-by date --partition-value 2023-01-01

# Process multiple partitions
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --partition-by date --partition-values 2023-01-01,2023-01-02,2023-01-03
```

### Filtering

```bash
# Apply filter
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --filter "region='us-west'"

# Apply complex filter
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --filter "region='us-west' AND created_at > '2023-01-01'"
```

### Performance Tuning

```bash
# Set performance preset
nessi --performance-preset high quality check --path /path/to/table --rules /path/to/quality_rules.yaml

# Set specific performance parameters
nessi quality check --path /path/to/table --rules /path/to/quality_rules.yaml --parallel --batch-size 10000 --memory-limit 4G
```

## Workflow Examples

### Data Pipeline Quality Gate

```bash
#!/bin/bash
# quality_gate.sh

# Process data
echo "Processing data..."
process_data.sh /path/to/input /path/to/output

# Run quality checks
echo "Running quality checks..."
nessi quality check --path /path/to/output --rules /path/to/quality_rules.yaml --output-format json --output /tmp/quality_results.json

# Check results
VALID=$(jq -r '.valid' /tmp/quality_results.json)
if [ "$VALID" = "true" ]; then
    echo "Quality checks passed. Publishing data..."
    publish_data.sh /path/to/output
    exit 0
else
    echo "Quality checks failed. See details below:"
    jq -r '.failures[]' /tmp/quality_results.json
    exit 1
fi
```

### Daily Quality Monitoring

```bash
#!/bin/bash
# daily_quality_monitoring.sh

# Set date
DATE=$(date +%Y-%m-%d)
REPORT_DIR="/path/to/reports/$DATE"
mkdir -p "$REPORT_DIR"

# List tables
TABLES=$(nessi tables list --path /path/to/data --output-format json | jq -r '.tables[].path')

# Check each table
for TABLE in $TABLES; do
    TABLE_NAME=$(basename "$TABLE")
    echo "Checking quality for $TABLE_NAME..."
    
    # Run quality checks
    nessi quality check --path "$TABLE" --rules /path/to/quality_rules.yaml
    
    # Generate report
    nessi report generate --path "$TABLE" --type quality --output "$REPORT_DIR/$TABLE_NAME-quality.html"
    
    # Collect metrics
    nessi metrics collect --path "$TABLE" --tags date=$DATE
done

# Generate summary report
nessi report generate-summary --input-dir "$REPORT_DIR" --output "$REPORT_DIR/summary.html"

# Send notification
echo "Daily quality monitoring completed. Reports available at $REPORT_DIR" | mail -s "Quality Monitoring Report $DATE" team@example.com
```

### Schema Evolution Monitoring

```bash
#!/bin/bash
# schema_evolution_monitoring.sh

# Set date
DATE=$(date +%Y-%m-%d)
REPORT_DIR="/path/to/reports/schema/$DATE"
mkdir -p "$REPORT_DIR"

# List tables
TABLES=$(nessi tables list --path /path/to/data --output-format json | jq -r '.tables[].path')

# Check each table
for TABLE in $TABLES; do
    TABLE_NAME=$(basename "$TABLE")
    echo "Checking schema for $TABLE_NAME..."
    
    # Get current schema
    CURRENT_SCHEMA="$REPORT_DIR/$TABLE_NAME-schema-current.json"
    nessi schema show --path "$TABLE" --output-format json --output "$CURRENT_SCHEMA"
    
    # Get previous schema (if exists)
    PREVIOUS_SCHEMA="/path/to/schemas/$TABLE_NAME-schema.json"
    if [ -f "$PREVIOUS_SCHEMA" ]; then
        # Compare schemas
        DIFF_FILE="$REPORT_DIR/$TABLE_NAME-schema-diff.json"
        nessi schema compare --source-file "$PREVIOUS_SCHEMA" --target-file "$CURRENT_SCHEMA" --output-format json --output "$DIFF_FILE"
        
        # Check if schema changed
        CHANGED=$(jq -r '.changed' "$DIFF_FILE")
        if [ "$CHANGED" = "true" ]; then
            echo "Schema changed for $TABLE_NAME"
            # Generate report
            nessi report generate --type schema-diff --source-file "$PREVIOUS_SCHEMA" --target-file "$CURRENT_SCHEMA" --output "$REPORT_DIR/$TABLE_NAME-schema-diff.html"
            # Send notification
            echo "Schema changed for $TABLE_NAME. See report at $REPORT_DIR/$TABLE_NAME-schema-diff.html" | mail -s "Schema Change Alert: $TABLE_NAME" team@example.com
        fi
    fi
    
    # Save current schema as previous for next run
    cp "$CURRENT_SCHEMA" "/path/to/schemas/$TABLE_NAME-schema.json"
done
```

## Scripting Examples

### Python Integration

```python
#!/usr/bin/env python3
# nessi_python_integration.py

import subprocess
import json
import os
import sys
from datetime import datetime

def run_nessi_command(command):
    """Run a Nessi command and return the output as JSON."""
    result = subprocess.run(command, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error running command: {' '.join(command)}")
        print(f"Error: {result.stderr}")
        sys.exit(1)
    return json.loads(result.stdout)

def check_quality(table_path, rules_path):
    """Run quality checks on a table."""
    command = [
        "nessi", "quality", "check",
        "--path", table_path,
        "--rules", rules_path,
        "--output-format", "json"
    ]
    return run_nessi_command(command)

def generate_report(table_path, report_path):
    """Generate a quality report for a table."""
    command = [
        "nessi", "report", "generate",
        "--path", table_path,
        "--type", "quality",
        "--format", "html",
        "--output", report_path
    ]
    subprocess.run(command, check=True)
    return report_path

def main():
    # Configuration
    data_dir = "/path/to/data"
    rules_path = "/path/to/quality_rules.yaml"
    report_dir = f"/path/to/reports/{datetime.now().strftime('%Y-%m-%d')}"
    
    # Create report directory
    os.makedirs(report_dir, exist_ok=True)
    
    # Get list of tables
    tables_command = [
        "nessi", "tables", "list",
        "--path", data_dir,
        "--output-format", "json"
    ]
    tables_result = run_nessi_command(tables_command)
    
    # Process each table
    results = []
    for table in tables_result["tables"]:
        table_path = table["path"]
        table_name = os.path.basename(table_path)
        print(f"Processing table: {table_name}")
        
        # Run quality checks
        quality_result = check_quality(table_path, rules_path)
        
        # Generate report
        report_path = os.path.join(report_dir, f"{table_name}-quality.html")
        generate_report(table_path, report_path)
        
        # Store results
        results.append({
            "table": table_name,
            "path": table_path,
            "valid": quality_result["valid"],
            "score": quality_result.get("score", 0),
            "report": report_path
        })
    
    # Generate summary
    summary = {
        "date": datetime.now().isoformat(),
        "tables_checked": len(results),
        "tables_passed": sum(1 for r in results if r["valid"]),
        "tables_failed": sum(1 for r in results if not r["valid"]),
        "average_score": sum(r["score"] for r in results) / len(results) if results else 0,
        "results": results
    }
    
    # Save summary
    summary_path = os.path.join(report_dir, "summary.json")
    with open(summary_path, "w") as f:
        json.dump(summary, f, indent=2)
    
    print(f"Summary: {summary['tables_passed']}/{summary['tables_checked']} tables passed quality checks")
    print(f"Average quality score: {summary['average_score']:.2f}")
    print(f"Reports saved to: {report_dir}")

if __name__ == "__main__":
    main()
```

### Shell Script for Batch Processing

```bash
#!/bin/bash
# batch_processing.sh

# Configuration
DATA_DIR="/path/to/data"
RULES_DIR="/path/to/rules"
REPORT_DIR="/path/to/reports/$(date +%Y-%m-%d)"
LOG_DIR="/path/to/logs"
LOG_FILE="$LOG_DIR/batch_processing_$(date +%Y%m%d_%H%M%S).log"

# Create directories
mkdir -p "$REPORT_DIR" "$LOG_DIR"

# Start logging
exec > >(tee -a "$LOG_FILE") 2>&1
echo "Starting batch processing at $(date)"

# Function to process a table
process_table() {
    local table_path="$1"
    local table_name=$(basename "$table_path")
    local rules_file="$RULES_DIR/${table_name}_rules.yaml"
    local default_rules="$RULES_DIR/default_rules.yaml"
    
    echo "Processing table: $table_name"
    
    # Use table-specific rules if available, otherwise use default
    if [ -f "$rules_file" ]; then
        echo "Using table-specific rules: $rules_file"
        RULES="$rules_file"
    else
        echo "Using default rules: $default_rules"
        RULES="$default_rules"
    fi
    
    # Run quality checks
    echo "Running quality checks..."
    if nessi quality check --path "$table_path" --rules "$RULES" --output-format json > "$REPORT_DIR/${table_name}_quality.json"; then
        echo "Quality checks completed"
    else
        echo "Error running quality checks for $table_name"
        return 1
    fi
    
    # Generate report
    echo "Generating report..."
    if nessi report generate --path "$table_path" --type quality --output "$REPORT_DIR/${table_name}_quality.html"; then
        echo "Report generated: $REPORT_DIR/${table_name}_quality.html"
    else
        echo "Error generating report for $table_name"
        return 1
    fi
    
    # Collect metrics
    echo "Collecting metrics..."
    if nessi metrics collect --path "$table_path" --tags batch=$(date +%Y%m%d),table=$table_name; then
        echo "Metrics collected"
    else
        echo "Error collecting metrics for $table_name"
        return 1
    fi
    
    return 0
}

# Get list of tables
echo "Listing tables in $DATA_DIR"
TABLES=$(nessi tables list --path "$DATA_DIR" --output-format json | jq -r '.tables[].path')

# Process each table
SUCCESS_COUNT=0
FAILURE_COUNT=0
for TABLE in $TABLES; do
    if process_table "$TABLE"; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        FAILURE_COUNT=$((FAILURE_COUNT + 1))
    fi
    echo "-----------------------------------"
done

# Generate summary report
echo "Generating summary report..."
nessi report generate-summary --input-dir "$REPORT_DIR" --output "$REPORT_DIR/summary.html"

# Print summary
echo "Batch processing completed at $(date)"
echo "Tables processed successfully: $SUCCESS_COUNT"
echo "Tables with errors: $FAILURE_COUNT"
echo "Reports saved to: $REPORT_DIR"
echo "Log file: $LOG_FILE"

# Exit with error if any table failed
if [ "$FAILURE_COUNT" -gt 0 ]; then
    exit 1
fi

exit 0
```

---

For more examples and detailed information, please visit [nessi.dev](https://nessi.dev) or contact support@nessi.dev.
