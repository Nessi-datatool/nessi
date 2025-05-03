FROM python:3.11-slim-bookworm

# Install system dependencies
RUN apt-get update && apt-get install -y \
    openjdk-17-jdk \
    procps \
    wget \
    gnupg \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set JAVA_HOME based on architecture
RUN if [ "$(uname -m)" = "aarch64" ]; then \
        echo "export JAVA_HOME=/usr/lib/jvm/java-17-openjdk-arm64" >> /etc/profile.d/java_home.sh; \
    else \
        echo "export JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64" >> /etc/profile.d/java_home.sh; \
    fi

# Source the JAVA_HOME and set PATH
ENV JAVA_HOME_SOURCE=/etc/profile.d/java_home.sh
RUN . $JAVA_HOME_SOURCE && \
    echo "export PATH=$JAVA_HOME/bin:$PATH" >> $JAVA_HOME_SOURCE

# Set working directory
WORKDIR /app

# Create a non-root user
RUN groupadd -r nessi && useradd -r -g nessi nessi

# Copy requirements first to leverage Docker cache
COPY backend/requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the backend application
COPY backend/ /app/backend/

# Set environment variables
ENV PYTHONPATH=/app/backend
ENV PYSPARK_PYTHON=/usr/local/bin/python
ENV PYSPARK_DRIVER_PYTHON=/usr/local/bin/python
ENV PYSPARK_SUBMIT_ARGS="--driver-java-options=-Xms1024M --driver-java-options=-Xmx2048M --driver-java-options=-Dlog4j.logLevel=info pyspark-shell"
ENV GRAFANA_URL=http://grafana:3000
ENV PROMETHEUS_URL=http://prometheus:9090

# Create necessary directories with proper permissions
RUN mkdir -p /app/data/parquet /app/data/csv /app/data/delta /app/metrics && \
    chown -R nessi:nessi /app && \
    chmod -R 755 /app/data /app/metrics

# Switch to non-root user
USER nessi

# Source JAVA_HOME at runtime
SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["/bin/bash", "-c", "source $JAVA_HOME_SOURCE && exec \"$@\""]
CMD ["python", "backend/src/main.py"] 