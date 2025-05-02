import pytest
import sys
import os

def run_tests():
    """Run pytest with the test directory"""
    test_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), "tests")
    return pytest.main(["-v", test_dir])

if __name__ == "__main__":
    sys.exit(run_tests()) 