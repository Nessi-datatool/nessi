package security

import (
	"os"
	"testing"
)

// ShouldSkipIntegrationTests returns true if integration tests should be skipped
// based on environment variables or test flags
func ShouldSkipIntegrationTests(t *testing.T) bool {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
		return true
	}
	
	// Check for environment variable to skip security integration tests
	if os.Getenv("SKIP_SECURITY_INTEGRATION") == "1" {
		t.Skip("Skipping security integration test due to SKIP_SECURITY_INTEGRATION=1")
		return true
	}
	
	return false
}
