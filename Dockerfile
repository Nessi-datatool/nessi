# Build stage
FROM golang:1.21-alpine AS builder

# Build arguments for version information
ARG VERSION=dev
ARG BUILD_DATE=unknown
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with version information and multiple architecture support
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-s -w -X github.com/nessi-dev/nessi/pkg/version.Version=${VERSION} \
    -X github.com/nessi-dev/nessi/pkg/version.BuildDate=${BUILD_DATE} \
    -X github.com/nessi-dev/nessi/pkg/version.GitCommit=${GIT_COMMIT} \
    -X github.com/nessi-dev/nessi/pkg/version.GitBranch=${GIT_BRANCH}" \
    -o bin/nessi cmd/nessi/main.go

# Final stage - using distroless for minimal attack surface
FROM gcr.io/distroless/static:nonroot

# Add labels for better container management
LABEL org.opencontainers.image.title="Nessi"
LABEL org.opencontainers.image.description="Delta Lake management and data quality tool"
LABEL org.opencontainers.image.source="https://github.com/nessi-dev/nessi"
LABEL org.opencontainers.image.vendor="Nessi Project"
LABEL org.opencontainers.image.licenses="Apache-2.0"

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/nessi .

# Copy configuration and example files
COPY --from=builder /app/config/ ./config/
COPY --from=builder /app/examples/ ./examples/

# Create necessary directories
RUN mkdir -p data/delta/tables \
    data/delta/logs \
    data/quality/reports \
    data/reports \
    logs \
    certs

# Set environment variables
ENV CONFIG_PATH=/app/config/config.yaml
ENV TZ=UTC

# Expose port
EXPOSE 8080

# Use non-root user for security
USER nonroot:nonroot

# Run the application
ENTRYPOINT ["/app/nessi"]
CMD ["--help"]