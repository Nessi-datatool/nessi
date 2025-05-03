"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 nessi.dev. All rights reserved.

nessi.dev is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of nessi.dev will require a license for all users.

1. PERSONAL USE LICENSE
   a) The Software is provided free of charge for personal use
   b) Personal use includes:
      - Individual data analysis and processing
      - Educational purposes
      - Non-commercial research
   c) Enterprise or consulting use requires a paid license

2. RESTRICTIONS
   You shall not:
   a) Copy, modify, adapt, translate, reverse engineer, decompile, or disassemble the Software
   b) Create derivative works based on the Software
   c) Rent, lease, loan, sell, sublicense, distribute, transmit, or otherwise transfer the Software
   d) Remove or alter any proprietary notices or labels on the Software
   e) Use the Software for enterprise or consulting purposes without a valid paid license

3. OWNERSHIP
   The Software is licensed, not sold. nessi.dev retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

"""License validation module."""

from .validator import LicenseValidator, LicenseError

__all__ = ['LicenseValidator', 'LicenseError']

import os
import json
import base64
import hashlib
import time
import logging
import platform
from pathlib import Path
from datetime import datetime
from typing import Dict, Any, Optional
from functools import wraps
from cryptography.fernet import Fernet

logger = logging.getLogger(__name__)

class LicenseError(Exception):
    """Exception raised for license-related errors."""
    pass

class LicenseValidator:
    """Class for validating and managing software licenses."""
    
    def __init__(self, config_dir: str = None):
        """Initialize the license validator.
        
        Args:
            config_dir (str): Directory to store license configuration.
        """
        self.config_dir = config_dir or os.path.expanduser("~/.nessi")
        self.config_file = os.path.join(self.config_dir, "config.json")
        self.key_file = os.path.join(self.config_dir, ".key")
        
        # Create config directory if it doesn't exist
        os.makedirs(self.config_dir, exist_ok=True)
        
        # Initialize or load config
        self.config = self._load_config()
    
    def _load_config(self) -> dict:
        """Load or initialize license configuration.
        
        Returns:
            dict: License configuration.
        """
        if os.path.exists(self.config_file):
            try:
                with open(self.config_file, 'r') as f:
                    return json.load(f)
            except Exception as e:
                logger.error(f"Failed to load config: {str(e)}")
                return self._init_config()
        return self._init_config()
    
    def _init_config(self) -> dict:
        """Initialize default license configuration.
        
        Returns:
            dict: Default configuration.
        """
        config = {
            "is_personal_use": True,
            "machine_id": self._get_machine_id()
        }
        self._save_config(config)
        return config
    
    def _save_config(self, config: dict = None) -> None:
        """Save license configuration.
        
        Args:
            config (dict): Configuration to save. If None, saves current config.
        """
        if config is None:
            config = self.config
        
        try:
            with open(self.config_file, 'w') as f:
                json.dump(config, f, indent=2)
        except Exception as e:
            logger.error(f"Failed to save config: {str(e)}")
            raise LicenseError(f"Failed to save configuration: {str(e)}")
    
    def _get_machine_id(self) -> str:
        """Get a unique machine identifier.
        
        Returns:
            str: Machine identifier.
        """
        try:
            # Get system information
            system = platform.system()
            machine = platform.machine()
            processor = platform.processor()
            node = platform.node()
            
            # Create a unique identifier
            identifier = f"{system}-{machine}-{processor}-{node}"
            return hashlib.sha256(identifier.encode()).hexdigest()
        except Exception as e:
            logger.error(f"Failed to get machine ID: {str(e)}")
            return None
    
    def get_status(self) -> dict:
        """Get current license status.
        
        Returns:
            dict: License status information.
        """
        try:
            return {
                "status": "active",
                "type": "personal" if self.config["is_personal_use"] else "enterprise",
                "message": "Free for personal use" if self.config["is_personal_use"] else "Enterprise license required"
            }
        except Exception as e:
            logger.error(f"Failed to get status: {str(e)}")
            raise LicenseError(f"Failed to get status: {str(e)}")
    
    def check_license(self) -> bool:
        """Check if the current usage is allowed.
        
        Returns:
            bool: True if usage is allowed.
        """
        try:
            return self.config["is_personal_use"]
        except Exception as e:
            logger.error(f"Failed to check license: {str(e)}")
            return False

def requires_personal_use(func):
    """Decorator to check if the usage is personal use."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        validator = LicenseValidator()
        if not validator.check_license():
            raise LicenseError("This feature requires personal use only. Enterprise use requires a paid license.")
        return func(*args, **kwargs)
    return wrapper