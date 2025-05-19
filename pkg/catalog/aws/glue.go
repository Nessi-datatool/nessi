package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
	nessitypes "github.com/nessi-dev/nessi/pkg/api/types"
)

// GlueAPI defines the interface for the AWS Glue client.
// This interface is used to abstract the AWS SDK's Glue client for testing purposes.
type GlueAPI interface {
	GetDatabases(ctx context.Context, params *glue.GetDatabasesInput, optFns ...func(*glue.Options)) (*glue.GetDatabasesOutput, error)
	GetTables(ctx context.Context, params *glue.GetTablesInput, optFns ...func(*glue.Options)) (*glue.GetTablesOutput, error)
	GetTable(ctx context.Context, params *glue.GetTableInput, optFns ...func(*glue.Options)) (*glue.GetTableOutput, error)
	GetTableDetails(ctx context.Context, database, table string) (*nessitypes.TableDetails, error)
	UpdateTable(ctx context.Context, params *glue.UpdateTableInput, optFns ...func(*glue.Options)) (*glue.UpdateTableOutput, error)
}

type glueAPIImpl struct {
	*glue.Client
}

func (g *glueAPIImpl) GetTableDetails(ctx context.Context, database, table string) (*nessitypes.TableDetails, error) {
	output, err := g.GetTable(ctx, &glue.GetTableInput{
		DatabaseName: aws.String(database),
		Name:         aws.String(table),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get table details: %w", err)
	}

	// Convert AWS Glue Table to nessitypes.TableDetails
	tableDetails := &nessitypes.TableDetails{
		Info: nessitypes.TableInfo{
			Name:        aws.ToString(output.Table.Name),
			Type:        output.Table.Parameters["table_type"],
			Description: aws.ToString(output.Table.Description),
			Location:    aws.ToString(output.Table.StorageDescriptor.Location),
		},
		Schema: &nessitypes.TableSchema{
			Format: aws.ToString(output.Table.StorageDescriptor.InputFormat),
		},
		Metadata: &nessitypes.TableMetadata{
			Owner:     aws.ToString(output.Table.Owner),
			CreatedAt: time.Now(), // AWS Glue doesn't provide creation time
			UpdatedAt: time.Now(),
		},
	}

	return tableDetails, nil
}

func init() {
	nessitypes.GetCatalogFactory().RegisterProvider(nessitypes.AWSGlue, NewGlueCatalog)
}

// GlueCatalog implements the DataCatalog interface for AWS Glue Data Catalog
type GlueCatalog struct {
	client    GlueAPI
	connected bool
	region    string
}

// NewGlueCatalog creates a new AWS Glue Data Catalog client
func NewGlueCatalog() nessitypes.DataCatalog {
	return &GlueCatalog{
		connected: false,
	}
}

// Name returns the name of the catalog
func (c *GlueCatalog) Name() string {
	return "AWS Glue Data Catalog"
}

// Connect establishes a connection to AWS Glue Data Catalog
func (c *GlueCatalog) Connect(ctx context.Context, config map[string]interface{}) error {
	// Extract configuration
	region, _ := config["region"].(string)
	accessKey, _ := config["access_key"].(string)
	secretKey, _ := config["secret_key"].(string)
	useIAMRole, _ := config["use_iam_role"].(bool)

	// Create AWS configuration
	var awsCfg aws.Config
	var err error

	if accessKey != "" && secretKey != "" {
		// Use access key and secret key
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(region),
			awsconfig.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
			),
		)
	} else if useIAMRole {
		// Use IAM role
		awsCfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(region),
		)
	} else {
		return fmt.Errorf("no valid authentication method provided")
	}

	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create Glue client
	c.client = &glueAPIImpl{glue.NewFromConfig(awsCfg)}
	c.connected = true
	c.region = region

	return nil
}

// Disconnect disconnects from AWS Glue Data Catalog
func (c *GlueCatalog) Disconnect(ctx context.Context) error {
	c.connected = false
	return nil
}

// ListDatabases lists all databases in the Glue Data Catalog
func (c *GlueCatalog) ListDatabases(ctx context.Context) ([]nessitypes.DatabaseInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to AWS Glue Data Catalog")
	}

	// Call Glue API to list databases
	result, err := c.client.GetDatabases(ctx, &glue.GetDatabasesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	// Convert to DatabaseInfo
	databases := make([]nessitypes.DatabaseInfo, 0, len(result.DatabaseList))
	for _, db := range result.DatabaseList {
		databases = append(databases, nessitypes.DatabaseInfo{
			Name:        *db.Name,
			Description: aws.ToString(db.Description),
			Properties:  convertMapToStringMap(db.Parameters),
		})
	}

	return databases, nil
}

// ListTables lists all tables in a database
func (c *GlueCatalog) ListTables(ctx context.Context, database string) ([]nessitypes.TableInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to AWS Glue Data Catalog")
	}

	// Call Glue API to list tables
	result, err := c.client.GetTables(ctx, &glue.GetTablesInput{
		DatabaseName: aws.String(database),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tables in database %s: %w", database, err)
	}

	// Convert to TableInfo
	tables := make([]nessitypes.TableInfo, 0, len(result.TableList))
	for _, table := range result.TableList {
		location := ""
		if table.StorageDescriptor != nil && table.StorageDescriptor.Location != nil {
			location = *table.StorageDescriptor.Location
		}

		tableType := ""
		if table.Parameters != nil {
			tableType = table.Parameters["table_type"]
		}

		tables = append(tables, nessitypes.TableInfo{
			Name:        *table.Name,
			Type:        tableType,
			Description: aws.ToString(table.Description),
			Location:    location,
			Properties:  convertMapToStringMap(table.Parameters),
		})
	}

	return tables, nil
}

// GetTableDetails retrieves detailed information about a table
func (c *GlueCatalog) GetTableDetails(ctx context.Context, database, table string) (*nessitypes.TableDetails, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to AWS Glue Data Catalog")
	}

	return c.client.GetTableDetails(ctx, database, table)
}

// GetTableMetadata gets metadata for a table
func (c *GlueCatalog) GetTableMetadata(ctx context.Context, database, table string) (*nessitypes.TableMetadata, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return nil, err
	}
	return details.Metadata, nil
}

// UpdateTableMetadata updates metadata for a table
func (c *GlueCatalog) UpdateTableMetadata(ctx context.Context, database, table string, metadata *nessitypes.TableMetadata) error {
	if !c.connected {
		return fmt.Errorf("not connected to AWS Glue Data Catalog")
	}

	// Get current table
	getResult, err := c.client.GetTable(ctx, &glue.GetTableInput{
		DatabaseName: aws.String(database),
		Name:         aws.String(table),
	})
	if err != nil {
		return fmt.Errorf("failed to get table %s in database %s: %w", table, database, err)
	}

	if getResult.Table == nil {
		return fmt.Errorf("table %s not found in database %s", table, database)
	}

	// Update table parameters with metadata
	parameters := make(map[string]string)
	for k, v := range metadata.Properties {
		parameters[k] = v
	}

	// Add tags as parameters with "tag:" prefix
	for _, tag := range metadata.Tags {
		parameters["tag:"+tag] = "true"
	}

	// Create update input
	updateInput := &glue.UpdateTableInput{
		DatabaseName: aws.String(database),
		TableInput: &gluetypes.TableInput{
			Name:        getResult.Table.Name,
			Description: aws.String(metadata.Properties["description"]),
			Owner:       aws.String(metadata.Owner),
			Parameters:  parameters,
			// Keep other fields from existing table
			StorageDescriptor: getResult.Table.StorageDescriptor,
			PartitionKeys:     getResult.Table.PartitionKeys,
			TableType:         getResult.Table.TableType,
		},
	}

	// Update table
	_, err = c.client.UpdateTable(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to update table %s in database %s: %w", table, database, err)
	}

	return nil
}

// GetTableLineage gets lineage information for a table
func (c *GlueCatalog) GetTableLineage(ctx context.Context, database, table string) (*nessitypes.LineageInfo, error) {
	// AWS Glue doesn't have built-in lineage capabilities
	// This would require integration with AWS Lake Formation or custom implementation
	return nil, fmt.Errorf("lineage information not available in AWS Glue Data Catalog")
}

// UpdateTableLineage updates lineage information for a table
func (c *GlueCatalog) UpdateTableLineage(ctx context.Context, database, table string, lineage *nessitypes.LineageInfo) error {
	// AWS Glue doesn't have built-in lineage capabilities
	return fmt.Errorf("lineage update not supported in AWS Glue Data Catalog")
}

// PublishQualityMetrics publishes data quality metrics for a table
func (c *GlueCatalog) PublishQualityMetrics(ctx context.Context, database, table string, metrics *nessitypes.QualityMetrics) error {
	if !c.connected {
		return fmt.Errorf("not connected to AWS Glue Data Catalog")
	}

	// Get current table
	getResult, err := c.client.GetTable(ctx, &glue.GetTableInput{
		DatabaseName: aws.String(database),
		Name:         aws.String(table),
	})
	if err != nil {
		return fmt.Errorf("failed to get table %s in database %s: %w", table, database, err)
	}

	if getResult.Table == nil {
		return fmt.Errorf("table %s not found in database %s", table, database)
	}

	// Create parameters map from existing parameters
	parameters := make(map[string]string)
	if getResult.Table.Parameters != nil {
		for k, v := range getResult.Table.Parameters {
			parameters[k] = v
		}
	}

	// Add quality metrics as parameters
	parameters["quality:total_rows"] = fmt.Sprintf("%d", metrics.TotalRows)
	parameters["quality:null_rows"] = fmt.Sprintf("%d", metrics.NullRows)
	parameters["quality:duplicate_rows"] = fmt.Sprintf("%d", metrics.DuplicateRows)
	parameters["quality:invalid_rows"] = fmt.Sprintf("%d", metrics.InvalidRows)
	parameters["quality:data_completeness"] = fmt.Sprintf("%.2f", metrics.DataCompleteness)
	parameters["quality:data_accuracy"] = fmt.Sprintf("%.2f", metrics.DataAccuracy)
	parameters["quality:data_consistency"] = fmt.Sprintf("%.2f", metrics.DataConsistency)
	parameters["quality:schema_version"] = metrics.SchemaVersion
	parameters["quality:last_updated"] = metrics.LastUpdated.Format(time.RFC3339)

	// Create update input
	updateInput := &glue.UpdateTableInput{
		DatabaseName: aws.String(database),
		TableInput: &gluetypes.TableInput{
			Name:              getResult.Table.Name,
			Description:       getResult.Table.Description,
			Owner:             getResult.Table.Owner,
			Parameters:        parameters,
			StorageDescriptor: getResult.Table.StorageDescriptor,
			PartitionKeys:     getResult.Table.PartitionKeys,
			TableType:         getResult.Table.TableType,
		},
	}

	// Update table
	_, err = c.client.UpdateTable(ctx, updateInput)
	if err != nil {
		return fmt.Errorf("failed to update table %s in database %s with quality metrics: %w", table, database, err)
	}

	return nil
}

// GetQualityMetrics gets data quality metrics for a table
func (c *GlueCatalog) GetQualityMetrics(ctx context.Context, database, table string) (*nessitypes.QualityMetrics, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to AWS Glue Data Catalog")
	}

	// Get table
	getResult, err := c.client.GetTable(ctx, &glue.GetTableInput{
		DatabaseName: aws.String(database),
		Name:         aws.String(table),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get table %s in database %s: %w", table, database, err)
	}

	if getResult.Table == nil || getResult.Table.Parameters == nil {
		return nil, fmt.Errorf("table %s not found in database %s or has no parameters", table, database)
	}

	// Extract quality metrics from parameters
	params := getResult.Table.Parameters
	metrics := &nessitypes.QualityMetrics{}

	if totalRows, ok := params["quality:total_rows"]; ok {
		fmt.Sscanf(totalRows, "%d", &metrics.TotalRows)
	}

	if nullRows, ok := params["quality:null_rows"]; ok {
		fmt.Sscanf(nullRows, "%d", &metrics.NullRows)
	}

	if duplicateRows, ok := params["quality:duplicate_rows"]; ok {
		fmt.Sscanf(duplicateRows, "%d", &metrics.DuplicateRows)
	}

	if invalidRows, ok := params["quality:invalid_rows"]; ok {
		fmt.Sscanf(invalidRows, "%d", &metrics.InvalidRows)
	}

	if dataCompleteness, ok := params["quality:data_completeness"]; ok {
		fmt.Sscanf(dataCompleteness, "%f", &metrics.DataCompleteness)
	}

	if dataAccuracy, ok := params["quality:data_accuracy"]; ok {
		fmt.Sscanf(dataAccuracy, "%f", &metrics.DataAccuracy)
	}

	if dataConsistency, ok := params["quality:data_consistency"]; ok {
		fmt.Sscanf(dataConsistency, "%f", &metrics.DataConsistency)
	}

	if schemaVersion, ok := params["quality:schema_version"]; ok {
		metrics.SchemaVersion = schemaVersion
	}

	if lastUpdated, ok := params["quality:last_updated"]; ok {
		metrics.LastUpdated, _ = time.Parse(time.RFC3339, lastUpdated)
	}

	return metrics, nil
}

// Helper function to convert AWS map to string map
func convertMapToStringMap(m map[string]string) map[string]string {
	if m == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for k, v := range m {
		result[k] = v
	}
	return result
}
