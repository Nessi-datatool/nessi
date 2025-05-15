#!/bin/bash

# Define colors for better readability
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

# Default test settings
PACKAGE="./..."
SKIP_CACHE=false
SPECIFIC_FILES=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --package)
      PACKAGE="$2"
      shift 2
      ;;
    --no-cache)
      SKIP_CACHE=true
      shift
      ;;
    --files)
      SPECIFIC_FILES="$2"
      shift 2
      ;;
    --help)
      echo -e "${BLUE}Nessi.dev Test Runner${NC}"
      echo "Usage: ./scripts/run_tests.sh [options]"
      echo ""
      echo "Options:"
      echo "  --package PKG  Specify a package to test (e.g., ./pkg/dbt/...)"
      echo "  --no-cache     Disable test caching (adds -count=1)"
      echo "  --files FILES  Run specific test files (comma-separated)"
      echo "  --help         Show this help message"
      echo ""
      echo "Default: Run all tests"
      exit 0
      ;;
    *)
      echo -e "${RED}Unknown option: $1${NC}"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

# Create cache-busting option if needed
CACHE_OPTION=""
if [ "$SKIP_CACHE" = true ]; then
  CACHE_OPTION="-count=1"
fi

# Set environment variables
echo -e "${YELLOW}Running tests...${NC}"
export NESSI_TEST_TIMEOUT_MS=30000
export NESSI_TEST_DEBUG=true
export GO_TEST=true

# Run tests
if [ -n "$SPECIFIC_FILES" ]; then
  # Run the specified test files
  echo -e "${CYAN}Running specific test files...${NC}"
  IFS=',' read -ra FILES <<< "$SPECIFIC_FILES"
  for file in "${FILES[@]}"; do
    echo -e "${GREEN}Running test file: $file${NC}"
    go test -v $CACHE_OPTION "$file"
    if [ $? -ne 0 ]; then
      echo -e "${RED}Tests failed for file: $file${NC}"
      exit 1
    fi
  done
else
  # Run tests for the specified package
  echo -e "${GREEN}Running tests for package: $PACKAGE${NC}"
  go test -v $CACHE_OPTION "$PACKAGE"
fi

# Check final test result
TEST_RESULT=$?
if [ $TEST_RESULT -eq 0 ]; then
  echo -e "${GREEN}All tests passed!${NC}"
  exit 0
else
  echo -e "${RED}Some tests failed!${NC}"
  exit 1
fi
