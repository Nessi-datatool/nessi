#!/bin/bash

# Build the test image
docker build -t nessi-test -f Dockerfile.test .

# Run tests in Docker
docker run --rm \
    -v $(pwd):/app \
    -w /app \
    nessi-test \
    pytest tests/ $@ 