# Custom Rule Extensions

## Overview

Nessi.dev implements a flexible custom rule extension system that allows users to create and manage their own validation rules beyond the built-in rule types. This system enables organizations to implement domain-specific validation logic and share rule libraries across teams.

## Core Components

### Rule Registry

The Rule Registry is responsible for:

- Registering custom rule types
- Managing rule lifecycle (creation, updating, deletion)
- Providing access to rule definitions
- Validating rule configurations

### Rule Library Manager

The Rule Library Manager handles:

- Organizing rules into reusable libraries
- Sharing rule libraries between teams
- Versioning rule libraries
- Importing and exporting rule libraries

## Usage

### CLI Commands

```bash
# List all custom rules
nessi rules list

# Create a new custom rule
nessi rules create --name "date_format_check" --type "regex" --pattern "^\d{4}-\d{2}-\d{2}$" --description "Validates ISO date format"

# Get details of a specific rule
nessi rules get date_format_check

# Update an existing rule
nessi rules update date_format_check --pattern "^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$"

# Delete a rule
nessi rules delete date_format_check

# List rule libraries
nessi rules libraries list

# Create a rule library
nessi rules libraries create --name "date_validation" --description "Date validation rules"

# Add a rule to a library
nessi rules libraries add --library date_validation --rule date_format_check

# Share a library with a team
nessi rules libraries share --library date_validation --team "Data Quality Team"

# Export a rule library
nessi rules libraries export date_validation --output date_validation_lib.json

# Import a rule library
nessi rules libraries import date_validation_lib.json
```

### API Endpoints

The custom rule extensions system exposes the following API endpoints:

- `GET /api/v1/rules` - List all custom rules
- `POST /api/v1/rules` - Create a new custom rule
- `GET /api/v1/rules/{name}` - Get a specific rule
- `PUT /api/v1/rules/{name}` - Update a rule
- `DELETE /api/v1/rules/{name}` - Delete a rule
- `GET /api/v1/rules/libraries` - List rule libraries
- `POST /api/v1/rules/libraries` - Create a rule library
- `PUT /api/v1/rules/libraries/{name}/rules` - Add a rule to a library
- `POST /api/v1/rules/libraries/export/{name}` - Export a rule library
- `POST /api/v1/rules/libraries/import` - Import a rule library

## Rule Types

The system supports the following rule types:

### Built-in Rule Types

- **regex**: Validates values against a regular expression
- **range**: Validates numeric values within a range
- **enum**: Validates values against a set of allowed values
- **length**: Validates string length
- **date_format**: Validates date format
- **null_check**: Validates presence/absence of null values
- **uniqueness**: Validates uniqueness of values
- **referential**: Validates references to other data

### Custom Rule Types

Users can create custom rule types by implementing the `RuleValidator` interface:

```go
type RuleValidator interface {
    Validate(value interface{}, params map[string]interface{}) (bool, error)
    GetName() string
    GetDescription() string
    GetParamSchema() map[string]ParamDefinition
}
```

## Rule Libraries

Rule libraries provide a way to organize and share rules:

- **Organization**: Group related rules together
- **Versioning**: Track changes to rules over time
- **Sharing**: Share rules between teams
- **Reuse**: Import rules from other libraries

## Implementation Details

### Rule Definition

A rule definition includes:

- **Name**: Unique identifier for the rule
- **Type**: The type of rule (regex, range, etc.)
- **Parameters**: Configuration parameters for the rule
- **Description**: Human-readable description of the rule
- **Created By**: The user who created the rule
- **Created At**: When the rule was created
- **Updated At**: When the rule was last updated

### Rule Execution

Rules are executed by the rule engine, which:

1. Loads the rule definition
2. Retrieves the appropriate validator
3. Applies the rule to the data
4. Returns validation results

### Rule Storage

Rules can be stored in:

- **File System**: JSON or YAML files
- **Database**: Relational or NoSQL database
- **Remote Repository**: Git or other version control system

## Integration with Other Systems

### RBAC Integration

Access to rules and rule libraries is controlled by the RBAC system, ensuring that only authorized users can create, modify, or delete rules.

### Audit Logging Integration

Rule creation, modification, and deletion are logged in the audit system, providing accountability and traceability.

### Webhook Integration

Rule events (creation, modification, deletion) can trigger webhooks to notify external systems.

## Best Practices

1. **Documentation**: Document the purpose and expected behavior of each rule
2. **Testing**: Test rules with both valid and invalid data
3. **Versioning**: Version rule libraries to track changes over time
4. **Reuse**: Create generic rules that can be reused across different datasets
5. **Performance**: Consider the performance impact of complex rules on large datasets

## Examples

### Creating a Custom Date Format Rule

```bash
nessi rules create --name "iso_date_check" --type "regex" --pattern "^\d{4}-\d{2}-\d{2}$" --description "Validates ISO date format (YYYY-MM-DD)"
```

### Creating a Rule Library for Financial Data

```bash
# Create the library
nessi rules libraries create --name "financial_validation" --description "Validation rules for financial data"

# Add rules to the library
nessi rules libraries add --library financial_validation --rule positive_amount
nessi rules libraries add --library financial_validation --rule valid_currency_code
nessi rules libraries add --library financial_validation --rule balanced_transaction

# Share the library with the finance team
nessi rules libraries share --library financial_validation --team "Finance Team"
```

## Troubleshooting

### Common Issues

1. **Rule Not Found**: Ensure the rule name is correct and the rule exists
2. **Invalid Rule Configuration**: Verify that the rule parameters are correctly configured
3. **Permission Denied**: Check that the user has the necessary permissions to access or modify the rule

### Debugging

Use the system logs to troubleshoot issues with the custom rule extensions:

```bash
nessi system logs --component rules
```

### Rule Testing

Test rules before applying them to production data:

```bash
nessi rules test iso_date_check --value "2025-01-01"
```
