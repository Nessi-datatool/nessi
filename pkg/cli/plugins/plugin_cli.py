"""
Plugin CLI for Nessi.dev.
"""

import os
import json
from typing import Dict, List, Any, Optional
from pkg.plugins import PluginManager


class PluginCLI:
    """Plugin CLI for Nessi.dev."""

    def __init__(self, plugin_manager=None):
        """Initialize the plugin CLI.

        Args:
            plugin_manager: Plugin manager instance.
        """
        self.plugin_manager = plugin_manager or PluginManager()

    def list_plugins(self) -> List[Dict[str, Any]]:
        """List installed plugins.

        Returns:
            List of plugin metadata.
        """
        return self.plugin_manager.list_plugins()

    def install_plugin(self, plugin_package: str) -> Dict[str, Any]:
        """Install a plugin.

        Args:
            plugin_package: Path to plugin package.

        Returns:
            Plugin metadata.

        Raises:
            Exception: If installation fails.
        """
        return self.plugin_manager.install_plugin(plugin_package)

    def uninstall_plugin(self, plugin_name: str) -> bool:
        """Uninstall a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was uninstalled, False otherwise.
        """
        return self.plugin_manager.uninstall_plugin(plugin_name)

    def enable_plugin(self, plugin_name: str) -> bool:
        """Enable a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was enabled, False otherwise.
        """
        return self.plugin_manager.enable_plugin(plugin_name)

    def disable_plugin(self, plugin_name: str) -> bool:
        """Disable a plugin.

        Args:
            plugin_name: Plugin name.

        Returns:
            True if plugin was disabled, False otherwise.
        """
        return self.plugin_manager.disable_plugin(plugin_name)

    def plugin_info(self, plugin_name: str) -> Dict[str, Any]:
        """Get plugin information.

        Args:
            plugin_name: Plugin name.

        Returns:
            Plugin metadata.

        Raises:
            Exception: If plugin not found.
        """
        plugins = self.plugin_manager.list_plugins()
        
        for plugin in plugins:
            if plugin['name'] == plugin_name:
                return {
                    'name': plugin['name'],
                    'version': plugin['version'],
                    'type': plugin['type'],
                    'enabled': plugin['enabled'],
                    'description': plugin.get('description', 'No description available'),
                    'author': plugin.get('author', 'Unknown'),
                    'license': plugin.get('license', 'Unknown'),
                }
        
        raise Exception(f"Plugin not found: {plugin_name}")
