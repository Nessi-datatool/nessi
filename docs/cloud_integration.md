# Cloud Integration Guide

Nessi.dev provides seamless integration with major cloud providers (AWS, Azure, and GCP) to enable monitoring and validation of Delta Lake tables stored in cloud environments.

## Overview

The cloud integration features allow you to:

1. Connect to cloud storage services (S3, Azure Blob Storage, Google Cloud Storage)
2. Access and analyze Delta Lake tables stored in the cloud
3. Perform time travel operations on cloud-hosted Delta Lake tables
4. Track schema evolution and changes across versions
5. Validate data quality in cloud environments

## Supported Cloud Providers

Nessi.dev supports the following cloud providers:

- **AWS**: Amazon S3 for storage, AWS Glue Data Catalog
- **Azure**: Azure Blob Storage, Azure Purview
- **GCP**: Google Cloud Storage, Google Cloud Data Catalog

## Data Catalog Integration

In addition to cloud storage integration, Nessi.dev also provides seamless integration with popular data catalogs:

- **AWS Glue Data Catalog**: Discover tables, leverage metadata, and publish quality metrics
- **Azure Purview**: Connect to Azure's unified data governance service for metadata management
- **Google Cloud Data Catalog**: Integrate with GCP's fully managed, scalable metadata management service

See the [Data Catalog Integration Guide](./data_catalog_integration.md) for more details.

## Configuration

### Using the CLI

The `nessi cloud` command provides tools for managing cloud provider connections:

```bash
# Configure a cloud provider interactively
nessi cloud configure --name aws-prod --provider aws

# Configure using a configuration file
nessi cloud configure --config-file aws-config.yaml --name aws-prod

# List configured cloud providers
nessi cloud list

# Test a cloud provider connection
nessi cloud test --name aws-prod

# Remove a cloud provider configuration
nessi cloud remove --name aws-prod
```

### Configuration File Format

You can define cloud provider configurations in YAML or JSON format:

```yaml
# AWS Example
provider: aws
region: us-west-2
credentials:
  access_key: YOUR_ACCESS_KEY
  secret_key: YOUR_SECRET_KEY
  use_iam_role: false
default_bucket: my-delta-bucket
endpoint_override: ""  # Optional, for S3-compatible services

# Azure Example
provider: azure
credentials:
  account_name: mystorageaccount
  # Authentication options (choose one):
  # Option 1: Account Key
  account_key: YOUR_ACCOUNT_KEY
  # Option 2: SAS Token
  # sas_token: YOUR_SAS_TOKEN
  # Option 3: Azure AD (recommended for production)
  # use_azure_ad: true
default_bucket: my-container
endpoint_override: ""  # Optional, for custom endpoints

# GCP Example
provider: gcp
credentials:
  credentials_file: /path/to/credentials.json
additional_options:
  project_id: my-gcp-project
default_bucket: my-gcs-bucket
```

### Environment Variables

You can also configure cloud providers using environment variables:

```bash
# AWS
export NESSI_AWS_ACCESS_KEY=YOUR_ACCESS_KEY
export NESSI_AWS_SECRET_KEY=YOUR_SECRET_KEY
export NESSI_AWS_REGION=us-west-2
export NESSI_AWS_BUCKET=my-delta-bucket

# Azure
export NESSI_AZURE_ACCOUNT_NAME=mystorageaccount
export NESSI_AZURE_ACCOUNT_KEY=YOUR_ACCOUNT_KEY
export NESSI_AZURE_CONTAINER=my-container

# GCP
export NESSI_GCP_PROJECT_ID=my-gcp-project
export NESSI_GCP_CREDENTIALS_FILE=/path/to/credentials.json
export NESSI_GCP_BUCKET=my-gcs-bucket
```

## Working with Delta Lake Tables in the Cloud

### Listing Delta Lake Tables

```bash
# List Delta Lake tables in cloud storage
nessi cloud delta-list --name aws-prod --bucket my-delta-bucket --prefix data/

# List tables with detailed information
nessi delta list --cloud aws-prod --bucket my-delta-bucket --verbose
```

### Viewing Table Information

```bash
# View table schema
nessi delta schema show --cloud aws-prod --path data/my_table

# View table history
nessi delta history --cloud aws-prod --path data/my_table

# View table statistics
nessi delta stats --cloud aws-prod --path data/my_table
```

### Time Travel Operations

```bash
# View table at a specific version
nessi delta time-travel --cloud aws-prod --path data/my_table --version 5

# View table at a specific timestamp
nessi delta time-travel --cloud aws-prod --path data/my_table --timestamp "2023-05-15T14:30:00"

# Compare versions
nessi delta diff --cloud aws-prod --path data/my_table --version1 3 --version2 5
```

### Data Validation

```bash
# Validate table against rules
nessi validate --cloud aws-prod --path data/my_table --rules finance_rules.yaml

# Validate specific version
nessi validate --cloud aws-prod --path data/my_table --version 3 --rules finance_rules.yaml

# Generate validation report
nessi validate --cloud aws-prod --path data/my_table --report-format html --output report.html
```

## Authentication Methods

### AWS Authentication

Nessi.dev supports the following AWS authentication methods:

1. **Access Key and Secret Key**: Provide explicit credentials
2. **IAM Role**: Use the IAM role attached to the instance (EC2, ECS, etc.)
3. **AWS Profiles**: Use profiles defined in `~/.aws/credentials`
4. **Environment Variables**: Use `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

### Azure Authentication

Nessi.dev supports the following Azure authentication methods (using the latest Azure SDK v1.6.1+):

1. **Account Key**: Use storage account key (simplest but least secure for production)
2. **SAS Token**: Use Shared Access Signature token
3. **Azure AD**: Use Azure Active Directory authentication

### GCP Authentication

Nessi.dev supports the following GCP authentication methods:

1. **Service Account Key File**: Provide a JSON key file
2. **Application Default Credentials**: Use the default credentials from the environment
3. **Environment Variables**: Use `GOOGLE_APPLICATION_CREDENTIALS`

## Security Best Practices

1. **Use IAM Roles/Managed Identities**: Whenever possible, use IAM roles (AWS), managed identities (Azure), or service accounts (GCP) instead of access keys
2. **Least Privilege**: Grant only the necessary permissions to the credentials used by Nessi.dev
3. **Rotate Credentials**: Regularly rotate access keys and other credentials
4. **Secure Configuration Files**: Store configuration files securely and with restricted permissions
5. **Use VPC Endpoints**: When running in cloud environments, use VPC endpoints (AWS), private endpoints (Azure), or VPC Service Controls (GCP) to restrict access

## Troubleshooting

### Common Issues

#### Connection Failures

If you encounter connection issues:

1. Verify your credentials are correct
2. Check network connectivity to the cloud provider
3. Ensure the required permissions are granted
4. Check if the bucket/container exists and is accessible

```bash
# Test connection
nessi cloud test --name aws-prod --verbose
```

#### Permission Errors

If you encounter permission errors:

1. Verify the IAM policy (AWS), access policy (Azure), or IAM role (GCP) has the necessary permissions
2. For AWS, ensure the policy includes `s3:ListBucket`, `s3:GetObject`, etc.
3. For Azure, ensure the role has at least "Storage Blob Data Reader" permissions
4. For GCP, ensure the service account has at least "Storage Object Viewer" permissions

#### Delta Lake Table Not Found

If a Delta Lake table is not found:

1. Verify the table path is correct
2. Check if the `_delta_log` directory exists
3. Ensure you have permissions to access the table files

```bash
# List objects to verify path
nessi cloud list-objects --name aws-prod --bucket my-delta-bucket --prefix data/my_table/
```

## Performance Considerations

When working with large Delta Lake tables in the cloud:

1. **Use Partitioning**: Leverage Delta Lake's partitioning for better performance
2. **Limit Time Travel Depth**: Specify version ranges to avoid processing the entire transaction log
3. **Use Sampling**: For large tables, enable sampling for faster profiling
4. **Consider Region**: Deploy Nessi.dev in the same region as your data to reduce latency and data transfer costs

## Cost Optimization

To optimize costs when using cloud integration:

1. **Minimize Data Transfer**: Deploy Nessi.dev in the same region as your data
2. **Use Caching**: Enable caching to reduce repeated cloud API calls
3. **Optimize API Calls**: Use batch operations and pagination to reduce the number of API calls
4. **Monitor Usage**: Track your cloud usage and costs using cloud provider tools

## Integration with Cloud-Native Services

### AWS Integration

- **CloudWatch**: Send metrics and alerts to CloudWatch
- **AWS Lambda**: Deploy Nessi.dev validation as Lambda functions
- **AWS Glue**: Integrate with Glue Data Catalog for metadata

### Azure Integration

- **Azure Monitor**: Send metrics and alerts to Azure Monitor
- **Azure Functions**: Deploy Nessi.dev validation as Azure Functions
- **Azure Data Factory**: Trigger validations from Data Factory pipelines
- **Azure Blob Storage**: Seamless integration with the latest Azure SDK (v1.6.1+)
- **Azure Key Vault**: Secure storage for connection credentials

### GCP Integration

- **Cloud Monitoring**: Send metrics and alerts to Cloud Monitoring
- **Cloud Functions**: Deploy Nessi.dev validation as Cloud Functions
- **Cloud Composer**: Integrate with Airflow workflows in Cloud Composer

## Example Use Cases

### Data Quality Monitoring in the Cloud

```bash
# Create a scheduled validation job
nessi job create --name daily-validation --cloud aws-prod --path data/sales_table --schedule "0 6 * * *" --rules sales_rules.yaml --alert-on-failure

# View job results
nessi job results --name daily-validation
```

### Multi-Cloud Delta Lake Management

```bash
# Compare tables across clouds
nessi delta compare --source-cloud aws-prod --source-path data/users --target-cloud azure-prod --target-path data/users

# Sync validation rules across clouds
nessi rules sync --source-cloud aws-prod --target-cloud gcp-prod
```

### CI/CD Integration

```bash
# Validate data as part of CI/CD pipeline
nessi validate --cloud aws-prod --path data/financial_data --rules finance_rules.yaml --fail-on-error
```
