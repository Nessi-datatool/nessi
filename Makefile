# Nessi.dev Makefile
# Create a Makefile for Nessi.dev with these targets:
# 1. build: Build the Go binary
#    - Set proper Go build flags
#    - Create output in ./bin directory
# 2. test: Run all tests
#    - Unit tests
#    - Integration tests
#    - Generate coverage report
# 3. lint: Run linters
#    - golangci-lint
#    - Format checking
# 4. docker: Build Docker images
#    - Go image
#    - Python image (if enabled)
# 5. run: Run locally for development
#    - Set development environment
#    - Watch for file changes
# 6. clean: Clean build artifacts
# 7. deps: Install development dependencies
# 8. docs: Generate documentation
# Include help target explaining all commands
# Use proper variable definitions at the top
# Include .PHONY declarations for non-file targets

# Variables
BINARY_NAME=nessi
GO=go
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
GOFILES=$(shell find . -name "*.go" -type f)
GOTEST=$(GO) test
GOVET=$(GO) vet
GOLINT=golangci-lint
DOCKER=docker
DOCKER_COMPOSE=docker-compose
PYTHON=python3
PIP=pip3
AIR=air
GOMODULES=on
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

# Version information
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

# Build flags
BUILD_FLAGS=-trimpath -mod=vendor

# Test flags
TEST_FLAGS=-v -race -coverprofile=coverage.out
INTEGRATION_TEST_FLAGS=-v -tags=integration

# Default target
.PHONY: all
all: format-check go-mod-check clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GO) build $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/nessi

# Run the application in development mode with hot reload
.PHONY: run
run:
	@echo "Starting development server with hot reload..."
	$(AIR) -c .air.toml

# Run the application in production mode
.PHONY: run-prod
run-prod: build
	@echo "Starting production server..."
	./bin/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out
	rm -rf vendor/
	find . -name "*.test" -type f -delete

# Run all tests
.PHONY: test
test: test-unit test-integration

# Run unit tests
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) $(TEST_FLAGS) ./...

# Run short tests (skips long-running tests)
.PHONY: test-short
test-short:
	@echo "Running short tests (skipping long-running tests)..."
	$(GOTEST) -short -v ./...

# Run integration tests
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) $(INTEGRATION_TEST_FLAGS) ./tests/integration/...

# Generate test coverage report
.PHONY: test-coverage
test-coverage: test-unit
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run linters
.PHONY: lint
lint: lint-go lint-python

# Run Go linters
.PHONY: lint-go
lint-go: format-check go-mod-check
	@echo "Running Go linters..."
	$(GOLINT) run

# Format Go files
.PHONY: format
format:
	@echo "Formatting Go files..."
	$(GO) fmt ./...

# Check Go formatting
.PHONY: format-check
format-check:
	@echo "Checking Go formatting..."
	@UNFORMATTED=$$(gofmt -l .) && [ -z "$$UNFORMATTED" ] || (echo "The following files need formatting:" && echo "$$UNFORMATTED" && exit 1)

# Format go.mod file
.PHONY: go-mod-format
go-mod-format:
	@echo "Formatting go.mod file..."
	$(GO) mod edit -fmt

# Check go.mod formatting
.PHONY: go-mod-check
go-mod-check:
	@echo "Checking go.mod formatting..."
	@cp go.mod go.mod.orig && $(GO) mod edit -fmt && diff -q go.mod go.mod.orig > /dev/null || (echo "go.mod file needs formatting" && rm go.mod.orig && exit 1)
	@rm go.mod.orig

# Run Python linters
.PHONY: lint-python
lint-python:
	@echo "Running Python linters..."
	$(PYTHON) -m flake8 python/
	$(PYTHON) -m black --check python/

# Build Docker images
.PHONY: docker
docker: docker-go docker-python

# Build Go Docker image
.PHONY: docker-go
docker-go:
	@echo "Building Go Docker image..."
	$(DOCKER) build -t nessi-dev:$(VERSION) -f docker/go/Dockerfile .

# Build Python Docker image
.PHONY: docker-python
docker-python:
	@echo "Building Python Docker image..."
	$(DOCKER) build -t nessi-dev-python:$(VERSION) -f docker/python/Dockerfile .

# Run with Docker Compose
.PHONY: compose-up
compose-up:
	@echo "Starting services with Docker Compose..."
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml up -d

# Stop Docker Compose services
.PHONY: compose-down
compose-down:
	@echo "Stopping Docker Compose services..."
	$(DOCKER_COMPOSE) -f docker/docker-compose.yml down

# Install development dependencies
.PHONY: deps
deps: deps-go deps-python

# Install Go dependencies
.PHONY: deps-go
deps-go:
	@echo "Installing Go dependencies..."
	$(GO) mod download
	$(GO) mod vendor
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/go-delve/delve/cmd/dlv@latest
	$(GO) install github.com/cosmtrek/air@latest

# Install Python dependencies
.PHONY: deps-python
deps-python:
	@echo "Installing Python dependencies..."
	$(PIP) install -r python/requirements.txt
	$(PIP) install -r python/requirements-dev.txt

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	$(GO) doc -all ./...
	$(PYTHON) -m pdoc --html python/nessi -o docs/python

# Help target
.PHONY: help
help:
	@echo "Nessi.dev Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Clean and build the application"
	@echo "  build          - Build the Go binary"
	@echo "  run            - Run with hot reload (development)"
	@echo "  run-prod       - Run in production mode"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests"
	@echo "  test-short     - Run tests with -short flag (skips long-running tests)"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage  - Generate test coverage report"
	@echo "  lint           - Run all linters"
	@echo "  lint-go        - Run Go linters"
	@echo "  lint-python    - Run Python linters"
	@echo "  format         - Format Go files"
	@echo "  format-check   - Check if Go files need formatting"
	@echo "  go-mod-format  - Format go.mod file"
	@echo "  go-mod-check   - Check if go.mod file needs formatting"
	@echo "  docker         - Build all Docker images"
	@echo "  docker-go      - Build Go Docker image"
	@echo "  docker-python  - Build Python Docker image"
	@echo "  compose-up     - Start services with Docker Compose"
	@echo "  compose-down   - Stop Docker Compose services"
	@echo "  deps           - Install all dependencies"
	@echo "  deps-go        - Install Go dependencies"
	@echo "  deps-python    - Install Python dependencies"
	@echo "  docs           - Generate documentation"
	@echo "  help           - Show this help message"

# Security targets
.PHONY: security-certs security-check security-test

# Generate development certificates
security-certs:
	@echo "Generating development certificates..."
	@./scripts/generate_certs.sh

# Run security checks
security-check:
	@echo "Running security checks..."
	@go vet ./...
	@go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

# Run security tests
security-test:
	@echo "Running security tests..."
	@go test -v ./internal/security/...

# New targets
.PHONY: all build test lint clean run docker-build docker-run generate-certs

# Variables
VERSION=0.1.0
CONFIG_PATH=config/config.yaml

# Default target
all: lint test build

# Build the application
build:
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/nessi

# Run tests
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

# Run linters
lint:
	$(GOLINT) run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Run the application
run:
	$(GO) run ./cmd/nessi

# Build Docker image
docker-build:
	$(DOCKER) build -t $(BINARY_NAME):$(VERSION) .

# Run Docker container
docker-run:
	$(DOCKER) run -p 8080:8080 -v $(PWD)/config:/app/config $(BINARY_NAME):$(VERSION)

# Generate self-signed certificates
generate-certs:
	./scripts/generate_certs.sh

# Install development tools
install-tools:
	# Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

	# Install other development tools
	go install github.com/golang/mock/mockgen@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/cosmtrek/air@latest

# Create necessary directories
setup:
	mkdir -p bin
	mkdir -p data/delta/tables
	mkdir -p data/delta/logs
	mkdir -p data/quality/reports
	mkdir -p data/reports
	mkdir -p logs
	mkdir -p certs

# Initialize development environment
init: install-tools setup generate-certs

# Help target
help:
	@echo "Available targets:"
	@echo "  all            - Run lint, test, and build"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests with coverage"
	@echo "  lint           - Run linters"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Run the application"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  generate-certs - Generate self-signed certificates"
	@echo "  install-tools  - Install development tools"
	@echo "  setup          - Create necessary directories"
	@echo "  init           - Initialize development environment"
	@echo "  help           - Show this help message" 