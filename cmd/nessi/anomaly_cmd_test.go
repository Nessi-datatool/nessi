package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/quality/anomaly"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnomalyCommands(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "anomaly-cmd-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Helper function to execute a command and capture its output
	executeCommand := func(cmd *cobra.Command, args ...string) (string, error) {
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs(args)
		err := cmd.Execute()
		return buf.String(), err
	}

	// Test generate-sample-data command
	t.Run("GenerateSampleData", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(generateSampleDataCmd)

		// Test generating JSON data
		jsonFile := filepath.Join(tempDir, "sample-spike.json")
		output, err := executeCommand(cmd, "generate-sample-data", jsonFile, "--pattern", "spike", "--points", "10")
		assert.NoError(t, err)
		assert.Contains(t, output, "Generated sample data with spike pattern")
		assert.FileExists(t, jsonFile)

		// Verify the generated JSON data
		jsonData, err := os.ReadFile(jsonFile)
		assert.NoError(t, err)
		var data []map[string]interface{}
		err = json.Unmarshal(jsonData, &data)
		assert.NoError(t, err)
		assert.Equal(t, 10, len(data))

		// Verify the spike
		middleIndex := len(data) / 2
		middleValue := data[middleIndex]["value"].(float64)
		firstValue := data[0]["value"].(float64)
		assert.Greater(t, middleValue, firstValue)

		// Test generating CSV data
		csvFile := filepath.Join(tempDir, "sample-dip.csv")
		output, err = executeCommand(cmd, "generate-sample-data", csvFile, "--pattern", "dip", "--points", "10")
		assert.NoError(t, err)
		assert.Contains(t, output, "Generated sample data with dip pattern")
		assert.FileExists(t, csvFile)

		// Verify the generated CSV data
		csvData, err := os.ReadFile(csvFile)
		assert.NoError(t, err)
		assert.Contains(t, string(csvData), "timestamp,value")

		// Test other patterns
		patterns := []string{"constant-change", "oscillation", "trend-deviation", "complex", "random"}
		for _, pattern := range patterns {
			file := filepath.Join(tempDir, "sample-"+pattern+".json")
			output, err = executeCommand(cmd, "generate-sample-data", file, "--pattern", pattern, "--points", "10")
			assert.NoError(t, err)
			assert.Contains(t, output, "Generated sample data with "+pattern+" pattern")
			assert.FileExists(t, file)
		}
	})

	// Test detect-patterns command
	t.Run("DetectPatterns", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(detectPatternsCmd)

		// Generate sample data files for testing
		patterns := []string{"spike", "dip", "constant-change", "oscillation", "trend-deviation", "complex"}
		for _, pattern := range patterns {
			// Generate sample data
			file := filepath.Join(tempDir, "sample-"+pattern+".json")
			generateCmd := &cobra.Command{Use: "test"}
			generateCmd.AddCommand(generateSampleDataCmd)
			_, err := executeCommand(generateCmd, "generate-sample-data", file, "--pattern", pattern, "--points", "30")
			assert.NoError(t, err)

			// Test detection with text output
			output, err := executeCommand(cmd, "detect-patterns", file)
			assert.NoError(t, err)
			assert.Contains(t, output, "Detected")

			// Test detection with JSON output
			output, err = executeCommand(cmd, "detect-patterns", file, "--format", "json")
			assert.NoError(t, err)

			// Verify JSON output
			var results []*anomaly.PatternDetectionResult
			err = json.Unmarshal([]byte(output), &results)
			assert.NoError(t, err)

			// For specific patterns, verify we detect the expected pattern type
			if pattern == "spike" {
				foundSpike := false
				for _, result := range results {
					if result.Pattern == anomaly.SpikePattern {
						foundSpike = true
						break
					}
				}
				assert.True(t, foundSpike, "Should detect a spike pattern")
			} else if pattern == "dip" {
				foundDip := false
				for _, result := range results {
					if result.Pattern == anomaly.DipPattern {
						foundDip = true
						break
					}
				}
				assert.True(t, foundDip, "Should detect a dip pattern")
			} else if pattern == "constant-change" {
				foundConstantChange := false
				for _, result := range results {
					if result.Pattern == anomaly.ConstantChangePattern {
						foundConstantChange = true
						break
					}
				}
				assert.True(t, foundConstantChange, "Should detect a constant change pattern")
			} else if pattern == "oscillation" {
				foundOscillation := false
				for _, result := range results {
					if result.Pattern == anomaly.OscillationPattern {
						foundOscillation = true
						break
					}
				}
				assert.True(t, foundOscillation, "Should detect an oscillation pattern")
			}
		}

		// Test with custom configuration
		file := filepath.Join(tempDir, "sample-spike.json")
		output, err := executeCommand(cmd, "detect-patterns", file,
			"--min-data-points", "5",
			"--spike-threshold", "0.1",
			"--dip-threshold", "0.1",
			"--constant-change-threshold", "0.01",
			"--oscillation-threshold", "0.05",
			"--trend-deviation-threshold", "0.1",
			"--historical-window-size", "10")
		assert.NoError(t, err)
		assert.Contains(t, output, "Detected")
	})

	// Test readTimeSeriesData function
	t.Run("ReadTimeSeriesData", func(t *testing.T) {
		// Create a JSON file with time series data
		jsonFile := filepath.Join(tempDir, "test-data.json")
		now := time.Now()
		jsonData := []map[string]interface{}{
			{
				"timestamp": now.Format(time.RFC3339),
				"value":     100.0,
			},
			{
				"timestamp": now.Add(time.Hour).Format(time.RFC3339),
				"value":     110.0,
			},
			{
				"timestamp": now.Add(2 * time.Hour).Format(time.RFC3339),
				"value":     120.0,
			},
		}
		jsonBytes, err := json.Marshal(jsonData)
		require.NoError(t, err)
		err = os.WriteFile(jsonFile, jsonBytes, 0644)
		require.NoError(t, err)

		// Test reading JSON data
		data, err := readTimeSeriesData(jsonFile)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(data))
		assert.Equal(t, 100.0, data[0].Value)
		assert.Equal(t, 110.0, data[1].Value)
		assert.Equal(t, 120.0, data[2].Value)

		// Create a CSV file with time series data
		csvFile := filepath.Join(tempDir, "test-data.csv")
		csvContent := "timestamp,value\n" +
			now.Format(time.RFC3339) + ",100.0\n" +
			now.Add(time.Hour).Format(time.RFC3339) + ",110.0\n" +
			now.Add(2*time.Hour).Format(time.RFC3339) + ",120.0\n"
		err = os.WriteFile(csvFile, []byte(csvContent), 0644)
		require.NoError(t, err)

		// Test reading CSV data
		data, err = readTimeSeriesData(csvFile)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(data))
		assert.Equal(t, 100.0, data[0].Value)
		assert.Equal(t, 110.0, data[1].Value)
		assert.Equal(t, 120.0, data[2].Value)

		// Test with invalid file format
		invalidFile := filepath.Join(tempDir, "invalid-data.txt")
		err = os.WriteFile(invalidFile, []byte("invalid data"), 0644)
		require.NoError(t, err)

		_, err = readTimeSeriesData(invalidFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported file format")

		// Test with invalid JSON data
		invalidJsonFile := filepath.Join(tempDir, "invalid-data.json")
		err = os.WriteFile(invalidJsonFile, []byte("invalid json"), 0644)
		require.NoError(t, err)

		_, err = readTimeSeriesData(invalidJsonFile)
		assert.Error(t, err)

		// Test with invalid CSV data
		invalidCsvFile := filepath.Join(tempDir, "invalid-data.csv")
		err = os.WriteFile(invalidCsvFile, []byte("invalid csv"), 0644)
		require.NoError(t, err)

		_, err = readTimeSeriesData(invalidCsvFile)
		assert.Error(t, err)
	})
}
