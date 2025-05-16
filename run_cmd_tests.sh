#!/bin/bash

# Run tests for individual command files in cmd/nessi
cd /Users/meisi/Documents/nessi-dev/cmd/nessi

# Function to run tests for a specific file pattern
run_test() {
    local pattern=$1
    echo "\nRunning tests matching pattern: $pattern"
    # Use a clean environment for each test run to avoid flag conflicts
    env -i PATH=$PATH HOME=$HOME GOPATH=$GOPATH go test -run "$pattern" -v 2>&1 | grep -v "no test files"
}

# Run tests for each major component
run_test "TestSchema"
run_test "TestAlerts"
run_test "TestTimeTravel"
run_test "TestFreshness"
run_test "TestValidate"
run_test "TestConfig"
run_test "TestPlugin"
run_test "TestWebhook"
run_test "TestRCA"

echo "\nAll tests completed."
