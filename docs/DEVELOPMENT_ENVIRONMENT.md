# Nessi Development Environment

This document provides instructions for setting up and using the Nessi development environment.

## Overview

The Nessi project includes a comprehensive development environment using Docker Compose, which provides:

- Hot-reloading for rapid development
- Test databases and services for integration testing
- S3-compatible storage for testing cloud storage features
- Grafana for visualization testing

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install) (version 1.21 or later)
- [Git](https://git-scm.com/downloads)

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
```

### Development Environment Setup

The project includes multiple Docker Compose profiles for different development scenarios:

1. **Production-like Environment**:
   ```bash
   docker-compose up
   ```

2. **Development Environment** (with hot-reloading):
   ```bash
   docker-compose --profile dev up
   ```

3. **Testing Environment** (for integration tests):
   ```bash
   docker-compose --profile test up
   ```

### Development Workflow

#### Hot-Reloading

The development environment uses [Air](https://github.com/cosmtrek/air) for hot-reloading. When you make changes to your Go code, the application will automatically rebuild and restart.

The Nessi development server will be available at: http://localhost:8081

#### Database Access

For development and testing, a PostgreSQL database is available:

- Host: `localhost`
- Port: `5432`
- Username: `nessi`
- Password: `nessi`
- Database: `nessi_test`

You can connect to it using any PostgreSQL client:

```bash
psql -h localhost -p 5432 -U nessi -d nessi_test
```

#### S3-Compatible Storage

For testing S3 storage features, MinIO is included:

- S3 Endpoint: `http://localhost:9000`
- Console: `http://localhost:9001`
- Access Key: `minioadmin`
- Secret Key: `minioadmin`
- Bucket: `nessi-test`

#### Grafana

For testing visualization integrations, Grafana is available at:

- URL: `http://localhost:3000`
- Username: `admin`
- Password: `admin`

### Running Tests

#### Unit Tests

To run unit tests locally:

```bash
go test -v -short ./...
```

#### Integration Tests

To run integration tests locally:

1. Start the test environment:
   ```bash
   docker-compose --profile test up -d
   ```

2. Run the integration tests:
   ```bash
   go test -v -tags=integration ./...
   ```

#### All Tests with Coverage

To run all tests and generate a coverage report:

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Then open `coverage.html` in your browser.

## CI/CD Integration

The development environment is designed to match the CI/CD pipeline as closely as possible. The same tests that run in your local environment will run in the CI pipeline.

For more information on the CI/CD workflows, see [CI_CD_WORKFLOWS.md](CI_CD_WORKFLOWS.md).

## Docker Compose Services

### Production Service (`nessi`)

This service runs Nessi in a production-like environment.

- Port: `8080`
- Configuration: `./config/config.yaml`
- Data: `./data`
- Logs: `./logs`
- Certificates: `./certs`

### Development Service (`nessi-dev`)

This service runs Nessi with hot-reloading for development.

- Port: `8081`
- Configuration: `./config/config.yaml`
- Source code mounted from host
- Hot-reloading enabled

### Test Database (`test-db`)

PostgreSQL database for integration testing.

- Port: `5432`
- Username: `nessi`
- Password: `nessi`
- Database: `nessi_test`

### MinIO (`minio`)

S3-compatible storage for testing cloud storage features.

- API Port: `9000`
- Console Port: `9001`
- Access Key: `minioadmin`
- Secret Key: `minioadmin`

### MinIO Init (`minio-init`)

Creates the necessary buckets in MinIO for testing.

- Bucket: `nessi-test`
- Policy: `public`

### Grafana (`grafana`)

Visualization platform for testing dashboard integrations.

- Port: `3000`
- Anonymous access enabled for testing

## Customizing the Environment

### Modifying the Docker Compose Configuration

You can modify the `docker-compose.yml` file to add or change services as needed for your development workflow.

### Adding New Services

To add a new service to the development environment:

1. Add the service configuration to `docker-compose.yml`
2. Add the service to the appropriate profiles (`dev`, `test`, or both)
3. Update this documentation to include information about the new service

### Environment Variables

You can customize the environment by setting environment variables in the `docker-compose.yml` file or by creating a `.env` file in the project root.

## Troubleshooting

### Common Issues

#### Port Conflicts

If you encounter port conflicts, you can change the port mappings in the `docker-compose.yml` file.

#### Docker Volume Permissions

If you encounter permission issues with Docker volumes, you may need to adjust the permissions on your host machine:

```bash
sudo chown -R $(id -u):$(id -g) ./data ./logs ./config ./certs
```

#### Hot-Reloading Not Working

If hot-reloading is not working:

1. Check the Air logs in the Docker container:
   ```bash
   docker logs nessi-dev
   ```

2. Ensure your `.air.toml` configuration is correct
3. Try restarting the development container:
   ```bash
   docker-compose --profile dev restart nessi-dev
   ```

## Best Practices

1. **Use the Development Environment**: Always use the provided development environment to ensure consistency across all developers.

2. **Write Tests**: Write unit and integration tests for all new features and bug fixes.

3. **Check Test Coverage**: Ensure your changes maintain or improve test coverage.

4. **Follow the CI/CD Workflow**: Make sure your changes pass all CI checks before submitting a pull request.

5. **Document Changes**: Update documentation when adding or changing features.

## Contributing

For more information on contributing to Nessi, see [CONTRIBUTING.md](../CONTRIBUTING.md).
