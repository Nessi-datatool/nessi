#!/bin/bash

# Run tests for core packages
echo "Running tests for core packages..."
go test ./pkg/... ./internal/... ./tests/... -v

# Run tests for individual command files with specific build tags
echo "\nRunning tests for individual command files..."
cd cmd/nessi

# Run schema tests
echo "\nRunning schema tests..."
go test -tags=schema_test -run "TestSchema" -v

# Run alerts tests
echo "\nRunning alerts tests..."
go test -tags=alerts_test -run "TestAlerts" -v

# Run time travel tests
echo "\nRunning time travel tests..."
go test -tags=time_travel_test -run "TestTimeTravel" -v
