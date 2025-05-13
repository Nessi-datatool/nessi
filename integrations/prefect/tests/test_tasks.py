"""
Tests for the Nessi.dev Prefect tasks.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

from prefect import flow, task, get_run_logger
from prefect.testing.utilities import prefect_test_harness

from nessi_prefect.tasks.quality import run_quality_check, wait_for_quality_results
from nessi_prefect.tasks.profile import run_profile, wait_for_profile
from nessi_prefect.tasks.validation import run_validation, wait_for_validation_results
from nessi_prefect.tasks.lineage import get_lineage, visualize_lineage


class TestNessiTasks(unittest.TestCase):
    """Test the Nessi.dev Prefect tasks."""

    def setUp(self):
        """Set up the test case."""
        # Set up the Prefect test harness
        self.harness = prefect_test_harness()
        self.harness.start()
        
        # Mock the NessiClient
        self.mock_client = mock.patch(
            "nessi_prefect.client.NessiClient"
        ).start()
        
        # Set up client instance
        self.client_instance = self.mock_client.return_value
    
    def tearDown(self):
        """Tear down the test case."""
        mock.patch.stopall()
        self.harness.stop()
    
    @mock.patch("nessi_prefect.tasks.quality.NessiClient")
    def test_run_quality_check(self, mock_client_class):
        """Test the run_quality_check task."""
        # Set up the mock client
        mock_client = mock_client_class.return_value
        mock_client.run_quality_check.return_value = {"check_id": "test-check-id"}
        
        # Create a test flow
        @flow
        def test_flow():
            return run_quality_check(
                table_name="test_table",
                rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
                profile=True,
                api_host="http://localhost:8080",
                api_key="test-api-key",
            )
        
        # Run the flow
        result = test_flow()
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        
        # Check that the client was created correctly
        mock_client_class.assert_called_once_with(
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret=None,
            timeout=60,
        )
        
        # Check that the client methods were called correctly
        mock_client.run_quality_check.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            timeout=None,
        )
    
    @mock.patch("nessi_prefect.tasks.quality.NessiClient")
    def test_wait_for_quality_results(self, mock_client_class):
        """Test the wait_for_quality_results task."""
        # Set up the mock client
        mock_client = mock_client_class.return_value
        mock_client.wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Create a test flow
        @flow
        def test_flow():
            return wait_for_quality_results(
                check_id="test-check-id",
                quality_threshold=0.9,
                fail_on_rule_failure=True,
                api_host="http://localhost:8080",
                api_key="test-api-key",
            )
        
        # Run the flow
        result = test_flow()
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that the client was created correctly
        mock_client_class.assert_called_once_with(
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret=None,
            timeout=60,
        )
        
        # Check that the client methods were called correctly
        mock_client.wait_for_quality_results.assert_called_once_with(
            check_id="test-check-id",
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            timeout=None,
            poll_interval=10,
        )
    
    @mock.patch("nessi_prefect.tasks.profile.NessiClient")
    def test_run_profile(self, mock_client_class):
        """Test the run_profile task."""
        # Set up the mock client
        mock_client = mock_client_class.return_value
        mock_client.run_profile.return_value = {"profile_id": "test-profile-id"}
        
        # Create a test flow
        @flow
        def test_flow():
            return run_profile(
                table_name="test_table",
                columns=["id", "name"],
                api_host="http://localhost:8080",
                api_key="test-api-key",
            )
        
        # Run the flow
        result = test_flow()
        
        # Check the result
        self.assertEqual(result["profile_id"], "test-profile-id")
        
        # Check that the client methods were called correctly
        mock_client.run_profile.assert_called_once_with(
            table_name="test_table",
            columns=["id", "name"],
            sample_size=None,
            timeout=None,
        )
    
    @mock.patch("nessi_prefect.tasks.validation.NessiClient")
    def test_run_validation(self, mock_client_class):
        """Test the run_validation task."""
        # Set up the mock client
        mock_client = mock_client_class.return_value
        mock_client.run_validation.return_value = {"validation_id": "test-validation-id"}
        
        # Create a test flow
        @flow
        def test_flow():
            return run_validation(
                table_name="test_table",
                rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
                api_host="http://localhost:8080",
                api_key="test-api-key",
            )
        
        # Run the flow
        result = test_flow()
        
        # Check the result
        self.assertEqual(result["validation_id"], "test-validation-id")
        
        # Check that the client methods were called correctly
        mock_client.run_validation.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            timeout=None,
        )
    
    @mock.patch("nessi_prefect.tasks.lineage.NessiClient")
    def test_get_lineage(self, mock_client_class):
        """Test the get_lineage task."""
        # Set up the mock client
        mock_client = mock_client_class.return_value
        mock_client.get_lineage.return_value = {
            "nodes": [{"id": "test-node-id", "name": "test_table"}],
            "edges": [],
        }
        
        # Create a test flow
        @flow
        def test_flow():
            return get_lineage(
                table_name="test_table",
                api_host="http://localhost:8080",
                api_key="test-api-key",
            )
        
        # Run the flow
        result = test_flow()
        
        # Check the result
        self.assertEqual(len(result["nodes"]), 1)
        self.assertEqual(result["nodes"][0]["name"], "test_table")
        
        # Check that the client methods were called correctly
        mock_client.get_lineage.assert_called_once_with(
            table_name="test_table",
            max_depth=None,
        )


if __name__ == "__main__":
    unittest.main()
