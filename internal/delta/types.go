package delta

import (
	"errors"
	"time"
)

// Config represents Delta table configuration
type Config struct {
	ReadOnly      bool
	MaxRetries    int
	RetryInterval time.Duration
}

// Action represents a Delta log action
type Action struct {
	Add    *AddAction    `json:"add,omitempty"`
	Remove *RemoveAction `json:"remove,omitempty"`
	Meta   *MetaAction   `json:"metaData,omitempty"`
}

// Validate checks if the Action is valid.
// An action is considered valid if at least one of its sub-actions (Add, Remove, Meta) is defined.
func (a Action) Validate() error {
	if a.Add == nil && a.Remove == nil && a.Meta == nil {
		return errors.New("action must have at least one of Add, Remove, or Meta defined")
	}
	// TODO: Consider adding more specific validation, e.g., ensuring only one sub-action is set,
	// or validating the contents of the populated sub-action.
	return nil
}

// AddAction represents a file addition action
type AddAction struct {
	Path             string            `json:"path"`
	Size             int64             `json:"size"`
	ModificationTime int64             `json:"modificationTime"`
	DataChange       bool              `json:"dataChange"`
	Stats            interface{}       `json:"stats,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	PartitionValues  map[string]string `json:"partitionValues,omitempty"`
}

// RemoveAction represents a file removal action
type RemoveAction struct {
	Path         string `json:"path"`
	DeletionTime int64  `json:"deletionTime"`
	DataChange   bool   `json:"dataChange"`
}

// MetaAction represents a metadata action
type MetaAction struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Format           Format            `json:"format"`
	SchemaString     string            `json:"schemaString"`
	PartitionColumns []string          `json:"partitionColumns"`
	Configuration    map[string]string `json:"configuration"`
	CreatedTime      int64             `json:"createdTime"`
}

// Format represents the table format
type Format struct {
	Provider string            `json:"provider"`
	Options  map[string]string `json:"options,omitempty"`
}

// Table is an alias for MetaAction, representing a Delta table's metadata and schema.
// This is used by handlers to return table information.
type Table MetaAction
