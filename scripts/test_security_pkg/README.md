# Test Security Script

This script runs a comprehensive set of tests on the Nessi.dev security package and provides detailed results.

## Functionality

The script tests:
1. Auth manager creation
2. Cert manager creation
3. User creation
4. Authentication
5. Token validation
6. API key generation
7. API key authentication
8. Admin user creation
9. Role middleware
10. SSL certificate generation

## Usage

```bash
cd /Users/meisi/Documents/nessi-dev
go run ./scripts/test_security_pkg/test_security.go
```

## Output

The script will output detailed test results for each test case, including:
- Test name
- Pass/Fail status
- Description of the test
- Error details (if any)

At the end, it will provide a summary of the total tests, passed tests, and failed tests.

## Integration with Freshness SLA and Intelligent Alerting

This script ensures that the security components used by the Freshness SLA monitoring and intelligent alerting systems are functioning correctly. Proper authentication and authorization are critical for these systems, especially when they're configured to send alerts or take automated actions based on data freshness status.
