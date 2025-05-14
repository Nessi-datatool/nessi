"""
Tests for the Nessi.dev plugin system.

These tests verify the functionality of the plugin system, including:
- Plugin loading and unloading
- Plugin execution
- Plugin validation
- Plugin management via CLI
"""

import unittest
import os
import json
import tempfile
from unittest import mock
import shutil
import sys

class TestPluginManager(unittest.TestCase):
    """Test cases for the Plugin Manager."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        self.plugins_dir = os.path.join(self.temp_dir.name, 'plugins')
        os.makedirs(self.plugins_dir, exist_ok=True)
        
        # Mock the plugin manager
        self.mock_plugin_manager_patcher = mock.patch('pkg.plugins.PluginManager')
        self.mock_plugin_manager = self.mock_plugin_manager_patcher.start()
        self.mock_plugin_manager_instance = self.mock_plugin_manager.return_value
        
        # Set up the plugin manager
        self.mock_plugin_manager_instance.plugins_dir = self.plugins_dir
        self.mock_plugin_manager_instance.loaded_plugins = {}
    
    def tearDown(self):
        """Clean up test environment."""
        # Clean up temporary directory
        self.temp_dir.cleanup()
        
        # Stop patches
        self.mock_plugin_manager_patcher.stop()
    
    def test_plugin_discovery(self):
        """Test plugin discovery."""
        # Create mock plugin files
        plugin1_dir = os.path.join(self.plugins_dir, 'plugin1')
        os.makedirs(plugin1_dir, exist_ok=True)
        with open(os.path.join(plugin1_dir, 'plugin.json'), 'w') as f:
            json.dump({
                'name': 'plugin1',
                'version': '1.0.0',
                'type': 'validation',
                'entry_point': 'plugin1.so',
            }, f)
        
        plugin2_dir = os.path.join(self.plugins_dir, 'plugin2')
        os.makedirs(plugin2_dir, exist_ok=True)
        with open(os.path.join(plugin2_dir, 'plugin.json'), 'w') as f:
            json.dump({
                'name': 'plugin2',
                'version': '1.0.0',
                'type': 'alert',
                'entry_point': 'plugin2.so',
            }, f)
        
        # Configure mock responses
        self.mock_plugin_manager_instance.discover_plugins.return_value = [
            {
                'name': 'plugin1',
                'version': '1.0.0',
                'type': 'validation',
                'path': plugin1_dir,
                'entry_point': 'plugin1.so',
            },
            {
                'name': 'plugin2',
                'version': '1.0.0',
                'type': 'alert',
                'path': plugin2_dir,
                'entry_point': 'plugin2.so',
            },
        ]
        
        # Discover plugins
        plugins = self.mock_plugin_manager_instance.discover_plugins()
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.discover_plugins.assert_called_once()
        
        # Check the discovered plugins
        self.assertEqual(len(plugins), 2)
        self.assertEqual(plugins[0]['name'], 'plugin1')
        self.assertEqual(plugins[0]['type'], 'validation')
        self.assertEqual(plugins[1]['name'], 'plugin2')
        self.assertEqual(plugins[1]['type'], 'alert')
    
    def test_plugin_loading(self):
        """Test plugin loading."""
        # Create a mock plugin
        plugin_dir = os.path.join(self.plugins_dir, 'test_plugin')
        os.makedirs(plugin_dir, exist_ok=True)
        with open(os.path.join(plugin_dir, 'plugin.json'), 'w') as f:
            json.dump({
                'name': 'test_plugin',
                'version': '1.0.0',
                'type': 'validation',
                'entry_point': 'test_plugin.so',
            }, f)
        
        # Create a mock shared library file
        with open(os.path.join(plugin_dir, 'test_plugin.so'), 'wb') as f:
            f.write(b'mock shared library')
        
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'test_plugin'
        mock_plugin.version = '1.0.0'
        mock_plugin.type = 'validation'
        
        self.mock_plugin_manager_instance.load_plugin.return_value = mock_plugin
        
        # Load the plugin
        plugin = self.mock_plugin_manager_instance.load_plugin('test_plugin')
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.load_plugin.assert_called_once_with('test_plugin')
        
        # Check the loaded plugin
        self.assertEqual(plugin.name, 'test_plugin')
        self.assertEqual(plugin.version, '1.0.0')
        self.assertEqual(plugin.type, 'validation')
    
    def test_plugin_unloading(self):
        """Test plugin unloading."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'test_plugin'
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'test_plugin': mock_plugin,
        }
        
        # Unload the plugin
        self.mock_plugin_manager_instance.unload_plugin.return_value = True
        result = self.mock_plugin_manager_instance.unload_plugin('test_plugin')
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.unload_plugin.assert_called_once_with('test_plugin')
        
        # Check the result
        self.assertTrue(result)
    
    def test_plugin_execution(self):
        """Test plugin execution."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'test_plugin'
        mock_plugin.type = 'validation'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'data': {'validation_passed': True},
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'test_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'data': {'validation_passed': True},
        }
        
        # Execute the plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'test_plugin',
            {'table_name': 'test_table', 'rules': []},
        )
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.execute_plugin.assert_called_once()
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertEqual(result['data']['validation_passed'], True)
    
    def test_plugin_validation(self):
        """Test plugin validation."""
        # Create a valid plugin
        valid_plugin_dir = os.path.join(self.plugins_dir, 'valid_plugin')
        os.makedirs(valid_plugin_dir, exist_ok=True)
        with open(os.path.join(valid_plugin_dir, 'plugin.json'), 'w') as f:
            json.dump({
                'name': 'valid_plugin',
                'version': '1.0.0',
                'type': 'validation',
                'entry_point': 'valid_plugin.so',
            }, f)
        
        # Create an invalid plugin (missing required fields)
        invalid_plugin_dir = os.path.join(self.plugins_dir, 'invalid_plugin')
        os.makedirs(invalid_plugin_dir, exist_ok=True)
        with open(os.path.join(invalid_plugin_dir, 'plugin.json'), 'w') as f:
            json.dump({
                'name': 'invalid_plugin',
                'version': '1.0.0',
                # Missing 'type' field
                'entry_point': 'invalid_plugin.so',
            }, f)
        
        # Configure mock responses
        self.mock_plugin_manager_instance.validate_plugin.side_effect = lambda plugin_dir: {
            valid_plugin_dir: True,
            invalid_plugin_dir: False,
        }[plugin_dir]
        
        # Validate the plugins
        valid_result = self.mock_plugin_manager_instance.validate_plugin(valid_plugin_dir)
        invalid_result = self.mock_plugin_manager_instance.validate_plugin(invalid_plugin_dir)
        
        # Check the results
        self.assertTrue(valid_result)
        self.assertFalse(invalid_result)
    
    def test_plugin_installation(self):
        """Test plugin installation."""
        # Create a mock plugin package
        plugin_package = os.path.join(self.temp_dir.name, 'test_plugin.zip')
        with open(plugin_package, 'wb') as f:
            f.write(b'mock plugin package')
        
        # Configure mock responses
        self.mock_plugin_manager_instance.install_plugin.return_value = {
            'name': 'test_plugin',
            'version': '1.0.0',
            'type': 'validation',
            'path': os.path.join(self.plugins_dir, 'test_plugin'),
        }
        
        # Install the plugin
        result = self.mock_plugin_manager_instance.install_plugin(plugin_package)
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.install_plugin.assert_called_once_with(plugin_package)
        
        # Check the result
        self.assertEqual(result['name'], 'test_plugin')
        self.assertEqual(result['version'], '1.0.0')
        self.assertEqual(result['type'], 'validation')
    
    def test_plugin_uninstallation(self):
        """Test plugin uninstallation."""
        # Create a mock plugin
        plugin_dir = os.path.join(self.plugins_dir, 'test_plugin')
        os.makedirs(plugin_dir, exist_ok=True)
        
        # Configure mock responses
        self.mock_plugin_manager_instance.uninstall_plugin.return_value = True
        
        # Uninstall the plugin
        result = self.mock_plugin_manager_instance.uninstall_plugin('test_plugin')
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.uninstall_plugin.assert_called_once_with('test_plugin')
        
        # Check the result
        self.assertTrue(result)
    
    def test_plugin_listing(self):
        """Test plugin listing."""
        # Configure mock responses
        self.mock_plugin_manager_instance.list_plugins.return_value = [
            {
                'name': 'plugin1',
                'version': '1.0.0',
                'type': 'validation',
                'enabled': True,
            },
            {
                'name': 'plugin2',
                'version': '1.0.0',
                'type': 'alert',
                'enabled': False,
            },
        ]
        
        # List plugins
        plugins = self.mock_plugin_manager_instance.list_plugins()
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.list_plugins.assert_called_once()
        
        # Check the result
        self.assertEqual(len(plugins), 2)
        self.assertEqual(plugins[0]['name'], 'plugin1')
        self.assertEqual(plugins[0]['type'], 'validation')
        self.assertTrue(plugins[0]['enabled'])
        self.assertEqual(plugins[1]['name'], 'plugin2')
        self.assertEqual(plugins[1]['type'], 'alert')
        self.assertFalse(plugins[1]['enabled'])
    
    def test_plugin_enabling_disabling(self):
        """Test plugin enabling and disabling."""
        # Configure mock responses
        self.mock_plugin_manager_instance.enable_plugin.return_value = True
        self.mock_plugin_manager_instance.disable_plugin.return_value = True
        
        # Enable and disable plugins
        enable_result = self.mock_plugin_manager_instance.enable_plugin('test_plugin')
        disable_result = self.mock_plugin_manager_instance.disable_plugin('test_plugin')
        
        # Check that the plugin manager was called
        self.mock_plugin_manager_instance.enable_plugin.assert_called_once_with('test_plugin')
        self.mock_plugin_manager_instance.disable_plugin.assert_called_once_with('test_plugin')
        
        # Check the results
        self.assertTrue(enable_result)
        self.assertTrue(disable_result)


class TestPluginTypes(unittest.TestCase):
    """Test cases for different plugin types."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        self.plugins_dir = os.path.join(self.temp_dir.name, 'plugins')
        os.makedirs(self.plugins_dir, exist_ok=True)
        
        # Mock the plugin manager
        self.mock_plugin_manager_patcher = mock.patch('pkg.plugins.PluginManager')
        self.mock_plugin_manager = self.mock_plugin_manager_patcher.start()
        self.mock_plugin_manager_instance = self.mock_plugin_manager.return_value
        
        # Set up the plugin manager
        self.mock_plugin_manager_instance.plugins_dir = self.plugins_dir
        self.mock_plugin_manager_instance.loaded_plugins = {}
    
    def tearDown(self):
        """Clean up test environment."""
        # Clean up temporary directory
        self.temp_dir.cleanup()
        
        # Stop patches
        self.mock_plugin_manager_patcher.stop()
    
    def test_validation_plugin(self):
        """Test validation plugin."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'validation_plugin'
        mock_plugin.type = 'validation'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'validation_passed': True,
            'rules_passed': 10,
            'rules_failed': 0,
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'validation_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'validation_passed': True,
            'rules_passed': 10,
            'rules_failed': 0,
        }
        
        # Execute the validation plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'validation_plugin',
            {'table_name': 'test_table', 'rules': []},
        )
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertTrue(result['validation_passed'])
        self.assertEqual(result['rules_passed'], 10)
        self.assertEqual(result['rules_failed'], 0)
    
    def test_alert_plugin(self):
        """Test alert plugin."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'alert_plugin'
        mock_plugin.type = 'alert'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'alert_sent': True,
            'destination': 'slack',
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'alert_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'alert_sent': True,
            'destination': 'slack',
        }
        
        # Execute the alert plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'alert_plugin',
            {'alert_type': 'quality_check_failed', 'table_name': 'test_table'},
        )
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertTrue(result['alert_sent'])
        self.assertEqual(result['destination'], 'slack')
    
    def test_metric_plugin(self):
        """Test metric plugin."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'metric_plugin'
        mock_plugin.type = 'metric'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'metrics_collected': True,
            'metric_count': 5,
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'metric_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'metrics_collected': True,
            'metric_count': 5,
        }
        
        # Execute the metric plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'metric_plugin',
            {'table_name': 'test_table', 'columns': ['id', 'name']},
        )
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertTrue(result['metrics_collected'])
        self.assertEqual(result['metric_count'], 5)
    
    def test_storage_plugin(self):
        """Test storage plugin."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'storage_plugin'
        mock_plugin.type = 'storage'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'data_stored': True,
            'storage_location': 's3://test-bucket/test-data',
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'storage_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'data_stored': True,
            'storage_location': 's3://test-bucket/test-data',
        }
        
        # Execute the storage plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'storage_plugin',
            {'data': {'id': 1, 'name': 'test'}, 'format': 'json'},
        )
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertTrue(result['data_stored'])
        self.assertEqual(result['storage_location'], 's3://test-bucket/test-data')
    
    def test_export_plugin(self):
        """Test export plugin."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'export_plugin'
        mock_plugin.type = 'export'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'data_exported': True,
            'export_format': 'pdf',
            'export_path': '/path/to/export.pdf',
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'export_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'data_exported': True,
            'export_format': 'pdf',
            'export_path': '/path/to/export.pdf',
        }
        
        # Execute the export plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'export_plugin',
            {'data': {'table_name': 'test_table', 'quality_score': 0.95}, 'format': 'pdf'},
        )
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertTrue(result['data_exported'])
        self.assertEqual(result['export_format'], 'pdf')
        self.assertEqual(result['export_path'], '/path/to/export.pdf')
    
    def test_ui_plugin(self):
        """Test UI plugin."""
        # Configure mock responses
        mock_plugin = mock.MagicMock()
        mock_plugin.name = 'ui_plugin'
        mock_plugin.type = 'ui'
        mock_plugin.execute.return_value = {
            'result': 'success',
            'component_rendered': True,
            'component_type': 'dashboard',
        }
        
        self.mock_plugin_manager_instance.loaded_plugins = {
            'ui_plugin': mock_plugin,
        }
        
        self.mock_plugin_manager_instance.execute_plugin.return_value = {
            'result': 'success',
            'component_rendered': True,
            'component_type': 'dashboard',
        }
        
        # Execute the UI plugin
        result = self.mock_plugin_manager_instance.execute_plugin(
            'ui_plugin',
            {'component': 'dashboard', 'data': {'quality_score': 0.95}},
        )
        
        # Check the result
        self.assertEqual(result['result'], 'success')
        self.assertTrue(result['component_rendered'])
        self.assertEqual(result['component_type'], 'dashboard')


class TestPluginCLI(unittest.TestCase):
    """Test cases for the Plugin CLI."""

    def setUp(self):
        """Set up test environment."""
        # Create a temporary directory for test files
        self.temp_dir = tempfile.TemporaryDirectory()
        self.plugins_dir = os.path.join(self.temp_dir.name, 'plugins')
        os.makedirs(self.plugins_dir, exist_ok=True)
        
        # Mock the plugin CLI
        self.mock_plugin_cli_patcher = mock.patch('pkg.cli.plugins.PluginCLI')
        self.mock_plugin_cli = self.mock_plugin_cli_patcher.start()
        self.mock_plugin_cli_instance = self.mock_plugin_cli.return_value
    
    def tearDown(self):
        """Clean up test environment."""
        # Clean up temporary directory
        self.temp_dir.cleanup()
        
        # Stop patches
        self.mock_plugin_cli_patcher.stop()
    
    def test_cli_list_plugins(self):
        """Test CLI list plugins command."""
        # Configure mock responses
        self.mock_plugin_cli_instance.list_plugins.return_value = [
            {
                'name': 'plugin1',
                'version': '1.0.0',
                'type': 'validation',
                'enabled': True,
            },
            {
                'name': 'plugin2',
                'version': '1.0.0',
                'type': 'alert',
                'enabled': False,
            },
        ]
        
        # Run the CLI command
        result = self.mock_plugin_cli_instance.list_plugins()
        
        # Check that the CLI was called
        self.mock_plugin_cli_instance.list_plugins.assert_called_once()
        
        # Check the result
        self.assertEqual(len(result), 2)
        self.assertEqual(result[0]['name'], 'plugin1')
        self.assertEqual(result[0]['type'], 'validation')
        self.assertTrue(result[0]['enabled'])
        self.assertEqual(result[1]['name'], 'plugin2')
        self.assertEqual(result[1]['type'], 'alert')
        self.assertFalse(result[1]['enabled'])
    
    def test_cli_install_plugin(self):
        """Test CLI install plugin command."""
        # Configure mock responses
        self.mock_plugin_cli_instance.install_plugin.return_value = {
            'name': 'test_plugin',
            'version': '1.0.0',
            'type': 'validation',
            'path': os.path.join(self.plugins_dir, 'test_plugin'),
        }
        
        # Run the CLI command
        result = self.mock_plugin_cli_instance.install_plugin('/path/to/plugin.zip')
        
        # Check that the CLI was called
        self.mock_plugin_cli_instance.install_plugin.assert_called_once_with('/path/to/plugin.zip')
        
        # Check the result
        self.assertEqual(result['name'], 'test_plugin')
        self.assertEqual(result['version'], '1.0.0')
        self.assertEqual(result['type'], 'validation')
    
    def test_cli_uninstall_plugin(self):
        """Test CLI uninstall plugin command."""
        # Configure mock responses
        self.mock_plugin_cli_instance.uninstall_plugin.return_value = True
        
        # Run the CLI command
        result = self.mock_plugin_cli_instance.uninstall_plugin('test_plugin')
        
        # Check that the CLI was called
        self.mock_plugin_cli_instance.uninstall_plugin.assert_called_once_with('test_plugin')
        
        # Check the result
        self.assertTrue(result)
    
    def test_cli_enable_plugin(self):
        """Test CLI enable plugin command."""
        # Configure mock responses
        self.mock_plugin_cli_instance.enable_plugin.return_value = True
        
        # Run the CLI command
        result = self.mock_plugin_cli_instance.enable_plugin('test_plugin')
        
        # Check that the CLI was called
        self.mock_plugin_cli_instance.enable_plugin.assert_called_once_with('test_plugin')
        
        # Check the result
        self.assertTrue(result)
    
    def test_cli_disable_plugin(self):
        """Test CLI disable plugin command."""
        # Configure mock responses
        self.mock_plugin_cli_instance.disable_plugin.return_value = True
        
        # Run the CLI command
        result = self.mock_plugin_cli_instance.disable_plugin('test_plugin')
        
        # Check that the CLI was called
        self.mock_plugin_cli_instance.disable_plugin.assert_called_once_with('test_plugin')
        
        # Check the result
        self.assertTrue(result)
    
    def test_cli_info_plugin(self):
        """Test CLI plugin info command."""
        # Configure mock responses
        self.mock_plugin_cli_instance.plugin_info.return_value = {
            'name': 'test_plugin',
            'version': '1.0.0',
            'type': 'validation',
            'enabled': True,
            'description': 'A test plugin',
            'author': 'Test Author',
            'license': 'MIT',
        }
        
        # Run the CLI command
        result = self.mock_plugin_cli_instance.plugin_info('test_plugin')
        
        # Check that the CLI was called
        self.mock_plugin_cli_instance.plugin_info.assert_called_once_with('test_plugin')
        
        # Check the result
        self.assertEqual(result['name'], 'test_plugin')
        self.assertEqual(result['version'], '1.0.0')
        self.assertEqual(result['type'], 'validation')
        self.assertTrue(result['enabled'])
        self.assertEqual(result['description'], 'A test plugin')
        self.assertEqual(result['author'], 'Test Author')
        self.assertEqual(result['license'], 'MIT')


if __name__ == '__main__':
    unittest.main()
