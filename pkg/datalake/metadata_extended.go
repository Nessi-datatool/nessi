package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/nessi-dev/nessi-dev/pkg/logging"
)

// getArrowDataType converts a Delta Lake type to an Arrow data type
func getArrowDataType(typeName string) (arrow.DataType, error) {
	logger := logging.GetLogger()
	logger.Debug("Converting Delta Lake type to Arrow type", "type", typeName)

	if typeName == "" {
		logger.Error("Empty type name provided")
		return nil, fmt.Errorf("type name cannot be empty")
	}

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
	case "date":
		return arrow.FixedWidthTypes.Date32, nil
	case "binary":
		return arrow.BinaryTypes.Binary, nil
	case "decimal":
		// Default precision and scale for decimal
		return &arrow.Decimal128Type{Precision: 38, Scale: 10}, nil
	default:
		logger.Error("Unsupported data type", "type", typeName)
		return nil, fmt.Errorf("unsupported data type: %s", typeName)
	}
}

// WriteMetadataFile writes a metadata file to the transaction log
func (m *MetadataManager) WriteMetadataFile(version int64, metadata map[string]interface{}) error {
	logger := logging.GetLogger()
	logger.Debug("Writing metadata file", "tablePath", m.tablePath, "version", version)

	if metadata == nil {
		logger.Error("Metadata cannot be nil")
		return fmt.Errorf("metadata cannot be nil")
	}

	if version < 0 {
		logger.Error("Version cannot be negative", "version", version)
		return fmt.Errorf("version cannot be negative: %d", version)
	}

	// Use mutex to ensure thread safety
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create the _delta_log directory if it doesn't exist
	logDir := filepath.Join(m.tablePath, "_delta_log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		logger.Error("Failed to create _delta_log directory", "error", err)
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Convert metadata to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal metadata to JSON", "error", err)
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	// Write to transaction log
	logFile := filepath.Join(logDir, fmt.Sprintf("%020d.json", version))
	err = os.WriteFile(logFile, data, 0644)
	if err != nil {
		logger.Error("Failed to write metadata file", "error", err, "file", logFile)
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Invalidate caches since we've written a new file
	m.cache = nil
	m.versions = nil

	logger.Info("Successfully wrote metadata file", "version", version, "file", logFile)
	return nil
}

// WriteCommitFile writes a commit file to the transaction log
func (m *MetadataManager) WriteCommitFile(version int64, commitInfo map[string]interface{}) error {
	logger := logging.GetLogger()
	logger.Debug("Writing commit file", "tablePath", m.tablePath, "version", version)

	if commitInfo == nil {
		logger.Error("Commit info cannot be nil")
		return fmt.Errorf("commit info cannot be nil")
	}

	if version < 0 {
		logger.Error("Version cannot be negative", "version", version)
		return fmt.Errorf("version cannot be negative: %d", version)
	}

	// Use mutex to ensure thread safety
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Create the _delta_log directory if it doesn't exist
	logDir := filepath.Join(m.tablePath, "_delta_log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		logger.Error("Failed to create _delta_log directory", "error", err)
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Add timestamp if not present
	if _, ok := commitInfo["timestamp"]; !ok {
		commitInfo["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000) // Convert to milliseconds
	}

	// Convert commit info to JSON
	data, err := json.MarshalIndent(commitInfo, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal commit info to JSON", "error", err)
		return fmt.Errorf("failed to marshal commit info to JSON: %w", err)
	}

	// Write commit info
	commitFile := filepath.Join(logDir, fmt.Sprintf("%020d.commit.json", version))
	err = os.WriteFile(commitFile, data, 0644)
	if err != nil {
		logger.Error("Failed to write commit file", "error", err, "file", commitFile)
		return fmt.Errorf("failed to write commit file: %w", err)
	}

	logger.Info("Successfully wrote commit file", "version", version, "file", commitFile)
	return nil
}

// getTypeString converts an Arrow data type to a Delta Lake type string
func getTypeString(dataType arrow.DataType) (string, error) {
	logger := logging.GetLogger()
	logger.Debug("Converting Arrow type to Delta Lake type", "type", fmt.Sprintf("%T", dataType))

	if dataType == nil {
		logger.Error("Data type cannot be nil")
		return "", fmt.Errorf("data type cannot be nil")
	}

	switch dataType.(type) {
	case *arrow.Int32Type:
		return "int32", nil
	case *arrow.Int64Type:
		return "int64", nil
	case *arrow.Float32Type:
		return "float32", nil
	case *arrow.Float64Type:
		return "float64", nil
	case *arrow.StringType:
		return "utf8", nil
	case *arrow.BinaryType:
		return "binary", nil
	case *arrow.BooleanType:
		return "boolean", nil
	case *arrow.TimestampType:
		return "timestamp", nil
	case *arrow.Date32Type:
		return "date", nil
	case *arrow.Decimal128Type:
		return "decimal", nil
	default:
		logger.Error("Unsupported Arrow data type", "type", fmt.Sprintf("%T", dataType))
		return "", fmt.Errorf("unsupported field type: %T", dataType)
	}
}
