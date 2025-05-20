# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/nessi .

# Copy configuration and scripts
COPY config/ ./config/
COPY scripts/ ./scripts/

# Create necessary directories
RUN mkdir -p data/delta/tables \
    data/delta/logs \
    data/quality/reports \
    data/reports \
    logs \
    certs

# Set environment variables
ENV CONFIG_PATH=/app/config/config.yaml

# Expose port
EXPOSE 8080

# Run the application
CMD ["./nessi"] 