#!/bin/bash

# CI Setup Script for Nessi
# This script ensures all dependencies are properly resolved for the CI environment

echo "Setting up CI environment for Nessi..."

# Downgrade the problematic dependency
echo "Downgrading klauspost/compress to a version compatible with Go 1.20..."
go mod edit -replace=github.com/klauspost/compress=github.com/klauspost/compress@v1.16.7

# Clean the module cache
go clean -modcache

# Download all dependencies
echo "Downloading dependencies..."
go mod download

# Run go mod tidy to ensure everything is consistent
echo "Running go mod tidy..."
go mod tidy

echo "CI setup complete!"
