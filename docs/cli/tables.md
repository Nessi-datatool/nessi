# Table Commands

This document provides detailed information about Nessi's table commands for working with Delta Lake tables.

## `tables list`

Lists all available Delta Lake tables.

### Usage

```bash
nessi tables list [--path <path>] [--format <format>] [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--path` | Root path to search for tables |
| `--format` | Filter tables by format (delta, parquet, csv) |
| `--recursive` | Search recursively in subdirectories |
| `--output` | Output format (text, json, csv) |
| `--sort` | Sort by column (name, path, size, modified) |
| `--reverse` | Reverse sort order |

### Examples

```bash
# List all tables in the default location
nessi tables list

# List tables in a specific location
nessi tables list --path /data/warehouse

# List only Delta Lake tables
nessi tables list --format delta

# List tables with JSON output
nessi tables list --output json

# List tables sorted by modification time
nessi tables list --sort modified --reverse
```

### Output

#### Text Format (Default)

```
NAME             PATH                       FORMAT  SIZE     MODIFIED
customers        /data/customers            delta   1.2 GB   2025-05-01 12:34:56
orders           /data/orders               delta   3.5 GB   2025-05-02 10:11:12
products         /data/products             delta   500 MB   2025-04-30 09:08:07
```

#### JSON Format

```json
{
  "tables": [
    {
      "name": "customers",
      "path": "/data/customers",
      "format": "delta",
      "size": 1288490188,
      "modified": "2025-05-01T12:34:56Z"
    },
    {
      "name": "orders",
      "path": "/data/orders",
      "format": "delta",
      "size": 3758096384,
      "modified": "2025-05-02T10:11:12Z"
    },
    {
      "name": "products",
      "path": "/data/products",
      "format": "delta",
      "size": 524288000,
      "modified": "2025-04-30T09:08:07Z"
    }
  ]
}
```

## `tables describe`

Shows detailed information about a specific table.

### Usage

```bash
nessi tables describe <table-name> [--version <version>] [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table to describe |

### Options

| Option | Description |
|--------|-------------|
| `--version` | Specific version of the table to describe |
| `--output` | Output format (text, json, csv) |
| `--show-stats` | Include column statistics |
| `--show-partitions` | Include partition information |
| `--show-history` | Include version history |

### Examples

```bash
# Describe a table
nessi tables describe customers

# Describe a specific version of a table
nessi tables describe customers --version 3

# Describe a table with JSON output
nessi tables describe customers --output json

# Describe a table with all details
nessi tables describe customers --show-stats --show-partitions --show-history
```

### Output

#### Text Format (Default)

```
Table: customers
Path: /data/customers
Format: delta
Size: 1.2 GB
Rows: 1,000,000
Created: 2025-04-15 08:30:45
Last Modified: 2025-05-01 12:34:56
Current Version: 5

Schema:
  id: string (not null)
  name: string (not null)
  email: string
  age: integer
  signup_date: date
  last_login: timestamp

Partition Columns:
  signup_date

Statistics:
  id: 1,000,000 values, 0 nulls, 100% distinct
  name: 1,000,000 values, 0 nulls, 95% distinct
  email: 1,000,000 values, 2,500 nulls, 99% distinct
  age: 1,000,000 values, 1,200 nulls, min=18, max=95, avg=42.5
  signup_date: 1,000,000 values, 0 nulls, min=2020-01-01, max=2025-05-01
  last_login: 1,000,000 values, 150,000 nulls, min=2020-01-01T00:00:00Z, max=2025-05-01T12:34:56Z

Recent History:
  Version 5: 2025-05-01 12:34:56 - UPDATE (affected rows: 15,000)
  Version 4: 2025-04-25 09:12:34 - INSERT (rows: 50,000)
  Version 3: 2025-04-20 14:25:36 - DELETE (affected rows: 5,000)
```

#### JSON Format

```json
{
  "table": {
    "name": "customers",
    "path": "/data/customers",
    "format": "delta",
    "size_bytes": 1288490188,
    "row_count": 1000000,
    "created_at": "2025-04-15T08:30:45Z",
    "modified_at": "2025-05-01T12:34:56Z",
    "current_version": 5,
    "schema": [
      {
        "name": "id",
        "type": "string",
        "nullable": false
      },
      {
        "name": "name",
        "type": "string",
        "nullable": false
      },
      {
        "name": "email",
        "type": "string",
        "nullable": true
      },
      {
        "name": "age",
        "type": "integer",
        "nullable": true
      },
      {
        "name": "signup_date",
        "type": "date",
        "nullable": true
      },
      {
        "name": "last_login",
        "type": "timestamp",
        "nullable": true
      }
    ],
    "partition_columns": [
      "signup_date"
    ],
    "statistics": {
      "id": {
        "count": 1000000,
        "null_count": 0,
        "distinct_ratio": 1.0
      },
      "name": {
        "count": 1000000,
        "null_count": 0,
        "distinct_ratio": 0.95
      },
      "email": {
        "count": 1000000,
        "null_count": 2500,
        "distinct_ratio": 0.99
      },
      "age": {
        "count": 1000000,
        "null_count": 1200,
        "min": 18,
        "max": 95,
        "mean": 42.5
      },
      "signup_date": {
        "count": 1000000,
        "null_count": 0,
        "min": "2020-01-01",
        "max": "2025-05-01"
      },
      "last_login": {
        "count": 1000000,
        "null_count": 150000,
        "min": "2020-01-01T00:00:00Z",
        "max": "2025-05-01T12:34:56Z"
      }
    },
    "history": [
      {
        "version": 5,
        "timestamp": "2025-05-01T12:34:56Z",
        "operation": "UPDATE",
        "affected_rows": 15000
      },
      {
        "version": 4,
        "timestamp": "2025-04-25T09:12:34Z",
        "operation": "INSERT",
        "affected_rows": 50000
      },
      {
        "version": 3,
        "timestamp": "2025-04-20T14:25:36Z",
        "operation": "DELETE",
        "affected_rows": 5000
      }
    ]
  }
}
```

## `tables schema`

Shows the schema of a table.

### Usage

```bash
nessi tables schema <table-name> [--version <version>] [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table |

### Options

| Option | Description |
|--------|-------------|
| `--version` | Specific version of the table |
| `--output` | Output format (text, json, csv) |
| `--format` | Schema format (simple, detailed, avro, json) |
| `--show-metadata` | Include schema metadata |

### Examples

```bash
# Show table schema
nessi tables schema customers

# Show schema for a specific version
nessi tables schema customers --version 3

# Show schema in Avro format
nessi tables schema customers --format avro

# Show schema with metadata
nessi tables schema customers --show-metadata
```

### Output

#### Text Format (Default)

```
Table: customers
Schema:
  id: string (not null)
  name: string (not null)
  email: string
  age: integer
  signup_date: date
  last_login: timestamp
```

#### JSON Format with Avro Schema

```json
{
  "table": "customers",
  "schema": {
    "type": "record",
    "name": "customers",
    "fields": [
      {"name": "id", "type": {"type": "string", "nullable": false}},
      {"name": "name", "type": {"type": "string", "nullable": false}},
      {"name": "email", "type": ["null", "string"]},
      {"name": "age", "type": ["null", "int"]},
      {"name": "signup_date", "type": ["null", {"type": "int", "logicalType": "date"}]},
      {"name": "last_login", "type": ["null", {"type": "long", "logicalType": "timestamp-micros"}]}
    ]
  }
}
```

## `tables history`

Shows the version history of a table.

### Usage

```bash
nessi tables history <table-name> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table |

### Options

| Option | Description |
|--------|-------------|
| `--limit` | Maximum number of versions to show |
| `--output` | Output format (text, json, csv) |
| `--detailed` | Show detailed information for each version |
| `--from` | Show versions from this timestamp |
| `--to` | Show versions up to this timestamp |

### Examples

```bash
# Show table history
nessi tables history customers

# Show limited history
nessi tables history customers --limit 5

# Show detailed history
nessi tables history customers --detailed

# Show history within a time range
nessi tables history customers --from 2025-04-01 --to 2025-05-01
```

### Output

#### Text Format (Default)

```
Table: customers
Version History:

VERSION  TIMESTAMP            OPERATION  USER        AFFECTED ROWS
5        2025-05-01 12:34:56  UPDATE     john.doe    15,000
4        2025-04-25 09:12:34  INSERT     jane.smith  50,000
3        2025-04-20 14:25:36  DELETE     john.doe    5,000
2        2025-04-18 10:45:23  UPDATE     jane.smith  25,000
1        2025-04-15 08:30:45  CREATE     john.doe    930,000
```

#### Detailed Text Format

```
Table: customers
Version History:

Version: 5
Timestamp: 2025-05-01 12:34:56
Operation: UPDATE
User: john.doe
Affected Rows: 15,000
Commit Message: "Update customer email addresses"
Changes:
  - Modified column: email

Version: 4
Timestamp: 2025-04-25 09:12:34
Operation: INSERT
User: jane.smith
Affected Rows: 50,000
Commit Message: "Add new customers from marketing campaign"
Changes:
  - Added 50,000 rows

...
```

## `tables profile`

Generates a profile of a table's data.

### Usage

```bash
nessi tables profile <table-name> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table |

### Options

| Option | Description |
|--------|-------------|
| `--sample` | Sampling percentage (0-100) |
| `--columns` | Comma-separated list of columns to profile |
| `--output` | Output format (text, json, csv, html) |
| `--include-histograms` | Include histograms for numeric columns |
| `--include-top-values` | Include top values for string columns |
| `--file` | Output file path |

### Examples

```bash
# Generate table profile
nessi tables profile customers

# Generate profile with sampling
nessi tables profile customers --sample 50

# Profile specific columns
nessi tables profile customers --columns id,name,email

# Generate HTML report
nessi tables profile customers --output html --file customers_profile.html

# Include detailed statistics
nessi tables profile customers --include-histograms --include-top-values
```

### Output

#### Text Format (Default)

```
Table Profile: customers
Sample: 100% (1,000,000 rows)
Generated: 2025-05-22 22:12:09

Column: id (string)
  Count: 1,000,000
  Distinct: 1,000,000 (100.00%)
  Nulls: 0 (0.00%)
  Min Length: 36
  Max Length: 36
  Pattern: UUID format (100.00%)

Column: name (string)
  Count: 1,000,000
  Distinct: 950,000 (95.00%)
  Nulls: 0 (0.00%)
  Min Length: 3
  Max Length: 50
  Average Length: 22.5
  Top Values:
    "John Smith": 1,200 (0.12%)
    "Mary Johnson": 1,100 (0.11%)
    "Robert Williams": 950 (0.10%)

Column: email (string)
  Count: 1,000,000
  Distinct: 997,500 (99.75%)
  Nulls: 2,500 (0.25%)
  Min Length: 10
  Max Length: 50
  Average Length: 25.3
  Pattern: Contains "@" (99.95%)

Column: age (integer)
  Count: 1,000,000
  Distinct: 78
  Nulls: 1,200 (0.12%)
  Min: 18
  Max: 95
  Mean: 42.5
  Median: 41
  StdDev: 15.2
  Histogram:
    [18-25]: 150,000 (15.00%)
    [26-35]: 250,000 (25.00%)
    [36-45]: 300,000 (30.00%)
    [46-55]: 150,000 (15.00%)
    [56-65]: 100,000 (10.00%)
    [66+]: 50,000 (5.00%)

Column: signup_date (date)
  Count: 1,000,000
  Distinct: 1,825
  Nulls: 0 (0.00%)
  Min: 2020-01-01
  Max: 2025-05-01
  Distribution:
    2020: 200,000 (20.00%)
    2021: 250,000 (25.00%)
    2022: 250,000 (25.00%)
    2023: 200,000 (20.00%)
    2024-2025: 100,000 (10.00%)

Column: last_login (timestamp)
  Count: 1,000,000
  Distinct: 950,000
  Nulls: 150,000 (15.00%)
  Min: 2020-01-01T00:00:00Z
  Max: 2025-05-01T12:34:56Z
  Inactive (>30 days): 250,000 (25.00%)
```

## `tables diff`

Compares two versions of a table or two different tables.

### Usage

```bash
nessi tables diff <table1> <table2> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table1` | First table name or path |
| `table2` | Second table name or path |

### Options

| Option | Description |
|--------|-------------|
| `--version1` | Version of the first table |
| `--version2` | Version of the second table |
| `--output` | Output format (text, json, csv, html) |
| `--columns` | Comma-separated list of columns to compare |
| `--schema-only` | Compare only schema, not data |
| `--data-only` | Compare only data, not schema |
| `--sample` | Sampling percentage for data comparison (0-100) |
| `--file` | Output file path |

### Examples

```bash
# Compare two versions of the same table
nessi tables diff customers customers --version1 3 --version2 5

# Compare two different tables
nessi tables diff customers customers_backup

# Compare only schema
nessi tables diff customers customers_backup --schema-only

# Compare specific columns
nessi tables diff customers customers_backup --columns id,name,email

# Generate HTML report
nessi tables diff customers customers_backup --output html --file diff_report.html
```

### Output

#### Text Format (Default)

```
Table Diff Report
Table 1: customers (version 3)
Table 2: customers (version 5)
Generated: 2025-05-22 22:12:09

Schema Differences:
  - No differences in schema

Data Summary:
  - Table 1: 950,000 rows
  - Table 2: 1,000,000 rows
  - Added rows: 50,000
  - Modified rows: 15,000
  - Deleted rows: 0

Column: email
  - Modified in 15,000 rows (1.58%)
  - Example changes:
    Row ID: 12345
      Old: "john.doe@example.com"
      New: "john.doe@newdomain.com"
    Row ID: 67890
      Old: "jane.smith@example.com"
      New: "jane.smith@newdomain.com"
```

## `tables export`

Exports a table to various formats.

### Usage

```bash
nessi tables export <table-name> <output-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `table-name` | Name or path of the table |
| `output-path` | Path to export the table to |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Export format (parquet, csv, json, avro) |
| `--version` | Specific version of the table to export |
| `--columns` | Comma-separated list of columns to export |
| `--where` | Filter condition for rows to export |
| `--sample` | Sampling percentage (0-100) |
| `--partitioned` | Maintain partition structure in output |
| `--compression` | Compression type (none, snappy, gzip, etc.) |

### Examples

```bash
# Export table to CSV
nessi tables export customers /output/customers.csv --format csv

# Export specific version to Parquet
nessi tables export customers /output/customers --format parquet --version 3

# Export specific columns
nessi tables export customers /output/customers.csv --format csv --columns id,name,email

# Export with filtering
nessi tables export customers /output/active_customers.csv --format csv --where "last_login > '2025-04-01'"

# Export with sampling
nessi tables export customers /output/customers_sample.csv --format csv --sample 10
```

### Output

```
Exporting table: customers
Format: csv
Output: /output/customers.csv
Columns: All
Version: 5 (latest)

Progress: 100% [==================================================]
Exported 1,000,000 rows (1.2 GB)
Export completed in 45.2s
```

## Error Handling

Table commands use the following error codes:

- `N500`: Delta Lake table not found
- `N501`: Invalid Delta Lake table format
- `N502`: Table version not found
- `N503`: Schema validation error
- `N504`: Export format not supported
- `N505`: Invalid filter condition
- `N506`: Invalid column name

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
