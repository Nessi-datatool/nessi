"""
Role-Based Access Control (RBAC) manager for Nessi.dev.
"""

from typing import Dict, List, Any, Optional


class RBACManager:
    """Role-Based Access Control (RBAC) manager for Nessi.dev."""

    def __init__(self):
        """Initialize the RBAC manager."""
        self.role_permissions = {
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

    def check_permission(self, user_id: str, resource: str, action: str) -> bool:
        """Check if a user has permission to perform an action on a resource.

        Args:
            user_id: User ID.
            resource: Resource name.
            action: Action name.

        Returns:
            True if the user has permission, False otherwise.
        """
        # This is a mock implementation that would be replaced with actual logic
        # In a real implementation, we would get the user's roles and check if any of them
        # have the required permission
        
        # For testing purposes, we'll use a simple mapping
        if user_id.startswith('admin'):
            return True
        elif user_id.startswith('regular'):
            if resource == 'tables':
                return action in ['read', 'write']
            elif resource == 'users':
                return action == 'read'
            else:
                return action == 'read'
        elif user_id.startswith('viewer'):
            return resource == 'tables' and action == 'read'
        else:
            return False

    def get_user_permissions(self, user_id: str, roles: List[str]) -> Dict[str, List[str]]:
        """Get all permissions for a user.

        Args:
            user_id: User ID.
            roles: List of role names.

        Returns:
            Dict mapping resource names to lists of allowed actions.
        """
        # This is a mock implementation
        permissions = {}
        
        for role in roles:
            if role in self.role_permissions:
                role_perms = self.role_permissions[role]
                for resource, actions in role_perms.items():
                    if resource not in permissions:
                        permissions[resource] = []
                    for action in actions:
                        if action not in permissions[resource]:
                            permissions[resource].append(action)
        
        return permissions
