package alerts

import (
	"fmt"
	"math"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
)

// IntelligentAlertingConfig represents the configuration for intelligent alerting
type IntelligentAlertingConfig struct {
	// MinimumDataPoints is the minimum number of data points required for analysis
	MinimumDataPoints int `json:"minimum_data_points"`
	
	// AnalysisPeriod is the time period to analyze for patterns
	AnalysisPeriod time.Duration `json:"analysis_period"`
	
	// UpdateFrequency is how often to update the intelligent alert rules
	UpdateFrequency time.Duration `json:"update_frequency"`
	
	// Sensitivity controls how sensitive the anomaly detection is (0.0-1.0)
	// Higher values mean more alerts will be generated
	Sensitivity float64 `json:"sensitivity"`
	
	// EnableOutlierDetection enables outlier-based alert rules
	EnableOutlierDetection bool `json:"enable_outlier_detection"`
	
	// EnableTrendDeviation enables trend deviation alert rules
	EnableTrendDeviation bool `json:"enable_trend_deviation"`
	
	// EnableSeasonalPatterns enables seasonal pattern alert rules
	EnableSeasonalPatterns bool `json:"enable_seasonal_patterns"`
	
	// AutoDisableUnusedRules automatically disables rules that haven't triggered in a while
	AutoDisableUnusedRules bool `json:"auto_disable_unused_rules"`
	
	// DisableThreshold is the number of days after which an unused rule is disabled
	DisableThreshold int `json:"disable_threshold"`
}

// DefaultIntelligentAlertingConfig returns the default configuration for intelligent alerting
func DefaultIntelligentAlertingConfig() *IntelligentAlertingConfig {
	return &IntelligentAlertingConfig{
		MinimumDataPoints:     100,
		AnalysisPeriod:        7 * 24 * time.Hour, // 1 week
		UpdateFrequency:       24 * time.Hour,     // 1 day
		Sensitivity:           0.7,
		EnableOutlierDetection: true,
		EnableTrendDeviation:   true,
		EnableSeasonalPatterns: true,
		AutoDisableUnusedRules: true,
		DisableThreshold:      30, // 30 days
	}
}

// IntelligentAlertManager manages intelligent alerting
type IntelligentAlertManager struct {
	config      *IntelligentAlertingConfig
	alertManager *AlertManager
	lastUpdate  time.Time
	metricStore MetricStore
}

// MetricStore is an interface for retrieving historical metrics
type MetricStore interface {
	// GetMetricValues retrieves historical values for a metric
	GetMetricValues(metricName string, start, end time.Time, labels map[string]string) ([]MetricDataPoint, error)
	
	// GetMetricNames returns all available metric names
	GetMetricNames() ([]string, error)
}

// MetricDataPoint represents a single data point for a metric
type MetricDataPoint struct {
	Timestamp time.Time
	Value     float64
	Labels    map[string]string
}

// GetConfig returns the current configuration
func (iam *IntelligentAlertManager) GetConfig() *IntelligentAlertingConfig {
	return iam.config
}

// GetAlertManager returns the alert manager
func (iam *IntelligentAlertManager) GetAlertManager() *AlertManager {
	return iam.alertManager
}

// GetMetricStore returns the metric store
func (iam *IntelligentAlertManager) GetMetricStore() MetricStore {
	return iam.metricStore
}

// NewIntelligentAlertManager creates a new intelligent alert manager
func NewIntelligentAlertManager(alertManager *AlertManager, metricStore MetricStore, config *IntelligentAlertingConfig) *IntelligentAlertManager {
	if config == nil {
		config = DefaultIntelligentAlertingConfig()
	}
	
	return &IntelligentAlertManager{
		config:       config,
		alertManager: alertManager,
		metricStore:  metricStore,
		lastUpdate:   time.Now().Add(-config.UpdateFrequency), // Force update on first run
	}
}

// Start starts the intelligent alert manager
func (iam *IntelligentAlertManager) Start() {
	go iam.run()
}

// run is the main loop for the intelligent alert manager
func (iam *IntelligentAlertManager) run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if time.Since(iam.lastUpdate) >= iam.config.UpdateFrequency {
				if err := iam.updateIntelligentRules(); err != nil {
					logging.Error("Failed to update intelligent alert rules", err)
				} else {
					iam.lastUpdate = time.Now()
				}
			}
			
			if iam.config.AutoDisableUnusedRules {
				if err := iam.cleanupUnusedRules(); err != nil {
					logging.Error("Failed to clean up unused alert rules", err)
				}
			}
		}
	}
}

// updateIntelligentRules updates the intelligent alert rules based on historical data
func (iam *IntelligentAlertManager) updateIntelligentRules() error {
	// Get all available metrics
	metricNames, err := iam.metricStore.GetMetricNames()
	if err != nil {
		return fmt.Errorf("failed to get metric names: %w", err)
	}
	
	// Process each metric
	for _, metricName := range metricNames {
		// Skip metrics that don't have enough data
		dataPoints, err := iam.metricStore.GetMetricValues(
			metricName,
			time.Now().Add(-iam.config.AnalysisPeriod),
			time.Now(),
			nil,
		)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to get metric values for %s: %v", metricName, err))
			continue
		}
		
		if len(dataPoints) < iam.config.MinimumDataPoints {
			logging.Debug(fmt.Sprintf("Not enough data points for metric %s: %d/%d", 
				metricName, len(dataPoints), iam.config.MinimumDataPoints))
			continue
		}
		
		// Process the metric with different strategies
		if iam.config.EnableOutlierDetection {
			if err := iam.createOutlierRules(metricName, dataPoints); err != nil {
				logging.Warn(fmt.Sprintf("Failed to create outlier rules for %s: %v", metricName, err))
			}
		}
		
		if iam.config.EnableTrendDeviation {
			if err := iam.createTrendDeviationRules(metricName, dataPoints); err != nil {
				logging.Warn(fmt.Sprintf("Failed to create trend deviation rules for %s: %v", metricName, err))
			}
		}
		
		if iam.config.EnableSeasonalPatterns {
			if err := iam.createSeasonalPatternRules(metricName, dataPoints); err != nil {
				logging.Warn(fmt.Sprintf("Failed to create seasonal pattern rules for %s: %v", metricName, err))
			}
		}
	}
	
	return nil
}

// createOutlierRules creates alert rules based on outlier detection
func (iam *IntelligentAlertManager) createOutlierRules(metricName string, dataPoints []MetricDataPoint) error {
	// Extract values
	values := make([]float64, len(dataPoints))
	for i, dp := range dataPoints {
		values[i] = dp.Value
	}
	
	// Calculate statistics
	mean, stdDev := calculateMeanAndStdDev(values)
	
	// Adjust sensitivity
	threshold := 3.0 * (1.0 - iam.config.Sensitivity) + 1.0 * iam.config.Sensitivity
	
	// Create upper bound rule
	upperRule := &AlertRule{
		ID:                 fmt.Sprintf("intelligent-outlier-upper-%s", metricName),
		Name:               fmt.Sprintf("Intelligent: %s Upper Bound", metricName),
		Description:        fmt.Sprintf("Alert when %s exceeds the upper bound based on historical data", metricName),
		Metric:             metricName,
		Threshold:          mean + threshold*stdDev,
		ComparisonOperator: ">",
		Severity:           SeverityWarning,
		Source:             "intelligent",
		Labels: map[string]string{
			"type":       "intelligent",
			"strategy":   "outlier",
			"bound":      "upper",
			"mean":       fmt.Sprintf("%.2f", mean),
			"std_dev":    fmt.Sprintf("%.2f", stdDev),
			"threshold":  fmt.Sprintf("%.2f", threshold),
			"generated":  "true",
		},
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Create lower bound rule
	createLower := mean > 1.0 // Only create lower bound if mean is significantly above zero
	lowerRule := &AlertRule{
		ID:                 fmt.Sprintf("intelligent-outlier-lower-%s", metricName),
		Name:               fmt.Sprintf("Intelligent: %s Lower Bound", metricName),
		Description:        fmt.Sprintf("Alert when %s falls below the lower bound based on historical data", metricName),
		Metric:             metricName,
		Threshold:          math.Max(0, mean-threshold*stdDev), // Ensure threshold is not negative
		ComparisonOperator: "<",
		Severity:           SeverityWarning,
		Source:             "intelligent",
		Labels: map[string]string{
			"type":       "intelligent",
			"strategy":   "outlier",
			"bound":      "lower",
			"mean":       fmt.Sprintf("%.2f", mean),
			"std_dev":    fmt.Sprintf("%.2f", stdDev),
			"threshold":  fmt.Sprintf("%.2f", threshold),
			"generated":  "true",
		},
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Check if rules already exist
	existingRules, err := iam.alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "outlier",
		"metric":   metricName,
	})
	if err != nil {
		return fmt.Errorf("failed to get existing rules: %w", err)
	}
	
	// Update or create rules
	upperExists := false
	lowerExists := false
	
	for _, rule := range existingRules {
		if rule.Labels["bound"] == "upper" {
			upperExists = true
			rule.Threshold = upperRule.Threshold
			rule.UpdatedAt = time.Now()
			if err := iam.updateRule(rule); err != nil {
				return fmt.Errorf("failed to update upper bound rule: %w", err)
			}
		} else if rule.Labels["bound"] == "lower" {
			lowerExists = true
			rule.Threshold = lowerRule.Threshold
			rule.UpdatedAt = time.Now()
			if err := iam.updateRule(rule); err != nil {
				return fmt.Errorf("failed to update lower bound rule: %w", err)
			}
		}
	}
	
	if !upperExists {
		if err := iam.createRule(upperRule); err != nil {
			return fmt.Errorf("failed to create upper bound rule: %w", err)
		}
	}
	
	if !lowerExists && createLower {
		if err := iam.createRule(lowerRule); err != nil {
			return fmt.Errorf("failed to create lower bound rule: %w", err)
		}
	}
	
	return nil
}

// createTrendDeviationRules creates alert rules based on trend deviation
func (iam *IntelligentAlertManager) createTrendDeviationRules(metricName string, dataPoints []MetricDataPoint) error {
	// Need at least 2 days of data for trend analysis
	if len(dataPoints) < 48 {
		return fmt.Errorf("not enough data points for trend analysis")
	}
	
	// Calculate linear regression
	xValues := make([]float64, len(dataPoints))
	yValues := make([]float64, len(dataPoints))
	
	for i, dp := range dataPoints {
		xValues[i] = float64(i)
		yValues[i] = dp.Value
	}
	
	slope, intercept := calculateLinearRegression(xValues, yValues)
	
	// Calculate residuals
	residuals := make([]float64, len(dataPoints))
	for i := range dataPoints {
		predicted := slope*xValues[i] + intercept
		residuals[i] = math.Abs(yValues[i] - predicted)
	}
	
	// Calculate mean and standard deviation of residuals
	resMean, resStdDev := calculateMeanAndStdDev(residuals)
	
	// Adjust sensitivity
	threshold := 3.0 * (1.0 - iam.config.Sensitivity) + 1.0 * iam.config.Sensitivity
	
	// Create trend deviation rule
	deviationRule := &AlertRule{
		ID:                 fmt.Sprintf("intelligent-trend-%s", metricName),
		Name:               fmt.Sprintf("Intelligent: %s Trend Deviation", metricName),
		Description:        fmt.Sprintf("Alert when %s deviates from the expected trend", metricName),
		Metric:             metricName,
		Threshold:          threshold * resStdDev,
		ComparisonOperator: ">",
		Severity:           SeverityWarning,
		Source:             "intelligent",
		Labels: map[string]string{
			"type":       "intelligent",
			"strategy":   "trend_deviation",
			"slope":      fmt.Sprintf("%.6f", slope),
			"intercept":  fmt.Sprintf("%.2f", intercept),
			"res_mean":   fmt.Sprintf("%.2f", resMean),
			"res_stddev": fmt.Sprintf("%.2f", resStdDev),
			"threshold":  fmt.Sprintf("%.2f", threshold),
			"generated":  "true",
		},
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Check if rule already exists
	existingRules, err := iam.alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "trend_deviation",
		"metric":   metricName,
	})
	if err != nil {
		return fmt.Errorf("failed to get existing rules: %w", err)
	}
	
	// Update or create rule
	if len(existingRules) > 0 {
		rule := existingRules[0]
		rule.Threshold = deviationRule.Threshold
		rule.Labels["slope"] = deviationRule.Labels["slope"]
		rule.Labels["intercept"] = deviationRule.Labels["intercept"]
		rule.Labels["res_mean"] = deviationRule.Labels["res_mean"]
		rule.Labels["res_stddev"] = deviationRule.Labels["res_stddev"]
		rule.UpdatedAt = time.Now()
		if err := iam.updateRule(rule); err != nil {
			return fmt.Errorf("failed to update trend deviation rule: %w", err)
		}
	} else {
		if err := iam.createRule(deviationRule); err != nil {
			return fmt.Errorf("failed to create trend deviation rule: %w", err)
		}
	}
	
	return nil
}

// createSeasonalPatternRules creates alert rules based on seasonal patterns
func (iam *IntelligentAlertManager) createSeasonalPatternRules(metricName string, dataPoints []MetricDataPoint) error {
	// Need at least 7 days of data for weekly seasonality
	if len(dataPoints) < 7*24 {
		return fmt.Errorf("not enough data points for seasonal analysis")
	}
	
	// Group data by hour of day
	hourlyData := make(map[int][]float64)
	for _, dp := range dataPoints {
		hour := dp.Timestamp.Hour()
		hourlyData[hour] = append(hourlyData[hour], dp.Value)
	}
	
	// Calculate statistics for each hour
	hourlyStats := make(map[int]struct {
		Mean   float64
		StdDev float64
	})
	
	for hour, values := range hourlyData {
		mean, stdDev := calculateMeanAndStdDev(values)
		hourlyStats[hour] = struct {
			Mean   float64
			StdDev float64
		}{
			Mean:   mean,
			StdDev: stdDev,
		}
	}
	
	// Create rules for each hour with significant patterns
	for hour, stats := range hourlyStats {
		// Skip if standard deviation is too small
		if stats.StdDev < 0.01*stats.Mean {
			continue
		}
		
		// Adjust sensitivity
		threshold := 3.0 * (1.0 - iam.config.Sensitivity) + 1.0 * iam.config.Sensitivity
		
		// Create upper bound rule
		upperRule := &AlertRule{
			ID:                 fmt.Sprintf("intelligent-seasonal-upper-%s-%d", metricName, hour),
			Name:               fmt.Sprintf("Intelligent: %s Hourly Pattern (Hour %d) Upper", metricName, hour),
			Description:        fmt.Sprintf("Alert when %s exceeds the upper bound for hour %d", metricName, hour),
			Metric:             metricName,
			Threshold:          stats.Mean + threshold*stats.StdDev,
			ComparisonOperator: ">",
			Severity:           SeverityWarning,
			Source:             "intelligent",
			Labels: map[string]string{
				"type":       "intelligent",
				"strategy":   "seasonal",
				"pattern":    "hourly",
				"hour":       fmt.Sprintf("%d", hour),
				"bound":      "upper",
				"mean":       fmt.Sprintf("%.2f", stats.Mean),
				"std_dev":    fmt.Sprintf("%.2f", stats.StdDev),
				"threshold":  fmt.Sprintf("%.2f", threshold),
				"generated":  "true",
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		// Create lower bound rule
		createLower := stats.Mean > 1.0 // Only create lower bound if mean is significantly above zero
		lowerRule := &AlertRule{
			ID:                 fmt.Sprintf("intelligent-seasonal-lower-%s-%d", metricName, hour),
			Name:               fmt.Sprintf("Intelligent: %s Hourly Pattern (Hour %d) Lower", metricName, hour),
			Description:        fmt.Sprintf("Alert when %s falls below the lower bound for hour %d", metricName, hour),
			Metric:             metricName,
			Threshold:          math.Max(0, stats.Mean-threshold*stats.StdDev), // Ensure threshold is not negative
			ComparisonOperator: "<",
			Severity:           SeverityWarning,
			Source:             "intelligent",
			Labels: map[string]string{
				"type":       "intelligent",
				"strategy":   "seasonal",
				"pattern":    "hourly",
				"hour":       fmt.Sprintf("%d", hour),
				"bound":      "lower",
				"mean":       fmt.Sprintf("%.2f", stats.Mean),
				"std_dev":    fmt.Sprintf("%.2f", stats.StdDev),
				"threshold":  fmt.Sprintf("%.2f", threshold),
				"generated":  "true",
			},
			Enabled:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		// Update createLower based on stats
		if stats.Mean <= 0 || stats.StdDev <= 0 {
			createLower = false
		}
		
		// Check if rules already exist
		existingRules, err := iam.alertManager.GetRulesByLabels(map[string]string{
			"type":     "intelligent",
			"strategy": "seasonal",
			"pattern":  "hourly",
			"hour":     fmt.Sprintf("%d", hour),
			"metric":   metricName,
		})
		if err != nil {
			return fmt.Errorf("failed to get existing rules: %w", err)
		}
		
		// Update or create rules
		upperExists := false
		lowerExists := false
		
		for _, rule := range existingRules {
			if rule.Labels["bound"] == "upper" {
				upperExists = true
				rule.Threshold = upperRule.Threshold
				rule.UpdatedAt = time.Now()
				if err := iam.updateRule(rule); err != nil {
					return fmt.Errorf("failed to update upper bound rule: %w", err)
				}
			} else if rule.Labels["bound"] == "lower" {
				lowerExists = true
				rule.Threshold = lowerRule.Threshold
				rule.UpdatedAt = time.Now()
				if err := iam.updateRule(rule); err != nil {
					return fmt.Errorf("failed to update lower bound rule: %w", err)
				}
			}
		}
		
		if !upperExists {
			if err := iam.createRule(upperRule); err != nil {
				return fmt.Errorf("failed to create upper bound rule: %w", err)
			}
		}
		
		if !lowerExists && createLower {
			if err := iam.createRule(lowerRule); err != nil {
				return fmt.Errorf("failed to create lower bound rule: %w", err)
			}
		}
	}
	
	return nil
}

// AnalyzeMetricData analyzes a specific metric and creates intelligent alert rules
// This method is primarily used for testing purposes
func (iam *IntelligentAlertManager) AnalyzeMetricData(metricName string) error {
	// Get historical data for the metric
	end := time.Now()
	start := end.Add(-iam.config.AnalysisPeriod)
	
	dataPoints, err := iam.metricStore.GetMetricValues(metricName, start, end, nil)
	if err != nil {
		return fmt.Errorf("failed to get metric values: %w", err)
	}
	
	// Check if we have enough data points
	if len(dataPoints) < iam.config.MinimumDataPoints {
		return fmt.Errorf("insufficient data points for metric %s: got %d, need %d", 
			metricName, len(dataPoints), iam.config.MinimumDataPoints)
	}
	
	// Create rules based on the configuration
	if iam.config.EnableOutlierDetection {
		if err := iam.createOutlierRules(metricName, dataPoints); err != nil {
			return fmt.Errorf("failed to create outlier rules: %w", err)
		}
	}
	
	if iam.config.EnableTrendDeviation {
		if err := iam.createTrendDeviationRules(metricName, dataPoints); err != nil {
			return fmt.Errorf("failed to create trend deviation rules: %w", err)
		}
	}
	
	if iam.config.EnableSeasonalPatterns {
		if err := iam.createSeasonalPatternRules(metricName, dataPoints); err != nil {
			return fmt.Errorf("failed to create seasonal pattern rules: %w", err)
		}
	}
	
	return nil
}

// cleanupUnusedRules disables rules that haven't triggered in a while
func (iam *IntelligentAlertManager) cleanupUnusedRules() error {
	// Get all intelligent rules
	rules, err := iam.alertManager.GetRulesByLabels(map[string]string{
		"type":      "intelligent",
		"generated": "true",
	})
	if err != nil {
		return fmt.Errorf("failed to get intelligent rules: %w", err)
	}
	
	// Check each rule's last trigger time
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		
		// Skip rules that don't have a last triggered time (never triggered)
		if rule.LastTriggered.IsZero() {
			// If the rule is older than the threshold, disable it
			if time.Since(rule.CreatedAt) > time.Duration(iam.config.DisableThreshold)*24*time.Hour {
				rule.Enabled = false
				rule.UpdatedAt = time.Now()
				if err := iam.updateRule(rule); err != nil {
					logging.Warn(fmt.Sprintf("Failed to disable unused rule %s: %v", rule.ID, err))
				} else {
					logging.Info(fmt.Sprintf("Disabled unused rule %s (never triggered)", rule.ID))
				}
			}
			continue
		}
		
		// Check if the rule hasn't triggered in a while
		if time.Since(rule.LastTriggered) > time.Duration(iam.config.DisableThreshold)*24*time.Hour {
			rule.Enabled = false
			rule.UpdatedAt = time.Now()
			if err := iam.updateRule(rule); err != nil {
				logging.Warn(fmt.Sprintf("Failed to disable unused rule %s: %v", rule.ID, err))
			} else {
				logging.Info(fmt.Sprintf("Disabled unused rule %s (last triggered: %s)", 
					rule.ID, rule.LastTriggered.Format(time.RFC3339)))
			}
		}
	}
	
	return nil
}

// Helper functions for statistical calculations

// calculateMeanAndStdDev calculates the mean and standard deviation of a slice of values
func calculateMeanAndStdDev(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	
	// Calculate mean
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))
	
	// Calculate standard deviation
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	stdDev := math.Sqrt(variance)
	
	return mean, stdDev
}

// calculateLinearRegression calculates the slope and intercept of a linear regression
func calculateLinearRegression(x, y []float64) (float64, float64) {
	if len(x) != len(y) || len(x) == 0 {
		return 0, 0
	}
	
	n := float64(len(x))
	
	// Calculate means
	sumX := 0.0
	sumY := 0.0
	for i := range x {
		sumX += x[i]
		sumY += y[i]
	}
	meanX := sumX / n
	meanY := sumY / n
	
	// Calculate slope
	numerator := 0.0
	denominator := 0.0
	for i := range x {
		xDiff := x[i] - meanX
		yDiff := y[i] - meanY
		numerator += xDiff * yDiff
		denominator += xDiff * xDiff
	}
	
	// Avoid division by zero
	if denominator == 0 {
		return 0, meanY
	}
	
	slope := numerator / denominator
	intercept := meanY - slope*meanX
	
	return slope, intercept
}
