package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nessi-dev/nessi/pkg/cloud"
	"github.com/nessi-dev/nessi/pkg/cloud/common"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var cloudManager *cloud.CloudManager

// cloudCmd represents the cloud command
var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Manage cloud provider integrations",
	Long: `Manage cloud provider integrations for Nessi.dev.

This command allows you to configure and manage connections to cloud providers
like AWS, Azure, and GCP for accessing Delta Lake tables and other data stored
in cloud storage.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize cloud manager if not already initialized
		if cloudManager == nil {
			cloudManager = cloud.NewCloudManager()
		}
	},
}

// cloudConfigureCmd represents the configure command for cloud providers
var cloudConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure a cloud provider",
	Long: `Configure a cloud provider for use with Nessi.dev.

This command allows you to set up a connection to a cloud provider like AWS, Azure, or GCP.
You can either provide a configuration file or enter the configuration interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config-file")
		providerName, _ := cmd.Flags().GetString("name")
		providerType, _ := cmd.Flags().GetString("provider")

		if configFile != "" {
			// Load configuration from file
			configData, err := os.ReadFile(configFile)
			if err != nil {
				fmt.Printf("Error reading config file: %v\n", err)
				os.Exit(1)
			}

			var config common.CloudConfig
			if strings.HasSuffix(configFile, ".yaml") || strings.HasSuffix(configFile, ".yml") {
				err = yaml.Unmarshal(configData, &config)
			} else {
				err = json.Unmarshal(configData, &config)
			}

			if err != nil {
				fmt.Printf("Error parsing config file: %v\n", err)
				os.Exit(1)
			}

			// Override provider type if specified
			if providerType != "" {
				config.Provider = providerType
			}

			// Create provider
			_, err = cloudManager.CreateProvider(providerName, config)
			if err != nil {
				fmt.Printf("Error creating cloud provider: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Cloud provider '%s' configured successfully.\n", providerName)
		} else {
			// Interactive configuration
			if providerName == "" {
				fmt.Print("Enter a name for this cloud provider configuration: ")
				fmt.Scanln(&providerName)
			}

			if providerType == "" {
				fmt.Println("Available provider types:")
				for _, p := range cloudManager.ListSupportedProviders() {
					fmt.Printf("  - %s\n", p)
				}
				fmt.Print("Enter provider type (aws, azure, gcp): ")
				fmt.Scanln(&providerType)
			}

			config := common.CloudConfig{
				Provider:    providerType,
				Credentials: make(map[string]interface{}),
			}

			// Provider-specific configuration
			switch providerType {
			case "aws":
				configureAWS(&config)
			case "azure":
				configureAzure(&config)
			case "gcp":
				configureGCP(&config)
			default:
				fmt.Printf("Unsupported provider type: %s\n", providerType)
				os.Exit(1)
			}

			// Create provider
			_, err := cloudManager.CreateProvider(providerName, config)
			if err != nil {
				fmt.Printf("Error creating cloud provider: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Cloud provider '%s' configured successfully.\n", providerName)

			// Save configuration to file
			saveConfig, _ := cmd.Flags().GetBool("save")
			if saveConfig {
				saveConfigFile, _ := cmd.Flags().GetString("save-file")
				if saveConfigFile == "" {
					saveConfigFile = fmt.Sprintf("%s_config.yaml", providerName)
				}

				configData, err := yaml.Marshal(config)
				if err != nil {
					fmt.Printf("Error serializing configuration: %v\n", err)
					os.Exit(1)
				}

				err = os.WriteFile(saveConfigFile, configData, 0600)
				if err != nil {
					fmt.Printf("Error saving configuration: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Configuration saved to %s\n", saveConfigFile)
			}
		}
	},
}

// cloudListCmd represents the list command for cloud providers
var cloudListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured cloud providers",
	Long:  `List all cloud providers configured for use with Nessi.dev.`,
	Run: func(cmd *cobra.Command, args []string) {
		providers := cloudManager.ListProviders()

		if len(providers) == 0 {
			fmt.Println("No cloud providers configured.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tSTATUS")

		for _, name := range providers {
			provider, _ := cloudManager.GetProvider(name)
			fmt.Fprintf(w, "%s\t%s\tConnected\n", name, provider.Name())
		}

		w.Flush()
	},
}

// cloudRemoveCmd represents the remove command for cloud providers
var cloudRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a cloud provider configuration",
	Long:  `Remove a cloud provider configuration from Nessi.dev.`,
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("name")

		if providerName == "" {
			fmt.Println("Provider name is required.")
			os.Exit(1)
		}

		err := cloudManager.RemoveProvider(providerName)
		if err != nil {
			fmt.Printf("Error removing cloud provider: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cloud provider '%s' removed successfully.\n", providerName)
	},
}

// cloudTestCmd represents the test command for cloud providers
var cloudTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test a cloud provider connection",
	Long:  `Test the connection to a cloud provider.`,
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("name")

		if providerName == "" {
			fmt.Println("Provider name is required.")
			os.Exit(1)
		}

		provider, exists := cloudManager.GetProvider(providerName)
		if !exists {
			fmt.Printf("Cloud provider '%s' not found.\n", providerName)
			os.Exit(1)
		}

		// Test connection by listing buckets
		buckets, err := provider.ListBuckets(cmd.Context())
		if err != nil {
			fmt.Printf("Connection test failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Connection test successful.")
		fmt.Printf("Found %d buckets:\n", len(buckets))

		for _, bucket := range buckets {
			fmt.Printf("  - %s (created: %s)\n", bucket.Name, bucket.CreationDate)
		}
	},
}

// cloudDeltaListCmd represents the delta list command for cloud providers
var cloudDeltaListCmd = &cobra.Command{
	Use:   "delta-list",
	Short: "List Delta Lake tables in cloud storage",
	Long:  `List Delta Lake tables stored in cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("name")
		bucket, _ := cmd.Flags().GetString("bucket")
		prefix, _ := cmd.Flags().GetString("prefix")

		if providerName == "" {
			fmt.Println("Provider name is required.")
			os.Exit(1)
		}

		if bucket == "" {
			fmt.Println("Bucket name is required.")
			os.Exit(1)
		}

		provider, exists := cloudManager.GetProvider(providerName)
		if !exists {
			fmt.Printf("Cloud provider '%s' not found.\n", providerName)
			os.Exit(1)
		}

		// List objects with the given prefix
		objects, err := provider.ListObjects(cmd.Context(), bucket, prefix)
		if err != nil {
			fmt.Printf("Error listing objects: %v\n", err)
			os.Exit(1)
		}

		// Find Delta Lake tables (look for _delta_log directories)
		deltaLogPaths := make(map[string]bool)
		for _, obj := range objects {
			if strings.Contains(obj.Key, "/_delta_log/") {
				// Extract table path (parent directory of _delta_log)
				tablePath := strings.Split(obj.Key, "/_delta_log/")[0]
				deltaLogPaths[tablePath] = true
			}
		}

		if len(deltaLogPaths) == 0 {
			fmt.Println("No Delta Lake tables found.")
			return
		}

		fmt.Printf("Found %d Delta Lake tables:\n", len(deltaLogPaths))

		for tablePath := range deltaLogPaths {
			fmt.Printf("  - %s\n", tablePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(cloudCmd)

	// Configure command
	cloudCmd.AddCommand(cloudConfigureCmd)
	cloudConfigureCmd.Flags().String("config-file", "", "Path to cloud provider configuration file")
	cloudConfigureCmd.Flags().String("name", "", "Name for this cloud provider configuration")
	cloudConfigureCmd.Flags().String("provider", "", "Cloud provider type (aws, azure, gcp)")
	cloudConfigureCmd.Flags().Bool("save", false, "Save configuration to file")
	cloudConfigureCmd.Flags().String("save-file", "", "Path to save configuration file")

	// List command
	cloudCmd.AddCommand(cloudListCmd)

	// Remove command
	cloudCmd.AddCommand(cloudRemoveCmd)
	cloudRemoveCmd.Flags().String("name", "", "Name of the cloud provider to remove")
	cloudRemoveCmd.MarkFlagRequired("name")

	// Test command
	cloudCmd.AddCommand(cloudTestCmd)
	cloudTestCmd.Flags().String("name", "", "Name of the cloud provider to test")
	cloudTestCmd.MarkFlagRequired("name")

	// Delta list command
	cloudCmd.AddCommand(cloudDeltaListCmd)
	cloudDeltaListCmd.Flags().String("name", "", "Name of the cloud provider")
	cloudDeltaListCmd.Flags().String("bucket", "", "Bucket/container name")
	cloudDeltaListCmd.Flags().String("prefix", "", "Object prefix (optional)")
	cloudDeltaListCmd.MarkFlagRequired("name")
	cloudDeltaListCmd.MarkFlagRequired("bucket")
}

// configureAWS configures AWS-specific settings
func configureAWS(config *common.CloudConfig) {
	var accessKey, secretKey, region, endpoint string
	var useIAMRole bool

	fmt.Print("Use IAM role for authentication? (y/n): ")
	var useIAMRoleStr string
	fmt.Scanln(&useIAMRoleStr)
	useIAMRole = strings.ToLower(useIAMRoleStr) == "y"

	if !useIAMRole {
		fmt.Print("AWS Access Key ID: ")
		fmt.Scanln(&accessKey)

		fmt.Print("AWS Secret Access Key: ")
		fmt.Scanln(&secretKey)
	}

	fmt.Print("AWS Region (e.g., us-west-2): ")
	fmt.Scanln(&region)

	fmt.Print("S3 Endpoint (optional, press Enter to skip): ")
	fmt.Scanln(&endpoint)

	fmt.Print("Default S3 Bucket: ")
	fmt.Scanln(&config.DefaultBucket)

	config.Region = region
	config.Credentials["access_key"] = accessKey
	config.Credentials["secret_key"] = secretKey
	config.Credentials["use_iam_role"] = useIAMRole

	if endpoint != "" {
		config.EndpointOverride = endpoint
	}
}

// configureAzure configures Azure-specific settings
func configureAzure(config *common.CloudConfig) {
	var accountName, accountKey, sasToken, endpoint string
	var useAzureAD bool

	fmt.Print("Azure Storage Account Name: ")
	fmt.Scanln(&accountName)

	fmt.Print("Authentication Method (key, sas, azuread): ")
	var authMethod string
	fmt.Scanln(&authMethod)

	switch strings.ToLower(authMethod) {
	case "key":
		fmt.Print("Azure Storage Account Key: ")
		fmt.Scanln(&accountKey)
	case "sas":
		fmt.Print("Azure SAS Token (without leading ?): ")
		fmt.Scanln(&sasToken)
	case "azuread":
		useAzureAD = true
		fmt.Println("Using Azure AD authentication.")
	default:
		fmt.Println("Invalid authentication method. Using Azure AD.")
		useAzureAD = true
	}

	fmt.Print("Azure Storage Endpoint (optional, press Enter to skip): ")
	fmt.Scanln(&endpoint)

	fmt.Print("Default Container: ")
	fmt.Scanln(&config.DefaultBucket)

	config.Credentials["account_name"] = accountName
	config.Credentials["account_key"] = accountKey
	config.Credentials["sas_token"] = sasToken
	config.Credentials["use_azure_ad"] = useAzureAD

	if endpoint != "" {
		config.EndpointOverride = endpoint
	}
}

// configureGCP configures GCP-specific settings
func configureGCP(config *common.CloudConfig) {
	var projectID, credentialsFile string

	fmt.Print("GCP Project ID: ")
	fmt.Scanln(&projectID)

	fmt.Print("GCP Credentials File Path (optional, press Enter to use default): ")
	fmt.Scanln(&credentialsFile)

	fmt.Print("Default GCS Bucket: ")
	fmt.Scanln(&config.DefaultBucket)

	config.AdditionalOptions = make(map[string]interface{})
	config.AdditionalOptions["project_id"] = projectID

	if credentialsFile != "" {
		config.Credentials["credentials_file"] = credentialsFile
	}
}
