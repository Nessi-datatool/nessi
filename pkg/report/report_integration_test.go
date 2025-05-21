package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Import the report types package to avoid name conflicts
import reporttypes "github.com/nessi-dev/nessi/pkg/api/types/report"

// TestReportIntegration tests the report generation functionality end-to-end
func TestReportIntegration(t *testing.T) {
	// Skip this test in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directories for templates and output
	templatesDir, err := os.MkdirTemp("", "report-templates")
	require.NoError(t, err)
	defer os.RemoveAll(templatesDir)

	outputDir, err := os.MkdirTemp("", "report-output")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// Create a report generator
	generator, err := NewReportGenerator(templatesDir, outputDir)
	require.NoError(t, err)

	// Create test data
	scanResult := &reporttypes.ScanResult{
		ScanID:      "test-scan-123",
		TablePath:   "/data/test_table",
		RowCount:    1000,
		ColumnCount: 5,
		Timestamp:   time.Now(),
		Duration:    2 * time.Second,
	}

	// Generate HTML report
	htmlOutputPath, err := generator.GenerateScanReport(scanResult, HTML, "")
	require.NoError(t, err)
	assert.NotEmpty(t, htmlOutputPath)

	// Verify HTML report content
	htmlContent, err := os.ReadFile(htmlOutputPath)
	require.NoError(t, err)
	htmlContentStr := string(htmlContent)

	// Check that the HTML report contains the expected data
	assert.Contains(t, htmlContentStr, "test-scan-123")
	assert.Contains(t, htmlContentStr, "/data/test_table")
	assert.Contains(t, htmlContentStr, "1000")
	assert.Contains(t, htmlContentStr, "5")

	// Generate JSON report
	jsonOutputPath, err := generator.GenerateScanReport(scanResult, JSON, "")
	require.NoError(t, err)
	assert.NotEmpty(t, jsonOutputPath)

	// Verify JSON report content
	jsonContent, err := os.ReadFile(jsonOutputPath)
	require.NoError(t, err)
	jsonContentStr := string(jsonContent)

	// Check that the JSON report contains the expected data
	assert.Contains(t, jsonContentStr, "test-scan-123")
	assert.Contains(t, jsonContentStr, "/data/test_table")
	assert.Contains(t, jsonContentStr, "1000")
	assert.Contains(t, jsonContentStr, "5")

	// Test with custom output path
	customOutputPath := filepath.Join(outputDir, "custom_report.html")
	outputPath, err := generator.GenerateScanReport(scanResult, HTML, customOutputPath)
	require.NoError(t, err)
	assert.Equal(t, customOutputPath, outputPath)

	// Verify the custom output file exists
	_, err = os.Stat(customOutputPath)
	assert.NoError(t, err)

	// Test error cases
	// 1. Invalid report format
	_, err = generator.GenerateScanReport(scanResult, "invalid-format", "")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unsupported report format"))

	// 2. Nil scan result
	_, err = generator.GenerateScanReport(nil, HTML, "")
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "scan result is nil"))
}
