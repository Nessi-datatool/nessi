"""Spark configuration module.

This module provides functionality to create and configure Spark sessions
with Delta Lake support.
"""

import os
import logging
from pyspark.sql import SparkSession
from pyspark import SparkConf

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def get_spark_session(app_name: str = "Nessi") -> SparkSession:
    """Get a configured Spark session.
    
    Args:
        app_name (str): Name of the Spark application.
        
    Returns:
        SparkSession: Configured Spark session.
    """
    # Set up Spark configuration
    conf = SparkConf()
    
    # Delta Lake configuration
    conf.set("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension")
    conf.set("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog")
    
    # Add Delta Lake jars to classpath
    delta_core_jar = os.path.join(os.environ.get("SPARK_JARS_DIR", "/app/jars"), "delta-core_2.12-2.4.0.jar")
    delta_storage_jar = os.path.join(os.environ.get("SPARK_JARS_DIR", "/app/jars"), "delta-storage-2.4.0.jar")
    conf.set("spark.jars", f"{delta_core_jar},{delta_storage_jar}")
    
    # Memory settings
    conf.set("spark.driver.memory", "2g")
    conf.set("spark.executor.memory", "2g")
    
    # Other settings
    conf.set("spark.sql.legacy.timeParserPolicy", "LEGACY")
    conf.set("spark.sql.legacy.parquet.int96RebaseModeInRead", "LEGACY")
    conf.set("spark.sql.legacy.parquet.int96RebaseModeInWrite", "LEGACY")
    conf.set("spark.sql.legacy.parquet.datetimeRebaseModeInRead", "LEGACY")
    conf.set("spark.sql.legacy.parquet.datetimeRebaseModeInWrite", "LEGACY")
    
    # Create Spark session
    spark = SparkSession.builder \
        .appName(app_name) \
        .config(conf=conf) \
        .enableHiveSupport() \
        .getOrCreate()
    
    # Set log level to WARN to reduce noise
    spark.sparkContext.setLogLevel("WARN")
    
    return spark 