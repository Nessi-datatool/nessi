package datalake

import (
	"time"

	"github.com/apache/arrow/go/v15/arrow"
)

// DeltaTableMetadata represents metadata information about a Delta table
type DeltaTableMetadata struct {
	Format           string    `json:"format"`
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	Location         string    `json:"location"`
	SchemaString     string    `json:"schemaString"`
	PartitionColumns []string  `json:"partitionColumns"`
	CreatedTime      time.Time `json:"createdTime"`
	LastModified     time.Time `json:"lastModified"`
}

// DeltaConnector defines the interface for interacting with Delta Lake tables
type DeltaConnector interface {
	// IsDeltaTable checks if the given path is a Delta table
	IsDeltaTable(path string) bool

	// Read reads data from a Delta table
	Read(path string) ([]map[string]interface{}, error)

	// ReadWithInference reads data from a Delta table and infers the Arrow schema
	ReadWithInference(path string) ([]map[string]interface{}, *arrow.Schema, error)

	// Write writes data to a Delta table
	Write(path string, data []map[string]interface{}, schema *arrow.Schema) error

	// GetDeltaTableMetadata gets metadata about a Delta table
	GetDeltaTableMetadata(path string) (*DeltaTableMetadata, error)

	// ReadAsOfVersion reads data from a Delta table as of a specific version
	ReadAsOfVersion(path string, version int) ([]map[string]interface{}, error)

	// ReadAsOfTimestamp reads data from a Delta table as of a specific timestamp
	ReadAsOfTimestamp(path string, timestamp time.Time) ([]map[string]interface{}, error)

	// GetVersionHistory gets the version history of a Delta table
	GetVersionHistory(path string) ([]DeltaVersion, error)
}
