#!/bin/bash

# Script to run all tests with the -short flag
# This script runs tests quickly by skipping long-running tests

echo "Running short tests for Nessi.dev..."
echo "===================================="

# List of key packages to test
PACKAGES=(
  "./pkg/freshness"
  "./pkg/quality/rules"
  "./pkg/lineage/visualization"
  "./pkg/webhook"
  "./pkg/rca"
  "./pkg/security/audit"
  "./pkg/security/rbac"
  "./pkg/monitoring/alerts"
  "./pkg/monitoring/tests"
  "./pkg/quality/anomaly"
  "./pkg/quality/engine"
  "./pkg/quality/profile"
  "./pkg/plugin"
  "./pkg/dbt"
  "./pkg/cloud"
  "./pkg/common"
  "./pkg/config"
)

# Run tests with -short flag
for pkg in "${PACKAGES[@]}"; do
  echo -e "\nRunning tests for package: $pkg"
  go test -short -v $pkg
  
  if [ $? -ne 0 ]; then
    echo "❌ Test failures detected in $pkg"
  else
    echo "✅ All tests in $pkg passed successfully"
  fi
done

echo -e "\nShort test run complete."
echo "For a full test run, use: go test ./..."
