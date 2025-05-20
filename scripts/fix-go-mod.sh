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
sed -i '' -e 's/go 1.23.0/go 1.20/' go.mod

# Verify the change was made
if grep -q 'go 1.23.0' go.mod; then
    echo "Warning: Failed to update Go version with sed, trying direct file manipulation"
    # Create a temporary file with the correct content
    awk '{if($1=="go" && $2=="1.23.0") {print "go 1.20"} else {print $0}}' go.mod > go.mod.tmp
    mv go.mod.tmp go.mod
fi

# Remove the toolchain directive
sed -i '' -e '/toolchain go1.24.3/d' go.mod

# Verify the toolchain directive was removed
if grep -q 'toolchain go1.24.3' go.mod; then
    echo "Warning: Failed to remove toolchain directive with sed, trying direct file manipulation"
    # Create a temporary file without the toolchain directive
    grep -v 'toolchain go1.24.3' go.mod > go.mod.tmp
    mv go.mod.tmp go.mod
fi

# Format the go.mod file
go mod edit -fmt

# Run go mod tidy to ensure dependencies are properly updated
echo "Running go mod tidy..."
go mod tidy

echo "go.mod file fixed successfully!"
echo "Current go.mod version:"
grep "^go " go.mod

echo "go.mod file is now ready for use!"
