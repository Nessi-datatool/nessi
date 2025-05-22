# Nessi API Documentation

## Overview

Nessi provides multiple API options to meet different integration needs. This document serves as an index to help you find the appropriate API documentation based on your requirements.

## Available APIs

Nessi offers three main API types:

### 1. REST API

Our RESTful API provides a straightforward HTTP-based interface for managing and analyzing Delta Lake tables. It's ideal for most integrations and follows standard REST principles.

- **Documentation**: [REST API Documentation](rest/README.md)
- **OpenAPI Specification**: [OpenAPI YAML](rest/openapi.yaml)
- **Base URL**: `http://localhost:8080/api/v1` (local) or `https://api.nessi.dev/v1` (production)
- **Authentication**: Bearer token

**Best for**:
- Web applications
- Scripts and automation
- General purpose integrations
- Clients in any programming language

### 2. gRPC API

Our gRPC API provides a high-performance, strongly-typed interface using Protocol Buffers. It's ideal for microservices and systems requiring efficient binary communication.

- **Documentation**: [gRPC API Documentation](grpc/README.md)
- **Proto Files**: Available in the [GitHub repository](https://github.com/nessi-dev/nessi/tree/main/api/proto)
- **Connection**: `localhost:9090` (local) or `grpc.nessi.dev:443` (production)
- **Authentication**: Bearer token via metadata

**Best for**:
- Microservices architectures
- High-performance systems
- Strongly-typed interfaces
- Streaming data
- Language-specific client generation

### 3. GraphQL API

Our GraphQL API provides a flexible, query-based interface that allows clients to request exactly the data they need. It's ideal for complex data requirements and frontend applications.

- **Documentation**: [GraphQL API Documentation](graphql/README.md)
- **Schema**: Available via introspection
- **Endpoint**: `http://localhost:8080/graphql` (local) or `https://api.nessi.dev/graphql` (production)
- **Explorer**: `http://localhost:8080/graphiql` (local) or `https://api.nessi.dev/graphiql` (production)
- **Authentication**: Bearer token

**Best for**:
- Frontend applications
- Dashboards and visualizations
- Complex data requirements
- Reducing over-fetching and under-fetching
- Aggregating multiple resources in a single request

## API Feature Comparison

| Feature | REST API | gRPC API | GraphQL API |
|---------|----------|----------|-------------|
| Performance | Good | Excellent | Good |
| Type Safety | No | Yes | Yes |
| Schema Definition | OpenAPI | Protocol Buffers | GraphQL Schema |
| Streaming | No | Yes | Yes (Subscriptions) |
| Browser Support | Native | Via gRPC-Web | Native |
| Tooling | Extensive | Good | Excellent |
| Learning Curve | Low | Medium | Medium |
| Payload Size | Larger (JSON) | Smaller (Binary) | Varies (JSON) |
| Community Edition | ✓ | ✓ | ✓ |
| Pro Edition Features | ✓ | ✓ | ✓ |

## Authentication

All APIs use the same authentication mechanism: Bearer tokens. Include the token in your requests:

- **REST**: HTTP header `Authorization: Bearer <your-token>`
- **gRPC**: Metadata `authorization: Bearer <your-token>`
- **GraphQL**: HTTP header `Authorization: Bearer <your-token>`

## Community and Pro Features

All APIs include both Community Edition and Pro Edition features. Pro features are marked in the documentation and require a valid Pro license to access.

For more information on Nessi's licensing, see the [License Management documentation](../LICENSE_MANAGEMENT.md).

## Client Libraries

Nessi provides official client libraries for several languages:

| Language | REST | gRPC | GraphQL |
|----------|------|------|---------|
| Go | [nessi-go-client](https://github.com/nessi-dev/nessi-go-client) | [nessi-go-client](https://github.com/nessi-dev/nessi-go-client) | [nessi-go-client](https://github.com/nessi-dev/nessi-go-client) |
| Python | [nessi-python-client](https://pypi.org/project/nessi-client/) | [nessi-python-client](https://pypi.org/project/nessi-client/) | [nessi-python-client](https://pypi.org/project/nessi-client/) |
| JavaScript | [nessi-js-client](https://www.npmjs.com/package/nessi-client) | [nessi-js-client](https://www.npmjs.com/package/nessi-client) | [nessi-js-client](https://www.npmjs.com/package/nessi-client) |
| Java | [nessi-java-client](https://github.com/nessi-dev/nessi-java-client) | [nessi-java-client](https://github.com/nessi-dev/nessi-java-client) | [nessi-java-client](https://github.com/nessi-dev/nessi-java-client) |

## API Versioning

All APIs follow semantic versioning:

- **REST**: Version in URL path (`/v1/...`)
- **gRPC**: Version in package name (`nessi.v1`)
- **GraphQL**: Version in schema types (`TableV1`)

Breaking changes will only be introduced in major version increments.

## Rate Limiting

All APIs implement rate limiting:

- 100 requests per minute per IP
- 10 concurrent requests per IP

Pro edition users have higher rate limits.

## Getting Started

1. Choose the API that best fits your needs based on the comparison above
2. Review the specific API documentation
3. Set up authentication
4. Use the appropriate client library or make direct API calls

## Support

If you need help with the APIs:

- **Community Support**: [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions)
- **Pro Support**: [Support Portal](https://nessi.dev/support)
- **Documentation Issues**: [File an Issue](https://github.com/nessi-dev/nessi/issues/new?template=documentation_improvement.md)
