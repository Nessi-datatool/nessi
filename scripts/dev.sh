#!/bin/bash

# Exit on error
set -e

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
mkdir -p ssl
mkdir -p prometheus
mkdir -p grafana

# Generate SSL certificates if they don't exist
if [ ! -f "ssl/cert.pem" ] || [ ! -f "ssl/key.pem" ]; then
    echo -e "${GREEN}Generating SSL certificates...${NC}"
    openssl req -x509 -newkey rsa:4096 -nodes -out ssl/cert.pem -keyout ssl/key.pem -days 365 -subj "/CN=localhost"
fi

# Build and start development containers
echo -e "${GREEN}Building and starting development containers...${NC}"
docker-compose -f docker-compose.dev.yml up -d --build

echo -e "${GREEN}Development environment is ready!${NC}"
echo -e "${YELLOW}Access the following services:${NC}"
echo -e "  - Backend API: https://localhost:8000"
echo -e "  - Prometheus: http://localhost:9090"
echo -e "  - Grafana: http://localhost:3000"
echo -e "\n${YELLOW}To run tests:${NC}"
echo -e "  docker-compose -f docker-compose.test.yml up --build"
echo -e "\n${YELLOW}To stop the development environment:${NC}"
echo -e "  docker-compose -f docker-compose.dev.yml down" 