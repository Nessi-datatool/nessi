#!/bin/bash

# Script to generate test coverage reports for the Nessi project
# This helps identify areas with insufficient test coverage

# Create output directory
mkdir -p ./coverage

# Run tests with coverage
echo "Running tests with coverage..."
go test ./... -coverprofile=./coverage/coverage.out

# Generate HTML report
echo "Generating HTML coverage report..."
go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

# Generate function-level coverage report
echo "Generating function-level coverage report..."
go tool cover -func=./coverage/coverage.out > ./coverage/coverage_func.txt

# Find packages with low coverage
echo "Packages with low test coverage:"
go tool cover -func=./coverage/coverage.out | grep -v "100.0%" | sort -k 3 -n | head -n 20 > ./coverage/low_coverage.txt
cat ./coverage/low_coverage.txt

echo "Coverage reports generated in ./coverage directory"
echo "Open ./coverage/coverage.html in a browser to view detailed coverage information"
