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

"""Pytest configuration file."""

import os
import tempfile
import pytest
import pandas as pd
import numpy as np
from pathlib import Path
from pyspark.sql import SparkSession
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, LongType, DoubleType, DateType, TimestampType

from backend.src.scanner.spark_config import get_spark_session

def assert_dataframes_equal(df1, df2):
    """Custom function to compare DataFrames with proper null value handling."""
    # Convert None to NaN for consistent comparison
    df1 = df1.replace({None: np.nan})
    df2 = df2.replace({None: np.nan})
    pd.testing.assert_frame_equal(df1, df2, check_dtype=False)

@pytest.fixture(scope="session", autouse=True)
def setup_test_environment():
    """Set up test environment and create test data"""
    # Create test data directory
    test_data_dir = Path("data/test_tables")
    test_data_dir.mkdir(parents=True, exist_ok=True)
    
    # Create sample data
    sample_data = pd.DataFrame({
        'id': range(1, 101),
        'name': [f'Name_{i}' if i % 10 != 0 else None for i in range(1, 101)],  # Add some nulls
        'value': np.random.rand(100),
        'category': np.random.choice(['A', 'B', 'C', None], 100),  # Add some nulls
        'date': pd.date_range('2023-01-01', periods=100).astype(str)  # Convert to string to avoid timestamp issues
    })
    
    # Create sample CSV table
    csv_path = test_data_dir / 'sample_csv_table'
    csv_path.mkdir(exist_ok=True)
    sample_data.to_csv(csv_path / 'data.csv', index=False)
    
    # Create sample Parquet table
    parquet_path = test_data_dir / 'sample_parquet_table'
    parquet_path.mkdir(exist_ok=True)
    sample_data.to_parquet(parquet_path / 'data.parquet', engine='pyarrow')
    
    # Create empty table
    empty_path = test_data_dir / 'empty_table'
    empty_path.mkdir(exist_ok=True)
    pd.DataFrame(columns=['id', 'value']).to_csv(empty_path / 'data.csv', index=False)
    
    # Create table with special characters
    special_chars_data = pd.DataFrame({
        'Column Name': range(1, 11),
        'Column$Value': np.random.rand(10),
        'Column-With-Hyphen': [f'Value_{i}' if i % 2 != 0 else None for i in range(1, 11)]  # Add some nulls
    })
    special_chars_path = test_data_dir / 'special_chars_table'
    special_chars_path.mkdir(exist_ok=True)
    special_chars_data.to_csv(special_chars_path / 'data.csv', index=False)
    
    # Create table with mixed types
    mixed_types_data = pd.DataFrame({
        'int_col': range(1, 11),
        'float_col': np.random.rand(10),
        'string_col': [f'String_{i}' if i % 3 != 0 else None for i in range(1, 11)],  # Add some nulls
        'bool_col': np.random.choice([True, False, None], 10),  # Add some nulls
        'date_col': pd.date_range('2023-01-01', periods=10).astype(str)  # Convert to string to avoid timestamp issues
    })
    mixed_types_path = test_data_dir / 'mixed_types_table'
    mixed_types_path.mkdir(exist_ok=True)
    mixed_types_data.to_csv(mixed_types_path / 'data.csv', index=False)
    
    # Create large table
    large_data = pd.DataFrame({
        'id': range(1, 1000001),
        'value': np.random.rand(1000000),
        'category': np.random.choice(['A', 'B', 'C', 'D', 'E', None], 1000000)  # Add some nulls
    })
    large_path = test_data_dir / 'large_table'
    large_path.mkdir(exist_ok=True)
    large_data.to_csv(large_path / 'data.csv', index=False)
    
    # Create compressed table
    compressed_path = test_data_dir / 'compressed_table'
    compressed_path.mkdir(exist_ok=True)
    sample_data.to_csv(compressed_path / 'data.csv.gz', index=False, compression='gzip')
    
    yield
    
    # Cleanup (optional)
    # shutil.rmtree(test_data_dir) 

@pytest.fixture(scope="session")
def spark_session():
    """Create a Spark session for testing"""
    spark = get_spark_session("TestSession")
    yield spark
    spark.stop()

@pytest.fixture(scope="function")
def test_data_dir():
    """Get the test data directory"""
    return Path("data/test_tables")

@pytest.fixture(scope="function")
def sample_schema():
    """Get a sample schema for testing"""
    return StructType([
        StructField("id", IntegerType(), True),
        StructField("name", StringType(), True),
        StructField("value", DoubleType(), True),
        StructField("category", StringType(), True),
        StructField("date", StringType(), True)
    ])

@pytest.fixture
def custom_schema():
    """Get a custom schema for testing"""
    return StructType([
        StructField("id", LongType(), True),
        StructField("name", StringType(), True),
        StructField("value", DoubleType(), True),
        StructField("category", StringType(), True),
        StructField("date", DateType(), True)
    ])