package plugin

import (
	"fmt"
	"plugin"
	"sync"
)

// PluginType represents the type of plugin
type PluginType string

const (
	// ValidationPlugin is a plugin that provides custom validation rules
	ValidationPlugin PluginType = "validation"
	// AlertPlugin is a plugin that provides custom alerting mechanisms
	AlertPlugin PluginType = "alert"
	// MetricPlugin is a plugin that provides custom metrics collection
	MetricPlugin PluginType = "metric"
	// StoragePlugin is a plugin that provides custom storage backends
	StoragePlugin PluginType = "storage"
	// ExportPlugin is a plugin that provides custom export formats
	ExportPlugin PluginType = "export"
	// UIPlugin is a plugin that provides custom UI components
	UIPlugin PluginType = "ui"
)

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Type        PluginType `json:"type"`
	Enabled     bool       `json:"enabled"`
}

// Plugin represents a loaded plugin
type Plugin struct {
	Info   PluginInfo     `json:"info"`
	Path   string         `json:"path"`
	Module *plugin.Plugin `json:"-"`
	Lookup map[string]any `json:"-"`
}

// PluginManager manages the loading and execution of plugins
type PluginManager struct {
	plugins map[string]*Plugin
	mu      sync.RWMutex
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*Plugin),
	}
}

// LoadPlugin loads a plugin from a shared library file
func (m *PluginManager) LoadPlugin(path string) (*Plugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %v", err)
	}

	// Get plugin info
	infoSymbol, err := p.Lookup("PluginInfo")
	if err != nil {
		return nil, fmt.Errorf("plugin does not export PluginInfo: %v", err)
	}

	info, ok := infoSymbol.(*PluginInfo)
	if !ok {
		return nil, fmt.Errorf("plugin info is not of type *PluginInfo")
	}

	// Check if plugin with the same name is already loaded
	if _, exists := m.plugins[info.Name]; exists {
		return nil, fmt.Errorf("plugin with name %s is already loaded", info.Name)
	}

	// Create plugin instance
	plugin := &Plugin{
		Info:   *info,
		Path:   path,
		Module: p,
		Lookup: make(map[string]any),
	}

	// Store the plugin
	m.plugins[info.Name] = plugin

	return plugin, nil
}

// UnloadPlugin unloads a plugin
func (m *PluginManager) UnloadPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[name]; !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	delete(m.plugins, name)
	return nil
}

// GetPlugin returns a loaded plugin by name
func (m *PluginManager) GetPlugin(name string) (*Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s is not loaded", name)
	}

	return plugin, nil
}

// ListPlugins returns a list of all loaded plugins
func (m *PluginManager) ListPlugins() []*Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]*Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}

// ListPluginsByType returns a list of loaded plugins of a specific type
func (m *PluginManager) ListPluginsByType(pluginType PluginType) []*Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]*Plugin, 0)
	for _, p := range m.plugins {
		if p.Info.Type == pluginType {
			plugins = append(plugins, p)
		}
	}

	return plugins
}

// EnablePlugin enables a plugin
func (m *PluginManager) EnablePlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	plugin.Info.Enabled = true
	return nil
}

// DisablePlugin disables a plugin
func (m *PluginManager) DisablePlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	plugin.Info.Enabled = false
	return nil
}

// LookupSymbol looks up a symbol in a plugin
func (p *Plugin) LookupSymbol(name string) (any, error) {
	// Check if we've already looked up this symbol
	if sym, exists := p.Lookup[name]; exists {
		return sym, nil
	}

	// Look up the symbol
	sym, err := p.Module.Lookup(name)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup symbol %s: %v", name, err)
	}

	// Cache the symbol
	p.Lookup[name] = sym

	return sym, nil
}

// ExecuteFunc executes a function in a plugin
func (p *Plugin) ExecuteFunc(name string, args ...any) ([]any, error) {
	// Look up the function
	symIface, err := p.LookupSymbol(name)
	if err != nil {
		return nil, err
	}

	// Check if it's a function
	fn, ok := symIface.(func(...any) ([]any, error))
	if !ok {
		return nil, fmt.Errorf("symbol %s is not a function", name)
	}

	// Execute the function
	return fn(args...)
}

// ValidationPluginInterface defines the interface for validation plugins
type ValidationPluginInterface interface {
	// Validate validates data against custom rules
	Validate(data any, options map[string]any) (bool, string, error)
}

// AlertPluginInterface defines the interface for alert plugins
type AlertPluginInterface interface {
	// SendAlert sends an alert to a custom destination
	SendAlert(alert any, options map[string]any) error
}

// MetricPluginInterface defines the interface for metric plugins
type MetricPluginInterface interface {
	// CollectMetrics collects custom metrics
	CollectMetrics(options map[string]any) (map[string]any, error)
}

// StoragePluginInterface defines the interface for storage plugins
type StoragePluginInterface interface {
	// Store stores data in a custom backend
	Store(key string, data any, options map[string]any) error
	// Retrieve retrieves data from a custom backend
	Retrieve(key string, options map[string]any) (any, error)
}

// ExportPluginInterface defines the interface for export plugins
type ExportPluginInterface interface {
	// Export exports data to a custom format
	Export(data any, options map[string]any) ([]byte, error)
}

// UIPluginInterface defines the interface for UI plugins
type UIPluginInterface interface {
	// GetComponent returns a UI component
	GetComponent(name string, options map[string]any) (any, error)
}
