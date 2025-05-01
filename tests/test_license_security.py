import pytest
import json
import os
import base64
from datetime import datetime, timedelta
from pathlib import Path
from src.license import LicenseValidator
from src.decorators import require_valid_license, LicenseError

@pytest.fixture
def temp_config_dir(tmp_path):
    """Create a temporary config directory for testing."""
    config_dir = tmp_path / ".nessi"
    config_dir.mkdir()
    return config_dir

@pytest.fixture
def validator(temp_config_dir, monkeypatch):
    """Create a LicenseValidator with a temporary config directory."""
    monkeypatch.setattr("src.license.Path.home", lambda: temp_config_dir.parent)
    return LicenseValidator()

def test_config_encryption(validator, temp_config_dir):
    """Test that config file is properly encrypted."""
    # Initial config should create encrypted file
    config_file = temp_config_dir / "config.json"
    assert config_file.exists()
    
    # Read raw file content
    with open(config_file, 'rb') as f:
        content = f.read()
    
    # Content should be encrypted (not plain JSON)
    with pytest.raises(json.JSONDecodeError):
        json.loads(content)

def test_machine_binding(validator):
    """Test that license is bound to machine."""
    # Activate license
    assert validator.validate_license_key("test-license")
    
    # Simulate different machine
    original_machine_id = validator._get_machine_id()
    validator.config["machine_id"] = "different-machine-id"
    validator._save_config()
    
    # License should be invalid on different machine
    status = validator.get_status()
    assert status["is_licensed"] is False
    assert status["status"] == "invalid"
    assert "not valid for this machine" in status["message"]

def test_config_tampering(validator, temp_config_dir):
    """Test that tampered config file is rejected."""
    config_file = temp_config_dir / "config.json"
    
    # Corrupt the config file
    with open(config_file, 'wb') as f:
        f.write(b'corrupted data')
    
    # Creating new validator should reset config
    new_validator = LicenseValidator()
    status = new_validator.get_status()
    assert status["status"] == "trial"

def test_key_file_security(validator, temp_config_dir):
    """Test encryption key file handling."""
    key_file = temp_config_dir / ".key"
    assert key_file.exists()
    
    # Key file should contain salt
    with open(key_file, 'rb') as f:
        salt = f.read()
    assert len(salt) == 16  # Standard salt length

def test_license_expiration(validator):
    """Test that expired license is properly handled."""
    # Activate license
    assert validator.validate_license_key("test-license")
    
    # Set expiration to past date
    validator.config["expiry_date"] = (datetime.now() - timedelta(days=1)).isoformat()
    validator._save_config()
    
    # License should be expired
    status = validator.get_status()
    assert status["is_licensed"] is False
    assert status["status"] == "expired"

def test_trial_manipulation(validator):
    """Test against trial period manipulation."""
    # Try to extend trial by changing start date
    original_start = validator.config["trial_start"]
    validator.config["trial_start"] = (datetime.now() - timedelta(days=1)).isoformat()
    validator._save_config()
    
    # Config should maintain encryption
    with open(validator.config_file, 'rb') as f:
        content = f.read()
    with pytest.raises(json.JSONDecodeError):
        json.loads(content)

def test_protected_function_security():
    """Test that protected functions enforce license check."""
    @require_valid_license
    def protected_function():
        return "success"
    
    # Function should raise LicenseError when license is invalid
    with pytest.raises(LicenseError) as exc_info:
        # Create validator with expired trial
        temp_dir = Path("/tmp/test_nessi")
        temp_dir.mkdir(exist_ok=True)
        os.environ["HOME"] = str(temp_dir)
        validator = LicenseValidator()
        validator.config["trial_start"] = (datetime.now() - timedelta(days=15)).isoformat()
        validator._save_config()
        
        # Try to call protected function
        protected_function()
    
    assert "License required" in str(exc_info.value)

def test_multiple_instances(validator, temp_config_dir):
    """Test that multiple validator instances use same config."""
    # First instance activates license
    assert validator.validate_license_key("test-license")
    
    # Second instance should see the same license
    validator2 = LicenseValidator()
    status = validator2.get_status()
    assert status["status"] == "licensed"
    assert status["is_licensed"] is True

def test_invalid_config_recovery(validator, temp_config_dir):
    """Test recovery from invalid config."""
    # Corrupt the config file with invalid encrypted data
    with open(validator.config_file, 'wb') as f:
        f.write(base64.b64encode(b'invalid encrypted data'))
    
    # New validator should recover and create fresh config
    new_validator = LicenseValidator()
    status = new_validator.get_status()
    assert status["status"] == "trial" 