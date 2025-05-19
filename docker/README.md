# Nessi.dev Docker Deployment

This document explains how to deploy Nessi.dev using Docker and Docker Compose.

## Quick Start

The easiest way to get started with Nessi.dev is to use the provided Docker Compose configuration:

```bash
# Clone the repository
git clone https://github.com/nessi-dev/nessi-dev.git
cd nessi-dev

# Start all services
docker-compose up -d

# Check the status of the services
docker-compose ps
```

This will start the following services:
- **Nessi**: The main application (accessible at http://localhost:8080)
- **Prometheus**: Metrics collection (accessible at http://localhost:9090)
- **Grafana**: Dashboards and visualization (accessible at http://localhost:3000)
- **Pushgateway**: Metrics pushing (accessible at http://localhost:9091)

## No JVM or Spark Dependencies

Nessi.dev is built with a Go-first architecture, which means you don't need to install or configure:
- Java Virtual Machine (JVM)
- Apache Spark
- Hadoop
- Any other Java-based dependencies

All necessary components are included in the Docker images, making deployment simple and lightweight.

## Databricks Integration

Nessi.dev includes integration with Databricks, allowing you to:
- Connect to Databricks workspaces via the REST API
- Browse Unity Catalog resources (catalogs, schemas, and tables)
- Read and write Delta Lake tables
- Perform time travel queries on Delta Lake tables

To use the Databricks integration, set the following environment variables:

```yaml
environment:
  - DATABRICKS_HOST=your-databricks-workspace.cloud.databricks.com
  - DATABRICKS_TOKEN=your-personal-access-token
  - DATABRICKS_WORKSPACE_ID=your-workspace-id
  - DATABRICKS_DEFAULT_CATALOG=main
  - DATABRICKS_DEFAULT_SCHEMA=default
```

## Configuration

### Environment Variables

You can customize the Nessi.dev deployment by setting environment variables in the `docker-compose.yml` file:

```yaml
environment:
  - CONFIG_PATH=/app/config/config.yaml
  - TZ=UTC
  # Add your custom environment variables here
```

### Volumes

The Docker Compose configuration mounts several volumes to persist data and configuration:

```yaml
volumes:
  - ./config:/app/config     # Configuration files
  - ./data:/app/data         # Data files and Delta tables
  - ./logs:/app/logs         # Log files
  - ./certs:/app/certs       # SSL/TLS certificates
```

## Multi-Service Orchestration

The Docker Compose configuration orchestrates multiple services that work together:

1. **Nessi**: The main application that provides the API and core functionality
2. **Prometheus**: Collects and stores metrics from Nessi
3. **Grafana**: Provides dashboards for visualizing metrics and alerts
4. **Pushgateway**: Allows pushing metrics to Prometheus

Each service is configured with:
- Health checks to ensure availability
- Restart policies for reliability
- Proper networking to communicate with other services
- Volume mounts for persistence

## Custom Builds

If you need to customize the Nessi.dev Docker image, you can modify the `Dockerfile` and build your own image:

```bash
# Build the Docker image
docker build -t nessi:custom .

# Update the docker-compose.yml file to use your custom image
# Then start the services
docker-compose up -d
```

## Production Deployment

For production deployments, consider the following:

1. Use a reverse proxy (like Nginx or Traefik) for TLS termination
2. Set up proper authentication for Grafana and Prometheus
3. Configure persistent volumes for data storage
4. Set up monitoring and alerting for the Docker containers themselves

Example production `docker-compose.override.yml`:

```yaml
version: '3.8'

services:
  nessi:
    restart: always
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G

  prometheus:
    restart: always
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 2G

  grafana:
    restart: always
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=your-secure-password
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
```

## Troubleshooting

If you encounter issues with the Docker deployment, check the following:

1. Container logs:
   ```bash
   docker-compose logs nessi
   ```

2. Container health:
   ```bash
   docker-compose ps
   ```

3. Network connectivity:
   ```bash
   docker network inspect nessi-net
   ```

4. Volume mounts:
   ```bash
   docker volume ls
   ```
