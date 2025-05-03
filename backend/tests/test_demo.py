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

"""Test the demo functionality."""

import os
from unittest.mock import patch, MagicMock
import pytest
from src.demo import main, wait_for_grafana
from src.scanner.spark_config import get_spark_session, stop_spark_session

@pytest.fixture(scope="session")
def spark_session():
    spark = get_spark_session("TestSession")
    yield spark
    stop_spark_session()

@pytest.fixture
def mock_grafana():
    """Mock Grafana responses."""
    with patch('requests.get') as mock_get:
        mock_get.return_value.status_code = 200
        yield mock_get

def test_wait_for_grafana(mock_grafana):
    """Test waiting for Grafana to be ready."""
    assert wait_for_grafana("http://localhost:3000")

def test_main_function(mock_grafana, spark_session):
    """Test main function"""
    # Set up environment variables
    os.environ["GRAFANA_URL"] = "http://localhost:3000"
    os.environ["GF_SECURITY_ADMIN_PASSWORD"] = "admin"
    
    # Mock Grafana dashboard creation
    mock_grafana.return_value.json.return_value = {"dashboard": {"uid": "test"}}
    
    # Run main function
    main()
    
    # Check if demo directory and files were created
    demo_dir = "data/demo"
    assert os.path.exists(demo_dir)
    assert os.path.exists(os.path.join(demo_dir, "sample.parquet"))
    assert os.path.exists(os.path.join(demo_dir, "sample.csv"))
    
    # Verify that Grafana was called
    mock_grafana.assert_any_call("http://localhost:3000/api/health")