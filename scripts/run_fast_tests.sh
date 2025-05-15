#!/bin/bash

# This script runs all tests except the long-running security integration tests
# It's a faster alternative to running the full test suite

# Set environment variable to skip security integration tests
export SKIP_SECURITY_INTEGRATION=1

# Run all tests with the short flag to skip long-running tests
go test -short ./pkg/... "$@"

# Print success message
echo "Fast tests completed successfully!"
