package delta

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	// Create temporary directories for source and target tables
	sourceDir := filepath.Join(t.TempDir(), "source")
	targetDir := filepath.Join(t.TempDir(), "target")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "_delta_log"), 0755)) // Create _delta_log for source
	require.NoError(t, os.MkdirAll(targetDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(targetDir, "_delta_log"), 0755)) // Create _delta_log for target

	// Create schema for both tables
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		{Name: "updated_at", Type: arrow.FixedWidthTypes.Timestamp_ms},
	}, nil)

	// Create source connector
	sourceConnector, err := NewConnector(sourceDir)
	require.NoError(t, err, "Failed to create source connector")
	sourceConnector.Schema = schema

	// Initial commit for source table (MetaAction)
	sourceTx := sourceConnector.log.BeginTransaction()
	
	// Create a properly formatted schema JSON string
	schemaFields := []map[string]interface{}{}
	for _, field := range schema.Fields() {
		typeStr := "string" // Default
		switch field.Type.ID() {
		case arrow.INT64:
			typeStr = "int64"
		case arrow.FLOAT64:
			typeStr = "float64"
		case arrow.TIMESTAMP:
			typeStr = "timestamp_ms"
		case arrow.STRING:
			typeStr = "string"
		case arrow.BOOL:
			typeStr = "boolean"
		}
		
		schemaFields = append(schemaFields, map[string]interface{}{
			"name": field.Name,
			"type": typeStr,
			"nullable": field.Nullable,
		})
	}
	
	schemaMap := map[string]interface{}{
		"fields": schemaFields,
	}
	
	sourceSchemaBytes, err := json.Marshal(schemaMap)
	require.NoError(t, err, "Failed to marshal source schema to JSON")
	sourceMetaAction := MetaAction{
		ID:               uuid.NewString(),
		Format:           Format{Provider: "parquet"},
		SchemaString:     string(sourceSchemaBytes),
		PartitionColumns: []string{},
		Configuration:    make(map[string]string), // Or nil, depending on how it's handled
		CreatedTime:      time.Now().UnixMilli(),
	}
	err = sourceTx.AddAction(Action{Meta: &sourceMetaAction})
	require.NoError(t, err, "Failed to add source MetaAction to transaction")
	sourceTx.CommitInfo.Operation = "CREATE TABLE"
	err = sourceConnector.CommitTransaction(context.Background(), sourceTx)
	require.NoError(t, err, "Failed to commit initial source transaction")

	// Create target connector
	targetConnector, err := NewConnector(targetDir)
	require.NoError(t, err, "Failed to create target connector")
	targetConnector.Schema = schema

	// Initial commit for target table (MetaAction)
	targetTx := targetConnector.log.BeginTransaction()
	
	// Create a properly formatted schema JSON string for target (same as source)
	targetSchemaFields := []map[string]interface{}{}
	for _, field := range schema.Fields() {
		typeStr := "string" // Default
		switch field.Type.ID() {
		case arrow.INT64:
			typeStr = "int64"
		case arrow.FLOAT64:
			typeStr = "float64"
		case arrow.TIMESTAMP:
			typeStr = "timestamp_ms"
		case arrow.STRING:
			typeStr = "string"
		case arrow.BOOL:
			typeStr = "boolean"
		}
		
		targetSchemaFields = append(targetSchemaFields, map[string]interface{}{
			"name": field.Name,
			"type": typeStr,
			"nullable": field.Nullable,
		})
	}
	
	targetSchemaMap := map[string]interface{}{
		"fields": targetSchemaFields,
	}
	
	targetSchemaBytes, err := json.Marshal(targetSchemaMap)
	require.NoError(t, err, "Failed to marshal target schema to JSON")
	targetMetaAction := MetaAction{
		ID:               uuid.NewString(),
		Format:           Format{Provider: "parquet"},
		SchemaString:     string(targetSchemaBytes),
		PartitionColumns: []string{},
		Configuration:    make(map[string]string),
		CreatedTime:      time.Now().UnixMilli(),
	}
	err = targetTx.AddAction(Action{Meta: &targetMetaAction})
	require.NoError(t, err, "Failed to add target MetaAction to transaction")
	targetTx.CommitInfo.Operation = "CREATE TABLE"
	err = targetConnector.CommitTransaction(context.Background(), targetTx)
	require.NoError(t, err, "Failed to commit initial target transaction")

	// Create source data
	sourceBuilder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer sourceBuilder.Release()

	// Add source data
	sourceBuilder.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3, 4}, nil)
	sourceBuilder.Field(1).(*array.StringBuilder).AppendValues([]string{"Alice", "Bob", "Charlie", "David"}, nil)
	sourceBuilder.Field(2).(*array.Float64Builder).AppendValues([]float64{10.5, 20.5, 30.5, 40.5}, nil)
	sourceBuilder.Field(3).(*array.TimestampBuilder).AppendValues([]arrow.Timestamp{
		arrow.Timestamp(time.Now().UnixMilli()),
		arrow.Timestamp(time.Now().UnixMilli()),
		arrow.Timestamp(time.Now().UnixMilli()),
		arrow.Timestamp(time.Now().UnixMilli()),
	}, nil)

	sourceRecord := sourceBuilder.NewRecord()
	defer sourceRecord.Release()

	// Write source data
	require.NoError(t, sourceConnector.WritePartition(context.Background(), "", sourceRecord))

	// Create target data
	targetBuilder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer targetBuilder.Release()

	// Add target data
	targetBuilder.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 5, 6}, nil)
	targetBuilder.Field(1).(*array.StringBuilder).AppendValues([]string{"Alice", "Bob", "Eve", "Frank"}, nil)
	targetBuilder.Field(2).(*array.Float64Builder).AppendValues([]float64{10.0, 20.0, 50.0, 60.0}, nil)
	targetBuilder.Field(3).(*array.TimestampBuilder).AppendValues([]arrow.Timestamp{
		arrow.Timestamp(time.Now().UnixMilli()),
		arrow.Timestamp(time.Now().UnixMilli()),
		arrow.Timestamp(time.Now().UnixMilli()),
		arrow.Timestamp(time.Now().UnixMilli()),
	}, nil)

	targetRecord := targetBuilder.NewRecord()
	defer targetRecord.Release()

	// Write target data
	require.NoError(t, targetConnector.WritePartition(context.Background(), "", targetRecord))

	// Define merge conditions and actions
	options := MergeOptions{
		SourceTable: sourceDir,
		Conditions: []MergeCondition{
			{
				LeftColumn:  "id",
				RightColumn: "id",
				Operator:    "=",
			},
		},
		WhenMatched: []MergeAction{
			{
				Type: "update",
				Values: map[string]string{
					"value": "source.value",
				},
			},
		},
		WhenNotMatched: []MergeAction{
			{
				Type: "insert",
			},
		},
	}

	// Perform merge
	stats, err := targetConnector.Merge(context.Background(), options)
	require.NoError(t, err)

	// Verify merge statistics
	assert.Equal(t, int64(4), stats.NumSourceRows)
	assert.Equal(t, int64(4), stats.NumTargetRows)
	assert.Equal(t, int64(2), stats.NumMatchedRows)    // IDs 1 and 2
	assert.Equal(t, int64(2), stats.NumNotMatchedRows) // IDs 3 and 4
	assert.Equal(t, int64(2), stats.NumUpdatedRows)    // Updated values for IDs 1 and 2
	assert.Equal(t, int64(2), stats.NumInsertedRows)   // Inserted IDs 3 and 4

	// Read merged data
	mergedRecord, err := targetConnector.ReadPartition(context.Background(), "")
	require.NoError(t, err)
	defer mergedRecord.Release()

	// Verify merged data
	assert.Equal(t, int64(6), mergedRecord.NumRows()) // 4 original + 2 inserted

	// Verify specific rows
	idCol := mergedRecord.Column(0).(*array.Int64)
	valueCol := mergedRecord.Column(2).(*array.Float64)

	// Check updated values
	for i := 0; i < idCol.Len(); i++ {
		id := idCol.Value(i)
		value := valueCol.Value(i)

		switch id {
		case 1:
			assert.Equal(t, 10.5, value) // Updated from source
		case 2:
			assert.Equal(t, 20.5, value) // Updated from source
		case 3:
			assert.Equal(t, 30.5, value) // Inserted from source
		case 4:
			assert.Equal(t, 40.5, value) // Inserted from source
		case 5:
			assert.Equal(t, 50.0, value) // Unchanged
		case 6:
			assert.Equal(t, 60.0, value) // Unchanged
		}
	}
}