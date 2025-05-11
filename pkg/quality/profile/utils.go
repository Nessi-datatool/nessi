package profile

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// CommonPrefix finds the longest common prefix between two strings
func CommonPrefix(s1, s2 string) string {
	if len(s1) == 0 || len(s2) == 0 {
		return ""
	}
	
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}
	
	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			return s1[:i]
		}
	}

	// If no mismatch found, return the entire prefix
	return s1[:minLen]
}

// CommonSuffix finds the longest common suffix between two strings
func CommonSuffix(s1, s2 string) string {
	if len(s1) == 0 || len(s2) == 0 {
		return ""
	}

	// Find the minimum length suffix to check
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	// Compare characters from the end of both strings
	for i := 0; i < minLen; i++ {
		if s1[len(s1)-1-i] != s2[len(s2)-1-i] {
			// Return the suffix up to the first mismatch
			return s1[len(s1)-i-1:]
		}
	}

	// If no mismatch found, return the entire suffix
	return s1[len(s1)-minLen:]
}

// DetectNumericPatterns finds common patterns in numeric data
func DetectNumericPatterns(values []float64) []string {
	var patterns []string
	if len(values) == 0 {
		return patterns
	}

	// Sort values to detect patterns
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Detect integer vs decimal
	allIntegers := true
	for _, v := range values {
		if math.Floor(v) != v {
			allIntegers = false
			break
		}
	}
	if allIntegers {
		patterns = append(patterns, "Integer values")
	} else {
		patterns = append(patterns, "Decimal values")
	}

	// Detect common ranges
	min := sorted[0]
	max := sorted[len(sorted)-1]
	if max-min > 0 {
		// Check if values are mostly in a specific range
		if min >= 0 && max <= 100 {
			patterns = append(patterns, "Values in range 0-100")
		} else if min >= 0 && max <= 1000 {
			patterns = append(patterns, "Values in range 0-1000")
		} else if min >= 0 && max <= 1000000 {
			patterns = append(patterns, "Values in range 0-1M")
		}
	}

	// Detect if values are mostly positive/negative
	positiveCount := 0
	negativeCount := 0
	for _, v := range values {
		if v > 0 {
			positiveCount++
		} else if v < 0 {
			negativeCount++
		}
	}
	if float64(positiveCount)/float64(len(values)) > 0.9 {
		patterns = append(patterns, "Mostly positive values")
	} else if float64(negativeCount)/float64(len(values)) > 0.9 {
		patterns = append(patterns, "Mostly negative values")
	}

	return patterns
}

// DetectDatePatterns finds common patterns in date data
func DetectDatePatterns(dates []time.Time) []string {
	var patterns []string
	if len(dates) == 0 {
		return patterns
	}

	// Sort dates
	sorted := make([]time.Time, len(dates))
	copy(sorted, dates)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Before(sorted[j])
	})

	// Detect common time ranges
	min := sorted[0]
	max := sorted[len(sorted)-1]
	duration := max.Sub(min)

	// Check if dates are mostly recent
	if duration.Hours() < 24*7 { // Less than a week
		patterns = append(patterns, "Recent dates (within a week)")
	} else if duration.Hours() < 24*30 { // Less than a month
		patterns = append(patterns, "Recent dates (within a month)")
	} else if duration.Hours() < 24*365 { // Less than a year
		patterns = append(patterns, "Recent dates (within a year)")
	}

	// Detect if dates are mostly in a specific year
	yearCounts := make(map[int]int)
	for _, d := range dates {
		yearCounts[d.Year()]++
	}
	mostCommonYear := 0
	maxCount := 0
	for year, count := range yearCounts {
		if count > maxCount {
			maxCount = count
			mostCommonYear = year
		}
	}
	if float64(maxCount)/float64(len(dates)) > 0.9 {
		patterns = append(patterns, fmt.Sprintf("Mostly from year %d", mostCommonYear))
	}

	return patterns
}

// DetectNumericAnomalies finds anomalies in numeric data using statistical methods
func DetectNumericAnomalies(values []float64, threshold float64) []Anomaly {
	var anomalies []Anomaly
	if len(values) == 0 {
		return anomalies
	}

	// Sort values
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate quartiles
	q1 := sorted[len(sorted)*1/4]
	q3 := sorted[len(sorted)*3/4]
	iqr := q3 - q1

	// Define outlier thresholds
	lowerBound := q1 - threshold*iqr
	upperBound := q3 + threshold*iqr

	// Find outliers
	outliers := []float64{}
	for _, v := range values {
		if v < lowerBound || v > upperBound {
			outliers = append(outliers, v)
		}
	}

	if len(outliers) > 0 {
		// Create examples string
		exampleStr := ""
		for i, v := range outliers[:min(len(outliers), 5)] {
			if i > 0 {
				exampleStr += ", "
			}
			exampleStr += fmt.Sprintf("%.2f", v)
		}
		anomalies = append(anomalies, Anomaly{
			Type:        "Outlier",
			Value:       outliers[0],
			Description: fmt.Sprintf("Values outside IQR range (%.2f to %.2f). Examples: %s", lowerBound, upperBound, exampleStr),
		})
	}

	// Check for unusual value distribution
	if len(values) > 1 {
		// Calculate standard deviation
		var sum, sumSquares float64
		for _, v := range values {
			sum += v
			sumSquares += v * v
		}
		mean := sum / float64(len(values))
		stdDev := math.Sqrt(sumSquares/float64(len(values)) - mean*mean)

		// Check for high standard deviation
		if stdDev > 100*mean && mean != 0 {
			anomalies = append(anomalies, Anomaly{
				Type:        "High Variability",
				Value:       stdDev,
				Description: fmt.Sprintf("Standard deviation (%.2f) is unusually high compared to mean (%.2f)", stdDev, mean),
			})
		}

		// Check for unusual value concentration
		valueCounts := make(map[float64]int)
		for _, v := range values {
			valueCounts[v]++
		}

		mostCommonValue := 0.0
		maxCount := 0
		for v, count := range valueCounts {
			if count > maxCount {
				mostCommonValue = v
				maxCount = count
			}
		}

		// If more than 80% of values are the same
		if float64(maxCount)/float64(len(values)) > 0.8 {
			anomalies = append(anomalies, Anomaly{
				Type:        "Value Concentration",
				Value:       mostCommonValue,
				Description: fmt.Sprintf("%.2f%% of values are %.2f", 100*float64(maxCount)/float64(len(values)), mostCommonValue),
			})
		}
	}

	return anomalies
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MatchesRegex checks if a string matches a regex pattern
func MatchesRegex(s, pattern string) bool {
	// Simple regex matching - could be enhanced
	switch pattern {
	case `^\d{4}-\d{2}-\d{2}$`:
		return len(s) == 10 && s[4] == '-' && s[7] == '-' && IsDigit(s[:4]) && IsDigit(s[5:7]) && IsDigit(s[8:])
	case `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`:
		atPos := strings.Index(s, "@")
		if atPos == -1 {
			return false
		}
		dotPos := strings.LastIndex(s, ".")
		return dotPos > atPos && dotPos < len(s)-1
	case `^\+?\d{1,3}[-.\s]?\d{3}[-.\s]?\d{3}[-.\s]?\d{4}$`:
		if len(s) == 0 {
			return false
		}
		firstChar := s[0]
		if firstChar == '+' {
			return len(s) >= 10 && IsDigit(s[1:])
		}
		return IsDigit(s)
	default:
		return false
	}
}

// IsDigit checks if a string consists only of digits
func IsDigit(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
