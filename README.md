# nessi.dev

A modern data tool for data quality monitoring and reporting.

## Prerequisites

- Docker (version 20.10.0 or higher)
- Docker Compose (version 2.0.0 or higher)
- Git

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/nessi.git
   cd nessi
   ```

2. Start the development environment:
   ```bash
   ./scripts/dev.sh
   ```

3. Access the services:
   - Backend API: https://localhost:8000
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000

## Development

### Running Tests

All tests must be run through Docker containers:

```bash
docker-compose -f docker-compose.test.yml up --build
```

### Code Style

We use Black for Python code formatting. The code style is enforced through Docker:

```bash
docker-compose -f docker-compose.dev.yml run backend black .
```

### Documentation

Documentation is built and served through Docker:

```bash
docker-compose -f docker-compose.dev.yml run backend mkdocs serve
```

## Important Notes

- The application is designed to run exclusively in Docker containers
- Local Python execution is disabled for security and consistency
- All development should occur within Docker containers
- The development environment includes hot-reloading for faster development

## Security

- All containers run as non-root users
- SSL/TLS is enabled by default
- Authentication is required for all services
- Resource limits are enforced through Docker

## Troubleshooting

### Common Issues

1. **Port conflicts**
   - Ensure ports 8000, 9090, and 3000 are not in use
   - Use `docker-compose -f docker-compose.dev.yml down` to stop all containers

2. **SSL certificate issues**
   - The development script automatically generates SSL certificates
   - For production, replace the certificates in the `ssl` directory

3. **Container startup issues**
   - Check Docker logs: `docker-compose -f docker-compose.dev.yml logs`
   - Ensure all required environment variables are set

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 