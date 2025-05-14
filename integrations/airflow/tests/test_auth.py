"""
Tests for the authentication utilities in the Nessi Airflow integration.
"""
import unittest
from unittest import mock
import json
import os
import time

import requests
from airflow.models.connection import Connection
from airflow.exceptions import AirflowException

from nessi_airflow.utils.auth import (
    NessiAuthBase,
    ApiKeyAuth,
    OAuth2Auth,
    get_auth_from_connection,
    get_auth_from_env,
)
from nessi_airflow.utils.error_handling import NessiAuthenticationError


class TestAuthUtils(unittest.TestCase):
    """Test the authentication utilities."""

    def setUp(self):
        """Set up the test case."""
        # Create a mock connection
        self.conn_id = "nessi_test"
        self.mock_connection = mock.patch(
            "airflow.hooks.base.BaseHook.get_connection",
            return_value=Connection(
                conn_id=self.conn_id,
                conn_type="http",
                host="http://localhost:8080",
                login="test-api-key",
                password="test-api-secret",
                extra=json.dumps({"auth_type": "api_key"}),
            ),
        )
        self.mock_connection.start()
        
        # Save environment variables
        self.original_env = dict(os.environ)
    
    def tearDown(self):
        """Tear down the test case."""
        self.mock_connection.stop()
        
        # Restore environment variables
        os.environ.clear()
        os.environ.update(self.original_env)
    
    def test_api_key_auth(self):
        """Test API key authentication."""
        auth = ApiKeyAuth(
            base_url="http://localhost:8080",
            api_key="test-api-key",
            api_secret="test-api-secret",
        )
        
        # Check auth headers
        headers = auth.get_auth_headers()
        self.assertEqual(headers["X-API-Key"], "test-api-key")
        self.assertEqual(headers["X-API-Secret"], "test-api-secret")
        
        # Check session update
        session = requests.Session()
        auth.update_session(session)
        self.assertEqual(session.headers["X-API-Key"], "test-api-key")
        self.assertEqual(session.headers["X-API-Secret"], "test-api-secret")
    
    @mock.patch("requests.post")
    def test_oauth2_auth_client_credentials(self, mock_post):
        """Test OAuth2 authentication with client credentials."""
        # Mock token response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "access_token": "test-access-token",
            "token_type": "bearer",
            "expires_in": 3600,
        }
        mock_post.return_value = mock_response
        
        # Create OAuth2 auth
        auth = OAuth2Auth(
            base_url="http://localhost:8080",
            client_id="test-client-id",
            client_secret="test-client-secret",
        )
        
        # Get auth headers (should trigger token request)
        headers = auth.get_auth_headers()
        self.assertEqual(headers["Authorization"], "Bearer test-access-token")
        
        # Check that token request was made correctly
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(args[0], "http://localhost:8080/oauth/token")
        self.assertEqual(kwargs["data"]["grant_type"], "client_credentials")
        self.assertEqual(kwargs["data"]["client_id"], "test-client-id")
        self.assertEqual(kwargs["data"]["client_secret"], "test-client-secret")
    
    @mock.patch("requests.post")
    def test_oauth2_auth_refresh_token(self, mock_post):
        """Test OAuth2 authentication with refresh token."""
        # Mock token response
        mock_response = mock.Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "access_token": "new-access-token",
            "refresh_token": "new-refresh-token",
            "token_type": "bearer",
            "expires_in": 3600,
        }
        mock_post.return_value = mock_response
        
        # Create OAuth2 auth with expired token
        auth = OAuth2Auth(
            base_url="http://localhost:8080",
            client_id="test-client-id",
            client_secret="test-client-secret",
            access_token="old-access-token",
            refresh_token="old-refresh-token",
            expires_at=time.time() - 100,  # Expired
        )
        
        # Get auth headers (should trigger token refresh)
        headers = auth.get_auth_headers()
        self.assertEqual(headers["Authorization"], "Bearer new-access-token")
        
        # Check that refresh token request was made correctly
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(args[0], "http://localhost:8080/oauth/token")
        self.assertEqual(kwargs["data"]["grant_type"], "refresh_token")
        self.assertEqual(kwargs["data"]["refresh_token"], "old-refresh-token")
        
        # Check that tokens were updated
        self.assertEqual(auth.access_token, "new-access-token")
        self.assertEqual(auth.refresh_token, "new-refresh-token")
    
    def test_get_auth_from_connection_api_key(self):
        """Test getting API key authentication from a connection."""
        auth = get_auth_from_connection(self.conn_id)
        
        self.assertIsInstance(auth, ApiKeyAuth)
        self.assertEqual(auth.base_url, "http://localhost:8080")
        self.assertEqual(auth.api_key, "test-api-key")
        self.assertEqual(auth.api_secret, "test-api-secret")
    
    def test_get_auth_from_connection_oauth2(self):
        """Test getting OAuth2 authentication from a connection."""
        # Update mock connection to use OAuth2
        self.mock_connection.stop()
        self.mock_connection = mock.patch(
            "airflow.hooks.base.BaseHook.get_connection",
            return_value=Connection(
                conn_id=self.conn_id,
                conn_type="http",
                host="http://localhost:8080",
                extra=json.dumps({
                    "auth_type": "oauth2",
                    "client_id": "test-client-id",
                    "client_secret": "test-client-secret",
                    "token_url": "http://localhost:8080/custom/token",
                    "access_token": "test-access-token",
                    "refresh_token": "test-refresh-token",
                    "expires_at": time.time() + 3600,
                }),
            ),
        )
        self.mock_connection.start()
        
        auth = get_auth_from_connection(self.conn_id)
        
        self.assertIsInstance(auth, OAuth2Auth)
        self.assertEqual(auth.base_url, "http://localhost:8080")
        self.assertEqual(auth.client_id, "test-client-id")
        self.assertEqual(auth.client_secret, "test-client-secret")
        self.assertEqual(auth.token_url, "http://localhost:8080/custom/token")
        self.assertEqual(auth.access_token, "test-access-token")
        self.assertEqual(auth.refresh_token, "test-refresh-token")
    
    def test_get_auth_from_env_api_key(self):
        """Test getting API key authentication from environment variables."""
        # Set environment variables
        os.environ["NESSI_API_HOST"] = "http://localhost:8080"
        os.environ["NESSI_API_KEY"] = "env-api-key"
        os.environ["NESSI_API_SECRET"] = "env-api-secret"
        
        auth, base_url = get_auth_from_env()
        
        self.assertIsInstance(auth, ApiKeyAuth)
        self.assertEqual(base_url, "http://localhost:8080")
        self.assertEqual(auth.base_url, "http://localhost:8080")
        self.assertEqual(auth.api_key, "env-api-key")
        self.assertEqual(auth.api_secret, "env-api-secret")
    
    def test_get_auth_from_env_oauth2(self):
        """Test getting OAuth2 authentication from environment variables."""
        # Set environment variables
        os.environ["NESSI_API_HOST"] = "http://localhost:8080"
        os.environ["NESSI_CLIENT_ID"] = "env-client-id"
        os.environ["NESSI_CLIENT_SECRET"] = "env-client-secret"
        os.environ["NESSI_TOKEN_URL"] = "http://localhost:8080/env/token"
        os.environ["NESSI_ACCESS_TOKEN"] = "env-access-token"
        os.environ["NESSI_REFRESH_TOKEN"] = "env-refresh-token"
        os.environ["NESSI_TOKEN_EXPIRES_AT"] = str(time.time() + 3600)
        
        auth, base_url = get_auth_from_env()
        
        self.assertIsInstance(auth, OAuth2Auth)
        self.assertEqual(base_url, "http://localhost:8080")
        self.assertEqual(auth.base_url, "http://localhost:8080")
        self.assertEqual(auth.client_id, "env-client-id")
        self.assertEqual(auth.client_secret, "env-client-secret")
        self.assertEqual(auth.token_url, "http://localhost:8080/env/token")
        self.assertEqual(auth.access_token, "env-access-token")
        self.assertEqual(auth.refresh_token, "env-refresh-token")
    
    def test_get_auth_from_env_not_configured(self):
        """Test getting authentication from environment variables when not configured."""
        # Clear relevant environment variables
        for var in ["NESSI_API_HOST", "NESSI_API_KEY", "NESSI_CLIENT_ID"]:
            if var in os.environ:
                del os.environ[var]
        
        auth, base_url = get_auth_from_env()
        
        self.assertIsNone(auth)
        self.assertIsNone(base_url)


if __name__ == "__main__":
    unittest.main()
