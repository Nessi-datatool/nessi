import pytest
import os
import tempfile
from pathlib import Path
from src.scanner.sample_data_generator import SampleDataGenerator

@pytest.fixture
def temp_dir():
    with tempfile.TemporaryDirectory() as tmpdirname:
        yield tmpdirname

def test_generate_sample_data():
    generator = SampleDataGenerator()
    df = generator.generate_sample_data()
    
    # Check if DataFrame is not empty
    assert df.count() > 0
    
    # Check if all required columns are present
    required_columns = ['id', 'name', 'age', 'email', 'created_at']
    assert all(col in df.columns for col in required_columns)
    
    # Check data types
    assert df.schema['id'].dataType.typeName() == 'integer'
    assert df.schema['name'].dataType.typeName() == 'string'
    assert df.schema['age'].dataType.typeName() == 'integer'
    assert df.schema['email'].dataType.typeName() == 'string'
    assert df.schema['created_at'].dataType.typeName() == 'timestamp'

def test_save_as_parquet(temp_dir):
    generator = SampleDataGenerator()
    df = generator.generate_sample_data()
    
    # Create a subdirectory for the parquet file
    parquet_dir = os.path.join(temp_dir, "test_parquet")
    os.makedirs(parquet_dir, exist_ok=True)
    
    # Save as parquet
    generator.save_as_parquet(df, parquet_dir)
    
    # Verify the parquet file was created
    assert os.path.exists(parquet_dir)
    assert len(os.listdir(parquet_dir)) > 0  # Should contain parquet files

def test_save_as_csv(temp_dir):
    generator = SampleDataGenerator()
    df = generator.generate_sample_data()
    
    # Create a subdirectory for the CSV file
    csv_dir = os.path.join(temp_dir, "test_csv")
    os.makedirs(csv_dir, exist_ok=True)
    
    # Save as CSV
    generator.save_as_csv(df, csv_dir)
    
    # Verify the CSV file was created
    assert os.path.exists(csv_dir)
    assert len(os.listdir(csv_dir)) > 0  # Should contain CSV files 