"""Security manager module."""

from typing import Dict, Any, Optional, List
from fastapi import HTTPException, Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
import jwt
import bcrypt
import logging
from datetime import datetime, timedelta
from .config import SecurityConfig
from .middleware import SecurityMiddleware
from .password import PasswordManager
from .audit_logger import AuditLogger
from .ip_manager import IPManager
from .rate_limiter import RateLimiter
from .ssl_manager import SSLManager

logger = logging.getLogger(__name__)


class SecurityManager:
    """Security manager class."""
    
    def __init__(self, config: SecurityConfig):
        """Initialize security manager."""
        self.config = config
        self.security = HTTPBearer()
        self.password_manager = PasswordManager(config)
        self.audit_logger = AuditLogger(config)
        self.ip_manager = IPManager(config)
        self.rate_limiter = RateLimiter(config)
        self.ssl_manager = SSLManager(config)
        self.middleware = SecurityMiddleware(config)
    
    def create_token(self, user_id: str, roles: List[str]) -> str:
        """Create a JWT token for a user."""
        try:
            payload = {
                "user_id": user_id,
                "roles": roles,
                "exp": datetime.utcnow() + timedelta(
                    hours=self.config.jwt_expiry_hours
                )
            }
            return jwt.encode(
                payload,
                self.config.jwt_secret,
                algorithm=self.config.jwt_algorithm
            )
        except Exception as e:
            logger.error(f"Error creating token: {str(e)}")
            raise HTTPException(
                status_code=500,
                detail="Error creating token"
            )
    
    def verify_token(self, token: str) -> Dict[str, Any]:
        """Verify a JWT token."""
        try:
            return jwt.decode(
                token,
                self.config.jwt_secret,
                algorithms=[self.config.jwt_algorithm]
            )
        except jwt.ExpiredSignatureError:
            raise HTTPException(
                status_code=401,
                detail="Token has expired"
            )
        except jwt.InvalidTokenError:
            raise HTTPException(
                status_code=401,
                detail="Invalid token"
            )
    
    async def authenticate(self, request: Request) -> Dict[str, Any]:
        """Authenticate a request."""
        try:
            # Get token from request
            credentials: HTTPAuthorizationCredentials = await self.security(request)
            token = credentials.credentials
            
            # Verify token
            payload = self.verify_token(token)
            
            # Check IP whitelist/blacklist
            if not self.ip_manager.is_allowed(request.client.host):
                raise HTTPException(
                    status_code=403,
                    detail="IP not allowed"
                )
            
            # Check rate limit
            if self.rate_limiter.is_rate_limited(request.client.host):
                raise HTTPException(
                    status_code=429,
                    detail="Too many requests"
                )
            
            # Log audit event
            self.audit_logger.log_auth_event(
                user_id=payload["user_id"],
                event_type="authentication",
                success=True,
                client_ip=request.client.host
            )
            
            return payload
        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"Error authenticating request: {str(e)}")
            raise HTTPException(
                status_code=500,
                detail="Error authenticating request"
            )
    
    def check_role(self, required_roles: List[str], user_roles: List[str]) -> bool:
        """Check if a user has the required roles."""
        return any(role in user_roles for role in required_roles)
    
    def hash_password(self, password: str) -> str:
        """Hash a password."""
        return self.password_manager.hash_password(password)
    
    def verify_password(self, password: str, hashed_password: str) -> bool:
        """Verify a password."""
        return self.password_manager.verify_password(password, hashed_password)
    
    def validate_password(self, password: str) -> bool:
        """Validate a password."""
        return self.password_manager.validate_password(password)
    
    def get_ssl_context(self):
        """Get SSL context."""
        return self.ssl_manager.get_ssl_context()
    
    def update_rate_limit(
        self,
        requests: int,
        window: int,
        block_duration: int = 300
    ) -> None:
        """Update rate limit settings."""
        self.rate_limiter.update_settings(requests, window, block_duration)
    
    def add_to_allowlist(self, ip_or_network: str) -> None:
        """Add IP or network to allowlist."""
        self.ip_manager.add_to_allowlist(ip_or_network)
    
    def remove_from_allowlist(self, ip_or_network: str) -> None:
        """Remove IP or network from allowlist."""
        self.ip_manager.remove_from_allowlist(ip_or_network)
    
    def add_to_blocklist(self, ip_or_network: str) -> None:
        """Add IP or network to blocklist."""
        self.ip_manager.add_to_blocklist(ip_or_network)
    
    def remove_from_blocklist(self, ip_or_network: str) -> None:
        """Remove IP or network from blocklist."""
        self.ip_manager.remove_from_blocklist(ip_or_network) 