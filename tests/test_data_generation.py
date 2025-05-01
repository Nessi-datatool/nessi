"""Tests for the data generation module."""

import os
import pytest
import pandas as pd
from datetime import date
from pyspark.sql import SparkSession
from src.sample_data_generator import generate_sample_data, save_as_delta, save_as_parquet, save_as_csv
from src.create_test_delta_table import create_test_delta_table, get_or_create_spark_session
from src.scanner.table_scanner import TableScanner
import numpy as np

@pytest.fixture(scope="session")
def spark_session():
    """Create a Spark session configured for Delta Lake."""
    spark = get_or_create_spark_session()
    yield spark
    spark.stop()

@pytest.fixture
def test_data_dir(tmp_path):
    """Create a temporary directory for test data."""
    os.makedirs(tmp_path, exist_ok=True)
    return str(tmp_path)

def assert_dataframes_equal(df1, df2):
    """Custom function to compare DataFrames with proper null value handling."""
    # Convert None to NaN for consistent comparison
    df1 = df1.replace({None: np.nan})
    df2 = df2.replace({None: np.nan})
    pd.testing.assert_frame_equal(df1, df2, check_dtype=False)

def test_generate_sample_data():
    """Test sample data generation with various row counts."""
    # Test with default number of rows
    data = generate_sample_data()
    assert len(data) == 1000
    
    # Test with custom number of rows
    data = generate_sample_data(100)
    assert len(data) == 100
    
    # Test data types and structure
    row = data[0]
    assert isinstance(row[0], int)  # id
    assert isinstance(row[1], str)  # name
    assert isinstance(row[2], (int, type(None)))  # age
    assert isinstance(row[3], (float, type(None)))  # salary
    assert isinstance(row[4], (str, type(None)))  # department
    assert isinstance(row[5], date)  # join_date

def test_save_as_delta(test_data_dir, spark_session):
    """Test saving data as Delta table."""
    data = generate_sample_data(10)
    delta_path = os.path.join(test_data_dir, "test_delta")
    save_as_delta(data, delta_path)
    
    # Verify Delta table was created
    assert os.path.exists(delta_path)
    assert os.path.exists(os.path.join(delta_path, "_delta_log"))
    
    # Verify data can be read back
    df = spark_session.read.format("delta").load(delta_path)
    assert df.count() == 10

def test_save_as_parquet(test_data_dir, spark_session):
    """Test saving data as Parquet file."""
    data = generate_sample_data(10)
    df = pd.DataFrame(data)
    parquet_path = os.path.join(test_data_dir, "test_parquet")
    save_as_parquet(df, parquet_path)
    
    # Verify Parquet file was created
    assert os.path.exists(parquet_path)
    
    # Verify data can be read back
    df = spark_session.read.format("parquet").load(parquet_path)
    assert df.count() == 10

def test_save_as_csv(test_data_dir):
    """Test saving data as CSV file."""
    data = generate_sample_data(10)
    df = pd.DataFrame(data)
    csv_path = os.path.join(test_data_dir, "test.csv")
    save_as_csv(df, csv_path)
    
    # Verify CSV file was created
    assert os.path.exists(csv_path)
    
    # Verify content
    df_read = pd.read_csv(csv_path)
    assert len(df_read) == 10

def test_create_test_delta_table(test_data_dir, spark_session):
    """Test creation of test Delta table with multiple versions."""
    # Create test table path
    table_path = os.path.join(test_data_dir, "test_delta")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create test table
    create_test_delta_table(table_path=table_path, num_rows=4)
    
    # Verify table was created
    assert os.path.exists(table_path)
    assert os.path.exists(os.path.join(table_path, "_delta_log"))
    
    # Verify version history
    df = spark_session.read.format("delta").load(table_path)
    
    # Check final state
    assert df.count() == 4  # Should have 4 records in final version
    assert df.filter("id = 3").count() == 0  # Record with id=3 should be deleted
    assert df.filter("id = 1").select("department").first()[0] == "Engineering"  # Department should be updated 