"""
Tests for the Nessi.dev Dagster jobs.
"""
import unittest
from unittest import mock
import json
from datetime import datetime

from dagster import build_op_context, ExecuteInProcessResult, DagsterInstance
from dagster._core.execution.api import create_execution_plan, execute_plan

from nessi_dagster.jobs.data_quality_job import data_quality_job
from nessi_dagster.jobs.data_pipeline_job import data_pipeline_job, process_data


class TestNessiJobs(unittest.TestCase):
    """Test the Nessi.dev Dagster jobs."""

    def setUp(self):
        """Set up the test case."""
        # Mock the ops
        self.mock_run_quality_check = mock.patch(
            "nessi_dagster.jobs.data_quality_job.run_quality_check"
        ).start()
        
        self.mock_wait_for_quality_results = mock.patch(
            "nessi_dagster.jobs.data_quality_job.wait_for_quality_results"
        ).start()
        
        self.mock_run_profile = mock.patch(
            "nessi_dagster.jobs.data_pipeline_job.run_profile"
        ).start()
        
        self.mock_wait_for_profile = mock.patch(
            "nessi_dagster.jobs.data_pipeline_job.wait_for_profile"
        ).start()
        
        self.mock_run_validation = mock.patch(
            "nessi_dagster.jobs.data_pipeline_job.run_validation"
        ).start()
        
        self.mock_wait_for_validation_results = mock.patch(
            "nessi_dagster.jobs.data_pipeline_job.wait_for_validation_results"
        ).start()
        
        self.mock_get_lineage = mock.patch(
            "nessi_dagster.jobs.data_pipeline_job.get_lineage"
        ).start()
    
    def tearDown(self):
        """Tear down the test case."""
        mock.patch.stopall()
    
    def test_process_data_op(self):
        """Test the process_data op."""
        # Create a context
        context = build_op_context()
        
        # Run the op
        result = process_data(
            context,
            input_table="test_input_table",
            output_table="test_output_table",
            quality_results={
                "check_id": "test-check-id",
                "status": "completed",
                "quality_score": 0.95,
            },
        )
        
        # Check the result
        self.assertEqual(result["input_table"], "test_input_table")
        self.assertEqual(result["output_table"], "test_output_table")
        self.assertEqual(result["quality_score"], 0.95)
        self.assertEqual(result["status"], "completed")
    
    @mock.patch("dagster._core.execution.api.execute_job")
    def test_data_quality_job(self, mock_execute_job):
        """Test the data_quality_job."""
        # Set up the mock op returns
        self.mock_run_quality_check.return_value = {"check_id": "test-check-id"}
        
        self.mock_wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Create a mock result
        mock_result = mock.Mock(spec=ExecuteInProcessResult)
        mock_result.success = True
        mock_execute_job.return_value = mock_result
        
        # Create a run config
        run_config = {
            "resources": {
                "nessi": {
                    "config": {
                        "api_host": "http://localhost:8080",
                        "api_key": "test-api-key",
                        "api_secret": "test-api-secret",
                        "timeout": 300,
                    }
                }
            },
            "ops": {
                "run_quality_check": {
                    "inputs": {
                        "table_name": "test_table",
                        "rules": [{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
                        "profile": True,
                    }
                },
                "wait_for_quality_results": {
                    "inputs": {
                        "quality_threshold": 0.9,
                        "fail_on_rule_failure": True,
                    }
                }
            }
        }
        
        # Execute the job
        with mock.patch("dagster._core.instance.DagsterInstance.ephemeral"):
            result = data_quality_job.execute_in_process(run_config=run_config)
        
        # Check that the job executed
        mock_execute_job.assert_called()
    
    @mock.patch("dagster._core.execution.api.execute_job")
    def test_data_pipeline_job(self, mock_execute_job):
        """Test the data_pipeline_job."""
        # Set up the mock op returns
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
        
        self.mock_run_validation.return_value = {"validation_id": "test-validation-id"}
        
        self.mock_wait_for_validation_results.return_value = {
            "validation_id": "test-validation-id",
            "status": "completed",
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Create a mock result
        mock_result = mock.Mock(spec=ExecuteInProcessResult)
        mock_result.success = True
        mock_execute_job.return_value = mock_result
        
        # Create a run config
        run_config = {
            "resources": {
                "nessi": {
                    "config": {
                        "api_host": "http://localhost:8080",
                        "api_key": "test-api-key",
                        "api_secret": "test-api-secret",
                        "timeout": 300,
                    }
                }
            },
            "ops": {
                "get_lineage": {
                    "inputs": {
                        "table_name": "test_input_table",
                    }
                },
                "run_quality_check": {
                    "inputs": {
                        "table_name": "test_input_table",
                        "rules": [{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
                        "profile": True,
                    }
                },
                "process_data": {
                    "inputs": {
                        "input_table": "test_input_table",
                        "output_table": "test_output_table",
                    }
                },
                "run_validation": {
                    "inputs": {
                        "table_name": "test_output_table",
                        "rules": [{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
                    }
                },
            }
        }
        
        # Execute the job
        with mock.patch("dagster._core.instance.DagsterInstance.ephemeral"):
            result = data_pipeline_job.execute_in_process(run_config=run_config)
        
        # Check that the job executed
        mock_execute_job.assert_called()


if __name__ == "__main__":
    unittest.main()
