package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the configuration directory for Nessi
// This is a shared utility function used by multiple commands
func GetConfigDir() string {
	// Try to get platform-specific config directory first
	var configDir string
	
	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "nessi")
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "nessi")
	default: // linux and others
		configDir = filepath.Join(os.Getenv("HOME"), ".config", "nessi")
	}
	
	// Fallback to home directory if needed
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		configDir = filepath.Join(homeDir, ".nessi")
	}
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}
	
	return configDir
}
