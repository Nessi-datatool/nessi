# Nessi Installation Guide

## Prerequisites

Before installing Nessi, ensure you have the following prerequisites:

- Docker Engine 20.10.0 or later
- Docker Compose v2.0.0 or later
- Git (for cloning the repository)
- At least 4GB of available RAM
- At least 10GB of free disk space

## Installation Steps

### 1. Clone the Repository

```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
```

### 2. Environment Setup

The project is designed to run entirely in Docker containers, so no local Python installation is required. However, if you want to develop or run tests locally, you'll need:

- Python 3.11 or later
- Java 17 (for Apache Spark)
- Apache Spark 3.5.0

### 3. Build and Start Services

```bash
# Build and start all services
docker-compose up -d

# Check if services are running
docker-compose ps
```

This will start:
- Nessi backend service
- Grafana for monitoring (accessible at http://localhost:3000)
- Prometheus for metrics collection

### 4. Verify Installation

To verify that everything is working correctly:

1. Access Grafana dashboard:
   - URL: http://localhost:3000
   - Username: admin
   - Password: admin

2. Run the demo script:
```bash
docker-compose exec backend python src/demo.py
```

## Configuration

### Environment Variables

Key environment variables that can be configured in `docker-compose.yml`:

```yaml
environment:
  - PYTHONPATH=/app/backend
  - JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64
  - PYSPARK_PYTHON=/usr/local/bin/python
  - PYSPARK_DRIVER_PYTHON=/usr/local/bin/python
  - GRAFANA_URL=http://grafana:3000
  - PROMETHEUS_URL=http://prometheus:9090
```

### Volume Mounts

Important volume mounts in `docker-compose.yml`:

```yaml
volumes:
  - ./backend:/app/backend
  - ./data:/app/data
  - ./metrics:/app/metrics
```

## Troubleshooting

### Common Issues

1. **Container fails to start**
   - Check Docker logs: `docker-compose logs`
   - Ensure all ports (3000, 9090) are available
   - Verify sufficient system resources

2. **Permission Issues**
   - Ensure proper permissions on mounted volumes
   - Check container logs for permission errors

3. **Memory Issues**
   - Increase Docker memory limit
   - Check Spark memory configuration

### Getting Help

If you encounter any issues:
1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Search existing GitHub issues
3. Create a new issue with detailed information about your problem

## Next Steps

- Follow the [Quick Start Guide](quickstart.md)
- Read the [Features Documentation](features.md)
- Explore the [API Documentation](../API.md)
- Learn about [Contributing](../CONTRIBUTING.md) 