"""
Tests for the Nessi.dev Airflow sensors.
"""
import unittest
from unittest import mock
from datetime import datetime

from airflow.models import DAG, TaskInstance
from airflow.utils.dates import days_ago

from nessi_airflow.sensors.data_quality_sensor import NessiDataQualitySensor
from nessi_airflow.sensors.validation_sensor import NessiValidationSensor
from nessi_airflow.sensors.profile_sensor import NessiProfileSensor


class TestNessiSensors(unittest.TestCase):
    """Test the Nessi.dev Airflow sensors."""

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
    
    def test_data_quality_sensor_running(self):
        """Test the NessiDataQualitySensor when the check is still running."""
        # Set up the sensor
        sensor = NessiDataQualitySensor(
            task_id="test_task",
            dag=self.dag,
            check_id="test-check-id",
            quality_threshold=0.9,
            conn_id="nessi_test",
        )
        
        # Mock the hook methods to return a running check
        self.hook_instance.get_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "running",
        }
        
        # Execute the sensor's poke method
        result = sensor.poke(context={})
        
        # Check that the sensor is still poking (not done)
        self.assertFalse(result)
        
        # Check that the hook methods were called correctly
        self.hook_instance.get_quality_results.assert_called_once_with(
            check_id="test-check-id",
        )
    
    def test_data_quality_sensor_completed(self):
        """Test the NessiDataQualitySensor when the check is completed."""
        # Set up the sensor
        sensor = NessiDataQualitySensor(
            task_id="test_task",
            dag=self.dag,
            check_id="test-check-id",
            quality_threshold=0.9,
            conn_id="nessi_test",
        )
        
        # Mock the hook methods to return a completed check
        self.hook_instance.get_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.95,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Execute the sensor's poke method
        result = sensor.poke(context={})
        
        # Check that the sensor is done poking
        self.assertTrue(result)
        
        # Check that the hook methods were called correctly
        self.hook_instance.get_quality_results.assert_called_once_with(
            check_id="test-check-id",
        )
    
    def test_data_quality_sensor_failed_quality(self):
        """Test the NessiDataQualitySensor when the quality threshold is not met."""
        # Set up the sensor
        sensor = NessiDataQualitySensor(
            task_id="test_task",
            dag=self.dag,
            check_id="test-check-id",
            quality_threshold=0.9,
            conn_id="nessi_test",
        )
        
        # Mock the hook methods to return a completed check with low quality
        self.hook_instance.get_quality_results.return_value = {
            "check_id": "test-check-id",
            "status": "completed",
            "quality_score": 0.8,
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Execute the sensor's poke method and expect an exception
        with self.assertRaises(Exception) as context:
            sensor.poke(context={})
        
        # Check the exception message
        self.assertIn("Quality score 0.8 is below threshold 0.9", str(context.exception))
    
    def test_validation_sensor_completed(self):
        """Test the NessiValidationSensor when the validation is completed."""
        # Set up the sensor
        sensor = NessiValidationSensor(
            task_id="test_task",
            dag=self.dag,
            validation_id="test-validation-id",
            conn_id="nessi_test",
        )
        
        # Mock the hook methods to return a completed validation
        self.hook_instance.get_validation_results.return_value = {
            "validation_id": "test-validation-id",
            "status": "completed",
            "rule_results": [
                {"name": "test_rule", "status": "passed"}
            ]
        }
        
        # Execute the sensor's poke method
        result = sensor.poke(context={})
        
        # Check that the sensor is done poking
        self.assertTrue(result)
        
        # Check that the hook methods were called correctly
        self.hook_instance.get_validation_results.assert_called_once_with(
            validation_id="test-validation-id",
        )
    
    def test_profile_sensor_completed(self):
        """Test the NessiProfileSensor when the profile is completed."""
        # Set up the sensor
        sensor = NessiProfileSensor(
            task_id="test_task",
            dag=self.dag,
            profile_id="test-profile-id",
            conn_id="nessi_test",
        )
        
        # Mock the hook methods to return a completed profile
        self.hook_instance.get_profile.return_value = {
            "profile_id": "test-profile-id",
            "status": "completed",
            "columns": [
                {"name": "id", "type": "integer"},
                {"name": "name", "type": "string"},
            ]
        }
        
        # Execute the sensor's poke method
        result = sensor.poke(context={})
        
        # Check that the sensor is done poking
        self.assertTrue(result)
        
        # Check that the hook methods were called correctly
        self.hook_instance.get_profile.assert_called_once_with(
            profile_id="test-profile-id",
        )


if __name__ == "__main__":
    unittest.main()
