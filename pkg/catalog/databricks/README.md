# Databricks Integration for Nessi CLI

This package provides integration with Databricks, allowing Nessi CLI to interact with Databricks catalogs, schemas, and tables. It supports reading and writing data to Delta Lake tables, with time travel capabilities.

## Features

- Connect to Databricks workspaces using personal access tokens
- Browse catalogs, schemas, and tables in Unity Catalog
- Read and write data to Delta Lake tables
- Support for time travel (read data as of a specific version or timestamp)
- Schema inference for Delta Lake tables

## Prerequisites

- Databricks workspace with Unity Catalog enabled
- Personal access token with appropriate permissions
- Environment variables set up for authentication

## Environment Variables

The following environment variables are required:

```
DATABRICKS_HOST=<your-databricks-host>
DATABRICKS_TOKEN=<your-personal-access-token>
DATABRICKS_WORKSPACE_ID=<your-workspace-id>
DATABRICKS_DEFAULT_CATALOG=<default-catalog-name>
DATABRICKS_DEFAULT_SCHEMA=<default-schema-name>
```

## Usage Examples

### Listing Catalogs

```bash
nessi catalog list --provider databricks
```

### Listing Schemas in a Catalog

```bash
nessi catalog list-schemas --provider databricks --catalog main
```

### Listing Tables in a Schema

```bash
nessi catalog list-tables --provider databricks --catalog main --schema default
```

### Reading Data from a Delta Table

```bash
nessi data read --provider databricks --catalog main --schema default --table customer_data
```

### Reading Data with Time Travel

```bash
# Read as of a specific version
nessi data read --provider databricks --catalog main --schema default --table customer_data --version 2

# Read as of a specific timestamp
nessi data read --provider databricks --catalog main --schema default --table customer_data --timestamp "2025-05-18T10:00:00Z"
```

### Writing Data to a Delta Table

```bash
nessi data write --provider databricks --catalog main --schema default --table customer_data --input-file data.csv
```

## Error Handling

The Databricks integration includes comprehensive error handling for various scenarios:

- API errors (4xx, 5xx status codes) with JSON error responses
- Network errors when connecting to Databricks
- JSON parsing errors for malformed responses
- Authentication errors with invalid tokens
- Resource not found errors for non-existent catalogs/schemas/tables
- Rate limiting errors when too many requests are made

## Development

### Running Tests

```bash
# Run all tests
go test -v ./pkg/catalog/databricks/...

# Run integration tests
go test -v ./pkg/catalog/databricks/... -run TestDatabricksIntegration

# Run with coverage
go test -coverprofile=coverage.out ./pkg/catalog/databricks/...
go tool cover -html=coverage.out
```

### Dependencies

This package depends on:

- `github.com/apache/arrow/go/v15` for Arrow data format handling
- `github.com/nessi-dev/nessi/pkg/datalake` for Delta Lake format handling

## Contributing

Contributions to improve the Databricks integration are welcome. Please ensure all tests pass before submitting a pull request.
