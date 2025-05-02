import pytest
from backend.src.scanner.decorators import require_valid_license
from backend.src.scanner.license import LicenseError

@require_valid_license
def dummy_protected_function():
    return "Success"

def test_valid_license(mocker):
    """Test function execution with valid license"""
    # Mock the license validator to return True
    mocker.patch('backend.src.scanner.license.LicenseValidator.check_license', return_value=True)
    
    # Function should execute successfully
    result = dummy_protected_function()
    assert result == "Success"

def test_invalid_license(mocker):
    """Test function execution with invalid license"""
    # Mock the license validator to return False
    mocker.patch('backend.src.scanner.license.LicenseValidator.check_license', return_value=False)
    
    # Function should raise LicenseError
    with pytest.raises(LicenseError) as exc_info:
        dummy_protected_function()
    assert "Invalid or expired license" in str(exc_info.value)

def test_license_check_called(mocker):
    """Test that license check is actually called"""
    # Create a spy on the check_license method
    mock_check = mocker.patch('backend.src.scanner.license.LicenseValidator.check_license', return_value=True)
    
    # Call the protected function
    dummy_protected_function()
    
    # Verify that check_license was called
    mock_check.assert_called_once() 