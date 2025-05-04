from fastapi import FastAPI, Request, HTTPException, Depends
from fastapi.security import OAuth2PasswordBearer, APIKeyHeader, HTTPBearer, HTTPAuthorizationCredentials
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.trustedhost import TrustedHostMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from starlette.middleware.sessions import SessionMiddleware
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.types import ASGIApp
from starlette.responses import Response
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address
from slowapi.errors import RateLimitExceeded
import jwt
import redis
import logging
from typing import Optional, List, Callable, Dict, Any
import time
from .config import SecurityConfig
from .rate_limiter import RateLimiter, RateLimit
from fastapi.responses import JSONResponse
from .ip_manager import IPManager
from .audit_logger import AuditLogger
import ssl
import bcrypt
from datetime import datetime, timedelta
from pathlib import Path
import json
from functools import wraps

logger = logging.getLogger(__name__)

class SecurityMiddleware:
    """Security middleware with RBAC and rate limiting."""
    
    def __init__(self, config_path: str = "config/security.json"):
        self.config = self._load_config(config_path)
        self.redis = self._setup_redis()
        self.security = HTTPBearer()
    
    def _load_config(self, config_path: str) -> Dict[str, Any]:
        """Load security configuration."""
        try:
            with open(config_path, 'r') as f:
                return json.load(f)
        except Exception as e:
            logger.error(f"Error loading security config: {str(e)}")
            return {}
    
    def _setup_redis(self) -> redis.Redis:
        """Set up Redis connection for rate limiting."""
        try:
            return redis.Redis(
                host=self.config.get("redis_host", "localhost"),
                port=self.config.get("redis_port", 6379),
                db=self.config.get("redis_db", 0),
                decode_responses=True
            )
        except Exception as e:
            logger.error(f"Error setting up Redis: {str(e)}")
            raise
    
    def create_token(self, user_id: str, roles: List[str]) -> str:
        """Create a JWT token for a user."""
        try:
            payload = {
                "user_id": user_id,
                "roles": roles,
                "exp": datetime.utcnow() + timedelta(
                    hours=self.config.get("token_expiry_hours", 24)
                )
            }
            return jwt.encode(
                payload,
                self.config["jwt_secret"],
                algorithm="HS256"
            )
        except Exception as e:
            logger.error(f"Error creating token: {str(e)}")
            raise
    
    def verify_token(self, token: str) -> Dict[str, Any]:
        """Verify a JWT token."""
        try:
            return jwt.decode(
                token,
                self.config["jwt_secret"],
                algorithms=["HS256"]
            )
        except jwt.ExpiredSignatureError:
            raise HTTPException(status_code=401, detail="Token has expired")
        except jwt.InvalidTokenError:
            raise HTTPException(status_code=401, detail="Invalid token")
    
    def rate_limit(self, key: str, limit: int, window: int) -> bool:
        """Check if a request should be rate limited."""
        try:
            current = self.redis.get(key)
            if current is None:
                self.redis.setex(key, window, 1)
                return False
            
            current = int(current)
            if current >= limit:
                return True
            
            self.redis.incr(key)
            return False
        except Exception as e:
            logger.error(f"Error checking rate limit: {str(e)}")
            return False
    
    def check_role(self, required_roles: List[str]) -> Callable:
        """Decorator to check if a user has the required roles."""
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(*args, **kwargs):
                try:
                    # Get token from request
                    credentials: HTTPAuthorizationCredentials = await self.security(*args, **kwargs)
                    token = credentials.credentials
                    
                    # Verify token and get user roles
                    payload = self.verify_token(token)
                    user_roles = payload.get("roles", [])
                    
                    # Check if user has any of the required roles
                    if not any(role in user_roles for role in required_roles):
                        raise HTTPException(
                            status_code=403,
                            detail="Insufficient permissions"
                        )
                    
                    return await func(*args, **kwargs)
                except Exception as e:
                    logger.error(f"Error checking roles: {str(e)}")
                    raise
            return wrapper
        return decorator
    
    def rate_limit_middleware(self, limit: int = 100, window: int = 60) -> Callable:
        """Middleware to enforce rate limiting."""
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(request: Request, *args, **kwargs):
                try:
                    # Get client IP
                    client_ip = request.client.host
                    
                    # Create rate limit key
                    key = f"rate_limit:{client_ip}:{func.__name__}"
                    
                    # Check rate limit
                    if self.rate_limit(key, limit, window):
                        raise HTTPException(
                            status_code=429,
                            detail="Too many requests"
                        )
                    
                    return await func(request, *args, **kwargs)
                except Exception as e:
                    logger.error(f"Error in rate limit middleware: {str(e)}")
                    raise
            return wrapper
        return decorator
    
    def ssl_middleware(self) -> Callable:
        """Middleware to enforce SSL/TLS."""
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(request: Request, *args, **kwargs):
                try:
                    if not request.url.scheme == "https":
                        raise HTTPException(
                            status_code=403,
                            detail="SSL/TLS required"
                        )
                    return await func(request, *args, **kwargs)
                except Exception as e:
                    logger.error(f"Error in SSL middleware: {str(e)}")
                    raise
            return wrapper
        return decorator
    
    def ip_whitelist_middleware(self) -> Callable:
        """Middleware to enforce IP whitelist."""
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(request: Request, *args, **kwargs):
                try:
                    client_ip = request.client.host
                    whitelist = self.config.get("ip_whitelist", [])
                    
                    if whitelist and client_ip not in whitelist:
                        raise HTTPException(
                            status_code=403,
                            detail="IP not allowed"
                        )
                    
                    return await func(request, *args, **kwargs)
                except Exception as e:
                    logger.error(f"Error in IP whitelist middleware: {str(e)}")
                    raise
            return wrapper
        return decorator
    
    def audit_log_middleware(self) -> Callable:
        """Middleware to log all requests for audit purposes."""
        def decorator(func: Callable) -> Callable:
            @wraps(func)
            async def wrapper(request: Request, *args, **kwargs):
                try:
                    start_time = time.time()
                    
                    # Get request details
                    client_ip = request.client.host
                    method = request.method
                    path = request.url.path
                    user_agent = request.headers.get("user-agent")
                    
                    # Execute the request
                    response = await func(request, *args, **kwargs)
                    
                    # Calculate duration
                    duration = time.time() - start_time
                    
                    # Log audit entry
                    audit_entry = {
                        "timestamp": datetime.utcnow().isoformat(),
                        "client_ip": client_ip,
                        "method": method,
                        "path": path,
                        "user_agent": user_agent,
                        "duration": duration,
                        "status_code": response.status_code
                    }
                    
                    logger.info(f"Audit log: {json.dumps(audit_entry)}")
                    
                    return response
                except Exception as e:
                    logger.error(f"Error in audit log middleware: {str(e)}")
                    raise
            return wrapper
        return decorator

class SecurityHeadersMiddleware(BaseHTTPMiddleware):
    """Middleware for adding security headers."""
    
    def __init__(self, app, config: SecurityConfig):
        super().__init__(app)
        self.config = config
    
    async def dispatch(self, request: Request, call_next):
        response = await call_next(request)
        for header, value in self.config.SECURITY_HEADERS.items():
            response.headers[header] = value
        return response

# OAuth2 scheme for token authentication
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

# API key scheme for API key authentication
api_key_scheme = APIKeyHeader(name="X-API-Key")

class SecurityMiddleware(BaseHTTPMiddleware):
    """Middleware for rate limiting and IP allowlisting."""
    
    def __init__(
        self,
        app: FastAPI,
        rate_limit: RateLimit = RateLimit(requests=100, window=60),
        exclude_paths: list[str] = None
    ):
        super().__init__(app)
        self.logger = logging.getLogger(__name__)
        self.rate_limiter = RateLimiter()
        self.rate_limit = rate_limit
        self.exclude_paths = exclude_paths or ["/health", "/metrics"]
    
    async def dispatch(
        self,
        request: Request,
        call_next: Callable
    ) -> Response:
        """Process each request through security checks."""
        # Skip security checks for excluded paths
        if request.url.path in self.exclude_paths:
            return await call_next(request)
        
        # Get client IP
        client_ip = request.client.host
        
        # Check rate limit
        if not self.rate_limiter.check_rate_limit(client_ip, self.rate_limit):
            self.logger.warning(f"Rate limit exceeded for IP: {client_ip}")
            return JSONResponse(
                status_code=429,
                content={
                    "detail": "Too many requests. Please try again later."
                }
            )
        
        # Process request
        try:
            response = await call_next(request)
            return response
        except Exception as e:
            self.logger.error(f"Error processing request: {str(e)}")
            return JSONResponse(
                status_code=500,
                content={
                    "detail": "Internal server error"
                }
            )
    
    def add_to_allowlist(self, ip_or_network: str) -> None:
        """Add an IP or network to the allowlist."""
        self.rate_limiter.add_to_allowlist(ip_or_network)
    
    def remove_from_allowlist(self, ip_or_network: str) -> None:
        """Remove an IP or network from the allowlist."""
        self.rate_limiter.remove_from_allowlist(ip_or_network)
    
    def update_rate_limit(self, requests: int, window: int, block_duration: int = 300) -> None:
        """Update rate limit configuration."""
        self.rate_limit = RateLimit(
            requests=requests,
            window=window,
            block_duration=block_duration
        )

class SecurityMiddleware(BaseHTTPMiddleware):
    """Middleware for handling security features including IP allowlisting and audit logging."""
    
    def __init__(
        self,
        app: ASGIApp,
        ip_manager: IPManager,
        audit_logger: AuditLogger,
        exclude_paths: Optional[list] = None
    ):
        super().__init__(app)
        self.ip_manager = ip_manager
        self.audit_logger = audit_logger
        self.exclude_paths = exclude_paths or []
    
    async def dispatch(self, request: Request, call_next: Callable) -> Response:
        # Skip security checks for excluded paths
        if request.url.path in self.exclude_paths:
            return await call_next(request)
        
        # Get client IP
        client_ip = request.client.host if request.client else None
        
        # Check IP allowlist
        if client_ip and not self.ip_manager.is_allowed(client_ip):
            self.audit_logger.log_security_event(
                user="anonymous",
                event_type="ip_blocked",
                details={
                    "ip": client_ip,
                    "path": request.url.path,
                    "method": request.method
                },
                ip_address=client_ip
            )
            return Response(
                content="Access denied",
                status_code=403
            )
        
        # Log request
        self.audit_logger.log_access(
            user=request.headers.get("X-User", "anonymous"),
            resource=request.url.path,
            action=request.method,
            success=True,
            ip_address=client_ip
        )
        
        # Process request
        response = await call_next(request)
        
        # Log response status
        if response.status_code >= 400:
            self.audit_logger.log_security_event(
                user=request.headers.get("X-User", "anonymous"),
                event_type="error",
                details={
                    "status_code": response.status_code,
                    "path": request.url.path,
                    "method": request.method
                },
                ip_address=client_ip
            )
        
        return response 