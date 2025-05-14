#!/bin/bash

# Set environment variables to minimize test time
export NESSI_TEST_TIMEOUT_MS=1000
export NESSI_SKIP_LONG_TESTS=true
export GO_TEST=1

# Run only the specific short test files we created
echo "Running short tests only..."

# Run dashboard short tests
go test -v ./pkg/monitoring/dashboard/alerts_short_test.go

# Run monitoring short tests
go test -v ./pkg/monitoring/export_short_test.go

echo "Short tests completed!"
