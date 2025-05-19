#!/bin/bash
set -e

# Ensure the script runs from the repository root
cd "$(dirname "$0")/.."

# Find all Go files in the repository
GO_FILES=$(find . -type f -name "*.go" | grep -v "/vendor/" | grep -v "/node_modules/")

if [ -z "$GO_FILES" ]; then
  echo "Error: No Go files found to analyze"
  exit 1
fi

# Run gofmt to fix formatting issues
echo "Running gofmt..."
gofmt -w -s $GO_FILES

# Run go vet for basic checks
echo "Running go vet..."
go vet ./...

# Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
  echo "Running golangci-lint..."
  golangci-lint run --timeout=5m
else
  echo "golangci-lint not found, skipping"
fi

echo "Linting completed successfully!"
