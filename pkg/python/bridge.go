package python

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nessi-dev/nessi/pkg/config"
)

var (
	// Global instance of the Python bridge
	bridge     *PythonBridge
	bridgeMu   sync.Mutex
	initialized bool
)

// PythonBridge provides a bridge between Go and Python
type PythonBridge struct {
	config         *config.PythonExtensionsConfig
	pythonPath     string
	scriptDir      string
	initialized    bool
	initializeMu   sync.Mutex
	featureFlagMgr *config.FeatureFlagManager
}

// NewPythonBridge creates a new PythonBridge
func NewPythonBridge(cfg *config.PythonExtensionsConfig, featureFlagMgr *config.FeatureFlagManager) *PythonBridge {
	return &PythonBridge{
		config:         cfg,
		pythonPath:     cfg.PythonPath,
		scriptDir:      filepath.Join(".", "scripts", "python"),
		initialized:    false,
		featureFlagMgr: featureFlagMgr,
	}
}

// GetPythonBridge returns the global PythonBridge instance
func GetPythonBridge() *PythonBridge {
	bridgeMu.Lock()
	defer bridgeMu.Unlock()

	if bridge == nil {
		// Create a new bridge with default config
		cfg := config.NewDefaultPythonExtensionsConfig()
		featureFlagMgr := config.NewFeatureFlagManager()
		bridge = NewPythonBridge(cfg, featureFlagMgr)
	}

	return bridge
}

// Initialize initializes the Python bridge
func (b *PythonBridge) Initialize() error {
	b.initializeMu.Lock()
	defer b.initializeMu.Unlock()

	if b.initialized {
		return nil
	}

	// Check if Python is available
	if err := b.checkPythonAvailable(); err != nil {
		return fmt.Errorf("Python is not available: %v", err)
	}

	// Check if required packages are installed
	if err := b.checkRequiredPackages(); err != nil {
		return fmt.Errorf("Required Python packages are not installed: %v", err)
	}

	// Create script directory if it doesn't exist
	if err := os.MkdirAll(b.scriptDir, 0755); err != nil {
		return fmt.Errorf("Failed to create script directory: %v", err)
	}

	b.initialized = true
	return nil
}

// IsInitialized returns whether the bridge is initialized
func (b *PythonBridge) IsInitialized() bool {
	b.initializeMu.Lock()
	defer b.initializeMu.Unlock()
	return b.initialized
}

// checkPythonAvailable checks if Python is available
func (b *PythonBridge) checkPythonAvailable() error {
	cmd := exec.Command(b.pythonPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Python is not available: %v", err)
	}
	return nil
}

// checkRequiredPackages checks if required Python packages are installed
func (b *PythonBridge) checkRequiredPackages() error {
	if !b.featureFlagMgr.IsEnabled(config.FeatureFlagPythonExtensions) {
		// Skip check if Python extensions are disabled
		return nil
	}

	for _, pkg := range b.config.RequiredPackages {
		cmd := exec.Command(b.pythonPath, "-c", fmt.Sprintf("import %s", pkg))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Required package '%s' is not installed", pkg)
		}
	}
	return nil
}

// ExecuteScript executes a Python script with the given arguments
func (b *PythonBridge) ExecuteScript(scriptName string, args map[string]interface{}) (string, error) {
	if !b.featureFlagMgr.IsEnabled(config.FeatureFlagPythonExtensions) {
		return "", fmt.Errorf("Python extensions are disabled")
	}

	if !b.IsInitialized() {
		if err := b.Initialize(); err != nil {
			return "", fmt.Errorf("Failed to initialize Python bridge: %v", err)
		}
	}

	scriptPath := filepath.Join(b.scriptDir, scriptName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Script '%s' not found", scriptName)
	}

	// Convert args to JSON
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal arguments: %v", err)
	}

	// Execute the script
	cmd := exec.Command(b.pythonPath, scriptPath, string(argsJSON))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Failed to execute script: %v\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// IsPythonFeatureEnabled checks if a specific Python feature is enabled
func (b *PythonBridge) IsPythonFeatureEnabled(feature config.FeatureFlag) bool {
	// First check if Python extensions are enabled globally
	if !b.featureFlagMgr.IsEnabled(config.FeatureFlagPythonExtensions) {
		return false
	}

	// Then check if the specific feature is enabled
	return b.featureFlagMgr.IsEnabled(feature)
}

// GetPythonVersion returns the Python version
func (b *PythonBridge) GetPythonVersion() (string, error) {
	cmd := exec.Command(b.pythonPath, "--version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Failed to get Python version: %v", err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// GetInstalledPackages returns a list of installed Python packages
func (b *PythonBridge) GetInstalledPackages() (map[string]string, error) {
	cmd := exec.Command(b.pythonPath, "-m", "pip", "list", "--format=json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Failed to get installed packages: %v", err)
	}

	var packages []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &packages); err != nil {
		return nil, fmt.Errorf("Failed to parse package list: %v", err)
	}

	result := make(map[string]string)
	for _, pkg := range packages {
		result[pkg.Name] = pkg.Version
	}
	return result, nil
}
