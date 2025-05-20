"""
User store for Nessi.dev.
"""

import os
import json
import uuid
import hashlib
from typing import Dict, List, Any, Optional


class UserStore:
    """User store for Nessi.dev."""

    def __init__(self, users_file=None):
        """Initialize the user store.

        Args:
            users_file: Path to the users file.
        """
        self.users_file = users_file
        self.users = {
            'test-user': {
                'user_id': 'test-user',
                'username': 'testuser',
                'password_hash': self._hash_password('testpassword'),
                'roles': ['user'],
                'api_key': 'test-api-key',
                'api_secret': 'test-api-secret',
            },
            'admin-user': {
                'user_id': 'admin-user',
                'username': 'adminuser',
                'password_hash': self._hash_password('adminpassword'),
                'roles': ['admin'],
                'api_key': 'admin-api-key',
                'api_secret': 'admin-api-secret',
            },
        }

    def _hash_password(self, password: str) -> str:
        """Hash a password.

        Args:
            password: Password to hash.

        Returns:
            Hashed password.
        """
        return hashlib.sha256(password.encode()).hexdigest()

    def get_user(self, username: str) -> Dict[str, Any]:
        """Get a user by username.

        Args:
            username: Username.

        Returns:
            User information.

        Raises:
            Exception: If user not found.
        """
        for user_id, user in self.users.items():
            if user['username'] == username:
                return user
        raise Exception(f"User not found: {username}")

    def get_user_by_id(self, user_id: str) -> Dict[str, Any]:
        """Get a user by ID.

        Args:
            user_id: User ID.

        Returns:
            User information.

        Raises:
            Exception: If user not found.
        """
        if user_id in self.users:
            return self.users[user_id]
        raise Exception(f"User not found: {user_id}")

    def get_user_by_api_key(self, api_key: str) -> Dict[str, Any]:
        """Get a user by API key.

        Args:
            api_key: API key.

        Returns:
            User information.

        Raises:
            Exception: If user not found.
        """
        for user_id, user in self.users.items():
            if user.get('api_key') == api_key:
                return user
        raise Exception(f"User not found for API key: {api_key}")

    def verify_user(self, username: str, password: str) -> bool:
        """Verify a user's credentials.

        Args:
            username: Username.
            password: Password.

        Returns:
            True if credentials are valid, False otherwise.
        """
        try:
            user = self.get_user(username)
            return user['password_hash'] == self._hash_password(password)
        except Exception:
            return False

    def verify_api_key(self, api_key: str, api_secret: str) -> bool:
        """Verify an API key.

        Args:
            api_key: API key.
            api_secret: API secret.

        Returns:
            True if API key is valid, False otherwise.
        """
        try:
            user = self.get_user_by_api_key(api_key)
            return user['api_secret'] == api_secret
        except Exception:
            return False

    def create_user(self, username: str, password: str, roles: List[str]) -> Dict[str, Any]:
        """Create a new user.

        Args:
            username: Username.
            password: Password.
            roles: List of role names.

        Returns:
            User information.

        Raises:
            Exception: If username already exists.
        """
        # Check if username already exists
        for user in self.users.values():
            if user['username'] == username:
                raise Exception(f"Username already exists: {username}")

        # Create new user
        user_id = f"user-{uuid.uuid4()}"
        api_key = f"api-{uuid.uuid4()}"
        api_secret = f"secret-{uuid.uuid4()}"

        user = {
            'user_id': user_id,
            'username': username,
            'password_hash': self._hash_password(password),
            'roles': roles,
            'api_key': api_key,
            'api_secret': api_secret,
        }

        self.users[user_id] = user
        return user

    def update_user(self, user_id: str, **kwargs) -> Dict[str, Any]:
        """Update a user.

        Args:
            user_id: User ID.
            **kwargs: Fields to update.

        Returns:
            Updated user information.

        Raises:
            Exception: If user not found.
        """
        if user_id not in self.users:
            raise Exception(f"User not found: {user_id}")

        user = self.users[user_id]
        
        # Update fields
        for key, value in kwargs.items():
            if key == 'password':
                user['password_hash'] = self._hash_password(value)
            else:
                user[key] = value

        return user

    def delete_user(self, user_id: str) -> bool:
        """Delete a user.

        Args:
            user_id: User ID.

        Returns:
            True if user was deleted, False otherwise.
        """
        if user_id in self.users:
            del self.users[user_id]
            return True
        return False
