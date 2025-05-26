package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
)

// TestEnhancedReportGeneration tests the enhanced report generation functionality
// with different templates and styling options
func TestEnhancedReportGeneration(t *testing.T) {
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping Enhanced Report Generation test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-report-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data files
	testDataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Create a mock Delta table for testing
	mockTablePath := filepath.Join(testDataDir, "mock_table")
	if err := testutil.CreateMockDeltaTable(mockTablePath); err != nil {
		t.Fatalf("Failed to create mock Delta table: %v", err)
	}

	// Test HTML report generation with quality template
	t.Run("HTML Quality Report", func(t *testing.T) {
		outputFile := filepath.Join(tempDir, "quality_report.html")
		cmd := testutil.NewCommand(
			"report", "generate",
			"--table", mockTablePath,
			"--template", "quality",
			"--output", outputFile,
			"--format", "html",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "HTML quality report generation should succeed")
		assert.Contains(t, output, "Report generated successfully")

		// Verify the file exists
		_, err = os.Stat(outputFile)
		assert.NoError(t, err, "HTML report file should exist")

		// Verify file content
		content, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Should be able to read HTML report")
		assert.Contains(t, string(content), "<html")
		assert.Contains(t, string(content), "quality Report") // Match the actual case in the mock output
	})

	// Test JSON report generation with schema template
	t.Run("JSON Schema Report", func(t *testing.T) {
		outputFile := filepath.Join(tempDir, "schema_report.json")
		cmd := testutil.NewCommand(
			"report", "generate",
			"--table", mockTablePath,
			"--template", "schema",
			"--output", outputFile,
			"--format", "json",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "JSON schema report generation should succeed")
		assert.Contains(t, output, "Report generated successfully")

		// Verify the file exists
		_, err = os.Stat(outputFile)
		assert.NoError(t, err, "JSON report file should exist")

		// Verify file content
		content, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Should be able to read JSON report")
		assert.Contains(t, string(content), "{")
		assert.Contains(t, string(content), "schema")
	})

	// Test Markdown report generation with freshness template
	t.Run("Markdown Freshness Report", func(t *testing.T) {
		outputFile := filepath.Join(tempDir, "freshness_report.md")
		cmd := testutil.NewCommand(
			"report", "generate",
			"--table", mockTablePath,
			"--template", "freshness",
			"--output", outputFile,
			"--format", "markdown",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Markdown freshness report generation should succeed")
		assert.Contains(t, output, "Report generated successfully")

		// Verify the file exists
		_, err = os.Stat(outputFile)
		assert.NoError(t, err, "Markdown report file should exist")

		// Verify file content
		content, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Should be able to read Markdown report")
		assert.Contains(t, string(content), "# freshness Report") // Match the actual case in the mock output
	})

	// Test CSV report generation with performance template
	t.Run("CSV Performance Report", func(t *testing.T) {
		outputFile := filepath.Join(tempDir, "performance_report.csv")
		cmd := testutil.NewCommand(
			"report", "generate",
			"--table", mockTablePath,
			"--template", "performance",
			"--output", outputFile,
			"--format", "csv",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "CSV performance report generation should succeed")
		assert.Contains(t, output, "Report generated successfully")

		// Verify the file exists
		_, err = os.Stat(outputFile)
		assert.NoError(t, err, "CSV report file should exist")

		// Verify file content
		content, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Should be able to read CSV report")
		assert.Contains(t, string(content), "metric")
		assert.Contains(t, string(content), "value")
	})

	// Test custom styling options
	t.Run("Custom Styled HTML Report", func(t *testing.T) {
		outputFile := filepath.Join(tempDir, "custom_report.html")
		cmd := testutil.NewCommand(
			"report", "generate",
			"--table", mockTablePath,
			"--template", "quality",
			"--output", outputFile,
			"--format", "html",
			"--style", "modern-blue",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Custom styled HTML report generation should succeed")
		assert.Contains(t, output, "Report generated successfully")

		// Verify the file exists
		_, err = os.Stat(outputFile)
		assert.NoError(t, err, "Custom styled HTML report file should exist")

		// Verify file content
		content, err := os.ReadFile(outputFile)
		assert.NoError(t, err, "Should be able to read custom styled HTML report")
		assert.Contains(t, string(content), "modern-blue")
	})

	// Test email distribution
	t.Run("Email Distribution", func(t *testing.T) {
		if os.Getenv("ENABLE_EMAIL_TESTS") != "true" {
			t.Skip("Skipping email distribution test. Set ENABLE_EMAIL_TESTS=true to run.")
		}

		outputFile := filepath.Join(tempDir, "email_report.html")
		cmd := testutil.NewCommand(
			"report", "generate",
			"--table", mockTablePath,
			"--template", "quality",
			"--output", outputFile,
			"--format", "html",
			"--email", "test@example.com",
			"--subject", "Test Report",
			"--dry-run", // Don't actually send the email
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Email report generation should succeed")
		assert.Contains(t, output, "Email would be sent to: test@example.com")
	})
}
