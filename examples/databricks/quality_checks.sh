#!/bin/bash

# Databricks Data Quality Checks Example
# This script demonstrates how to run data quality checks on Databricks tables using Nessi

# Ensure environment variables are set
if [ -z "$DATABRICKS_HOST" ] || [ -z "$DATABRICKS_TOKEN" ]; then
  echo "Error: DATABRICKS_HOST and DATABRICKS_TOKEN must be set"
  echo "Example:"
  echo "export DATABRICKS_HOST=\"https://your-workspace.cloud.databricks.com\""
  echo "export DATABRICKS_TOKEN=\"your-personal-access-token\""
  exit 1
fi

# Set default catalog and schema if not specified
DATABRICKS_DEFAULT_CATALOG=${DATABRICKS_DEFAULT_CATALOG:-"main"}
DEFAULT_SCHEMA=${DATABRICKS_DEFAULT_SCHEMA:-"default"}

# Prompt for table name
echo "Enter a Delta Lake table name from your Databricks catalog:"
read -p "Table name: " TABLE_NAME

if [ -z "$TABLE_NAME" ]; then
  echo "Error: Table name is required"
  exit 1
fi

# Run basic data quality profile
echo "\n=== Running data quality profile on: $TABLE_NAME ==="
nessi quality profile --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME

# Run data quality checks
echo "\n=== Running data quality checks ==="
nessi quality check --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME

# Generate quality report
echo "\n=== Generating quality report ==="
nessi quality report --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME \
  --output-format json --output-file "${TABLE_NAME}_quality_report.json"

echo "\nData quality checks complete!"
echo "Quality report saved to: ${TABLE_NAME}_quality_report.json"
