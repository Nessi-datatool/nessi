//go:build fast_tests

package security

import (
	"time"
)

// FastTestTimeout returns a shorter timeout for fast tests
func FastTestTimeout() time.Duration {
	return 2 * time.Second
}

// FastSkipIntegrationTests returns true if integration tests should be skipped in fast mode
func FastSkipIntegrationTests() bool {
	return true
}
