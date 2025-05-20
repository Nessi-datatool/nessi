#!/bin/bash

# Script to fix go.mod file format issues
# Created on 2025-05-20

echo "Fixing go.mod file format..."

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod file not found!"
    exit 1
fi

# Fix the invalid Go version format
sed -i 's/go 1.23.0/go 1.20/' go.mod

# Remove the toolchain directive
sed -i '/toolchain go1.24.3/d' go.mod

echo "go.mod file fixed successfully!"
echo "Current go.mod version:"
grep "^go " go.mod

# Run go mod tidy to ensure dependencies are properly updated
echo "Running go mod tidy..."
go mod tidy

echo "go.mod file is now ready for use!"
