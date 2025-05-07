// Implement an extension management system with these components:
// 1. Extension interface that all extensions must implement:
//    - Start() error - Start the extension
//    - Stop() error - Stop the extension
//    - Status() string - Get current status
//    - IsRunning() bool - Check if running
//    - Name() string - Get extension name
//    - Version() string - Get extension version
//    - Dependencies() []string - List dependencies
// 2. Manager struct that maintains:
//    - Map of all registered extensions
//    - Methods to enable/disable extensions
//    - Methods to start/stop extensions
//    - Method to check status of all extensions
//    - Method to resolve dependencies between extensions
// 3. PythonExtension struct that implements Extension interface:
//    - Communicates with Python via a socket or pipe
//    - Manages Python process lifecycle
//    - Handles errors and recovery
// 4. Core Go extensions:
//    - DeltaBasicExtension - Basic Delta Lake operations
//    - QualityCoreExtension - Core data quality
//    - MonitoringExtension - Prometheus metrics
// Extensions should be startable/stoppable independently when possible
// Handle dependency resolution (don't start extensions with missing dependencies)
// Provide details about why extensions failed to start

package extensions

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/nessi-dev/nessi-dev/pkg"
)

// ExtensionInterface defines the interface that all extensions must implement
type ExtensionInterface interface {
	Start() error
	Stop() error
	Status() string
	IsRunning() bool
	Name() string
	Version() string
	Dependencies() []string
}

// Extension represents a loaded extension
type Extension struct {
	Name         string
	Version      string
	Description  string
	Plugin       *plugin.Plugin
	Enabled      bool
	Interface    ExtensionInterface
	Dependencies []string
	Status       string
	mu           sync.RWMutex
}

// Manager handles the loading and management of extensions
type Manager struct {
	extensions map[string]*Extension
	extPath    string
	mu         sync.RWMutex
}

var (
	manager *Manager
)

// InitExtensionManager initializes the extension manager
func InitExtensionManager() error {
	extConfig := pkg.GetExtensionsConfig()
	manager = &Manager{
		extensions: make(map[string]*Extension),
		extPath:    extConfig.Path,
	}

	// Create extensions directory if it doesn't exist
	if err := os.MkdirAll(manager.extPath, 0755); err != nil {
		return fmt.Errorf("failed to create extensions directory: %w", err)
	}

	// Load all extensions
	if err := manager.loadExtensions(); err != nil {
		return fmt.Errorf("failed to load extensions: %w", err)
	}

	return nil
}

// loadExtensions loads all available extensions
func (m *Manager) loadExtensions() error {
	files, err := ioutil.ReadDir(m.extPath)
	if err != nil {
		return fmt.Errorf("failed to read extensions directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			ext, err := m.loadExtension(file.Name())
			if err != nil {
				fmt.Printf("Warning: Failed to load extension %s: %v\n", file.Name(), err)
				continue
			}
			m.extensions[ext.Name] = ext
		}
	}

	return nil
}

// loadExtension loads a single extension from a .so file
func (m *Manager) loadExtension(filename string) (*Extension, error) {
	pluginPath := filepath.Join(m.extPath, filename)
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up required symbols
	nameSym, err := p.Lookup("Name")
	if err != nil {
		return nil, fmt.Errorf("extension missing Name symbol: %w", err)
	}

	versionSym, err := p.Lookup("Version")
	if err != nil {
		return nil, fmt.Errorf("extension missing Version symbol: %w", err)
	}

	descSym, err := p.Lookup("Description")
	if err != nil {
		return nil, fmt.Errorf("extension missing Description symbol: %w", err)
	}

	// Create extension
	ext := &Extension{
		Name:        *nameSym.(*string),
		Version:     *versionSym.(*string),
		Description: *descSym.(*string),
		Plugin:      p,
		Enabled:     pkg.IsExtensionEnabled(*nameSym.(*string)),
		Status:      "stopped",
	}

	// Load extension interface
	ifaceSym, err := p.Lookup("NewExtension")
	if err != nil {
		return nil, fmt.Errorf("extension missing NewExtension symbol: %w", err)
	}

	newExt := ifaceSym.(func() (ExtensionInterface, error))
	iface, err := newExt()
	if err != nil {
		return nil, fmt.Errorf("failed to create extension interface: %w", err)
	}

	ext.Interface = iface
	ext.Dependencies = iface.Dependencies()

	return ext, nil
}

// GetExtension returns an extension by name
func GetExtension(name string) (*Extension, error) {
	if manager == nil {
		return nil, fmt.Errorf("extension manager not initialized")
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	ext, exists := manager.extensions[name]
	if !exists {
		return nil, fmt.Errorf("extension %s not found", name)
	}

	return ext, nil
}

// ListExtensions returns all available extensions
func ListExtensions() ([]*Extension, error) {
	if manager == nil {
		return nil, fmt.Errorf("extension manager not initialized")
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	extensions := make([]*Extension, 0, len(manager.extensions))
	for _, ext := range manager.extensions {
		extensions = append(extensions, ext)
	}

	return extensions, nil
}

// InstallExtension installs a new extension
func InstallExtension(name string) error {
	if manager == nil {
		return fmt.Errorf("extension manager not initialized")
	}

	// TODO: Implement extension installation logic
	// This would typically involve:
	// 1. Downloading the extension from a repository
	// 2. Verifying the extension
	// 3. Installing it to the extensions directory
	// 4. Loading it into the manager

	return fmt.Errorf("extension installation not implemented")
}

// UninstallExtension removes an installed extension
func UninstallExtension(name string) error {
	if manager == nil {
		return fmt.Errorf("extension manager not initialized")
	}

	ext, exists := manager.extensions[name]
	if !exists {
		return fmt.Errorf("extension %s not found", name)
	}

	// Remove from enabled list if it was enabled
	if ext.Enabled {
		if err := pkg.DisableExtension(name); err != nil {
			return fmt.Errorf("failed to disable extension: %w", err)
		}
	}

	// Remove from manager
	delete(manager.extensions, name)

	// TODO: Implement actual file removal
	// This would involve:
	// 1. Finding the .so file
	// 2. Removing it from the filesystem
	// 3. Cleaning up any extension-specific resources

	return fmt.Errorf("extension uninstallation not implemented")
}

// StartExtension starts an extension
func StartExtension(name string) error {
	ext, err := GetExtension(name)
	if err != nil {
		return err
	}

	ext.mu.Lock()
	defer ext.mu.Unlock()

	if ext.Status == "running" {
		return nil // Already running
	}

	// Check dependencies
	for _, dep := range ext.Dependencies {
		depExt, err := GetExtension(dep)
		if err != nil {
			return fmt.Errorf("dependency %s not found: %w", dep, err)
		}
		if !depExt.IsRunning() {
			return fmt.Errorf("dependency %s is not running", dep)
		}
	}

	// Start extension
	if err := ext.Interface.Start(); err != nil {
		return fmt.Errorf("failed to start extension: %w", err)
	}

	ext.Status = "running"
	return nil
}

// StopExtension stops an extension
func StopExtension(name string) error {
	ext, err := GetExtension(name)
	if err != nil {
		return err
	}

	ext.mu.Lock()
	defer ext.mu.Unlock()

	if ext.Status != "running" {
		return nil // Not running
	}

	// Stop extension
	if err := ext.Interface.Stop(); err != nil {
		return fmt.Errorf("failed to stop extension: %w", err)
	}

	ext.Status = "stopped"
	return nil
}

// EnableExtension enables an extension
func EnableExtension(name string) error {
	if manager == nil {
		return fmt.Errorf("extension manager not initialized")
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	ext, exists := manager.extensions[name]
	if !exists {
		return fmt.Errorf("extension %s not found", name)
	}

	if ext.Enabled {
		return nil // Already enabled
	}

	if err := pkg.EnableExtension(name); err != nil {
		return fmt.Errorf("failed to enable extension: %w", err)
	}

	ext.Enabled = true
	return nil
}

// DisableExtension disables an extension
func DisableExtension(name string) error {
	if manager == nil {
		return fmt.Errorf("extension manager not initialized")
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	ext, exists := manager.extensions[name]
	if !exists {
		return fmt.Errorf("extension %s not found", name)
	}

	if !ext.Enabled {
		return nil // Already disabled
	}

	// Stop extension if running
	if ext.Status == "running" {
		if err := ext.Interface.Stop(); err != nil {
			return fmt.Errorf("failed to stop extension: %w", err)
		}
		ext.Status = "stopped"
	}

	if err := pkg.DisableExtension(name); err != nil {
		return fmt.Errorf("failed to disable extension: %w", err)
	}

	ext.Enabled = false
	return nil
}

// GetExtensionStatus returns the status of all extensions
func GetExtensionStatus() map[string]string {
	if manager == nil {
		return nil
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	status := make(map[string]string)
	for name, ext := range manager.extensions {
		ext.mu.RLock()
		status[name] = ext.Status
		ext.mu.RUnlock()
	}

	return status
}

// CheckDependencies checks if all dependencies are satisfied
func CheckDependencies(name string) error {
	ext, err := GetExtension(name)
	if err != nil {
		return err
	}

	for _, dep := range ext.Dependencies {
		depExt, err := GetExtension(dep)
		if err != nil {
			return fmt.Errorf("dependency %s not found: %w", dep, err)
		}
		if !depExt.Enabled {
			return fmt.Errorf("dependency %s is not enabled", dep)
		}
	}

	return nil
}
