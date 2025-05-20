"""
Tests for the Nessi Kubernetes client module.
"""

import os
import json
import unittest
from unittest import mock

from nessi_k8s.client import NessiK8sClient


class TestNessiK8sClient(unittest.TestCase):
    """Test cases for the NessiK8sClient class."""

    def setUp(self):
        """Set up test environment."""
        # Mock the operators
        self.mock_quality_operator_patcher = mock.patch('nessi_k8s.client.NessiDataQualityOperator')
        self.mock_quality_operator = self.mock_quality_operator_patcher.start()
        self.mock_quality_operator_instance = self.mock_quality_operator.return_value
        
        self.mock_profile_operator_patcher = mock.patch('nessi_k8s.client.NessiProfileOperator')
        self.mock_profile_operator = self.mock_profile_operator_patcher.start()
        self.mock_profile_operator_instance = self.mock_profile_operator.return_value
        
        self.mock_validation_operator_patcher = mock.patch('nessi_k8s.client.NessiValidationOperator')
        self.mock_validation_operator = self.mock_validation_operator_patcher.start()
        self.mock_validation_operator_instance = self.mock_validation_operator.return_value
        
        # Mock the load_config function
        self.mock_load_config_patcher = mock.patch('nessi_k8s.client.load_config')
        self.mock_load_config = self.mock_load_config_patcher.start()
        self.mock_load_config.return_value = {
            "namespace": "default",
            "image": "nessi/nessi:latest",
            "service_account_name": None,
            "api_host": None,
            "api_key": None,
            "api_secret": None,
        }
        
        # Create client instance
        self.client = NessiK8sClient(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            api_key="test-key",
            api_secret="test-secret",
        )
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_quality_operator_patcher.stop()
        self.mock_profile_operator_patcher.stop()
        self.mock_validation_operator_patcher.stop()
        self.mock_load_config_patcher.stop()
    
    def test_init(self):
        """Test initialization of client."""
        # Check that load_config was called
        self.mock_load_config.assert_called_once_with(None)
        
        # Check that config was updated with provided values
        self.assertEqual(self.client.config["namespace"], "test-namespace")
        self.assertEqual(self.client.config["image"], "test-image:latest")
        self.assertEqual(self.client.config["api_host"], "http://test-host")
        self.assertEqual(self.client.config["api_key"], "test-key")
        self.assertEqual(self.client.config["api_secret"], "test-secret")
        
        # Check that operators were created
        self.mock_quality_operator.assert_called_once()
        self.mock_profile_operator.assert_called_once()
        self.mock_validation_operator.assert_called_once()
        
        # Check that operator instances were assigned
        self.assertEqual(self.client.quality_operator, self.mock_quality_operator_instance)
        self.assertEqual(self.client.profile_operator, self.mock_profile_operator_instance)
        self.assertEqual(self.client.validation_operator, self.mock_validation_operator_instance)
    
    def test_init_with_config_file(self):
        """Test initialization of client with config file."""
        # Reset mocks
        self.mock_load_config.reset_mock()
        self.mock_quality_operator.reset_mock()
        self.mock_profile_operator.reset_mock()
        self.mock_validation_operator.reset_mock()
        
        # Create client with config file
        client = NessiK8sClient(config_file="test-config.yaml")
        
        # Check that load_config was called with config file
        self.mock_load_config.assert_called_once_with("test-config.yaml")
        
        # Check that operators were created
        self.mock_quality_operator.assert_called_once()
        self.mock_profile_operator.assert_called_once()
        self.mock_validation_operator.assert_called_once()
    
    def test_run_quality_check(self):
        """Test running a data quality check."""
        # Mock operator response
        self.mock_quality_operator_instance.run_quality_check.return_value = {
            "status": "Succeeded",
            "job_name": "test-quality-job",
            "logs": "Quality check results",
        }
        
        # Run quality check
        result = self.client.run_quality_check(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            profile=True,
            output_format="json",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
            job_name="test-job",
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-quality-job")
        self.assertEqual(result["logs"], "Quality check results")
        
        # Check that operator method was called correctly
        self.mock_quality_operator_instance.run_quality_check.assert_called_once_with(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            profile=True,
            output_format="json",
            output_path=None,
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
            job_name="test-job",
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
    
    def test_get_quality_results(self):
        """Test getting data quality results."""
        # Mock operator response
        self.mock_quality_operator_instance.get_quality_results.return_value = {
            "status": "Succeeded",
            "job_name": "test-quality-job",
            "logs": "Quality check results",
        }
        
        # Get quality results
        result = self.client.get_quality_results(
            job_name="test-quality-job",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-quality-job")
        self.assertEqual(result["logs"], "Quality check results")
        
        # Check that operator method was called correctly
        self.mock_quality_operator_instance.get_quality_results.assert_called_once_with(
            job_name="test-quality-job",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
        )
    
    def test_run_profile(self):
        """Test running a data profile."""
        # Mock operator response
        self.mock_profile_operator_instance.run_profile.return_value = {
            "status": "Succeeded",
            "job_name": "test-profile-job",
            "logs": "Profile results",
        }
        
        # Run profile
        result = self.client.run_profile(
            table_name="test-table",
            columns=["id", "name"],
            sample_size=1000,
            output_format="json",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
            job_name="test-job",
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-profile-job")
        self.assertEqual(result["logs"], "Profile results")
        
        # Check that operator method was called correctly
        self.mock_profile_operator_instance.run_profile.assert_called_once_with(
            table_name="test-table",
            columns=["id", "name"],
            sample_size=1000,
            output_format="json",
            output_path=None,
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
            job_name="test-job",
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
    
    def test_get_profile_results(self):
        """Test getting data profile results."""
        # Mock operator response
        self.mock_profile_operator_instance.get_profile_results.return_value = {
            "status": "Succeeded",
            "job_name": "test-profile-job",
            "logs": "Profile results",
        }
        
        # Get profile results
        result = self.client.get_profile_results(
            job_name="test-profile-job",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-profile-job")
        self.assertEqual(result["logs"], "Profile results")
        
        # Check that operator method was called correctly
        self.mock_profile_operator_instance.get_profile_results.assert_called_once_with(
            job_name="test-profile-job",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
        )
    
    def test_run_validation(self):
        """Test running a data validation."""
        # Mock operator response
        self.mock_validation_operator_instance.run_validation.return_value = {
            "status": "Succeeded",
            "job_name": "test-validation-job",
            "logs": "Validation results",
        }
        
        # Run validation
        result = self.client.run_validation(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            output_format="json",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
            job_name="test-job",
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-validation-job")
        self.assertEqual(result["logs"], "Validation results")
        
        # Check that operator method was called correctly
        self.mock_validation_operator_instance.run_validation.assert_called_once_with(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            output_format="json",
            output_path=None,
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
            job_name="test-job",
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
    
    def test_get_validation_results(self):
        """Test getting data validation results."""
        # Mock operator response
        self.mock_validation_operator_instance.get_validation_results.return_value = {
            "status": "Succeeded",
            "job_name": "test-validation-job",
            "logs": "Validation results",
        }
        
        # Get validation results
        result = self.client.get_validation_results(
            job_name="test-validation-job",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-validation-job")
        self.assertEqual(result["logs"], "Validation results")
        
        # Check that operator method was called correctly
        self.mock_validation_operator_instance.get_validation_results.assert_called_once_with(
            job_name="test-validation-job",
            wait_for_completion=True,
            timeout=600,
            poll_interval=10,
        )


if __name__ == "__main__":
    unittest.main()
