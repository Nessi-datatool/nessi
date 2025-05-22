# Quality Commands

This document provides detailed information about Nessi's quality commands for performing data quality checks and generating quality reports.

## `quality check`

Runs quality checks on a table against predefined or custom rules.

### Usage

```bash
nessi quality check <table-name> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table to check |

### Options

| Option | Description |
|--------|-------------|
| `--rules` | Comma-separated list of rule names or path to rules file |
| `--rule-type` | Type of rules to apply (completeness, uniqueness, pattern, range, all) |
| `--columns` | Comma-separated list of columns to check |
| `--severity` | Minimum severity level to check (critical, high, medium, low) |
| `--threshold` | Override default threshold for all rules (0.0-1.0) |
| `--sample` | Sampling percentage (0-100) |
| `--output` | Output format (text, json, csv, html, pdf) |
| `--file` | Output file path |
| `--fail-on` | Exit with error if checks fail at this severity or higher |

### Examples

```bash
# Run all quality checks on a table
nessi quality check customers

# Run specific rules
nessi quality check customers --rules email_validation,age_range

# Run checks from a rules file
nessi quality check customers --rules /path/to/rules.json

# Run checks on specific columns
nessi quality check customers --columns email,age

# Run only completeness checks
nessi quality check customers --rule-type completeness

# Generate HTML report
nessi quality check customers --output html --file quality_report.html

# Fail if critical issues are found
nessi quality check customers --fail-on critical
```

### Rules File Format

Rules can be defined in a JSON file:

```json
{
  "rules": [
    {
      "name": "email_validation",
      "description": "Validates email format",
      "type": "pattern",
      "column": "email",
      "pattern": "@",
      "threshold": 0.99,
      "severity": "high"
    },
    {
      "name": "age_range",
      "description": "Validates age is within reasonable range",
      "type": "range",
      "column": "age",
      "min": 0,
      "max": 120,
      "threshold": 1.0,
      "severity": "medium"
    },
    {
      "name": "name_completeness",
      "description": "Checks for missing names",
      "type": "completeness",
      "column": "name",
      "threshold": 0.99,
      "severity": "high"
    },
    {
      "name": "id_uniqueness",
      "description": "Ensures IDs are unique",
      "type": "uniqueness",
      "column": "id",
      "threshold": 1.0,
      "severity": "critical"
    }
  ]
}
```

### Output

#### Text Format (Default)

```
Quality Check Results: customers
Timestamp: 2025-05-22 22:12:09
Rules: 4, Passed: 3, Failed: 1

RULE               TYPE          COLUMN  THRESHOLD  ACTUAL    STATUS  SEVERITY
email_validation   pattern       email   0.99       0.998     PASS    high
age_range          range         age     1.00       0.998     FAIL    medium
name_completeness  completeness  name    0.99       1.000     PASS    high
id_uniqueness      uniqueness    id      1.00       1.000     PASS    critical

Details:
- age_range: 2 violations found (values: 150, -5)
  Rows: 12345, 67890

Summary:
- Critical issues: 0
- High issues: 0
- Medium issues: 1
- Low issues: 0
```

#### HTML Format

When using the `--output html` option, a comprehensive HTML report is generated with:
- Modern blue color scheme
- Summary cards showing pass/fail metrics
- Detailed tables of rule results
- Visualizations of data quality metrics
- Interactive elements with tab navigation
- Responsive design for different screen sizes

## `quality rules list`

Lists all available quality rules.

### Usage

```bash
nessi quality rules list [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--type` | Filter by rule type (completeness, uniqueness, pattern, range) |
| `--severity` | Filter by severity (critical, high, medium, low) |
| `--table` | Filter by table name |
| `--output` | Output format (text, json, csv) |
| `--detailed` | Show detailed rule information |

### Examples

```bash
# List all rules
nessi quality rules list

# List rules by type
nessi quality rules list --type pattern

# List rules by severity
nessi quality rules list --severity critical

# List rules for a specific table
nessi quality rules list --table customers

# List detailed rule information
nessi quality rules list --detailed
```

### Output

#### Text Format (Default)

```
Available Quality Rules:

NAME               TYPE          SEVERITY  TABLE       DESCRIPTION
email_validation   pattern       high      customers   Validates email format
age_range          range         medium    customers   Validates age is within reasonable range
name_completeness  completeness  high      customers   Checks for missing names
id_uniqueness      uniqueness    critical  customers   Ensures IDs are unique
order_id_format    pattern       medium    orders      Validates order ID format
price_range        range         high      orders      Validates price is positive
```

#### Detailed Text Format

```
Available Quality Rules:

Rule: email_validation
  Type: pattern
  Severity: high
  Table: customers
  Column: email
  Pattern: @
  Threshold: 0.99
  Description: Validates email format
  Created: 2025-04-01
  Last Modified: 2025-04-15
  Last Run: 2025-05-22
  Last Status: PASS

Rule: age_range
  Type: range
  Severity: medium
  Table: customers
  Column: age
  Range: 0 to 120
  Threshold: 1.0
  Description: Validates age is within reasonable range
  Created: 2025-04-01
  Last Modified: 2025-04-10
  Last Run: 2025-05-22
  Last Status: FAIL
```

## `quality rules create`

Creates a new quality rule.

### Usage

```bash
nessi quality rules create [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--name` | Rule name |
| `--type` | Rule type (completeness, uniqueness, pattern, range) |
| `--table` | Table name |
| `--column` | Column name |
| `--description` | Rule description |
| `--severity` | Rule severity (critical, high, medium, low) |
| `--threshold` | Rule threshold (0.0-1.0) |
| `--pattern` | Pattern for pattern rules |
| `--min` | Minimum value for range rules |
| `--max` | Maximum value for range rules |
| `--file` | Path to JSON file containing rule definition |

### Examples

```bash
# Create a pattern rule
nessi quality rules create --name email_validation --type pattern --table customers --column email --pattern "@" --threshold 0.99 --severity high --description "Validates email format"

# Create a range rule
nessi quality rules create --name age_range --type range --table customers --column age --min 0 --max 120 --threshold 1.0 --severity medium --description "Validates age is within reasonable range"

# Create a completeness rule
nessi quality rules create --name name_completeness --type completeness --table customers --column name --threshold 0.99 --severity high --description "Checks for missing names"

# Create a rule from a file
nessi quality rules create --file /path/to/rule.json
```

### Rule File Format

```json
{
  "name": "email_validation",
  "type": "pattern",
  "table": "customers",
  "column": "email",
  "pattern": "@",
  "threshold": 0.99,
  "severity": "high",
  "description": "Validates email format"
}
```

### Output

```
Rule created successfully:
  Name: email_validation
  Type: pattern
  Table: customers
  Column: email
  Pattern: @
  Threshold: 0.99
  Severity: high
  Description: Validates email format
```

## `quality rules delete`

Deletes a quality rule.

### Usage

```bash
nessi quality rules delete <rule-name> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `rule-name` | Name of the rule to delete |

### Options

| Option | Description |
|--------|-------------|
| `--table` | Table name (required if rule name is not unique) |
| `--force` | Delete without confirmation |

### Examples

```bash
# Delete a rule
nessi quality rules delete email_validation

# Delete a rule for a specific table
nessi quality rules delete email_validation --table customers

# Delete without confirmation
nessi quality rules delete email_validation --force
```

### Output

```
Rule deleted successfully:
  Name: email_validation
  Table: customers
```

## `quality report`

Generates a comprehensive quality report for a table.

### Usage

```bash
nessi quality report <table-name> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table |

### Options

| Option | Description |
|--------|-------------|
| `--rules` | Comma-separated list of rule names or path to rules file |
| `--include-profile` | Include table profile in the report |
| `--include-schema` | Include schema information in the report |
| `--include-history` | Include quality history in the report |
| `--sample` | Sampling percentage (0-100) |
| `--output` | Output format (html, pdf, json, csv) |
| `--template` | Report template to use (default, enhanced, minimal) |
| `--file` | Output file path |

### Examples

```bash
# Generate a quality report
nessi quality report customers

# Generate a comprehensive report
nessi quality report customers --include-profile --include-schema --include-history

# Generate a PDF report
nessi quality report customers --output pdf --file customers_quality.pdf

# Use the enhanced template
nessi quality report customers --template enhanced
```

### Report Templates

The report generation system supports multiple templates:

1. **Default Template**: Standard report with quality metrics and rule results
2. **Enhanced Template**: Modern blue color scheme with improved typography, interactive elements, and visualizations
3. **Minimal Template**: Simplified report with just the essential information

### Output Formats

The report generation system supports multiple output formats:

1. **HTML**: Interactive web-based report with visualizations and navigation
2. **PDF**: Printable document format generated using wkhtmltopdf or reportlab
3. **JSON**: Machine-readable format with all report data
4. **CSV**: Tabular format for importing into spreadsheets

### Enhanced Report Features

The enhanced report template includes:

- Modern blue color scheme with improved typography and layout
- Summary cards showing overall quality metrics
- Progress indicators for each quality dimension
- Tabbed navigation for different report sections
- Detailed tables of rule results with highlighting
- Visualizations of data quality metrics
- Responsive design that works well on different screen sizes

## `quality metrics`

Retrieves and analyzes quality metrics over time.

### Usage

```bash
nessi quality metrics <table-name> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table |

### Options

| Option | Description |
|--------|-------------|
| `--metric` | Specific metric to retrieve (completeness, accuracy, consistency, uniqueness, timeliness) |
| `--column` | Filter by column name |
| `--from` | Start date for metrics (YYYY-MM-DD) |
| `--to` | End date for metrics (YYYY-MM-DD) |
| `--output` | Output format (text, json, csv, html) |
| `--file` | Output file path |
| `--trend` | Show trend analysis |

### Examples

```bash
# Get all quality metrics
nessi quality metrics customers

# Get specific metric
nessi quality metrics customers --metric completeness

# Get metrics for a specific column
nessi quality metrics customers --column email

# Get metrics for a date range
nessi quality metrics customers --from 2025-04-01 --to 2025-05-01

# Get metrics with trend analysis
nessi quality metrics customers --trend
```

### Output

#### Text Format (Default)

```
Quality Metrics: customers
Period: 2025-04-01 to 2025-05-22

METRIC        CURRENT  7-DAY AVG  30-DAY AVG  TREND
completeness  0.998    0.997      0.995       ↑
accuracy      0.995    0.994      0.992       ↑
consistency   0.990    0.989      0.985       ↑
uniqueness    1.000    1.000      1.000       →
timeliness    0.985    0.980      0.975       ↑

Column Metrics:
COLUMN  METRIC        CURRENT  7-DAY AVG  30-DAY AVG  TREND
id      completeness  1.000    1.000      1.000       →
id      uniqueness    1.000    1.000      1.000       →
name    completeness  1.000    1.000      0.999       ↑
email   completeness  0.997    0.996      0.995       ↑
email   accuracy      0.995    0.994      0.992       ↑
age     completeness  0.998    0.997      0.996       ↑
age     accuracy      0.998    0.997      0.995       ↑
```

## Error Handling

Quality commands use the following error codes:

- `N800`: Validation rule not found
- `N801`: Invalid rule type
- `N802`: Invalid rule definition
- `N803`: Rule creation failed
- `N804`: Rule deletion failed
- `N805`: Quality check failed
- `N806`: Report generation failed
- `N807`: Invalid metric name
- `N808`: Metrics retrieval failed

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
