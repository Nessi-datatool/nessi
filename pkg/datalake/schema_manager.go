package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/nessi-dev/nessi/pkg/logging"
)

// DeltaSchemaManager handles Delta Lake schema operations with caching and thread safety
type DeltaSchemaManager struct {
	tablePath      string
	metaDir        string
	mutex          sync.RWMutex
	historyCache   *SchemaHistory // Cache for schema history
	currentCache   *SchemaVersion // Cache for current schema
	versionCache   map[int64]*SchemaVersion // Cache for schema versions by version number
	// Used to ensure arrow package is referenced
	arrowSchema    *arrow.Schema
}

// NewDeltaSchemaManager creates a new schema manager with initialized caches
func NewDeltaSchemaManager(tablePath string) *DeltaSchemaManager {
	logger := logging.GetLogger()
	logger.Debug("Creating new SchemaManager", "tablePath", tablePath)
	
	return &DeltaSchemaManager{
		tablePath:    tablePath,
		metaDir:      filepath.Join(tablePath, "_delta_log", "schema_history"),
		versionCache: make(map[int64]*SchemaVersion),
	}
}

// InitializeSchema initializes a schema for a table and updates caches
func (sm *DeltaSchemaManager) InitializeSchema(schema *DeltaSchema, partitionBy []string, zOrderBy []string) (*SchemaVersion, error) {
	logger := logging.GetLogger()
	logger.Debug("Initializing schema", "tablePath", sm.tablePath, "fields", len(schema.Fields))
	
	// Lock for thread safety
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Check if schema already exists by reading directly from disk
	// We don't use GetCurrentSchema to avoid potential cache issues during initialization
	existingHistory, err := sm.readHistory()
	if err == nil && len(existingHistory.Versions) > 0 {
		logger.Warn("Schema already initialized", "tablePath", sm.tablePath)
		return nil, fmt.Errorf("schema already initialized")
	}

	// Convert Schema to arrow.Schema using helper function
	arrowSchema, arrowFields, err := ConvertToArrowSchema(schema)
	if err != nil {
		logger.Error("Failed to convert schema to arrow schema", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Create initial schema version
	schemaVersion := &SchemaVersion{
		Version:       1,
		Timestamp:     time.Now(),
		Schema:        arrowSchema,
		SchemaFields:  arrowFields,
	}

	// Create schema history
	history := &SchemaHistory{
		Versions: []SchemaVersion{*schemaVersion},
	}

	// Write schema history to file
	if err := sm.writeHistory(history); err != nil {
		logger.Error("Failed to write schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Update caches
	sm.historyCache = history
	sm.currentCache = schemaVersion
	sm.versionCache[schemaVersion.Version] = schemaVersion
	logger.Debug("Schema initialized and caches updated", "tablePath", sm.tablePath, "version", schemaVersion.Version)

	return schemaVersion, nil
}

// GetCurrentSchema gets the current schema with caching
func (sm *DeltaSchemaManager) GetCurrentSchema() (*SchemaVersion, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting current schema", "tablePath", sm.tablePath)
	
	// Check cache first
	sm.mutex.RLock()
	if sm.currentCache != nil {
		defer sm.mutex.RUnlock()
		logger.Debug("Returning cached current schema", "tablePath", sm.tablePath, "version", sm.currentCache.Version)
		return sm.currentCache, nil
	}
	sm.mutex.RUnlock()
	
	// Cache miss, get from history
	history, err := sm.GetSchemaHistory()
	if err != nil {
		logger.Error("Failed to get schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to get current schema: %w", err)
	}

	if len(history.Versions) == 0 {
		logger.Warn("No schema found", "tablePath", sm.tablePath)
		return nil, fmt.Errorf("no schema found")
	}

	// Get the latest schema version (last in the list)
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Update cache
	sm.currentCache = &history.Versions[len(history.Versions)-1]
	logger.Debug("Updated current schema cache", "tablePath", sm.tablePath, "version", sm.currentCache.Version)
	
	return sm.currentCache, nil
}

// GetSchemaVersion gets a specific schema version with caching
func (sm *DeltaSchemaManager) GetSchemaVersion(version int64) (*SchemaVersion, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting schema version", "tablePath", sm.tablePath, "version", version)
	
	// Check cache first
	sm.mutex.RLock()
	if cachedVersion, exists := sm.versionCache[version]; exists {
		defer sm.mutex.RUnlock()
		logger.Debug("Returning cached schema version", "tablePath", sm.tablePath, "version", version)
		return cachedVersion, nil
	}
	sm.mutex.RUnlock()
	
	// Cache miss, get from history
	history, err := sm.GetSchemaHistory()
	if err != nil {
		logger.Error("Failed to get schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to get schema version %d: %w", version, err)
	}

	// Find the requested version
	for _, schema := range history.Versions {
		if schema.Version == version {
			// Update cache
			sm.mutex.Lock()
			defer sm.mutex.Unlock()
			
			// Create a copy to avoid modifying the original
			schemaCopy := schema
			sm.versionCache[version] = &schemaCopy
			logger.Debug("Updated schema version cache", "tablePath", sm.tablePath, "version", version)
			
			return &schemaCopy, nil
		}
	}

	logger.Warn("Schema version not found", "tablePath", sm.tablePath, "version", version)
	return nil, fmt.Errorf("schema version %d not found", version)
}

// GetSchemaHistory gets the schema history with caching
func (sm *DeltaSchemaManager) GetSchemaHistory() (*SchemaHistory, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting schema history", "tablePath", sm.tablePath)
	
	// Check cache first
	sm.mutex.RLock()
	if sm.historyCache != nil {
		defer sm.mutex.RUnlock()
		logger.Debug("Returning cached schema history", "tablePath", sm.tablePath, "versions", len(sm.historyCache.Versions))
		return sm.historyCache, nil
	}
	sm.mutex.RUnlock()
	
	// Cache miss, read from disk
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Double-check after acquiring write lock
	if sm.historyCache != nil {
		logger.Debug("Returning cached schema history (after lock)", "tablePath", sm.tablePath)
		return sm.historyCache, nil
	}
	
	// Read from disk
	history, err := sm.readHistory()
	if err != nil {
		logger.Error("Failed to read schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to get schema history: %w", err)
	}
	
	// Update cache
	sm.historyCache = history
	logger.Debug("Updated schema history cache", "tablePath", sm.tablePath, "versions", len(history.Versions))
	
	return history, nil
}

// UpdateSchema updates the schema and invalidates caches
func (sm *DeltaSchemaManager) UpdateSchema(newSchema *DeltaSchema, commitInfo map[string]string) (*SchemaVersion, error) {
	logger := logging.GetLogger()
	logger.Debug("Updating schema", "tablePath", sm.tablePath, "fields", len(newSchema.Fields))
	
	// Get history with thread safety
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	// Read fresh history from disk to ensure we have the latest
	history, err := sm.readHistory()
	if err != nil {
		logger.Error("Failed to read schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to update schema: %w", err)
	}

	if len(history.Versions) == 0 {
		logger.Warn("No schema found for update", "tablePath", sm.tablePath)
		return nil, fmt.Errorf("no schema found, use InitializeSchema instead")
	}

	// Get the latest schema version
	currentSchema := history.Versions[len(history.Versions)-1]
	logger.Debug("Current schema version", "tablePath", sm.tablePath, "version", currentSchema.Version)
	
	// Convert Schema to arrow.Schema using helper function
	arrowSchema, arrowFields, err := ConvertToArrowSchema(newSchema)
	if err != nil {
		logger.Error("Failed to convert schema to arrow schema", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to update schema: %w", err)
	}

	// Get the current schema as a Schema object for comparison
	currentSchemaObj, err := ConvertFromArrowSchema(currentSchema.Schema)
	if err != nil {
		logger.Error("Failed to convert current arrow schema to schema", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to convert current schema for comparison: %w", err)
	}
	
	// Compare schemas to detect changes
	schemaChanges, err := sm.CompareSchemas(currentSchemaObj, newSchema)
	if err != nil {
		logger.Error("Failed to compare schemas", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to detect schema changes: %w", err)
	}
	
	// Log the detected changes
	if len(schemaChanges) > 0 {
		logger.Info("Schema changes detected", "tablePath", sm.tablePath, "changeCount", len(schemaChanges))
		for i, change := range schemaChanges {
			logger.Info("Schema change", "tablePath", sm.tablePath, "index", i, 
				"type", change.Type, "field", change.FieldName)
		}
	} else {
		logger.Info("No schema changes detected", "tablePath", sm.tablePath)
	}

	// Create new schema version with change information
	newVersion := &SchemaVersion{
		Version:       currentSchema.Version + 1,
		Timestamp:     time.Now(),
		Schema:        arrowSchema,
		SchemaFields:  arrowFields,
	}

	// Add new schema to history
	history.Versions = append(history.Versions, *newVersion)

	// Write schema history to file
	if err := sm.writeHistory(history); err != nil {
		logger.Error("Failed to write schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to save updated schema: %w", err)
	}

	// Update caches
	sm.historyCache = history
	sm.currentCache = newVersion
	sm.versionCache[newVersion.Version] = newVersion
	logger.Debug("Updated schema caches", "tablePath", sm.tablePath, "newVersion", newVersion.Version)

	return newVersion, nil
}

// ValidateSchemaCompatibility checks if a new schema is compatible with the current schema
func (sm *DeltaSchemaManager) ValidateSchemaCompatibility(newSchema *DeltaSchema) (bool, []string, error) {
	logger := logging.GetLogger()
	logger.Debug("Validating schema compatibility", "tablePath", sm.tablePath)
	
	// Get current schema
	currentVersion, err := sm.GetCurrentSchema()
	if err != nil {
		logger.Error("Failed to get current schema", "tablePath", sm.tablePath, "error", err)
		return false, nil, fmt.Errorf("failed to validate schema compatibility: %w", err)
	}
	
	// Convert current schema to Schema object
	currentSchema, err := ConvertFromArrowSchema(currentVersion.Schema)
	if err != nil {
		logger.Error("Failed to convert current schema", "tablePath", sm.tablePath, "error", err)
		return false, nil, fmt.Errorf("failed to convert current schema: %w", err)
	}
	
	// Compare schemas
	changes, err := sm.CompareSchemas(currentSchema, newSchema)
	if err != nil {
		logger.Error("Failed to compare schemas", "tablePath", sm.tablePath, "error", err)
		return false, nil, fmt.Errorf("failed to compare schemas: %w", err)
	}
	
	// Check for incompatible changes
	compatible := true
	var incompatibleChanges []string
	
	for _, change := range changes {
		switch change.Type {
		case "remove":
			// Removing fields is incompatible
			compatible = false
			incompatibleChanges = append(incompatibleChanges, 
				fmt.Sprintf("Cannot remove field '%s'", change.FieldName))
			
		case "type_change":
			// Type changes are incompatible
			compatible = false
			incompatibleChanges = append(incompatibleChanges, 
				fmt.Sprintf("Cannot change type of field '%s'", change.FieldName))
			
		case "nullability_change":
			// Nullability changes are incompatible
			compatible = false
			incompatibleChanges = append(incompatibleChanges, 
				fmt.Sprintf("Cannot change nullability of field '%s'", change.FieldName))
		}
	}
	
	if compatible {
		logger.Info("Schema is compatible", "tablePath", sm.tablePath)
	} else {
		logger.Warn("Schema is incompatible", "tablePath", sm.tablePath, 
			"incompatibleChanges", incompatibleChanges)
	}
	
	return compatible, incompatibleChanges, nil
}

// GetSchemaChangesBetweenVersions returns the schema changes between two versions
func (sm *DeltaSchemaManager) GetSchemaChangesBetweenVersions(fromVersion, toVersion int64) ([]SchemaChange, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting schema changes between versions", "tablePath", sm.tablePath, 
		"fromVersion", fromVersion, "toVersion", toVersion)
	
	// Validate version numbers
	if fromVersion >= toVersion {
		return nil, fmt.Errorf("fromVersion must be less than toVersion")
	}
	
	// Get schema history
	history, err := sm.GetSchemaHistory()
	if err != nil {
		logger.Error("Failed to get schema history", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to get schema changes: %w", err)
	}
	
	// Find the versions
	var fromSchema, toSchema *SchemaVersion
	for i := range history.Versions {
		version := history.Versions[i].Version
		if version == fromVersion {
			fromSchema = &history.Versions[i]
		}
		if version == toVersion {
			toSchema = &history.Versions[i]
		}
	}
	
	// Check if versions exist
	if fromSchema == nil {
		return nil, fmt.Errorf("schema version %d not found", fromVersion)
	}
	if toSchema == nil {
		return nil, fmt.Errorf("schema version %d not found", toVersion)
	}
	
	// Convert schemas to Schema objects
	fromSchemaObj, err := ConvertFromArrowSchema(fromSchema.Schema)
	if err != nil {
		logger.Error("Failed to convert from schema", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to convert from schema: %w", err)
	}
	
	toSchemaObj, err := ConvertFromArrowSchema(toSchema.Schema)
	if err != nil {
		logger.Error("Failed to convert to schema", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to convert to schema: %w", err)
	}
	
	// Compare schemas
	changes, err := sm.CompareSchemas(fromSchemaObj, toSchemaObj)
	if err != nil {
		logger.Error("Failed to compare schemas", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to compare schemas: %w", err)
	}
	
	logger.Debug("Found schema changes", "tablePath", sm.tablePath, "changes", len(changes))
	return changes, nil
}

// CompareSchemas compares two schemas and returns the changes
func (sm *DeltaSchemaManager) CompareSchemas(oldSchema, newSchema *DeltaSchema) ([]SchemaChange, error) {
	logger := logging.GetLogger()
	logger.Debug("Comparing schemas", "tablePath", sm.tablePath)
	
	if oldSchema == nil || newSchema == nil {
		return nil, fmt.Errorf("cannot compare nil schemas")
	}
	
	// Create maps for quick lookup
	oldFields := make(map[string]DeltaField, len(oldSchema.Fields))
	for _, field := range oldSchema.Fields {
		oldFields[field.Name] = field
	}
	
	newFields := make(map[string]DeltaField, len(newSchema.Fields))
	for _, field := range newSchema.Fields {
		newFields[field.Name] = field
	}
	
	// Track changes
	var changes []SchemaChange
	
	// Check for added fields
	for name, _ := range newFields {
		if _, exists := oldFields[name]; !exists {
			changes = append(changes, SchemaChange{
				Type:      "add",
				FieldName: name,
			})
			logger.Debug("Field added", "tablePath", sm.tablePath, "field", name)
		}
	}
	
	// Check for removed fields
	for name, _ := range oldFields {
		if _, exists := newFields[name]; !exists {
			changes = append(changes, SchemaChange{
				Type:      "remove",
				FieldName: name,
			})
			logger.Debug("Field removed", "tablePath", sm.tablePath, "field", name)
		}
	}
	
	// Check for modified fields
	for name, oldField := range oldFields {
		if newField, exists := newFields[name]; exists {
			// Check if field type changed
			if oldField.Type != newField.Type {
				changes = append(changes, SchemaChange{
					Type:      "type_change",
					FieldName: name,
				})
				logger.Debug("Field type changed", "tablePath", sm.tablePath, "field", name, 
					"oldType", oldField.Type, "newType", newField.Type)
			}
			
			// Check if nullability changed
			if oldField.Nullable != newField.Nullable {
				changes = append(changes, SchemaChange{
					Type:      "nullability_change",
					FieldName: name,
				})
				logger.Debug("Field nullability changed", "tablePath", sm.tablePath, "field", name, 
					"oldNullable", oldField.Nullable, "newNullable", newField.Nullable)
			}
			
			// Check if metadata changed
			if !reflect.DeepEqual(oldField.Metadata, newField.Metadata) {
				changes = append(changes, SchemaChange{
					Type:      "metadata_change",
					FieldName: name,
				})
				logger.Debug("Field metadata changed", "tablePath", sm.tablePath, "field", name)
			}
		}
	}
	
	logger.Debug("Schema comparison complete", "tablePath", sm.tablePath, "changes", len(changes))
	return changes, nil
}

// ClearCache clears all caches in the schema manager
func (sm *DeltaSchemaManager) ClearCache() {
	logger := logging.GetLogger()
	logger.Debug("Clearing schema manager cache", "tablePath", sm.tablePath)
	
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.historyCache = nil
	sm.currentCache = nil
	
	// Clear version cache
	for k := range sm.versionCache {
		delete(sm.versionCache, k)
	}
	
	logger.Debug("Schema manager cache cleared", "tablePath", sm.tablePath)
}

// Helper functions

// readHistory reads the schema history from file
func (sm *DeltaSchemaManager) readHistory() (*SchemaHistory, error) {
	logger := logging.GetLogger()
	logger.Debug("Reading schema history from disk", "tablePath", sm.tablePath)
	
	// Create schema history directory if it doesn't exist
	if err := os.MkdirAll(sm.metaDir, 0755); err != nil {
		logger.Error("Failed to create schema history directory", "tablePath", sm.tablePath, "metaDir", sm.metaDir, "error", err)
		return nil, fmt.Errorf("failed to create schema history directory: %w", err)
	}

	// Check if schema history file exists
	historyPath := filepath.Join(sm.metaDir, "history.json")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		// Return empty history if file doesn't exist
		logger.Debug("Schema history file does not exist, returning empty history", "tablePath", sm.tablePath, "historyPath", historyPath)
		return &SchemaHistory{
			Versions: []SchemaVersion{},
		}, nil
	}

	// Read schema history file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		logger.Error("Failed to read schema history file", "tablePath", sm.tablePath, "historyPath", historyPath, "error", err)
		return nil, fmt.Errorf("failed to read schema history: %w", err)
	}

	// Parse schema history
	var history SchemaHistory
	if err := json.Unmarshal(data, &history); err != nil {
		logger.Error("Failed to parse schema history JSON", "tablePath", sm.tablePath, "error", err)
		return nil, fmt.Errorf("failed to parse schema history: %w", err)
	}

	logger.Debug("Successfully read schema history", "tablePath", sm.tablePath, "versions", len(history.Versions))
	return &history, nil
}

// writeHistory writes the schema history to file
func (sm *DeltaSchemaManager) writeHistory(history *SchemaHistory) error {
	logger := logging.GetLogger()
	logger.Debug("Writing schema history to disk", "tablePath", sm.tablePath, "versions", len(history.Versions))
	
	// Create schema history directory if it doesn't exist
	if err := os.MkdirAll(sm.metaDir, 0755); err != nil {
		logger.Error("Failed to create schema history directory", "tablePath", sm.tablePath, "metaDir", sm.metaDir, "error", err)
		return fmt.Errorf("failed to create schema history directory: %w", err)
	}

	// Marshal schema history
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal schema history to JSON", "tablePath", sm.tablePath, "error", err)
		return fmt.Errorf("failed to marshal schema history: %w", err)
	}

	// Write schema history file
	historyPath := filepath.Join(sm.metaDir, "history.json")
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		logger.Error("Failed to write schema history file", "tablePath", sm.tablePath, "historyPath", historyPath, "error", err)
		return fmt.Errorf("failed to write schema history: %w", err)
	}

	logger.Debug("Successfully wrote schema history", "tablePath", sm.tablePath, "historyPath", historyPath)
	return nil
}
