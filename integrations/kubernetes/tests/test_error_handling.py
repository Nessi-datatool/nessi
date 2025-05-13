"""
Tests for error handling and retry mechanisms in the Nessi Kubernetes integration.
"""

import os
import json
import unittest
from unittest import mock
import time
import kubernetes
from kubernetes.client.rest import ApiException

from nessi_k8s.operator import NessiOperatorBase, NessiDataQualityOperator
from nessi_k8s.client import NessiK8sClient


class TestErrorHandlingAndRetries(unittest.TestCase):
    """Test cases for error handling and retry mechanisms."""

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
            retries=3,
            retry_backoff=2.0,
            initial_retry_delay=1.0,
            max_retry_delay=30.0,
        )
        
        # Create client instance
        self.client = NessiK8sClient(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            api_key="test-key",
            api_secret="test-secret",
            retries=3,
            retry_backoff=2.0,
            initial_retry_delay=1.0,
            max_retry_delay=30.0,
        )
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_batch_v1_api_patcher.stop()
        self.mock_core_v1_api_patcher.stop()
    
    @mock.patch('time.sleep')
    def test_api_error_retry_success(self, mock_sleep):
        """Test retrying after API errors with eventual success."""
        # Configure mock to fail twice then succeed
        api_exception = ApiException(status=500, reason="Internal Server Error")
        self.mock_batch_v1_api_instance.create_namespaced_job.side_effect = [
            api_exception,
            api_exception,
            mock.MagicMock(metadata=mock.MagicMock(name="test-job")),
        ]
        
        # Configure job status
        self.mock_batch_v1_api_instance.read_namespaced_job.return_value = mock.MagicMock(
            status=mock.MagicMock(
                active=0,
                succeeded=1,
                failed=0,
                completion_time="2023-01-01T00:00:00Z",
            ),
        )
        
        # Configure pod logs
        self.mock_pod_list = mock.MagicMock()
        self.mock_pod_list.items = [
            mock.MagicMock(
                metadata=mock.MagicMock(
                    name="test-job-pod",
                ),
            ),
        ]
        self.mock_core_v1_api_instance.list_namespaced_pod.return_value = self.mock_pod_list
        self.mock_core_v1_api_instance.read_namespaced_pod_log.return_value = "Job completed successfully"
        
        # Run job with retries
        result = self.operator.run_job(
            job_name="test-job",
            command=["nessi", "test"],
            wait_for_completion=True,
        )
        
        # Check result
        self.assertEqual(result["status"], "Succeeded")
        
        # Check that create_namespaced_job was called 3 times (2 failures + 1 success)
        self.assertEqual(self.mock_batch_v1_api_instance.create_namespaced_job.call_count, 3)
        
        # Check that sleep was called with exponential backoff
        mock_sleep.assert_has_calls([
            mock.call(1.0),  # Initial retry delay
            mock.call(2.0),  # Second retry delay (initial * backoff)
        ])
    
    @mock.patch('time.sleep')
    def test_api_error_retry_exhaustion(self, mock_sleep):
        """Test exhausting all retries after API errors."""
        # Configure mock to always fail
        api_exception = ApiException(status=500, reason="Internal Server Error")
        self.mock_batch_v1_api_instance.create_namespaced_job.side_effect = api_exception
        
        # Run job with retries
        with self.assertRaises(ApiException):
            self.operator.run_job(
                job_name="test-job",
                command=["nessi", "test"],
                wait_for_completion=True,
            )
        
        # Check that create_namespaced_job was called 4 times (initial + 3 retries)
        self.assertEqual(self.mock_batch_v1_api_instance.create_namespaced_job.call_count, 4)
        
        # Check that sleep was called with exponential backoff
        mock_sleep.assert_has_calls([
            mock.call(1.0),  # Initial retry delay
            mock.call(2.0),  # Second retry delay (initial * backoff)
            mock.call(4.0),  # Third retry delay (second * backoff)
        ])
    
    def test_error_categorization(self):
        """Test error categorization for different error types."""
        # API error
        api_error = ApiException(status=500, reason="Internal Server Error")
        error_info = self.operator.categorize_error(api_error)
        self.assertEqual(error_info["error_type"], "ApiError")
        self.assertEqual(error_info["status_code"], 500)
        self.assertEqual(error_info["reason"], "Internal Server Error")
        
        # Authentication error
        auth_error = ApiException(status=401, reason="Unauthorized")
        error_info = self.operator.categorize_error(auth_error)
        self.assertEqual(error_info["error_type"], "AuthenticationError")
        self.assertEqual(error_info["status_code"], 401)
        
        # Resource error
        resource_error = ApiException(status=409, reason="Conflict")
        error_info = self.operator.categorize_error(resource_error)
        self.assertEqual(error_info["error_type"], "ResourceError")
        self.assertEqual(error_info["status_code"], 409)
        
        # Timeout error
        timeout_error = TimeoutError("Operation timed out")
        error_info = self.operator.categorize_error(timeout_error)
        self.assertEqual(error_info["error_type"], "TimeoutError")
        self.assertEqual(error_info["message"], "Operation timed out")
        
        # Generic error
        generic_error = Exception("Unknown error")
        error_info = self.operator.categorize_error(generic_error)
        self.assertEqual(error_info["error_type"], "UnknownError")
        self.assertEqual(error_info["message"], "Unknown error")
    
    @mock.patch('time.sleep')
    def test_timeout_handling(self, mock_sleep):
        """Test handling of timeout errors."""
        # Configure mock to succeed job creation
        self.mock_batch_v1_api_instance.create_namespaced_job.return_value = mock.MagicMock(
            metadata=mock.MagicMock(name="test-job"),
        )
        
        # Configure job status to always be active (never completes)
        self.mock_batch_v1_api_instance.read_namespaced_job.return_value = mock.MagicMock(
            status=mock.MagicMock(
                active=1,
                succeeded=0,
                failed=0,
                completion_time=None,
            ),
        )
        
        # Run job with a short timeout
        result = self.operator.run_job(
            job_name="test-job",
            command=["nessi", "test"],
            wait_for_completion=True,
            timeout=2,
            poll_interval=1,
        )
        
        # Check result
        self.assertEqual(result["status"], "Timeout")
        self.assertEqual(result["job_name"], "test-job")
        self.assertIn("timeout", result["message"].lower())
    
    def test_structured_logging(self):
        """Test structured logging for operations."""
        # Mock the logger
        mock_logger = mock.MagicMock()
        self.operator.logger = mock_logger
        
        # Perform an operation that logs
        try:
            # Trigger an API exception
            api_exception = ApiException(status=500, reason="Internal Server Error")
            self.mock_batch_v1_api_instance.create_namespaced_job.side_effect = api_exception
            
            self.operator.run_job(
                job_name="test-job",
                command=["nessi", "test"],
                wait_for_completion=False,
            )
        except ApiException:
            pass
        
        # Check that error was logged with structured data
        mock_logger.error.assert_called()
        
        # Get the call arguments
        args, kwargs = mock_logger.error.call_args
        
        # Check that extra data was included
        self.assertIn("extra", kwargs)
        extra = kwargs["extra"]
        
        # Check that structured data was included
        self.assertIn("error_type", extra)
        self.assertIn("status_code", extra)
        self.assertIn("job_name", extra)
        self.assertIn("namespace", extra)
        self.assertIn("timestamp", extra)


class TestClientErrorHandling(unittest.TestCase):
    """Test cases for client-level error handling."""

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
            "retries": 3,
            "retry_backoff": 2.0,
            "initial_retry_delay": 1.0,
            "max_retry_delay": 30.0,
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
    
    def test_quality_check_error_handling(self):
        """Test error handling during quality check."""
        # Configure mock to raise an exception
        api_exception = ApiException(status=500, reason="Internal Server Error")
        self.mock_quality_operator_instance.run_quality_check.side_effect = api_exception
        
        # Run quality check with error handling
        with self.assertRaises(ApiException):
            self.client.run_quality_check(
                table_name="test-table",
                rules=[{"name": "test-rule", "rule_type": "not_null", "column": "id"}],
            )
        
        # Check that the operator was called
        self.mock_quality_operator_instance.run_quality_check.assert_called_once()
    
    def test_oauth_authentication(self):
        """Test OAuth 2.0 authentication support."""
        # Create client with OAuth config
        client = NessiK8sClient(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            auth_type="oauth2",
            oauth_token="test-token",
            oauth_refresh_token="test-refresh-token",
            oauth_token_url="http://test-token-url",
            oauth_client_id="test-client-id",
            oauth_client_secret="test-client-secret",
        )
        
        # Check that OAuth config was set
        self.assertEqual(client.config["auth_type"], "oauth2")
        self.assertEqual(client.config["oauth_token"], "test-token")
        self.assertEqual(client.config["oauth_refresh_token"], "test-refresh-token")
        self.assertEqual(client.config["oauth_token_url"], "http://test-token-url")
        self.assertEqual(client.config["oauth_client_id"], "test-client-id")
        self.assertEqual(client.config["oauth_client_secret"], "test-client-secret")
        
        # Check that operators were created with OAuth config
        self.mock_quality_operator.assert_called_once()
        quality_operator_kwargs = self.mock_quality_operator.call_args[1]
        self.assertEqual(quality_operator_kwargs["auth_type"], "oauth2")
        self.assertEqual(quality_operator_kwargs["oauth_token"], "test-token")
    
    @mock.patch('nessi_k8s.client.refresh_oauth_token')
    def test_token_refresh(self, mock_refresh_token):
        """Test token refresh capability."""
        # Configure mock to return new tokens
        mock_refresh_token.return_value = {
            "access_token": "new-token",
            "refresh_token": "new-refresh-token",
            "expires_in": 3600,
        }
        
        # Create client with OAuth config
        client = NessiK8sClient(
            namespace="test-namespace",
            image="test-image:latest",
            api_host="http://test-host",
            auth_type="oauth2",
            oauth_token="expired-token",
            oauth_refresh_token="test-refresh-token",
            oauth_token_url="http://test-token-url",
            oauth_client_id="test-client-id",
            oauth_client_secret="test-client-secret",
        )
        
        # Simulate token refresh
        client.refresh_token()
        
        # Check that tokens were updated
        self.assertEqual(client.config["oauth_token"], "new-token")
        self.assertEqual(client.config["oauth_refresh_token"], "new-refresh-token")
        
        # Check that refresh_oauth_token was called with correct parameters
        mock_refresh_token.assert_called_once_with(
            token_url="http://test-token-url",
            client_id="test-client-id",
            client_secret="test-client-secret",
            refresh_token="test-refresh-token",
        )


if __name__ == "__main__":
    unittest.main()
