#!/bin/bash

# Script to fix dependency issues with dependencies requiring newer Go versions
# Created on 2025-05-19, updated for slices package compatibility

set -e

echo "Fixing dependency issues for Databricks integration..."

# Explicitly downgrade dependencies that require newer Go versions
echo "Downgrading dependencies that require Go 1.23+..."
go get github.com/apache/arrow/go/v15/arrow/ipc@v15.0.2
go get github.com/apache/arrow/go/v15/arrow/memory@v15.0.2
go get github.com/klauspost/compress/zstd@v1.16.7
go get golang.org/x/sys/cpu@v0.15.0
go get golang.org/x/sys/unix@v0.15.0
go get github.com/sirupsen/logrus@v1.9.3

# Pin specific versions of problematic dependencies
echo "Pinning specific versions of dependencies..."
go mod edit -replace=golang.org/x/sys=golang.org/x/sys@v0.15.0
go mod edit -replace=github.com/klauspost/compress=github.com/klauspost/compress@v1.16.7

# Clean the module cache
echo "Cleaning module cache..."
go clean -modcache

# Download all dependencies
echo "Downloading dependencies..."
go mod download

# Run go mod tidy to ensure everything is consistent
echo "Running go mod tidy..."
go mod tidy

# Create vendor directory to ensure all dependencies are included
echo "Creating vendor directory..."
rm -rf vendor
go mod vendor

echo "Dependency fix completed successfully!"

# Run tests to verify the fix worked
echo "Running tests to verify the fix..."
go test -v ./pkg/catalog/databricks/... -short
go test -v ./pkg/datalake/... -short

echo "All tests passed successfully!"
