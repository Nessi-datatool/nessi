#!/bin/bash

# Exit on error
set -e

# Function to display usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help        Show this help message"
    echo "  -v, --verbose     Run tests in verbose mode"
    echo "  -k, --keep        Keep containers running after tests"
    echo "  -f, --file FILE   Run specific test file"
    echo "  -m, --mark MARK   Run tests with specific marker"
    echo "  -c, --coverage    Generate coverage report"
    exit 1
}

# Parse command line arguments
VERBOSE=""
KEEP=""
TEST_FILE=""
TEST_MARK=""
COVERAGE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        -k|--keep)
            KEEP="true"
            shift
            ;;
        -f|--file)
            TEST_FILE="$2"
            shift 2
            ;;
        -m|--mark)
            TEST_MARK="-m $2"
            shift 2
            ;;
        -c|--coverage)
            COVERAGE="--cov=src/nessi --cov-report=html --cov-report=xml"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Build test container
echo "Building test container..."
docker-compose build nessi-test

# Run tests
if [ -n "$TEST_FILE" ]; then
    echo "Running specific test file: $TEST_FILE"
    docker-compose run --rm nessi-test pytest $VERBOSE $TEST_MARK $COVERAGE $TEST_FILE
else
    echo "Running all tests..."
    docker-compose run --rm nessi-test pytest $VERBOSE $TEST_MARK $COVERAGE tests/
fi

# Clean up if not keeping containers
if [ -z "$KEEP" ]; then
    echo "Cleaning up containers..."
    docker-compose down
fi 