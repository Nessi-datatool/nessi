from pyspark.sql import SparkSession
import os

def get_spark_session(app_name="NessiApp"):
    """
    Get or create a Spark session with the specified configuration.
    
    Args:
        app_name (str): Name of the Spark application
        
    Returns:
        SparkSession: Configured Spark session
    """
    return SparkSession.builder \
        .appName(app_name) \
        .config("spark.driver.host", "0.0.0.0") \
        .config("spark.driver.bindAddress", "0.0.0.0") \
        .config("spark.executor.memory", "1g") \
        .config("spark.driver.memory", "1g") \
        .config("spark.sql.shuffle.partitions", "1") \
        .config("spark.default.parallelism", "1") \
        .config("spark.dynamicAllocation.enabled", "false") \
        .config("spark.shuffle.service.enabled", "false") \
        .config("spark.ui.enabled", "false") \
        .config("spark.sql.warehouse.dir", "/tmp/spark-warehouse") \
        .config("spark.local.dir", "/tmp/spark-local") \
        .config("spark.sql.execution.arrow.pyspark.enabled", "true") \
        .config("spark.sql.execution.arrow.pyspark.fallback.enabled", "true") \
        .master("local[1]") \
        .getOrCreate() 