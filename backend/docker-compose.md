# Docker Compose Configuration Documentation

## Overview

This document describes the Docker Compose configuration for the nessi.dev monitoring system. The configuration sets up a containerized environment with all necessary services for metrics collection, visualization, and testing.

## Configuration File

The main configuration file is located at `docker-compose.test.yml` and contains the following services:

### Metrics Test Service
```yaml
metrics-test:
  build:
    context: .
    dockerfile: Dockerfile.test
  ports:
    - "8000:8000"
  networks:
    - monitoring
  environment:
    - PYTHONUNBUFFERED=1
  volumes:
    - ./src/reports:/app/src/reports
  command: python src/reports/metrics_test.py
```

- **Build**: Uses custom Dockerfile for testing
- **Ports**: Exposes metrics endpoint on port 8000
- **Networks**: Connects to monitoring network
- **Environment**: Sets Python unbuffered output
- **Volumes**: Mounts reports directory
- **Command**: Runs metrics test script

### Prometheus Service
```yaml
prometheus:
  image: prom/prometheus:latest
  ports:
    - "9090:9090"
  volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml
  networks:
    - monitoring
  command:
    - '--config.file=/etc/prometheus/prometheus.yml'
    - '--storage.tsdb.path=/prometheus'
    - '--web.console.libraries=/usr/share/prometheus/console_libraries'
    - '--web.console.templates=/usr/share/prometheus/consoles'
```

- **Image**: Uses latest Prometheus image
- **Ports**: Exposes Prometheus UI on port 9090
- **Volumes**: Mounts configuration file
- **Networks**: Connects to monitoring network
- **Command**: Configures Prometheus with custom settings

### Grafana Service
```yaml
grafana:
  image: grafana/grafana:latest
  ports:
    - "3000:3000"
  volumes:
    - ./src/reports/grafana_dashboard.json:/etc/grafana/provisioning/dashboards/nessi_dashboard.json
  networks:
    - monitoring
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
    - GF_SECURITY_ADMIN_USER=admin
    - GF_USERS_ALLOW_SIGN_UP=false
    - GF_AUTH_ANONYMOUS_ENABLED=true
    - GF_AUTH_ANONYMOUS_ORG_ROLE=Viewer
  depends_on:
    - prometheus
```

- **Image**: Uses latest Grafana image
- **Ports**: Exposes Grafana UI on port 3000
- **Volumes**: Mounts dashboard configuration
- **Networks**: Connects to monitoring network
- **Environment**: Configures security settings
- **Depends_on**: Ensures Prometheus starts first

### Network Configuration
```yaml
networks:
  monitoring:
    driver: bridge
```

- **Driver**: Uses bridge network driver
- **Purpose**: Enables service communication

## Service Dependencies

### Metrics Flow
1. Metrics Test generates metrics
2. Prometheus scrapes metrics
3. Grafana visualizes metrics

### Startup Order
1. Prometheus starts first
2. Grafana starts after Prometheus
3. Metrics Test can start anytime

## Volume Mounts

### Configuration Files
- `prometheus.yml` → `/etc/prometheus/prometheus.yml`
- `grafana_dashboard.json` → `/etc/grafana/provisioning/dashboards/nessi_dashboard.json`
- `./src/reports` → `/app/src/reports`

## Environment Variables

### Grafana Configuration
- `GF_SECURITY_ADMIN_PASSWORD`: Admin password
- `GF_SECURITY_ADMIN_USER`: Admin username
- `GF_USERS_ALLOW_SIGN_UP`: Disables user signup
- `GF_AUTH_ANONYMOUS_ENABLED`: Enables anonymous access
- `GF_AUTH_ANONYMOUS_ORG_ROLE`: Sets anonymous role

### Metrics Test Configuration
- `PYTHONUNBUFFERED`: Ensures real-time logging

## Port Mapping

### Service Ports
- Metrics Test: 8000
- Prometheus: 9090
- Grafana: 3000

## Best Practices

### Configuration
1. Use version-specific images
2. Configure proper networks
3. Set appropriate volumes
4. Define service dependencies

### Security
1. Use secure passwords
2. Configure proper access
3. Limit exposed ports
4. Set up network isolation

### Performance
1. Optimize resource limits
2. Configure proper volumes
3. Set appropriate timeouts
4. Monitor container health

## Troubleshooting

### Common Issues

1. **Container Startup Failures**
   - Check image availability
   - Verify port availability
   - Check volume permissions
   - Validate network configuration

2. **Service Communication Issues**
   - Verify network connectivity
   - Check service dependencies
   - Validate DNS resolution
   - Test port accessibility

3. **Volume Mount Issues**
   - Check file permissions
   - Verify mount paths
   - Validate file existence
   - Test volume access

### Logs
- View all logs: `docker-compose logs`
- View specific service: `docker-compose logs <service>`
- Follow logs: `docker-compose logs -f`

## Maintenance

### Regular Tasks
1. Update container images
2. Check volume usage
3. Monitor network health
4. Validate configurations

### Backup
1. Export configurations
2. Backup volumes
3. Document changes
4. Version control files

## Security

### Access Control
1. Configure service authentication
2. Set up proper authorization
3. Manage API access
4. Secure endpoints

### Network Security
1. Use internal networks
2. Configure firewalls
3. Monitor access
4. Log security events

## Examples

### Example 1: Basic Service Configuration
```yaml
services:
  metrics-test:
    build:
      context: .
      dockerfile: Dockerfile.test
    ports:
      - "8000:8000"
    networks:
      - monitoring
    environment:
      - PYTHONUNBUFFERED=1
    volumes:
      - ./src/reports:/app/src/reports
    command: python src/reports/metrics_test.py
```

### Example 2: Prometheus Service
```yaml
services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - monitoring
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
```

### Example 3: Grafana Service
```yaml
services:
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - ./src/reports/grafana_dashboard.json:/etc/grafana/provisioning/dashboards/nessi.json
    networks:
      - monitoring
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    depends_on:
      - prometheus
```

### Example 4: Network Configuration
```yaml
networks:
  monitoring:
    driver: bridge
```

### Example 5: Volume Configuration
```yaml
volumes:
  prometheus_data:
    driver: local
  grafana_data:
    driver: local
```

### Example 6: Environment Variables
```yaml
services:
  metrics-test:
    environment:
      - METRICS_INTERVAL=5s
      - SCENARIO=normal_operation
      - LOG_LEVEL=INFO
```

### Example 7: Health Checks
```yaml
services:
  prometheus:
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:9090/-/healthy"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Example 8: Resource Limits
```yaml
services:
  metrics-test:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

### Example 9: Service Dependencies
```yaml
services:
  grafana:
    depends_on:
      prometheus:
        condition: service_healthy
```

### Example 10: Custom Labels
```yaml
services:
  metrics-test:
    labels:
      - "com.nessi.monitoring=true"
      - "com.nessi.environment=development" 