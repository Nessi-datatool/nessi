package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Extensions ExtensionsConfig `mapstructure:"extensions"`
	Telemetry  TelemetryConfig  `mapstructure:"telemetry"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ExtensionsConfig represents the extensions configuration
type ExtensionsConfig struct {
	Enabled []string `mapstructure:"enabled"`
	Path    string   `mapstructure:"path"`
}

var (
	config      *Config
	viperConfig *viper.Viper
)

// InitConfig initializes the configuration
func InitConfig(configFile string) error {
	viperConfig = viper.New()

	if configFile != "" {
		viperConfig.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		// Search config in home directory with name ".nessi" (without extension)
		viperConfig.AddConfigPath(home)
		viperConfig.SetConfigName(".nessi")
	}

	// Read in environment variables that match
	viperConfig.AutomaticEnv()

	// If a config file is found, read it in
	if err := viperConfig.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viperConfig.ConfigFileUsed())
	} else {
		// Create default config if it doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(); err != nil {
				return fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config into struct
	config = &Config{}
	if err := viperConfig.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	defaultConfig := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
		"extensions": map[string]interface{}{
			"enabled": []string{},
			"path":    filepath.Join(home, ".nessi", "extensions"),
		},
	}

	// Set default values
	for k, v := range defaultConfig {
		viperConfig.Set(k, v)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(filepath.Join(home, ".nessi.yaml"))
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	if err := viperConfig.SafeWriteConfig(); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return config
}

// GetServerConfig returns the server configuration
func GetServerConfig() ServerConfig {
	return config.Server
}

// GetExtensionsConfig returns the extensions configuration
func GetExtensionsConfig() ExtensionsConfig {
	return config.Extensions
}

// GetTelemetryConfig returns the telemetry configuration
func GetTelemetryConfig() TelemetryConfig {
	return config.Telemetry
}

// UpdateConfig updates the configuration with new values
func UpdateConfig(updates map[string]interface{}) error {
	for k, v := range updates {
		viperConfig.Set(k, v)
	}

	if err := viperConfig.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Reload config
	if err := viperConfig.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal updated config: %w", err)
	}

	return nil
}

// EnableExtension adds an extension to the enabled list
func EnableExtension(name string) error {
	enabled := config.Extensions.Enabled
	for _, ext := range enabled {
		if ext == name {
			return nil // Already enabled
		}
	}

	enabled = append(enabled, name)
	return UpdateConfig(map[string]interface{}{
		"extensions.enabled": enabled,
	})
}

// DisableExtension removes an extension from the enabled list
func DisableExtension(name string) error {
	enabled := config.Extensions.Enabled
	for i, ext := range enabled {
		if ext == name {
			enabled = append(enabled[:i], enabled[i+1:]...)
			return UpdateConfig(map[string]interface{}{
				"extensions.enabled": enabled,
			})
		}
	}

	return nil // Already disabled
}

// IsExtensionEnabled checks if an extension is enabled
func IsExtensionEnabled(name string) bool {
	for _, ext := range config.Extensions.Enabled {
		if ext == name {
			return true
		}
	}
	return false
}
