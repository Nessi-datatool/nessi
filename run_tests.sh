#!/bin/bash

# This is a wrapper script that calls the main test runner in the scripts directory
# It's kept in the root directory for backward compatibility and convenience

# Forward all arguments to the main script
./scripts/run_tests.sh "$@"
