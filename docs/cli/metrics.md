# Metrics Commands

This document provides detailed information about Nessi's metrics commands for collecting, analyzing, and visualizing data quality metrics.

## `metrics collect`

Collects metrics from a table or file.

### Usage

```bash
nessi metrics collect <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table or file to collect metrics from |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Metrics type (quality, performance, freshness, all) |
| `--format` | Source format (delta, parquet, csv, json, auto) |
| `--sample` | Sampling percentage (0-100) |
| `--columns` | Comma-separated list of columns to collect metrics for |
| `--output` | Output format (text, json, csv) |
| `--file` | Output file path |
| `--store` | Store metrics in the metrics store |
| `--tag` | Tags to associate with the metrics (key=value format) |

### Examples

```bash
# Collect all metrics
nessi metrics collect /data/customers

# Collect specific metric types
nessi metrics collect /data/customers --type quality,performance

# Collect metrics for specific columns
nessi metrics collect /data/customers --columns id,name,email

# Collect metrics with sampling
nessi metrics collect /data/customers --sample 50

# Store metrics in the metrics store
nessi metrics collect /data/customers --store

# Add tags to metrics
nessi metrics collect /data/customers --store --tag environment=production --tag region=us-west
```

### Output

#### Text Format (Default)

```
Metrics Collection: /data/customers
Timestamp: 2025-05-22 22:12:09
Sample: 100% (1,000,000 rows)

Quality Metrics:
  Completeness: 0.9744 (97.44%)
  Accuracy: 0.9950 (99.50%)
  Consistency: 0.9900 (99.00%)
  Uniqueness: 0.9950 (99.50%)
  Timeliness: 0.9850 (98.50%)
  Overall Quality Score: 0.9879 (98.79%)

Performance Metrics:
  Scan Duration: 45.2s
  Memory Usage: 1.2 GB
  CPU Usage: 75%
  I/O Operations: 15,000
  Rows Processed: 1,000,000
  Bytes Processed: 1.2 GB

Freshness Metrics:
  Last Updated: 2025-05-01 12:34:56 (21 days ago)
  Update Frequency: 1.2 days
  Staleness Score: 0.85 (85%)

Column Metrics:
  id:
    Completeness: 1.0000 (100.00%)
    Uniqueness: 1.0000 (100.00%)
    Pattern Match: 1.0000 (100.00%)

  name:
    Completeness: 1.0000 (100.00%)
    Uniqueness: 0.9500 (95.00%)
    Pattern Match: 1.0000 (100.00%)

  email:
    Completeness: 0.9975 (99.75%)
    Uniqueness: 0.9975 (99.75%)
    Pattern Match: 0.9995 (99.95%)

  age:
    Completeness: 0.9988 (99.88%)
    Range Match: 0.9998 (99.98%)
    Distribution: Normal (mean=42.5, stddev=15.2)

  signup_date:
    Completeness: 1.0000 (100.00%)
    Range Match: 1.0000 (100.00%)
    Timeliness: 0.9500 (95.00%)

  last_login:
    Completeness: 0.8500 (85.00%)
    Range Match: 1.0000 (100.00%)
    Timeliness: 0.7500 (75.00%)

Metrics stored with ID: met_12345678
Tags: environment=production, region=us-west
```

## `metrics list`

Lists stored metrics.

### Usage

```bash
nessi metrics list [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--source` | Filter by source path |
| `--type` | Filter by metrics type (quality, performance, freshness) |
| `--from` | Start date for metrics (YYYY-MM-DD) |
| `--to` | End date for metrics (YYYY-MM-DD) |
| `--tag` | Filter by tag (key=value format) |
| `--limit` | Maximum number of metrics to list |
| `--output` | Output format (text, json, csv) |

### Examples

```bash
# List all metrics
nessi metrics list

# Filter by source
nessi metrics list --source /data/customers

# Filter by type
nessi metrics list --type quality

# Filter by date range
nessi metrics list --from 2025-04-01 --to 2025-05-01

# Filter by tag
nessi metrics list --tag environment=production

# Limit results
nessi metrics list --limit 10
```

### Output

#### Text Format (Default)

```
Stored Metrics:

ID           SOURCE             TYPE       TIMESTAMP            TAGS
met_12345678 /data/customers    all        2025-05-22 22:12:09  environment=production,region=us-west
met_23456789 /data/customers    all        2025-05-15 22:00:00  environment=production,region=us-west
met_34567890 /data/customers    all        2025-05-08 22:00:00  environment=production,region=us-west
met_45678901 /data/customers    all        2025-05-01 22:00:00  environment=production,region=us-west
met_56789012 /data/orders       all        2025-05-22 23:00:00  environment=production,region=us-west
met_67890123 /data/orders       all        2025-05-15 23:00:00  environment=production,region=us-west
met_78901234 /data/products     all        2025-05-22 00:00:00  environment=production,region=us-west
met_89012345 /data/products     all        2025-05-15 00:00:00  environment=production,region=us-west
```

## `metrics show`

Shows detailed information about specific metrics.

### Usage

```bash
nessi metrics show <metrics-id> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `metrics-id` | ID of the metrics to show |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Filter by metrics type (quality, performance, freshness) |
| `--column` | Filter by column name |
| `--output` | Output format (text, json, csv, html) |
| `--file` | Output file path |
| `--detailed` | Show detailed metrics information |

### Examples

```bash
# Show metrics details
nessi metrics show met_12345678

# Show specific metric type
nessi metrics show met_12345678 --type quality

# Show metrics for a specific column
nessi metrics show met_12345678 --column email

# Show detailed metrics
nessi metrics show met_12345678 --detailed

# Export metrics to HTML
nessi metrics show met_12345678 --output html --file metrics_report.html
```

### Output

#### Text Format (Default)

```
Metrics Details: met_12345678
Source: /data/customers
Timestamp: 2025-05-22 22:12:09
Sample: 100% (1,000,000 rows)
Tags: environment=production, region=us-west

Quality Metrics:
  Completeness: 0.9744 (97.44%)
  Accuracy: 0.9950 (99.50%)
  Consistency: 0.9900 (99.00%)
  Uniqueness: 0.9950 (99.50%)
  Timeliness: 0.9850 (98.50%)
  Overall Quality Score: 0.9879 (98.79%)

Performance Metrics:
  Scan Duration: 45.2s
  Memory Usage: 1.2 GB
  CPU Usage: 75%
  I/O Operations: 15,000
  Rows Processed: 1,000,000
  Bytes Processed: 1.2 GB

Freshness Metrics:
  Last Updated: 2025-05-01 12:34:56 (21 days ago)
  Update Frequency: 1.2 days
  Staleness Score: 0.85 (85%)

Column Metrics:
  id:
    Completeness: 1.0000 (100.00%)
    Uniqueness: 1.0000 (100.00%)
    Pattern Match: 1.0000 (100.00%)

  name:
    Completeness: 1.0000 (100.00%)
    Uniqueness: 0.9500 (95.00%)
    Pattern Match: 1.0000 (100.00%)

  email:
    Completeness: 0.9975 (99.75%)
    Uniqueness: 0.9975 (99.75%)
    Pattern Match: 0.9995 (99.95%)

  age:
    Completeness: 0.9988 (99.88%)
    Range Match: 0.9998 (99.98%)
    Distribution: Normal (mean=42.5, stddev=15.2)

  signup_date:
    Completeness: 1.0000 (100.00%)
    Range Match: 1.0000 (100.00%)
    Timeliness: 0.9500 (95.00%)

  last_login:
    Completeness: 0.8500 (85.00%)
    Range Match: 1.0000 (100.00%)
    Timeliness: 0.7500 (75.00%)
```

## `metrics compare`

Compares metrics between different sources or time periods.

### Usage

```bash
nessi metrics compare <metrics-id1> <metrics-id2> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `metrics-id1` | ID of the first metrics |
| `metrics-id2` | ID of the second metrics |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Filter by metrics type (quality, performance, freshness) |
| `--column` | Filter by column name |
| `--output` | Output format (text, json, csv, html) |
| `--file` | Output file path |
| `--threshold` | Difference threshold for highlighting (0.0-1.0) |

### Examples

```bash
# Compare two metrics
nessi metrics compare met_12345678 met_23456789

# Compare specific metric types
nessi metrics compare met_12345678 met_23456789 --type quality

# Compare specific columns
nessi metrics compare met_12345678 met_23456789 --column email

# Generate HTML comparison report
nessi metrics compare met_12345678 met_23456789 --output html --file metrics_comparison.html

# Set difference threshold
nessi metrics compare met_12345678 met_23456789 --threshold 0.05
```

### Output

#### Text Format (Default)

```
Metrics Comparison
Metrics 1: met_12345678 (2025-05-22 22:12:09)
Metrics 2: met_23456789 (2025-05-15 22:00:00)
Source: /data/customers

Quality Metrics:
  Completeness: 0.9744 vs 0.9730 (+0.0014, +0.14%) 
  Accuracy: 0.9950 vs 0.9940 (+0.0010, +0.10%)
  Consistency: 0.9900 vs 0.9880 (+0.0020, +0.20%)
  Uniqueness: 0.9950 vs 0.9950 (0.0000, 0.00%)
  Timeliness: 0.9850 vs 0.9820 (+0.0030, +0.31%)
  Overall Quality Score: 0.9879 vs 0.9864 (+0.0015, +0.15%)

Performance Metrics:
  Scan Duration: 45.2s vs 44.8s (+0.4s, +0.89%)
  Memory Usage: 1.2 GB vs 1.2 GB (0.0 GB, 0.00%)
  CPU Usage: 75% vs 74% (+1%, +1.35%)
  I/O Operations: 15,000 vs 14,800 (+200, +1.35%)
  Rows Processed: 1,000,000 vs 980,000 (+20,000, +2.04%)
  Bytes Processed: 1.2 GB vs 1.18 GB (+0.02 GB, +1.69%)

Freshness Metrics:
  Last Updated: 2025-05-01 vs 2025-04-24 (+7 days)
  Update Frequency: 1.2 days vs 1.2 days (0.0 days, 0.00%)
  Staleness Score: 0.85 vs 0.82 (+0.03, +3.66%)

Column Metrics:
  email:
    Completeness: 0.9975 vs 0.9970 (+0.0005, +0.05%)
    Uniqueness: 0.9975 vs 0.9975 (0.0000, 0.00%)
    Pattern Match: 0.9995 vs 0.9990 (+0.0005, +0.05%)

Legend:
  Values shown as: Metrics 1 vs Metrics 2 (absolute difference, percentage difference)
  Positive differences indicate improvement in Metrics 1 compared to Metrics 2
```

## `metrics trend`

Analyzes trends in metrics over time.

### Usage

```bash
nessi metrics trend <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table or file |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Metrics type (quality, performance, freshness, all) |
| `--column` | Filter by column name |
| `--from` | Start date for trend analysis (YYYY-MM-DD) |
| `--to` | End date for trend analysis (YYYY-MM-DD) |
| `--period` | Aggregation period (day, week, month) |
| `--output` | Output format (text, json, csv, html) |
| `--file` | Output file path |
| `--visualize` | Include visualizations in the output |

### Examples

```bash
# Analyze all metrics trends
nessi metrics trend /data/customers

# Analyze specific metric types
nessi metrics trend /data/customers --type quality

# Analyze trends for a specific column
nessi metrics trend /data/customers --column email

# Analyze trends for a date range
nessi metrics trend /data/customers --from 2025-04-01 --to 2025-05-22

# Aggregate by week
nessi metrics trend /data/customers --period week

# Generate HTML report with visualizations
nessi metrics trend /data/customers --output html --file trend_report.html --visualize
```

### Output

#### Text Format (Default)

```
Metrics Trend Analysis: /data/customers
Period: 2025-04-01 to 2025-05-22
Aggregation: Weekly

Quality Metrics Trend:

DATE        COMPLETENESS  ACCURACY  CONSISTENCY  UNIQUENESS  TIMELINESS  OVERALL
2025-04-01  0.9700        0.9920    0.9850       0.9950      0.9780      0.9840
2025-04-08  0.9710        0.9925    0.9860       0.9950      0.9790      0.9847
2025-04-15  0.9720        0.9930    0.9870       0.9950      0.9800      0.9854
2025-04-22  0.9730        0.9940    0.9880       0.9950      0.9820      0.9864
2025-04-29  0.9735        0.9945    0.9890       0.9950      0.9830      0.9870
2025-05-06  0.9740        0.9945    0.9895       0.9950      0.9840      0.9874
2025-05-13  0.9742        0.9948    0.9898       0.9950      0.9845      0.9877
2025-05-20  0.9744        0.9950    0.9900       0.9950      0.9850      0.9879

CHANGE      +0.0044       +0.0030   +0.0050      +0.0000     +0.0070     +0.0039
% CHANGE    +0.45%        +0.30%    +0.51%       +0.00%      +0.72%      +0.40%
TREND       ↑             ↑         ↑            →           ↑           ↑

Performance Metrics Trend:

DATE        DURATION (s)  MEMORY (GB)  CPU (%)  I/O OPS   ROWS (M)  BYTES (GB)
2025-04-01  41.0          1.08         72       14,000    0.90      1.08
2025-04-08  41.8          1.10         72       14,200    0.92      1.10
2025-04-15  42.5          1.12         73       14,400    0.93      1.12
2025-04-22  43.1          1.15         74       14,800    0.98      1.18
2025-04-29  43.8          1.16         74       14,850    0.99      1.19
2025-05-06  44.2          1.18         74       14,900    0.99      1.19
2025-05-13  44.8          1.20         75       15,000    1.00      1.20
2025-05-20  45.2          1.20         75       15,000    1.00      1.20

CHANGE      +4.2          +0.12        +3        +1,000    +0.10     +0.12
% CHANGE    +10.24%       +11.11%      +4.17%    +7.14%    +11.11%   +11.11%
TREND       ↑             ↑            ↑         ↑         ↑         ↑

Column Metrics Trend (email):

DATE        COMPLETENESS  UNIQUENESS  PATTERN MATCH
2025-04-01  0.9950        0.9970      0.9980
2025-04-08  0.9955        0.9970      0.9982
2025-04-15  0.9960        0.9972      0.9985
2025-04-22  0.9965        0.9975      0.9988
2025-04-29  0.9968        0.9975      0.9990
2025-05-06  0.9970        0.9975      0.9992
2025-05-13  0.9972        0.9975      0.9993
2025-05-20  0.9975        0.9975      0.9995

CHANGE      +0.0025       +0.0005     +0.0015
% CHANGE    +0.25%        +0.05%      +0.15%
TREND       ↑             ↑           ↑

Trend Analysis:
- Quality metrics show consistent improvement over the period
- Performance metrics show increasing resource usage with data growth
- Column-level metrics for 'email' show improving data quality
```

## `metrics alert`

Manages metric-based alerts.

### Usage

```bash
nessi metrics alert [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `list` | List all alerts |
| `create` | Create a new alert |
| `update` | Update an existing alert |
| `delete` | Delete an alert |
| `enable` | Enable an alert |
| `disable` | Disable an alert |
| `history` | Show alert history |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--source` | Filter by source path |
| `--type` | Filter by metrics type (quality, performance, freshness) |
| `--status` | Filter by status (enabled, disabled, triggered) |
| `--output` | Output format (text, json, csv) |

### Options for `create` and `update`

| Option | Description |
|--------|-------------|
| `--id` | Alert ID (for update only) |
| `--name` | Alert name |
| `--source` | Path to the table or file to monitor |
| `--metric` | Metric to monitor (e.g., quality.completeness, performance.duration) |
| `--column` | Column name (for column-specific metrics) |
| `--condition` | Alert condition (e.g., "<", ">", "<=", ">=", "==", "!=") |
| `--threshold` | Alert threshold value |
| `--severity` | Alert severity (critical, high, medium, low) |
| `--notify` | Notification settings (email, webhook, none) |
| `--description` | Alert description |

### Options for `delete`, `enable`, and `disable`

| Option | Description |
|--------|-------------|
| `--id` | Alert ID |
| `--force` | Force operation without confirmation (for delete) |

### Options for `history`

| Option | Description |
|--------|-------------|
| `--id` | Alert ID |
| `--from` | Start date for history (YYYY-MM-DD) |
| `--to` | End date for history (YYYY-MM-DD) |
| `--limit` | Maximum number of history entries to show |
| `--output` | Output format (text, json, csv) |

### Examples

```bash
# List all alerts
nessi metrics alert list

# Create a new alert
nessi metrics alert create --name "Low Completeness Alert" --source /data/customers --metric quality.completeness --condition "<" --threshold 0.95 --severity high --notify email --description "Alert when data completeness falls below 95%"

# Update an existing alert
nessi metrics alert update --id 12345 --threshold 0.97 --severity critical

# Delete an alert
nessi metrics alert delete --id 12345

# Enable an alert
nessi metrics alert enable --id 12345

# Disable an alert
nessi metrics alert disable --id 12345

# Show alert history
nessi metrics alert history --id 12345 --from 2025-04-01 --to 2025-05-22
```

### Output for `list`

```
Metric Alerts:

ID     NAME                   SOURCE            METRIC                 CONDITION  THRESHOLD  SEVERITY  STATUS
12345  Low Completeness Alert /data/customers   quality.completeness   <          0.95       high      enabled
67890  High CPU Usage Alert   /data/customers   performance.cpu_usage  >          90         critical  enabled
24680  Data Freshness Alert   /data/customers   freshness.staleness    >          0.8        medium    disabled
```

### Output for `create` and `update`

```
Alert Created:
  ID: 12345
  Name: Low Completeness Alert
  Source: /data/customers
  Metric: quality.completeness
  Condition: < 0.95
  Severity: high
  Notification: email
  Status: enabled
```

### Output for `history`

```
Alert History: Low Completeness Alert (ID: 12345)
Period: 2025-04-01 to 2025-05-22

TIMESTAMP            STATUS     VALUE    THRESHOLD  MESSAGE
2025-05-10 22:15:30  TRIGGERED  0.9480   0.9500     Completeness below threshold: 94.80% (threshold: 95.00%)
2025-05-03 22:10:15  TRIGGERED  0.9490   0.9500     Completeness below threshold: 94.90% (threshold: 95.00%)
2025-04-26 22:05:45  OK         0.9520   0.9500     Completeness above threshold: 95.20% (threshold: 95.00%)
2025-04-19 22:00:30  TRIGGERED  0.9470   0.9500     Completeness below threshold: 94.70% (threshold: 95.00%)
2025-04-12 22:00:15  TRIGGERED  0.9450   0.9500     Completeness below threshold: 94.50% (threshold: 95.00%)
2025-04-05 22:00:00  TRIGGERED  0.9430   0.9500     Completeness below threshold: 94.30% (threshold: 95.00%)
```

## `metrics visualize`

Visualizes metrics data.

### Usage

```bash
nessi metrics visualize <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table, file, or metrics ID |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Visualization type (line, bar, heatmap, radar, gauge) |
| `--metric` | Metric to visualize (e.g., quality.completeness, performance.duration) |
| `--column` | Column name (for column-specific metrics) |
| `--from` | Start date for visualization (YYYY-MM-DD) |
| `--to` | End date for visualization (YYYY-MM-DD) |
| `--output` | Output format (html, png, pdf) |
| `--file` | Output file path |
| `--title` | Visualization title |
| `--theme` | Visualization theme (default, dark, light, colorful) |
| `--interactive` | Enable interactive visualizations (for HTML output) |

### Examples

```bash
# Visualize quality metrics
nessi metrics visualize /data/customers --type radar --metric quality --output html --file quality_radar.html

# Visualize completeness trend
nessi metrics visualize /data/customers --type line --metric quality.completeness --from 2025-04-01 --to 2025-05-22 --output html --file completeness_trend.html

# Visualize column metrics
nessi metrics visualize /data/customers --type heatmap --metric quality --column all --output html --file column_quality_heatmap.html

# Create an interactive dashboard
nessi metrics visualize /data/customers --type line,radar,gauge --metric quality --interactive --output html --file metrics_dashboard.html

# Visualize with custom theme and title
nessi metrics visualize /data/customers --type radar --metric quality --theme dark --title "Data Quality Radar" --output png --file quality_radar.png
```

### Output

```
Metrics Visualization:
  Source: /data/customers
  Metrics: quality
  Visualization Type: radar
  Period: 2025-04-01 to 2025-05-22
  Output: HTML to quality_radar.html

Visualization created successfully.
```

## `metrics export`

Exports metrics data to various formats.

### Usage

```bash
nessi metrics export <source> <output-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table, file, or metrics ID |
| `output-path` | Path to export the metrics to |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Export format (json, csv, xlsx, html, pdf) |
| `--type` | Metrics type to export (quality, performance, freshness, all) |
| `--column` | Column name (for column-specific metrics) |
| `--from` | Start date for metrics (YYYY-MM-DD) |
| `--to` | End date for metrics (YYYY-MM-DD) |
| `--include-visualizations` | Include visualizations in the export |
| `--template` | Report template to use (default, enhanced, minimal) |

### Examples

```bash
# Export metrics to JSON
nessi metrics export /data/customers /exports/customers_metrics.json --format json

# Export specific metric types
nessi metrics export /data/customers /exports/customers_quality.csv --format csv --type quality

# Export metrics for a date range
nessi metrics export /data/customers /exports/customers_metrics.xlsx --format xlsx --from 2025-04-01 --to 2025-05-22

# Export with visualizations
nessi metrics export /data/customers /exports/customers_metrics.html --format html --include-visualizations

# Export using enhanced template
nessi metrics export /data/customers /exports/customers_metrics.html --format html --template enhanced
```

### Output

```
Exporting metrics:
  Source: /data/customers
  Format: html
  Metrics Types: all
  Period: 2025-04-01 to 2025-05-22
  Visualizations: included
  Template: enhanced
  Output: /exports/customers_metrics.html

Export completed successfully.
```

## Grafana Integration

Nessi's metrics commands include integration with Grafana for advanced visualization and monitoring. The integration allows you to:

1. **Create Dashboards**: Automatically create Grafana dashboards from Nessi metrics
2. **Update Dashboards**: Update existing dashboards with new metrics data
3. **Create Alerts**: Set up Grafana alerts based on Nessi metrics
4. **Time-Series Analysis**: Perform advanced time-series analysis on metrics data

### Using Grafana Integration

```bash
# Create a Grafana dashboard for quality metrics
nessi metrics grafana create-dashboard --source /data/customers --type quality --name "Customers Quality Dashboard"

# Update an existing dashboard
nessi metrics grafana update-dashboard --dashboard-id 12345 --source /data/customers

# Create a Grafana alert
nessi metrics grafana create-alert --dashboard-id 12345 --panel-id 1 --name "Low Completeness Alert" --metric quality.completeness --condition "<" --threshold 0.95

# Perform time-series analysis
nessi metrics grafana analyze --source /data/customers --metric quality.completeness --analysis trend,seasonal,anomaly --period 90d
```

For more information on Grafana integration, see the [Grafana Integration documentation](../integrations/GRAFANA.md).

## Error Handling

Metrics commands use the following error codes:

- `N900`: Metrics collection failed
- `N901`: Metrics store connection failed
- `N902`: Metrics retrieval failed
- `N903`: Metrics comparison failed
- `N904`: Trend analysis failed
- `N905`: Alert configuration failed
- `N906`: Visualization failed
- `N907`: Export failed
- `N908`: Grafana integration failed

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
