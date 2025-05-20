"""
Tests for the Nessi.dev Airflow hook.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

import requests
from airflow.models.connection import Connection

from nessi_airflow.hooks.nessi_hook import NessiHook


class TestNessiHook(unittest.TestCase):
    """Test the NessiHook class."""

    def setUp(self):
        """Set up the test case."""
        self.conn_id = "nessi_test"
        self.mock_connection = Connection(
            conn_id=self.conn_id,
            conn_type="http",
            host="http://localhost:8080",
            login="test-api-key",
            password="test-api-secret",
        )
        
        # Create a mock connection object
        self.conn = mock.patch(
            "airflow.hooks.base.BaseHook.get_connection",
            return_value=self.mock_connection,
        )
        self.conn.start()
        
        # Create a hook
        self.hook = NessiHook(conn_id=self.conn_id)
    
    def tearDown(self):
        """Tear down the test case."""
        self.conn.stop()
    
    def test_hook_initialization(self):
        """Test that the hook is initialized correctly."""
        self.assertEqual(self.hook.conn_id, self.conn_id)
        self.assertEqual(self.hook.api_host, "http://localhost:8080")
        self.assertEqual(self.hook.api_key, "test-api-key")
        self.assertEqual(self.hook.api_secret, "test-api-secret")
    
    @mock.patch("requests.Session.post")
    def test_run_quality_check(self, mock_post):
        """Test running a quality check."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"check_id": "test-check-id"}
        mock_post.return_value = mock_response
        
        # Run the quality check
        result = self.hook.run_quality_check(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        
        # Check that the request was made correctly
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(kwargs["url"], "http://localhost:8080/quality/check")
        self.assertEqual(kwargs["headers"]["X-API-Key"], "test-api-key")
        self.assertEqual(kwargs["headers"]["X-API-Secret"], "test-api-secret")
        
        # Check the request payload
        payload = json.loads(kwargs["data"])
        self.assertEqual(payload["table_name"], "test_table")
        self.assertEqual(len(payload["rules"]), 1)
        self.assertEqual(payload["rules"][0]["name"], "test_rule")
        self.assertEqual(payload["profile"], True)
    
    @mock.patch("requests.Session.get")
    def test_get_quality_results(self, mock_get):
        """Test getting quality check results."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        mock_get.return_value = mock_response
        
        # Get the quality results
        result = self.hook.get_quality_results(check_id="test-check-id")
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        self.assertEqual(len(result["rule_results"]), 1)
        
        # Check that the request was made correctly
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        self.assertEqual(kwargs["url"], "http://localhost:8080/quality/results/test-check-id")
        self.assertEqual(kwargs["headers"]["X-API-Key"], "test-api-key")
        self.assertEqual(kwargs["headers"]["X-API-Secret"], "test-api-secret")
    
    @mock.patch("nessi_airflow.hooks.nessi_hook.NessiHook.get_quality_results")
    def test_wait_for_quality_results(self, mock_get_results):
        """Test waiting for quality check results."""
        # Mock the response
        mock_get_results.side_effect = [
            # First call - still running
            {
                "check_id": "test-check-id",
                "status": "running",
            },
            # Second call - completed
            {
                "check_id": "test-check-id",
                "status": "completed",
                "quality_score": 0.95,
                "rule_results": [
                    {"name": "test_rule", "status": "passed"}
                ]
            }
        ]
        
        # Wait for the quality results
        result = self.hook.wait_for_quality_results(
            check_id="test-check-id",
            poke_interval=0.01,
            timeout=1,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that get_quality_results was called twice
        self.assertEqual(mock_get_results.call_count, 2)
    
    @mock.patch("requests.Session.post")
    def test_run_profile(self, mock_post):
        """Test running a profile."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"profile_id": "test-profile-id"}
        mock_post.return_value = mock_response
        
        # Run the profile
        result = self.hook.run_profile(
            table_name="test_table",
            columns=["id", "name"],
        )
        
        # Check the result
        self.assertEqual(result["profile_id"], "test-profile-id")
        
        # Check that the request was made correctly
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(kwargs["url"], "http://localhost:8080/profile/run")
        
        # Check the request payload
        payload = json.loads(kwargs["data"])
        self.assertEqual(payload["table_name"], "test_table")
        self.assertEqual(payload["columns"], ["id", "name"])
    
    @mock.patch("requests.Session.post")
    def test_run_validation(self, mock_post):
        """Test running a validation."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"validation_id": "test-validation-id"}
        mock_post.return_value = mock_response
        
        # Run the validation
        result = self.hook.run_validation(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
        )
        
        # Check the result
        self.assertEqual(result["validation_id"], "test-validation-id")
        
        # Check that the request was made correctly
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(kwargs["url"], "http://localhost:8080/validation/run")
        
        # Check the request payload
        payload = json.loads(kwargs["data"])
        self.assertEqual(payload["table_name"], "test_table")
        self.assertEqual(len(payload["rules"]), 1)
        self.assertEqual(payload["rules"][0]["name"], "test_rule")
    
    @mock.patch("requests.Session.get")
    def test_get_lineage(self, mock_get):
        """Test getting lineage."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "nodes": [{"id": "test-node-id", "name": "test_table"}],
            "edges": [],
        }
        mock_get.return_value = mock_response
        
        # Get the lineage
        result = self.hook.get_lineage(table_name="test_table")
        
        # Check the result
        self.assertEqual(len(result["nodes"]), 1)
        self.assertEqual(result["nodes"][0]["name"], "test_table")
        
        # Check that the request was made correctly
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        self.assertEqual(kwargs["url"], "http://localhost:8080/lineage/table/test_table")


if __name__ == "__main__":
    unittest.main()
