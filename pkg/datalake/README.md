# Delta Lake Format Handler for Nessi CLI

This package provides support for reading and writing data in Delta Lake format, with time travel capabilities. It's designed to work seamlessly with the Databricks integration in the CLI-only version of Nessi.

## Features

- Read data from Delta Lake tables
- Write data to Delta Lake tables
- Time travel capabilities:
  - Read data as of a specific version
  - Read data as of a specific timestamp
  - Get version history of a Delta table
- Schema inference for Delta Lake tables
- Support for various data types including nested structures

## Usage Examples

### Reading from a Delta Table

```go
import (
    "context"
    "github.com/nessi-dev/nessi/pkg/datalake"
)

func readDeltaTable() {
    // Initialize Delta format handler
    deltaHandler := datalake.NewDeltaFormatHandler()
    
    // Read data from Delta table
    data, err := deltaHandler.Read("/path/to/delta/table")
    if err != nil {
        // Handle error
    }
    
    // Process data
    for _, record := range data {
        // Each record is a map[string]interface{}
    }
}
```

### Writing to a Delta Table

```go
import (
    "context"
    "github.com/nessi-dev/nessi/pkg/datalake"
)

func writeDeltaTable() {
    // Create data to write
    data := []map[string]interface{}{
        {"id": int64(1), "name": "John Doe", "age": int32(30), "active": true},
        {"id": int64(2), "name": "Jane Smith", "age": int32(25), "active": true},
    }
    
    // Create schema
    schema := datalake.NewSchema([]datalake.Field{
        {Name: "id", Type: datalake.FieldTypeInt64},
        {Name: "name", Type: datalake.FieldTypeString},
        {Name: "age", Type: datalake.FieldTypeInt32},
        {Name: "active", Type: datalake.FieldTypeBool},
    })
    
    // Initialize Delta format handler
    deltaHandler := datalake.NewDeltaFormatHandler()
    
    // Write data to Delta table
    err := deltaHandler.Write("/path/to/delta/table", data, schema)
    if err != nil {
        // Handle error
    }
}
```

### Time Travel

```go
import (
    "context"
    "time"
    "github.com/nessi-dev/nessi/pkg/datalake"
)

func timeTravel() {
    // Initialize Delta format handler
    deltaHandler := datalake.NewDeltaFormatHandler()
    
    // Read as of version 2
    dataV2, err := deltaHandler.ReadAsOfVersion("/path/to/delta/table", 2)
    if err != nil {
        // Handle error
    }
    
    // Read as of timestamp
    timestamp := time.Date(2025, 5, 18, 10, 0, 0, 0, time.UTC)
    dataTS, err := deltaHandler.ReadAsOfTimestamp("/path/to/delta/table", timestamp)
    if err != nil {
        // Handle error
    }
    
    // Get version history
    history, err := deltaHandler.GetVersionHistory("/path/to/delta/table")
    if err != nil {
        // Handle error
    }
    
    for _, version := range history {
        // Process version information
        // version.Version, version.Timestamp, version.Operation
    }
}
```

## Schema Handling

The Delta Lake format handler supports the following field types:

- `FieldTypeInt32` - 32-bit integer
- `FieldTypeInt64` - 64-bit integer
- `FieldTypeFloat32` - 32-bit floating point
- `FieldTypeFloat64` - 64-bit floating point
- `FieldTypeString` - String
- `FieldTypeBool` - Boolean
- `FieldTypeTimestamp` - Timestamp
- `FieldTypeDate` - Date
- `FieldTypeStruct` - Nested structure
- `FieldTypeList` - List/array

## Error Handling

The Delta Lake format handler includes comprehensive error handling for various scenarios:

- File system errors when reading/writing Delta files
- Schema mismatch errors when writing data
- Version not found errors for time travel
- Timestamp not found errors for time travel
- Invalid Delta table errors

## Development

### Running Tests

```bash
# Run all tests
go test -v ./pkg/datalake/...

# Run with coverage
go test -coverprofile=coverage.out ./pkg/datalake/...
go tool cover -html=coverage.out
```

### Dependencies

This package depends on:

- `github.com/apache/arrow/go/v15` for Arrow data format handling

## Contributing

Contributions to improve the Delta Lake format handler are welcome. Please ensure all tests pass before submitting a pull request.
