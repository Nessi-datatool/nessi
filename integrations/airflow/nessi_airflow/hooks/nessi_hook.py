"""
Nessi.dev hook for Apache Airflow.

This module provides a hook for connecting to Nessi.dev API.
"""

from typing import Dict, Any, Optional, List, Union
import json

from airflow.hooks.base import BaseHook
from airflow.exceptions import AirflowException

import requests


class NessiHook(BaseHook):
    """
    Hook for Nessi.dev API.

    This hook handles authentication and provides methods for interacting
    with the Nessi.dev API.

    :param nessi_conn_id: Connection ID for Nessi.dev
    :type nessi_conn_id: str
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
                "extra": "Additional configuration in JSON format",
            },
        }

    def __init__(
        self,
        nessi_conn_id: str = default_conn_name,
        timeout: int = 60,
    ) -> None:
        super().__init__()
        self.nessi_conn_id = nessi_conn_id
        self.timeout = timeout
        self.base_url = None
        self.session = None

    def get_conn(self) -> requests.Session:
        """
        Returns a requests session for making API calls to Nessi.dev.
        """
        if self.session is not None:
            return self.session

        conn = self.get_connection(self.nessi_conn_id)
        
        # Get connection details
        self.base_url = conn.host.rstrip('/')
        api_key = conn.login
        api_secret = conn.password
        
        # Parse extra configuration
        extra_config = {}
        if conn.extra:
            try:
                extra_config = json.loads(conn.extra)
            except json.JSONDecodeError:
                self.log.warning(
                    "Failed to parse extra config for connection %s", self.nessi_conn_id
                )
        
        # Create session
        self.session = requests.Session()
        
        # Set authentication
        if api_key:
            self.session.headers.update({"X-API-Key": api_key})
            if api_secret:
                self.session.headers.update({"X-API-Secret": api_secret})
        else:
            raise AirflowException(
                f"No API Key provided for Nessi connection ID '{self.nessi_conn_id}'"
            )
        
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

    def _do_api_call(
        self,
        endpoint: str,
        method: str = "GET",
        data: Optional[Union[Dict[str, Any], List[Any]]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Performs an API call to Nessi.dev.

        :param endpoint: API endpoint (without base URL)
        :param method: HTTP method
        :param data: Request payload
        :param params: Query parameters
        :return: Response as dictionary
        """
        session = self.get_conn()
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        
        try:
            if method == "GET":
                response = session.get(
                    url,
                    params=params,
                    timeout=self.timeout,
                )
            elif method == "POST":
                response = session.post(
                    url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                )
            elif method == "PUT":
                response = session.put(
                    url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                )
            elif method == "DELETE":
                response = session.delete(
                    url,
                    json=data,
                    params=params,
                    timeout=self.timeout,
                )
            else:
                raise AirflowException(f"Unsupported HTTP method: {method}")
            
            response.raise_for_status()
            
            if response.content:
                return response.json()
            return {}
            
        except requests.exceptions.HTTPError as e:
            self.log.error("HTTP error: %s", e)
            if e.response.content:
                try:
                    error_response = e.response.json()
                    self.log.error("Error response: %s", error_response)
                except ValueError:
                    self.log.error("Error response: %s", e.response.content)
            raise AirflowException(f"Nessi API HTTP error: {e}")
        except requests.exceptions.RequestException as e:
            self.log.error("Request error: %s", e)
            raise AirflowException(f"Nessi API request error: {e}")

    # Health and status methods
    
    def get_health(self) -> Dict[str, Any]:
        """
        Get Nessi.dev API health status.
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

        :param table_name: Name of the table to check
        :param rules: List of quality rules to apply
        :param profile: Whether to generate a profile
        :param timeout: Timeout in seconds for the quality check
        :return: Quality check results
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

        :param check_id: ID of the quality check
        :return: Quality check results
        """
        return self._do_api_call(f"quality/results/{check_id}")
    
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

        :param table_name: Name of the table to profile
        :param columns: List of columns to profile (all if not specified)
        :param sample_size: Number of rows to sample
        :param timeout: Timeout in seconds for the profiling
        :return: Profile results
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

        :param profile_id: ID of the profile
        :return: Profile data
        """
        return self._do_api_call(f"profile/{profile_id}")
    
    # Validation methods
    
    def run_validation(
        self,
        table_name: str,
        rules: List[Dict[str, Any]],
        timeout: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Run a validation on a table.

        :param table_name: Name of the table to validate
        :param rules: List of validation rules
        :param timeout: Timeout in seconds for the validation
        :return: Validation results
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

        :param validation_id: ID of the validation
        :return: Validation results
        """
        return self._do_api_call(f"validation/results/{validation_id}")
    
    # Lineage methods
    
    def get_lineage(
        self,
        table_name: str,
        max_depth: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Get lineage for a table.

        :param table_name: Name of the table
        :param max_depth: Maximum depth of lineage graph
        :return: Lineage graph
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

        :param table_name: Name of the table
        :param format: Visualization format (html, json, svg, dot)
        :param max_depth: Maximum depth of lineage graph
        :return: Lineage visualization
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
