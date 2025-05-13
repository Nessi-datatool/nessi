"""
Tests for the Nessi Kubernetes operator module.
"""

import os
import json
import uuid
import unittest
from unittest import mock
import kubernetes
from kubernetes.client.rest import ApiException

from nessi_k8s.operator import (
    NessiOperatorBase,
    NessiDataQualityOperator,
    NessiProfileOperator,
    NessiValidationOperator,
)


class TestNessiOperatorBase(unittest.TestCase):
    """Test cases for the NessiOperatorBase class."""

    def setUp(self):
        """Set up test environment."""
        # Mock Kubernetes client
        self.mock_batch_v1_api_patcher = mock.patch('kubernetes.client.BatchV1Api')
        self.mock_batch_v1_api = self.mock_batch_v1_api_patcher.start()
        self.mock_batch_v1_api_instance = self.mock_batch_v1_api.return_value
        
        self.mock_core_v1_api_patcher = mock.patch('kubernetes.client.CoreV1Api')
        self.mock_core_v1_api = self.mock_core_v1_api_patcher.start()
        self.mock_core_v1_api_instance = self.mock_core_v1_api.return_value
        
        # Create operator instance
        self.operator = NessiOperatorBase(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            api_key="test-key",
            api_secret="test-secret",
        )
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_batch_v1_api_patcher.stop()
        self.mock_core_v1_api_patcher.stop()
    
    def test_init(self):
        """Test initialization of operator."""
        # Check that attributes are set correctly
        self.assertEqual(self.operator.namespace, "test-namespace")
        self.assertEqual(self.operator.image, "test-image:latest")
        self.assertEqual(self.operator.api_host, "http://test-host")
        self.assertEqual(self.operator.api_key, "test-key")
        self.assertEqual(self.operator.api_secret, "test-secret")
        
        # Check that Kubernetes clients are created
        self.assertEqual(self.operator.batch_v1_api, self.mock_batch_v1_api_instance)
        self.assertEqual(self.operator.core_v1_api, self.mock_core_v1_api_instance)
    
    def test_generate_job_name(self):
        """Test generation of job names."""
        # Test with prefix
        job_name = self.operator._generate_job_name(prefix="test")
        self.assertTrue(job_name.startswith("test-"))
        self.assertTrue(len(job_name) > 5)  # Prefix + some random chars
        
        # Test with different prefix
        job_name = self.operator._generate_job_name(prefix="other")
        self.assertTrue(job_name.startswith("other-"))
        
        # Test with no prefix
        job_name = self.operator._generate_job_name()
        self.assertTrue(job_name.startswith("nessi-"))
    
    def test_create_job_object(self):
        """Test creation of Kubernetes job object."""
        # Create job object
        job = self.operator._create_job_object(
            job_name="test-job",
            command=["nessi", "test"],
            env_vars=[{"name": "TEST_VAR", "value": "test-value"}],
        )
        
        # Check job metadata
        self.assertEqual(job.metadata.name, "test-job")
        self.assertEqual(job.metadata.namespace, "test-namespace")
        self.assertEqual(job.metadata.labels["app"], "nessi")
        
        # Check job spec
        self.assertEqual(job.spec.template.spec.containers[0].name, "nessi")
        self.assertEqual(job.spec.template.spec.containers[0].image, "test-image:latest")
        self.assertEqual(job.spec.template.spec.containers[0].command, ["nessi", "test"])
        
        # Check environment variables
        env = job.spec.template.spec.containers[0].env
        self.assertEqual(len(env), 4)  # 3 default + 1 custom
        
        env_names = [e.name for e in env]
        self.assertIn("NESSI_API_HOST", env_names)
        self.assertIn("NESSI_API_KEY", env_names)
        self.assertIn("NESSI_API_SECRET", env_names)
        self.assertIn("TEST_VAR", env_names)
        
        # Find the custom env var and check its value
        test_var = next(e for e in env if e.name == "TEST_VAR")
        self.assertEqual(test_var.value, "test-value")
    
    @mock.patch('time.sleep', return_value=None)
    def test_wait_for_job_completion_success(self, mock_sleep):
        """Test waiting for job completion with success."""
        # Mock job status responses
        job_responses = [
            # First response: job is running
            mock.MagicMock(
                status=mock.MagicMock(
                    active=1,
                    succeeded=0,
                    failed=0,
                    completion_time=None,
                )
            ),
            # Second response: job is still running
            mock.MagicMock(
                status=mock.MagicMock(
                    active=1,
                    succeeded=0,
                    failed=0,
                    completion_time=None,
                )
            ),
            # Third response: job has succeeded
            mock.MagicMock(
                status=mock.MagicMock(
                    active=0,
                    succeeded=1,
                    failed=0,
                    completion_time="2023-01-01T00:00:00Z",
                )
            ),
        ]
        
        self.mock_batch_v1_api_instance.read_namespaced_job.side_effect = job_responses
        
        # Wait for job completion
        result = self.operator._wait_for_job_completion(
            job_name="test-job",
            timeout=10,
            poll_interval=1,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-job")
        
        # Check that read_namespaced_job was called multiple times
        self.assertEqual(self.mock_batch_v1_api_instance.read_namespaced_job.call_count, 3)
    
    @mock.patch('time.sleep', return_value=None)
    def test_wait_for_job_completion_failure(self, mock_sleep):
        """Test waiting for job completion with failure."""
        # Mock job status responses
        job_responses = [
            # First response: job is running
            mock.MagicMock(
                status=mock.MagicMock(
                    active=1,
                    succeeded=0,
                    failed=0,
                    completion_time=None,
                )
            ),
            # Second response: job has failed
            mock.MagicMock(
                status=mock.MagicMock(
                    active=0,
                    succeeded=0,
                    failed=1,
                    completion_time="2023-01-01T00:00:00Z",
                )
            ),
        ]
        
        self.mock_batch_v1_api_instance.read_namespaced_job.side_effect = job_responses
        
        # Wait for job completion
        result = self.operator._wait_for_job_completion(
            job_name="test-job",
            timeout=10,
            poll_interval=1,
        )
        
        # Check result
        self.assertEqual(result["status"], "Failed")
        self.assertEqual(result["job_name"], "test-job")
        
        # Check that read_namespaced_job was called multiple times
        self.assertEqual(self.mock_batch_v1_api_instance.read_namespaced_job.call_count, 2)
    
    @mock.patch('time.sleep', return_value=None)
    def test_wait_for_job_completion_timeout(self, mock_sleep):
        """Test waiting for job completion with timeout."""
        # Mock job status response - job is always running
        job_response = mock.MagicMock(
            status=mock.MagicMock(
                active=1,
                succeeded=0,
                failed=0,
                completion_time=None,
            )
        )
        
        self.mock_batch_v1_api_instance.read_namespaced_job.return_value = job_response
        
        # Wait for job completion with a short timeout
        result = self.operator._wait_for_job_completion(
            job_name="test-job",
            timeout=2,
            poll_interval=1,
        )
        
        # Check result
        self.assertEqual(result["status"], "Running")
        self.assertEqual(result["job_name"], "test-job")
        self.assertIn("timeout", result["message"].lower())
        
        # Check that read_namespaced_job was called multiple times
        self.assertTrue(self.mock_batch_v1_api_instance.read_namespaced_job.call_count >= 2)
    
    def test_get_pod_logs(self):
        """Test getting pod logs."""
        # Mock pod list response
        pod_list = mock.MagicMock()
        pod_list.items = [
            mock.MagicMock(
                metadata=mock.MagicMock(
                    name="test-job-pod",
                )
            ),
        ]
        
        self.mock_core_v1_api_instance.list_namespaced_pod.return_value = pod_list
        
        # Mock pod logs response
        self.mock_core_v1_api_instance.read_namespaced_pod_log.return_value = "Test logs"
        
        # Get pod logs
        logs = self.operator._get_pod_logs(job_name="test-job")
        
        # Check logs
        self.assertEqual(logs, "Test logs")
        
        # Check that list_namespaced_pod was called correctly
        self.mock_core_v1_api_instance.list_namespaced_pod.assert_called_once_with(
            namespace="test-namespace",
            label_selector="job-name=test-job",
        )
        
        # Check that read_namespaced_pod_log was called correctly
        self.mock_core_v1_api_instance.read_namespaced_pod_log.assert_called_once_with(
            name="test-job-pod",
            namespace="test-namespace",
        )
    
    def test_get_pod_logs_no_pods(self):
        """Test getting pod logs when no pods are found."""
        # Mock pod list response with no pods
        pod_list = mock.MagicMock()
        pod_list.items = []
        
        self.mock_core_v1_api_instance.list_namespaced_pod.return_value = pod_list
        
        # Get pod logs
        logs = self.operator._get_pod_logs(job_name="test-job")
        
        # Check logs
        self.assertEqual(logs, "No pods found for job test-job")
        
        # Check that list_namespaced_pod was called correctly
        self.mock_core_v1_api_instance.list_namespaced_pod.assert_called_once_with(
            namespace="test-namespace",
            label_selector="job-name=test-job",
        )
        
        # Check that read_namespaced_pod_log was not called
        self.mock_core_v1_api_instance.read_namespaced_pod_log.assert_not_called()


class TestNessiDataQualityOperator(unittest.TestCase):
    """Test cases for the NessiDataQualityOperator class."""

    def setUp(self):
        """Set up test environment."""
        # Mock Kubernetes client
        self.mock_batch_v1_api_patcher = mock.patch('kubernetes.client.BatchV1Api')
        self.mock_batch_v1_api = self.mock_batch_v1_api_patcher.start()
        self.mock_batch_v1_api_instance = self.mock_batch_v1_api.return_value
        
        self.mock_core_v1_api_patcher = mock.patch('kubernetes.client.CoreV1Api')
        self.mock_core_v1_api = self.mock_core_v1_api_patcher.start()
        self.mock_core_v1_api_instance = self.mock_core_v1_api.return_value
        
        # Create operator instance
        self.operator = NessiDataQualityOperator(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            api_key="test-key",
            api_secret="test-secret",
        )
        
        # Mock job creation response
        self.mock_job_response = mock.MagicMock(
            metadata=mock.MagicMock(
                name="test-quality-job",
            ),
        )
        self.mock_batch_v1_api_instance.create_namespaced_job.return_value = self.mock_job_response
        
        # Mock job status response
        self.mock_job_status = mock.MagicMock(
            status=mock.MagicMock(
                active=0,
                succeeded=1,
                failed=0,
                completion_time="2023-01-01T00:00:00Z",
            ),
        )
        self.mock_batch_v1_api_instance.read_namespaced_job.return_value = self.mock_job_status
        
        # Mock pod logs response
        self.mock_pod_list = mock.MagicMock()
        self.mock_pod_list.items = [
            mock.MagicMock(
                metadata=mock.MagicMock(
                    name="test-quality-job-pod",
                ),
            ),
        ]
        self.mock_core_v1_api_instance.list_namespaced_pod.return_value = self.mock_pod_list
        self.mock_core_v1_api_instance.read_namespaced_pod_log.return_value = "Quality check results"
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_batch_v1_api_patcher.stop()
        self.mock_core_v1_api_patcher.stop()
    
    def test_run_quality_check(self):
        """Test running a data quality check."""
        # Run quality check
        result = self.operator.run_quality_check(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            profile=True,
            output_format="json",
            wait_for_completion=True,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-quality-job")
        self.assertEqual(result["logs"], "Quality check results")
        
        # Check that create_namespaced_job was called correctly
        self.mock_batch_v1_api_instance.create_namespaced_job.assert_called_once()
        
        # Check job arguments
        job_arg = self.mock_batch_v1_api_instance.create_namespaced_job.call_args[0][1]
        self.assertEqual(job_arg.metadata.namespace, "test-namespace")
        
        # Check container command
        container = job_arg.spec.template.spec.containers[0]
        self.assertEqual(container.image, "test-image:latest")
        
        # Command should include 'quality check' and the table name
        command = container.command
        self.assertIn("quality", command)
        self.assertIn("check", command)
        self.assertIn("test-table", command)
        
        # Check that wait_for_job_completion was called
        self.mock_batch_v1_api_instance.read_namespaced_job.assert_called_once_with(
            name="test-quality-job",
            namespace="test-namespace",
        )
    
    def test_run_quality_check_async(self):
        """Test running a data quality check asynchronously."""
        # Run quality check without waiting
        result = self.operator.run_quality_check(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            profile=True,
            output_format="json",
            wait_for_completion=False,
        )
        
        # Check result
        self.assertEqual(result["name"], "test-quality-job")
        self.assertEqual(result["namespace"], "test-namespace")
        
        # Check that create_namespaced_job was called correctly
        self.mock_batch_v1_api_instance.create_namespaced_job.assert_called_once()
        
        # Check that wait_for_job_completion was not called
        self.mock_batch_v1_api_instance.read_namespaced_job.assert_not_called()
    
    def test_get_quality_results(self):
        """Test getting data quality results."""
        # Get quality results
        result = self.operator.get_quality_results(
            job_name="test-quality-job",
            wait_for_completion=True,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-quality-job")
        self.assertEqual(result["logs"], "Quality check results")
        
        # Check that read_namespaced_job was called correctly
        self.mock_batch_v1_api_instance.read_namespaced_job.assert_called_once_with(
            name="test-quality-job",
            namespace="test-namespace",
        )
        
        # Check that list_namespaced_pod was called correctly
        self.mock_core_v1_api_instance.list_namespaced_pod.assert_called_once_with(
            namespace="test-namespace",
            label_selector="job-name=test-quality-job",
        )
        
        # Check that read_namespaced_pod_log was called correctly
        self.mock_core_v1_api_instance.read_namespaced_pod_log.assert_called_once_with(
            name="test-quality-job-pod",
            namespace="test-namespace",
        )


class TestNessiProfileOperator(unittest.TestCase):
    """Test cases for the NessiProfileOperator class."""

    def setUp(self):
        """Set up test environment."""
        # Mock Kubernetes client
        self.mock_batch_v1_api_patcher = mock.patch('kubernetes.client.BatchV1Api')
        self.mock_batch_v1_api = self.mock_batch_v1_api_patcher.start()
        self.mock_batch_v1_api_instance = self.mock_batch_v1_api.return_value
        
        self.mock_core_v1_api_patcher = mock.patch('kubernetes.client.CoreV1Api')
        self.mock_core_v1_api = self.mock_core_v1_api_patcher.start()
        self.mock_core_v1_api_instance = self.mock_core_v1_api.return_value
        
        # Create operator instance
        self.operator = NessiProfileOperator(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            api_key="test-key",
            api_secret="test-secret",
        )
        
        # Mock job creation response
        self.mock_job_response = mock.MagicMock(
            metadata=mock.MagicMock(
                name="test-profile-job",
            ),
        )
        self.mock_batch_v1_api_instance.create_namespaced_job.return_value = self.mock_job_response
        
        # Mock job status response
        self.mock_job_status = mock.MagicMock(
            status=mock.MagicMock(
                active=0,
                succeeded=1,
                failed=0,
                completion_time="2023-01-01T00:00:00Z",
            ),
        )
        self.mock_batch_v1_api_instance.read_namespaced_job.return_value = self.mock_job_status
        
        # Mock pod logs response
        self.mock_pod_list = mock.MagicMock()
        self.mock_pod_list.items = [
            mock.MagicMock(
                metadata=mock.MagicMock(
                    name="test-profile-job-pod",
                ),
            ),
        ]
        self.mock_core_v1_api_instance.list_namespaced_pod.return_value = self.mock_pod_list
        self.mock_core_v1_api_instance.read_namespaced_pod_log.return_value = "Profile results"
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_batch_v1_api_patcher.stop()
        self.mock_core_v1_api_patcher.stop()
    
    def test_run_profile(self):
        """Test running a data profile."""
        # Run profile
        result = self.operator.run_profile(
            table_name="test-table",
            columns=["id", "name"],
            sample_size=1000,
            output_format="json",
            wait_for_completion=True,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-profile-job")
        self.assertEqual(result["logs"], "Profile results")
        
        # Check that create_namespaced_job was called correctly
        self.mock_batch_v1_api_instance.create_namespaced_job.assert_called_once()
        
        # Check job arguments
        job_arg = self.mock_batch_v1_api_instance.create_namespaced_job.call_args[0][1]
        self.assertEqual(job_arg.metadata.namespace, "test-namespace")
        
        # Check container command
        container = job_arg.spec.template.spec.containers[0]
        self.assertEqual(container.image, "test-image:latest")
        
        # Command should include 'profile' and the table name
        command = container.command
        self.assertIn("profile", command)
        self.assertIn("test-table", command)
        
        # Check that columns and sample_size are included
        command_str = " ".join(command)
        self.assertIn("--columns", command_str)
        self.assertIn("id,name", command_str)
        self.assertIn("--sample-size", command_str)
        self.assertIn("1000", command_str)


class TestNessiValidationOperator(unittest.TestCase):
    """Test cases for the NessiValidationOperator class."""

    def setUp(self):
        """Set up test environment."""
        # Mock Kubernetes client
        self.mock_batch_v1_api_patcher = mock.patch('kubernetes.client.BatchV1Api')
        self.mock_batch_v1_api = self.mock_batch_v1_api_patcher.start()
        self.mock_batch_v1_api_instance = self.mock_batch_v1_api.return_value
        
        self.mock_core_v1_api_patcher = mock.patch('kubernetes.client.CoreV1Api')
        self.mock_core_v1_api = self.mock_core_v1_api_patcher.start()
        self.mock_core_v1_api_instance = self.mock_core_v1_api.return_value
        
        # Create operator instance
        self.operator = NessiValidationOperator(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            api_key="test-key",
            api_secret="test-secret",
        )
        
        # Mock job creation response
        self.mock_job_response = mock.MagicMock(
            metadata=mock.MagicMock(
                name="test-validation-job",
            ),
        )
        self.mock_batch_v1_api_instance.create_namespaced_job.return_value = self.mock_job_response
        
        # Mock job status response
        self.mock_job_status = mock.MagicMock(
            status=mock.MagicMock(
                active=0,
                succeeded=1,
                failed=0,
                completion_time="2023-01-01T00:00:00Z",
            ),
        )
        self.mock_batch_v1_api_instance.read_namespaced_job.return_value = self.mock_job_status
        
        # Mock pod logs response
        self.mock_pod_list = mock.MagicMock()
        self.mock_pod_list.items = [
            mock.MagicMock(
                metadata=mock.MagicMock(
                    name="test-validation-job-pod",
                ),
            ),
        ]
        self.mock_core_v1_api_instance.list_namespaced_pod.return_value = self.mock_pod_list
        self.mock_core_v1_api_instance.read_namespaced_pod_log.return_value = "Validation results"
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_batch_v1_api_patcher.stop()
        self.mock_core_v1_api_patcher.stop()
    
    def test_run_validation(self):
        """Test running a data validation."""
        # Run validation
        result = self.operator.run_validation(
            table_name="test-table",
            rules=[
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ],
            output_format="json",
            wait_for_completion=True,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        self.assertEqual(result["job_name"], "test-validation-job")
        self.assertEqual(result["logs"], "Validation results")
        
        # Check that create_namespaced_job was called correctly
        self.mock_batch_v1_api_instance.create_namespaced_job.assert_called_once()
        
        # Check job arguments
        job_arg = self.mock_batch_v1_api_instance.create_namespaced_job.call_args[0][1]
        self.assertEqual(job_arg.metadata.namespace, "test-namespace")
        
        # Check container command
        container = job_arg.spec.template.spec.containers[0]
        self.assertEqual(container.image, "test-image:latest")
        
        # Command should include 'validate' and the table name
        command = container.command
        self.assertIn("validate", command)
        self.assertIn("test-table", command)


if __name__ == "__main__":
    unittest.main()
