package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigFormat represents the format of a configuration file
type ConfigFormat string

const (
	// JSONFormat represents JSON configuration format
	JSONFormat ConfigFormat = "json"
	// YAMLFormat represents YAML configuration format
	YAMLFormat ConfigFormat = "yaml"
)

// DynamicConfig represents a configuration that can be dynamically loaded and reloaded
type DynamicConfig struct {
	data       map[string]interface{}
	filePath   string
	format     ConfigFormat
	lastLoaded time.Time
	watchers   []ConfigChangeWatcher
	mu         sync.RWMutex
}

// ConfigChangeWatcher is an interface for objects that want to be notified of config changes
type ConfigChangeWatcher interface {
	OnConfigChange(newConfig map[string]interface{})
}

// NewDynamicConfig creates a new dynamic configuration
func NewDynamicConfig(filePath string, format ConfigFormat) (*DynamicConfig, error) {
	config := &DynamicConfig{
		data:     make(map[string]interface{}),
		filePath: filePath,
		format:   format,
		watchers: []ConfigChangeWatcher{},
	}

	// Load the initial configuration
	if err := config.Load(); err != nil {
		return nil, err
	}

	return config, nil
}

// Load loads the configuration from the file
func (c *DynamicConfig) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var newConfig map[string]interface{}

	switch c.format {
	case JSONFormat:
		if err := json.Unmarshal(data, &newConfig); err != nil {
			return fmt.Errorf("failed to parse JSON config: %v", err)
		}
	case YAMLFormat:
		if err := yaml.Unmarshal(data, &newConfig); err != nil {
			return fmt.Errorf("failed to parse YAML config: %v", err)
		}
	default:
		return fmt.Errorf("unsupported config format: %s", c.format)
	}

	// Update the configuration
	c.data = newConfig
	c.lastLoaded = time.Now()

	// Notify watchers
	for _, watcher := range c.watchers {
		go watcher.OnConfigChange(c.data)
	}

	return nil
}

// Save saves the configuration to the file
func (c *DynamicConfig) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var data []byte
	var err error

	switch c.format {
	case JSONFormat:
		data, err = json.MarshalIndent(c.data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %v", err)
		}
	case YAMLFormat:
		data, err = yaml.Marshal(c.data)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %v", err)
		}
	default:
		return fmt.Errorf("unsupported config format: %s", c.format)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Write the file
	if err := ioutil.WriteFile(c.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Get gets a configuration value
func (c *DynamicConfig) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	return value, exists
}

// GetString gets a string configuration value
func (c *DynamicConfig) GetString(key string, defaultValue string) string {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	strValue, ok := value.(string)
	if !ok {
		return defaultValue
	}

	return strValue
}

// GetInt gets an integer configuration value
func (c *DynamicConfig) GetInt(key string, defaultValue int) int {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	// Handle different number types
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return defaultValue
	}
}

// GetBool gets a boolean configuration value
func (c *DynamicConfig) GetBool(key string, defaultValue bool) bool {
	value, exists := c.Get(key)
	if !exists {
		return defaultValue
	}

	boolValue, ok := value.(bool)
	if !ok {
		return defaultValue
	}

	return boolValue
}

// Set sets a configuration value
func (c *DynamicConfig) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// Delete deletes a configuration value
func (c *DynamicConfig) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// GetAll gets all configuration values
func (c *DynamicConfig) GetAll() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a copy to avoid race conditions
	config := make(map[string]interface{})
	for k, v := range c.data {
		config[k] = v
	}

	return config
}

// SetAll sets all configuration values
func (c *DynamicConfig) SetAll(data map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = data

	// Notify watchers
	for _, watcher := range c.watchers {
		go watcher.OnConfigChange(c.data)
	}
}

// AddWatcher adds a configuration change watcher
func (c *DynamicConfig) AddWatcher(watcher ConfigChangeWatcher) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.watchers = append(c.watchers, watcher)
}

// RemoveWatcher removes a configuration change watcher
func (c *DynamicConfig) RemoveWatcher(watcher ConfigChangeWatcher) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, w := range c.watchers {
		if w == watcher {
			c.watchers = append(c.watchers[:i], c.watchers[i+1:]...)
			break
		}
	}
}

// StartWatcher starts a goroutine that watches for configuration changes
func (c *DynamicConfig) StartWatcher(interval time.Duration) chan struct{} {
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Check if the file has been modified
				info, err := os.Stat(c.filePath)
				if err != nil {
					continue
				}

				c.mu.RLock()
				lastLoaded := c.lastLoaded
				c.mu.RUnlock()

				if info.ModTime().After(lastLoaded) {
					if err := c.Load(); err != nil {
						// Log the error but continue
						fmt.Printf("Error reloading config: %v\n", err)
					}
				}
			case <-done:
				return
			}
		}
	}()

	return done
}
