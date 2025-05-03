"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""Tests for the data generation module."""

import os
import pytest
import tempfile
import shutil
from datetime import date
import pandas as pd
from pathlib import Path
from backend.src.scanner.sample_data_generator import SampleDataGenerator

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

def test_generate_sample_data():
    generator = SampleDataGenerator()
    df = generator.generate_sample_data()
    
    assert df is not None
    assert df.count() == 5  # We know we generated 5 rows
    assert len(df.columns) == 5  # id, name, age, email, created_at

def test_save_as_parquet(sample_data):
    assert os.path.exists(sample_data["parquet_path"])
    assert len(os.listdir(sample_data["parquet_path"])) > 0

def test_save_as_csv(sample_data):
    assert os.path.exists(sample_data["csv_path"])
    assert len(os.listdir(sample_data["csv_path"])) > 0

def test_save_empty_data(test_data_dir):
    """Test saving empty data."""
    # Generate empty data
    data = generate_sample_data(num_rows=0)
    
    # Save as Parquet
    table_path = os.path.join(test_data_dir, "empty_parquet")
    save_as_parquet(data, table_path)
    
    # Verify
    df = pd.read_parquet(table_path)
    assert len(df) == 0
    assert len(df.columns) == 6
    
    # Save as CSV
    file_path = os.path.join(test_data_dir, "empty.csv")
    save_as_csv(data, file_path)
    
    # Verify
    df = pd.read_csv(file_path)
    assert len(df) == 0
    assert len(df.columns) == 6

def test_save_invalid_path(test_data_dir):
    """Test saving to invalid path."""
    data = generate_sample_data(num_rows=10)
    
    # Try to save to invalid path
    with pytest.raises(RuntimeError):
        save_as_parquet(data, "/invalid/path/table")
    
    with pytest.raises(RuntimeError):
        save_as_csv(data, "/invalid/path/file.csv")

def test_save_large_data(test_data_dir):
    """Test saving large amount of data."""
    # Generate large dataset
    data = generate_sample_data(num_rows=10000)
    
    # Save as Parquet
    table_path = os.path.join(test_data_dir, "large_parquet")
    save_as_parquet(data, table_path)
    
    # Verify
    df = pd.read_parquet(table_path)
    assert len(df) == 10000
    
    # Save as CSV
    file_path = os.path.join(test_data_dir, "large.csv")
    save_as_csv(data, file_path)
    
    # Verify
    df = pd.read_csv(file_path)
    assert len(df) == 10000

def test_save_data_with_null_values(test_data_dir):
    """Test saving data with null values."""
    # Create data with null values
    data = pd.DataFrame({
        'id': [1, 2, 3, 4, 5],
        'name': ['Alice', None, 'Charlie', None, 'Eve'],
        'age': [25, 30, None, 40, None],
        'salary': [50000, None, 70000, None, 90000],
        'department': ['HR', 'Engineering', None, 'Sales', None],
        'join_date': [date(2020, 1, 1), None, date(2020, 3, 1), None, date(2020, 5, 1)]
    })
    
    # Save as Parquet
    table_path = os.path.join(test_data_dir, "null_parquet")
    save_as_parquet(data, table_path)
    
    # Verify
    df = pd.read_parquet(table_path)
    assert len(df) == 5
    assert df['name'].isna().sum() == 2
    assert df['age'].isna().sum() == 2
    assert df['salary'].isna().sum() == 2
    assert df['department'].isna().sum() == 2
    assert df['join_date'].isna().sum() == 2
    
    # Save as CSV
    file_path = os.path.join(test_data_dir, "null.csv")
    save_as_csv(data, file_path)
    
    # Verify
    df = pd.read_csv(file_path)
    assert len(df) == 5
    assert df['name'].isna().sum() == 2
    assert df['age'].isna().sum() == 2
    assert df['salary'].isna().sum() == 2
    assert df['department'].isna().sum() == 2
    assert df['join_date'].isna().sum() == 2