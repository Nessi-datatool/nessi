package delta

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
)

// Action represents a Delta Lake transaction log action
type Action struct {
	Meta     *MetaAction     `json:"meta,omitempty"`
	Add      *AddAction      `json:"add,omitempty"`
	Remove   *RemoveAction   `json:"remove,omitempty"`
	Protocol *ProtocolAction `json:"protocol,omitempty"`
}

// MetaAction represents metadata changes
type MetaAction struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Format          Format            `json:"format"`
	SchemaString    string            `json:"schemaString"`
	PartitionColumns []string         `json:"partitionColumns"`
	Configuration   map[string]string `json:"configuration"`
	CreatedTime     int64             `json:"createdTime"`
}

// Format represents the Delta table format
type Format struct {
	Provider string `json:"provider"`
	Options  map[string]string `json:"options"`
}

// AddAction represents adding a file
type AddAction struct {
	Path             string            `json:"path"`
	PartitionValues  map[string]string `json:"partitionValues"`
	Size             int64             `json:"size"`
	ModificationTime int64             `json:"modificationTime"`
	DataChange       bool              `json:"dataChange"`
	Stats            string            `json:"stats,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

// RemoveAction represents removing a file
type RemoveAction struct {
	Path             string            `json:"path"`
	DeletionTimestamp int64            `json:"deletionTimestamp"`
	DataChange       bool              `json:"dataChange"`
	ExtendedFileMetadata bool          `json:"extendedFileMetadata"`
	PartitionValues  map[string]string `json:"partitionValues,omitempty"`
	Size             int64             `json:"size,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

// ProtocolAction represents protocol changes
type ProtocolAction struct {
	MinReaderVersion int `json:"minReaderVersion"`
	MinWriterVersion int `json:"minWriterVersion"`
}

// NewMetaAction creates a new metadata action
func NewMetaAction(schema *arrow.Schema, partitionColumns []string) *Action {
	schemaJSON, _ := json.Marshal(schema)
	return &Action{
		Meta: &MetaAction{
			ID:              generateUUID(),
			Name:            "",
			Description:     "",
			Format: Format{
				Provider: "parquet",
				Options:  make(map[string]string),
			},
			SchemaString:     string(schemaJSON),
			PartitionColumns: partitionColumns,
			Configuration:    make(map[string]string),
			CreatedTime:      time.Now().UnixMilli(),
		},
	}
}

// NewAddAction creates a new add action
func NewAddAction(path string, partitionValues map[string]string, size int64, numRecords int64) *Action {
	return &Action{
		Add: &AddAction{
			Path:             path,
			PartitionValues:  partitionValues,
			Size:             size,
			ModificationTime: time.Now().UnixMilli(),
			DataChange:       true,
			Stats:            fmt.Sprintf(`{"numRecords":%d}`, numRecords),
			Tags:             make(map[string]string),
		},
	}
}

// NewRemoveAction creates a new remove action
func NewRemoveAction(path string, partitionValues map[string]string, size int64) *Action {
	return &Action{
		Remove: &RemoveAction{
			Path:              path,
			DeletionTimestamp: time.Now().UnixMilli(),
			DataChange:        true,
			ExtendedFileMetadata: true,
			PartitionValues:   partitionValues,
			Size:              size,
			Tags:              make(map[string]string),
		},
	}
}

// NewProtocolAction creates a new protocol action
func NewProtocolAction(minReaderVersion, minWriterVersion int) *Action {
	return &Action{
		Protocol: &ProtocolAction{
			MinReaderVersion: minReaderVersion,
			MinWriterVersion: minWriterVersion,
		},
	}
}

// generateUUID generates a unique identifier
func generateUUID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

// Validate validates an action
func (a *Action) Validate() error {
	if a.Meta == nil && a.Add == nil && a.Remove == nil && a.Protocol == nil {
		return fmt.Errorf("action must have at least one of: meta, add, remove, or protocol")
	}

	if a.Meta != nil {
		if a.Meta.ID == "" {
			return fmt.Errorf("meta action must have an ID")
		}
		if a.Meta.SchemaString == "" {
			return fmt.Errorf("meta action must have a schema")
		}
	}

	if a.Add != nil {
		if a.Add.Path == "" {
			return fmt.Errorf("add action must have a path")
		}
		if a.Add.Size <= 0 {
			return fmt.Errorf("add action must have a positive size")
		}
	}

	if a.Remove != nil {
		if a.Remove.Path == "" {
			return fmt.Errorf("remove action must have a path")
		}
	}

	if a.Protocol != nil {
		if a.Protocol.MinReaderVersion < 1 {
			return fmt.Errorf("protocol action must have a minimum reader version >= 1")
		}
		if a.Protocol.MinWriterVersion < 1 {
			return fmt.Errorf("protocol action must have a minimum writer version >= 1")
		}
	}

	return nil
}

// IsDataChange returns whether the action represents a data change
func (a *Action) IsDataChange() bool {
	if a.Add != nil {
		return a.Add.DataChange
	}
	if a.Remove != nil {
		return a.Remove.DataChange
	}
	return false
}

// GetPath returns the path associated with the action
func (a *Action) GetPath() string {
	if a.Add != nil {
		return a.Add.Path
	}
	if a.Remove != nil {
		return a.Remove.Path
	}
	return ""
}

// GetPartitionValues returns the partition values associated with the action
func (a *Action) GetPartitionValues() map[string]string {
	if a.Add != nil {
		return a.Add.PartitionValues
	}
	if a.Remove != nil {
		return a.Remove.PartitionValues
	}
	return nil
}

// GetSize returns the size associated with the action
func (a *Action) GetSize() int64 {
	if a.Add != nil {
		return a.Add.Size
	}
	if a.Remove != nil {
		return a.Remove.Size
	}
	return 0
}

// GetModificationTime returns the modification time associated with the action
func (a *Action) GetModificationTime() time.Time {
	if a.Add != nil {
		return time.Unix(0, a.Add.ModificationTime*int64(time.Millisecond))
	}
	if a.Remove != nil {
		return time.Unix(0, a.Remove.DeletionTimestamp*int64(time.Millisecond))
	}
	return time.Time{}
} 