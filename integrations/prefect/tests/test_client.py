"""
Tests for the Nessi.dev Prefect client.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

import requests

from nessi_prefect.client import NessiClient


class TestNessiClient(unittest.TestCase):
    """Test the NessiClient class."""

    def setUp(self):
        """Set up the test case."""
        # Create a client
        self.client = NessiClient(
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret="test-api-secret",
        )
        
        # Mock requests.Session
        self.mock_session = mock.patch.object(
            self.client, "session", autospec=True
        ).start()
    
    def tearDown(self):
        """Tear down the test case."""
        mock.patch.stopall()
    
    def test_client_initialization(self):
        """Test that the client is initialized correctly."""
        self.assertEqual(self.client.api_host, "http://localhost:8080")
        self.assertEqual(self.client.api_key, "test-api-key")
        self.assertEqual(self.client.api_secret, "test-api-secret")
        self.assertEqual(self.client.timeout, 60)
    
    def test_run_quality_check(self):
        """Test running a quality check."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"check_id": "test-check-id"}
        self.mock_session.post.return_value = mock_response
        
        # Run the quality check
        result = self.client.run_quality_check(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        
        # Check that the request was made correctly
        self.mock_session.post.assert_called_once()
        args, kwargs = self.mock_session.post.call_args
        self.assertEqual(args[0], "http://localhost:8080/quality/check")
        
        # Check the request payload
        payload = kwargs["json"]
        self.assertEqual(payload["table_name"], "test_table")
        self.assertEqual(len(payload["rules"]), 1)
        self.assertEqual(payload["rules"][0]["name"], "test_rule")
        self.assertEqual(payload["profile"], True)
    
    def test_get_quality_results(self):
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
        self.mock_session.get.return_value = mock_response
        
        # Get the quality results
        result = self.client.get_quality_results(check_id="test-check-id")
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        self.assertEqual(len(result["rule_results"]), 1)
        
        # Check that the request was made correctly
        self.mock_session.get.assert_called_once()
        args, kwargs = self.mock_session.get.call_args
        self.assertEqual(args[0], "http://localhost:8080/quality/results/test-check-id")
    
    @mock.patch("nessi_prefect.client.NessiClient.get_quality_results")
    @mock.patch("time.sleep")
    def test_wait_for_quality_results(self, mock_sleep, mock_get_results):
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
        result = self.client.wait_for_quality_results(
            check_id="test-check-id",
            poll_interval=0.01,
            timeout=1,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that get_quality_results was called twice
        self.assertEqual(mock_get_results.call_count, 2)
        
        # Check that sleep was called once
        mock_sleep.assert_called_once_with(0.01)
    
    def test_run_profile(self):
        """Test running a profile."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"profile_id": "test-profile-id"}
        self.mock_session.post.return_value = mock_response
        
        # Run the profile
        result = self.client.run_profile(
            table_name="test_table",
            columns=["id", "name"],
        )
        
        # Check the result
        self.assertEqual(result["profile_id"], "test-profile-id")
        
        # Check that the request was made correctly
        self.mock_session.post.assert_called_once()
        args, kwargs = self.mock_session.post.call_args
        self.assertEqual(args[0], "http://localhost:8080/profile/run")
        
        # Check the request payload
        payload = kwargs["json"]
        self.assertEqual(payload["table_name"], "test_table")
        self.assertEqual(payload["columns"], ["id", "name"])
    
    def test_run_validation(self):
        """Test running a validation."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"validation_id": "test-validation-id"}
        self.mock_session.post.return_value = mock_response
        
        # Run the validation
        result = self.client.run_validation(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
        )
        
        # Check the result
        self.assertEqual(result["validation_id"], "test-validation-id")
        
        # Check that the request was made correctly
        self.mock_session.post.assert_called_once()
        args, kwargs = self.mock_session.post.call_args
        self.assertEqual(args[0], "http://localhost:8080/validation/run")
        
        # Check the request payload
        payload = kwargs["json"]
        self.assertEqual(payload["table_name"], "test_table")
        self.assertEqual(len(payload["rules"]), 1)
        self.assertEqual(payload["rules"][0]["name"], "test_rule")
    
    def test_get_lineage(self):
        """Test getting lineage."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "nodes": [{"id": "test-node-id", "name": "test_table"}],
            "edges": [],
        }
        self.mock_session.get.return_value = mock_response
        
        # Get the lineage
        result = self.client.get_lineage(table_name="test_table")
        
        # Check the result
        self.assertEqual(len(result["nodes"]), 1)
        self.assertEqual(result["nodes"][0]["name"], "test_table")
        
        # Check that the request was made correctly
        self.mock_session.get.assert_called_once()
        args, kwargs = self.mock_session.get.call_args
        self.assertEqual(args[0], "http://localhost:8080/lineage/table/test_table")


if __name__ == "__main__":
    unittest.main()
