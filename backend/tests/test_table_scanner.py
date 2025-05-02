import pytest
import os
import tempfile
from pathlib import Path
from src.scanner.table_scanner import TableScanner
from src.scanner.sample_data_generator import SampleDataGenerator

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

def test_scan_invalid_format(sample_data):
    scanner = TableScanner()
    with pytest.raises(ValueError):
        scanner.scan_table(sample_data["parquet_path"], "invalid_format")

def test_scan_nonexistent_path():
    scanner = TableScanner()
    with pytest.raises(ValueError):
        scanner.scan_table("/nonexistent/path", "parquet") 