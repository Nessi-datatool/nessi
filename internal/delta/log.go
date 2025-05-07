package delta

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
)

// Action represents a Delta log action
type Action struct {
	Add    *AddFile    `json:"add,omitempty"`
	Remove *RemoveFile `json:"remove,omitempty"`
	Meta   *MetaData   `json:"metaData,omitempty"`
}

// AddFile represents an add file action
type AddFile struct {
	Path             string            `json:"path"`
	PartitionValues  map[string]string `json:"partitionValues"`
	Size             int64             `json:"size"`
	ModificationTime int64             `json:"modificationTime"`
	DataChange       bool              `json:"dataChange"`
	Stats            string            `json:"stats,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

// RemoveFile represents a remove file action
type RemoveFile struct {
	Path             string            `json:"path"`
	DeletionTimestamp int64            `json:"deletionTimestamp,omitempty"`
	DataChange       bool              `json:"dataChange"`
	ExtendedFileMetadata bool          `json:"extendedFileMetadata,omitempty"`
	PartitionValues  map[string]string `json:"partitionValues,omitempty"`
	Size             int64             `json:"size,omitempty"`
}

// MetaData represents table metadata
type MetaData struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Format      Format            `json:"format"`
	Schema      Schema            `json:"schema"`
	PartitionColumns []string     `json:"partitionColumns"`
	Configuration map[string]string `json:"configuration"`
	CreatedTime int64             `json:"createdTime"`
}

// Format represents the table format
type Format struct {
	Provider string            `json:"provider"`
	Options  map[string]string `json:"options,omitempty"`
}

// Schema represents the table schema
type Schema struct {
	Type        string   `json:"type"`
	Fields      []Field  `json:"fields"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Field represents a schema field
type Field struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Nullable bool     `json:"nullable"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DeltaLog represents the Delta log parser
type DeltaLog struct {
	tablePath string
	version   int64
	actions   []Action
	metadata  *MetaData
}

// NewDeltaLog creates a new Delta log parser
func NewDeltaLog(tablePath string) (*DeltaLog, error) {
	log := &DeltaLog{
		tablePath: tablePath,
	}

	if err := log.load(); err != nil {
		return nil, fmt.Errorf("failed to load Delta log: %w", err)
	}

	return log, nil
}

// load loads the Delta log from the table path
func (l *DeltaLog) load() error {
	logPath := filepath.Join(l.tablePath, "_delta_log")
	if _, err := os.Stat(logPath); err != nil {
		return fmt.Errorf("Delta log directory not found: %w", err)
	}

	// Read all JSON files in the log directory
	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	// Sort files by version number
	var versions []int64
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			version, err := strconv.ParseInt(strings.TrimSuffix(file.Name(), ".json"), 10, 64)
			if err != nil {
				continue
			}
			versions = append(versions, version)
		}
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i] < versions[j] })

	// Process each version in order
	for _, version := range versions {
		if err := l.processVersion(version); err != nil {
			return fmt.Errorf("failed to process version %d: %w", version, err)
		}
	}

	return nil
}

// processVersion processes a single version of the Delta log
func (l *DeltaLog) processVersion(version int64) error {
	filePath := filepath.Join(l.tablePath, "_delta_log", fmt.Sprintf("%d.json", version))
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	var actions []Action
	if err := json.Unmarshal(data, &actions); err != nil {
		return fmt.Errorf("failed to parse log file: %w", err)
	}

	// Process each action
	for _, action := range actions {
		if action.Meta != nil {
			l.metadata = action.Meta
		}
		l.actions = append(l.actions, action)
	}

	l.version = version
	return nil
}

// GetVersion returns the current version
func (l *DeltaLog) GetVersion() int64 {
	return l.version
}

// GetMetadata returns the current metadata
func (l *DeltaLog) GetMetadata() *MetaData {
	return l.metadata
}

// GetActions returns all actions
func (l *DeltaLog) GetActions() []Action {
	return l.actions
}

// GetActiveFiles returns the list of active files
func (l *DeltaLog) GetActiveFiles() []AddFile {
	var activeFiles []AddFile
	removedPaths := make(map[string]bool)

	// Process actions in reverse to find the latest state
	for i := len(l.actions) - 1; i >= 0; i-- {
		action := l.actions[i]
		if action.Remove != nil {
			removedPaths[action.Remove.Path] = true
		} else if action.Add != nil && !removedPaths[action.Add.Path] {
			activeFiles = append(activeFiles, *action.Add)
		}
	}

	return activeFiles
}

// GetPartitionValues returns all unique partition values
func (l *DeltaLog) GetPartitionValues() map[string]map[string]bool {
	partitions := make(map[string]map[string]bool)

	for _, action := range l.actions {
		if action.Add != nil {
			for col, value := range action.Add.PartitionValues {
				if _, ok := partitions[col]; !ok {
					partitions[col] = make(map[string]bool)
				}
				partitions[col][value] = true
			}
		}
	}

	return partitions
}

// GetSchema returns the Arrow schema
func (l *DeltaLog) GetSchema() (*arrow.Schema, error) {
	if l.metadata == nil {
		return nil, fmt.Errorf("no metadata available")
	}

	// Convert Delta schema to Arrow schema
	fields := make([]arrow.Field, 0, len(l.metadata.Schema.Fields))
	for _, field := range l.metadata.Schema.Fields {
		arrowType, err := convertToArrowType(field.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field type: %w", err)
		}

		fields = append(fields, arrow.Field{
			Name:     field.Name,
			Type:     arrowType,
			Nullable: field.Nullable,
			Metadata: arrow.NewMetadata(nil, nil),
		})
	}

	return arrow.NewSchema(fields, nil), nil
}

// convertToArrowType converts a Delta type to an Arrow type
func convertToArrowType(deltaType string) (arrow.DataType, error) {
	switch deltaType {
	case "string":
		return arrow.BinaryTypes.String, nil
	case "integer":
		return arrow.PrimitiveTypes.Int32, nil
	case "long":
		return arrow.PrimitiveTypes.Int64, nil
	case "float":
		return arrow.PrimitiveTypes.Float32, nil
	case "double":
		return arrow.PrimitiveTypes.Float64, nil
	case "boolean":
		return arrow.FixedWidthTypes.Boolean, nil
	case "date":
		return arrow.FixedWidthTypes.Date32, nil
	case "timestamp":
		return arrow.FixedWidthTypes.Timestamp_ns, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", deltaType)
	}
} 