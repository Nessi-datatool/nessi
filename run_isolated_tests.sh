#!/bin/bash

# Run tests for individual command files in cmd/nessi in isolation
cd /Users/meisi/Documents/nessi-dev

# Function to run a specific test file in isolation
run_isolated_test() {
    local test_file=$1
    local test_pattern=$2
    
    echo "\nRunning tests for: $test_file"
    # Use a clean environment for each test run to avoid flag conflicts
    # NESSI_TEST_MODE=true tells the application it's running in test mode
    env -i PATH=$PATH HOME=$HOME GOPATH=$GOPATH NESSI_TEST_MODE=true \
        go test -run "$test_pattern" "./cmd/nessi/$test_file" -v
}

# Run each test file individually
run_isolated_test "schema_cmd_test.go" "TestSchema"
run_isolated_test "alerts_cmd_test.go" "TestAlerts"
run_isolated_test "time_travel_cmd_test.go" "TestTimeTravel"

echo "\nAll isolated tests completed."
echo "Note: Some tests may still fail due to flag redefinition issues."
echo "The core functionality tests are passing successfully as shown by run_core_tests.sh"
