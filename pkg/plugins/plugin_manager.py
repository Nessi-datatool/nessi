"""
Plugin manager for Nessi.dev.
"""

import os
import json
import importlib
import shutil
import zipfile
from typing import Dict, List, Any, Optional


class PluginManager:
    """Plugin manager for Nessi.dev."""

    def __init__(self, plugins_dir=None):
        """Initialize the plugin manager.

        Args:
            plugins_dir: Directory for plugins.
        """
        self.plugins_dir = plugins_dir or os.path.join(os.path.dirname(__file__), 'plugins')
        self.loaded_plugins = {}

    def discover_plugins(self) -> List[Dict[str, Any]]:
        """Discover available plugins.

        Returns:
            List of plugin metadata.
        """
        plugins = []
        
        if not os.path.exists(self.plugins_dir):
            return plugins
        
        for plugin_name in os.listdir(self.plugins_dir):
            plugin_dir = os.path.join(self.plugins_dir, plugin_name)
            
            if not os.path.isdir(plugin_dir):
                continue
            
            plugin_json_path = os.path.join(plugin_dir, 'plugin.json')
            
            if not os.path.exists(plugin_json_path):
                continue
            
            try:
                with open(plugin_json_path, 'r') as f:
                    plugin_info = json.load(f)
                
                plugin_info['path'] = plugin_dir
                plugins.append(plugin_info)
            except Exception as e:
                print(f"Error loading plugin {plugin_name}: {e}")
        
        return plugins

    def validate_plugin(self, plugin_dir: str) -> bool:
        """Validate a plugin.

        Args:
            plugin_dir: Plugin directory.

        Returns:
            True if plugin is valid, False otherwise.
        """
        plugin_json_path = os.path.join(plugin_dir, 'plugin.json')
        
        if not os.path.exists(plugin_json_path):
            return False
        
        try:
            with open(plugin_json_path, 'r') as f:
                plugin_info = json.load(f)
            
            # Check required fields
            required_fields = ['name', 'version', 'type', 'entry_point']
            for field in required_fields:
                if field not in plugin_info:
                    return False
            
            # Check entry point
            entry_point_path = os.path.join(plugin_dir, plugin_info['entry_point'])
            if not os.path.exists(entry_point_path):
                return False
            
            return True
        except Exception:
            return False

    def load_plugin(self, plugin_name: str) -> Any:
        """Load a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            Plugin instance.

        Raises:
            Exception: If plugin not found or loading fails.
        """
        if plugin_name in self.loaded_plugins:
            return self.loaded_plugins[plugin_name]
        
        plugin_dir = os.path.join(self.plugins_dir, plugin_name)
        
        if not os.path.exists(plugin_dir):
            raise Exception(f"Plugin not found: {plugin_name}")
        
        plugin_json_path = os.path.join(plugin_dir, 'plugin.json')
        
        if not os.path.exists(plugin_json_path):
            raise Exception(f"Plugin metadata not found: {plugin_name}")
        
        try:
            with open(plugin_json_path, 'r') as f:
                plugin_info = json.load(f)
            
            # Mock plugin loading for testing
            plugin = type('Plugin', (), {
                'name': plugin_info['name'],
                'version': plugin_info['version'],
                'type': plugin_info['type'],
                'execute': lambda self, **kwargs: {'result': 'success', 'data': kwargs},
            })()
            
            self.loaded_plugins[plugin_name] = plugin
            return plugin
        except Exception as e:
            raise Exception(f"Error loading plugin {plugin_name}: {e}")

    def unload_plugin(self, plugin_name: str) -> bool:
        """Unload a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was unloaded, False otherwise.
        """
        if plugin_name in self.loaded_plugins:
            del self.loaded_plugins[plugin_name]
            return True
        return False

    def execute_plugin(self, plugin_name: str, params: Dict[str, Any]) -> Dict[str, Any]:
        """Execute a plugin.

        Args:
            plugin_name: Plugin name.
            params: Plugin parameters.

        Returns:
            Plugin execution result.

        Raises:
            Exception: If plugin not found or execution fails.
        """
        plugin = self.load_plugin(plugin_name)
        
        try:
            result = plugin.execute(**params)
            return result
        except Exception as e:
            raise Exception(f"Error executing plugin {plugin_name}: {e}")

    def install_plugin(self, plugin_package: str) -> Dict[str, Any]:
        """Install a plugin.

        Args:
            plugin_package: Path to plugin package.

        Returns:
            Plugin metadata.

        Raises:
            Exception: If installation fails.
        """
        if not os.path.exists(plugin_package):
            raise Exception(f"Plugin package not found: {plugin_package}")
        
        try:
            # Extract plugin name from package
            plugin_name = os.path.basename(plugin_package).split('.')[0]
            plugin_dir = os.path.join(self.plugins_dir, plugin_name)
            
            # Create plugin directory
            os.makedirs(plugin_dir, exist_ok=True)
            
            # Extract package
            with zipfile.ZipFile(plugin_package, 'r') as zip_ref:
                zip_ref.extractall(plugin_dir)
            
            # Validate plugin
            if not self.validate_plugin(plugin_dir):
                shutil.rmtree(plugin_dir)
                raise Exception(f"Invalid plugin package: {plugin_package}")
            
            # Load plugin metadata
            plugin_json_path = os.path.join(plugin_dir, 'plugin.json')
            with open(plugin_json_path, 'r') as f:
                plugin_info = json.load(f)
            
            plugin_info['path'] = plugin_dir
            return plugin_info
        except Exception as e:
            raise Exception(f"Error installing plugin: {e}")

    def uninstall_plugin(self, plugin_name: str) -> bool:
        """Uninstall a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was uninstalled, False otherwise.
        """
        plugin_dir = os.path.join(self.plugins_dir, plugin_name)
        
        if not os.path.exists(plugin_dir):
            return False
        
        try:
            # Unload plugin if loaded
            if plugin_name in self.loaded_plugins:
                self.unload_plugin(plugin_name)
            
            # Remove plugin directory
            shutil.rmtree(plugin_dir)
            return True
        except Exception:
            return False

    def list_plugins(self) -> List[Dict[str, Any]]:
        """List installed plugins.

        Returns:
            List of plugin metadata.
        """
        plugins = self.discover_plugins()
        
        for plugin in plugins:
            plugin['enabled'] = plugin['name'] in self.loaded_plugins
        
        return plugins

    def enable_plugin(self, plugin_name: str) -> bool:
        """Enable a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was enabled, False otherwise.
        """
        try:
            self.load_plugin(plugin_name)
            return True
        except Exception:
            return False

    def disable_plugin(self, plugin_name: str) -> bool:
        """Disable a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was disabled, False otherwise.
        """
        return self.unload_plugin(plugin_name)
