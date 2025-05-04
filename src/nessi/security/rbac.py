from typing import List, Dict, Optional
from enum import Enum
import logging
from functools import wraps
from fastapi import HTTPException, Security
from fastapi.security import APIKeyHeader

class Role(Enum):
    """User roles with different permission levels."""
    ADMIN = "admin"
    MANAGER = "manager"
    ANALYST = "analyst"
    VIEWER = "viewer"

class Permission(Enum):
    """Available permissions in the system."""
    READ_DATA = "read_data"
    WRITE_DATA = "write_data"
    DELETE_DATA = "delete_data"
    MANAGE_USERS = "manage_users"
    VIEW_REPORTS = "view_reports"
    GENERATE_REPORTS = "generate_reports"
    MANAGE_ALERTS = "manage_alerts"
    CONFIGURE_SYSTEM = "configure_system"

class RBAC:
    """Role-Based Access Control system."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.role_permissions = {
            Role.ADMIN: [
                Permission.READ_DATA,
                Permission.WRITE_DATA,
                Permission.DELETE_DATA,
                Permission.MANAGE_USERS,
                Permission.VIEW_REPORTS,
                Permission.GENERATE_REPORTS,
                Permission.MANAGE_ALERTS,
                Permission.CONFIGURE_SYSTEM
            ],
            Role.MANAGER: [
                Permission.READ_DATA,
                Permission.WRITE_DATA,
                Permission.VIEW_REPORTS,
                Permission.GENERATE_REPORTS,
                Permission.MANAGE_ALERTS
            ],
            Role.ANALYST: [
                Permission.READ_DATA,
                Permission.VIEW_REPORTS,
                Permission.GENERATE_REPORTS
            ],
            Role.VIEWER: [
                Permission.READ_DATA,
                Permission.VIEW_REPORTS
            ]
        }
        
    def has_permission(self, role: Role, permission: Permission) -> bool:
        """Check if a role has a specific permission."""
        return permission in self.role_permissions.get(role, [])
    
    def get_role_permissions(self, role: Role) -> List[Permission]:
        """Get all permissions for a specific role."""
        return self.role_permissions.get(role, [])
    
    def require_permission(self, permission: Permission):
        """Decorator to require a specific permission for an endpoint."""
        def decorator(func):
            @wraps(func)
            async def wrapper(*args, **kwargs):
                api_key = kwargs.get("api_key")
                if not api_key:
                    raise HTTPException(status_code=401, detail="API key required")
                
                # Get user role from API key (implementation depends on your auth system)
                user_role = self._get_role_from_api_key(api_key)
                if not self.has_permission(user_role, permission):
                    raise HTTPException(
                        status_code=403,
                        detail=f"Permission {permission.value} required"
                    )
                
                return await func(*args, **kwargs)
            return wrapper
        return decorator
    
    def _get_role_from_api_key(self, api_key: str) -> Role:
        """Get user role from API key (to be implemented based on your auth system)."""
        # This is a placeholder - implement based on your authentication system
        # For example, you might look up the role in a database
        return Role.VIEWER  # Default to viewer role for now

# API key security
api_key_header = APIKeyHeader(name="X-API-Key")

def get_api_key(api_key: str = Security(api_key_header)) -> str:
    """Dependency to get and validate API key."""
    return api_key 