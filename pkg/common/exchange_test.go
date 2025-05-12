package common

import (
	"testing"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/stretchr/testify/assert"
)

func TestExchangeManager(t *testing.T) {
	config := ExchangeConfig{
		ArrowThreshold: 1024 * 1024, // 1MB
		MaxMemoryUsage: 1024 * 1024 * 1024, // 1GB
		UseCompression: false,
	}

	manager := NewExchangeManager(config)
	assert.NotNil(t, manager)

	// Register strategies
	manager.RegisterStrategy(NewArrowStrategy())
	manager.RegisterStrategy(NewJSONStrategy())

	// Create test Arrow record
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	b.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3}, nil)
	b.Field(1).(*array.Float64Builder).AppendValues([]float64{1.1, 2.2, 3.3}, nil)

	record := b.NewRecord()
	defer record.Release()

	// Test Arrow strategy selection
	strategy, err := manager.SelectStrategy(record)
	assert.NoError(t, err)
	assert.IsType(t, &ArrowStrategy{}, strategy)

	// Test Arrow serialization
	data, err := strategy.Serialize(record)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test Arrow deserialization
	var result arrow.Record
	err = strategy.Deserialize(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, record.NumRows(), result.NumRows())
}

func TestJSONStrategy(t *testing.T) {
	strategy := NewJSONStrategy()
	assert.NotNil(t, strategy)

	// Test data
	testData := map[string]interface{}{
		"name": "test",
		"age":  30,
	}

	// Test serialization
	data, err := strategy.Serialize(testData)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test deserialization
	var result map[string]interface{}
	err = strategy.Deserialize(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, testData["name"], result["name"])
	// JSON unmarshals numbers as float64 by default
	if age, ok := result["age"].(float64); ok {
		assert.Equal(t, float64(testData["age"].(int)), age)
	} else {
		t.Errorf("expected float64 for age, got %T", result["age"])
	}
}
