# Nessi

Nessi is a powerful data analysis and processing tool, free for personal use. It provides comprehensive data scanning, reporting, and visualization capabilities.

## License

Nessi is free for personal use. A paid license is required for enterprise or consulting use. This version is free, while future versions will require a license for all users.

For more details, see [LICENSE](LICENSE).

## Quick Start

### Prerequisites

- Python 3.8 or higher
- Docker and Docker Compose (for running with containers)
- Java 8 or higher (for Spark support)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Nessi-datatool/nessi.git
cd nessi
```

2. Create and activate a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

3. Install dependencies:
```bash
pip install -e .
```

### Running with Docker

1. Build and start the services:
```bash
docker-compose up -d
```

2. Access the web interface at `http://localhost:8080`

### Basic Usage

1. Start the Nessi CLI:
```bash
nessi
```

2. Check your usage status:
```bash
nessi status
```

3. Verify personal use:
```bash
nessi check
```

## Features

- Data scanning and analysis
- Custom report generation
- Data quality metrics
- Performance monitoring
- Grafana dashboard integration
- Delta Lake support
- Parquet and CSV file handling

## Documentation

- [API Documentation](API.md)
- [Security Guidelines](SECURITY.md)
- [Contributing Guide](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

## Support

For personal use support, please open an issue in the GitHub repository.

For enterprise support, please contact us at support@nessi.dev

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## Security

Please report any security issues to security@nessi.dev. See our [Security Policy](SECURITY.md) for more information. 