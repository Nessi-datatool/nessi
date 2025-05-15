package testutil

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

// GetTestTimeout returns an appropriate timeout duration for tests
// Uses a consistent timeout that's fast enough for CI but allows tests to complete
func GetTestTimeout() time.Duration {
	// 2 seconds is a good balance - fast enough for CI but allows tests to complete
	return 2 * time.Second
}

// ShouldSkipIntegrationTests returns true if integration tests should be skipped
// based on environment variables or test flags
func ShouldSkipIntegrationTests(t *testing.T) bool {
	// Skip in short mode (go test -short)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
		return true
	}
	
	// Check for environment variable to skip integration tests
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "1" {
		t.Skip("Skipping integration test due to SKIP_INTEGRATION_TESTS=1")
		return true
	}
	
	return false
}

// RunWithTimeout runs a test function with a timeout
// If no timeout is provided, it uses the default from GetTestTimeout
func RunWithTimeout(t *testing.T, testFunc func(), timeout ...time.Duration) {
	t.Helper() // Mark as test helper for better error reporting

	done := make(chan struct{})
	go func() {
		defer close(done)
		testFunc()
	}()

	// Use provided timeout or default
	actualTimeout := GetTestTimeout()
	if len(timeout) > 0 {
		actualTimeout = timeout[0]
	}

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(actualTimeout):
		t.Fatal("Test timed out after", actualTimeout)
	}
}

// RunTestWithContext runs a test function with a context that has a timeout
// This is useful for tests that need to pass the context to functions being tested
func RunTestWithContext(t *testing.T, testFunc func(ctx context.Context)) {
	t.Helper() // Mark as test helper for better error reporting

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), GetTestTimeout())
	defer cancel()

	// Create a WaitGroup to wait for the test to complete
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the test in a goroutine
	go func() {
		defer wg.Done()
		testFunc(ctx)
	}()

	// Wait for the test to complete or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-ctx.Done():
		t.Fatal("Test timed out after", GetTestTimeout())
	}
}

// RunInParallel marks a test to run in parallel with other tests
// and runs it with a timeout to prevent hanging
func RunInParallel(t *testing.T) {
	t.Helper() // Mark as test helper for better error reporting
	t.Parallel()
	// Note: This doesn't actually run the test, just marks it as parallel
	// You should use RunWithTimeout or RunTestWithContext after this
}
