package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigExportImport(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test viper instance
	v := viper.New()
	v.Set("server.host", "localhost")
	v.Set("server.port", 8080)
	v.Set("server.tls", false)
	v.Set("extensions.enabled", []string{"ext1", "ext2"})
	v.Set("extensions.path", "/path/to/extensions")
	v.Set("output.format", "json")
	v.Set("output.path", "/path/to/output")
	v.Set("batch.parallel", 4)
	v.Set("batch.output", "batch_results")
	v.Set("logging.level", "info")
	v.Set("logging.file", "/path/to/logs")
	v.Set("logging.stdout", true)

	// Create test config files
	jsonFile := filepath.Join(tempDir, "config.json")

	// Test exporting to JSON
	jsonData, err := json.MarshalIndent(v.AllSettings(), "", "  ")
	require.NoError(t, err)
	err = os.WriteFile(jsonFile, jsonData, 0644)
	require.NoError(t, err)

	// Test importing from JSON
	newConfig := viper.New()
	newConfig.SetConfigFile(jsonFile)
	err = newConfig.ReadInConfig()
	require.NoError(t, err)

	// Verify imported settings
	assert.Equal(t, "localhost", newConfig.GetString("server.host"))
	assert.Equal(t, 8080, newConfig.GetInt("server.port"))
	assert.Equal(t, false, newConfig.GetBool("server.tls"))
	assert.Equal(t, []interface{}{"ext1", "ext2"}, newConfig.Get("extensions.enabled"))
	assert.Equal(t, "/path/to/extensions", newConfig.GetString("extensions.path"))
	assert.Equal(t, "json", newConfig.GetString("output.format"))
	assert.Equal(t, "/path/to/output", newConfig.GetString("output.path"))
	assert.Equal(t, 4, newConfig.GetInt("batch.parallel"))
	assert.Equal(t, "batch_results", newConfig.GetString("batch.output"))
	assert.Equal(t, "info", newConfig.GetString("logging.level"))
	assert.Equal(t, "/path/to/logs", newConfig.GetString("logging.file"))
	assert.Equal(t, true, newConfig.GetBool("logging.stdout"))
}

func TestResetConfig(t *testing.T) {
	// Save the original config
	originalConfig := viper.AllSettings()
	
	// Create a temporary directory for the config file
	tempDir, err := os.MkdirTemp("", "config_reset_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Set a temporary config file path
	configFile := filepath.Join(tempDir, ".nessi.yaml")
	viper.SetConfigFile(configFile)
	
	// Set some custom values
	viper.Set("server.host", "custom-host")
	viper.Set("server.port", 9090)
	viper.Set("logging.level", "debug")
	
	// Reset the config
	err = resetConfig()
	require.NoError(t, err)
	
	// Verify reset values
	assert.Equal(t, "localhost", viper.GetString("server.host"))
	assert.Equal(t, 8080, viper.GetInt("server.port"))
	assert.Equal(t, false, viper.GetBool("server.tls"))
	assert.Equal(t, "info", viper.GetString("logging.level"))
	
	// Restore the original config
	for k, v := range originalConfig {
		viper.Set(k, v)
	}
}

func TestPrintValue(t *testing.T) {
	// Test cases for printValue function
	testCases := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "string",
			value:    "test",
			expected: "\"test\"",
		},
		{
			name:     "int",
			value:    42,
			expected: "42",
		},
		{
			name:     "bool",
			value:    true,
			expected: "true",
		},
		{
			name:     "slice",
			value:    []interface{}{"a", "b", "c"},
			expected: "[\"a\", \"b\", \"c\"]",
		},
		{
			name:     "map",
			value:    map[string]interface{}{"key": "value"},
			expected: "map[key:value]",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			
			// Call printValue
			printValue(tc.value)
			
			// Restore stdout
			w.Close()
			os.Stdout = oldStdout
			
			// Read captured output
			var buf []byte
			buf = make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])
			
			// Verify output
			assert.Equal(t, tc.expected, output)
		})
	}
}

func TestSetConfigValue(t *testing.T) {
	// Save the original config
	originalConfig := viper.AllSettings()
	
	// Create a temporary directory for the config file
	tempDir, err := os.MkdirTemp("", "config_set_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Set a temporary config file path
	configFile := filepath.Join(tempDir, ".nessi.yaml")
	viper.SetConfigFile(configFile)
	
	// Test setting a string value
	err = setConfigValue("server.host", "test-host")
	require.NoError(t, err)
	assert.Equal(t, "test-host", viper.GetString("server.host"))
	
	// Test setting a numeric value
	err = setConfigValue("server.port", "9090")
	require.NoError(t, err)
	assert.Equal(t, "9090", viper.GetString("server.port"))
	
	// Test setting a boolean value
	err = setConfigValue("server.tls", "true")
	require.NoError(t, err)
	assert.Equal(t, "true", viper.GetString("server.tls"))
	
	// Test setting a nested value
	err = setConfigValue("logging.level", "debug")
	require.NoError(t, err)
	assert.Equal(t, "debug", viper.GetString("logging.level"))
	
	// Restore the original config
	for k, v := range originalConfig {
		viper.Set(k, v)
	}
}
