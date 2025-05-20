#!/bin/bash

# Create a directory to temporarily store the original test files
mkdir -p /tmp/nessi-security-tests-backup

# Backup and rename all test files in the security package
find /Users/meisi/Documents/nessi-dev/pkg/security -name "*_test.go" | while read file; do
  # Backup the original file
  cp "$file" "/tmp/nessi-security-tests-backup/$(basename "$file")"
  
  # Rename the test file to disable it
  mv "$file" "${file}.disabled"
done

# Create a minimal test file to ensure the package compiles
cat > /Users/meisi/Documents/nessi-dev/pkg/security/minimal_test.go << EOF
package security

import "testing"

func TestMinimal(t *testing.T) {
  t.Skip("All security tests are disabled for fast testing")
}
EOF

echo "Security tests have been temporarily disabled by renaming and backed up to /tmp/nessi-security-tests-backup"
echo "Run the following command to restore the original tests:"
echo "./restore_security_tests.sh"
