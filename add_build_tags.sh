#!/bin/bash

# Add build tags to all integration test files
find /Users/meisi/Documents/nessi-dev/pkg/security -name "*_integration_test.go" | while read file; do
  # Check if the file already has a build tag
  if ! grep -q "//go:build" "$file"; then
    # Add the build tag at the top of the file
    sed -i '' '1s/^/\/\/go:build !fast_tests\n\npackage/' "$file"
    # Fix the package line (remove the added "package" from the sed command)
    sed -i '' 's/package security_test/package security_test/' "$file"
  fi
done

echo "Build tags added to all integration test files"
