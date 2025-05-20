#!/usr/bin/env python3
"""
Integration tests for the Nessi client security features.

These tests verify that the Python client correctly interacts with
the security features of the Nessi monitoring system.

Note: These tests require a running Nessi server with security features enabled.
"""

import unittest
import os
import sys
import time
import uuid
import requests
from datetime import datetime

sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from nessi_client import NessiClient, Metric


class TestSecurityIntegration(unittest.TestCase):
    """Integration tests for security features."""

    @classmethod
    def setUpClass(cls):
        """Set up test fixtures that are reused across tests."""
        # Server configuration
        cls.base_url = os.environ.get("NESSI_TEST_URL", "http://localhost:8080")
        cls.admin_username = os.environ.get("NESSI_TEST_ADMIN_USER", "admin")
        cls.admin_password = os.environ.get("NESSI_TEST_ADMIN_PASSWORD", "admin")
        
        # Test user information
        cls.test_username = f"testuser_{uuid.uuid4().hex[:8]}"
        cls.test_password = "testpassword123"
        cls.test_email = f"{cls.test_username}@example.com"
        
        # Skip tests if server is not available
        try:
            response = requests.get(f"{cls.base_url}/health")
            if response.status_code != 200:
                raise unittest.SkipTest("Nessi server is not available")
        except requests.RequestException:
            raise unittest.SkipTest("Nessi server is not available")
        
        # Create admin client
        cls.admin_client = NessiClient(
            base_url=cls.base_url,
            username=cls.admin_username,
            password=cls.admin_password
        )
        
        # Create test user
        try:
            user_data = {
                "username": cls.test_username,
                "password": cls.test_password,
                "email": cls.test_email,
                "role": "user"
            }
            cls.admin_client._request("post", "api/users", json=user_data)
            print(f"Created test user: {cls.test_username}")
        except Exception as e:
            print(f"Failed to create test user: {e}")
            raise

    @classmethod
    def tearDownClass(cls):
        """Clean up after all tests have run."""
        # Delete test user
        try:
            cls.admin_client._request("delete", f"api/users/{cls.test_username}")
            print(f"Deleted test user: {cls.test_username}")
        except Exception as e:
            print(f"Failed to delete test user: {e}")

    def test_01_username_password_authentication(self):
        """Test authentication with username and password."""
        # Create client with username/password
        client = NessiClient(
            base_url=self.base_url,
            username=self.test_username,
            password=self.test_password
        )
        
        # Verify token was obtained
        self.assertIsNotNone(client.token)
        self.assertIn("Authorization", client.session.headers)
        
        # Test a simple API call
        try:
            response = client._request("get", "api/metrics", params={"limit": 1})
            self.assertEqual(response.status_code, 200)
        except Exception as e:
            self.fail(f"API call failed: {e}")

    def test_02_api_key_generation_and_authentication(self):
        """Test API key generation and authentication."""
        # Create client with username/password
        client = NessiClient(
            base_url=self.base_url,
            username=self.test_username,
            password=self.test_password
        )
        
        # Generate API key
        try:
            response = client._request("post", "api/apikeys", json={"description": "Test API key"})
            api_key_data = response.json()
            api_key = api_key_data.get("key")
            self.assertIsNotNone(api_key)
            print(f"Generated API key: {api_key[:8]}...{api_key[-4:]}")
        except Exception as e:
            self.fail(f"Failed to generate API key: {e}")
        
        # Create client with API key
        api_client = NessiClient(
            base_url=self.base_url,
            api_key=api_key
        )
        
        # Test a simple API call
        try:
            response = api_client._request("get", "api/metrics", params={"limit": 1})
            self.assertEqual(response.status_code, 200)
        except Exception as e:
            self.fail(f"API call with API key failed: {e}")
        
        # List API keys
        try:
            response = client._request("get", "api/apikeys")
            api_keys = response.json()
            self.assertGreaterEqual(len(api_keys), 1)
            
            # Revoke the API key
            key_id = api_keys[0].get('id')
            client._request("delete", f"api/apikeys/{key_id}")
            print(f"Revoked API key with ID: {key_id}")
        except Exception as e:
            self.fail(f"Failed to manage API keys: {e}")

    def test_03_role_based_access_control(self):
        """Test role-based access control."""
        # Create admin client
        admin_client = NessiClient(
            base_url=self.base_url,
            username=self.admin_username,
            password=self.admin_password
        )
        
        # Create user client
        user_client = NessiClient(
            base_url=self.base_url,
            username=self.test_username,
            password=self.test_password
        )
        
        # Test admin-only operation: listing users
        try:
            # Admin should be able to list users
            response = admin_client._request("get", "api/users")
            self.assertEqual(response.status_code, 200)
            users = response.json()
            self.assertGreaterEqual(len(users), 2)  # At least admin and test user
        except Exception as e:
            self.fail(f"Admin failed to list users: {e}")
        
        # Regular user should not be able to list users
        try:
            user_client._request("get", "api/users")
            self.fail("Regular user should not be able to list users")
        except requests.RequestException as e:
            # Expected to fail with 403 Forbidden
            self.assertIn(e.response.status_code, [401, 403])
    
    def test_04_token_refresh(self):
        """Test token refresh functionality."""
        # Create client with username/password
        client = NessiClient(
            base_url=self.base_url,
            username=self.test_username,
            password=self.test_password
        )
        
        # Store original token
        original_token = client.token
        self.assertIsNotNone(original_token)
        
        # Force token refresh
        client._refresh_token_if_needed()
        
        # Token should be the same (not expired yet)
        self.assertEqual(client.token, original_token)
        
        # Simulate token expiration
        client.token_expiry = datetime(2000, 1, 1)  # Set to past date
        client._refresh_token_if_needed()
        
        # Token should be refreshed
        self.assertIsNotNone(client.token)
        self.assertNotEqual(client.token, original_token)

    def test_05_metric_operations_with_authentication(self):
        """Test metric operations with authentication."""
        # Create client with username/password
        client = NessiClient(
            base_url=self.base_url,
            username=self.test_username,
            password=self.test_password
        )
        
        # Create a test metric
        metric_name = f"test_metric_{uuid.uuid4().hex[:8]}"
        metric = Metric(
            name=metric_name,
            value=42.0,
            timestamp=datetime.now(),
            tags={"test": "security_integration"},
            metadata={"source": "python_test"}
        )
        
        # Send the metric
        try:
            result = client.send_metric(metric)
            self.assertIsNotNone(result)
            print(f"Sent metric: {metric_name}")
        except Exception as e:
            self.fail(f"Failed to send metric: {e}")
        
        # Wait for the metric to be processed
        time.sleep(1)
        
        # Get the metric
        try:
            metrics = client.get_metrics(name=metric_name, limit=1)
            self.assertGreaterEqual(len(metrics), 1)
            self.assertEqual(metrics[0].name, metric_name)
            self.assertEqual(metrics[0].value, 42.0)
            self.assertEqual(metrics[0].tags.get("test"), "security_integration")
        except Exception as e:
            self.fail(f"Failed to get metric: {e}")

    def test_06_invalid_credentials(self):
        """Test behavior with invalid credentials."""
        # Try to authenticate with invalid username
        with self.assertRaises(requests.RequestException):
            NessiClient(
                base_url=self.base_url,
                username="nonexistent_user",
                password="wrong_password"
            )
        
        # Try to authenticate with invalid API key
        invalid_client = NessiClient(
            base_url=self.base_url,
            api_key="invalid_api_key"
        )
        
        # API call should fail
        with self.assertRaises(requests.RequestException):
            invalid_client.get_metrics(limit=1)

    def test_07_ssl_verification(self):
        """Test SSL verification options."""
        # This test is conditional on whether HTTPS is available
        https_url = self.base_url.replace("http://", "https://")
        if "localhost" in https_url:
            https_url = https_url.replace("localhost:8080", "localhost:8443")
        
        # Try with SSL verification (expected to fail with self-signed cert)
        try:
            client = NessiClient(
                base_url=https_url,
                username=self.test_username,
                password=self.test_password,
                verify_ssl=True
            )
            # If this succeeds, the server has a valid certificate
            print("Server has a valid SSL certificate")
        except requests.RequestException as e:
            # This is expected with self-signed certificates
            print(f"SSL verification failed as expected: {e}")
        
        # Try without SSL verification (should work with self-signed cert)
        try:
            client = NessiClient(
                base_url=https_url,
                username=self.test_username,
                password=self.test_password,
                verify_ssl=False
            )
            # Test a simple API call
            response = client._request("get", "api/metrics", params={"limit": 1})
            self.assertEqual(response.status_code, 200)
            print("Successfully connected to HTTPS endpoint without SSL verification")
        except requests.RequestException as e:
            # This might happen if HTTPS is not enabled
            print(f"HTTPS connection failed: {e}")
            print("This test is skipped if HTTPS is not enabled")


if __name__ == "__main__":
    unittest.main()
