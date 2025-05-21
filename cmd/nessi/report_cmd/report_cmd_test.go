package report_cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportCommand(t *testing.T) {
	// Test that the command is created correctly
	cmd := NewReportCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "report [type] [data_path]", cmd.Use)
	// Verify the command requires exactly 2 arguments
	args := []string{"quality", "/test/path"}
	err := cmd.Args(cmd, args)
	assert.NoError(t, err)

	// Test that the flags are set correctly
	outputFlag := cmd.Flag("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "o", outputFlag.Shorthand)

	formatFlag := cmd.Flag("format")
	assert.NotNil(t, formatFlag)
	assert.Equal(t, "f", formatFlag.Shorthand)
	assert.Equal(t, string(HTML), formatFlag.DefValue)

	templatesFlag := cmd.Flag("templates")
	assert.NotNil(t, templatesFlag)
	assert.Equal(t, "t", templatesFlag.Shorthand)
}

func TestGenerateMockData(t *testing.T) {
	// Test that the mock data is generated correctly for each report type
	qualityData := generateMockData("quality", "/test/path")
	assert.NotNil(t, qualityData)
	qualityMap, ok := qualityData.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "/test/path", qualityMap["table_path"])

	schemaData := generateMockData("schema", "/test/path")
	assert.NotNil(t, schemaData)
	schemaMap, ok := schemaData.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "/test/path", schemaMap["table_path"])

	freshnessData := generateMockData("freshness", "/test/path")
	assert.NotNil(t, freshnessData)
	freshnessMap, ok := freshnessData.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "/test/path", freshnessMap["table_path"])

	performanceData := generateMockData("performance", "/test/path")
	assert.NotNil(t, performanceData)
	performanceMap, ok := performanceData.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "/test/path", performanceMap["table_path"])

	// Test with an invalid report type
	invalidData := generateMockData("invalid", "/test/path")
	assert.NotNil(t, invalidData)
	invalidMap, ok := invalidData.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Unknown report type", invalidMap["error"])
}

func TestReportCommandExecution(t *testing.T) {
	// Create temporary directories for testing
	templateDir, err := os.MkdirTemp("", "report-templates")
	require.NoError(t, err)
	defer os.RemoveAll(templateDir)

	outputDir, err := os.MkdirTemp("", "report-output")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	// Create a simple HTML template for testing
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Test Report</title>
</head>
<body>
    <h1>Test Report</h1>
    <p>Table Path: {{ .table_path }}</p>
</body>
</html>`

	// Create template directories and files
	qualityTemplateDir := filepath.Join(templateDir, "quality")
	err = os.MkdirAll(qualityTemplateDir, 0755)
	require.NoError(t, err)

	qualityTemplatePath := filepath.Join(qualityTemplateDir, "quality_report.html")
	err = os.WriteFile(qualityTemplatePath, []byte(htmlTemplate), 0644)
	require.NoError(t, err)

	// Test the command execution (we can't actually execute the command in a unit test,
	// but we can test the individual functions it calls)
	data := generateQualityReportData("/test/path")
	assert.NotNil(t, data)
	assert.Equal(t, "/test/path", data["table_path"])

	// Verify that the schema report data is generated correctly
	schemaData := generateSchemaReportData("/test/path")
	assert.NotNil(t, schemaData)
	assert.Equal(t, "/test/path", schemaData["table_path"])

	// Verify that the freshness report data is generated correctly
	freshnessData := generateFreshnessReportData("/test/path")
	assert.NotNil(t, freshnessData)
	assert.Equal(t, "/test/path", freshnessData["table_path"])

	// Verify that the performance report data is generated correctly
	performanceData := generatePerformanceReportData("/test/path")
	assert.NotNil(t, performanceData)
	assert.Equal(t, "/test/path", performanceData["table_path"])
}
