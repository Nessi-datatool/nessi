package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

func main() {
	fmt.Println("Consistency Check Demo")
	fmt.Println("=====================")
	fmt.Println()

	// Create a temporary directory for the demo
	tempDir, err := os.MkdirTemp("", "delta-consistency-demo-*")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	fmt.Printf("Created temporary Delta table at: %s\n", tempDir)
	fmt.Println()

	// Create a Delta table with multiple versions
	fmt.Println("Creating a Delta table with multiple versions...")
	err = createDeltaTableWithVersions(tempDir)
	if err != nil {
		fmt.Printf("Error creating Delta table: %v\n", err)
		return
	}

	// Create metadata manager
	manager := datalake.NewMetadataManager(tempDir)

	// Get version history
	fmt.Println("Getting version history...")
	history, err := manager.GetVersionHistory()
	if err != nil {
		fmt.Printf("Error getting version history: %v\n", err)
		return
	}

	fmt.Println("Version History:")
	fmt.Println("----------------")
	for _, version := range history {
		fmt.Printf("Version: %d\n", version.Version)
		// Convert timestamp to time.Time
		timestamp := time.Unix(version.Timestamp, 0)
		fmt.Printf("Timestamp: %s\n", timestamp.Format(time.RFC3339))
		fmt.Printf("Operation: %s\n", version.Operation)
		fmt.Println("----------------")
	}
	fmt.Println()

	// Run schema consistency check
	fmt.Println("Running schema consistency check...")
	schemaResult, err := manager.CheckConsistency(datalake.SchemaConsistencyCheck, []int64{0, 1, 2, 3})
	if err != nil {
		fmt.Printf("Error running schema consistency check: %v\n", err)
		return
	}
	printConsistencyResult(schemaResult)
	fmt.Println()

	// Run type consistency check
	fmt.Println("Running type consistency check...")
	typeResult, err := manager.CheckConsistency(datalake.TypeConsistencyCheck, []int64{0, 1, 2, 3})
	if err != nil {
		fmt.Printf("Error running type consistency check: %v\n", err)
		return
	}
	printConsistencyResult(typeResult)
	fmt.Println()

	// Run field consistency check
	fmt.Println("Running field consistency check...")
	fieldResult, err := manager.CheckConsistency(datalake.FieldConsistencyCheck, []int64{0, 1, 2, 3})
	if err != nil {
		fmt.Printf("Error running field consistency check: %v\n", err)
		return
	}
	printConsistencyResult(fieldResult)
	fmt.Println()

	fmt.Println("Demo completed!")
}

// createDeltaTableWithVersions creates a Delta table with multiple versions with schema changes
func createDeltaTableWithVersions(tablePath string) error {
	// Create _delta_log directory
	logDir := filepath.Join(tablePath, "_delta_log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Create test versions with schema changes
	now := time.Now()

	// Version 0: Initial schema with 3 fields
	err = createTestVersion(logDir, 0, "CREATE_TABLE", now.Add(-3*time.Hour), []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		{"name": "value", "type": "float64"},
	}, []string{"part-00000.parquet"})
	if err != nil {
		return err
	}

	// Version 1: Added a field
	err = createTestVersion(logDir, 1, "ALTER_TABLE", now.Add(-2*time.Hour), []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		{"name": "value", "type": "float64"},
		{"name": "timestamp", "type": "timestamp"},
	}, []string{"part-00000.parquet", "part-00001.parquet"})
	if err != nil {
		return err
	}

	// Version 2: Changed a field type
	err = createTestVersion(logDir, 2, "ALTER_TABLE", now.Add(-1*time.Hour), []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		{"name": "value", "type": "int32"}, // Changed from float64 to int32
		{"name": "timestamp", "type": "timestamp"},
	}, []string{"part-00000.parquet", "part-00001.parquet", "part-00002.parquet"})
	if err != nil {
		return err
	}

	// Version 3: Removed a field and added a new one
	err = createTestVersion(logDir, 3, "ALTER_TABLE", now, []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		// "value" field removed
		{"name": "timestamp", "type": "timestamp"},
		{"name": "is_active", "type": "boolean"}, // Added new field
	}, []string{"part-00000.parquet", "part-00001.parquet", "part-00002.parquet", "part-00003.parquet"})
	if err != nil {
		return err
	}

	return nil
}

// createTestVersion creates a test version with a specific schema and files
func createTestVersion(logDir string, version int64, operation string, timestamp time.Time, fields []map[string]string, files []string) error {
	// Create metadata manager
	manager := datalake.NewMetadataManager(filepath.Dir(logDir))

	// Create schema
	schema := map[string]interface{}{
		"fields": fields,
	}

	// Create metadata
	metadata := map[string]interface{}{
		"version":    version,
		"timestamp":  timestamp.Unix(),
		"schema":     schema,
		"files":      files,
		"partitions": map[string][]string{},
		"metadata":   map[string]interface{}{},
	}

	// Write metadata to transaction log
	err := manager.WriteMetadataFile(version, metadata)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Create commit info
	commitInfo := map[string]interface{}{
		"timestamp":           timestamp.UnixMilli(),
		"operation":           operation,
		"operationParameters": map[string]string{},
		"isBlindAppend":       operation == "APPEND",
		"isolationLevel":      "Serializable",
	}

	// Write commit info
	err = manager.WriteCommitFile(version, commitInfo)
	if err != nil {
		return fmt.Errorf("failed to write commit file: %w", err)
	}

	return nil
}

// printConsistencyResult prints the consistency check result
func printConsistencyResult(result *datalake.ConsistencyCheckResult) {
	fmt.Printf("Consistency Check Result: %s\n", result.CheckType)
	fmt.Printf("Passed: %t\n", result.Passed)
	
	if len(result.Issues) > 0 {
		fmt.Printf("Issues Found (%d):\n", len(result.Issues))
		
		// Group issues by severity
		errorIssues := []datalake.ConsistencyIssue{}
		warningIssues := []datalake.ConsistencyIssue{}
		infoIssues := []datalake.ConsistencyIssue{}
		
		for _, issue := range result.Issues {
			switch issue.Severity {
			case "error":
				errorIssues = append(errorIssues, issue)
			case "warning":
				warningIssues = append(warningIssues, issue)
			case "info":
				infoIssues = append(infoIssues, issue)
			}
		}
		
		// Print errors first
		if len(errorIssues) > 0 {
			fmt.Printf("\nErrors (%d):\n", len(errorIssues))
			for _, issue := range errorIssues {
				fmt.Printf("  - %s: %s\n", issue.Field, issue.Description)
			}
		}
		
		// Then warnings
		if len(warningIssues) > 0 {
			fmt.Printf("\nWarnings (%d):\n", len(warningIssues))
			for _, issue := range warningIssues {
				fmt.Printf("  - %s: %s\n", issue.Field, issue.Description)
			}
		}
		
		// Then info
		if len(infoIssues) > 0 {
			fmt.Printf("\nInfo (%d):\n", len(infoIssues))
			for _, issue := range infoIssues {
				fmt.Printf("  - %s: %s\n", issue.Field, issue.Description)
			}
		}
	} else {
		fmt.Println("No issues found.")
	}
}
