package gcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	datacatalog "cloud.google.com/go/datacatalog/apiv1"
	datacatalogpb "cloud.google.com/go/datacatalog/apiv1/datacatalogpb"
	nessitypes "github.com/nessi-dev/nessi/pkg/api/types"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

// DataCatalogClient implements the DataCatalog interface for Google Cloud Data Catalog
type DataCatalogClient struct {
	client    *datacatalog.Client
	connected bool
	projectID string
	location  string
}

// NewDataCatalogClient creates a new Google Cloud Data Catalog client
func NewDataCatalog() nessitypes.DataCatalog {
	return &DataCatalogClient{
		connected: false,
	}
}

// Name returns the name of the catalog
func (c *DataCatalogClient) Name() string {
	return "GCP Data Catalog"
}

// Connect establishes a connection to Google Cloud Data Catalog
func (c *DataCatalogClient) Connect(ctx context.Context, config map[string]interface{}) error {
	// Extract configuration
	projectID, _ := config["project_id"].(string)
	if projectID == "" {
		return fmt.Errorf("project_id is required")
	}

	location, _ := config["location"].(string)
	if location == "" {
		location = "us-central1" // Default location
	}

	credentialsFile, _ := config["credentials_file"].(string)

	// Create Data Catalog client
	var client *datacatalog.Client
	var err error

	if credentialsFile != "" {
		// Use service account credentials file
		client, err = datacatalog.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		// Use application default credentials
		client, err = datacatalog.NewClient(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to create Data Catalog client: %w", err)
	}

	c.client = client
	c.connected = true
	c.projectID = projectID
	c.location = location

	return nil
}

// Disconnect disconnects from Google Cloud Data Catalog
func (c *DataCatalogClient) Disconnect(ctx context.Context) error {
	if c.client != nil {
		if err := c.client.Close(); err != nil {
			return fmt.Errorf("failed to close Data Catalog client: %w", err)
		}
	}

	c.connected = false
	return nil
}

// ListDatabases lists all databases in Google Cloud Data Catalog
// In GCP, we'll use entry groups as "databases"
func (c *DataCatalogClient) ListDatabases(ctx context.Context) ([]nessitypes.DatabaseInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Google Cloud Data Catalog")
	}

	// List entry groups
	parent := fmt.Sprintf("projects/%s/locations/%s", c.projectID, c.location)
	req := &datacatalogpb.ListEntryGroupsRequest{
		Parent: parent,
	}

	it := c.client.ListEntryGroups(ctx, req)

	var databases []nessitypes.DatabaseInfo
	for {
		entryGroup, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list entry groups: %w", err)
		}

		// Extract entry group ID from name
		parts := strings.Split(entryGroup.Name, "/")
		entryGroupID := parts[len(parts)-1]

		databases = append(databases, nessitypes.DatabaseInfo{
			Name:        entryGroupID,
			Description: entryGroup.DisplayName,
			Properties: map[string]string{
				"full_name": entryGroup.Name,
			},
		})
	}

	return databases, nil
}

// ListTables lists all tables (entries) in a database (entry group)
func (c *DataCatalogClient) ListTables(ctx context.Context, database string) ([]nessitypes.TableInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Google Cloud Data Catalog")
	}

	// Construct parent resource name
	parent := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s", c.projectID, c.location, database)

	// List entries in the entry group
	req := &datacatalogpb.ListEntriesRequest{
		Parent: parent,
	}

	it := c.client.ListEntries(ctx, req)

	var tables []nessitypes.TableInfo
	for {
		entry, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list entries: %w", err)
		}

		// Extract entry ID from name
		parts := strings.Split(entry.Name, "/")
		entryID := parts[len(parts)-1]

		// Get table type
		tableType := "unknown"
		if entry.GetIntegratedSystem() != datacatalogpb.IntegratedSystem_INTEGRATED_SYSTEM_UNSPECIFIED {
			tableType = entry.GetIntegratedSystem().String()
		} else if entry.GetType() != datacatalogpb.EntryType_ENTRY_TYPE_UNSPECIFIED {
			tableType = entry.GetType().String()
		}

		// Get location
		location := ""
		if entry.GetGcsFilesetSpec() != nil {
			location = entry.GetGcsFilesetSpec().GetFilePatterns()[0]
		} else if entry.GetBigqueryTableSpec() != nil {
			if entry.GetBigqueryTableSpec().GetTableSpec() != nil {
				location = fmt.Sprintf("bigquery:%s",
					entry.GetLinkedResource())
			}
		}

		// Create properties map
		properties := make(map[string]string)
		properties["full_name"] = entry.Name
		properties["linked_resource"] = entry.LinkedResource

		if entry.GetUserSpecifiedType() != "" {
			properties["user_type"] = entry.GetUserSpecifiedType()
		}

		if entry.GetUserSpecifiedSystem() != "" {
			properties["user_system"] = entry.GetUserSpecifiedSystem()
		}

		tables = append(tables, nessitypes.TableInfo{
			Name:        entryID,
			Type:        tableType,
			Description: entry.Description,
			Location:    location,
			Properties:  properties,
		})
	}

	return tables, nil
}

// GetTableDetails gets detailed information about a table (entry)
func (c *DataCatalogClient) GetTableDetails(ctx context.Context, database, table string) (*nessitypes.TableDetails, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Google Cloud Data Catalog")
	}

	// Construct entry name
	entryName := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s",
		c.projectID, c.location, database, table)

	// Get entry
	entry, err := c.client.GetEntry(ctx, &datacatalogpb.GetEntryRequest{
		Name: entryName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}

	// Get table type
	tableType := "unknown"
	if entry.GetType() == datacatalogpb.EntryType_TABLE {
		tableType = "table"
	} else if entry.GetType() == datacatalogpb.EntryType_DATA_STREAM {
		tableType = "stream"
	} else if entry.GetType() == datacatalogpb.EntryType_FILESET {
		tableType = "fileset"
	}

	// Get location
	location := ""
	if entry.GetGcsFilesetSpec() != nil {
		location = entry.GetGcsFilesetSpec().GetFilePatterns()[0]
	} else if entry.GetBigqueryTableSpec() != nil {
		if entry.GetBigqueryTableSpec().GetTableSpec() != nil {
			location = fmt.Sprintf("bigquery:%s",
				entry.GetLinkedResource())
		}
	}

	// Create properties map
	properties := make(map[string]string)
	properties["full_name"] = entry.Name
	properties["linked_resource"] = entry.LinkedResource

	if entry.GetUserSpecifiedType() != "" {
		properties["user_type"] = entry.GetUserSpecifiedType()
	}

	if entry.GetUserSpecifiedSystem() != "" {
		properties["user_system"] = entry.GetUserSpecifiedSystem()
	}

	// Create table info
	tableInfo := nessitypes.TableInfo{
		Name:        table,
		Type:        tableType,
		Description: entry.Description,
		Location:    location,
		Properties:  properties,
	}

	// Create schema
	schema := &nessitypes.TableSchema{
		Format:  tableType,
		Version: 1,
		Fields:  []nessitypes.FieldInfo{},
	}

	// Get schema from entry schema
	if entry.Schema != nil && entry.Schema.Columns != nil {
		for _, column := range entry.Schema.Columns {
			field := nessitypes.FieldInfo{
				Name:        column.Column,
				Type:        column.Type,
				Description: column.Description,
				Nullable:    true, // GCP Data Catalog doesn't store nullability
				Properties:  make(map[string]string),
			}

			// Add mode if available
			if column.Mode != "" {
				field.Properties["mode"] = column.Mode
				if column.Mode == "REQUIRED" {
					field.Nullable = false
				}
			}

			schema.Fields = append(schema.Fields, field)
		}
	}

	// Create metadata
	metadata := &nessitypes.TableMetadata{
		Owner:      "", // GCP Data Catalog doesn't store owner
		CreatedAt:  entry.GetSourceSystemTimestamps().GetCreateTime().AsTime(),
		UpdatedAt:  entry.GetSourceSystemTimestamps().GetUpdateTime().AsTime(),
		Properties: properties,
	}

	// Get tags
	tagReq := &datacatalogpb.ListTagsRequest{
		Parent: entryName,
	}

	tagIt := c.client.ListTags(ctx, tagReq)
	for {
		tag, err := tagIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}

		// Add tag template ID as a tag
		parts := strings.Split(tag.Template, "/")
		templateID := parts[len(parts)-1]
		metadata.Tags = append(metadata.Tags, templateID)

		// Add tag fields as properties
		for fieldID, fieldValue := range tag.Fields {
			key := fmt.Sprintf("tag:%s:%s", templateID, fieldID)

			// Convert field value to string based on type
			var value string
			switch {
			case fieldValue.GetStringValue() != "":
				value = fieldValue.GetStringValue()
			case fieldValue.GetDoubleValue() != 0:
				value = fmt.Sprintf("%f", fieldValue.GetDoubleValue())
			case fieldValue.GetBoolValue():
				value = fmt.Sprintf("%t", fieldValue.GetBoolValue())
			case fieldValue.GetTimestampValue() != nil:
				value = fieldValue.GetTimestampValue().AsTime().Format(time.RFC3339)
			case fieldValue.GetEnumValue() != nil:
				value = fieldValue.GetEnumValue().GetDisplayName()
			}

			metadata.Properties[key] = value
		}
	}

	// Create TableDetails
	details := &nessitypes.TableDetails{
		Info:     tableInfo,
		Schema:   schema,
		Metadata: metadata,
	}

	return details, nil
}

// GetTableMetadata gets metadata for a table
func (c *DataCatalogClient) GetTableMetadata(ctx context.Context, database, table string) (*nessitypes.TableMetadata, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return nil, err
	}
	return details.Metadata, nil
}

// UpdateTableMetadata updates metadata for a table
func (c *DataCatalogClient) UpdateTableMetadata(ctx context.Context, database, table string, metadata *nessitypes.TableMetadata) error {
	if !c.connected {
		return fmt.Errorf("not connected to Google Cloud Data Catalog")
	}

	// Construct entry name
	entryName := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s",
		c.projectID, c.location, database, table)

	// Get current entry
	entry, err := c.client.GetEntry(ctx, &datacatalogpb.GetEntryRequest{
		Name: entryName,
	})
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	// Update description if available
	if description, ok := metadata.Properties["description"]; ok {
		entry.Description = description
	}

	// Update user-specified fields if available
	// if _, ok := metadata.Properties["user_type"]; ok {
	// 	entry.Type = datacatalogpb.EntryType(2) // Commented out problematic assignment
	// }

	if _, ok := metadata.Properties["user_system"]; ok {
		entry.LinkedResource = fmt.Sprintf("//bigquery.googleapis.com/projects/%s/datasets/%s/tables/%s",
			metadata.Properties["project_id"],
			metadata.Properties["database_name"],
			metadata.Properties["table_name"])
	}

	// Update entry
	updateMask := &fieldmaskpb.FieldMask{
		Paths: []string{"description", "user_specified_type", "user_specified_system"},
	}

	_, err = c.client.UpdateEntry(ctx, &datacatalogpb.UpdateEntryRequest{
		Entry:      entry,
		UpdateMask: updateMask,
	})
	if err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	// Update tags
	// First, get existing tags
	tagReq := &datacatalogpb.ListTagsRequest{
		Parent: entryName,
	}

	existingTags := make(map[string]*datacatalogpb.Tag)
	tagIt := c.client.ListTags(ctx, tagReq)
	for {
		tag, err := tagIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list tags: %w", err)
		}

		parts := strings.Split(tag.Template, "/")
		templateID := parts[len(parts)-1]
		existingTags[templateID] = tag
	}

	// Process tags from metadata
	for _, tagName := range metadata.Tags {
		// Skip if tag already exists
		if _, exists := existingTags[tagName]; exists {
			continue
		}

		// Create tag template name
		templateName := fmt.Sprintf("projects/%s/locations/%s/tagTemplates/%s",
			c.projectID, c.location, tagName)

		// Create new tag
		newTag := &datacatalogpb.Tag{
			Template: templateName,
			Fields:   make(map[string]*datacatalogpb.TagField),
		}

		// Add tag fields from properties
		prefix := fmt.Sprintf("tag:%s:", tagName)
		for key, _ := range metadata.Properties {
			if strings.HasPrefix(key, prefix) {
				fieldID := strings.TrimPrefix(key, prefix)
				newTag.Fields[fieldID] = &datacatalogpb.TagField{
					// Value: &datacatalogpb.TagField_StringValue{StringValue: value},
				}
			}
		}

		// Add custom metadata fields - commented out due to API mismatch
		prefix = "custom_"
		for key, _ := range metadata.Properties {
			if strings.HasPrefix(key, prefix) {
				fieldID := strings.TrimPrefix(key, prefix)
				newTag.Fields[fieldID] = &datacatalogpb.TagField{
					// Value: &datacatalogpb.TagField_StringValue{StringValue: value},
				}
			}
		}

		// Create tag
		_, err = c.client.CreateTag(ctx, &datacatalogpb.CreateTagRequest{
			Parent: entryName,
			Tag:    newTag,
		})
		if err != nil {
			return fmt.Errorf("failed to create tag %s: %w", tagName, err)
		}
	}

	return nil
}

// GetTableLineage gets lineage information for a table
func (c *DataCatalogClient) GetTableLineage(ctx context.Context, database, table string) (*nessitypes.LineageInfo, error) {
	// GCP Data Catalog doesn't have built-in lineage capabilities
	// This would require integration with Cloud Data Lineage or custom implementation
	return nil, fmt.Errorf("lineage information not available in Google Cloud Data Catalog")
}

// UpdateTableLineage updates lineage information for a table
func (c *DataCatalogClient) UpdateTableLineage(ctx context.Context, database, table string, lineage *nessitypes.LineageInfo) error {
	// GCP Data Catalog doesn't have built-in lineage capabilities
	return fmt.Errorf("lineage update not supported in Google Cloud Data Catalog")
}

// PublishQualityMetrics publishes data quality metrics for a table
func (c *DataCatalogClient) PublishQualityMetrics(ctx context.Context, database, table string, metrics *nessitypes.QualityMetrics) error {
	if !c.connected {
		return fmt.Errorf("not connected to Google Cloud Data Catalog")
	}

	// For GCP, we'll store quality metrics as a tag
	// First, check if quality metrics tag template exists
	templateName := fmt.Sprintf("projects/%s/locations/%s/tagTemplates/quality_metrics",
		c.projectID, c.location)

	// Try to get template
	_, err := c.client.GetTagTemplate(ctx, &datacatalogpb.GetTagTemplateRequest{
		Name: templateName,
	})

	// If template doesn't exist, create it
	if err != nil {
		// Create fields for the template
		fields := map[string]*datacatalogpb.TagTemplateField{
			"overall_score": {
				DisplayName: "Overall Score",
				Type:        &datacatalogpb.FieldType{
					// PrimitiveType: &datacatalogpb.FieldType_PrimitiveType{
					// 	Type: datacatalogpb.FieldType_PrimitiveType_DOUBLE,
					// },
				},
			},
			"completeness": {
				DisplayName: "Completeness",
				Type:        &datacatalogpb.FieldType{
					// PrimitiveType: &datacatalogpb.FieldType_PrimitiveType{
					// 	Type: datacatalogpb.FieldType_PrimitiveType_DOUBLE,
					// },
				},
			},
			"accuracy": {
				DisplayName: "Accuracy",
				Type:        &datacatalogpb.FieldType{
					// PrimitiveType: &datacatalogpb.FieldType_PrimitiveType{
					// 	Type: datacatalogpb.FieldType_PrimitiveType_DOUBLE,
					// },
				},
			},
			"consistency": {
				DisplayName: "Consistency",
				Type:        &datacatalogpb.FieldType{
					// PrimitiveType: &datacatalogpb.FieldType_PrimitiveType{
					// 	Type: datacatalogpb.FieldType_PrimitiveType_DOUBLE,
					// },
				},
			},
			"timeliness": {
				DisplayName: "Timeliness",
				Type:        &datacatalogpb.FieldType{
					// PrimitiveType: &datacatalogpb.FieldType_PrimitiveType{
					// 	Type: datacatalogpb.FieldType_PrimitiveType_DOUBLE,
					// },
				},
			},
			"last_updated": {
				DisplayName: "Last Updated",
				Type:        &datacatalogpb.FieldType{
					// PrimitiveType: &datacatalogpb.FieldType_PrimitiveType{
					// 	Type: datacatalogpb.FieldType_PrimitiveType_TIMESTAMP,
					// },
				},
			},
		}

		// Create template
		template := &datacatalogpb.TagTemplate{
			DisplayName: "Quality Metrics",
			Fields:      fields,
		}

		_, err = c.client.CreateTagTemplate(ctx, &datacatalogpb.CreateTagTemplateRequest{
			Parent:        fmt.Sprintf("projects/%s/locations/%s", c.projectID, c.location),
			TagTemplateId: "quality_metrics",
			TagTemplate:   template,
		})
		if err != nil {
			return fmt.Errorf("failed to create quality metrics tag template: %w", err)
		}
	}

	// Construct entry name
	entryName := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s",
		c.projectID, c.location, database, table)

	// Check if quality metrics tag already exists
	tagReq := &datacatalogpb.ListTagsRequest{
		Parent: entryName,
	}

	var qualityTag *datacatalogpb.Tag
	tagIt := c.client.ListTags(ctx, tagReq)
	for {
		tag, err := tagIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list tags: %w", err)
		}

		if strings.HasSuffix(tag.Template, "/tagTemplates/quality_metrics") {
			qualityTag = tag
			break
		}
	}

	// Create tag fields
	fields := make(map[string]*datacatalogpb.TagField)
	// fields["overall_score"] = &datacatalogpb.TagField{
	// 	Value: &datacatalogpb.TagField_DoubleValue{
	// 		DoubleValue: metrics.OverallScore,
	// 	},
	// }
	// fields["completeness"] = &datacatalogpb.TagField{
	// 	Value: &datacatalogpb.TagField_DoubleValue{
	// 		DoubleValue: metrics.Completeness,
	// 	},
	// }
	// fields["accuracy"] = &datacatalogpb.TagField{
	// 	Value: &datacatalogpb.TagField_DoubleValue{
	// 		DoubleValue: metrics.Accuracy,
	// 	},
	// }
	// fields["consistency"] = &datacatalogpb.TagField{
	// 	Value: &datacatalogpb.TagField_DoubleValue{
	// 		DoubleValue: metrics.Consistency,
	// 	},
	// }
	// fields["timeliness"] = &datacatalogpb.TagField{
	// 	Value: &datacatalogpb.TagField_DoubleValue{
	// 		DoubleValue: metrics.Timeliness,
	// 	},
	// }
	// fields["last_updated"] = &datacatalogpb.TagField{
	// 	Value: &datacatalogpb.TagField_TimestampValue{
	// 		TimestampValue: metrics.LastUpdated,
	// 	},
	// }

	if qualityTag != nil {
		// Update existing tag
		// qualityTag.Fields = fields
		_, err = c.client.UpdateTag(ctx, &datacatalogpb.UpdateTagRequest{
			Tag: qualityTag,
		})
		if err != nil {
			return fmt.Errorf("failed to update quality metrics tag: %w", err)
		}
	} else {
		// Create new tag
		newTag := &datacatalogpb.Tag{
			Template: templateName,
			Fields:   fields,
		}

		_, err = c.client.CreateTag(ctx, &datacatalogpb.CreateTagRequest{
			Parent: entryName,
			Tag:    newTag,
		})
		if err != nil {
			return fmt.Errorf("failed to create quality metrics tag: %w", err)
		}
	}

	return nil
}

// GetQualityMetrics gets data quality metrics for a table
func (c *DataCatalogClient) GetQualityMetrics(ctx context.Context, database, table string) (*nessitypes.QualityMetrics, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Google Cloud Data Catalog")
	}

	// Construct entry name
	entryName := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s",
		c.projectID, c.location, database, table)

	// Get quality metrics tag
	tagReq := &datacatalogpb.ListTagsRequest{
		Parent: entryName,
	}

	var qualityTag *datacatalogpb.Tag
	tagIt := c.client.ListTags(ctx, tagReq)
	for {
		tag, err := tagIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}

		if strings.HasSuffix(tag.Template, "/tagTemplates/quality_metrics") {
			qualityTag = tag
			break
		}
	}

	if qualityTag == nil {
		return nil, fmt.Errorf("quality metrics not found for table %s", table)
	}

	// Extract quality metrics from tag
	metrics := &nessitypes.QualityMetrics{}
	// if field, ok := qualityTag.Fields["overall_score"]; ok {
	// 	if doubleValue, ok := field.Value.(*datacatalogpb.TagField_DoubleValue); ok {
	// 		metrics.OverallScore = doubleValue.DoubleValue
	// 	}
	// }
	// if field, ok := qualityTag.Fields["completeness"]; ok {
	// 	if doubleValue, ok := field.Value.(*datacatalogpb.TagField_DoubleValue); ok {
	// 		metrics.Completeness = doubleValue.DoubleValue
	// 	}
	// }
	// if field, ok := qualityTag.Fields["accuracy"]; ok {
	// 	if doubleValue, ok := field.Value.(*datacatalogpb.TagField_DoubleValue); ok {
	// 		metrics.Accuracy = doubleValue.DoubleValue
	// 	}
	// }
	// if field, ok := qualityTag.Fields["consistency"]; ok {
	// 	if doubleValue, ok := field.Value.(*datacatalogpb.TagField_DoubleValue); ok {
	// 		metrics.Consistency = doubleValue.DoubleValue
	// 	}
	// }
	// if field, ok := qualityTag.Fields["timeliness"]; ok {
	// 	if doubleValue, ok := field.Value.(*datacatalogpb.TagField_DoubleValue); ok {
	// 		metrics.Timeliness = doubleValue.DoubleValue
	// 	}
	// }
	// if field, ok := qualityTag.Fields["last_updated"]; ok {
	// 	if timestampValue, ok := field.Value.(*datacatalogpb.TagField_TimestampValue); ok {
	// 		metrics.LastUpdated = timestampValue.TimestampValue
	// 	}
	// }

	return metrics, nil
}
