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

### 1. Container Startup Issues

#### Containers Won't Start

**Symptoms:**
- `docker-compose up -d` fails
- Containers exit immediately
- Services are not accessible

**Solutions:**

1. Check port conflicts:
```bash
# Check if ports 3000 or 9090 are in use
lsof -i :3000
lsof -i :9090
```

2. Check Docker logs:
```bash
docker-compose logs backend
docker-compose logs grafana
docker-compose logs prometheus
```

3. Verify Docker resources:
```bash
# Check Docker disk space
docker system df

# Check memory limits
docker stats
```

#### Permission Issues

**Symptoms:**
- "Permission denied" errors in logs
- Can't write to mounted volumes

**Solutions:**

1. Check volume permissions:
```bash
# Fix permissions on data directories
sudo chown -R 1000:1000 ./data
sudo chown -R 1000:1000 ./metrics
```

2. Verify Docker socket permissions:
```bash
sudo chmod 666 /var/run/docker.sock
```

### 2. Memory Issues

**Symptoms:**
- Container crashes with OOM (Out of Memory)
- Spark jobs fail
- Slow performance

**Solutions:**

1. Check container memory usage:
```bash
docker stats
```

2. Adjust memory limits in docker-compose.yml:
```yaml
services:
  backend:
    mem_limit: 4g
    mem_reservation: 2g
```

3. Modify Spark memory settings:
```bash
docker-compose exec backend pyspark-shell --driver-memory 2g --executor-memory 2g
```

### 3. Network Issues

**Symptoms:**
- Services can't communicate
- "Connection refused" errors
- Grafana can't reach Prometheus

**Solutions:**

1. Check network connectivity:
```bash
# Verify network exists
docker network ls

# Inspect network
docker network inspect nessi-network
```

2. Test inter-container communication:
```bash
docker-compose exec backend ping grafana
docker-compose exec backend ping prometheus
```

3. Verify service URLs in environment:
```bash
docker-compose exec backend env | grep URL
```

### 4. Data Persistence Issues

**Symptoms:**
- Data disappears after container restart
- Volume mount errors
- Can't write to data directories

**Solutions:**

1. Check volume mounts:
```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect nessi_grafana-storage
```

2. Verify directory structure:
```bash
# Check directory permissions
ls -la ./data
ls -la ./metrics

# Create directories if missing
mkdir -p ./data/{parquet,csv,delta}
mkdir -p ./metrics
chmod -R 777 ./data ./metrics
```

### 5. Grafana Issues

**Symptoms:**
- Can't access Grafana UI
- Dashboard not loading
- Authentication issues

**Solutions:**

1. Reset Grafana:
```bash
# Remove Grafana volume
docker-compose down
docker volume rm nessi_grafana-storage
docker-compose up -d
```

2. Check Grafana logs:
```bash
docker-compose logs grafana
```

3. Verify Grafana configuration:
```bash
docker-compose exec grafana cat /etc/grafana/grafana.ini
```

### 6. Prometheus Issues

**Symptoms:**
- Metrics not showing
- Target scraping fails
- Data gaps in graphs

**Solutions:**

1. Check Prometheus targets:
```bash
curl http://localhost:9090/api/v1/targets
```

2. Verify scrape configuration:
```bash
docker-compose exec prometheus cat /etc/prometheus/prometheus.yml
```

3. Check Prometheus logs:
```bash
docker-compose logs prometheus
```

## Development Issues

### 1. Testing Problems

**Symptoms:**
- Tests fail in container
- Import errors
- Path issues

**Solutions:**

1. Run tests with proper Python path:
```bash
docker-compose exec backend pytest -v tests/
```

2. Check test environment:
```bash
docker-compose exec backend python -c "import sys; print(sys.path)"
```

### 2. Code Changes Not Reflected

**Symptoms:**
- Code updates not visible in container
- Old version still running

**Solutions:**

1. Rebuild container:
```bash
docker-compose build backend
docker-compose up -d backend
```

2. Verify file mounting:
```bash
docker-compose exec backend ls -la /app/backend
```

## Maintenance Tasks

### 1. Cleaning Up

```bash
# Remove unused containers
docker-compose down

# Clean up volumes
docker volume prune

# Remove unused images
docker image prune

# Full cleanup
docker system prune -a
```

### 2. Updating Services

```bash
# Pull latest images
docker-compose pull

# Rebuild with no cache
docker-compose build --no-cache

# Restart services
docker-compose down && docker-compose up -d
```

## Getting Help

If you're still experiencing issues:

1. Check the logs:
```bash
docker-compose logs --tail=100
```

2. Gather system information:
```bash
docker version
docker-compose version
docker info
```

3. Create an issue with:
- Full error message
- Container logs
- System information
- Steps to reproduce 