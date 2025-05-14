package profile

import (
	"fmt"
	"math"
	"sort"
)

// GenerateDistributionAnalysis creates a detailed distribution analysis for a profile
func GenerateDistributionAnalysis(profile *Profile) *DistributionAnalysis {
	analysis := &DistributionAnalysis{
		Histogram:   make(map[string]int),
		Frequencies: make(map[string]float64),
		TopValues:   []ValueCount{},
		BottomValues: []ValueCount{},
	}

	// Skip if no value counts
	if profile.ValueCounts == nil || len(profile.ValueCounts) == 0 {
		return analysis
	}

	// Convert to a slice for sorting
	var valueCounts []ValueCount
	totalCount := 0
	for value, count := range profile.ValueCounts {
		valueCounts = append(valueCounts, ValueCount{
			Value: value,
			Count: count,
		})
		totalCount += count
	}

	// Sort by count (descending)
	sort.Slice(valueCounts, func(i, j int) bool {
		return valueCounts[i].Count > valueCounts[j].Count
	})

	// Get top values (up to 10)
	topCount := int(math.Min(10, float64(len(valueCounts))))
	analysis.TopValues = valueCounts[:topCount]

	// Get bottom values (up to 10)
	if len(valueCounts) > 10 {
		analysis.BottomValues = valueCounts[len(valueCounts)-10:]
	} else {
		analysis.BottomValues = valueCounts
	}

	// Create histogram and frequencies
	for _, vc := range valueCounts {
		valueStr := formatValue(vc.Value)
		analysis.Histogram[valueStr] = vc.Count
		analysis.Frequencies[valueStr] = float64(vc.Count) / float64(totalCount)
	}

	// Calculate statistical measures for numeric data
	if profile.Type == "integer" || profile.Type == "float" {
		// Extract numeric values
		var numericValues []float64
		for value, count := range profile.ValueCounts {
			// Convert to float64
			var floatVal float64
			switch v := value.(type) {
			case int:
				floatVal = float64(v)
			case int32:
				floatVal = float64(v)
			case int64:
				floatVal = float64(v)
			case float32:
				floatVal = float64(v)
			case float64:
				floatVal = v
			default:
				continue
			}

			// Add value multiple times based on count
			for i := 0; i < count; i++ {
				numericValues = append(numericValues, floatVal)
			}
		}

		if len(numericValues) > 0 {
			// Calculate mean
			mean := calculateMean(numericValues)

			// Calculate standard deviation
			stdDev := calculateStdDev(numericValues, mean)

			// Calculate skewness and kurtosis
			analysis.Skewness = calculateSkewness(numericValues, mean, stdDev)
			analysis.Kurtosis = calculateKurtosis(numericValues, mean, stdDev)

			// Check if distribution is approximately normal
			analysis.IsNormal = isNormalDistribution(analysis.Skewness, analysis.Kurtosis)

			// Calculate entropy
			analysis.Entropy = calculateEntropy(analysis.Frequencies)
		}
	}

	return analysis
}

// formatValue converts a value to a string representation
func formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", value)
}

// calculateMean calculates the mean of a slice of values
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateStdDev calculates the standard deviation
func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	return math.Sqrt(sumSquaredDiff / float64(len(values)))
}

// isNormalDistribution checks if a distribution is approximately normal
func isNormalDistribution(skewness, kurtosis float64) bool {
	// A normal distribution has skewness close to 0 and kurtosis close to 3
	return math.Abs(skewness) < 0.5 && math.Abs(kurtosis-3) < 1
}
