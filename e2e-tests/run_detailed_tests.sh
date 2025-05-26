#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
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

# Verify test data
echo -e "${YELLOW}Verifying test data...${NC}"
TEST_DATA_DIR="$PROJECT_ROOT/test_data"
if [ -d "$TEST_DATA_DIR/delta_tables/sample_table" ]; then
    echo -e "${GREEN}Found sample Delta Lake table${NC}"
else
    echo -e "${RED}Sample Delta Lake table not found!${NC}"
    exit 1
fi

if [ -f "$TEST_DATA_DIR/schemas/sample_schema.json" ]; then
    echo -e "${GREEN}Found sample schema file${NC}"
else
    echo -e "${RED}Sample schema file not found!${NC}"
    exit 1
fi

if [ -f "$TEST_DATA_DIR/config/sample_config.yaml" ]; then
    echo -e "${GREEN}Found sample config file${NC}"
else
    echo -e "${RED}Sample config file not found!${NC}"
    exit 1
fi

# Check if we need to enable tests
if [ "$1" == "--enable" ]; then
    echo -e "${YELLOW}Enabling tests (removing skip directives)...${NC}"
    # Remove skip directives from test files
    sed -i '' 's/t.Skip("Skipping test until CLI commands are fully implemented")/\/\/ Test is now enabled/' "$PROJECT_ROOT/e2e-tests/scenarios/user-journeys/new_user_onboarding_test.go"
    sed -i '' 's/t.Skip("Skipping test until CLI commands are fully implemented")/\/\/ Test is now enabled/' "$PROJECT_ROOT/e2e-tests/scenarios/integrations/delta_lake_test.go"
    sed -i '' 's/t.Skip("Skipping test until CLI commands are fully implemented")/\/\/ Test is now enabled/' "$PROJECT_ROOT/e2e-tests/scenarios/stress-tests/error_handling_test.go"
    echo -e "${GREEN}Tests enabled!${NC}"
fi

# Run the tests with verbose output
echo -e "${BLUE}=======================================================${NC}"
echo -e "${BLUE}             RUNNING E2E TESTS WITH DETAILS            ${NC}"
echo -e "${BLUE}=======================================================${NC}"

# Run user journey tests
echo -e "${YELLOW}Running user journey tests...${NC}"
echo -e "${BLUE}-------------------------------------------------------${NC}"
go test -v ./e2e-tests/scenarios/user-journeys/... || true
echo -e "${BLUE}-------------------------------------------------------${NC}"

# Run integration tests
echo -e "${YELLOW}Running integration tests...${NC}"
echo -e "${BLUE}-------------------------------------------------------${NC}"
go test -v ./e2e-tests/scenarios/integrations/... || true
echo -e "${BLUE}-------------------------------------------------------${NC}"

# Run stress tests
echo -e "${YELLOW}Running stress tests...${NC}"
echo -e "${BLUE}-------------------------------------------------------${NC}"
go test -v ./e2e-tests/scenarios/stress-tests/... || true
echo -e "${BLUE}-------------------------------------------------------${NC}"

echo -e "${GREEN}All tests completed!${NC}"

# Check if tests are skipped
if grep -q "t.Skip" "$PROJECT_ROOT/e2e-tests/scenarios/user-journeys/new_user_onboarding_test.go"; then
    echo -e "${YELLOW}Note: Tests are currently skipped. Run with --enable to enable them.${NC}"
    echo -e "${YELLOW}Warning: Only enable tests when CLI commands are fully implemented.${NC}"
fi
