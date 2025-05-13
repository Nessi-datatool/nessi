package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/nessi-dev/nessi-dev/pkg/quality/anomaly"
)

// anomalyCmd represents the anomaly command
var anomalyCmd = &cobra.Command{
	Use:   "anomaly",
	Short: "Anomaly detection operations",
	Long:  `Perform anomaly detection operations on time series data, including pattern detection and outlier identification.`,
}

// detectPatternsCmd represents the detect-patterns command
var detectPatternsCmd = &cobra.Command{
	Use:   "detect-patterns [data_file]",
	Short: "Detect patterns in time series data",
	Long:  `Detect patterns in time series data, including constant changes, spikes, dips, oscillations, and trend deviations.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dataFile := args[0]
		
		// Read data from file
		data, err := readTimeSeriesData(dataFile)
		if err != nil {
			fmt.Printf("Error reading data: %v\n", err)
			os.Exit(1)
		}
		
		// Create pattern detector
		config := anomaly.DefaultPatternDetectionConfig()
		
		// Override config with command line flags
		if minDataPoints, _ := cmd.Flags().GetInt("min-data-points"); minDataPoints > 0 {
			config.MinDataPoints = minDataPoints
		}
		
		if constantChangeThreshold, _ := cmd.Flags().GetFloat64("constant-change-threshold"); constantChangeThreshold > 0 {
			config.ConstantChangeThreshold = constantChangeThreshold
		}
		
		if spikeThreshold, _ := cmd.Flags().GetFloat64("spike-threshold"); spikeThreshold > 0 {
			config.SpikeThreshold = spikeThreshold
		}
		
		if dipThreshold, _ := cmd.Flags().GetFloat64("dip-threshold"); dipThreshold > 0 {
			config.DipThreshold = dipThreshold
		}
		
		if oscillationThreshold, _ := cmd.Flags().GetFloat64("oscillation-threshold"); oscillationThreshold > 0 {
			config.OscillationThreshold = oscillationThreshold
		}
		
		if trendDeviationThreshold, _ := cmd.Flags().GetFloat64("trend-deviation-threshold"); trendDeviationThreshold > 0 {
			config.TrendDeviationThreshold = trendDeviationThreshold
		}
		
		if historicalWindowSize, _ := cmd.Flags().GetInt("historical-window-size"); historicalWindowSize > 0 {
			config.HistoricalWindowSize = historicalWindowSize
		}
		
		detector := anomaly.NewPatternDetector(config)
		
		// Detect patterns
		patterns := detector.DetectPatterns(data)
		
		// Get output format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(patterns, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling patterns to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		if len(patterns) == 0 {
			fmt.Println("No patterns detected in the data.")
			return
		}
		
		fmt.Printf("Detected %d patterns in the data:\n\n", len(patterns))
		
		for i, pattern := range patterns {
			fmt.Printf("%d. Pattern: %s\n", i+1, pattern.Pattern)
			fmt.Printf("   Description: %s\n", pattern.Description)
			fmt.Printf("   Confidence: %.2f\n", pattern.Confidence)
			fmt.Printf("   Magnitude: %.2f%%\n", pattern.Magnitude*100)
			fmt.Printf("   Range: [%d, %d]\n", pattern.StartIndex, pattern.EndIndex)
			fmt.Println()
		}
		
		// Print summary
		patternCounts := make(map[anomaly.PatternType]int)
		for _, pattern := range patterns {
			patternCounts[pattern.Pattern]++
		}
		
		fmt.Println("Summary:")
		for patternType, count := range patternCounts {
			fmt.Printf("- %s: %d\n", patternType, count)
		}
	},
}

// readTimeSeriesData reads time series data from a file
func readTimeSeriesData(filePath string) ([]anomaly.TimeSeriesDataPoint, error) {
	// Read file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Check file extension to determine format
	if strings.HasSuffix(filePath, ".json") {
		// Parse JSON
		var jsonData []map[string]interface{}
		if err := json.Unmarshal(fileData, &jsonData); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		
		// Convert to TimeSeriesDataPoint
		result := make([]anomaly.TimeSeriesDataPoint, 0, len(jsonData))
		for _, item := range jsonData {
			// Get timestamp
			var timestamp time.Time
			if ts, ok := item["timestamp"].(string); ok {
				timestamp, err = time.Parse(time.RFC3339, ts)
				if err != nil {
					return nil, fmt.Errorf("failed to parse timestamp: %w", err)
				}
			} else {
				return nil, fmt.Errorf("missing or invalid timestamp field")
			}
			
			// Get value
			var value float64
			switch v := item["value"].(type) {
			case float64:
				value = v
			case int:
				value = float64(v)
			case string:
				value, err = strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("failed to parse value: %w", err)
				}
			default:
				return nil, fmt.Errorf("missing or invalid value field")
			}
			
			result = append(result, anomaly.TimeSeriesDataPoint{
				Timestamp: timestamp,
				Value:     value,
			})
		}
		
		return result, nil
	} else if strings.HasSuffix(filePath, ".csv") {
		// Parse CSV
		lines := strings.Split(string(fileData), "\n")
		
		// Skip header
		if len(lines) <= 1 {
			return nil, fmt.Errorf("CSV file is empty or has only a header")
		}
		
		// Convert to TimeSeriesDataPoint
		result := make([]anomaly.TimeSeriesDataPoint, 0, len(lines)-1)
		for i, line := range lines {
			if i == 0 || line == "" {
				continue
			}
			
			fields := strings.Split(line, ",")
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d has fewer than 2 fields", i+1)
			}
			
			// Parse timestamp
			timestamp, err := time.Parse(time.RFC3339, strings.TrimSpace(fields[0]))
			if err != nil {
				return nil, fmt.Errorf("failed to parse timestamp on line %d: %w", i+1, err)
			}
			
			// Parse value
			value, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse value on line %d: %w", i+1, err)
			}
			
			result = append(result, anomaly.TimeSeriesDataPoint{
				Timestamp: timestamp,
				Value:     value,
			})
		}
		
		return result, nil
	}
	
	return nil, fmt.Errorf("unsupported file format: must be .json or .csv")
}

// generateSampleDataCmd represents the generate-sample-data command
var generateSampleDataCmd = &cobra.Command{
	Use:   "generate-sample-data [output_file]",
	Short: "Generate sample time series data",
	Long:  `Generate sample time series data with various patterns for testing.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputFile := args[0]
		
		// Get pattern type
		patternType, _ := cmd.Flags().GetString("pattern")
		
		// Get number of points
		numPoints, _ := cmd.Flags().GetInt("points")
		if numPoints <= 0 {
			numPoints = 30
		}
		
		// Generate data
		data := generateSampleData(patternType, numPoints)
		
		// Get output format
		format := "json"
		if strings.HasSuffix(outputFile, ".csv") {
			format = "csv"
		}
		
		// Write to file
		if format == "json" {
			// Convert to JSON-friendly format
			jsonData := make([]map[string]interface{}, len(data))
			for i, point := range data {
				jsonData[i] = map[string]interface{}{
					"timestamp": point.Timestamp.Format(time.RFC3339),
					"value":     point.Value,
				}
			}
			
			// Marshal to JSON
			fileData, err := json.MarshalIndent(jsonData, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling data to JSON: %v\n", err)
				os.Exit(1)
			}
			
			// Write to file
			if err := os.WriteFile(outputFile, fileData, 0644); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Convert to CSV
			var sb strings.Builder
			sb.WriteString("timestamp,value\n")
			
			for _, point := range data {
				sb.WriteString(fmt.Sprintf("%s,%.2f\n", point.Timestamp.Format(time.RFC3339), point.Value))
			}
			
			// Write to file
			if err := os.WriteFile(outputFile, []byte(sb.String()), 0644); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				os.Exit(1)
			}
		}
		
		fmt.Printf("Generated sample data with %s pattern and wrote to %s\n", patternType, outputFile)
	},
}

// generateSampleData generates sample time series data with the specified pattern
func generateSampleData(patternType string, numPoints int) []anomaly.TimeSeriesDataPoint {
	result := make([]anomaly.TimeSeriesDataPoint, numPoints)
	now := time.Now()
	
	// Base value
	baseValue := 100.0
	
	switch patternType {
	case "constant-change":
		// Generate data with constant increase
		for i := 0; i < numPoints; i++ {
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     baseValue + float64(i)*5.0, // 5% increase each time
			}
		}
	case "spike":
		// Generate data with a spike in the middle
		spikeIndex := numPoints / 2
		for i := 0; i < numPoints; i++ {
			value := baseValue
			if i == spikeIndex {
				value = baseValue * 1.5 // 50% spike
			}
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     value,
			}
		}
	case "dip":
		// Generate data with a dip in the middle
		dipIndex := numPoints / 2
		for i := 0; i < numPoints; i++ {
			value := baseValue
			if i == dipIndex {
				value = baseValue * 0.5 // 50% dip
			}
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     value,
			}
		}
	case "oscillation":
		// Generate data with oscillation
		for i := 0; i < numPoints; i++ {
			value := baseValue
			if i%2 == 0 {
				value = baseValue * 1.2 // 20% oscillation
			} else {
				value = baseValue * 0.8
			}
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     value,
			}
		}
	case "trend-deviation":
		// Generate data with a trend and then a deviation
		for i := 0; i < numPoints; i++ {
			value := baseValue + float64(i)*2.0 // Linear trend
			if i == numPoints-1 {
				value = baseValue + float64(i)*2.0 * 1.5 // 50% deviation
			}
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     value,
			}
		}
	case "complex":
		// Generate data with multiple patterns
		for i := 0; i < numPoints; i++ {
			// Base value with linear trend
			value := baseValue + float64(i)*2.0
			
			// Add some noise (±5%)
			noise := value * (0.05 * (float64(i%3) - 1.0))
			value += noise
			
			// Add a spike around 1/3 of the way through
			if i == numPoints/3 {
				value *= 1.3 // 30% spike
			}
			
			// Add a dip around 2/3 of the way through
			if i == 2*numPoints/3 {
				value *= 0.7 // 30% dip
			}
			
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     value,
			}
		}
	default:
		// Generate random data
		for i := 0; i < numPoints; i++ {
			// Random value between 90 and 110
			value := baseValue * (0.9 + 0.2*float64(i%10)/10.0)
			result[i] = anomaly.TimeSeriesDataPoint{
				Timestamp: now.Add(time.Duration(i) * time.Hour),
				Value:     value,
			}
		}
	}
	
	return result
}

func init() {
	rootCmd.AddCommand(anomalyCmd)
	
	// Add subcommands
	anomalyCmd.AddCommand(detectPatternsCmd)
	anomalyCmd.AddCommand(generateSampleDataCmd)
	
	// Add flags for detect-patterns
	detectPatternsCmd.Flags().String("format", "text", "Output format (text, json)")
	detectPatternsCmd.Flags().Int("min-data-points", 0, "Minimum number of data points required for detection")
	detectPatternsCmd.Flags().Float64("constant-change-threshold", 0, "Minimum percentage change for constant change detection")
	detectPatternsCmd.Flags().Float64("spike-threshold", 0, "Minimum percentage change for spike detection")
	detectPatternsCmd.Flags().Float64("dip-threshold", 0, "Minimum percentage change for dip detection")
	detectPatternsCmd.Flags().Float64("oscillation-threshold", 0, "Minimum percentage change for oscillation detection")
	detectPatternsCmd.Flags().Float64("trend-deviation-threshold", 0, "Minimum percentage deviation from trend")
	detectPatternsCmd.Flags().Int("historical-window-size", 0, "Number of data points to use for trend analysis")
	
	// Add flags for generate-sample-data
	generateSampleDataCmd.Flags().String("pattern", "complex", "Pattern type (constant-change, spike, dip, oscillation, trend-deviation, complex, random)")
	generateSampleDataCmd.Flags().Int("points", 30, "Number of data points to generate")
}
