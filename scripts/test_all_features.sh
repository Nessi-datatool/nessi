#!/bin/bash

# Comprehensive test script for Nessi error handling and usability features
# This script runs all the test scripts to verify that all features work correctly

# Function to print section headers
print_header() {
  echo "\n=============================================="
  echo "$1"
  echo "=============================================="
}

# Function to run a test script and check the result
run_test_script() {
  local script=$1
  local description=$2
  
  print_header "Running $description"
  echo "Script: $script"
  echo "-------------------------------------------"
  
  # Run the script
  bash "$script"
  local result=$?
  
  # Check the result
  if [ $result -eq 0 ]; then
    echo -e "\033[32m✅ Test passed\033[0m"
    return 0
  else
    echo -e "\033[31m❌ Test failed\033[0m"
    return 1
  fi
}

# Initialize counters
PASSED_TESTS=0
FAILED_TESTS=0
TOTAL_TESTS=0

print_header "NESSI COMPREHENSIVE TEST SUITE"
echo "Testing all error handling and usability features"

# Test 1: Error Handling
run_test_script "./scripts/test_error_handling.sh" "Error Handling Tests"
if [ $? -eq 0 ]; then
  ((PASSED_TESTS++))
else
  ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))

# Test 2: Usability Features
run_test_script "./scripts/test_usability_features.sh" "Usability Features Tests"
if [ $? -eq 0 ]; then
  ((PASSED_TESTS++))
else
  ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))

# Test 3: Validation
run_test_script "./scripts/test_validation.sh" "Validation Tests"
if [ $? -eq 0 ]; then
  ((PASSED_TESTS++))
else
  ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))

# Test 4: Dry Run Mode
run_test_script "./scripts/test_dry_run.sh" "Dry Run Mode Tests"
if [ $? -eq 0 ]; then
  ((PASSED_TESTS++))
else
  ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))

# Test 5: Setup Wizard
run_test_script "./scripts/test_setup_wizard.sh" "Setup Wizard Tests"
if [ $? -eq 0 ]; then
  ((PASSED_TESTS++))
else
  ((FAILED_TESTS++))
fi
((TOTAL_TESTS++))

# Test Summary
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
