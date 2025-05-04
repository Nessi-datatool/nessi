# API Reference

## Authentication

All API endpoints require authentication using an API key. Include the API key in the `X-API-Key` header:

```bash
curl -H "X-API-Key: your-api-key" https://api.nessi.dev/v1/endpoint
```

## Endpoints

### Table Management

#### Create Table
```http
POST /v1/tables
```

Create a new Delta Lake table with specified format and optimizations.

**Request Body:**
```json
{
  "path": "s3://bucket/table",
  "format": "delta",
  "partition_by": ["date", "region"],
  "z_order_by": ["id"]
}
```

#### Time Travel
```http
GET /v1/tables/{path}/versions/{version}
```

Access a specific version of a Delta Lake table.

**Parameters:**
- `path`: Table path
- `version`: Version number or timestamp

### Data Quality

#### Run Quality Check
```http
POST /v1/quality/check
```

Run data quality checks on a table.

**Request Body:**
```json
{
  "table_path": "s3://bucket/table",
  "rules": {
    "completeness": 0.95,
    "accuracy": 0.98,
    "consistency": 0.97
  }
}
```

#### Get Quality Report
```http
GET /v1/quality/report/{table_path}
```

Get a quality report for a table.

**Parameters:**
- `table_path`: Path to the table
- `format`: Report format (html, pdf, json)

### Monitoring

#### Get Metrics
```http
GET /v1/metrics
```

Get current metric values.

**Response:**
```json
{
  "data_quality_score": 0.95,
  "processing_time": 1.2,
  "error_rate": 0.01
}
```

#### Get Alerts
```http
GET /v1/alerts
```

Get current alerts.

**Response:**
```json
{
  "alerts": [
    {
      "name": "HighErrorRate",
      "severity": "critical",
      "message": "Error rate exceeds threshold"
    }
  ]
}
```

### User Management

#### Create User
```http
POST /v1/users
```

Create a new user with specified role.

**Request Body:**
```json
{
  "username": "user1",
  "email": "user1@example.com",
  "role": "analyst"
}
```

#### List Users
```http
GET /v1/users
```

List all users.

### Pipeline Management

#### Create Pipeline
```http
POST /v1/pipelines
```

Create a new transformation pipeline.

**Request Body:**
```json
{
  "name": "clean_data",
  "transformations": [
    {
      "type": "filter",
      "condition": "age > 18"
    },
    {
      "type": "rename",
      "mapping": {
        "old_name": "new_name"
      }
    }
  ]
}
```

#### Run Pipeline
```http
POST /v1/pipelines/{name}/run
```

Run a transformation pipeline on a table.

**Request Body:**
```json
{
  "input_path": "s3://bucket/input",
  "output_path": "s3://bucket/output"
}
```

## Error Responses

All endpoints may return the following error responses:

- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Invalid or missing API key
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Rate Limiting

API requests are limited to 100 requests per minute per IP address. Allowlisted IPs bypass rate limiting. 