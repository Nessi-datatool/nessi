#!/bin/bash

# Run only fast/unit tests, skipping integration/slow tests
# Integration tests are run in CI only

echo "Running fast/unit tests (integration tests are skipped)..."
go test -short $(go list ./... | grep -v '/scripts')

# This script runs tests in fast mode, skipping long-running tests
# Usage: ./scripts/run_fast_tests.sh [package]
# Example: ./scripts/run_fast_tests.sh ./pkg/security/...

# Set environment variable to enable fast test mode
export NESSI_FAST_TESTS=1

# Default to running all tests if no package is specified
PACKAGE=${1:-"./..."}

# Run tests with verbose output and without caching
go test -v -count=1 $PACKAGE

# Return the exit code from the tests
exit $?
