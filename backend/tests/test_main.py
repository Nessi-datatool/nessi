import pytest
import sys
import os
from src.main import run_tests

def test_run_tests():
    """Test run_tests function"""
    # Mock pytest.main to return success
    original_pytest_main = pytest.main
    pytest.main = lambda args: 0
    
    try:
        result = run_tests()
        assert result == 0
    finally:
        # Restore original pytest.main
        pytest.main = original_pytest_main 