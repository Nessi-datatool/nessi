import pytest
import os
import sys
from backend.src.main import run_tests

def test_run_tests(mocker):
    """Test that run_tests function executes pytest correctly"""
    # Mock pytest.main to avoid actually running tests
    mock_pytest_main = mocker.patch('pytest.main', return_value=0)
    
    # Call run_tests
    result = run_tests()
    
    # Verify pytest.main was called with correct arguments
    test_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), "tests")
    mock_pytest_main.assert_called_once_with(["-v", test_dir])
    
    # Verify the return value
    assert result == 0

def test_main_execution(mocker):
    """Test the main execution block"""
    # Mock sys.exit to avoid actually exiting
    mock_exit = mocker.patch('sys.exit')
    
    # Mock pytest.main to avoid actually running tests
    mock_pytest_main = mocker.patch('pytest.main', return_value=42)
    
    # Mock __name__ to trigger main block
    mocker.patch.object(backend.src.main, '__name__', '__main__')
    
    # Import main module to trigger __main__ block
    import importlib
    importlib.reload(backend.src.main)
    
    # Verify sys.exit was called with the return value from run_tests
    mock_exit.assert_called_once_with(42)

def test_main_execution_direct(mocker):
    """Test the main execution block directly"""
    # Mock sys.exit to avoid actually exiting
    mock_exit = mocker.patch('sys.exit')
    
    # Mock pytest.main to avoid actually running tests
    mock_pytest_main = mocker.patch('pytest.main', return_value=42)
    
    # Execute the main block directly
    backend.src.main.sys.exit(backend.src.main.run_tests())
    
    # Verify sys.exit was called with the return value from run_tests
    mock_exit.assert_called_once_with(42)

def test_main_execution_import(mocker):
    """Test the main execution block by importing"""
    # Mock sys.exit to avoid actually exiting
    mock_exit = mocker.patch('sys.exit')
    
    # Mock pytest.main to avoid actually running tests
    mock_pytest_main = mocker.patch('pytest.main', return_value=42)
    
    # Import main module with __name__ set to __main__
    import importlib
    import backend.src.main
    backend.src.main.__name__ = '__main__'
    importlib.reload(backend.src.main)
    
    # Verify sys.exit was called with the return value from run_tests
    mock_exit.assert_called_once_with(42) 