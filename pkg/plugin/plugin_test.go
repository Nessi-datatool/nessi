package plugin

import (
	"testing"
)

func TestPluginManager(t *testing.T) {
	// Create a new plugin manager
	manager := NewPluginManager()
	if manager == nil {
		t.Fatal("Failed to create plugin manager")
	}

	// Test ListPlugins with no plugins
	plugins := manager.ListPlugins()
	if len(plugins) != 0 {
		t.Errorf("Expected 0 plugins, got %d", len(plugins))
	}

	// Test ListPluginsByType with no plugins
	validationPlugins := manager.ListPluginsByType(ValidationPlugin)
	if len(validationPlugins) != 0 {
		t.Errorf("Expected 0 validation plugins, got %d", len(validationPlugins))
	}

	// Test GetPlugin with non-existent plugin
	_, err := manager.GetPlugin("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent plugin, got nil")
	}

	// Test UnloadPlugin with non-existent plugin
	err = manager.UnloadPlugin("non-existent")
	if err == nil {
		t.Error("Expected error when unloading non-existent plugin, got nil")
	}

	// Test EnablePlugin with non-existent plugin
	err = manager.EnablePlugin("non-existent")
	if err == nil {
		t.Error("Expected error when enabling non-existent plugin, got nil")
	}

	// Test DisablePlugin with non-existent plugin
	err = manager.DisablePlugin("non-existent")
	if err == nil {
		t.Error("Expected error when disabling non-existent plugin, got nil")
	}

	// Note: We can't test LoadPlugin in a unit test because it requires
	// a compiled shared library. This would be better tested in an integration test.
}

func TestPluginTypes(t *testing.T) {
	// Test plugin types
	if ValidationPlugin != "validation" {
		t.Errorf("Expected ValidationPlugin to be 'validation', got '%s'", ValidationPlugin)
	}

	if AlertPlugin != "alert" {
		t.Errorf("Expected AlertPlugin to be 'alert', got '%s'", AlertPlugin)
	}

	if MetricPlugin != "metric" {
		t.Errorf("Expected MetricPlugin to be 'metric', got '%s'", MetricPlugin)
	}

	if StoragePlugin != "storage" {
		t.Errorf("Expected StoragePlugin to be 'storage', got '%s'", StoragePlugin)
	}

	if ExportPlugin != "export" {
		t.Errorf("Expected ExportPlugin to be 'export', got '%s'", ExportPlugin)
	}

	if UIPlugin != "ui" {
		t.Errorf("Expected UIPlugin to be 'ui', got '%s'", UIPlugin)
	}
}

func TestMockPlugin(t *testing.T) {
	// Create a mock plugin for testing
	mockPlugin := &Plugin{
		Info: PluginInfo{
			Name:        "mock-plugin",
			Version:     "1.0.0",
			Description: "Mock plugin for testing",
			Author:      "Test Author",
			Type:        ValidationPlugin,
			Enabled:     true,
		},
		Path:   "/path/to/mock/plugin.so",
		Module: nil,
		Lookup: make(map[string]any),
	}

	// Test plugin info
	if mockPlugin.Info.Name != "mock-plugin" {
		t.Errorf("Expected plugin name to be 'mock-plugin', got '%s'", mockPlugin.Info.Name)
	}

	if mockPlugin.Info.Version != "1.0.0" {
		t.Errorf("Expected plugin version to be '1.0.0', got '%s'", mockPlugin.Info.Version)
	}

	if mockPlugin.Info.Description != "Mock plugin for testing" {
		t.Errorf("Expected plugin description to be 'Mock plugin for testing', got '%s'", mockPlugin.Info.Description)
	}

	if mockPlugin.Info.Author != "Test Author" {
		t.Errorf("Expected plugin author to be 'Test Author', got '%s'", mockPlugin.Info.Author)
	}

	if mockPlugin.Info.Type != ValidationPlugin {
		t.Errorf("Expected plugin type to be 'validation', got '%s'", mockPlugin.Info.Type)
	}

	if !mockPlugin.Info.Enabled {
		t.Error("Expected plugin to be enabled, but it is disabled")
	}

	// Test plugin path
	if mockPlugin.Path != "/path/to/mock/plugin.so" {
		t.Errorf("Expected plugin path to be '/path/to/mock/plugin.so', got '%s'", mockPlugin.Path)
	}

	// Note: We can't test LookupSymbol or ExecuteFunc in a unit test because they
	// require a real plugin module. These would be better tested in an integration test.
}
