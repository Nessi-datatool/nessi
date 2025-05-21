// Package plugins provides a plugin system for Nessi
package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

// This package should not import other Nessi packages to avoid circular dependencies

// PluginManager manages the loading and execution of plugins
type PluginManager struct {
	PluginsDir  string
	Plugins     map[string]*PluginInfo
	Initialized bool
	debugMode   bool
	verbose     bool
}

// PluginInfo contains information about a plugin
type PluginInfo struct {
	Name        string
	Version     string
	Description string
	Author      string
	Commands    []PluginCommand
	Hooks       map[string]string
	Path        string
	Instance    interface{}
}

// PluginCommand represents a command provided by a plugin
type PluginCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

// PluginMetadata represents the metadata of a plugin
type PluginMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Commands    []PluginCommand   `json:"commands"`
	Hooks       map[string]string `json:"hooks"`
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
		PluginsDir:  pluginsDir,
		Plugins:     make(map[string]*PluginInfo),
		Initialized: false,
		debugMode:   false,
		verbose:     false,
	}
}

// SetDebugMode sets the debug mode
func (pm *PluginManager) SetDebugMode(debug bool) {
	pm.debugMode = debug
}

// SetVerbose sets the verbose mode
func (pm *PluginManager) SetVerbose(verbose bool) {
	pm.verbose = verbose
}

// Initialize loads all plugins from the plugins directory
func (pm *PluginManager) Initialize() error {
	// Clear existing plugins
	pm.Plugins = make(map[string]*PluginInfo)

	// Check if the plugins directory exists
	if _, err := os.Stat(pm.PluginsDir); os.IsNotExist(err) {
		// Create the plugins directory
		err := os.MkdirAll(pm.PluginsDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create plugins directory: %v", err)
		}

		// No plugins to load yet
		pm.Initialized = true
		return nil
	}

	// Get all .so and .dll files in the plugins directory
	files, err := filepath.Glob(filepath.Join(pm.PluginsDir, "*.so"))
	if err != nil {
		return fmt.Errorf("failed to list plugins: %v", err)
	}

	// Add .dll files on Windows
	dllFiles, err := filepath.Glob(filepath.Join(pm.PluginsDir, "*.dll"))
	if err == nil {
		files = append(files, dllFiles...)
	}

	// Load each plugin
	for _, file := range files {
		err := pm.loadPlugin(file)
		if err != nil && pm.verbose {
			fmt.Printf("Warning: Failed to load plugin %s: %v\n", file, err)
		}
	}

	pm.Initialized = true
	return nil
}

// loadPlugin loads a plugin from a file
func (pm *PluginManager) loadPlugin(pluginPath string) error {
	// Get the plugin name from the file name
	pluginName := filepath.Base(pluginPath)
	pluginName = strings.TrimSuffix(pluginName, filepath.Ext(pluginName))

	// Check if the metadata file exists
	metadataPath := strings.TrimSuffix(pluginPath, filepath.Ext(pluginPath)) + ".json"
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return fmt.Errorf("metadata file not found: %s", metadataPath)
	}

	// Read the metadata file
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file: %v", err)
	}

	// Parse the metadata
	var metadata PluginMetadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return fmt.Errorf("failed to parse metadata: %v", err)
	}

	// Validate the metadata
	if metadata.Name == "" {
		return fmt.Errorf("plugin name is required in metadata")
	}

	if metadata.Version == "" {
		return fmt.Errorf("plugin version is required in metadata")
	}

	// Debug output
	if pm.debugMode {
		fmt.Printf("Loading plugin: %s (version %s)\n", metadata.Name, metadata.Version)
		fmt.Printf("  Description: %s\n", metadata.Description)
		fmt.Printf("  Author: %s\n", metadata.Author)
		fmt.Printf("  Commands: %d\n", len(metadata.Commands))
		fmt.Printf("  Hooks: %d\n", len(metadata.Hooks))
	}

	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %v", err)
	}

	// Look up the Plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'Plugin' symbol: %v", err)
	}

	// Check if the symbol implements the PluginInterface
	pluginInstance, ok := sym.(PluginInterface)
	if !ok {
		return fmt.Errorf("plugin does not implement PluginInterface")
	}

	// Initialize the plugin
	err = pluginInstance.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin: %v", err)
	}

	// Create the plugin info
	pluginInfo := &PluginInfo{
		Name:        metadata.Name,
		Version:     metadata.Version,
		Description: metadata.Description,
		Author:      metadata.Author,
		Commands:    metadata.Commands,
		Hooks:       metadata.Hooks,
		Path:        pluginPath,
		Instance:    pluginInstance,
	}

	// Add the plugin to the map
	pm.Plugins[metadata.Name] = pluginInfo

	// Debug output
	if pm.debugMode {
		fmt.Printf("Plugin loaded successfully: %s\n", metadata.Name)
	}

	return nil
}

// ListPlugins returns a list of all loaded plugins
func (pm *PluginManager) ListPlugins() []*PluginInfo {
	plugins := make([]*PluginInfo, 0, len(pm.Plugins))
	for _, p := range pm.Plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// GetPlugin returns a plugin by name
func (pm *PluginManager) GetPlugin(name string) (*PluginInfo, bool) {
	plugin, exists := pm.Plugins[name]
	return plugin, exists
}

// ExecuteCommand executes a command on a plugin
func (pm *PluginManager) ExecuteCommand(plugin *PluginInfo, cmd string, args []string) error {
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
