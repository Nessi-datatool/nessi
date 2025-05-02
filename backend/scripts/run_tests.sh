#!/bin/bash

# Exit on error
set -e

# Function to display usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help           Show this help message"
    echo "  -t, --test-file      Run specific test file"
    echo "  -k, --test-keyword   Run tests matching keyword"
    echo "  -v, --verbose        Enable verbose output"
    echo "  --no-cache           Rebuild Docker image without cache"
}

# Default values
TEST_FILE=""
TEST_KEYWORD=""
VERBOSE=""
NO_CACHE=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -t|--test-file)
            TEST_FILE="$2"
            shift 2
            ;;
        -k|--test-keyword)
            TEST_KEYWORD="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE="-v"
            shift
            ;;
        --no-cache)
            NO_CACHE="--no-cache"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Build Docker image
echo "Building test image..."
docker-compose -f docker-compose.test.yml build $NO_CACHE test-runner

# Prepare test command
TEST_CMD="pytest $VERBOSE --cov=src"
if [ -n "$TEST_FILE" ]; then
    TEST_CMD="$TEST_CMD $TEST_FILE"
fi
if [ -n "$TEST_KEYWORD" ]; then
    TEST_CMD="$TEST_CMD -k $TEST_KEYWORD"
fi

# Run tests
echo "Running tests..."
docker-compose -f docker-compose.test.yml run --rm test-runner $TEST_CMD 