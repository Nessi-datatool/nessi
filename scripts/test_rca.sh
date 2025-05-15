#!/bin/bash

# Define colors for better readability
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
NC="\033[0m" # No Color

echo -e "${BLUE}Running Root Cause Analysis (RCA) Tests${NC}"
echo "================================================"

# Set environment variables for testing
export NESSI_TEST_DEBUG=true
export GO_TEST=true

# Run unit tests for the RCA package
echo -e "${YELLOW}Running RCA unit tests...${NC}"
go test -v ./pkg/rca/analyzer_test.go ./pkg/rca/analyzer.go
if [ $? -ne 0 ]; then
  echo -e "${RED}RCA unit tests failed!${NC}"
  exit 1
fi
echo -e "${GREEN}RCA unit tests passed!${NC}"
echo ""

# Run integration tests for the RCA package
echo -e "${YELLOW}Running RCA integration tests...${NC}"
go test -v ./pkg/rca/integration_test.go ./pkg/rca/analyzer.go
if [ $? -ne 0 ]; then
  echo -e "${RED}RCA integration tests failed!${NC}"
  exit 1
fi
echo -e "${GREEN}RCA integration tests passed!${NC}"
echo ""

# Run CLI command tests
echo -e "${YELLOW}Running RCA CLI command tests...${NC}"
go test -v ./cmd/nessi/rca_test.go ./cmd/nessi/rca.go
if [ $? -ne 0 ]; then
  echo -e "${RED}RCA CLI command tests failed!${NC}"
  exit 1
fi
echo -e "${GREEN}RCA CLI command tests passed!${NC}"
echo ""

# Run end-to-end test with sample data
echo -e "${YELLOW}Running RCA end-to-end test with sample data...${NC}"
echo "This test requires the nessi CLI to be built."
echo "Building nessi CLI..."
go build -o nessi ./cmd/nessi
if [ $? -ne 0 ]; then
  echo -e "${RED}Failed to build nessi CLI!${NC}"
  exit 1
fi

# Run the RCA command with sample data
./nessi rca anom-20250514-001 --format=json --output=rca_result.json
if [ $? -ne 0 ]; then
  echo -e "${RED}RCA end-to-end test failed!${NC}"
  exit 1
fi
echo -e "${GREEN}RCA end-to-end test passed!${NC}"
echo "RCA result saved to rca_result.json"
echo ""

echo -e "${GREEN}All RCA tests passed successfully!${NC}"
exit 0
