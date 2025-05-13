"""
Tests for the Nessi Kubernetes configuration module.
"""

import os
import json
import tempfile
import unittest
from unittest import mock
import yaml

from nessi_k8s.config import NessiK8sConfig, load_config, create_default_config


class TestNessiK8sConfig(unittest.TestCase):
    """Test cases for the NessiK8sConfig class."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        
        # Clear environment variables that might interfere with tests
        self.env_backup = {}
        for key in list(os.environ.keys()):
            if key.startswith("NESSI_K8S_"):
                self.env_backup[key] = os.environ[key]
                del os.environ[key]
    
    def tearDown(self):
        """Clean up test environment."""
        # Restore environment variables
        for key, value in self.env_backup.items():
            os.environ[key] = value
        
        # Clean up temporary directory
        self.temp_dir.cleanup()
    
    def test_default_config(self):
        """Test default configuration values."""
        config = NessiK8sConfig()
        
        # Check default values
        self.assertEqual(config.config["namespace"], "default")
        self.assertEqual(config.config["image"], "nessi/nessi:latest")
        self.assertEqual(config.config["job_ttl_seconds_after_finished"], 3600)
        self.assertEqual(config.config["job_backoff_limit"], 3)
        self.assertEqual(config.config["job_timeout_seconds"], 600)
        
        # Check default resources
        self.assertEqual(config.config["resources"]["requests"]["cpu"], "100m")
        self.assertEqual(config.config["resources"]["requests"]["memory"], "128Mi")
        self.assertEqual(config.config["resources"]["limits"]["cpu"], "500m")
        self.assertEqual(config.config["resources"]["limits"]["memory"], "512Mi")
    
    def test_load_from_json_file(self):
        """Test loading configuration from a JSON file."""
        # Create a test JSON config file
        config_data = {
            "namespace": "test-namespace",
            "image": "test-image:latest",
            "api_host": "http://test-host",
            "api_key": "test-key",
        }
        
        config_file = os.path.join(self.temp_dir.name, "test-config.json")
        with open(config_file, "w") as f:
            json.dump(config_data, f)
        
        # Load config from file
        config = NessiK8sConfig(config_file=config_file)
        
        # Check loaded values
        self.assertEqual(config.config["namespace"], "test-namespace")
        self.assertEqual(config.config["image"], "test-image:latest")
        self.assertEqual(config.config["api_host"], "http://test-host")
        self.assertEqual(config.config["api_key"], "test-key")
        
        # Check that defaults are preserved for unspecified values
        self.assertEqual(config.config["job_ttl_seconds_after_finished"], 3600)
    
    def test_load_from_yaml_file(self):
        """Test loading configuration from a YAML file."""
        # Create a test YAML config file
        config_data = """
        namespace: test-namespace-yaml
        image: test-image-yaml:latest
        api_host: http://test-host-yaml
        api_key: test-key-yaml
        resources:
          requests:
            cpu: 200m
            memory: 256Mi
        """
        
        config_file = os.path.join(self.temp_dir.name, "test-config.yaml")
        with open(config_file, "w") as f:
            f.write(config_data)
        
        # Load config from file
        config = NessiK8sConfig(config_file=config_file)
        
        # Check loaded values
        self.assertEqual(config.config["namespace"], "test-namespace-yaml")
        self.assertEqual(config.config["image"], "test-image-yaml:latest")
        self.assertEqual(config.config["api_host"], "http://test-host-yaml")
        self.assertEqual(config.config["api_key"], "test-key-yaml")
        
        # Check that resource values are loaded correctly
        self.assertEqual(config.config["resources"]["requests"]["cpu"], "200m")
        self.assertEqual(config.config["resources"]["requests"]["memory"], "256Mi")
    
    def test_load_from_env_variables(self):
        """Test loading configuration from environment variables."""
        # Set environment variables
        os.environ["NESSI_K8S_NAMESPACE"] = "env-namespace"
        os.environ["NESSI_K8S_IMAGE"] = "env-image:latest"
        os.environ["NESSI_K8S_API_HOST"] = "http://env-host"
        os.environ["NESSI_K8S_API_KEY"] = "env-key"
        os.environ["NESSI_K8S_JOB_TTL_SECONDS"] = "7200"
        os.environ["NESSI_K8S_RESOURCES"] = json.dumps({
            "requests": {"cpu": "300m", "memory": "512Mi"},
            "limits": {"cpu": "1000m", "memory": "1Gi"}
        })
        
        # Load config from environment
        config = NessiK8sConfig()
        
        # Check loaded values
        self.assertEqual(config.config["namespace"], "env-namespace")
        self.assertEqual(config.config["image"], "env-image:latest")
        self.assertEqual(config.config["api_host"], "http://env-host")
        self.assertEqual(config.config["api_key"], "env-key")
        self.assertEqual(config.config["job_ttl_seconds_after_finished"], 7200)
        
        # Check that resource values are loaded correctly
        self.assertEqual(config.config["resources"]["requests"]["cpu"], "300m")
        self.assertEqual(config.config["resources"]["requests"]["memory"], "512Mi")
        self.assertEqual(config.config["resources"]["limits"]["cpu"], "1000m")
        self.assertEqual(config.config["resources"]["limits"]["memory"], "1Gi")
    
    def test_precedence_env_over_file(self):
        """Test that environment variables take precedence over file configuration."""
        # Create a test config file
        config_data = {
            "namespace": "file-namespace",
            "image": "file-image:latest",
            "api_host": "http://file-host",
            "api_key": "file-key",
        }
        
        config_file = os.path.join(self.temp_dir.name, "test-config.json")
        with open(config_file, "w") as f:
            json.dump(config_data, f)
        
        # Set environment variables
        os.environ["NESSI_K8S_NAMESPACE"] = "env-namespace"
        os.environ["NESSI_K8S_API_HOST"] = "http://env-host"
        
        # Load config
        config = NessiK8sConfig(config_file=config_file)
        
        # Check that env vars override file values
        self.assertEqual(config.config["namespace"], "env-namespace")
        self.assertEqual(config.config["api_host"], "http://env-host")
        
        # Check that file values are used for unspecified env vars
        self.assertEqual(config.config["image"], "file-image:latest")
        self.assertEqual(config.config["api_key"], "file-key")
    
    def test_update_config(self):
        """Test updating configuration."""
        config = NessiK8sConfig()
        
        # Update config
        config.update_config({
            "namespace": "updated-namespace",
            "image": "updated-image:latest",
            "resources": {
                "requests": {"cpu": "400m"},
            },
        })
        
        # Check updated values
        self.assertEqual(config.config["namespace"], "updated-namespace")
        self.assertEqual(config.config["image"], "updated-image:latest")
        self.assertEqual(config.config["resources"]["requests"]["cpu"], "400m")
        
        # Check that unspecified values are preserved
        self.assertEqual(config.config["job_ttl_seconds_after_finished"], 3600)
        self.assertEqual(config.config["resources"]["requests"]["memory"], "128Mi")
    
    def test_save_config_json(self):
        """Test saving configuration to a JSON file."""
        config = NessiK8sConfig()
        
        # Update config
        config.update_config({
            "namespace": "save-namespace",
            "image": "save-image:latest",
        })
        
        # Save config
        config_file = os.path.join(self.temp_dir.name, "save-config.json")
        config.save_config(config_file)
        
        # Load saved config
        with open(config_file, "r") as f:
            saved_config = json.load(f)
        
        # Check saved values
        self.assertEqual(saved_config["namespace"], "save-namespace")
        self.assertEqual(saved_config["image"], "save-image:latest")
    
    def test_save_config_yaml(self):
        """Test saving configuration to a YAML file."""
        config = NessiK8sConfig()
        
        # Update config
        config.update_config({
            "namespace": "save-namespace-yaml",
            "image": "save-image-yaml:latest",
        })
        
        # Save config
        config_file = os.path.join(self.temp_dir.name, "save-config.yaml")
        config.save_config(config_file)
        
        # Load saved config
        with open(config_file, "r") as f:
            saved_config = yaml.safe_load(f)
        
        # Check saved values
        self.assertEqual(saved_config["namespace"], "save-namespace-yaml")
        self.assertEqual(saved_config["image"], "save-image-yaml:latest")


class TestConfigHelpers(unittest.TestCase):
    """Test cases for the configuration helper functions."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
    
    def tearDown(self):
        """Clean up test environment."""
        # Clean up temporary directory
        self.temp_dir.cleanup()
    
    def test_load_config(self):
        """Test the load_config helper function."""
        # Create a test config file
        config_data = {
            "namespace": "helper-namespace",
            "image": "helper-image:latest",
        }
        
        config_file = os.path.join(self.temp_dir.name, "helper-config.json")
        with open(config_file, "w") as f:
            json.dump(config_data, f)
        
        # Load config
        config = load_config(config_file)
        
        # Check loaded values
        self.assertEqual(config["namespace"], "helper-namespace")
        self.assertEqual(config["image"], "helper-image:latest")
        
        # Check that defaults are included
        self.assertEqual(config["job_ttl_seconds_after_finished"], 3600)
    
    def test_create_default_config_json(self):
        """Test the create_default_config helper function with JSON format."""
        # Create default config
        config_file = os.path.join(self.temp_dir.name, "default-config.json")
        create_default_config(config_file, format="json")
        
        # Check that file exists
        self.assertTrue(os.path.exists(config_file))
        
        # Load created config
        with open(config_file, "r") as f:
            created_config = json.load(f)
        
        # Check default values
        self.assertEqual(created_config["namespace"], "default")
        self.assertEqual(created_config["image"], "nessi/nessi:latest")
    
    def test_create_default_config_yaml(self):
        """Test the create_default_config helper function with YAML format."""
        # Create default config
        config_file = os.path.join(self.temp_dir.name, "default-config.yaml")
        create_default_config(config_file, format="yaml")
        
        # Check that file exists
        self.assertTrue(os.path.exists(config_file))
        
        # Load created config
        with open(config_file, "r") as f:
            created_config = yaml.safe_load(f)
        
        # Check default values
        self.assertEqual(created_config["namespace"], "default")
        self.assertEqual(created_config["image"], "nessi/nessi:latest")


if __name__ == "__main__":
    unittest.main()
