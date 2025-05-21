# Nessi Error Handling Cheat Sheet

## Quick Reference

| Error Category | Code Range | Common Resolution Steps |
|---------------|------------|-------------------------|
| Path Errors | N1XX | Check if path exists and you have permissions |
| Configuration Errors | N2XX | Verify config file syntax and required fields |
| Authentication Errors | N3XX | Check credentials and token expiration |
| Connection Errors | N4XX | Verify network connectivity and server status |
| Delta Lake Errors | N5XX | Ensure path is a valid Delta table with _delta_log |
| Databricks Errors | N6XX | Verify Databricks workspace and credentials |
| Schema Errors | N7XX | Check schema compatibility and data types |
| Validation Errors | N8XX | Verify data meets quality and constraint requirements |
| Internal Errors | N9XX | Check logs and report issue if persistent |

## Common Error Codes and Solutions

### N101: Invalid Path
```
❌ Error: Invalid path: /path/to/file does not exist or is not accessible
```
**Solutions:**
- Check if the path exists: `ls -la /path/to/file`
- Verify you have permissions: `ls -la /path`
- Create the directory if needed: `mkdir -p /path/to/file`
- Use absolute paths instead of relative paths

### N201: Invalid Configuration
```
❌ Error: Invalid configuration: host is required
```
**Solutions:**
- Check your configuration file for syntax errors
- Ensure all required fields are present
- Verify environment variables aren't overriding config
- Run with `--debug` flag for more details

### N301: Authentication Failed
```
❌ Error: Authentication failed: Invalid Databricks token
```
**Solutions:**
- Verify your token is correct
- Check if token has expired
- Regenerate token if necessary
- Ensure environment variables are set correctly

### N401: Connection Failed
```
❌ Error: Connection failed: Could not connect to server
```
**Solutions:**
- Check network connectivity
- Verify server is running
- Check firewall settings
- Try increasing timeout with `--timeout` flag

### N501: Invalid Delta Table
```
❌ Error: Not a Delta table: /path/to/table is not a Delta Lake table
```
**Solutions:**
- Verify path points to a Delta table
- Check for _delta_log directory
- Ensure you have read permissions
- Try using `nessi repair` command

## Interactive Error Resolution

Many errors can be resolved interactively by running commands with the `--interactive` flag:

```bash
nessi schema show --path /nonexistent/path --interactive
```

This will prompt you to create the directory if it doesn't exist.

## Error Telemetry Commands

```bash
# View error telemetry status
nessi telemetry status

# View error statistics
nessi telemetry report-errors

# Export error report
nessi telemetry export-errors /path/to/report.json
```

## Testing Error Handling

```bash
# List all error codes
nessi test-error --list

# Generate a specific error
nessi test-error N101

# Test interactive resolution
nessi test-error N101 --resolvable
```

## Getting Help

```bash
# Get help for a specific command
nessi <command> --help

# Run with debug logging
NESSI_LOG_LEVEL=debug nessi <command>

# Run the error handling test script
./scripts/test_error_handling.sh
```

For more detailed information, see the [Error Handling Documentation](ERROR_HANDLING.md) and [Error Codes Reference](ERROR_CODES_REFERENCE.md).
