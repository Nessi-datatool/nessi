"""
Tests for the Nessi Kubernetes CLI module.
"""

import os
import sys
import json
import tempfile
import unittest
from unittest import mock
from io import StringIO

from nessi_k8s.cli import main


class TestNessiK8sCLI(unittest.TestCase):
    """Test cases for the Nessi Kubernetes CLI."""

    def setUp(self):
        """Set up test environment."""
        # Save original sys.argv and stdout
        self.original_argv = sys.argv
        self.original_stdout = sys.stdout
        self.original_stderr = sys.stderr
        
        # Redirect stdout and stderr
        sys.stdout = StringIO()
        sys.stderr = StringIO()
        
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        
        # Mock NessiK8sClient
        self.mock_client_patcher = mock.patch('nessi_k8s.cli.NessiK8sClient')
        self.mock_client = self.mock_client_patcher.start()
        self.mock_client_instance = self.mock_client.return_value
        
        # Mock create_default_config
        self.mock_create_config_patcher = mock.patch('nessi_k8s.cli.create_default_config')
        self.mock_create_config = self.mock_create_config_patcher.start()
    
    def tearDown(self):
        """Clean up test environment."""
        # Restore original sys.argv and stdout
        sys.argv = self.original_argv
        sys.stdout = self.original_stdout
        sys.stderr = self.original_stderr
        
        # Clean up temporary directory
        self.temp_dir.cleanup()
        
        # Stop mocks
        self.mock_client_patcher.stop()
        self.mock_create_config_patcher.stop()
    
    def test_create_config_command(self):
        """Test the create-config command."""
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'create-config',
            '--output', os.path.join(self.temp_dir.name, 'test-config.yaml'),
            '--format', 'yaml',
        ]
        
        # Run CLI
        main()
        
        # Check that create_default_config was called correctly
        self.mock_create_config.assert_called_once_with(
            os.path.join(self.temp_dir.name, 'test-config.yaml'),
            'yaml',
        )
    
    def test_quality_check_command(self):
        """Test the quality-check command."""
        # Create a test rules file
        rules_file = os.path.join(self.temp_dir.name, 'test-rules.json')
        with open(rules_file, 'w') as f:
            json.dump([
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ], f)
        
        # Mock client response
        self.mock_client_instance.run_quality_check.return_value = {
            "status": "Succeeded",
            "job_name": "test-quality-job",
            "logs": "Quality check results",
        }
        
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'quality-check',
            '--table', 'test-table',
            '--rules', rules_file,
            '--profile',
            '--output-format', 'json',
            '--wait',
            '--timeout', '600',
            '--job-name', 'test-job',
            '--namespace', 'test-namespace',
            '--image', 'test-image:latest',
            '--api-host', 'http://test-host',
            '--api-key', 'test-key',
            '--api-secret', 'test-secret',
        ]
        
        # Run CLI
        main()
        
        # Check that client was created with correct arguments
        self.mock_client.assert_called_once_with(
            config_file=None,
            namespace='test-namespace',
            image='test-image:latest',
            api_host='http://test-host',
            api_key='test-key',
            api_secret='test-secret',
            service_account_name=None,
        )
        
        # Check that run_quality_check was called correctly
        self.mock_client_instance.run_quality_check.assert_called_once_with(
            table_name='test-table',
            rules=[{"name": "test-rule", "rule_type": "not_null", "column": "id"}],
            profile=True,
            output_format='json',
            output_path=None,
            wait_for_completion=True,
            timeout=600,
            job_name='test-job',
        )
    
    def test_profile_command(self):
        """Test the profile command."""
        # Mock client response
        self.mock_client_instance.run_profile.return_value = {
            "status": "Succeeded",
            "job_name": "test-profile-job",
            "logs": "Profile results",
        }
        
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'profile',
            '--table', 'test-table',
            '--columns', 'id,name,email',
            '--sample-size', '1000',
            '--output-format', 'json',
            '--wait',
            '--timeout', '600',
            '--job-name', 'test-job',
            '--namespace', 'test-namespace',
            '--image', 'test-image:latest',
            '--api-host', 'http://test-host',
            '--api-key', 'test-key',
            '--api-secret', 'test-secret',
        ]
        
        # Run CLI
        main()
        
        # Check that run_profile was called correctly
        self.mock_client_instance.run_profile.assert_called_once_with(
            table_name='test-table',
            columns=['id', 'name', 'email'],
            sample_size=1000,
            output_format='json',
            output_path=None,
            wait_for_completion=True,
            timeout=600,
            job_name='test-job',
        )
    
    def test_validate_command(self):
        """Test the validate command."""
        # Create a test rules file
        rules_file = os.path.join(self.temp_dir.name, 'test-rules.json')
        with open(rules_file, 'w') as f:
            json.dump([
                {"name": "test-rule", "rule_type": "not_null", "column": "id"},
            ], f)
        
        # Mock client response
        self.mock_client_instance.run_validation.return_value = {
            "status": "Succeeded",
            "job_name": "test-validation-job",
            "logs": "Validation results",
        }
        
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'validate',
            '--table', 'test-table',
            '--rules', rules_file,
            '--output-format', 'json',
            '--wait',
            '--timeout', '600',
            '--job-name', 'test-job',
            '--namespace', 'test-namespace',
            '--image', 'test-image:latest',
            '--api-host', 'http://test-host',
            '--api-key', 'test-key',
            '--api-secret', 'test-secret',
        ]
        
        # Run CLI
        main()
        
        # Check that run_validation was called correctly
        self.mock_client_instance.run_validation.assert_called_once_with(
            table_name='test-table',
            rules=[{"name": "test-rule", "rule_type": "not_null", "column": "id"}],
            output_format='json',
            output_path=None,
            wait_for_completion=True,
            timeout=600,
            job_name='test-job',
        )
    
    def test_get_results_command(self):
        """Test the get-results command."""
        # Mock client response
        self.mock_client_instance.get_quality_results.return_value = {
            "status": "Succeeded",
            "job_name": "test-quality-job",
            "logs": "Quality check results",
        }
        
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'get-results',
            '--job-name', 'test-quality-job',
            '--job-type', 'quality',
            '--wait',
            '--timeout', '600',
            '--namespace', 'test-namespace',
            '--image', 'test-image:latest',
            '--api-host', 'http://test-host',
            '--api-key', 'test-key',
            '--api-secret', 'test-secret',
        ]
        
        # Run CLI
        main()
        
        # Check that get_quality_results was called correctly
        self.mock_client_instance.get_quality_results.assert_called_once_with(
            job_name='test-quality-job',
            wait_for_completion=True,
            timeout=600,
        )
    
    def test_get_results_profile_command(self):
        """Test the get-results command for profile jobs."""
        # Mock client response
        self.mock_client_instance.get_profile_results.return_value = {
            "status": "Succeeded",
            "job_name": "test-profile-job",
            "logs": "Profile results",
        }
        
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'get-results',
            '--job-name', 'test-profile-job',
            '--job-type', 'profile',
            '--wait',
            '--timeout', '600',
        ]
        
        # Run CLI
        main()
        
        # Check that get_profile_results was called correctly
        self.mock_client_instance.get_profile_results.assert_called_once_with(
            job_name='test-profile-job',
            wait_for_completion=True,
            timeout=600,
        )
    
    def test_get_results_validation_command(self):
        """Test the get-results command for validation jobs."""
        # Mock client response
        self.mock_client_instance.get_validation_results.return_value = {
            "status": "Succeeded",
            "job_name": "test-validation-job",
            "logs": "Validation results",
        }
        
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
            'get-results',
            '--job-name', 'test-validation-job',
            '--job-type', 'validation',
            '--wait',
            '--timeout', '600',
        ]
        
        # Run CLI
        main()
        
        # Check that get_validation_results was called correctly
        self.mock_client_instance.get_validation_results.assert_called_once_with(
            job_name='test-validation-job',
            wait_for_completion=True,
            timeout=600,
        )
    
    def test_no_command(self):
        """Test CLI with no command."""
        # Set up command line arguments
        sys.argv = [
            'nessi_k8s',
        ]
        
        # Run CLI and check exit code
        exit_code = main()
        self.assertEqual(exit_code, 1)
        
        # Check that error message was printed
        stderr = sys.stderr.getvalue()
        self.assertIn("error", stderr.lower())


if __name__ == "__main__":
    unittest.main()
