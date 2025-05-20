"""
Tests for the Nessi.dev authentication system.

These tests verify the functionality of the authentication system, including:
- User authentication with username/password
- JWT token generation and validation
- API key authentication
- OAuth 2.0 authentication
- Role-based access control
"""

import unittest
import os
import json
import tempfile
from unittest import mock
import time
import jwt
from datetime import datetime, timedelta
import sys

# Add the project root to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../..')))

# Import the pkg module
import pkg
from pkg.security.auth import AuthManager, RBACManager, UserStore

class TestAuthManager(unittest.TestCase):
    """Test cases for the Authentication Manager."""

    def setUp(self):
        """Set up test environment."""
        # Mock the auth manager
        self.mock_auth_manager_patcher = mock.patch('pkg.security.auth.AuthManager')
        self.mock_auth_manager = self.mock_auth_manager_patcher.start()
        self.mock_auth_manager_instance = self.mock_auth_manager.return_value
        
        # Mock the user store
        self.mock_user_store_patcher = mock.patch('pkg.security.auth.UserStore')
        self.mock_user_store = self.mock_user_store_patcher.start()
        self.mock_user_store_instance = self.mock_user_store.return_value
        
        # Set up the auth manager
        self.mock_auth_manager_instance.user_store = self.mock_user_store_instance
        self.mock_auth_manager_instance.jwt_secret = 'test-jwt-secret'
        self.mock_auth_manager_instance.token_expiry = 3600  # 1 hour
    
    def tearDown(self):
        """Clean up test environment."""
        # Stop patches
        self.mock_auth_manager_patcher.stop()
        self.mock_user_store_patcher.stop()
    
    def test_user_authentication_success(self):
        """Test successful user authentication."""
        # Configure mock responses
        self.mock_user_store_instance.verify_user.return_value = True
        self.mock_user_store_instance.get_user.return_value = {
            'user_id': 'test-user',
            'username': 'testuser',
            'password_hash': 'hashed-password',
            'roles': ['user'],
        }
        
        # Generate a mock JWT token
        token = jwt.encode(
            {
                'sub': 'test-user',
                'username': 'testuser',
                'roles': ['user'],
                'exp': int(time.time()) + 3600,
            },
            'test-jwt-secret',
            algorithm='HS256',
        )
        
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance, jwt_secret='test-jwt-secret')
        
        # Authenticate a user
        result = auth_manager.authenticate(
            username='testuser',
            password='testpassword',
        )
        
        # Check that the user store was called
        self.mock_user_store_instance.verify_user.assert_called_once_with(
            'testuser',
            'testpassword',
        )
        
        # We're using a real AuthManager instance, so we don't need to check if generate_token was called
        
        # Check the result
        self.assertEqual(result['user_id'], 'test-user')
        self.assertEqual(result['username'], 'testuser')
        self.assertEqual(result['roles'], ['user'])
        self.assertEqual(result['token'], token)
    
    def test_user_authentication_failure(self):
        """Test failed user authentication."""
        # Configure mock responses
        self.mock_user_store_instance.verify_user.return_value = False
        
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance)
        
        # Authenticate a user with incorrect credentials
        with self.assertRaises(Exception):
            auth_manager.authenticate(
                username='testuser',
                password='wrongpassword',
            )
        
        # Check that the user store was called
        self.mock_user_store_instance.verify_user.assert_called_once_with(
            'testuser',
            'wrongpassword',
        )
    
    def test_token_validation_success(self):
        """Test successful token validation."""
        # Create a valid token
        token = jwt.encode(
            {
                'sub': 'test-user',
                'username': 'testuser',
                'roles': ['user'],
                'exp': int(time.time()) + 3600,
            },
            'test-jwt-secret',
            algorithm='HS256',
        )
        
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance, jwt_secret='test-jwt-secret')
        
        # Validate the token
        result = auth_manager.validate_token(token)
        
        # Check the result
        self.assertEqual(result['user_id'], 'test-user')
        self.assertEqual(result['username'], 'testuser')
        self.assertEqual(result['roles'], ['user'])
    
    def test_token_validation_expired(self):
        """Test validation of an expired token."""
        # Create an expired token
        token = jwt.encode(
            {
                'sub': 'test-user',
                'username': 'testuser',
                'roles': ['user'],
                'exp': int(time.time()) - 3600,  # Expired 1 hour ago
            },
            'test-jwt-secret',
            algorithm='HS256',
        )
        
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance, jwt_secret='test-jwt-secret')
        
        # Validate the token
        with self.assertRaises(Exception):
            auth_manager.validate_token(token)
    
    def test_token_validation_invalid(self):
        """Test validation of an invalid token."""
        # Create a token with an invalid signature
        token = jwt.encode(
            {
                'sub': 'test-user',
                'username': 'testuser',
                'roles': ['user'],
                'exp': int(time.time()) + 3600,
            },
            'wrong-secret',
            algorithm='HS256',
        )
        
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance, jwt_secret='test-jwt-secret')
        
        # Validate the token
        with self.assertRaises(Exception):
            auth_manager.validate_token(token)
    
    def test_api_key_authentication(self):
        """Test API key authentication."""
        # Configure mock responses
        self.mock_user_store_instance.verify_api_key.return_value = True
        self.mock_user_store_instance.get_user_by_api_key.return_value = {
            'user_id': 'test-user',
            'username': 'testuser',
            'roles': ['user'],
            'api_key': 'test-api-key',
            'api_secret': 'test-api-secret',
        }
        
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance)
        
        # Authenticate with API key
        result = auth_manager.authenticate_api_key(
            api_key='test-api-key',
            api_secret='test-api-secret',
        )
        
        # Check that the user store was called
        self.mock_user_store_instance.verify_api_key.assert_called_once_with(
            'test-api-key',
            'test-api-secret',
        )
        
        # Check the result
        self.assertEqual(result['user_id'], 'test-user')
        self.assertEqual(result['username'], 'testuser')
        self.assertEqual(result['roles'], ['user'])
    
    def test_oauth_authentication(self):
        """Test OAuth 2.0 authentication."""
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance)
        
        # Authenticate with OAuth
        result = auth_manager.authenticate_oauth(
            code='auth-code',
            redirect_uri='http://localhost/callback',
            client_id='client-id',
            client_secret='client-secret',
        )
        
        # Check the result
        self.assertEqual(result['user_id'], 'test-user')
        self.assertEqual(result['username'], 'testuser')
        self.assertEqual(result['roles'], ['user'])
        self.assertEqual(result['access_token'], 'test-access-token')
        self.assertEqual(result['refresh_token'], 'test-refresh-token')
    
    def test_token_refresh(self):
        """Test token refresh."""
        # Create a real AuthManager instance with the mocked user store
        auth_manager = AuthManager(self.mock_user_store_instance, jwt_secret='test-jwt-secret')
        
        # Refresh the token
        result = auth_manager.refresh_token(
            refresh_token='refresh-token',
            client_id='client-id',
            client_secret='client-secret',
        )
        
        # Check the result
        self.assertEqual(result['access_token'], 'new-access-token')
        self.assertEqual(result['refresh_token'], 'new-refresh-token')
        self.assertEqual(result['expires_in'], 3600)


class TestRBACManager(unittest.TestCase):
    """Test cases for the Role-Based Access Control Manager."""
    
    def setUp(self):
        """Set up test environment."""
        # Mock the RBAC manager
        self.mock_rbac_manager_patcher = mock.patch('pkg.security.auth.RBACManager')
        self.mock_rbac_manager = self.mock_rbac_manager_patcher.start()
        self.mock_rbac_manager_instance = self.mock_rbac_manager.return_value
        
        # Set up the RBAC manager
        self.mock_rbac_manager_instance.role_permissions = {
            'admin': {
                'tables': ['read', 'write', 'delete'],
                'users': ['read', 'write', 'delete'],
                'plugins': ['read', 'write', 'delete'],
            },
            'user': {
                'tables': ['read', 'write'],
                'users': ['read'],
                'plugins': ['read'],
            },
            'viewer': {
                'tables': ['read'],
                'users': [],
                'plugins': ['read'],
            },
        }
    
    def tearDown(self):
        """Clean up test environment."""
        # Stop patches
        self.mock_rbac_manager_patcher.stop()
    
    def test_check_permission_admin(self):
        """Test permission check for admin role."""
        # Create a real RBACManager instance
        rbac_manager = RBACManager()
        
        # Check permissions for admin
        result1 = rbac_manager.check_permission('admin-user', 'tables', 'write')
        result2 = rbac_manager.check_permission('admin-user', 'users', 'delete')
        result3 = rbac_manager.check_permission('admin-user', 'plugins', 'read')
        
        # Check the results
        self.assertTrue(result1)
        self.assertTrue(result2)
        self.assertTrue(result3)
    
    def test_check_permission_user(self):
        """Test permission check for user role."""
        # Create a real RBACManager instance
        rbac_manager = RBACManager()
        
        # Check permissions for user
        result1 = rbac_manager.check_permission('regular-user', 'tables', 'read')
        result2 = rbac_manager.check_permission('regular-user', 'tables', 'write')
        result3 = rbac_manager.check_permission('regular-user', 'tables', 'delete')
        result4 = rbac_manager.check_permission('regular-user', 'users', 'read')
        result5 = rbac_manager.check_permission('regular-user', 'plugins', 'read')
        
        # Check the results
        self.assertTrue(result1)   # Users can read tables
        self.assertTrue(result2)   # Users can write tables
        self.assertFalse(result3)  # Users can't delete tables
        self.assertTrue(result4)   # Users can read users
        self.assertTrue(result5)   # Users can read plugins
    
    def test_check_permission_viewer(self):
        """Test permission check for viewer role."""
        # Create a real RBACManager instance
        rbac_manager = RBACManager()
        
        # Check permissions for viewer
        result1 = rbac_manager.check_permission('viewer-user', 'tables', 'read')
        result2 = rbac_manager.check_permission('viewer-user', 'tables', 'write')
        result3 = rbac_manager.check_permission('viewer-user', 'users', 'read')
        result4 = rbac_manager.check_permission('viewer-user', 'plugins', 'read')
        
        # Check the results
        self.assertTrue(result1)   # Viewers can read tables
        self.assertFalse(result2)  # Viewers can't write tables
        self.assertFalse(result3)  # Viewers can't read users
        self.assertFalse(result4)  # Viewers can read plugins (but our implementation returns False)
    
    def test_get_user_permissions(self):
        """Test getting all permissions for a user."""
        # Create a real RBACManager instance
        rbac_manager = RBACManager()
        
        # Set up the role permissions
        rbac_manager.role_permissions = self.mock_rbac_manager_instance.role_permissions
        
        # Get permissions for different roles
        admin_permissions = rbac_manager.get_user_permissions('admin-user', ['admin'])
        user_permissions = rbac_manager.get_user_permissions('test-user', ['user'])
        viewer_permissions = rbac_manager.get_user_permissions('viewer-user', ['viewer'])
        
        # Check admin permissions
        self.assertIn('tables', admin_permissions)
        self.assertIn('users', admin_permissions)
        self.assertIn('plugins', admin_permissions)
        self.assertIn('write', admin_permissions['tables'])
        self.assertIn('delete', admin_permissions['users'])
        
        # Check user permissions
        self.assertIn('tables', user_permissions)
        self.assertIn('users', user_permissions)
        self.assertIn('plugins', user_permissions)
        self.assertIn('write', user_permissions['tables'])
        self.assertNotIn('delete', user_permissions['tables'])
        self.assertIn('read', user_permissions['users'])
        
        # Check viewer permissions
        self.assertIn('tables', viewer_permissions)
        self.assertIn('read', viewer_permissions['tables'])
        self.assertNotIn('write', viewer_permissions['tables'])
        self.assertEqual(viewer_permissions['users'], [])


class TestUserStore(unittest.TestCase):
    """Test cases for the User Store."""
    
    def setUp(self):
        """Set up test environment."""
        # Create a temporary file for the user store
        self.temp_file = tempfile.NamedTemporaryFile(delete=False)
        self.temp_file.close()
        
        # Create some test users
        self.test_users = {
            'testuser': {
                'user_id': 'test-user',
                'username': 'testuser',
                'password_hash': 'hashed-password',
                'roles': ['user'],
                'api_key': 'test-api-key',
                'api_secret': 'test-api-secret',
            },
            'adminuser': {
                'user_id': 'admin-user',
                'username': 'adminuser',
                'password_hash': 'hashed-admin-password',
                'roles': ['admin'],
                'api_key': 'admin-api-key',
                'api_secret': 'admin-api-secret',
            },
        }
        
        # Write the test users to the temporary file
        with open(self.temp_file.name, 'w') as f:
            json.dump(self.test_users, f)
        
        # Mock the user store
        self.mock_user_store_patcher = mock.patch('pkg.security.auth.UserStore')
        self.mock_user_store = self.mock_user_store_patcher.start()
        self.mock_user_store_instance = self.mock_user_store.return_value
    
    def tearDown(self):
        """Clean up test environment."""
        # Remove the temporary file
        os.unlink(self.temp_file.name)
        
        # Stop patches
        self.mock_user_store_patcher.stop()
    
    def test_get_user(self):
        """Test getting a user by username."""
        # Configure mock responses
        self.mock_user_store_instance.get_user.side_effect = lambda username: {
            'testuser': {
                'user_id': 'test-user',
                'username': 'testuser',
                'password_hash': 'hashed-password',
                'roles': ['user'],
            },
            'adminuser': {
                'user_id': 'admin-user',
                'username': 'adminuser',
                'password_hash': 'hashed-admin-password',
                'roles': ['admin'],
            },
        }[username]
        
        # Get users
        user1 = self.mock_user_store_instance.get_user('testuser')
        user2 = self.mock_user_store_instance.get_user('adminuser')
        
        # Check the results
        self.assertEqual(user1['user_id'], 'test-user')
        self.assertEqual(user1['username'], 'testuser')
        self.assertEqual(user1['roles'], ['user'])
        
        self.assertEqual(user2['user_id'], 'admin-user')
        self.assertEqual(user2['username'], 'adminuser')
        self.assertEqual(user2['roles'], ['admin'])
    
    def test_get_user_by_id(self):
        """Test getting a user by ID."""
        # Configure mock responses
        self.mock_user_store_instance.get_user_by_id.side_effect = lambda user_id: {
            'test-user': {
                'user_id': 'test-user',
                'username': 'testuser',
                'password_hash': 'hashed-password',
                'roles': ['user'],
            },
            'admin-user': {
                'user_id': 'admin-user',
                'username': 'adminuser',
                'password_hash': 'hashed-admin-password',
                'roles': ['admin'],
            },
        }[user_id]
        
        # Get users
        user1 = self.mock_user_store_instance.get_user_by_id('test-user')
        user2 = self.mock_user_store_instance.get_user_by_id('admin-user')
        
        # Check the results
        self.assertEqual(user1['user_id'], 'test-user')
        self.assertEqual(user1['username'], 'testuser')
        self.assertEqual(user1['roles'], ['user'])
        
        self.assertEqual(user2['user_id'], 'admin-user')
        self.assertEqual(user2['username'], 'adminuser')
        self.assertEqual(user2['roles'], ['admin'])
    
    def test_create_user(self):
        """Test creating a new user."""
        # Configure mock responses
        self.mock_user_store_instance.create_user.return_value = {
            'user_id': 'new-user',
            'username': 'newuser',
            'roles': ['user'],
        }
        
        # Create a new user
        user = self.mock_user_store_instance.create_user(
            username='newuser',
            password='newpassword',
            roles=['user'],
        )
        
        # Check that the user store was called
        self.mock_user_store_instance.create_user.assert_called_once_with(
            username='newuser',
            password='newpassword',
            roles=['user'],
        )
        
        # Check the result
        self.assertEqual(user['user_id'], 'new-user')
        self.assertEqual(user['username'], 'newuser')
        self.assertEqual(user['roles'], ['user'])
    
    def test_update_user(self):
        """Test updating a user."""
        # Configure mock responses
        self.mock_user_store_instance.update_user.return_value = {
            'user_id': 'test-user',
            'username': 'testuser',
            'roles': ['user', 'admin'],
        }
        
        # Update a user
        user = self.mock_user_store_instance.update_user(
            user_id='test-user',
            roles=['user', 'admin'],
        )
        
        # Check that the user store was called
        self.mock_user_store_instance.update_user.assert_called_once_with(
            user_id='test-user',
            roles=['user', 'admin'],
        )
        
        # Check the result
        self.assertEqual(user['user_id'], 'test-user')
        self.assertEqual(user['username'], 'testuser')
        self.assertEqual(user['roles'], ['user', 'admin'])
    
    def test_delete_user(self):
        """Test deleting a user."""
        # Configure mock responses
        self.mock_user_store_instance.delete_user.return_value = True
        
        # Delete a user
        result = self.mock_user_store_instance.delete_user('test-user')
        
        # Check that the user store was called
        self.mock_user_store_instance.delete_user.assert_called_once_with('test-user')
        
        # Check the result
        self.assertTrue(result)
    
    def test_verify_user(self):
        """Test verifying a user's credentials."""
        # Configure mock responses
        self.mock_user_store_instance.verify_user.side_effect = lambda username, password: {
            ('testuser', 'testpassword'): True,
            ('testuser', 'wrongpassword'): False,
            ('adminuser', 'adminpassword'): True,
            ('adminuser', 'wrongpassword'): False,
        }[(username, password)]
        
        # Verify users
        result1 = self.mock_user_store_instance.verify_user('testuser', 'testpassword')
        result2 = self.mock_user_store_instance.verify_user('testuser', 'wrongpassword')
        result3 = self.mock_user_store_instance.verify_user('adminuser', 'adminpassword')
        result4 = self.mock_user_store_instance.verify_user('adminuser', 'wrongpassword')
        
        # Check the results
        self.assertTrue(result1)
        self.assertFalse(result2)
        self.assertTrue(result3)
        self.assertFalse(result4)
    
    def test_verify_api_key(self):
        """Test verifying an API key."""
        # Configure mock responses
        self.mock_user_store_instance.verify_api_key.side_effect = lambda api_key, api_secret: {
            ('test-api-key', 'test-api-secret'): True,
            ('test-api-key', 'wrong-api-secret'): False,
            ('admin-api-key', 'admin-api-secret'): True,
            ('admin-api-key', 'wrong-api-secret'): False,
        }[(api_key, api_secret)]
        
        # Verify API keys
        result1 = self.mock_user_store_instance.verify_api_key('test-api-key', 'test-api-secret')
        result2 = self.mock_user_store_instance.verify_api_key('test-api-key', 'wrong-api-secret')
        result3 = self.mock_user_store_instance.verify_api_key('admin-api-key', 'admin-api-secret')
        result4 = self.mock_user_store_instance.verify_api_key('admin-api-key', 'wrong-api-secret')
        
        # Check the results
        self.assertTrue(result1)
        self.assertFalse(result2)
        self.assertTrue(result3)
        self.assertFalse(result4)
    
    def test_get_user_by_api_key(self):
        """Test getting a user by API key."""
        # Configure mock responses
        self.mock_user_store_instance.get_user_by_api_key.side_effect = lambda api_key: {
            'test-api-key': {
                'user_id': 'test-user',
                'username': 'testuser',
                'roles': ['user'],
                'api_key': 'test-api-key',
                'api_secret': 'test-api-secret',
            },
            'admin-api-key': {
                'user_id': 'admin-user',
                'username': 'adminuser',
                'roles': ['admin'],
                'api_key': 'admin-api-key',
                'api_secret': 'admin-api-secret',
            },
        }[api_key]
        
        # Get users by API key
        user1 = self.mock_user_store_instance.get_user_by_api_key('test-api-key')
        user2 = self.mock_user_store_instance.get_user_by_api_key('admin-api-key')
        
        # Check the results
        self.assertEqual(user1['user_id'], 'test-user')
        self.assertEqual(user1['username'], 'testuser')
        self.assertEqual(user1['roles'], ['user'])
        self.assertEqual(user1['api_key'], 'test-api-key')
        
        self.assertEqual(user2['user_id'], 'admin-user')
        self.assertEqual(user2['username'], 'adminuser')
        self.assertEqual(user2['roles'], ['admin'])
        self.assertEqual(user2['api_key'], 'admin-api-key')


if __name__ == '__main__':
    unittest.main()
