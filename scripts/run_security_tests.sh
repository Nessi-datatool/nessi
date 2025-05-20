#!/bin/bash

# Script to run all security tests with optimized test helpers
# Created as part of the security test optimization project

set -e

echo "Running security tests with optimized test helpers..."

# Run the security package tests
echo "Testing pkg/security..."
go test -v ./pkg/security/...

# Run the webhook package tests (which use security components)
echo "Testing pkg/webhook..."
go test -v ./pkg/webhook/...

# Run the monitoring package tests (which were also optimized)
echo "Testing pkg/monitoring..."
go test -v ./pkg/monitoring/...

echo "All tests completed successfully!"
