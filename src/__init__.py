"""Nessi backend package."""

from .license import LicenseValidator
from .decorators import require_valid_license, LicenseError

__version__ = '1.0.0'
__all__ = ['scanner', 'LicenseValidator', 'require_valid_license', 'LicenseError'] 