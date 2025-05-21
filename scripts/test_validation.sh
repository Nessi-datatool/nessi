#!/bin/bash

# Test script for Nessi validation features
# This script tests the validation utilities to ensure they work correctly

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
    echo -e "\033[32mu2705 Test passed\033[0m"
    ((PASSED_TESTS++))
  else
    echo -e "\033[31mu274c Test failed\033[0m"
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

print_header "NESSI VALIDATION FEATURES TEST"

# Test 1: Path Validation
print_header "1. PATH VALIDATION"

# Create test files and directories
TEST_SUBDIR="$TEST_DIR/subdir"
mkdir -p "$TEST_SUBDIR"
TEST_FILE="$TEST_DIR/test_file.txt"
touch "$TEST_FILE"

print_test "Valid path detection" "nessi schema show --path $TEST_DIR"
((TOTAL_TESTS++))
nessi schema show --path $TEST_DIR
check_result $?

print_test "Invalid path detection" "nessi schema show --path $TEST_DIR/nonexistent"
((TOTAL_TESTS++))
nessi schema show --path $TEST_DIR/nonexistent
check_result $?

print_test "Path type validation - directory" "nessi schema show --path $TEST_FILE"
((TOTAL_TESTS++))
nessi schema show --path $TEST_FILE
check_result $?

print_test "Path type validation - file" "nessi config validate --file $TEST_DIR"
((TOTAL_TESTS++))
nessi config validate --file $TEST_DIR
check_result $?

print_test "Path normalization - relative path" "nessi schema show --path ./relative/path"
((TOTAL_TESTS++))
nessi schema show --path ./relative/path
check_result $?

# Test 2: Configuration Validation
print_header "2. CONFIGURATION VALIDATION"

print_test "Valid configuration" "nessi config validate --file $TEST_DIR/valid_config.yaml"
((TOTAL_TESTS++))

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

nessi config validate --file $TEST_DIR/valid_config.yaml
check_result $?

print_test "Invalid configuration - type error" "nessi config validate --file $TEST_DIR/invalid_config_type.yaml"
((TOTAL_TESTS++))

# Create an invalid config file with type error
cat > "$TEST_DIR/invalid_config_type.yaml" << EOF
error_handling:
  interactive_resolution: "not_a_boolean"
EOF

nessi config validate --file $TEST_DIR/invalid_config_type.yaml
check_result $?

print_test "Invalid configuration - missing required field" "nessi config validate --file $TEST_DIR/invalid_config_missing.yaml"
((TOTAL_TESTS++))

# Create an invalid config file with missing required field
cat > "$TEST_DIR/invalid_config_missing.yaml" << EOF
error_handling:
  telemetry:
    max_errors: 100
EOF

nessi config validate --file $TEST_DIR/invalid_config_missing.yaml
check_result $?

# Test 3: Output Format Validation
print_header "3. OUTPUT FORMAT VALIDATION"

print_test "Valid output format - html" "nessi report generate --table $TEST_DIR --format html"
((TOTAL_TESTS++))
nessi report generate --table $TEST_DIR --format html
check_result $?

print_test "Valid output format - json" "nessi report generate --table $TEST_DIR --format json"
((TOTAL_TESTS++))
nessi report generate --table $TEST_DIR --format json
check_result $?

print_test "Invalid output format" "nessi report generate --table $TEST_DIR --format invalid"
((TOTAL_TESTS++))
nessi report generate --table $TEST_DIR --format invalid
check_result $?

# Test Summary
print_header "TEST SUMMARY"

echo "Total tests: $TOTAL_TESTS"
echo -e "Passed tests: \033[32m$PASSED_TESTS\033[0m"
echo -e "Failed tests: \033[31m$FAILED_TESTS\033[0m"

if [ $FAILED_TESTS -eq 0 ]; then
  echo -e "\n\033[32mu2705 All tests passed!\033[0m"
  exit 0
else
  echo -e "\n\033[31mu274c Some tests failed!\033[0m"
  exit 1
fi
