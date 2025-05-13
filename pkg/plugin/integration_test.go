package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginManagerIntegration(t *testing.T) {
	// Create a plugin manager
	manager := NewPluginManager()
	require.NotNil(t, manager)

	// Test ListPlugins with no plugins
	plugins := manager.ListPlugins()
	assert.Empty(t, plugins)

	// Test ListPluginsByType with no plugins
	validationPlugins := manager.ListPluginsByType(ValidationPlugin)
	assert.Empty(t, validationPlugins)

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
	assert.Equal(t, "mock-plugin", mockPlugin.Info.Name)
	assert.Equal(t, "1.0.0", mockPlugin.Info.Version)
	assert.Equal(t, "Mock plugin for testing", mockPlugin.Info.Description)
	assert.Equal(t, "Test Author", mockPlugin.Info.Author)
	assert.Equal(t, ValidationPlugin, mockPlugin.Info.Type)
	assert.True(t, mockPlugin.Info.Enabled)
	assert.Equal(t, "/path/to/mock/plugin.so", mockPlugin.Path)

	// Test plugin types
	assert.Equal(t, string(ValidationPlugin), "validation")
	assert.Equal(t, string(AlertPlugin), "alert")
	assert.Equal(t, string(MetricPlugin), "metric")
	assert.Equal(t, string(StoragePlugin), "storage")
	assert.Equal(t, string(ExportPlugin), "export")
	assert.Equal(t, string(UIPlugin), "ui")

	// Test plugin interfaces
	var _ ValidationPluginInterface = (*mockValidationPlugin)(nil)
	var _ AlertPluginInterface = (*mockAlertPlugin)(nil)
	var _ MetricPluginInterface = (*mockMetricPlugin)(nil)
	var _ StoragePluginInterface = (*mockStoragePlugin)(nil)
	var _ ExportPluginInterface = (*mockExportPlugin)(nil)
	var _ UIPluginInterface = (*mockUIPlugin)(nil)
}

// Mock implementations of plugin interfaces for testing
type mockValidationPlugin struct{}

func (m *mockValidationPlugin) Validate(data any, options map[string]any) (bool, string, error) {
	return true, "Valid", nil
}

type mockAlertPlugin struct{}

func (m *mockAlertPlugin) SendAlert(alert any, options map[string]any) error {
	return nil
}

type mockMetricPlugin struct{}

func (m *mockMetricPlugin) CollectMetrics(options map[string]any) (map[string]any, error) {
	return map[string]any{"metric": 42}, nil
}

type mockStoragePlugin struct{}

func (m *mockStoragePlugin) Store(key string, data any, options map[string]any) error {
	return nil
}

func (m *mockStoragePlugin) Retrieve(key string, options map[string]any) (any, error) {
	return nil, nil
}

type mockExportPlugin struct{}

func (m *mockExportPlugin) Export(data any, options map[string]any) ([]byte, error) {
	return []byte("exported data"), nil
}

type mockUIPlugin struct{}

func (m *mockUIPlugin) GetComponent(name string, options map[string]any) (any, error) {
	return nil, nil
}
