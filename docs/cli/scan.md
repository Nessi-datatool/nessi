# Scan Commands

This document provides detailed information about Nessi's scan commands for analyzing, profiling, and validating data.

## `scan run`

Runs a data scan on a table or file.

### Usage

```bash
nessi scan run <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table, file, or directory to scan |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Scan type (profile, quality, schema, performance, all) |
| `--format` | Source format (delta, parquet, csv, json, auto) |
| `--sample` | Sampling percentage (0-100) |
| `--columns` | Comma-separated list of columns to scan |
| `--rules` | Path to rules file for quality scans |
| `--output` | Output format (text, json, csv, html, pdf) |
| `--file` | Output file path |
| `--timeout` | Scan timeout in seconds |
| `--verbose` | Show detailed progress information |

### Examples

```bash
# Run a full scan on a Delta table
nessi scan run /data/customers

# Run a profile scan
nessi scan run /data/customers --type profile

# Run a scan with sampling
nessi scan run /data/customers --sample 50

# Scan specific columns
nessi scan run /data/customers --columns id,name,email

# Generate an HTML report
nessi scan run /data/customers --output html --file scan_report.html

# Scan a CSV file
nessi scan run /data/exports/customers.csv --format csv
```

### Output

#### Text Format (Default)

```
Scan Results: /data/customers
Timestamp: 2025-05-22 22:12:09
Format: delta
Rows: 1,000,000 (100% sample)
Scan Duration: 45.2s

Profile Summary:
  Columns: 6
  Data Size: 1.2 GB
  Null Values: 153,700 (2.56%)
  Distinct Values Ratio: 0.82

Quality Summary:
  Rules Checked: 4
  Rules Passed: 3
  Rules Failed: 1
  Overall Quality Score: 0.95

Schema Summary:
  Data Types: string (3), integer (1), date (1), timestamp (1)
  Primary Keys: id
  Foreign Keys: 0
  Constraints: 2

Performance Summary:
  Scan Duration: 45.2s
  Memory Usage: 1.2 GB
  CPU Usage: 75%
  I/O Operations: 15,000

Column Details:
  id (string):
    Completeness: 100.00%
    Uniqueness: 100.00%
    Min Length: 36
    Max Length: 36
    Pattern: UUID format (100.00%)

  name (string):
    Completeness: 100.00%
    Uniqueness: 95.00%
    Min Length: 3
    Max Length: 50
    Average Length: 22.5

  email (string):
    Completeness: 99.75%
    Uniqueness: 99.75%
    Min Length: 10
    Max Length: 50
    Average Length: 25.3
    Pattern: Contains "@" (99.95%)

  age (integer):
    Completeness: 99.88%
    Min: 18
    Max: 95
    Mean: 42.5
    Median: 41
    StdDev: 15.2

  signup_date (date):
    Completeness: 100.00%
    Min: 2020-01-01
    Max: 2025-05-01

  last_login (timestamp):
    Completeness: 85.00%
    Min: 2020-01-01T00:00:00Z
    Max: 2025-05-01T12:34:56Z
```

#### HTML Report

When using the `--output html` option, a comprehensive HTML report is generated with:
- Modern blue color scheme
- Summary cards showing key metrics
- Detailed tables of column statistics
- Visualizations of data distributions
- Interactive elements with tab navigation
- Responsive design for different screen sizes

## `scan compare`

Compares scan results between two sources or two versions of the same source.

### Usage

```bash
nessi scan compare <source1> <source2> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source1` | Path to the first table, file, or scan result |
| `source2` | Path to the second table, file, or scan result |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Comparison type (profile, quality, schema, performance, all) |
| `--format1` | Format of the first source (delta, parquet, csv, json, auto) |
| `--format2` | Format of the second source (delta, parquet, csv, json, auto) |
| `--sample` | Sampling percentage (0-100) |
| `--columns` | Comma-separated list of columns to compare |
| `--output` | Output format (text, json, csv, html, pdf) |
| `--file` | Output file path |
| `--threshold` | Difference threshold for highlighting (0.0-1.0) |

### Examples

```bash
# Compare two tables
nessi scan compare /data/customers /data/customers_backup

# Compare scan results
nessi scan compare /results/scan_20250501.json /results/scan_20250522.json

# Compare specific columns
nessi scan compare /data/customers /data/customers_backup --columns id,name,email

# Generate an HTML report
nessi scan compare /data/customers /data/customers_backup --output html --file comparison_report.html

# Compare with a specific threshold
nessi scan compare /data/customers /data/customers_backup --threshold 0.05
```

### Output

#### Text Format (Default)

```
Scan Comparison Results
Source 1: /data/customers
Source 2: /data/customers_backup
Timestamp: 2025-05-22 22:12:09

General Comparison:
  Rows: 1,000,000 vs 950,000 (-5.00%) ⚠️
  Size: 1.2 GB vs 1.15 GB (-4.17%)
  Columns: 6 vs 6 (0.00%)
  Null Values: 2.56% vs 2.60% (+0.04%)
  Distinct Values Ratio: 0.82 vs 0.81 (-0.01%)

Quality Comparison:
  Overall Quality Score: 0.95 vs 0.94 (-0.01%)
  Rules Passed: 3 vs 3 (0.00%)
  Rules Failed: 1 vs 1 (0.00%)

Schema Comparison:
  Data Types: Identical
  Primary Keys: Identical
  Foreign Keys: Identical
  Constraints: 2 vs 1 (-50.00%) ⚠️

Performance Comparison:
  Scan Duration: 45.2s vs 43.1s (-4.65%)
  Memory Usage: 1.2 GB vs 1.15 GB (-4.17%)
  CPU Usage: 75% vs 72% (-4.00%)
  I/O Operations: 15,000 vs 14,250 (-5.00%)

Column Comparisons:
  id (string):
    Completeness: 100.00% vs 100.00% (0.00%)
    Uniqueness: 100.00% vs 100.00% (0.00%)
    Pattern Match: 100.00% vs 100.00% (0.00%)

  name (string):
    Completeness: 100.00% vs 100.00% (0.00%)
    Uniqueness: 95.00% vs 94.50% (-0.53%)
    Average Length: 22.5 vs 22.3 (-0.89%)

  email (string):
    Completeness: 99.75% vs 99.70% (-0.05%)
    Uniqueness: 99.75% vs 99.70% (-0.05%)
    Pattern Match: 99.95% vs 99.90% (-0.05%)

  age (integer):
    Completeness: 99.88% vs 99.85% (-0.03%)
    Min: 18 vs 18 (0.00%)
    Max: 95 vs 95 (0.00%)
    Mean: 42.5 vs 42.3 (-0.47%)
    StdDev: 15.2 vs 15.1 (-0.66%)

  signup_date (date):
    Completeness: 100.00% vs 100.00% (0.00%)
    Min: 2020-01-01 vs 2020-01-01 (0.00%)
    Max: 2025-05-01 vs 2025-05-01 (0.00%)

  last_login (timestamp):
    Completeness: 85.00% vs 84.50% (-0.59%)
    Min: 2020-01-01T00:00:00Z vs 2020-01-01T00:00:00Z (0.00%)
    Max: 2025-05-01T12:34:56Z vs 2025-05-01T12:34:56Z (0.00%)

Legend:
  ⚠️ Difference exceeds threshold (0.05)
```

## `scan history`

Shows the history of scans for a specific source.

### Usage

```bash
nessi scan history <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table or file |

### Options

| Option | Description |
|--------|-------------|
| `--limit` | Maximum number of historical scans to show |
| `--from` | Start date for history (YYYY-MM-DD) |
| `--to` | End date for history (YYYY-MM-DD) |
| `--type` | Filter by scan type (profile, quality, schema, performance) |
| `--output` | Output format (text, json, csv, html) |
| `--file` | Output file path |
| `--trend` | Show trend analysis |

### Examples

```bash
# Show scan history
nessi scan history /data/customers

# Show limited history
nessi scan history /data/customers --limit 5

# Show history for a date range
nessi scan history /data/customers --from 2025-04-01 --to 2025-05-01

# Show history for a specific scan type
nessi scan history /data/customers --type quality

# Show history with trend analysis
nessi scan history /data/customers --trend
```

### Output

#### Text Format (Default)

```
Scan History: /data/customers
Period: 2025-04-01 to 2025-05-22

DATE        TYPE      ROWS      SIZE    QUALITY  DURATION  STATUS
2025-05-22  all       1,000,000 1.2 GB  0.95     45.2s     SUCCESS
2025-05-15  all       980,000   1.18 GB 0.94     44.5s     SUCCESS
2025-05-08  all       975,000   1.17 GB 0.94     44.1s     SUCCESS
2025-05-01  all       950,000   1.15 GB 0.94     43.1s     SUCCESS
2025-04-24  all       930,000   1.12 GB 0.93     42.5s     SUCCESS
2025-04-17  all       925,000   1.11 GB 0.93     42.2s     SUCCESS
2025-04-10  all       920,000   1.10 GB 0.92     41.8s     SUCCESS
2025-04-03  all       900,000   1.08 GB 0.92     41.0s     SUCCESS

Trends:
  Rows: +11.11% (900,000 → 1,000,000)
  Size: +11.11% (1.08 GB → 1.2 GB)
  Quality Score: +3.26% (0.92 → 0.95)
  Scan Duration: +10.24% (41.0s → 45.2s)
```

## `scan schedule`

Manages scheduled scans.

### Usage

```bash
nessi scan schedule [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `list` | List all scheduled scans |
| `create` | Create a new scheduled scan |
| `update` | Update an existing scheduled scan |
| `delete` | Delete a scheduled scan |
| `enable` | Enable a scheduled scan |
| `disable` | Disable a scheduled scan |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--source` | Filter by source path |
| `--status` | Filter by status (enabled, disabled) |
| `--output` | Output format (text, json, csv) |

### Options for `create` and `update`

| Option | Description |
|--------|-------------|
| `--id` | Schedule ID (for update only) |
| `--source` | Path to the table or file to scan |
| `--type` | Scan type (profile, quality, schema, performance, all) |
| `--schedule` | Schedule expression (cron format) |
| `--sample` | Sampling percentage (0-100) |
| `--output` | Output format for scan results (text, json, csv, html, pdf) |
| `--file` | Output file path pattern |
| `--notify` | Notification settings (email, webhook, none) |
| `--description` | Description of the scheduled scan |

### Options for `delete`, `enable`, and `disable`

| Option | Description |
|--------|-------------|
| `--id` | Schedule ID |
| `--force` | Force operation without confirmation (for delete) |

### Examples

```bash
# List all scheduled scans
nessi scan schedule list

# Create a scheduled scan
nessi scan schedule create --source /data/customers --type all --schedule "0 0 * * *" --output html --file "/reports/customers_scan_%Y%m%d.html" --notify email --description "Daily customer data scan"

# Update a scheduled scan
nessi scan schedule update --id 12345 --schedule "0 0 * * 1" --description "Weekly customer data scan"

# Delete a scheduled scan
nessi scan schedule delete --id 12345

# Enable a scheduled scan
nessi scan schedule enable --id 12345

# Disable a scheduled scan
nessi scan schedule disable --id 12345
```

### Output for `list`

```
Scheduled Scans:

ID     SOURCE            TYPE  SCHEDULE    STATUS   LAST RUN            NEXT RUN
12345  /data/customers   all   0 0 * * *   enabled  2025-05-22 00:00:00 2025-05-23 00:00:00
67890  /data/orders      all   0 0 * * 1   enabled  2025-05-20 00:00:00 2025-05-27 00:00:00
24680  /data/products    all   0 0 1 * *   disabled 2025-05-01 00:00:00 N/A
```

### Output for `create` and `update`

```
Scheduled Scan Created:
  ID: 12345
  Source: /data/customers
  Type: all
  Schedule: 0 0 * * * (Daily at midnight)
  Output: html to "/reports/customers_scan_%Y%m%d.html"
  Notification: email
  Status: enabled
  Next Run: 2025-05-23 00:00:00
```

## `scan export`

Exports scan results to various formats.

### Usage

```bash
nessi scan export <source> <output-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the scan result file |
| `output-path` | Path to export the results to |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Export format (html, pdf, json, csv, xlsx) |
| `--template` | Report template to use (default, enhanced, minimal) |
| `--include-visualizations` | Include data visualizations |
| `--include-recommendations` | Include improvement recommendations |
| `--title` | Custom report title |
| `--logo` | Path to custom logo image |

### Examples

```bash
# Export scan results to HTML
nessi scan export /results/scan_20250522.json /reports/scan_report.html --format html

# Export with enhanced template
nessi scan export /results/scan_20250522.json /reports/scan_report.html --format html --template enhanced

# Export to PDF with visualizations
nessi scan export /results/scan_20250522.json /reports/scan_report.pdf --format pdf --include-visualizations

# Export with custom title and logo
nessi scan export /results/scan_20250522.json /reports/scan_report.html --format html --title "Customer Data Quality Report" --logo /branding/logo.png
```

### Output

```
Exporting scan results:
  Source: /results/scan_20250522.json
  Format: html
  Template: enhanced
  Output: /reports/scan_report.html
  Visualizations: included
  Recommendations: included

Export completed successfully.
```

## Report Generation System

The scan commands utilize Nessi's robust report generation system with the following features:

1. **Multiple Output Formats**:
   - HTML reports using Jinja2 templates
   - PDF generation with wkhtmltopdf and reportlab fallback
   - JSON serialization with custom serializer for datetime objects
   - CSV export functionality

2. **Customizable Templates**:
   - Default template with standard layout
   - Enhanced template with modern blue color scheme, improved typography, and interactive elements
   - Minimal template for simplified reports

3. **Visualization Components**:
   - Data distribution charts
   - Quality score gauges
   - Trend line charts
   - Heatmaps for correlation analysis

4. **Interactive Features** (HTML reports):
   - Tabbed navigation
   - Collapsible sections
   - Sortable tables
   - Filterable content

5. **Error Handling**:
   - Validation for required fields
   - Fallback mechanisms for missing templates
   - Error reporting for failed report generation

For more information on the report generation system, see the [Report Generation documentation](../reports/REPORT_GENERATION.md).

## Error Handling

Scan commands use the following error codes:

- `N600`: Scan source not found
- `N601`: Invalid scan type
- `N602`: Scan timeout
- `N603`: Insufficient permissions
- `N604`: Invalid scan result format
- `N605`: Scan comparison failed
- `N606`: Schedule creation failed
- `N607`: Schedule not found
- `N608`: Export failed

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
