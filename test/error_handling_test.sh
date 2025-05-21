#!/bin/bash

# Test script for error handling in Nessi CLI

# Set up colors for output
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
NC="\033[0m" # No Color

# Initialize counters
failed_count=0
passed_count=0

# Build the CLI
echo -e "${YELLOW}Building Nessi CLI...${NC}"
cd "$(dirname "$0")/.." || exit 1
go build -o ./bin/nessi ./cmd/nessi

if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}Build successful!${NC}"

# Function to run a test case
run_test() {
    local test_name=$1
    local command=$2
    local expected_error_code=$3
    local expected_output=$4
    
    echo -e "\n${YELLOW}Running test: ${test_name}${NC}"
    echo -e "Command: ${command}"
    
    # Run the command and capture output
    output=$(eval "$command" 2>&1)
    exit_code=$?
    
    echo -e "Output:\n$output"
    
    # Check exit code
    if [ -n "$expected_error_code" ] && [ $exit_code -eq $expected_error_code ]; then
        echo -e "${GREEN}✓ Exit code matches expected: $exit_code${NC}"
    elif [ -n "$expected_error_code" ]; then
        echo -e "${RED}✗ Exit code does not match expected: got $exit_code, expected $expected_error_code${NC}"
        failed_count=$((failed_count + 1))
        return 1
    fi
    
    # Check output
    if [ -n "$expected_output" ] && echo "$output" | grep -q "$expected_output"; then
        echo -e "${GREEN}✓ Output contains expected text${NC}"
    elif [ -n "$expected_output" ]; then
        echo -e "${RED}✗ Output does not contain expected text: $expected_output${NC}"
        failed_count=$((failed_count + 1))
        return 1
    fi
    
    passed_count=$((passed_count + 1))
    return 0
}

# Test cases

# Test 1: Invalid path error
run_test "Invalid path error" "./bin/nessi profile --table /nonexistent/path --interactive=false" 1 "path '/nonexistent/path' does not exist"

# Test 2: Invalid command error
run_test "Invalid command error" "./bin/nessi invalid-command --interactive=false" 1 "unknown command"

# Test 3: Missing required argument
run_test "Missing required argument" "./bin/nessi profile --interactive=false" 1 "table path is required"

# Test 4: Help command (should succeed)
run_test "Help command" "./bin/nessi --help" 0 "Nessi is a powerful tool for managing data quality"

# Test 5: Version command (should succeed)
run_test "Version command" "./bin/nessi --version" 0 ""

# Test 6: Info command (should succeed and show error handling)
run_test "Info command" "./bin/nessi info" 0 "Error handling: enhanced"

# Summary
echo -e "\n${YELLOW}Test Summary:${NC}"
echo -e "Total tests: 6"
echo -e "${GREEN}Passed: $passed_count${NC}"
if [ $failed_count -gt 0 ]; then
    echo -e "${RED}Failed: $failed_count${NC}"
    exit 1
else
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
fi
