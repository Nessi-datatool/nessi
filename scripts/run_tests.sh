#!/bin/bash

# Define colors for better readability
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

# Default test mode
TEST_MODE="standard"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --fast)
      TEST_MODE="fast"
      shift
      ;;
    --all)
      TEST_MODE="all"
      shift
      ;;
    --help)
      echo -e "${BLUE}Nessi.dev Test Runner${NC}"
      echo "Usage: ./scripts/run_tests.sh [options]"
      echo ""
      echo "Options:"
      echo "  --fast     Run only fast tests (< 5 seconds)"
      echo "  --all      Run all tests including long-running ones"
      echo "  --help     Show this help message"
      echo ""
      echo "Default: Run standard tests (< 30 seconds)"
      exit 0
      ;;
    *)
      echo -e "${RED}Unknown option: $1${NC}"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

# Set environment variables based on test mode
if [ "$TEST_MODE" = "fast" ]; then
  echo -e "${YELLOW}Running FAST tests only...${NC}"
  export NESSI_SKIP_LONG_TESTS=true
  export NESSI_TEST_TIMEOUT_MS=5000
elif [ "$TEST_MODE" = "all" ]; then
  echo -e "${YELLOW}Running ALL tests including long-running ones...${NC}"
  export NESSI_SKIP_LONG_TESTS=false
  export NESSI_TEST_TIMEOUT_MS=300000
else
  echo -e "${YELLOW}Running STANDARD tests...${NC}"
  export NESSI_SKIP_LONG_TESTS=true
  export NESSI_TEST_TIMEOUT_MS=30000
fi

# Common environment variables
export NESSI_TEST_DEBUG=true
export GO_TEST=true

# Run tests with appropriate tags
if [ "$TEST_MODE" = "fast" ]; then
  echo -e "${GREEN}Running tests with -tags=skiplong...${NC}"
  go test -tags=skiplong ./... -v
elif [ "$TEST_MODE" = "all" ]; then
  echo -e "${GREEN}Running all tests...${NC}"
  go test ./... -v
else
  echo -e "${GREEN}Running standard tests...${NC}"
  go test -tags=skiplong ./... -v
fi

# Check test result
if [ $? -eq 0 ]; then
  echo -e "${GREEN}All tests passed!${NC}"
  exit 0
else
  echo -e "${RED}Some tests failed!${NC}"
  exit 1
fi
