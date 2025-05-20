# Verify Security Script

This script provides a simple verification of the SSL certificate generation functionality in the Nessi.dev security package.

## Functionality

The script:
1. Creates a temporary directory for test files
2. Configures the SSL certificate manager with auto-generation enabled
3. Generates a TLS configuration
4. Verifies that the certificate and key files were created successfully

## Usage

```bash
cd /Users/meisi/Documents/nessi-dev
go run ./scripts/verify_security_pkg/verify_security.go
```

## Output

The script will output the status of each verification step, with checkmarks (✅) for successful steps and error messages for any failures.

## Integration with Freshness SLA

This script can be used alongside the Freshness SLA monitoring system to verify that the security components are working correctly. The Freshness SLA system requires secure connections when `RequireHTTPS` is enabled in the configuration, and this script helps ensure that the certificate generation is functioning properly.
