# Cloud Integration Quick Start Guide

This guide will help you quickly get started with Nessi.dev's cloud integration features for AWS, Azure, and GCP.

## Prerequisites

Before starting, ensure you have:
- Nessi.dev installed
- Access to at least one cloud provider (AWS, Azure, or GCP)
- Appropriate credentials for your cloud provider

## AWS Integration

### Step 1: Configure AWS Provider

```bash
# Interactive configuration
nessi cloud configure --name aws-prod --provider aws

# Or using a configuration file
cat > aws-config.yaml << EOF
provider: aws
region: us-west-2
credentials:
  access_key: YOUR_ACCESS_KEY
  secret_key: YOUR_SECRET_KEY
  # Or use IAM role
  use_iam_role: true
default_bucket: your-delta-bucket
EOF

nessi cloud configure --config-file aws-config.yaml --name aws-prod
```

### Step 2: Test AWS Connection

```bash
nessi cloud test --name aws-prod
```

### Step 3: List Delta Tables in S3

```bash
nessi cloud delta-list --name aws-prod --bucket your-delta-bucket
```

### Step 4: Analyze Delta Table in S3

```bash
# View table schema
nessi delta schema show --cloud aws-prod --path data/your_table

# View table history
nessi delta history --cloud aws-prod --path data/your_table

# Perform time travel
nessi delta time-travel --cloud aws-prod --path data/your_table --version 5
```

## Azure Integration

### Step 1: Configure Azure Provider

```bash
# Interactive configuration
nessi cloud configure --name azure-prod --provider azure

# Or using a configuration file
cat > azure-config.yaml << EOF
provider: azure
credentials:
  account_name: yourstorageaccount
  # Authentication options (choose one):
  # Option 1: Account Key
  account_key: YOUR_ACCOUNT_KEY
  # Option 2: SAS Token
  # sas_token: YOUR_SAS_TOKEN
  # Option 3: Azure AD (recommended for production)
  # use_azure_ad: true
default_bucket: your-container
# Using latest Azure SDK v1.6.1+
EOF

nessi cloud configure --config-file azure-config.yaml --name azure-prod
```

### Step 2: Test Azure Connection

```bash
nessi cloud test --name azure-prod
```

### Step 3: List Delta Tables in Azure Blob Storage

```bash
nessi cloud delta-list --name azure-prod --bucket your-container
```

### Step 4: Analyze Delta Table in Azure Blob Storage

```bash
# View table schema
nessi delta schema show --cloud azure-prod --path data/your_table

# View table history
nessi delta history --cloud azure-prod --path data/your_table

# Perform time travel
nessi delta time-travel --cloud azure-prod --path data/your_table --version 5
```

## GCP Integration

### Step 1: Configure GCP Provider

```bash
# Interactive configuration
nessi cloud configure --name gcp-prod --provider gcp

# Or using a configuration file
cat > gcp-config.yaml << EOF
provider: gcp
credentials:
  credentials_file: /path/to/credentials.json
  # Or use application default credentials
additional_options:
  project_id: your-gcp-project
default_bucket: your-gcs-bucket
EOF

nessi cloud configure --config-file gcp-config.yaml --name gcp-prod
```

### Step 2: Test GCP Connection

```bash
nessi cloud test --name gcp-prod
```

### Step 3: List Delta Tables in GCP Cloud Storage

```bash
nessi cloud delta-list --name gcp-prod --bucket your-gcs-bucket
```

### Step 4: Analyze Delta Table in GCP Cloud Storage

```bash
# View table schema
nessi delta schema show --cloud gcp-prod --path data/your_table

# View table history
nessi delta history --cloud gcp-prod --path data/your_table

# Perform time travel
nessi delta time-travel --cloud gcp-prod --path data/your_table --version 5
```

## Data Validation in the Cloud

### Validate Delta Table Against Rules

```bash
# Create validation rules
cat > finance_rules.yaml << EOF
rules:
  - name: amount_positive
    field: amount
    type: range
    config:
      min: 0.01
    severity: error
  - name: transaction_date_valid
    field: transaction_date
    type: date_format
    config:
      format: "yyyy-MM-dd"
    severity: warning
EOF

# Run validation
nessi validate --cloud aws-prod --path data/transactions --rules finance_rules.yaml
```

### Generate Validation Report

```bash
nessi validate --cloud aws-prod --path data/transactions --rules finance_rules.yaml --report-format html --output report.html
```

## Troubleshooting

### Connection Issues

If you encounter connection issues:

```bash
# Test connection with verbose output
nessi cloud test --name aws-prod --verbose

# Check environment variables
env | grep AWS_
env | grep AZURE_
env | grep GOOGLE_
```

### Permission Issues

If you encounter permission errors:

```bash
# For AWS, ensure your IAM policy includes:
# - s3:ListBucket
# - s3:GetObject
# - s3:PutObject (if writing)

# For Azure, ensure your role has at least:
# - Storage Blob Data Reader

# For GCP, ensure your service account has at least:
# - Storage Object Viewer
```

## Next Steps

- Explore the [Cloud Integration Guide](cloud_integration.md) for more details
- Set up [Scheduled Validation Jobs](scheduled_jobs.md) for your cloud data
- Configure [Cloud Alerting](cloud_alerting.md) for data quality issues
