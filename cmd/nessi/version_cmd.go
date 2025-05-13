package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage Delta Lake versions",
	Long:  `Manage Delta Lake versions, including transaction history, rollback, and version comparison.`,
}

// logTransactionCmd represents the log-transaction command
var logTransactionCmd = &cobra.Command{
	Use:   "log [table_path]",
	Short: "Log a transaction for a Delta Lake table",
	Long:  `Log a transaction for a Delta Lake table, recording file changes and metadata changes.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]

		// Get operation and commit message
		operation, _ := cmd.Flags().GetString("operation")
		message, _ := cmd.Flags().GetString("message")
		
		// Get added and removed files
		addedFiles, _ := cmd.Flags().GetStringSlice("added-files")
		removedFiles, _ := cmd.Flags().GetStringSlice("removed-files")
		
		// Get metadata changes
		schemaChange, _ := cmd.Flags().GetBool("schema-change")
		partitionChange, _ := cmd.Flags().GetBool("partition-change")
		propertiesChange, _ := cmd.Flags().GetBool("properties-change")
		addedColumns, _ := cmd.Flags().GetStringSlice("added-columns")
		removedColumns, _ := cmd.Flags().GetStringSlice("removed-columns")
		modifiedColumns, _ := cmd.Flags().GetStringSlice("modified-columns")
		
		// Get stats
		numRecords, _ := cmd.Flags().GetInt64("num-records")
		bytesAdded, _ := cmd.Flags().GetInt64("bytes-added")
		bytesRemoved, _ := cmd.Flags().GetInt64("bytes-removed")
		
		// Create commit info
		commitInfo := map[string]string{
			"message":   message,
			"timestamp": time.Now().Format(time.RFC3339),
			"user":      os.Getenv("USER"),
		}
		
		// Create metadata change
		metadataChange := &datalake.MetadataChange{
			SchemaChange:    schemaChange,
			PartitionChange: partitionChange,
			PropertiesChange: propertiesChange,
			AddedColumns:    addedColumns,
			RemovedColumns:  removedColumns,
			ModifiedColumns: modifiedColumns,
		}
		
		// Create stats
		stats := map[string]interface{}{}
		if numRecords > 0 {
			stats["numRecords"] = float64(numRecords)
		}
		if bytesAdded > 0 {
			stats["bytesAdded"] = float64(bytesAdded)
		}
		if bytesRemoved > 0 {
			stats["bytesRemoved"] = float64(bytesRemoved)
		}
		
		// Log transaction
		vm := datalake.NewVersionManager(tablePath)
		tx, err := vm.RecordTransaction(
			operation,
			commitInfo,
			addedFiles,
			removedFiles,
			metadataChange,
			stats,
		)
		if err != nil {
			fmt.Printf("Error logging transaction: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Transaction logged successfully (version %d)\n", tx.Version)
		fmt.Printf("Operation: %s\n", tx.Operation)
		fmt.Printf("Timestamp: %s\n", tx.Timestamp.Format(time.RFC3339))
		
		if len(tx.AddedFiles) > 0 {
			fmt.Printf("Added files: %s\n", strings.Join(tx.AddedFiles, ", "))
		}
		
		if len(tx.RemovedFiles) > 0 {
			fmt.Printf("Removed files: %s\n", strings.Join(tx.RemovedFiles, ", "))
		}
		
		if tx.MetadataChange != nil {
			if tx.MetadataChange.SchemaChange {
				fmt.Println("Schema changed: yes")
			}
			
			if len(tx.MetadataChange.AddedColumns) > 0 {
				fmt.Printf("Added columns: %s\n", strings.Join(tx.MetadataChange.AddedColumns, ", "))
			}
			
			if len(tx.MetadataChange.RemovedColumns) > 0 {
				fmt.Printf("Removed columns: %s\n", strings.Join(tx.MetadataChange.RemovedColumns, ", "))
			}
			
			if len(tx.MetadataChange.ModifiedColumns) > 0 {
				fmt.Printf("Modified columns: %s\n", strings.Join(tx.MetadataChange.ModifiedColumns, ", "))
			}
		}
		
		if tx.Stats != nil && len(tx.Stats) > 0 {
			fmt.Println("Stats:")
			for k, v := range tx.Stats {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	},
}

// historyCmd represents the history command
var versionHistoryCmd = &cobra.Command{
	Use:   "history [table_path]",
	Short: "Show transaction history",
	Long:  `Show the history of transactions for a Delta Lake table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		
		// Get version manager
		vm := datalake.NewVersionManager(tablePath)
		history, err := vm.GetTransactionHistory()
		if err != nil {
			fmt.Printf("Error getting transaction history: %v\n", err)
			os.Exit(1)
		}
		
		// Get format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(history, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling history to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		fmt.Printf("Transaction history for %s:\n", tablePath)
		fmt.Printf("Total versions: %d\n", len(history.Transactions))
		
		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 || limit > len(history.Transactions) {
			limit = len(history.Transactions)
		}
		
		// Show most recent transactions first
		fmt.Printf("Showing %d most recent transactions:\n\n", limit)
		for i := len(history.Transactions) - 1; i >= len(history.Transactions)-limit; i-- {
			tx := history.Transactions[i]
			
			fmt.Printf("Version %d: %s (%s)\n", tx.Version, tx.Operation, tx.Timestamp.Format(time.RFC3339))
			
			if tx.CommitInfo != nil {
				if message, ok := tx.CommitInfo["message"]; ok {
					fmt.Printf("  Message: %s\n", message)
				}
				
				if user, ok := tx.CommitInfo["user"]; ok {
					fmt.Printf("  User: %s\n", user)
				}
			}
			
			if len(tx.AddedFiles) > 0 {
				fmt.Printf("  Added files: %d\n", len(tx.AddedFiles))
			}
			
			if len(tx.RemovedFiles) > 0 {
				fmt.Printf("  Removed files: %d\n", len(tx.RemovedFiles))
			}
			
			if tx.MetadataChange != nil {
				if tx.MetadataChange.SchemaChange {
					fmt.Println("  Schema changed: yes")
					
					if len(tx.MetadataChange.AddedColumns) > 0 {
						fmt.Printf("  Added columns: %s\n", strings.Join(tx.MetadataChange.AddedColumns, ", "))
					}
					
					if len(tx.MetadataChange.RemovedColumns) > 0 {
						fmt.Printf("  Removed columns: %s\n", strings.Join(tx.MetadataChange.RemovedColumns, ", "))
					}
					
					if len(tx.MetadataChange.ModifiedColumns) > 0 {
						fmt.Printf("  Modified columns: %s\n", strings.Join(tx.MetadataChange.ModifiedColumns, ", "))
					}
				}
			}
			
			if tx.Stats != nil {
				if numRecords, ok := tx.Stats["numRecords"].(float64); ok {
					fmt.Printf("  Records: %.0f\n", numRecords)
				}
			}
			
			fmt.Println()
		}
	},
}

// showVersionCmd represents the show command
var showVersionCmd = &cobra.Command{
	Use:   "show [table_path] [version]",
	Short: "Show details of a specific version",
	Long:  `Show detailed information about a specific version of a Delta Lake table.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		versionStr := args[1]
		
		// Parse version
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			fmt.Printf("Error parsing version: %v\n", err)
			os.Exit(1)
		}
		
		// Get version manager
		vm := datalake.NewVersionManager(tablePath)
		tx, err := vm.GetTransaction(version)
		if err != nil {
			fmt.Printf("Error getting version: %v\n", err)
			os.Exit(1)
		}
		
		// Get format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(tx, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling transaction to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		fmt.Printf("Version %d: %s\n", tx.Version, tx.Operation)
		fmt.Printf("Timestamp: %s\n", tx.Timestamp.Format(time.RFC3339))
		fmt.Printf("ID: %s\n", tx.ID)
		
		if tx.CommitInfo != nil {
			fmt.Println("Commit info:")
			for k, v := range tx.CommitInfo {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
		
		if len(tx.AddedFiles) > 0 {
			fmt.Println("Added files:")
			for _, file := range tx.AddedFiles {
				fmt.Printf("  %s\n", file)
			}
		}
		
		if len(tx.RemovedFiles) > 0 {
			fmt.Println("Removed files:")
			for _, file := range tx.RemovedFiles {
				fmt.Printf("  %s\n", file)
			}
		}
		
		if tx.MetadataChange != nil {
			fmt.Println("Metadata changes:")
			
			if tx.MetadataChange.SchemaChange {
				fmt.Println("  Schema changed: yes")
			}
			
			if tx.MetadataChange.PartitionChange {
				fmt.Println("  Partition changed: yes")
			}
			
			if tx.MetadataChange.PropertiesChange {
				fmt.Println("  Properties changed: yes")
			}
			
			if len(tx.MetadataChange.AddedColumns) > 0 {
				fmt.Println("  Added columns:")
				for _, col := range tx.MetadataChange.AddedColumns {
					fmt.Printf("    %s\n", col)
				}
			}
			
			if len(tx.MetadataChange.RemovedColumns) > 0 {
				fmt.Println("  Removed columns:")
				for _, col := range tx.MetadataChange.RemovedColumns {
					fmt.Printf("    %s\n", col)
				}
			}
			
			if len(tx.MetadataChange.ModifiedColumns) > 0 {
				fmt.Println("  Modified columns:")
				for _, col := range tx.MetadataChange.ModifiedColumns {
					fmt.Printf("    %s\n", col)
				}
			}
		}
		
		if tx.Stats != nil && len(tx.Stats) > 0 {
			fmt.Println("Stats:")
			for k, v := range tx.Stats {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
		
		fmt.Printf("Isolation level: %s\n", tx.IsolationLevel)
		fmt.Printf("Read version: %d\n", tx.ReadVersion)
		
		if tx.UserID != "" {
			fmt.Printf("User ID: %s\n", tx.UserID)
		}
		
		if tx.ClientInfo != nil && len(tx.ClientInfo) > 0 {
			fmt.Println("Client info:")
			for k, v := range tx.ClientInfo {
				fmt.Printf("  %s: %s\n", k, v)
			}
		}
	},
}

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback [table_path] [version]",
	Short: "Rollback to a specific version",
	Long:  `Rollback a Delta Lake table to a specific version.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		versionStr := args[1]
		
		// Parse version
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			fmt.Printf("Error parsing version: %v\n", err)
			os.Exit(1)
		}
		
		// Get version manager
		vm := datalake.NewVersionManager(tablePath)
		
		// Get current version
		latestTx, err := vm.GetLatestTransaction()
		if err != nil {
			fmt.Printf("Error getting latest version: %v\n", err)
			os.Exit(1)
		}
		
		// Check if already at target version
		if latestTx.Version == version {
			fmt.Printf("Already at version %d\n", version)
			return
		}
		
		// Get target version
		targetTx, err := vm.GetTransaction(version)
		if err != nil {
			fmt.Printf("Error getting target version: %v\n", err)
			os.Exit(1)
		}
		
		// Confirm rollback
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Are you sure you want to rollback from version %d to version %d? [y/N] ", latestTx.Version, version)
			var response string
			fmt.Scanln(&response)
			
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("Rollback cancelled")
				return
			}
		}
		
		// Rollback
		err = vm.RollbackToVersion(version)
		if err != nil {
			fmt.Printf("Error rolling back: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Successfully rolled back from version %d to version %d\n", latestTx.Version, version)
		fmt.Printf("Timestamp: %s\n", targetTx.Timestamp.Format(time.RFC3339))
		
		// Get new latest transaction (the rollback transaction)
		newLatestTx, err := vm.GetLatestTransaction()
		if err != nil {
			fmt.Printf("Error getting new latest version: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Rollback transaction recorded as version %d\n", newLatestTx.Version)
	},
}

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare [table_path] [from_version] [to_version]",
	Short: "Compare two versions",
	Long:  `Compare two versions of a Delta Lake table and show the differences.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		fromVersionStr := args[1]
		toVersionStr := args[2]
		
		// Parse versions
		fromVersion, err := strconv.Atoi(fromVersionStr)
		if err != nil {
			fmt.Printf("Error parsing from_version: %v\n", err)
			os.Exit(1)
		}
		
		toVersion, err := strconv.Atoi(toVersionStr)
		if err != nil {
			fmt.Printf("Error parsing to_version: %v\n", err)
			os.Exit(1)
		}
		
		// Get version manager
		vm := datalake.NewVersionManager(tablePath)
		
		// Compare versions
		diff, err := vm.CompareVersions(fromVersion, toVersion)
		if err != nil {
			fmt.Printf("Error comparing versions: %v\n", err)
			os.Exit(1)
		}
		
		// Get format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(diff, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling diff to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		fmt.Printf("Comparing version %d to version %d:\n", fromVersion, toVersion)
		fmt.Printf("Time span: %s\n", formatDuration(diff.TimeSpan))
		fmt.Printf("From: %s\n", diff.FromTimestamp.Format(time.RFC3339))
		fmt.Printf("To: %s\n", diff.ToTimestamp.Format(time.RFC3339))
		
		fmt.Println("\nChanges:")
		fmt.Printf("  Added rows: %d\n", diff.AddedRows)
		fmt.Printf("  Removed rows: %d\n", diff.RemovedRows)
		fmt.Printf("  Modified rows: %d\n", diff.ModifiedRows)
		
		if len(diff.SchemaChanges) > 0 {
			fmt.Println("\nSchema changes:")
			for _, change := range diff.SchemaChanges {
				fmt.Printf("  - %s\n", change.Description)
			}
		}
		
		if len(diff.Operations) > 0 {
			fmt.Println("\nOperations:")
			for i, op := range diff.Operations {
				fmt.Printf("  %d. %s\n", i+1, op)
			}
		}
	},
}

// timeVersionCmd represents the time-version command
var timeVersionCmd = &cobra.Command{
	Use:   "time-version [table_path] [timestamp]",
	Short: "Get version at a specific timestamp",
	Long:  `Get the version of a Delta Lake table that was current at a specific timestamp.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		timestampStr := args[1]
		
		// Parse timestamp
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			fmt.Printf("Error parsing timestamp: %v\n", err)
			fmt.Println("Timestamp must be in RFC3339 format (e.g., 2023-01-01T12:00:00Z)")
			os.Exit(1)
		}
		
		// Get version manager
		vm := datalake.NewVersionManager(tablePath)
		
		// Get version at timestamp
		tx, err := vm.GetVersionAtTimestamp(timestamp)
		if err != nil {
			fmt.Printf("Error getting version at timestamp: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Version at %s: %d\n", timestamp.Format(time.RFC3339), tx.Version)
		fmt.Printf("Actual timestamp: %s\n", tx.Timestamp.Format(time.RFC3339))
		fmt.Printf("Operation: %s\n", tx.Operation)
		
		if tx.CommitInfo != nil {
			if message, ok := tx.CommitInfo["message"]; ok {
				fmt.Printf("Message: %s\n", message)
			}
		}
	},
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d.Hours() > 24 {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%d days, %d hours", days, hours)
	}
	
	if d.Hours() >= 1 {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%d hours, %d minutes", hours, minutes)
	}
	
	if d.Minutes() >= 1 {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%d minutes, %d seconds", minutes, seconds)
	}
	
	return fmt.Sprintf("%.2f seconds", d.Seconds())
}

func init() {
	rootCmd.AddCommand(versionCmd)
	
	// Add subcommands
	versionCmd.AddCommand(logTransactionCmd)
	versionCmd.AddCommand(versionHistoryCmd)
	versionCmd.AddCommand(showVersionCmd)
	versionCmd.AddCommand(rollbackCmd)
	versionCmd.AddCommand(compareCmd)
	versionCmd.AddCommand(timeVersionCmd)
	
	// Add flags for log-transaction
	logTransactionCmd.Flags().String("operation", "WRITE", "Operation type (WRITE, UPDATE, DELETE, MERGE)")
	logTransactionCmd.Flags().String("message", "", "Commit message")
	logTransactionCmd.MarkFlagRequired("message")
	
	logTransactionCmd.Flags().StringSlice("added-files", []string{}, "Files added in this transaction")
	logTransactionCmd.Flags().StringSlice("removed-files", []string{}, "Files removed in this transaction")
	
	logTransactionCmd.Flags().Bool("schema-change", false, "Whether the transaction includes schema changes")
	logTransactionCmd.Flags().Bool("partition-change", false, "Whether the transaction includes partition changes")
	logTransactionCmd.Flags().Bool("properties-change", false, "Whether the transaction includes properties changes")
	logTransactionCmd.Flags().StringSlice("added-columns", []string{}, "Columns added in this transaction")
	logTransactionCmd.Flags().StringSlice("removed-columns", []string{}, "Columns removed in this transaction")
	logTransactionCmd.Flags().StringSlice("modified-columns", []string{}, "Columns modified in this transaction")
	
	logTransactionCmd.Flags().Int64("num-records", 0, "Number of records affected")
	logTransactionCmd.Flags().Int64("bytes-added", 0, "Bytes added in this transaction")
	logTransactionCmd.Flags().Int64("bytes-removed", 0, "Bytes removed in this transaction")
	
	// Add flags for history
	versionHistoryCmd.Flags().String("format", "text", "Output format (text, json)")
	versionHistoryCmd.Flags().Int("limit", 10, "Maximum number of transactions to show")
	
	// Add flags for show
	showVersionCmd.Flags().String("format", "text", "Output format (text, json)")
	
	// Add flags for rollback
	rollbackCmd.Flags().Bool("force", false, "Skip confirmation prompt")
	
	// Add flags for compare
	compareCmd.Flags().String("format", "text", "Output format (text, json)")
}
