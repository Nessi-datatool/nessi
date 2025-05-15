#!/bin/bash

# Define colors for better readability
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
NC="\033[0m" # No Color

# Default test mode and package
TEST_MODE="standard"
PACKAGE="./..."
SKIP_CACHE=false
SPECIFIC_FILES=""

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
    --short)
      TEST_MODE="short"
      shift
      ;;
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
      echo "  --fast         Run only fast tests (< 5 seconds)"
      echo "  --short        Run only specific short test files"
      echo "  --all          Run all tests including long-running ones"
      echo "  --package PKG  Specify a package to test (e.g., ./pkg/dbt/...)"
      echo "  --no-cache     Disable test caching (adds -count=1)"
      echo "  --files FILES  Run specific test files (comma-separated)"
      echo "  --help         Show this help message"
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

# Create cache-busting option if needed
CACHE_OPTION=""
if [ "$SKIP_CACHE" = true ]; then
  CACHE_OPTION="-count=1"
fi

# Set environment variables based on test mode
if [ "$TEST_MODE" = "fast" ]; then
  echo -e "${YELLOW}Running FAST tests only...${NC}"
  export NESSI_SKIP_LONG_TESTS=true
  export NESSI_TEST_TIMEOUT_MS=5000
  export NESSI_FAST_TESTS=1
elif [ "$TEST_MODE" = "short" ]; then
  echo -e "${YELLOW}Running specific SHORT test files only...${NC}"
  export NESSI_SKIP_LONG_TESTS=true
  export NESSI_TEST_TIMEOUT_MS=1000
  export NESSI_FAST_TESTS=1
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

# Run tests based on the selected mode
if [ "$TEST_MODE" = "short" ]; then
  # Run specific short test files (similar to the old run_short_tests.sh)
  echo -e "${CYAN}Running specific short test files...${NC}"
  
  if [ -n "$SPECIFIC_FILES" ]; then
    # Run the specified test files
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
    # Default short test files if none specified
    echo -e "${GREEN}Running dashboard short tests...${NC}"
    go test -v $CACHE_OPTION ./pkg/monitoring/dashboard/alerts_short_test.go
    
    echo -e "${GREEN}Running monitoring short tests...${NC}"
    go test -v $CACHE_OPTION ./pkg/monitoring/export_short_test.go
  fi
elif [ "$TEST_MODE" = "fast" ]; then
  # Create a list of packages with long-running tests to skip
  echo -e "${CYAN}Running fast tests, skipping long-running test packages...${NC}"
  
  LONG_TEST_PACKAGES=(
    "./pkg/monitoring/dashboard"
    "./pkg/monitoring"
    "./pkg/datalake"
    "./pkg/cloud"
    "./pkg/catalog"
  )
  
  if [ "$PACKAGE" = "./..." ]; then
    # Run tests for all packages, skipping the ones with long tests
    TEST_PACKAGES=$(go list ./pkg/... | grep -v vendor)
    
    for pkg in $TEST_PACKAGES; do
      # Check if this package should be skipped
      skip=false
      for long_pkg in "${LONG_TEST_PACKAGES[@]}"; do
        if [[ "$pkg" == *"$long_pkg"* ]]; then
          skip=true
          break
        fi
      done
      
      if [ "$skip" = true ]; then
        echo -e "${YELLOW}Skipping package with long tests: $pkg${NC}"
      else
        echo -e "${GREEN}Running tests for package: $pkg${NC}"
        go test -v $CACHE_OPTION "$pkg" -timeout=5s
        if [ $? -ne 0 ]; then
          echo -e "${RED}Tests failed for package: $pkg${NC}"
          exit 1
        fi
      fi
    done
  else
    # Run tests for a specific package with the -short flag
    echo -e "${GREEN}Running tests for package: $PACKAGE${NC}"
    go test -short -v $CACHE_OPTION "$PACKAGE"
  fi
else
  # Run standard or all tests
  TAG_OPTION=""
  if [ "$TEST_MODE" != "all" ]; then
    TAG_OPTION="-tags=skiplong"
  fi
  
  echo -e "${GREEN}Running tests: go test $TAG_OPTION $CACHE_OPTION $PACKAGE -v${NC}"
  go test $TAG_OPTION $CACHE_OPTION $PACKAGE -v
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
