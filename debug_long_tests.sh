#!/bin/bash

# Script to debug long-running tests
echo "Debugging long-running tests with verbose output and timeouts..."

# Set environment variables for debugging
export NESSI_TEST_DEBUG=true
export NESSI_TEST_TIMEOUT_MS=10000  # 10 second timeout for each test

# Run the specific long-running packages with verbose output and race detection
go test -v -race -timeout=30s ./pkg/monitoring/dashboard/... 2>&1 | tee dashboard_test_debug.log
go test -v -race -timeout=30s ./pkg/monitoring/... 2>&1 | tee monitoring_test_debug.log

echo "Debug logs saved to dashboard_test_debug.log and monitoring_test_debug.log"
