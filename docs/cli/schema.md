# Schema Commands

This document provides detailed information about Nessi's schema commands for managing, validating, and evolving data schemas.

## `schema show`

Shows the schema of a table or file.

### Usage

```bash
nessi schema show <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table or file |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Source format (delta, parquet, csv, json, auto) |
| `--output` | Output format (text, json, avro, ddl, html) |
| `--file` | Output file path |
| `--include-metadata` | Include schema metadata |
| `--version` | Specific version of the table (for Delta tables) |

### Examples

```bash
# Show schema of a Delta table
nessi schema show /data/customers

# Show schema in JSON format
nessi schema show /data/customers --output json

# Show schema with metadata
nessi schema show /data/customers --include-metadata

# Show schema of a specific version
nessi schema show /data/customers --version 3

# Show schema of a CSV file
nessi schema show /data/exports/customers.csv --format csv

# Export schema to a file
nessi schema show /data/customers --output avro --file customers_schema.avsc
```

### Output

#### Text Format (Default)

```
Schema: /data/customers
Format: delta
Version: 5 (latest)

COLUMN       TYPE        NULLABLE  DESCRIPTION
id           string      false     Unique customer identifier
name         string      false     Customer's full name
email        string      true      Customer's email address
age          integer     true      Customer's age in years
signup_date  date        true      Date when customer signed up
last_login   timestamp   true      Timestamp of customer's last login

Metadata:
  Created: 2025-04-15 08:30:45
  Last Modified: 2025-05-01 12:34:56
  Primary Key: id
  Partition Columns: signup_date
```

#### JSON Format

```json
{
  "name": "customers",
  "fields": [
    {
      "name": "id",
      "type": "string",
      "nullable": false,
      "metadata": {
        "description": "Unique customer identifier"
      }
    },
    {
      "name": "name",
      "type": "string",
      "nullable": false,
      "metadata": {
        "description": "Customer's full name"
      }
    },
    {
      "name": "email",
      "type": "string",
      "nullable": true,
      "metadata": {
        "description": "Customer's email address"
      }
    },
    {
      "name": "age",
      "type": "integer",
      "nullable": true,
      "metadata": {
        "description": "Customer's age in years"
      }
    },
    {
      "name": "signup_date",
      "type": "date",
      "nullable": true,
      "metadata": {
        "description": "Date when customer signed up"
      }
    },
    {
      "name": "last_login",
      "type": "timestamp",
      "nullable": true,
      "metadata": {
        "description": "Timestamp of customer's last login"
      }
    }
  ],
  "metadata": {
    "created": "2025-04-15T08:30:45Z",
    "modified": "2025-05-01T12:34:56Z",
    "primaryKey": ["id"],
    "partitionColumns": ["signup_date"]
  }
}
```

## `schema validate`

Validates data against a schema.

### Usage

```bash
nessi schema validate <source> [schema-file] [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table or file to validate |
| `schema-file` | Optional path to schema file (if not provided, infers from source) |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Source format (delta, parquet, csv, json, auto) |
| `--schema-format` | Schema format (json, avro, ddl) |
| `--output` | Output format (text, json, html, pdf) |
| `--file` | Output file path |
| `--sample` | Sampling percentage (0-100) |
| `--strict` | Enable strict validation mode |
| `--ignore-columns` | Comma-separated list of columns to ignore |

### Examples

```bash
# Validate a table against its own schema
nessi schema validate /data/customers

# Validate against a specific schema file
nessi schema validate /data/customers /schemas/customers_schema.json

# Validate with sampling
nessi schema validate /data/customers --sample 50

# Validate in strict mode
nessi schema validate /data/customers --strict

# Generate validation report
nessi schema validate /data/customers --output html --file validation_report.html

# Ignore specific columns
nessi schema validate /data/customers --ignore-columns created_at,updated_at
```

### Output

#### Text Format (Default)

```
Schema Validation Results: /data/customers
Timestamp: 2025-05-22 22:12:09
Sample: 100% (1,000,000 rows)
Mode: standard

Overall Result: PASSED
  Valid Rows: 998,500 (99.85%)
  Invalid Rows: 1,500 (0.15%)
  Validation Score: 0.9985

Column Validation:
  id (string, not null): PASSED
    Valid: 1,000,000 (100.00%)
    Invalid: 0 (0.00%)
    Issues: None

  name (string, not null): PASSED
    Valid: 1,000,000 (100.00%)
    Invalid: 0 (0.00%)
    Issues: None

  email (string): PASSED
    Valid: 997,500 (99.75%)
    Invalid: 0 (0.00%)
    Null: 2,500 (0.25%)
    Issues: None

  age (integer): FAILED
    Valid: 998,800 (99.88%)
    Invalid: 200 (0.02%)
    Null: 1,000 (0.10%)
    Issues:
      - Type mismatch: 200 values (expected integer, found string)
      - Example rows: 12345, 67890

  signup_date (date): PASSED
    Valid: 1,000,000 (100.00%)
    Invalid: 0 (0.00%)
    Issues: None

  last_login (timestamp): PASSED
    Valid: 850,000 (85.00%)
    Invalid: 0 (0.00%)
    Null: 150,000 (15.00%)
    Issues: None

Summary:
  Columns Validated: 6
  Columns Passed: 5
  Columns Failed: 1
  Invalid Rows: 200 (0.02%)
```

## `schema diff`

Compares two schemas and shows the differences.

### Usage

```bash
nessi schema diff <source1> <source2> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source1` | Path to the first table, file, or schema file |
| `source2` | Path to the second table, file, or schema file |

### Options

| Option | Description |
|--------|-------------|
| `--format1` | Format of the first source (delta, parquet, csv, json, schema, auto) |
| `--format2` | Format of the second source (delta, parquet, csv, json, schema, auto) |
| `--output` | Output format (text, json, html) |
| `--file` | Output file path |
| `--include-metadata` | Include metadata differences |
| `--version1` | Version of the first table (for Delta tables) |
| `--version2` | Version of the second table (for Delta tables) |

### Examples

```bash
# Compare schemas of two tables
nessi schema diff /data/customers /data/customers_backup

# Compare schemas of different versions
nessi schema diff /data/customers /data/customers --version1 3 --version2 5

# Compare schema files
nessi schema diff /schemas/customers_v1.json /schemas/customers_v2.json --format1 schema --format2 schema

# Generate HTML diff report
nessi schema diff /data/customers /data/customers_backup --output html --file schema_diff.html

# Include metadata differences
nessi schema diff /data/customers /data/customers_backup --include-metadata
```

### Output

#### Text Format (Default)

```
Schema Diff Results
Schema 1: /data/customers (version 3)
Schema 2: /data/customers (version 5)
Timestamp: 2025-05-22 22:12:09

Summary:
  Added Columns: 1
  Removed Columns: 0
  Modified Columns: 1
  Unchanged Columns: 4
  Compatibility: BACKWARD

Added Columns:
  + subscription_tier (string, nullable)
      Description: Customer's subscription tier

Modified Columns:
  ~ email (string, nullable) → email (string, not null)
      Nullability: true → false

Unchanged Columns:
  = id (string, not null)
  = name (string, not null)
  = age (integer, nullable)
  = signup_date (date, nullable)
  = last_login (timestamp, nullable)

Metadata Differences:
  ~ Primary Key: [id] → [id, email]
  = Partition Columns: [signup_date]
```

## `schema evolve`

Evolves a schema with new columns or changes.

### Usage

```bash
nessi schema evolve <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table or schema file to evolve |

### Options

| Option | Description |
|--------|-------------|
| `--add-column` | Add a column in format "name:type[:nullable][:description]" |
| `--drop-column` | Drop a column by name |
| `--rename-column` | Rename a column in format "old_name:new_name" |
| `--modify-column` | Modify a column in format "name:type[:nullable][:description]" |
| `--format` | Source format (delta, parquet, schema) |
| `--output` | Output format (text, json, avro, ddl) |
| `--file` | Output file path |
| `--apply` | Apply changes to the source (for Delta tables) |
| `--dry-run` | Show changes without applying them |

### Examples

```bash
# Add a new column
nessi schema evolve /data/customers --add-column "subscription_tier:string:true:Customer's subscription tier" --dry-run

# Drop a column
nessi schema evolve /data/customers --drop-column "age" --dry-run

# Rename a column
nessi schema evolve /data/customers --rename-column "email:contact_email" --dry-run

# Modify a column
nessi schema evolve /data/customers --modify-column "email:string:false:Customer's verified email address" --dry-run

# Apply multiple changes
nessi schema evolve /data/customers --add-column "subscription_tier:string:true:Customer's subscription tier" --modify-column "email:string:false" --apply

# Export evolved schema
nessi schema evolve /data/customers --add-column "subscription_tier:string:true" --output json --file evolved_schema.json
```

### Output

#### Text Format (Default)

```
Schema Evolution Plan: /data/customers
Timestamp: 2025-05-22 22:12:09
Mode: dry-run

Current Schema:
  id: string (not null)
  name: string (not null)
  email: string
  age: integer
  signup_date: date
  last_login: timestamp

Planned Changes:
  + ADD COLUMN subscription_tier: string
      Description: Customer's subscription tier
  ~ MODIFY COLUMN email: string → string (not null)
      Description: Customer's verified email address

Evolved Schema:
  id: string (not null)
  name: string (not null)
  email: string (not null)
      Description: Customer's verified email address
  age: integer
  signup_date: date
  last_login: timestamp
  subscription_tier: string
      Description: Customer's subscription tier

Compatibility Analysis:
  Backward Compatible: YES
  Forward Compatible: NO
  Full Compatible: NO

Impact Assessment:
  Rows Requiring Updates: 2,500 (0.25%)
  Estimated Time: 5.2s
  Requires Backfill: NO
```

## `schema tree`

Visualizes schema as a tree structure.

### Usage

```bash
nessi schema tree <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the table, file, or schema file |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Source format (delta, parquet, csv, json, schema, auto) |
| `--output` | Output format (text, json, html, png) |
| `--file` | Output file path |
| `--include-metadata` | Include schema metadata |
| `--max-depth` | Maximum depth for nested structures |
| `--theme` | Visual theme (default, dark, light, colorful) |

### Examples

```bash
# Visualize schema as a tree
nessi schema tree /data/customers

# Export tree visualization
nessi schema tree /data/customers --output png --file schema_tree.png

# Visualize with a specific theme
nessi schema tree /data/customers --theme colorful

# Limit depth for nested structures
nessi schema tree /data/nested_data --max-depth 3

# Include metadata in visualization
nessi schema tree /data/customers --include-metadata
```

### Output

#### Text Format (Default)

```
Schema Tree: /data/customers
Format: delta
Version: 5 (latest)

customers
├── id: string (not null)
│   └── description: Unique customer identifier
├── name: string (not null)
│   └── description: Customer's full name
├── email: string
│   └── description: Customer's email address
├── age: integer
│   └── description: Customer's age in years
├── signup_date: date
│   └── description: Date when customer signed up
└── last_login: timestamp
    └── description: Timestamp of customer's last login

Metadata:
├── Primary Key: id
└── Partition Columns: signup_date
```

## `schema infer`

Infers schema from data.

### Usage

```bash
nessi schema infer <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the file or data to infer schema from |

### Options

| Option | Description |
|--------|-------------|
| `--format` | Source format (csv, json, parquet, auto) |
| `--output` | Output format (text, json, avro, ddl) |
| `--file` | Output file path |
| `--sample` | Sampling percentage (0-100) |
| `--include-metadata` | Include inferred metadata |
| `--strict-types` | Use strict type inference |
| `--header` | Whether the file has a header row (for CSV) |
| `--delimiter` | Delimiter character (for CSV) |

### Examples

```bash
# Infer schema from a CSV file
nessi schema infer /data/exports/customers.csv --format csv

# Infer schema with sampling
nessi schema infer /data/exports/customers.csv --format csv --sample 50

# Export inferred schema
nessi schema infer /data/exports/customers.csv --format csv --output json --file inferred_schema.json

# Infer schema with strict type inference
nessi schema infer /data/exports/customers.csv --format csv --strict-types

# Specify CSV options
nessi schema infer /data/exports/customers.csv --format csv --header true --delimiter ","
```

### Output

#### Text Format (Default)

```
Inferred Schema: /data/exports/customers.csv
Format: csv
Sample: 100% (1,000,000 rows)
Timestamp: 2025-05-22 22:12:09

COLUMN       INFERRED TYPE  CONFIDENCE  NULL %  DESCRIPTION
id           string         100.00%     0.00%   Unique identifier (UUID pattern)
name         string         100.00%     0.00%   Full name (word pattern)
email        string         100.00%     0.25%   Email address (email pattern)
age          integer        99.98%      0.10%   Numeric value (range: 18-95)
signup_date  date           100.00%     0.00%   Date value (ISO format)
last_login   timestamp      100.00%     15.00%  Timestamp value (ISO format)

Inferred Metadata:
  Primary Key Candidates: id
  Foreign Key Candidates: None
  Partition Column Candidates: signup_date
```

## `schema registry`

Manages schemas in a schema registry.

### Usage

```bash
nessi schema registry [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `list` | List schemas in the registry |
| `get` | Get a specific schema |
| `register` | Register a new schema |
| `update` | Update an existing schema |
| `delete` | Delete a schema |
| `validate` | Validate data against a registered schema |
| `versions` | List versions of a schema |
| `compatibility` | Check compatibility between schema versions |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Filter by subject |
| `--output` | Output format (text, json) |

### Options for `get`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Subject name |
| `--version` | Schema version (latest, all, or specific version) |
| `--output` | Output format (text, json, avro) |
| `--file` | Output file path |

### Options for `register` and `update`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Subject name |
| `--schema` | Schema file path or inline schema |
| `--schema-type` | Schema type (avro, json, protobuf) |
| `--compatibility` | Compatibility level (backward, forward, full, none) |
| `--description` | Schema description |

### Options for `delete`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Subject name |
| `--version` | Schema version (latest, all, or specific version) |
| `--force` | Force deletion without confirmation |

### Options for `validate`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Subject name |
| `--version` | Schema version (latest or specific version) |
| `--data` | Data file path to validate |
| `--format` | Data format (json, avro, csv) |
| `--output` | Output format (text, json) |

### Options for `versions`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Subject name |
| `--output` | Output format (text, json) |

### Options for `compatibility`

| Option | Description |
|--------|-------------|
| `--registry` | Registry URL or name |
| `--subject` | Subject name |
| `--schema1` | First schema version or file |
| `--schema2` | Second schema version or file |
| `--type` | Compatibility type (backward, forward, full) |
| `--output` | Output format (text, json) |

### Examples

```bash
# List schemas in registry
nessi schema registry list --registry http://schema-registry:8081

# Get a specific schema
nessi schema registry get --registry http://schema-registry:8081 --subject customers --version latest

# Register a new schema
nessi schema registry register --registry http://schema-registry:8081 --subject customers --schema /schemas/customers.avsc --schema-type avro --compatibility backward

# Update an existing schema
nessi schema registry update --registry http://schema-registry:8081 --subject customers --schema /schemas/customers_v2.avsc --schema-type avro

# Delete a schema
nessi schema registry delete --registry http://schema-registry:8081 --subject customers --version latest

# Validate data against a registered schema
nessi schema registry validate --registry http://schema-registry:8081 --subject customers --version latest --data /data/exports/customers.json --format json

# List versions of a schema
nessi schema registry versions --registry http://schema-registry:8081 --subject customers

# Check compatibility between schema versions
nessi schema registry compatibility --registry http://schema-registry:8081 --subject customers --schema1 1 --schema2 2 --type backward
```

### Output for `list`

```
Schema Registry: http://schema-registry:8081
Subjects:

SUBJECT     VERSIONS  LATEST VERSION  COMPATIBILITY  LAST UPDATED
customers   5         5               BACKWARD       2025-05-01 12:34:56
orders      3         3               BACKWARD       2025-04-15 09:23:45
products    2         2               FULL           2025-03-20 14:56:23
```

### Output for `get`

```
Schema: customers (version 5)
Registry: http://schema-registry:8081
Type: AVRO
Compatibility: BACKWARD
Created: 2025-05-01 12:34:56

{
  "type": "record",
  "name": "customers",
  "fields": [
    {"name": "id", "type": {"type": "string", "nullable": false}},
    {"name": "name", "type": {"type": "string", "nullable": false}},
    {"name": "email", "type": ["null", "string"]},
    {"name": "age", "type": ["null", "int"]},
    {"name": "signup_date", "type": ["null", {"type": "int", "logicalType": "date"}]},
    {"name": "last_login", "type": ["null", {"type": "long", "logicalType": "timestamp-micros"}]},
    {"name": "subscription_tier", "type": ["null", "string"]}
  ]
}
```

## Error Handling

Schema commands use the following error codes:

- `N700`: Schema not found
- `N701`: Invalid schema format
- `N702`: Schema validation failed
- `N703`: Schema evolution failed
- `N704`: Schema registry connection failed
- `N705`: Schema compatibility check failed
- `N706`: Schema inference failed
- `N707`: Schema registration failed
- `N708`: Schema version not found

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
