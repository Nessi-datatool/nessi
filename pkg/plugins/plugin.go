package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
)

// Plugin represents a Nessi plugin with metadata and functionality
type Plugin struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Commands    []PluginCommand   `json:"commands"`
	Hooks       map[string]string `json:"hooks"`
	Path        string            `json:"path"`
	Instance    interface{}       `json:"-"`
}

// PluginCommand represents a command provided by a plugin
type PluginCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

// PluginMetadata represents the metadata of a plugin
type PluginMetadata struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	Commands     []PluginCommand   `json:"commands"`
	Hooks        map[string]string `json:"hooks"`
	Capabilities []string          `json:"capabilities,omitempty"`
}

// PluginManager handles plugin discovery, loading, and execution
type PluginManager struct {
	Plugins     map[string]*Plugin
	PluginsDir  string
	Enabled     bool
	Initialized bool
	debugMode   bool
	verbose     bool
}

// PluginInterface defines the basic interface that all plugins must implement
type PluginInterface interface {
	Initialize() error
}

// CommandExecutor defines the interface for plugins that provide commands
type CommandExecutor interface {
	ExecuteCommand(cmd string, args []string) error
}

// QualityRuleProvider defines the interface for plugins that provide custom quality rules
type QualityRuleProvider interface {
	GetCustomRules() []Rule
}

// Rule represents a custom quality rule
type Rule struct {
	Name        string
	Description string
	Validator   func(interface{}, map[string]interface{}) (bool, string)
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginsDir string) *PluginManager {
	return &PluginManager{
		Plugins:    make(map[string]*Plugin),
		PluginsDir: pluginsDir,
		Enabled:    true,
	}
}

// Initialize discovers and loads all available plugins
func (pm *PluginManager) Initialize() error {
	if pm.Initialized {
		return nil
	}

	if !pm.Enabled {
		return nil
	}

	// Create plugins directory if it doesn't exist
	if _, err := os.Stat(pm.PluginsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(pm.PluginsDir, 0755); err != nil {
			return fmt.Errorf("failed to create plugins directory: %w", err)
		}
	}

	// Discover plugins
	files, err := os.ReadDir(pm.PluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	// Load each plugin
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Check if it's a plugin file
		if !strings.HasSuffix(file.Name(), ".so") {
			continue
		}

		// Load the plugin
		pluginPath := filepath.Join(pm.PluginsDir, file.Name())
		metadataPath := strings.TrimSuffix(pluginPath, ".so") + ".json"

		// Check if metadata exists
		if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
			fmt.Printf("Warning: Plugin metadata not found for %s\n", file.Name())
			continue
		}

		// Read metadata
		metadataBytes, err := os.ReadFile(metadataPath)
		if err != nil {
			fmt.Printf("Warning: Failed to read plugin metadata for %s: %v\n", file.Name(), err)
			continue
		}

		// Parse metadata
		var p Plugin
		if err := json.Unmarshal(metadataBytes, &p); err != nil {
			fmt.Printf("Warning: Failed to parse plugin metadata for %s: %v\n", file.Name(), err)
			continue
		}

		// Set plugin path
		p.Path = pluginPath

		// Load the plugin
		if err := pm.loadPlugin(&p); err != nil {
			fmt.Printf("Warning: Failed to load plugin %s: %v\n", p.Name, err)
			continue
		}

		// Add to plugins map
		pm.Plugins[p.Name] = &p
	}

	pm.Initialized = true
	return nil
}

// loadPlugin loads a plugin from its file path
func (pm *PluginManager) loadPlugin(p *Plugin) error {
	// Open the plugin
	plug, err := plugin.Open(p.Path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the Plugin symbol
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'Plugin' symbol: %w", err)
	}

	// Store the plugin instance
	p.Instance = sym
	return nil
}

// GetPlugin returns a plugin by name
func (pm *PluginManager) GetPlugin(name string) (*Plugin, bool) {
	plugin, exists := pm.Plugins[name]
	return plugin, exists
}

// ListPlugins returns a list of all loaded plugins
func (pm *PluginManager) ListPlugins() []*Plugin {
	plugins := make([]*Plugin, 0, len(pm.Plugins))
	for _, p := range pm.Plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// ExecuteHook executes a plugin hook
func (pm *PluginManager) ExecuteHook(hookName string, args ...interface{}) error {
	if !pm.Enabled || !pm.Initialized {
		return nil
	}

	for _, p := range pm.Plugins {
		if methodName, exists := p.Hooks[hookName]; exists {
			// Get the method
			plug, ok := p.Instance.(*plugin.Plugin)
			if !ok {
				return fmt.Errorf("plugin instance is not a *plugin.Plugin")
			}

			method, err := plug.Lookup(methodName)
			if err != nil {
				return fmt.Errorf("failed to lookup hook method %s in plugin %s: %w", methodName, p.Name, err)
			}

			// Call the method
			if err := callMethod(method, args...); err != nil {
				return fmt.Errorf("failed to execute hook %s in plugin %s: %w", hookName, p.Name, err)
			}
		}
	}

	return nil
}

// SetDebugMode sets the debug mode
func (pm *PluginManager) SetDebugMode(debug bool) {
	pm.debugMode = debug
}

// SetVerbose sets the verbose mode
func (pm *PluginManager) SetVerbose(verbose bool) {
	pm.verbose = verbose
}

// ExecuteCommand executes a command on a plugin
func (pm *PluginManager) ExecuteCommand(plugin *Plugin, cmd string, args []string) error {
	// Check if the plugin implements the CommandExecutor interface
	executor, ok := plugin.Instance.(CommandExecutor)
	if !ok {
		return fmt.Errorf("plugin does not implement CommandExecutor interface")
	}

	// Execute the command
	return executor.ExecuteCommand(cmd, args)
}

// GetCustomRules returns all custom quality rules from all plugins
func (pm *PluginManager) GetCustomRules() []Rule {
	rules := []Rule{}

	// Iterate through all plugins
	for _, p := range pm.Plugins {
		// Check if the plugin implements the QualityRuleProvider interface
		provider, ok := p.Instance.(QualityRuleProvider)
		if !ok {
			continue
		}

		// Get the custom rules from the plugin
		pluginRules := provider.GetCustomRules()
		rules = append(rules, pluginRules...)
	}

	return rules
}

// callMethod calls a plugin method with the given arguments
func callMethod(method interface{}, args ...interface{}) error {
	// Use reflection to call the method with the correct arguments
	methodValue := reflect.ValueOf(method)
	if !methodValue.IsValid() {
		return fmt.Errorf("invalid method")
	}

	// Check if the method is callable
	if methodValue.Kind() != reflect.Func {
		return fmt.Errorf("method is not a function")
	}

	// Prepare arguments
	methodType := methodValue.Type()
	numIn := methodType.NumIn()
	if len(args) < numIn {
		return fmt.Errorf("not enough arguments: expected %d, got %d", numIn, len(args))
	}

	// Convert arguments to the correct types
	callArgs := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		argValue := reflect.ValueOf(args[i])
		paramType := methodType.In(i)

		// Check if the argument can be converted to the parameter type
		if !argValue.Type().AssignableTo(paramType) {
			return fmt.Errorf("argument %d has wrong type: expected %v, got %v", i, paramType, argValue.Type())
		}

		callArgs[i] = argValue
	}

	// Call the method
	retValues := methodValue.Call(callArgs)

	// Check for error return value
	if methodType.NumOut() > 0 && methodType.Out(methodType.NumOut()-1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		last := retValues[len(retValues)-1].Interface()
		if last != nil {
			return last.(error)
		}
	}

	return nil
}
