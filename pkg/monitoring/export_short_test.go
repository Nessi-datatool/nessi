//go:build always
// +build always

package monitoring

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExportMetricsInvalidFormatShort is a short version of the export test
func TestExportMetricsInvalidFormatShort(t *testing.T) {
	// Set up test timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Create a monitor with minimal configuration
	options := MonitorOptions{
		MetricsPort: 9090,
	}
	monitor, err := New(options)
	require.NoError(t, err)

	// Try to export with an invalid format (this should be fast)
	exportOptions := ExportOptions{
		Format:     "invalid_format",
		OutputPath: "test_output.txt",
	}

	_, err = monitor.ExportMetrics(exportOptions)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported export format")
}
