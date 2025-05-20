# Create a Dockerfile for Nessi.dev Go services with these requirements:
# 1. Use multi-stage build:
#    - Stage 1: Build the Go binary
#      - Use golang:1.22-alpine as base
#      - Install build dependencies
#      - Copy source code
#      - Set up Go modules
#      - Build static binary with proper flags
#    - Stage 2: Create runtime image
#      - Use alpine:latest
#      - Install runtime dependencies
#      - Copy binary from build stage
#      - Set up non-root user
# 2. Configure defaults:
#    - Working directory: /app
#    - Expose port 8080
#    - Set environment variables
#    - Volume mount point for data: /data
#    - Volume mount point for config: /config
# 3. Add metadata:
#    - Labels for version, maintainer, etc.
#    - Health check using nessi health command
# 4. Set entrypoint to the nessi binary
# Optimize for small image size
# Include only necessary runtime dependencies
# Set up proper permissions for security

# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# - CGO_ENABLED=0 for static linking
# - -ldflags="-s -w" for smaller binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/nessi ./cmd/nessi

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget

# Create non-root user
RUN adduser -D -g '' nessi

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/nessi /app/nessi

# Create necessary directories
RUN mkdir -p /data /config && \
    chown -R nessi:nessi /app /data /config

# Switch to non-root user
USER nessi

# Set environment variables
ENV TZ=UTC \
    NESSI_CONFIG_PATH=/config \
    NESSI_DATA_PATH=/data \
    NESSI_EXTENSIONS_PATH=/config/extensions

# Expose port
EXPOSE 8080

# Set volume mount points
VOLUME ["/data", "/config"]

# Set healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Add metadata
LABEL org.opencontainers.image.title="Nessi.dev" \
      org.opencontainers.image.description="Data quality and profiling tool for Delta tables" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.source="https://github.com/meisi/nessi-dev" \
      org.opencontainers.image.licenses="MIT"

# Run the application
ENTRYPOINT ["/app/nessi"]

# Default command
CMD ["serve"]