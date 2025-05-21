#!/bin/bash

# Test script for Nessi setup wizard
# This script tests the setup wizard functionality to ensure it works correctly

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

print_header "NESSI SETUP WIZARD TEST"

# Test 1: Non-Interactive Setup
print_header "1. NON-INTERACTIVE SETUP"

print_test "Non-interactive setup" "nessi setup --non-interactive --config-dir $TEST_DIR/nessi_config"
((TOTAL_TESTS++))

# Run setup in non-interactive mode
nessi setup --non-interactive --config-dir $TEST_DIR/nessi_config
check_result $?

# Check that the config directory was created
if [ -d "$TEST_DIR/nessi_config" ]; then
  echo -e "\033[32mu2705 Config directory created\033[0m"
else
  echo -e "\033[31mu274c Config directory not created\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

# Check that the config file was created
if [ -f "$TEST_DIR/nessi_config/config.yaml" ]; then
  echo -e "\033[32mu2705 Config file created\033[0m"
else
  echo -e "\033[31mu274c Config file not created\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

# Test 2: Config File Content
print_header "2. CONFIG FILE CONTENT"

print_test "Config file content validation" "cat $TEST_DIR/nessi_config/config.yaml"
((TOTAL_TESTS++))

# Check if the config file contains the expected content
if grep -q "interactive_resolution: true" "$TEST_DIR/nessi_config/config.yaml" && \
   grep -q "telemetry:" "$TEST_DIR/nessi_config/config.yaml" && \
   grep -q "retry:" "$TEST_DIR/nessi_config/config.yaml" && \
   grep -q "logging:" "$TEST_DIR/nessi_config/config.yaml"; then
  echo -e "\033[32mu2705 Config file contains expected content\033[0m"
else
  echo -e "\033[31mu274c Config file does not contain expected content\033[0m"
  ((FAILED_TESTS++))
fi

# Display config file content
echo "\nConfig file content:"
cat "$TEST_DIR/nessi_config/config.yaml"

# Test 3: Config Validation
print_header "3. CONFIG VALIDATION"

print_test "Config validation" "nessi config validate --file $TEST_DIR/nessi_config/config.yaml"
((TOTAL_TESTS++))

# Validate the generated config file
nessi config validate --file $TEST_DIR/nessi_config/config.yaml
check_result $?

# Test 4: Setup with Dry Run
print_header "4. SETUP WITH DRY RUN"

print_test "Setup with dry run" "nessi setup --non-interactive --config-dir $TEST_DIR/dry_run_config --dry-run"
((TOTAL_TESTS++))

# Run setup with dry run
nessi setup --non-interactive --config-dir $TEST_DIR/dry_run_config --dry-run
check_result $?

# Check that the config directory was not created (dry run)
if [ ! -d "$TEST_DIR/dry_run_config" ]; then
  echo -e "\033[32mu2705 Config directory not created (as expected in dry run)\033[0m"
else
  echo -e "\033[31mu274c Config directory was created (should not happen in dry run)\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

# Test 5: Setup with Custom Config
print_header "5. SETUP WITH CUSTOM CONFIG"

# Create a custom config directory
CUSTOM_CONFIG_DIR="$TEST_DIR/custom_config"
mkdir -p "$CUSTOM_CONFIG_DIR"

print_test "Setup with existing directory" "nessi setup --non-interactive --config-dir $CUSTOM_CONFIG_DIR"
((TOTAL_TESTS++))

# Run setup with existing directory
nessi setup --non-interactive --config-dir $CUSTOM_CONFIG_DIR
check_result $?

# Check that the config file was created in the custom directory
if [ -f "$CUSTOM_CONFIG_DIR/config.yaml" ]; then
  echo -e "\033[32mu2705 Config file created in custom directory\033[0m"
else
  echo -e "\033[31mu274c Config file not created in custom directory\033[0m"
  ((FAILED_TESTS++))
  ((TOTAL_TESTS++))
fi

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
