package datalake

// DetectSuddenChanges implements the sudden change detection for the MockMetadataManager
func (mm *MockMetadataManager) DetectSuddenChanges(options ChangeDetectionOptions) (*ChangeDetectionResult, error) {
	// Get current metrics using the mock function
	var currentMetrics map[string]float64
	var err error

	if mm.MockCalculateFieldMetrics != nil {
		currentMetrics, err = mm.MockCalculateFieldMetrics(options.Field, options.MetricTypes)
		if err != nil {
			return nil, err
		}
	} else if mm.MetadataManager != nil {
		// Fall back to the real implementation if available
		currentMetrics, err = mm.MetadataManager.calculateFieldMetrics(options.Field, options.MetricTypes)
		if err != nil {
			return nil, err
		}
	} else {
		// Default empty metrics if neither mock nor real implementation is available
		currentMetrics = make(map[string]float64)
	}

	// Create a simple detector with the current metrics
	detector := &SimpleSuddenChangeDetector{
		MetricsDir: options.MetricsDir,
		Metrics: map[string]map[string]float64{
			options.Field: currentMetrics,
		},
	}

	// Use a run number that ensures we have enough previous runs
	runNumber := options.MinRuns + 2

	return detector.DetectSuddenChanges(options.Field, runNumber)
}
