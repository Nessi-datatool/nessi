# Integration Commands

This document provides detailed information about Nessi's integration commands for connecting with external systems and services.

## `integration databricks`

Commands for Databricks integration.

### Usage

```bash
nessi integration databricks [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `connect` | Configure connection to Databricks |
| `list` | List tables in Databricks catalog |
| `import` | Import a table from Databricks |
| `export` | Export a table to Databricks |
| `execute` | Execute a query on Databricks |
| `status` | Check Databricks connection status |

### Options for `connect`

| Option | Description |
|--------|-------------|
| `--host` | Databricks host URL |
| `--token` | Databricks access token |
| `--workspace-id` | Databricks workspace ID |
| `--catalog` | Default catalog name |
| `--schema` | Default schema name |
| `--config` | Path to configuration file |
| `--profile` | Configuration profile name |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--catalog` | Catalog name |
| `--schema` | Schema name |
| `--filter` | Filter pattern for table names |
| `--output` | Output format (text, json, csv) |
| `--include-views` | Include views in the results |
| `--include-details` | Include table details |

### Options for `import`

| Option | Description |
|--------|-------------|
| `--table` | Databricks table name (catalog.schema.table) |
| `--path` | Local path to save the table |
| `--format` | Import format (delta, parquet, csv) |
| `--query` | SQL query to execute instead of table name |
| `--sample` | Sampling percentage (0-100) |
| `--columns` | Comma-separated list of columns to import |
| `--where` | Filter condition for rows to import |

### Options for `export`

| Option | Description |
|--------|-------------|
| `--path` | Local path of the table or file to export |
| `--table` | Databricks table name (catalog.schema.table) |
| `--format` | Export format (delta, parquet, csv) |
| `--mode` | Write mode (error, overwrite, append, ignore) |
| `--partition-by` | Comma-separated list of partition columns |
| `--options` | Additional options in key=value format |

### Options for `execute`

| Option | Description |
|--------|-------------|
| `--query` | SQL query to execute |
| `--file` | Path to SQL file |
| `--output` | Output format (text, json, csv) |
| `--result-path` | Path to save query results |
| `--async` | Execute query asynchronously |
| `--timeout` | Query timeout in seconds |

### Options for `status`

| Option | Description |
|--------|-------------|
| `--verbose` | Show detailed status information |
| `--output` | Output format (text, json) |
| `--test-query` | Execute a test query to verify connection |

### Examples

```bash
# Configure Databricks connection
nessi integration databricks connect --host https://your-workspace.cloud.databricks.com --token your-token

# List tables in Databricks catalog
nessi integration databricks list --catalog main --schema default

# Import a table from Databricks
nessi integration databricks import --table main.default.customers --path /data/customers

# Import using a query
nessi integration databricks import --query "SELECT * FROM main.default.customers WHERE region = 'US'" --path /data/us_customers

# Export a table to Databricks
nessi integration databricks export --path /data/customers --table main.default.customers_copy --mode overwrite

# Execute a query on Databricks
nessi integration databricks execute --query "SELECT COUNT(*) FROM main.default.customers" --output json

# Check connection status
nessi integration databricks status --verbose
```

### Output for `list`

```
Databricks Tables:
Catalog: main
Schema: default

NAME                TYPE    ROWS        SIZE       LAST MODIFIED
customers           TABLE   1,000,000   1.2 GB     2025-05-01 12:34:56
orders              TABLE   5,000,000   3.5 GB     2025-05-10 09:23:45
products            TABLE   50,000      500 MB     2025-04-15 14:56:23
customer_view       VIEW    -           -          2025-05-05 10:11:12
```

### Output for `status`

```
Databricks Connection Status:

Connection: ACTIVE
Host: https://your-workspace.cloud.databricks.com
Workspace ID: 1234567890
Default Catalog: main
Default Schema: default
API Version: 2.0
Client Version: 1.2.3
Server Version: 14.3
Cluster Status: RUNNING
Authentication: TOKEN (valid)
Permissions: READ, WRITE, EXECUTE

Test Query: SUCCESS (executed in 0.35s)
```

## `integration aws`

Commands for AWS integration.

### Usage

```bash
nessi integration aws [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `configure` | Configure AWS credentials and settings |
| `s3-list` | List objects in an S3 bucket |
| `s3-import` | Import data from S3 |
| `s3-export` | Export data to S3 |
| `s3-sync` | Sync data between local and S3 |
| `status` | Check AWS connection status |

### Options for `configure`

| Option | Description |
|--------|-------------|
| `--access-key` | AWS access key ID |
| `--secret-key` | AWS secret access key |
| `--region` | AWS region |
| `--profile` | AWS profile name |
| `--config` | Path to configuration file |
| `--role-arn` | AWS role ARN for assuming role |
| `--session-token` | AWS session token |

### Options for `s3-list`

| Option | Description |
|--------|-------------|
| `--bucket` | S3 bucket name |
| `--prefix` | Object prefix (folder path) |
| `--recursive` | List objects recursively |
| `--output` | Output format (text, json, csv) |
| `--filter` | Filter pattern for object names |
| `--include-details` | Include object details |

### Options for `s3-import`

| Option | Description |
|--------|-------------|
| `--bucket` | S3 bucket name |
| `--key` | S3 object key |
| `--path` | Local path to save the data |
| `--format` | Import format (delta, parquet, csv, auto) |
| `--recursive` | Import objects recursively |
| `--include` | Pattern to include objects |
| `--exclude` | Pattern to exclude objects |

### Options for `s3-export`

| Option | Description |
|--------|-------------|
| `--path` | Local path of the data to export |
| `--bucket` | S3 bucket name |
| `--key` | S3 object key |
| `--format` | Export format (delta, parquet, csv) |
| `--acl` | S3 object ACL |
| `--storage-class` | S3 storage class |
| `--encryption` | Server-side encryption method |
| `--metadata` | Object metadata in key=value format |

### Options for `s3-sync`

| Option | Description |
|--------|-------------|
| `--source` | Source path (local or s3://) |
| `--destination` | Destination path (local or s3://) |
| `--delete` | Delete files that exist in the destination but not in the source |
| `--include` | Pattern to include files |
| `--exclude` | Pattern to exclude files |
| `--acl` | S3 object ACL |
| `--storage-class` | S3 storage class |

### Options for `status`

| Option | Description |
|--------|-------------|
| `--verbose` | Show detailed status information |
| `--output` | Output format (text, json) |
| `--test` | Test connection by listing buckets |

### Examples

```bash
# Configure AWS credentials
nessi integration aws configure --access-key AKIAIOSFODNN7EXAMPLE --secret-key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY --region us-west-2

# List objects in an S3 bucket
nessi integration aws s3-list --bucket my-data-bucket --prefix data/customers/

# Import data from S3
nessi integration aws s3-import --bucket my-data-bucket --key data/customers/customer_data.parquet --path /data/customers

# Import recursively from S3
nessi integration aws s3-import --bucket my-data-bucket --prefix data/customers/ --path /data/customers --recursive

# Export data to S3
nessi integration aws s3-export --path /data/customers --bucket my-data-bucket --key data/customers/customer_data.parquet

# Sync data between local and S3
nessi integration aws s3-sync --source /data/customers --destination s3://my-data-bucket/data/customers/

# Check AWS connection status
nessi integration aws status --verbose
```

### Output for `s3-list`

```
S3 Objects:
Bucket: my-data-bucket
Prefix: data/customers/

KEY                                     SIZE        LAST MODIFIED           STORAGE CLASS
data/customers/customer_data.parquet    1.2 GB      2025-05-01 12:34:56     STANDARD
data/customers/orders.parquet           3.5 GB      2025-05-10 09:23:45     STANDARD
data/customers/products.parquet         500 MB      2025-04-15 14:56:23     STANDARD
```

### Output for `status`

```
AWS Connection Status:

Connection: ACTIVE
Region: us-west-2
Access Key ID: AKIA************MPLE
Account ID: 123456789012
IAM User: nessi-service
Available Buckets: 3
Permissions: s3:ListBucket, s3:GetObject, s3:PutObject

Test: SUCCESS (listed buckets in 0.25s)
```

## `integration grafana`

Commands for Grafana integration.

### Usage

```bash
nessi integration grafana [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `configure` | Configure Grafana connection |
| `dashboards` | List, create, update, or delete dashboards |
| `datasources` | List, create, update, or delete data sources |
| `alerts` | List, create, update, or delete alerts |
| `panels` | Create or update dashboard panels |
| `status` | Check Grafana connection status |

### Options for `configure`

| Option | Description |
|--------|-------------|
| `--url` | Grafana URL |
| `--api-key` | Grafana API key |
| `--username` | Grafana username (for basic auth) |
| `--password` | Grafana password (for basic auth) |
| `--config` | Path to configuration file |
| `--profile` | Configuration profile name |
| `--org-id` | Grafana organization ID |

### Options for `dashboards list`

| Option | Description |
|--------|-------------|
| `--folder` | Filter by folder |
| `--tag` | Filter by tag |
| `--starred` | Show only starred dashboards |
| `--output` | Output format (text, json, csv) |
| `--include-details` | Include dashboard details |

### Options for `dashboards create`

| Option | Description |
|--------|-------------|
| `--title` | Dashboard title |
| `--folder` | Dashboard folder |
| `--template` | Path to dashboard template JSON file |
| `--source` | Data source for metrics (path or metrics ID) |
| `--metrics` | Metrics to include (quality, performance, freshness, all) |
| `--time-range` | Default time range (e.g., "now-30d to now") |
| `--variables` | Dashboard variables in key=value format |
| `--tags` | Comma-separated list of tags |

### Options for `panels create`

| Option | Description |
|--------|-------------|
| `--dashboard` | Dashboard UID or ID |
| `--title` | Panel title |
| `--type` | Panel type (graph, gauge, stat, table, heatmap) |
| `--metric` | Metric to visualize |
| `--source` | Data source for metrics (path or metrics ID) |
| `--position` | Panel position (x,y,w,h) |
| `--description` | Panel description |
| `--thresholds` | Panel thresholds in key=value format |

### Examples

```bash
# Configure Grafana connection
nessi integration grafana configure --url http://grafana:3000 --api-key eyJrIjoiT0tTcG1pUlY2RnVKZTFVaDFsNFZXdE9ZWmNrMkZYbk

# List dashboards
nessi integration grafana dashboards list

# Create a dashboard
nessi integration grafana dashboards create --title "Data Quality Dashboard" --folder "Nessi" --source /data/customers --metrics quality,performance

# Create a panel
nessi integration grafana panels create --dashboard abc123 --title "Data Completeness" --type gauge --metric quality.completeness --source /data/customers

# Create an alert
nessi integration grafana alerts create --name "Low Completeness Alert" --metric quality.completeness --source /data/customers --condition "<" --threshold 0.95 --severity high

# Check connection status
nessi integration grafana status
```

### Output for `dashboards list`

```
Grafana Dashboards:

UID     TITLE                   FOLDER    TAGS                  LAST MODIFIED
abc123  Data Quality Dashboard  Nessi     quality,nessi         2025-05-22 15:30:45
def456  Performance Metrics     Nessi     performance,nessi     2025-05-20 10:15:30
ghi789  Freshness Monitoring    Nessi     freshness,nessi       2025-05-18 09:45:15
```

### Output for `status`

```
Grafana Connection Status:

Connection: ACTIVE
URL: http://grafana:3000
Version: 10.2.3
Organization: Nessi
User: api_key
Permissions: Admin
Datasources: 3
Dashboards: 5

Test: SUCCESS (API responded in 0.15s)
```

## `integration kafka`

Commands for Apache Kafka integration.

### Usage

```bash
nessi integration kafka [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `configure` | Configure Kafka connection |
| `topics` | List, create, or delete topics |
| `produce` | Produce messages to a topic |
| `consume` | Consume messages from a topic |
| `schema-registry` | Interact with Schema Registry |
| `status` | Check Kafka connection status |

### Options for `configure`

| Option | Description |
|--------|-------------|
| `--bootstrap-servers` | Kafka bootstrap servers |
| `--security-protocol` | Security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL) |
| `--sasl-mechanism` | SASL mechanism (PLAIN, SCRAM-SHA-256, SCRAM-SHA-512) |
| `--sasl-username` | SASL username |
| `--sasl-password` | SASL password |
| `--ssl-ca-cert` | SSL CA certificate path |
| `--ssl-client-cert` | SSL client certificate path |
| `--ssl-client-key` | SSL client key path |
| `--config` | Path to configuration file |
| `--profile` | Configuration profile name |

### Options for `topics list`

| Option | Description |
|--------|-------------|
| `--filter` | Filter pattern for topic names |
| `--output` | Output format (text, json, csv) |
| `--include-details` | Include topic details |
| `--include-configs` | Include topic configurations |

### Options for `produce`

| Option | Description |
|--------|-------------|
| `--topic` | Topic name |
| `--message` | Message content |
| `--file` | Path to file containing messages |
| `--key` | Message key |
| `--headers` | Message headers in key=value format |
| `--format` | Message format (json, avro, string) |
| `--schema` | Schema for Avro messages |
| `--batch-size` | Number of messages to send in a batch |

### Options for `consume`

| Option | Description |
|--------|-------------|
| `--topic` | Topic name |
| `--group` | Consumer group ID |
| `--from` | Offset to start consuming from (earliest, latest, or specific offset) |
| `--limit` | Maximum number of messages to consume |
| `--timeout` | Consume timeout in seconds |
| `--format` | Message format (json, avro, string) |
| `--output` | Output format (text, json, csv) |
| `--file` | Path to save consumed messages |

### Examples

```bash
# Configure Kafka connection
nessi integration kafka configure --bootstrap-servers kafka1:9092,kafka2:9092 --security-protocol SASL_SSL --sasl-mechanism PLAIN --sasl-username admin --sasl-password admin-secret

# List topics
nessi integration kafka topics list

# Create a topic
nessi integration kafka topics create --topic nessi-metrics --partitions 3 --replication-factor 2

# Produce messages
nessi integration kafka produce --topic nessi-metrics --file /data/metrics.json --format json

# Consume messages
nessi integration kafka consume --topic nessi-metrics --group nessi-consumer --from earliest --limit 100 --output json

# Check connection status
nessi integration kafka status
```

### Output for `topics list`

```
Kafka Topics:

NAME                PARTITIONS  REPLICATION  RETENTION     MESSAGES
nessi-metrics       3           2            7 days        15,230
nessi-alerts        3           2            7 days        1,245
nessi-reports       3           2            7 days        523
```

### Output for `status`

```
Kafka Connection Status:

Connection: ACTIVE
Bootstrap Servers: kafka1:9092,kafka2:9092
Security Protocol: SASL_SSL
SASL Mechanism: PLAIN
Cluster ID: XYZ-123456
Broker Count: 3
Controller: kafka1:9092
Topics: 15

Test: SUCCESS (connected in 0.22s)
```

## `integration snowflake`

Commands for Snowflake integration.

### Usage

```bash
nessi integration snowflake [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `connect` | Configure connection to Snowflake |
| `list` | List tables in Snowflake |
| `import` | Import a table from Snowflake |
| `export` | Export a table to Snowflake |
| `execute` | Execute a query on Snowflake |
| `status` | Check Snowflake connection status |

### Options for `connect`

| Option | Description |
|--------|-------------|
| `--account` | Snowflake account identifier |
| `--user` | Snowflake username |
| `--password` | Snowflake password |
| `--warehouse` | Snowflake warehouse |
| `--database` | Default database |
| `--schema` | Default schema |
| `--role` | Snowflake role |
| `--private-key` | Path to private key file |
| `--config` | Path to configuration file |
| `--profile` | Configuration profile name |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--database` | Database name |
| `--schema` | Schema name |
| `--filter` | Filter pattern for table names |
| `--output` | Output format (text, json, csv) |
| `--include-views` | Include views in the results |
| `--include-details` | Include table details |

### Options for `import`

| Option | Description |
|--------|-------------|
| `--table` | Snowflake table name (database.schema.table) |
| `--path` | Local path to save the table |
| `--format` | Import format (parquet, csv) |
| `--query` | SQL query to execute instead of table name |
| `--sample` | Sampling percentage (0-100) |
| `--columns` | Comma-separated list of columns to import |
| `--where` | Filter condition for rows to import |

### Options for `export`

| Option | Description |
|--------|-------------|
| `--path` | Local path of the table or file to export |
| `--table` | Snowflake table name (database.schema.table) |
| `--format` | Export format (parquet, csv) |
| `--mode` | Write mode (error, overwrite, append, ignore) |
| `--create-table` | Create table if it doesn't exist |
| `--stage` | Snowflake stage name |
| `--options` | Additional options in key=value format |

### Examples

```bash
# Configure Snowflake connection
nessi integration snowflake connect --account xy12345.us-east-1 --user admin --password admin-secret --warehouse compute_wh --database analytics --schema public

# List tables in Snowflake
nessi integration snowflake list --database analytics --schema public

# Import a table from Snowflake
nessi integration snowflake import --table analytics.public.customers --path /data/customers

# Import using a query
nessi integration snowflake import --query "SELECT * FROM analytics.public.customers WHERE region = 'US'" --path /data/us_customers

# Export a table to Snowflake
nessi integration snowflake export --path /data/customers --table analytics.public.customers_copy --mode overwrite

# Execute a query on Snowflake
nessi integration snowflake execute --query "SELECT COUNT(*) FROM analytics.public.customers" --output json

# Check connection status
nessi integration snowflake status --verbose
```

### Output for `list`

```
Snowflake Tables:
Database: analytics
Schema: public

NAME                TYPE    ROWS        SIZE       LAST MODIFIED
customers           TABLE   1,000,000   1.2 GB     2025-05-01 12:34:56
orders              TABLE   5,000,000   3.5 GB     2025-05-10 09:23:45
products            TABLE   50,000      500 MB     2025-04-15 14:56:23
customer_view       VIEW    -           -          2025-05-05 10:11:12
```

### Output for `status`

```
Snowflake Connection Status:

Connection: ACTIVE
Account: xy12345.us-east-1
User: admin
Warehouse: compute_wh
Database: analytics
Schema: public
Role: ACCOUNTADMIN
Version: 7.12.0
Edition: Enterprise
Region: US East (N. Virginia)

Test Query: SUCCESS (executed in 0.42s)
```

## `integration list`

Lists all configured integrations.

### Usage

```bash
nessi integration list [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--type` | Filter by integration type |
| `--status` | Filter by status (active, inactive) |
| `--output` | Output format (text, json, csv) |
| `--detailed` | Show detailed integration information |

### Examples

```bash
# List all integrations
nessi integration list

# List active integrations
nessi integration list --status active

# List specific integration type
nessi integration list --type databricks

# Show detailed information
nessi integration list --detailed
```

### Output

#### Text Format (Default)

```
Configured Integrations:

TYPE        NAME            STATUS    LAST USED
databricks  prod-workspace  active    2025-05-22 15:30:45
aws         data-lake       active    2025-05-20 10:15:30
grafana     monitoring      active    2025-05-18 09:45:15
kafka       events          inactive  2025-05-01 14:20:10
snowflake   analytics       active    2025-05-15 11:30:25
```

#### Detailed Text Format

```
Integration: databricks (prod-workspace)
  Type: databricks
  Status: active
  Last Used: 2025-05-22 15:30:45
  Host: https://your-workspace.cloud.databricks.com
  Workspace ID: 1234567890
  Default Catalog: main
  Default Schema: default
  API Version: 2.0
  Client Version: 1.2.3
  Configuration File: /Users/username/.nessi/integrations/databricks-prod.yaml
```

## `integration test`

Tests an integration connection.

### Usage

```bash
nessi integration test <type> [name] [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `type` | Integration type (databricks, aws, grafana, kafka, snowflake) |
| `name` | Integration name (optional) |

### Options

| Option | Description |
|--------|-------------|
| `--verbose` | Show detailed test information |
| `--output` | Output format (text, json) |
| `--timeout` | Test timeout in seconds |

### Examples

```bash
# Test all Databricks integrations
nessi integration test databricks

# Test a specific integration
nessi integration test databricks prod-workspace

# Show detailed test information
nessi integration test databricks prod-workspace --verbose
```

### Output

```
Integration Test Results:

Integration: databricks (prod-workspace)
Status: PASSED

Connection: SUCCESS (connected in 0.35s)
Authentication: SUCCESS
Permissions: SUCCESS (READ, WRITE, EXECUTE)
API Version: SUCCESS (2.0)
Query Execution: SUCCESS (executed in 0.42s)

All tests passed successfully.
```

## `integration export`

Exports integration configuration.

### Usage

```bash
nessi integration export <type> <name> <output-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `type` | Integration type (databricks, aws, grafana, kafka, snowflake) |
| `name` | Integration name |
| `output-path` | Path to export the configuration to |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Export format (yaml, json) |
| `--include-secrets` | Include secrets in the export (not recommended) |
| `--encrypt` | Encrypt the exported configuration |
| `--password` | Password for encryption |

### Examples

```bash
# Export integration configuration
nessi integration export databricks prod-workspace /exports/databricks-prod.yaml

# Export without secrets
nessi integration export databricks prod-workspace /exports/databricks-prod.yaml --include-secrets=false

# Export with encryption
nessi integration export databricks prod-workspace /exports/databricks-prod.yaml --encrypt --password my-secure-password
```

### Output

```
Exporting integration configuration:
  Type: databricks
  Name: prod-workspace
  Format: yaml
  Output Path: /exports/databricks-prod.yaml
  Include Secrets: false
  Encryption: none

Export completed successfully.
```

## `integration import`

Imports integration configuration.

### Usage

```bash
nessi integration import <config-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `config-path` | Path to the configuration file |

### Options

| Option | Description |
|--------|-------------|
| `--name` | Name for the imported integration |
| `--overwrite` | Overwrite existing integration with the same name |
| `--validate` | Validate configuration before importing |
| `--decrypt` | Decrypt the configuration file |
| `--password` | Password for decryption |

### Examples

```bash
# Import integration configuration
nessi integration import /exports/databricks-prod.yaml

# Import with a new name
nessi integration import /exports/databricks-prod.yaml --name databricks-staging

# Import and overwrite existing integration
nessi integration import /exports/databricks-prod.yaml --overwrite

# Import encrypted configuration
nessi integration import /exports/databricks-prod.yaml --decrypt --password my-secure-password
```

### Output

```
Importing integration configuration:
  Path: /exports/databricks-prod.yaml
  Type: databricks
  Name: prod-workspace
  Validation: passed

Integration imported successfully.
```

## Error Handling

Integration commands use the following error codes:

- `N1100`: Integration configuration failed
- `N1101`: Connection failed
- `N1102`: Authentication failed
- `N1103`: Permission denied
- `N1104`: Resource not found
- `N1105`: Import/export operation failed
- `N1106`: Query execution failed
- `N1107`: Data transfer failed
- `N1108`: Rate limit exceeded
- `N1109`: Service unavailable

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
