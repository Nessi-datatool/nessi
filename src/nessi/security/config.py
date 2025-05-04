from typing import Dict, List, Optional
from pydantic import BaseModel
import os
from datetime import timedelta

class SecurityConfig(BaseModel):
    """Security configuration for the application."""
    # JWT Configuration
    JWT_SECRET_KEY: str = os.getenv("JWT_SECRET_KEY", "your-secret-key")
    JWT_ALGORITHM: str = "HS256"
    JWT_ACCESS_TOKEN_EXPIRES: timedelta = timedelta(minutes=30)
    JWT_REFRESH_TOKEN_EXPIRES: timedelta = timedelta(days=7)

    # Rate Limiting
    RATE_LIMIT_ENABLED: bool = True
    RATE_LIMIT_DEFAULT: str = "100/minute"
    RATE_LIMIT_STORAGE_URL: str = "redis://localhost:6379/0"

    # IP Allowlisting
    IP_ALLOWLIST_ENABLED: bool = True
    IP_ALLOWLIST: List[str] = ["127.0.0.1", "::1"]

    # SSL/TLS Configuration
    SSL_ENABLED: bool = True
    SSL_CERT_PATH: Optional[str] = os.getenv("SSL_CERT_PATH")
    SSL_KEY_PATH: Optional[str] = os.getenv("SSL_KEY_PATH")

    # RBAC Configuration
    ROLES: Dict[str, List[str]] = {
        "admin": ["*"],
        "user": [
            "read:table",
            "read:metrics",
            "read:reports"
        ],
        "viewer": [
            "read:metrics",
            "read:reports"
        ]
    }

    # API Key Configuration
    API_KEY_HEADER: str = "X-API-Key"
    API_KEY_LENGTH: int = 32

    # Password Policy
    PASSWORD_MIN_LENGTH: int = 12
    PASSWORD_REQUIRE_UPPERCASE: bool = True
    PASSWORD_REQUIRE_LOWERCASE: bool = True
    PASSWORD_REQUIRE_NUMBERS: bool = True
    PASSWORD_REQUIRE_SPECIAL: bool = True

    # Session Configuration
    SESSION_TIMEOUT: int = 3600  # 1 hour
    SESSION_COOKIE_SECURE: bool = True
    SESSION_COOKIE_HTTPONLY: bool = True
    SESSION_COOKIE_SAMESITE: str = "Lax"

    # Audit Logging
    AUDIT_LOG_ENABLED: bool = True
    AUDIT_LOG_PATH: str = "/var/log/nessi/audit.log"

    # Security Headers
    SECURITY_HEADERS: Dict[str, str] = {
        "X-Content-Type-Options": "nosniff",
        "X-Frame-Options": "DENY",
        "X-XSS-Protection": "1; mode=block",
        "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
        "Content-Security-Policy": "default-src 'self'",
        "Referrer-Policy": "strict-origin-when-cross-origin"
    } 