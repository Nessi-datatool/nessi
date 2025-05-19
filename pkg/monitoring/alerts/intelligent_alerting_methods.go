// Implementation of monitoring and alerting features.

package alerts

import (
	"fmt"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
)

// StartAnalysis starts the background analysis of metrics
func (iam *IntelligentAlertManager) StartAnalysis() {
	// Start a ticker for periodic analysis
	ticker := time.NewTicker(iam.config.UpdateFrequency)
	defer ticker.Stop()

	// Run initial analysis
	iam.analyzeMetrics()

	// Run analysis on ticker
	go func() {
		for range ticker.C {
			iam.analyzeMetrics()
		}
	}()
}

// analyzeMetrics analyzes all metrics and creates intelligent alert rules
func (iam *IntelligentAlertManager) analyzeMetrics() {
	logging.Info("Starting intelligent metrics analysis")
	startTime := time.Now()
	
	// Get all metric names
	metricNames, err := iam.metricStore.GetMetricNames()
	if err != nil {
		logging.Error("Failed to get metric names", err)
		return
	}
	
	// Analyze each metric
	for _, metricName := range metricNames {
		// Get historical data for the metric
		end := time.Now()
		start := end.Add(-iam.config.AnalysisPeriod)
		
		dataPoints, err := iam.metricStore.GetMetricValues(metricName, start, end, nil)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to get data for metric %s: %v", metricName, err))
			continue
		}
		
		// Skip metrics with insufficient data
		if len(dataPoints) < iam.config.MinimumDataPoints {
			logging.Debug(fmt.Sprintf("Skipping metric %s: insufficient data points (%d < %d)", 
				metricName, len(dataPoints), iam.config.MinimumDataPoints))
			continue
		}
		
		// Create rules based on enabled strategies
		if iam.config.EnableOutlierDetection {
			if err := iam.createOutlierRules(metricName, dataPoints); err != nil {
				logging.Warn(fmt.Sprintf("Failed to create outlier rules for metric %s: %v", metricName, err))
			}
		}
		
		if iam.config.EnableTrendDeviation {
			if err := iam.createTrendDeviationRules(metricName, dataPoints); err != nil {
				logging.Warn(fmt.Sprintf("Failed to create trend deviation rules for metric %s: %v", metricName, err))
			}
		}
		
		if iam.config.EnableSeasonalPatterns {
			if err := iam.createSeasonalPatternRules(metricName, dataPoints); err != nil {
				logging.Warn(fmt.Sprintf("Failed to create seasonal pattern rules for metric %s: %v", metricName, err))
			}
		}
	}
	
	// Clean up unused rules if enabled
	if iam.config.AutoDisableUnusedRules {
		if err := iam.cleanupUnusedRules(); err != nil {
			logging.Warn(fmt.Sprintf("Failed to clean up unused rules: %v", err))
		}
	}
	
	duration := time.Since(startTime)
	logging.Info(fmt.Sprintf("Completed intelligent metrics analysis in %v", duration))
}

// These methods are now defined in intelligent_alerting.go
// GetConfig() *IntelligentAlertingConfig
// GetMetricStore() MetricStore
// GetAlertManager() *AlertManager
