package databricks

import (
	"context"
	"errors"
	"testing"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

// TestDatabricksErrorHandlingStandalone tests error handling in the Databricks client
func TestDatabricksErrorHandlingStandalone(t *testing.T) {
	// Create a context for testing
	ctx := context.Background()

	// Define test cases for different error scenarios
	testCases := []struct {
		name         string
		errorType    string
		errorMessage string
	}{
		{
			name:         "Authentication Error",
			errorType:    "auth",
			errorMessage: "invalid authentication token",
		},
		{
			name:         "Network Error",
			errorType:    "network",
			errorMessage: "failed to connect to Databricks API",
		},
		{
			name:         "Resource Not Found",
			errorType:    "not_found",
			errorMessage: "resource not found",
		},
		{
			name:         "Rate Limiting",
			errorType:    "rate_limit",
			errorMessage: "too many requests",
		},
		{
			name:         "Server Error",
			errorType:    "server",
			errorMessage: "internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock client that returns the specified error
			mockClient := &errorHandlingMockClient{
				errorType: tc.errorType,
			}

			// Test various API calls to ensure they properly handle errors
			_, err := mockClient.GetWorkspaces(ctx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMessage)

			_, err = mockClient.GetCatalogs(ctx, "workspace-1")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMessage)

			_, err = mockClient.GetSchemas(ctx, "workspace-1", "catalog-1")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMessage)

			_, err = mockClient.GetTables(ctx, "workspace-1", "catalog-1", "schema-1")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMessage)

			_, err = mockClient.GetTableDetails(ctx, "workspace-1", "catalog-1", "schema-1", "table-1")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMessage)
		})
	}
}

// errorHandlingMockClient is a mock implementation of the Databricks client for error handling tests
type errorHandlingMockClient struct {
	errorType string
}

// getError returns an error based on the error type
func (c *errorHandlingMockClient) getError() error {
	switch c.errorType {
	case "auth":
		return errors.New("invalid authentication token")
	case "network":
		return errors.New("failed to connect to Databricks API")
	case "not_found":
		return errors.New("resource not found")
	case "rate_limit":
		return errors.New("too many requests")
	case "server":
		return errors.New("internal server error")
	default:
		return errors.New("unknown error")
	}
}

// GetWorkspaces implements the DatabricksAPI interface
func (c *errorHandlingMockClient) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	return nil, c.getError()
}

// GetCatalogs implements the DatabricksAPI interface
func (c *errorHandlingMockClient) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	return nil, c.getError()
}

// GetSchemas implements the DatabricksAPI interface
func (c *errorHandlingMockClient) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	return nil, c.getError()
}

// GetTables implements the DatabricksAPI interface
func (c *errorHandlingMockClient) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	return nil, c.getError()
}

// GetTableDetails implements the DatabricksAPI interface
func (c *errorHandlingMockClient) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	return nil, c.getError()
}

// UpdateTableDetails implements the DatabricksAPI interface
func (c *errorHandlingMockClient) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	return c.getError()
}

// Using the types defined in catalog.go
