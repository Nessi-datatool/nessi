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