"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.
"""

from typing import List, Optional
from .authentication import Authentication

class Authorization:
    """Handles user authorization and role-based access control."""
    
    def __init__(self, auth: Authentication):
        """Initialize the authorization system.
        
        Args:
            auth: Authentication instance
        """
        self.auth = auth
        
    def has_role(self, token: str, role: str) -> bool:
        """Check if a user has a specific role.
        
        Args:
            token: JWT token
            role: Role to check
            
        Returns:
            True if user has the role, False otherwise
        """
        payload = self.auth.verify_token(token)
        if not payload:
            return False
            
        return role in payload.get('roles', [])
        
    def has_any_role(self, token: str, roles: List[str]) -> bool:
        """Check if a user has any of the specified roles.
        
        Args:
            token: JWT token
            roles: List of roles to check
            
        Returns:
            True if user has any of the roles, False otherwise
        """
        payload = self.auth.verify_token(token)
        if not payload:
            return False
            
        user_roles = payload.get('roles', [])
        return any(role in user_roles for role in roles)
        
    def has_all_roles(self, token: str, roles: List[str]) -> bool:
        """Check if a user has all of the specified roles.
        
        Args:
            token: JWT token
            roles: List of roles to check
            
        Returns:
            True if user has all of the roles, False otherwise
        """
        payload = self.auth.verify_token(token)
        if not payload:
            return False
            
        user_roles = payload.get('roles', [])
        return all(role in user_roles for role in roles) 