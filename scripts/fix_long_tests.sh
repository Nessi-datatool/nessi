#!/bin/bash

# Script to identify and fix long-running tests
# This script runs tests with a timeout and identifies which ones are taking too long

echo "Identifying long-running tests in Nessi.dev..."
echo "=============================================="

# List of packages to check (add more as needed)
PACKAGES=(
  "./pkg/webhook"
  "./pkg/rca"
  "./pkg/monitoring/alerts"
  "./pkg/quality/anomaly"
)

# Run tests with timeout
for pkg in "${PACKAGES[@]}"; do
  echo -e "\nChecking package: $pkg"
  timeout 5s go test -v $pkg
  
  if [ $? -eq 124 ]; then
    echo "⚠️  Test timeout detected in $pkg"
    echo "Running with -short flag to identify specific long-running tests:"
    go test -v -short $pkg
  elif [ $? -ne 0 ]; then
    echo "❌ Test failures detected in $pkg"
  else
    echo "✅ All tests in $pkg completed successfully within timeout"
  fi
done

echo -e "\nTest analysis complete. Review the output above to identify long-running tests."
echo "For any packages with timeouts, examine the test files and add appropriate skips for long-running tests when the -short flag is used."
echo "Example: if testing.Short() { t.Skip(\"Skipping long-running test in short mode\") }"
