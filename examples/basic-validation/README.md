# Basic Data Quality Validation Example

This example demonstrates how to use Nessi for basic data quality validation on a Delta Lake table.

## Overview

In this example, we will:

1. Create a sample Delta Lake table
2. Define data quality rules
3. Run validation using Nessi
4. Generate an HTML report

## Prerequisites

- Nessi installed (see [Getting Started Guide](../../docs/GETTING_STARTED.md))
- Access to a Delta Lake table (we'll create one in this example)

## Step 1: Create a Sample Delta Table

Use the provided script to create a sample Delta table:

```bash
# Run the sample data generation script
./generate_sample_data.sh
```

This will create a Delta table at `./data/sample_table` with the following schema:

```
root
 |-- user_id: string
 |-- name: string
 |-- age: integer
 |-- email: string
 |-- signup_date: date
 |-- last_login: timestamp
 |-- active: boolean
 |-- score: double
```

## Step 2: Define Data Quality Rules

Create a file named `rules.yaml` with the following content:

```yaml
rules:
  - name: "user_id_not_null"
    description: "User ID should not be null"
    column: "user_id"
    type: "not_null"
    
  - name: "valid_email"
    description: "Email should be in valid format"
    column: "email"
    type: "regex"
    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    
  - name: "age_range"
    description: "Age should be between 18 and 120"
    column: "age"
    type: "range"
    min: 18
    max: 120
    
  - name: "score_range"
    description: "Score should be between 0 and 100"
    column: "score"
    type: "range"
    min: 0
    max: 100
    
  - name: "recent_login"
    description: "Last login should be within the last year"
    column: "last_login"
    type: "recency"
    max_age: "365d"
```

## Step 3: Run Validation

Run the validation using Nessi:

```bash
nessi validate ./data/sample_table --rules ./rules.yaml
```

This will output the validation results to the console.

## Step 4: Generate a Report

Generate an HTML report of the validation results:

```bash
nessi report ./data/sample_table --rules ./rules.yaml --format html --output ./reports
```

This will create an HTML report in the `./reports` directory.

## Step 5: Explore Additional Features

Try these additional commands:

```bash
# Get table profile
nessi profile ./data/sample_table

# View schema
nessi schema info ./data/sample_table

# Check table health
nessi health ./data/sample_table
```

## Next Steps

- Try modifying the rules to see how validation results change
- Create your own Delta table and apply these validation techniques
- Explore more advanced validation options in the [documentation](../../docs/)
