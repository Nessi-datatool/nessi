#!/bin/bash

# Test script for Nessi usability features
# This script tests all the new usability features to ensure they work correctly

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

print_header "NESSI USABILITY FEATURES TEST"

# Test 1: Input Validation - Path Validation
print_header "1. INPUT VALIDATION - PATH VALIDATION"

print_test "Invalid path detection" "nessi schema show --path $TEST_DIR/nonexistent/path"
((TOTAL_TESTS++))
nessi schema show --path $TEST_DIR/nonexistent/path
check_result $?

print_test "Path normalization" "nessi schema show --path ./relative/path"
((TOTAL_TESTS++))
nessi schema show --path ./relative/path
check_result $?

# Test 2: Input Validation - Configuration Validation
print_header "2. INPUT VALIDATION - CONFIGURATION VALIDATION"

print_test "Invalid configuration detection" "nessi config validate --file $TEST_DIR/invalid_config.yaml"
((TOTAL_TESTS++))

# Create an invalid config file
cat > "$TEST_DIR/invalid_config.yaml" << EOF
error_handling:
  interactive_resolution: not_a_boolean
EOF

nessi config validate --file $TEST_DIR/invalid_config.yaml
check_result $?

print_test "Valid configuration validation" "nessi config validate --file $TEST_DIR/valid_config.yaml"
((TOTAL_TESTS++))

# Create a valid config file
cat > "$TEST_DIR/valid_config.yaml" << EOF
error_handling:
  interactive_resolution: true
EOF

nessi config validate --file $TEST_DIR/valid_config.yaml
check_result $?

# Test 3: Guided Setup Wizard
print_header "3. GUIDED SETUP WIZARD"

print_test "Setup wizard execution" "nessi setup --non-interactive --config-dir $TEST_DIR/nessi_config"
((TOTAL_TESTS++))

# Run setup wizard in non-interactive mode
nessi setup --non-interactive --config-dir $TEST_DIR/nessi_config
check_result $?

# Check if config file was created
if [ -f "$TEST_DIR/nessi_config/config.yaml" ]; then
  echo -e "\033[32m✅ Config file created successfully\033[0m"
else
  echo -e "\033[31m❌ Config file not created\033[0m"
fi

# Test 4: Dry Run Mode
print_header "4. DRY RUN MODE"

print_test "Dry run mode" "nessi schema create --path $TEST_DIR/new_schema --dry-run"
((TOTAL_TESTS++))

# Run a command in dry run mode
nessi schema create --path $TEST_DIR/new_schema --dry-run
check_result $?

# Check that the directory was not actually created
if [ ! -d "$TEST_DIR/new_schema" ]; then
  echo -e "\033[32m✅ Directory not created (as expected in dry run)\033[0m"
else
  echo -e "\033[31m❌ Directory was created (should not happen in dry run)\033[0m"
fi

# Test 5: Error Handler Integration
print_header "5. ERROR HANDLER INTEGRATION"

print_test "Error handler with validation" "nessi test-error N101 --path $TEST_DIR/nonexistent/path"
((TOTAL_TESTS++))

# Test error handler with validation
nessi test-error N101 --path $TEST_DIR/nonexistent/path
check_result $?

print_test "Error handler with suggestions" "nessi test-error N201 --details 'Invalid configuration'"
((TOTAL_TESTS++))

# Test error handler with suggestions
nessi test-error N201 --details 'Invalid configuration'
check_result $?

# Test 6: Interactive Error Resolution
print_header "6. INTERACTIVE ERROR RESOLUTION"

print_test "Interactive error resolution" "nessi test-error N101 --resolvable --interactive --path $TEST_DIR/resolvable"
((TOTAL_TESTS++))

# Test interactive error resolution
echo "y" | nessi test-error N101 --resolvable --interactive --path $TEST_DIR/resolvable
check_result $?

# Check if the directory was created
if [ -d "$TEST_DIR/resolvable" ]; then
  echo -e "\033[32m✅ Directory created through interactive resolution\033[0m"
else
  echo -e "\033[31m❌ Directory not created\033[0m"
fi

# Test 7: Configuration System
print_header "7. CONFIGURATION SYSTEM"

print_test "Configuration override with flags" "nessi schema show --interactive=false --path $TEST_DIR"
((TOTAL_TESTS++))

# Test configuration override with flags
nessi schema show --interactive=false --path $TEST_DIR
check_result $?

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
