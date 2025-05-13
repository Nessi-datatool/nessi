# Delta Lake Features

Nessi.dev provides comprehensive support for Delta Lake, enabling advanced data quality monitoring and management capabilities for Delta Lake tables.

## Overview

Delta Lake is an open-source storage layer that brings ACID transactions to Apache Spark and big data workloads. Nessi.dev extends Delta Lake with:

- Schema management and evolution tracking
- Time travel capabilities
- Version control and history
- Data quality validation

## Schema Management

### Viewing Schema Information

```bash
# View schema for a Delta Lake table
nessi schema show --path /path/to/delta/table

# Export schema to JSON
nessi schema export --path /path/to/delta/table --output schema.json

# Compare schemas between versions
nessi schema diff --path /path/to/delta/table --version1 0 --version2 1
```

### Schema Evolution Tracking

Nessi.dev tracks schema changes over time, allowing you to:

- View the history of schema changes
- Understand how your schema has evolved
- Identify potential compatibility issues

```bash
# View schema history
nessi schema history --path /path/to/delta/table

# Get schema at a specific version
nessi schema show --path /path/to/delta/table --version 5
```

## Time Travel

Time travel allows you to access previous versions of your Delta Lake tables.

### Viewing Previous Versions

```bash
# List available versions
nessi time-travel versions --path /path/to/delta/table

# View data at a specific version
nessi time-travel show --path /path/to/delta/table --version 3

# View data at a specific timestamp
nessi time-travel show --path /path/to/delta/table --timestamp "2023-01-15T14:30:00"
```

### Data Comparison

Compare data between different versions to understand what changed:

```bash
# Compare data between versions
nessi time-travel diff --path /path/to/delta/table --version1 2 --version2 3

# Get summary statistics of changes
nessi time-travel diff-stats --path /path/to/delta/table --version1 2 --version2 3
```

## Version Control

Nessi.dev provides version control capabilities for Delta Lake tables:

### Version History

```bash
# View version history with commit information
nessi version history --path /path/to/delta/table

# Get detailed information about a specific version
nessi version info --path /path/to/delta/table --version 2
```

### Rollback Capabilities

```bash
# Rollback to a previous version (creates a new commit)
nessi version rollback --path /path/to/delta/table --version 5

# Create a new branch from a specific version
nessi version branch --path /path/to/delta/table --version 3 --name "fix-branch"
```

## Data Quality Validation

Validate Delta Lake tables against quality rules:

```bash
# Validate current version
nessi validate --path /path/to/delta/table

# Validate a specific version
nessi validate --path /path/to/delta/table --version 2

# Apply specific rule set
nessi validate --path /path/to/delta/table --rules finance_rules.yaml
```

## Integration with Python Extensions

When Python extensions are enabled, additional Delta Lake capabilities become available:

### Advanced Schema Analysis

```bash
# Generate schema recommendations
nessi schema recommend --path /path/to/delta/table

# Detect schema drift over time
nessi schema drift --path /path/to/delta/table
```

### Performance Optimization

```bash
# Analyze table performance
nessi delta analyze --path /path/to/delta/table

# Get optimization recommendations
nessi delta optimize-recommend --path /path/to/delta/table
```

## Configuration

Delta Lake features can be configured in your Nessi.dev configuration file:

```yaml
delta_lake:
  default_version_limit: 10  # Number of versions to fetch by default
  time_travel_enabled: true
  schema_tracking_enabled: true
  performance_analysis_enabled: true
```

## Requirements

- Delta Lake 1.0.0 or later
- For advanced features, Python extensions must be enabled

## Troubleshooting

### Common Issues

**Q: Why can't I access older versions?**  
A: Check if the Delta Lake table has retention policies that remove older versions.

**Q: Schema tracking shows incorrect information**  
A: Ensure you have read permissions for the _delta_log directory.

**Q: Performance is slow when analyzing large tables**  
A: Use sampling with the `--sample-ratio` option to improve performance.

### Getting Help

For more information, use the built-in help:

```bash
nessi schema --help
nessi time-travel --help
nessi version --help
```
