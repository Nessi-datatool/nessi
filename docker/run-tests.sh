#!/bin/bash

# Stop on error
set -e

# Build and run test containers
docker-compose -f docker-compose.test.yml down -v
docker-compose -f docker-compose.test.yml build
docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# Copy test results from volume to host
mkdir -p ../test-results
docker run --rm -v nessi_test-results:/app/test-results -v $(pwd)/../test-results:/host alpine sh -c "cp -r /app/test-results/* /host/"

# Clean up
docker-compose -f docker-compose.test.yml down -v 