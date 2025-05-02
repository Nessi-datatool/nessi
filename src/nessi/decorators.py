"""Decorators for license validation."""

from functools import wraps
from .license import LicenseValidator, LicenseError

def require_valid_license(func):
    """Decorator to check if the license is valid before executing a function."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        validator = LicenseValidator()
        if not validator.check_license():
            raise LicenseError("Invalid or expired license")
        return func(*args, **kwargs)
    return wrapper 