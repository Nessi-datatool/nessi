import pytest
import os
import shutil
from pathlib import Path
from unittest.mock import patch, MagicMock
from src.scanner.table_scanner import TableScanner
from src.scanner.sample_data_generator import SampleDataGenerator
from src.scanner.spark_config import get_spark_session

@pytest.fixture(scope="session")
def spark():
    spark = get_spark_session("TestSession")
    yield spark
    spark.stop()

@pytest.fixture
def test_data_dir():
    test_dir = Path("data/test")
    test_dir.mkdir(parents=True, exist_ok=True)
    yield test_dir
    if test_dir.exists():
        shutil.rmtree(test_dir)

@pytest.fixture
def data_generator():
    generator = SampleDataGenerator()
    yield generator
    if hasattr(generator, 'spark'):
        generator.spark.stop()

@pytest.fixture
def table_scanner():
    scanner = TableScanner()
    yield scanner
    if hasattr(scanner, 'spark'):
        scanner.spark.stop()

@pytest.fixture
def sample_data(test_data_dir, data_generator):
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
    assert df.count() == 5  # Check number of rows
    assert len(df.columns) == 5  # Check number of columns
    
    # Verify column names
    expected_columns = {'id', 'name', 'age', 'email', 'created_at'}
    assert set(df.columns) == expected_columns

def test_parquet_scanning(table_scanner, sample_data):
    """Test scanning Parquet files"""
    # Get the actual parquet file path (Spark creates a directory)
    parquet_file = str(sample_data["parquet_path"])
    
    # Scan the parquet table
    result = table_scanner.scan_table(parquet_file, "parquet")
    
    # Verify the scan results
    assert 'schema' in result
    assert 'row_count' in result
    assert 'column_stats' in result
    
    assert result['row_count'] == 5
    assert len(result['schema']) == 5
    
    # Check column statistics
    stats = result['column_stats']
    assert 'age' in stats
    assert 'count' in stats['age']
    assert 'mean' in stats['age']

def test_csv_scanning(table_scanner, sample_data):
    """Test scanning CSV files"""
    # Get the CSV file path
    csv_file = str(sample_data["csv_path"])
    
    # Scan the CSV table
    result = table_scanner.scan_table(csv_file, "csv")
    
    # Verify the scan results
    assert 'schema' in result
    assert 'row_count' in result
    assert 'column_stats' in result
    
    assert result['row_count'] == 5
    assert len(result['schema']) == 5
    
    # Check column statistics
    stats = result['column_stats']
    assert 'age' in stats
    assert 'count' in stats['age']
    assert 'mean' in stats['age']

def test_invalid_path(table_scanner):
    """Test handling of invalid file paths"""
    with pytest.raises(ValueError):
        table_scanner.scan_table("/nonexistent/path", "parquet")

def test_invalid_format(table_scanner, sample_data):
    """Test handling of invalid format types"""
    with pytest.raises(ValueError):
        table_scanner.scan_table(str(sample_data["parquet_path"]), "invalid_format") 