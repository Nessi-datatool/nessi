"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

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
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

import pytest
import os
from backend.src.scanner.spark_config import get_spark_session

def test_spark_session_creation():
    """Test basic Spark session creation"""
    spark = get_spark_session("TestApp")
    assert spark is not None
    assert spark.conf.get("spark.app.name") == "TestApp"
    spark.stop()

def test_delta_lake_configuration():
    """Test Delta Lake related configurations"""
    spark = get_spark_session()
    
    # Check Delta Lake extensions
    assert spark.conf.get("spark.sql.extensions") == "io.delta.sql.DeltaSparkSessionExtension"
    assert spark.conf.get("spark.sql.catalog.spark_catalog") == "org.apache.spark.sql.delta.catalog.DeltaCatalog"
    
    # Check Delta Lake jars
    jars = spark.conf.get("spark.jars")
    assert "delta-core" in jars
    assert "delta-storage" in jars
    
    spark.stop()

def test_memory_settings():
    """Test memory configurations"""
    spark = get_spark_session()
    
    assert spark.conf.get("spark.driver.memory") == "2g"
    assert spark.conf.get("spark.executor.memory") == "2g"
    
    spark.stop()

def test_legacy_settings():
    """Test legacy mode settings"""
    spark = get_spark_session()
    
    assert spark.conf.get("spark.sql.legacy.timeParserPolicy") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.int96RebaseModeInRead") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.int96RebaseModeInWrite") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.datetimeRebaseModeInRead") == "LEGACY"
    assert spark.conf.get("spark.sql.legacy.parquet.datetimeRebaseModeInWrite") == "LEGACY"
    
    spark.stop()

def test_log_level():
    """Test log level configuration"""
    spark = get_spark_session()
    assert spark.sparkContext.getLogLevel() == "WARN"
    spark.stop()

def test_custom_jars_directory(monkeypatch):
    """Test custom jars directory configuration"""
    custom_jars_dir = "/custom/jars/path"
    monkeypatch.setenv("SPARK_JARS_DIR", custom_jars_dir)
    
    spark = get_spark_session()
    jars = spark.conf.get("spark.jars")
    
    assert custom_jars_dir in jars
    spark.stop()