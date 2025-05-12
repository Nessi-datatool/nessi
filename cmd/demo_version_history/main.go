package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

func main() {
	fmt.Println("Version History & Rollback Demo")
	fmt.Println("==============================")
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
		if entry.UserName != "" {
			fmt.Printf("User: %s\n", entry.UserName)
		}
		fmt.Println("----------------")
	}

	// If we have at least two versions, compare them
	if len(history) >= 2 {
		version1 := history[len(history)-1].Version // Oldest version
		version2 := history[0].Version              // Newest version

		fmt.Printf("\nComparing Version %d with Version %d...\n", version1, version2)
		comparison, err := manager.CompareVersions(version1, version2)
		if err != nil {
			fmt.Printf("Error comparing versions: %v\n", err)
		} else {
			fmt.Println("Schema Changes:")
			if len(comparison.SchemaChanges) == 0 {
				fmt.Println("  No schema changes")
			} else {
				for _, change := range comparison.SchemaChanges {
					var typeDesc string
					switch change.Type {
					case "added":
						typeDesc = fmt.Sprintf("Added field '%s' with type %v", change.FieldName, change.NewType)
					case "removed":
						typeDesc = fmt.Sprintf("Removed field '%s' with type %v", change.FieldName, change.OldType)
					case "type_changed":
						typeDesc = fmt.Sprintf("Changed type of field '%s' from %v to %v", change.FieldName, change.OldType, change.NewType)
					default:
						typeDesc = fmt.Sprintf("Unknown change to field '%s'", change.FieldName)
					}
					fmt.Printf("  %s\n", typeDesc)
				}
			}

			fmt.Println("Data Changes:")
			fmt.Printf("  Files Added: %d\n", comparison.FilesAdded)
			fmt.Printf("  Files Removed: %d\n", comparison.FilesRemoved)
			fmt.Printf("  Records Added (estimate): %d\n", comparison.RecordsAdded)
			fmt.Printf("  Records Removed (estimate): %d\n", comparison.RecordsRemoved)
		}
	}

	// Demonstrate rollback (only if we have versions)
	if len(history) > 0 {
		oldestVersion := history[len(history)-1].Version
		fmt.Printf("\nDemonstrating rollback to version %d...\n", oldestVersion)
		
		// Rollback to oldest version (force=true to bypass 30-day restriction for demo)
		err := manager.RollbackToVersion(oldestVersion, true)
		if err != nil {
			fmt.Printf("Error rolling back: %v\n", err)
		} else {
			fmt.Printf("Successfully rolled back to version %d\n", oldestVersion)
			
			// Get updated version history
			updatedHistory, err := manager.GetVersionHistory()
			if err != nil {
				fmt.Printf("Error getting updated version history: %v\n", err)
			} else {
				fmt.Println("\nUpdated Version History:")
				fmt.Println("------------------------")
				// Just show the latest version (the rollback)
				if len(updatedHistory) > 0 {
					entry := updatedHistory[0]
					timestamp := time.Unix(0, entry.Timestamp*int64(time.Millisecond))
					fmt.Printf("Version: %d\n", entry.Version)
					fmt.Printf("Timestamp: %s\n", timestamp.Format(time.RFC3339))
					fmt.Printf("Operation: %s\n", entry.Operation)
					if entry.Parameters != nil && entry.Parameters["targetVersion"] != "" {
						fmt.Printf("Target Version: %s\n", entry.Parameters["targetVersion"])
					}
				}
			}
		}
	}

	fmt.Println("\nDemo completed!")
}
