#!/usr/bin/env python3
"""
Unit tests for the Nessi client models.
"""

import unittest
from datetime import datetime
import json

import sys
import os
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from nessi_client.models import Metric, Alert, Rule, Profile, ValidationResult


class TestModels(unittest.TestCase):
    """Test cases for the data models."""

    def test_metric_model(self):
        """Test the Metric model."""
        # Create a metric
        timestamp = datetime(2025, 5, 12, 12, 0, 0)
        metric = Metric(
            name="cpu_usage",
            value=45.2,
            timestamp=timestamp,
            tags={"host": "server-01", "environment": "production"},
            metadata={"source": "system_monitor"}
        )
        
        # Verify attributes
        self.assertEqual(metric.name, "cpu_usage")
        self.assertEqual(metric.value, 45.2)
        self.assertEqual(metric.timestamp, timestamp)
        self.assertEqual(metric.tags, {"host": "server-01", "environment": "production"})
        self.assertEqual(metric.metadata, {"source": "system_monitor"})
        
        # Test default values
        default_metric = Metric(name="test", value=1)
        self.assertIsNotNone(default_metric.timestamp)
        self.assertEqual(default_metric.tags, {})
        self.assertEqual(default_metric.metadata, {})
        
        # Test with different value types
        int_metric = Metric(name="count", value=100)
        self.assertEqual(int_metric.value, 100)
        
        str_metric = Metric(name="status", value="ok")
        self.assertEqual(str_metric.value, "ok")

    def test_alert_model(self):
        """Test the Alert model."""
        # Create an alert
        timestamp = datetime(2025, 5, 12, 12, 0, 0)
        alert = Alert(
            id="alert-123",
            name="High CPU Usage",
            description="CPU usage exceeded threshold",
            severity="warning",
            triggered=True,
            timestamp=timestamp,
            metric_name="cpu_usage",
            metric_value=85.0,
            threshold=80.0,
            tags={"host": "server-01"},
            metadata={"duration": "5m"}
        )
        
        # Verify attributes
        self.assertEqual(alert.id, "alert-123")
        self.assertEqual(alert.name, "High CPU Usage")
        self.assertEqual(alert.description, "CPU usage exceeded threshold")
        self.assertEqual(alert.severity, "warning")
        self.assertTrue(alert.triggered)
        self.assertEqual(alert.timestamp, timestamp)
        self.assertEqual(alert.metric_name, "cpu_usage")
        self.assertEqual(alert.metric_value, 85.0)
        self.assertEqual(alert.threshold, 80.0)
        self.assertEqual(alert.tags, {"host": "server-01"})
        self.assertEqual(alert.metadata, {"duration": "5m"})
        
        # Test default values
        default_alert = Alert(
            id="alert-456",
            name="Test Alert",
            description="Test description",
            severity="info",
            triggered=False
        )
        self.assertIsNotNone(default_alert.timestamp)
        self.assertIsNone(default_alert.metric_name)
        self.assertIsNone(default_alert.metric_value)
        self.assertIsNone(default_alert.threshold)
        self.assertEqual(default_alert.tags, {})
        self.assertEqual(default_alert.metadata, {})

    def test_rule_model(self):
        """Test the Rule model."""
        # Create a rule
        rule = Rule(
            id="price-range",
            name="Price Range Check",
            description="Validates that product prices are within an acceptable range",
            severity="error",
            field="price",
            rule_type="range",
            config={"min": 0.01, "max": 9999.99},
            tags=["product", "validation"],
            enabled=True
        )
        
        # Verify attributes
        self.assertEqual(rule.id, "price-range")
        self.assertEqual(rule.name, "Price Range Check")
        self.assertEqual(rule.description, "Validates that product prices are within an acceptable range")
        self.assertEqual(rule.severity, "error")
        self.assertEqual(rule.field, "price")
        self.assertEqual(rule.rule_type, "range")
        self.assertEqual(rule.config, {"min": 0.01, "max": 9999.99})
        self.assertEqual(rule.tags, ["product", "validation"])
        self.assertTrue(rule.enabled)
        
        # Test default values
        default_rule = Rule(
            id="test-rule",
            name="Test Rule",
            description="Test description",
            severity="warning"
        )
        self.assertIsNone(default_rule.field)
        self.assertEqual(default_rule.rule_type, "custom")
        self.assertEqual(default_rule.config, {})
        self.assertEqual(default_rule.tags, [])
        self.assertTrue(default_rule.enabled)

    def test_validation_result_model(self):
        """Test the ValidationResult model."""
        # Create a validation result
        timestamp = datetime(2025, 5, 12, 12, 0, 0)
        result = ValidationResult(
            rule_id="price-range",
            field="price",
            value=-10.0,
            message="Price must be greater than 0.01",
            row_index=1,
            passed=False,
            timestamp=timestamp
        )
        
        # Verify attributes
        self.assertEqual(result.rule_id, "price-range")
        self.assertEqual(result.field, "price")
        self.assertEqual(result.value, -10.0)
        self.assertEqual(result.message, "Price must be greater than 0.01")
        self.assertEqual(result.row_index, 1)
        self.assertFalse(result.passed)
        self.assertEqual(result.timestamp, timestamp)
        
        # Test default values
        default_result = ValidationResult(
            rule_id="test-rule",
            field="test_field",
            value="test_value",
            message="Test message"
        )
        self.assertIsNone(default_result.row_index)
        self.assertFalse(default_result.passed)
        self.assertIsNotNone(default_result.timestamp)

    def test_profile_model(self):
        """Test the Profile model."""
        # Create a profile
        timestamp = datetime(2025, 5, 12, 12, 0, 0)
        profile = Profile(
            id="profile-123",
            name="Customer Data Profile",
            dataset_name="customers",
            timestamp=timestamp,
            row_count=1000,
            column_count=10,
            column_stats={
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
            outliers={
                "age": [17, 86, 90]
            },
            metadata={"source": "database"},
            detection_method="zscore"
        )
        
        # Verify attributes
        self.assertEqual(profile.id, "profile-123")
        self.assertEqual(profile.name, "Customer Data Profile")
        self.assertEqual(profile.dataset_name, "customers")
        self.assertEqual(profile.timestamp, timestamp)
        self.assertEqual(profile.row_count, 1000)
        self.assertEqual(profile.column_count, 10)
        self.assertEqual(profile.column_stats["customer_id"]["type"], "string")
        self.assertEqual(profile.column_stats["age"]["mean"], 42.5)
        self.assertEqual(profile.outliers["age"], [17, 86, 90])
        self.assertEqual(profile.metadata, {"source": "database"})
        self.assertEqual(profile.detection_method, "zscore")
        
        # Test default values
        default_profile = Profile(
            id="profile-456",
            name="Test Profile",
            dataset_name="test_dataset"
        )
        self.assertIsNotNone(default_profile.timestamp)
        self.assertEqual(default_profile.row_count, 0)
        self.assertEqual(default_profile.column_count, 0)
        self.assertEqual(default_profile.column_stats, {})
        self.assertEqual(default_profile.outliers, {})
        self.assertEqual(default_profile.metadata, {})
        self.assertEqual(default_profile.detection_method, "zscore")


if __name__ == "__main__":
    unittest.main()
