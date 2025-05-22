# Nessi CLI Reference Guide

This reference guide provides a quick overview of all available Nessi CLI commands. For detailed documentation on each command, please refer to the linked documentation files.

## Global Options

The following options are available for all Nessi commands:

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help for a command |
| `--verbose`, `-v` | Enable verbose output |
| `--quiet`, `-q` | Suppress all output except errors |
| `--json` | Output in JSON format |
| `--config` | Path to configuration file |
| `--log-level` | Set log level (debug, info, warn, error) |
| `--no-color` | Disable colored output |
| `--license` | Path to license file |

## Basic Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi --help` | Show help information | [index.md](index.md) |
| `nessi version` | Show version information | [index.md](index.md) |
| `nessi info` | Show environment information | [index.md](index.md) |
| `nessi completion [shell]` | Generate shell completion scripts | [index.md](index.md) |
| `nessi license [command]` | Manage license information | [index.md](index.md) |

## Table Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi tables list` | List all available Delta Lake tables | [tables.md](tables.md#tables-list) |
| `nessi tables describe <table-name>` | Show detailed information about a table | [tables.md](tables.md#tables-describe) |
| `nessi tables schema <table-name>` | Show the schema of a table | [tables.md](tables.md#tables-schema) |
| `nessi tables history <table-name>` | Show the version history of a table | [tables.md](tables.md#tables-history) |
| `nessi tables profile <table-name>` | Generate a profile of a table's data | [tables.md](tables.md#tables-profile) |
| `nessi tables diff <table1> <table2>` | Compare two versions of a table or two different tables | [tables.md](tables.md#tables-diff) |
| `nessi tables export <table-name> <output-path>` | Export a table to various formats | [tables.md](tables.md#tables-export) |

## Quality Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi quality check <table-name>` | Run quality checks on a table | [quality.md](quality.md#quality-check) |
| `nessi quality rules list` | List all available quality rules | [quality.md](quality.md#quality-rules-list) |
| `nessi quality rules create` | Create a new quality rule | [quality.md](quality.md#quality-rules-create) |
| `nessi quality rules delete <rule-name>` | Delete a quality rule | [quality.md](quality.md#quality-rules-delete) |
| `nessi quality report <table-name>` | Generate a comprehensive quality report | [quality.md](quality.md#quality-report) |
| `nessi quality metrics <table-name>` | Retrieve and analyze quality metrics | [quality.md](quality.md#quality-metrics) |

## Scan Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi scan run <source>` | Run a data scan on a table or file | [scan.md](scan.md#scan-run) |
| `nessi scan compare <source1> <source2>` | Compare scan results between two sources | [scan.md](scan.md#scan-compare) |
| `nessi scan history <source>` | Show the history of scans for a specific source | [scan.md](scan.md#scan-history) |
| `nessi scan schedule [command]` | Manage scheduled scans | [scan.md](scan.md#scan-schedule) |
| `nessi scan export <source> <output-path>` | Export scan results to various formats | [scan.md](scan.md#scan-export) |

## Schema Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi schema show <source>` | Show the schema of a table or file | [schema.md](schema.md#schema-show) |
| `nessi schema validate <source> [schema-file]` | Validate data against a schema | [schema.md](schema.md#schema-validate) |
| `nessi schema diff <source1> <source2>` | Compare two schemas | [schema.md](schema.md#schema-diff) |
| `nessi schema evolve <source>` | Evolve a schema with new columns or changes | [schema.md](schema.md#schema-evolve) |
| `nessi schema tree <source>` | Visualize schema as a tree structure | [schema.md](schema.md#schema-tree) |
| `nessi schema infer <source>` | Infer schema from data | [schema.md](schema.md#schema-infer) |
| `nessi schema registry [command]` | Manage schemas in a schema registry | [schema.md](schema.md#schema-registry) |

## Metrics Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi metrics collect <source>` | Collect metrics from a table or file | [metrics.md](metrics.md#metrics-collect) |
| `nessi metrics list` | List stored metrics | [metrics.md](metrics.md#metrics-list) |
| `nessi metrics show <metrics-id>` | Show detailed information about specific metrics | [metrics.md](metrics.md#metrics-show) |
| `nessi metrics compare <metrics-id1> <metrics-id2>` | Compare metrics between different sources or time periods | [metrics.md](metrics.md#metrics-compare) |
| `nessi metrics trend <source>` | Analyze trends in metrics over time | [metrics.md](metrics.md#metrics-trend) |
| `nessi metrics alert [command]` | Manage metric-based alerts | [metrics.md](metrics.md#metrics-alert) |
| `nessi metrics visualize <source>` | Visualize metrics data | [metrics.md](metrics.md#metrics-visualize) |
| `nessi metrics export <source> <output-path>` | Export metrics data to various formats | [metrics.md](metrics.md#metrics-export) |
| `nessi metrics grafana [command]` | Grafana integration commands | [metrics.md](metrics.md#grafana-integration) |

## Report Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi report generate <source>` | Generate a report from scan results, quality checks, or metrics | [report.md](report.md#report-generate) |
| `nessi report template list` | List available report templates | [report.md](report.md#report-template-list) |
| `nessi report template create` | Create a new report template | [report.md](report.md#report-template-create) |
| `nessi report template export <template-name> <output-path>` | Export a report template | [report.md](report.md#report-template-export) |
| `nessi report template import <template-path>` | Import a report template | [report.md](report.md#report-template-import) |
| `nessi report schedule [command]` | Manage scheduled reports | [report.md](report.md#report-schedule) |
| `nessi report history` | Show the history of generated reports | [report.md](report.md#report-history) |

## Integration Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi integration databricks [command]` | Databricks integration commands | [integration.md](integration.md#integration-databricks) |
| `nessi integration aws [command]` | AWS integration commands | [integration.md](integration.md#integration-aws) |
| `nessi integration grafana [command]` | Grafana integration commands | [integration.md](integration.md#integration-grafana) |
| `nessi integration kafka [command]` | Apache Kafka integration commands | [integration.md](integration.md#integration-kafka) |
| `nessi integration snowflake [command]` | Snowflake integration commands | [integration.md](integration.md#integration-snowflake) |
| `nessi integration list` | List all configured integrations | [integration.md](integration.md#integration-list) |
| `nessi integration test <type> [name]` | Test an integration connection | [integration.md](integration.md#integration-test) |
| `nessi integration export <type> <name> <output-path>` | Export integration configuration | [integration.md](integration.md#integration-export) |
| `nessi integration import <config-path>` | Import integration configuration | [integration.md](integration.md#integration-import) |

## Configuration Commands

| Command | Description | Documentation |
|---------|-------------|---------------|
| `nessi config set <key> <value>` | Set a configuration value | [config.md](config.md#config-set) |
| `nessi config get <key>` | Get a configuration value | [config.md](config.md#config-get) |
| `nessi config list` | List all configuration values | [config.md](config.md#config-list) |
| `nessi config unset <key>` | Unset a configuration value | [config.md](config.md#config-unset) |
| `nessi config init` | Initialize a new configuration file | [config.md](config.md#config-init) |
| `nessi config profile [command]` | Manage configuration profiles | [config.md](config.md#config-profile) |
| `nessi config import <file-path>` | Import configuration from a file | [config.md](config.md#config-import) |
| `nessi config export <file-path>` | Export configuration to a file | [config.md](config.md#config-export) |
| `nessi config validate` | Validate configuration | [config.md](config.md#config-validate) |
| `nessi config env [command]` | Manage environment variables for configuration | [config.md](config.md#config-env) |

## Common Workflows

### Basic Data Quality Check

```bash
# List available tables
nessi tables list

# Check the schema of a table
nessi schema show /data/customers

# Run quality checks
nessi quality check /data/customers

# Generate a quality report
nessi quality report /data/customers --output html --file quality_report.html
```

### Comprehensive Data Scan

```bash
# Run a full scan
nessi scan run /data/customers --type all

# Export scan results
nessi scan export /results/scan_20250522.json /reports/scan_report.html --format html

# Schedule regular scans
nessi scan schedule create --source /data/customers --type all --schedule "0 0 * * *" --output html --file "/reports/customers_scan_%Y%m%d.html"
```

### Metrics Collection and Analysis

```bash
# Collect metrics
nessi metrics collect /data/customers --store

# Analyze trends
nessi metrics trend /data/customers --from 2025-04-01 --to 2025-05-22

# Visualize metrics
nessi metrics visualize /data/customers --type radar --metric quality --output html --file quality_radar.html

# Create a Grafana dashboard
nessi metrics grafana create-dashboard --source /data/customers --type quality --name "Data Quality Dashboard"
```

### Schema Management

```bash
# Show schema
nessi schema show /data/customers

# Validate against a schema
nessi schema validate /data/customers /schemas/customers_schema.json

# Compare schemas
nessi schema diff /data/customers /data/customers_backup

# Evolve schema
nessi schema evolve /data/customers --add-column "subscription_tier:string:true:Customer's subscription tier" --apply
```

### Integration with External Systems

```bash
# Configure Databricks connection
nessi integration databricks connect --host https://your-workspace.cloud.databricks.com --token your-token

# Import a table from Databricks
nessi integration databricks import --table main.default.customers --path /data/customers

# Export a table to S3
nessi integration aws s3-export --path /data/customers --bucket my-data-bucket --key data/customers/customer_data.parquet
```

## Feature Availability

| Feature | Community Edition | Pro Edition |
|---------|-------------------|-------------|
| Basic Table Operations | ✓ | ✓ |
| Basic Schema Operations | ✓ | ✓ |
| Basic Quality Checks | ✓ | ✓ |
| Basic Data Scanning | ✓ | ✓ |
| Basic Metrics Collection | ✓ | ✓ |
| Simple Reports | ✓ | ✓ |
| Advanced Schema Operations | | ✓ |
| Advanced Quality Checks | | ✓ |
| Comprehensive Metrics | | ✓ |
| Enhanced Reports | | ✓ |
| Grafana Integration | | ✓ |
| Workflow Orchestration | | ✓ |
| Databricks Integration | | ✓ |
| AWS S3 Storage | | ✓ |

For more information on licensing, see the [License Management documentation](../LICENSE_MANAGEMENT.md).

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NESSI_CONFIG` | Path to configuration file |
| `NESSI_LOG_LEVEL` | Log level (debug, info, warn, error) |
| `NESSI_LICENSE_PATH` | Path to license file |
| `NESSI_NO_COLOR` | Disable colored output if set to true |
| `NESSI_DEFAULT_FORMAT` | Default output format |
| `NESSI_METRICS_STORE` | Path or URL to metrics store |
| `NESSI_SCHEMA_REGISTRY` | URL to schema registry |
| `NESSI_DATABRICKS_TOKEN` | Databricks access token |
| `NESSI_DATABRICKS_HOST` | Databricks host URL |
| `NESSI_AWS_REGION` | AWS region for S3 access |

## Error Codes

| Error Code Range | Component |
|------------------|-----------|
| N500-N599 | Table Operations |
| N600-N699 | Scan Operations |
| N700-N799 | Schema Operations |
| N800-N899 | Quality Operations |
| N900-N999 | Metrics Operations |
| N1000-N1099 | Report Operations |
| N1100-N1199 | Integration Operations |
| N1200-N1299 | Configuration Operations |

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).

## Getting Help

If you need help with Nessi CLI, you can:

- Use the `--help` flag with any command to see usage information
- Visit the [Nessi website](https://nessi.dev) for documentation and tutorials
- Join the [Nessi community forum](https://community.nessi.dev) for support and discussions
- Report issues on [GitHub](https://github.com/nessi-dev/nessi/issues)
