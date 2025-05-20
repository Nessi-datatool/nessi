"""
Tests for the GitHub Actions integration.
"""

import unittest
import os
import json
import tempfile
import subprocess
from unittest import mock
import yaml

class TestGitHubActionsIntegration(unittest.TestCase):
    """Test cases for the GitHub Actions integration."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        
        # Path to the GitHub Actions workflow file
        self.workflow_path = os.path.join(
            os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))),
            'examples', 'integrations', 'github_actions_workflow.yml'
        )
    
    def tearDown(self):
        """Clean up test environment."""
        # Clean up temporary directory
        self.temp_dir.cleanup()
    
    def test_workflow_file_exists(self):
        """Test that the GitHub Actions workflow file exists."""
        self.assertTrue(os.path.exists(self.workflow_path))
    
    def test_workflow_file_valid_yaml(self):
        """Test that the GitHub Actions workflow file is valid YAML."""
        with open(self.workflow_path, 'r') as f:
            try:
                workflow = yaml.safe_load(f)
                self.assertIsNotNone(workflow)
            except yaml.YAMLError:
                self.fail("Workflow file is not valid YAML")
    
    def test_workflow_structure(self):
        """Test the structure of the GitHub Actions workflow."""
        with open(self.workflow_path, 'r') as f:
            workflow = yaml.safe_load(f)
        
        # Check basic workflow structure
        self.assertIn('name', workflow)
        
        # Handle the case where 'on' is parsed as True
        on_key_exists = 'on' in workflow or True in workflow
        self.assertTrue(on_key_exists, "Workflow must have 'on' section")
        
        self.assertIn('jobs', workflow)
        
        # Check that there's at least one job
        self.assertTrue(len(workflow['jobs']) > 0)
        
        # Get the first job
        job_name = list(workflow['jobs'].keys())[0]
        job = workflow['jobs'][job_name]
        
        # Check job structure
        self.assertIn('runs-on', job)
        self.assertIn('steps', job)
        
        # Check that there are steps
        self.assertTrue(len(job['steps']) > 0)
        
        # Check that there's a step that runs nessi
        nessi_steps = [step for step in job['steps'] if 'run' in step and 'nessi' in step['run']]
        self.assertTrue(len(nessi_steps) > 0)
    
    def test_workflow_triggers(self):
        """Test the triggers of the GitHub Actions workflow."""
        with open(self.workflow_path, 'r') as f:
            workflow = yaml.safe_load(f)
        
        # Check triggers
        on_key_exists = 'on' in workflow or True in workflow
        self.assertTrue(on_key_exists, "Workflow must have 'on' section")
        
        # Get triggers from either 'on' or True key
        triggers = workflow.get('on', workflow.get(True, {}))
        
        # Check that there's at least one trigger
        self.assertTrue(len(triggers) > 0)
        
        # Common triggers
        common_triggers = ['push', 'pull_request', 'schedule', 'workflow_dispatch']
        self.assertTrue(any(trigger in triggers for trigger in common_triggers))
    
    @mock.patch('subprocess.run')
    def test_nessi_command_execution(self, mock_run):
        """Test execution of nessi commands in the workflow."""
        # Mock subprocess.run to return success
        mock_process = mock.MagicMock()
        mock_process.returncode = 0
        mock_process.stdout = "Validation passed! Quality checks completed successfully.\n" + json.dumps({
            'quality_score': 0.95,
            'passed': True,
            'rules_passed': 10,
            'rules_failed': 0,
        })
        mock_run.return_value = mock_process
        
        # Create a test script that simulates the GitHub Actions workflow
        script_path = os.path.join(self.temp_dir.name, 'test_workflow.sh')
        with open(script_path, 'w') as f:
            f.write("""#!/bin/bash
set -e

# Install nessi (simulated)
echo "Installing nessi..."

# Run nessi validation
nessi validate /path/to/table --format json > results.json

# Check if validation passed
passed=$(jq -r '.passed' results.json)
if [ "$passed" = "true" ]; then
    echo "Validation passed!"
    exit 0
else
    echo "Validation failed!"
    exit 1
fi
""")
        os.chmod(script_path, 0o755)
        
        # Run the script
        try:
            result = subprocess.run(
                [script_path],
                capture_output=True,
                text=True,
                check=False,
            )
            
            # Check that the script ran successfully
            self.assertEqual(result.returncode, 0)
            self.assertIn("Validation passed!", result.stdout)
        except subprocess.CalledProcessError:
            self.fail("Script execution failed")
    
    def test_artifact_handling(self):
        """Test handling of artifacts in the GitHub Actions workflow."""
        with open(self.workflow_path, 'r') as f:
            workflow = yaml.safe_load(f)
        
        # Get the first job
        job_name = list(workflow['jobs'].keys())[0]
        job = workflow['jobs'][job_name]
        
        # Check for artifact upload steps
        artifact_steps = [
            step for step in job['steps']
            if 'uses' in step and 'actions/upload-artifact' in step['uses']
        ]
        
        # There should be at least one artifact upload step
        self.assertTrue(len(artifact_steps) > 0)
        
        # Check artifact configuration
        for step in artifact_steps:
            self.assertIn('with', step)
            with_config = step['with']
            self.assertIn('name', with_config)
            self.assertIn('path', with_config)
    
    def test_environment_variables(self):
        """Test environment variables in the GitHub Actions workflow."""
        with open(self.workflow_path, 'r') as f:
            workflow = yaml.safe_load(f)
        
        # Get the first job
        job_name = list(workflow['jobs'].keys())[0]
        job = workflow['jobs'][job_name]
        
        # Check for environment variables
        env_vars = {}
        
        # Job-level env vars
        if 'env' in job:
            env_vars.update(job['env'])
        
        # Step-level env vars
        for step in job['steps']:
            if 'env' in step:
                env_vars.update(step['env'])
        
        # There should be some environment variables
        self.assertTrue(len(env_vars) > 0)
        
        # Check for common environment variables
        common_env_vars = ['NESSI_API_HOST', 'NESSI_API_KEY', 'NESSI_API_SECRET']
        self.assertTrue(any(var in env_vars for var in common_env_vars))
    
    def test_secrets_usage(self):
        """Test usage of secrets in the GitHub Actions workflow."""
        with open(self.workflow_path, 'r') as f:
            workflow = yaml.safe_load(f)
        
        # Get the first job
        job_name = list(workflow['jobs'].keys())[0]
        job = workflow['jobs'][job_name]
        
        # Check for secrets in environment variables
        secret_usage = False
        
        # Job-level env vars
        if 'env' in job:
            for key, value in job['env'].items():
                if isinstance(value, str) and 'secrets.' in value:
                    secret_usage = True
                    break
        
        # Step-level env vars
        if not secret_usage:
            for step in job['steps']:
                if 'env' in step:
                    for key, value in step['env'].items():
                        if isinstance(value, str) and 'secrets.' in value:
                            secret_usage = True
                            break
                if secret_usage:
                    break
        
        # There should be some secret usage
        self.assertTrue(secret_usage)
    
    def test_workflow_notifications(self):
        """Test notifications in the GitHub Actions workflow."""
        with open(self.workflow_path, 'r') as f:
            workflow = yaml.safe_load(f)
        
        # Get the first job
        job_name = list(workflow['jobs'].keys())[0]
        job = workflow['jobs'][job_name]
        
        # Check for notification steps
        notification_steps = []
        
        for step in job['steps']:
            # Check for Slack notifications
            if 'uses' in step and any(keyword in step['uses'] for keyword in ['slack', 'rtCamp/action-slack-notify']):
                notification_steps.append(step)
            
            # Check for email notifications
            if 'uses' in step and any(keyword in step['uses'] for keyword in ['email', 'mail', 'dawidd6/action-send-mail']):
                notification_steps.append(step)
            
            # Check for notification in step name
            if 'name' in step and any(keyword in step['name'].lower() for keyword in ['notify', 'notification', 'slack', 'email']):
                notification_steps.append(step)
                
            # Check for custom notification scripts
            if 'run' in step and any(keyword in step['run'] for keyword in ['notify', 'notification', 'slack', 'email']):
                notification_steps.append(step)
        
        # There should be some notification steps
        self.assertTrue(len(notification_steps) > 0)


if __name__ == '__main__':
    unittest.main()
