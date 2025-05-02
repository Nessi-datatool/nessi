import pytest
import json
from unittest.mock import patch, MagicMock
from backend.src.scanner.grafana_dashboard import GrafanaDashboard

@pytest.fixture
def mock_response():
    mock = MagicMock()
    mock.json.return_value = {
        "dashboard": {
            "uid": "test-uid",
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

def test_add_column_stats_panel(grafana_dashboard, mock_response):
    column_stats = {
        "id": {"mean": 50, "stddev": 10},
        "age": {"mean": 30, "stddev": 5}
    }
    
    with patch('requests.get', return_value=mock_response), \
         patch('requests.post', return_value=mock_response):
        response = grafana_dashboard.add_column_stats_panel("test-uid", column_stats)
        assert response["dashboard"]["uid"] == "test-uid"

def test_export_dashboard(grafana_dashboard, mock_response, tmp_path):
    output_path = str(tmp_path / "dashboard.json")
    
    with patch('requests.get', return_value=mock_response):
        result_path = grafana_dashboard.export_dashboard("test-uid", output_path)
        assert result_path == output_path
        
        with open(output_path) as f:
            dashboard = json.load(f)
            assert dashboard["uid"] == "test-uid" 