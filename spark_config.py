from pyspark.sql import SparkSession
from pyspark import SparkConf
import os

def get_spark_session(app_name: str = "Nessi") -> SparkSession:
    """Create and configure a Spark session.
    
    Args:
        app_name (str): Name of the Spark application.
        
    Returns:
        SparkSession: Configured Spark session.
    """
    # Get the directory where this script is located
    current_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Set up Spark configuration
    conf = SparkConf() \
        .setAppName(app_name) \
        .set("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
        .set("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
        .set("spark.delta.logStore.class", "org.apache.spark.sql.delta.storage.S3SingleDriverLogStore") \
        .set("spark.sql.warehouse.dir", "/tmp/spark-warehouse") \
        .set("spark.driver.memory", "2g") \
        .set("spark.executor.memory", "2g") \
        .set("spark.driver.maxResultSize", "2g") \
        .set("spark.sql.shuffle.partitions", "4") \
        .set("spark.default.parallelism", "4") \
        .set("spark.sql.legacy.timeParserPolicy", "LEGACY") \
        .set("spark.sql.legacy.parquet.int96RebaseModeInRead", "LEGACY") \
        .set("spark.sql.legacy.parquet.int96RebaseModeInWrite", "LEGACY") \
        .set("spark.sql.legacy.parquet.datetimeRebaseModeInRead", "LEGACY") \
        .set("spark.sql.legacy.parquet.datetimeRebaseModeInWrite", "LEGACY") \
        .set("spark.sql.legacy.json.allowEmptyString.enabled", "true") \
        .set("spark.sql.legacy.createHiveTableByDefault.enabled", "false") \
        .set("spark.sql.legacy.allowNonEmptyLocationInCTAS", "true") \
        .set("spark.sql.legacy.allowCreatingManagedTableUsingNonemptyLocation", "true") \
        .set("spark.sql.legacy.parquet.datetimeRebaseModeInRead", "CORRECTED") \
        .set("spark.sql.legacy.parquet.datetimeRebaseModeInWrite", "CORRECTED") \
        .set("spark.sql.legacy.parquet.int96RebaseModeInRead", "CORRECTED") \
        .set("spark.sql.legacy.parquet.int96RebaseModeInWrite", "CORRECTED") \
        .set("spark.sql.legacy.timeParserPolicy", "CORRECTED") \
        .set("spark.sql.legacy.json.allowEmptyString.enabled", "false") \
        .set("spark.sql.legacy.createHiveTableByDefault.enabled", "true") \
        .set("spark.sql.legacy.allowNonEmptyLocationInCTAS", "false") \
        .set("spark.sql.legacy.allowCreatingManagedTableUsingNonemptyLocation", "false")

    # Create Spark session
    spark = SparkSession.builder \
        .config(conf=conf) \
        .enableHiveSupport() \
        .getOrCreate()

    # Set log level to WARN to reduce noise
    spark.sparkContext.setLogLevel("WARN")

    return spark 