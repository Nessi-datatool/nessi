"""
Nessi.dev resource for Dagster.

This module provides a resource for interacting with the Nessi.dev API.
"""

from typing import Dict, Any, Optional, List, Union
import json
import time
import logging

import requests
from dagster import resource, ConfigResource


class NessiClient:
    """
    Client for Nessi.dev API.

    This client handles authentication and provides methods for interacting
    with the Nessi.dev API.
    """

    def __init__(
        self,
        api_host: str,
        api_key: str,
        api_secret: Optional[str] = None,
        timeout: int = 60,
    ):
        """
        Initialize the Nessi client.

        Args:
            api_host: Host URL for the Nessi.dev API
            api_key: API key for authentication
            api_secret: API secret for authentication (optional)
            timeout: Timeout for API requests in seconds
        """
        self.api_host = api_host.rstrip('/')
        self.api_key = api_key
        self.api_secret = api_secret
        self.timeout = timeout
        self.session = self._create_session()
        self.logger = logging.getLogger("dagster")

    def _create_session(self) -> requests.Session:
        """
        Create a requests session for making API calls.

        Returns:
            requests.Session instance
        """
        session = requests.Session()
        
        # Set authentication headers
        if self.api_key:
            session.headers.update({"X-API-Key": self.api_key})
            if self.api_secret:
                session.headers.update({"X-API-Secret": self.api_secret})
        else:
            raise ValueError("API key is required")
        
        # Set additional headers
        session.headers.update({
            "Content-Type": "application/json",
            "Accept": "application/json",
            "User-Agent": f"nessi-dagster/{__name__}",
        })
        
        return session

    def _do_api_call(
        self,
        endpoint: str,
        method: str = "GET",
        data: Optional[Union[Dict[str, Any], List[Any]]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Perform an API call to Nessi.dev.

        Args:
            endpoint: API endpoint (without base URL)
            method: HTTP method
            data: Request payload
            params: Query parameters

        Returns:
            Response as dictionary
        """
        url = f"{self.api_host}/{endpoint.lstrip('/')}"
        
        try:
            if method == "GET":
                response = self.session.get(
                    url,
                    params=params,
                    timeout=self.timeout,
                )
            elif method == "POST":
                response = self.session.post(
                    url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                )
            elif method == "PUT":
                response = self.session.put(
                    url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                )
            elif method == "DELETE":
                response = self.session.delete(
                    url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                )
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")
            
            response.raise_for_status()
            
            if response.content:
                return response.json()
            return {}
            
        except requests.exceptions.HTTPError as e:
            self.logger.error("HTTP error: %s", e)
            if e.response.content:
                try:
                    error_response = e.response.json()
                    self.logger.error("Error response: %s", error_response)
                except ValueError:
                    self.logger.error("Error response: %s", e.response.content)
            raise RuntimeError(f"Nessi API HTTP error: {e}")
        except requests.exceptions.RequestException as e:
            self.logger.error("Request error: %s", e)
            raise RuntimeError(f"Nessi API request error: {e}")

    # Health and status methods
    
    def get_health(self) -> Dict[str, Any]:
        """
        Get Nessi.dev API health status.

        Returns:
            Health status information
        """
        return self._do_api_call("health")
    
    # Data quality methods
    
    def run_quality_check(
        self,
        table_name: str,
        rules: Optional[List[Dict[str, Any]]] = None,
        profile: bool = True,
        timeout: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Run a data quality check on a table.

        Args:
            table_name: Name of the table to check
            rules: List of quality rules to apply
            profile: Whether to generate a profile
            timeout: Timeout in seconds for the quality check

        Returns:
            Quality check results
        """
        data = {
            "table_name": table_name,
            "profile": profile,
        }
        
        if rules:
            data["rules"] = rules
            
        if timeout:
            data["timeout"] = timeout
            
        return self._do_api_call(
            "quality/check",
            method="POST",
            data=data,
        )
    
    def get_quality_results(
        self,
        check_id: str,
    ) -> Dict[str, Any]:
        """
        Get results of a data quality check.

        Args:
            check_id: ID of the quality check

        Returns:
            Quality check results
        """
        return self._do_api_call(f"quality/results/{check_id}")
    
    def wait_for_quality_results(
        self,
        check_id: str,
        timeout: Optional[int] = None,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Wait for quality check results.

        Args:
            check_id: ID of the quality check
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds

        Returns:
            Quality check results
        """
        start_time = time.time()
        max_time = None if timeout is None else start_time + timeout
        
        while True:
            # Check timeout
            if max_time and time.time() > max_time:
                raise TimeoutError(
                    f"Timeout reached while waiting for quality check results: {check_id}"
                )
            
            # Get results
            results = self.get_quality_results(check_id)
            status = results.get('status')
            
            if status == 'completed':
                self.logger.info("Quality check completed: %s", check_id)
                return results
            elif status == 'failed':
                raise RuntimeError(f"Quality check failed: {results.get('error')}")
            elif status == 'running':
                self.logger.info("Quality check still running, waiting...")
                time.sleep(poll_interval)
            else:
                raise RuntimeError(f"Unknown quality check status: {status}")
    
    # Profile methods
    
    def run_profile(
        self,
        table_name: str,
        columns: Optional[List[str]] = None,
        sample_size: Optional[int] = None,
        timeout: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Run a data profile on a table.

        Args:
            table_name: Name of the table to profile
            columns: List of columns to profile (all if not specified)
            sample_size: Number of rows to sample
            timeout: Timeout in seconds for the profiling

        Returns:
            Profile results
        """
        data = {
            "table_name": table_name,
        }
        
        if columns:
            data["columns"] = columns
            
        if sample_size:
            data["sample_size"] = sample_size
            
        if timeout:
            data["timeout"] = timeout
            
        return self._do_api_call(
            "profile/run",
            method="POST",
            data=data,
        )
    
    def get_profile(
        self,
        profile_id: str,
    ) -> Dict[str, Any]:
        """
        Get a data profile.

        Args:
            profile_id: ID of the profile

        Returns:
            Profile data
        """
        return self._do_api_call(f"profile/{profile_id}")
    
    def wait_for_profile(
        self,
        profile_id: str,
        timeout: Optional[int] = None,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Wait for profile results.

        Args:
            profile_id: ID of the profile
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds

        Returns:
            Profile results
        """
        start_time = time.time()
        max_time = None if timeout is None else start_time + timeout
        
        while True:
            # Check timeout
            if max_time and time.time() > max_time:
                raise TimeoutError(
                    f"Timeout reached while waiting for profile results: {profile_id}"
                )
            
            # Get profile
            profile = self.get_profile(profile_id)
            status = profile.get('status')
            
            if status == 'completed':
                self.logger.info("Profiling completed: %s", profile_id)
                return profile
            elif status == 'failed':
                raise RuntimeError(f"Profiling failed: {profile.get('error')}")
            elif status == 'running':
                self.logger.info("Profiling still running, waiting...")
                time.sleep(poll_interval)
            else:
                raise RuntimeError(f"Unknown profiling status: {status}")
    
    # Validation methods
    
    def run_validation(
        self,
        table_name: str,
        rules: List[Dict[str, Any]],
        timeout: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Run a validation on a table.

        Args:
            table_name: Name of the table to validate
            rules: List of validation rules
            timeout: Timeout in seconds for the validation

        Returns:
            Validation results
        """
        data = {
            "table_name": table_name,
            "rules": rules,
        }
            
        if timeout:
            data["timeout"] = timeout
            
        return self._do_api_call(
            "validation/run",
            method="POST",
            data=data,
        )
    
    def get_validation_results(
        self,
        validation_id: str,
    ) -> Dict[str, Any]:
        """
        Get results of a validation.

        Args:
            validation_id: ID of the validation

        Returns:
            Validation results
        """
        return self._do_api_call(f"validation/results/{validation_id}")
    
    def wait_for_validation_results(
        self,
        validation_id: str,
        timeout: Optional[int] = None,
        poll_interval: int = 10,
    ) -> Dict[str, Any]:
        """
        Wait for validation results.

        Args:
            validation_id: ID of the validation
            timeout: Maximum time to wait in seconds
            poll_interval: Time between polls in seconds

        Returns:
            Validation results
        """
        start_time = time.time()
        max_time = None if timeout is None else start_time + timeout
        
        while True:
            # Check timeout
            if max_time and time.time() > max_time:
                raise TimeoutError(
                    f"Timeout reached while waiting for validation results: {validation_id}"
                )
            
            # Get results
            results = self.get_validation_results(validation_id)
            status = results.get('status')
            
            if status == 'completed':
                self.logger.info("Validation completed: %s", validation_id)
                return results
            elif status == 'failed':
                raise RuntimeError(f"Validation failed: {results.get('error')}")
            elif status == 'running':
                self.logger.info("Validation still running, waiting...")
                time.sleep(poll_interval)
            else:
                raise RuntimeError(f"Unknown validation status: {status}")
    
    # Lineage methods
    
    def get_lineage(
        self,
        table_name: str,
        max_depth: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Get lineage for a table.

        Args:
            table_name: Name of the table
            max_depth: Maximum depth of lineage graph

        Returns:
            Lineage graph
        """
        params = {}
        if max_depth:
            params["max_depth"] = max_depth
            
        return self._do_api_call(
            f"lineage/table/{table_name}",
            params=params,
        )
    
    def visualize_lineage(
        self,
        table_name: str,
        format: str = "html",
        max_depth: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Get a visualization of lineage for a table.

        Args:
            table_name: Name of the table
            format: Visualization format (html, json, svg, dot)
            max_depth: Maximum depth of lineage graph

        Returns:
            Lineage visualization
        """
        params = {
            "format": format,
        }
        
        if max_depth:
            params["max_depth"] = max_depth
            
        return self._do_api_call(
            f"lineage/visualize/{table_name}",
            params=params,
        )


class NessiResource(ConfigResource):
    """
    Dagster resource for Nessi.dev.

    This resource provides a client for interacting with the Nessi.dev API.
    """

    @property
    def client(self) -> NessiClient:
        """
        Get the Nessi client.

        Returns:
            NessiClient instance
        """
        return self._client

    def setup_for_execution(self, context) -> None:
        """
        Set up the resource for execution.

        Args:
            context: Dagster resource context
        """
        self._client = NessiClient(
            api_host=self.api_host,
            api_key=self.api_key,
            api_secret=self.api_secret,
            timeout=self.timeout,
        )

    @property
    def api_host(self) -> str:
        """
        Get the API host.

        Returns:
            API host URL
        """
        return self.config["api_host"]

    @property
    def api_key(self) -> str:
        """
        Get the API key.

        Returns:
            API key
        """
        return self.config["api_key"]

    @property
    def api_secret(self) -> Optional[str]:
        """
        Get the API secret.

        Returns:
            API secret or None
        """
        return self.config.get("api_secret")

    @property
    def timeout(self) -> int:
        """
        Get the timeout.

        Returns:
            Timeout in seconds
        """
        return self.config.get("timeout", 60)


@resource(
    config_schema={
        "api_host": str,
        "api_key": str,
        "api_secret": str,
        "timeout": int,
    },
    description="A resource for interacting with the Nessi.dev API",
)
def nessi_resource(context):
    """
    Create a Nessi resource.

    Args:
        context: Dagster resource context

    Returns:
        NessiResource instance
    """
    return NessiResource(context.resource_config)
