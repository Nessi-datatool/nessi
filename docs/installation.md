# Installation Guide

## Prerequisites

- Docker Engine (version 20.10.0 or later)
- Docker Compose (version 2.0.0 or later)
- Git
- 4GB RAM minimum (8GB recommended)
- 10GB free disk space

## Quick Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Nessi-datatool/nessi.git
   cd nessi
   ```

2. Build and start the containers:
   ```bash
   docker-compose build
   docker-compose up -d
   ```

3. Verify the installation:
   ```bash
   docker-compose ps
   ```

## Detailed Installation

### 1. System Requirements

- Operating System: Linux, macOS, or Windows (with WSL2)
- CPU: 2 cores minimum (4 cores recommended)
- Memory: 4GB minimum (8GB recommended)
- Storage: 10GB free space
- Network: Internet access for downloading images

### 2. Docker Setup

1. Install Docker Engine:
   - [Linux](https://docs.docker.com/engine/install/)
   - [macOS](https://docs.docker.com/desktop/install/mac-install/)
   - [Windows](https://docs.docker.com/desktop/install/windows-install/)

2. Install Docker Compose:
   ```bash
   # Linux
   sudo apt-get install docker-compose-plugin

   # macOS/Windows (included with Docker Desktop)
   ```

3. Verify Docker installation:
   ```bash
   docker --version
   docker-compose --version
   ```

### 3. Project Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/Nessi-datatool/nessi.git
   cd nessi
   ```

2. Configure environment (optional):
   ```bash
   cp .env.example .env
   # Edit .env file with your settings
   ```

3. Build the containers:
   ```bash
   docker-compose build
   ```

4. Start the services:
   ```bash
   docker-compose up -d
   ```

### 4. Verification

1. Check container status:
   ```bash
   docker-compose ps
   ```

2. Check logs:
   ```bash
   docker-compose logs -f
   ```

3. Test the API:
   ```bash
   curl http://localhost:8000/health
   ```

## Development Setup

For development purposes, use the development configuration:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

This will:
- Mount the source code as a volume
- Enable hot-reloading
- Provide additional debugging tools
- Run tests automatically

## Testing

All tests are run in Docker containers. Use the provided script:

```bash
# Make the script executable
chmod +x docker/run-tests.sh

# Run all tests
./docker/run-tests.sh
```

Test results will be available in the `test-results` directory.

## Updating

To update to the latest version:

```bash
git pull
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Troubleshooting

Common issues and solutions:

1. **Port Conflicts**
   - Check if ports 8000, 3000, 9090 are available
   - Modify ports in docker-compose.yml if needed

2. **Build Failures**
   - Clear Docker cache: `docker system prune -a`
   - Check network connectivity
   - Verify Docker daemon is running

3. **Container Issues**
   - Check logs: `docker-compose logs -f`
   - Restart containers: `docker-compose restart`
   - Rebuild containers: `docker-compose build --no-cache`

4. **Test Issues**
   - Check test-results directory for detailed reports
   - Verify service connectivity in Docker network
   - Ensure sufficient resources for Spark

## Next Steps

- [Quick Start Guide](quickstart.md)
- [Configuration Guide](configuration.md)
- [API Reference](api-reference.md) 