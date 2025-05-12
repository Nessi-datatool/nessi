"""
Nessi Monitoring Client

Client for interacting with the Nessi monitoring system API.
"""

import json
import logging
import requests
from datetime import datetime
from typing import Dict, List, Optional, Any, Union
from urllib.parse import urljoin

from .models import Metric, Alert, Rule, Profile, ValidationResult


class NessiClient:
    """
    Client for interacting with the Nessi monitoring system.
    
    This client provides a Python interface to the Nessi monitoring system,
    allowing you to send metrics, retrieve alerts, manage rules, and access
    data quality profiles programmatically.
    """
    
    def __init__(
        self, 
        base_url: str = "http://localhost:8080", 
        api_key: Optional[str] = None,
        username: Optional[str] = None,
        password: Optional[str] = None,
        verify_ssl: bool = True
    ):
        """
        Initialize the Nessi client.
        
        Args:
            base_url: Base URL of the Nessi monitoring API
            api_key: API key for authentication (preferred over username/password)
            username: Username for authentication
            password: Password for authentication
            verify_ssl: Whether to verify SSL certificates
        """
        self.base_url = base_url.rstrip("/") + "/"
        self.api_key = api_key
        self.username = username
        self.password = password
        self.verify_ssl = verify_ssl
        self.token = None
        self.token_expiry = None
        self.logger = logging.getLogger("nessi_client")
        
        # Set up session
        self.session = requests.Session()
        self.session.verify = verify_ssl
        
        # Authenticate if credentials are provided
        if api_key:
            self.session.headers.update({"X-API-Key": api_key})
        elif username and password:
            self._authenticate()
    
    def _authenticate(self) -> None:
        """Authenticate with the API using username and password."""
        auth_url = urljoin(self.base_url, "api/login")
        payload = {
            "username": self.username,
            "password": self.password
        }
        
        try:
            response = self.session.post(auth_url, json=payload)
            response.raise_for_status()
            data = response.json()
            
            self.token = data.get("token")
            if self.token:
                self.session.headers.update({"Authorization": f"Bearer {self.token}"})
                self.logger.info("Successfully authenticated with username and password")
            else:
                self.logger.error("Authentication failed: No token received")
                
        except requests.RequestException as e:
            self.logger.error(f"Authentication failed: {str(e)}")
            raise
    
    def _refresh_token_if_needed(self) -> None:
        """Refresh the authentication token if it's expired."""
        if not self.token or (self.token_expiry and datetime.now() > self.token_expiry):
            if self.username and self.password:
                self._authenticate()
    
    def _request(self, method: str, endpoint: str, **kwargs) -> requests.Response:
        """
        Make a request to the API.
        
        Args:
            method: HTTP method (get, post, put, delete)
            endpoint: API endpoint (without base URL)
            **kwargs: Additional arguments to pass to requests
            
        Returns:
            Response object
        """
        self._refresh_token_if_needed()
        
        url = urljoin(self.base_url, endpoint)
        try:
            response = self.session.request(method, url, **kwargs)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            self.logger.error(f"API request failed: {str(e)}")
            raise
    
    # Metrics API
    
    def send_metric(self, metric: Metric) -> Dict[str, Any]:
        """
        Send a metric to the monitoring system.
        
        Args:
            metric: Metric object to send
            
        Returns:
            Response data
        """
        payload = {
            "name": metric.name,
            "value": metric.value,
            "timestamp": metric.timestamp.isoformat(),
            "tags": metric.tags,
            "metadata": metric.metadata
        }
        
        response = self._request("post", "api/metrics", json=payload)
        return response.json()
    
    def get_metrics(
        self, 
        name: Optional[str] = None,
        start_time: Optional[datetime] = None,
        end_time: Optional[datetime] = None,
        tags: Optional[Dict[str, str]] = None,
        limit: int = 100
    ) -> List[Metric]:
        """
        Get metrics from the monitoring system.
        
        Args:
            name: Filter by metric name
            start_time: Filter by start time
            end_time: Filter by end time
            tags: Filter by tags
            limit: Maximum number of metrics to return
            
        Returns:
            List of Metric objects
        """
        params = {"limit": limit}
        
        if name:
            params["name"] = name
        if start_time:
            params["start"] = start_time.isoformat()
        if end_time:
            params["end"] = end_time.isoformat()
        if tags:
            params["tags"] = json.dumps(tags)
        
        response = self._request("get", "api/metrics", params=params)
        data = response.json()
        
        metrics = []
        for item in data:
            timestamp = datetime.fromisoformat(item["timestamp"]) if "timestamp" in item else datetime.now()
            metrics.append(Metric(
                name=item["name"],
                value=item["value"],
                timestamp=timestamp,
                tags=item.get("tags", {}),
                metadata=item.get("metadata", {})
            ))
        
        return metrics
    
    # Alerts API
    
    def get_alerts(
        self,
        severity: Optional[str] = None,
        triggered: Optional[bool] = None,
        start_time: Optional[datetime] = None,
        end_time: Optional[datetime] = None,
        limit: int = 100
    ) -> List[Alert]:
        """
        Get alerts from the monitoring system.
        
        Args:
            severity: Filter by severity
            triggered: Filter by triggered status
            start_time: Filter by start time
            end_time: Filter by end time
            limit: Maximum number of alerts to return
            
        Returns:
            List of Alert objects
        """
        params = {"limit": limit}
        
        if severity:
            params["severity"] = severity
        if triggered is not None:
            params["triggered"] = str(triggered).lower()
        if start_time:
            params["start"] = start_time.isoformat()
        if end_time:
            params["end"] = end_time.isoformat()
        
        response = self._request("get", "api/alerts", params=params)
        data = response.json()
        
        alerts = []
        for item in data:
            timestamp = datetime.fromisoformat(item["timestamp"]) if "timestamp" in item else datetime.now()
            alerts.append(Alert(
                id=item["id"],
                name=item["name"],
                description=item["description"],
                severity=item["severity"],
                triggered=item["triggered"],
                timestamp=timestamp,
                metric_name=item.get("metric_name"),
                metric_value=item.get("metric_value"),
                threshold=item.get("threshold"),
                tags=item.get("tags", {}),
                metadata=item.get("metadata", {})
            ))
        
        return alerts
    
    # Rules API
    
    def get_rules(self) -> List[Rule]:
        """
        Get data quality rules from the monitoring system.
        
        Returns:
            List of Rule objects
        """
        response = self._request("get", "api/rules")
        data = response.json()
        
        rules = []
        for item in data:
            rules.append(Rule(
                id=item["id"],
                name=item["name"],
                description=item["description"],
                severity=item["severity"],
                field=item.get("field"),
                rule_type=item.get("rule_type", "custom"),
                config=item.get("config", {}),
                tags=item.get("tags", []),
                enabled=item.get("enabled", True)
            ))
        
        return rules
    
    def create_rule(self, rule: Rule) -> Dict[str, Any]:
        """
        Create a new data quality rule.
        
        Args:
            rule: Rule object to create
            
        Returns:
            Response data
        """
        payload = {
            "id": rule.id,
            "name": rule.name,
            "description": rule.description,
            "severity": rule.severity,
            "field": rule.field,
            "rule_type": rule.rule_type,
            "config": rule.config,
            "tags": rule.tags,
            "enabled": rule.enabled
        }
        
        response = self._request("post", "api/rules", json=payload)
        return response.json()
    
    def update_rule(self, rule_id: str, updates: Dict[str, Any]) -> Dict[str, Any]:
        """
        Update an existing data quality rule.
        
        Args:
            rule_id: ID of the rule to update
            updates: Dictionary of fields to update
            
        Returns:
            Response data
        """
        response = self._request("put", f"api/rules/{rule_id}", json=updates)
        return response.json()
    
    def delete_rule(self, rule_id: str) -> Dict[str, Any]:
        """
        Delete a data quality rule.
        
        Args:
            rule_id: ID of the rule to delete
            
        Returns:
            Response data
        """
        response = self._request("delete", f"api/rules/{rule_id}")
        return response.json()
    
    # Profiles API
    
    def get_profiles(
        self,
        dataset_name: Optional[str] = None,
        start_time: Optional[datetime] = None,
        end_time: Optional[datetime] = None,
        limit: int = 10
    ) -> List[Profile]:
        """
        Get data profiles from the monitoring system.
        
        Args:
            dataset_name: Filter by dataset name
            start_time: Filter by start time
            end_time: Filter by end time
            limit: Maximum number of profiles to return
            
        Returns:
            List of Profile objects
        """
        params = {"limit": limit}
        
        if dataset_name:
            params["dataset"] = dataset_name
        if start_time:
            params["start"] = start_time.isoformat()
        if end_time:
            params["end"] = end_time.isoformat()
        
        response = self._request("get", "api/profiles", params=params)
        data = response.json()
        
        profiles = []
        for item in data:
            timestamp = datetime.fromisoformat(item["timestamp"]) if "timestamp" in item else datetime.now()
            profiles.append(Profile(
                id=item["id"],
                name=item["name"],
                dataset_name=item["dataset_name"],
                timestamp=timestamp,
                row_count=item.get("row_count", 0),
                column_count=item.get("column_count", 0),
                column_stats=item.get("column_stats", {}),
                outliers=item.get("outliers", {}),
                metadata=item.get("metadata", {}),
                detection_method=item.get("detection_method", "zscore")
            ))
        
        return profiles
    
    def get_profile(self, profile_id: str) -> Profile:
        """
        Get a specific data profile.
        
        Args:
            profile_id: ID of the profile to retrieve
            
        Returns:
            Profile object
        """
        response = self._request("get", f"api/profiles/{profile_id}")
        item = response.json()
        
        timestamp = datetime.fromisoformat(item["timestamp"]) if "timestamp" in item else datetime.now()
        return Profile(
            id=item["id"],
            name=item["name"],
            dataset_name=item["dataset_name"],
            timestamp=timestamp,
            row_count=item.get("row_count", 0),
            column_count=item.get("column_count", 0),
            column_stats=item.get("column_stats", {}),
            outliers=item.get("outliers", {}),
            metadata=item.get("metadata", {}),
            detection_method=item.get("detection_method", "zscore")
        )
    
    # Validation API
    
    def validate_data(
        self, 
        data: Union[Dict[str, Any], List[Dict[str, Any]]],
        rules: Optional[List[str]] = None
    ) -> List[ValidationResult]:
        """
        Validate data against quality rules.
        
        Args:
            data: Data to validate (single record or list of records)
            rules: List of rule IDs to validate against (optional)
            
        Returns:
            List of ValidationResult objects
        """
        payload = {
            "data": data if isinstance(data, list) else [data]
        }
        
        if rules:
            payload["rules"] = rules
        
        response = self._request("post", "api/validate", json=payload)
        data = response.json()
        
        results = []
        for item in data:
            timestamp = datetime.fromisoformat(item["timestamp"]) if "timestamp" in item else datetime.now()
            results.append(ValidationResult(
                rule_id=item["rule_id"],
                field=item.get("field"),
                value=item.get("value"),
                message=item["message"],
                row_index=item.get("row_index"),
                passed=item.get("passed", False),
                timestamp=timestamp
            ))
        
        return results
