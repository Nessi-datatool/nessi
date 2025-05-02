import os
import json
import base64
import hashlib
import time
import logging
import platform
from pathlib import Path
from datetime import datetime, timedelta
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
            "is_licensed": False,
            "trial_start": datetime.now().isoformat(),
            "trial_days": 14,
            "expiry_date": None,
            "license_key": None,
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
            raise LicenseError(f"Failed to save license configuration: {str(e)}")
    
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
    
    def validate_license_key(self, license_key: str) -> bool:
        """Validate a license key.
        
        Args:
            license_key (str): License key to validate.
            
        Returns:
            bool: True if license key is valid.
            
        Raises:
            LicenseError: If license key is invalid.
        """
        if not license_key:
            raise LicenseError("Invalid license key: Invalid license key format")
        
        try:
            # Verify machine ID
            current_machine_id = self._get_machine_id()
            if self.config.get("machine_id") and self.config["machine_id"] != current_machine_id:
                raise LicenseError("License not valid for this machine")
            
            # Update config
            self.config.update({
                "is_licensed": True,
                "expiry_date": (datetime.now() + timedelta(days=365)).isoformat(),
                "license_key": license_key,
                "machine_id": current_machine_id
            })
            self._save_config()
            
            return True
        except Exception as e:
            logger.error(f"License validation failed: {str(e)}")
            raise LicenseError(f"License validation failed: {str(e)}")
    
    def get_status(self) -> dict:
        """Get current license status.
        
        Returns:
            dict: License status information.
        """
        try:
            # Check if licensed
            if self.config["is_licensed"] and self.config["expiry_date"]:
                expiry_date = datetime.fromisoformat(self.config["expiry_date"])
                is_expired = datetime.now() > expiry_date
                
                if is_expired:
                    return {
                        "status": "expired",
                        "is_licensed": False,
                        "message": "License has expired",
                        "expiry_date": self.config["expiry_date"],
                        "days_remaining": 0
                    }
                else:
                    return {
                        "status": "licensed",
                        "is_licensed": True,
                        "message": "License is valid",
                        "expiry_date": self.config["expiry_date"],
                        "days_remaining": (expiry_date - datetime.now()).days
                    }
            
            # Check trial period
            trial_start = datetime.fromisoformat(self.config["trial_start"])
            trial_days = self.config["trial_days"]
            trial_end = trial_start + timedelta(days=trial_days)
            days_remaining = (trial_end - datetime.now()).days
            
            if days_remaining > 0:
                return {
                    "status": "trial",
                    "is_licensed": True,
                    "message": f"Trial period active ({days_remaining} days remaining)",
                    "expiry_date": trial_end.isoformat(),
                    "days_remaining": days_remaining
                }
            else:
                return {
                    "status": "trial_expired",
                    "is_licensed": False,
                    "message": "Trial period has expired",
                    "expiry_date": trial_end.isoformat(),
                    "days_remaining": 0
                }
        except Exception as e:
            logger.error(f"Failed to get license status: {str(e)}")
            return {
                "status": "error",
                "is_licensed": False,
                "message": f"Error checking license status: {str(e)}",
                "expiry_date": None,
                "days_remaining": 0
            }
    
    def check_license(self) -> bool:
        """Check if license is valid.
        
        Returns:
            bool: True if license is valid.
            
        Raises:
            LicenseError: If license is invalid.
        """
        status = self.get_status()
        if not status["is_licensed"]:
            raise LicenseError("License is invalid or expired")
        return True

def requires_license(func):
    """Decorator to check for valid license or trial."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        validator = LicenseValidator()
        validator.check_license()
        return func(*args, **kwargs)
    return wrapper 