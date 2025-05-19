package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
	"github.com/golang/mock/gomock"
	nessitypes "github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination mock_glue_client.go -package aws -source glue.go GlueAPI

// TestGlueCatalog_Name tests the Name method
func TestGlueCatalog_Name(t *testing.T) {
	catalog := NewGlueCatalog()
	assert.Equal(t, "AWS Glue Data Catalog", catalog.Name())
}

// TestGlueCatalog_ListDatabases tests the ListDatabases method
func TestGlueCatalog_ListDatabases(t *testing.T) {
	// Create mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	m := NewMockGlueAPI(ctrl)

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
	m.EXPECT().GetDatabases(gomock.Any(), gomock.Any()).Return(&glue.GetDatabasesOutput{
		DatabaseList: testDatabases,
	}, nil)

	// Create catalog with mock client
	catalog := &GlueCatalog{}
	catalog.client = m
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
}

// TestGlueCatalog_ListTables tests the ListTables method
func TestGlueCatalog_ListTables(t *testing.T) {
	// Create mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockGlueAPI(ctrl)

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
	m.EXPECT().GetTables(gomock.Any(), gomock.Any()).Return(&glue.GetTablesOutput{
		TableList: testTables,
	}, nil)

	// Create catalog with mock client
	catalog := &GlueCatalog{}
	catalog.client = m
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

	// No need to verify mock expectations with gomock
}

// TestGlueCatalog_GetTableDetails tests the GetTableDetails method
func TestGlueCatalog_GetTableDetails(t *testing.T) {
	// Create mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := NewMockGlueAPI(ctrl)

	// Create test data
	createTime := time.Now().Add(-24 * time.Hour)
	updateTime := time.Now()

	expectedDetails := &nessitypes.TableDetails{
		Info: nessitypes.TableInfo{
			Name:     "table1",
			Type:     "delta",
			Location: "s3://test-bucket/table1",
		},
		Metadata: &nessitypes.TableMetadata{
			CreatedAt: createTime,
			UpdatedAt: updateTime,
		},
	}

	// Set up mock expectations
	mockClient.EXPECT().GetTableDetails(gomock.Any(), "testdb", "table1").Return(expectedDetails, nil)

	// Create catalog with mock client
	c := GlueCatalog{client: mockClient, connected: true, region: "us-west-2"}

	// Call GetTableDetails
	ctx := context.Background()
	details, err := c.GetTableDetails(ctx, "testdb", "table1")

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, details)
	assert.Equal(t, "table1", details.Info.Name)
	assert.Equal(t, "delta", details.Info.Type)
	assert.Equal(t, "s3://test-bucket/table1", details.Info.Location)

	// Verify metadata
	assert.NotNil(t, details.Metadata)
	assert.WithinDuration(t, createTime, details.Metadata.CreatedAt, time.Second)
	assert.WithinDuration(t, updateTime, details.Metadata.UpdatedAt, time.Second)

	// No need to verify mock expectations with gomock
}

// TestGlueCatalog_PublishQualityMetrics tests the PublishQualityMetrics method
func TestGlueCatalog_PublishQualityMetrics(t *testing.T) {
	// Create mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockGlueAPI(ctrl)

	// Create test data

	// Set up mock expectations
	m.EXPECT().GetTable(gomock.Any(), gomock.Any()).Return(&glue.GetTableOutput{
		Table: &types.Table{
			Name:        aws.String("table1"),
			Description: aws.String("Test Table 1"),
			Parameters: map[string]string{
				"table_type": "delta",
			},
			StorageDescriptor: &types.StorageDescriptor{
				Location: aws.String("s3://bucket/path/to/table1"),
			},
		},
	}, nil)
	m.EXPECT().UpdateTable(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, input *glue.UpdateTableInput, optFns ...func(*glue.Options)) (*glue.UpdateTableOutput, error) {
		// Verify that quality metrics are added to parameters
		params := input.TableInput.Parameters
		_, hasCompleteness := params["quality:data_completeness"]
		_, hasAccuracy := params["quality:data_accuracy"]

		if !hasCompleteness || !hasAccuracy {
			t.Error("Missing required quality metrics parameters")
		}
		return &glue.UpdateTableOutput{}, nil
	})

	// Create catalog with mock client
	catalog := &GlueCatalog{}
	catalog.client = m
	catalog.connected = true

	// Call PublishQualityMetrics
	err := catalog.PublishQualityMetrics(context.Background(), "testdb", "table1", &nessitypes.QualityMetrics{
		DataCompleteness: 0.95,
		DataAccuracy:     0.98,
	})

	// Verify results
	assert.NoError(t, err)

	// No need to verify mock expectations with gomock
}

// TestGlueCatalog_GetQualityMetrics tests the GetQualityMetrics method
func TestGlueCatalog_GetQualityMetrics(t *testing.T) {
	// Create mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockGlueAPI(ctrl)

	// Create test data

	testTable := &types.Table{
		Name:        aws.String("table1"),
		Description: aws.String("Test Table 1"),
		Parameters: map[string]string{
			"quality:data_completeness": "0.90",
			"quality:data_accuracy":     "0.85",
		},
	}

	// Set up mock expectations
	m.EXPECT().GetTable(gomock.Any(), gomock.Any()).Return(&glue.GetTableOutput{
		Table: testTable,
	}, nil)

	// Create catalog with mock client
	catalog := &GlueCatalog{}
	catalog.client = m
	catalog.connected = true

	// Call GetQualityMetrics
	ctx := context.Background()
	metrics, err := catalog.GetQualityMetrics(ctx, "testdb", "table1")

	// Verify results
	assert.NoError(t, err)
	assert.NotNil(t, metrics)

	// Verify metrics
	assert.InDelta(t, 0.90, metrics.DataCompleteness, 0.001)
	assert.InDelta(t, 0.85, metrics.DataAccuracy, 0.001)

	// No need to verify mock expectations with gomock
}
