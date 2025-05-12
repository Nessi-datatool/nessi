# Nessi.dev Implementation Plan

## Architecture Overview

### Core Architecture
```
nessi/
├── cmd/
│   └── nessi/          # Main CLI application
├── pkg/
│   ├── datalake/       # Delta Lake & multi-format operations (Go)
│   ├── quality/        # Data quality & profiling (Go)
│   │   ├── profile/    # Data profiling & statistics
│   │   ├── rules/      # Rule validation engine
│   │   └── anomaly/    # Anomaly detection algorithms
│   ├── monitoring/     # Metrics & monitoring (Go)
│   │   ├── metrics/    # Metrics collection & storage
│   │   ├── alerts/     # Alerting system
│   │   └── dashboard/  # Web dashboard
│   ├── security/       # Authentication & authorization (Go)
│   │   ├── auth/       # Authentication system
│   │   ├── ssl/        # SSL/TLS support
│   │   └── audit/      # Audit logging
│   └── api/            # API endpoints (Go)
└── scripts/
    ├── delta/          # Delta Lake operations (Python)
    ├── profiling/      # Advanced profiling (Python)
    ├── viz/            # Data visualization (Python)
    └── ml/             # Machine learning for anomaly detection (Python)
```

## Implementation Phases

### Phase 1: Core Infrastructure (2 weeks)
- [x] Basic project structure
- [x] Security implementation
- [x] Monitoring setup
- [x] Go-Python integration layer
  - Implementation: Use `os/exec` to call Python scripts
  - Data exchange via JSON/Arrow
  - Error handling and logging

### Phase 2: Full Schema Management & Version Control (2 weeks)

#### Delta Lake Interface
```go
// pkg/datalake/delta.go
type DeltaTable interface {
    // Core operations
    Initialize(schema *arrow.Schema) error
    Read(partition string, options *ReadOptions) (arrow.Record, error)
    Write(partition string, record arrow.Record) error
    GetSchema() (*arrow.Schema, error)
    GetStats() (*TableStats, error)

    // Schema management
    ValidateSchema(schema *arrow.Schema) error
    EvolveSchemaSafely(newSchema *arrow.Schema) error
    GetSchemaHistory() ([]*SchemaVersion, error)
    GetFieldMetadata(field string) (map[string]interface{}, error)
    UpdateFieldMetadata(field string, metadata map[string]interface{}) error

    // Version control
    GetVersion() (int64, error)
    GetVersionAt(timestamp time.Time) (int64, error)
    GetVersionHistory(limit int) ([]*CommitInfo, error)
    CompareVersions(v1, v2 int64) (*VersionDiff, error)
    Commit(metadata map[string]interface{}) error
    Rollback(version int64) error

    // Time travel
    ReadAsOfVersion(version int64, partition string) (arrow.Record, error)
    ReadAsOfTimestamp(timestamp time.Time, partition string) (arrow.Record, error)
    GetStateAt(timestamp time.Time) (*TableState, error)

    // Metadata operations
    GetMetadata() (map[string]interface{}, error)
    UpdateMetadata(metadata map[string]interface{}) error
}

type ReadOptions struct {
    Filters        []Filter
    Columns        []string
    MaxRowCount    int64
    BatchSize      int64
    UseCache       bool
}

type SchemaVersion struct {
    Version     int64
    Schema      *arrow.Schema
    Timestamp   time.Time
    CommitInfo  *CommitInfo
    Changes     []SchemaChange
}

type SchemaChange struct {
    Type        string // added, removed, modified
    Field       string
    OldType     string
    NewType     string
    Description string
}

type CommitInfo struct {
    Version     int64
    Timestamp   time.Time
    Operation   string
    User        string
    Changes     *ChangeStats
    Message     string
    Metadata    map[string]interface{}
}

type ChangeStats struct {
    FilesAdded      int
    FilesRemoved    int
    RowsAdded       int64
    RowsRemoved     int64
    BytesChanged    int64
}

type VersionDiff struct {
    SchemaChanges   []SchemaChange
    DataChanges     *ChangeStats
    MetadataChanges map[string]interface{}
}

type TableState struct {
    Version     int64
    Schema      *arrow.Schema
    Metadata    map[string]interface{}
    Stats       *TableStats
    Timestamp   time.Time
}
```

#### Python Delta Lake Scripts
```python
# scripts/delta/table.py
from deltalake import DeltaTable, write_deltalake
from pyarrow import Table, Schema
import pandas as pd
from datetime import datetime
import json
import difflib

class DeltaManager:
    def __init__(self, table_path: str):
        self.table_path = table_path
        self.table = DeltaTable(table_path)
    
    # Core operations
    def read_data(self, filters=None, columns=None, limit=None, partition=None) -> Table:
        """Read data with optional filtering and column selection"""
        return self.table.to_pyarrow_table(filters=filters, columns=columns, limit=limit)
    
    def write_records(self, table: Table, partition_cols=None, mode="append") -> None:
        """Write records to the Delta table"""
        write_deltalake(
            self.table_path, 
            table, 
            mode=mode,
            partition_by=partition_cols
        )
        # Refresh table reference after write
        self.table = DeltaTable(self.table_path)
    
    # Schema management
    def get_schema_history(self) -> list:
        """Get the history of schema changes"""
        history = self.table.history()
        schema_versions = []
        
        for i, commit in enumerate(history):
            if i == 0:
                # First commit establishes the schema
                schema_versions.append({
                    'version': commit.version,
                    'timestamp': commit.timestamp,
                    'schema': self.table.schema(commit.version).to_pyarrow(),
                    'changes': [],
                    'commit_info': commit.to_dict()
                })
            else:
                # Compare with previous schema
                current_schema = self.table.schema(commit.version).to_pyarrow()
                prev_schema = self.table.schema(history[i-1].version).to_pyarrow()
                
                # Detect schema changes
                changes = self._detect_schema_changes(prev_schema, current_schema)
                
                if changes:
                    schema_versions.append({
                        'version': commit.version,
                        'timestamp': commit.timestamp,
                        'schema': current_schema,
                        'changes': changes,
                        'commit_info': commit.to_dict()
                    })
        
        return schema_versions
    
    def _detect_schema_changes(self, old_schema, new_schema) -> list:
        """Detect changes between two schemas"""
        changes = []
        old_fields = {field.name: field for field in old_schema}
        new_fields = {field.name: field for field in new_schema}
        
        # Find added fields
        for name, field in new_fields.items():
            if name not in old_fields:
                changes.append({
                    'type': 'added',
                    'field': name,
                    'new_type': str(field.type),
                    'description': f"Added field '{name}' with type {field.type}"
                })
        
        # Find removed fields
        for name, field in old_fields.items():
            if name not in new_fields:
                changes.append({
                    'type': 'removed',
                    'field': name,
                    'old_type': str(field.type),
                    'description': f"Removed field '{name}' with type {field.type}"
                })
        
        # Find modified fields
        for name, old_field in old_fields.items():
            if name in new_fields:
                new_field = new_fields[name]
                if str(old_field.type) != str(new_field.type):
                    changes.append({
                        'type': 'modified',
                        'field': name,
                        'old_type': str(old_field.type),
                        'new_type': str(new_field.type),
                        'description': f"Modified field '{name}' from {old_field.type} to {new_field.type}"
                    })
        
        return changes
    
    def validate_schema(self, schema: Schema) -> dict:
        """Validate if a schema is compatible with the current schema"""
        current_schema = self.table.schema().to_pyarrow()
        changes = self._detect_schema_changes(current_schema, schema)
        
        # Check for breaking changes
        breaking_changes = []
        for change in changes:
            if change['type'] == 'removed':
                breaking_changes.append(change)
            elif change['type'] == 'modified':
                # Check if type conversion is safe
                if not self._is_safe_type_conversion(change['old_type'], change['new_type']):
                    breaking_changes.append(change)
        
        return {
            'is_valid': len(breaking_changes) == 0,
            'changes': changes,
            'breaking_changes': breaking_changes
        }
    
    def _is_safe_type_conversion(self, old_type, new_type) -> bool:
        """Check if a type conversion is safe (non-breaking)"""
        # Safe conversions (examples)
        safe_conversions = {
            'int32': ['int64', 'float32', 'float64'],
            'int64': ['float64'],
            'float32': ['float64'],
            'string': ['large_string']
        }
        
        if old_type in safe_conversions and new_type in safe_conversions[old_type]:
            return True
        return old_type == new_type
    
    # Version control
    def get_version_history(self, limit=None) -> list:
        """Get the history of versions with detailed commit information"""
        history = self.table.history(limit=limit)
        return [commit.to_dict() for commit in history]
    
    def compare_versions(self, v1, v2) -> dict:
        """Compare two versions of the table"""
        # Get commit info
        history = self.table.history()
        commit1 = next((c for c in history if c.version == v1), None)
        commit2 = next((c for c in history if c.version == v2), None)
        
        if not commit1 or not commit2:
            raise ValueError(f"Version {v1 if not commit1 else v2} not found")
        
        # Compare schemas
        schema1 = self.table.schema(v1).to_pyarrow()
        schema2 = self.table.schema(v2).to_pyarrow()
        schema_changes = self._detect_schema_changes(schema1, schema2)
        
        # Compare data stats (if available)
        stats1 = self._get_version_stats(v1)
        stats2 = self._get_version_stats(v2)
        
        # Compare metadata
        metadata1 = self.table.metadata(v1)
        metadata2 = self.table.metadata(v2)
        metadata_changes = self._dict_diff(metadata1, metadata2)
        
        return {
            'schema_changes': schema_changes,
            'data_changes': {
                'files_added': len(stats2.get('files', [])) - len(stats1.get('files', [])),
                'rows_added': stats2.get('num_records', 0) - stats1.get('num_records', 0),
                'bytes_changed': stats2.get('size_bytes', 0) - stats1.get('size_bytes', 0)
            },
            'metadata_changes': metadata_changes,
            'v1_commit': commit1.to_dict(),
            'v2_commit': commit2.to_dict()
        }
    
    def _get_version_stats(self, version) -> dict:
        """Get statistics for a specific version"""
        try:
            # This is a simplified approach - actual implementation would depend on
            # what statistics are available in your Delta Lake implementation
            files = self.table.files(version)
            return {
                'files': files,
                'num_records': sum(self._get_file_stats(f).get('num_records', 0) for f in files),
                'size_bytes': sum(self._get_file_stats(f).get('size_bytes', 0) for f in files)
            }
        except:
            return {}
    
    def _get_file_stats(self, file_path) -> dict:
        """Get statistics for a specific file"""
        # This would be implemented based on how your Delta Lake implementation
        # provides file-level statistics
        return {}
    
    def _dict_diff(self, d1, d2) -> dict:
        """Compare two dictionaries and return the differences"""
        diff = {}
        all_keys = set(d1.keys()) | set(d2.keys())
        
        for key in all_keys:
            if key not in d1:
                diff[key] = {'added': d2[key]}
            elif key not in d2:
                diff[key] = {'removed': d1[key]}
            elif d1[key] != d2[key]:
                diff[key] = {'old': d1[key], 'new': d2[key]}
        
        return diff
    
    # Time travel
    def read_as_of_version(self, version, filters=None, columns=None) -> Table:
        """Read data as of a specific version"""
        return self.table.to_pyarrow_table(version=version, filters=filters, columns=columns)
    
    def read_as_of_timestamp(self, timestamp, filters=None, columns=None) -> Table:
        """Read data as of a specific timestamp"""
        # Find the version at or before the timestamp
        history = self.table.history()
        target_version = None
        
        for commit in history:
            if commit.timestamp <= timestamp:
                target_version = commit.version
                break
        
        if target_version is None:
            raise ValueError(f"No version found at or before timestamp {timestamp}")
        
        return self.read_as_of_version(target_version, filters, columns)
    
    def get_state_at(self, timestamp) -> dict:
        """Get the complete state of the table at a specific timestamp"""
        # Find the version at or before the timestamp
        history = self.table.history()
        target_version = None
        commit_info = None
        
        for commit in history:
            if commit.timestamp <= timestamp:
                target_version = commit.version
                commit_info = commit
                break
        
        if target_version is None:
            raise ValueError(f"No version found at or before timestamp {timestamp}")
        
        return {
            'version': target_version,
            'schema': self.table.schema(target_version).to_pyarrow(),
            'metadata': self.table.metadata(target_version),
            'stats': self._get_version_stats(target_version),
            'timestamp': commit_info.timestamp,
            'commit_info': commit_info.to_dict()
        }
    
    # Metadata operations
    def get_metadata(self, version=None) -> dict:
        """Get metadata for a specific version or the latest version"""
        if version is not None:
            return self.table.metadata(version)
        return self.table.metadata()
    
    def update_metadata(self, metadata: dict) -> None:
        """Update metadata for the table"""
        # This would typically involve a write operation with metadata
        # The exact implementation depends on the Delta Lake library
        current_metadata = self.table.metadata()
        merged_metadata = {**current_metadata, **metadata}
        
        # In a real implementation, you would write this metadata back to the table
        # This is a simplified placeholder
        print(f"Updating metadata: {json.dumps(merged_metadata, indent=2)}")
        
    def get_field_metadata(self, field: str) -> dict:
        """Get metadata for a specific field"""
        schema = self.table.schema().to_pyarrow()
        for f in schema:
            if f.name == field:
                return f.metadata or {}
        raise ValueError(f"Field '{field}' not found in schema")
    
    def update_field_metadata(self, field: str, metadata: dict) -> None:
        """Update metadata for a specific field"""
        # This would typically involve a schema evolution operation
        # The exact implementation depends on the Delta Lake library
        print(f"Updating metadata for field '{field}': {json.dumps(metadata, indent=2)}")

```

### Phase 3: Data Quality & Profiling (2 weeks)

#### Enhanced Rule Engine
```go
// pkg/quality/rules/engine.go
type Rule interface {
    Validate(record arrow.Record) []ValidationResult
    GetMetadata() RuleMetadata
    IsEnabled() bool
    SetEnabled(enabled bool)
    GetHistory() []RuleExecutionRecord
}

type RuleMetadata struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Severity    string                 `json:"severity"`
    RuleType    string                 `json:"rule_type"`
    Field       string                 `json:"field,omitempty"`
    Tags        []string               `json:"tags"`
    Config      map[string]interface{} `json:"config"`
    CreatedAt   time.Time             `json:"created_at"`
    UpdatedAt   time.Time             `json:"updated_at"`
    CreatedBy   string                 `json:"created_by,omitempty"`
}

type ValidationResult struct {
    RuleID      string      `json:"rule_id"`
    Field       string      `json:"field"`
    Value       interface{} `json:"value"`
    Message     string      `json:"message"`
    RowIndex    int64       `json:"row_index"`
    Passed      bool        `json:"passed"`
    Timestamp   time.Time   `json:"timestamp"`
    ExecutionID string      `json:"execution_id"`
}

type RuleExecutionRecord struct {
    ExecutionID   string    `json:"execution_id"`
    RuleID        string    `json:"rule_id"`
    Timestamp     time.Time `json:"timestamp"`
    RecordsChecked int64     `json:"records_checked"`
    Failures      int64     `json:"failures"`
    ExecutionTimeMs int64   `json:"execution_time_ms"`
    DatasetID     string    `json:"dataset_id,omitempty"`
}

type RuleManager struct {
    rules       map[string]Rule
    storage     RuleStorage
    executions  ExecutionStorage
    validators  map[string]RuleValidator
}

type RuleStorage interface {
    GetRule(id string) (Rule, error)
    GetRules() ([]Rule, error)
    SaveRule(rule Rule) error
    DeleteRule(id string) error
    GetRulesByTag(tag string) ([]Rule, error)
    GetRulesByType(ruleType string) ([]Rule, error)
}

type ExecutionStorage interface {
    SaveExecution(record RuleExecutionRecord) error
    GetExecutions(ruleID string, limit int) ([]RuleExecutionRecord, error)
    GetExecutionsByDataset(datasetID string, limit int) ([]RuleExecutionRecord, error)
    GetExecutionTrends(ruleID string, days int) ([]ExecutionTrend, error)
}

type ExecutionTrend struct {
    Date            time.Time `json:"date"`
    ExecutionCount  int      `json:"execution_count"`
    FailureRate     float64  `json:"failure_rate"`
    AverageFailures float64  `json:"average_failures"`
}

type RuleValidator interface {
    ValidateRule(rule Rule) error
    ValidateConfig(config map[string]interface{}, ruleType string) error
    GetSupportedRuleTypes() []string
}
```

#### Extended Rule Types
```go
// pkg/quality/rules/types.go
type RangeRule struct {
    metadata RuleMetadata
    min      float64
    max      float64
    enabled  bool
    history  []RuleExecutionRecord
}

type RegexRule struct {
    metadata RuleMetadata
    pattern  *regexp.Regexp
    enabled  bool
    history  []RuleExecutionRecord
}

type EnumRule struct {
    metadata   RuleMetadata
    allowedValues []interface{}
    caseSensitive bool
    enabled    bool
    history    []RuleExecutionRecord
}

type LengthRule struct {
    metadata RuleMetadata
    minLength int
    maxLength int
    enabled   bool
    history   []RuleExecutionRecord
}

type DateFormatRule struct {
    metadata  RuleMetadata
    format    string
    enabled   bool
    history   []RuleExecutionRecord
}

type CustomRule struct {
    metadata  RuleMetadata
    script    string
    language  string // "python", "sql", etc.
    enabled   bool
    history   []RuleExecutionRecord
}
```

#### YAML Rule Configuration
```yaml
# config/rules.yaml
rules:
  - id: null_check
    name: Null Value Check
    description: Checks for null values in required fields
    severity: error
    rule_type: null_check
    tags: [data_quality, nulls]
    config:
      fields: [id, name, value]

  - id: range_check
    name: Value Range Check
    description: Validates numeric values are within expected ranges
    severity: warning
    rule_type: range
    tags: [data_quality, range]
    config:
      fields:
        age: {min: 0, max: 120}
        score: {min: 0, max: 100}
        
  - id: email_format
    name: Email Format Check
    description: Validates email addresses have correct format
    severity: warning
    rule_type: regex
    field: email
    tags: [data_quality, format]
    config:
      pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
      
  - id: country_code
    name: Country Code Check
    description: Validates country codes are from allowed list
    severity: error
    rule_type: enum
    field: country
    tags: [data_quality, reference]
    config:
      values: ["US", "CA", "UK", "AU", "DE", "FR", "JP"]
      case_sensitive: true
      
  - id: username_length
    name: Username Length Check
    description: Validates username length is between 3 and 20 characters
    severity: error
    rule_type: length
    field: username
    tags: [data_quality, length]
    config:
      min_length: 3
      max_length: 20
      
  - id: date_format
    name: Date Format Check
    description: Validates dates are in ISO format
    severity: warning
    rule_type: date_format
    field: created_at
    tags: [data_quality, date]
    config:
      format: "2006-01-02T15:04:05Z07:00"
      
  - id: custom_validation
    name: Custom Python Validation
    description: Custom validation using Python script
    severity: warning
    rule_type: custom
    tags: [data_quality, custom]
    config:
      language: python
      script: |
        def validate(record):
            # Custom validation logic
            if record['total'] != sum(record['items']):
                return False, "Total doesn't match sum of items"
            return True, ""
```

### Phase 4: Advanced Anomaly Detection (2 weeks)

#### Anomaly Detection Engine
```go
// pkg/quality/anomaly/detector.go
type AnomalyDetector interface {
    DetectOutliers(data []float64, options OutlierDetectionOptions) []Outlier
    DetectSuddenChanges(timeSeries []TimeSeriesPoint, options ChangeDetectionOptions) []ChangePoint
    DetectPatterns(timeSeries []TimeSeriesPoint) []Pattern
}

type OutlierDetectionOptions struct {
    Method          string  `json:"method"` // "zscore" or "iqr"
    ZScoreThreshold float64 `json:"zscore_threshold,omitempty"`
    IQRMultiplier   float64 `json:"iqr_multiplier,omitempty"`
}

type ChangeDetectionOptions struct {
    WindowSize       int     `json:"window_size"`
    Threshold        float64 `json:"threshold"`
    MinConsecutive   int     `json:"min_consecutive"`
    DetectConstant   bool    `json:"detect_constant"`
    DetectSpikes     bool    `json:"detect_spikes"`
    DetectOscillation bool    `json:"detect_oscillation"`
}

type Outlier struct {
    Value       float64   `json:"value"`
    Index       int      `json:"index"`
    Score       float64   `json:"score"`
    Method      string    `json:"method"`
    Timestamp   time.Time `json:"timestamp,omitempty"`
}

type TimeSeriesPoint struct {
    Value     float64   `json:"value"`
    Timestamp time.Time `json:"timestamp"`
}

type ChangePoint struct {
    Index       int       `json:"index"`
    Timestamp   time.Time `json:"timestamp"`
    OldValue    float64   `json:"old_value"`
    NewValue    float64   `json:"new_value"`
    ChangeType  string    `json:"change_type"` // "spike", "drop", "level_shift", "trend_change"
    Magnitude   float64   `json:"magnitude"`
    Confidence  float64   `json:"confidence"`
}

type Pattern struct {
    Type        string    `json:"type"` // "constant", "increasing", "decreasing", "oscillating"
    StartIndex  int       `json:"start_index"`
    EndIndex    int       `json:"end_index"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Magnitude   float64   `json:"magnitude"`
    Period      float64   `json:"period,omitempty"` // For oscillating patterns
}
```

#### Outlier Detection Implementation
```go
// pkg/quality/anomaly/outlier.go
func (d *DefaultAnomalyDetector) DetectOutliersZScore(data []float64, threshold float64) []Outlier {
    // Calculate mean and standard deviation
    mean := calculateMean(data)
    stdDev := calculateStdDev(data, mean)
    
    outliers := []Outlier{}
    
    // Detect outliers using Z-score
    for i, value := range data {
        zScore := math.Abs((value - mean) / stdDev)
        if zScore > threshold {
            outliers = append(outliers, Outlier{
                Value:  value,
                Index:  i,
                Score:  zScore,
                Method: "zscore",
            })
        }
    }
    
    return outliers
}

func (d *DefaultAnomalyDetector) DetectOutliersIQR(data []float64, multiplier float64) []Outlier {
    // Sort data to calculate quartiles
    sortedData := make([]float64, len(data))
    copy(sortedData, data)
    sort.Float64s(sortedData)
    
    // Calculate Q1 and Q3
    q1 := calculatePercentile(sortedData, 25)
    q3 := calculatePercentile(sortedData, 75)
    
    // Calculate IQR and bounds
    iqr := q3 - q1
    lowerBound := q1 - (multiplier * iqr)
    upperBound := q3 + (multiplier * iqr)
    
    outliers := []Outlier{}
    
    // Detect outliers using IQR
    for i, value := range data {
        if value < lowerBound || value > upperBound {
            // Calculate how many IQRs away from the bounds
            var score float64
            if value < lowerBound {
                score = math.Abs((lowerBound - value) / iqr)
            } else {
                score = math.Abs((value - upperBound) / iqr)
            }
            
            outliers = append(outliers, Outlier{
                Value:  value,
                Index:  i,
                Score:  score,
                Method: "iqr",
            })
        }
    }
    
    return outliers
}
```

#### Sudden Change Detection Implementation
```go
// pkg/quality/anomaly/change.go
func (d *DefaultAnomalyDetector) DetectSuddenChanges(timeSeries []TimeSeriesPoint, options ChangeDetectionOptions) []ChangePoint {
    if len(timeSeries) < options.WindowSize*2 {
        return []ChangePoint{}
    }
    
    changes := []ChangePoint{}
    
    // Sliding window approach
    for i := options.WindowSize; i < len(timeSeries)-options.WindowSize; i++ {
        // Get windows before and after current point
        before := getValues(timeSeries[i-options.WindowSize:i])
        after := getValues(timeSeries[i:i+options.WindowSize])
        
        // Calculate statistics for both windows
        beforeMean := calculateMean(before)
        afterMean := calculateMean(after)
        beforeStd := calculateStdDev(before, beforeMean)
        afterStd := calculateStdDev(after, afterMean)
        
        // Calculate change magnitude
        meanDiff := math.Abs(afterMean - beforeMean)
        stdDiff := math.Abs(afterStd - beforeStd)
        
        // Detect level shifts (mean changes significantly)
        if meanDiff > options.Threshold * math.Max(beforeStd, afterStd) {
            // Check if this is part of a consistent pattern
            consecutive := 1
            for j := 1; j < options.MinConsecutive && i+j < len(timeSeries)-options.WindowSize; j++ {
                nextAfter := getValues(timeSeries[i+j:i+j+options.WindowSize])
                nextAfterMean := calculateMean(nextAfter)
                
                if (afterMean > beforeMean && nextAfterMean > afterMean) || 
                   (afterMean < beforeMean && nextAfterMean < afterMean) {
                    consecutive++
                }
            }
            
            // Determine change type
            changeType := "level_shift"
            if consecutive >= options.MinConsecutive && options.DetectConstant {
                if afterMean > beforeMean {
                    changeType = "increasing"
                } else {
                    changeType = "decreasing"
                }
            } else if stdDiff > options.Threshold * beforeStd && options.DetectSpikes {
                changeType = "spike"
            }
            
            changes = append(changes, ChangePoint{
                Index:      i,
                Timestamp:  timeSeries[i].Timestamp,
                OldValue:   beforeMean,
                NewValue:   afterMean,
                ChangeType: changeType,
                Magnitude:  meanDiff / beforeMean, // Relative change
                Confidence: calculateConfidence(meanDiff, beforeStd, afterStd),
            })
        }
        
        // Detect oscillations if enabled
        if options.DetectOscillation {
            oscillation := detectOscillation(timeSeries, i, options.WindowSize*2)
            if oscillation != nil {
                changes = append(changes, *oscillation)
            }
        }
    }
    
    return changes
}

func detectOscillation(timeSeries []TimeSeriesPoint, centerIndex, windowSize int) *ChangePoint {
    // Implementation of oscillation detection algorithm
    // This would analyze frequency components and detect regular patterns
    // Simplified placeholder implementation
    return nil
}
```

#### Advanced Profiling Integration
```python
# scripts/profiling/profile.py
from ydata_profiling import ProfileReport
import pyarrow as pa
import pandas as pd
import numpy as np
from scipy import stats
from statsmodels.tsa.seasonal import seasonal_decompose

def generate_profile(table: pa.Table, config: dict) -> dict:
    # Convert Arrow table to Pandas
    df = table.to_pandas()
    
    # Generate profile
    profile = ProfileReport(
        df,
        title=config.get('title', 'Data Profile'),
        minimal=config.get('minimal', False),
        explorative=config.get('explorative', True)
    )
    
    # Get profile data
    return {
        'summary': profile.get_description(),
        'warnings': profile.get_warnings(),
        'correlations': profile.get_correlations(),
        'missing': profile.missing_diagrams,
        'stats': profile.calculate_stats()
    }
```

### Phase 4: Visualization & Reporting (1 week)

#### Report Templates
```go
// pkg/reporting/templates.go
type ReportTemplate struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Format      string                 `json:"format"` // html, pdf, md
    Variables   map[string]string      `json:"variables"`
    Sections    []ReportSection        `json:"sections"`
    Assets      map[string]AssetConfig `json:"assets"`
}

type ReportSection struct {
    Title      string       `json:"title"`
    Type       string       `json:"type"` // text, chart, table, profile
    Content    string       `json:"content"`
    ChartType  string       `json:"chart_type,omitempty"`
    DataSource DataSource   `json:"data_source"`
    Options    interface{}  `json:"options,omitempty"`
}

type AssetConfig struct {
    Type     string `json:"type"` // font, image, style
    Path     string `json:"path"`
    MimeType string `json:"mime_type"`
}
```

#### Visualization Engine
```python
# scripts/viz/charts.py
import plotly.graph_objects as go
import plotly.express as px
from plotly.subplots import make_subplots

def create_chart(data: dict, config: dict) -> dict:
    chart_type = config['type']
    
    if chart_type == 'distribution':
        fig = px.histogram(
            data['values'],
            title=config.get('title'),
            nbins=config.get('bins', 50),
            marginal=config.get('marginal', 'box')
        )
    elif chart_type == 'correlation':
        fig = px.scatter_matrix(
            data['frame'],
            dimensions=config.get('columns'),
            title=config.get('title')
        )
    elif chart_type == 'time_series':
        fig = px.line(
            data['frame'],
            x=config['x'],
            y=config['y'],
            title=config.get('title')
        )
    
    # Apply theme and layout options
    fig.update_layout(
        template=config.get('template', 'plotly_white'),
        **config.get('layout', {})
    )
    
    return {
        'plot': fig.to_json(),
        'image': fig.to_image(format='png') if config.get('export_image') else None
    }
```

#### Report Generation
```go
// pkg/reporting/generator.go
type ReportGenerator struct {
    templates map[string]*ReportTemplate
    assets   *AssetManager
}

func (g *ReportGenerator) Generate(templateID string, data map[string]interface{}) (*Report, error) {
    template := g.templates[templateID]
    if template == nil {
        return nil, fmt.Errorf("template not found: %s", templateID)
    }
    
    report := &Report{
        ID:        uuid.New().String(),
        Template:  template,
        Data:      data,
        CreatedAt: time.Now(),
    }
    
    // Generate sections
    for _, section := range template.Sections {
        content, err := g.generateSection(section, data)
        if err != nil {
            return nil, fmt.Errorf("failed to generate section %s: %w", section.Title, err)
        }
        report.Sections = append(report.Sections, content)
    }
    
    // Apply template
    if err := g.applyTemplate(report); err != nil {
        return nil, fmt.Errorf("failed to apply template: %w", err)
    }
    
    return report, nil
}
```

## Integration Details

### Go-Python Integration

#### Process Pool Management
```go
// pkg/common/python.go
type PythonProcess struct {
    Cmd    *exec.Cmd
    Stdin  io.WriteCloser
    Stdout io.ReadCloser
    Stderr io.ReadCloser
}

type ProcessPool struct {
    processes chan *PythonProcess
    size      int
}

func NewProcessPool(size int) *ProcessPool {
    pool := &ProcessPool{
        processes: make(chan *PythonProcess, size),
        size:      size,
    }
    pool.initialize()
    return pool
}

func (p *ProcessPool) Execute(data []byte) ([]byte, error) {
    process := <-p.processes
    defer func() { p.processes <- process }()

    // Send data via stdin
    if _, err := process.Stdin.Write(data); err != nil {
        return nil, fmt.Errorf("failed to write to stdin: %w", err)
    }

    // Read result with timeout
    resultCh := make(chan []byte)
    errCh := make(chan error)
    go func() {
        result, err := io.ReadAll(process.Stdout)
        if err != nil {
            errCh <- err
            return
        }
        resultCh <- result
    }()

    select {
    case result := <-resultCh:
        return result, nil
    case err := <-errCh:
        return nil, err
    case <-time.After(30 * time.Second):
        return nil, fmt.Errorf("execution timed out")
    }
}
```

#### Error Handling
```go
type PythonError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Stack   string `json:"stack,omitempty"`
}

func (e *PythonError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (p *ProcessPool) handleError(stderr io.Reader) error {
    errOutput, err := io.ReadAll(stderr)
    if err != nil {
        return fmt.Errorf("failed to read stderr: %w", err)
    }

    var pyErr PythonError
    if err := json.Unmarshal(errOutput, &pyErr); err != nil {
        return fmt.Errorf("python error: %s", string(errOutput))
    }
    return &pyErr
}
```

### Data Exchange Format

#### Data Exchange Strategy
```go
// pkg/common/exchange.go
type ExchangeStrategy interface {
    Serialize(data interface{}) ([]byte, error)
    Deserialize(data []byte, target interface{}) error
    SupportsType(data interface{}) bool
}

type ExchangeManager struct {
    strategies []ExchangeStrategy
    config    ExchangeConfig
}

type ExchangeConfig struct {
    // Size thresholds in bytes
    ArrowThreshold int64 `json:"arrow_threshold"`
    // Memory limits
    MaxMemoryUsage int64 `json:"max_memory_usage"`
    // Compression settings
    UseCompression bool `json:"use_compression"`
}

// Arrow-based Exchange for Large Datasets
type ArrowExchange struct {
    Schema  *arrow.Schema
    Records []arrow.Record
    stats   ExchangeStats
}

type ExchangeStats struct {
    RowCount     int64
    SizeBytes    int64
    MemoryUsage  int64
    SerializeMs  int64
    DeserializeMs int64
}

func (a *ArrowExchange) ToIPC() ([]byte, error) {
    start := time.Now()
    defer func() {
        a.stats.SerializeMs = time.Since(start).Milliseconds()
    }()

    // Check memory limits
    if err := a.checkMemoryUsage(); err != nil {
        return nil, fmt.Errorf("memory limit exceeded: %w", err)
    }

    buf := new(bytes.Buffer)
    writer := ipc.NewWriter(buf, ipc.WithSchema(a.Schema))
    defer writer.Close()

    for _, record := range a.Records {
        if err := writer.Write(record); err != nil {
            return nil, fmt.Errorf("failed to write record: %w", err)
        }
        a.stats.RowCount += record.NumRows()
    }

    data := buf.Bytes()
    a.stats.SizeBytes = int64(len(data))
    return data, nil
}

func (a *ArrowExchange) FromIPC(data []byte) error {
    start := time.Now()
    defer func() {
        a.stats.DeserializeMs = time.Since(start).Milliseconds()
    }()

    if len(data) == 0 {
        return fmt.Errorf("empty data received")
    }

    reader := ipc.NewReader(bytes.NewReader(data))
    defer reader.Close()

    if reader.Schema() == nil {
        return fmt.Errorf("invalid Arrow IPC format: missing schema")
    }

    var records []arrow.Record
    for reader.Next() {
        record := reader.Record()
        records = append(records, record)
    }

    if err := reader.Err(); err != nil {
        return fmt.Errorf("error reading Arrow IPC: %w", err)
    }

    a.Schema = reader.Schema()
    a.Records = records
    return nil
}
```

#### JSON Exchange (Small Datasets)
```go
// pkg/common/json.go
type JSONExchange struct {
    stats ExchangeStats
}

type JSONData struct {
    Status string          `json:"status"`
    Data   JSONDataContent `json:"data"`
    Error  *JSONError      `json:"error,omitempty"`
}

type JSONDataContent struct {
    Records  []map[string]interface{} `json:"records"`
    Metadata JSONMetadata            `json:"metadata"`
    Schema   JSONSchema              `json:"schema"`
}

type JSONMetadata struct {
    RowCount   int       `json:"row_count"`
    Timestamp  time.Time `json:"timestamp"`
    Version    string    `json:"version"`
    Source     string    `json:"source"`
    Properties map[string]interface{} `json:"properties,omitempty"`
}

type JSONError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

## Testing Strategy

### Unit Tests
- Go packages: Standard Go testing with testify assertions
- Python scripts: pytest with coverage and hypothesis
- Integration tests: End-to-end workflows using Docker Compose

### Performance Tests
- Large dataset handling (100GB+)
  - Memory usage profiling
  - CPU profiling
  - I/O profiling
- Python script execution overhead
  - Process pool efficiency
  - Data serialization benchmarks
  - Arrow vs JSON comparison
- Memory usage monitoring
  - Resource limits enforcement
  - Memory leak detection
  - Garbage collection metrics

### Load Testing
- Concurrent request handling
- Process pool scaling
- Resource utilization under load

### Security Testing
- Input validation
- Resource isolation
- Timeout enforcement
- Error propagation

## Dependencies

### Go Dependencies
```go
require (
    "github.com/prometheus/client_golang"
    "github.com/spf13/cobra"
    "github.com/stretchr/testify"
    // Add other Go dependencies
)
```

### Python Dependencies
```python
deltalake>=0.15.0
pyarrow>=14.0.1
ydata-profiling>=4.5.1
pandas>=2.0.0
plotly>=5.18.0
kaleido>=0.2.1
```

## Deployment Requirements

### Container-based Deployment
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o nessi ./cmd/nessi

FROM python:3.9-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    libgomp1 \
    && rm -rf /var/lib/apt/lists/*

# Create Python virtual environment
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy Go binary and scripts
COPY --from=builder /app/nessi /usr/local/bin/
COPY scripts/ /app/scripts/
COPY config/ /app/config/

# Set environment variables
ENV NESSI_CONFIG=/app/config
ENV NESSI_SCRIPTS=/app/scripts
ENV PYTHONPATH=/app/scripts

# Set resource limits
ENV NESSI_MAX_MEMORY=4G
ENV NESSI_MAX_CPU=2

CMD ["nessi", "serve"]
```

### Docker Compose Setup
```yaml
# docker-compose.yml
version: '3.8'

services:
  nessi:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - NESSI_ENV=production
      - NESSI_LOG_LEVEL=info
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '1'
          memory: 2G
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Rule Validation Implementation
```go
// pkg/quality/rules/validator.go
type RuleValidator struct {
    rules []Rule
    stats ValidationStats
}

type ValidationStats struct {
    TotalRecords   int64
    ValidRecords   int64
    InvalidRecords int64
    ErrorsByRule   map[string]int64
    StartTime      time.Time
    EndTime        time.Time
}

// Example rule implementation
type NullCheckRule struct {
    fields []string
    config RuleConfig
}

func (r *NullCheckRule) Validate(record arrow.Record) []ValidationError {
    var errors []ValidationError
    
    for _, field := range r.fields {
        colIdx := record.Schema().FieldIndices(field)[0]
        if colIdx < 0 {
            continue
        }
        
        col := record.Column(colIdx)
        for i := 0; i < int(record.NumRows()); i++ {
            if col.IsNull(i) {
                errors = append(errors, ValidationError{
                    RuleID:   "null_check",
                    Field:    field,
                    Value:    "null",
                    Message:  fmt.Sprintf("Field '%s' cannot be null", field),
                    RowIndex: int64(i),
                })
            }
        }
    }
    
    return errors
}

// Example range check rule
type RangeCheckRule struct {
    fields map[string]RangeConfig
    config RuleConfig
}

type RangeConfig struct {
    Min float64 `json:"min"`
    Max float64 `json:"max"`
}

func (r *RangeCheckRule) Validate(record arrow.Record) []ValidationError {
    var errors []ValidationError
    
    for field, rangeConfig := range r.fields {
        colIdx := record.Schema().FieldIndices(field)[0]
        if colIdx < 0 {
            continue
        }
        
        col := record.Column(colIdx)
        for i := 0; i < int(record.NumRows()); i++ {
            if col.IsNull(i) {
                continue
            }
            
            value := col.(*array.Float64).Value(i)
            if value < rangeConfig.Min || value > rangeConfig.Max {
                errors = append(errors, ValidationError{
                    RuleID:   "range_check",
                    Field:    field,
                    Value:    fmt.Sprintf("%f", value),
                    Message:  fmt.Sprintf("Value %f is outside range [%f, %f]",
                                         value, rangeConfig.Min, rangeConfig.Max),
                    RowIndex: int64(i),
                })
            }
        }
    }
    
    return errors
}
```

## Implementation Decisions

### Data Exchange Strategy
- **Arrow vs. JSON Selection:**
  - Use Arrow for datasets > 1MB (configurable via `ArrowThreshold`)
  - Use JSON for smaller datasets and control messages
  - Track memory usage and enforce limits
  - Compress data when `UseCompression` is enabled

### Error Handling
- **Arrow Operations:**
  - Validate schema before serialization/deserialization
  - Track operation timing and memory usage
  - Provide detailed error messages with context
  - Implement graceful fallback to JSON if needed

### Rule Validation
- **Implementation Pattern:**
  - Each rule is a separate type implementing the `Rule` interface
  - Rules are configured via YAML files
  - Rules operate on Arrow records for efficiency
  - Track validation stats for monitoring

### Report Generation
- **Component Responsibilities:**
  - `ReportGenerator`: Orchestrates report creation
  - `ReportSection`: Handles specific content types
  - `AssetManager`: Manages resources and templates
  - Support multiple output formats (HTML, PDF, MD)

### Deployment Strategy
- **Container-based Deployment:**
  - Multi-stage builds for minimal image size
  - Python environment isolation
  - Resource limits and monitoring
  - Health checks and auto-recovery

## Timeline and Milestones

### Week 1-2: Core Infrastructure
- Project setup
- Basic functionality
- Integration layer

### Week 3: Delta Lake Integration
- Python scripts
- Go wrappers
- Basic operations

### Week 4: Quality & Profiling
- Data profiling
- Quality rules
- Metrics collection
- Anomaly detection (z-score, IQR, sudden change detection)

### Week 5: Visualization & Polish
- Reporting
- Documentation
- Testing

## Progress Update (May 2025)

### Completed Features
- **IQR-based Outlier Detection**: Implemented a more robust outlier detection method using Interquartile Range (IQR), which works better for skewed distributions compared to z-score method.
- **Sudden Change Detection**: Added capability to detect various patterns of changes in metrics over time:
  - Constant changes: Gradual increases or decreases over multiple runs
  - Sudden spikes/dips: Temporary anomalies that return to normal
  - Oscillations: Alternating patterns in metric values

### Phase 5: Advanced Monitoring & Alerting (2 weeks)

#### Monitoring System
```go
// pkg/monitoring/metrics/collector.go
type MetricsCollector interface {
    RecordMetric(metric *Metric) error
    GetMetrics(options *MetricQuery) ([]*Metric, error)
    GetMetricHistory(name string, options *HistoryOptions) ([]*MetricPoint, error)
    GetMetricStats(name string) (*MetricStats, error)
    RegisterMetricSource(source MetricSource) error
    GetDashboardData(timeRange *TimeRange) (*DashboardData, error)
}

type Metric struct {
    Name        string                 `json:"name"`
    Value       float64                `json:"value"`
    Timestamp   time.Time              `json:"timestamp"`
    Tags        map[string]string      `json:"tags"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    Source      string                 `json:"source"`
}

type MetricPoint struct {
    Value       float64   `json:"value"`
    Timestamp   time.Time `json:"timestamp"`
}

type MetricQuery struct {
    Names       []string               `json:"names,omitempty"`
    Tags        map[string]string      `json:"tags,omitempty"`
    TimeRange   *TimeRange             `json:"time_range,omitempty"`
    Limit       int                    `json:"limit,omitempty"`
    Aggregation string                 `json:"aggregation,omitempty"` // sum, avg, min, max
    GroupBy     []string               `json:"group_by,omitempty"`
}

type TimeRange struct {
    Start       time.Time `json:"start"`
    End         time.Time `json:"end"`
    Resolution  string    `json:"resolution,omitempty"` // 1m, 5m, 1h, 1d
}

type MetricStats struct {
    Count       int64     `json:"count"`
    Min         float64   `json:"min"`
    Max         float64   `json:"max"`
    Mean        float64   `json:"mean"`
    Median      float64   `json:"median"`
    StdDev      float64   `json:"std_dev"`
    LastValue   float64   `json:"last_value"`
    LastUpdate  time.Time `json:"last_update"`
}

type DashboardData struct {
    Metrics     map[string][]*MetricPoint `json:"metrics"`
    Alerts      []*Alert                 `json:"alerts"`
    Stats       map[string]*MetricStats   `json:"stats"`
    TimeRange   *TimeRange               `json:"time_range"`
}
```

#### Advanced Alerting System
```go
// pkg/monitoring/alerts/manager.go
type AlertManager interface {
    CreateAlert(alert *Alert) error
    UpdateAlert(id string, alert *Alert) error
    DeleteAlert(id string) error
    GetAlert(id string) (*Alert, error)
    GetAlerts(options *AlertQuery) ([]*Alert, error)
    ProcessMetric(metric *Metric) ([]*AlertInstance, error)
    SendNotification(alert *AlertInstance) error
    GetAlertHistory(id string, limit int) ([]*AlertInstance, error)
    RegisterNotifier(notifier Notifier) error
}

type Alert struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    MetricName  string                 `json:"metric_name"`
    Condition   *AlertCondition        `json:"condition"`
    Severity    string                 `json:"severity"` // info, warning, error, critical
    Status      string                 `json:"status"`   // active, inactive, triggered
    Tags        map[string]string      `json:"tags,omitempty"`
    Notifiers   []string               `json:"notifiers"` // email, slack, webhook, etc.
    Silenced    bool                   `json:"silenced"`
    SilencedUntil time.Time           `json:"silenced_until,omitempty"`
    CreatedAt   time.Time             `json:"created_at"`
    UpdatedAt   time.Time             `json:"updated_at"`
    CreatedBy   string                 `json:"created_by,omitempty"`
    SnoozeOptions *SnoozeOptions      `json:"snooze_options,omitempty"`
    NotificationConfig map[string]interface{} `json:"notification_config,omitempty"`
}

type AlertCondition struct {
    Type        string    `json:"type"` // threshold, change, anomaly
    Threshold   float64   `json:"threshold,omitempty"`
    Operator    string    `json:"operator,omitempty"` // gt, lt, eq, neq, gte, lte
    Duration    string    `json:"duration,omitempty"` // Time period condition must be true
    Percentage  float64   `json:"percentage,omitempty"` // For change conditions
    Direction   string    `json:"direction,omitempty"` // up, down, both
    Sensitivity float64   `json:"sensitivity,omitempty"` // For anomaly conditions
    WindowSize  int       `json:"window_size,omitempty"` // For moving averages
}

type AlertInstance struct {
    ID          string    `json:"id"`
    AlertID     string    `json:"alert_id"`
    Timestamp   time.Time `json:"timestamp"`
    Value       float64   `json:"value"`
    Threshold   float64   `json:"threshold,omitempty"`
    Message     string    `json:"message"`
    Status      string    `json:"status"` // triggered, resolved
    ResolvedAt  time.Time `json:"resolved_at,omitempty"`
    NotifiedAt  time.Time `json:"notified_at,omitempty"`
    AcknowledgedBy string  `json:"acknowledged_by,omitempty"`
    AcknowledgedAt time.Time `json:"acknowledged_at,omitempty"`
}

type AlertQuery struct {
    IDs         []string          `json:"ids,omitempty"`
    Severities  []string          `json:"severities,omitempty"`
    Statuses    []string          `json:"statuses,omitempty"`
    Tags        map[string]string `json:"tags,omitempty"`
    TimeRange   *TimeRange        `json:"time_range,omitempty"`
    Limit       int               `json:"limit,omitempty"`
}

type SnoozeOptions struct {
    Duration    string    `json:"duration"` // 1h, 1d, etc.
    Until       time.Time `json:"until,omitempty"`
    Reason      string    `json:"reason,omitempty"`
    SnoozedBy   string    `json:"snoozed_by,omitempty"`
}

type Notifier interface {
    GetType() string
    Notify(alert *AlertInstance, config map[string]interface{}) error
    ValidateConfig(config map[string]interface{}) error
}
```

#### Email Notifier Implementation
```go
// pkg/monitoring/alerts/notifiers/email.go
type EmailNotifier struct {
    smtpConfig *SMTPConfig
    templates  map[string]*template.Template
}

type SMTPConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    From     string `json:"from"`
    UseTLS   bool   `json:"use_tls"`
}

func (n *EmailNotifier) GetType() string {
    return "email"
}

func (n *EmailNotifier) Notify(alert *AlertInstance, config map[string]interface{}) error {
    recipients, ok := config["recipients"].([]string)
    if !ok || len(recipients) == 0 {
        return errors.New("missing or invalid recipients")
    }
    
    templateName, _ := config["template"].(string)
    if templateName == "" {
        templateName = "default"
    }
    
    tmpl, ok := n.templates[templateName]
    if !ok {
        return fmt.Errorf("template not found: %s", templateName)
    }
    
    var body bytes.Buffer
    if err := tmpl.Execute(&body, alert); err != nil {
        return err
    }
    
    subject := fmt.Sprintf("[%s] Alert: %s", strings.ToUpper(alert.Status), alert.Message)
    
    return n.sendEmail(recipients, subject, body.String())
}

func (n *EmailNotifier) ValidateConfig(config map[string]interface{}) error {
    recipients, ok := config["recipients"].([]string)
    if !ok || len(recipients) == 0 {
        return errors.New("missing or invalid recipients")
    }
    
    for _, recipient := range recipients {
        if !isValidEmail(recipient) {
            return fmt.Errorf("invalid email address: %s", recipient)
        }
    }
    
    return nil
}
```

#### Slack Notifier Implementation
```go
// pkg/monitoring/alerts/notifiers/slack.go
type SlackNotifier struct {
    httpClient *http.Client
}

func (n *SlackNotifier) GetType() string {
    return "slack"
}

func (n *SlackNotifier) Notify(alert *AlertInstance, config map[string]interface{}) error {
    webhookURL, ok := config["webhook_url"].(string)
    if !ok || webhookURL == "" {
        return errors.New("missing or invalid webhook URL")
    }
    
    channel, _ := config["channel"].(string)
    
    // Build Slack message payload
    color := "#ff0000" // Default to red for alerts
    if alert.Status == "resolved" {
        color = "#36a64f" // Green for resolved
    }
    
    payload := map[string]interface{}{
        "attachments": []map[string]interface{}{
            {
                "fallback":    alert.Message,
                "color":       color,
                "title":       fmt.Sprintf("Alert: %s", alert.Message),
                "text":        fmt.Sprintf("Value: %.2f\nThreshold: %.2f\nStatus: %s", 
                                alert.Value, alert.Threshold, strings.ToUpper(alert.Status)),
                "footer":      "Nessi Monitoring",
                "ts":          alert.Timestamp.Unix(),
            },
        },
    }
    
    if channel != "" {
        payload["channel"] = channel
    }
    
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := n.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("slack API error: %s, %s", resp.Status, string(body))
    }
    
    return nil
}

func (n *SlackNotifier) ValidateConfig(config map[string]interface{}) error {
    webhookURL, ok := config["webhook_url"].(string)
    if !ok || webhookURL == "" {
        return errors.New("missing or invalid webhook URL")
    }
    
    // Validate webhook URL format
    if !strings.HasPrefix(webhookURL, "https://hooks.slack.com/") {
        return errors.New("invalid Slack webhook URL format")
    }
    
    return nil
}
```

### Phase 6: Collaboration & Governance (2 weeks)

#### Role-Based Access Control
```go
// pkg/security/rbac/manager.go
type RBACManager interface {
    // Role management
    CreateRole(role *Role) error
    UpdateRole(id string, role *Role) error
    DeleteRole(id string) error
    GetRole(id string) (*Role, error)
    GetRoles() ([]*Role, error)
    
    // Permission management
    AddPermissionToRole(roleID string, permission Permission) error
    RemovePermissionFromRole(roleID string, permission Permission) error
    HasPermission(userID string, permission Permission) (bool, error)
    
    // User-role assignment
    AssignRoleToUser(userID string, roleID string) error
    RemoveRoleFromUser(userID string, roleID string) error
    GetUserRoles(userID string) ([]*Role, error)
    GetUsersWithRole(roleID string) ([]string, error)
}

type Role struct {
    ID          string       `json:"id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Permissions []Permission `json:"permissions"`
    IsSystem    bool         `json:"is_system"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type Permission struct {
    Resource    string `json:"resource"`
    Action      string `json:"action"` // create, read, update, delete, execute
    Constraints map[string]interface{} `json:"constraints,omitempty"`
}

type UserRoleMapping struct {
    UserID      string    `json:"user_id"`
    RoleID      string    `json:"role_id"`
    AssignedAt  time.Time `json:"assigned_at"`
    AssignedBy  string    `json:"assigned_by"`
}
```

#### Audit Logging
```go
// pkg/security/audit/logger.go
type AuditLogger interface {
    LogEvent(event *AuditEvent) error
    GetEvents(options *AuditQuery) ([]*AuditEvent, error)
    GetEventsByUser(userID string, limit int) ([]*AuditEvent, error)
    GetEventsByResource(resource string, limit int) ([]*AuditEvent, error)
}

type AuditEvent struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    UserID      string                 `json:"user_id"`
    Action      string                 `json:"action"`
    Resource    string                 `json:"resource"`
    ResourceID  string                 `json:"resource_id,omitempty"`
    Status      string                 `json:"status"` // success, failure
    Message     string                 `json:"message"`
    IPAddress   string                 `json:"ip_address,omitempty"`
    UserAgent   string                 `json:"user_agent,omitempty"`
    Changes     map[string]interface{} `json:"changes,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type AuditQuery struct {
    UserIDs     []string  `json:"user_ids,omitempty"`
    Actions     []string  `json:"actions,omitempty"`
    Resources   []string  `json:"resources,omitempty"`
    Statuses    []string  `json:"statuses,omitempty"`
    TimeRange   *TimeRange `json:"time_range,omitempty"`
    Limit       int       `json:"limit,omitempty"`
}
```

#### Completed Features
- **IQR-based Outlier Detection**: Implemented a more robust outlier detection method using Interquartile Range (IQR), which works better for skewed distributions compared to z-score method.
- **Sudden Change Detection**: Added capability to detect various patterns of changes in metrics over time:
  - Constant changes: Gradual increases or decreases over multiple runs
  - Sudden spikes/dips: Temporary anomalies that return to normal
  - Oscillations: Alternating patterns in metric values
- **Enhanced Rule Validation**:
  - Expanded predefined rule library (length, regex, enum validations)
  - Implemented custom rule authoring via YAML/SQL/Python
  - Added rule execution history tracking
- **Advanced Monitoring & Alerts**:
  - Developed real-time dashboard with interactive reports
  - Implemented metric retention system with 30-day history
  - Added comprehensive alerting capabilities (email, Slack, webhook notifications)
  - Created alert acknowledgment and silencing features



## Notes

### Performance Considerations
- Cache Python script results where possible
- Minimize data serialization/deserialization
- Use streaming for large datasets
- Monitor Python process memory usage

### Security Implementation

#### Authentication System
- **User Management**:
  - Role-based access control with predefined roles (admin, user, viewer)
  - Password hashing using bcrypt for secure storage
  - User creation, listing, updating, and deletion APIs
  - Default admin user creation for initial setup

#### Authorization
- **JWT-based Authentication**:
  - JSON Web Token (JWT) generation and validation
  - Token expiration and refresh mechanisms
  - Role-based middleware for protecting routes
  - Context-based user information storage

#### API Security
- **API Key Support**:
  - API key generation for programmatic access
  - API key validation and user association
  - API key regeneration functionality

#### Transport Security
- **SSL/TLS Support**:
  - Certificate management with auto-generation capability
  - TLS configuration with modern security settings
  - Option to require HTTPS for all connections
  - Self-signed certificate generation for development

#### Dashboard Security
- **Secure Web Interface**:
  - Login page with token storage
  - Protected API endpoints
  - Authentication middleware for web routes
  - Secure file downloads with authentication

#### Testing
- **Security Verification**:
  - Integration tests for authentication flows
  - SSL certificate generation testing
  - Role-based access control verification
  - API key authentication testing

### Additional Security Considerations

#### Input Validation
- Strict schema validation for all data exchanges
- File path sanitization and validation
- Command injection prevention
- Size limits for inputs

#### Process Isolation
- Resource limits (CPU, memory) for Python processes
- Temporary file cleanup
- Restricted file system access
- Network access control

#### Error Handling
- Timeouts for all operations
- Graceful process termination
- Error message sanitization
- Audit logging

#### Resource Management
- Process pool limits
- Disk space monitoring
- Memory usage tracking
- CPU usage restrictions
