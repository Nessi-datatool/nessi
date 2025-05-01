import os
import json
import base64
from datetime import datetime, timedelta
from pathlib import Path
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC

class LicenseValidator:
    def __init__(self):
        self.config_dir = Path.home() / ".nessi"
        self.config_file = self.config_dir / "config.json"
        self._ensure_config_dir()
        self._init_encryption()
        self._load_config()

    def _ensure_config_dir(self):
        """Ensure the config directory exists."""
        self.config_dir.mkdir(exist_ok=True)
        
    def _init_encryption(self):
        """Initialize encryption key."""
        key_file = self.config_dir / ".key"
        if not key_file.exists():
            # Generate a unique salt based on machine-specific data
            salt = os.urandom(16)
            with open(key_file, 'wb') as f:
                f.write(salt)
        else:
            with open(key_file, 'rb') as f:
                salt = f.read()

        # Generate key from machine-specific data
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA256(),
            length=32,
            salt=salt,
            iterations=100000,
        )
        # Use hardware info and user info to make key unique per machine
        machine_data = f"{os.getlogin()}:{os.uname().nodename}:{os.uname().machine}"
        key = base64.urlsafe_b64encode(kdf.derive(machine_data.encode()))
        self.fernet = Fernet(key)

    def _load_config(self):
        """Load and decrypt the configuration file."""
        if self.config_file.exists():
            with open(self.config_file, 'rb') as f:
                encrypted_data = f.read()
                try:
                    decrypted_data = self.fernet.decrypt(encrypted_data)
                    self.config = json.loads(decrypted_data)
                except Exception:
                    # If decryption fails, start fresh
                    self._init_new_config()
        else:
            self._init_new_config()

    def _init_new_config(self):
        """Initialize new configuration with trial period."""
        self.config = {
            "trial_start": datetime.now().isoformat(),
            "license_key": None,
            "expiry_date": None,
            "machine_id": self._get_machine_id()
        }
        self._save_config()

    def _get_machine_id(self):
        """Generate a unique machine identifier."""
        machine_data = f"{os.getlogin()}:{os.uname().nodename}:{os.uname().machine}"
        return base64.b64encode(machine_data.encode()).decode()

    def _save_config(self):
        """Encrypt and save the configuration file."""
        encrypted_data = self.fernet.encrypt(json.dumps(self.config).encode())
        with open(self.config_file, 'wb') as f:
            f.write(encrypted_data)

    def get_status(self):
        """Get the current license status."""
        if self._get_machine_id() != self.config.get("machine_id"):
            return {
                "is_licensed": False,
                "status": "invalid",
                "message": "License not valid for this machine"
            }

        trial_start = datetime.fromisoformat(self.config["trial_start"])
        trial_end = trial_start + timedelta(days=14)
        
        if self.config["license_key"]:
            # Check if license is valid
            if self.config["expiry_date"]:
                expiry = datetime.fromisoformat(self.config["expiry_date"])
                if datetime.now() > expiry:
                    return {
                        "is_licensed": False,
                        "status": "expired",
                        "message": "License has expired"
                    }
            return {
                "is_licensed": True,
                "status": "licensed",
                "message": "Valid license"
            }
        
        # Check trial period
        if datetime.now() > trial_end:
            return {
                "is_licensed": False,
                "status": "trial_expired",
                "message": "Trial period has expired"
            }
        
        remaining_days = (trial_end - datetime.now()).days
        return {
            "is_licensed": True,
            "status": "trial",
            "message": f"Trial period: {remaining_days} days remaining"
        }

    def validate_license_key(self, license_key):
        """Validate a license key and update the configuration if valid."""
        if not license_key or not isinstance(license_key, str):
            return False
            
        # Update the configuration with the new license key
        self.config["license_key"] = license_key
        self.config["expiry_date"] = (datetime.now() + timedelta(days=365)).isoformat()
        self.config["machine_id"] = self._get_machine_id()
        self._save_config()
        return True

    def is_valid(self):
        """Check if the current license is valid."""
        status = self.get_status()
        return status["is_licensed"] 