package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
)

func init() {
	// Version history command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Delta table version management",
		Long:  "Commands for managing Delta table versions, including history, comparison, and rollback",
	}

	// List version history
	historyCmd := &cobra.Command{
		Use:   "history [table_path]",
		Short: "Show version history of a Delta table",
		Long:  "Display the commit history of a Delta table, including version, timestamp, and operation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid table path: %w", err)
			}

			manager := datalake.NewMetadataManager(tablePath)
			history, err := manager.GetVersionHistory()
			if err != nil {
				return fmt.Errorf("failed to get version history: %w", err)
			}

			// Print version history
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

			return nil
		},
	}

	// Compare versions
	compareCmd := &cobra.Command{
		Use:   "compare [table_path] [version1] [version2]",
		Short: "Compare two versions of a Delta table",
		Long:  "Compare schema and data changes between two versions of a Delta table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid table path: %w", err)
			}

			// Parse versions
			var version1, version2 int64
			_, err = fmt.Sscanf(args[1], "%d", &version1)
			if err != nil {
				return fmt.Errorf("invalid version1: %w", err)
			}
			_, err = fmt.Sscanf(args[2], "%d", &version2)
			if err != nil {
				return fmt.Errorf("invalid version2: %w", err)
			}

			manager := datalake.NewMetadataManager(tablePath)
			comparison, err := manager.CompareVersions(version1, version2)
			if err != nil {
				return fmt.Errorf("failed to compare versions: %w", err)
			}

			// Print comparison
			fmt.Printf("Comparing Version %d with Version %d:\n", version1, version2)
			fmt.Println("Schema Changes:")
			if len(comparison.SchemaChanges) == 0 {
				fmt.Println("  No schema changes")
			} else {
				for _, change := range comparison.SchemaChanges {
					fmt.Printf("  %s: %s\n", change.Type, change.Description)
				}
			}

			fmt.Println("Data Changes:")
			fmt.Printf("  Files Added: %d\n", comparison.FilesAdded)
			fmt.Printf("  Files Removed: %d\n", comparison.FilesRemoved)
			fmt.Printf("  Records Added: %d\n", comparison.RecordsAdded)
			fmt.Printf("  Records Removed: %d\n", comparison.RecordsRemoved)

			return nil
		},
	}

	// Rollback command
	rollbackCmd := &cobra.Command{
		Use:   "rollback [table_path] [version]",
		Short: "Rollback a Delta table to a previous version",
		Long:  "Rollback a Delta table to a specific version (up to 30 days old)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid table path: %w", err)
			}

			// Parse version
			var version int64
			_, err = fmt.Sscanf(args[1], "%d", &version)
			if err != nil {
				return fmt.Errorf("invalid version: %w", err)
			}

			// Get force flag
			force, _ := cmd.Flags().GetBool("force")

			manager := datalake.NewMetadataManager(tablePath)
			err = manager.RollbackToVersion(version, force)
			if err != nil {
				return fmt.Errorf("failed to rollback: %w", err)
			}

			fmt.Printf("Successfully rolled back to version %d\n", version)
			return nil
		},
	}

	// Add force flag to rollback command
	rollbackCmd.Flags().BoolP("force", "f", false, "Force rollback even if version is older than 30 days")

	versionCmd.AddCommand(historyCmd, compareCmd, rollbackCmd)
	rootCmd.AddCommand(versionCmd)
}
