package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var (
	setupOnce sync.Once
	tempDir   string
)

// TestFixtures provides common test data and utilities
type TestFixtures struct {
	TempDir string
	T       *testing.T
}

// NewTestFixtures creates a new test fixtures instance
func NewTestFixtures(t *testing.T) *TestFixtures {
	setupOnce.Do(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "nessi-test-fixtures-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
	})

	return &TestFixtures{
		TempDir: tempDir,
		T:       t,
	}
}

// CreateTempDir creates a new temporary directory within the fixtures directory
func (f *TestFixtures) CreateTempDir(prefix string) string {
	dir, err := os.MkdirTemp(f.TempDir, prefix)
	if err != nil {
		f.T.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// CreateMetricsData creates test metrics data for trend analysis
func (f *TestFixtures) CreateMetricsData(metricsDir, field string, version int, mean, stdDev, min, max float64) {
	// Create run metrics
	runMetrics := map[string]interface{}{
		"timestamp": time.Now().Add(-time.Duration(version) * 24 * time.Hour),
		"field_metrics": map[string]map[string]float64{
			field: {
				"mean":         mean,
				"standard_dev": stdDev,
				"min":          min,
				"max":          max,
				"record_count": 100.0,
			},
		},
	}

	// Save to file
	filename := filepath.Join(metricsDir, filepath.Clean(field), filepath.Clean(field+"_"+time.Now().Format("20060102")+".json"))
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		f.T.Fatalf("Failed to create metrics directory: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		f.T.Fatalf("Failed to create metrics file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(runMetrics); err != nil {
		f.T.Fatalf("Failed to write metrics: %v", err)
	}
}

// Cleanup removes all test fixtures
func (f *TestFixtures) Cleanup() {
	if err := os.RemoveAll(f.TempDir); err != nil {
		f.T.Errorf("Failed to cleanup test fixtures: %v", err)
	}
}
