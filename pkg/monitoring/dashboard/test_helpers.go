package dashboard

import (
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

// TestCleanup is a helper function to clean up resources after tests
// It sets a short timeout for any goroutines that might be waiting
func TestCleanup() {
	// Set environment variable to indicate we're in test mode
	os.Setenv("GO_TEST", "1")
	
	// Set a default value for test timeout if not already set
	if os.Getenv("NESSI_TEST_TIMEOUT_MS") == "" {
		os.Setenv("NESSI_TEST_TIMEOUT_MS", "500") // 0.5 seconds default
	}
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
func SetupTestTimeout(t *testing.T, defaultDuration time.Duration) {
	t.Helper()
	
	// Get timeout from environment variable or use default
	duration := defaultDuration
	if timeoutStr := os.Getenv("NESSI_TEST_TIMEOUT_MS"); timeoutStr != "" {
		if timeoutMs, err := strconv.Atoi(timeoutStr); err == nil {
			duration = time.Duration(timeoutMs) * time.Millisecond
		}
	}
	
	// Create a channel to signal test completion
	done := make(chan bool)
	
	// Start a goroutine that will fail the test if it takes too long
	go func() {
		select {
		case <-done:
			// Test completed normally
			return
		case <-time.After(duration):
			// Test took too long, fail it
			t.Error("Test timed out after", duration)
			t.FailNow()
		}
	}()
	
	// Register cleanup to signal completion
	t.Cleanup(func() {
		close(done)
	})
}

// ShouldSkipLongTests returns true if long tests should be skipped
func ShouldSkipLongTests() bool {
	// Check if NESSI_SKIP_LONG_TESTS is set to true/1/yes
	skipStr := os.Getenv("NESSI_SKIP_LONG_TESTS")
	if skipStr == "" {
		return false
	}
	
	// Convert to boolean
	skip, err := strconv.ParseBool(skipStr)
	if err != nil {
		// If we can't parse it, default to not skipping
		return false
	}
	
	return skip
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
