"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

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
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

import os
import shutil
import tempfile
import pytest
from pyspark.sql import SparkSession
from delta.tables import DeltaTable
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, DoubleType, DateType
from src.scanner.create_test_delta_table import create_test_delta_table

def _make_temp_dir():
    temp_dir = tempfile.mkdtemp()
    yield temp_dir
    shutil.rmtree(temp_dir)

@pytest.fixture(scope="function")
def test_data_dir():
    yield from _make_temp_dir()

@pytest.fixture(scope="session")
def spark_session():
    return SparkSession.builder.master("local[*]").appName("NessiTest").getOrCreate()

def test_delta_table_time_travel(spark_session, test_data_dir):
    table_path = os.path.join(test_data_dir, "test_delta_time_travel")
    create_test_delta_table(table_path)
    # Append new data
    schema = StructType([
        StructField("id", IntegerType(), False),
        StructField("name", StringType(), True),
        StructField("age", IntegerType(), True),
        StructField("salary", DoubleType(), True),
        StructField("department", StringType(), True),
        StructField("join_date", StringType(), True)
    ])
    df_new = spark_session.createDataFrame(
        [(999, "New", 99, 9999.0, "Dept", "2024-01-01")],
        schema
    )
    df_new.write.format("delta").mode("append").save(table_path)
    # Check history
    delta_table = DeltaTable.forPath(spark_session, table_path)
    history = delta_table.history().collect()
    assert len(history) >= 2
    # Time travel to version 0
    old_df = spark_session.read.format("delta").option("versionAsOf", 0).load(table_path)
    assert old_df.count() == 4
    # Latest version should have 5 rows
    latest_df = spark_session.read.format("delta").load(table_path)
    assert latest_df.count() == 5
    # Verify data types
    assert old_df.schema["join_date"].dataType == StringType()
    assert latest_df.schema["join_date"].dataType == StringType()

def test_delta_table_schema_evolution(spark_session, test_data_dir):
    table_path = os.path.join(test_data_dir, "test_delta_schema_evolution")
    create_test_delta_table(table_path)
    # Add a new column
    schema = StructType([
        StructField("id", IntegerType(), False),
        StructField("name", StringType(), True),
        StructField("age", IntegerType(), True),
        StructField("salary", DoubleType(), True),
        StructField("department", StringType(), True),
        StructField("join_date", StringType(), True),
        StructField("extra_col", StringType(), True)
    ])
    df_new = spark_session.createDataFrame(
        [(1, "A", 10, 1000.0, "Dept", "2024-01-01", "extra")],
        schema
    )
    df_new.write.format("delta").option("mergeSchema", "true").mode("append").save(table_path)
    df = spark_session.read.format("delta").load(table_path)
    assert "extra_col" in df.columns
    assert df.filter(df.extra_col.isNotNull()).count() == 1

def test_delta_table_corrupted_log(spark_session, test_data_dir):
    table_path = os.path.join(test_data_dir, "test_delta_corrupt_log")
    create_test_delta_table(table_path)
    # Corrupt the _delta_log
    log_dir = os.path.join(table_path, "_delta_log")
    for f in os.listdir(log_dir):
        os.remove(os.path.join(log_dir, f))
    # Should raise error when reading
    with pytest.raises(Exception):
        spark_session.read.format("delta").load(table_path).count()

def test_delta_table_unsupported_version(spark_session, test_data_dir):
    table_path = os.path.join(test_data_dir, "test_delta_version")
    # Create a new Delta table with unsupported version
    log_dir = os.path.join(table_path, "_delta_log")
    os.makedirs(log_dir, exist_ok=True)
    
    # Create a new log file with unsupported version
    import json
    protocol = {
        "protocol": {
            "minReaderVersion": 9999,
            "minWriterVersion": 9999,
            "readerFeatures": ["unsupportedFeature1", "unsupportedFeature2"],
            "writerFeatures": ["unsupportedFeature1", "unsupportedFeature2"]
        }
    }
    metadata = {
        "metaData": {
            "id": "test-delta-table",
            "format": {"provider": "parquet", "options": {}},
            "schemaString": '{"type":"struct","fields":[{"name":"id","type":"integer","nullable":true,"metadata":{}},{"name":"name","type":"string","nullable":true,"metadata":{}},{"name":"age","type":"integer","nullable":true,"metadata":{}},{"name":"salary","type":"double","nullable":true,"metadata":{}},{"name":"department","type":"string","nullable":true,"metadata":{}},{"name":"join_date","type":"string","nullable":true,"metadata":{}}]}',
            "partitionColumns": [],
            "configuration": {},
            "createdTime": 1746280867021
        }
    }
    
    # Write the log file
    log_file = os.path.join(log_dir, "00000000000000000000.json")
    with open(log_file, "w") as f:
        json.dump(protocol, f)
        f.write("\n")
        json.dump(metadata, f)
        f.write("\n")
    
    # Try to read the table - should raise an exception
    with pytest.raises(Exception) as exc_info:
        spark_session.read.format("delta").load(table_path).count()
    assert "Delta protocol version is not supported by this version of Delta Lake" in str(exc_info.value)