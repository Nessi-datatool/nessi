#!/bin/bash

# Example script demonstrating the error handling capabilities of Nessi
# This script shows how to use different error handling features

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

print_header "NESSI ERROR HANDLING EXAMPLES"

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

# Example 1: Basic Error Handling
print_header "BASIC ERROR HANDLING"

print_test "Standard error handling" "nessi schema show --path $TEST_DIR/nonexistent/path"
echo "This will show a standard error message with error code N101"
nessi schema show --path $TEST_DIR/nonexistent/path

# Example 2: Interactive Error Resolution
print_header "INTERACTIVE ERROR RESOLUTION"

print_test "Interactive error resolution" "nessi schema show --path $TEST_DIR/resolvable --interactive"
echo "This will offer to create the directory interactively"
nessi schema show --path $TEST_DIR/resolvable --interactive

# Example 3: Error Suggestions
print_header "ERROR SUGGESTIONS"

print_test "Error suggestions" "nessi test-error N201"
echo "This will show suggestions for fixing a configuration error"
nessi test-error N201

# Example 4: Error Telemetry
print_header "ERROR TELEMETRY"

print_test "Error telemetry recording" "nessi test-error N301 --telemetry"
echo "This will record the error in telemetry"
nessi test-error N301 --telemetry

print_test "Error telemetry status" "nessi telemetry status"
echo "This will show the current telemetry status"
nessi telemetry status

print_test "Error telemetry export" "nessi telemetry export-errors $TEST_DIR/error_report.json"
echo "This will export error telemetry to a file"
nessi telemetry export-errors $TEST_DIR/error_report.json

# Example 5: Testing Different Error Types
print_header "TESTING DIFFERENT ERROR TYPES"

print_test "Path error" "nessi test-error N101"
nessi test-error N101

print_test "Configuration error" "nessi test-error N201"
nessi test-error N201

print_test "Authentication error" "nessi test-error N301"
nessi test-error N301

print_test "Connection error" "nessi test-error N401"
nessi test-error N401

print_test "Delta Lake error" "nessi test-error N501"
nessi test-error N501

# Example 6: Custom Error Details and Suggestions
print_header "CUSTOM ERROR DETAILS AND SUGGESTIONS"

print_test "Custom error details" "nessi test-error N201 --details 'Custom configuration error details'"
nessi test-error N201 --details "Custom configuration error details"

print_test "Custom error suggestion" "nessi test-error N201 --suggestion 'Try updating your config file'"
nessi test-error N201 --suggestion "Try updating your config file"

# Example 7: Resolvable Errors
print_header "RESOLVABLE ERRORS"

print_test "Resolvable path error" "nessi test-error N101 --resolvable"
echo "This will create a resolvable path error"
nessi test-error N101 --resolvable

print_test "Resolvable configuration error" "nessi test-error N201 --resolvable"
echo "This will create a resolvable configuration error"
nessi test-error N201 --resolvable

print_header "ERROR HANDLING EXAMPLES COMPLETED"
echo "These examples demonstrate the comprehensive error handling capabilities of Nessi"
echo "For more information, see the documentation at: https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md"
