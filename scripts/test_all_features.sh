#!/bin/bash

# Comprehensive test script for Nessi CLI commands
# This script runs the integration test to verify that all available commands work correctly

# Function to print section headers
print_header() {
  echo "\n=============================================="
  echo "$1"
  echo "=============================================="
}

# Initialize counters
PASSED_TESTS=0
FAILED_TESTS=0
TOTAL_TESTS=0

print_header "NESSI CLI COMMANDS TEST SUITE"
echo "Testing all available commands in the Nessi CLI"

# Run the integration test script
echo -e "\n=============================================="
echo "Running Integration Tests"
echo "=============================================="
echo "Script: ./scripts/test_integration.sh"
echo "-------------------------------------------"
./scripts/test_integration.sh
INTEGRATION_RESULT=$?
if [ $INTEGRATION_RESULT -eq 0 ]; then
  echo -e "\033[32m✅ Test passed\033[0m"
  ((PASSED_TESTS++))
else
  echo -e "\033[31m❌ Test failed\033[0m"
  ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))

# Print test summary
print_header "TEST SUMMARY"
echo "Total tests: $TOTAL_TESTS"
echo -e "Passed tests: \033[32m$PASSED_TESTS\033[0m"
echo -e "Failed tests: \033[31m$FAILED_TESTS\033[0m"

if [ $FAILED_TESTS -eq 0 ]; then
  echo -e "\n\033[32m✅ All tests passed!\033[0m"
  exit 0
else
  echo -e "\n\033[31m❌ Some tests failed!\033[0m"
  exit 1
fi
