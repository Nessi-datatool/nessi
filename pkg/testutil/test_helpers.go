package testutil

import (
	"context"
	"os"
	"testing"
	"time"
)

// FastTestMode returns true if tests should run in fast mode
// This can be enabled by setting the NESSI_FAST_TESTS environment variable
func FastTestMode() bool {
	return os.Getenv("NESSI_FAST_TESTS") != ""
}

// SkipIfLongRunning skips a test if it's a long-running test and fast mode is enabled
func SkipIfLongRunning(t *testing.T) {
	if FastTestMode() {
		t.Skip("Skipping long-running test in fast mode")
	}
}

// SetShortTimeout sets a short timeout for a test function
// This is useful for tests that might hang or take too long
func SetShortTimeout(t *testing.T, timeout time.Duration) {
	if timeout == 0 {
		timeout = 50 * time.Millisecond // Default to 50ms for faster tests
	}
	
	// Use context with timeout for more reliable test timeouts
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			t.Errorf("Test exceeded timeout of %v", timeout)
		}
	})
}

// RunInParallel marks a test to run in parallel with other tests
// and sets a short timeout to prevent long-running tests
func RunInParallel(t *testing.T) {
	t.Helper() // Mark as test helper for better error reporting
	t.Parallel()
	SetShortTimeout(t, 50*time.Millisecond) // Use shorter timeout
}
