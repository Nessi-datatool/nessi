"""
Tests for OAuth 2.0 authentication in the Nessi Airflow integration.
"""

import unittest
from unittest import mock
import json
import os
from datetime import datetime, timedelta

from airflow.models import Connection
from airflow.utils.session import provide_session
import requests

from nessi_airflow.hooks.nessi_hook import NessiHook
from nessi_airflow.auth.oauth_handler import OAuthHandler, refresh_oauth_token


class TestOAuthAuthentication(unittest.TestCase):
    """Test cases for OAuth 2.0 authentication in Airflow integration."""

    def setUp(self):
        """Set up test environment."""
        # Mock the requests module
        self.mock_requests_patcher = mock.patch('nessi_airflow.auth.oauth_handler.requests')
        self.mock_requests = self.mock_requests_patcher.start()
        
        # Mock the airflow connection
        self.mock_get_connection_patcher = mock.patch('nessi_airflow.hooks.nessi_hook.BaseHook.get_connection')
        self.mock_get_connection = self.mock_get_connection_patcher.start()
        
        # Create a mock connection
        self.mock_conn = mock.MagicMock(spec=Connection)
        self.mock_conn.host = 'http://test-host'
        self.mock_conn.login = 'test-client-id'
        self.mock_conn.password = 'test-client-secret'
        self.mock_conn.extra_dejson = {
            'auth_type': 'oauth2',
            'oauth_token': 'test-token',
            'oauth_refresh_token': 'test-refresh-token',
            'oauth_token_url': 'http://test-token-url',
            'oauth_token_expiry': (datetime.now() + timedelta(hours=1)).isoformat(),
        }
        self.mock_get_connection.return_value = self.mock_conn
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_requests_patcher.stop()
        self.mock_get_connection_patcher.stop()
    
    def test_oauth_handler_init(self):
        """Test initialization of OAuth handler."""
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='test-token',
            refresh_token='test-refresh-token',
            token_expiry=(datetime.now() + timedelta(hours=1)).isoformat(),
        )
        
        self.assertEqual(handler.token_url, 'http://test-token-url')
        self.assertEqual(handler.client_id, 'test-client-id')
        self.assertEqual(handler.client_secret, 'test-client-secret')
        self.assertEqual(handler.access_token, 'test-token')
        self.assertEqual(handler.refresh_token, 'test-refresh-token')
        self.assertFalse(handler.is_token_expired())
    
    def test_is_token_expired(self):
        """Test token expiry check."""
        # Token expired
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='test-token',
            refresh_token='test-refresh-token',
            token_expiry=(datetime.now() - timedelta(hours=1)).isoformat(),
        )
        self.assertTrue(handler.is_token_expired())
        
        # Token not expired
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='test-token',
            refresh_token='test-refresh-token',
            token_expiry=(datetime.now() + timedelta(hours=1)).isoformat(),
        )
        self.assertFalse(handler.is_token_expired())
    
    def test_refresh_token(self):
        """Test token refresh."""
        # Mock the token refresh response
        mock_response = mock.MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            'access_token': 'new-token',
            'refresh_token': 'new-refresh-token',
            'expires_in': 3600,
        }
        self.mock_requests.post.return_value = mock_response
        
        # Create handler with expired token
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='test-token',
            refresh_token='test-refresh-token',
            token_expiry=(datetime.now() - timedelta(hours=1)).isoformat(),
        )
        
        # Refresh token
        handler.refresh_token()
        
        # Check that token was refreshed
        self.assertEqual(handler.access_token, 'new-token')
        self.assertEqual(handler.refresh_token, 'new-refresh-token')
        self.assertFalse(handler.is_token_expired())
        
        # Check that the request was made correctly
        self.mock_requests.post.assert_called_once_with(
            'http://test-token-url',
            data={
                'grant_type': 'refresh_token',
                'refresh_token': 'test-refresh-token',
                'client_id': 'test-client-id',
                'client_secret': 'test-client-secret',
            },
            headers={'Content-Type': 'application/x-www-form-urlencoded'},
        )
    
    def test_refresh_token_error(self):
        """Test token refresh with error."""
        # Mock the token refresh response with error
        mock_response = mock.MagicMock()
        mock_response.status_code = 400
        mock_response.json.return_value = {
            'error': 'invalid_grant',
            'error_description': 'Invalid refresh token',
        }
        self.mock_requests.post.return_value = mock_response
        
        # Create handler with expired token
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='test-token',
            refresh_token='test-refresh-token',
            token_expiry=(datetime.now() - timedelta(hours=1)).isoformat(),
        )
        
        # Refresh token should raise exception
        with self.assertRaises(Exception):
            handler.refresh_token()
        
        # Check that the request was made correctly
        self.mock_requests.post.assert_called_once()
    
    def test_get_auth_headers(self):
        """Test getting authentication headers."""
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='test-token',
            refresh_token='test-refresh-token',
            token_expiry=(datetime.now() + timedelta(hours=1)).isoformat(),
        )
        
        headers = handler.get_auth_headers()
        self.assertEqual(headers['Authorization'], 'Bearer test-token')
    
    def test_nessi_hook_with_oauth(self):
        """Test NessiHook with OAuth authentication."""
        # Create hook
        hook = NessiHook(conn_id='nessi_default')
        
        # Mock the requests session
        mock_session = mock.MagicMock()
        hook.session = mock_session
        
        # Call a method that makes an API request
        hook.get_quality_check('test-check-id')
        
        # Check that the request was made with OAuth headers
        mock_session.request.assert_called_once()
        args, kwargs = mock_session.request.call_args
        self.assertEqual(kwargs['headers']['Authorization'], 'Bearer test-token')
    
    @mock.patch('nessi_airflow.hooks.nessi_hook.OAuthHandler')
    def test_nessi_hook_with_expired_token(self, mock_oauth_handler):
        """Test NessiHook with expired OAuth token."""
        # Mock the OAuth handler
        mock_handler_instance = mock.MagicMock()
        mock_handler_instance.is_token_expired.return_value = True
        mock_handler_instance.get_auth_headers.return_value = {'Authorization': 'Bearer new-token'}
        mock_oauth_handler.return_value = mock_handler_instance
        
        # Create hook
        hook = NessiHook(conn_id='nessi_default')
        
        # Mock the requests session
        mock_session = mock.MagicMock()
        hook.session = mock_session
        
        # Call a method that makes an API request
        hook.get_quality_check('test-check-id')
        
        # Check that token was refreshed
        mock_handler_instance.refresh_token.assert_called_once()
        
        # Check that the request was made with new OAuth headers
        mock_session.request.assert_called_once()
        args, kwargs = mock_session.request.call_args
        self.assertEqual(kwargs['headers']['Authorization'], 'Bearer new-token')
    
    @provide_session
    def test_update_connection_with_new_tokens(self, session=None):
        """Test updating Airflow connection with new tokens."""
        # Mock the session
        mock_session = mock.MagicMock()
        
        # Create handler with new tokens
        handler = OAuthHandler(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            access_token='new-token',
            refresh_token='new-refresh-token',
            token_expiry=(datetime.now() + timedelta(hours=1)).isoformat(),
        )
        
        # Update connection
        handler.update_connection('nessi_default', session=mock_session)
        
        # Check that the connection was updated
        mock_session.query.assert_called_once()
        mock_session.commit.assert_called_once()


class TestRefreshOAuthToken(unittest.TestCase):
    """Test cases for the refresh_oauth_token function."""

    def setUp(self):
        """Set up test environment."""
        # Mock the requests module
        self.mock_requests_patcher = mock.patch('nessi_airflow.auth.oauth_handler.requests')
        self.mock_requests = self.mock_requests_patcher.start()
    
    def tearDown(self):
        """Clean up test environment."""
        self.mock_requests_patcher.stop()
    
    def test_refresh_oauth_token(self):
        """Test refreshing OAuth token."""
        # Mock the token refresh response
        mock_response = mock.MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            'access_token': 'new-token',
            'refresh_token': 'new-refresh-token',
            'expires_in': 3600,
        }
        self.mock_requests.post.return_value = mock_response
        
        # Refresh token
        result = refresh_oauth_token(
            token_url='http://test-token-url',
            client_id='test-client-id',
            client_secret='test-client-secret',
            refresh_token='test-refresh-token',
        )
        
        # Check result
        self.assertEqual(result['access_token'], 'new-token')
        self.assertEqual(result['refresh_token'], 'new-refresh-token')
        self.assertEqual(result['expires_in'], 3600)
        
        # Check that the request was made correctly
        self.mock_requests.post.assert_called_once_with(
            'http://test-token-url',
            data={
                'grant_type': 'refresh_token',
                'refresh_token': 'test-refresh-token',
                'client_id': 'test-client-id',
                'client_secret': 'test-client-secret',
            },
            headers={'Content-Type': 'application/x-www-form-urlencoded'},
        )
    
    def test_refresh_oauth_token_error(self):
        """Test refreshing OAuth token with error."""
        # Mock the token refresh response with error
        mock_response = mock.MagicMock()
        mock_response.status_code = 400
        mock_response.json.return_value = {
            'error': 'invalid_grant',
            'error_description': 'Invalid refresh token',
        }
        self.mock_requests.post.return_value = mock_response
        
        # Refresh token should raise exception
        with self.assertRaises(Exception):
            refresh_oauth_token(
                token_url='http://test-token-url',
                client_id='test-client-id',
                client_secret='test-client-secret',
                refresh_token='test-refresh-token',
            )
        
        # Check that the request was made correctly
        self.mock_requests.post.assert_called_once()


if __name__ == '__main__':
    unittest.main()
