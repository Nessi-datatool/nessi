"""
NESSI DATA TOOL - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi Data Tool. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

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
   The Software is licensed, not sold. Nessi Data Tool retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI DATA TOOL BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

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