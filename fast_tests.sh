#!/bin/bash

# Script to run only fast tests, skipping long-running ones
echo "Running fast tests only..."

# Set environment variables to skip long tests and use short timeouts
export NESSI_SKIP_LONG_TESTS=true
export NESSI_TEST_TIMEOUT_MS=1000

# Create a list of packages with long-running tests
LONG_TEST_PACKAGES=(
  "./pkg/monitoring/dashboard"
  "./pkg/monitoring"
)

# Create a list of packages with tests
TEST_PACKAGES=$(go list ./pkg/... | grep -v vendor)

# Run tests for each package, skipping the ones with long tests
for pkg in $TEST_PACKAGES; do
  # Check if this package should be skipped
  skip=false
  for long_pkg in "${LONG_TEST_PACKAGES[@]}"; do
    if [[ "$pkg" == *"$long_pkg"* ]]; then
      skip=true
      break
    fi
  done
  
  if [ "$skip" = true ]; then
    echo "Skipping package with long tests: $pkg"
  else
    echo "Running tests for package: $pkg"
    go test -v "$pkg" -count=1 -timeout=5s
    if [ $? -ne 0 ]; then
      echo "Tests failed for package: $pkg"
      exit 1
    fi
  fi
done

echo "All fast tests completed successfully!"
