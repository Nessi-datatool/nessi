from pyspark.sql.types import StructType, StructField, IntegerType, StringType, TimestampType
from pyspark.sql.functions import current_timestamp
from .spark_config import get_spark_session
import os
from pathlib import Path

class SampleDataGenerator:
    def __init__(self):
        self.spark = get_spark_session("SampleDataGenerator")

    def generate_sample_data(self):
        # Define schema
        schema = StructType([
            StructField("id", IntegerType(), False),
            StructField("name", StringType(), True),
            StructField("age", IntegerType(), True),
            StructField("email", StringType(), True),
            StructField("created_at", TimestampType(), True)
        ])

        # Generate sample data
        data = [
            (1, "John Doe", 30, "john@example.com", None),
            (2, "Jane Smith", 25, "jane@example.com", None),
            (3, "Bob Johnson", 35, "bob@example.com", None),
            (4, "Alice Brown", 28, "alice@example.com", None),
            (5, "Charlie Wilson", 32, "charlie@example.com", None)
        ]

        # Create DataFrame
        df = self.spark.createDataFrame(data, schema)
        
        # Add current timestamp to created_at column
        df = df.withColumn("created_at", current_timestamp())
        
        return df

    def save_as_parquet(self, df, output_path):
        """Save DataFrame as Parquet format"""
        # Ensure the directory exists
        os.makedirs(output_path, exist_ok=True)
        
        # Save as Parquet
        df.write.mode("overwrite").parquet(output_path)

    def save_as_csv(self, df, output_path):
        """Save DataFrame as CSV format"""
        # Ensure the directory exists
        os.makedirs(output_path, exist_ok=True)
        
        # Save as CSV with specific options for proper handling
        df.coalesce(1).write \
            .mode("overwrite") \
            .option("header", "true") \
            .option("escape", '"') \
            .option("quote", '"') \
            .option("quoteAll", "true") \
            .csv(output_path)

    def __del__(self):
        """Cleanup Spark session"""
        if hasattr(self, 'spark'):
            self.spark.stop() 