# NESSI - Docker Development Environment

This project uses Docker for development and testing. Local Python environments are not supported.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. Build and run the Docker container:
```bash
docker-compose up --build
```

2. To run tests:
```bash
docker-compose up
```

3. To run a specific test:
```bash
docker-compose run nessi python -m pytest backend/tests/test_file.py -v
```

## Development

All development should be done within the Docker container. The container mounts the current directory, so changes to files will be reflected immediately.

### Running Tests

Tests are automatically run when starting the container. To run specific tests:

```bash
docker-compose run nessi python -m pytest backend/tests/test_file.py -v
```

### Accessing the Container

To access the container shell:
```bash
docker-compose run nessi /bin/bash
```

## Important Notes

- Do not use local Python environments or virtual environments
- All development and testing must be done within Docker
- The container mounts the current directory, so changes are reflected immediately
- Use `docker-compose` commands for all operations 