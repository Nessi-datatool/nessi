package main

import (
	"os"
	"testing"

	"github.com/nessi-dev/nessi/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginTypes(t *testing.T) {
	// Test plugin types
	assert.Equal(t, plugin.ValidationPlugin, "validation")
	assert.Equal(t, plugin.AlertPlugin, "alert")
	assert.Equal(t, plugin.MetricPlugin, "metric")
	assert.Equal(t, plugin.StoragePlugin, "storage")
	assert.Equal(t, plugin.ExportPlugin, "export")
	assert.Equal(t, plugin.UIPlugin, "ui")
}

func TestPluginManager(t *testing.T) {
	// Create a plugin manager
	manager := getPluginManager()
	require.NotNil(t, manager)

	// Test ListPlugins with no plugins
	plugins := manager.ListPlugins()
	assert.Empty(t, plugins)

	// Create a mock plugin for testing
	mockPlugin := &plugin.Plugin{
		Info: plugin.PluginInfo{
			Name:        "mock-plugin",
			Version:     "1.0.0",
			Description: "Mock plugin for testing",
			Author:      "Test Author",
			Type:        plugin.ValidationPlugin,
			Enabled:     true,
		},
		Path:   "/path/to/mock/plugin.so",
		Module: nil,
		Lookup: make(map[string]any),
	}

	// Test plugin info
	assert.Equal(t, "mock-plugin", mockPlugin.Info.Name)
	assert.Equal(t, "1.0.0", mockPlugin.Info.Version)
	assert.Equal(t, "Mock plugin for testing", mockPlugin.Info.Description)
	assert.Equal(t, "Test Author", mockPlugin.Info.Author)
	assert.Equal(t, plugin.ValidationPlugin, mockPlugin.Info.Type)
	assert.True(t, mockPlugin.Info.Enabled)
}

func TestPluginCommands(t *testing.T) {
	// Skip this test if we're not running in an environment where we can build plugins
	if os.Getenv("SKIP_PLUGIN_TESTS") == "true" {
		t.Skip("Skipping plugin tests")
	}

	// This test would normally build and load a real plugin, but for unit testing
	// we'll just verify that the commands are registered correctly

	// Verify that the plugin command is registered
	cmd, _, err := rootCmd.Find([]string{"plugin"})
	require.NoError(t, err)
	assert.Equal(t, "plugin", cmd.Name())

	// Verify that the plugin subcommands are registered
	subcommands := []string{
		"list",
		"install",
		"uninstall",
		"info",
		"enable",
		"disable",
		"types",
		"exec",
	}

	for _, subcommand := range subcommands {
		cmd, _, err := rootCmd.Find([]string{"plugin", subcommand})
		require.NoError(t, err)
		assert.Equal(t, subcommand, cmd.Name())
	}
}
