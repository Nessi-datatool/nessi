# Nessi gRPC API Documentation

## Overview

Nessi provides a gRPC API for high-performance, strongly-typed interactions with the Nessi platform. This API is ideal for programmatic access and integration with other systems.

## Connection Details

```
Host: localhost
Port: 9090
```

## Authentication

All gRPC methods require authentication using a token. Include the token in your requests using the metadata key `authorization` with the value `Bearer <your-token>`.

```go
// Example in Go
ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
response, err := client.GetTable(ctx, request)
```

## Service Definitions

Nessi's gRPC API is defined using Protocol Buffers. The main services are:

### TableService

Manages Delta Lake tables and their metadata.

```protobuf
service TableService {
  // Lists all available tables
  rpc ListTables(ListTablesRequest) returns (ListTablesResponse);
  
  // Gets details for a specific table
  rpc GetTable(GetTableRequest) returns (GetTableResponse);
  
  // Gets profile information for a table
  rpc GetTableProfile(GetTableProfileRequest) returns (GetTableProfileResponse);
  
  // Checks quality of a table against rules
  rpc CheckTableQuality(CheckTableQualityRequest) returns (CheckTableQualityResponse);
  
  // Lists all versions of a table
  rpc ListTableVersions(ListTableVersionsRequest) returns (ListTableVersionsResponse);
}
```

### MetricsService

Provides access to metrics and monitoring data.

```protobuf
service MetricsService {
  // Gets metrics for tables
  rpc GetTableMetrics(GetTableMetricsRequest) returns (GetTableMetricsResponse);
  
  // Gets quality metrics
  rpc GetQualityMetrics(GetQualityMetricsRequest) returns (GetQualityMetricsResponse);
  
  // Gets performance metrics
  rpc GetPerformanceMetrics(GetPerformanceMetricsRequest) returns (GetPerformanceMetricsResponse);
  
  // Streams real-time metrics updates
  rpc StreamMetrics(StreamMetricsRequest) returns (stream MetricUpdate);
}
```

### ValidationService

Handles data validation and quality rules.

```protobuf
service ValidationService {
  // Creates a validation rule
  rpc CreateRule(CreateRuleRequest) returns (CreateRuleResponse);
  
  // Updates an existing validation rule
  rpc UpdateRule(UpdateRuleRequest) returns (UpdateRuleResponse);
  
  // Deletes a validation rule
  rpc DeleteRule(DeleteRuleRequest) returns (DeleteRuleResponse);
  
  // Lists all validation rules
  rpc ListRules(ListRulesRequest) returns (ListRulesResponse);
  
  // Applies validation rules to a table
  rpc ValidateTable(ValidateTableRequest) returns (ValidateTableResponse);
}
```

## Message Definitions

### Common Types

```protobuf
message Table {
  string name = 1;
  string path = 2;
  int64 version = 3;
  repeated string columns = 4;
  TableMetadata metadata = 5;
}

message TableMetadata {
  int64 created_time = 1;
  string location = 2;
  repeated string partition_columns = 3;
}

message Column {
  string name = 1;
  string type = 2;
  ColumnMetrics metrics = 3;
}

message ColumnMetrics {
  int64 count = 1;
  int64 distinct_count = 2;
  int64 null_count = 3;
  string min = 4;
  string max = 5;
  double mean = 6;
  double stddev = 7;
}

message ValidationRule {
  string name = 1;
  string description = 2;
  RuleType type = 3;
  string column = 4;
  double threshold = 5;
  Severity severity = 6;
  string pattern = 7;  // For pattern rules
  double min = 8;      // For range rules
  double max = 9;      // For range rules
}

enum RuleType {
  COMPLETENESS = 0;
  UNIQUENESS = 1;
  PATTERN = 2;
  RANGE = 3;
}

enum Severity {
  CRITICAL = 0;
  HIGH = 1;
  MEDIUM = 2;
  LOW = 3;
}

message ValidationResult {
  string rule_name = 1;
  bool passed = 2;
  int64 violations = 3;
}
```

### Request/Response Types

```protobuf
// ListTables
message ListTablesRequest {
  string path = 1;  // Optional root path
}

message ListTablesResponse {
  repeated Table tables = 1;
}

// GetTable
message GetTableRequest {
  string table_name = 1;
  optional int64 version = 2;  // Optional specific version
}

message GetTableResponse {
  Table table = 1;
}

// GetTableProfile
message GetTableProfileRequest {
  string table_name = 1;
  optional int32 sample = 2;  // Optional sampling percentage (0-100)
}

message GetTableProfileResponse {
  repeated Column columns = 1;
}

// CheckTableQuality
message CheckTableQualityRequest {
  string table_name = 1;
  repeated ValidationRule rules = 2;
}

message CheckTableQualityResponse {
  string table_name = 1;
  string timestamp = 2;
  int64 version = 3;
  repeated Metric metrics = 4;
  repeated ValidationResult rule_results = 5;
}

message Metric {
  string name = 1;
  double value = 2;
  bool passed = 3;
}

// ListTableVersions
message ListTableVersionsRequest {
  string table_name = 1;
}

message ListTableVersionsResponse {
  repeated int64 versions = 1;
}
```

## Error Handling

gRPC errors are returned using standard gRPC status codes:

- `INVALID_ARGUMENT (3)`: Invalid request parameters
- `UNAUTHENTICATED (16)`: Invalid or missing authentication
- `NOT_FOUND (5)`: Resource not found
- `INTERNAL (13)`: Server error occurred
- `PERMISSION_DENIED (7)`: Insufficient permissions
- `RESOURCE_EXHAUSTED (8)`: Rate limit exceeded

Each error includes a detailed message and may include additional metadata.

## Client Libraries

Nessi provides client libraries for several languages:

- Go: `github.com/nessi-dev/nessi-go-client`
- Python: `nessi-python-client`
- Java: `com.nessi.client`

## Example Usage

### Go Client Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nessi-dev/nessi-go-client/nessi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := nessi.NewTableServiceClient(conn)

	// Add authentication token
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer your-token-here")

	// List tables
	resp, err := client.ListTables(ctx, &nessi.ListTablesRequest{})
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	// Print tables
	for _, table := range resp.Tables {
		fmt.Printf("Table: %s, Path: %s, Version: %d\n", table.Name, table.Path, table.Version)
	}
}
```

### Python Client Example

```python
import grpc
from nessi_python_client import nessi_pb2, nessi_pb2_grpc

def main():
    # Create a gRPC channel
    channel = grpc.insecure_channel('localhost:9090')
    
    # Create a stub (client)
    stub = nessi_pb2_grpc.TableServiceStub(channel)
    
    # Add authentication token
    metadata = [('authorization', 'Bearer your-token-here')]
    
    # List tables
    response = stub.ListTables(nessi_pb2.ListTablesRequest(), metadata=metadata)
    
    # Print tables
    for table in response.tables:
        print(f"Table: {table.name}, Path: {table.path}, Version: {table.version}")

if __name__ == '__main__':
    main()
```

## Streaming Example

Nessi supports streaming for real-time metrics updates:

```go
// Go example
req := &nessi.StreamMetricsRequest{
    TableName: "customers",
    MetricTypes: []nessi.MetricType{nessi.MetricType_QUALITY, nessi.MetricType_PERFORMANCE},
}

stream, err := client.StreamMetrics(ctx, req)
if err != nil {
    log.Fatalf("Failed to start streaming: %v", err)
}

for {
    update, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("Error while streaming: %v", err)
    }
    
    fmt.Printf("Received metric update: %s = %f\n", update.Name, update.Value)
}
```

## Community and Pro Features

The gRPC API includes both Community Edition and Pro Edition features. Pro features are marked with a `[PRO]` tag in the documentation and require a valid Pro license to access.

For more information on Nessi's licensing, see the [License Management documentation](../../LICENSE_MANAGEMENT.md).
