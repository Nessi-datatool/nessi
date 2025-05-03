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
import shutil
from pathlib import Path
from unittest.mock import patch, MagicMock
from src.scanner.scanner import TableScanner
from src.scanner.sample_data_generator import SampleDataGenerator
from pyspark.sql import SparkSession
from pyspark import SparkContext

@pytest.fixture(scope="session")
def spark():
    """Create a SparkSession for testing."""
    # Stop any existing SparkContext
    try:
        sc = SparkContext.getOrCreate()
        sc.stop()
    except:
        pass
    
    # Create new SparkSession
    spark = SparkSession.builder \
        .appName("NessiTest") \
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
        .getOrCreate()
    print(f"[DEBUG] Created SparkSession in fixture: id={id(spark)}, isStopped={getattr(spark.sparkContext, 'isStopped', None)}")
    
    yield spark
    
    # Cleanup
    spark.stop()
    try:
        sc = SparkContext.getOrCreate()
        sc.stop()
    except:
        pass

@pytest.fixture
def test_data_dir():
    """Create a temporary directory for test data."""
    test_dir = Path("data/test")
    test_dir.mkdir(parents=True, exist_ok=True)
    yield test_dir
    if test_dir.exists():
        shutil.rmtree(test_dir)

@pytest.fixture
def data_generator(spark):
    """Create a SampleDataGenerator instance with the test Spark session."""
    return SampleDataGenerator(spark)

@pytest.fixture
def table_scanner(spark):
    """Create a TableScanner instance with the test Spark session."""
    return TableScanner(spark)

@pytest.fixture
def sample_data(test_data_dir, data_generator):
    """Generate and save sample data for testing."""
    parquet_path = test_data_dir / "test_parquet"
    csv_path = test_data_dir / "test_csv"
    
    # Generate and save sample data
    df = data_generator.generate_sample_data()
    data_generator.save_as_parquet(df, str(parquet_path))
    data_generator.save_as_csv(df, str(csv_path))
    
    return {
        "parquet_path": parquet_path,
        "csv_path": csv_path,
        "dataframe": df
    }

def test_sample_data_generation(data_generator):
    """Test sample data generation"""
    df = data_generator.generate_sample_data()
    assert df is not None
    assert df.count() > 0
    assert "id" in df.columns
    assert "name" in df.columns
    assert "age" in df.columns
    assert "email" in df.columns
    assert "created_at" in df.columns

def test_parquet_scanning(data_generator, table_scanner, test_data_dir):
    """Test scanning Parquet files"""
    # Generate and save sample data
    df = data_generator.generate_sample_data()
    parquet_path = test_data_dir / "test_parquet"
    data_generator.save_as_parquet(df, str(parquet_path))
    
    # Scan the saved data
    results = table_scanner.scan_table(str(parquet_path), "parquet")
    
    assert results is not None
    assert "row_count" in results
    assert "schema" in results
    assert "column_stats" in results
    assert results["row_count"] == df.count()

def test_csv_scanning(data_generator, table_scanner, test_data_dir):
    """Test scanning CSV files"""
    # Generate and save sample data
    df = data_generator.generate_sample_data()
    csv_path = test_data_dir / "test_csv"
    data_generator.save_as_csv(df, str(csv_path))
    
    # Scan the saved data
    results = table_scanner.scan_table(str(csv_path), "csv")
    
    assert results is not None
    assert "row_count" in results
    assert "schema" in results
    assert "column_stats" in results
    assert results["row_count"] == df.count()

def test_invalid_path(table_scanner):
    """Test scanning invalid path"""
    with pytest.raises(ValueError):
        table_scanner.scan_table("invalid_path", "parquet")

def test_invalid_format(table_scanner):
    """Test scanning with invalid format"""
    with pytest.raises(ValueError):
        table_scanner.scan_table("test_data.txt", "txt")