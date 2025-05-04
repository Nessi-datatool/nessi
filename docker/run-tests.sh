#!/bin/bash

# Stop on error
set -e

# Build and run test containers
docker-compose -f docker-compose.test.yml down -v
docker-compose -f docker-compose.test.yml build
docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# Copy test results to host
mkdir -p ../test-results
docker cp $(docker-compose -f docker-compose.test.yml ps -q test-runner):/app/test-results/ ../test-results/

# Clean up
docker-compose -f docker-compose.test.yml down -v 