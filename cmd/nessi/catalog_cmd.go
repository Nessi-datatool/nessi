package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/catalog"
	"github.com/spf13/cobra"
)

// catalogCmd represents the catalog command
var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Interact with data catalogs",
	Long:  `Commands for interacting with data catalogs like AWS Glue, Azure Purview, and Google Cloud Data Catalog.`,
}

// catalogConnectCmd represents the catalog connect command
var catalogConnectCmd = &cobra.Command{
	Use:   "connect [catalog_type] [config_file]",
	Short: "Connect to a data catalog",
	Long: `Connect to a data catalog using the specified configuration.
Supported catalog types: aws_glue, azure_purview, gcp_data_catalog.

Example:
  nessi catalog connect aws_glue config.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		catalogType := args[0]
		configFile := args[1]

		// Read configuration file
		configData, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		// Parse configuration
		var config map[string]interface{}
		if err := json.Unmarshal(configData, &config); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}

		// Create catalog factory
		factory := catalog.NewCatalogFactory()

		// Create catalog manager
		manager := catalog.NewCatalogManager()

		// Create catalog
		var cat catalog.DataCatalog
		switch catalog.CatalogType(catalogType) {
		case catalog.AWSGlue:
			cat, err = factory.CreateCatalog(catalog.AWSGlue)
		case catalog.AzurePurview:
			cat, err = factory.CreateCatalog(catalog.AzurePurview)
		case catalog.GCPDataCatalog:
			cat, err = factory.CreateCatalog(catalog.GCPDataCatalog)
		default:
			return fmt.Errorf("unsupported catalog type: %s", catalogType)
		}

		if err != nil {
			return fmt.Errorf("failed to create catalog: %w", err)
		}

		// Register catalog
		manager.RegisterCatalog(cat)

		// Connect to catalog
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := cat.Connect(ctx, config); err != nil {
			return fmt.Errorf("failed to connect to catalog: %w", err)
		}

		fmt.Printf("Successfully connected to %s\n", cat.Name())
		return nil
	},
}

// catalogListDatabasesCmd represents the catalog list-databases command
var catalogListDatabasesCmd = &cobra.Command{
	Use:   "list-databases [catalog_type]",
	Short: "List databases in a data catalog",
	Long: `List databases in a data catalog.
Supported catalog types: aws_glue, azure_purview, gcp_data_catalog.

Example:
  nessi catalog list-databases aws_glue`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		catalogType := args[0]

		// Create catalog factory
		factory := catalog.NewCatalogFactory()

		// Create catalog
		var cat catalog.DataCatalog
		var err error
		switch catalog.CatalogType(catalogType) {
		case catalog.AWSGlue:
			cat, err = factory.CreateCatalog(catalog.AWSGlue)
		case catalog.AzurePurview:
			cat, err = factory.CreateCatalog(catalog.AzurePurview)
		case catalog.GCPDataCatalog:
			cat, err = factory.CreateCatalog(catalog.GCPDataCatalog)
		default:
			return fmt.Errorf("unsupported catalog type: %s", catalogType)
		}

		if err != nil {
			return fmt.Errorf("failed to create catalog: %w", err)
		}

		// Get configuration from flags
		config := make(map[string]interface{})
		if cmd.Flag("config").Changed {
			configFile := cmd.Flag("config").Value.String()
			configData, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			if err := json.Unmarshal(configData, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}
		} else {
			// Try to get configuration from environment variables
			switch catalog.CatalogType(catalogType) {
			case catalog.AWSGlue:
				config["region"] = os.Getenv("AWS_REGION")
				config["use_iam_role"] = true
			case catalog.AzurePurview:
				config["account_name"] = os.Getenv("AZURE_PURVIEW_ACCOUNT")
				config["use_azure_ad"] = true
			case catalog.GCPDataCatalog:
				config["project_id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
			}
		}

		// Connect to catalog
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := cat.Connect(ctx, config); err != nil {
			return fmt.Errorf("failed to connect to catalog: %w", err)
		}

		// List databases
		databases, err := cat.ListDatabases(ctx)
		if err != nil {
			return fmt.Errorf("failed to list databases: %w", err)
		}

		// Print databases
		fmt.Printf("Databases in %s:\n", cat.Name())
		for i, db := range databases {
			fmt.Printf("%d. %s - %s\n", i+1, db.Name, db.Description)
		}

		return nil
	},
}

// catalogListTablesCmd represents the catalog list-tables command
var catalogListTablesCmd = &cobra.Command{
	Use:   "list-tables [catalog_type] [database]",
	Short: "List tables in a database",
	Long: `List tables in a database in a data catalog.
Supported catalog types: aws_glue, azure_purview, gcp_data_catalog.

Example:
  nessi catalog list-tables aws_glue my_database`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		catalogType := args[0]
		database := args[1]

		// Create catalog factory
		factory := catalog.NewCatalogFactory()

		// Create catalog
		var cat catalog.DataCatalog
		var err error
		switch catalog.CatalogType(catalogType) {
		case catalog.AWSGlue:
			cat, err = factory.CreateCatalog(catalog.AWSGlue)
		case catalog.AzurePurview:
			cat, err = factory.CreateCatalog(catalog.AzurePurview)
		case catalog.GCPDataCatalog:
			cat, err = factory.CreateCatalog(catalog.GCPDataCatalog)
		default:
			return fmt.Errorf("unsupported catalog type: %s", catalogType)
		}

		if err != nil {
			return fmt.Errorf("failed to create catalog: %w", err)
		}

		// Get configuration from flags
		config := make(map[string]interface{})
		if cmd.Flag("config").Changed {
			configFile := cmd.Flag("config").Value.String()
			configData, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			if err := json.Unmarshal(configData, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}
		} else {
			// Try to get configuration from environment variables
			switch catalog.CatalogType(catalogType) {
			case catalog.AWSGlue:
				config["region"] = os.Getenv("AWS_REGION")
				config["use_iam_role"] = true
			case catalog.AzurePurview:
				config["account_name"] = os.Getenv("AZURE_PURVIEW_ACCOUNT")
				config["use_azure_ad"] = true
			case catalog.GCPDataCatalog:
				config["project_id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
			}
		}

		// Connect to catalog
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := cat.Connect(ctx, config); err != nil {
			return fmt.Errorf("failed to connect to catalog: %w", err)
		}

		// List tables
		tables, err := cat.ListTables(ctx, database)
		if err != nil {
			return fmt.Errorf("failed to list tables: %w", err)
		}

		// Print tables
		fmt.Printf("Tables in %s.%s:\n", cat.Name(), database)
		for i, table := range tables {
			fmt.Printf("%d. %s - %s (%s)\n", i+1, table.Name, table.Description, table.Type)
			fmt.Printf("   Location: %s\n", table.Location)
		}

		return nil
	},
}

// catalogGetTableCmd represents the catalog get-table command
var catalogGetTableCmd = &cobra.Command{
	Use:   "get-table [catalog_type] [database] [table]",
	Short: "Get table details",
	Long: `Get detailed information about a table in a data catalog.
Supported catalog types: aws_glue, azure_purview, gcp_data_catalog.

Example:
  nessi catalog get-table aws_glue my_database my_table`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		catalogType := args[0]
		database := args[1]
		table := args[2]

		// Create catalog factory
		factory := catalog.NewCatalogFactory()

		// Create catalog
		var cat catalog.DataCatalog
		var err error
		switch catalog.CatalogType(catalogType) {
		case catalog.AWSGlue:
			cat, err = factory.CreateCatalog(catalog.AWSGlue)
		case catalog.AzurePurview:
			cat, err = factory.CreateCatalog(catalog.AzurePurview)
		case catalog.GCPDataCatalog:
			cat, err = factory.CreateCatalog(catalog.GCPDataCatalog)
		default:
			return fmt.Errorf("unsupported catalog type: %s", catalogType)
		}

		if err != nil {
			return fmt.Errorf("failed to create catalog: %w", err)
		}

		// Get configuration from flags
		config := make(map[string]interface{})
		if cmd.Flag("config").Changed {
			configFile := cmd.Flag("config").Value.String()
			configData, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			if err := json.Unmarshal(configData, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}
		} else {
			// Try to get configuration from environment variables
			switch catalog.CatalogType(catalogType) {
			case catalog.AWSGlue:
				config["region"] = os.Getenv("AWS_REGION")
				config["use_iam_role"] = true
			case catalog.AzurePurview:
				config["account_name"] = os.Getenv("AZURE_PURVIEW_ACCOUNT")
				config["use_azure_ad"] = true
			case catalog.GCPDataCatalog:
				config["project_id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
			}
		}

		// Connect to catalog
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := cat.Connect(ctx, config); err != nil {
			return fmt.Errorf("failed to connect to catalog: %w", err)
		}

		// Get table details
		details, err := cat.GetTableDetails(ctx, database, table)
		if err != nil {
			return fmt.Errorf("failed to get table details: %w", err)
		}

		// Print table details
		fmt.Printf("Table: %s.%s.%s\n", cat.Name(), database, table)
		fmt.Printf("Description: %s\n", details.Info.Description)
		fmt.Printf("Type: %s\n", details.Info.Type)
		fmt.Printf("Location: %s\n", details.Info.Location)
		fmt.Printf("Owner: %s\n", details.Metadata.Owner)
		fmt.Printf("Created: %s\n", details.Metadata.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Updated: %s\n", details.Metadata.UpdatedAt.Format(time.RFC3339))
		
		// Print tags
		if len(details.Metadata.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(details.Metadata.Tags, ", "))
		}
		
		// Print schema
		if details.Schema != nil && len(details.Schema.Fields) > 0 {
			fmt.Printf("\nSchema:\n")
			for i, field := range details.Schema.Fields {
				nullableStr := "NULLABLE"
				if !field.Nullable {
					nullableStr = "NOT NULL"
				}
				fmt.Printf("%d. %s: %s %s", i+1, field.Name, field.Type, nullableStr)
				if field.Description != "" {
					fmt.Printf(" - %s", field.Description)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

// catalogPublishQualityCmd represents the catalog publish-quality command
var catalogPublishQualityCmd = &cobra.Command{
	Use:   "publish-quality [catalog_type] [database] [table] [profile_file] [validation_file]",
	Short: "Publish quality metrics to a data catalog",
	Long: `Publish data quality metrics to a data catalog.
Supported catalog types: aws_glue, azure_purview, gcp_data_catalog.

Example:
  nessi catalog publish-quality aws_glue my_database my_table profile.json validation.json`,
	Args: cobra.ExactArgs(5),
	RunE: func(cmd *cobra.Command, args []string) error {
		catalogType := args[0]
		database := args[1]
		table := args[2]
		profileFile := args[3]
		validationFile := args[4]

		// Read profile file
		profileData, err := os.ReadFile(profileFile)
		if err != nil {
			return fmt.Errorf("failed to read profile file: %w", err)
		}

		// Parse profile
		var profile struct {
			Timestamp string `json:"timestamp"`
			Columns   []struct {
				Name  string `json:"name"`
				Stats struct {
					Count     int64 `json:"count"`
					NullCount int64 `json:"null_count"`
				} `json:"stats"`
			} `json:"columns"`
		}
		if err := json.Unmarshal(profileData, &profile); err != nil {
			return fmt.Errorf("failed to parse profile file: %w", err)
		}

		// Read validation file
		validationData, err := os.ReadFile(validationFile)
		if err != nil {
			return fmt.Errorf("failed to read validation file: %w", err)
		}

		// Parse validation results
		var validation struct {
			RuleResults []struct {
				Rule struct {
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"rule"`
				Passed  bool    `json:"passed"`
				Score   float64 `json:"score"`
				Details string  `json:"details"`
			} `json:"rule_results"`
		}
		if err := json.Unmarshal(validationData, &validation); err != nil {
			return fmt.Errorf("failed to parse validation file: %w", err)
		}

		// Create catalog factory
		factory := catalog.NewCatalogFactory()

		// Create catalog
		var cat catalog.DataCatalog
		switch catalog.CatalogType(catalogType) {
		case catalog.AWSGlue:
			cat, err = factory.CreateCatalog(catalog.AWSGlue)
		case catalog.AzurePurview:
			cat, err = factory.CreateCatalog(catalog.AzurePurview)
		case catalog.GCPDataCatalog:
			cat, err = factory.CreateCatalog(catalog.GCPDataCatalog)
		default:
			return fmt.Errorf("unsupported catalog type: %s", catalogType)
		}

		if err != nil {
			return fmt.Errorf("failed to create catalog: %w", err)
		}

		// Get configuration from flags
		config := make(map[string]interface{})
		if cmd.Flag("config").Changed {
			configFile := cmd.Flag("config").Value.String()
			configData, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			if err := json.Unmarshal(configData, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}
		} else {
			// Try to get configuration from environment variables
			switch catalog.CatalogType(catalogType) {
			case catalog.AWSGlue:
				config["region"] = os.Getenv("AWS_REGION")
				config["use_iam_role"] = true
			case catalog.AzurePurview:
				config["account_name"] = os.Getenv("AZURE_PURVIEW_ACCOUNT")
				config["use_azure_ad"] = true
			case catalog.GCPDataCatalog:
				config["project_id"] = os.Getenv("GOOGLE_CLOUD_PROJECT")
			}
		}

		// Connect to catalog
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := cat.Connect(ctx, config); err != nil {
			return fmt.Errorf("failed to connect to catalog: %w", err)
		}

		// Create quality metrics
		metrics := &catalog.QualityMetrics{
			OverallScore: 0.0,
			Completeness: 0.0,
			Accuracy:     0.0,
			Consistency:  0.0,
			Timeliness:   0.0,
			LastUpdated:  time.Now(),
		}

		// Calculate completeness
		if len(profile.Columns) > 0 {
			var totalCompleteness float64
			for _, col := range profile.Columns {
				if col.Stats.Count > 0 {
					completeness := 1.0 - (float64(col.Stats.NullCount) / float64(col.Stats.Count))
					totalCompleteness += completeness
				}
			}
			metrics.Completeness = totalCompleteness / float64(len(profile.Columns))
		}

		// Calculate accuracy and consistency
		if len(validation.RuleResults) > 0 {
			var accuracyScore, consistencyScore float64
			var accuracyRules, consistencyRules int

			for _, result := range validation.RuleResults {
				// Add rule result
				metrics.RuleResults = append(metrics.RuleResults, catalog.RuleResult{
					RuleName: result.Rule.Name,
					RuleType: result.Rule.Type,
					Passed:   result.Passed,
					Score:    result.Score,
					Details:  result.Details,
				})

				// Calculate accuracy and consistency
				switch result.Rule.Type {
				case "range", "enum", "regex", "format":
					accuracyScore += result.Score
					accuracyRules++
				case "unique", "relationship", "referential_integrity":
					consistencyScore += result.Score
					consistencyRules++
				}
			}

			if accuracyRules > 0 {
				metrics.Accuracy = accuracyScore / float64(accuracyRules)
			} else {
				metrics.Accuracy = 1.0
			}

			if consistencyRules > 0 {
				metrics.Consistency = consistencyScore / float64(consistencyRules)
			} else {
				metrics.Consistency = 1.0
			}
		}

		// Calculate timeliness
		if profile.Timestamp != "" {
			timestamp, err := time.Parse(time.RFC3339, profile.Timestamp)
			if err == nil {
				ageHours := time.Since(timestamp).Hours()
				if ageHours <= 24.0 {
					metrics.Timeliness = 1.0
				} else if ageHours <= 168.0 {
					metrics.Timeliness = 0.8
				} else if ageHours <= 720.0 {
					metrics.Timeliness = 0.6
				} else {
					metrics.Timeliness = 0.4
				}
			}
		} else {
			metrics.Timeliness = 1.0
		}

		// Calculate overall score
		metrics.OverallScore = (metrics.Completeness*0.25 + metrics.Accuracy*0.35 +
			metrics.Consistency*0.25 + metrics.Timeliness*0.15)

		// Publish quality metrics
		if err := cat.PublishQualityMetrics(ctx, database, table, metrics); err != nil {
			return fmt.Errorf("failed to publish quality metrics: %w", err)
		}

		fmt.Printf("Successfully published quality metrics to %s for %s.%s\n", cat.Name(), database, table)
		fmt.Printf("Overall Score: %.2f\n", metrics.OverallScore)
		fmt.Printf("Completeness: %.2f\n", metrics.Completeness)
		fmt.Printf("Accuracy: %.2f\n", metrics.Accuracy)
		fmt.Printf("Consistency: %.2f\n", metrics.Consistency)
		fmt.Printf("Timeliness: %.2f\n", metrics.Timeliness)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catalogCmd)
	catalogCmd.AddCommand(catalogConnectCmd)
	catalogCmd.AddCommand(catalogListDatabasesCmd)
	catalogCmd.AddCommand(catalogListTablesCmd)
	catalogCmd.AddCommand(catalogGetTableCmd)
	catalogCmd.AddCommand(catalogPublishQualityCmd)

	// Add flags
	catalogListDatabasesCmd.Flags().String("config", "", "Path to configuration file")
	catalogListTablesCmd.Flags().String("config", "", "Path to configuration file")
	catalogGetTableCmd.Flags().String("config", "", "Path to configuration file")
	catalogPublishQualityCmd.Flags().String("config", "", "Path to configuration file")
}
