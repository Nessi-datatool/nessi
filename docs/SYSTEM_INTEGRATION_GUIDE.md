# Nessi System Integration Guide

This guide provides comprehensive information on integrating Nessi with other data tools, platforms, and systems in your data ecosystem.

## Table of Contents

- [Overview](#overview)
- [Integration Patterns](#integration-patterns)
- [Data Pipeline Integration](#data-pipeline-integration)
- [Cloud Platform Integration](#cloud-platform-integration)
- [Databricks Integration](#databricks-integration)
- [AWS Integration](#aws-integration)
- [CI/CD Integration](#cicd-integration)
- [Monitoring System Integration](#monitoring-system-integration)
- [Authentication Integration](#authentication-integration)
- [API Integration](#api-integration)
- [Custom Integrations](#custom-integrations)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

Nessi is designed to integrate seamlessly with your existing data ecosystem. This guide covers various integration patterns and specific integrations with popular data platforms and tools.

### Integration Capabilities

#### Core Integration Capabilities (Community Edition)

- Command-line interface for scripting and automation
- Local file system integration
- Delta Lake table support
- JSON/CSV/Parquet file support
- HTML/PDF report generation
- Basic authentication support

#### Premium Integration Capabilities (Pro Edition)

- Databricks integration
- AWS S3 integration
- Advanced authentication methods
- Workflow orchestration
- Team collaboration features

## Integration Patterns

Nessi supports several integration patterns:

### Command-line Integration

The simplest integration pattern is to use Nessi's command-line interface in scripts and automation workflows:

```bash
# Example script integrating Nessi
#!/bin/bash

# Run data quality checks
nessi quality check --path /path/to/table --rules /path/to/rules.yaml

# Generate a report
nessi report generate --path /path/to/table --output /path/to/report.html

# Take action based on quality results
if [ $? -eq 0 ]; then
    echo "Quality checks passed"
    # Proceed with downstream processes
else
    echo "Quality checks failed"
    # Send notification or take remedial action
fi
```

### File-based Integration

Nessi can read and write files in various formats, enabling integration with systems that share file storage:

```bash
# Export quality results to JSON
nessi quality check --path /path/to/table --rules /path/to/rules.yaml --output-format json --output /path/to/results.json

# Export metrics to CSV
nessi metrics show --path /path/to/table --output-format csv --output /path/to/metrics.csv
```

### API Integration (Pro Edition)

Nessi provides APIs for programmatic integration:

```python
# Python example using Nessi API client
from nessi import NessiClient

# Initialize client
client = NessiClient(api_key="your_api_key")

# Run quality checks
results = client.quality.check(
    path="/path/to/table",
    rules_file="/path/to/rules.yaml"
)

# Process results
if results.valid:
    print("Quality checks passed")
else:
    print("Quality checks failed")
    for failure in results.failures:
        print(f"- {failure}")
```

## Data Pipeline Integration

Nessi can be integrated into data pipelines to validate data quality at various stages of the pipeline.

### Apache Airflow Integration

Integrate Nessi with Apache Airflow using the BashOperator:

```python
from airflow import DAG
from airflow.operators.bash import BashOperator
from datetime import datetime, timedelta

default_args = {
    'owner': 'airflow',
    'depends_on_past': False,
    'start_date': datetime(2023, 1, 1),
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 1,
    'retry_delay': timedelta(minutes=5),
}

dag = DAG(
    'nessi_quality_checks',
    default_args=default_args,
    description='Run Nessi quality checks on data',
    schedule_interval=timedelta(days=1),
)

# Task to run quality checks
run_quality_checks = BashOperator(
    task_id='run_quality_checks',
    bash_command='nessi quality check --path /path/to/table --rules /path/to/rules.yaml',
    dag=dag,
)

# Task to generate report
generate_report = BashOperator(
    task_id='generate_report',
    bash_command='nessi report generate --path /path/to/table --output /path/to/report.html',
    dag=dag,
)

# Define task dependencies
run_quality_checks >> generate_report
```

### Apache Spark Integration

Integrate Nessi with Apache Spark using PySpark:

```python
from pyspark.sql import SparkSession
import subprocess
import json

# Initialize Spark session
spark = SparkSession.builder \
    .appName("NessiSparkIntegration") \
    .getOrCreate()

# Load data
df = spark.read.format("delta").load("/path/to/table")

# Process data with Spark
processed_df = df.filter(df.value > 0).cache()

# Save processed data
processed_df.write.format("delta").mode("overwrite").save("/path/to/processed_table")

# Run Nessi quality checks
result = subprocess.run(
    ["nessi", "quality", "check", "--path", "/path/to/processed_table", "--rules", "/path/to/rules.yaml", "--output-format", "json"],
    capture_output=True,
    text=True
)

# Parse results
quality_results = json.loads(result.stdout)

# Take action based on results
if quality_results["valid"]:
    print("Quality checks passed")
    # Proceed with downstream processes
else:
    print("Quality checks failed")
    # Log failures
    for failure in quality_results["failures"]:
        print(f"- {failure}")
```

### dbt Integration

Integrate Nessi with dbt (data build tool):

```yaml
# dbt_project.yml
name: 'my_dbt_project'
version: '1.0.0'
config-version: 2

# Define a post-hook to run Nessi quality checks
models:
  my_dbt_project:
    post-hook:
      - "{{ nessi_quality_check(this) }}"
```

Create a dbt macro for Nessi integration:

```sql
-- macros/nessi_quality_check.sql
{% macro nessi_quality_check(model) %}
  {% if target.name == 'prod' %}
    {% set table_path = target.schema ~ '.' ~ model.name %}
    {% set cmd %}
      nessi quality check --path {{ table_path }} --rules /path/to/rules.yaml
    {% endset %}
    {% do run_command(cmd) %}
  {% endif %}
{% endmacro %}
```

## Cloud Platform Integration

### Databricks Integration (Pro Edition)

Nessi integrates with Databricks to provide data quality monitoring for Databricks tables.

#### Configuration

Configure Databricks integration:

```yaml
# config.yaml
integrations:
  databricks:
    host: your-databricks-instance.cloud.databricks.com
    token: dapi_xxxxxxxxxxxxxxxxxxxxxxxx
    workspace_id: 1234567890
```

Or using environment variables:

```bash
export DATABRICKS_HOST=your-databricks-instance.cloud.databricks.com
export DATABRICKS_TOKEN=dapi_xxxxxxxxxxxxxxxxxxxxxxxx
export DATABRICKS_WORKSPACE_ID=1234567890
```

#### Usage

Use Nessi with Databricks:

```bash
# List Databricks catalogs
nessi integration databricks list-catalogs

# List schemas in a catalog
nessi integration databricks list-schemas --catalog your_catalog

# List tables in a schema
nessi integration databricks list-tables --catalog your_catalog --schema your_schema

# Run quality checks on a Databricks table
nessi quality check --integration databricks --catalog your_catalog --schema your_schema --table your_table --rules /path/to/rules.yaml
```

#### Databricks Notebook Integration

Integrate Nessi in a Databricks notebook:

```python
# Databricks notebook
import subprocess
import json
import os

# Set environment variables
os.environ["DATABRICKS_TOKEN"] = dbutils.secrets.get("nessi", "databricks_token")
os.environ["DATABRICKS_HOST"] = "your-databricks-instance.cloud.databricks.com"

# Run Nessi quality checks
result = subprocess.run(
    ["nessi", "quality", "check", "--integration", "databricks", "--catalog", "your_catalog", "--schema", "your_schema", "--table", "your_table", "--rules", "/path/to/rules.yaml", "--output-format", "json"],
    capture_output=True,
    text=True
)

# Parse results
quality_results = json.loads(result.stdout)

# Display results
display(quality_results)
```

### AWS Integration (Pro Edition)

Nessi integrates with AWS services, particularly S3 for storage.

#### Configuration

Configure AWS integration:

```yaml
# config.yaml
integrations:
  aws:
    region: us-west-2
    access_key_id: AKIAXXXXXXXXXXXXXXXX
    secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Or using environment variables:

```bash
export AWS_REGION=us-west-2
export AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

#### Usage

Use Nessi with AWS S3:

```bash
# List S3 buckets
nessi integration aws list-buckets

# List objects in a bucket
nessi integration aws list-objects --bucket your_bucket --prefix your_prefix

# Run quality checks on a Delta Lake table in S3
nessi quality check --integration aws --bucket your_bucket --key your_key --rules /path/to/rules.yaml
```

#### AWS Lambda Integration

Create an AWS Lambda function to run Nessi:

```python
# Lambda function
import json
import subprocess
import os

def lambda_handler(event, context):
    # Set environment variables
    os.environ["AWS_REGION"] = "us-west-2"
    
    # Get parameters from event
    bucket = event.get("bucket")
    key = event.get("key")
    rules = event.get("rules", "/tmp/rules.yaml")
    
    # Write rules to file if provided in event
    if "rules_content" in event:
        with open(rules, "w") as f:
            f.write(event["rules_content"])
    
    # Run Nessi quality checks
    result = subprocess.run(
        ["nessi", "quality", "check", "--integration", "aws", "--bucket", bucket, "--key", key, "--rules", rules, "--output-format", "json"],
        capture_output=True,
        text=True
    )
    
    # Parse results
    quality_results = json.loads(result.stdout)
    
    return {
        "statusCode": 200,
        "body": quality_results
    }
```

## CI/CD Integration

Integrate Nessi into your CI/CD pipelines to ensure data quality throughout your development lifecycle.

### GitHub Actions Integration

Create a GitHub Actions workflow for Nessi:

```yaml
# .github/workflows/data-quality.yml
name: Data Quality Checks

on:
  schedule:
    - cron: '0 0 * * *'  # Run daily at midnight
  workflow_dispatch:  # Allow manual triggering

jobs:
  quality-checks:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
      
    - name: Install Nessi
      run: |
        curl -sSL https://nessi.dev/install.sh | bash
        
    - name: Run quality checks
      run: |
        nessi quality check --path /path/to/table --rules ./quality_rules.yaml
        
    - name: Generate report
      run: |
        nessi report generate --path /path/to/table --output ./report.html
        
    - name: Upload report artifact
      uses: actions/upload-artifact@v3
      with:
        name: quality-report
        path: ./report.html
```

### GitLab CI Integration

Create a GitLab CI pipeline for Nessi:

```yaml
# .gitlab-ci.yml
stages:
  - quality

quality-checks:
  stage: quality
  image: ubuntu:latest
  script:
    - apt-get update && apt-get install -y curl
    - curl -sSL https://nessi.dev/install.sh | bash
    - nessi quality check --path /path/to/table --rules ./quality_rules.yaml
    - nessi report generate --path /path/to/table --output ./report.html
  artifacts:
    paths:
      - report.html
```

### Jenkins Integration

Create a Jenkins pipeline for Nessi:

```groovy
// Jenkinsfile
pipeline {
    agent any
    
    stages {
        stage('Install Nessi') {
            steps {
                sh 'curl -sSL https://nessi.dev/install.sh | bash'
            }
        }
        
        stage('Run Quality Checks') {
            steps {
                sh 'nessi quality check --path /path/to/table --rules ./quality_rules.yaml'
            }
        }
        
        stage('Generate Report') {
            steps {
                sh 'nessi report generate --path /path/to/table --output ./report.html'
                archiveArtifacts artifacts: 'report.html', fingerprint: true
            }
        }
    }
    
    post {
        failure {
            mail to: 'team@example.com',
                 subject: "Failed Pipeline: ${currentBuild.fullDisplayName}",
                 body: "Quality checks failed. See ${env.BUILD_URL} for details."
        }
    }
}
```

## Monitoring System Integration

Integrate Nessi with monitoring systems to track data quality metrics over time.

### Prometheus Integration

Export Nessi metrics to Prometheus:

```bash
# Export metrics in Prometheus format
nessi metrics export --path /path/to/table --format prometheus --output /path/to/metrics
```

Create a Prometheus scrape configuration:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'nessi'
    static_configs:
      - targets: ['localhost:9090']
    file_sd_configs:
      - files:
        - '/path/to/metrics'
```

### Grafana Integration

Create a Grafana dashboard for Nessi metrics:

```bash
# Export metrics in JSON format for Grafana
nessi metrics export --path /path/to/table --format json --output /path/to/metrics.json
```

Example Grafana dashboard configuration:

```json
{
  "dashboard": {
    "id": null,
    "title": "Nessi Data Quality Dashboard",
    "tags": ["nessi", "data-quality"],
    "timezone": "browser",
    "panels": [
      {
        "title": "Overall Quality Score",
        "type": "gauge",
        "datasource": "Prometheus",
        "targets": [
          {
            "expr": "nessi_quality_score{table=\"your_table\"}",
            "refId": "A"
          }
        ],
        "options": {
          "thresholds": [
            {
              "color": "red",
              "value": 0
            },
            {
              "color": "yellow",
              "value": 70
            },
            {
              "color": "green",
              "value": 90
            }
          ]
        }
      },
      {
        "title": "Quality Metrics Over Time",
        "type": "graph",
        "datasource": "Prometheus",
        "targets": [
          {
            "expr": "nessi_quality_completeness{table=\"your_table\"}",
            "refId": "A",
            "legendFormat": "Completeness"
          },
          {
            "expr": "nessi_quality_accuracy{table=\"your_table\"}",
            "refId": "B",
            "legendFormat": "Accuracy"
          },
          {
            "expr": "nessi_quality_consistency{table=\"your_table\"}",
            "refId": "C",
            "legendFormat": "Consistency"
          }
        ]
      }
    ]
  }
}
```

### ELK Stack Integration

Export Nessi metrics to Elasticsearch:

```bash
# Export metrics in JSON format for Elasticsearch
nessi metrics export --path /path/to/table --format json --output /path/to/metrics.json

# Index metrics in Elasticsearch
curl -X POST "localhost:9200/nessi-metrics/_doc" -H "Content-Type: application/json" -d @/path/to/metrics.json
```

Example Kibana dashboard configuration:

```json
{
  "attributes": {
    "title": "Nessi Data Quality Dashboard",
    "hits": 0,
    "description": "Data quality metrics from Nessi",
    "panelsJSON": "[{\"embeddableConfig\":{},\"gridData\":{\"h\":15,\"i\":\"1\",\"w\":24,\"x\":0,\"y\":0},\"id\":\"nessi-quality-score\",\"panelIndex\":\"1\",\"type\":\"visualization\",\"version\":\"7.10.0\"},{\"embeddableConfig\":{},\"gridData\":{\"h\":15,\"i\":\"2\",\"w\":24,\"x\":24,\"y\":0},\"id\":\"nessi-quality-metrics\",\"panelIndex\":\"2\",\"type\":\"visualization\",\"version\":\"7.10.0\"}]",
    "optionsJSON": "{\"hidePanelTitles\":false,\"useMargins\":true}",
    "version": 1,
    "timeRestore": false,
    "kibanaSavedObjectMeta": {
      "searchSourceJSON": "{\"query\":{\"language\":\"kuery\",\"query\":\"\"},\"filter\":[]}"
    }
  }
}
```

## Authentication Integration

Integrate Nessi with authentication systems to secure access.

### Environment Variable Authentication

Use environment variables for authentication:

```bash
# Set authentication credentials
export NESSI_API_KEY=your_api_key
export DATABRICKS_TOKEN=your_databricks_token
export AWS_ACCESS_KEY_ID=your_aws_access_key
export AWS_SECRET_ACCESS_KEY=your_aws_secret_key

# Run Nessi commands
nessi quality check --path /path/to/table --rules /path/to/rules.yaml
```

### Credential File Authentication

Use credential files for authentication:

```yaml
# ~/.nessi/credentials.yaml
api_key: your_api_key
integrations:
  databricks:
    token: your_databricks_token
  aws:
    access_key_id: your_aws_access_key
    secret_access_key: your_aws_secret_key
```

```bash
# Run Nessi commands with credentials file
nessi --credentials-file ~/.nessi/credentials.yaml quality check --path /path/to/table --rules /path/to/rules.yaml
```

### OAuth Integration (Pro Edition)

Integrate Nessi with OAuth for authentication:

```yaml
# config.yaml
auth:
  type: oauth
  client_id: your_client_id
  client_secret: your_client_secret
  auth_url: https://auth.example.com/oauth/authorize
  token_url: https://auth.example.com/oauth/token
  redirect_url: http://localhost:8080/callback
  scopes: ["read", "write"]
```

```bash
# Authenticate with OAuth
nessi auth login

# Run Nessi commands with OAuth authentication
nessi quality check --path /path/to/table --rules /path/to/rules.yaml
```

## API Integration

Integrate with Nessi's APIs for programmatic access.

### REST API Integration (Pro Edition)

Use Nessi's REST API:

```python
# Python example
import requests
import json

# API configuration
api_url = "https://api.nessi.dev/v1"
api_key = "your_api_key"
headers = {
    "Authorization": f"Bearer {api_key}",
    "Content-Type": "application/json"
}

# Run quality checks
response = requests.post(
    f"{api_url}/quality/check",
    headers=headers,
    json={
        "path": "/path/to/table",
        "rules_file": "/path/to/rules.yaml"
    }
)

# Process response
if response.status_code == 200:
    results = response.json()
    print(f"Quality check status: {'Passed' if results['valid'] else 'Failed'}")
    if not results['valid']:
        for failure in results['failures']:
            print(f"- {failure}")
else:
    print(f"Error: {response.status_code} - {response.text}")
```

### gRPC API Integration (Pro Edition)

Use Nessi's gRPC API:

```python
# Python example
import grpc
import nessi_pb2
import nessi_pb2_grpc

# Create gRPC channel
channel = grpc.secure_channel("api.nessi.dev:443", grpc.ssl_channel_credentials())
stub = nessi_pb2_grpc.NessiServiceStub(channel)

# Create request
request = nessi_pb2.QualityCheckRequest(
    path="/path/to/table",
    rules_file="/path/to/rules.yaml"
)

# Add API key metadata
metadata = [("authorization", f"Bearer your_api_key")]

# Call API
response = stub.QualityCheck(request, metadata=metadata)

# Process response
print(f"Quality check status: {'Passed' if response.valid else 'Failed'}")
if not response.valid:
    for failure in response.failures:
        print(f"- {failure}")
```

### GraphQL API Integration (Pro Edition)

Use Nessi's GraphQL API:

```python
# Python example
import requests
import json

# API configuration
api_url = "https://api.nessi.dev/graphql"
api_key = "your_api_key"
headers = {
    "Authorization": f"Bearer {api_key}",
    "Content-Type": "application/json"
}

# GraphQL query
query = """
query QualityCheck($path: String!, $rulesFile: String!) {
  qualityCheck(path: $path, rulesFile: $rulesFile) {
    valid
    score
    failures
    metrics {
      name
      score
      description
    }
  }
}
"""

# Variables
variables = {
    "path": "/path/to/table",
    "rulesFile": "/path/to/rules.yaml"
}

# Send request
response = requests.post(
    api_url,
    headers=headers,
    json={
        "query": query,
        "variables": variables
    }
)

# Process response
if response.status_code == 200:
    results = response.json()
    if "errors" in results:
        print(f"GraphQL errors: {results['errors']}")
    else:
        data = results["data"]["qualityCheck"]
        print(f"Quality check status: {'Passed' if data['valid'] else 'Failed'}")
        print(f"Overall score: {data['score']}")
        if not data['valid']:
            for failure in data['failures']:
                print(f"- {failure}")
else:
    print(f"Error: {response.status_code} - {response.text}")
```

## Custom Integrations

Develop custom integrations for Nessi using the extension system.

### Custom Integration Example

Create a custom integration for a proprietary data system:

```go
package main

import (
    "github.com/nessi-dev/nessi/pkg/integration"
)

// MySystemIntegration is a custom integration for MySystem
type MySystemIntegration struct {
    client *MySystemClient
    connected bool
}

// Name returns the integration name
func (i *MySystemIntegration) Name() string {
    return "mysystem"
}

// Description returns the integration description
func (i *MySystemIntegration) Description() string {
    return "Integration with MySystem"
}

// Connect establishes a connection to MySystem
func (i *MySystemIntegration) Connect(parameters map[string]interface{}) error {
    // Extract parameters
    host, ok := parameters["host"].(string)
    if !ok {
        return fmt.Errorf("host parameter is required")
    }
    
    apiKey, ok := parameters["api_key"].(string)
    if !ok {
        return fmt.Errorf("api_key parameter is required")
    }
    
    // Create client
    client, err := NewMySystemClient(host, apiKey)
    if err != nil {
        return err
    }
    
    // Test connection
    if err := client.Ping(); err != nil {
        return err
    }
    
    // Store client
    i.client = client
    i.connected = true
    
    return nil
}

// Disconnect closes the connection
func (i *MySystemIntegration) Disconnect() error {
    if !i.connected {
        return nil
    }
    
    if err := i.client.Close(); err != nil {
        return err
    }
    
    i.connected = false
    return nil
}

// ListResources lists available resources
func (i *MySystemIntegration) ListResources() ([]integration.Resource, error) {
    if !i.connected {
        return nil, fmt.Errorf("not connected")
    }
    
    // Get resources from MySystem
    myResources, err := i.client.ListResources()
    if err != nil {
        return nil, err
    }
    
    // Convert to Nessi resources
    resources := make([]integration.Resource, len(myResources))
    for j, res := range myResources {
        resources[j] = integration.Resource{
            ID:   res.ID,
            Name: res.Name,
            Type: res.Type,
            Metadata: map[string]interface{}{
                "created_at": res.CreatedAt,
                "owner":      res.Owner,
                "size":       res.Size,
            },
        }
    }
    
    return resources, nil
}

// GetResource gets a specific resource
func (i *MySystemIntegration) GetResource(id string) (integration.Resource, error) {
    if !i.connected {
        return integration.Resource{}, fmt.Errorf("not connected")
    }
    
    // Get resource from MySystem
    myResource, err := i.client.GetResource(id)
    if err != nil {
        return integration.Resource{}, err
    }
    
    // Convert to Nessi resource
    resource := integration.Resource{
        ID:   myResource.ID,
        Name: myResource.Name,
        Type: myResource.Type,
        Metadata: map[string]interface{}{
            "created_at": myResource.CreatedAt,
            "owner":      myResource.Owner,
            "size":       myResource.Size,
        },
        Data: myResource.Data,
    }
    
    return resource, nil
}
```

For more information on developing custom integrations, see the [Extension Guide](EXTENSION_GUIDE.md).

## Best Practices

### Integration Best Practices

- **Use Environment Variables**: Store sensitive information like API keys and tokens in environment variables rather than configuration files
- **Implement Proper Error Handling**: Handle integration errors gracefully with appropriate error messages and recovery mechanisms
- **Monitor Integration Health**: Set up monitoring for integration health to detect issues early
- **Implement Retries**: Use retry mechanisms with exponential backoff for transient failures
- **Limit Permissions**: Use the principle of least privilege when setting up integrations
- **Test Integrations**: Thoroughly test integrations before deploying to production
- **Document Integrations**: Document integration details, including configuration, usage, and troubleshooting

### Security Best Practices

- **Secure Credentials**: Protect authentication credentials using secure storage
- **Use Encryption**: Encrypt sensitive data in transit and at rest
- **Audit Access**: Implement audit logging for integration access
- **Rotate Credentials**: Regularly rotate API keys and tokens
- **Validate Input**: Validate all input data before processing
- **Use Secure Connections**: Always use secure connections (HTTPS, TLS) for API calls

## Troubleshooting

### Common Integration Issues

#### Authentication Failures

**Issue**: Unable to authenticate with external system

**Solution**:
1. Verify credentials are correct
2. Check that the authentication method is supported
3. Ensure network connectivity to the authentication endpoint
4. Check for expired tokens or credentials

#### Connection Timeouts

**Issue**: Connection to external system times out

**Solution**:
1. Check network connectivity
2. Verify firewall rules allow the connection
3. Increase connection timeout settings
4. Check if the external system is experiencing issues

#### Resource Not Found

**Issue**: External system resource not found

**Solution**:
1. Verify the resource exists in the external system
2. Check for typos in resource identifiers
3. Ensure you have permission to access the resource
4. Check if the resource path or format is correct

#### Integration Compatibility

**Issue**: Integration not compatible with external system version

**Solution**:
1. Check the supported versions for the integration
2. Update the integration or external system to a compatible version
3. Use a compatibility layer if available

### Debugging Integrations

Enable integration debugging:

```bash
# Enable integration debugging
export NESSI_INTEGRATION_DEBUG=1
nessi integration databricks list-catalogs
```

View integration logs:

```bash
# View integration logs
cat ~/.nessi/logs/integration.log
```

Test integration connectivity:

```bash
# Test Databricks connectivity
nessi integration databricks test-connection

# Test AWS connectivity
nessi integration aws test-connection
```

---

For more information on Nessi integrations, please visit [nessi.dev/integrations](https://nessi.dev/integrations) or contact support@nessi.dev.
