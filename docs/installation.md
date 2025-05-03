# nessi-dev Installation Guide

## Prerequisites

Before installing nessi-dev, ensure you have the following prerequisites:

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

### 2. Build and Start Services

```bash
# Build and start all services
docker-compose up -d

# Check if services are running
docker-compose ps
```

This will start:
- nessi-dev backend service
- Grafana for monitoring (accessible at http://localhost:3000)
- Prometheus for metrics collection

### 3. Verify Installation

To verify that everything is working correctly:

1. Access Grafana dashboard:
   - URL: http://localhost:3000
   - Username: admin
   - Password: admin

2. Run the demo script:
```bash
docker-compose exec backend python src/demo.py
```

## Important Notes

- The application is designed to run exclusively in Docker containers
- Local Python execution is disabled for security and consistency
- All development and testing should be done through Docker containers
- Use `docker-compose exec` to run any Python commands

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

1. **Container startup issues**
   - Check Docker logs: `docker-compose logs`
   - Verify container status: `docker-compose ps`
   - Check resource limits: `docker stats`

2. **Port conflicts**
   - Ensure ports 8000, 9090, and 3000 are available
   - Use `docker-compose down` to stop all containers

3. **Permission issues**
   - Check container user permissions
   - Verify volume mount permissions
   - Ensure proper SELinux/AppArmor settings

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