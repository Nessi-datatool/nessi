package main

import (
	"os"
	"testing"

	"github.com/nessi-dev/nessi-dev/cmd/nessi/testing"
)

// TestMain is the main entry point for all tests in the cmd/nessi package
func TestMain(m *testing.M) {
	// Set test environment variable
	os.Setenv("NESSI_TEST_MODE", "true")

	// Reset flags to avoid conflicts
	testing.ResetFlags()

	// Run tests
	code := m.Run()

	// Unset test environment variable
	os.Unsetenv("NESSI_TEST_MODE")

	os.Exit(code)
}
