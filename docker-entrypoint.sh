#!/bin/bash
set -e

# Install the package in development mode
pip install -e .

# Set up Spark environment
export SPARK_CONF_DIR=/app/conf
export SPARK_JARS_DIR=/app/jars
export PYTHONPATH=/app:$PYTHONPATH

# Create necessary directories
mkdir -p /root/.nessi
mkdir -p /app/data/test_tables
mkdir -p /app/data/reports
mkdir -p /tmp/spark-warehouse
mkdir -p $SPARK_JARS_DIR
mkdir -p $SPARK_CONF_DIR

# Set permissions
chmod -R 777 /root/.nessi
chmod -R 777 /app/data
chmod -R 777 /tmp/spark-warehouse
chmod -R 777 $SPARK_JARS_DIR
chmod -R 777 $SPARK_CONF_DIR

# Verify Java installation
if [ ! -d "$JAVA_HOME" ]; then
    echo "Java installation not found at $JAVA_HOME"
    exit 1
fi

# Verify Spark installation
if [ ! -d "$SPARK_HOME" ]; then
    echo "Spark installation not found at $SPARK_HOME"
    exit 1
fi

# Initialize Delta Lake
echo "spark.sql.extensions=io.delta.sql.DeltaSparkSessionExtension" > $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.sql.catalog.spark_catalog=org.apache.spark.sql.delta.catalog.DeltaCatalog" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.delta.logStore.class=org.apache.spark.sql.delta.storage.S3SingleDriverLogStore" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.sql.warehouse.dir=/tmp/spark-warehouse" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.driver.memory=2g" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.executor.memory=2g" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.driver.maxResultSize=2g" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.sql.shuffle.partitions=4" >> $SPARK_CONF_DIR/spark-defaults.conf
echo "spark.default.parallelism=4" >> $SPARK_CONF_DIR/spark-defaults.conf

# Print environment information
echo "JAVA_HOME: $JAVA_HOME"
echo "SPARK_HOME: $SPARK_HOME"
echo "PYTHONPATH: $PYTHONPATH"
echo "SPARK_JARS_DIR: $SPARK_JARS_DIR"
echo "SPARK_CONF_DIR: $SPARK_CONF_DIR"

# Execute the command
exec "$@" 