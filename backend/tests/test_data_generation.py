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