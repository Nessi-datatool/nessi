package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
	"github.com/nessi-dev/nessi-dev/pkg/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGlueClient is a mock implementation of the AWS Glue client
type MockGlueClient struct {
	mock.Mock
}

func (m *MockGlueClient) GetDatabases(ctx context.Context, params *glue.GetDatabasesInput, optFns ...func(*glue.Options)) (*glue.GetDatabasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetDatabasesOutput), args.Error(1)
}

func (m *MockGlueClient) GetTables(ctx context.Context, params *glue.GetTablesInput, optFns ...func(*glue.Options)) (*glue.GetTablesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetTablesOutput), args.Error(1)
}

func (m *MockGlueClient) GetTable(ctx context.Context, params *glue.GetTableInput, optFns ...func(*glue.Options)) (*glue.GetTableOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.GetTableOutput), args.Error(1)
}

func (m *MockGlueClient) UpdateTable(ctx context.Context, params *glue.UpdateTableInput, optFns ...func(*glue.Options)) (*glue.UpdateTableOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*glue.UpdateTableOutput), args.Error(1)
}

// TestGlueCatalog_Name tests the Name method
func TestGlueCatalog_Name(t *testing.T) {
	catalog := NewGlueCatalog()
	assert.Equal(t, "AWS Glue Data Catalog", catalog.Name())
}

// TestGlueCatalog_ListDatabases tests the ListDatabases method
func TestGlueCatalog_ListDatabases(t *testing.T) {
	// Create mock client
	mockClient := new(MockGlueClient)
	
	// Create test data
	testDatabases := []types.Database{
		{
			Name:        aws.String("database1"),
			Description: aws.String("Test Database 1"),
			Parameters: map[string]string{
				"param1": "value1",
			},
		},
		{
			Name:        aws.String("database2"),
			Description: aws.String("Test Database 2"),
			Parameters: map[string]string{
				"param2": "value2",
			},
		},
	}
	
	// Set up mock expectations
	mockClient.On("GetDatabases", mock.Anything, mock.Anything).Return(&glue.GetDatabasesOutput{
		DatabaseList: testDatabases,
	}, nil)
	
	// Create catalog with mock client
	catalog := NewGlueCatalog()
	catalog.client = mockClient
	catalog.connected = true
	
	// Call ListDatabases
	ctx := context.Background()
	databases, err := catalog.ListDatabases(ctx)
	
	// Verify results
	assert.NoError(t, err)
	assert.Len(t, databases, 2)
	assert.Equal(t, "database1", databases[0].Name)
	assert.Equal(t, "Test Database 1", databases[0].Description)
	assert.Equal(t, "value1", databases[0].Properties["param1"])
	assert.Equal(t, "database2", databases[1].Name)
	assert.Equal(t, "Test Database 2", databases[1].Description)
	assert.Equal(t, "value2", databases[1].Properties["param2"])
	
	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// TestGlueCatalog_ListTables tests the ListTables method
func TestGlueCatalog_ListTables(t *testing.T) {
	// Create mock client
	mockClient := new(MockGlueClient)
	
	// Create test data
	testTables := []types.Table{
		{
			Name:        aws.String("table1"),
			Description: aws.String("Test Table 1"),
			Parameters: map[string]string{
				"table_type": "delta",
			},
			StorageDescriptor: &types.StorageDescriptor{
				Location: aws.String("s3://bucket/path/to/table1"),
			},
		},
		{
			Name:        aws.String("table2"),
			Description: aws.String("Test Table 2"),
			Parameters: map[string]string{
				"table_type": "parquet",
			},
			StorageDescriptor: &types.StorageDescriptor{
				Location: aws.String("s3://bucket/path/to/table2"),
			},
		},
	}
	
	// Set up mock expectations
	mockClient.On("GetTables", mock.Anything, mock.Anything).Return(&glue.GetTablesOutput{
		TableList: testTables,
	}, nil)
	
	// Create catalog with mock client
	catalog := NewGlueCatalog()
	catalog.client = mockClient
	catalog.connected = true
	
	// Call ListTables
	ctx := context.Background()
	tables, err := catalog.ListTables(ctx, "testdb")
	
	// Verify results
	assert.NoError(t, err)
	assert.Len(t, tables, 2)
	assert.Equal(t, "table1", tables[0].Name)
	assert.Equal(t, "delta", tables[0].Type)
	assert.Equal(t, "Test Table 1", tables[0].Description)
	assert.Equal(t, "s3://bucket/path/to/table1", tables[0].Location)
	assert.Equal(t, "table2", tables[1].Name)
	assert.Equal(t, "parquet", tables[1].Type)
	assert.Equal(t, "Test Table 2", tables[1].Description)
	assert.Equal(t, "s3://bucket/path/to/table2", tables[1].Location)
	
	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// TestGlueCatalog_GetTableDetails tests the GetTableDetails method
func TestGlueCatalog_GetTableDetails(t *testing.T) {
	// Create mock client
	mockClient := new(MockGlueClient)
	
	// Create test data
	createTime := time.Now().Add(-24 * time.Hour)
	updateTime := time.Now()
	
	testTable := &types.Table{
		Name:        aws.String("table1"),
		Description: aws.String("Test Table 1"),
		Owner:       aws.String("testuser"),
		CreateTime:  &createTime,
		UpdateTime:  &updateTime,
		Parameters: map[string]string{
			"table_type": "delta",
			"param1":     "value1",
		},
		StorageDescriptor: &types.StorageDescriptor{
			Location: aws.String("s3://bucket/path/to/table1"),
			Columns: []types.Column{
				{
					Name:    aws.String("id"),
					Type:    aws.String("int"),
					Comment: aws.String("ID column"),
					Parameters: map[string]string{
						"primary_key": "true",
					},
				},
				{
					Name:    aws.String("name"),
					Type:    aws.String("string"),
					Comment: aws.String("Name column"),
				},
			},
		},
	}
	
	// Set up mock expectations
	mockClient.On("GetTable", mock.Anything, mock.Anything).Return(&glue.GetTableOutput{
		Table: testTable,
	}, nil)
	
	// Create catalog with mock client
	catalog := NewGlueCatalog()
	catalog.client = mockClient
	catalog.connected = true
	
	// Call GetTableDetails
	ctx := context.Background()
	details, err := catalog.GetTableDetails(ctx, "testdb", "table1")
	
	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, "table1", details.Info.Name)
	assert.Equal(t, "delta", details.Info.Type)
	assert.Equal(t, "Test Table 1", details.Info.Description)
	assert.Equal(t, "s3://bucket/path/to/table1", details.Info.Location)
	assert.Equal(t, "value1", details.Info.Properties["param1"])
	
	// Verify schema
	assert.NotNil(t, details.Schema)
	assert.Equal(t, "delta", details.Schema.Format)
	assert.Len(t, details.Schema.Fields, 2)
	assert.Equal(t, "id", details.Schema.Fields[0].Name)
	assert.Equal(t, "int", details.Schema.Fields[0].Type)
	assert.Equal(t, "ID column", details.Schema.Fields[0].Description)
	assert.Equal(t, "true", details.Schema.Fields[0].Properties["primary_key"])
	assert.Equal(t, "name", details.Schema.Fields[1].Name)
	assert.Equal(t, "string", details.Schema.Fields[1].Type)
	assert.Equal(t, "Name column", details.Schema.Fields[1].Description)
	
	// Verify metadata
	assert.NotNil(t, details.Metadata)
	assert.Equal(t, "testuser", details.Metadata.Owner)
	assert.Equal(t, createTime, details.Metadata.CreatedAt)
	assert.Equal(t, updateTime, details.Metadata.UpdatedAt)
	assert.Equal(t, "value1", details.Metadata.Properties["param1"])
	
	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// TestGlueCatalog_PublishQualityMetrics tests the PublishQualityMetrics method
func TestGlueCatalog_PublishQualityMetrics(t *testing.T) {
	// Create mock client
	mockClient := new(MockGlueClient)
	
	// Create test data
	testTable := &types.Table{
		Name:        aws.String("table1"),
		Description: aws.String("Test Table 1"),
		Parameters: map[string]string{
			"table_type": "delta",
		},
		StorageDescriptor: &types.StorageDescriptor{
			Location: aws.String("s3://bucket/path/to/table1"),
		},
	}
	
	// Set up mock expectations
	mockClient.On("GetTable", mock.Anything, mock.Anything).Return(&glue.GetTableOutput{
		Table: testTable,
	}, nil)
	mockClient.On("UpdateTable", mock.Anything, mock.MatchedBy(func(input *glue.UpdateTableInput) bool {
		// Verify that quality metrics are added to parameters
		params := input.TableInput.Parameters
		_, hasOverallScore := params["quality:overall_score"]
		_, hasCompleteness := params["quality:completeness"]
		_, hasAccuracy := params["quality:accuracy"]
		_, hasConsistency := params["quality:consistency"]
		_, hasTimeliness := params["quality:timeliness"]
		_, hasLastUpdated := params["quality:last_updated"]
		
		return hasOverallScore && hasCompleteness && hasAccuracy && 
			hasConsistency && hasTimeliness && hasLastUpdated
	})).Return(&glue.UpdateTableOutput{}, nil)
	
	// Create catalog with mock client
	catalog := NewGlueCatalog()
	catalog.client = mockClient
	catalog.connected = true
	
	// Create quality metrics
	metrics := &catalog.QualityMetrics{
		OverallScore: 0.85,
		Completeness: 0.90,
		Accuracy:     0.85,
		Consistency:  0.80,
		Timeliness:   0.75,
		LastUpdated:  time.Now(),
		RuleResults: []catalog.RuleResult{
			{
				RuleName: "test_rule",
				RuleType: "range",
				Passed:   true,
				Score:    1.0,
				Details:  "Rule passed",
			},
		},
	}
	
	// Call PublishQualityMetrics
	ctx := context.Background()
	err := catalog.PublishQualityMetrics(ctx, "testdb", "table1", metrics)
	
	// Verify results
	assert.NoError(t, err)
	
	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// TestGlueCatalog_GetQualityMetrics tests the GetQualityMetrics method
func TestGlueCatalog_GetQualityMetrics(t *testing.T) {
	// Create mock client
	mockClient := new(MockGlueClient)
	
	// Create test data
	now := time.Now()
	nowStr := now.Format(time.RFC3339)
	
	testTable := &types.Table{
		Name:        aws.String("table1"),
		Description: aws.String("Test Table 1"),
		Parameters: map[string]string{
			"table_type":           "delta",
			"quality:overall_score": "0.85",
			"quality:completeness":  "0.90",
			"quality:accuracy":      "0.85",
			"quality:consistency":   "0.80",
			"quality:timeliness":    "0.75",
			"quality:last_updated":  nowStr,
		},
	}
	
	// Set up mock expectations
	mockClient.On("GetTable", mock.Anything, mock.Anything).Return(&glue.GetTableOutput{
		Table: testTable,
	}, nil)
	
	// Create catalog with mock client
	catalog := NewGlueCatalog()
	catalog.client = mockClient
	catalog.connected = true
	
	// Call GetQualityMetrics
	ctx := context.Background()
	metrics, err := catalog.GetQualityMetrics(ctx, "testdb", "table1")
	
	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.InDelta(t, 0.85, metrics.OverallScore, 0.001)
	assert.InDelta(t, 0.90, metrics.Completeness, 0.001)
	assert.InDelta(t, 0.85, metrics.Accuracy, 0.001)
	assert.InDelta(t, 0.80, metrics.Consistency, 0.001)
	assert.InDelta(t, 0.75, metrics.Timeliness, 0.001)
	
	// Verify last updated time (just check the date part since time parsing might have small differences)
	assert.Equal(t, now.Format("2006-01-02"), metrics.LastUpdated.Format("2006-01-02"))
	
	// Verify mock expectations
	mockClient.AssertExpectations(t)
}
