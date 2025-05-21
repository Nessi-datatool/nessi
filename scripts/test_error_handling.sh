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

# ===== RESTORE ENVIRONMENT =====
print_header "TESTS COMPLETED"

# Restore environment variables
export DATABRICKS_HOST=$SAVED_HOST
export DATABRICKS_TOKEN=$SAVED_TOKEN
export DATABRICKS_WORKSPACE_ID=$SAVED_WORKSPACE_ID

echo "Environment variables restored to original values."
echo "All error handling tests completed."
