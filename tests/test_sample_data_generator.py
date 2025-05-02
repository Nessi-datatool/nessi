"""Tests for the sample data generator module."""

import os
import pytest
import tempfile
import shutil
import pandas as pd
import numpy as np
from datetime import date
from backend.src.scanner.sample_data_generator import generate_sample_data, save_as_parquet, save_as_csv

@pytest.fixture
def temp_dir():
    """Create a temporary directory for test data."""
    temp_dir = tempfile.mkdtemp()
    yield temp_dir
    shutil.rmtree(temp_dir)

def test_generate_sample_data():
    """Test generating sample data."""
    data = generate_sample_data(100)
    
    # Check data structure
    assert len(data) == 100
    assert all(col in data.columns for col in ['id', 'name', 'age', 'salary', 'department', 'join_date'])
    
    # Check data types
    assert data['id'].dtype == np.int64
    assert data['name'].dtype == object
    assert data['age'].dtype == float  # float because of NaN values
    assert data['salary'].dtype == float
    assert data['department'].dtype == object
    assert isinstance(data['join_date'].iloc[0], date)
    
    # Check value ranges
    valid_ages = data['age'].dropna()
    assert valid_ages.between(18, 65).all()
    
    valid_salaries = data['salary'].dropna()
    assert valid_salaries.min() >= 0
    
    # Check null values
    assert data['age'].isnull().sum() > 0  # Some null values in age
    assert data['salary'].isnull().sum() > 0  # Some null values in salary
    assert data['department'].isnull().sum() > 0  # Some null values in department

def test_save_empty_data(temp_dir):
    """Test saving empty data."""
    with pytest.raises(ValueError, match="Number of rows must be at least 1"):
        data = generate_sample_data(0)

def test_save_invalid_path():
    """Test saving to invalid path."""
    data = generate_sample_data(100)
    with pytest.raises(ValueError, match="Invalid output path"):
        save_as_parquet(data, "")

def test_save_as_parquet(temp_dir):
    """Test saving data as Parquet."""
    data = generate_sample_data(100)
    output_path = os.path.join(temp_dir, "test.parquet")
    save_as_parquet(data, output_path)
    assert os.path.exists(output_path)
    
    # Read back and verify
    df = pd.read_parquet(output_path)
    assert len(df) == 100
    assert all(col in df.columns for col in ['id', 'name', 'age', 'salary', 'department', 'join_date'])

def test_save_as_csv(temp_dir):
    """Test saving data as CSV."""
    data = generate_sample_data(100)
    output_path = os.path.join(temp_dir, "test.csv")
    save_as_csv(data, output_path)
    assert os.path.exists(output_path)
    
    # Read back and verify
    df = pd.read_csv(output_path)
    assert len(df) == 100
    assert all(col in df.columns for col in ['id', 'name', 'age', 'salary', 'department', 'join_date'])

def test_save_large_data(temp_dir):
    """Test saving large dataset."""
    data = generate_sample_data(10000)
    output_path = os.path.join(temp_dir, "large.parquet")
    save_as_parquet(data, output_path)
    assert os.path.exists(output_path)
    
    # Read back and verify
    df = pd.read_parquet(output_path)
    assert len(df) == 10000

def test_save_data_with_nulls(temp_dir):
    """Test saving data with null values."""
    data = generate_sample_data(100)
    output_path = os.path.join(temp_dir, "nulls.parquet")
    save_as_parquet(data, output_path)
    assert os.path.exists(output_path)
    
    # Read back and verify
    df = pd.read_parquet(output_path)
    assert df['age'].isnull().sum() > 0
    assert df['salary'].isnull().sum() > 0
    assert df['department'].isnull().sum() > 0 