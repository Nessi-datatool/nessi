#!/usr/bin/env python3
"""
Unit tests for the Nessi client library.
"""

import unittest
from unittest.mock import patch, MagicMock
import json
from datetime import datetime
import requests

import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nessi_client import NessiClient, Metric, Alert, Rule, Profile, ValidationResult


class TestNessiClient(unittest.TestCase):
    """Test cases for the NessiClient class."""

    def setUp(self):
        """Set up test fixtures."""
        self.base_url = "http://localhost:8080"
        self.api_key = "test_api_key"
        self.username = "test_user"
        self.password = "test_password"
        
        # Create a client with API key
        self.api_client = NessiClient(
            base_url=self.base_url,
            api_key=self.api_key
        )
        
        # Create a client with username/password
        self.user_client = NessiClient(
            base_url=self.base_url,
            username=self.username,
            password=self.password
        )

    @patch('requests.Session.post')
    def test_authenticate(self, mock_post):
        """Test authentication with username and password."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = {"token": "test_token"}
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response
        
        # Create a new client to trigger authentication
        client = NessiClient(
            base_url=self.base_url,
            username=self.username,
            password=self.password
        )
        
        # Verify authentication request
        mock_post.assert_called_once_with(
            f"{self.base_url}/api/login",
            json={"username": self.username, "password": self.password}
        )
        
        # Verify token is set
        self.assertEqual(client.token, "test_token")
        self.assertIn("Authorization", client.session.headers)
        self.assertEqual(client.session.headers["Authorization"], "Bearer test_token")

    @patch('requests.Session.post')
    def test_send_metric(self, mock_post):
        """Test sending a metric."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = {"id": "metric_id"}
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response
        
        # Create a metric
        metric = Metric(
            name="test_metric",
            value=42.0,
            timestamp=datetime(2025, 5, 12, 12, 0, 0),
            tags={"tag1": "value1"},
            metadata={"meta1": "value1"}
        )
        
        # Send the metric
        result = self.api_client.send_metric(metric)
        
        # Verify request
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(kwargs["json"]["name"], "test_metric")
        self.assertEqual(kwargs["json"]["value"], 42.0)
        self.assertEqual(kwargs["json"]["timestamp"], "2025-05-12T12:00:00")
        self.assertEqual(kwargs["json"]["tags"], {"tag1": "value1"})
        self.assertEqual(kwargs["json"]["metadata"], {"meta1": "value1"})
        
        # Verify result
        self.assertEqual(result, {"id": "metric_id"})

    @patch('requests.Session.get')
    def test_get_metrics(self, mock_get):
        """Test retrieving metrics."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = [
            {
                "name": "test_metric1",
                "value": 42.0,
                "timestamp": "2025-05-12T12:00:00",
                "tags": {"tag1": "value1"},
                "metadata": {"meta1": "value1"}
            },
            {
                "name": "test_metric2",
                "value": 43.0,
                "timestamp": "2025-05-12T12:01:00",
                "tags": {"tag2": "value2"},
                "metadata": {"meta2": "value2"}
            }
        ]
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response
        
        # Get metrics
        metrics = self.api_client.get_metrics(
            name="test_metric",
            start_time=datetime(2025, 5, 12, 12, 0, 0),
            end_time=datetime(2025, 5, 12, 12, 30, 0),
            tags={"tag1": "value1"},
            limit=10
        )
        
        # Verify request
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        self.assertEqual(kwargs["params"]["name"], "test_metric")
        self.assertEqual(kwargs["params"]["start"], "2025-05-12T12:00:00")
        self.assertEqual(kwargs["params"]["end"], "2025-05-12T12:30:00")
        self.assertEqual(kwargs["params"]["tags"], json.dumps({"tag1": "value1"}))
        self.assertEqual(kwargs["params"]["limit"], 10)
        
        # Verify result
        self.assertEqual(len(metrics), 2)
        self.assertEqual(metrics[0].name, "test_metric1")
        self.assertEqual(metrics[0].value, 42.0)
        self.assertEqual(metrics[0].timestamp, datetime.fromisoformat("2025-05-12T12:00:00"))
        self.assertEqual(metrics[0].tags, {"tag1": "value1"})
        self.assertEqual(metrics[0].metadata, {"meta1": "value1"})
        
        self.assertEqual(metrics[1].name, "test_metric2")
        self.assertEqual(metrics[1].value, 43.0)
        self.assertEqual(metrics[1].timestamp, datetime.fromisoformat("2025-05-12T12:01:00"))
        self.assertEqual(metrics[1].tags, {"tag2": "value2"})
        self.assertEqual(metrics[1].metadata, {"meta2": "value2"})

    @patch('requests.Session.get')
    def test_get_alerts(self, mock_get):
        """Test retrieving alerts."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = [
            {
                "id": "alert1",
                "name": "High CPU Usage",
                "description": "CPU usage is too high",
                "severity": "warning",
                "triggered": True,
                "timestamp": "2025-05-12T12:00:00",
                "metric_name": "cpu_usage",
                "metric_value": 85.0,
                "threshold": 80.0,
                "tags": {"host": "server1"},
                "metadata": {"meta1": "value1"}
            }
        ]
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response
        
        # Get alerts
        alerts = self.api_client.get_alerts(
            severity="warning",
            triggered=True,
            start_time=datetime(2025, 5, 12, 12, 0, 0),
            end_time=datetime(2025, 5, 12, 12, 30, 0),
            limit=10
        )
        
        # Verify request
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        self.assertEqual(kwargs["params"]["severity"], "warning")
        self.assertEqual(kwargs["params"]["triggered"], "true")
        self.assertEqual(kwargs["params"]["start"], "2025-05-12T12:00:00")
        self.assertEqual(kwargs["params"]["end"], "2025-05-12T12:30:00")
        self.assertEqual(kwargs["params"]["limit"], 10)
        
        # Verify result
        self.assertEqual(len(alerts), 1)
        self.assertEqual(alerts[0].id, "alert1")
        self.assertEqual(alerts[0].name, "High CPU Usage")
        self.assertEqual(alerts[0].description, "CPU usage is too high")
        self.assertEqual(alerts[0].severity, "warning")
        self.assertEqual(alerts[0].triggered, True)
        self.assertEqual(alerts[0].timestamp, datetime.fromisoformat("2025-05-12T12:00:00"))
        self.assertEqual(alerts[0].metric_name, "cpu_usage")
        self.assertEqual(alerts[0].metric_value, 85.0)
        self.assertEqual(alerts[0].threshold, 80.0)
        self.assertEqual(alerts[0].tags, {"host": "server1"})
        self.assertEqual(alerts[0].metadata, {"meta1": "value1"})

    @patch('requests.Session.get')
    def test_get_rules(self, mock_get):
        """Test retrieving rules."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = [
            {
                "id": "rule1",
                "name": "Price Range Check",
                "description": "Validates that product prices are within an acceptable range",
                "severity": "error",
                "field": "price",
                "rule_type": "range",
                "config": {"min": 0.01, "max": 9999.99},
                "tags": ["product", "validation"],
                "enabled": True
            }
        ]
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response
        
        # Get rules
        rules = self.api_client.get_rules()
        
        # Verify request
        mock_get.assert_called_once()
        
        # Verify result
        self.assertEqual(len(rules), 1)
        self.assertEqual(rules[0].id, "rule1")
        self.assertEqual(rules[0].name, "Price Range Check")
        self.assertEqual(rules[0].description, "Validates that product prices are within an acceptable range")
        self.assertEqual(rules[0].severity, "error")
        self.assertEqual(rules[0].field, "price")
        self.assertEqual(rules[0].rule_type, "range")
        self.assertEqual(rules[0].config, {"min": 0.01, "max": 9999.99})
        self.assertEqual(rules[0].tags, ["product", "validation"])
        self.assertEqual(rules[0].enabled, True)

    @patch('requests.Session.post')
    def test_create_rule(self, mock_post):
        """Test creating a rule."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = {"id": "rule1"}
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response
        
        # Create a rule
        rule = Rule(
            id="rule1",
            name="Price Range Check",
            description="Validates that product prices are within an acceptable range",
            severity="error",
            field="price",
            rule_type="range",
            config={"min": 0.01, "max": 9999.99},
            tags=["product", "validation"],
            enabled=True
        )
        
        # Send the rule
        result = self.api_client.create_rule(rule)
        
        # Verify request
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(kwargs["json"]["id"], "rule1")
        self.assertEqual(kwargs["json"]["name"], "Price Range Check")
        self.assertEqual(kwargs["json"]["description"], "Validates that product prices are within an acceptable range")
        self.assertEqual(kwargs["json"]["severity"], "error")
        self.assertEqual(kwargs["json"]["field"], "price")
        self.assertEqual(kwargs["json"]["rule_type"], "range")
        self.assertEqual(kwargs["json"]["config"], {"min": 0.01, "max": 9999.99})
        self.assertEqual(kwargs["json"]["tags"], ["product", "validation"])
        self.assertEqual(kwargs["json"]["enabled"], True)
        
        # Verify result
        self.assertEqual(result, {"id": "rule1"})

    @patch('requests.Session.put')
    def test_update_rule(self, mock_put):
        """Test updating a rule."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = {"id": "rule1", "updated": True}
        mock_response.raise_for_status.return_value = None
        mock_put.return_value = mock_response
        
        # Update a rule
        updates = {
            "config": {"min": 0.01, "max": 5000.00},
            "severity": "warning"
        }
        
        # Send the update
        result = self.api_client.update_rule("rule1", updates)
        
        # Verify request
        mock_put.assert_called_once()
        args, kwargs = mock_put.call_args
        self.assertEqual(args[0], f"{self.base_url}/api/rules/rule1")
        self.assertEqual(kwargs["json"], updates)
        
        # Verify result
        self.assertEqual(result, {"id": "rule1", "updated": True})

    @patch('requests.Session.delete')
    def test_delete_rule(self, mock_delete):
        """Test deleting a rule."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = {"id": "rule1", "deleted": True}
        mock_response.raise_for_status.return_value = None
        mock_delete.return_value = mock_response
        
        # Delete a rule
        result = self.api_client.delete_rule("rule1")
        
        # Verify request
        mock_delete.assert_called_once()
        args, kwargs = mock_delete.call_args
        self.assertEqual(args[0], f"{self.base_url}/api/rules/rule1")
        
        # Verify result
        self.assertEqual(result, {"id": "rule1", "deleted": True})

    @patch('requests.Session.get')
    def test_get_profiles(self, mock_get):
        """Test retrieving profiles."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = [
            {
                "id": "profile1",
                "name": "Customer Data Profile",
                "dataset_name": "customers",
                "timestamp": "2025-05-12T12:00:00",
                "row_count": 1000,
                "column_count": 10,
                "column_stats": {
                    "customer_id": {
                        "type": "string",
                        "null_percentage": 0.0,
                        "unique_percentage": 100.0
                    },
                    "age": {
                        "type": "numeric",
                        "null_percentage": 2.0,
                        "min": 18,
                        "max": 85,
                        "mean": 42.5,
                        "median": 41.0
                    }
                },
                "outliers": {
                    "age": [17, 86, 90]
                },
                "metadata": {"source": "database"},
                "detection_method": "zscore"
            }
        ]
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response
        
        # Get profiles
        profiles = self.api_client.get_profiles(
            dataset_name="customers",
            start_time=datetime(2025, 5, 12, 12, 0, 0),
            end_time=datetime(2025, 5, 12, 12, 30, 0),
            limit=10
        )
        
        # Verify request
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        self.assertEqual(kwargs["params"]["dataset"], "customers")
        self.assertEqual(kwargs["params"]["start"], "2025-05-12T12:00:00")
        self.assertEqual(kwargs["params"]["end"], "2025-05-12T12:30:00")
        self.assertEqual(kwargs["params"]["limit"], 10)
        
        # Verify result
        self.assertEqual(len(profiles), 1)
        self.assertEqual(profiles[0].id, "profile1")
        self.assertEqual(profiles[0].name, "Customer Data Profile")
        self.assertEqual(profiles[0].dataset_name, "customers")
        self.assertEqual(profiles[0].timestamp, datetime.fromisoformat("2025-05-12T12:00:00"))
        self.assertEqual(profiles[0].row_count, 1000)
        self.assertEqual(profiles[0].column_count, 10)
        self.assertEqual(profiles[0].column_stats["customer_id"]["type"], "string")
        self.assertEqual(profiles[0].column_stats["age"]["mean"], 42.5)
        self.assertEqual(profiles[0].outliers["age"], [17, 86, 90])
        self.assertEqual(profiles[0].metadata, {"source": "database"})
        self.assertEqual(profiles[0].detection_method, "zscore")

    @patch('requests.Session.get')
    def test_get_profile(self, mock_get):
        """Test retrieving a specific profile."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "id": "profile1",
            "name": "Customer Data Profile",
            "dataset_name": "customers",
            "timestamp": "2025-05-12T12:00:00",
            "row_count": 1000,
            "column_count": 10,
            "column_stats": {
                "customer_id": {
                    "type": "string",
                    "null_percentage": 0.0,
                    "unique_percentage": 100.0
                },
                "age": {
                    "type": "numeric",
                    "null_percentage": 2.0,
                    "min": 18,
                    "max": 85,
                    "mean": 42.5,
                    "median": 41.0
                }
            },
            "outliers": {
                "age": [17, 86, 90]
            },
            "metadata": {"source": "database"},
            "detection_method": "zscore"
        }
        mock_response.raise_for_status.return_value = None
        mock_get.return_value = mock_response
        
        # Get a specific profile
        profile = self.api_client.get_profile("profile1")
        
        # Verify request
        mock_get.assert_called_once()
        args, kwargs = mock_get.call_args
        self.assertEqual(args[0], f"{self.base_url}/api/profiles/profile1")
        
        # Verify result
        self.assertEqual(profile.id, "profile1")
        self.assertEqual(profile.name, "Customer Data Profile")
        self.assertEqual(profile.dataset_name, "customers")
        self.assertEqual(profile.timestamp, datetime.fromisoformat("2025-05-12T12:00:00"))
        self.assertEqual(profile.row_count, 1000)
        self.assertEqual(profile.column_count, 10)
        self.assertEqual(profile.column_stats["customer_id"]["type"], "string")
        self.assertEqual(profile.column_stats["age"]["mean"], 42.5)
        self.assertEqual(profile.outliers["age"], [17, 86, 90])
        self.assertEqual(profile.metadata, {"source": "database"})
        self.assertEqual(profile.detection_method, "zscore")

    @patch('requests.Session.post')
    def test_validate_data(self, mock_post):
        """Test validating data against rules."""
        # Mock response
        mock_response = MagicMock()
        mock_response.json.return_value = [
            {
                "rule_id": "price_range",
                "field": "price",
                "value": -10.0,
                "message": "Price must be greater than 0.01",
                "row_index": 1,
                "passed": False,
                "timestamp": "2025-05-12T12:00:00"
            },
            {
                "rule_id": "price_range",
                "field": "price",
                "value": 99.99,
                "message": "Price is within acceptable range",
                "row_index": 0,
                "passed": True,
                "timestamp": "2025-05-12T12:00:00"
            }
        ]
        mock_response.raise_for_status.return_value = None
        mock_post.return_value = mock_response
        
        # Test data
        data = [
            {"id": "12345", "name": "Product A", "price": 99.99},
            {"id": "67890", "name": "Product B", "price": -10.00}
        ]
        
        # Validate data
        results = self.api_client.validate_data(data, rules=["price_range"])
        
        # Verify request
        mock_post.assert_called_once()
        args, kwargs = mock_post.call_args
        self.assertEqual(kwargs["json"]["data"], data)
        self.assertEqual(kwargs["json"]["rules"], ["price_range"])
        
        # Verify result
        self.assertEqual(len(results), 2)
        
        self.assertEqual(results[0].rule_id, "price_range")
        self.assertEqual(results[0].field, "price")
        self.assertEqual(results[0].value, -10.0)
        self.assertEqual(results[0].message, "Price must be greater than 0.01")
        self.assertEqual(results[0].row_index, 1)
        self.assertEqual(results[0].passed, False)
        self.assertEqual(results[0].timestamp, datetime.fromisoformat("2025-05-12T12:00:00"))
        
        self.assertEqual(results[1].rule_id, "price_range")
        self.assertEqual(results[1].field, "price")
        self.assertEqual(results[1].value, 99.99)
        self.assertEqual(results[1].message, "Price is within acceptable range")
        self.assertEqual(results[1].row_index, 0)
        self.assertEqual(results[1].passed, True)
        self.assertEqual(results[1].timestamp, datetime.fromisoformat("2025-05-12T12:00:00"))

    @patch('requests.Session.request')
    def test_request_error_handling(self, mock_request):
        """Test error handling in the _request method."""
        # Mock a request exception
        mock_request.side_effect = requests.RequestException("Connection error")
        
        # Attempt to get metrics
        with self.assertRaises(requests.RequestException):
            self.api_client.get_metrics()
        
        # Verify request was attempted
        mock_request.assert_called_once()


if __name__ == "__main__":
    unittest.main()
