# Nessi Troubleshooting Guide

This guide provides solutions for common issues you might encounter when using Nessi. If you're experiencing problems, follow the appropriate troubleshooting steps below.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Configuration Issues](#configuration-issues)
- [Delta Lake Issues](#delta-lake-issues)
- [Databricks Integration Issues](#databricks-integration-issues)
- [AWS S3 Integration Issues](#aws-s3-integration-issues)
- [Schema Validation Issues](#schema-validation-issues)
- [Data Quality Issues](#data-quality-issues)
- [Performance Issues](#performance-issues)
- [Report Generation Issues](#report-generation-issues)
- [License Issues](#license-issues)
- [Common Error Codes](#common-error-codes)
- [Getting Additional Help](#getting-additional-help)

## Installation Issues

### Command Not Found

**Issue**: After installation, running `nessi` results in "command not found" error.

**Solutions**:
1. Ensure the installation completed successfully
2. Verify the installation path is in your system PATH
3. For manual installations, check that the binary has execute permissions:
   ```bash
   chmod +x /path/to/nessi
   ```
4. Try using the full path to the binary:
   ```bash
   /path/to/nessi --version
   ```

### Permission Denied

**Issue**: Permission denied when trying to run Nessi.

**Solutions**:
1. Ensure the binary has execute permissions:
   ```bash
   chmod +x /path/to/nessi
   ```
2. If using Docker, ensure you have permissions to run Docker containers

### Incorrect Version

**Issue**: The installed version doesn't match the expected version.

**Solutions**:
1. Verify you downloaded the correct version
2. Check if multiple versions are installed on your system
3. Update to the latest version:
   ```bash
   curl -sSL https://nessi.dev/install.sh | bash
   ```

### Installation Script Fails

**Issue**: The installation script fails to complete.

**Solutions**:
1. Check your internet connection
2. Ensure you have sufficient permissions
3. Try manual installation by downloading the binary directly
4. Check the installation logs for specific errors

## Configuration Issues

### Configuration File Not Found

**Issue**: Nessi reports that it cannot find the configuration file.

**Solutions**:
1. Create a default configuration file:
   ```bash
   nessi config init
   ```
2. Specify the configuration file path explicitly:
   ```bash
   nessi --config /path/to/config.yaml command
   ```
3. Check that the default configuration directory exists:
   ```bash
   mkdir -p ~/.nessi
   ```

### Invalid Configuration

**Issue**: Nessi reports that the configuration is invalid.

**Solutions**:
1. Check the configuration file syntax
2. Ensure all required fields are present
3. Reset to the default configuration:
   ```bash
   nessi config reset
   ```
4. Check the logs for specific validation errors:
   ```bash
   nessi --log-level debug config validate
   ```

### Environment Variables Not Applied

**Issue**: Environment variables don't seem to override configuration settings.

**Solutions**:
1. Ensure environment variables are correctly formatted (e.g., `NESSI_LOG_LEVEL=debug`)
2. Verify the environment variables are set in the current shell session
3. Check the precedence order: command-line flags override environment variables, which override configuration file settings

## Delta Lake Issues

### Unable to Detect Delta Lake Table

**Issue**: Nessi cannot detect a Delta Lake table at the specified path.

**Solutions**:
1. Verify the path points to a valid Delta Lake table (contains `_delta_log` directory)
2. Check file permissions for the directory and files
3. Ensure the path is absolute or correctly relative to the current directory
4. Try running with debug logs:
   ```bash
   nessi --log-level debug tables describe --path /path/to/table
   ```

### Transaction Log Errors

**Issue**: Errors when reading the Delta Lake transaction log.

**Solutions**:
1. Check if the Delta Lake table is corrupted
2. Verify compatibility with the Delta Lake version
3. Ensure all transaction log files are accessible
4. Try running with debug logs:
   ```bash
   nessi --log-level debug tables history --path /path/to/table
   ```

### Time Travel Issues

**Issue**: Time travel to a specific version or timestamp fails.

**Solutions**:
1. Verify the version exists in the table history
2. Check the timestamp format (ISO 8601: `YYYY-MM-DDTHH:MM:SSZ`)
3. Ensure the version or timestamp is within the available history range
4. List available versions:
   ```bash
   nessi tables history --path /path/to/table
   ```

## Databricks Integration Issues

### Authentication Failures

**Issue**: Cannot authenticate with Databricks.

**Solutions**:
1. Verify your Databricks token is correct and not expired
2. Check that the Databricks host URL is correct
3. Ensure network connectivity to the Databricks instance
4. Verify the environment variables are set correctly:
   ```bash
   export DATABRICKS_HOST=your_databricks_host
   export DATABRICKS_TOKEN=your_databricks_token
   ```

### Resource Not Found

**Issue**: Databricks resource (catalog, schema, table) not found.

**Solutions**:
1. Verify the resource exists in Databricks
2. Check for typos in resource names
3. Ensure you have permissions to access the resource
4. List available resources:
   ```bash
   nessi integration databricks list-catalogs
   nessi integration databricks list-schemas --catalog your_catalog
   nessi integration databricks list-tables --catalog your_catalog --schema your_schema
   ```

### Rate Limiting

**Issue**: Experiencing rate limiting when accessing Databricks.

**Solutions**:
1. Reduce the frequency of requests
2. Implement retry logic with exponential backoff
3. Contact your Databricks administrator to increase rate limits
4. Use batch operations where possible

## AWS S3 Integration Issues

### Authentication Failures

**Issue**: Cannot authenticate with AWS S3.

**Solutions**:
1. Verify your AWS credentials are correct
2. Check that the AWS region is set correctly
3. Ensure the IAM role or user has appropriate permissions
4. Verify the environment variables are set correctly:
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_REGION=your_region
   ```

### Bucket or Object Not Found

**Issue**: S3 bucket or object not found.

**Solutions**:
1. Verify the bucket and object exist
2. Check for typos in bucket and object names
3. Ensure you have permissions to access the bucket and object
4. List available objects:
   ```bash
   nessi integration s3 list-objects --bucket your_bucket --prefix your_prefix
   ```

### Access Denied

**Issue**: Access denied when accessing S3 resources.

**Solutions**:
1. Verify your IAM permissions
2. Check bucket policies and ACLs
3. Ensure the bucket is in the correct region
4. Check if the bucket requires specific encryption settings

## Schema Validation Issues

### Schema Validation Failures

**Issue**: Schema validation fails for a table.

**Solutions**:
1. Compare the actual schema with the expected schema
2. Check for data type mismatches
3. Verify required fields are present
4. Check for case sensitivity issues in field names
5. View the current schema:
   ```bash
   nessi schema show --path /path/to/table
   ```

### Schema Evolution Issues

**Issue**: Problems tracking schema evolution.

**Solutions**:
1. Ensure schema history is being recorded
2. Check for compatibility between schema versions
3. Verify schema changes follow the expected patterns
4. View schema history:
   ```bash
   nessi schema history --path /path/to/table
   ```

### Schema Comparison Discrepancies

**Issue**: Unexpected differences in schema comparison.

**Solutions**:
1. Check for whitespace or case sensitivity differences
2. Verify field order is consistent (or use order-insensitive comparison)
3. Check for metadata differences that might not be relevant
4. Use detailed comparison:
   ```bash
   nessi schema compare --source /path/to/table1 --target /path/to/table2 --detailed
   ```

## Data Quality Issues

### Quality Check Failures

**Issue**: Data quality checks fail unexpectedly.

**Solutions**:
1. Verify the quality rules are correctly defined
2. Check the data for actual quality issues
3. Adjust quality thresholds if necessary
4. Run with detailed output:
   ```bash
   nessi quality check --path /path/to/table --rules /path/to/rules.yaml --detailed
   ```

### Rule Definition Errors

**Issue**: Errors in quality rule definitions.

**Solutions**:
1. Check the YAML syntax of the rules file
2. Verify all required fields are present for each rule
3. Ensure column names match the actual schema
4. Validate the rules file:
   ```bash
   nessi quality validate-rules --rules /path/to/rules.yaml
   ```

### False Positives/Negatives

**Issue**: Quality checks produce false positives or negatives.

**Solutions**:
1. Refine rule definitions to be more precise
2. Adjust thresholds based on data characteristics
3. Consider adding or modifying rules for specific edge cases
4. Test rules on known good and bad data samples

## Performance Issues

### Slow Processing for Large Tables

**Issue**: Performance issues when processing large tables.

**Solutions**:
1. Use sampling to reduce processing time:
   ```bash
   nessi quality check --path /path/to/table --rules /path/to/rules.yaml --sample-size 1000
   ```
2. Increase memory allocation if possible
3. Process tables incrementally based on partitions
4. Run operations during off-peak hours

### High Memory Usage

**Issue**: Nessi uses excessive memory.

**Solutions**:
1. Use streaming operations where available
2. Process data in smaller batches
3. Reduce concurrency settings
4. Monitor memory usage with profiling tools

### Timeout Errors

**Issue**: Operations timeout before completion.

**Solutions**:
1. Increase timeout settings if available
2. Break operations into smaller chunks
3. Optimize query patterns
4. Check for resource contention on the system

## Report Generation Issues

### Report Generation Failures

**Issue**: Unable to generate reports.

**Solutions**:
1. Check write permissions for the output directory
2. Verify the report template exists
3. Ensure all required data is available
4. Check for disk space issues
5. Run with debug logs:
   ```bash
   nessi --log-level debug report generate --path /path/to/table --output /path/to/report.html
   ```

### Missing Report Elements

**Issue**: Generated reports are missing expected elements.

**Solutions**:
1. Verify all required data was collected
2. Check the report template configuration
3. Ensure quality checks and metrics collection ran successfully
4. Try a different report template:
   ```bash
   nessi report generate --path /path/to/table --output /path/to/report.html --template alternative
   ```

### Formatting Issues

**Issue**: Report formatting is incorrect.

**Solutions**:
1. Check for CSS or template issues
2. Verify the browser or PDF viewer is compatible
3. Try a different output format
4. Update to the latest version of Nessi

## License Issues

### License Activation Failures

**Issue**: Unable to activate Pro Edition license.

**Solutions**:
1. Verify the license key is correct
2. Check internet connectivity for license validation
3. Ensure the machine ID hasn't changed
4. Contact support if the issue persists

### Feature Not Available

**Issue**: Attempting to use a feature that requires Pro Edition.

**Solutions**:
1. Verify your license status:
   ```bash
   nessi config license status
   ```
2. Upgrade to Pro Edition if necessary
3. Activate your license if already purchased
4. Start a trial if eligible:
   ```bash
   nessi config license trial
   ```

### Trial Expiration

**Issue**: Pro Edition trial has expired.

**Solutions**:
1. Purchase a Pro Edition license
2. Check if you're eligible for an extension
3. Continue using Community Edition features
4. Contact sales for special arrangements

## Common Error Codes

Nessi uses standardized error codes to help diagnose issues:

### N1XX: Input/Output Errors

- **N100**: File not found
- **N101**: Permission denied
- **N102**: I/O error
- **N103**: Network error

### N2XX: Configuration Errors

- **N200**: Configuration file not found
- **N201**: Invalid configuration
- **N202**: Missing required configuration
- **N203**: Environment variable error

### N3XX: Delta Lake Errors

- **N300**: Not a Delta Lake table
- **N301**: Transaction log error
- **N302**: Version not found
- **N303**: Timestamp not found

### N4XX: Schema Errors

- **N400**: Schema validation error
- **N401**: Schema evolution error
- **N402**: Schema comparison error
- **N403**: Schema not found

### N5XX: Quality Errors

- **N500**: Quality rule definition error
- **N501**: Quality check failure
- **N502**: Quality threshold exceeded
- **N503**: Quality metric calculation error

### N6XX: Metrics Errors

- **N600**: Metrics collection error
- **N601**: Metrics storage error
- **N602**: Metrics analysis error
- **N603**: Metrics visualization error

### N7XX: Report Errors

- **N700**: Report generation error
- **N701**: Report template error
- **N702**: Report output error
- **N703**: Report data error

### N8XX: Integration Errors

- **N800**: Integration authentication error
- **N801**: Integration resource error
- **N802**: Integration rate limiting
- **N803**: Integration feature not available

### N9XX: Internal Errors

- **N900**: Unexpected error
- **N901**: Feature not implemented
- **N902**: Internal assertion error
- **N903**: Dependency error

## Getting Additional Help

If you've tried the troubleshooting steps above and still have issues, you can get additional help:

### Community Support

- **GitHub Issues**: Report bugs and request features on the [GitHub Issues](https://github.com/nessi-dev/nessi/issues) page
- **GitHub Discussions**: Ask questions and discuss ideas on the [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions) page
- **Stack Overflow**: Ask questions with the `nessi` tag on Stack Overflow

### Pro Edition Support

Pro Edition customers have access to additional support options:

- **Email Support**: Contact support@nessi.dev for assistance
- **Priority Issue Resolution**: Get priority attention for reported issues
- **Hotfixes**: Receive hotfixes for critical issues

### Providing Information for Support

When seeking help, provide the following information:

1. Nessi version: `nessi --version`
2. Operating system and version
3. Detailed error message and error code
4. Steps to reproduce the issue
5. Relevant configuration settings (with sensitive information redacted)
6. Debug logs: `nessi --log-level debug command > nessi-debug.log 2>&1`

### Self-Help Resources

- [Documentation](https://nessi.dev/docs): Comprehensive documentation
- [FAQ](https://nessi.dev/faq): Frequently asked questions
- [Blog](https://nessi.dev/blog): Articles and tutorials
- [Release Notes](https://nessi.dev/releases): Information about recent releases

---

For more information, visit [nessi.dev](https://nessi.dev) or contact support@nessi.dev.
