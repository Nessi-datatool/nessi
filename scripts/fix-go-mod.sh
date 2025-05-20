#!/bin/bash

# Script to fix go.mod format issues
set -e

echo "Fixing go.mod format issues..."

# Check if go.mod exists
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod file not found"
    exit 1
fi

# Backup the original go.mod
cp go.mod go.mod.bak

# Fix the Go version if needed (replace 1.23.0 with 1.20)
awk '{
    if ($0 ~ /^go 1\.23\.0/) {
        print "go 1.20"
    } else if ($0 ~ /^toolchain/) {
        # Skip toolchain directive
    } else {
        print $0
    }
}' go.mod.bak > go.mod

# Format the go.mod file
go mod edit -fmt

# Run go mod tidy to clean up dependencies
go mod tidy

# Clean up backup
rm go.mod.bak

echo "go.mod has been fixed and formatted."
