# Data Catalog Integration

This document describes how to use Nessi.dev's Data Catalog Integration feature to connect with popular data catalogs like AWS Glue, Azure Purview, and Google Cloud Data Catalog.

## Overview

Nessi.dev's Data Catalog Integration enables seamless connections to popular data catalogs, allowing you to:

1. **Discover data assets** in your existing data catalogs
2. **Leverage existing metadata** for data profiling and validation
3. **Publish data quality metrics** back to your data catalog
4. **Share lineage information** between Nessi.dev and your data catalog
5. **Enhance collaboration** between data teams

## Supported Data Catalogs

Nessi.dev currently supports the following data catalogs:

- **AWS Glue Data Catalog**
- **Azure Purview**
- **Google Cloud Data Catalog**

Support for the following catalogs is planned for future releases:

- Apache Atlas
- Collibra

## Prerequisites

Before using the Data Catalog Integration feature, ensure you have:

1. Appropriate access credentials for your data catalog
2. Necessary permissions to read and write metadata
3. Nessi.dev installed and configured

## Configuration

### AWS Glue Data Catalog

To connect to AWS Glue Data Catalog, you need to provide the following configuration:

```json
{
  "region": "us-east-1",
  "access_key": "YOUR_ACCESS_KEY",
  "secret_key": "YOUR_SECRET_KEY",
  "use_iam_role": false
}
```

If you're using IAM roles (recommended for production), you can omit the access_key and secret_key:

```json
{
  "region": "us-east-1",
  "use_iam_role": true
}
```

### Azure Purview

To connect to Azure Purview, you need to provide the following configuration:

```json
{
  "account_name": "your-purview-account",
  "use_azure_ad": true
}
```

Azure Purview requires Azure AD authentication. Make sure you have logged in using the Azure CLI or have the appropriate environment variables set.

### Google Cloud Data Catalog

To connect to Google Cloud Data Catalog, you need to provide the following configuration:

```json
{
  "project_id": "your-gcp-project",
  "location": "us-central1",
  "credentials_file": "/path/to/credentials.json"
}
```

If you're using application default credentials, you can omit the credentials_file:

```json
{
  "project_id": "your-gcp-project",
  "location": "us-central1"
}
```

## Using the CLI

Nessi.dev provides a set of CLI commands to interact with data catalogs.

### Connect to a Data Catalog

```bash
nessi catalog connect aws_glue config.json
```

### List Databases

```bash
nessi catalog list-databases aws_glue
```

### List Tables in a Database

```bash
nessi catalog list-tables aws_glue my_database
```

### Get Table Details

```bash
nessi catalog get-table aws_glue my_database my_table
```

### Publish Quality Metrics

```bash
nessi catalog publish-quality aws_glue my_database my_table profile.json validation.json
```

## Programmatic Usage

You can also use the Data Catalog Integration programmatically in your Go code:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nessi-dev/nessi-dev/pkg/catalog"
)

func main() {
	// Create catalog factory
	factory := catalog.NewCatalogFactory()

	// Create catalog manager
	manager := catalog.NewCatalogManager()

	// Create AWS Glue catalog
	awsGlue, err := factory.CreateCatalog(catalog.AWSGlue)
	if err != nil {
		log.Fatalf("Failed to create AWS Glue catalog: %v", err)
	}

	// Register catalog
	manager.RegisterCatalog(awsGlue)

	// Connect to catalog
	config := map[string]interface{}{
		"region":      "us-east-1",
		"use_iam_role": true,
	}

	ctx := context.Background()
	if err := awsGlue.Connect(ctx, config); err != nil {
		log.Fatalf("Failed to connect to AWS Glue: %v", err)
	}

	// List databases
	databases, err := awsGlue.ListDatabases(ctx)
	if err != nil {
		log.Fatalf("Failed to list databases: %v", err)
	}

	// Print databases
	fmt.Println("Databases in AWS Glue:")
	for _, db := range databases {
		fmt.Printf("- %s: %s\n", db.Name, db.Description)
	}
}
```

## Publishing Data Quality Metrics

One of the key features of the Data Catalog Integration is the ability to publish data quality metrics back to your data catalog. This allows you to:

1. Track data quality over time
2. Share quality information with other teams
3. Generate comprehensive data quality reports

To publish quality metrics, you need to:

1. Run a data profile and validation using Nessi.dev
2. Use the `publish-quality` command to send the metrics to your data catalog

Example:

```bash
# First, profile your data
nessi profile --table my_table --output profile.json

# Then, validate it against rules
nessi validate --table my_table --rules rules.yaml --output validation.json

# Finally, publish the quality metrics
nessi catalog publish-quality aws_glue my_database my_table profile.json validation.json
```

## Quality Metrics

The following quality metrics are published to your data catalog:

- **Overall Score**: A weighted average of all quality dimensions
- **Completeness**: Measures the presence of null values
- **Accuracy**: Measures how well the data conforms to expected formats and ranges
- **Consistency**: Measures the uniformity of data across records
- **Timeliness**: Measures how recent the data is

Each metric is a score between 0.0 and 1.0, where 1.0 represents perfect quality.

## Integration with Data Validation

The Data Catalog Integration works seamlessly with Nessi.dev's data validation features. You can:

1. Discover tables in your data catalog
2. Generate validation rules based on catalog metadata
3. Run validations against the tables
4. Publish the validation results back to the catalog

This creates a continuous feedback loop that helps maintain and improve data quality over time.

## Best Practices

1. **Use IAM roles or managed identities** instead of access keys for authentication
2. **Publish quality metrics regularly** to track trends over time
3. **Use tags in your data catalog** to categorize and filter data assets
4. **Leverage existing metadata** to generate validation rules
5. **Share quality reports** with stakeholders to promote data quality awareness

## Troubleshooting

### Common Issues

1. **Connection failures**: Ensure your credentials are correct and you have the necessary permissions
2. **Missing tables**: Check if you have access to the database and table in the data catalog
3. **Failed to publish metrics**: Verify that you have write permissions in the data catalog

### Logging

To enable debug logging for the Data Catalog Integration, set the following environment variable:

```bash
export NESSI_LOG_LEVEL=debug
```

## Next Steps

- Learn how to [create custom validation rules](./validation_rules.md)
- Explore [data profiling capabilities](./data_profiling.md)
- Set up [automated quality monitoring](./monitoring.md)
