package dbt

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Helper function for tests
func now() time.Time {
	return time.Now()
}

func TestGenerateDBTArtifacts(t *testing.T) {
	// Create a temporary directory for test artifacts
	tempDir, err := os.MkdirTemp("", "dbt-artifacts-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with quality_score artifact type

	artifactsPath := filepath.Join(tempDir, "artifacts")
	err = GenerateDBTArtifacts("quality_score", artifactsPath)
	if err != nil {
		t.Fatalf("GenerateDBTArtifacts() error = %v", err)
	}

	// Test with different artifact types

	// Test with invalid results type
	err = GenerateDBTArtifacts("invalid", artifactsPath)
	if err == nil {
		t.Errorf("GenerateDBTArtifacts() with invalid type should return an error")
	}
}
