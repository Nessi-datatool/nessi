# Nessi.dev Docker Deployment

This document explains how to deploy Nessi.dev using Docker and Docker Compose.

## Quick Start

The easiest way to get started with Nessi.dev is to use the provided Docker Compose configuration:

```bash
# Clone the repository
git clone git@github.com:nessi-dev/nessi.git
cd nessi

# Start all services
docker-compose up -d

# Check the status of the services
docker-compose ps
```

This will start the Nessi CLI application which provides:
- Data catalog browsing and management
- Delta Lake format handling for reading and writing data
- Data quality validation and monitoring
- Integration with Databricks and other data sources

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

You can customize the Nessi CLI by setting environment variables in the `docker-compose.yml` file:

```yaml
environment:
  - DATABRICKS_HOST=your-databricks-workspace.cloud.databricks.com
  - DATABRICKS_TOKEN=your-personal-access-token
  - DATABRICKS_WORKSPACE_ID=your-workspace-id
  - DATABRICKS_DEFAULT_CATALOG=main
  - DATABRICKS_DEFAULT_SCHEMA=default
  - LOG_LEVEL=info
```

### Volumes

The Docker Compose configuration mounts several volumes to persist data and configuration:

```yaml
volumes:
  - ./config:/app/config     # Configuration files
  - ./data:/app/data         # Data files and Delta tables
  - ./logs:/app/logs         # Log files
```

## CLI-Only Approach

The Docker configuration provides a lightweight CLI-only approach that doesn't require additional services:

1. **Nessi CLI**: The main application that provides all functionality
2. **No external dependencies**: No need for separate metrics or visualization services
3. **Minimal resource usage**: Runs efficiently even on limited hardware

The container is configured with:
- Proper volume mounts for data persistence
- Environment variable configuration
- Lightweight footprint for easy deployment

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

1. Use a reverse proxy (like Nginx or Traefik) for TLS termination if needed
2. Configure persistent volumes for data storage
3. Set up appropriate environment variables for your deployment

Example production `docker-compose.override.yml`:

```yaml
version: '3.8'

services:
  nessi:
    restart: always
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
      - DATABRICKS_HOST=your-databricks-workspace.cloud.databricks.com
      - DATABRICKS_TOKEN=your-personal-access-token
      - DATABRICKS_WORKSPACE_ID=your-workspace-id
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

## Troubleshooting

If you encounter issues with the Docker deployment, check the following:

1. Container logs:
   ```bash
   docker-compose logs
   ```

2. Container status:
   ```bash
   docker-compose ps
   ```

3. Check environment variables:
   ```bash
   docker-compose exec nessi env | grep DATABRICKS
   ```

4. Verify volume mounts:
   ```bash
   docker inspect nessi | grep Mounts -A 10
   ```
