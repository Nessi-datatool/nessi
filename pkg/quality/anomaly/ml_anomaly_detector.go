package anomaly

import (
	"encoding/json"
	"fmt"

	"github.com/nessi-dev/nessi-dev/pkg/config"
	"github.com/nessi-dev/nessi-dev/pkg/python"
)

// MLAnomalyDetector provides machine learning based anomaly detection
type MLAnomalyDetector struct {
	bridge         *python.PythonBridge
	featureFlagMgr *config.FeatureFlagManager
}

// AnomalyDetectionResult represents the result of anomaly detection
type AnomalyDetectionResult struct {
	AnomalyIndices []int     `json:"anomaly_indices"`
	AnomalyValues  []float64 `json:"anomaly_values"`
	Visualization  string    `json:"visualization"` // Base64 encoded PNG
}

// SeasonalDecompositionResult represents the result of seasonal decomposition
type SeasonalDecompositionResult struct {
	Trend         []float64                       `json:"trend"`
	Seasonal      []float64                       `json:"seasonal"`
	Residual      []float64                       `json:"residual"`
	Visualizations map[string]string              `json:"visualizations"` // Base64 encoded PNGs
	Error         string                          `json:"error,omitempty"`
}

// NewMLAnomalyDetector creates a new MLAnomalyDetector
func NewMLAnomalyDetector() *MLAnomalyDetector {
	return &MLAnomalyDetector{
		bridge:         python.GetPythonBridge(),
		featureFlagMgr: config.NewFeatureFlagManager(),
	}
}

// IsAvailable checks if ML-based anomaly detection is available
func (d *MLAnomalyDetector) IsAvailable() bool {
	return d.bridge.IsPythonFeatureEnabled(config.FeatureFlagMLAnomalyDetection) &&
		d.bridge.IsPythonFeatureEnabled(config.FeatureFlagPythonExtensions)
}

// DetectAnomalies detects anomalies in the given data
func (d *MLAnomalyDetector) DetectAnomalies(data []float64, algorithm string, contamination float64) (*AnomalyDetectionResult, error) {
	if !d.IsAvailable() {
		return nil, fmt.Errorf("ML-based anomaly detection is not available. Enable Python extensions to use this feature.")
	}

	// Prepare arguments for Python script
	args := map[string]interface{}{
		"action":        "detect_anomalies",
		"data":          data,
		"algorithm":     algorithm,
		"contamination": contamination,
	}

	// Execute Python script
	output, err := d.bridge.ExecuteScript("anomaly_detection.py", args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute anomaly detection script: %v", err)
	}

	// Parse result
	var result AnomalyDetectionResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse anomaly detection result: %v", err)
	}

	return &result, nil
}

// DecomposeTimeSeries decomposes a time series into trend, seasonal, and residual components
func (d *MLAnomalyDetector) DecomposeTimeSeries(data []float64, period int) (*SeasonalDecompositionResult, error) {
	if !d.IsAvailable() {
		return nil, fmt.Errorf("ML-based time series decomposition is not available. Enable Python extensions to use this feature.")
	}

	// Prepare arguments for Python script
	args := map[string]interface{}{
		"action": "seasonal_decomposition",
		"data":   data,
		"period": period,
	}

	// Execute Python script
	output, err := d.bridge.ExecuteScript("anomaly_detection.py", args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute seasonal decomposition script: %v", err)
	}

	// Parse result
	var result SeasonalDecompositionResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse seasonal decomposition result: %v", err)
	}

	return &result, nil
}

// GetSupportedAlgorithms returns the list of supported anomaly detection algorithms
func (d *MLAnomalyDetector) GetSupportedAlgorithms() []string {
	return []string{
		"isolation_forest",
		"one_class_svm",
		"local_outlier_factor",
	}
}
