#!/bin/bash

# Test script for Nessi dry run mode
# This script tests the dry run functionality to ensure it works correctly

# Function to print section headers
print_header() {
  echo "\n=============================================="
  echo "$1"
  echo "=============================================="
}

# Function to print test case
print_test() {
  echo "\n>> Testing: $1"
  echo "Command: $2"
  echo "-------------------------------------------"
}

# Function to check if a test passed or failed
check_result() {
  if [ $1 -eq 0 ]; then
    echo -e "\033[32m✅ Test passed\033[0m"
    ((PASSED_TESTS++))
  else
    echo -e "\033[31m❌ Test failed\033[0m"
    ((FAILED_TESTS++))
  fi
}

# Initialize counters
PASSED_TESTS=0
FAILED_TESTS=0
TOTAL_TESTS=0

# Create a temporary directory for test files
TEST_DIR="$(mktemp -d)"
echo "Created temporary test directory: $TEST_DIR"

# Register cleanup function to run on script exit
cleanup() {
  echo "\nCleaning up test directory..."
  rm -rf "$TEST_DIR"
  echo "Done."
}
trap cleanup EXIT

print_header "NESSI DRY RUN MODE TEST"

# Test 1: Schema Create with Dry Run
print_header "1. SCHEMA CREATE WITH DRY RUN"

print_test "Dry run mode - schema create" "nessi schema create --path $TEST_DIR/new_schema --dry-run"
((TOTAL_TESTS++))

# Run schema create with dry run
nessi schema create --path $TEST_DIR/new_schema --dry-run
check_result $?

# Check that the directory was not actually created
if [ ! -d "$TEST_DIR/new_schema" ]; then
  echo -e "\033[32m✅ Directory not created (as expected in dry run)\033[0m"
else
  echo -e "\033[31m❌ Directory was created (should not happen in dry run)\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

# Test 2: Schema Create without Dry Run
print_header "2. SCHEMA CREATE WITHOUT DRY RUN"

print_test "Normal mode - schema create" "nessi schema create --path $TEST_DIR/new_schema"
((TOTAL_TESTS++))

# Run schema create without dry run
nessi schema create --path $TEST_DIR/new_schema
check_result $?

# Check that the directory was actually created
if [ -d "$TEST_DIR/new_schema" ]; then
  echo -e "\033[32m✅ Directory created (as expected in normal mode)\033[0m"
else
  echo -e "\033[31m❌ Directory not created (should happen in normal mode)\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

# Test 3: Configuration Validation with Dry Run
print_header "3. CONFIGURATION VALIDATION WITH DRY RUN"

# Create a valid config file
cat > "$TEST_DIR/valid_config.yaml" << EOF
error_handling:
  interactive_resolution: true
  telemetry:
    enabled: true
    max_errors: 100
  retry:
    enabled: true
    max_attempts: 3
EOF

print_test "Dry run mode - config validate" "nessi config validate --file $TEST_DIR/valid_config.yaml --dry-run"
((TOTAL_TESTS++))

# Run config validate with dry run
nessi config validate --file $TEST_DIR/valid_config.yaml --dry-run
check_result $?

# Test 4: Dry Run with Error
print_header "4. DRY RUN WITH ERROR"

print_test "Dry run mode with error" "nessi schema create --path /nonexistent/path --dry-run"
((TOTAL_TESTS++))

# Run schema create with dry run on a path that would cause an error
nessi schema create --path /nonexistent/path --dry-run
check_result $?

# Test 5: Dry Run with Multiple Actions
print_header "5. DRY RUN WITH MULTIPLE ACTIONS"

print_test "Dry run mode with multiple actions" "nessi schema create --path $TEST_DIR/multi_action --dry-run"
((TOTAL_TESTS++))

# Run schema create with dry run
nessi schema create --path $TEST_DIR/multi_action --dry-run
check_result $?

# Check that the directory was not actually created
if [ ! -d "$TEST_DIR/multi_action" ]; then
  echo -e "\033[32m✅ Directory not created (as expected in dry run)\033[0m"
else
  echo -e "\033[31m❌ Directory was created (should not happen in dry run)\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

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
