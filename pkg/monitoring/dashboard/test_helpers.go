// [REMOVED FOR OSS]: alerting and advanced dashboards are only available in LakeDiff Enterprise.

package dashboard

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/monitoring"
	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi/pkg/testutil"
)

// TestCleanup is a helper function to clean up resources after tests
func TestCleanup() {
	// Set environment variable to indicate we're in test mode
	os.Setenv("GO_TEST", "1")
}

// CloseTestServer is a helper function to close a test server with proper cleanup
func CloseTestServer(server *httptest.Server) {
	if server != nil {
		server.Close()
	}
}

// CleanupTempDir is a helper function to clean up temporary directories
func CleanupTempDir(path string) {
	if path != "" {
		os.RemoveAll(path)
	}
}

// SetupTestTimeout sets up a timeout for tests to prevent them from hanging
func SetupTestTimeout(t *testing.T, timeout time.Duration) {
	t.Helper()
	// No longer calling testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Instead, just set a timeout for the test
	if timeout == 0 {
		timeout = 5 * time.Second // Default timeout
	}
	t.Cleanup(func() {
		// This is just a placeholder for the cleanup function
		// The actual timeout is handled by the Go test runner
	})
}

// RunWithTimeout runs a test function with a timeout
func RunWithTimeout(t *testing.T, testFunc func()) {
	t.Helper()
	testutil.RunWithTimeout(t, testFunc)
}

// ShouldSkipIntegrationTests returns true if integration tests should be skipped
func ShouldSkipIntegrationTests(t *testing.T) bool {
	return testutil.ShouldSkipIntegrationTests(t)
}

// DebugTimer is a helper for tracking time spent in test sections
type DebugTimer struct {
	t        *testing.T
	startTime time.Time
	section   string
	enabled   bool
}

// NewDebugTimer creates a new debug timer
func NewDebugTimer(t *testing.T, section string) *DebugTimer {
	enabled := os.Getenv("NESSI_TEST_DEBUG") != ""
	dt := &DebugTimer{
		t:        t,
		startTime: time.Now(),
		section:   section,
		enabled:   enabled,
	}
	
	if enabled {
		t.Logf("[DEBUG] Starting section: %s", section)
	}
	
	return dt
}

// End stops the timer and logs the elapsed time
func (dt *DebugTimer) End() {
	if !dt.enabled {
		return
	}
	
	elapsed := time.Since(dt.startTime)
	dt.t.Logf("[DEBUG] Section '%s' took %v", dt.section, elapsed)
}

// alertsMockMonitor implements the monitoring.Monitor interface for testing
type alertsMockMonitor struct {
	metricsPort  int
	alertManager *alerts.AlertManager
}

// GetAlertManager returns the alert manager
func (m *alertsMockMonitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

// GetMetricsPort returns the metrics port
func (m *alertsMockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

// RecordMetric implements the Monitor interface
func (m *alertsMockMonitor) RecordMetric(name string, value float64, labels map[string]string) {
	// No-op for testing
}

// ExportMetrics implements the Monitor interface
func (m *alertsMockMonitor) ExportMetrics(options monitoring.ExportOptions) (string, error) {
	return "", nil
}

// GetIntelligentAlertManager implements the Monitor interface
func (m *alertsMockMonitor) GetIntelligentAlertManager() *alerts.IntelligentAlertManager {
	return nil
}

// SkipTest is a no-op function now that security tests are fixed
// It's kept for backward compatibility but doesn't skip tests anymore
func SkipTest(t *testing.T) {
	// No longer skipping tests
}
