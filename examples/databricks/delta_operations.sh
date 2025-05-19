#!/bin/bash

# Databricks Delta Lake Operations Example
# This script demonstrates how to work with Delta Lake tables using Nessi

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

# Read data from the Delta Lake table
echo "\n=== Reading data from Delta Lake table: $TABLE_NAME ==="
nessi datalake read --format delta --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME --limit 10

# Get table metadata
echo "\n=== Getting table metadata ==="
nessi datalake metadata --format delta --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME

# Perform time travel (if the table has multiple versions)
echo "\n=== Performing time travel (reading previous version) ==="
nessi datalake timetravel --format delta --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME \
  --version 0 --limit 5

# Get version history
echo "\n=== Getting version history ==="
nessi datalake history --format delta --catalog databricks \
  --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME

echo "\nDelta Lake operations complete!"
