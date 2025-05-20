"""
Tests for the Nessi.dev Prefect flows.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

from prefect import flow, task, get_run_logger
from prefect.testing.utilities import prefect_test_harness

from nessi_prefect.flows.data_quality_flow import data_quality_flow
from nessi_prefect.flows.data_pipeline_flow import data_pipeline_flow


class TestNessiFlows(unittest.TestCase):
    """Test the Nessi.dev Prefect flows."""

    def setUp(self):
        """Set up the test case."""
        # Set up the Prefect test harness
        self.harness = prefect_test_harness()
        self.harness.start()
        
        # Mock the tasks
        self.mock_run_quality_check = mock.patch(
            "nessi_prefect.flows.data_quality_flow.run_quality_check"
        ).start()
        
        self.mock_wait_for_quality_results = mock.patch(
            "nessi_prefect.flows.data_quality_flow.wait_for_quality_results"
        ).start()
        
        self.mock_run_profile = mock.patch(
            "nessi_prefect.flows.data_quality_flow.run_profile"
        ).start()
        
        self.mock_wait_for_profile = mock.patch(
            "nessi_prefect.flows.data_quality_flow.wait_for_profile"
        ).start()
        
        self.mock_process_data = mock.patch(
            "nessi_prefect.flows.data_pipeline_flow.process_data"
        ).start()
        
        self.mock_run_validation = mock.patch(
            "nessi_prefect.flows.data_pipeline_flow.run_validation"
        ).start()
        
        self.mock_wait_for_validation_results = mock.patch(
            "nessi_prefect.flows.data_pipeline_flow.wait_for_validation_results"
        ).start()
        
        self.mock_get_lineage = mock.patch(
            "nessi_prefect.flows.data_pipeline_flow.get_lineage"
        ).start()
    
    def tearDown(self):
        """Tear down the test case."""
        mock.patch.stopall()
        self.harness.stop()
    
    def test_data_quality_flow(self):
        """Test the data_quality_flow."""
        # Set up the mock task returns
        self.mock_run_quality_check.return_value = {"check_id": "test-check-id"}
        
        self.mock_wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Run the flow
        result = data_quality_flow(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            api_host="http://localhost:8080",
            api_key="test-api-key",
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that the tasks were called correctly
        self.mock_run_quality_check.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            timeout=None,
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret=None,
            secret_block_name=None,
        )
        
        self.mock_wait_for_quality_results.assert_called_once_with(
            check_id="test-check-id",
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            timeout=None,
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret=None,
            secret_block_name=None,
        )
    
    def test_data_quality_flow_with_separate_profile(self):
        """Test the data_quality_flow with a separate profile."""
        # Set up the mock task returns
        self.mock_run_quality_check.return_value = {"check_id": "test-check-id"}
        
        self.mock_wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ],
            # No profile in the quality results
        }
        
        self.mock_run_profile.return_value = {"profile_id": "test-profile-id"}
        
        self.mock_wait_for_profile.return_value = {
            "profile_id": "test-profile-id",
            "status": "completed",
            "columns": [
                {"name": "id", "type": "integer"},
                {"name": "name", "type": "string"},
            ]
        }
        
        # Run the flow
        result = data_quality_flow(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            api_host="http://localhost:8080",
            api_key="test-api-key",
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        self.assertIn("profile", result)
        self.assertEqual(result["profile"]["profile_id"], "test-profile-id")
        
        # Check that the profile tasks were called
        self.mock_run_profile.assert_called_once_with(
            table_name="test_table",
            timeout=None,
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret=None,
            secret_block_name=None,
        )
        
        self.mock_wait_for_profile.assert_called_once_with(
            profile_id="test-profile-id",
            timeout=None,
            api_host="http://localhost:8080",
            api_key="test-api-key",
            api_secret=None,
            secret_block_name=None,
        )
    
    def test_data_pipeline_flow(self):
        """Test the data_pipeline_flow."""
        # Set up the mock task returns
        self.mock_get_lineage.return_value = {
            "nodes": [{"id": "test-node-id", "name": "test_input_table"}],
            "edges": [],
        }
        
        self.mock_run_quality_check.return_value = {"check_id": "test-check-id"}
        
        self.mock_wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        self.mock_process_data.return_value = {
            "input_table": "test_input_table",
            "output_table": "test_output_table",
            "processed_at": "2023-01-01T00:00:00",
            "status": "completed",
        }
        
        self.mock_run_validation.return_value = {"validation_id": "test-validation-id"}
        
        self.mock_wait_for_validation_results.return_value = {
            "validation_id": "test-validation-id",
            "status": "completed",
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Run the flow
        result = data_pipeline_flow(
            input_table="test_input_table",
            output_table="test_output_table",
            quality_rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            validation_rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            capture_lineage=True,
            api_host="http://localhost:8080",
            api_key="test-api-key",
        )
        
        # Check the result
        self.assertEqual(result["input_table"], "test_input_table")
        self.assertEqual(result["output_table"], "test_output_table")
        self.assertIn("quality_results", result)
        self.assertIn("validation_results", result)
        self.assertIn("processing_results", result)
        
        # Check that the tasks were called correctly
        self.mock_get_lineage.assert_called()
        self.mock_run_quality_check.assert_called_once()
        self.mock_wait_for_quality_results.assert_called_once()
        self.mock_process_data.assert_called_once()
        self.mock_run_validation.assert_called_once()
        self.mock_wait_for_validation_results.assert_called_once()


if __name__ == "__main__":
    unittest.main()
