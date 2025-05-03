# Nessi Features

Nessi is a powerful data processing and analysis tool built on Apache Spark, designed for efficient table scanning, data generation, and analysis. It runs entirely in Docker containers for easy deployment and management.

## Key Features

### 🎯 Table Scanning
- **Multi-format Support**: Parquet and CSV
- **Schema Analysis**: Automatic detection and type information
- **Statistical Analysis**: Row counts and column-level statistics
- **Smart Column Handling**: Automatic special character management

### 🚀 Sample Data Generation
- **Flexible Data Types**: Integer, String, Timestamp support
- **Multiple Export Formats**: Parquet and CSV
- **Automatic Timestamps**: Built-in timestamp generation
- **Optimized Storage**: Efficient data export options

### ⚡ Spark Integration
- **Optimized Configuration**: Memory and performance settings
- **Resource Management**: Automatic session cleanup
- **Docker-based Execution**: Containerized Spark environment

### 📊 Grafana Monitoring
- **Real-time Metrics**: Table and column statistics visualization
- **Custom Dashboards**: Create and manage monitoring dashboards
- **Prometheus Integration**: Native metrics collection
- **Alert Configuration**: Set up monitoring alerts

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/nessi-dev/nessi.git
cd nessi
```

2. Start the services:
```bash
docker-compose up -d
```

3. Run the demo:
```bash
docker-compose exec backend python src/demo.py
```

4. Access the dashboard:
- Open http://localhost:3000
- Log in with username: admin, password: admin
- Navigate to "Nessi Demo Dashboard"

## System Requirements
- Docker Engine 20.10.0+
- Docker Compose v2.0.0+
- 4GB RAM minimum
- 10GB disk space minimum

## Docker Services
- **Backend Service**: Data processing and analysis
- **Grafana**: Monitoring and visualization
- **Prometheus**: Metrics collection

## Get Started
1. Follow the [Installation Guide](/docs/installation)
2. Run the [Demo Tutorial](/docs/demo-tutorial)
3. Set up [Monitoring](/docs/monitoring)
4. Explore [Advanced Features](/docs/features)

[View Documentation →](/docs)
[Check API Reference →](/docs/api)
[Learn About Monitoring →](/docs/monitoring) 