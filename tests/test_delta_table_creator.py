"""Tests for the Delta table creator module."""

import os
import shutil
import tempfile
import pytest
from datetime import date
from pyspark.sql import SparkSession
from pyspark.sql.types import StructType, StructField, StringType, IntegerType

from backend.src.scanner.create_test_delta_table import (
    create_test_delta_table,
    add_column_to_delta_table
)

@pytest.fixture(scope="function")
def test_data_dir():
    """Create a temporary directory for test data."""
    temp_dir = tempfile.mkdtemp()
    yield temp_dir
    if os.path.exists(temp_dir):
        shutil.rmtree(temp_dir)

def test_create_test_delta_table(test_data_dir, spark_session):
    """Test creating a test Delta table."""
    # Create test table
    table_path = os.path.join(test_data_dir, "test_delta")
    os.makedirs(test_data_dir, exist_ok=True)
    create_test_delta_table(table_path)
    
    # Verify table was created
    assert os.path.exists(table_path)
    
    # Read table and verify schema
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() == 4
    assert "id" in df.columns
    assert "name" in df.columns
    assert "age" in df.columns
    assert "salary" in df.columns
    assert "department" in df.columns
    assert "join_date" in df.columns

def test_create_test_delta_table_with_schema(test_data_dir, spark_session):
    """Test creating a test Delta table with custom schema."""
    # Define custom schema
    custom_schema = StructType([
        StructField("id", IntegerType(), False),
        StructField("name", StringType(), True),
        StructField("department", StringType(), True)
    ])
    
    # Create table with custom schema
    table_path = os.path.join(test_data_dir, "test_delta_schema")
    os.makedirs(test_data_dir, exist_ok=True)
    create_test_delta_table(table_path, schema=custom_schema)
    assert os.path.exists(table_path)
    
    # Read table and verify schema
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() == 4
    assert "id" in df.columns
    assert "name" in df.columns
    assert "department" in df.columns
    assert "age" not in df.columns
    assert "salary" not in df.columns
    assert "join_date" not in df.columns

def test_create_test_delta_table_with_partitioning(test_data_dir, spark_session):
    """Test creating a test Delta table with partitioning."""
    # Create partitioned table
    table_path = os.path.join(test_data_dir, "test_delta_partitioned")
    os.makedirs(test_data_dir, exist_ok=True)
    create_test_delta_table(table_path, partition_by=["department"])
    assert os.path.exists(table_path)
    
    # Read table and verify partitioning
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() == 4
    assert "_delta_log" in os.listdir(table_path)
    assert "department=Department_1" in os.listdir(table_path)

def test_create_test_delta_table_with_num_rows(test_data_dir, spark_session):
    """Test creating a test Delta table with custom number of rows."""
    # Create table with custom number of rows
    table_path = os.path.join(test_data_dir, "test_delta_rows")
    os.makedirs(test_data_dir, exist_ok=True)
    create_test_delta_table(table_path, num_rows=100)
    assert os.path.exists(table_path)
    
    # Read table and verify row count
    df = spark_session.read.format("delta").load(table_path)
    assert df.count() == 100

def test_create_test_delta_table_invalid_args(test_data_dir):
    """Test creating a test Delta table with invalid arguments."""
    # Test with invalid path
    with pytest.raises(ValueError):
        create_test_delta_table("")
    
    # Test with invalid number of rows
    with pytest.raises(ValueError):
        create_test_delta_table(test_data_dir, num_rows=0)
    with pytest.raises(ValueError):
        create_test_delta_table(test_data_dir, num_rows=-1)

def test_add_column_to_delta_table(test_data_dir, spark_session):
    """Test adding a column to a Delta table."""
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
    
    # Add string column with default value
    add_column_to_delta_table(
        table_path,
        column_name="status",
        data_type="string",
        default_value="active"
    )
    
    # Read table and verify columns
    df = spark_session.read.format("delta").load(table_path)
    assert "is_active" in df.columns
    assert "status" in df.columns
    assert df.select("is_active").first()[0] == True
    assert df.select("status").first()[0] == "active"

def test_add_column_to_delta_table_invalid_args(test_data_dir, spark_session):
    """Test adding a column to a Delta table with invalid arguments."""
    # Create initial table
    table_path = os.path.join(test_data_dir, "test_delta_invalid")
    os.makedirs(test_data_dir, exist_ok=True)
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
            data_type="integer"
        )
    
    # Test with invalid data type
    with pytest.raises(ValueError):
        add_column_to_delta_table(
            table_path,
            column_name="test",
            data_type="invalid_type"
        )

def test_add_column_to_delta_table_null_values(test_data_dir, spark_session):
    """Test adding a column to a Delta table with null values."""
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
    
    # Read table and verify null values
    df = spark_session.read.format("delta").load(table_path)
    assert "notes" in df.columns
    assert df.select("notes").first()[0] is None 