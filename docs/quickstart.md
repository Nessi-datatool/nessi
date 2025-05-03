# Quick Start Guide

This guide will help you get started with nessi.dev quickly. nessi.dev is free for personal use, while enterprise or consulting use requires a paid license.

## Installation

### Prerequisites

- Python 3.8 or higher
- Docker and Docker Compose (optional, for containerized deployment)
- Java 8 or higher (for Spark support)

### Quick Installation

1. Clone the repository:
```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
```

2. Create and activate a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
```

3. Install nessi.dev:
```bash
pip install -e .
```

### Docker Installation (Optional)

If you prefer using Docker:

1. Build and start the services:
```bash
docker-compose up -d
```

2. Access the web interface at `http://localhost:8080`

## Basic Usage

### Command Line Interface

1. Start the nessi.dev CLI:
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

### Data Scanning

1. Scan a Delta Lake table:
```python
from nessi.scanner import TableScanner

scanner = TableScanner()
results = scanner.scan_table("path/to/your/table")
```

2. Generate a report:
```python
report = scanner.generate_report(results)
print(report)
```

### Data Quality Checks

1. Run data quality checks:
```python
from nessi.scanner import TableScanner

scanner = TableScanner()
quality_metrics = scanner.check_data_quality("path/to/your/table")
```

2. View quality metrics:
```python
print(quality_metrics)
```

## Next Steps

- Read the [Features Guide](features.md) for detailed feature information
- Check out the [Reports Guide](reports.md) for report customization
- Visit the [Troubleshooting Guide](troubleshooting.md) if you encounter any issues
- Review the [Monitoring Guide](monitoring.md) for performance monitoring

## License Information

Remember:
- nessi.dev is free for personal use
- Enterprise or consulting use requires a paid license
- This version is free
- Future versions will require a license for all users

For more details, see [LICENSE](../LICENSE). 