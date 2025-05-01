from functools import wraps
from .license import LicenseValidator

class LicenseError(Exception):
    """Exception raised when license validation fails."""
    pass

def require_valid_license(func):
    """Decorator to ensure a valid license is present before executing a function."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        validator = LicenseValidator()
        status = validator.get_status()
        if not status["is_licensed"]:
            raise LicenseError(f"License required: {status['message']}")
        return func(*args, **kwargs)
    return wrapper 