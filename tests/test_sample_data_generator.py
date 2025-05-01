"""Tests for the sample data generator module."""

import os
import pytest
from datetime import date
from typing import List, Tuple
import pandas as pd
import numpy as np
from src.sample_data_generator import (
    generate_sample_data,
    save_as_delta,
    save_as_parquet,
    save_as_csv
)
from pyspark.sql import SparkSession

@pytest.fixture(scope="session")
def spark_session():
    """Create a Spark session for testing."""
    spark = SparkSession.builder \
        .appName("test_sample_data_generator") \
        .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
        .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
        .getOrCreate()
    yield spark
    spark.stop()

@pytest.fixture
def test_data_dir(tmp_path):
    """Create a temporary directory for test data."""
    return str(tmp_path)

def test_generate_sample_data():
    """Test sample data generation with various row counts."""
    # Test with default number of rows
    data = generate_sample_data()
    assert isinstance(data, list)
    assert len(data) == 1000
    assert all(isinstance(row, tuple) for row in data)
    
    # Test with custom number of rows
    data = generate_sample_data(num_rows=100)
    assert len(data) == 100
    
    # Test with minimum number of rows
    data = generate_sample_data(num_rows=1)
    assert len(data) == 1
    
    # Test with invalid number of rows
    with pytest.raises(ValueError):
        generate_sample_data(num_rows=0)
    with pytest.raises(ValueError):
        generate_sample_data(num_rows=-1)

def test_generate_sample_data_types():
    """Test data types in generated sample data."""
    data = generate_sample_data(num_rows=10)
    
    for row in data:
        # Check each column's type
        assert isinstance(row[0], int)  # id
        assert isinstance(row[1], str)  # name
        assert isinstance(row[2], (int, type(None)))  # age
        assert isinstance(row[3], (float, type(None)))  # salary
        assert isinstance(row[4], (str, type(None)))  # department
        assert isinstance(row[5], date)  # join_date

def test_generate_sample_data_null_values():
    """Test null value distribution in generated data."""
    data = generate_sample_data(num_rows=100)
    df = pd.DataFrame(data, columns=['id', 'name', 'age', 'salary', 'department', 'join_date'])
    
    # Check null value distribution
    assert df['age'].isna().sum() > 0  # Every 10th row should be null
    assert df['salary'].isna().sum() > 0  # Every 5th row should be null
    assert df['department'].isna().sum() > 0  # Every 7th row should be null

def test_save_as_delta(test_data_dir):
    """Test saving data as Delta table."""
    data = generate_sample_data(num_rows=10)
    table_path = os.path.join(test_data_dir, "test_delta")
    os.makedirs(test_data_dir, exist_ok=True)

    # Test saving list of tuples
    save_as_delta(data, table_path)
    assert os.path.exists(table_path)
    assert os.path.exists(os.path.join(table_path, "_delta_log"))

    # Test saving pandas DataFrame
    df = pd.DataFrame(data, columns=['id', 'name', 'age', 'salary', 'department', 'join_date'])
    save_as_delta(df, table_path)
    assert os.path.exists(table_path)
    assert os.path.exists(os.path.join(table_path, "_delta_log"))

    # Test with empty data
    with pytest.raises(ValueError, match="Data cannot be empty"):
        save_as_delta([], table_path)
    with pytest.raises(ValueError, match="DataFrame cannot be empty"):
        save_as_delta(pd.DataFrame(), table_path)

    # Test with invalid path
    with pytest.raises(ValueError, match="Table path cannot be empty"):
        save_as_delta(data, "")

def test_save_as_parquet(test_data_dir, spark_session):
    """Test saving data as Parquet file."""
    data = generate_sample_data(num_rows=10)
    parquet_path = os.path.join(test_data_dir, "test.parquet")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Test saving list of tuples
    save_as_parquet(data, parquet_path)
    assert os.path.exists(parquet_path)
    
    # Test saving pandas DataFrame
    df = pd.DataFrame(data, columns=['id', 'name', 'age', 'salary', 'department', 'join_date'])
    save_as_parquet(df, parquet_path)
    assert os.path.exists(parquet_path)
    
    # Test with empty data
    with pytest.raises(ValueError, match="Data cannot be empty"):
        save_as_parquet([], parquet_path)
    with pytest.raises(ValueError, match="DataFrame cannot be empty"):
        save_as_parquet(pd.DataFrame(), parquet_path)
    
    # Test with invalid path
    with pytest.raises(ValueError, match="Parquet path cannot be empty"):
        save_as_parquet(data, "")

def test_save_as_csv(test_data_dir):
    """Test saving data as CSV file."""
    # Generate sample data
    data = generate_sample_data(10)
    df = pd.DataFrame(data, columns=["id", "name", "age", "salary", "department", "join_date"])
    
    # Save as CSV
    csv_path = os.path.join(test_data_dir, "test.csv")
    save_as_csv(df, csv_path)
    
    # Verify CSV file was created
    assert os.path.exists(csv_path)
    
    # Verify content
    df_read = pd.read_csv(csv_path)
    assert len(df_read) == 10

def assert_dataframes_equal(df1, df2):
    """Custom function to compare DataFrames with proper null value handling."""
    # Convert None to NaN for consistent comparison
    df1 = df1.replace({None: np.nan})
    df2 = df2.replace({None: np.nan})
    pd.testing.assert_frame_equal(df1, df2, check_dtype=False)

def test_data_consistency_across_formats(test_data_dir, spark_session):
    """Test data consistency across different file formats."""
    # Generate sample data
    data = generate_sample_data(num_rows=10)
    df = pd.DataFrame(data, columns=['id', 'name', 'age', 'salary', 'department', 'join_date'])
    
    # Convert data types consistently
    df['id'] = df['id'].astype('int64')
    df['age'] = df['age'].fillna(-1).astype('int64')  # Fill NA with -1 for integers
    df['salary'] = df['salary'].astype('float64')
    df['join_date'] = pd.to_datetime(df['join_date']).dt.date
    
    # Save in different formats
    delta_path = os.path.join(test_data_dir, "test_consistency.delta")
    parquet_path = os.path.join(test_data_dir, "test_consistency.parquet")
    csv_path = os.path.join(test_data_dir, "test_consistency.csv")
    
    save_as_delta(df, delta_path)
    save_as_parquet(df, parquet_path)
    save_as_csv(df, csv_path)
    
    # Read back and verify
    df_delta = spark_session.read.format("delta").load(delta_path).toPandas()
    df_parquet = pd.read_parquet(parquet_path)
    df_csv = pd.read_csv(csv_path)
    
    # Convert data types consistently for comparison
    df_delta['id'] = df_delta['id'].astype('int64')
    df_delta['age'] = df_delta['age'].fillna(-1).astype('int64')
    df_delta['salary'] = df_delta['salary'].astype('float64')
    df_delta['join_date'] = pd.to_datetime(df_delta['join_date']).dt.date
    
    df_parquet['id'] = df_parquet['id'].astype('int64')
    df_parquet['age'] = df_parquet['age'].fillna(-1).astype('int64')
    df_parquet['salary'] = df_parquet['salary'].astype('float64')
    df_parquet['join_date'] = pd.to_datetime(df_parquet['join_date']).dt.date
    
    df_csv['id'] = df_csv['id'].astype('int64')
    df_csv['age'] = df_csv['age'].fillna(-1).astype('int64')
    df_csv['salary'] = df_csv['salary'].astype('float64')
    df_csv['join_date'] = pd.to_datetime(df_csv['join_date']).dt.date
    
    # Compare data using custom comparison function
    assert_dataframes_equal(df_delta.sort_values('id').reset_index(drop=True),
                          df_parquet.sort_values('id').reset_index(drop=True))
    assert_dataframes_equal(df_delta.sort_values('id').reset_index(drop=True),
                          df_csv.sort_values('id').reset_index(drop=True)) 