# Nessi Security Guide

This guide provides comprehensive information about security considerations, best practices, and features in the Nessi project.

## Table of Contents

- [Overview](#overview)
- [Security Features](#security-features)
- [Authentication](#authentication)
- [Authorization](#authorization)
- [Data Protection](#data-protection)
- [Network Security](#network-security)
- [License Security](#license-security)
- [Secure Configuration](#secure-configuration)
- [Security Best Practices](#security-best-practices)
- [Vulnerability Management](#vulnerability-management)
- [Audit and Compliance](#audit-and-compliance)
- [Security FAQs](#security-faqs)

## Overview

Security is a fundamental aspect of the Nessi project. This guide outlines the security features, considerations, and best practices for using Nessi in a secure manner.

## Security Features

### Core Security Features (Community Edition)

- Basic access controls for local resources
- Secure configuration management
- Credential protection mechanisms
- Audit logging for operations
- Input validation and sanitization
- Secure error handling

### Premium Security Features (Pro Edition)

- Integration with enterprise authentication systems
- Enhanced audit logging and reporting
- Advanced access controls for cloud resources
- Secure integration with Databricks and AWS S3
- Compliance reporting capabilities

## Authentication

Nessi supports various authentication mechanisms depending on the integration:

### Local Authentication

For local operations, Nessi relies on the operating system's authentication and file permissions.

### Databricks Authentication (Pro Edition)

For Databricks integration, Nessi supports:

- Personal Access Tokens (PAT)
- OAuth 2.0 authentication
- Service Principal authentication

Configuration example:
```yaml
integrations:
  databricks:
    auth_type: pat
    host: your-databricks-instance.cloud.databricks.com
    token: dapi_xxxxxxxxxxxxxxxxxxxxxxxx
```

Security best practices:
- Use environment variables for tokens instead of configuration files
- Create dedicated tokens with minimal permissions
- Rotate tokens regularly
- Use service principals for automated workflows

### AWS Authentication (Pro Edition)

For AWS S3 integration, Nessi supports:

- AWS Access Keys
- IAM Roles
- AWS Profiles
- Environment variables

Configuration example:
```yaml
integrations:
  aws:
    auth_type: access_key
    region: us-west-2
    access_key_id: AKIAXXXXXXXXXXXXXXXX
    secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Security best practices:
- Use IAM roles instead of access keys when possible
- Create dedicated IAM users with minimal permissions
- Enable MFA for IAM users
- Use temporary credentials when possible

## Authorization

Nessi implements authorization controls to ensure that operations are performed with appropriate permissions:

### File System Authorization

For local file operations, Nessi respects the operating system's file permissions.

### Databricks Authorization (Pro Edition)

When integrating with Databricks, Nessi respects the Databricks permission model:

- Table ACLs
- Workspace permissions
- Cluster permissions

### AWS Authorization (Pro Edition)

When integrating with AWS S3, Nessi respects the AWS permission model:

- IAM policies
- Bucket policies
- ACLs
- S3 Block Public Access settings

## Data Protection

Nessi includes several features to protect sensitive data:

### Credential Protection

- Credentials are never stored in plain text
- Configuration files with credentials are protected with appropriate permissions
- Environment variables are recommended for credential storage
- Memory containing credentials is securely wiped after use

### Data Encryption

- Support for encrypted Delta Lake tables
- Integration with AWS S3 server-side encryption
- TLS for all network communications

### Data Masking

- Support for masking sensitive data in reports and logs
- Configuration options for controlling data visibility

Example configuration:
```yaml
security:
  data_masking:
    enabled: true
    patterns:
      - type: regex
        pattern: "\d{4}-\d{4}-\d{4}-\d{4}"
        replacement: "XXXX-XXXX-XXXX-XXXX"
      - type: column
        names: ["email", "phone", "ssn"]
        replacement: "[REDACTED]"
```

## Network Security

Nessi implements network security best practices:

### TLS/SSL

- All network communications use TLS 1.2 or higher
- Certificate validation is enforced
- Support for custom certificate authorities

### Firewall Considerations

- Minimal network footprint
- Documented network requirements for integrations
- Support for proxy configurations

Example proxy configuration:
```yaml
network:
  proxy:
    http_proxy: http://proxy.example.com:8080
    https_proxy: http://proxy.example.com:8080
    no_proxy: localhost,127.0.0.1
```

## License Security

The license management system includes several security measures:

### License Validation

- License keys are validated cryptographically
- Online validation with the license server
- Offline validation with pre-generated license files

### Machine Binding

- Licenses can be bound to specific machines
- Machine IDs are generated using hardware characteristics
- Limited number of activations per license

### Anti-Tampering

- Binary integrity verification
- Runtime integrity checks
- Protection against license bypass attempts

## Secure Configuration

Nessi provides several options for secure configuration:

### Configuration File Security

- Configuration files should have restricted permissions
- Sensitive configuration should be stored in environment variables
- Support for encrypted configuration values

### Environment Variables

Recommended environment variables for secure configuration:

```bash
# Databricks credentials
export DATABRICKS_HOST=your-databricks-instance.cloud.databricks.com
export DATABRICKS_TOKEN=dapi_xxxxxxxxxxxxxxxxxxxxxxxx

# AWS credentials
export AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AWS_REGION=us-west-2

# Nessi license
export NESSI_LICENSE_KEY=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

### Secrets Management Integration

For enterprise environments, consider integrating with secrets management solutions:

- HashiCorp Vault
- AWS Secrets Manager
- Azure Key Vault
- Google Secret Manager

## Security Best Practices

### Installation Security

- Verify download integrity using checksums
- Install from trusted sources only
- Keep Nessi updated to the latest version

### Operational Security

- Run Nessi with minimal privileges
- Use dedicated service accounts for automated operations
- Implement proper log management and monitoring
- Regularly review audit logs

### Integration Security

- Use dedicated credentials for each integration
- Implement proper credential rotation
- Monitor integration access patterns
- Use network segmentation where appropriate

### Report Security

- Control access to generated reports
- Consider sensitivity of data in reports
- Implement proper retention policies for reports
- Use secure channels for report distribution

## Vulnerability Management

Nessi follows a structured approach to vulnerability management:

### Security Updates

- Regular security updates
- Critical vulnerabilities are patched promptly
- Security advisories are published for known issues

### Dependency Management

- Regular updates of dependencies
- Vulnerability scanning of dependencies
- Dependency lockfiles to prevent unexpected changes

### Reporting Vulnerabilities

If you discover a security vulnerability in Nessi:

1. **Do not** disclose it publicly
2. Email security@nessi.dev with details
3. Include steps to reproduce the vulnerability
4. If possible, include a proposed fix

## Audit and Compliance

Nessi includes features to support audit and compliance requirements:

### Audit Logging

- Comprehensive audit logging of all operations
- Configurable log levels and destinations
- Structured logs for easy parsing and analysis

Example configuration:
```yaml
logging:
  audit:
    enabled: true
    level: info
    destination: file
    file_path: /var/log/nessi/audit.log
    format: json
    retention:
      days: 90
      max_size_mb: 1024
```

### Compliance Reporting (Pro Edition)

- Pre-built compliance reports
- Custom compliance report templates
- Integration with compliance frameworks

### Data Lineage

- Tracking of data sources and transformations
- Documentation of data flows
- Support for data governance requirements

## Security FAQs

### Q: Is Nessi secure by default?

A: Yes, Nessi follows secure-by-default principles. Default configurations prioritize security over convenience, and sensitive operations require explicit configuration.

### Q: How does Nessi protect sensitive data?

A: Nessi implements multiple layers of protection for sensitive data, including credential protection, data encryption, and data masking. Configuration files with sensitive information are protected with appropriate permissions, and environment variables are recommended for credential storage.

### Q: Does Nessi support encryption?

A: Yes, Nessi supports encryption for data at rest (through Delta Lake encryption and AWS S3 server-side encryption) and data in transit (through TLS for all network communications).

### Q: How does Nessi handle authentication for cloud services?

A: Nessi supports various authentication mechanisms for cloud services, including API tokens, access keys, and IAM roles. These credentials are handled securely and can be provided through environment variables or secure configuration files.

### Q: Is Nessi compliant with industry standards?

A: Nessi is designed with security best practices in mind and can be used in environments that require compliance with various standards. However, compliance ultimately depends on how Nessi is deployed and used in your environment.

### Q: How are security vulnerabilities handled?

A: Security vulnerabilities are treated with the highest priority. Critical vulnerabilities are patched promptly, and security advisories are published for known issues. Users are encouraged to keep Nessi updated to the latest version.

### Q: Can I run Nessi in a restricted network environment?

A: Yes, Nessi can be run in restricted network environments. It has minimal network requirements and supports proxy configurations. The Pro Edition license can be activated offline using pre-generated license files.

---

For more information on Nessi security, please visit [nessi.dev/security](https://nessi.dev/security) or contact security@nessi.dev.
