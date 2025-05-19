package datalake

import (
	"fmt"
	"time"
)

// ConsistencyCheckType defines the type of consistency check to perform
type ConsistencyCheckType string

const (
	// SchemaConsistencyCheck checks for schema consistency between versions
	SchemaConsistencyCheck ConsistencyCheckType = "schema"

	// TypeConsistencyCheck checks for type consistency between versions
	TypeConsistencyCheck ConsistencyCheckType = "type"

	// FieldConsistencyCheck checks for field consistency between versions
	FieldConsistencyCheck ConsistencyCheckType = "field"
)

// ConsistencyCheckResult represents the result of a consistency check
type ConsistencyCheckResult struct {
	// The type of consistency check performed
	CheckType ConsistencyCheckType

	// Whether the check passed
	Passed bool

	// Issues found during the check
	Issues []ConsistencyIssue

	// The versions that were compared
	VersionsCompared []int64

	// The timestamp when the check was performed
	Timestamp time.Time
}

// ConsistencyIssue represents an issue found during a consistency check
type ConsistencyIssue struct {
	// The type of issue
	Type string

	// The field that has the issue
	Field string

	// A description of the issue
	Description string

	// The severity of the issue (warning, error)
	Severity string

	// The versions where the issue was found
	Versions []int64
}

// CheckConsistency performs a consistency check between versions of a Delta table
func (m *MetadataManager) CheckConsistency(checkType ConsistencyCheckType, versions []int64) (*ConsistencyCheckResult, error) {
	if len(versions) < 2 {
		return nil, fmt.Errorf("at least two versions are required for consistency check")
	}

	// Sort versions in ascending order
	sortVersions(versions)

	// Get schema fields for each version
	schemasByVersion := make(map[int64][]SchemaField)
	for _, version := range versions {
		fields, err := m.GetSchemaFieldsAtVersion(version)
		if err != nil {
			return nil, fmt.Errorf("failed to get schema for version %d: %w", version, err)
		}
		schemasByVersion[version] = fields
	}

	// Create result
	result := &ConsistencyCheckResult{
		CheckType:        checkType,
		Passed:           true,
		Issues:           []ConsistencyIssue{},
		VersionsCompared: versions,
		Timestamp:        time.Now(),
	}

	// Perform the appropriate check
	switch checkType {
	case SchemaConsistencyCheck:
		checkSchemaConsistency(schemasByVersion, result)
	case TypeConsistencyCheck:
		checkTypeConsistency(schemasByVersion, result)
	case FieldConsistencyCheck:
		checkFieldConsistency(schemasByVersion, result)
	default:
		return nil, fmt.Errorf("unsupported check type: %s", checkType)
	}

	return result, nil
}

// checkSchemaConsistency checks for schema consistency between versions
func checkSchemaConsistency(schemasByVersion map[int64][]SchemaField, result *ConsistencyCheckResult) {
	// Get all versions in ascending order
	versions := make([]int64, 0, len(schemasByVersion))
	for version := range schemasByVersion {
		versions = append(versions, version)
	}
	sortVersions(versions)

	// Compare schemas between consecutive versions
	for i := 1; i < len(versions); i++ {
		prevVersion := versions[i-1]
		currVersion := versions[i]

		prevSchema := schemasByVersion[prevVersion]
		currSchema := schemasByVersion[currVersion]

		// Check for added fields
		for _, currField := range currSchema {
			found := false
			for _, prevField := range prevSchema {
				if currField.Name == prevField.Name {
					found = true
					break
				}
			}

			if !found {
				// Field was added
				result.Issues = append(result.Issues, ConsistencyIssue{
					Type:        "field_added",
					Field:       currField.Name,
					Description: fmt.Sprintf("Field '%s' was added in version %d", currField.Name, currVersion),
					Severity:    "info",
					Versions:    []int64{currVersion},
				})
			}
		}

		// Check for removed fields
		for _, prevField := range prevSchema {
			found := false
			for _, currField := range currSchema {
				if prevField.Name == currField.Name {
					found = true
					break
				}
			}

			if !found {
				// Field was removed
				result.Issues = append(result.Issues, ConsistencyIssue{
					Type:        "field_removed",
					Field:       prevField.Name,
					Description: fmt.Sprintf("Field '%s' was removed in version %d", prevField.Name, currVersion),
					Severity:    "warning",
					Versions:    []int64{prevVersion, currVersion},
				})
				result.Passed = false
			}
		}
	}
}

// checkTypeConsistency checks for type consistency between versions
func checkTypeConsistency(schemasByVersion map[int64][]SchemaField, result *ConsistencyCheckResult) {
	// Get all versions in ascending order
	versions := make([]int64, 0, len(schemasByVersion))
	for version := range schemasByVersion {
		versions = append(versions, version)
	}
	sortVersions(versions)

	// Create a map of field names to their types in each version
	fieldTypes := make(map[string]map[int64]string)

	// Collect field types for each version
	for _, version := range versions {
		schema := schemasByVersion[version]
		for _, field := range schema {
			if _, ok := fieldTypes[field.Name]; !ok {
				fieldTypes[field.Name] = make(map[int64]string)
			}
			fieldTypes[field.Name][version] = field.Type
		}
	}

	// Check for type changes
	for fieldName, typesByVersion := range fieldTypes {
		if len(typesByVersion) < 2 {
			// Field only exists in one version, already handled by schema consistency check
			continue
		}

		// Check if the type changed across versions
		var prevType string
		var prevVersion int64

		for i, version := range versions {
			fieldType, exists := typesByVersion[version]
			if !exists {
				// Field doesn't exist in this version, already handled by schema consistency check
				continue
			}

			if i > 0 && prevType != "" && fieldType != prevType {
				// Type changed
				result.Issues = append(result.Issues, ConsistencyIssue{
					Type:        "type_changed",
					Field:       fieldName,
					Description: fmt.Sprintf("Field '%s' changed type from '%s' in version %d to '%s' in version %d", fieldName, prevType, prevVersion, fieldType, version),
					Severity:    "error",
					Versions:    []int64{prevVersion, version},
				})
				result.Passed = false
			}

			prevType = fieldType
			prevVersion = version
		}
	}
}

// checkFieldConsistency checks for field consistency between versions
func checkFieldConsistency(schemasByVersion map[int64][]SchemaField, result *ConsistencyCheckResult) {
	// Get all versions in ascending order
	versions := make([]int64, 0, len(schemasByVersion))
	for version := range schemasByVersion {
		versions = append(versions, version)
	}
	sortVersions(versions)

	// Get the latest version schema
	latestVersion := versions[len(versions)-1]
	latestSchema := schemasByVersion[latestVersion]

	// Track which fields have been reported as missing to avoid duplicates
	reportedMissingFields := make(map[string]bool)

	// Check if all fields in the latest version exist in all previous versions
	for _, latestField := range latestSchema {
		// Skip fields that we've already reported as missing
		if reportedMissingFields[latestField.Name] {
			continue
		}

		// Check if this field is new (added in the latest version)
		isNewField := false
		for i := 0; i < len(versions)-1; i++ {
			version := versions[i]
			schema := schemasByVersion[version]

			found := false
			for _, field := range schema {
				if field.Name == latestField.Name {
					found = true
					break
				}
			}

			if !found {
				isNewField = true
				break
			}
		}

		// If this is a new field, report it once
		if isNewField {
			result.Issues = append(result.Issues, ConsistencyIssue{
				Type:        "field_missing",
				Field:       latestField.Name,
				Description: fmt.Sprintf("Field '%s' in latest version %d is not present in all previous versions", latestField.Name, latestVersion),
				Severity:    "warning",
				Versions:    []int64{latestVersion},
			})
			result.Passed = false
			reportedMissingFields[latestField.Name] = true
		}
	}
}

// sortVersions sorts versions in ascending order
func sortVersions(versions []int64) {
	for i := 0; i < len(versions); i++ {
		for j := i + 1; j < len(versions); j++ {
			if versions[i] > versions[j] {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}
}
