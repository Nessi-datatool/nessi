import pytest
import os
import pandas as pd
import numpy as np
from pathlib import Path
from pyspark.sql import SparkSession

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
    test_data_dir = Path(__file__).parent / '..' / 'data' / 'test_tables'
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
    """Create a Spark session for testing."""
    spark = SparkSession.builder \
        .appName("test_session") \
        .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
        .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
        .getOrCreate()
    yield spark
    spark.stop()

@pytest.fixture
def test_data_dir(tmp_path):
    """Create a temporary directory for test data."""
    return str(tmp_path) 