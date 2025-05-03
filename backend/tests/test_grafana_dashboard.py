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

import unittest
import os
import pytest
import json
from unittest.mock import patch, MagicMock
from src.scanner.grafana_dashboard import GrafanaDashboard

@pytest.fixture
def mock_response():
    mock = MagicMock()
    mock.json.return_value = {
        "dashboard": {
            "uid": "test-uid",
            "id": 1,
            "title": "Test Dashboard",
            "version": 1,
            "panels": []
        }
    }
    mock.raise_for_status.return_value = None
    return mock

@pytest.fixture
def grafana_dashboard():
    return GrafanaDashboard("http://localhost:3000", "test-api-key")

def test_create_dashboard(grafana_dashboard, mock_response):
    with patch('requests.post', return_value=mock_response):
        response = grafana_dashboard.create_dashboard("Test Dashboard")
        assert response["dashboard"]["uid"] == "test-uid"
        assert response["dashboard"]["id"] == 1
        assert response["dashboard"]["title"] == "Test Dashboard"
        assert response["dashboard"]["version"] == 1
        assert response["dashboard"]["panels"] == []

def test_add_table_metrics_panel(grafana_dashboard, mock_response):
    scan_results = {
        "file_info": {"path": "test.parquet"},
        "row_count": 100,
        "schema": {"id": "integer", "name": "string"}
    }
    
    with patch('requests.get', return_value=mock_response), \
         patch('requests.post', return_value=mock_response):
        response = grafana_dashboard.add_table_metrics_panel("test-uid", scan_results)
        assert response["dashboard"]["uid"] == "test-uid"
        assert response["dashboard"]["id"] == 1
        assert response["dashboard"]["title"] == "Test Dashboard"
        assert response["dashboard"]["version"] == 1

def test_add_column_stats_panel(grafana_dashboard, mock_response):
    column_stats = {
        "id": {"mean": 50, "stddev": 10},
        "age": {"mean": 30, "stddev": 5}
    }
    
    with patch('requests.get', return_value=mock_response), \
         patch('requests.post', return_value=mock_response):
        response = grafana_dashboard.add_column_stats_panel("test-uid", column_stats)
        assert response["dashboard"]["uid"] == "test-uid"
        assert response["dashboard"]["id"] == 1
        assert response["dashboard"]["title"] == "Test Dashboard"
        assert response["dashboard"]["version"] == 1

def test_export_dashboard(grafana_dashboard, mock_response, tmp_path):
    output_path = str(tmp_path / "dashboard.json")
    
    with patch('requests.get', return_value=mock_response):
        result_path = grafana_dashboard.export_dashboard("test-uid", output_path)
        assert result_path == output_path
        
        with open(output_path) as f:
            dashboard = json.load(f)
            assert dashboard["uid"] == "test-uid"
            assert dashboard["id"] == 1
            assert dashboard["title"] == "Test Dashboard"
            assert dashboard["version"] == 1