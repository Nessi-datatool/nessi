#!/bin/bash

# Check if running in Docker
if [ ! -f /.dockerenv ]; then
    echo "Error: This application must be run in a Docker container."
    echo "Please use docker-compose to run the application:"
    echo "docker-compose up -d"
    exit 1
fi

# Check for required environment variables
required_vars=(
    "PYTHONPATH"
    "JAVA_HOME"
    "PYSPARK_PYTHON"
    "PYSPARK_DRIVER_PYTHON"
)

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set."
        echo "This suggests the container is not properly configured."
        exit 1
    fi
done

# Check if running as non-root user
if [ "$(id -u)" = "0" ]; then
    echo "Error: This application should not be run as root."
    echo "Please ensure the container is configured to run as a non-root user."
    exit 1
fi

# If all checks pass, continue with the original command
exec "$@" 