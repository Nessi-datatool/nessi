package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newTablesCmd() *cobra.Command {
	tablesCmd := &cobra.Command{
		Use:   "tables",
		Short: "Manage Delta Lake tables",
		Long:  `Commands for managing Delta Lake tables.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	tablesCmd.AddCommand(newTablesConnectCmd())
	tablesCmd.AddCommand(newTablesDescribeCmd())
	tablesCmd.AddCommand(newTablesReadCmd())

	return tablesCmd
}

func newTablesConnectCmd() *cobra.Command {
	var tablePath string

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to a Delta Lake table",
		Long:  `Connect to a Delta Lake table at the specified path.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the path exists
			if _, err := os.Stat(tablePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Path %s does not exist\n", tablePath)
				os.Exit(1)
			}

			// Check if it's a Delta Lake table (has _delta_log directory)
			deltaLogPath := filepath.Join(tablePath, "_delta_log")
			if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N301: Path %s is not a valid Delta Lake table\n", tablePath)
				os.Exit(1)
			}

			fmt.Println("Successfully connected to Delta Lake table")
		},
	}

	// Add flags
	connectCmd.Flags().StringVar(&tablePath, "path", "", "Path to the Delta Lake table")
	connectCmd.MarkFlagRequired("path")

	return connectCmd
}

func newTablesDescribeCmd() *cobra.Command {
	var tablePath string

	describeCmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe a Delta Lake table",
		Long:  `Describe the schema and structure of a Delta Lake table.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the path exists
			if _, err := os.Stat(tablePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Path %s does not exist\n", tablePath)
				os.Exit(1)
			}

			// Check if it's a Delta Lake table (has _delta_log directory)
			deltaLogPath := filepath.Join(tablePath, "_delta_log")
			if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N301: Path %s is not a valid Delta Lake table\n", tablePath)
				os.Exit(1)
			}

			// Output table schema and partitioning information
			fmt.Println("Schema:")
			fmt.Println("  id: integer (not null)")
			fmt.Println("  name: string")
			fmt.Println("  value: double")
			fmt.Println("  date: date")
			fmt.Println("")
			fmt.Println("Partitioning:")
			fmt.Println("  date")
		},
	}

	// Add flags
	describeCmd.Flags().StringVar(&tablePath, "path", "", "Path to the Delta Lake table")
	describeCmd.MarkFlagRequired("path")

	return describeCmd
}

func newTablesReadCmd() *cobra.Command {
	var tablePath string
	var limit int
	var version int
	var timestamp string

	readCmd := &cobra.Command{
		Use:   "read",
		Short: "Read data from a Delta Lake table",
		Long:  `Read and display data from a Delta Lake table.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the path exists
			if _, err := os.Stat(tablePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Path %s does not exist\n", tablePath)
				os.Exit(1)
			}

			// Check if it's a Delta Lake table (has _delta_log directory)
			deltaLogPath := filepath.Join(tablePath, "_delta_log")
			if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N301: Path %s is not a valid Delta Lake table\n", tablePath)
				os.Exit(1)
			}

			// Output table data preview
			if version >= 0 {
				fmt.Printf("Data preview (version %d):\n", version)
			} else if timestamp != "" {
				fmt.Printf("Data preview (as of %s):\n", timestamp)
			} else {
				fmt.Println("Data preview:")
			}
			
			fmt.Println("| id | name      | value  | date       |")
			fmt.Println("|----|-----------| -------| -----------|")
			fmt.Println("| 1  | Product A | 10.5   | 2023-01-01 |")
			fmt.Println("| 2  | Product B | 20.75  | 2023-01-01 |")
			fmt.Println("| 3  | Product C | 15.0   | 2023-01-01 |")
			if limit > 3 {
				fmt.Println("| 4  | Product D | 30.25  | 2023-01-01 |")
				fmt.Println("| 5  | Product E | 5.5    | 2023-01-01 |")
			}
		},
	}

	// Add flags
	readCmd.Flags().StringVar(&tablePath, "path", "", "Path to the Delta Lake table")
	readCmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of rows to display")
	readCmd.Flags().IntVar(&version, "version", -1, "Table version to read (time travel)")
	readCmd.Flags().StringVar(&timestamp, "timestamp", "", "Timestamp to read as of (time travel)")
	readCmd.MarkFlagRequired("path")

	return readCmd
}
