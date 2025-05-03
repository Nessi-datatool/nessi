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
    
    def get_session(self, app_name="NessiApp"):
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

def get_spark_session(app_name="NessiApp"):
    """Get the Spark session from the manager."""
    return spark_manager.get_session(app_name)

def stop_spark_session():
    """Stop the Spark session."""
    spark_manager.stop_session()