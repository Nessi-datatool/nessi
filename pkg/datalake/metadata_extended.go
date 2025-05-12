package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/apache/arrow/go/v15/arrow"
)

// getArrowDataType converts a Delta Lake type to an Arrow data type
func getArrowDataType(typeName string) (arrow.DataType, error) {
	switch typeName {
	case "int32":
		return arrow.PrimitiveTypes.Int32, nil
	case "int64":
		return arrow.PrimitiveTypes.Int64, nil
	case "float32":
		return arrow.PrimitiveTypes.Float32, nil
	case "float64":
		return arrow.PrimitiveTypes.Float64, nil
	case "boolean":
		return arrow.FixedWidthTypes.Boolean, nil
	case "utf8", "string":
		return arrow.BinaryTypes.String, nil
	case "timestamp":
		return arrow.FixedWidthTypes.Timestamp_us, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %s", typeName)
	}
}

// WriteMetadataFile writes a metadata file to the transaction log
func (m *MetadataManager) WriteMetadataFile(version int64, metadata map[string]interface{}) error {
	// Create the _delta_log directory if it doesn't exist
	logDir := filepath.Join(m.tablePath, "_delta_log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Convert metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	// Write to transaction log
	logFile := filepath.Join(logDir, fmt.Sprintf("%020d.json", version))
	err = os.WriteFile(logFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	return nil
}

// WriteCommitFile writes a commit file to the transaction log
func (m *MetadataManager) WriteCommitFile(version int64, commitInfo map[string]interface{}) error {
	// Create the _delta_log directory if it doesn't exist
	logDir := filepath.Join(m.tablePath, "_delta_log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Convert commit info to JSON
	data, err := json.MarshalIndent(commitInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal commit info to JSON: %w", err)
	}

	// Write commit info
	commitFile := filepath.Join(logDir, fmt.Sprintf("%020d.commit.json", version))
	err = os.WriteFile(commitFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write commit file: %w", err)
	}

	return nil
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
