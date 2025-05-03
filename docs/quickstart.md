# Quick Start Guide

This guide will help you get started with nessi-dev quickly. nessi-dev is free for personal use, while enterprise or consulting use requires a paid license.

## Installation

### Prerequisites

- Docker Engine 20.10.0 or later
- Docker Compose v2.0.0 or later
- Git (for cloning the repository)
- At least 4GB of available RAM
- At least 10GB of free disk space

### Quick Installation

1. Clone the repository:
```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
```

2. Start the services:
```bash
docker-compose up -d
```

3. Access the services:
- Backend API: https://localhost:8000
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

### Docker Installation

If you prefer using Docker:

1. Build and start the services:
```bash
docker-compose up -d
```

2. Access the web interface at `http://localhost:8080`

## Basic Usage

### Command Line Interface

1. Start the nessi-dev CLI:
```bash
docker-compose exec backend nessi
```

2. Check your usage status:
```bash
docker-compose exec backend nessi status
```

3. Verify personal use:
```bash
docker-compose exec backend nessi check
```

### Data Scanning

1. Scan a Delta Lake table:
```bash
docker-compose exec backend python -c "
from nessi.scanner import TableScanner
scanner = TableScanner()
results = scanner.scan_table('path/to/your/table')
"
```

2. Generate a report:
```bash
docker-compose exec backend python -c "
from nessi.scanner import TableScanner
scanner = TableScanner()
report = scanner.generate_report(results)
print(report)
"
```

### Data Quality Checks

1. Run data quality checks:
```bash
docker-compose exec backend python -c "
from nessi.scanner import TableScanner
scanner = TableScanner()
quality_metrics = scanner.check_data_quality('path/to/your/table')
"
```

2. View quality metrics:
```bash
docker-compose exec backend python -c "
print(quality_metrics)
"
```

## Next Steps

- Read the [Features Guide](features.md) for detailed feature information
- Check out the [Reports Guide](reports.md) for report customization
- Visit the [Troubleshooting Guide](troubleshooting.md) if you encounter any issues
- Review the [Monitoring Guide](monitoring.md) for performance monitoring

## License Information

Remember:
- nessi-dev is free for personal use
- Enterprise or consulting use requires a paid license
- Future versions will require a license for all users

For more details, see [LICENSE](../LICENSE). 