"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi. All rights reserved.

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
   The Software is licensed, not sold. Nessi retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

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