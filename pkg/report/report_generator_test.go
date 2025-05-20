package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportGenerator_GenerateHTMLReport(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "report-generator-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a report generator
	generator, err := NewReportGenerator(tempDir, tempDir)
	require.NoError(t, err)

	// Create test data
	data := map[string]interface{}{
		"title":     "Test Report",
		"timestamp": time.Now().Format(time.RFC3339),
		"metrics": map[string]interface{}{
			"completeness": 0.95,
			"accuracy":     0.98,
			"consistency":  0.92,
		},
	}

	// Generate an HTML report
	outputPath := filepath.Join(tempDir, "test-report.html")
	resultPath, err := generator.GenerateReport(data, "test", HTML, outputPath)
	require.NoError(t, err)
	assert.Equal(t, outputPath, resultPath)

	// Check that the file was created
	_, err = os.Stat(outputPath)
	require.NoError(t, err)

	// Check the content of the file
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "<html>")
	assert.Contains(t, string(content), "Test Report")
}

func TestReportGenerator_GenerateJSONReport(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "report-generator-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a report generator
	generator, err := NewReportGenerator(tempDir, tempDir)
	require.NoError(t, err)

	// Create test data
	data := map[string]interface{}{
		"title":     "Test Report",
		"timestamp": time.Now().Format(time.RFC3339),
		"metrics": map[string]interface{}{
			"completeness": 0.95,
			"accuracy":     0.98,
			"consistency":  0.92,
		},
	}

	// Generate a JSON report
	outputPath := filepath.Join(tempDir, "test-report.json")
	resultPath, err := generator.GenerateReport(data, "test", JSON, outputPath)
	require.NoError(t, err)
	assert.Equal(t, outputPath, resultPath)

	// Check that the file was created
	_, err = os.Stat(outputPath)
	require.NoError(t, err)

	// Check the content of the file
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	// Unmarshal the JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	require.NoError(t, err)

	// Check the data
	assert.Equal(t, "Test Report", jsonData["title"])
	assert.NotNil(t, jsonData["timestamp"])
	assert.NotNil(t, jsonData["metrics"])
}

func TestReportGenerator_GenerateCSVReport(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "report-generator-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a report generator
	generator, err := NewReportGenerator(tempDir, tempDir)
	require.NoError(t, err)

	// Create test data
	data := []map[string]interface{}{
		{
			"id":     1,
			"name":   "John Doe",
			"age":    30,
			"active": true,
		},
		{
			"id":     2,
			"name":   "Jane Smith",
			"age":    25,
			"active": true,
		},
		{
			"id":     3,
			"name":   "Bob Johnson",
			"age":    40,
			"active": false,
		},
	}

	// Generate a CSV report
	outputPath := filepath.Join(tempDir, "test-report.csv")
	resultPath, err := generator.GenerateReport(data, "test", CSV, outputPath)
	require.NoError(t, err)
	assert.Equal(t, outputPath, resultPath)

	// Check that the file was created
	_, err = os.Stat(outputPath)
	require.NoError(t, err)

	// Check the content of the file
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	// Check the CSV content
	assert.Contains(t, string(content), "name")
	assert.Contains(t, string(content), "age")
	assert.Contains(t, string(content), "active")
	assert.Contains(t, string(content), "id")
	assert.Contains(t, string(content), "John Doe")
	assert.Contains(t, string(content), "Jane Smith")
	assert.Contains(t, string(content), "Bob Johnson")
}

func TestReportGenerator_GenerateScanReport(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "report-generator-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a report generator
	generator, err := NewReportGenerator(tempDir, tempDir)
	require.NoError(t, err)

	// Create a scan result
	scanResult := &report.ScanResult{
		ScanID:      "test-scan-id",
		TablePath:   "/path/to/table",
		Timestamp:   time.Now(),
		Duration:    time.Second * 10,
		RowCount:    1000,
		ColumnCount: 5,
		QualityMetrics: &report.QualityMetrics{
			Completeness: map[string]float64{"overall": 0.95},
			Accuracy:     map[string]float64{"overall": 0.98},
			Consistency:  map[string]float64{"overall": 0.92},
			Uniqueness:   map[string]float64{"overall": 0.99},
			Timeliness:   map[string]float64{"overall": 0.90},
		},
		PerformanceMetrics: &report.PerformanceMetrics{
			ScanDurationMs:  10000,
			MemoryUsageMb:   100,
			CpuUsagePercent: 50,
			IoOperations:    1000,
			RowsProcessed:   1000,
			BytesProcessed:  1000000,
			StartTime:       time.Now().Add(-time.Second * 10),
			EndTime:         time.Now(),
		},
	}

	// Generate a scan report in JSON format
	outputPath := filepath.Join(tempDir, "scan-report.json")
	resultPath, err := generator.GenerateScanReport(scanResult, JSON, outputPath)
	require.NoError(t, err)
	assert.Equal(t, outputPath, resultPath)

	// Check that the file was created
	_, err = os.Stat(outputPath)
	require.NoError(t, err)

	// Check the content of the file
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	// Unmarshal the JSON
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	require.NoError(t, err)

	// Check the data
	assert.Equal(t, "test-scan-id", jsonData["scan_id"])
	assert.Equal(t, "/path/to/table", jsonData["table_path"])
	assert.NotNil(t, jsonData["quality_metrics"])
	assert.NotNil(t, jsonData["performance_metrics"])
}

func TestValidateReportData(t *testing.T) {
	// Test valid scan result
	scanResult := &report.ScanResult{
		ScanID:      "test-scan-id",
		TablePath:   "/path/to/table",
		Timestamp:   time.Now(),
		Duration:    time.Second * 10,
		RowCount:    1000,
		ColumnCount: 5,
		QualityMetrics: &report.QualityMetrics{
			Completeness: map[string]float64{"overall": 0.95},
			Accuracy:     map[string]float64{"overall": 0.98},
			Consistency:  map[string]float64{"overall": 0.92},
			Uniqueness:   map[string]float64{"overall": 0.99},
			Timeliness:   map[string]float64{"overall": 0.90},
		},
		PerformanceMetrics: &report.PerformanceMetrics{
			ScanDurationMs:  10000,
			MemoryUsageMb:   100,
			CpuUsagePercent: 50,
			IoOperations:    1000,
			RowsProcessed:   1000,
			BytesProcessed:  1000000,
			StartTime:       time.Now().Add(-time.Second * 10),
			EndTime:         time.Now(),
		},
	}

	// Validate the scan result
	err := ValidateReportData(scanResult)
	require.NoError(t, err)

	// Test invalid scan result
	invalidScanResult := &report.ScanResult{
		// Missing required fields
	}

	// Validate the invalid scan result
	err = ValidateReportData(invalidScanResult)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "table path is required")

	// Test valid quality metrics
	qualityMetrics := &report.QualityMetrics{
		Completeness: map[string]float64{"overall": 0.95},
		Accuracy:     map[string]float64{"overall": 0.98},
		Consistency:  map[string]float64{"overall": 0.92},
		Uniqueness:   map[string]float64{"overall": 0.99},
		Timeliness:   map[string]float64{"overall": 0.90},
	}

	// Validate the quality metrics
	err = ValidateReportData(qualityMetrics)
	require.NoError(t, err)

	// Test invalid quality metrics
	invalidQualityMetrics := &report.QualityMetrics{
		// Missing required fields
	}

	// Validate the invalid quality metrics
	err = ValidateReportData(invalidQualityMetrics)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "completeness metrics are required")

	// Test valid performance metrics
	performanceMetrics := &report.PerformanceMetrics{
		ScanDurationMs:  10000,
		MemoryUsageMb:   100,
		CpuUsagePercent: 50,
		IoOperations:    1000,
		RowsProcessed:   1000,
		BytesProcessed:  1000000,
		StartTime:       time.Now().Add(-time.Second * 10),
		EndTime:         time.Now(),
	}

	// Validate the performance metrics
	err = ValidateReportData(performanceMetrics)
	require.NoError(t, err)

	// Test invalid performance metrics
	invalidPerformanceMetrics := &report.PerformanceMetrics{
		// Missing required fields
	}

	// Validate the invalid performance metrics
	err = ValidateReportData(invalidPerformanceMetrics)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scan duration is required")
}
