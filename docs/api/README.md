# Nessi API Documentation

## Overview

Nessi provides a RESTful API for managing and analyzing Delta Lake tables. The API is built using Gin and follows REST principles.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All endpoints require authentication using a Bearer token. Include the token in your requests:

```http
Authorization: Bearer <your-token>
```

## API Endpoints

### 1. Tables Management

#### List Tables
```http
GET /tables
```

**Query Parameters**
- `path` (optional): Root path to search for tables

**Response**
```json
{
  "tables": [
    {
      "name": "string",
      "path": "string",
      "version": int64,
      "columns": ["string"]
    }
  ]
}
```

#### Get Table Details
```http
GET /tables/{tableName}
```

**Query Parameters**
- `version` (optional): Specific version of the table to retrieve

**Response**
```json
{
  "name": "string",
  "path": "string",
  "version": int64,
  "columns": ["string"],
  "metadata": {
    "created_time": int64,
    "location": "string",
    "partition_columns": ["string"]
  }
}
```

### 2. Table Operations

#### Get Table Profile
```http
GET /tables/{tableName}/profile
```

**Query Parameters**
- `sample` (optional): Sampling percentage (0-100)

**Response**
```json
{
  "columns": [
    {
      "name": "string",
      "type": "string",
      "metrics": {
        "count": int64,
        "distinct_count": int64,
        "null_count": int64,
        "min": "string",
        "max": "string",
        "mean": float64,
        "stddev": float64
      }
    }
  ]
}
```

#### Check Table Quality
```http
POST /tables/{tableName}/check
```

**Request Body**
```json
{
  "rules": [
    {
      "name": "string",
      "description": "string",
      "type": "completeness|uniqueness|pattern|range",
      "column": "string",
      "threshold": float64,
      "severity": "critical|high|medium|low",
      "pattern": "string" (for pattern rules),
      "min": float64 (for range rules),
      "max": float64 (for range rules)
    }
  ]
}
```

**Response**
```json
{
  "table_name": "string",
  "timestamp": "string",
  "version": int64,
  "metrics": [
    {
      "name": "string",
      "value": float64,
      "status": "pass|fail"
    }
  ],
  "rules": [
    {
      "name": "string",
      "status": "pass|fail",
      "violations": int64
    }
  ]
}
```

#### List Table Versions
```http
GET /tables/{tableName}/versions
```

**Response**
```json
{
  "versions": [int64]
}
```

## Error Responses

All endpoints return appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid request parameters or body
- `401 Unauthorized`: Invalid or missing authentication
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error occurred

```json
{
  "error": "string"
}
```

## Monitoring

The API provides metrics that can be exported to various formats (JSON, CSV) for monitoring and analysis. Key metrics include:

- Table metrics: size, count, etc.
- Quality metrics: checks performed, violations
- Performance metrics: query latencies
- Error metrics: error counts by type

## Rate Limiting

The API implements rate limiting:
- 100 requests per minute per IP
- 10 concurrent requests per IP

## Versioning

The API follows semantic versioning with the format `v1`. Breaking changes will be introduced in major version increments.

## Example Usage

### List all tables
```bash
curl -X GET "http://localhost:8080/api/v1/tables" \
  -H "Authorization: Bearer <your-token>"
```

### Get table profile with sampling
```bash
curl -X GET "http://localhost:8080/api/v1/tables/customers/profile?sample=50" \
  -H "Authorization: Bearer <your-token>"
```

### Check table quality
```bash
curl -X POST "http://localhost:8080/api/v1/tables/customers/check" \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "rules": [
      {
        "name": "email_validation",
        "type": "pattern",
        "column": "email",
        "pattern": "@",
        "severity": "high"
      }
    ]
  }'
```
