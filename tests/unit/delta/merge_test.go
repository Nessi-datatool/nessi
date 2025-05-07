package delta_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi-dev/internal/delta"
)

func TestMergeOperation(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Initialize Delta connector
	connector := delta.NewConnector(tempDir, "test_table")
	require.NoError(t, connector.Validate())

	// Create test schemas
	sourceSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
			{Name: "active", Type: arrow.FixedWidthTypes.Boolean},
			{Name: "updated_at", Type: arrow.FixedWidthTypes.Timestamp_ms},
		},
		nil,
	)

	targetSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
			{Name: "active", Type: arrow.FixedWidthTypes.Boolean},
			{Name: "updated_at", Type: arrow.FixedWidthTypes.Timestamp_ms},
		},
		nil,
	)

	// Create source data
	sourceBuilder := array.NewRecordBuilder(memory.DefaultAllocator, sourceSchema)
	defer sourceBuilder.Release()

	// Add source data
	sourceBuilder.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3}, nil)
	sourceBuilder.Field(1).(*array.StringBuilder).AppendValues([]string{"Alice", "Bob", "Charlie"}, nil)
	sourceBuilder.Field(2).(*array.Float64Builder).AppendValues([]float64{10.5, 20.5, 30.5}, nil)
	sourceBuilder.Field(3).(*array.BooleanBuilder).AppendValues([]bool{true, false, true}, nil)
	sourceBuilder.Field(4).(*array.TimestampBuilder).AppendValues(
		[]arrow.Timestamp{
			arrow.Timestamp(time.Now().UnixNano() / 1000000),
			arrow.Timestamp(time.Now().UnixNano() / 1000000),
			arrow.Timestamp(time.Now().UnixNano() / 1000000),
		},
		nil,
	)

	sourceRecord := sourceBuilder.NewRecord()
	defer sourceRecord.Release()

	// Create target data
	targetBuilder := array.NewRecordBuilder(memory.DefaultAllocator, targetSchema)
	defer targetBuilder.Release()

	// Add target data
	targetBuilder.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 4}, nil)
	targetBuilder.Field(1).(*array.StringBuilder).AppendValues([]string{"Alice Old", "Bob Old", "David"}, nil)
	targetBuilder.Field(2).(*array.Float64Builder).AppendValues([]float64{10.0, 20.0, 40.0}, nil)
	targetBuilder.Field(3).(*array.BooleanBuilder).AppendValues([]bool{false, false, true}, nil)
	targetBuilder.Field(4).(*array.TimestampBuilder).AppendValues(
		[]arrow.Timestamp{
			arrow.Timestamp(time.Now().UnixNano() / 1000000),
			arrow.Timestamp(time.Now().UnixNano() / 1000000),
			arrow.Timestamp(time.Now().UnixNano() / 1000000),
		},
		nil,
	)

	targetRecord := targetBuilder.NewRecord()
	defer targetRecord.Release()

	// Write initial data
	require.NoError(t, connector.WriteRecord(context.Background(), targetRecord))

	// Define merge options
	options := delta.MergeOptions{
		SourceTable: "source_table",
		Conditions: []delta.MergeCondition{
			{
				LeftColumn:  "id",
				RightColumn: "id",
				Operator:    "=",
			},
		},
		WhenMatched: []delta.MergeAction{
			{
				Type: "update",
				Values: map[string]string{
					"name":   "source.name",
					"value":  "source.value",
					"active": "source.active",
				},
			},
		},
		WhenNotMatched: []delta.MergeAction{
			{
				Type: "insert",
				Values: map[string]string{
					"id":        "source.id",
					"name":      "source.name",
					"value":     "source.value",
					"active":    "source.active",
					"updated_at": "source.updated_at",
				},
			},
		},
	}

	// Perform merge operation
	stats, err := connector.Merge(context.Background(), options)
	require.NoError(t, err)

	// Verify merge statistics
	assert.Equal(t, int64(3), stats.NumSourceRows)
	assert.Equal(t, int64(3), stats.NumTargetRows)
	assert.Equal(t, int64(2), stats.NumMatchedRows) // IDs 1 and 2
	assert.Equal(t, int64(1), stats.NumNotMatchedRows) // ID 3
	assert.Equal(t, int64(2), stats.NumUpdatedRows) // IDs 1 and 2
	assert.Equal(t, int64(1), stats.NumInsertedRows) // ID 3

	// Read and verify merged data
	mergedRecord, err := connector.ReadRecord(context.Background())
	require.NoError(t, err)
	defer mergedRecord.Release()

	// Verify ID 1 (matched and updated)
	id1Idx := findRowIndex(mergedRecord, "id", int64(1))
	require.NotEqual(t, -1, id1Idx)
	assert.Equal(t, "Alice", getStringValue(mergedRecord, "name", id1Idx))
	assert.Equal(t, 10.5, getFloat64Value(mergedRecord, "value", id1Idx))
	assert.Equal(t, true, getBooleanValue(mergedRecord, "active", id1Idx))

	// Verify ID 2 (matched and updated)
	id2Idx := findRowIndex(mergedRecord, "id", int64(2))
	require.NotEqual(t, -1, id2Idx)
	assert.Equal(t, "Bob", getStringValue(mergedRecord, "name", id2Idx))
	assert.Equal(t, 20.5, getFloat64Value(mergedRecord, "value", id2Idx))
	assert.Equal(t, false, getBooleanValue(mergedRecord, "active", id2Idx))

	// Verify ID 3 (inserted)
	id3Idx := findRowIndex(mergedRecord, "id", int64(3))
	require.NotEqual(t, -1, id3Idx)
	assert.Equal(t, "Charlie", getStringValue(mergedRecord, "name", id3Idx))
	assert.Equal(t, 30.5, getFloat64Value(mergedRecord, "value", id3Idx))
	assert.Equal(t, true, getBooleanValue(mergedRecord, "active", id3Idx))

	// Verify ID 4 (unchanged)
	id4Idx := findRowIndex(mergedRecord, "id", int64(4))
	require.NotEqual(t, -1, id4Idx)
	assert.Equal(t, "David", getStringValue(mergedRecord, "name", id4Idx))
	assert.Equal(t, 40.0, getFloat64Value(mergedRecord, "value", id4Idx))
	assert.Equal(t, true, getBooleanValue(mergedRecord, "active", id4Idx))
}

// Helper functions for test verification
func findRowIndex(record arrow.Record, colName string, value interface{}) int {
	col := record.ColumnByName(colName)
	if col == nil {
		return -1
	}

	switch v := value.(type) {
	case int64:
		arr := col.(*array.Int64)
		for i := 0; i < arr.Len(); i++ {
			if !arr.IsNull(i) && arr.Value(i) == v {
				return i
			}
		}
	}
	return -1
}

func getStringValue(record arrow.Record, colName string, rowIndex int) string {
	col := record.ColumnByName(colName)
	if col == nil || col.IsNull(rowIndex) {
		return ""
	}
	return col.(*array.String).Value(rowIndex)
}

func getFloat64Value(record arrow.Record, colName string, rowIndex int) float64 {
	col := record.ColumnByName(colName)
	if col == nil || col.IsNull(rowIndex) {
		return 0
	}
	return col.(*array.Float64).Value(rowIndex)
}

func getBooleanValue(record arrow.Record, colName string, rowIndex int) bool {
	col := record.ColumnByName(colName)
	if col == nil || col.IsNull(rowIndex) {
		return false
	}
	return col.(*array.Boolean).Value(rowIndex)
} 