"""
Tests for the Nessi.dev Airflow operators.
"""
import unittest
from unittest import mock
from datetime import datetime

from airflow.models import DAG, TaskInstance
from airflow.utils.dates import days_ago

from nessi_airflow.operators.data_quality_operator import NessiDataQualityOperator
from nessi_airflow.operators.profile_operator import NessiProfileOperator
from nessi_airflow.operators.validation_operator import NessiValidationOperator
from nessi_airflow.operators.lineage_operator import NessiLineageOperator


class TestNessiOperators(unittest.TestCase):
    """Test the Nessi.dev Airflow operators."""

    def setUp(self):
        """Set up the test case."""
        self.dag = DAG(
            "test_dag",
            default_args={"owner": "airflow", "start_date": days_ago(1)},
            schedule_interval=None,
        )
        
        # Mock the NessiHook
        self.mock_hook = mock.patch(
            "nessi_airflow.hooks.nessi_hook.NessiHook"
        ).start()
        
        # Set up hook instance
        self.hook_instance = self.mock_hook.return_value
    
    def tearDown(self):
        """Tear down the test case."""
        mock.patch.stopall()
    
    def test_data_quality_operator(self):
        """Test the NessiDataQualityOperator."""
        # Set up the operator
        operator = NessiDataQualityOperator(
            task_id="test_task",
            dag=self.dag,
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            quality_threshold=0.9,
            conn_id="nessi_test",
        )
        
        # Mock the hook methods
        self.hook_instance.run_quality_check.return_value = {"check_id": "test-check-id"}
        self.hook_instance.wait_for_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Execute the operator
        result = operator.execute(context={})
        
        # Check the result
        self.assertEqual(result["check_id"], "test-check-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(result["quality_score"], 0.95)
        
        # Check that the hook methods were called correctly
        self.hook_instance.run_quality_check.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
        )
        self.hook_instance.wait_for_quality_results.assert_called_once_with(
            check_id="test-check-id",
            quality_threshold=0.9,
            fail_on_rule_failure=True,
            poke_interval=60,
            timeout=3600,
        )
    
    def test_profile_operator(self):
        """Test the NessiProfileOperator."""
        # Set up the operator
        operator = NessiProfileOperator(
            task_id="test_task",
            dag=self.dag,
            table_name="test_table",
            columns=["id", "name"],
            conn_id="nessi_test",
        )
        
        # Mock the hook methods
        self.hook_instance.run_profile.return_value = {"profile_id": "test-profile-id"}
        self.hook_instance.wait_for_profile.return_value = {
            "profile_id": "test-profile-id",
            "status": "completed",
            "columns": [
                {"name": "id", "type": "integer"},
                {"name": "name", "type": "string"},
            ]
        }
        
        # Execute the operator
        result = operator.execute(context={})
        
        # Check the result
        self.assertEqual(result["profile_id"], "test-profile-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(len(result["columns"]), 2)
        
        # Check that the hook methods were called correctly
        self.hook_instance.run_profile.assert_called_once_with(
            table_name="test_table",
            columns=["id", "name"],
            sample_size=None,
        )
        self.hook_instance.wait_for_profile.assert_called_once_with(
            profile_id="test-profile-id",
            poke_interval=60,
            timeout=3600,
        )
    
    def test_validation_operator(self):
        """Test the NessiValidationOperator."""
        # Set up the operator
        operator = NessiValidationOperator(
            task_id="test_task",
            dag=self.dag,
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
            conn_id="nessi_test",
        )
        
        # Mock the hook methods
        self.hook_instance.run_validation.return_value = {"validation_id": "test-validation-id"}
        self.hook_instance.wait_for_validation_results.return_value = {
            "validation_id": "test-validation-id",
            "status": "completed",
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Execute the operator
        result = operator.execute(context={})
        
        # Check the result
        self.assertEqual(result["validation_id"], "test-validation-id")
        self.assertEqual(result["status"], "completed")
        self.assertEqual(len(result["rule_results"]), 1)
        
        # Check that the hook methods were called correctly
        self.hook_instance.run_validation.assert_called_once_with(
            table_name="test_table",
            rules=[{"name": "test_rule", "rule_type": "not_null", "column": "id"}],
        )
        self.hook_instance.wait_for_validation_results.assert_called_once_with(
            validation_id="test-validation-id",
            fail_on_validation_failure=True,
            poke_interval=60,
            timeout=3600,
        )
    
    def test_lineage_operator(self):
        """Test the NessiLineageOperator."""
        # Set up the operator
        operator = NessiLineageOperator(
            task_id="test_task",
            dag=self.dag,
            table_name="test_table",
            conn_id="nessi_test",
        )
        
        # Mock the hook methods
        self.hook_instance.get_lineage.return_value = {
            "nodes": [{"id": "test-node-id", "name": "test_table"}],
            "edges": [],
        }
        
        # Execute the operator
        result = operator.execute(context={})
        
        # Check the result
        self.assertEqual(len(result["nodes"]), 1)
        self.assertEqual(result["nodes"][0]["name"], "test_table")
        
        # Check that the hook methods were called correctly
        self.hook_instance.get_lineage.assert_called_once_with(
            table_name="test_table",
            max_depth=None,
        )


if __name__ == "__main__":
    unittest.main()
