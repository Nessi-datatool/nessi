#!/bin/bash

# Script to debug long-running tests
echo "Debugging long-running tests with verbose output and timeouts..."

# Set environment variables for debugging
export NESSI_TEST_DEBUG=true
export NESSI_TEST_TIMEOUT_MS=10000  # 10 second timeout for each test

# Create log directory if it doesn't exist
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

# Run the specific long-running packages with verbose output and race detection
echo "Running dashboard tests..."
./scripts/run_tests.sh --package="./pkg/monitoring/dashboard/..." --no-cache 2>&1 | tee "$LOG_DIR/dashboard_test_debug.log"

echo "Running monitoring tests..."
./scripts/run_tests.sh --package="./pkg/monitoring/..." --no-cache 2>&1 | tee "$LOG_DIR/monitoring_test_debug.log"

# Run with race detection for specific packages
echo "Running with race detection..."
go test -v -race -timeout=30s ./pkg/monitoring/dashboard/... 2>&1 | tee "$LOG_DIR/dashboard_race_debug.log"
go test -v -race -timeout=30s ./pkg/monitoring/... 2>&1 | tee "$LOG_DIR/monitoring_race_debug.log"

echo "Debug logs saved to $LOG_DIR directory"
