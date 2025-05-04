# Demo Tutorial

## Overview

This tutorial demonstrates the capabilities of nessi.dev through practical examples:
- Monitoring setup
- Security features
- Basic metrics collection

## Prerequisites

- Docker (version 20.10.0 or higher)
- Docker Compose (version 2.0.0 or higher)
- Git
- Minimum 4GB RAM
- 10GB free disk space
- Python 3.8 or higher (for development)
- Virtual environment (recommended)

## Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-org/nessi-dev.git
   cd nessi-dev
   ```

2. **Start Services**:
   ```bash
   docker-compose up -d
   ```

3. **Verify Installation**:
   ```bash
   docker-compose ps
   docker-compose logs -f
   ```

## Demo Steps

### 1. Monitoring Setup

1. **Access Grafana Dashboard**:
   - URL: http://localhost:3000
   - Default credentials: admin/admin

2. **Import Dashboards**:
   ```bash
   # Copy dashboard JSON files
   cp grafana/dashboards/*.json /path/to/grafana/provisioning/dashboards/
   ```

3. **Configure Data Sources**:
   - Add Prometheus as a data source
   - Configure alert rules
   - Set up notification channels

### 2. Security Configuration

1. **Basic Authentication Setup**:
   ```bash
   # Configure basic authentication for metrics
   python scripts/configure_auth.py --type basic \
     --username admin \
     --password admin
   ```

2. **Rate Limiting Setup**:
   ```bash
   # Configure rate limiting
   python scripts/configure_rate_limit.py \
     --rate 10 \
     --burst 20
   ```

3. **SSL Configuration**:
   ```bash
   # Generate SSL certificates
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
     -keyout private.key -out certificate.crt
   ```

## Next Steps

- [Quick Start Guide](quickstart.md)
- [Features Overview](features.md)
- [API Documentation](../API.md)
- [Troubleshooting Guide](troubleshooting.md) 