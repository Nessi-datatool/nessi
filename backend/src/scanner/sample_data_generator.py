"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi Data Tool. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. Nessi Data Tool retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Sample data generator for testing and demonstration."""

import logging
from typing import Dict, List, Optional
from pyspark.sql import SparkSession, DataFrame
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, TimestampType
from .spark_config import get_spark_session
import os
from pathlib import Path

logger = logging.getLogger(__name__)

class SampleDataGenerator:
    """Generate sample data for testing and demonstration."""
    
    def __init__(self, spark: Optional[SparkSession] = None):
        """Initialize the sample data generator.
        
        Args:
            spark: Optional SparkSession instance.
        """
        self.spark = spark or get_spark_session("SampleDataGenerator")
        logger.info("Initialized SampleDataGenerator with Spark session")
        
    def generate_sample_data(self, num_rows: int = 5) -> DataFrame:
        """Generate sample data with predefined schema.
        
        Args:
            num_rows: Number of rows to generate.
            
        Returns:
            DataFrame with sample data.
        """
        try:
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
            return df
            
        except Exception as e:
            logger.error(f"Failed to generate sample data: {str(e)}")
            raise
            
    def save_as_parquet(self, df: DataFrame, path: str) -> None:
        """Save DataFrame as Parquet file.
        
        Args:
            df: DataFrame to save.
            path: Output path.
        """
        try:
            df.write.parquet(path, mode="overwrite")
        except Exception as e:
            logger.error(f"Failed to save Parquet file: {str(e)}")
            raise
            
    def save_as_csv(self, df: DataFrame, path: str) -> None:
        """Save DataFrame as CSV file.
        
        Args:
            df: DataFrame to save.
            path: Output path.
        """
        try:
            df.write.csv(path, mode="overwrite", header=True)
        except Exception as e:
            logger.error(f"Failed to save CSV file: {str(e)}")
            raise