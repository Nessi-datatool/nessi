# Databricks Integration Examples

This directory contains examples of using Nessi with Databricks.

## Prerequisites

Before running these examples, make sure you have:

1. A Databricks workspace with access to Unity Catalog
2. A personal access token for authentication
3. Nessi CLI installed

## Environment Setup

Set up your environment variables:

```bash
export DATABRICKS_HOST="https://your-workspace.cloud.databricks.com"
export DATABRICKS_TOKEN="your-personal-access-token"
export DATABRICKS_WORKSPACE_ID="your-workspace-id"  # Optional, defaults to "0"
export DATABRICKS_DEFAULT_SCHEMA="default"          # Optional, defaults to "default"
export DATABRICKS_DEFAULT_CATALOG="main"            # Optional, defaults to "hive_metastore"
```

## Examples

### 1. Browsing Unity Catalog

The `catalog_browser.sh` script demonstrates how to browse the Databricks Unity Catalog using Nessi.

```bash
./catalog_browser.sh
```

### 2. Working with Delta Lake Tables

The `delta_operations.sh` script shows how to work with Delta Lake tables, including reading data, performing time travel, and writing data.

```bash
./delta_operations.sh
```

### 3. Data Quality Checks

The `quality_checks.sh` script demonstrates how to run data quality checks on Databricks tables.

```bash
./quality_checks.sh
```

## Troubleshooting

If you encounter issues:

1. Verify your environment variables are set correctly
2. Ensure your Databricks token has the necessary permissions
3. Check network connectivity to your Databricks workspace
