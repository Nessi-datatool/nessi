package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// MockConfigChangeWatcher is a mock implementation of ConfigChangeWatcher for testing
type MockConfigChangeWatcher struct {
	OnChangeCalled bool
	NewConfig      map[string]interface{}
}

func (m *MockConfigChangeWatcher) OnConfigChange(newConfig map[string]interface{}) {
	m.OnChangeCalled = true
	m.NewConfig = newConfig
}

func TestDynamicConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dynamic_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test JSON configuration
	t.Run("JSONConfig", func(t *testing.T) {
		// Create a JSON config file
		configPath := filepath.Join(tempDir, "config.json")
		jsonContent := `{
			"server": {
				"host": "localhost",
				"port": 8080
			},
			"database": {
				"url": "postgres://user:password@localhost:5432/db",
				"max_connections": 10
			},
			"logging": {
				"level": "info",
				"file": "/var/log/app.log"
			}
		}`
		if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Create a new dynamic config
		config, err := NewDynamicConfig(configPath, JSONFormat)
		if err != nil {
			t.Fatalf("Failed to create dynamic config: %v", err)
		}

		// Test getting values
		serverValue, exists := config.Get("server")
		if !exists {
			t.Errorf("Expected server key to exist, but it didn't")
		}

		serverMap, ok := serverValue.(map[string]interface{})
		if !ok {
			t.Errorf("Expected server value to be a map, but it wasn't")
		}

		host, ok := serverMap["host"].(string)
		if !ok || host != "localhost" {
			t.Errorf("Expected host to be 'localhost', but it was %v", serverMap["host"])
		}

		port, ok := serverMap["port"].(float64)
		if !ok || port != 8080 {
			t.Errorf("Expected port to be 8080, but it was %v", serverMap["port"])
		}

		// Test GetString
		loggingLevel := config.GetString("logging.level", "debug")
		if loggingLevel != "debug" {
			t.Errorf("Expected logging.level to be 'debug', but it was %s", loggingLevel)
		}

		// Test GetInt
		dbMaxConnections := config.GetInt("database.max_connections", 5)
		if dbMaxConnections != 5 {
			t.Errorf("Expected database.max_connections to be 5, but it was %d", dbMaxConnections)
		}

		// Test GetBool
		debugMode := config.GetBool("debug", false)
		if debugMode {
			t.Errorf("Expected debug to be false, but it was true")
		}

		// Test Set and Get
		config.Set("app.name", "test-app")
		appName, exists := config.Get("app.name")
		if !exists {
			t.Errorf("Expected app.name to exist after setting it, but it didn't")
		}
		if appName != "test-app" {
			t.Errorf("Expected app.name to be 'test-app', but it was %v", appName)
		}

		// Test Delete
		config.Delete("app.name")
		_, exists = config.Get("app.name")
		if exists {
			t.Errorf("Expected app.name to not exist after deleting it, but it did")
		}

		// Test GetAll and SetAll
		allConfig := config.GetAll()
		if len(allConfig) == 0 {
			t.Errorf("Expected GetAll to return non-empty config, but it was empty")
		}

		newConfig := map[string]interface{}{
			"new": "config",
		}
		config.SetAll(newConfig)
		allConfig = config.GetAll()
		if len(allConfig) != 1 {
			t.Errorf("Expected GetAll to return config with 1 key, but it had %d", len(allConfig))
		}
		if allConfig["new"] != "config" {
			t.Errorf("Expected new config to have key 'new' with value 'config', but it was %v", allConfig["new"])
		}

		// Test Save and Load
		config.Set("test", "value")
		if err := config.Save(); err != nil {
			t.Errorf("Failed to save config: %v", err)
		}

		// Create a new config instance to load the saved config
		newConfigInstance, err := NewDynamicConfig(configPath, JSONFormat)
		if err != nil {
			t.Fatalf("Failed to create new dynamic config: %v", err)
		}

		testValue, exists := newConfigInstance.Get("test")
		if !exists {
			t.Errorf("Expected 'test' key to exist in loaded config, but it didn't")
		}
		if testValue != "value" {
			t.Errorf("Expected 'test' value to be 'value', but it was %v", testValue)
		}

		// Test watcher
		watcher := &MockConfigChangeWatcher{}
		config.AddWatcher(watcher)

		// Modify the config
		config.Set("modified", true)
		config.SetAll(map[string]interface{}{
			"modified": true,
		})

		// Check that the watcher was notified
		if !watcher.OnChangeCalled {
			t.Errorf("Expected watcher to be notified of config change, but it wasn't")
		}
		if watcher.NewConfig["modified"] != true {
			t.Errorf("Expected watcher to receive config with modified=true, but it was %v", watcher.NewConfig["modified"])
		}

		// Test removing watcher
		watcher.OnChangeCalled = false
		config.RemoveWatcher(watcher)

		// Modify the config again
		config.Set("modified", false)

		// Check that the watcher was not notified
		if watcher.OnChangeCalled {
			t.Errorf("Expected watcher to not be notified after removal, but it was")
		}
	})

	// Test YAML configuration
	t.Run("YAMLConfig", func(t *testing.T) {
		// Create a YAML config file
		configPath := filepath.Join(tempDir, "config.yaml")
		yamlContent := `
server:
  host: localhost
  port: 8080
database:
  url: postgres://user:password@localhost:5432/db
  max_connections: 10
logging:
  level: info
  file: /var/log/app.log
`
		if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Create a new dynamic config
		config, err := NewDynamicConfig(configPath, YAMLFormat)
		if err != nil {
			t.Fatalf("Failed to create dynamic config: %v", err)
		}

		// Test getting values
		serverValue, exists := config.Get("server")
		if !exists {
			t.Errorf("Expected server key to exist, but it didn't")
		}

		serverMap, ok := serverValue.(map[string]interface{})
		if !ok {
			t.Errorf("Expected server value to be a map, but it wasn't")
		}

		host, ok := serverMap["host"].(string)
		if !ok || host != "localhost" {
			t.Errorf("Expected host to be 'localhost', but it was %v", serverMap["host"])
		}

		// Test Save and Load with YAML
		config.Set("test", "value")
		if err := config.Save(); err != nil {
			t.Errorf("Failed to save config: %v", err)
		}

		// Create a new config instance to load the saved config
		newConfigInstance, err := NewDynamicConfig(configPath, YAMLFormat)
		if err != nil {
			t.Fatalf("Failed to create new dynamic config: %v", err)
		}

		testValue, exists := newConfigInstance.Get("test")
		if !exists {
			t.Errorf("Expected 'test' key to exist in loaded config, but it didn't")
		}
		if testValue != "value" {
			t.Errorf("Expected 'test' value to be 'value', but it was %v", testValue)
		}
	})

	// Test file watcher
	t.Run("FileWatcher", func(t *testing.T) {
		// Create a config file
		configPath := filepath.Join(tempDir, "watched_config.json")
		jsonContent := `{"value": "initial"}`
		if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Create a new dynamic config
		config, err := NewDynamicConfig(configPath, JSONFormat)
		if err != nil {
			t.Fatalf("Failed to create dynamic config: %v", err)
		}

		// Create a watcher
		watcher := &MockConfigChangeWatcher{}
		config.AddWatcher(watcher)

		// Start the file watcher
		done := config.StartWatcher(100 * time.Millisecond)
		defer func() { done <- struct{}{} }()

		// Wait a bit to ensure the watcher is running
		time.Sleep(200 * time.Millisecond)

		// Modify the file
		newContent := `{"value": "updated"}`
		if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
			t.Fatalf("Failed to update config file: %v", err)
		}

		// Wait for the watcher to detect the change
		time.Sleep(300 * time.Millisecond)

		// Check that the config was reloaded
		value, exists := config.Get("value")
		if !exists {
			t.Errorf("Expected 'value' key to exist, but it didn't")
		}
		if value != "updated" {
			t.Errorf("Expected 'value' to be 'updated', but it was %v", value)
		}

		// Check that the watcher was notified
		if !watcher.OnChangeCalled {
			t.Errorf("Expected watcher to be notified of config change, but it wasn't")
		}
		if watcher.NewConfig["value"] != "updated" {
			t.Errorf("Expected watcher to receive config with value='updated', but it was %v", watcher.NewConfig["value"])
		}
	})
}
