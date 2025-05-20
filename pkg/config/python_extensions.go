package config

// Python extensions related feature flags
const (
	// FeatureFlagPythonExtensions enables all Python-based features
	FeatureFlagPythonExtensions FeatureFlag = "python_extensions"

	// FeatureFlagMLAnomalyDetection enables ML-based anomaly detection
	FeatureFlagMLAnomalyDetection FeatureFlag = "ml_anomaly_detection"

	// FeatureFlagAdvancedDeltaLake enables advanced Delta Lake features
	FeatureFlagAdvancedDeltaLake FeatureFlag = "advanced_delta_lake"

	// FeatureFlagAdvancedVisualization enables advanced visualization components
	FeatureFlagAdvancedVisualization FeatureFlag = "advanced_visualization"
)

// PythonExtensionsConfig represents configuration for Python extensions
type PythonExtensionsConfig struct {
	// Enabled indicates if Python extensions are enabled
	Enabled bool `json:"enabled" yaml:"enabled"`

	// PythonPath is the path to the Python executable
	PythonPath string `json:"python_path" yaml:"python_path"`

	// RequiredPackages is a list of required Python packages
	RequiredPackages []string `json:"required_packages" yaml:"required_packages"`

	// Features is a map of feature names to their enabled status
	Features map[string]bool `json:"features" yaml:"features"`
}

// NewDefaultPythonExtensionsConfig creates a new PythonExtensionsConfig with default values
func NewDefaultPythonExtensionsConfig() *PythonExtensionsConfig {
	return &PythonExtensionsConfig{
		Enabled:    false,
		PythonPath: "python", // Default to system Python
		RequiredPackages: []string{
			"pandas",
			"numpy",
			"scikit-learn",
			"matplotlib",
			"pyarrow",
			"deltalake",
		},
		Features: map[string]bool{
			string(FeatureFlagMLAnomalyDetection):    false,
			string(FeatureFlagAdvancedDeltaLake):     false,
			string(FeatureFlagAdvancedVisualization): false,
		},
	}
}

// IsPythonExtensionEnabled checks if a specific Python extension feature is enabled
func (c *PythonExtensionsConfig) IsPythonExtensionEnabled(feature string) bool {
	if !c.Enabled {
		return false
	}

	enabled, exists := c.Features[feature]
	return exists && enabled
}

// EnablePythonExtension enables a specific Python extension feature
func (c *PythonExtensionsConfig) EnablePythonExtension(feature string) {
	if c.Features == nil {
		c.Features = make(map[string]bool)
	}
	c.Features[feature] = true
}

// DisablePythonExtension disables a specific Python extension feature
func (c *PythonExtensionsConfig) DisablePythonExtension(feature string) {
	if c.Features == nil {
		c.Features = make(map[string]bool)
	}
	c.Features[feature] = false
}

// EnableAllPythonExtensions enables all Python extension features
func (c *PythonExtensionsConfig) EnableAllPythonExtensions() {
	c.Enabled = true
	if c.Features == nil {
		c.Features = make(map[string]bool)
	}
	c.Features[string(FeatureFlagMLAnomalyDetection)] = true
	c.Features[string(FeatureFlagAdvancedDeltaLake)] = true
	c.Features[string(FeatureFlagAdvancedVisualization)] = true
}

// DisableAllPythonExtensions disables all Python extension features
func (c *PythonExtensionsConfig) DisableAllPythonExtensions() {
	c.Enabled = false
}
