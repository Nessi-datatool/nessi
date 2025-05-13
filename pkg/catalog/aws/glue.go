package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
	"github.com/nessi-dev/nessi-dev/pkg/catalog"
)

// GlueCatalog implements the DataCatalog interface for AWS Glue Data Catalog
type GlueCatalog struct {
	client      *glue.Client
	connected   bool
	region      string
}

// NewGlueCatalog creates a new AWS Glue Data Catalog client
func NewGlueCatalog() *GlueCatalog {
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
	var cfg aws.Config
	var err error
	
	if accessKey != "" && secretKey != "" {
		// Use access key and secret key
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				accessKey, secretKey, "",
			)),
		)
	} else if useIAMRole {
		// Use IAM role
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	} else {
		return fmt.Errorf("no valid authentication method provided")
	}
	
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}
	
	// Create Glue client
	c.client = glue.NewFromConfig(cfg)
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
func (c *GlueCatalog) ListDatabases(ctx context.Context) ([]catalog.DatabaseInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to AWS Glue Data Catalog")
	}
	
	// Call Glue API to list databases
	result, err := c.client.GetDatabases(ctx, &glue.GetDatabasesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	
	// Convert to DatabaseInfo
	databases := make([]catalog.DatabaseInfo, 0, len(result.DatabaseList))
	for _, db := range result.DatabaseList {
		databases = append(databases, catalog.DatabaseInfo{
			Name:        *db.Name,
			Description: aws.ToString(db.Description),
			Properties:  convertMapToStringMap(db.Parameters),
		})
	}
	
	return databases, nil
}

// ListTables lists all tables in a database
func (c *GlueCatalog) ListTables(ctx context.Context, database string) ([]catalog.TableInfo, error) {
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
	tables := make([]catalog.TableInfo, 0, len(result.TableList))
	for _, table := range result.TableList {
		location := ""
		if table.StorageDescriptor != nil && table.StorageDescriptor.Location != nil {
			location = *table.StorageDescriptor.Location
		}
		
		tableType := "unknown"
		if table.Parameters != nil {
			if format, ok := table.Parameters["table_type"]; ok {
				tableType = format
			}
		}
		
		tables = append(tables, catalog.TableInfo{
			Name:        *table.Name,
			Type:        tableType,
			Description: aws.ToString(table.Description),
			Location:    location,
			Properties:  convertMapToStringMap(table.Parameters),
		})
	}
	
	return tables, nil
}

// GetTableDetails gets detailed information about a table
func (c *GlueCatalog) GetTableDetails(ctx context.Context, database, table string) (*catalog.TableDetails, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to AWS Glue Data Catalog")
	}
	
	// Call Glue API to get table
	result, err := c.client.GetTable(ctx, &glue.GetTableInput{
		DatabaseName: aws.String(database),
		Name:         aws.String(table),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get table %s in database %s: %w", table, database, err)
	}
	
	if result.Table == nil {
		return nil, fmt.Errorf("table %s not found in database %s", table, database)
	}
	
	// Get table info
	location := ""
	if result.Table.StorageDescriptor != nil && result.Table.StorageDescriptor.Location != nil {
		location = *result.Table.StorageDescriptor.Location
	}
	
	tableType := "unknown"
	if result.Table.Parameters != nil {
		if format, ok := result.Table.Parameters["table_type"]; ok {
			tableType = format
		}
	}
	
	tableInfo := catalog.TableInfo{
		Name:        *result.Table.Name,
		Type:        tableType,
		Description: aws.ToString(result.Table.Description),
		Location:    location,
		Properties:  convertMapToStringMap(result.Table.Parameters),
	}
	
	// Get schema
	schema := &catalog.TableSchema{
		Format:  tableType,
		Version: 1, // Glue doesn't have schema versioning
	}
	
	if result.Table.StorageDescriptor != nil && result.Table.StorageDescriptor.Columns != nil {
		for _, col := range result.Table.StorageDescriptor.Columns {
			field := catalog.FieldInfo{
				Name:        *col.Name,
				Type:        *col.Type,
				Description: aws.ToString(col.Comment),
				Nullable:    true, // Glue doesn't store nullability
				Properties:  convertMapToStringMap(col.Parameters),
			}
			schema.Fields = append(schema.Fields, field)
		}
	}
	
	// Get metadata
	metadata := &catalog.TableMetadata{
		Owner:      aws.ToString(result.Table.Owner),
		CreatedAt:  aws.ToTime(result.Table.CreateTime),
		UpdatedAt:  aws.ToTime(result.Table.UpdateTime),
		Properties: convertMapToStringMap(result.Table.Parameters),
	}
	
	// Create TableDetails
	details := &catalog.TableDetails{
		Info:     tableInfo,
		Schema:   schema,
		Metadata: metadata,
	}
	
	return details, nil
}

// GetTableMetadata gets metadata for a table
func (c *GlueCatalog) GetTableMetadata(ctx context.Context, database, table string) (*catalog.TableMetadata, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return nil, err
	}
	return details.Metadata, nil
}

// UpdateTableMetadata updates metadata for a table
func (c *GlueCatalog) UpdateTableMetadata(ctx context.Context, database, table string, metadata *catalog.TableMetadata) error {
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
		TableInput: &types.TableInput{
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
func (c *GlueCatalog) GetTableLineage(ctx context.Context, database, table string) (*catalog.LineageInfo, error) {
	// AWS Glue doesn't have built-in lineage capabilities
	// This would require integration with AWS Lake Formation or custom implementation
	return nil, fmt.Errorf("lineage information not available in AWS Glue Data Catalog")
}

// UpdateTableLineage updates lineage information for a table
func (c *GlueCatalog) UpdateTableLineage(ctx context.Context, database, table string, lineage *catalog.LineageInfo) error {
	// AWS Glue doesn't have built-in lineage capabilities
	return fmt.Errorf("lineage update not supported in AWS Glue Data Catalog")
}

// PublishQualityMetrics publishes data quality metrics for a table
func (c *GlueCatalog) PublishQualityMetrics(ctx context.Context, database, table string, metrics *catalog.QualityMetrics) error {
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
	parameters["quality:overall_score"] = fmt.Sprintf("%.2f", metrics.OverallScore)
	parameters["quality:completeness"] = fmt.Sprintf("%.2f", metrics.Completeness)
	parameters["quality:accuracy"] = fmt.Sprintf("%.2f", metrics.Accuracy)
	parameters["quality:consistency"] = fmt.Sprintf("%.2f", metrics.Consistency)
	parameters["quality:timeliness"] = fmt.Sprintf("%.2f", metrics.Timeliness)
	parameters["quality:last_updated"] = metrics.LastUpdated.Format(time.RFC3339)
	
	// Create update input
	updateInput := &glue.UpdateTableInput{
		DatabaseName: aws.String(database),
		TableInput: &types.TableInput{
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
func (c *GlueCatalog) GetQualityMetrics(ctx context.Context, database, table string) (*catalog.QualityMetrics, error) {
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
	metrics := &catalog.QualityMetrics{}
	
	if score, ok := params["quality:overall_score"]; ok {
		fmt.Sscanf(score, "%f", &metrics.OverallScore)
	}
	
	if completeness, ok := params["quality:completeness"]; ok {
		fmt.Sscanf(completeness, "%f", &metrics.Completeness)
	}
	
	if accuracy, ok := params["quality:accuracy"]; ok {
		fmt.Sscanf(accuracy, "%f", &metrics.Accuracy)
	}
	
	if consistency, ok := params["quality:consistency"]; ok {
		fmt.Sscanf(consistency, "%f", &metrics.Consistency)
	}
	
	if timeliness, ok := params["quality:timeliness"]; ok {
		fmt.Sscanf(timeliness, "%f", &metrics.Timeliness)
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
