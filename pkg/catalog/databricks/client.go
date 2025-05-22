package databricks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/nessi-dev/nessi/pkg/security"
)

// API endpoints for Databricks Unity Catalog
const (
	catalogsEndpoint = "/api/2.1/unity-catalog/catalogs"
	schemasEndpoint  = "/api/2.1/unity-catalog/schemas"
	tablesEndpoint   = "/api/2.1/unity-catalog/tables"
)

// makeRequest makes an HTTP request to the Databricks API
func (c *DatabricksClient) makeRequest(ctx context.Context, method, endpoint string, queryParams map[string]string, body interface{}) ([]byte, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	// Add query parameters
	q := u.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s (status code: %d)", string(respBody), resp.StatusCode)
	}

	return respBody, nil
}

// GetCatalogs returns a list of catalogs in the specified workspace
func (c *DatabricksClient) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	// Check for Databricks integration license
	if err := security.RequireFeature(security.FeatureDatabricksIntegration); err != nil {
		return nil, err
	}

	respBody, err := c.makeRequest(ctx, http.MethodGet, catalogsEndpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalogs: %w", err)
	}

	var response struct {
		Catalogs []Catalog `json:"catalogs"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal catalogs response: %w", err)
	}

	return response.Catalogs, nil
}

// GetSchemas returns a list of schemas in the specified catalog
func (c *DatabricksClient) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	// Check for Databricks integration license
	if err := security.RequireFeature(security.FeatureDatabricksIntegration); err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"catalog_name": catalogName,
	}

	respBody, err := c.makeRequest(ctx, http.MethodGet, schemasEndpoint, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas: %w", err)
	}

	var response struct {
		Schemas []Schema `json:"schemas"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schemas response: %w", err)
	}

	return response.Schemas, nil
}

// GetTables returns a list of tables in the specified schema
func (c *DatabricksClient) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	// Check for Databricks integration license
	if err := security.RequireFeature(security.FeatureDatabricksIntegration); err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"catalog_name": catalogName,
		"schema_name":  schemaName,
	}

	respBody, err := c.makeRequest(ctx, http.MethodGet, tablesEndpoint, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tables: %w", err)
	}

	var response struct {
		Tables []Table `json:"tables"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tables response: %w", err)
	}

	return response.Tables, nil
}

// GetTableDetails returns details of a specific table
func (c *DatabricksClient) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	endpoint := fmt.Sprintf("%s/%s.%s.%s", tablesEndpoint, catalogName, schemaName, tableName)

	respBody, err := c.makeRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get table details: %w", err)
	}

	var response struct {
		Name        string `json:"name"`
		CatalogName string `json:"catalog_name"`
		SchemaName  string `json:"schema_name"`
		TableType   string `json:"table_type"`
		Columns     []struct {
			Name        string `json:"name"`
			TypeText    string `json:"type_text"`
			TypeName    string `json:"type_name"`
			Description string `json:"description"`
			Nullable    bool   `json:"nullable"`
		} `json:"columns"`
		Properties  map[string]string `json:"properties"`
		Description string            `json:"description"`
		Owner       string            `json:"owner"`
		Format      string            `json:"format"`
		Location    string            `json:"storage_location"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal table details response: %w", err)
	}

	// Convert Databricks table details to Nessi TableDetails
	tableDetails := &types.TableDetails{
		Info: types.TableInfo{
			Name:        response.Name,
			Type:        response.Format, // Using format as type
			Description: response.Description,
			Location:    response.Location,
			Properties:  response.Properties,
		},
		Schema: &types.TableSchema{
			Format:  response.Format,
			Version: 1, // Default version
			Fields:  make([]types.FieldInfo, 0, len(response.Columns)),
		},
		Metadata: &types.TableMetadata{
			Owner:      response.Owner,
			Properties: response.Properties,
		},
	}

	for _, col := range response.Columns {
		tableDetails.Schema.Fields = append(tableDetails.Schema.Fields, types.FieldInfo{
			Name:        col.Name,
			Type:        col.TypeName,
			Description: col.Description,
			Nullable:    col.Nullable,
		})
	}

	return tableDetails, nil
}

// UpdateTableDetails updates the details of a specific table
func (c *DatabricksClient) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	endpoint := fmt.Sprintf("%s/%s.%s.%s", tablesEndpoint, catalogName, schemaName, tableName)

	// Prepare the request body
	requestBody := map[string]interface{}{
		"description": details.Info.Description,
		"properties":  details.Info.Properties,
	}

	// If fields have descriptions, update them
	if details.Schema != nil && len(details.Schema.Fields) > 0 {
		columns := make([]map[string]interface{}, 0, len(details.Schema.Fields))
		for _, field := range details.Schema.Fields {
			columns = append(columns, map[string]interface{}{
				"name":        field.Name,
				"description": field.Description,
			})
		}
		requestBody["columns"] = columns
	}

	_, err := c.makeRequest(ctx, http.MethodPatch, endpoint, nil, requestBody)
	if err != nil {
		return fmt.Errorf("failed to update table details: %w", err)
	}

	return nil
}
