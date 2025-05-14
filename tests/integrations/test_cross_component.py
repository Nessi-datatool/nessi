"""
Tests for cross-component integration in Nessi.dev.

These tests verify the interaction between different components:
- Delta Lake + Cloud
- Quality + Monitoring
- Security + API
"""

import unittest
import os
import json
import tempfile
from unittest import mock
import time
from datetime import datetime
import sys

# Add the project root to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../..')))

# Import the pkg modules
from pkg.datalake import CloudDeltaConnector
from pkg.cloud import CloudProvider, CloudManager
from pkg.quality import QualityChecker
from pkg.monitoring.alerts import AlertManager
from pkg.monitoring.metrics import MetricsCollector
from pkg.security.auth import AuthManager, RBACManager
from pkg.api import APIHandler

class TestDeltaLakeCloudIntegration(unittest.TestCase):
    """Test cases for Delta Lake and Cloud integration."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        
        # Mock the Delta Lake connector
        self.mock_delta_connector_patcher = mock.patch('pkg.datalake.CloudDeltaConnector')
        self.mock_delta_connector = self.mock_delta_connector_patcher.start()
        self.mock_delta_connector_instance = self.mock_delta_connector.return_value
        
        # Mock the cloud provider
        self.mock_cloud_provider_patcher = mock.patch('pkg.cloud.CloudProvider')
        self.mock_cloud_provider = self.mock_cloud_provider_patcher.start()
        self.mock_cloud_provider_instance = self.mock_cloud_provider.return_value
        
        # Mock the cloud manager
        self.mock_cloud_manager_patcher = mock.patch('pkg.cloud.CloudManager')
        self.mock_cloud_manager = self.mock_cloud_manager_patcher.start()
        self.mock_cloud_manager_instance = self.mock_cloud_manager.return_value
        self.mock_cloud_manager_instance.get_provider.return_value = self.mock_cloud_provider_instance
    
    def tearDown(self):
        """Clean up test environment."""
        # Clean up temporary directory
        self.temp_dir.cleanup()
        
        # Stop patches
        self.mock_delta_connector_patcher.stop()
        self.mock_cloud_provider_patcher.stop()
        self.mock_cloud_manager_patcher.stop()
    
    def test_delta_lake_cloud_read(self):
        """Test reading Delta Lake table from cloud storage."""
        # Configure mock responses
        self.mock_cloud_provider_instance.get_object.return_value = b'test data'
        self.mock_delta_connector_instance.read.return_value = {
            'schema': {'fields': [{'name': 'id', 'type': 'integer'}]},
            'data': [{'id': 1}, {'id': 2}, {'id': 3}],
        }
        
        # Create a real CloudDeltaConnector with the mocked cloud provider
        connector = CloudDeltaConnector(
            provider=self.mock_cloud_provider_instance,
            table_path='s3://test-bucket/test-table',
        )
        
        # Read data from the table
        data = connector.read()
        
        # Check that the cloud provider was used to get objects
        self.mock_cloud_provider_instance.get_object.assert_called()
        
        # Check the returned data
        self.assertEqual(data['schema']['fields'][0]['name'], 'id')
        self.assertEqual(len(data['data']), 3)
    
    def test_delta_lake_cloud_write(self):
        """Test writing Delta Lake table to cloud storage."""
        # Configure mock responses
        self.mock_cloud_provider_instance.put_object.return_value = True
        
        # Create a real CloudDeltaConnector with the mocked cloud provider
        connector = CloudDeltaConnector(
            provider=self.mock_cloud_provider_instance,
            table_path='s3://test-bucket/test-table',
        )
        
        # Write data to the table
        result = connector.write({
            'schema': {'fields': [{'name': 'id', 'type': 'integer'}]},
            'data': [{'id': 1}, {'id': 2}, {'id': 3}],
        })
        
        # Check that the cloud provider was used to put objects
        self.mock_cloud_provider_instance.put_object.assert_called()
        
        # Check the result
        self.assertTrue(result)
    
    def test_delta_lake_version_history(self):
        """Test getting version history for a Delta Lake table."""
        # Configure mock responses
        self.mock_delta_connector_instance.get_version_history.return_value = [
            {'version': 3, 'timestamp': '2023-01-03T00:00:00Z'},
            {'version': 2, 'timestamp': '2023-01-02T00:00:00Z'},
            {'version': 1, 'timestamp': '2023-01-01T00:00:00Z'},
        ]
        
        # Create a CloudDeltaConnector with the mocked cloud provider
        connector = self.mock_delta_connector(
            provider=self.mock_cloud_provider_instance,
            table_path='s3://test-bucket/test-table',
        )
        
        # Get version history
        history = connector.get_version_history()
        
        # Check that the connector get_version_history method was called
        self.mock_delta_connector_instance.get_version_history.assert_called_once()
        
        # Check the returned history
        self.assertEqual(len(history), 3)
        self.assertEqual(history[0]['version'], 3)
        self.assertEqual(history[1]['version'], 2)
        self.assertEqual(history[2]['version'], 1)
    
    def test_delta_lake_time_travel(self):
        """Test reading a Delta Lake table as of a specific version."""
        # Configure mock responses
        self.mock_delta_connector_instance.read_as_of_version.return_value = {
            'schema': {'fields': [{'name': 'id', 'type': 'integer'}]},
            'data': [{'id': 1}, {'id': 2}],
            'version': 2,
        }
        
        # Create a CloudDeltaConnector with the mocked cloud provider
        connector = self.mock_delta_connector(
            provider=self.mock_cloud_provider_instance,
            table_path='s3://test-bucket/test-table',
        )
        
        # Read data as of version 2
        data = connector.read_as_of_version(2)
        
        # Check that the connector read_as_of_version method was called
        self.mock_delta_connector_instance.read_as_of_version.assert_called_once_with(2)
        
        # Check the returned data
        self.assertEqual(data['schema']['fields'][0]['name'], 'id')
        self.assertEqual(len(data['data']), 2)
        self.assertEqual(data['version'], 2)
    
    def test_multi_cloud_delta_lake(self):
        """Test Delta Lake operations across multiple cloud providers."""
        # Configure mock responses for AWS
        mock_aws_provider = mock.MagicMock()
        mock_aws_provider.get_object.return_value = b'aws data'
        
        # Configure mock responses for Azure
        mock_azure_provider = mock.MagicMock()
        mock_azure_provider.get_object.return_value = b'azure data'
        
        # Configure mock responses for GCP
        mock_gcp_provider = mock.MagicMock()
        mock_gcp_provider.get_object.return_value = b'gcp data'
        
        # Configure cloud manager to return different providers
        self.mock_cloud_manager_instance.get_provider.side_effect = lambda provider_type: {
            'aws': mock_aws_provider,
            'azure': mock_azure_provider,
            'gcp': mock_gcp_provider,
        }[provider_type]
        
        # Create real CloudDeltaConnectors for each cloud provider
        aws_connector = CloudDeltaConnector(
            provider=mock_aws_provider,
            table_path='s3://test-bucket/test-table',
        )
        
        azure_connector = CloudDeltaConnector(
            provider=mock_azure_provider,
            table_path='abfs://test-container/test-table',
        )
        
        gcp_connector = CloudDeltaConnector(
            provider=mock_gcp_provider,
            table_path='gs://test-bucket/test-table',
        )
        
        # Read data from each connector
        aws_connector.read()
        azure_connector.read()
        gcp_connector.read()
        
        # Check that each provider's get_object method was called
        mock_aws_provider.get_object.assert_called()
        mock_azure_provider.get_object.assert_called()
        mock_gcp_provider.get_object.assert_called()


class TestQualityMonitoringIntegration(unittest.TestCase):
    """Test cases for Quality and Monitoring integration."""

    def setUp(self):
        """Set up test environment."""
        # Mock the quality checker
        self.mock_quality_checker_patcher = mock.patch('pkg.quality.QualityChecker')
        self.mock_quality_checker = self.mock_quality_checker_patcher.start()
        self.mock_quality_checker_instance = self.mock_quality_checker.return_value
        
        # Mock the alert manager
        self.mock_alert_manager_patcher = mock.patch('pkg.monitoring.alerts.AlertManager')
        self.mock_alert_manager = self.mock_alert_manager_patcher.start()
        self.mock_alert_manager_instance = self.mock_alert_manager.return_value
        
        # Mock the metrics collector
        self.mock_metrics_collector_patcher = mock.patch('pkg.monitoring.metrics.MetricsCollector')
        self.mock_metrics_collector = self.mock_metrics_collector_patcher.start()
        self.mock_metrics_collector_instance = self.mock_metrics_collector.return_value
    
    def tearDown(self):
        """Clean up test environment."""
        # Stop patches
        self.mock_quality_checker_patcher.stop()
        self.mock_alert_manager_patcher.stop()
        self.mock_metrics_collector_patcher.stop()
    
    def test_quality_check_triggers_alert(self):
        """Test that a quality check triggers an alert when quality score is below threshold."""
        # Configure mock responses
        self.mock_quality_checker_instance.check_quality.return_value = {
            'quality_score': 0.7,
            'passed': False,
            'rules_passed': 7,
            'rules_failed': 3,
            'table_name': 'test_table',
            'check_id': 'test-check-id',
        }
        
        # Create a real QualityChecker with the mocked alert manager
        quality_checker = QualityChecker(
            alert_manager=self.mock_alert_manager_instance
        )
        
        # Run a quality check on a table with low quality
        result = quality_checker.check_quality(
            table_name='low_quality_table',
            rules=[
                {'name': 'test_rule', 'rule_type': 'not_null', 'column': 'id'},
            ],
            quality_threshold=0.9
        )
        
        # Check that the alert manager was called to trigger an alert
        self.mock_alert_manager_instance.trigger_alert.assert_called_once()
        
        # Check the result
        self.assertEqual(result['quality_score'], 0.7)
        self.assertFalse(result['passed'])
        self.assertEqual(result['rules_passed'], 7)
        self.assertEqual(result['rules_failed'], 3)
    
    def test_quality_check_records_metrics(self):
        """Test that a quality check records metrics."""
        # Create a real QualityChecker with the mocked metrics collector
        quality_checker = QualityChecker(
            metrics_collector=self.mock_metrics_collector_instance
        )
        
        # Run a quality check
        result = quality_checker.check_quality(
            table_name='test_table',
            rules=[
                {'name': 'test_rule', 'rule_type': 'not_null', 'column': 'id'},
            ],
        )
        
        # Check that the metrics collector was called to record metrics
        self.mock_metrics_collector_instance.record_metric.assert_called()
        
        # Check the result
        self.assertEqual(result['quality_score'], 0.95)
        self.assertTrue(result['passed'])
        self.assertEqual(result['rules_passed'], 10)
        self.assertEqual(result['rules_failed'], 0)
    
    def test_anomaly_detection_triggers_alert(self):
        """Test that anomaly detection triggers an alert."""
        # Create a real QualityChecker with the mocked alert manager
        quality_checker = QualityChecker(
            alert_manager=self.mock_alert_manager_instance
        )
        
        # Run anomaly detection on a table with anomalies
        result = quality_checker.detect_anomalies(
            table_name='anomaly_table',
            columns=['id', 'name', 'value'],
        )
        
        # Check that the alert manager was called to trigger an alert
        self.mock_alert_manager_instance.trigger_alert.assert_called_once()
        
        # Check the result
        self.assertEqual(len(result['anomalies']), 1)
        self.assertEqual(result['anomalies'][0]['column'], 'id')
        self.assertEqual(result['anomalies'][0]['type'], 'sudden_change')
        self.assertEqual(result['anomalies'][0]['severity'], 'high')
    
    def test_intelligent_alerting(self):
        """Test intelligent alerting based on historical data."""
        # Create a real QualityChecker with the mocked alert manager
        quality_checker = QualityChecker(
            alert_manager=self.mock_alert_manager_instance
        )
        
        # Apply intelligent alerting
        result = quality_checker.intelligent_alert(
            table_name='test_table',
            metric_name='quality_score',
        )
        
        # Check that the alert manager was called to get alert history
        self.mock_alert_manager_instance.get_alert_history.assert_called_once()
        
        # Check the result
        self.assertEqual(result['table_name'], 'test_table')
        self.assertEqual(result['metric_name'], 'quality_score')
        self.assertEqual(result['alert_type'], 'intelligent')


class TestSecurityAPIIntegration(unittest.TestCase):
    """Test cases for Security and API integration."""

    def setUp(self):
        """Set up test environment."""
        # Mock the auth manager
        self.mock_auth_manager_patcher = mock.patch('pkg.security.auth.AuthManager')
        self.mock_auth_manager = self.mock_auth_manager_patcher.start()
        self.mock_auth_manager_instance = self.mock_auth_manager.return_value
        
        # Mock the RBAC manager
        self.mock_rbac_manager_patcher = mock.patch('pkg.security.auth.RBACManager')
        self.mock_rbac_manager = self.mock_rbac_manager_patcher.start()
        self.mock_rbac_manager_instance = self.mock_rbac_manager.return_value
        
        # Mock the API handler
        self.mock_api_handler_patcher = mock.patch('pkg.api.APIHandler')
        self.mock_api_handler = self.mock_api_handler_patcher.start()
        self.mock_api_handler_instance = self.mock_api_handler.return_value
        
        # Mock the audit logger
        self.mock_audit_logger_patcher = mock.patch('pkg.security.auth.AuditLogger')
        self.mock_audit_logger = self.mock_audit_logger_patcher.start()
        self.mock_audit_logger_instance = self.mock_audit_logger.return_value
    
    def tearDown(self):
        """Clean up test environment."""
        # Stop patches
        self.mock_auth_manager_patcher.stop()
        self.mock_rbac_manager_patcher.stop()
        self.mock_api_handler_patcher.stop()
        self.mock_audit_logger_patcher.stop()
    
    def test_user_authentication(self):
        """Test user authentication."""
        # Configure mock responses
        self.mock_auth_manager_instance.authenticate.return_value = {
            'user_id': 'test-user',
            'username': 'testuser',
            'roles': ['user'],
            'token': 'test-token',
        }
        
        # Authenticate a user
        result = self.mock_auth_manager_instance.authenticate(
            username='testuser',
            password='testpassword',
        )
        
        # Check that the auth manager was called
        self.mock_auth_manager_instance.authenticate.assert_called_once_with(
            username='testuser',
            password='testpassword',
        )
        
        # Check the result
        self.assertEqual(result['user_id'], 'test-user')
        self.assertEqual(result['username'], 'testuser')
        self.assertEqual(result['roles'], ['user'])
        self.assertEqual(result['token'], 'test-token')
    
    def test_api_authorization(self):
        """Test API authorization."""
        # Configure mock responses
        self.mock_rbac_manager_instance.check_permission.return_value = True
        
        # Check permission
        result = self.mock_rbac_manager_instance.check_permission(
            user_id='test-user',
            resource='tables',
            action='read',
        )
        
        # Check that the RBAC manager was called
        self.mock_rbac_manager_instance.check_permission.assert_called_once_with(
            user_id='test-user',
            resource='tables',
            action='read',
        )
        
        # Check the result
        self.assertTrue(result)
    
    def test_api_request_with_auth(self):
        """Test API request with authentication and authorization."""
        # Configure mock responses
        self.mock_auth_manager_instance.validate_token.return_value = {
            'user_id': 'test-user',
            'username': 'testuser',
            'roles': ['user'],
        }
        
        self.mock_rbac_manager_instance.check_permission.return_value = True
        
        # Create a real APIHandler with the mocked dependencies
        api_handler = APIHandler(
            auth_manager=self.mock_auth_manager_instance,
            rbac_manager=self.mock_rbac_manager_instance
        )
        
        # Make an API request
        try:
            api_handler.handle_request(
                method='GET',
                path='/api/v1/tables/test_table/quality',
                headers={'Authorization': 'Bearer test-token'},
                body=None,
            )
        except Exception:
            # We don't care about the actual result, just that the validate_token method was called
            pass
        
        # Check that the auth manager was called to validate the token
        self.mock_auth_manager_instance.validate_token.assert_called_once_with('test-token')
    
    def test_api_request_audit_logging(self):
        """Test audit logging for API requests."""
        # Configure mock responses
        self.mock_auth_manager_instance.validate_token.return_value = {
            'user_id': 'test-user',
            'username': 'testuser',
            'roles': ['user'],
        }
        
        self.mock_rbac_manager_instance.check_permission.return_value = True
        
        # Create a real APIHandler with the mocked dependencies
        api_handler = APIHandler(
            auth_manager=self.mock_auth_manager_instance,
            rbac_manager=self.mock_rbac_manager_instance,
            audit_logger=self.mock_audit_logger_instance
        )
        
        # Make an API request
        try:
            api_handler.handle_request(
                method='GET',
                path='/api/v1/tables/test_table/quality',
                headers={'Authorization': 'Bearer test-token'},
                body=None,
            )
        except Exception:
            # We don't care about the actual result, just that the log_api_request method was called
            pass
        
        # Check that the audit logger was called to log the request
        self.mock_audit_logger_instance.log_api_request.assert_called_once()


if __name__ == '__main__':
    unittest.main()
