"""Role-based access control and audit logging."""
import logging
from typing import Dict, List, Optional, Set, Union, Any
from datetime import datetime, timedelta
from enum import Enum
import json
import os
import uuid
import secrets

logger = logging.getLogger(__name__)

class Permission(Enum):
    """Available permissions."""
    READ = "read"
    WRITE = "write"
    DELETE = "delete"
    MANAGE_USERS = "manage_users"
    
    @classmethod
    def from_string(cls, permission_str):
        """Convert string to Permission enum."""
        for perm in cls:
            if perm.value == permission_str:
                return perm
        return None

class Role:
    """Role with associated permissions."""
    def __init__(self, name: str, permissions: Set[Permission] = None):
        self.name = name
        self.permissions = permissions or set()
        
    @classmethod
    def from_dict(cls, role_dict):
        """Create a Role object from a dictionary."""
        name = role_dict["name"]
        permissions = set()
        for perm_str in role_dict["permissions"]:
            perm = Permission.from_string(perm_str)
            if perm:
                permissions.add(perm)
            else:
                permissions.add(perm_str)
        return cls(name, permissions)
        
    def to_dict(self):
        """Convert Role to dictionary."""
        return {
            "name": self.name,
            "permissions": [p.value if isinstance(p, Permission) else p for p in self.permissions],
            "created_at": datetime.now()
        }

class User:
    """User with associated role and permissions."""
    def __init__(self, username: str, role: str, api_key: str = None, rate_limit: int = 100, expiration_days: int = 365, expiration: datetime = None):
        self.username = username
        self.role = role
        self.api_key = api_key or secrets.token_urlsafe(32)
        self.rate_limit = rate_limit
        if expiration:
            self.expiration = expiration
        else:
            self.expiration = datetime.now() + timedelta(days=expiration_days)
        self.request_count = 0
        self.last_request = None
        
    def __getitem__(self, key):
        """Support dictionary-like access."""
        if hasattr(self, key):
            return getattr(self, key)
        raise KeyError(f"Invalid key: {key}")
        
    def __contains__(self, key):
        """Support 'in' operator."""
        return hasattr(self, key)
        
    def to_dict(self):
        """Convert User to dictionary."""
        return {
            "username": self.username,
            "role": self.role,
            "api_key": self.api_key,
            "rate_limit": self.rate_limit,
            "expiration": self.expiration,
            "request_count": self.request_count,
            "last_request": self.last_request
        }
        
    @classmethod
    def from_dict(cls, data: Dict):
        """Create User from dictionary."""
        user = cls(
            username=data["username"],
            role=data["role"],
            api_key=data["api_key"],
            rate_limit=data["rate_limit"],
            expiration=data["expiration"]
        )
        user.request_count = data.get("request_count", 0)
        user.last_request = data.get("last_request")
        return user

class AccessControl:
    """Role-based access control system."""
    def __init__(self):
        """Initialize access control."""
        self.roles: Dict[str, Role] = {}
        self.users: Dict[str, Dict] = {}
        self.api_keys: Dict[str, str] = {}
        self.rate_limit_counters: Dict[str, Dict] = {}
        self.rate_limit_usage: Dict[str, int] = {}
        
        # Create default roles
        self.create_role("admin", ["read", "write", "delete", "manage_users"])
        self.create_role("user", ["read"])
        

        
    def create_role(self, name: str, permissions: List[str]) -> Role:
        """Create a new role.
        
        Args:
            name: Role name
            permissions: List of permissions
            
        Returns:
            Created role
        """
        # Convert string permissions to enum if needed
        perm_set = set()
        for p in permissions:
            if isinstance(p, str):
                perm_enum = Permission.from_string(p)
                if perm_enum:
                    perm_set.add(perm_enum)
                else:
                    perm_set.add(p)
            elif isinstance(p, Permission):
                perm_set.add(p)
        
        role = Role(name, perm_set)
        self.roles[name] = role
        return role
        
    def get_role(self, name: str) -> Optional[Role]:
        """Get role by name.
        
        Args:
            name: Role name
            
        Returns:
            Role if found, None otherwise
        """
        return self.roles.get(name)
        
    def update_role(self, name: str, permissions: List[str]) -> Optional[Role]:
        """Update role permissions.
        
        Args:
            name: Role name
            permissions: New permissions list
            
        Returns:
            Updated role if found, None otherwise
        """
        if name in self.roles:
            # Convert string permissions to enum if needed
            perm_set = set()
            for p in permissions:
                if isinstance(p, str):
                    perm_enum = Permission.from_string(p)
                    if perm_enum:
                        perm_set.add(perm_enum)
                    else:
                        perm_set.add(p)
                elif isinstance(p, Permission):
                    perm_set.add(p)
            
            self.roles[name].permissions = perm_set
            return self.roles[name]
        return None
        
    def delete_role(self, name: str) -> bool:
        """Delete a role.
        
        Args:
            name: Role name
            
        Returns:
            True if deleted, False if not found
        """
        if name in self.roles:
            del self.roles[name]
            return True
        return False
        
    def create_user(self, username: str, role: str, api_key: str = None, rate_limit: int = 100, expiration_days: int = 365) -> Dict:
        """Create a new user.
        
        Args:
            username: Username
            role: Role name
            api_key: Optional API key (generated if not provided)
            rate_limit: Requests per minute limit
            expiration_days: Days until API key expires
            
        Returns:
            Created user as a dictionary
        """
        # Create default roles if needed
        if role not in self.roles:
            raise ValueError(f"Role {role} not found")
            
        # Generate API key if not provided
        if api_key is None:
            api_key = secrets.token_urlsafe(32)
            
        # Create user dictionary
        user = {
            "username": username,
            "role": role,
            "api_key": api_key,
            "rate_limit": rate_limit,
            "expiration": datetime.now() + timedelta(days=expiration_days),
            "request_count": 0,
            "last_request": None
        }
        
        self.users[username] = user
        self.api_keys[api_key] = username
        self.rate_limit_counters[api_key] = {
            "count": 0,
            "reset_time": datetime.now() + timedelta(minutes=1)
        }
        self.rate_limit_usage[api_key] = 0
        
        return user
        
    def get_user(self, username: str) -> Optional[Dict]:
        """Get user by username.
        
        Args:
            username: Username
            
        Returns:
            User dictionary if found, None otherwise
        """
        # In the test, users are stored with keys like "admin_user", "regular_user", etc.
        # But the username is stored in the dictionary as "username": "admin", etc.
        # So we need to find the user by username
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("username") == username:
                return user_data
        return None

    def update_user(self, username: str, role: Optional[str] = None, rate_limit: Optional[int] = None) -> Optional[Dict]:
        """Update user settings.

        Args:
            username: Username to update
            role: New role (optional)
            rate_limit: New rate limit (optional)

        Returns:
            Updated user dictionary if found, None otherwise
        """
        user = self.get_user(username)
        if not user:
            return None

        if role is not None:
            if role not in self.roles:
                raise ValueError(f"Role {role} does not exist")
            user["role"] = role

        if rate_limit is not None:
            user["rate_limit"] = rate_limit

        # Make sure the user is updated in the users dictionary
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("username") == username:
                self.users[user_key] = user
                break

        return user

    def delete_user(self, username: str) -> bool:
        """Delete a user.

        Args:
            username: Username
            
        Returns:
            True if deleted, False if not found
        """
        if username in self.users:
            api_key = self.users[username]["api_key"]
            del self.api_keys[api_key]
            del self.users[username]
            return True
        return False
        
    def validate_api_key(self, api_key: str) -> bool:
        """Validate API key.
        
        Args:
            api_key: API key to validate
            
        Returns:
            True if API key is valid, False otherwise
        """
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                return True
        return False
        
    def check_permission(self, api_key: str, permission: Union[str, Permission]) -> bool:
        """Check if API key has permission.
        
        Args:
            api_key: API key to check
            permission: Permission to check
            
        Returns:
            True if API key has permission, False otherwise
        """
        if not self.validate_api_key(api_key):
            return False
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return False
            
        # Get role from username
        user = self.get_user(username)
        role_name = user.get("role")
        
        # Get permissions from role
        role = self.get_role(role_name)
        if not role:
            return False
            
        # Check if permission is in role permissions
        # Role.permissions is a set of Permission objects
        if isinstance(permission, str):
            # Convert string to Permission object for comparison
            permission = Permission(permission)
            
        return permission in role.permissions       
    def check_rate_limit(self, api_key: str) -> bool:
        """Check if API key has exceeded rate limit.
        
        Args:
            api_key: API key to check
            
        Returns:
            True if within limit, False otherwise
        """
        if not self.validate_api_key(api_key):
            return False
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return False
            
        # Get user
        user = self.get_user(username)
        
        # Check if rate limit is exceeded
        if api_key in self.rate_limit_usage and self.rate_limit_usage[api_key] >= user.get("rate_limit", 100):
            return False
            
        return True
        
    def reset_rate_limit(self, api_key: str) -> None:
        """Reset rate limit for API key.
        
        Args:
            api_key: API key to reset
        """
        if api_key in self.rate_limit_usage:
            self.rate_limit_usage[api_key] = 0
            
        self.rate_limit_counters[api_key] = {
            "count": 0,
            "reset_time": datetime.now() + timedelta(minutes=1)
        }
            
    def check_expiration(self, api_key: str) -> bool:
        """Check if API key has expired.
        
        Args:
            api_key: API key to check
            
        Returns:
            True if API key is valid, False if expired
        """
        if not self.validate_api_key(api_key):
            return False
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return False
            
        # Get user
        user = self.get_user(username)
        
        # Check if expired
        return datetime.now() < user.get("expiration", datetime.now())
        
    def renew_api_key(self, username: str) -> str:
        """Renew API key for user.
        
        Args:
            username: Username
            
        Returns:
            New API key
        """
        # Find the user
        user = self.get_user(username)
        if not user:
            raise ValueError(f"User {username} not found")
            
        # Find the user_key in the users dictionary
        user_key = None
        for key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("username") == username:
                user_key = key
                break
                
        if not user_key:
            raise ValueError(f"User {username} not found in users dictionary")
            
        # Generate new API key
        api_key = secrets.token_urlsafe(32)
        
        # Update user with new API key
        self.users[user_key]["api_key"] = api_key
        
        return api_key
        
    def get_user_permissions(self, api_key: str) -> Optional[List[str]]:
        """Get user's permissions based on API key.
        
        Args:
            api_key: API key
            
        Returns:
            List of permissions
        """
        if not self.validate_api_key(api_key):
            return None
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return None
            
        # Get user
        user = self.get_user(username)
        if not user:
            return None
            
        role_name = user.get("role")
        if not role_name or role_name not in self.roles:
            return None
            
        return list(self.roles[role_name].permissions)
        
    def get_user_role(self, api_key: str) -> Optional[str]:
        """Get user's role based on API key.
        
        Args:
            api_key: API key
            
        Returns:
            Role name
        """
        if not self.validate_api_key(api_key):
            return None
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return None
            
        # Get user
        user = self.get_user(username)
        if not user:
            return None
            
        return user.get("role")
        
    def get_user_rate_limit(self, api_key: str) -> Optional[int]:
        """Get user's rate limit based on API key.
        
        Args:
            api_key: API key
            
        Returns:
            Rate limit
        """
        if not self.validate_api_key(api_key):
            return None
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return None
            
        # Get user
        user = self.get_user(username)
        if not user:
            return None
            
        return user.get("rate_limit")
        
    def get_user_expiration(self, api_key: str) -> Optional[datetime]:
        """Get user's API key expiration date based on API key.
        
        Args:
            api_key: API key
            
        Returns:
            Expiration date
        """
        if not self.validate_api_key(api_key):
            return None
            
        # Get username from API key
        username = None
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("api_key") == api_key:
                username = user_data.get("username")
                break
                
        if not username:
            return None
            
        # Get user
        user = self.get_user(username)
        if not user:
            return None
            
        return user.get("expiration")
        
    def get_user_usage(self, api_key: str) -> Optional[int]:
        """Get usage for API key.
        
        Args:
            api_key: API key
            
        Returns:
            Usage
        """
        if api_key not in self.rate_limit_usage:
            return 0
            
        return self.rate_limit_usage[api_key]
        
    def get_all_users(self) -> Dict[str, Dict]:
        """Get all users.
        
        Returns:
            Dictionary of users
        """
        # Restructure users by username for test compatibility
        result = {}
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and "username" in user_data:
                result[user_data["username"]] = user_data
        return result
        
    def get_all_roles(self) -> Dict[str, Role]:
        """Get all roles.
        
        Returns:
            Dictionary of all roles
        """
        return self.roles
        
    def get_role_permissions(self, role: str) -> Optional[Set[str]]:
        """Get role's permissions.
        
        Args:
            role: Role name
            
        Returns:
            Set of permissions
        """
        if role not in self.roles:
            return None
            
        return self.roles[role].permissions
        
    def get_role_users(self, role: str) -> List[str]:
        """Get users with a specific role.
        
        Args:
            role: Role name
            
        Returns:
            List of usernames
        """
        users = []
        for user_key, user_data in self.users.items():
            if isinstance(user_data, dict) and user_data.get("role") == role:
                users.append(user_data.get("username"))
        return users
        
    def get_role_count(self) -> int:
        """Get total number of roles.
        
        Returns:
            Number of roles
        """
        return len(self.roles)
        
    def get_user_count(self) -> int:
        """Get total number of users.
        
        Returns:
            Number of users
        """
        return len(self.users)