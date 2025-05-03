"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi. All rights reserved.

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
   The Software is licensed, not sold. Nessi retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

import os
import json
import time
import pytest
from unittest.mock import patch, MagicMock
from src.reports.setup_grafana import GrafanaSetup

@pytest.fixture
def mock_requests():
    """Mock requests module for testing."""
    with patch('src.reports.setup_grafana.requests') as mock_requests:
        yield mock_requests

@pytest.fixture
def mock_time():
    with patch("src.reports.setup_grafana.time") as mock:
        # Make time.time() return values that will cause a timeout
        mock.time.side_effect = [0] * 5 + [1]  # Return 0 five times, then 1 to trigger timeout
        mock.sleep.return_value = None
        yield mock

@pytest.fixture
def mock_dashboard_file(tmp_path):
    """Create a mock dashboard file for testing."""
    dashboard = {
        "dashboard": {
            "id": 1,
            "uid": "test-uid",
            "title": "Test Dashboard"
        }
    }
    dashboard_path = tmp_path / "grafana_dashboard.json"
    with open(dashboard_path, 'w') as f:
        json.dump(dashboard, f)
    return dashboard_path

@pytest.fixture
def setup(mock_requests, mock_dashboard_file):
    """Create a GrafanaSetup instance with mocked dependencies."""
    with patch('src.reports.setup_grafana.Path') as mock_path:
        mock_path.return_value.parent = mock_dashboard_file.parent
        setup = GrafanaSetup()
        setup.dashboard_path = mock_dashboard_file
        return setup

def test_wait_for_grafana_success(setup, mock_requests):
    """Test successful Grafana readiness check."""
    mock_response = MagicMock()
    mock_response.status_code = 200
    mock_requests.get.return_value = mock_response
    
    assert setup.wait_for_grafana(timeout=1) is True
    mock_requests.get.assert_called_once_with("http://localhost:3000/api/health")

def test_wait_for_grafana_timeout(setup, mock_requests, mock_time):
    """Test Grafana readiness check timeout."""
    # Set up mock to raise exception
    mock_requests.get.side_effect = [Exception("Connection error")] * 5
    
    # Set a very short timeout
    result = setup.wait_for_grafana(timeout=0.1)
    
    # Verify that we tried multiple times before timing out
    assert result is False
    assert mock_requests.get.call_count > 1

def test_setup_prometheus_datasource_success(setup, mock_requests):
    """Test successful Prometheus data source setup."""
    mock_response = MagicMock()
    mock_response.status_code = 200
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_prometheus_datasource() is True
    mock_requests.post.assert_called_once()
    call_args = mock_requests.post.call_args[1]
    assert call_args['json']['name'] == "Prometheus"
    assert call_args['json']['type'] == "prometheus"
    assert call_args['json']['url'] == "http://prometheus:9090"

def test_setup_prometheus_datasource_exists(setup, mock_requests):
    """Test Prometheus data source setup when it already exists."""
    mock_response = MagicMock()
    mock_response.status_code = 409  # Already exists
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_prometheus_datasource() is True

def test_setup_prometheus_datasource_failure(setup, mock_requests):
    """Test failed Prometheus data source setup."""
    mock_response = MagicMock()
    mock_response.status_code = 500
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_prometheus_datasource() is False

def test_setup_dashboard_success(setup, mock_requests):
    """Test successful dashboard setup."""
    mock_response = MagicMock()
    mock_response.status_code = 200
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_dashboard() is True
    mock_requests.post.assert_called_once()
    call_args = mock_requests.post.call_args[1]
    assert call_args['json']['dashboard']['uid'] == "nessi-reports"
    assert call_args['json']['dashboard']['title'] == "Nessi Reports Dashboard"

def test_setup_dashboard_failure(setup, mock_requests):
    """Test failed dashboard setup."""
    mock_response = MagicMock()
    mock_response.status_code = 500
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_dashboard() is False

def test_setup_alert_channels_success(setup, mock_requests):
    """Test successful alert channels setup."""
    mock_response = MagicMock()
    mock_response.status_code = 200
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_alert_channels() is True
    mock_requests.post.assert_called_once()
    call_args = mock_requests.post.call_args[1]
    assert call_args['json']['name'] == "Email Alerts"
    assert call_args['json']['type'] == "email"

def test_setup_alert_channels_exists(setup, mock_requests):
    """Test alert channels setup when they already exist."""
    mock_response = MagicMock()
    mock_response.status_code = 409  # Already exists
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_alert_channels() is True

def test_setup_alert_channels_failure(setup, mock_requests):
    """Test failed alert channels setup."""
    mock_response = MagicMock()
    mock_response.status_code = 500
    mock_requests.post.return_value = mock_response
    
    assert setup.setup_alert_channels() is False

def test_run_success(setup, mock_requests):
    """Test successful complete setup process."""
    # Mock all required responses
    mock_health = MagicMock()
    mock_health.status_code = 200
    mock_prometheus = MagicMock()
    mock_prometheus.status_code = 200
    mock_dashboard = MagicMock()
    mock_dashboard.status_code = 200
    mock_alerts = MagicMock()
    mock_alerts.status_code = 200
    
    mock_requests.get.return_value = mock_health
    mock_requests.post.side_effect = [mock_prometheus, mock_dashboard, mock_alerts]
    
    assert setup.run() is True
    assert mock_requests.get.call_count == 1
    assert mock_requests.post.call_count == 3

def test_run_grafana_not_ready(setup, mock_requests):
    """Test setup process when Grafana is not ready."""
    mock_requests.get.side_effect = Exception("Connection error")
    
    assert setup.run() is False
    assert mock_requests.post.call_count == 0

def test_run_prometheus_failure(setup, mock_requests):
    """Test setup process when Prometheus setup fails."""
    mock_health = MagicMock()
    mock_health.status_code = 200
    mock_prometheus = MagicMock()
    mock_prometheus.status_code = 500
    
    mock_requests.get.return_value = mock_health
    mock_requests.post.return_value = mock_prometheus
    
    assert setup.run() is False
    assert mock_requests.post.call_count == 1