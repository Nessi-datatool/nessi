package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	nessitypes "github.com/nessi-dev/nessi/pkg/api/types"
)

// PurviewClient is a simple client for Azure Purview API
type PurviewClient struct {
	endpoint   string
	credential azcore.TokenCredential
	httpClient *http.Client
}

// NewPurviewClient creates a new Azure Purview client
func NewPurviewClient(endpoint string, credential azcore.TokenCredential) *PurviewClient {
	return &PurviewClient{
		endpoint:   endpoint,
		credential: credential,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// sendRequest sends a request to the Purview API
func (c *PurviewClient) sendRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.endpoint, path)
	
	var req *http.Request
	var err error
	
	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(bodyJSON)))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}
	
	// Get token from credential
	token, err := c.credential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://purview.azure.net/.default"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	
	// Add token to request
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	
	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	return resp, nil
}

// PurviewCatalog implements the DataCatalog interface for Azure Purview
type PurviewCatalog struct {
	client      *PurviewClient
	connected   bool
	accountName string
}

// NewPurviewCatalog creates a new Azure Purview Data Catalog client
func NewPurviewCatalog() nessitypes.DataCatalog {
	return &PurviewCatalog{
		connected: false,
	}
}

// Name returns the name of the catalog
func (c *PurviewCatalog) Name() string {
	return "Azure Purview Data Catalog"
}

// Connect establishes a connection to Azure Purview
func (c *PurviewCatalog) Connect(ctx context.Context, config map[string]interface{}) error {
	// Extract configuration
	accountName, _ := config["account_name"].(string)
	if accountName == "" {
		return fmt.Errorf("account_name is required")
	}
	
	// Create credential
	var credential azcore.TokenCredential
	var err error
	
	// Use Azure AD authentication
	credential, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure AD credential: %w", err)
	}
	
	// Create Purview client
	endpoint := fmt.Sprintf("https://%s.purview.azure.com", accountName)
	c.client = NewPurviewClient(endpoint, credential)
	c.connected = true
	c.accountName = accountName
	
	return nil
}

// Disconnect disconnects from Azure Purview
func (c *PurviewCatalog) Disconnect(ctx context.Context) error {
	c.connected = false
	return nil
}

// ListDatabases lists all databases (collections) in Azure Purview
func (c *PurviewCatalog) ListDatabases(ctx context.Context) ([]nessitypes.DatabaseInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Azure Purview")
	}
	
	// Purview uses collections instead of databases
	// Call Purview API to list collections
	resp, err := c.client.sendRequest(ctx, "GET", "/catalog/api/atlas/v2/glossary", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list collections: %s", resp.Status)
	}
	
	// Parse response
	var result struct {
		GlossaryInfo []struct {
			GUID        string `json:"guid"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"glossaryInfo"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Convert to DatabaseInfo
	databases := make([]nessitypes.DatabaseInfo, 0, len(result.GlossaryInfo))
	for _, glossary := range result.GlossaryInfo {
		databases = append(databases, nessitypes.DatabaseInfo{
			Name:        glossary.Name,
			Description: glossary.Description,
			Properties: map[string]string{
				"guid": glossary.GUID,
			},
		})
	}
	
	return databases, nil
}

// ListTables lists all tables (entities) in a database (collection)
func (c *PurviewCatalog) ListTables(ctx context.Context, database string) ([]nessitypes.TableInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Azure Purview")
	}
	
	// In Purview, we need to search for entities of type "Table"
	requestBody := map[string]interface{}{
		"keywords": "",
		"limit":    1000,
		"offset":   0,
		"filter": map[string]interface{}{
			"typeName": "Table",
		},
	}
	
	// Call Purview API to search for tables
	resp, err := c.client.sendRequest(ctx, "POST", "/catalog/api/search/query", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list tables: %s", resp.Status)
	}
	
	// Parse response
	var result struct {
		Value []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Properties  map[string]interface{} `json:"properties"`
		} `json:"value"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Convert to TableInfo
	tables := make([]nessitypes.TableInfo, 0, len(result.Value))
	for _, entity := range result.Value {
		tableType := "unknown"
		location := ""
		
		if entity.Properties != nil {
			if format, ok := entity.Properties["format"]; ok {
				tableType = fmt.Sprintf("%v", format)
			}
			if loc, ok := entity.Properties["location"]; ok {
				location = fmt.Sprintf("%v", loc)
			}
		}
		
		properties := make(map[string]string)
		for k, v := range entity.Properties {
			properties[k] = fmt.Sprintf("%v", v)
		}
		
		tables = append(tables, nessitypes.TableInfo{
			Name:        entity.Name,
			Type:        tableType,
			Description: entity.Description,
			Location:    location,
			Properties:  properties,
		})
	}
	
	return tables, nil
}

// GetTableDetails gets detailed information about a table
func (c *PurviewCatalog) GetTableDetails(ctx context.Context, database, table string) (*nessitypes.TableDetails, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Azure Purview")
	}
	
	// In Purview, we need to search for the specific table by name
	requestBody := map[string]interface{}{
		"keywords": table,
		"limit":    1,
		"filter": map[string]interface{}{
			"typeName": "Table",
		},
	}
	
	// Call Purview API to search for the table
	resp, err := c.client.sendRequest(ctx, "POST", "/catalog/api/search/query", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get table details: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get table details: %s", resp.Status)
	}
	
	// Parse response
	var searchResult struct {
		Value []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Properties  map[string]interface{} `json:"properties"`
		} `json:"value"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(searchResult.Value) == 0 {
		return nil, fmt.Errorf("table %s not found", table)
	}
	
	entity := searchResult.Value[0]
	
	// Get table info
	tableType := "unknown"
	location := ""
	
	if entity.Properties != nil {
		if format, ok := entity.Properties["format"]; ok {
			tableType = fmt.Sprintf("%v", format)
		}
		if loc, ok := entity.Properties["location"]; ok {
			location = fmt.Sprintf("%v", loc)
		}
	}
	
	properties := make(map[string]string)
	for k, v := range entity.Properties {
		properties[k] = fmt.Sprintf("%v", v)
	}
	
	tableInfo := nessitypes.TableInfo{
		Name:        entity.Name,
		Type:        tableType,
		Description: entity.Description,
		Location:    location,
		Properties:  properties,
	}
	
	// Get entity details
	resp, err = c.client.sendRequest(ctx, "GET", fmt.Sprintf("/catalog/api/atlas/v2/entity/guid/%s", entity.ID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity details: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get entity details: %s", resp.Status)
	}
	
	// Parse entity details
	var entityResult struct {
		Entity struct {
			Attributes map[string]interface{} `json:"attributes"`
			CreateTime int64                  `json:"createTime"`
			UpdateTime int64                  `json:"updateTime"`
			Owner      string                 `json:"owner"`
		} `json:"entity"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&entityResult); err != nil {
		return nil, fmt.Errorf("failed to decode entity details: %w", err)
	}
	
	// Get schema
	schema := &nessitypes.TableSchema{
		Format:  tableType,
		Version: 1,
		Fields:  []nessitypes.FieldInfo{},
	}
	
	// Get columns
	if columns, ok := entityResult.Entity.Attributes["columns"]; ok {
		if columnsList, ok := columns.([]interface{}); ok {
			for _, col := range columnsList {
				if colMap, ok := col.(map[string]interface{}); ok {
					field := nessitypes.FieldInfo{
						Name:        fmt.Sprintf("%v", colMap["name"]),
						Type:        fmt.Sprintf("%v", colMap["dataType"]),
						Description: fmt.Sprintf("%v", colMap["description"]),
						Nullable:    true,
						Properties:  make(map[string]string),
					}
					
					for k, v := range colMap {
						if k != "name" && k != "dataType" && k != "description" {
							field.Properties[k] = fmt.Sprintf("%v", v)
						}
					}
					
					schema.Fields = append(schema.Fields, field)
				}
			}
		}
	}
	
	// Get metadata
	metadata := &nessitypes.TableMetadata{
		Owner:      entityResult.Entity.Owner,
		CreatedAt:  time.Unix(entityResult.Entity.CreateTime/1000, 0),
		UpdatedAt:  time.Unix(entityResult.Entity.UpdateTime/1000, 0),
		Properties: properties,
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
func (c *PurviewCatalog) GetTableMetadata(ctx context.Context, database, table string) (*nessitypes.TableMetadata, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return nil, err
	}
	return details.Metadata, nil
}

// UpdateTableMetadata updates metadata for a table
func (c *PurviewCatalog) UpdateTableMetadata(ctx context.Context, database, table string, metadata *nessitypes.TableMetadata) error {
	if !c.connected {
		return fmt.Errorf("not connected to Azure Purview")
	}
	
	// First, get the table to get its GUID
	requestBody := map[string]interface{}{
		"keywords": table,
		"limit":    1,
		"filter": map[string]interface{}{
			"typeName": "Table",
		},
	}
	
	resp, err := c.client.sendRequest(ctx, "POST", "/catalog/api/search/query", requestBody)
	if err != nil {
		return fmt.Errorf("failed to find table: %w", err)
	}
	defer resp.Body.Close()
	
	var searchResult struct {
		Value []struct {
			ID string `json:"id"`
		} `json:"value"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(searchResult.Value) == 0 {
		return fmt.Errorf("table %s not found", table)
	}
	
	guid := searchResult.Value[0].ID
	
	// Update entity
	updateBody := map[string]interface{}{
		"entity": map[string]interface{}{
			"guid": guid,
			"attributes": map[string]interface{}{
				"description": metadata.Properties["description"],
				"owner":       metadata.Owner,
			},
		},
	}
	
	// Add custom attributes
	for k, v := range metadata.Properties {
		if k != "description" {
			updateBody["entity"].(map[string]interface{})["attributes"].(map[string]interface{})[k] = v
		}
	}
	
	// Update entity
	resp, err = c.client.sendRequest(ctx, "PUT", "/catalog/api/atlas/v2/entity", updateBody)
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update entity: %s", resp.Status)
	}
	
	// Add tags
	if len(metadata.Tags) > 0 {
		tagsBody := map[string]interface{}{
			"classification": map[string]interface{}{
				"entityGuid": guid,
				"typeName":   "tags",
			},
		}
		
		for _, tag := range metadata.Tags {
			resp, err = c.client.sendRequest(ctx, "POST", "/catalog/api/atlas/v2/entity/guid/"+guid+"/classifications", tagsBody)
			if err != nil {
				return fmt.Errorf("failed to add tag %s: %w", tag, err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				return fmt.Errorf("failed to add tag %s: %s", tag, resp.Status)
			}
		}
	}
	
	return nil
}

// GetTableLineage gets lineage information for a table
func (c *PurviewCatalog) GetTableLineage(ctx context.Context, database, table string) (*nessitypes.LineageInfo, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Azure Purview")
	}
	
	// First, get the table to get its GUID
	requestBody := map[string]interface{}{
		"keywords": table,
		"limit":    1,
		"filter": map[string]interface{}{
			"typeName": "Table",
		},
	}
	
	resp, err := c.client.sendRequest(ctx, "POST", "/catalog/api/search/query", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to find table: %w", err)
	}
	defer resp.Body.Close()
	
	var searchResult struct {
		Value []struct {
			ID string `json:"id"`
		} `json:"value"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(searchResult.Value) == 0 {
		return nil, fmt.Errorf("table %s not found", table)
	}
	
	guid := searchResult.Value[0].ID
	
	// Get lineage
	resp, err = c.client.sendRequest(ctx, "GET", fmt.Sprintf("/catalog/api/atlas/v2/lineage/%s", guid), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineage: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get lineage: %s", resp.Status)
	}
	
	// Parse lineage
	var lineageResult struct {
		BaseEntityGuid string `json:"baseEntityGuid"`
		LineageDirection string `json:"lineageDirection"`
		LineageDepth int `json:"lineageDepth"`
		GuidEntityMap map[string]struct {
			TypeName string `json:"typeName"`
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"guidEntityMap"`
		Relations []struct {
			FromEntityId string `json:"fromEntityId"`
			ToEntityId string `json:"toEntityId"`
		} `json:"relations"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&lineageResult); err != nil {
		return nil, fmt.Errorf("failed to decode lineage: %w", err)
	}
	
	// Create LineageInfo
	lineageInfo := &nessitypes.LineageInfo{
		Upstream:   []nessitypes.TableReference{},
		Downstream: []nessitypes.TableReference{},
		Properties: map[string]string{
			"provider": "Azure Purview",
			"guid":    guid,
		},
	}
	
	// Process relations
	for _, relation := range lineageResult.Relations {
		// Upstream: relation.ToEntityId == guid
		if relation.ToEntityId == guid {
			fromEntity, ok := lineageResult.GuidEntityMap[relation.FromEntityId]
			if ok && fromEntity.TypeName == "Table" {
				ref := nessitypes.TableReference{
					Table:    fmt.Sprintf("%v", fromEntity.Attributes["name"]),
					Database: fmt.Sprintf("%v", fromEntity.Attributes["qualifiedName"]),
					Catalog:  "Azure Purview",
				}
				lineageInfo.Upstream = append(lineageInfo.Upstream, ref)
			}
		}
		
		// Downstream: relation.FromEntityId == guid
		if relation.FromEntityId == guid {
			toEntity, ok := lineageResult.GuidEntityMap[relation.ToEntityId]
			if ok && toEntity.TypeName == "Table" {
				ref := nessitypes.TableReference{
					Table:    fmt.Sprintf("%v", toEntity.Attributes["name"]),
					Database: fmt.Sprintf("%v", toEntity.Attributes["qualifiedName"]),
					Catalog:  "Azure Purview",
				}
				lineageInfo.Downstream = append(lineageInfo.Downstream, ref)
			}
		}
	}
	
	return lineageInfo, nil
}

// UpdateTableLineage updates lineage information for a table
func (c *PurviewCatalog) UpdateTableLineage(ctx context.Context, database, table string, lineage *nessitypes.LineageInfo) error {
	// Azure Purview doesn't support direct lineage updates through the API
	// Lineage is typically created through scanning or process registration
	return fmt.Errorf("direct lineage update not supported in Azure Purview")
}

// PublishQualityMetrics publishes data quality metrics for a table
func (c *PurviewCatalog) PublishQualityMetrics(ctx context.Context, database, table string, metrics *nessitypes.QualityMetrics) error {
	if !c.connected {
		return fmt.Errorf("not connected to Azure Purview")
	}
	
	// First, get the table to get its GUID
	requestBody := map[string]interface{}{
		"keywords": table,
		"limit":    1,
		"filter": map[string]interface{}{
			"typeName": "Table",
		},
	}
	
	resp, err := c.client.sendRequest(ctx, "POST", "/catalog/api/search/query", requestBody)
	if err != nil {
		return fmt.Errorf("failed to find table: %w", err)
	}
	defer resp.Body.Close()
	
	var searchResult struct {
		Value []struct {
			ID string `json:"id"`
		} `json:"value"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(searchResult.Value) == 0 {
		return fmt.Errorf("table %s not found", table)
	}
	
	guid := searchResult.Value[0].ID

	// Update entity with quality metrics
	updateBody := map[string]interface{}{
		"entity": map[string]interface{}{
			"guid": guid,
			"attributes": map[string]interface{}{
				"quality:total_rows":        fmt.Sprintf("%d", metrics.TotalRows),
				"quality:null_rows":         fmt.Sprintf("%d", metrics.NullRows),
				"quality:duplicate_rows":    fmt.Sprintf("%d", metrics.DuplicateRows),
				"quality:invalid_rows":      fmt.Sprintf("%d", metrics.InvalidRows),
				"quality:data_completeness": fmt.Sprintf("%.2f", metrics.DataCompleteness),
				"quality:data_accuracy":     fmt.Sprintf("%.2f", metrics.DataAccuracy),
				"quality:data_consistency":  fmt.Sprintf("%.2f", metrics.DataConsistency),
				"quality:schema_version":    metrics.SchemaVersion,
				"quality:last_updated":      metrics.LastUpdated.Format(time.RFC3339),
			},
		},
	}

	// Update entity
	resp, err = c.client.sendRequest(ctx, "PUT", "/catalog/api/atlas/v2/entity", updateBody)
	if err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update entity: %s", resp.Status)
	}
	
	return nil
}

// GetQualityMetrics gets data quality metrics for a table
func (c *PurviewCatalog) GetQualityMetrics(ctx context.Context, database, table string) (*nessitypes.QualityMetrics, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to Azure Purview")
	}
	
	// Get table details
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return nil, err
	}
	
	// Extract quality metrics from properties
	metrics := &nessitypes.QualityMetrics{}

	if totalRows, ok := details.Metadata.Properties["quality:total_rows"]; ok {
		fmt.Sscanf(totalRows, "%d", &metrics.TotalRows)
	}

	if nullRows, ok := details.Metadata.Properties["quality:null_rows"]; ok {
		fmt.Sscanf(nullRows, "%d", &metrics.NullRows)
	}

	if duplicateRows, ok := details.Metadata.Properties["quality:duplicate_rows"]; ok {
		fmt.Sscanf(duplicateRows, "%d", &metrics.DuplicateRows)
	}

	if invalidRows, ok := details.Metadata.Properties["quality:invalid_rows"]; ok {
		fmt.Sscanf(invalidRows, "%d", &metrics.InvalidRows)
	}

	if dataCompleteness, ok := details.Metadata.Properties["quality:data_completeness"]; ok {
		fmt.Sscanf(dataCompleteness, "%f", &metrics.DataCompleteness)
	}

	if dataAccuracy, ok := details.Metadata.Properties["quality:data_accuracy"]; ok {
		fmt.Sscanf(dataAccuracy, "%f", &metrics.DataAccuracy)
	}

	if dataConsistency, ok := details.Metadata.Properties["quality:data_consistency"]; ok {
		fmt.Sscanf(dataConsistency, "%f", &metrics.DataConsistency)
	}

	if schemaVersion, ok := details.Metadata.Properties["quality:schema_version"]; ok {
		metrics.SchemaVersion = schemaVersion
	}

	if lastUpdated, ok := details.Metadata.Properties["quality:last_updated"]; ok {
		metrics.LastUpdated, _ = time.Parse(time.RFC3339, lastUpdated)
	}
	
	return metrics, nil
}
