package databricks

import (
	"context"
	"net/http"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
)

// DatabricksAPI defines the interface for the Databricks catalog client.
// This interface is used to abstract the Databricks API client for testing purposes.
type DatabricksAPI interface {
	GetWorkspaces(ctx context.Context) ([]Workspace, error)
	GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error)
	GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error)
	GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error)
	GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error)
	UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error
}

// Workspace represents a Databricks workspace
type Workspace struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Region      string `json:"region"`
	Description string `json:"description,omitempty"`
}

// Catalog represents a Databricks Unity Catalog
type Catalog struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Type        string `json:"type,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// Schema represents a schema (database) in Databricks
type Schema struct {
	Name        string `json:"name"`
	CatalogName string `json:"catalog_name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// Table represents a table in Databricks
type Table struct {
	Name        string `json:"name"`
	CatalogName string `json:"catalog_name"`
	SchemaName  string `json:"schema_name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Format      string `json:"format,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// Column represents a column in a Databricks table
type Column struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Nullable    bool   `json:"nullable"`
}

// DatabricksClient implements the DatabricksAPI interface
type DatabricksClient struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	workspaceID string
}

// NewDatabricksClient creates a new Databricks client
func NewDatabricksClient(baseURL, token, workspaceID string) *DatabricksClient {
	return &DatabricksClient{
		baseURL:     baseURL,
		token:       token,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		workspaceID: workspaceID,
	}
}

// GetWorkspaces returns a list of available Databricks workspaces
func (c *DatabricksClient) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	// In a real implementation, this would make an API call to Databricks Account API
	// For now, we'll return the configured workspace
	return []Workspace{
		{
			ID:     c.workspaceID,
			Name:   "Default Workspace",
			URL:    c.baseURL,
			Region: "us-west-2", // Default region
		},
	}, nil
}

// Implementation of DatabricksAPI interface is in client.go
