import pytest
from datetime import datetime, timedelta
from src.license import LicenseValidator
from src.decorators import require_valid_license, LicenseError
import os
import shutil

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

def test_initial_trial_status(validator):
    """Test that a new installation starts with a trial period."""
    status = validator.get_status()
    assert status["status"] == "trial"
    assert "days remaining" in status["message"]
    assert status["is_licensed"] is True

def test_license_activation(validator):
    """Test license activation with a valid key."""
    assert validator.validate_license_key("test-license-key") is True
    status = validator.get_status()
    assert status["status"] == "licensed"
    assert status["is_licensed"] is True

def test_invalid_license_key(validator):
    """Test license activation with an invalid key."""
    assert validator.validate_license_key("") is False
    assert validator.validate_license_key(None) is False

def test_license_decorator(validator):
    """Test the license decorator."""
    @require_valid_license
    def test_function():
        return "success"
    
    # Should work during trial period
    assert test_function() == "success"
    
    # Activate a license
    validator.validate_license_key("test-license-key")
    assert test_function() == "success"
    
    # Simulate expired license
    validator.config["expiry_date"] = (datetime.now() - timedelta(days=1)).isoformat()
    validator._save_config()  # Make sure to save the config
    
    with pytest.raises(LicenseError) as exc_info:
        test_function()
    assert "License required" in str(exc_info.value)

def test_trial_expiration(validator):
    """Test trial period expiration."""
    # Set trial start to 15 days ago
    validator.config["trial_start"] = (datetime.now() - timedelta(days=15)).isoformat()
    validator._save_config()
    
    status = validator.get_status()
    assert status["status"] == "trial_expired"
    assert status["is_licensed"] is False
    assert "Trial period has expired" in status["message"]

def test_license_expiration(validator):
    """Test license expiration."""
    # Activate a license
    validator.validate_license_key("test-license-key")
    
    # Set expiry to yesterday
    validator.config["expiry_date"] = (datetime.now() - timedelta(days=1)).isoformat()
    validator._save_config()
    
    status = validator.get_status()
    assert status["status"] == "expired"
    assert status["is_licensed"] is False
    assert "License has expired" in status["message"] 