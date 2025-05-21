#!/bin/bash

# Test script for error handling in Nessi
# This script intentionally triggers various error conditions to verify error handling

# Function to print section headers
print_header() {
  echo "\n============================================="
  echo "$1"
  echo "============================================="
}

# Function to print test case
print_test() {
  echo "\n>> Testing: $1"
  echo "Command: $2"
  echo "-------------------------------------------"
}

# Create a temporary directory for test files
TEST_DIR="$(mktemp -d)"
echo "Created temporary test directory: $TEST_DIR"

# Clean up function
cleanup() {
  echo "\nCleaning up test directory..."
  rm -rf "$TEST_DIR"
  echo "Done."
}

# Register cleanup function to run on script exit
trap cleanup EXIT

# Save current environment variables
SAVED_HOST=$DATABRICKS_HOST
SAVED_TOKEN=$DATABRICKS_TOKEN
SAVED_WORKSPACE_ID=$DATABRICKS_WORKSPACE_ID

# ===== DELTA LAKE ERROR HANDLING TESTS =====
print_header "DELTA LAKE ERROR HANDLING TESTS"

# Test invalid path error
print_test "Invalid path error" "nessi schema show --path $TEST_DIR/nonexistent/path"
nessi schema show --path $TEST_DIR/nonexistent/path

# Test not a Delta table error
print_test "Not a Delta table error" "nessi schema show --path $TEST_DIR"
mkdir -p $TEST_DIR/regular_dir
nessi schema show --path $TEST_DIR/regular_dir

# Test corrupted Delta table
print_test "Corrupted Delta table error" "nessi schema show --path $TEST_DIR/corrupted_delta"
mkdir -p $TEST_DIR/corrupted_delta/_delta_log
touch $TEST_DIR/corrupted_delta/_delta_log/corrupted_file.json
nessi schema show --path $TEST_DIR/corrupted_delta

# ===== DATABRICKS INTEGRATION ERROR HANDLING TESTS =====
print_header "DATABRICKS INTEGRATION ERROR HANDLING TESTS"

# Test authentication error
print_test "Authentication error" "nessi catalog list --type databricks (with invalid token)"
export DATABRICKS_TOKEN="invalid_token"
nessi catalog list --type databricks

# Test workspace ID validation
print_test "Workspace ID validation" "nessi catalog list --type databricks (with missing workspace ID)"
unset DATABRICKS_WORKSPACE_ID
nessi catalog list --type databricks

# Test connection error
print_test "Connection error" "nessi catalog list --type databricks (with invalid host)"
export DATABRICKS_HOST="https://nonexistent-workspace.cloud.databricks.com"
nessi catalog list --type databricks

# Test resource not found error
print_test "Resource not found error" "nessi catalog schemas --type databricks --catalog nonexistent_catalog"
export DATABRICKS_HOST=$SAVED_HOST
export DATABRICKS_TOKEN=$SAVED_TOKEN
export DATABRICKS_WORKSPACE_ID=$SAVED_WORKSPACE_ID
nessi catalog schemas --type databricks --catalog nonexistent_catalog

# ===== INTEGRATION ERROR HANDLING TESTS =====
print_header "INTEGRATION ERROR HANDLING TESTS"

# Test Delta Lake and Databricks integration error
print_test "Delta Lake and Databricks integration error" "nessi catalog import --type databricks --path $TEST_DIR/not_a_delta_table"
nessi catalog import --type databricks --path $TEST_DIR/regular_dir

# ===== SCHEMA VALIDATION ERROR HANDLING TESTS =====
print_header "SCHEMA VALIDATION ERROR HANDLING TESTS"

# Test schema validation error
print_test "Schema validation error" "nessi validate-schema --path $TEST_DIR/invalid_schema.json"
cat > $TEST_DIR/invalid_schema.json << EOF
{
  "fields": [
    {"name": "id", "type": "invalid_type"},
    {"name": "name", "type": "string"}
  ]
}
EOF
nessi validate-schema --path $TEST_DIR/invalid_schema.json

# ===== ERROR TELEMETRY TESTS =====
print_header "ERROR TELEMETRY TESTS"

# Test error telemetry recording
print_test "Error telemetry recording" "nessi test-error N101 --telemetry"
nessi test-error N101 --telemetry

# Test error telemetry export
print_test "Error telemetry export" "nessi test-error N201 --telemetry --export $TEST_DIR/error_telemetry.json"
nessi test-error N201 --telemetry --export $TEST_DIR/error_telemetry.json

# Check if telemetry file was created
if [ -f "$TEST_DIR/error_telemetry.json" ]; then
  echo "✅ Telemetry file created successfully"
  echo "Telemetry file contents:"
  cat "$TEST_DIR/error_telemetry.json"
else
  echo "❌ Telemetry file was not created"
fi

# ===== ERROR SUGGESTIONS TESTS =====
print_header "ERROR SUGGESTIONS TESTS"

# Test error suggestions for path errors
print_test "Error suggestions for path errors" "nessi test-error N101"
nessi test-error N101

# Test error suggestions for configuration errors
print_test "Error suggestions for configuration errors" "nessi test-error N201"
nessi test-error N201

# Test error suggestions for authentication errors
print_test "Error suggestions for authentication errors" "nessi test-error N301"
nessi test-error N301

# ===== RESOLVABLE ERRORS TESTS =====
print_header "RESOLVABLE ERRORS TESTS"

# Test resolvable path error
print_test "Resolvable path error" "nessi test-error N101 --resolvable"
nessi test-error N101 --resolvable

# Test resolvable configuration error
print_test "Resolvable configuration error" "nessi test-error N201 --resolvable"
nessi test-error N201 --resolvable

# ===== RESTORE ENVIRONMENT =====
print_header "TESTS COMPLETED"

# Restore environment variables
export DATABRICKS_HOST=$SAVED_HOST
export DATABRICKS_TOKEN=$SAVED_TOKEN
export DATABRICKS_WORKSPACE_ID=$SAVED_WORKSPACE_ID

echo "Environment variables restored to original values."
echo "All error handling tests completed."
