"""Tests for the table scanner module."""

import os
import pytest
import tempfile
import shutil
from datetime import date
import pandas as pd
from pyspark.sql import SparkSession
from backend.src.scanner.table_scanner import TableScanner
from backend.src.scanner.sample_data_generator import SampleDataGenerator
from pathlib import Path
import mocker

@pytest.fixture(scope="function")
def test_data_dir():
    """Create a temporary directory for test data."""
    temp_dir = tempfile.mkdtemp()
    yield temp_dir
    shutil.rmtree(temp_dir, ignore_errors=True)

@pytest.fixture
def sample_parquet_table(test_data_dir):
    """Create a sample Parquet table for testing."""
    table_path = os.path.join(test_data_dir, "sample_parquet")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create sample data
    data = generate_sample_data(num_rows=100)
    save_as_parquet(data, table_path)
    
    return table_path

@pytest.fixture
def sample_csv_table(test_data_dir):
    """Create a sample CSV table for testing."""
    table_path = os.path.join(test_data_dir, "sample.csv")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create sample data
    data = pd.DataFrame({
        'id': [1, 2, 3],
        'name': ['Alice', 'Bob', 'Charlie'],
        'age': [25, 30, 35],
        'salary': [50000, 60000, 70000],
        'department': ['HR', 'Engineering', 'Sales'],
        'join_date': [date(2020, 1, 1), date(2020, 2, 1), date(2020, 3, 1)]
    })
    save_as_csv(data, table_path)
    
    return table_path

def test_table_scanner_initialization():
    """Test TableScanner initialization."""
    scanner = TableScanner()
    assert scanner is not None
    assert hasattr(scanner, 'spark')

def test_scan_parquet_table(sample_parquet_table):
    """Test scanning a Parquet table."""
    scanner = TableScanner()
    result = scanner.scan_table(sample_parquet_table)
    
    # Check result structure
    assert isinstance(result, dict)
    assert 'schema' in result
    assert 'row_count' in result
    assert 'column_stats' in result
    assert 'file_info' in result
    
    # Check schema
    assert len(result['schema']) > 0
    for field in result['schema']:
        assert 'name' in field
        assert 'type' in field
    
    # Check row count
    assert result['row_count'] == 100
    
    # Check column stats
    assert len(result['column_stats']) > 0
    for col, stats in result['column_stats'].items():
        assert 'count' in stats
        assert 'null_count' in stats
        assert 'min' in stats
        assert 'max' in stats
        if col in ['age', 'salary']:
            assert 'mean' in stats
            assert 'stddev' in stats
    
    # Check file info
    assert result['file_info']['format'] == 'parquet'
    assert result['file_info']['path'] == sample_parquet_table
    assert result['file_info']['size'] > 0

def test_scan_csv_table(sample_csv_table):
    """Test scanning a CSV table."""
    scanner = TableScanner()
    result = scanner.scan_table(sample_csv_table)
    
    # Check result structure
    assert isinstance(result, dict)
    assert 'schema' in result
    assert 'row_count' in result
    assert 'column_stats' in result
    assert 'file_info' in result
    
    # Check file info
    assert result['file_info']['format'] == 'csv'
    assert result['file_info']['path'] == sample_csv_table
    assert result['file_info']['size'] > 0
    
    # Check row count
    assert result['row_count'] == 3
    
    # Check schema
    assert len(result['schema']) == 6
    assert any(field['name'] == 'id' for field in result['schema'])
    assert any(field['name'] == 'name' for field in result['schema'])
    
    # Check column stats
    assert len(result['column_stats']) == 6
    assert 'id' in result['column_stats']
    assert 'name' in result['column_stats']
    assert result['column_stats']['id']['count'] == 3
    assert result['column_stats']['name']['count'] == 3

def test_scan_invalid_table(test_data_dir):
    """Test scanning an invalid table path."""
    scanner = TableScanner()
    invalid_path = os.path.join(test_data_dir, "nonexistent")
    
    with pytest.raises(ValueError):
        scanner.scan_table(invalid_path)

def test_scan_table_with_large_data(test_data_dir):
    """Test scanning a table with large amount of data."""
    table_path = os.path.join(test_data_dir, "large_parquet")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create large sample data
    data = generate_sample_data(num_rows=10000)
    save_as_parquet(data, table_path)
    
    # Scan the table
    scanner = TableScanner()
    result = scanner.scan_table(table_path)
    
    # Verify results
    assert result['row_count'] == 10000
    assert result['file_info']['format'] == 'parquet'
    assert result['file_info']['size'] > 0

def test_scan_table_with_null_values(test_data_dir):
    """Test scanning a table with null values."""
    table_path = os.path.join(test_data_dir, "null_parquet")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create data with null values
    data = pd.DataFrame({
        'id': [1, 2, 3, 4, 5],
        'name': ['Alice', None, 'Charlie', None, 'Eve'],
        'age': [25, 30, None, 40, None],
        'salary': [50000, None, 70000, None, 90000],
        'department': ['HR', 'Engineering', None, 'Sales', None],
        'join_date': [date(2020, 1, 1), None, date(2020, 3, 1), None, date(2020, 5, 1)]
    })
    save_as_parquet(data, table_path)
    
    # Scan the table
    scanner = TableScanner()
    result = scanner.scan_table(table_path)
    
    # Verify null value handling
    assert result['row_count'] == 5
    assert result['column_stats']['name']['null_count'] == 2
    assert result['column_stats']['age']['null_count'] == 2
    assert result['column_stats']['salary']['null_count'] == 2
    assert result['column_stats']['department']['null_count'] == 2
    assert result['column_stats']['join_date']['null_count'] == 2

def test_generate_report(test_data_dir):
    """Test generating a report from scan results."""
    # Generate and save test data
    table_path = os.path.join(test_data_dir, "report_parquet")
    os.makedirs(test_data_dir, exist_ok=True)
    
    data = generate_sample_data(num_rows=100)
    save_as_parquet(data, table_path)
    
    # Scan the table
    scanner = TableScanner()
    result = scanner.scan_table(table_path)
    
    # Generate report
    report_path = os.path.join(test_data_dir, "report.md")
    report = scanner.generate_report(result, report_path)
    
    # Verify report was generated
    assert os.path.exists(report_path)
    assert os.path.getsize(report_path) > 0
    
    # Check report content
    with open(report_path, 'r') as f:
        content = f.read()
        assert "Table Scan Report" in content
        assert "Schema" in content
        assert "Statistics" in content
        assert "File Information" in content

@pytest.fixture
def temp_dir():
    with tempfile.TemporaryDirectory() as tmpdirname:
        yield tmpdirname

@pytest.fixture
def sample_data(temp_dir):
    generator = SampleDataGenerator()
    df = generator.generate_sample_data()
    
    # Save sample data in both formats
    parquet_path = os.path.join(temp_dir, "sample_parquet")
    csv_path = os.path.join(temp_dir, "sample_csv")
    
    generator.save_as_parquet(df, parquet_path)
    generator.save_as_csv(df, csv_path)
    
    return {
        "parquet_path": parquet_path,
        "csv_path": csv_path,
        "dataframe": df
    }

def test_scan_parquet_table(sample_data):
    scanner = TableScanner()
    result = scanner.scan_table(sample_data["parquet_path"], "parquet")
    
    assert result is not None
    assert "schema" in result
    assert "row_count" in result
    assert "column_stats" in result
    assert result["row_count"] == 5  # We know we generated 5 rows

def test_scan_csv_table(sample_data):
    scanner = TableScanner()
    result = scanner.scan_table(sample_data["csv_path"], "csv")
    
    assert result is not None
    assert "schema" in result
    assert "row_count" in result
    assert "column_stats" in result
    assert result["row_count"] == 5  # We know we generated 5 rows

def test_scan_csv_special_chars(temp_dir):
    """Test scanning CSV with special characters in column names"""
    scanner = TableScanner()
    
    # Create a test CSV file with special characters in column names
    csv_path = os.path.join(temp_dir, "test.csv")
    with open(csv_path, "w") as f:
        f.write("id,user@domain.com,test.field\n")
        f.write("1,value1,value2\n")
    
    result = scanner.scan_table(csv_path, "csv")
    
    # Check if special characters were properly handled
    assert "user_at_domain_dot_com" in result["column_stats"]
    assert "test_dot_field" in result["column_stats"]

def test_scan_invalid_format(sample_data):
    scanner = TableScanner()
    with pytest.raises(ValueError):
        scanner.scan_table(sample_data["parquet_path"], "invalid_format")

def test_scan_nonexistent_path():
    scanner = TableScanner()
    with pytest.raises(ValueError):
        scanner.scan_table("/nonexistent/path", "parquet")

def test_spark_session_cleanup(mocker):
    """Test that Spark session is properly cleaned up"""
    # Create a mock Spark session
    mock_spark = mocker.Mock()
    mock_stop = mocker.Mock()
    mock_spark.stop = mock_stop
    
    # Create scanner and replace its Spark session with our mock
    scanner = TableScanner()
    scanner.spark = mock_spark
    
    # Test cleanup when spark attribute exists
    scanner.__del__()
    mock_stop.assert_called_once()
    
    # Test cleanup when spark attribute doesn't exist
    scanner = TableScanner()
    delattr(scanner, 'spark')  # Remove spark attribute
    scanner.__del__()  # Should not raise any error
    
    # Test with a new scanner
    scanner = TableScanner()
    assert scanner.spark is not None
    scanner.__del__() 