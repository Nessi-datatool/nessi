# Nessi - Data Quality and Profiling Tool

Nessi is a powerful data quality and profiling tool that helps you analyze and understand your data. It supports various data formats including Delta Lake, Parquet, and CSV.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/nessi.git
cd nessi
```

2. Build and run the Docker container:
```bash
docker-compose up --build
```

This will:
- Build the Nessi Docker image
- Start the container
- Run the test suite
- Make the service available on port 8080

## Configuration

The following environment variables can be configured in `docker-compose.yml`:

- `SPARK_CONF_DIR`: Directory containing Spark configuration files
- `SPARK_JARS_DIR`: Directory containing Spark and Delta Lake JARs
- `PYTHONPATH`: Python path configuration

## Volumes

The following directories are mounted as volumes:

- `/app/data`: Data directory
- `/app/jars`: JAR files directory
- `/app/conf`: Configuration files directory

## Development

To run tests:
```bash
docker-compose run nessi pytest tests/
```

To run a specific test:
```bash
docker-compose run nessi pytest tests/test_file.py::test_function
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 