#!/bin/bash

# Restore all security test files from the backup
if [ -d "/tmp/nessi-security-tests-backup" ]; then
  # Remove the minimal test file
  rm -f /Users/meisi/Documents/nessi-dev/pkg/security/minimal_test.go
  
  # Restore all disabled test files
  find /Users/meisi/Documents/nessi-dev/pkg/security -name "*_test.go.disabled" | while read file; do
    # Get the original filename
    original_name="${file%.disabled}"
    # Rename back to original
    mv "$file" "$original_name"
  done
  
  echo "Security tests have been restored"
else
  echo "Error: Backup directory not found at /tmp/nessi-security-tests-backup"
  echo "You may need to manually restore the tests from backup"
  exit 1
fi
