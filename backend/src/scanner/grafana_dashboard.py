"""
Grafana Dashboard Integration Module

This module provides functionality to create and manage Grafana dashboards
for monitoring Nessi's data processing operations.
"""

import json
import os
import time
from datetime import datetime
from typing import Dict, Any, List, Optional
import requests
from pyspark.sql import DataFrame

class GrafanaDashboard:
    """Class for managing Grafana dashboards and metrics."""
    
    def __init__(self, grafana_url: str, password: str, username: str = "admin"):
        """
        Initialize the Grafana dashboard manager.
        
        Args:
            grafana_url (str): URL of the Grafana instance
            password (str): Password for authentication
            username (str): Username for authentication (default: admin)
        """
        self.grafana_url = grafana_url.rstrip('/')
        self.auth = (username, password)
        self.headers = {
            'Content-Type': 'application/json'
        }
    
    def create_dashboard(self, title: str, folder_id: int = 0) -> Dict[str, Any]:
        """
        Create a new Grafana dashboard.
        
        Args:
            title (str): Title of the dashboard
            folder_id (int): ID of the folder to store the dashboard
            
        Returns:
            Dict[str, Any]: Dashboard creation response
        """
        dashboard = {
            "dashboard": {
                "id": None,
                "uid": None,
                "title": title,
                "tags": ["nessi", "data-processing"],
                "timezone": "browser",
                "schemaVersion": 16,
                "version": 0,
                "refresh": "5s",
                "panels": []
            },
            "folderId": folder_id,
            "overwrite": False
        }
        
        try:
            response = requests.post(
                f"{self.grafana_url}/api/dashboards/db",
                headers=self.headers,
                auth=self.auth,
                json=dashboard
            )
            response.raise_for_status()
            result = response.json()
            return {
                "dashboard": {
                    "uid": result.get("dashboard", {}).get("uid"),
                    "id": result.get("dashboard", {}).get("id"),
                    "title": title,
                    "version": 1,
                    "panels": []
                }
            }
        except requests.exceptions.RequestException as e:
            raise Exception(f"Failed to create dashboard: {str(e)}")
    
    def add_table_metrics_panel(self, dashboard_uid: str, scan_results: Dict[str, Any]) -> Dict[str, Any]:
        """
        Add a panel for table metrics to the dashboard.
        
        Args:
            dashboard_uid (str): UID of the dashboard
            scan_results (Dict[str, Any]): Results from table scanning
            
        Returns:
            Dict[str, Any]: Panel creation response
        """
        try:
            # Create metrics panel
            panel = {
                "title": "Table Metrics",
                "type": "table",
                "gridPos": {
                    "h": 8,
                    "w": 12,
                    "x": 0,
                    "y": 0
                },
                "options": {
                    "showHeader": True
                },
                "fieldConfig": {
                    "defaults": {
                        "custom": {
                            "align": "auto",
                            "displayMode": "auto"
                        }
                    }
                },
                "targets": [
                    {
                        "refId": "A",
                        "datasource": "Prometheus",
                        "expr": "nessi_table_metrics",
                        "format": "table",
                        "instant": True
                    }
                ],
                "transformations": [
                    {
                        "id": "organize",
                        "options": {
                            "excludeByName": {},
                            "indexByName": {},
                            "renameByName": {
                                "Value": "Count",
                                "__name__": "Metric"
                            }
                        }
                    }
                ]
            }
            
            # Add panel to dashboard
            response = requests.get(
                f"{self.grafana_url}/api/dashboards/uid/{dashboard_uid}",
                headers=self.headers,
                auth=self.auth
            )
            response.raise_for_status()
            dashboard = response.json()['dashboard']
            
            # Add panel to dashboard
            dashboard['panels'].append(panel)
            
            # Update dashboard
            update_response = requests.post(
                f"{self.grafana_url}/api/dashboards/db",
                headers=self.headers,
                auth=self.auth,
                json={"dashboard": dashboard, "overwrite": True}
            )
            update_response.raise_for_status()
            
            return response.json()
        except requests.exceptions.RequestException as e:
            raise Exception(f"Failed to add table metrics panel: {str(e)}")
    
    def add_column_stats_panel(self, dashboard_uid: str, column_stats: Dict[str, Any]) -> Dict[str, Any]:
        """
        Add a panel for column statistics to the dashboard.
        
        Args:
            dashboard_uid (str): UID of the dashboard
            column_stats (Dict[str, Any]): Column statistics from table scanning
            
        Returns:
            Dict[str, Any]: Panel creation response
        """
        try:
            # Create statistics panel
            panel = {
                "title": "Column Statistics",
                "type": "stat",
                "gridPos": {
                    "h": 8,
                    "w": 12,
                    "x": 12,
                    "y": 0
                },
                "options": {
                    "colorMode": "value",
                    "graphMode": "area",
                    "justifyMode": "auto",
                    "orientation": "auto",
                    "reduceOptions": {
                        "calcs": ["mean"],
                        "fields": "",
                        "values": False
                    },
                    "textMode": "auto"
                },
                "fieldConfig": {
                    "defaults": {
                        "mappings": [],
                        "thresholds": {
                            "mode": "absolute",
                            "steps": [
                                {"color": "green", "value": None},
                                {"color": "red", "value": 80}
                            ]
                        }
                    }
                },
                "targets": [
                    {
                        "refId": "A",
                        "datasource": "Prometheus",
                        "expr": "nessi_column_stats",
                        "format": "time_series",
                        "instant": True
                    }
                ]
            }
            
            # Add panel to dashboard
            response = requests.get(
                f"{self.grafana_url}/api/dashboards/uid/{dashboard_uid}",
                headers=self.headers,
                auth=self.auth
            )
            response.raise_for_status()
            dashboard = response.json()['dashboard']
            
            # Add panel to dashboard
            dashboard['panels'].append(panel)
            
            # Update dashboard
            update_response = requests.post(
                f"{self.grafana_url}/api/dashboards/db",
                headers=self.headers,
                auth=self.auth,
                json={"dashboard": dashboard, "overwrite": True}
            )
            update_response.raise_for_status()
            
            return response.json()
        except requests.exceptions.RequestException as e:
            raise Exception(f"Failed to add column statistics panel: {str(e)}")
    
    def export_dashboard(self, dashboard_uid: str, output_path: str) -> str:
        """
        Export a dashboard to a file.
        
        Args:
            dashboard_uid (str): UID of the dashboard
            output_path (str): Path to save the dashboard file
            
        Returns:
            str: Path to the exported dashboard file
        """
        try:
            response = requests.get(
                f"{self.grafana_url}/api/dashboards/uid/{dashboard_uid}",
                headers=self.headers,
                auth=self.auth
            )
            response.raise_for_status()
            dashboard = response.json()['dashboard']
            
            # Save dashboard to file
            with open(output_path, 'w') as f:
                json.dump(dashboard, f, indent=2)
            
            return output_path
        except requests.exceptions.RequestException as e:
            raise Exception(f"Failed to export dashboard: {str(e)}")
        except IOError as e:
            raise Exception(f"Failed to write dashboard to file: {str(e)}") 