# Configuration Guide

## Overview

This guide covers all configuration options available in Nessi, including security settings, monitoring configuration, and service-specific settings.

## Configuration Files

### 1. Security Configuration (`config/security.json`)

```json
{
  "jwt": {
    "secret": "your-secret-key",
    "algorithm": "HS256",
    "expiry_hours": 24
  },
  "rate_limit": {
    "enabled": true,
    "requests": 100,
    "window": 60,
    "block_duration": 300
  },
  "ip_whitelist": {
    "enabled": true,
    "networks": ["192.168.1.0/24"]
  },
  "roles": {
    "admin": ["*"],
    "analyst": ["read", "analyze"],
    "viewer": ["read"]
  },
  "ssl": {
    "enabled": true,
    "cert_path": "/app/ssl/cert.pem",
    "key_path": "/app/ssl/key.pem"
  }
}
```

### 2. Monitoring Configuration (`config/monitoring.json`)

```json
{
  "metrics": {
    "port": 9090,
    "path": "/metrics",
    "interval": "15s"
  },
  "grafana": {
    "host": "localhost",
    "port": 3000,
    "admin_user": "admin",
    "admin_password": "your-password"
  },
  "prometheus": {
    "scrape_interval": "15s",
    "evaluation_interval": "15s",
    "retention_days": 15
  },
  "alertmanager": {
    "host": "localhost",
    "port": 9093,
    "config_path": "/etc/alertmanager/alertmanager.yml"
  },
  "alerts": {
    "log_path": "logs/alerts.log",
    "email": {
      "enabled": true,
      "smtp_host": "smtp.example.com",
      "smtp_port": 587,
      "from": "alerts@nessi.dev",
      "to": ["admin@example.com"]
    },
    "slack": {
      "enabled": true,
      "webhook_url": "https://hooks.slack.com/services/..."
    }
  }
}
```

### 3. Delta Lake Configuration (`config/delta.json`)

```json
{
  "optimization": {
    "auto_optimize": true,
    "optimize_interval": "1d",
    "z_order_columns": ["timestamp", "user_id"],
    "vacuum_retention_hours": 168
  },
  "schema_evolution": {
    "auto_merge": true,
    "strict_mode": false
  },
  "time_travel": {
    "retention_days": 30,
    "log_retention_days": 90
  }
}
```

### 4. Quality Analysis Configuration (`config/quality.json`)

```json
{
  "profiling": {
    "sample_size": 10000,
    "numeric_precision": 2,
    "datetime_format": "%Y-%m-%d %H:%M:%S"
  },
  "anomaly_detection": {
    "z_score_threshold": 3.0,
    "iqr_multiplier": 1.5,
    "min_samples": 100
  },
  "pattern_recognition": {
    "email": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
    "phone": "^\\+?[\\d\\s-]{10,}$",
    "date": "^\\d{4}-\\d{2}-\\d{2}$",
    "time": "^\\d{2}:\\d{2}:\\d{2}$",
    "url": "^https?://(?:www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b(?:[-a-zA-Z0-9()@:%_\\+.~#?&//=]*)$"
  },
  "validation": {
    "rules_path": "config/validation_rules.json",
    "strict_mode": false
  }
}
```

## Environment Variables

### 1. Core Configuration

```bash
# Application
NESSI_ENV=production
NESSI_LOG_LEVEL=INFO
NESSI_LOG_FORMAT=json

# Database
NESSI_DB_HOST=localhost
NESSI_DB_PORT=5432
NESSI_DB_NAME=nessi
NESSI_DB_USER=nessi
NESSI_DB_PASSWORD=your-password

# Redis
NESSI_REDIS_HOST=localhost
NESSI_REDIS_PORT=6379
NESSI_REDIS_DB=0
```

### 2. Security Configuration

```bash
# JWT
NESSI_JWT_SECRET=your-secret-key
NESSI_JWT_ALGORITHM=HS256
NESSI_JWT_EXPIRY_HOURS=24

# SSL
NESSI_SSL_ENABLED=true
NESSI_SSL_CERT_PATH=/app/ssl/cert.pem
NESSI_SSL_KEY_PATH=/app/ssl/key.pem
```

### 3. Monitoring Configuration

```bash
# Prometheus
NESSI_PROMETHEUS_PORT=9090
NESSI_PROMETHEUS_PATH=/metrics
NESSI_PROMETHEUS_INTERVAL=15s

# Grafana
NESSI_GRAFANA_HOST=localhost
NESSI_GRAFANA_PORT=3000
NESSI_GRAFANA_ADMIN_USER=admin
NESSI_GRAFANA_ADMIN_PASSWORD=your-password
```

## Service Configuration

### 1. Nginx Configuration (`nginx/conf.d/nessi.conf`)

```nginx
server {
    listen 443 ssl;
    server_name nessi.dev;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://nessi:8443;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /metrics {
        proxy_pass http://prometheus:9090;
        auth_basic "Prometheus";
        auth_basic_user_file /etc/nginx/.htpasswd;
    }
}
```

### 2. Prometheus Configuration (`prometheus/prometheus.yml`)

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: 'nessi'
    static_configs:
      - targets: ['nessi:8443']
    scheme: https
    tls_config:
      insecure_skip_verify: true
```

### 3. Alertmanager Configuration (`alertmanager/alertmanager.yml`)

```yaml
global:
  resolve_timeout: 5m
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alerts@nessi.dev'
  smtp_auth_username: 'user'
  smtp_auth_password: 'password'

route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'email'

receivers:
  - name: 'email'
    email_configs:
      - to: 'admin@example.com'
```

## Configuration Management

### 1. Configuration Updates

```bash
# Reload configuration
nessi config reload

# Validate configuration
nessi config validate

# Show current configuration
nessi config show
```

### 2. Configuration Backup

```bash
# Backup configuration
nessi config backup --output config_backup.tar.gz

# Restore configuration
nessi config restore --input config_backup.tar.gz
```

## Best Practices

### 1. Security Configuration

- Use strong, unique secrets for JWT
- Enable SSL/TLS for all services
- Implement proper role hierarchy
- Regular security audits

### 2. Monitoring Configuration

- Set appropriate scrape intervals
- Configure meaningful alert thresholds
- Use multiple notification channels
- Regular log rotation

### 3. Performance Configuration

- Optimize Delta Lake settings
- Configure appropriate retention periods
- Set up proper partitioning
- Regular table optimization

## Next Steps

- [Security Best Practices](best_practices/security.md)
- [Monitoring Setup](admin/monitoring.md)
- [Troubleshooting Guide](troubleshooting/common_issues.md) 