# Best Practices

This guide provides best practices for using Nessi.dev effectively in your data pipeline.

## Table Management

### Partitioning Strategy
- Use date-based partitioning for time-series data
- Keep partition cardinality reasonable (100-1000 partitions)
- Use Z-ordering for frequently queried columns
- Example:
```sql
CREATE TABLE events
USING DELTA
PARTITIONED BY (date)
ZORDER BY (user_id, event_type)
```

### Schema Evolution
- Use `mergeSchema` option when writing to existing tables
- Add new columns as nullable when possible
- Document schema changes in version control
- Example:
```python
df.write.format("delta") \
    .option("mergeSchema", "true") \
    .mode("append") \
    .save("s3://bucket/table")
```

## Data Quality

### Rule Configuration
- Start with basic rules and gradually add complexity
- Use descriptive rule names
- Group related rules together
- Example:
```json
{
  "rules": {
    "completeness": {
      "name": "required_fields",
      "fields": ["id", "timestamp", "user_id"],
      "threshold": 0.99
    },
    "accuracy": {
      "name": "age_range",
      "field": "age",
      "min": 0,
      "max": 120
    }
  }
}
```

### Monitoring Thresholds
- Set realistic thresholds based on historical data
- Use different thresholds for different environments
- Implement gradual threshold changes
- Example:
```yaml
thresholds:
  production:
    completeness: 0.99
    accuracy: 0.98
  staging:
    completeness: 0.95
    accuracy: 0.90
```

## Security

### API Key Management
- Rotate API keys regularly
- Use different keys for different environments
- Store keys in secure vaults
- Example:
```bash
# Rotate API key
nessi user rotate-key --username admin
```

### Access Control
- Follow principle of least privilege
- Use role-based access control
- Regularly audit permissions
- Example:
```python
# Create user with limited permissions
nessi user create \
  --username analyst \
  --role data_analyst \
  --permissions read_data,view_reports
```

## Monitoring

### Alert Configuration
- Use appropriate severity levels
- Configure meaningful alert messages
- Set up alert routing
- Example:
```yaml
alerts:
  - name: high_error_rate
    condition: error_rate > 0.1
    severity: critical
    message: "Error rate exceeds 10%"
    route: operations_team
```

### Dashboard Design
- Group related metrics
- Use appropriate visualization types
- Include historical context
- Example:
```json
{
  "panels": [
    {
      "title": "Data Quality Score",
      "type": "gauge",
      "metrics": ["data_quality_score"],
      "thresholds": [0.9, 0.95]
    },
    {
      "title": "Processing Time Trend",
      "type": "line",
      "metrics": ["processing_time"],
      "time_range": "7d"
    }
  ]
}
```

## Performance Optimization

### Caching Strategy
- Cache frequently accessed data
- Use appropriate TTL values
- Implement cache invalidation
- Example:
```python
@cache(ttl=3600)  # Cache for 1 hour
def get_table_stats(table_path):
    return nessi.table.stats(table_path)
```

### Query Optimization
- Use appropriate file formats
- Optimize partition pruning
- Leverage Z-ordering
- Example:
```python
# Optimize query with partition pruning
df = spark.read.format("delta") \
    .load("s3://bucket/table") \
    .filter("date = '2024-01-01'") \
    .filter("user_id = 123")
```

## Integration

### CI/CD Pipeline
- Run quality checks in PRs
- Generate reports automatically
- Enforce quality gates
- Example:
```yaml
# GitHub Actions workflow
- name: Quality Check
  run: |
    nessi quality check \
      --table-path ${{ env.TABLE_PATH }} \
      --rules .nessi/rules.json \
      --threshold 0.9
```

### Error Handling
- Implement graceful degradation
- Log errors appropriately
- Set up error notifications
- Example:
```python
try:
    result = nessi.quality.check(table_path)
except NessiError as e:
    logger.error(f"Quality check failed: {str(e)}")
    notify_slack(f"Quality check failed: {str(e)}")
    raise
```

## Maintenance

### Backup Strategy
- Regular table backups
- Version control for configurations
- Document recovery procedures
- Example:
```bash
# Backup table
nessi table backup \
  --source s3://bucket/table \
  --target s3://backup/table/$(date +%Y%m%d)
```

### Logging
- Use structured logging
- Include relevant context
- Set appropriate log levels
- Example:
```python
logger.info("Quality check completed", extra={
    "table": table_path,
    "score": result.score,
    "duration": result.duration
})
``` 