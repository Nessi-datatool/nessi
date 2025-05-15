#!/bin/bash

# Define colors for better readability
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

# Common environment variables
export NESSI_TEST_DEBUG=true
export GO_TEST=true
export NESSI_SKIP_LONG_TESTS=true
export NESSI_TEST_TIMEOUT_MS=30000

# Skip the dashboard package but run the catalog package
echo -e "${CYAN}Running tests for catalog package while skipping dashboard package...${NC}"

# Run tests for the catalog package
echo -e "${GREEN}Running tests for package: ./pkg/catalog/...${NC}"
go test -v ./pkg/catalog/...

# Check final test result
TEST_RESULT=$?
if [ $TEST_RESULT -eq 0 ]; then
  echo -e "${GREEN}All catalog tests passed!${NC}"
  exit 0
else
  echo -e "${RED}Some catalog tests failed!${NC}"
  exit 1
fi
