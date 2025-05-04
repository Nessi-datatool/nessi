"""Role-based access control and audit logging."""
import logging
from typing import Dict, List, Optional, Set, Any
from datetime import datetime
from enum import Enum
import json
import os

logger = logging.getLogger(__name__)

class Permission(Enum):
    """Available permissions."""
    READ = "read"
    WRITE = "write"
    DELETE = "delete"
    ADMIN = "admin"

class Role:
    """Role with associated permissions."""
    def __init__(self, name: str, permissions: Set[Permission]):
        self.name = name
        self.permissions = permissions

class User:
    """User with associated roles."""
    def __init__(self, username: str, roles: List[Role]):
        self.username = username
        self.roles = roles

class AccessControl:
    """Role-based access control system."""
    def __init__(self, audit_log_path: str = "logs/audit.log"):
        self.users: Dict[str, User] = {}
        self.roles: Dict[str, Role] = {}
        self.audit_log_path = audit_log_path
        
        # Create audit log directory if it doesn't exist
        os.makedirs(os.path.dirname(audit_log_path), exist_ok=True)
        
        # Initialize default roles
        self._init_default_roles()

    def _init_default_roles(self) -> None:
        """Initialize default roles."""
        # Admin role
        admin_role = Role("admin", {
            Permission.READ,
            Permission.WRITE,
            Permission.DELETE,
            Permission.ADMIN
        })
        self.roles["admin"] = admin_role
        
        # Editor role
        editor_role = Role("editor", {
            Permission.READ,
            Permission.WRITE
        })
        self.roles["editor"] = editor_role
        
        # Viewer role
        viewer_role = Role("viewer", {
            Permission.READ
        })
        self.roles["viewer"] = viewer_role

    def create_role(self, name: str, permissions: Set[Permission]) -> Role:
        """Create a new role.
        
        Args:
            name: Name of the role
            permissions: Set of permissions for the role
            
        Returns:
            Created role
        """
        try:
            role = Role(name, permissions)
            self.roles[name] = role
            
            self._log_audit(
                "create_role",
                {"role": name, "permissions": [p.value for p in permissions]}
            )
            
            return role
            
        except Exception as e:
            logger.error(f"Error creating role: {str(e)}")
            raise

    def create_user(self, username: str, roles: List[str]) -> User:
        """Create a new user.
        
        Args:
            username: Username
            roles: List of role names
            
        Returns:
            Created user
        """
        try:
            user_roles = [self.roles[role] for role in roles]
            user = User(username, user_roles)
            self.users[username] = user
            
            self._log_audit(
                "create_user",
                {"username": username, "roles": roles}
            )
            
            return user
            
        except Exception as e:
            logger.error(f"Error creating user: {str(e)}")
            raise

    def check_permission(self, username: str, permission: Permission) -> bool:
        """Check if a user has a specific permission.
        
        Args:
            username: Username
            permission: Permission to check
            
        Returns:
            True if user has permission, False otherwise
        """
        try:
            if username not in self.users:
                return False
                
            user = self.users[username]
            
            # Check if any of the user's roles have the permission
            for role in user.roles:
                if permission in role.permissions:
                    return True
                    
            return False
            
        except Exception as e:
            logger.error(f"Error checking permission: {str(e)}")
            raise

    def get_user_permissions(self, username: str) -> Set[Permission]:
        """Get all permissions for a user.
        
        Args:
            username: Username
            
        Returns:
            Set of permissions
        """
        try:
            if username not in self.users:
                return set()
                
            user = self.users[username]
            permissions = set()
            
            for role in user.roles:
                permissions.update(role.permissions)
                
            return permissions
            
        except Exception as e:
            logger.error(f"Error getting user permissions: {str(e)}")
            raise

    def update_user_roles(self, username: str, roles: List[str]) -> None:
        """Update roles for a user.
        
        Args:
            username: Username
            roles: List of new role names
        """
        try:
            if username not in self.users:
                raise ValueError(f"User {username} not found")
                
            user_roles = [self.roles[role] for role in roles]
            self.users[username].roles = user_roles
            
            self._log_audit(
                "update_user_roles",
                {"username": username, "roles": roles}
            )
            
        except Exception as e:
            logger.error(f"Error updating user roles: {str(e)}")
            raise

    def _log_audit(self, action: str, details: Dict[str, Any]) -> None:
        """Log an audit event.
        
        Args:
            action: Action performed
            details: Details of the action
        """
        try:
            log_entry = {
                "timestamp": datetime.now().isoformat(),
                "action": action,
                "details": details
            }
            
            with open(self.audit_log_path, "a") as f:
                f.write(json.dumps(log_entry) + "\n")
                
        except Exception as e:
            logger.error(f"Error logging audit event: {str(e)}")
            raise

    def get_audit_log(self, start_time: Optional[datetime] = None,
                     end_time: Optional[datetime] = None) -> List[Dict[str, Any]]:
        """Get audit log entries.
        
        Args:
            start_time: Start time for filtering
            end_time: End time for filtering
            
        Returns:
            List of audit log entries
        """
        try:
            entries = []
            
            with open(self.audit_log_path, "r") as f:
                for line in f:
                    entry = json.loads(line)
                    entry_time = datetime.fromisoformat(entry["timestamp"])
                    
                    if start_time and entry_time < start_time:
                        continue
                        
                    if end_time and entry_time > end_time:
                        continue
                        
                    entries.append(entry)
                    
            return entries
            
        except Exception as e:
            logger.error(f"Error getting audit log: {str(e)}")
            raise 