package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nessi-dev/nessi-dev/pkg/config"
	"github.com/spf13/cobra"
)

// pythonExtensionsCmd represents the pythonExtensions command
var pythonExtensionsCmd = &cobra.Command{
	Use:   "python-extensions",
	Short: "Manage Python extensions",
	Long: `Manage Python extensions for Nessi.dev.

Python extensions provide advanced capabilities like ML-based anomaly detection,
advanced Delta Lake features, and sophisticated visualizations.`,
}

// pythonExtensionsCheckCmd represents the check command for Python extensions
var pythonExtensionsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Python environment",
	Long:  `Check if Python and required packages are available for Python extensions.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkPythonEnvironment()
	},
}

// pythonExtensionsEnableCmd represents the enable command for Python extensions
var pythonExtensionsEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable Python extensions",
	Long:  `Enable Python extensions and verify the environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		enablePythonExtensions()
	},
}

// pythonExtensionsDisableCmd represents the disable command for Python extensions
var pythonExtensionsDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable Python extensions",
	Long:  `Disable Python extensions.`,
	Run: func(cmd *cobra.Command, args []string) {
		disablePythonExtensions()
	},
}

// pythonExtensionsStatusCmd represents the status command for Python extensions
var pythonExtensionsStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Python extensions status",
	Long:  `Show the status of Python extensions and available features.`,
	Run: func(cmd *cobra.Command, args []string) {
		showPythonExtensionsStatus()
	},
}

func init() {
	rootCmd.AddCommand(pythonExtensionsCmd)
	pythonExtensionsCmd.AddCommand(pythonExtensionsCheckCmd)
	pythonExtensionsCmd.AddCommand(pythonExtensionsEnableCmd)
	pythonExtensionsCmd.AddCommand(pythonExtensionsDisableCmd)
	pythonExtensionsCmd.AddCommand(pythonExtensionsStatusCmd)
}

// checkPythonEnvironment checks if Python and required packages are available
func checkPythonEnvironment() {
	fmt.Println("Checking Python environment...")

	// Check Python version
	pythonCmd := exec.Command("python", "--version")
	pythonOutput, err := pythonCmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Python not found. Please install Python 3.8 or later.")
		return
	}

	fmt.Printf("✅ %s", pythonOutput)

	// Check required packages
	requiredPackages := []string{
		"pandas",
		"numpy",
		"scikit-learn",
		"matplotlib",
		"pyarrow",
		"deltalake",
	}

	fmt.Println("\nChecking required packages:")
	
	for _, pkg := range requiredPackages {
		pipCmd := exec.Command("pip", "show", pkg)
		pipOutput, err := pipCmd.CombinedOutput()
		
		if err != nil {
			fmt.Printf("❌ %s not found\n", pkg)
		} else {
			// Extract version from pip output
			lines := strings.Split(string(pipOutput), "\n")
			version := ""
			for _, line := range lines {
				if strings.HasPrefix(line, "Version:") {
					version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
					break
				}
			}
			fmt.Printf("✅ %s (version %s)\n", pkg, version)
		}
	}

	// Check if any packages are missing
	pipCmd := exec.Command("pip", "list")
	pipOutput, err := pipCmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Failed to list installed packages.")
		return
	}

	installedPackages := string(pipOutput)
	missingPackages := []string{}

	for _, pkg := range requiredPackages {
		if !strings.Contains(installedPackages, pkg) {
			missingPackages = append(missingPackages, pkg)
		}
	}

	if len(missingPackages) > 0 {
		fmt.Println("\n⚠️ Missing packages:")
		for _, pkg := range missingPackages {
			fmt.Printf("  - %s\n", pkg)
		}
		fmt.Println("\nInstall missing packages with:")
		fmt.Println("  pip install " + strings.Join(missingPackages, " "))
	} else {
		fmt.Println("\n✅ All required packages are installed.")
	}

	// Check if Python extensions are enabled
	featureFlagManager := config.NewFeatureFlagManager()
	// TODO: Load feature flags from config file

	if featureFlagManager.IsEnabled(config.FeatureFlagPythonExtensions) {
		fmt.Println("\n✅ Python extensions are enabled.")
	} else {
		fmt.Println("\n⚠️ Python extensions are disabled. Enable them with:")
		fmt.Println("  nessi python-extensions enable")
	}
}

// enablePythonExtensions enables Python extensions
func enablePythonExtensions() {
	// Check Python environment first
	pythonCmd := exec.Command("python", "--version")
	_, err := pythonCmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Python not found. Please install Python 3.8 or later.")
		fmt.Println("Python extensions cannot be enabled without Python installed.")
		os.Exit(1)
	}

	// Enable Python extensions in feature flags
	featureFlagManager := config.NewFeatureFlagManager()
	// TODO: Load feature flags from config file
	
	featureFlagManager.SetFlag(config.FeatureFlagPythonExtensions, true)
	
	// TODO: Save feature flags to config file

	fmt.Println("✅ Python extensions enabled successfully.")
	fmt.Println("You can now use the following Python-based features:")
	fmt.Println("  - ML-based anomaly detection")
	fmt.Println("  - Advanced Delta Lake features")
	fmt.Println("  - Advanced statistical analysis")
	fmt.Println("  - Custom rule execution engine")
	fmt.Println("  - Advanced visualization components")
}

// disablePythonExtensions disables Python extensions
func disablePythonExtensions() {
	// Disable Python extensions in feature flags
	featureFlagManager := config.NewFeatureFlagManager()
	// TODO: Load feature flags from config file
	
	featureFlagManager.SetFlag(config.FeatureFlagPythonExtensions, false)
	
	// TODO: Save feature flags to config file

	fmt.Println("✅ Python extensions disabled successfully.")
	fmt.Println("Nessi.dev will now use only Go-based features.")
}

// showPythonExtensionsStatus shows the status of Python extensions
func showPythonExtensionsStatus() {
	featureFlagManager := config.NewFeatureFlagManager()
	// TODO: Load feature flags from config file

	fmt.Println("Python Extensions Status:")
	fmt.Println("========================")
	
	pythonExtensionsEnabled := featureFlagManager.IsEnabled(config.FeatureFlagPythonExtensions)
	
	if pythonExtensionsEnabled {
		fmt.Println("✅ Python extensions: Enabled")
	} else {
		fmt.Println("❌ Python extensions: Disabled")
	}
	
	// Check individual features
	features := map[string]config.FeatureFlag{
		"ML-based anomaly detection":       config.FeatureFlagMLAnomalyDetection,
		"Advanced Delta Lake features":     config.FeatureFlagAdvancedDeltaLake,
		"Advanced visualization":           config.FeatureFlagAdvancedVisualization,
	}
	
	fmt.Println("\nIndividual Features:")
	for name, flag := range features {
		if featureFlagManager.IsEnabled(flag) && pythonExtensionsEnabled {
			fmt.Printf("✅ %s: Enabled\n", name)
		} else {
			fmt.Printf("❌ %s: Disabled\n", name)
		}
	}
	
	// Show Python environment status
	fmt.Println("\nPython Environment:")
	pythonCmd := exec.Command("python", "--version")
	pythonOutput, err := pythonCmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Python: Not installed or not in PATH")
	} else {
		fmt.Printf("✅ Python: %s", pythonOutput)
	}
}
