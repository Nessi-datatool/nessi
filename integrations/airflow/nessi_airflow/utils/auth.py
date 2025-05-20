"""
Authentication utilities for Nessi Airflow integration.

This module provides utilities for authenticating with the Nessi.dev API
using various authentication methods.
"""

import os
import json
import time
import logging
from typing import Dict, Any, Optional, Tuple

import requests
from airflow.exceptions import AirflowException
from airflow.hooks.base import BaseHook
from airflow.configuration import conf

from nessi_airflow.utils.error_handling import (
    NessiAuthenticationError,
    handle_api_exceptions,
    with_exponential_backoff,
)

logger = logging.getLogger(__name__)


class NessiAuthBase:
    """Base class for Nessi authentication methods."""
    
    def __init__(self, base_url: str):
        """
        Initialize the authentication handler.
        
        Args:
            base_url: Base URL for the Nessi.dev API
        """
        self.base_url = base_url.rstrip('/')
    
    def get_auth_headers(self) -> Dict[str, str]:
        """
        Get authentication headers for API requests.
        
        Returns:
            Dictionary of authentication headers
        """
        raise NotImplementedError("Subclasses must implement get_auth_headers")
    
    def update_session(self, session: requests.Session) -> None:
        """
        Update a requests session with authentication headers.
        
        Args:
            session: Requests session to update
        """
        session.headers.update(self.get_auth_headers())


class ApiKeyAuth(NessiAuthBase):
    """API key authentication for Nessi.dev."""
    
    def __init__(self, base_url: str, api_key: str, api_secret: Optional[str] = None):
        """
        Initialize API key authentication.
        
        Args:
            base_url: Base URL for the Nessi.dev API
            api_key: API key for authentication
            api_secret: API secret for authentication (optional)
        """
        super().__init__(base_url)
        self.api_key = api_key
        self.api_secret = api_secret
    
    def get_auth_headers(self) -> Dict[str, str]:
        """
        Get authentication headers for API requests.
        
        Returns:
            Dictionary of authentication headers
        """
        headers = {"X-API-Key": self.api_key}
        if self.api_secret:
            headers["X-API-Secret"] = self.api_secret
        return headers


class OAuth2Auth(NessiAuthBase):
    """OAuth 2.0 authentication for Nessi.dev."""
    
    def __init__(
        self,
        base_url: str,
        client_id: str,
        client_secret: str,
        token_url: Optional[str] = None,
        refresh_token: Optional[str] = None,
        access_token: Optional[str] = None,
        expires_at: Optional[float] = None,
    ):
        """
        Initialize OAuth 2.0 authentication.
        
        Args:
            base_url: Base URL for the Nessi.dev API
            client_id: OAuth client ID
            client_secret: OAuth client secret
            token_url: URL for token endpoint (defaults to {base_url}/oauth/token)
            refresh_token: Refresh token (if available)
            access_token: Access token (if available)
            expires_at: Timestamp when the access token expires (if available)
        """
        super().__init__(base_url)
        self.client_id = client_id
        self.client_secret = client_secret
        self.token_url = token_url or f"{base_url}/oauth/token"
        self.refresh_token = refresh_token
        self.access_token = access_token
        self.expires_at = expires_at
    
    @with_exponential_backoff()
    @handle_api_exceptions
    def _request_token(self, grant_type: str, **kwargs) -> Dict[str, Any]:
        """
        Request a token from the OAuth server.
        
        Args:
            grant_type: OAuth grant type
            **kwargs: Additional parameters for the token request
            
        Returns:
            Token response
        """
        data = {
            "grant_type": grant_type,
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            **kwargs,
        }
        
        response = requests.post(
            self.token_url,
            data=data,
            headers={"Content-Type": "application/x-www-form-urlencoded"},
            timeout=30,
        )
        response.raise_for_status()
        return response.json()
    
    def get_client_credentials_token(self) -> Dict[str, Any]:
        """
        Get a token using the client credentials grant.
        
        Returns:
            Token response
        """
        return self._request_token("client_credentials")
    
    def refresh_access_token(self) -> Dict[str, Any]:
        """
        Refresh the access token using the refresh token.
        
        Returns:
            Token response
        """
        if not self.refresh_token:
            raise NessiAuthenticationError("No refresh token available")
        
        token_response = self._request_token(
            "refresh_token",
            refresh_token=self.refresh_token,
        )
        
        # Update tokens
        self.access_token = token_response.get("access_token")
        if "refresh_token" in token_response:
            self.refresh_token = token_response.get("refresh_token")
        if "expires_in" in token_response:
            self.expires_at = time.time() + token_response.get("expires_in", 3600)
        
        return token_response
    
    def ensure_valid_token(self) -> None:
        """
        Ensure that we have a valid access token.
        
        If the token is expired or not available, request a new one.
        """
        # Check if we need a new token
        if not self.access_token or (self.expires_at and time.time() >= self.expires_at):
            logger.info("Access token expired or not available, requesting new token")
            
            # Try to refresh if we have a refresh token
            if self.refresh_token:
                try:
                    self.refresh_access_token()
                    return
                except Exception as e:
                    logger.warning("Failed to refresh token: %s", e)
            
            # Fall back to client credentials
            token_response = self.get_client_credentials_token()
            self.access_token = token_response.get("access_token")
            self.refresh_token = token_response.get("refresh_token")
            if "expires_in" in token_response:
                self.expires_at = time.time() + token_response.get("expires_in", 3600)
    
    def get_auth_headers(self) -> Dict[str, str]:
        """
        Get authentication headers for API requests.
        
        Returns:
            Dictionary of authentication headers
        """
        self.ensure_valid_token()
        return {"Authorization": f"Bearer {self.access_token}"}


def get_auth_from_connection(
    conn_id: str,
    base_url: Optional[str] = None,
) -> NessiAuthBase:
    """
    Create an authentication handler from an Airflow connection.
    
    Args:
        conn_id: Airflow connection ID
        base_url: Override base URL (optional)
        
    Returns:
        Authentication handler
    """
    conn = BaseHook.get_connection(conn_id)
    
    # Get base URL
    if not base_url:
        base_url = conn.host
        if not base_url:
            raise AirflowException(f"No host specified for connection {conn_id}")
    
    # Parse extra configuration
    extra_config = {}
    if conn.extra:
        try:
            extra_config = json.loads(conn.extra)
        except json.JSONDecodeError:
            logger.warning(
                "Failed to parse extra config for connection %s", conn_id
            )
    
    # Check for OAuth configuration
    auth_type = extra_config.get("auth_type", "api_key").lower()
    
    if auth_type == "oauth2":
        client_id = extra_config.get("client_id")
        client_secret = extra_config.get("client_secret")
        token_url = extra_config.get("token_url")
        refresh_token = extra_config.get("refresh_token")
        access_token = extra_config.get("access_token")
        expires_at = extra_config.get("expires_at")
        
        if not client_id or not client_secret:
            raise AirflowException(
                f"OAuth2 authentication requires client_id and client_secret in connection {conn_id}"
            )
        
        return OAuth2Auth(
            base_url=base_url,
            client_id=client_id,
            client_secret=client_secret,
            token_url=token_url,
            refresh_token=refresh_token,
            access_token=access_token,
            expires_at=expires_at,
        )
    else:
        # Default to API key authentication
        api_key = conn.login
        api_secret = conn.password
        
        if not api_key:
            raise AirflowException(
                f"API key authentication requires login (API key) in connection {conn_id}"
            )
        
        return ApiKeyAuth(
            base_url=base_url,
            api_key=api_key,
            api_secret=api_secret,
        )


def get_auth_from_env() -> Tuple[Optional[NessiAuthBase], Optional[str]]:
    """
    Create an authentication handler from environment variables.
    
    Returns:
        Tuple of (authentication handler, base URL) or (None, None) if not configured
    """
    # Check for base URL
    base_url = os.environ.get("NESSI_API_HOST")
    if not base_url:
        return None, None
    
    # Check for OAuth configuration
    client_id = os.environ.get("NESSI_CLIENT_ID")
    client_secret = os.environ.get("NESSI_CLIENT_SECRET")
    
    if client_id and client_secret:
        token_url = os.environ.get("NESSI_TOKEN_URL")
        refresh_token = os.environ.get("NESSI_REFRESH_TOKEN")
        access_token = os.environ.get("NESSI_ACCESS_TOKEN")
        expires_at_str = os.environ.get("NESSI_TOKEN_EXPIRES_AT")
        
        expires_at = None
        if expires_at_str:
            try:
                expires_at = float(expires_at_str)
            except ValueError:
                logger.warning("Invalid NESSI_TOKEN_EXPIRES_AT value: %s", expires_at_str)
        
        return OAuth2Auth(
            base_url=base_url,
            client_id=client_id,
            client_secret=client_secret,
            token_url=token_url,
            refresh_token=refresh_token,
            access_token=access_token,
            expires_at=expires_at,
        ), base_url
    
    # Check for API key configuration
    api_key = os.environ.get("NESSI_API_KEY")
    api_secret = os.environ.get("NESSI_API_SECRET")
    
    if api_key:
        return ApiKeyAuth(
            base_url=base_url,
            api_key=api_key,
            api_secret=api_secret,
        ), base_url
    
    return None, None
