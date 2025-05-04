"""Security configuration module."""

from typing import Dict, Any, Optional, List
from pydantic import BaseModel, Field


class SecurityConfig(BaseModel):
    """Security configuration model."""
    
    # JWT settings
    jwt_secret: str = Field(default="your-secret-key", description="JWT secret key")
    jwt_algorithm: str = Field(default="HS256", description="JWT algorithm")
    jwt_expiry_hours: int = Field(default=24, description="JWT token expiry in hours")
    
    # Rate limiting
    rate_limit_requests: int = Field(default=100, description="Number of requests allowed per window")
    rate_limit_window: int = Field(default=60, description="Rate limit window in seconds")
    rate_limit_block_duration: int = Field(default=300, description="Block duration in seconds")
    
    # IP whitelist/blacklist
    ip_whitelist: Optional[List[str]] = Field(default=None, description="List of allowed IPs")
    ip_blacklist: Optional[List[str]] = Field(default=None, description="List of blocked IPs")
    
    # Password settings
    min_password_length: int = Field(default=8, description="Minimum password length")
    require_uppercase: bool = Field(default=True, description="Require uppercase letters")
    require_lowercase: bool = Field(default=True, description="Require lowercase letters")
    require_numbers: bool = Field(default=True, description="Require numbers")
    require_special_chars: bool = Field(default=True, description="Require special characters")
    
    # Session settings
    session_expiry: int = Field(default=3600, description="Session expiry in seconds")
    session_cookie_name: str = Field(default="session", description="Session cookie name")
    session_cookie_secure: bool = Field(default=True, description="Require HTTPS for cookies")
    session_cookie_httponly: bool = Field(default=True, description="Prevent JavaScript access")
    
    # SSL/TLS settings
    ssl_enabled: bool = Field(default=True, description="Enable SSL/TLS")
    ssl_cert_path: Optional[str] = Field(default=None, description="SSL certificate path")
    ssl_key_path: Optional[str] = Field(default=None, description="SSL private key path")
    
    # CORS settings
    cors_origins: List[str] = Field(default=["*"], description="Allowed CORS origins")
    cors_methods: List[str] = Field(
        default=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
        description="Allowed CORS methods"
    )
    cors_headers: List[str] = Field(default=["*"], description="Allowed CORS headers")
    
    # Audit logging
    audit_log_enabled: bool = Field(default=True, description="Enable audit logging")
    audit_log_path: str = Field(default="logs/audit.log", description="Audit log file path")
    audit_log_format: str = Field(default="json", description="Audit log format")
    
    # Redis settings
    redis_host: str = Field(default="localhost", description="Redis host")
    redis_port: int = Field(default=6379, description="Redis port")
    redis_db: int = Field(default=0, description="Redis database number")
    redis_password: Optional[str] = Field(default=None, description="Redis password")
    
    class Config:
        """Pydantic model configuration."""
        
        env_prefix = "NESSI_SECURITY_"
        case_sensitive = False 