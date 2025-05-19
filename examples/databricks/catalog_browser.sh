#!/bin/bash

# Databricks Unity Catalog Browser Example
# This script demonstrates how to browse Databricks Unity Catalog using Nessi

# Ensure environment variables are set
if [ -z "$DATABRICKS_HOST" ] || [ -z "$DATABRICKS_TOKEN" ]; then
  echo "Error: DATABRICKS_HOST and DATABRICKS_TOKEN must be set"
  echo "Example:"
  echo "export DATABRICKS_HOST=\"https://your-workspace.cloud.databricks.com\""
  echo "export DATABRICKS_TOKEN=\"your-personal-access-token\""
  exit 1
fi

# Set default catalog if not specified
DATABRICKS_DEFAULT_CATALOG=${DATABRICKS_DEFAULT_CATALOG:-"main"}

# List all catalogs
echo "\n=== Listing all catalogs ==="
nessi catalog list --type databricks

# List all schemas in the default catalog
echo "\n=== Listing schemas in catalog: $DATABRICKS_DEFAULT_CATALOG ==="
nessi catalog schemas --type databricks --catalog $DATABRICKS_DEFAULT_CATALOG

# List tables in the default schema
DEFAULT_SCHEMA=${DATABRICKS_DEFAULT_SCHEMA:-"default"}
echo "\n=== Listing tables in schema: $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA ==="
nessi catalog tables --type databricks --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA

# Get details of a specific table (if available)
echo "\n=== Enter a table name to view details (or press Enter to skip) ==="
read -p "Table name: " TABLE_NAME

if [ ! -z "$TABLE_NAME" ]; then
  echo "\n=== Table details for: $TABLE_NAME ==="
  nessi catalog describe --type databricks --database $DATABRICKS_DEFAULT_CATALOG.$DEFAULT_SCHEMA --table $TABLE_NAME
fi

echo "\nDatabricks Unity Catalog browsing complete!"
