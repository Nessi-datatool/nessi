"""Spark session configuration and management."""

from pyspark.sql import SparkSession
import os
import logging
from pyspark import SparkContext, SparkConf

logger = logging.getLogger(__name__)

class SparkSessionManager:
    _instance = None
    _spark = None
    
    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(SparkSessionManager, cls).__new__(cls)
        return cls._instance
    
    def get_session(self, app_name="nessi.dev"):
        """
        Get or create a Spark session with the specified configuration.
        
        Args:
            app_name (str): Name of the Spark application
            
        Returns:
            SparkSession: Configured Spark session
        """
        if self._spark is None:
            try:
                # Create SparkSession with configuration
                self._spark = SparkSession.builder \
                    .appName(app_name) \
                    .master("local[*]") \
                    .config("spark.driver.host", "localhost") \
                    .config("spark.driver.bindAddress", "localhost") \
                    .config("spark.executor.memory", "1g") \
                    .config("spark.driver.memory", "1g") \
                    .config("spark.sql.shuffle.partitions", "2") \
                    .config("spark.default.parallelism", "2") \
                    .config("spark.dynamicAllocation.enabled", "false") \
                    .config("spark.shuffle.service.enabled", "false") \
                    .config("spark.ui.enabled", "false") \
                    .config("spark.sql.warehouse.dir", "/tmp/spark-warehouse") \
                    .config("spark.local.dir", "/tmp/spark-local") \
                    .config("spark.sql.execution.arrow.pyspark.enabled", "true") \
                    .config("spark.sql.execution.arrow.pyspark.fallback.enabled", "true") \
                    .config("spark.jars.packages", "io.delta:delta-core_2.12:2.4.0") \
                    .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
                    .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
                    .getOrCreate()
                    
                logger.info(f"Created new Spark session for {app_name}")
            except Exception as e:
                logger.error(f"Failed to create Spark session: {str(e)}")
                raise
        return self._spark
    
    def stop_session(self):
        """Stop the Spark session if it exists."""
        if self._spark is not None:
            try:
                # Get the current SparkContext
                sc = self._spark.sparkContext
                # Stop the SparkContext
                sc.stop()
                # Clear the session
                self._spark = None
                logger.info("Stopped Spark session")
            except Exception as e:
                logger.error(f"Failed to stop Spark session: {str(e)}")
                raise

# Global instance
spark_manager = SparkSessionManager()

def get_spark_session(app_name="nessi.dev"):
    """Get the Spark session from the manager."""
    return spark_manager.get_session(app_name)

def stop_spark_session():
    """Stop the Spark session."""
    spark_manager.stop_session() 