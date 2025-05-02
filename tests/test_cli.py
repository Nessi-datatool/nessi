import pytest
from click.testing import CliRunner
from backend.src.scanner.cli import cli
from backend.src.scanner import __version__

@pytest.fixture
def runner():
    return CliRunner()

def test_version_command(runner):
    """Test the version command"""
    result = runner.invoke(cli, ['version'])
    assert result.exit_code == 0
    assert f"Nessi version {__version__}" in result.output

def test_status_command(runner):
    """Test the status command"""
    result = runner.invoke(cli, ['status'])
    assert result.exit_code == 0
    assert "message" in result.output.lower()

def test_activate_valid_license(runner):
    """Test activating with a valid license key"""
    result = runner.invoke(cli, ['activate', 'valid-key-123'])
    assert result.exit_code == 0
    assert "License activated successfully" in result.output

def test_activate_invalid_license(runner):
    """Test activating with an invalid license key"""
    result = runner.invoke(cli, ['activate', 'invalid-key'])
    assert result.exit_code == 0
    assert "Invalid license key" in result.output 