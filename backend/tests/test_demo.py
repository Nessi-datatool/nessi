"""Test the demo functionality."""

import os
from unittest.mock import patch, MagicMock
import pytest
from src.demo import main, wait_for_grafana

@pytest.fixture
def mock_grafana():
    """Mock Grafana responses."""
    with patch('requests.get') as mock_get:
        mock_get.return_value.status_code = 200
        yield mock_get

def test_wait_for_grafana(mock_grafana):
    """Test waiting for Grafana to be ready."""
    assert wait_for_grafana("http://localhost:3000")

def test_main_function(mock_grafana):
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