from pyspark.sql import SparkSession
import os
from pathlib import Path
import logging

logger = logging.getLogger(__name__)

class TableScanner:
    def __init__(self, spark: SparkSession):
        """Initialize TableScanner with a Spark session.
        
        Args:
            spark: SparkSession instance.
        """
        self.spark = spark
        logger.info("Initialized TableScanner with Spark session")

    def scan_table(self, table_path: str, format_type: str) -> dict:
        """
        Scan a table and return its metadata and statistics
        
        Args:
            table_path (str): Path to the table
            format_type (str): Format of the table (parquet, csv, etc.)
            
        Returns:
            dict: Dictionary containing table metadata and statistics
        """
        if not table_path or not os.path.exists(table_path):
            raise ValueError(f"Table path does not exist: {table_path}")

        # Read the table based on format
        if format_type.lower() == "parquet":
            df = self.spark.read.parquet(table_path)
        elif format_type.lower() == "csv":
            # For CSV, we need to handle special characters in column names
            df = self.spark.read.csv(table_path, header=True, inferSchema=True)
            # Rename columns to remove special characters
            for col in df.columns:
                new_col = col.replace('@', '_at_').replace('.', '_dot_')
                if new_col != col:
                    df = df.withColumnRenamed(col, new_col)
        else:
            raise ValueError(f"Unsupported format type: {format_type}")

        # Get schema information
        schema = [{"name": field.name, "type": str(field.dataType)} 
                 for field in df.schema.fields]

        # Get row count
        row_count = df.count()

        # Get column statistics
        column_stats = {}
        for col in df.columns:
            stats = df.select(col).describe().collect()
            
            # Convert stats to dictionary
            stats_dict = {}
            for row in stats:
                if len(row) >= 2:  # Ensure we have both summary and value
                    stats_dict[row[0]] = row[1]
            
            column_stats[col] = {
                "count": stats_dict.get("count", None),
                "mean": stats_dict.get("mean", None),
                "stddev": stats_dict.get("stddev", None),
                "min": stats_dict.get("min", None),
                "max": stats_dict.get("max", None)
            }

        return {
            "schema": schema,
            "row_count": row_count,
            "column_stats": column_stats
        } 