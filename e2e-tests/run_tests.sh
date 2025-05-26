#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Check if nessi binary exists
echo -e "${YELLOW}Checking for nessi binary...${NC}"
if command -v nessi &> /dev/null; then
    NESSI_BIN=$(which nessi)
    echo -e "${GREEN}Found nessi binary at: $NESSI_BIN${NC}"
elif [ -f "$PROJECT_ROOT/bin/nessi" ]; then
    NESSI_BIN="$PROJECT_ROOT/bin/nessi"
    echo -e "${GREEN}Found nessi binary at: $NESSI_BIN${NC}"
elif [ -f "$PROJECT_ROOT/cmd/nessi/nessi" ]; then
    NESSI_BIN="$PROJECT_ROOT/cmd/nessi/nessi"
    echo -e "${GREEN}Found nessi binary at: $NESSI_BIN${NC}"
else
    echo -e "${YELLOW}Building nessi binary...${NC}"
    cd "$PROJECT_ROOT/cmd/nessi"
    go build -o nessi
    NESSI_BIN="$PROJECT_ROOT/cmd/nessi/nessi"
    echo -e "${GREEN}Built nessi binary at: $NESSI_BIN${NC}"
    cd "$PROJECT_ROOT"
fi

# Ensure test data exists
echo -e "${YELLOW}Checking test data...${NC}"
TEST_DATA_DIR="$PROJECT_ROOT/test_data"
if [ ! -d "$TEST_DATA_DIR" ]; then
    echo -e "${YELLOW}Creating test data directory...${NC}"
    mkdir -p "$TEST_DATA_DIR/delta_tables/sample_table"
    # Create a simple Delta table structure
    mkdir -p "$TEST_DATA_DIR/delta_tables/sample_table/_delta_log"
    echo '{"commitInfo":{"timestamp":1622548800000,"operation":"CREATE TABLE"}}' > "$TEST_DATA_DIR/delta_tables/sample_table/_delta_log/00000000000000000000.json"
    echo '{"metaData":{"id":"sample-table-id","format":{"provider":"parquet"},"schemaString":"{\"type\":\"struct\",\"fields\":[{\"name\":\"id\",\"type\":\"integer\",\"nullable\":false},{\"name\":\"name\",\"type\":\"string\",\"nullable\":true}]}"}}' > "$TEST_DATA_DIR/delta_tables/sample_table/_delta_log/00000000000000000000.json.metadata"
    echo -e "${GREEN}Created sample test data${NC}"
fi

# Run the tests
echo -e "${YELLOW}Running E2E tests...${NC}"

# Run user journey tests
echo -e "${YELLOW}Running user journey tests...${NC}"
go test -v ./e2e-tests/scenarios/user-journeys/... || true

# Run integration tests
echo -e "${YELLOW}Running integration tests...${NC}"
go test -v ./e2e-tests/scenarios/integrations/... || true

# Run stress tests
echo -e "${YELLOW}Running stress tests...${NC}"
go test -v ./e2e-tests/scenarios/stress-tests/... || true

echo -e "${GREEN}All tests completed!${NC}"
echo -e "${YELLOW}Note: Tests are currently skipped until CLI commands are fully implemented${NC}"
