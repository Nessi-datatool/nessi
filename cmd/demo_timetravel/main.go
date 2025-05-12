package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

func main() {
	fmt.Println("Time Travel Demo")
	fmt.Println("===============")
	fmt.Println()

	// Get test table path
	tablePath := "test_data/version_history_table"
	absPath, err := filepath.Abs(tablePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create metadata manager
	manager := datalake.NewMetadataManager(absPath)

	// Get version history
	fmt.Println("Getting version history...")
	history, err := manager.GetVersionHistory()
	if err != nil {
		fmt.Printf("Error getting version history: %v\n", err)
		os.Exit(1)
	}

	// Display version history
	fmt.Println("Version History:")
	fmt.Println("----------------")
	for _, entry := range history {
		timestamp := time.Unix(0, entry.Timestamp*int64(time.Millisecond))
		fmt.Printf("Version: %d\n", entry.Version)
		fmt.Printf("Timestamp: %s\n", timestamp.Format(time.RFC3339))
		fmt.Printf("Operation: %s\n", entry.Operation)
		fmt.Println("----------------")
	}
	fmt.Println()

	// Demonstrate time travel by version
	if len(history) > 0 {
		// Time travel to first version
		firstVersion := history[len(history)-1].Version
		fmt.Printf("Time traveling to version %d...\n", firstVersion)
		
		version := firstVersion
		result, err := manager.TimeTravel(datalake.TimeTravelOptions{Version: &version})
		if err != nil {
			fmt.Printf("Error time traveling to version %d: %v\n", version, err)
		} else {
			printTimeTravelResult(result)
		}
		fmt.Println()

		// Time travel to latest version
		latestVersion := history[0].Version
		fmt.Printf("Time traveling to version %d...\n", latestVersion)
		
		version = latestVersion
		result, err = manager.TimeTravel(datalake.TimeTravelOptions{Version: &version})
		if err != nil {
			fmt.Printf("Error time traveling to version %d: %v\n", version, err)
		} else {
			printTimeTravelResult(result)
		}
		fmt.Println()
	}

	// Demonstrate time travel by timestamp
	if len(history) > 1 {
		// Get timestamps of first and latest versions
		firstTimestamp := time.Unix(0, history[len(history)-1].Timestamp*int64(time.Millisecond))
		latestTimestamp := time.Unix(0, history[0].Timestamp*int64(time.Millisecond))
		
		// Calculate a timestamp between first and latest versions
		midTimestamp := firstTimestamp.Add(latestTimestamp.Sub(firstTimestamp) / 2)
		
		fmt.Printf("Time traveling to timestamp %s...\n", midTimestamp.Format(time.RFC3339))
		
		result, err := manager.TimeTravel(datalake.TimeTravelOptions{Timestamp: &midTimestamp})
		if err != nil {
			fmt.Printf("Error time traveling to timestamp %s: %v\n", midTimestamp.Format(time.RFC3339), err)
		} else {
			printTimeTravelResult(result)
		}
		fmt.Println()
	}

	// Demonstrate getting schema at version
	if len(history) > 0 {
		version := history[0].Version
		fmt.Printf("Getting schema at version %d...\n", version)
		
		schemaFields, err := manager.GetSchemaFieldsAtVersion(version)
		if err != nil {
			fmt.Printf("Error getting schema at version %d: %v\n", version, err)
		} else {
			fmt.Printf("Schema at version %d (%d fields):\n", version, len(schemaFields))
			for _, field := range schemaFields {
				fmt.Printf("  %s: %s\n", field.Name, field.Type)
			}
		}
		fmt.Println()
	}

	// Demonstrate getting files at version
	if len(history) > 0 {
		version := history[0].Version
		fmt.Printf("Getting files at version %d...\n", version)
		
		files, err := manager.GetFilesAtVersion(version)
		if err != nil {
			fmt.Printf("Error getting files at version %d: %v\n", version, err)
		} else {
			fmt.Printf("Files at version %d (%d files):\n", version, len(files))
			for _, file := range files {
				fmt.Printf("  %s\n", file)
			}
		}
		fmt.Println()
	}

	fmt.Println("Demo completed!")
}

// printTimeTravelResult prints the time travel result
func printTimeTravelResult(result *datalake.TimeTravelResult) {
	fmt.Printf("Time Travel Result:\n")
	fmt.Printf("Version: %d\n", result.Version)
	fmt.Printf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
	
	fmt.Printf("Schema (%d fields):\n", len(result.SchemaFields))
	for _, field := range result.SchemaFields {
		fmt.Printf("  %s: %s\n", field.Name, field.Type)
	}
	
	fmt.Printf("Files (%d):\n", len(result.Files))
	for _, file := range result.Files {
		fmt.Printf("  %s\n", file)
	}
}
