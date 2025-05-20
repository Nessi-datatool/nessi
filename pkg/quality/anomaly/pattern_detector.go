package anomaly

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// PatternType represents the type of pattern detected
type PatternType string

const (
	// ConstantChangePattern represents a consistent increase or decrease over time
	ConstantChangePattern PatternType = "constant_change"

	// SpikePattern represents a sudden increase followed by a return to normal
	SpikePattern PatternType = "spike"

	// DipPattern represents a sudden decrease followed by a return to normal
	DipPattern PatternType = "dip"

	// OscillationPattern represents alternating values
	OscillationPattern PatternType = "oscillation"

	// TrendDeviationPattern represents a deviation from a historical trend
	TrendDeviationPattern PatternType = "trend_deviation"
)

// PatternDetectionConfig contains configuration for pattern detection
type PatternDetectionConfig struct {
	// MinDataPoints is the minimum number of data points required for detection
	MinDataPoints int

	// ConstantChangeThreshold is the minimum percentage change for constant change detection
	ConstantChangeThreshold float64

	// SpikeThreshold is the minimum percentage change for spike detection
	SpikeThreshold float64

	// DipThreshold is the minimum percentage change for dip detection
	DipThreshold float64

	// OscillationThreshold is the minimum percentage change for oscillation detection
	OscillationThreshold float64

	// TrendDeviationThreshold is the minimum percentage deviation from trend
	TrendDeviationThreshold float64

	// HistoricalWindowSize is the number of data points to use for trend analysis
	HistoricalWindowSize int
}

// DefaultPatternDetectionConfig returns the default configuration
func DefaultPatternDetectionConfig() *PatternDetectionConfig {
	return &PatternDetectionConfig{
		MinDataPoints:           5,
		ConstantChangeThreshold: 0.05, // 5%
		SpikeThreshold:          0.20, // 20%
		DipThreshold:            0.20, // 20%
		OscillationThreshold:    0.10, // 10%
		TrendDeviationThreshold: 0.15, // 15%
		HistoricalWindowSize:    30,
	}
}

// TimeSeriesDataPoint represents a data point in a time series
type TimeSeriesDataPoint struct {
	Timestamp time.Time
	Value     float64
}

// PatternDetectionResult represents the result of pattern detection
type PatternDetectionResult struct {
	// Pattern is the detected pattern type
	Pattern PatternType

	// Confidence is the confidence level of the detection (0-1)
	Confidence float64

	// Description is a human-readable description of the pattern
	Description string

	// StartIndex is the index where the pattern starts
	StartIndex int

	// EndIndex is the index where the pattern ends
	EndIndex int

	// Magnitude is the magnitude of the pattern (e.g., percentage change)
	Magnitude float64
}

// PatternDetector detects patterns in time series data
type PatternDetector struct {
	config *PatternDetectionConfig
}

// NewPatternDetector creates a new pattern detector with the given configuration
func NewPatternDetector(config *PatternDetectionConfig) *PatternDetector {
	if config == nil {
		config = DefaultPatternDetectionConfig()
	}
	return &PatternDetector{
		config: config,
	}
}

// DetectPatterns detects all patterns in the given time series data
func (pd *PatternDetector) DetectPatterns(data []TimeSeriesDataPoint) []*PatternDetectionResult {
	if len(data) < pd.config.MinDataPoints {
		return nil
	}

	// Sort data by timestamp
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp.Before(data[j].Timestamp)
	})

	// Extract values
	values := make([]float64, len(data))
	for i, point := range data {
		values[i] = point.Value
	}

	var results []*PatternDetectionResult

	// Detect constant changes
	if constantChange := pd.detectConstantChange(values); constantChange != nil {
		results = append(results, constantChange)
	}

	// Detect spikes
	if spikes := pd.detectSpikes(values); len(spikes) > 0 {
		results = append(results, spikes...)
	}

	// Detect dips
	if dips := pd.detectDips(values); len(dips) > 0 {
		results = append(results, dips...)
	}

	// Detect oscillations
	if oscillation := pd.detectOscillation(values); oscillation != nil {
		results = append(results, oscillation)
	}

	// Detect trend deviations
	if trendDeviations := pd.detectTrendDeviations(values); len(trendDeviations) > 0 {
		results = append(results, trendDeviations...)
	}

	return results
}

// detectConstantChange detects a consistent increase or decrease over time
func (pd *PatternDetector) detectConstantChange(values []float64) *PatternDetectionResult {
	if len(values) < pd.config.MinDataPoints {
		return nil
	}

	// Calculate the slope using linear regression
	n := float64(len(values))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumXX := 0.0

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	// Calculate slope
	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	// Calculate R-squared (coefficient of determination)
	meanY := sumY / n
	totalSS := 0.0
	residualSS := 0.0

	for i, y := range values {
		x := float64(i)
		predicted := slope*x + (sumY-slope*sumX)/n
		totalSS += math.Pow(y-meanY, 2)
		residualSS += math.Pow(y-predicted, 2)
	}

	rSquared := 1 - residualSS/totalSS

	// Calculate total percentage change
	totalChange := (values[len(values)-1] - values[0]) / values[0]

	// Check if the change is significant and consistent
	if math.Abs(totalChange) >= pd.config.ConstantChangeThreshold && rSquared >= 0.7 {
		direction := "increase"
		if slope < 0 {
			direction = "decrease"
		}

		return &PatternDetectionResult{
			Pattern:     ConstantChangePattern,
			Confidence:  rSquared,
			Description: fmt.Sprintf("Constant %s of %.2f%% over %d data points", direction, totalChange*100, len(values)),
			StartIndex:  0,
			EndIndex:    len(values) - 1,
			Magnitude:   math.Abs(totalChange),
		}
	}

	return nil
}

// detectSpikes detects sudden increases followed by returns to normal
func (pd *PatternDetector) detectSpikes(values []float64) []*PatternDetectionResult {
	if len(values) < pd.config.MinDataPoints {
		return nil
	}

	var results []*PatternDetectionResult

	// Calculate moving average
	windowSize := 3
	if windowSize > len(values) {
		windowSize = len(values)
	}

	movingAvg := make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		sum := 0.0
		count := 0

		for j := maxInt(0, i-windowSize+1); j <= minInt(len(values)-1, i+windowSize-1); j++ {
			if j != i { // Exclude the current point
				sum += values[j]
				count++
			}
		}

		if count > 0 {
			movingAvg[i] = sum / float64(count)
		} else {
			movingAvg[i] = values[i]
		}
	}

	// Detect spikes
	for i := 1; i < len(values)-1; i++ {
		// Calculate percentage change from moving average
		change := (values[i] - movingAvg[i]) / movingAvg[i]

		// Check if it's a spike (value higher than neighbors and above threshold)
		if change >= pd.config.SpikeThreshold &&
			values[i] > values[i-1] &&
			values[i] > values[i+1] {

			// Find the extent of the spike
			start := i
			for start > 0 && values[start-1] < values[start] {
				start--
			}

			end := i
			for end < len(values)-1 && values[end] > values[end+1] {
				end++
			}

			results = append(results, &PatternDetectionResult{
				Pattern:     SpikePattern,
				Confidence:  minFloat64(1.0, change/pd.config.SpikeThreshold),
				Description: fmt.Sprintf("Spike of %.2f%% at index %d", change*100, i),
				StartIndex:  start,
				EndIndex:    end,
				Magnitude:   change,
			})
		}
	}

	return results
}

// detectDips detects sudden decreases followed by returns to normal
func (pd *PatternDetector) detectDips(values []float64) []*PatternDetectionResult {
	if len(values) < pd.config.MinDataPoints {
		return nil
	}

	var results []*PatternDetectionResult

	// Calculate moving average
	windowSize := 3
	if windowSize > len(values) {
		windowSize = len(values)
	}

	movingAvg := make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		sum := 0.0
		count := 0

		for j := maxInt(0, i-windowSize+1); j <= minInt(len(values)-1, i+windowSize-1); j++ {
			if j != i { // Exclude the current point
				sum += values[j]
				count++
			}
		}

		if count > 0 {
			movingAvg[i] = sum / float64(count)
		} else {
			movingAvg[i] = values[i]
		}
	}

	// Detect dips
	for i := 1; i < len(values)-1; i++ {
		// Calculate percentage change from moving average
		change := (movingAvg[i] - values[i]) / movingAvg[i]

		// Check if it's a dip (value lower than neighbors and above threshold)
		if change >= pd.config.DipThreshold &&
			values[i] < values[i-1] &&
			values[i] < values[i+1] {

			// Find the extent of the dip
			start := i
			for start > 0 && values[start-1] > values[start] {
				start--
			}

			end := i
			for end < len(values)-1 && values[end] < values[end+1] {
				end++
			}

			results = append(results, &PatternDetectionResult{
				Pattern:     DipPattern,
				Confidence:  minFloat64(1.0, change/pd.config.DipThreshold),
				Description: fmt.Sprintf("Dip of %.2f%% at index %d", change*100, i),
				StartIndex:  start,
				EndIndex:    end,
				Magnitude:   change,
			})
		}
	}

	return results
}

// detectOscillation detects alternating patterns in the data
func (pd *PatternDetector) detectOscillation(values []float64) *PatternDetectionResult {
	if len(values) < pd.config.MinDataPoints {
		return nil
	}

	// Count direction changes
	directionChanges := 0
	direction := 0 // 0 = undefined, 1 = increasing, -1 = decreasing

	for i := 1; i < len(values); i++ {
		if values[i] > values[i-1] {
			// Increasing
			if direction == -1 {
				directionChanges++
			}
			direction = 1
		} else if values[i] < values[i-1] {
			// Decreasing
			if direction == 1 {
				directionChanges++
			}
			direction = -1
		}
	}

	// Calculate average magnitude of changes
	totalMagnitude := 0.0
	for i := 1; i < len(values); i++ {
		if values[i-1] != 0 {
			totalMagnitude += math.Abs((values[i] - values[i-1]) / values[i-1])
		}
	}
	avgMagnitude := totalMagnitude / float64(len(values)-1)

	// Check if we have enough direction changes and the magnitude is significant
	if directionChanges >= len(values)/3 && avgMagnitude >= pd.config.OscillationThreshold {
		confidence := minFloat64(1.0, float64(directionChanges)/float64(len(values)/2))

		return &PatternDetectionResult{
			Pattern:     OscillationPattern,
			Confidence:  confidence,
			Description: fmt.Sprintf("Oscillation with %d direction changes and average magnitude of %.2f%%", directionChanges, avgMagnitude*100),
			StartIndex:  0,
			EndIndex:    len(values) - 1,
			Magnitude:   avgMagnitude,
		}
	}

	return nil
}

// detectTrendDeviations detects deviations from historical trends
func (pd *PatternDetector) detectTrendDeviations(values []float64) []*PatternDetectionResult {
	if len(values) < pd.config.MinDataPoints {
		return nil
	}

	var results []*PatternDetectionResult

	// We need enough historical data to establish a trend
	if len(values) < pd.config.HistoricalWindowSize {
		return nil
	}

	// Calculate trend using linear regression on historical window
	for i := pd.config.HistoricalWindowSize; i < len(values); i++ {
		// Get historical window
		historicalValues := values[i-pd.config.HistoricalWindowSize : i]

		// Calculate linear regression
		n := float64(len(historicalValues))
		sumX := 0.0
		sumY := 0.0
		sumXY := 0.0
		sumXX := 0.0

		for j, y := range historicalValues {
			x := float64(j)
			sumX += x
			sumY += y
			sumXY += x * y
			sumXX += x * x
		}

		// Calculate slope and intercept
		slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
		intercept := (sumY - slope*sumX) / n

		// Predict next value
		predicted := slope*float64(len(historicalValues)) + intercept

		// Calculate deviation
		actual := values[i]
		deviation := math.Abs((actual - predicted) / predicted)

		// Check if deviation exceeds threshold
		if deviation >= pd.config.TrendDeviationThreshold {
			direction := "above"
			if actual < predicted {
				direction = "below"
			}

			results = append(results, &PatternDetectionResult{
				Pattern:     TrendDeviationPattern,
				Confidence:  minFloat64(1.0, deviation/pd.config.TrendDeviationThreshold),
				Description: fmt.Sprintf("Value at index %d is %.2f%% %s predicted trend", i, deviation*100, direction),
				StartIndex:  i - pd.config.HistoricalWindowSize,
				EndIndex:    i,
				Magnitude:   deviation,
			})
		}
	}

	return results
}

// Helper functions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
