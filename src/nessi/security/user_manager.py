from typing import Optional, Dict, List
import logging
from datetime import datetime
import secrets
import hashlib
from .rbac import Role, RBAC

class User:
    """User model with authentication and role information."""
    
    def __init__(self, username: str, email: str, role: Role, api_key: Optional[str] = None):
        self.username = username
        self.email = email
        self.role = role
        self.api_key = api_key or self._generate_api_key()
        self.created_at = datetime.utcnow()
        self.last_login = None
        self.is_active = True
    
    def _generate_api_key(self) -> str:
        """Generate a secure API key."""
        return secrets.token_urlsafe(32)
    
    def to_dict(self) -> Dict:
        """Convert user to dictionary representation."""
        return {
            "username": self.username,
            "email": self.email,
            "role": self.role.value,
            "api_key": self.api_key,
            "created_at": self.created_at.isoformat(),
            "last_login": self.last_login.isoformat() if self.last_login else None,
            "is_active": self.is_active
        }

class UserManager:
    """Manages user creation, authentication, and role assignment."""
    
    def __init__(self):
        self.logger = logging.getLogger(__name__)
        self.users: Dict[str, User] = {}  # username -> User
        self.api_keys: Dict[str, str] = {}  # api_key -> username
        self.rbac = RBAC()
    
    def create_user(self, username: str, email: str, role: Role) -> User:
        """Create a new user with the specified role."""
        if username in self.users:
            raise ValueError(f"User {username} already exists")
        
        user = User(username, email, role)
        self.users[username] = user
        self.api_keys[user.api_key] = username
        return user
    
    def get_user(self, username: str) -> Optional[User]:
        """Get user by username."""
        return self.users.get(username)
    
    def get_user_by_api_key(self, api_key: str) -> Optional[User]:
        """Get user by API key."""
        username = self.api_keys.get(api_key)
        if username:
            return self.users.get(username)
        return None
    
    def update_user_role(self, username: str, new_role: Role) -> User:
        """Update a user's role."""
        user = self.get_user(username)
        if not user:
            raise ValueError(f"User {username} not found")
        
        user.role = new_role
        return user
    
    def deactivate_user(self, username: str) -> None:
        """Deactivate a user account."""
        user = self.get_user(username)
        if user:
            user.is_active = False
            # Optionally, invalidate API key
            if user.api_key in self.api_keys:
                del self.api_keys[user.api_key]
    
    def list_users(self) -> List[Dict]:
        """List all users with their roles."""
        return [user.to_dict() for user in self.users.values()]
    
    def validate_api_key(self, api_key: str) -> bool:
        """Validate an API key."""
        return api_key in self.api_keys
    
    def get_user_permissions(self, username: str) -> List[str]:
        """Get all permissions for a user."""
        user = self.get_user(username)
        if not user:
            return []
        return [p.value for p in self.rbac.get_role_permissions(user.role)] 