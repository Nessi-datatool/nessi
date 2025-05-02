# Nessi - Data Quality and Profiling Tool

Nessi is a powerful data quality and profiling tool that helps you analyze and understand your data. It supports various data formats including Delta Lake, Parquet, and CSV.

## Prerequisites

- Docker Engine 20.10.0 or later
- Docker Compose v2.0.0 or later
- At least 4GB of available RAM
- At least 10GB of free disk space

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/Nessi-datatool/nessi.git
cd nessi
```

2. Build and start the services:
```bash
docker-compose up -d
```

This will start:
- Nessi backend service
- Grafana for monitoring (http://localhost:3000)
- Prometheus for metrics collection

3. Verify the installation:
```bash
docker-compose exec backend python src/demo.py
```

## Services

The following services are available:

- **Nessi Backend**: Data processing and analysis
- **Grafana**: Monitoring and visualization (http://localhost:3000)
- **Prometheus**: Metrics collection

## Configuration

The following environment variables can be configured in `docker-compose.yml`:

```yaml
environment:
  - PYTHONPATH=/app/backend
  - JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64
  - PYSPARK_PYTHON=/usr/local/bin/python
  - PYSPARK_DRIVER_PYTHON=/usr/local/bin/python
  - GRAFANA_URL=http://grafana:3000
  - PROMETHEUS_URL=http://prometheus:9090
```

## Volumes

The following directories are mounted as volumes:

- `/app/backend`: Backend code
- `/app/data`: Data directory
- `/app/metrics`: Metrics storage
- `/app/conf`: Configuration files

## Development

All development and testing should be done within Docker containers:

1. Run tests:
```bash
docker-compose exec backend pytest tests/
```

2. Run a specific test:
```bash
docker-compose exec backend pytest tests/test_file.py::test_function
```

3. Access logs:
```bash
docker-compose logs -f backend
```

## Documentation

For detailed documentation, see the [docs](docs/) directory:
- [Installation Guide](docs/installation.md)
- [Quick Start Guide](docs/quickstart.md)
- [Demo Tutorial](docs/demo-tutorial.md)
- [Monitoring Guide](docs/monitoring.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 