# Nessi.dev

Nessi.dev is a powerful data quality and Delta Lake management tool that helps organizations maintain high-quality data and efficiently manage their Delta Lake tables.

## Features

- **Delta Lake Support**
  - Native Delta Lake table management
  - Transaction log parsing and analysis
  - Schema evolution tracking
  - Partition management
  - Version control and time travel
  - Multi-format support (Delta, Parquet, CSV)
  - Automatic schema inference and validation
  - Cloud integration with AWS, Azure (SDK v1.6.1+), and GCP
  - **Data Catalog Integration**: Connect with popular data catalogs like AWS Glue, Azure Purview, and Google Cloud Data Catalog to discover data assets and publish quality metrics

- **Data Quality Intelligence**
  - Automated data quality checks
  - Custom quality rules
  - Data profiling and statistics
  - Anomaly detection
  - Quality scoring and metrics

- **Monitoring and Alerting**
  - Real-time monitoring with Prometheus
  - Custom metrics and alerts
  - Performance tracking
  - Resource utilization monitoring
  - Health checks
  - Intelligent alerting with pattern recognition
  - Automated rule creation based on historical data
  - Anomaly and trend deviation detection

- **Report Generation**
  - Customizable report templates
  - Multiple output formats (PDF, CSV, HTML, JSON)
  - Scheduled report generation
  - Report archiving and retention
  - Data quality dashboards

- **Security**
  - JWT-based authentication
  - Role-based access control (RBAC)
  - TLS encryption
  - Rate limiting
  - IP allowlist
  - Security headers

## Developer Experience

- **CLI-First Approach**
  - One-line validation commands
  - Batch processing scripts
  - Comprehensive configuration management
  - Intuitive command structure

- **Developer Experience**
  - Intuitive CLI interface
  - Comprehensive documentation
  - Python API for programmatic access
  - Extensible architecture with plugin system
  - Custom extension support via plugins
  - Embeddable in Airflow, GitHub Actions, Kubernetes, and custom scripts
  - Interactive examples and integration guides
  - Observability features with metrics collection and structured logging

- **Integration Capabilities**
  - Prometheus metrics integration
  - Grafana dashboards
  - Secure API access
  - Comprehensive webhook integration
  - Event-based notifications
  - Cloud provider integration (AWS, Azure, GCP)
  - Modern Azure SDK (v1.6.1+) support
  - **Workflow Orchestration Integration**: Seamless integration with Apache Airflow, Prefect, Dagster, and Kubernetes for incorporating data quality checks into data pipelines
  - **Kubernetes Integration**: Native operators for running Nessi operations as Kubernetes jobs with configurable resources and monitoring

- **Documentation**
  - Comprehensive user guides
  - API reference
  - Example scripts
  - Best practices

## Requirements

- Go 1.21 or later
- Docker (optional)
- Make
- OpenSSL (for certificate generation)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/nessi-dev/nessi-dev.git
   cd nessi-dev
   ```

2. Install development tools:
   ```bash
   make install-tools
   ```

3. Initialize the development environment:
   ```bash
   make init
   ```

4. Build the application:
   ```bash
   make build
   ```

## Configuration

The application is configured using a YAML file located at `config/config.yaml`. The configuration includes settings for:

- Server configuration (host, port, timeouts)
- Security settings (JWT, TLS, rate limiting)
- Delta Lake paths and properties
- Data quality rules and thresholds
- Monitoring and alerting
- Report generation
- Logging

See `config/config.yaml` for detailed configuration options.

## Usage

### Running the Application

1. Start the server:
   ```bash
   make run
   ```

2. Or run with Docker:
   ```bash
   make docker-build
   make docker-run
   ```

### API Endpoints

The application exposes the following REST API endpoints:

- `GET /api/v1/tables` - List Delta tables
- `GET /api/v1/tables/{name}` - Get table details
- `POST /api/v1/tables/{name}/quality` - Run quality checks
- `GET /api/v1/tables/{name}/quality` - Get quality metrics
- `GET /api/v1/tables/{name}/reports` - Get quality reports
- `POST /api/v1/tables/{name}/reports` - Generate new report

### Development

1. Run tests:
   ```bash
   make test
   ```

2. Run linters:
   ```bash
   make lint
   ```

3. Clean build artifacts:
   ```bash
   make clean
   ```

## Security

### TLS Configuration

1. Generate self-signed certificates:
   ```bash
   make generate-certs
   ```

2. Update the configuration in `config/config.yaml`:
   ```yaml
   security:
     tls:
       enabled: true
       cert_file: "certs/server.crt"
       key_file: "certs/server.key"
   ```

### Authentication

The application uses JWT for authentication. Configure the JWT settings in `config/config.yaml`:

```yaml
security:
  jwt:
    secret: "your-secret-key"
    expiration: 24h
```

## Monitoring

The application integrates with Prometheus for monitoring. Configure the monitoring settings in `config/config.yaml`:

```yaml
monitoring:
  prometheus:
    enabled: true
    push_gateway: "http://localhost:9091"
    job_name: "nessi"
    interval: 15s
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers. 