"""
Tests for the Nessi.dev Dagster ops.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

from dagster import build_op_context, op

from nessi_dagster.ops.quality_ops import run_quality_check, wait_for_quality_results
from nessi_dagster.ops.profile_ops import run_profile, wait_for_profile
from nessi_dagster.ops.validation_ops import run_validation, wait_for_validation_results
from nessi_dagster.ops.lineage_ops import get_lineage, visualize_lineage


class TestNessiOps(unittest.TestCase):
    """Test the Nessi.dev Dagster ops."""

    def setUp(self):
        """Set up the test case."""
        # Mock the NessiClient
        self.mock_client = mock.Mock()
        
        # Create a mock resource
        self.mock_resource = mock.Mock()
        self.mock_resource.client = self.mock_client
        
        # Create a context with the mock resource
        self.context = build_op_context(
            resources={"nessi": self.mock_resource}
        )
    
    def test_run_quality_check(self):
        """Test the run_quality_check op."""
        # Set up the mock client
        self.mock_client.run_quality_check.return_value = {"check_id": "test-check-id"}
        
        # Run the op
        result = run_quality_check(
            self.context,
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        
        # Check that the client methods were called correctly
        self.mock_client.run_quality_check.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            timeout=None,
        )
    
    def test_wait_for_quality_results(self):
        """Test the wait_for_quality_results op."""
        # Set up the mock client
        self.mock_client.wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Run the op
        result = wait_for_quality_results(
            self.context,
            check_id="test-check-id",
            quality_threshold=0.9,
            fail_on_rule_failure=True,
        )
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that the client methods were called correctly
        self.mock_client.wait_for_quality_results.assert_called_once_with(
            check_id="test-check-id",
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            timeout=None,
        )
    
    def test_wait_for_quality_results_below_threshold(self):
        """Test the wait_for_quality_results op when quality is below threshold."""
        # Set up the mock client
        self.mock_client.wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.8,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Run the op and expect an exception
        with self.assertRaises(RuntimeError) as context:
            wait_for_quality_results(
                self.context,
                check_id="test-check-id",
                quality_threshold=0.9,
                fail_on_rule_failure=True,
            )
        
        # Check the exception message
        self.assertIn("Quality score 0.8 is below threshold 0.9", str(context.exception))
    
    def test_wait_for_quality_results_rule_failure(self):
        """Test the wait_for_quality_results op when a rule fails."""
        # Set up the mock client
        self.mock_client.wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "failed"}
            ]
        }
        
        # Run the op and expect an exception
        with self.assertRaises(RuntimeError) as context:
            wait_for_quality_results(
                self.context,
                check_id="test-check-id",
                quality_threshold=0.9,
                fail_on_rule_failure=True,
            )
        
        # Check the exception message
        self.assertIn("The following rules failed: test_rule", str(context.exception))
    
    def test_run_profile(self):
        """Test the run_profile op."""
        # Set up the mock client
        self.mock_client.run_profile.return_value = {"profile_id": "test-profile-id"}
        
        # Run the op
        result = run_profile(
            self.context,
            table_name="test_table",
            columns=["id", "name"],
        )
        
        # Check the result
        self.assertEqual(result["profile_id"], "test-profile-id")
        
        # Check that the client methods were called correctly
        self.mock_client.run_profile.assert_called_once_with(
            table_name="test_table",
            columns=["id", "name"],
            sample_size=None,
            timeout=None,
        )
    
    def test_wait_for_profile(self):
        """Test the wait_for_profile op."""
        # Set up the mock client
        self.mock_client.wait_for_profile.return_value = {
            "profile_id": "test-profile-id",
            "status": "completed",
            "columns": [
                {"name": "id", "type": "integer"},
                {"name": "name", "type": "string"},
            ]
        }
        
        # Run the op
        result = wait_for_profile(
            self.context,
            profile_id="test-profile-id",
        )
        
        # Check the result
        self.assertEqual(result["profile_id"], "test-profile-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(len(result["columns"]), 2)
        
        # Check that the client methods were called correctly
        self.mock_client.wait_for_profile.assert_called_once_with(
            profile_id="test-profile-id",
            timeout=None,
        )
    
    def test_run_validation(self):
        """Test the run_validation op."""
        # Set up the mock client
        self.mock_client.run_validation.return_value = {"validation_id": "test-validation-id"}
        
        # Run the op
        result = run_validation(
            self.context,
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
        )
        
        # Check the result
        self.assertEqual(result["validation_id"], "test-validation-id")
        
        # Check that the client methods were called correctly
        self.mock_client.run_validation.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            timeout=None,
        )
    
    def test_wait_for_validation_results(self):
        """Test the wait_for_validation_results op."""
        # Set up the mock client
        self.mock_client.wait_for_validation_results.return_value = {
            "validation_id": "test-validation-id",
            "status": "completed",
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Run the op
        result = wait_for_validation_results(
            self.context,
            validation_id="test-validation-id",
            fail_on_validation_failure=True,
        )
        
        # Check the result
        self.assertEqual(result["validation_id"], "test-validation-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(len(result["rule_results"]), 1)
        
        # Check that the client methods were called correctly
        self.mock_client.wait_for_validation_results.assert_called_once_with(
            validation_id="test-validation-id",
            fail_on_validation_failure=True,
            timeout=None,
        )
    
    def test_get_lineage(self):
        """Test the get_lineage op."""
        # Set up the mock client
        self.mock_client.get_lineage.return_value = {
            "nodes": [{"id": "test-node-id", "name": "test_table"}],
            "edges": [],
        }
        
        # Run the op
        result = get_lineage(
            self.context,
            table_name="test_table",
        )
        
        # Check the result
        self.assertEqual(len(result["nodes"]), 1)
        self.assertEqual(result["nodes"][0]["name"], "test_table")
        
        # Check that the client methods were called correctly
        self.mock_client.get_lineage.assert_called_once_with(
            table_name="test_table",
            max_depth=None,
        )
    
    def test_visualize_lineage(self):
        """Test the visualize_lineage op."""
        # Set up the mock client
        self.mock_client.visualize_lineage.return_value = {
            "format": "html",
            "content": "<html>...</html>",
        }
        
        # Run the op
        result = visualize_lineage(
            self.context,
            table_name="test_table",
            format="html",
        )
        
        # Check the result
        self.assertEqual(result["format"], "html")
        self.assertEqual(result["content"], "<html>...</html>")
        
        # Check that the client methods were called correctly
        self.mock_client.visualize_lineage.assert_called_once_with(
            table_name="test_table",
            format="html",
            max_depth=None,
        )


if __name__ == "__main__":
    unittest.main()
