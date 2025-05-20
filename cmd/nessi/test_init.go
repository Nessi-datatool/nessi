package main

import (
	"flag"
	"os"
	"testing"
)

var testMode bool = true

// TestMain is used to set up any global test configuration
func TestMain(m *testing.M) {
	// Reset flags before running tests to avoid redefinition errors
	flag.CommandLine = flag.NewFlagSet("nessi-test", flag.ExitOnError)

	// Create a clean root command for tests to avoid flag conflicts
	origRootCmd := rootCmd
	defer func() { rootCmd = origRootCmd }()

	// Run all tests
	code := m.Run()

	// Exit with the test status code
	os.Exit(code)
}
