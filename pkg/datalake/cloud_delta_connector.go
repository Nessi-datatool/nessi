package datalake

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/cloud/common"
)

// CloudDeltaConnector provides access to Delta Lake tables stored in cloud storage
type CloudDeltaConnector struct {
	provider      common.CloudProvider
	bucket        string
	defaultPrefix string
}

// NewCloudDeltaConnector creates a new CloudDeltaConnector
func NewCloudDeltaConnector(provider common.CloudProvider, bucket, defaultPrefix string) *CloudDeltaConnector {
	return &CloudDeltaConnector{
		provider:      provider,
		bucket:        bucket,
		defaultPrefix: defaultPrefix,
	}
}

// DeltaLogEntry represents an entry in the Delta Lake transaction log
type DeltaLogEntry struct {
	Add           []DeltaAddFile         `json:"add,omitempty"`
	Remove        []DeltaRemoveFile      `json:"remove,omitempty"`
	Metadata      *DeltaMetadata         `json:"metaData,omitempty"`
	Protocol      *DeltaProtocol         `json:"protocol,omitempty"`
	CommitInfo    map[string]interface{} `json:"commitInfo,omitempty"`
	Timestamp     int64                  `json:"timestamp,omitempty"`
	Operation     string                 `json:"operation,omitempty"`
	SchemaChange  *DeltaSchemaChange     `json:"schemaChange,omitempty"`
	TransactionId int64                  `json:"txn,omitempty"`
	DependsOn     []int64                `json:"dependsOn,omitempty"`
	AppMetadata   map[string]interface{} `json:"appMetadata,omitempty"`
}

// DeltaAddFile represents a file added to a Delta Lake table
type DeltaAddFile struct {
	Path             string                 `json:"path"`
	Size             int64                  `json:"size"`
	ModificationTime int64                  `json:"modificationTime"`
	DataChange       bool                   `json:"dataChange"`
	Stats            string                 `json:"stats,omitempty"`
	Partition        map[string]interface{} `json:"partitionValues,omitempty"`
	Tags             map[string]string      `json:"tags,omitempty"`
}

// DeltaRemoveFile represents a file removed from a Delta Lake table
type DeltaRemoveFile struct {
	Path             string                 `json:"path"`
	Size             int64                  `json:"size,omitempty"`
	ModificationTime int64                  `json:"modificationTime,omitempty"`
	DataChange       bool                   `json:"dataChange"`
	Partition        map[string]interface{} `json:"partitionValues,omitempty"`
	Tags             map[string]string      `json:"tags,omitempty"`
}

// DeltaMetadata represents the metadata of a Delta Lake table
type DeltaMetadata struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	Format           map[string]interface{} `json:"format"`
	SchemaString     string                 `json:"schemaString"`
	PartitionColumns []string               `json:"partitionColumns"`
	Configuration    map[string]string      `json:"configuration,omitempty"`
	CreatedTime      int64                  `json:"createdTime,omitempty"`
	Properties       map[string]string      `json:"properties,omitempty"`
}

// DeltaProtocol represents the protocol of a Delta Lake table
type DeltaProtocol struct {
	MinReaderVersion int `json:"minReaderVersion"`
	MinWriterVersion int `json:"minWriterVersion"`
}

// DeltaSchemaChange represents a schema change in a Delta Lake table
type DeltaSchemaChange struct {
	Schema        string `json:"schema"`
	SchemaChanges []struct {
		Column     string `json:"column"`
		ChangeType string `json:"changeType"`
	} `json:"schemaChanges"`
}

// DeltaTableInfo represents information about a Delta Lake table
type DeltaTableInfo struct {
	Path           string
	Version        int64
	Metadata       *DeltaMetadata
	Protocol       *DeltaProtocol
	Files          []DeltaAddFile
	LastCommitTime time.Time
	LastCommitInfo map[string]interface{}
	SchemaChanges  []DeltaSchemaChange
}

// GetTableInfo retrieves information about a Delta Lake table
func (c *CloudDeltaConnector) GetTableInfo(ctx context.Context, tablePath string) (*DeltaTableInfo, error) {
	// Normalize table path
	tablePath = c.normalizePath(tablePath)

	// Get the latest version of the table
	version, err := c.getLatestVersion(ctx, tablePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	// Get the transaction log entries
	entries, err := c.getTransactionLog(ctx, tablePath, 0, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction log: %w", err)
	}

	// Process transaction log entries
	tableInfo := &DeltaTableInfo{
		Path:    tablePath,
		Version: version,
		Files:   make([]DeltaAddFile, 0),
	}

	// Track active files (added but not removed)
	activeFiles := make(map[string]DeltaAddFile)

	// Process each log entry
	for _, entry := range entries {
		// Process metadata
		if entry.Metadata != nil {
			tableInfo.Metadata = entry.Metadata
		}

		// Process protocol
		if entry.Protocol != nil {
			tableInfo.Protocol = entry.Protocol
		}

		// Process commit info
		if entry.CommitInfo != nil {
			tableInfo.LastCommitInfo = entry.CommitInfo
			tableInfo.LastCommitTime = time.Unix(0, entry.Timestamp*1000000)
		}

		// Process schema changes
		if entry.SchemaChange != nil {
			tableInfo.SchemaChanges = append(tableInfo.SchemaChanges, *entry.SchemaChange)
		}

		// Process add files
		for _, addFile := range entry.Add {
			activeFiles[addFile.Path] = addFile
		}

		// Process remove files
		for _, removeFile := range entry.Remove {
			delete(activeFiles, removeFile.Path)
		}
	}

	// Convert active files map to slice
	for _, file := range activeFiles {
		tableInfo.Files = append(tableInfo.Files, file)
	}

	return tableInfo, nil
}

// getLatestVersion returns the latest version of a Delta Lake table
func (c *CloudDeltaConnector) getLatestVersion(ctx context.Context, tablePath string) (int64, error) {
	// Check for _delta_log/_last_checkpoint file
	checkpointPath := path.Join(tablePath, "_delta_log", "_last_checkpoint")
	checkpointExists, err := c.objectExists(ctx, checkpointPath)
	if err != nil {
		return 0, fmt.Errorf("failed to check for checkpoint: %w", err)
	}

	if checkpointExists {
		// Read checkpoint file
		checkpointData, err := c.readObject(ctx, checkpointPath)
		if err != nil {
			return 0, fmt.Errorf("failed to read checkpoint: %w", err)
		}

		// Parse checkpoint
		var checkpoint struct {
			Version int64 `json:"version"`
		}
		if err := json.Unmarshal(checkpointData, &checkpoint); err != nil {
			return 0, fmt.Errorf("failed to parse checkpoint: %w", err)
		}

		return checkpoint.Version, nil
	}

	// No checkpoint, scan log files
	logDir := path.Join(tablePath, "_delta_log")
	objects, err := c.provider.ListObjects(ctx, c.bucket, logDir)
	if err != nil {
		return 0, fmt.Errorf("failed to list log files: %w", err)
	}

	// Find the highest version
	var maxVersion int64 = -1
	for _, obj := range objects {
		// Check if it's a JSON log file
		if !strings.HasSuffix(obj.Key, ".json") {
			continue
		}

		// Extract version from filename
		filename := path.Base(obj.Key)
		var version int64
		if _, err := fmt.Sscanf(filename, "%020d.json", &version); err != nil {
			continue
		}

		if version > maxVersion {
			maxVersion = version
		}
	}

	if maxVersion == -1 {
		return 0, fmt.Errorf("no valid log files found")
	}

	return maxVersion, nil
}

// getTransactionLog returns the transaction log entries for a Delta Lake table
func (c *CloudDeltaConnector) getTransactionLog(ctx context.Context, tablePath string, startVersion, endVersion int64) ([]DeltaLogEntry, error) {
	var entries []DeltaLogEntry

	// Read log files from startVersion to endVersion
	for version := startVersion; version <= endVersion; version++ {
		logPath := path.Join(tablePath, "_delta_log", fmt.Sprintf("%020d.json", version))

		// Check if log file exists
		exists, err := c.objectExists(ctx, logPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check log file %d: %w", version, err)
		}

		if !exists {
			// Skip non-existent log files
			continue
		}

		// Read log file
		logData, err := c.readObject(ctx, logPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read log file %d: %w", version, err)
		}

		// Parse log entries (each line is a separate JSON object)
		lines := strings.Split(string(logData), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			var entry DeltaLogEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				return nil, fmt.Errorf("failed to parse log entry: %w", err)
			}

			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// objectExists checks if an object exists in cloud storage
func (c *CloudDeltaConnector) objectExists(ctx context.Context, objectPath string) (bool, error) {
	// Extract key from full path
	key := c.getObjectKey(objectPath)

	// List objects with the given key as prefix
	objects, err := c.provider.ListObjects(ctx, c.bucket, key)
	if err != nil {
		return false, fmt.Errorf("failed to list objects: %w", err)
	}

	// Check if the object exists
	for _, obj := range objects {
		if obj.Key == key {
			return true, nil
		}
	}

	return false, nil
}

// readObject reads an object from cloud storage
func (c *CloudDeltaConnector) readObject(ctx context.Context, objectPath string) ([]byte, error) {
	// Extract key from full path
	key := c.getObjectKey(objectPath)

	// Get object
	reader, err := c.provider.GetObject(ctx, c.bucket, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer reader.Close()

	// Read object data
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	return data, nil
}

// normalizePath normalizes a table path
func (c *CloudDeltaConnector) normalizePath(tablePath string) string {
	// If the path doesn't start with the default prefix, prepend it
	if c.defaultPrefix != "" && !strings.HasPrefix(tablePath, c.defaultPrefix) {
		return path.Join(c.defaultPrefix, tablePath)
	}
	return tablePath
}

// getObjectKey extracts the object key from a full path
func (c *CloudDeltaConnector) getObjectKey(fullPath string) string {
	// Remove bucket name if present
	key := fullPath
	if strings.HasPrefix(key, c.bucket+"/") {
		key = strings.TrimPrefix(key, c.bucket+"/")
	}
	return key
}

// GetTableSchema returns the schema of a Delta Lake table
func (c *CloudDeltaConnector) GetTableSchema(ctx context.Context, tablePath string) (string, error) {
	// Get table info
	tableInfo, err := c.GetTableInfo(ctx, tablePath)
	if err != nil {
		return "", fmt.Errorf("failed to get table info: %w", err)
	}

	// Return schema
	if tableInfo.Metadata != nil {
		return tableInfo.Metadata.SchemaString, nil
	}

	return "", fmt.Errorf("no schema found for table")
}

// GetTableVersions returns the available versions of a Delta Lake table
func (c *CloudDeltaConnector) GetTableVersions(ctx context.Context, tablePath string) ([]int64, error) {
	// Normalize table path
	tablePath = c.normalizePath(tablePath)

	// List log files
	logDir := path.Join(tablePath, "_delta_log")
	objects, err := c.provider.ListObjects(ctx, c.bucket, logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list log files: %w", err)
	}

	// Extract versions from log files
	var versions []int64
	for _, obj := range objects {
		// Check if it's a JSON log file
		if !strings.HasSuffix(obj.Key, ".json") {
			continue
		}

		// Extract version from filename
		filename := path.Base(obj.Key)
		var version int64
		if _, err := fmt.Sscanf(filename, "%020d.json", &version); err != nil {
			continue
		}

		versions = append(versions, version)
	}

	return versions, nil
}

// GetTableAtVersion returns the state of a Delta Lake table at a specific version
func (c *CloudDeltaConnector) GetTableAtVersion(ctx context.Context, tablePath string, version int64) (*DeltaTableInfo, error) {
	// Normalize table path
	tablePath = c.normalizePath(tablePath)

	// Check if version exists
	versions, err := c.GetTableVersions(ctx, tablePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get table versions: %w", err)
	}

	versionExists := false
	for _, v := range versions {
		if v == version {
			versionExists = true
			break
		}
	}

	if !versionExists {
		return nil, fmt.Errorf("version %d does not exist", version)
	}

	// Get transaction log entries up to the specified version
	entries, err := c.getTransactionLog(ctx, tablePath, 0, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction log: %w", err)
	}

	// Process transaction log entries
	tableInfo := &DeltaTableInfo{
		Path:    tablePath,
		Version: version,
		Files:   make([]DeltaAddFile, 0),
	}

	// Track active files (added but not removed)
	activeFiles := make(map[string]DeltaAddFile)

	// Process each log entry
	for _, entry := range entries {
		// Process metadata
		if entry.Metadata != nil {
			tableInfo.Metadata = entry.Metadata
		}

		// Process protocol
		if entry.Protocol != nil {
			tableInfo.Protocol = entry.Protocol
		}

		// Process commit info
		if entry.CommitInfo != nil {
			tableInfo.LastCommitInfo = entry.CommitInfo
			tableInfo.LastCommitTime = time.Unix(0, entry.Timestamp*1000000)
		}

		// Process schema changes
		if entry.SchemaChange != nil {
			tableInfo.SchemaChanges = append(tableInfo.SchemaChanges, *entry.SchemaChange)
		}

		// Process add files
		for _, addFile := range entry.Add {
			activeFiles[addFile.Path] = addFile
		}

		// Process remove files
		for _, removeFile := range entry.Remove {
			delete(activeFiles, removeFile.Path)
		}
	}

	// Convert active files map to slice
	for _, file := range activeFiles {
		tableInfo.Files = append(tableInfo.Files, file)
	}

	return tableInfo, nil
}

// GetTableAtTimestamp returns the state of a Delta Lake table at a specific timestamp
func (c *CloudDeltaConnector) GetTableAtTimestamp(ctx context.Context, tablePath string, timestamp time.Time) (*DeltaTableInfo, error) {
	// Normalize table path
	tablePath = c.normalizePath(tablePath)

	// Get the latest version of the table
	latestVersion, err := c.getLatestVersion(ctx, tablePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	// Get transaction log entries
	entries, err := c.getTransactionLog(ctx, tablePath, 0, latestVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction log: %w", err)
	}

	// Find the version at or before the specified timestamp
	targetVersion := int64(-1)
	targetTimestamp := timestamp.UnixNano() / 1000000 // Convert to milliseconds

	for _, entry := range entries {
		if entry.CommitInfo != nil && entry.Timestamp <= targetTimestamp {
			targetVersion = entry.TransactionId
		}
	}

	if targetVersion == -1 {
		return nil, fmt.Errorf("no version found at or before the specified timestamp")
	}

	// Get the table at the target version
	return c.GetTableAtVersion(ctx, tablePath, targetVersion)
}
