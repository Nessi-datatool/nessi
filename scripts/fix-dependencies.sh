#!/bin/bash

# Script to fix dependency issues with github.com/klauspost/compress/zstd
# Created on 2025-05-19

set -e

echo "Fixing dependency issues for Databricks integration..."

# Explicitly get the required dependencies
echo "Getting required dependencies..."
go get github.com/apache/arrow/go/v15/arrow/ipc@v15.0.2
go get github.com/klauspost/compress/zstd@latest

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
go mod vendor

echo "Dependency fix completed successfully!"

# Run tests to verify the fix worked
echo "Running tests to verify the fix..."
go test -v ./pkg/catalog/databricks/... -short
go test -v ./pkg/datalake/... -short

echo "All tests passed successfully!"
