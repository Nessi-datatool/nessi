from pyspark.sql import SparkSession

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
        .config("spark.driver.host", "localhost") \
        .config("spark.driver.bindAddress", "127.0.0.1") \
        .config("spark.executor.memory", "1g") \
        .config("spark.driver.memory", "1g") \
        .config("spark.sql.shuffle.partitions", "1") \
        .config("spark.default.parallelism", "1") \
        .config("spark.dynamicAllocation.enabled", "false") \
        .config("spark.shuffle.service.enabled", "false") \
        .config("spark.ui.enabled", "false") \
        .master("local[1]") \
        .getOrCreate() 