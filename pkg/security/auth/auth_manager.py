"""
Authentication manager for Nessi.dev.
"""

import time
import jwt
from typing import Dict, Any, Optional


class AuthManager:
    """Authentication manager for Nessi.dev."""

    def __init__(self, user_store=None, jwt_secret="default-secret", token_expiry=3600):
        """Initialize the authentication manager.

        Args:
            user_store: User store for user management.
            jwt_secret: Secret for JWT token generation and validation.
            token_expiry: Token expiry time in seconds.
        """
        self.user_store = user_store
        self.jwt_secret = jwt_secret
        self.token_expiry = token_expiry

    def authenticate(self, username: str, password: str) -> Dict[str, Any]:
        """Authenticate a user with username and password.

        Args:
            username: Username.
            password: Password.

        Returns:
            Dict with user information and token.

        Raises:
            Exception: If authentication fails.
        """
        if not self.user_store.verify_user(username, password):
            raise Exception("Authentication failed")

        user = self.user_store.get_user(username)
        token = self.generate_token(user)

        return {
            "user_id": user["user_id"],
            "username": user["username"],
            "roles": user["roles"],
            "token": token,
        }

    def generate_token(self, user: Dict[str, Any]) -> str:
        """Generate a JWT token for a user.

        Args:
            user: User information.

        Returns:
            JWT token.
        """
        payload = {
            "sub": user["user_id"],
            "username": user["username"],
            "roles": user["roles"],
            "exp": int(time.time()) + self.token_expiry,
        }
        return jwt.encode(payload, self.jwt_secret, algorithm="HS256")

    def validate_token(self, token: str) -> Dict[str, Any]:
        """Validate a JWT token.

        Args:
            token: JWT token.

        Returns:
            Dict with user information.

        Raises:
            jwt.ExpiredSignatureError: If token is expired.
            jwt.InvalidTokenError: If token is invalid.
        """
        try:
            payload = jwt.decode(token, self.jwt_secret, algorithms=["HS256"])
            return {
                "user_id": payload["sub"],
                "username": payload["username"],
                "roles": payload["roles"],
            }
        except jwt.ExpiredSignatureError:
            raise jwt.ExpiredSignatureError("Token expired")
        except jwt.InvalidTokenError:
            raise jwt.InvalidTokenError("Invalid token")

    def authenticate_api_key(self, api_key: str, api_secret: str) -> Dict[str, Any]:
        """Authenticate a user with API key and secret.

        Args:
            api_key: API key.
            api_secret: API secret.

        Returns:
            Dict with user information.

        Raises:
            Exception: If authentication fails.
        """
        if not self.user_store.verify_api_key(api_key, api_secret):
            raise Exception("API key authentication failed")

        user = self.user_store.get_user_by_api_key(api_key)
        return {
            "user_id": user["user_id"],
            "username": user["username"],
            "roles": user["roles"],
        }

    def authenticate_oauth(self, code: str, redirect_uri: str, client_id: str, client_secret: str) -> Dict[str, Any]:
        """Authenticate a user with OAuth 2.0.

        Args:
            code: Authorization code.
            redirect_uri: Redirect URI.
            client_id: Client ID.
            client_secret: Client secret.

        Returns:
            Dict with user information and tokens.
        """
        # This is a mock implementation
        return {
            "user_id": "test-user",
            "username": "testuser",
            "roles": ["user"],
            "access_token": "test-access-token",
            "refresh_token": "test-refresh-token",
            "expires_in": 3600,
        }

    def refresh_token(self, refresh_token: str, client_id: str, client_secret: str) -> Dict[str, Any]:
        """Refresh an OAuth 2.0 token.

        Args:
            refresh_token: Refresh token.
            client_id: Client ID.
            client_secret: Client secret.

        Returns:
            Dict with new tokens.
        """
        # This is a mock implementation
        return {
            "access_token": "new-access-token",
            "refresh_token": "new-refresh-token",
            "expires_in": 3600,
        }
