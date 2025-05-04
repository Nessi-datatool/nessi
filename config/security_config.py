from typing import Dict, Any, Optional, List
from dataclasses import dataclass, field
from pathlib import Path
import json
import logging

logger = logging.getLogger(__name__)

@dataclass
class SSLSettings:
    """SSL/TLS configuration settings."""
    enabled: bool = True
    cert_path: Optional[str] = None
    key_path: Optional[str] = None
    ca_path: Optional[str] = None
    verify_client: bool = False
    min_version: str = "TLSv1.2"
    cipher_suite: List[str] = field(default_factory=lambda: [
        "ECDHE-ECDSA-AES256-GCM-SHA384",
        "ECDHE-RSA-AES256-GCM-SHA384",
        "ECDHE-ECDSA-CHACHA20-POLY1305",
        "ECDHE-RSA-CHACHA20-POLY1305",
        "ECDHE-ECDSA-AES128-GCM-SHA256",
        "ECDHE-RSA-AES128-GCM-SHA256"
    ])

@dataclass
class Role:
    """Role definition for RBAC."""
    name: str
    permissions: List[str]
    description: Optional[str] = None

@dataclass
class User:
    """User definition for RBAC."""
    username: str
    password_hash: str
    roles: List[str]
    is_active: bool = True
    last_login: Optional[str] = None

@dataclass
class SecurityConfig:
    """Security configuration for the application."""
    ssl: SSLSettings = field(default_factory=SSLSettings)
    roles: Dict[str, Role] = field(default_factory=dict)
    users: Dict[str, User] = field(default_factory=dict)
    session_timeout: int = 3600  # 1 hour
    max_login_attempts: int = 5
    password_policy: Dict[str, Any] = field(default_factory=lambda: {
        "min_length": 12,
        "require_uppercase": True,
        "require_lowercase": True,
        "require_numbers": True,
        "require_special_chars": True,
        "expiry_days": 90
    })
    audit_log_path: str = "logs/security_audit.log"
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'SecurityConfig':
        """Create a SecurityConfig instance from a dictionary."""
        try:
            ssl_data = data.get('ssl', {})
            ssl = SSLSettings(
                enabled=ssl_data.get('enabled', True),
                cert_path=ssl_data.get('cert_path'),
                key_path=ssl_data.get('key_path'),
                ca_path=ssl_data.get('ca_path'),
                verify_client=ssl_data.get('verify_client', False),
                min_version=ssl_data.get('min_version', "TLSv1.2"),
                cipher_suite=ssl_data.get('cipher_suite', [])
            )
            
            roles = {}
            for role_data in data.get('roles', []):
                role = Role(
                    name=role_data['name'],
                    permissions=role_data['permissions'],
                    description=role_data.get('description')
                )
                roles[role.name] = role
            
            users = {}
            for user_data in data.get('users', []):
                user = User(
                    username=user_data['username'],
                    password_hash=user_data['password_hash'],
                    roles=user_data['roles'],
                    is_active=user_data.get('is_active', True),
                    last_login=user_data.get('last_login')
                )
                users[user.username] = user
            
            return cls(
                ssl=ssl,
                roles=roles,
                users=users,
                session_timeout=data.get('session_timeout', 3600),
                max_login_attempts=data.get('max_login_attempts', 5),
                password_policy=data.get('password_policy', {}),
                audit_log_path=data.get('audit_log_path', "logs/security_audit.log")
            )
        
        except Exception as e:
            logger.error(f"Error creating SecurityConfig from dictionary: {str(e)}")
            raise
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert the SecurityConfig instance to a dictionary."""
        return {
            'ssl': {
                'enabled': self.ssl.enabled,
                'cert_path': self.ssl.cert_path,
                'key_path': self.ssl.key_path,
                'ca_path': self.ssl.ca_path,
                'verify_client': self.ssl.verify_client,
                'min_version': self.ssl.min_version,
                'cipher_suite': self.ssl.cipher_suite
            },
            'roles': [
                {
                    'name': role.name,
                    'permissions': role.permissions,
                    'description': role.description
                }
                for role in self.roles.values()
            ],
            'users': [
                {
                    'username': user.username,
                    'password_hash': user.password_hash,
                    'roles': user.roles,
                    'is_active': user.is_active,
                    'last_login': user.last_login
                }
                for user in self.users.values()
            ],
            'session_timeout': self.session_timeout,
            'max_login_attempts': self.max_login_attempts,
            'password_policy': self.password_policy,
            'audit_log_path': self.audit_log_path
        }
    
    def save(self, path: str):
        """Save the configuration to a file."""
        try:
            with open(path, 'w') as f:
                json.dump(self.to_dict(), f, indent=2)
            logger.info(f"Saved security configuration to {path}")
        except Exception as e:
            logger.error(f"Error saving security configuration: {str(e)}")
            raise
    
    @classmethod
    def load(cls, path: str) -> 'SecurityConfig':
        """Load the configuration from a file."""
        try:
            with open(path, 'r') as f:
                data = json.load(f)
            return cls.from_dict(data)
        except Exception as e:
            logger.error(f"Error loading security configuration: {str(e)}")
            raise 