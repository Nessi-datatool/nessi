# Security Policy

## Grant:  
- Free for personal, non-commercial use  
- Prohibited in enterprises, consulting, or paid environments  

## Commercial Use:  
Requires a paid license from [LakeDiff](https://lakediff.com).  

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

Please report security vulnerabilities to security@nessi.dev. We will acknowledge receipt of your vulnerability report and provide a more detailed response within 48 hours.

## Docker Security

### Container Security

- All containers run as non-root users
- Containers use read-only filesystems where possible
- Resource limits are enforced through Docker
- Regular security updates are applied to base images

### Network Security

- Containers use internal networks
- External access is restricted to necessary ports
- TLS/SSL is enforced for all external connections
- Rate limiting is implemented at the container level

### Data Security

- Sensitive data is encrypted at rest
- Credentials are managed through Docker secrets
- Temporary files are stored in memory
- Data volumes are properly mounted with appropriate permissions

## Security Features

### Authentication

- JWT-based authentication
- OAuth 2.0 support
- API key authentication
- Multi-factor authentication support

### Authorization

- Role-based access control (RBAC)
- Fine-grained permissions
- API endpoint access control
- Resource-level permissions

### Encryption

- TLS 1.3 for all communications
- AES-256 for data at rest
- Secure key management
- Certificate rotation

### Monitoring

- Security event logging
- Intrusion detection
- Anomaly detection
- Audit logging

## Best Practices

1. **Container Security**
   - Use official base images
   - Keep images updated
   - Scan for vulnerabilities
   - Use minimal images

2. **Network Security**
   - Use internal networks
   - Implement firewalls
   - Use TLS everywhere
   - Monitor network traffic

3. **Data Security**
   - Encrypt sensitive data
   - Use secure storage
   - Implement backups
   - Regular security audits

4. **Access Control**
   - Principle of least privilege
   - Regular access reviews
   - Strong authentication
   - Session management

## Security Updates

Security updates are released as needed. Critical updates are released within 24 hours of discovery. All updates are tested in Docker environments before release.

## Contact

- Security Team: security@nessi.dev
- Emergency Contact: security-emergency@nessi.dev
- PGP Key: [Available upon request]

## Security Measures

### General Security
- All code changes require internal review and approval
- Automated security scanning is performed on all changes
- Dependencies are regularly audited for security vulnerabilities
- Access to the codebase is restricted to authorized personnel
- Regular security audits are conducted

### Data Security
- All sensitive data is encrypted at rest
- Secure configuration storage
- Protection against unauthorized access
- Regular security updates

## Security Updates

Security updates will be released as patch versions (e.g., 1.0.1) and will be clearly marked in the changelog.

## Responsible Disclosure

We follow responsible disclosure practices:
1. Report the vulnerability privately
2. Allow time for the fix to be developed
3. Coordinate the release of the fix and disclosure
4. Credit the reporter (if desired)

## Security Best Practices

- Use strong, unique passwords
- Enable two-factor authentication
- Keep dependencies up to date
- Follow the principle of least privilege
- Review code changes carefully
- Maintain secure development environments
- Regular security training for development team 