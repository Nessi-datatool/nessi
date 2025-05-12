# Nessi.dev Implementation Plan

## Architecture Overview

### Core Architecture
```
nessi/
├── cmd/
│   └── nessi/          # Main CLI application
├── pkg/
│   ├── datalake/       # Delta Lake operations (Go)
│   ├── quality/        # Data quality & profiling (Go)
│   ├── monitoring/     # Metrics & monitoring (Go)
│   ├── security/       # Authentication & authorization (Go)
│   └── api/           # API endpoints (Go)
└── scripts/
    ├── delta/         # Delta Lake operations (Python)
    ├── profiling/     # Advanced profiling (Python)
    └── viz/           # Data visualization (Python)
```

## Implementation Phases

### Phase 1: Core Infrastructure (2 weeks)
- [x] Basic project structure
- [x] Security implementation
- [x] Monitoring setup
- [ ] Go-Python integration layer
  - Implementation: Use `os/exec` to call Python scripts
  - Data exchange via JSON
  - Error handling and logging

### Phase 2: Delta Lake Integration (1 week)

#### Delta Lake Interface
```go
// pkg/datalake/delta.go
type DeltaTable interface {
    // Core operations
    Initialize(schema *arrow.Schema) error
    Read(partition string) (arrow.Record, error)
    Write(partition string, record arrow.Record) error
    GetSchema() (*arrow.Schema, error)
    GetStats() (*TableStats, error)

    // Version management
    GetVersion() (int64, error)
    GetVersionAt(timestamp time.Time) (int64, error)
    Commit(metadata map[string]interface{}) error
    Rollback(version int64) error

    // Metadata operations
    GetMetadata() (map[string]interface{}, error)
    UpdateMetadata(metadata map[string]interface{}) error
}
```

#### Python Delta Lake Scripts
```python
# scripts/delta/table.py
from deltalake import DeltaTable
from pyarrow import Table, Schema

class DeltaManager:
    def __init__(self, table_path: str):
        self.table = DeltaTable(table_path)
    
    def read_version(self, version: int = None) -> Table:
        if version is not None:
            return self.table.to_pyarrow_table(version=version)
        return self.table.to_pyarrow_table()
    
    def write_records(self, table: Table, partition_cols: list = None) -> None:
        self.table.write_deltalake(table, partition_cols=partition_cols)
    
    def get_metadata(self) -> dict:
        return {
            'version': self.table.version(),
            'metadata': self.table.metadata(),
            'schema': self.table.schema().to_pyarrow(),
            'files': self.table.files(),
            'stats': self.table.stats()
        }
```

### Phase 3: Data Quality & Profiling (1 week)

#### Rule Engine
```go
// pkg/quality/rules.go
type Rule interface {
    Validate(record arrow.Record) []ValidationError
    GetMetadata() RuleMetadata
}

type RuleMetadata struct {
    ID          string                 `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Severity    string                 `json:"severity"`
    Tags        []string               `json:"tags"`
    Config      map[string]interface{} `json:"config"`
}

type ValidationError struct {
    RuleID      string `json:"rule_id"`
    Field       string `json:"field"`
    Value       string `json:"value"`
    Message     string `json:"message"`
    RowIndex    int64  `json:"row_index"`
}
```

#### Rule Configuration
```yaml
# config/rules.yaml
rules:
  - id: null_check
    name: Null Value Check
    description: Checks for null values in required fields
    severity: error
    tags: [data_quality, nulls]
    config:
      fields: [id, name, value]

  - id: range_check
    name: Value Range Check
    description: Validates numeric values are within expected ranges
    severity: warning
    tags: [data_quality, range]
    config:
      fields:
        age: {min: 0, max: 120}
        score: {min: 0, max: 100}
```

#### Profiling Integration
```python
# scripts/profiling/profile.py
from ydata_profiling import ProfileReport
import pyarrow as pa
import pandas as pd

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

### Next Steps
- **Rule Validation Enhancements**:
  - Expand predefined rule library (length, regex, enum validations)
  - Implement custom rule authoring via YAML/SQL/Python
  - Add rule execution history tracking
- **Monitoring & Alerts**:
  - Develop real-time dashboard (web UI or Grafana integration)
  - Implement metric retention system (30-day history)
  - Add alerting capabilities (email/CLI notifications)

## Future Enhancements

### Phase 2 Features
- Advanced Delta Lake operations
- ML-based anomaly detection
- Advanced visualization options
- API extensions

## Notes

### Performance Considerations
- Cache Python script results where possible
- Minimize data serialization/deserialization
- Use streaming for large datasets
- Monitor Python process memory usage

### Security Considerations

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
