"""
Tests for the error handling utilities in the Nessi Airflow integration.
"""
import unittest
from unittest import mock
import time
import json
import requests

from airflow.exceptions import AirflowException

from nessi_airflow.utils.error_handling import (
    NessiApiError,
    NessiAuthenticationError,
    NessiResourceNotFoundError,
    NessiServerError,
    NessiClientError,
    NessiTimeoutError,
    NessiConnectionError,
    with_exponential_backoff,
    handle_api_exceptions,
    map_http_error,
)
from nessi_airflow.hooks.nessi_hook import NessiHook


class TestErrorHandling(unittest.TestCase):
    """Test the error handling utilities."""

    def setUp(self):
        """Set up the test case."""
        self.conn_id = "nessi_test"
        self.mock_connection = mock.patch(
            "airflow.hooks.base.BaseHook.get_connection",
            return_value=mock.MagicMock(
                conn_id=self.conn_id,
                conn_type="http",
                host="http://localhost:8080",
                login="test-api-key",
                password="test-api-secret",
            ),
        )
        self.mock_connection.start()
        
        # Create a hook
        self.hook = NessiHook(conn_id=self.conn_id)
    
    def tearDown(self):
        """Tear down the test case."""
        self.mock_connection.stop()
    
    def test_map_http_error(self):
        """Test mapping HTTP errors to exception types."""
        # Authentication errors
        self.assertEqual(map_http_error(401, {}), NessiAuthenticationError)
        self.assertEqual(map_http_error(403, {}), NessiAuthenticationError)
        
        # Resource not found
        self.assertEqual(map_http_error(404, {}), NessiResourceNotFoundError)
        
        # Client errors
        self.assertEqual(map_http_error(400, {}), NessiClientError)
        self.assertEqual(map_http_error(422, {}), NessiClientError)
        
        # Server errors
        self.assertEqual(map_http_error(500, {}), NessiServerError)
        self.assertEqual(map_http_error(502, {}), NessiServerError)
        self.assertEqual(map_http_error(503, {}), NessiServerError)
        
        # Other errors
        self.assertEqual(map_http_error(300, {}), NessiApiError)
    
    @mock.patch("requests.Session.get")
    def test_handle_api_exceptions_timeout(self, mock_get):
        """Test handling timeout exceptions."""
        # Mock a timeout exception
        mock_get.side_effect = requests.exceptions.Timeout("Connection timed out")
        
        # Create a test function
        @handle_api_exceptions
        def test_func():
            return mock_get()
        
        # Test that the exception is properly handled
        with self.assertRaises(NessiTimeoutError):
            test_func()
    
    @mock.patch("requests.Session.get")
    def test_handle_api_exceptions_connection_error(self, mock_get):
        """Test handling connection errors."""
        # Mock a connection error
        mock_get.side_effect = requests.exceptions.ConnectionError("Connection refused")
        
        # Create a test function
        @handle_api_exceptions
        def test_func():
            return mock_get()
        
        # Test that the exception is properly handled
        with self.assertRaises(NessiConnectionError):
            test_func()
    
    @mock.patch("requests.Session.get")
    def test_handle_api_exceptions_http_error(self, mock_get):
        """Test handling HTTP errors."""
        # Mock an HTTP error response
        mock_response = mock.Mock()
        mock_response.status_code = 404
        mock_response.json.return_value = {"error": "Resource not found"}
        mock_response.raise_for_status.side_effect = requests.exceptions.HTTPError(
            "404 Client Error", response=mock_response
        )
        mock_get.return_value = mock_response
        
        # Create a test function
        @handle_api_exceptions
        def test_func():
            response = mock_get()
            response.raise_for_status()
            return response
        
        # Test that the exception is properly handled
        with self.assertRaises(NessiResourceNotFoundError) as context:
            test_func()
        
        # Check the exception details
        self.assertEqual(context.exception.status_code, 404)
        self.assertEqual(context.exception.response, {"error": "Resource not found"})
    
    def test_with_exponential_backoff_success(self):
        """Test successful execution with exponential backoff."""
        mock_func = mock.Mock(return_value="success")
        
        # Create a decorated function
        decorated_func = with_exponential_backoff()(mock_func)
        
        # Call the function
        result = decorated_func()
        
        # Check that the function was called once and returned the expected result
        mock_func.assert_called_once()
        self.assertEqual(result, "success")
    
    def test_with_exponential_backoff_retry_success(self):
        """Test successful retry with exponential backoff."""
        # Mock a function that fails twice and then succeeds
        side_effects = [
            requests.exceptions.ConnectionError("Connection refused"),
            requests.exceptions.ConnectionError("Connection refused"),
            "success",
        ]
        mock_func = mock.Mock(side_effect=side_effects)
        
        # Create a decorated function with a small backoff
        decorated_func = with_exponential_backoff(
            max_retries=3,
            initial_backoff=0.01,
            backoff_factor=1.0,
        )(mock_func)
        
        # Call the function
        start_time = time.time()
        result = decorated_func()
        elapsed_time = time.time() - start_time
        
        # Check that the function was called three times and returned the expected result
        self.assertEqual(mock_func.call_count, 3)
        self.assertEqual(result, "success")
        
        # Check that the backoff was applied (at least 0.01 + 0.01 seconds)
        self.assertGreaterEqual(elapsed_time, 0.02)
    
    def test_with_exponential_backoff_max_retries(self):
        """Test reaching maximum retries with exponential backoff."""
        # Mock a function that always fails
        mock_func = mock.Mock(side_effect=requests.exceptions.ConnectionError("Connection refused"))
        
        # Create a decorated function with a small backoff
        decorated_func = with_exponential_backoff(
            max_retries=2,
            initial_backoff=0.01,
            backoff_factor=1.0,
        )(mock_func)
        
        # Call the function and expect it to raise an exception
        with self.assertRaises(requests.exceptions.ConnectionError):
            decorated_func()
        
        # Check that the function was called three times (initial + 2 retries)
        self.assertEqual(mock_func.call_count, 3)
    
    @mock.patch("requests.Session.get")
    def test_hook_retry_on_server_error(self, mock_get):
        """Test that the hook retries on server errors."""
        # Mock responses: first two fail with 503, third succeeds
        mock_responses = [
            mock.Mock(
                status_code=503,
                raise_for_status=mock.Mock(
                    side_effect=requests.exceptions.HTTPError(
                        "503 Server Error", 
                        response=mock.Mock(
                            status_code=503,
                            json=mock.Mock(return_value={"error": "Service Unavailable"})
                        )
                    )
                )
            ),
            mock.Mock(
                status_code=503,
                raise_for_status=mock.Mock(
                    side_effect=requests.exceptions.HTTPError(
                        "503 Server Error", 
                        response=mock.Mock(
                            status_code=503,
                            json=mock.Mock(return_value={"error": "Service Unavailable"})
                        )
                    )
                )
            ),
            mock.Mock(
                status_code=200,
                raise_for_status=mock.Mock(return_value=None),
                json=mock.Mock(return_value={"status": "ok"})
            ),
        ]
        mock_get.side_effect = mock_responses
        
        # Patch time.sleep to avoid waiting during tests
        with mock.patch("time.sleep"):
            # Call the hook method
            result = self.hook.get_health()
        
        # Check that the request was made three times
        self.assertEqual(mock_get.call_count, 3)
        
        # Check the result
        self.assertEqual(result, {"status": "ok"})


if __name__ == "__main__":
    unittest.main()
