package common

import (
	"time"
)

// This file contains additional error telemetry functionality that complements the existing implementation

// GetTelemetryStats returns statistics about error occurrences
func (et *ErrorTelemetry) GetTelemetryStats() map[string]interface{} {
	et.mu.Lock()
	defer et.mu.Unlock()

	stats := make(map[string]interface{})
	stats["error_counts"] = et.ErrorCounts

	// Calculate error frequencies (errors per hour) for the last 24 hours
	now := time.Now()
	hourlyFrequencies := make(map[ErrorCode]float64)

	// This is a simplified version since we don't store timestamps per error in the existing implementation
	// We can estimate based on the total errors and last reported time
	if !et.LastReported.IsZero() {
		hoursSinceLastReport := now.Sub(et.LastReported).Hours()
		if hoursSinceLastReport > 0 {
			for code, count := range et.ErrorCounts {
				hourlyFrequencies[code] = float64(count) / hoursSinceLastReport
			}
		}
	}

	stats["hourly_frequencies"] = hourlyFrequencies
	stats["total_errors"] = et.TotalErrors

	return stats
}

// RecordErrorWithCode records an error with a specific error code
func (et *ErrorTelemetry) RecordErrorWithCode(code ErrorCode) {
	if !et.Enabled {
		return
	}

	et.mu.Lock()
	defer et.mu.Unlock()

	// Increment total errors
	et.TotalErrors++

	// Increment error count for this error code
	et.ErrorCounts[code]++

	// Update last reported time
	et.LastReported = time.Now()

	// Save telemetry data periodically (every 10 errors)
	if et.TotalErrors%10 == 0 {
		et.Save()
	}
}

// GenerateTelemetryReport generates a comprehensive report of error telemetry
func (et *ErrorTelemetry) GenerateTelemetryReport() map[string]interface{} {
	stats := et.GetTelemetryStats()

	// Add additional report data
	stats["report_time"] = time.Now().Format(time.RFC3339)
	stats["anonymous"] = et.Anonymous

	// Calculate most frequent errors
	errorCounts := stats["error_counts"].(map[ErrorCode]int)
	mostFrequent := make([]map[string]interface{}, 0)

	// Convert to slice for sorting
	for code, count := range errorCounts {
		if count > 0 {
			mostFrequent = append(mostFrequent, map[string]interface{}{
				"code":        code,
				"count":       count,
				"description": GetErrorDescription(code),
			})
		}
	}

	// Sort by count (simple bubble sort for small lists)
	for i := 0; i < len(mostFrequent)-1; i++ {
		for j := 0; j < len(mostFrequent)-i-1; j++ {
			if mostFrequent[j]["count"].(int) < mostFrequent[j+1]["count"].(int) {
				mostFrequent[j], mostFrequent[j+1] = mostFrequent[j+1], mostFrequent[j]
			}
		}
	}

	// Take top 5 or less
	if len(mostFrequent) > 5 {
		mostFrequent = mostFrequent[:5]
	}

	stats["most_frequent_errors"] = mostFrequent

	return stats
}
