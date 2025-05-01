"""Tests for the test delta table creator module."""

import os
import pytest
from datetime import date
import pandas as pd
from pyspark.sql.types import StructType, StructField, StringType, IntegerType, LongType, DoubleType, DateType, TimestampType
from pyspark.sql import SparkSession
from pyspark.sql.functions import col
from src.create_test_delta_table import (
    create_test_delta_table,
    add_column_to_delta_table,
    get_or_create_spark_session
)

@pytest.fixture(scope="session")
def spark_session():
    """Create a Spark session configured for Delta Lake."""
    spark = get_or_create_spark_session()
    yield spark
    spark.stop()

@pytest.fixture
def test_data_dir(tmp_path):
    """Create a temporary directory for test data."""
    return str(tmp_path)

@pytest.fixture
def custom_schema():
    """Custom schema for testing."""
    return StructType([
        StructField("id", LongType(), False),
        StructField("name", StringType(), True),
        StructField("value", DoubleType(), True),
        StructField("timestamp", TimestampType(), True)
    ])

def get_or_create_spark_session():
    """Get or create a Spark session."""
    return SparkSession.builder.appName("test_delta_table_creator").getOrCreate()

def test_create_test_delta_table_basic(test_data_dir, spark_session):
    """Test basic Delta table creation."""
    # Create test table path
    table_path = os.path.join(test_data_dir, "test_delta_basic")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create test table
    create_test_delta_table(table_path)
    
    # Verify table was created
    assert os.path.exists(table_path)
    assert os.path.exists(os.path.join(table_path, "_delta_log"))
    
    # Read and verify data using Spark
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() > 0
    assert set(df.columns) == {'id', 'name', 'age', 'salary', 'department', 'join_date'}

def test_create_test_delta_table_custom_schema(test_data_dir, spark_session, custom_schema):
    """Test Delta table creation with custom schema."""
    table_path = os.path.join(test_data_dir, "test_delta_custom")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create table with custom schema
    create_test_delta_table(table_path, schema=custom_schema)
    assert os.path.exists(table_path)
    
    # Verify schema
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() > 0
    assert set(df.columns) == {'id', 'name', 'value', 'timestamp'}

def test_create_test_delta_table_partitioned(test_data_dir, spark_session):
    """Test Delta table creation with partitioning."""
    table_path = os.path.join(test_data_dir, "test_delta_partitioned")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create partitioned table
    create_test_delta_table(
        table_path,
        partition_by=["department"]
    )
    assert os.path.exists(table_path)
    
    # Verify partitioning
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() > 0
    assert "department" in df.columns

def test_create_test_delta_table_custom_rows(test_data_dir, spark_session):
    """Test Delta table creation with custom number of rows."""
    table_path = os.path.join(test_data_dir, "test_delta_custom_rows")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create table with custom number of rows
    create_test_delta_table(table_path, num_rows=100)
    assert os.path.exists(table_path)
    
    # Verify row count
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() == 100

def test_create_test_delta_table_invalid_inputs(test_data_dir):
    """Test Delta table creation with invalid inputs."""
    # Test with invalid path
    with pytest.raises(ValueError):
        create_test_delta_table("")
    
    # Test with invalid number of rows
    with pytest.raises(ValueError):
        create_test_delta_table(test_data_dir, num_rows=0)
    with pytest.raises(ValueError):
        create_test_delta_table(test_data_dir, num_rows=-1)

def test_add_column_to_delta_table(test_data_dir, spark_session):
    """Test adding columns to Delta table."""
    # Create initial table
    table_path = os.path.join(test_data_dir, "test_delta_add_column")
    os.makedirs(test_data_dir, exist_ok=True)
    create_test_delta_table(table_path)
    
    # Add boolean column
    add_column_to_delta_table(
        table_path,
        column_name="is_active",
        data_type="boolean",
        default_value=True
    )
    
    # Verify new column
    df = spark_session.read.format("delta").load(table_path)
    assert "is_active" in df.columns
    assert df.filter(col("is_active") == True).count() == df.count()  # All values should be True
    
    # Add string column with default value
    add_column_to_delta_table(
        table_path,
        column_name="status",
        data_type="string",
        default_value="active"
    )
    
    # Verify new column
    df = spark_session.read.format("delta").load(table_path)
    assert "status" in df.columns
    assert df.filter(col("status") == "active").count() == df.count()  # All values should be "active"

def test_add_column_to_delta_table_invalid_inputs(test_data_dir, spark_session):
    """Test adding columns with invalid inputs."""
    # Create initial table
    table_path = os.path.join(test_data_dir, "test_delta_invalid")
    create_test_delta_table(table_path)
    
    # Test with invalid path
    with pytest.raises(ValueError):
        add_column_to_delta_table(
            "",
            column_name="test",
            data_type="string"
        )
    
    # Test with existing column
    with pytest.raises(ValueError):
        add_column_to_delta_table(
            table_path,
            column_name="id",  # Already exists
            data_type="string"
        )
    
    # Test with invalid data type
    with pytest.raises(ValueError):
        add_column_to_delta_table(
            table_path,
            column_name="test",
            data_type="invalid_type"
        )

def test_add_column_to_delta_table_null_values(test_data_dir, spark_session):
    """Test adding columns with null values."""
    # Create initial table
    table_path = os.path.join(test_data_dir, "test_delta_null")
    os.makedirs(test_data_dir, exist_ok=True)
    create_test_delta_table(table_path)
    
    # Add column without default value
    add_column_to_delta_table(
        table_path,
        column_name="notes",
        data_type="string"
    )
    
    # Verify new column
    df = spark_session.read.format("delta").load(table_path)
    assert "notes" in df.columns
    assert df.filter(col("notes").isNull()).count() == df.count()  # All values should be null 