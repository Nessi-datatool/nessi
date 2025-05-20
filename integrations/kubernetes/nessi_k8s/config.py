"""
Configuration utilities for Nessi Kubernetes integration.

This module provides utilities for configuring the Nessi Kubernetes integration.
"""

import os
import json
import logging
import yaml
from typing import Dict, Any, Optional, List, Union, Tuple

logger = logging.getLogger(__name__)


class NessiK8sConfig:
    """
    Configuration manager for Nessi Kubernetes integration.
    
    This class provides utilities for loading and managing configuration
    for the Nessi Kubernetes integration.
    """
    
    DEFAULT_CONFIG = {
        "namespace": "default",
        "image": "nessi/nessi:latest",
        "service_account_name": None,
        "api_host": None,
        "api_key": None,
        "api_secret": None,
        "job_ttl_seconds_after_finished": 3600,
        "job_backoff_limit": 3,
        "job_timeout_seconds": 600,
        "resources": {
            "requests": {"cpu": "100m", "memory": "128Mi"},
            "limits": {"cpu": "500m", "memory": "512Mi"},
        },
        "node_selector": None,
        "tolerations": None,
        "labels": {"app": "nessi"},
        "annotations": {},
        "env_from": [],
        "volumes": [],
        "volume_mounts": [],
        "image_pull_secrets": [],
        "config_file": None,
        "context": None,
    }
    
    def __init__(
        self,
        config_file: Optional[str] = None,
        env_prefix: str = "NESSI_K8S_",
    ):
        """
        Initialize the configuration manager.
        
        Args:
            config_file: Path to configuration file
            env_prefix: Prefix for environment variables
        """
        self.config_file = config_file
        self.env_prefix = env_prefix
        self.config = self.DEFAULT_CONFIG.copy()
        
        # Load configuration
        self._load_config()
    
    def _load_config(self) -> None:
        """
        Load configuration from file and environment variables.
        """
        # Load from file if provided
        if self.config_file and os.path.exists(self.config_file):
            self._load_from_file(self.config_file)
        
        # Load from environment variables
        self._load_from_env()
    
    def _load_from_file(self, config_file: str) -> None:
        """
        Load configuration from file.
        
        Args:
            config_file: Path to configuration file
        """
        try:
            with open(config_file, "r") as f:
                if config_file.endswith(".json"):
                    file_config = json.load(f)
                elif config_file.endswith((".yaml", ".yml")):
                    file_config = yaml.safe_load(f)
                else:
                    logger.warning(f"Unsupported config file format: {config_file}")
                    return
                
                # Update configuration
                self.config.update(file_config)
                logger.info(f"Loaded configuration from {config_file}")
        except Exception as e:
            logger.error(f"Error loading configuration from {config_file}: {e}")
    
    def _load_from_env(self) -> None:
        """
        Load configuration from environment variables.
        """
        # Map of environment variable names to config keys
        env_map = {
            "NAMESPACE": "namespace",
            "IMAGE": "image",
            "SERVICE_ACCOUNT_NAME": "service_account_name",
            "API_HOST": "api_host",
            "API_KEY": "api_key",
            "API_SECRET": "api_secret",
            "JOB_TTL_SECONDS": "job_ttl_seconds_after_finished",
            "JOB_BACKOFF_LIMIT": "job_backoff_limit",
            "JOB_TIMEOUT_SECONDS": "job_timeout_seconds",
            "CONFIG_FILE": "config_file",
            "CONTEXT": "context",
        }
        
        # Load simple values
        for env_key, config_key in env_map.items():
            env_var = f"{self.env_prefix}{env_key}"
            if env_var in os.environ:
                value = os.environ[env_var]
                
                # Convert numeric values
                if config_key in [
                    "job_ttl_seconds_after_finished",
                    "job_backoff_limit",
                    "job_timeout_seconds",
                ]:
                    try:
                        value = int(value)
                    except ValueError:
                        logger.warning(f"Invalid numeric value for {env_var}: {value}")
                        continue
                
                self.config[config_key] = value
                logger.debug(f"Loaded {config_key} from {env_var}")
        
        # Load complex values
        for env_key in ["RESOURCES", "NODE_SELECTOR", "TOLERATIONS", "LABELS", "ANNOTATIONS"]:
            env_var = f"{self.env_prefix}{env_key}"
            if env_var in os.environ:
                try:
                    value = json.loads(os.environ[env_var])
                    config_key = env_key.lower()
                    self.config[config_key] = value
                    logger.debug(f"Loaded {config_key} from {env_var}")
                except json.JSONDecodeError:
                    logger.warning(f"Invalid JSON in {env_var}")
    
    def get_config(self) -> Dict[str, Any]:
        """
        Get the current configuration.
        
        Returns:
            Current configuration
        """
        return self.config.copy()
    
    def update_config(self, config: Dict[str, Any]) -> None:
        """
        Update the configuration.
        
        Args:
            config: Configuration to update
        """
        self.config.update(config)
    
    def save_config(self, config_file: Optional[str] = None) -> None:
        """
        Save the configuration to a file.
        
        Args:
            config_file: Path to configuration file
        """
        config_file = config_file or self.config_file
        if not config_file:
            logger.warning("No config file specified")
            return
        
        try:
            with open(config_file, "w") as f:
                if config_file.endswith(".json"):
                    json.dump(self.config, f, indent=2)
                elif config_file.endswith((".yaml", ".yml")):
                    yaml.dump(self.config, f)
                else:
                    logger.warning(f"Unsupported config file format: {config_file}")
                    return
                
                logger.info(f"Saved configuration to {config_file}")
        except Exception as e:
            logger.error(f"Error saving configuration to {config_file}: {e}")


def load_config(config_file: Optional[str] = None) -> Dict[str, Any]:
    """
    Load configuration for Nessi Kubernetes integration.
    
    Args:
        config_file: Path to configuration file
        
    Returns:
        Configuration dictionary
    """
    config_manager = NessiK8sConfig(config_file=config_file)
    return config_manager.get_config()


def create_default_config(config_file: str, format: str = "yaml") -> None:
    """
    Create a default configuration file.
    
    Args:
        config_file: Path to configuration file
        format: Configuration format (yaml or json)
    """
    config_manager = NessiK8sConfig()
    
    # Ensure the file has the correct extension
    if format.lower() == "json" and not config_file.endswith(".json"):
        config_file = f"{config_file}.json"
    elif format.lower() == "yaml" and not config_file.endswith((".yaml", ".yml")):
        config_file = f"{config_file}.yaml"
    
    config_manager.save_config(config_file)
