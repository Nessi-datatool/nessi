#!/bin/bash

# Exit on error
set -e

# Function to display usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help        Show this help message"
    echo "  -b, --build       Build containers"
    echo "  -r, --run         Run containers"
    echo "  -s, --stop        Stop containers"
    echo "  -l, --logs        Show container logs"
    echo "  -t, --test        Run tests"
    echo "  -c, --clean       Clean up containers and volumes"
    exit 1
}

# Parse command line arguments
BUILD=false
RUN=false
STOP=false
LOGS=false
TEST=false
CLEAN=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -b|--build)
            BUILD=true
            shift
            ;;
        -r|--run)
            RUN=true
            shift
            ;;
        -s|--stop)
            STOP=true
            shift
            ;;
        -l|--logs)
            LOGS=true
            shift
            ;;
        -t|--test)
            TEST=true
            shift
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Build containers
if [ "$BUILD" = true ]; then
    echo "Building containers..."
    docker-compose -f docker/docker-compose.dev.yml build
fi

# Run containers
if [ "$RUN" = true ]; then
    echo "Starting containers..."
    docker-compose -f docker/docker-compose.dev.yml up -d
fi

# Stop containers
if [ "$STOP" = true ]; then
    echo "Stopping containers..."
    docker-compose -f docker/docker-compose.dev.yml down
fi

# Show logs
if [ "$LOGS" = true ]; then
    echo "Showing container logs..."
    docker-compose -f docker/docker-compose.dev.yml logs -f
fi

# Run tests
if [ "$TEST" = true ]; then
    echo "Running tests..."
    docker-compose -f docker/docker-compose.dev.yml run --rm nessi-test pytest tests/
fi

# Clean up
if [ "$CLEAN" = true ]; then
    echo "Cleaning up..."
    docker-compose -f docker/docker-compose.dev.yml down -v
    docker system prune -f
fi

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Starting development environment setup...${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

# Create necessary directories
mkdir -p config/ssl
mkdir -p config/prometheus
mkdir -p config/grafana

# Generate SSL certificates if they don't exist
if [ ! -f "config/ssl/cert.pem" ] || [ ! -f "config/ssl/key.pem" ]; then
    echo -e "${GREEN}Generating SSL certificates...${NC}"
    openssl req -x509 -newkey rsa:4096 -nodes -out config/ssl/cert.pem -keyout config/ssl/key.pem -days 365 -subj "/CN=localhost"
fi

# Build and start development containers
echo -e "${GREEN}Building and starting development containers...${NC}"
docker-compose -f docker/docker-compose.dev.yml up -d --build

echo -e "${GREEN}Development environment is ready!${NC}"
echo -e "${YELLOW}Access the following services:${NC}"
echo -e "  - Backend API: https://localhost:8443"
echo -e "  - Prometheus: http://localhost:9090"
echo -e "  - Grafana: http://localhost:3000"
echo -e "\n${YELLOW}To run tests:${NC}"
echo -e "  docker-compose -f docker/docker-compose.test.yml up --build"
echo -e "\n${YELLOW}To stop the development environment:${NC}"
echo -e "  docker-compose -f docker/docker-compose.dev.yml down" 