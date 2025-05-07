#!/bin/bash

# Exit on any error
set -e

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "Installing air for hot reloading..."
    go install github.com/cosmtrek/air@latest
fi

# Create necessary directories if they don't exist
mkdir -p data/delta/tables
mkdir -p data/delta/logs
mkdir -p data/quality/reports
mkdir -p data/reports
mkdir -p logs
mkdir -p certs

# Generate certificates if they don't exist
if [ ! -f "certs/server.crt" ]; then
    echo "Generating self-signed certificates..."
    ./scripts/generate_certs.sh
fi

# Start monitoring stack if docker-compose is available
if command -v docker-compose &> /dev/null; then
    echo "Starting monitoring stack..."
    docker-compose up -d prometheus grafana pushgateway
fi

# Set development environment variables
export CONFIG_PATH=config/config.yaml
export NESSI_ENV=development
export NESSI_LOG_LEVEL=debug

# Run the application with hot reloading
echo "Starting Nessi.dev in development mode..."
air 