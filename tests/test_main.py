"""
NESSI - FREE FOR PERSONAL USE LICENSE

Copyright (c) 2025 Nessi. All rights reserved.

Nessi is free for personal use.

A paid license is required for enterprise or consulting use.

This version is free. Future versions of Nessi will require a license for all users.

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
   The Software is licensed, not sold. Nessi retains all right, title, and interest in and to the Software, including all intellectual property rights.

4. DISCLAIMER OF WARRANTY
   THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT.

5. LIMITATION OF LIABILITY
   IN NO EVENT SHALL NESSI BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
"""

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