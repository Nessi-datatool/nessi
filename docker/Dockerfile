FROM python:3.13-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    default-jdk \
    wget \
    libssl-dev \
    libffi-dev \
    python3-dev \
    libsnappy-dev \
    libzstd-dev \
    liblz4-dev \
    libbrotli-dev \
    && rm -rf /var/lib/apt/lists/*

# Download and install Spark
RUN wget -q https://archive.apache.org/dist/spark/spark-3.5.0/spark-3.5.0-bin-hadoop3.tgz \
    && tar xzf spark-3.5.0-bin-hadoop3.tgz \
    && mv spark-3.5.0-bin-hadoop3 /opt/spark \
    && rm spark-3.5.0-bin-hadoop3.tgz

# Download Delta Lake JARs
RUN wget -P /opt/spark/jars https://repo1.maven.org/maven2/io/delta/delta-core_2.12/3.0.0/delta-core_2.12-3.0.0.jar \
    && wget -P /opt/spark/jars https://repo1.maven.org/maven2/io/delta/delta-storage/3.0.0/delta-storage-3.0.0.jar

# Set environment variables
ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64
ENV SPARK_HOME=/opt/spark
ENV PYTHONPATH=/app/src
ENV PATH=$SPARK_HOME/bin:$PATH
ENV PYTHONUNBUFFERED=1

# Set working directory
WORKDIR /app

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY src /app/src
COPY config /app/config

# Create necessary directories
RUN mkdir -p /app/data /app/logs /app/metrics /app/reports

# Expose ports
EXPOSE 8443

# Start the application
CMD ["python", "-m", "nessi.main"] 