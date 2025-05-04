# Docker Troubleshooting Guide

This guide helps you troubleshoot common issues when running nessi.dev in Docker containers.

## Quick Checks

Before diving into specific issues, try these quick checks:

```bash
# Check if containers are running
docker-compose ps

# View container logs
docker-compose logs

# Check container resource usage
docker stats
```

## Common Issues

### 1. Installation Problems

#### Docker Issues
- **Symptoms**:
  - Container fails to start
  - Port conflicts
  - Permission errors
  - Resource constraints
  - Network issues
  - Volume mounting problems
  - Image pull failures
  - Configuration errors

- **Solution Steps**:
  1. Check Docker status:
     ```bash
     docker info
     docker ps -a
     ```
  2. Verify resource limits:
     ```bash
     docker stats
     ```
  3. Check logs:
     ```bash
     docker-compose logs -f
     ```
  4. Verify network:
     ```bash
     docker network ls
     ```
  5. Check volumes:
     ```bash
     docker volume ls
     ```
  6. Verify configuration:
     ```bash
     docker-compose config
     ```
  7. Check permissions:
     ```bash
     ls -la /var/run/docker.sock
     ```
  8. Monitor resources:
     ```bash
     top
     free -h
     ```

#### Development Environment Issues
- **Symptoms**:
  - Python package conflicts
  - Virtual environment problems
  - Dependency issues
  - Path configuration errors
  - Permission problems
  - Service startup failures
  - Configuration errors

- **Solution Steps**:
  1. Check Python environment:
     ```bash
     python --version
     pip list
     ```
  2. Verify virtual environment:
     ```bash
     source venv/bin/activate
     which python
     ```
  3. Check dependencies:
     ```bash
     pip check
     pip freeze > requirements.txt
     ```
  4. Verify paths:
     ```bash
     echo $PATH
     which python
     ```
  5. Check permissions:
     ```bash
     ls -la venv/
     ```
  6. Check services:
     ```bash
     ps aux | grep python
     ```
  7. Verify configuration:
     ```bash
     cat .env
     ```

### 2. Monitoring Issues

#### Grafana Issues
- **Symptoms**:
  - Dashboard not loading
  - Data source connection problems
  - Alert configuration errors
  - Permission issues
  - Configuration errors

- **Solution Steps**:
  1. Check Grafana logs:
     ```bash
     docker-compose logs grafana
     ```
  2. Verify data source:
     ```bash
     curl -u admin:admin http://localhost:3000/api/datasources
     ```
  3. Check permissions:
     ```bash
     ls -la /var/lib/grafana
     ```
  4. Verify configuration:
     ```bash
     cat /etc/grafana/grafana.ini
     ```

#### Prometheus Issues
- **Symptoms**:
  - Metric collection problems
  - Storage issues
  - Query errors
  - Configuration problems

- **Solution Steps**:
  1. Check Prometheus logs:
     ```bash
     docker-compose logs prometheus
     ```
  2. Verify targets:
     ```bash
     curl http://localhost:9090/api/v1/targets
     ```
  3. Check storage:
     ```bash
     ls -la /prometheus
     ```
  4. Verify configuration:
     ```bash
     cat /etc/prometheus/prometheus.yml
     ```

### 3. Security Issues

#### Authentication Issues
- **Symptoms**:
  - Login failures
  - Permission denied errors
  - Configuration problems

- **Solution Steps**:
  1. Check auth logs:
     ```bash
     docker-compose logs auth
     ```
  2. Verify credentials:
     ```bash
     cat .env | grep AUTH
     ```
  3. Check configuration:
     ```bash
     cat /etc/auth/config.yml
     ```

#### SSL Issues
- **Symptoms**:
  - Certificate errors
  - Connection problems
  - Configuration issues

- **Solution Steps**:
  1. Check SSL logs:
     ```bash
     docker-compose logs nginx
     ```
  2. Verify certificates:
     ```bash
     openssl x509 -in certificate.crt -text -noout
     ```
  3. Check configuration:
     ```bash
     cat /etc/nginx/nginx.conf
     ```

## Getting Help

- **Documentation**: https://docs.nessi.dev
- **Issue Tracker**: https://github.com/your-org/nessi-dev/issues
- **Community Forum**: https://community.nessi.dev
- **Support Email**: support@nessi.dev

## Next Steps

- [Quick Start Guide](quickstart.md)
- [Features Overview](features.md)
- [API Documentation](../API.md)
- [Installation Guide](installation.md)

# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with Nessi.dev.

## Common Issues

### Connection Issues

#### SSL/TLS Errors
**Symptoms:**
- `SSL: CERTIFICATE_VERIFY_FAILED` errors
- Connection refused errors

**Solutions:**
1. Verify SSL certificate:
```bash
nessi ssl verify --cert-path certs/certificate.crt
```

2. Generate new certificate:
```bash
nessi ssl generate --common-name your-domain.com
```

3. Update SSL configuration:
```yaml
ssl:
  verify: false  # For development only
  cert_path: /path/to/cert.crt
  key_path: /path/to/key.key
```

#### API Connection Issues
**Symptoms:**
- `Connection refused` errors
- Timeout errors

**Solutions:**
1. Check service status:
```bash
nessi status
```

2. Verify network connectivity:
```bash
curl -v https://api.nessi.dev/health
```

3. Check firewall rules:
```bash
# Allow Nessi.dev ports
sudo ufw allow 8443/tcp
sudo ufw allow 9090/tcp  # Prometheus
sudo ufw allow 3000/tcp  # Grafana
```

### Data Quality Issues

#### Low Quality Scores
**Symptoms:**
- Quality scores below thresholds
- Frequent quality check failures

**Solutions:**
1. Check quality rules:
```bash
nessi quality rules list --table-path s3://bucket/table
```

2. Analyze quality metrics:
```bash
nessi quality analyze --table-path s3://bucket/table
```

3. Adjust quality thresholds:
```yaml
quality:
  thresholds:
    completeness: 0.95
    accuracy: 0.90
```

#### Schema Evolution Issues
**Symptoms:**
- Schema merge failures
- Incompatible data types

**Solutions:**
1. Check schema history:
```bash
nessi table schema-history --table-path s3://bucket/table
```

2. Force schema evolution:
```python
df.write.format("delta") \
    .option("mergeSchema", "true") \
    .option("overwriteSchema", "true") \
    .mode("append") \
    .save("s3://bucket/table")
```

### Performance Issues

#### Slow Queries
**Symptoms:**
- High query latency
- Resource exhaustion

**Solutions:**
1. Check query plan:
```bash
nessi table explain --query "SELECT * FROM table"
```

2. Optimize partitioning:
```sql
OPTIMIZE table ZORDER BY (frequently_queried_column)
```

3. Adjust cache settings:
```yaml
cache:
  enabled: true
  ttl: 3600
  max_size: 1000
```

#### High Resource Usage
**Symptoms:**
- High CPU/memory usage
- Slow response times

**Solutions:**
1. Monitor resource usage:
```bash
nessi metrics system
```

2. Adjust resource limits:
```yaml
resources:
  cpu: 4
  memory: 8G
  disk: 100G
```

3. Scale services:
```bash
docker-compose scale api=3
```

### Monitoring Issues

#### Alert Configuration
**Symptoms:**
- Missing alerts
- False positives

**Solutions:**
1. Check alert rules:
```bash
nessi alerts rules list
```

2. Test alert conditions:
```bash
nessi alerts test --rule high_error_rate
```

3. Adjust alert thresholds:
```yaml
alerts:
  high_error_rate:
    threshold: 0.1
    window: 5m
```

#### Dashboard Issues
**Symptoms:**
- Missing metrics
- Incorrect visualizations

**Solutions:**
1. Check metric collection:
```bash
nessi metrics list
```

2. Verify dashboard configuration:
```bash
nessi dashboard validate --name data_quality
```

3. Update dashboard settings:
```json
{
  "refresh_interval": "1m",
  "time_range": "24h"
}
```

## Debugging Tools

### Log Analysis
```bash
# View application logs
nessi logs --service api

# Filter logs by level
nessi logs --level error

# Search logs
nessi logs --search "quality check failed"
```

### Metrics Analysis
```bash
# View system metrics
nessi metrics system

# View quality metrics
nessi metrics quality

# Export metrics
nessi metrics export --format csv
```

### Diagnostic Tools
```bash
# Run diagnostics
nessi diagnose

# Check system health
nessi health

# Generate debug report
nessi debug --output debug_report.zip
```

## Support

### Getting Help
1. Check the documentation
2. Search existing issues
3. Contact support:
   - Email: support@nessi.dev
   - Slack: #nessi-support

### Reporting Issues
When reporting issues, include:
1. Nessi.dev version
2. Environment details
3. Error messages
4. Steps to reproduce
5. Relevant logs

Example:
```bash
nessi version
nessi env
nessi logs --service api --level error
``` 