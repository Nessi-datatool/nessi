"""Tests for the table scanner module."""

import os
import pytest
from datetime import date
import pandas as pd
from pyspark.sql import SparkSession
from src.scanner.table_scanner import TableScanner
from src.create_test_delta_table import create_test_delta_table, get_or_create_spark_session, generate_sample_data

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
def sample_delta_table(test_data_dir, spark_session):
    """Create a sample Delta table for testing."""
    table_path = os.path.join(test_data_dir, "sample_delta")
    create_test_delta_table(table_path, num_rows=100)
    return table_path

@pytest.fixture
def sample_parquet_table(test_data_dir, spark_session):
    """Create a sample Parquet table for testing."""
    table_path = os.path.join(test_data_dir, "sample_parquet")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create sample data
    data = generate_sample_data(num_rows=10)
    df = spark_session.createDataFrame(data, ["id", "name", "age", "salary", "department", "join_date"])
    
    # Save as Parquet
    df.write.parquet(table_path, mode="overwrite")
    
    return table_path

@pytest.fixture
def sample_csv_table(test_data_dir, spark_session):
    """Create a sample CSV table for testing."""
    table_path = os.path.join(test_data_dir, "sample.csv")
    os.makedirs(test_data_dir, exist_ok=True)
    
    # Create sample data
    data = {
        'id': [1, 2, 3],
        'name': ['Alice', 'Bob', 'Charlie'],
        'age': [25, 30, 35],
        'salary': [50000, 60000, 70000],
        'department': ['HR', 'Engineering', 'Sales'],
        'join_date': [date(2020, 1, 1), date(2020, 2, 1), date(2020, 3, 1)]
    }
    df = pd.DataFrame(data)
    df.to_csv(table_path, index=False)
    
    return table_path

def test_table_scanner_initialization():
    """Test TableScanner initialization."""
    scanner = TableScanner()
    assert scanner is not None
    assert hasattr(scanner, 'spark')

def test_scan_delta_table(sample_delta_table):
    """Test scanning a Delta table."""
    scanner = TableScanner()
    result = scanner.scan_table(sample_delta_table)
    
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
    assert result['file_info']['format'] == 'delta'
    assert result['file_info']['path'] == sample_delta_table
    assert result['file_info']['size'] > 0

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
    
    # Check file info
    assert result['file_info']['format'] == 'parquet'
    assert result['file_info']['path'] == sample_parquet_table
    assert result['file_info']['size'] > 0
    
    # Check schema
    assert len(result['schema']) == 6
    for field in result['schema']:
        assert 'name' in field
        assert 'type' in field
        assert 'nullable' in field
    
    # Check row count
    assert result['row_count'] > 0
    
    # Check column stats
    assert len(result['column_stats']) == 6
    for col, stats in result['column_stats'].items():
        assert 'count' in stats
        assert 'null_count' in stats
        assert 'min' in stats
        assert 'max' in stats
        if col in ['age', 'salary']:
            assert 'mean' in stats
            assert 'stddev' in stats

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

def test_generate_report(sample_delta_table, test_data_dir):
    """Test report generation."""
    scanner = TableScanner()
    scan_result = scanner.scan_table(sample_delta_table)
    
    # Test saving report to file
    report_path = os.path.join(test_data_dir, "report.md")
    report = scanner.generate_report(scan_result, report_path)
    assert os.path.exists(report_path)
    
    # Test returning report as string
    report = scanner.generate_report(scan_result)
    assert isinstance(report, str)
    assert len(report) > 0
    
    # Test with invalid inputs
    with pytest.raises(ValueError):
        scanner.generate_report({})

def test_scan_invalid_table(test_data_dir):
    """Test scanning invalid table."""
    scanner = TableScanner()
    
    # Test with non-existent path
    with pytest.raises(ValueError):
        scanner.scan_table(os.path.join(test_data_dir, "non_existent"))
    
    # Test with empty path
    with pytest.raises(ValueError):
        scanner.scan_table("")

def test_scan_table_with_large_data(test_data_dir, spark_session):
    """Test scanning a table with large amount of data."""
    # Generate and save large dataset
    table_path = os.path.join(test_data_dir, "large_delta")
    create_test_delta_table(table_path, num_rows=10000)
    
    # Scan the table
    scanner = TableScanner()
    result = scanner.scan_table(table_path)
    
    # Verify results
    assert result['row_count'] == 10000
    assert result['file_info']['size'] > 0

def test_scan_table_with_null_values(test_data_dir, spark_session):
    """Test scanning a table with null values."""
    # Generate data with null values
    table_path = os.path.join(test_data_dir, "null_delta")
    create_test_delta_table(table_path, num_rows=100)
    
    # Scan the table
    scanner = TableScanner()
    result = scanner.scan_table(table_path)
    
    # Verify null value statistics
    for col, stats in result['column_stats'].items():
        if col in ['age', 'salary', 'department']:
            assert stats['null_count'] > 0  # Some values should be null 