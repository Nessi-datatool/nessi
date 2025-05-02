FROM python:3.9-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    default-jdk \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables
ENV JAVA_HOME=/usr/lib/jvm/default-java
ENV PYTHONPATH=/app
ENV PYSPARK_PYTHON=/usr/local/bin/python
ENV PYSPARK_DRIVER_PYTHON=/usr/local/bin/python

# Set working directory
WORKDIR /app

# Create necessary directories for Spark
RUN mkdir -p /tmp/spark-warehouse /tmp/spark-local

# Copy requirements first to leverage Docker cache
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir prometheus-client==0.19.0 pyyaml==6.0.1 requests==2.31.0 jinja2==3.1.2 pyarrow==14.0.1 pandas==2.1.4 numpy==1.26.3

# Copy application code
COPY src/ src/
COPY tests/ tests/

# Create directory for SSL certificates
RUN mkdir -p /app/ssl

# Expose port for metrics
EXPOSE 8000

# Command to run tests
CMD ["pytest", "tests/", "--cov=src", "--cov-report=term-missing"] 