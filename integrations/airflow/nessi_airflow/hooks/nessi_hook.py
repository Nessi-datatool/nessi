"""
Nessi.dev hook for Apache Airflow.

This module provides a hook for connecting to Nessi.dev API.
"""

from typing import Dict, Any, Optional, List, Union, Tuple
import json
import logging
import os
import uuid

from airflow.hooks.base import BaseHook
from airflow.exceptions import AirflowException

import requests

from nessi_airflow.utils.error_handling import (
    NessiApiError,
    with_exponential_backoff,
    handle_api_exceptions,
)
from nessi_airflow.utils.auth import (
    NessiAuthBase,
    ApiKeyAuth,
    OAuth2Auth,
    get_auth_from_connection,
    get_auth_from_env,
)
from nessi_airflow.utils.observability import (
    MetricsCollector,
    StructuredLogger,
    with_correlation_id,
    trace_api_call,
)

# Create a structured logger for this module
logger = StructuredLogger(__name__)


class NessiHook(BaseHook):
    """
    Hook for Nessi.dev API.

    This hook handles authentication and provides methods for interacting
    with the Nessi.dev API.

    :param nessi_conn_id: Connection ID for Nessi.dev, set to None to use environment variables
    :type nessi_conn_id: Optional[str]
    :param timeout: Timeout for API requests in seconds
    :type timeout: int
    """

    conn_name_attr = 'nessi_conn_id'
    default_conn_name = 'nessi_default'
    conn_type = 'nessi'
    hook_name = 'Nessi.dev'

    @staticmethod
    def get_ui_field_behavior() -> Dict:
        """
        Returns custom field behavior for the Airflow connection UI.
        """
        return {
            "hidden_fields": ["port", "schema"],
            "relabeling": {
                "host": "API Host",
                "login": "API Key",
                "password": "API Secret",
            },
            "placeholders": {
                "host": "https://api.nessi.dev",
                "login": "Your API Key",
                "password": "Your API Secret (leave empty if using API Key only)",
                "extra": """
{
  "auth_type": "api_key|oauth2",
  "client_id": "your-oauth-client-id",
  "client_secret": "your-oauth-client-secret",
  "token_url": "https://api.nessi.dev/oauth/token"
}
                """,
            },
        }

    def __init__(
        self,
        nessi_conn_id: Optional[str] = default_conn_name,
        timeout: int = 60,
    ) -> None:
        super().__init__()
        self.nessi_conn_id = nessi_conn_id
        self.timeout = timeout
        self.base_url = None
        self.session = None
        self.auth_handler = None

    def get_conn(self) -> requests.Session:
        """
        Returns a requests session for making API calls to Nessi.dev.
        """
        if self.session is not None:
            return self.session

        # Create session
        self.session = requests.Session()
        
        # Set up authentication
        if self.nessi_conn_id:
            # Use connection-based authentication
            logger.info("Using connection-based authentication: %s", self.nessi_conn_id)
            self.auth_handler = get_auth_from_connection(self.nessi_conn_id)
            self.base_url = self.auth_handler.base_url
        else:
            # Try environment variable-based authentication
            logger.info("Trying environment variable-based authentication")
            self.auth_handler, self.base_url = get_auth_from_env()
            
            if not self.auth_handler or not self.base_url:
                raise AirflowException(
                    "No authentication configuration found in connection or environment variables"
                )
        
        # Apply authentication to session
        self.auth_handler.update_session(self.session)
        
        # Set additional headers
        self.session.headers.update({
            "Content-Type": "application/json",
            "Accept": "application/json",
            "User-Agent": f"nessi-airflow/{self.__module__}",
        })
        
        return self.session

    def test_connection(self) -> tuple[bool, str]:
        """
        Test the Nessi.dev connection by making a request to the health endpoint.
        """
        try:
            response = self.get_health()
            if response.get("status") == "ok":
                return True, "Connection successful"
            return False, f"Connection failed: {response}"
        except Exception as e:
            return False, str(e)

    @trace_api_call("api_call")
    @handle_api_exceptions
    @with_exponential_backoff(
        max_retries=3,
        initial_backoff=1.0,
        max_backoff=60.0,
        backoff_factor=2.0,
    )
    def _do_api_call(
        self,
        endpoint: str,
        method: str = "GET",
        data: Optional[Union[Dict[str, Any], List[Any]]] = None,
        params: Optional[Dict[str, Any]] = None,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Performs an API call to Nessi.dev.

        :param endpoint: API endpoint (without base URL)
        :param method: HTTP method
        :param data: Request payload
        :param params: Query parameters
        :param correlation_id: Correlation ID for tracing
        :return: Response as dictionary
        """
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        request_id = str(uuid.uuid4())
        
        # Add correlation ID to request headers if provided
        headers = {}
        if correlation_id:
            headers["X-Correlation-ID"] = correlation_id
        headers["X-Request-ID"] = request_id
        
        # Log the API request with structured logging
        logger.info(
            f"API request: {method} {endpoint}",
            method=method,
            url=url,
            endpoint=endpoint,
            params=params,
            data_size=len(json.dumps(data)) if data else 0,
            request_id=request_id,
            correlation_id=correlation_id,
        )
        
        # Record metrics for the API request
        MetricsCollector.count(f"api.{endpoint}.request")
        
        # Make the API request
        if method == "GET":
            response = self.get_conn().get(
                url,
                params=params,
                headers=headers,
                timeout=self.timeout,
            )
        elif method == "POST":
            response = self.get_conn().post(
                url,
                json=data,
                params=params,
                headers=headers,
                timeout=self.timeout,
            )
        elif method == "PUT":
            response = self.get_conn().put(
                url,
                json=data,
                params=params,
                headers=headers,
                timeout=self.timeout,
            )
        elif method == "DELETE":
            response = self.get_conn().delete(
                url,
                params=params,
                headers=headers,
                timeout=self.timeout,
            )
        else:
            raise AirflowException(f"Unsupported HTTP method: {method}")
        
        # Log the API response
        logger.info(
            f"API response: {response.status_code} {method} {endpoint}",
            method=method,
            url=url,
            endpoint=endpoint,
            status_code=response.status_code,
            response_size=len(response.content),
            request_id=request_id,
            correlation_id=correlation_id,
        )
        
        # Record metrics for the API response
        MetricsCollector.count(f"api.{endpoint}.response.{response.status_code}")
        
        response.raise_for_status()
        return response.json()

    # Health and status methods
    
    @MetricsCollector.hook_timer("get_health")
    @with_correlation_id
    def get_health(self, correlation_id: Optional[str] = None) -> Dict[str, Any]:
        """
        Get Nessi.dev API health status.
        
        :param correlation_id: Correlation ID for tracing
        :return: Health status information
        """
        return self._do_api_call("health", correlation_id=correlation_id)
    
    # Quality methods
    
    @MetricsCollector.hook_timer("run_quality_check")
    @with_correlation_id
    def run_quality_check(
        self,
        table_name: str,
        rules: Optional[List[Dict[str, Any]]] = None,
        profile: bool = True,
        timeout: Optional[int] = None,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Run a data quality check on a table.

        :param table_name: Name of the table to check
        :param rules: List of quality rules to apply
        :param profile: Whether to generate a profile
        :param timeout: Timeout in seconds for the quality check
        :param correlation_id: Correlation ID for tracing
        :return: Quality check results
        """
        # Log the quality check request
        logger.info(
            f"Running quality check on table: {table_name}",
            table_name=table_name,
            rule_count=len(rules) if rules else 0,
            profile=profile,
            timeout=timeout,
            correlation_id=correlation_id,
        )
        
        # Record metrics for the quality check
        MetricsCollector.count("quality.check.request")
        
        data = {
            "table_name": table_name,
            "profile": profile,
        }
        
        if rules:
            data["rules"] = rules
            
        if timeout:
            data["timeout"] = timeout
            
        result = self._do_api_call(
            "quality/check",
            method="POST",
            data=data,
            correlation_id=correlation_id,
        )
        
        # Log the quality check result
        logger.info(
            f"Quality check initiated: {result.get('check_id')}",
            table_name=table_name,
            check_id=result.get('check_id'),
            correlation_id=correlation_id,
        )
        
        return result
    
    @MetricsCollector.hook_timer("get_quality_results")
    @with_correlation_id
    def get_quality_results(
        self,
        check_id: str,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Get results of a data quality check.

        :param check_id: ID of the quality check
        :param correlation_id: Correlation ID for tracing
        :return: Quality check results
        """
        # Log the quality results request
        logger.info(
            f"Getting quality check results: {check_id}",
            check_id=check_id,
            correlation_id=correlation_id,
        )
        
        result = self._do_api_call(
            f"quality/results/{check_id}",
            correlation_id=correlation_id,
        )
        
        # Log the quality results status
        logger.info(
            f"Quality check status: {result.get('status')}",
            check_id=check_id,
            status=result.get('status'),
            correlation_id=correlation_id,
        )
        
        # Record metrics for the quality results
        MetricsCollector.count(f"quality.results.status.{result.get('status', 'unknown')}")
        if 'quality_score' in result:
            MetricsCollector.gauge("quality.score", result['quality_score'])
        
        return result
    
    # Profile methods
    
    @MetricsCollector.hook_timer("run_profile")
    @with_correlation_id
    def run_profile(
        self,
        table_name: str,
        columns: Optional[List[str]] = None,
        sample_size: Optional[int] = None,
        timeout: Optional[int] = None,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Run a data profile on a table.

        :param table_name: Name of the table to profile
        :param columns: List of columns to profile (all if not specified)
        :param sample_size: Number of rows to sample
        :param timeout: Timeout in seconds for the profiling
        :param correlation_id: Correlation ID for tracing
        :return: Profile results
        """
        # Log the profile request
        logger.info(
            f"Running profile on table: {table_name}",
            table_name=table_name,
            column_count=len(columns) if columns else 0,
            sample_size=sample_size,
            timeout=timeout,
            correlation_id=correlation_id,
        )
        
        # Record metrics for the profile
        MetricsCollector.count("profile.run.request")
        
        data = {
            "table_name": table_name,
        }
        
        if columns:
            data["columns"] = columns
            
        if sample_size:
            data["sample_size"] = sample_size
            
        if timeout:
            data["timeout"] = timeout
            
        result = self._do_api_call(
            "profile/run",
            method="POST",
            data=data,
            correlation_id=correlation_id,
        )
        
        # Log the profile result
        logger.info(
            f"Profile initiated: {result.get('profile_id')}",
            table_name=table_name,
            profile_id=result.get('profile_id'),
            correlation_id=correlation_id,
        )
        
        return result
    
    @MetricsCollector.hook_timer("get_profile")
    @with_correlation_id
    def get_profile(
        self,
        profile_id: str,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Get a data profile.

        :param profile_id: ID of the profile
        :param correlation_id: Correlation ID for tracing
        :return: Profile data
        """
        # Log the profile request
        logger.info(
            f"Getting profile: {profile_id}",
            profile_id=profile_id,
            correlation_id=correlation_id,
        )
        
        result = self._do_api_call(
            f"profile/{profile_id}",
            correlation_id=correlation_id,
        )
        
        # Log the profile result
        logger.info(
            f"Profile retrieved: {profile_id}",
            profile_id=profile_id,
            correlation_id=correlation_id,
        )
        
        return result
    
    # Validation methods
    
    @MetricsCollector.hook_timer("run_validation")
    @with_correlation_id
    def run_validation(
        self,
        table_name: str,
        rules: List[Dict[str, Any]],
        timeout: Optional[int] = None,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Run a validation on a table.

        :param table_name: Name of the table to validate
        :param rules: List of validation rules
        :param timeout: Timeout in seconds for the validation
        :param correlation_id: Correlation ID for tracing
        :return: Validation results
        """
        # Log the validation request
        logger.info(
            f"Running validation on table: {table_name}",
            table_name=table_name,
            rule_count=len(rules),
            timeout=timeout,
            correlation_id=correlation_id,
        )
        
        # Record metrics for the validation
        MetricsCollector.count("validation.run.request")
        
        data = {
            "table_name": table_name,
            "rules": rules,
        }
            
        if timeout:
            data["timeout"] = timeout
            
        result = self._do_api_call(
            "validation/run",
            method="POST",
            data=data,
            correlation_id=correlation_id,
        )
        
        # Log the validation result
        logger.info(
            f"Validation initiated: {result.get('validation_id')}",
            table_name=table_name,
            validation_id=result.get('validation_id'),
            correlation_id=correlation_id,
        )
        
        return result
    
    @MetricsCollector.hook_timer("get_validation_results")
    @with_correlation_id
    def get_validation_results(
        self,
        validation_id: str,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Get results of a validation.

        :param validation_id: ID of the validation
        :param correlation_id: Correlation ID for tracing
        :return: Validation results
        """
        # Log the validation results request
        logger.info(
            f"Getting validation results: {validation_id}",
            validation_id=validation_id,
            correlation_id=correlation_id,
        )
        
        result = self._do_api_call(
            f"validation/results/{validation_id}",
            correlation_id=correlation_id,
        )
        
        # Log the validation results status
        logger.info(
            f"Validation status: {result.get('status')}",
            validation_id=validation_id,
            status=result.get('status'),
            correlation_id=correlation_id,
        )
        
        # Record metrics for the validation results
        MetricsCollector.count(f"validation.results.status.{result.get('status', 'unknown')}")
        
        return result
    
    # Lineage methods
    
    @MetricsCollector.hook_timer("get_lineage")
    @with_correlation_id
    def get_lineage(
        self,
        table_name: str,
        max_depth: Optional[int] = None,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Get lineage for a table.

        :param table_name: Name of the table
        :param max_depth: Maximum depth of lineage graph
        :param correlation_id: Correlation ID for tracing
        :return: Lineage graph
        """
        # Log the lineage request
        logger.info(
            f"Getting lineage for table: {table_name}",
            table_name=table_name,
            max_depth=max_depth,
            correlation_id=correlation_id,
        )
        
        params = {}
        if max_depth:
            params["max_depth"] = max_depth
            
        result = self._do_api_call(
            f"lineage/table/{table_name}",
            params=params,
            correlation_id=correlation_id,
        )
        
        # Log the lineage result
        logger.info(
            f"Lineage retrieved for table: {table_name}",
            table_name=table_name,
            node_count=len(result.get("nodes", [])),
            edge_count=len(result.get("edges", [])),
            correlation_id=correlation_id,
        )
        
        return result
    
    @MetricsCollector.hook_timer("visualize_lineage")
    @with_correlation_id
    def visualize_lineage(
        self,
        table_name: str,
        format: str = "html",
        max_depth: Optional[int] = None,
        correlation_id: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Get a visualization of lineage for a table.

        :param table_name: Name of the table
        :param format: Visualization format (html, json, svg, dot)
        :param max_depth: Maximum depth of lineage graph
        :param correlation_id: Correlation ID for tracing
        :return: Lineage visualization
        """
        # Log the lineage visualization request
        logger.info(
            f"Visualizing lineage for table: {table_name}",
            table_name=table_name,
            format=format,
            max_depth=max_depth,
            correlation_id=correlation_id,
        )
        
        params = {
            "format": format,
        }
        
        if max_depth:
            params["max_depth"] = max_depth
            
        result = self._do_api_call(
            f"lineage/visualize/{table_name}",
            params=params,
            correlation_id=correlation_id,
        )
        
        # Log the lineage visualization result
        logger.info(
            f"Lineage visualization created for table: {table_name}",
            table_name=table_name,
            format=format,
            correlation_id=correlation_id,
        )
        
        return result