"""
Tests for the Nessi.dev Dagster resource.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

import requests
from dagster import build_init_resource_context

from nessi_dagster.resources.nessi_resource import NessiClient, nessi_resource


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
    
    def test_create_session(self):
        """Test creating a session."""
        # Create a new client to test session creation
        with mock.patch("requests.Session") as mock_session:
            mock_session_instance = mock_session.return_value
            
            client = NessiClient(
                api_host="http://localhost:8080",
                api_key="test-api-key",
                api_secret="test-api-secret",
            )
            
            # Check that the session was created correctly
            mock_session.assert_called_once()
            self.assertEqual(
                mock_session_instance.headers["X-API-Key"], "test-api-key"
            )
            self.assertEqual(
                mock_session_instance.headers["X-API-Secret"], "test-api-secret"
            )
            self.assertEqual(
                mock_session_instance.headers["Content-Type"], "application/json"
            )
    
    def test_do_api_call_get(self):
        """Test making a GET API call."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.content = b'{"key": "value"}'
        mock_response.json.return_value = {"key": "value"}
        self.mock_session.get.return_value = mock_response
        
        # Make the API call
        result = self.client._do_api_call(
            endpoint="test/endpoint",
            method="GET",
            params={"param": "value"},
        )
        
        # Check the result
        self.assertEqual(result, {"key": "value"})
        
        # Check that the request was made correctly
        self.mock_session.get.assert_called_once_with(
            "http://localhost:8080/test/endpoint",
            params={"param": "value"},
            timeout=60,
        )
    
    def test_do_api_call_post(self):
        """Test making a POST API call."""
        # Mock the response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.content = b'{"key": "value"}'
        mock_response.json.return_value = {"key": "value"}
        self.mock_session.post.return_value = mock_response
        
        # Make the API call
        result = self.client._do_api_call(
            endpoint="test/endpoint",
            method="POST",
            data={"data": "value"},
            params={"param": "value"},
        )
        
        # Check the result
        self.assertEqual(result, {"key": "value"})
        
        # Check that the request was made correctly
        self.mock_session.post.assert_called_once_with(
            "http://localhost:8080/test/endpoint",
            json={"data": "value"},
            params={"param": "value"},
            timeout=60,
        )
    
    def test_run_quality_check(self):
        """Test running a quality check."""
        # Mock the API call
        mock_do_api_call = mock.patch.object(
            self.client, "_do_api_call", return_value={"check_id": "test-check-id"}
        ).start()
        
        # Run the quality check
        result = self.client.run_quality_check(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        
        # Check that the API call was made correctly
        mock_do_api_call.assert_called_once_with(
            "quality/check",
            method="POST",
            data={
                "table_name": "test_table",
                "rules": [{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
                "profile": True,
            },
        )
    
    def test_get_quality_results(self):
        """Test getting quality check results."""
        # Mock the API call
        mock_do_api_call = mock.patch.object(
            self.client, "_do_api_call", return_value={
                "check_id": "test-check-id",
                "status": "completed",
                "quality_score": 0.95,
                "rule_results": [
                    {"name": "test_rule", "status": "passed"}
                ]
            }
        ).start()
        
        # Get the quality results
        result = self.client.get_quality_results(check_id="test-check-id")
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that the API call was made correctly
        mock_do_api_call.assert_called_once_with(
            "quality/results/test-check-id",
        )
    
    @mock.patch("nessi_dagster.resources.nessi_resource.NessiClient.get_quality_results")
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


class TestNessiResource(unittest.TestCase):
    """Test the NessiResource."""

    def test_nessi_resource(self):
        """Test the nessi_resource factory."""
        # Create a resource context
        init_context = build_init_resource_context(
            config={
                "api_host": "http://localhost:8080",
                "api_key": "test-api-key",
                "api_secret": "test-api-secret",
                "timeout": 60,
            }
        )
        
        # Create the resource
        with mock.patch("nessi_dagster.resources.nessi_resource.NessiClient") as mock_client:
            resource = nessi_resource(init_context)
            
            # Set up the resource for execution
            resource.setup_for_execution(init_context)
            
            # Check that the client was created correctly
            mock_client.assert_called_once_with(
                api_host="http://localhost:8080",
                api_key="test-api-key",
                api_secret="test-api-secret",
                timeout=60,
            )
            
            # Check the resource properties
            self.assertEqual(resource.api_host, "http://localhost:8080")
            self.assertEqual(resource.api_key, "test-api-key")
            self.assertEqual(resource.api_secret, "test-api-secret")
            self.assertEqual(resource.timeout, 60)


if __name__ == "__main__":
    unittest.main()
