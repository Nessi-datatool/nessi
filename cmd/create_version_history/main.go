package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Create a test Delta table with multiple versions for demonstrating version history
func main() {
	fmt.Println("Creating test Delta table with version history...")

	// Create test directory
	testDir := "test_data/version_history_table"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("Error creating test directory: %v\n", err)
		os.Exit(1)
	}

	// Create _delta_log directory
	logDir := filepath.Join(testDir, "_delta_log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Error creating _delta_log directory: %v\n", err)
		os.Exit(1)
	}

	// Create version 0 - Initial schema with 3 fields
	createVersion0(logDir)
	fmt.Println("Created version 0: Initial table creation")

	// Create version 1 - Add a new field
	createVersion1(logDir)
	fmt.Println("Created version 1: Added 'timestamp' field")

	// Create version 2 - Add another field
	createVersion2(logDir)
	fmt.Println("Created version 2: Added 'is_active' field")

	fmt.Println("Test Delta table created successfully!")
	fmt.Printf("Table path: %s\n", testDir)
}

func createVersion0(logDir string) {
	// Create schema with 3 fields
	schema := map[string]interface{}{
		"fields": []map[string]interface{}{
			{"name": "id", "type": "int32"},
			{"name": "name", "type": "utf8"},
			{"name": "value", "type": "float64"},
		},
	}

	// Create metadata
	metadata := map[string]interface{}{
		"version":    0,
		"timestamp":  time.Now().Add(-2 * time.Hour).Unix(),
		"schema":     schema,
		"files":      []string{"part-00000.parquet"},
		"partitions": map[string][]string{},
		"metadata":   map[string]interface{}{},
	}

	// Write transaction log
	writeTransactionLog(logDir, 0, metadata, "CREATE_TABLE", nil)
}

func createVersion1(logDir string) {
	// Create schema with 4 fields (added timestamp)
	schema := map[string]interface{}{
		"fields": []map[string]interface{}{
			{"name": "id", "type": "int32"},
			{"name": "name", "type": "utf8"},
			{"name": "value", "type": "float64"},
			{"name": "timestamp", "type": "timestamp"},
		},
	}

	// Create metadata
	metadata := map[string]interface{}{
		"version":    1,
		"timestamp":  time.Now().Add(-1 * time.Hour).Unix(),
		"schema":     schema,
		"files":      []string{"part-00000.parquet", "part-00001.parquet"},
		"partitions": map[string][]string{},
		"metadata":   map[string]interface{}{},
	}

	// Write transaction log
	params := map[string]string{
		"addedFields": "timestamp",
	}
	writeTransactionLog(logDir, 1, metadata, "ALTER_TABLE", params)
}

func createVersion2(logDir string) {
	// Create schema with 5 fields (added is_active)
	schema := map[string]interface{}{
		"fields": []map[string]interface{}{
			{"name": "id", "type": "int32"},
			{"name": "name", "type": "utf8"},
			{"name": "value", "type": "float64"},
			{"name": "timestamp", "type": "timestamp"},
			{"name": "is_active", "type": "boolean"},
		},
	}

	// Create metadata
	metadata := map[string]interface{}{
		"version":    2,
		"timestamp":  time.Now().Unix(),
		"schema":     schema,
		"files":      []string{"part-00000.parquet", "part-00001.parquet", "part-00002.parquet"},
		"partitions": map[string][]string{},
		"metadata":   map[string]interface{}{},
	}

	// Write transaction log
	params := map[string]string{
		"addedFields": "is_active",
	}
	writeTransactionLog(logDir, 2, metadata, "ALTER_TABLE", params)
}

func writeTransactionLog(logDir string, version int64, metadata map[string]interface{}, operation string, params map[string]string) {
	// Convert metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling metadata: %v\n", err)
		os.Exit(1)
	}

	// Write transaction log
	logFile := filepath.Join(logDir, fmt.Sprintf("%020d.json", version))
	if err := os.WriteFile(logFile, data, 0644); err != nil {
		fmt.Printf("Error writing transaction log: %v\n", err)
		os.Exit(1)
	}

	// Create commit info
	commitInfo := map[string]interface{}{
		"timestamp":           time.Now().UnixMilli(),
		"operation":           operation,
		"operationParameters": params,
		"userName":            "demo_user",
		"userId":              "demo_user_id",
		"isBlindAppend":       false,
		"isolationLevel":      "Serializable",
	}

	// Convert commit info to JSON
	commitData, err := json.MarshalIndent(commitInfo, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling commit info: %v\n", err)
		os.Exit(1)
	}

	// Write commit info
	commitFile := filepath.Join(logDir, fmt.Sprintf("%020d.commit.json", version))
	if err := os.WriteFile(commitFile, commitData, 0644); err != nil {
		fmt.Printf("Error writing commit info: %v\n", err)
		os.Exit(1)
	}
}
