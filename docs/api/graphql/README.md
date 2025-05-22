# Nessi GraphQL API Documentation

## Overview

Nessi provides a GraphQL API for flexible, client-driven data access. GraphQL allows clients to request exactly the data they need, making it ideal for dashboard applications and custom integrations.

## Endpoint

```
http://localhost:8080/graphql
```

## Authentication

All GraphQL operations require authentication using a Bearer token. Include the token in your requests:

```http
Authorization: Bearer <your-token>
```

## Schema

The Nessi GraphQL API provides the following main types and operations:

### Types

#### Table

```graphql
type Table {
  id: ID!
  name: String!
  path: String!
  version: Int!
  columns: [String!]!
  metadata: TableMetadata!
  profile: TableProfile
  versions: [Int!]!
  metrics: [Metric!]!
}

type TableMetadata {
  createdTime: DateTime!
  location: String!
  partitionColumns: [String!]!
  format: String!
  properties: JSONObject
}

type TableProfile {
  columns: [ColumnProfile!]!
  rowCount: Int!
  sizeBytes: Int!
  lastUpdated: DateTime!
}

type ColumnProfile {
  name: String!
  type: String!
  metrics: ColumnMetrics!
  histogram: Histogram
  topValues: [ValueCount!]
}

type ColumnMetrics {
  count: Int!
  distinctCount: Int!
  nullCount: Int!
  min: String
  max: String
  mean: Float
  stddev: Float
  percentiles: Percentiles
}

type Percentiles {
  p25: Float
  p50: Float
  p75: Float
  p90: Float
  p95: Float
  p99: Float
}

type Histogram {
  buckets: [HistogramBucket!]!
}

type HistogramBucket {
  min: Float!
  max: Float!
  count: Int!
}

type ValueCount {
  value: String!
  count: Int!
}

type Metric {
  name: String!
  value: Float!
  timestamp: DateTime!
  status: MetricStatus!
}

enum MetricStatus {
  PASS
  FAIL
  WARNING
}

type ValidationRule {
  id: ID!
  name: String!
  description: String
  type: RuleType!
  column: String!
  threshold: Float
  severity: Severity!
  pattern: String
  min: Float
  max: Float
}

enum RuleType {
  COMPLETENESS
  UNIQUENESS
  PATTERN
  RANGE
}

enum Severity {
  CRITICAL
  HIGH
  MEDIUM
  LOW
}

type ValidationResult {
  rule: ValidationRule!
  passed: Boolean!
  violations: Int!
  details: String
}

type QualityCheckResult {
  tableName: String!
  timestamp: DateTime!
  version: Int!
  metrics: [Metric!]!
  ruleResults: [ValidationResult!]!
}

# Custom scalar types
scalar DateTime
scalar JSONObject
```

### Queries

```graphql
type Query {
  # Get a list of all tables
  tables(path: String): [Table!]!
  
  # Get a specific table by name
  table(name: String!, version: Int): Table
  
  # Get table profile
  tableProfile(name: String!, sample: Int): TableProfile
  
  # Get table quality check results
  tableQuality(name: String!, rules: [ValidationRuleInput!]!): QualityCheckResult
  
  # Get validation rules
  validationRules(tablePattern: String): [ValidationRule!]!
  
  # Get metrics
  metrics(
    tableName: String, 
    metricTypes: [MetricType!], 
    from: DateTime, 
    to: DateTime
  ): [Metric!]!
}

enum MetricType {
  QUALITY
  PERFORMANCE
  FRESHNESS
  VOLUME
}

input ValidationRuleInput {
  name: String!
  description: String
  type: RuleType!
  column: String!
  threshold: Float
  severity: Severity!
  pattern: String
  min: Float
  max: Float
}
```

### Mutations

```graphql
type Mutation {
  # Create a validation rule
  createValidationRule(input: ValidationRuleInput!): ValidationRule!
  
  # Update a validation rule
  updateValidationRule(id: ID!, input: ValidationRuleInput!): ValidationRule!
  
  # Delete a validation rule
  deleteValidationRule(id: ID!): Boolean!
  
  # Run quality check on a table
  runQualityCheck(tableName: String!, rules: [ValidationRuleInput!]!): QualityCheckResult!
}
```

### Subscriptions

```graphql
type Subscription {
  # Subscribe to metric updates
  metricUpdates(tableName: String, metricTypes: [MetricType!]): Metric!
  
  # Subscribe to table updates
  tableUpdates(name: String): TableUpdate!
}

type TableUpdate {
  table: Table!
  updateType: UpdateType!
  timestamp: DateTime!
}

enum UpdateType {
  CREATED
  UPDATED
  DELETED
}
```

## Example Queries

### Get All Tables

```graphql
query GetTables {
  tables {
    id
    name
    path
    version
    columns
  }
}
```

### Get Table Details with Profile

```graphql
query GetTableDetails($name: String!) {
  table(name: $name) {
    id
    name
    path
    version
    columns
    metadata {
      createdTime
      location
      partitionColumns
      format
      properties
    }
    profile {
      columns {
        name
        type
        metrics {
          count
          distinctCount
          nullCount
          min
          max
          mean
          stddev
        }
        topValues {
          value
          count
        }
      }
      rowCount
      sizeBytes
      lastUpdated
    }
    versions
  }
}
```

### Run Quality Check

```graphql
mutation RunQualityCheck($tableName: String!, $rules: [ValidationRuleInput!]!) {
  runQualityCheck(tableName: $tableName, rules: $rules) {
    tableName
    timestamp
    version
    metrics {
      name
      value
      status
    }
    ruleResults {
      rule {
        name
        type
        severity
      }
      passed
      violations
      details
    }
  }
}

# Variables
{
  "tableName": "customers",
  "rules": [
    {
      "name": "email_validation",
      "type": "PATTERN",
      "column": "email",
      "pattern": "@",
      "severity": "HIGH"
    },
    {
      "name": "age_range",
      "type": "RANGE",
      "column": "age",
      "min": 0,
      "max": 120,
      "severity": "MEDIUM"
    }
  ]
}
```

### Subscribe to Metric Updates

```graphql
subscription WatchMetrics($tableName: String!, $metricTypes: [MetricType!]!) {
  metricUpdates(tableName: $tableName, metricTypes: $metricTypes) {
    name
    value
    timestamp
    status
  }
}

# Variables
{
  "tableName": "customers",
  "metricTypes": ["QUALITY", "PERFORMANCE"]
}
```

## Error Handling

GraphQL errors are returned in the `errors` array of the response:

```json
{
  "errors": [
    {
      "message": "Table not found: customers",
      "locations": [{"line": 2, "column": 3}],
      "path": ["table"],
      "extensions": {
        "code": "NOT_FOUND",
        "classification": "DataFetchingException"
      }
    }
  ],
  "data": {
    "table": null
  }
}
```

Common error codes include:

- `VALIDATION_ERROR`: Invalid input
- `AUTHENTICATION_ERROR`: Invalid or missing authentication
- `AUTHORIZATION_ERROR`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `INTERNAL_ERROR`: Server error occurred
- `RATE_LIMIT_EXCEEDED`: Too many requests

## Client Libraries

Nessi's GraphQL API can be used with any GraphQL client library:

- JavaScript/TypeScript: [Apollo Client](https://www.apollographql.com/docs/react/)
- Python: [gql](https://gql.readthedocs.io/)
- Java: [Apollo Android](https://www.apollographql.com/docs/android/)
- Go: [graphql-go-client](https://github.com/machinebox/graphql)

## Example Usage

### JavaScript (Apollo Client)

```javascript
import { ApolloClient, InMemoryCache, gql } from '@apollo/client';

// Create a client
const client = new ApolloClient({
  uri: 'http://localhost:8080/graphql',
  cache: new InMemoryCache(),
  headers: {
    authorization: 'Bearer your-token-here'
  }
});

// Query tables
client.query({
  query: gql`
    query GetTables {
      tables {
        id
        name
        path
        version
      }
    }
  `
})
.then(result => console.log(result.data.tables))
.catch(error => console.error(error));
```

### Python (gql)

```python
from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport

# Create a transport with authentication
transport = AIOHTTPTransport(
    url="http://localhost:8080/graphql",
    headers={"Authorization": "Bearer your-token-here"}
)

# Create a client
client = Client(transport=transport, fetch_schema_from_transport=True)

# Define the query
query = gql("""
    query GetTableProfile($name: String!, $sample: Int) {
        tableProfile(name: $name, sample: $sample) {
            columns {
                name
                type
                metrics {
                    count
                    nullCount
                    distinctCount
                }
            }
            rowCount
        }
    }
""")

# Execute the query
variables = {"name": "customers", "sample": 50}
result = client.execute(query, variable_values=variables)
print(result)
```

## GraphiQL Explorer

Nessi provides a GraphiQL explorer interface for interactive API exploration at:

```
http://localhost:8080/graphiql
```

This interface allows you to:
- Browse the schema documentation
- Compose and execute queries
- View query results
- Explore available types and fields

## Community and Pro Features

The GraphQL API includes both Community Edition and Pro Edition features. Pro features are marked with a `[PRO]` tag in the schema documentation and require a valid Pro license to access.

For more information on Nessi's licensing, see the [License Management documentation](../../LICENSE_MANAGEMENT.md).
