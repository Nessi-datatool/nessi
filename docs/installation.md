# Installation Guide

## Prerequisites

- Docker (version 20.10.0 or higher)
- Docker Compose (version 2.0.0 or higher)
- Git
- Minimum 4GB RAM
- 10GB free disk space
- Python 3.8 or higher (for development installation)
- Virtual environment (recommended for development)

## Installation Methods

### 1. Docker Installation (Recommended)

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-org/nessi-dev.git
   cd nessi-dev
   ```

2. **Configure Environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Services**:
   ```bash
   docker-compose up -d
   ```

4. **Verify Installation**:
   ```bash
   docker-compose ps
   docker-compose logs -f
   ```

### 2. Development Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-org/nessi-dev.git
   cd nessi-dev
   ```

2. **Create Virtual Environment**:
   ```bash
   python -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```

3. **Install Dependencies**:
   ```bash
   pip install -r requirements.txt
   pip install -r requirements-dev.txt  # For development
   ```

4. **Configure Environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Initialize Database**:
   ```bash
   python scripts/init_db.py
   ```

6. **Start Services**:
   ```bash
   python scripts/start_services.py
   ```

## Configuration

### 1. Environment Variables

```bash
# Spark Configuration
SPARK_MASTER=local[*]
SPARK_DRIVER_MEMORY=2g
SPARK_EXECUTOR_MEMORY=2g
SPARK_EXECUTOR_CORES=2

# Monitoring Configuration
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=admin

# Security Configuration
AUTH_ENABLED=true
SSL_ENABLED=true
METRICS_USERNAME=admin
METRICS_PASSWORD=admin

# Storage Configuration
STORAGE_TYPE=s3  # Options: s3, local, hdfs
S3_BUCKET=your-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key

# Logging Configuration
LOG_LEVEL=INFO
LOG_FORMAT=json
LOG_PATH=/var/log/nessi
```

### 2. SSL Configuration

1. **Generate SSL Certificates**:
   ```bash
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
     -keyout private.key -out certificate.crt
   ```

2. **Update Docker Configuration**:
   ```yaml
   # docker-compose.yml
   services:
     nginx:
       volumes:
         - ./ssl/private.key:/etc/nginx/ssl/private.key
         - ./ssl/certificate.crt:/etc/nginx/ssl/certificate.crt
   ```

3. **Configure Nginx**:
   ```nginx
   server {
       listen 443 ssl;
       ssl_certificate /etc/nginx/ssl/certificate.crt;
       ssl_certificate_key /etc/nginx/ssl/private.key;
       # ... other configuration
   }
   ```

### 3. Security Configuration

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

## Post-Installation Steps

### 1. Access Services

- **Web Interface**: https://localhost:8443
- **API Documentation**: https://localhost:8443/api/docs
- **Grafana Dashboard**: https://localhost:3000
- **Prometheus**: https://localhost:9090

### 2. Grafana Configuration

1. **Import Dashboards**:
   ```bash
   # Copy dashboard JSON files
   cp grafana/dashboards/*.json /path/to/grafana/provisioning/dashboards/
   ```

2. **Configure Data Sources**:
   - Add Prometheus as a data source
   - Configure alert rules
   - Set up notification channels

### 3. Verification Script

```bash
# Run verification script
python scripts/verify_installation.py

# Expected output:
# ✓ Docker containers running
# ✓ Services accessible
# ✓ Database initialized
# ✓ SSL configured
# ✓ Authentication working
# ✓ Authorization configured
# ✓ Encryption enabled
# ✓ Monitoring setup
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**:
   - Check running services: `netstat -tuln`
   - Update port configurations in `.env`
   - Restart services

2. **Memory Issues**:
   - Increase Docker memory limit
   - Adjust Spark memory settings
   - Monitor resource usage

3. **SSL Certificate Problems**:
   - Verify certificate paths
   - Check file permissions
   - Validate certificate chain

4. **Authentication Issues**:
   - Check LDAP/OAuth configuration
   - Verify user credentials
   - Review authentication logs

5. **Authorization Issues**:
   - Verify role assignments
   - Check permission settings
   - Review access logs

6. **Encryption Issues**:
   - Verify key configuration
   - Check key rotation settings
   - Review encryption logs

## Getting Help

- **Documentation**: https://docs.nessi.dev
- **Issue Tracker**: https://github.com/your-org/nessi-dev/issues
- **Community Forum**: https://community.nessi.dev
- **Support Email**: support@nessi.dev

## Next Steps

- [Quick Start Guide](quickstart.md)
- [Features Overview](features.md)
- [API Documentation](../API.md)
- [Troubleshooting Guide](troubleshooting.md) 