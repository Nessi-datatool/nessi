"""
Client for Nessi Kubernetes integration.

This module provides a high-level client for interacting with
Nessi operations in Kubernetes.
"""

import os
import logging
from typing import Dict, Any, Optional, List, Union, Tuple

from nessi_k8s.operator import (
    NessiDataQualityOperator,
    NessiProfileOperator,
    NessiValidationOperator,
)
from nessi_k8s.config import load_config

logger = logging.getLogger(__name__)


class NessiK8sClient:
    """
    Client for Nessi Kubernetes integration.
    
    This class provides a high-level API for running Nessi operations
    in Kubernetes.
    """
    
    def __init__(
        self,
        config_file: Optional[str] = None,
        namespace: Optional[str] = None,
        image: Optional[str] = None,
        service_account_name: Optional[str] = None,
        api_host: Optional[str] = None,
        api_key: Optional[str] = None,
        api_secret: Optional[str] = None,
        **kwargs,
    ):
        """
        Initialize the Nessi Kubernetes client.
        
        Args:
            config_file: Path to configuration file
            namespace: Kubernetes namespace to run jobs in
            image: Docker image to use for jobs
            service_account_name: Service account to use for jobs
            api_host: Nessi API host URL
            api_key: Nessi API key
            api_secret: Nessi API secret
            **kwargs: Additional configuration options
        """
        # Load configuration
        self.config = load_config(config_file)
        
        # Override configuration with arguments
        if namespace:
            self.config["namespace"] = namespace
        if image:
            self.config["image"] = image
        if service_account_name:
            self.config["service_account_name"] = service_account_name
        if api_host:
            self.config["api_host"] = api_host
        if api_key:
            self.config["api_key"] = api_key
        if api_secret:
            self.config["api_secret"] = api_secret
        
        # Update with additional kwargs
        self.config.update(kwargs)
        
        # Initialize operators
        self.quality_operator = NessiDataQualityOperator(**self.config)
        self.profile_operator = NessiProfileOperator(**self.config)
        self.validation_operator = NessiValidationOperator(**self.config)
    
    def run_quality_check(
        self,
        table_name: str,
        rules: Optional[List[Dict[str, Any]]] = None,
        profile: bool = True,
        output_format: str = "json",
        output_path: Optional[str] = None,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
        job_name: Optional[str] = None,
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> Dict[str, Any]:
        """
        Run a data quality check as a Kubernetes job.
        
        Args:
            table_name: Name of the table to check
            rules: List of quality rules to apply
            profile: Whether to generate a profile
            output_format: Output format (json, yaml, text)
            output_path: Path to write output to
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            job_name: Name for the job (generated if not provided)
            env_vars: Additional environment variables
            
        Returns:
            Job status or job name if not waiting for completion
        """
        return self.quality_operator.run_quality_check(
            table_name=table_name,
            rules=rules,
            profile=profile,
            output_format=output_format,
            output_path=output_path,
            wait_for_completion=wait_for_completion,
            timeout=timeout,
            poll_interval=poll_interval,
            job_name=job_name,
            env_vars=env_vars,
        )
    
    def get_quality_results(
        self,
        job_name: str,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Get results of a data quality check job.
        
        Args:
            job_name: Name of the job
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status and logs
        """
        return self.quality_operator.get_quality_results(
            job_name=job_name,
            wait_for_completion=wait_for_completion,
            timeout=timeout,
            poll_interval=poll_interval,
        )
    
    def run_profile(
        self,
        table_name: str,
        columns: Optional[List[str]] = None,
        sample_size: Optional[int] = None,
        output_format: str = "json",
        output_path: Optional[str] = None,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
        job_name: Optional[str] = None,
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> Dict[str, Any]:
        """
        Run a data profile as a Kubernetes job.
        
        Args:
            table_name: Name of the table to profile
            columns: List of columns to profile
            sample_size: Number of rows to sample
            output_format: Output format (json, yaml, text)
            output_path: Path to write output to
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            job_name: Name for the job (generated if not provided)
            env_vars: Additional environment variables
            
        Returns:
            Job status or job name if not waiting for completion
        """
        return self.profile_operator.run_profile(
            table_name=table_name,
            columns=columns,
            sample_size=sample_size,
            output_format=output_format,
            output_path=output_path,
            wait_for_completion=wait_for_completion,
            timeout=timeout,
            poll_interval=poll_interval,
            job_name=job_name,
            env_vars=env_vars,
        )
    
    def get_profile_results(
        self,
        job_name: str,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Get results of a data profile job.
        
        Args:
            job_name: Name of the job
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status and logs
        """
        return self.profile_operator.get_profile_results(
            job_name=job_name,
            wait_for_completion=wait_for_completion,
            timeout=timeout,
            poll_interval=poll_interval,
        )
    
    def run_validation(
        self,
        table_name: str,
        rules: List[Dict[str, Any]],
        output_format: str = "json",
        output_path: Optional[str] = None,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
        job_name: Optional[str] = None,
        env_vars: Optional[List[Dict[str, str]]] = None,
    ) -> Dict[str, Any]:
        """
        Run a data validation as a Kubernetes job.
        
        Args:
            table_name: Name of the table to validate
            rules: List of validation rules
            output_format: Output format (json, yaml, text)
            output_path: Path to write output to
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            job_name: Name for the job (generated if not provided)
            env_vars: Additional environment variables
            
        Returns:
            Job status or job name if not waiting for completion
        """
        return self.validation_operator.run_validation(
            table_name=table_name,
            rules=rules,
            output_format=output_format,
            output_path=output_path,
            wait_for_completion=wait_for_completion,
            timeout=timeout,
            poll_interval=poll_interval,
            job_name=job_name,
            env_vars=env_vars,
        )
    
    def get_validation_results(
        self,
        job_name: str,
        wait_for_completion: bool = True,
        timeout: int = 600,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Get results of a data validation job.
        
        Args:
            job_name: Name of the job
            wait_for_completion: Whether to wait for job completion
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds
            
        Returns:
            Job status and logs
        """
        return self.validation_operator.get_validation_results(
            job_name=job_name,
            wait_for_completion=wait_for_completion,
            timeout=timeout,
            poll_interval=poll_interval,
        )
