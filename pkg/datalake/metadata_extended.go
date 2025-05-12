package datalake

import (
	"fmt"
	"github.com/apache/arrow/go/v15/arrow"
)

// getArrowDataType converts a Delta Lake type string to an Arrow data type
func getArrowDataType(typeStr string) (arrow.DataType, error) {
	switch typeStr {
	case "int32":
		return &arrow.Int32Type{}, nil
	case "string", "utf8":
		return &arrow.StringType{}, nil
	case "double", "float64":
		return &arrow.Float64Type{}, nil
	case "timestamp":
		return &arrow.TimestampType{Unit: arrow.Second}, nil
	case "boolean":
		return &arrow.BooleanType{}, nil
	default:
		return nil, fmt.Errorf("unsupported field type: %s", typeStr)
	}
}

// getTypeString converts an Arrow data type to a Delta Lake type string
func getTypeString(dataType arrow.DataType) (string, error) {
	switch dataType.(type) {
	case *arrow.Int32Type:
		return "int32", nil
	case *arrow.StringType:
		return "utf8", nil
	case *arrow.Float64Type:
		return "float64", nil
	case *arrow.TimestampType:
		return "timestamp", nil
	case *arrow.BooleanType:
		return "boolean", nil
	default:
		return "", fmt.Errorf("unsupported field type: %T", dataType)
	}
}
