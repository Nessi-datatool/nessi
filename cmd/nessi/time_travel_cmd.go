package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

// timeTravelCmd represents the time-travel command
var timeTravelCmd = &cobra.Command{
	Use:   "time-travel",
	Short: "Time travel operations for Delta Lake tables",
	Long:  `Perform time travel operations on Delta Lake tables, including querying data at specific versions or timestamps.`,
}

// queryVersionCmd represents the query-version command
var queryVersionCmd = &cobra.Command{
	Use:   "query-version [table_path] [version]",
	Short: "Query data at a specific version",
	Long:  `Query data at a specific version of a Delta Lake table.`,
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
		
		// Create time travel manager
		tt := datalake.NewTimeTravel(tablePath)
		
		// Query at version
		result, err := tt.QueryAtVersion(version)
		if err != nil {
			fmt.Printf("Error querying version: %v\n", err)
			os.Exit(1)
		}
		
		// Get format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling result to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		fmt.Printf("Query result for version %d:\n", version)
		fmt.Printf("Timestamp: %s\n", result.Transaction.Timestamp.Format(time.RFC3339))
		fmt.Printf("Operation: %s\n", result.Transaction.Operation)
		
		if result.Transaction.CommitInfo != nil {
			if message, ok := result.Transaction.CommitInfo["message"]; ok {
				fmt.Printf("Message: %s\n", message)
			}
		}
		
		fmt.Printf("Files: %d\n", len(result.Files))
		if showFiles, _ := cmd.Flags().GetBool("show-files"); showFiles && len(result.Files) > 0 {
			fmt.Println("Files:")
			for _, file := range result.Files {
				fmt.Printf("  %s\n", file)
			}
		}
		
		fmt.Printf("Schema fields: %d\n", len(result.Schema.Fields))
		if showSchema, _ := cmd.Flags().GetBool("show-schema"); showSchema && len(result.Schema.Fields) > 0 {
			fmt.Println("Schema:")
			for _, field := range result.Schema.Fields {
				fmt.Printf("  %s: %s", field.Name, field.Type)
				if !field.Nullable {
					fmt.Print(" (not null)")
				}
				fmt.Println()
			}
		}
	},
}

// queryTimestampCmd represents the query-timestamp command
var queryTimestampCmd = &cobra.Command{
	Use:   "query-timestamp [table_path] [timestamp]",
	Short: "Query data at a specific timestamp",
	Long:  `Query data at a specific timestamp of a Delta Lake table. Timestamp must be in RFC3339 format (e.g., 2023-01-01T12:00:00Z).`,
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
		
		// Create time travel manager
		tt := datalake.NewTimeTravel(tablePath)
		
		// Query at timestamp
		result, err := tt.QueryAtTimestamp(timestamp)
		if err != nil {
			fmt.Printf("Error querying timestamp: %v\n", err)
			os.Exit(1)
		}
		
		// Get format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling result to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		fmt.Printf("Query result for timestamp %s:\n", timestampStr)
		fmt.Printf("Actual version: %d\n", result.Transaction.Version)
		fmt.Printf("Actual timestamp: %s\n", result.Transaction.Timestamp.Format(time.RFC3339))
		fmt.Printf("Operation: %s\n", result.Transaction.Operation)
		
		if result.Transaction.CommitInfo != nil {
			if message, ok := result.Transaction.CommitInfo["message"]; ok {
				fmt.Printf("Message: %s\n", message)
			}
		}
		
		fmt.Printf("Files: %d\n", len(result.Files))
		if showFiles, _ := cmd.Flags().GetBool("show-files"); showFiles && len(result.Files) > 0 {
			fmt.Println("Files:")
			for _, file := range result.Files {
				fmt.Printf("  %s\n", file)
			}
		}
		
		fmt.Printf("Schema fields: %d\n", len(result.Schema.Fields))
		if showSchema, _ := cmd.Flags().GetBool("show-schema"); showSchema && len(result.Schema.Fields) > 0 {
			fmt.Println("Schema:")
			for _, field := range result.Schema.Fields {
				fmt.Printf("  %s: %s", field.Name, field.Type)
				if !field.Nullable {
					fmt.Print(" (not null)")
				}
				fmt.Println()
			}
		}
	},
}

// versionsInRangeCmd represents the versions-in-range command
var versionsInRangeCmd = &cobra.Command{
	Use:   "versions-in-range [table_path] [start_timestamp] [end_timestamp]",
	Short: "List versions in a time range",
	Long:  `List versions of a Delta Lake table that fall within a specific time range. Timestamps must be in RFC3339 format (e.g., 2023-01-01T12:00:00Z).`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		startTimestampStr := args[1]
		endTimestampStr := args[2]
		
		// Parse timestamps
		startTimestamp, err := time.Parse(time.RFC3339, startTimestampStr)
		if err != nil {
			fmt.Printf("Error parsing start timestamp: %v\n", err)
			fmt.Println("Timestamp must be in RFC3339 format (e.g., 2023-01-01T12:00:00Z)")
			os.Exit(1)
		}
		
		endTimestamp, err := time.Parse(time.RFC3339, endTimestampStr)
		if err != nil {
			fmt.Printf("Error parsing end timestamp: %v\n", err)
			fmt.Println("Timestamp must be in RFC3339 format (e.g., 2023-01-01T12:00:00Z)")
			os.Exit(1)
		}
		
		// Create time travel manager
		tt := datalake.NewTimeTravel(tablePath)
		
		// Get versions in range
		versions, err := tt.GetVersionsInTimeRange(startTimestamp, endTimestamp)
		if err != nil {
			fmt.Printf("Error getting versions in range: %v\n", err)
			os.Exit(1)
		}
		
		// Get format
		format, _ := cmd.Flags().GetString("format")
		
		if format == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(versions, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling versions to JSON: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(string(jsonData))
			return
		}
		
		// Output as text
		fmt.Printf("Versions between %s and %s:\n", startTimestampStr, endTimestampStr)
		fmt.Printf("Found %d versions\n\n", len(versions))
		
		for i, tx := range versions {
			fmt.Printf("%d. Version %d: %s (%s)\n", i+1, tx.Version, tx.Operation, tx.Timestamp.Format(time.RFC3339))
			
			if tx.CommitInfo != nil {
				if message, ok := tx.CommitInfo["message"]; ok {
					fmt.Printf("   Message: %s\n", message)
				}
			}
			
			fmt.Println()
		}
	},
}

// exportSnapshotCmd represents the export-snapshot command
var exportSnapshotCmd = &cobra.Command{
	Use:   "export-snapshot [table_path] [version] [output_dir]",
	Short: "Export a snapshot of a specific version",
	Long:  `Export a snapshot of a Delta Lake table at a specific version to a directory.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		versionStr := args[1]
		outputDir := args[2]
		
		// Parse version
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			fmt.Printf("Error parsing version: %v\n", err)
			os.Exit(1)
		}
		
		// Create time travel manager
		tt := datalake.NewTimeTravel(tablePath)
		
		// Export snapshot
		err = tt.ExportVersionSnapshot(version, outputDir)
		if err != nil {
			fmt.Printf("Error exporting snapshot: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Successfully exported snapshot of version %d to %s\n", version, outputDir)
		fmt.Println("The snapshot contains:")
		fmt.Println("- metadata.json: Table metadata")
		fmt.Println("- schema.json: Table schema")
		fmt.Println("- files.json: List of data files")
		fmt.Println("- transaction.json: Transaction details")
		fmt.Println("- data/: Directory that would contain data files")
	},
}

// reconstructStateCmd represents the reconstruct-state command
var reconstructStateCmd = &cobra.Command{
	Use:   "reconstruct-state [table_path] [version_or_timestamp] [output_dir]",
	Short: "Reconstruct table state at a specific version or timestamp",
	Long:  `Reconstruct the state of a Delta Lake table at a specific version or timestamp.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		versionOrTimestamp := args[1]
		outputDir := args[2]
		
		// Create time travel manager
		tt := datalake.NewTimeTravel(tablePath)
		
		// Check if version or timestamp
		isTimestamp, _ := cmd.Flags().GetBool("timestamp")
		
		var err error
		if isTimestamp {
			// Parse timestamp
			timestamp, err := time.Parse(time.RFC3339, versionOrTimestamp)
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				fmt.Println("Timestamp must be in RFC3339 format (e.g., 2023-01-01T12:00:00Z)")
				os.Exit(1)
			}
			
			// Reconstruct state at timestamp
			err = tt.ReconstructStateAtTimestamp(timestamp, outputDir)
			if err != nil {
				fmt.Printf("Error reconstructing state: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Successfully reconstructed state at timestamp %s to %s\n", versionOrTimestamp, outputDir)
		} else {
			// Parse version
			version, err := strconv.Atoi(versionOrTimestamp)
			if err != nil {
				fmt.Printf("Error parsing version: %v\n", err)
				os.Exit(1)
			}
			
			// Reconstruct state at version
			err = tt.ReconstructStateAtVersion(version, outputDir)
			if err != nil {
				fmt.Printf("Error reconstructing state: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Successfully reconstructed state at version %s to %s\n", versionOrTimestamp, outputDir)
		}
		
		fmt.Println("The reconstructed state contains:")
		fmt.Println("- metadata.json: Table metadata")
		fmt.Println("- schema.json: Table schema")
		fmt.Println("- files.json: List of data files")
		fmt.Println("- transaction.json: Transaction details")
		fmt.Println("- data/: Directory that would contain data files")
	},
}

func init() {
	rootCmd.AddCommand(timeTravelCmd)
	
	// Add subcommands
	timeTravelCmd.AddCommand(queryVersionCmd)
	timeTravelCmd.AddCommand(queryTimestampCmd)
	timeTravelCmd.AddCommand(versionsInRangeCmd)
	timeTravelCmd.AddCommand(exportSnapshotCmd)
	timeTravelCmd.AddCommand(reconstructStateCmd)
	
	// Add flags for query-version
	queryVersionCmd.Flags().String("format", "text", "Output format (text, json)")
	queryVersionCmd.Flags().Bool("show-files", false, "Show list of files")
	queryVersionCmd.Flags().Bool("show-schema", false, "Show schema details")
	
	// Add flags for query-timestamp
	queryTimestampCmd.Flags().String("format", "text", "Output format (text, json)")
	queryTimestampCmd.Flags().Bool("show-files", false, "Show list of files")
	queryTimestampCmd.Flags().Bool("show-schema", false, "Show schema details")
	
	// Add flags for versions-in-range
	versionsInRangeCmd.Flags().String("format", "text", "Output format (text, json)")
	
	// Add flags for reconstruct-state
	reconstructStateCmd.Flags().Bool("timestamp", false, "Interpret the second argument as a timestamp instead of a version")
}
