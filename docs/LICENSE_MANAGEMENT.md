# License Management System

Nessi includes a robust license management system that controls access to premium features while providing a free trial option for evaluation.

## License Tiers

Nessi offers the following license tiers:

| Feature | Community | Pro |
|---------|-----------|-----|
| Core Delta Lake Management | ✅ | ✅ |
| Data Quality Checks | ✅ | ✅ |
| HTML/PDF Reporting | ✅ | ✅ |
| CLI Monitoring | ✅ | ✅ |
| Databricks Integration | ❌ | ✅ |
| AWS S3 Storage | ❌ | ✅ |
| Azure Blob Storage | ❌ | ✅ |
| Google Cloud Storage | ❌ | ✅ |
| Data Catalog Integration | ❌ | ✅ |
| dbt Integration | ❌ | ✅ |
| Workflow Orchestration | ❌ | ✅ |

> **Note**: For Enterprise-level support and custom features, please visit the [Nessi website](https://nessi.dev) for contact information.

## Free Trial

Nessi offers a 1-month free trial that provides access to all premium features, including those in the Enterprise tier. Key aspects of the trial system:

- **Duration**: 30 days from activation
- **Limitations**: Maximum of 2 trials per machine
- **Security**: Trial information is cryptographically signed to prevent tampering

## License Commands

### Starting a Trial

```bash
# Start a free trial to access premium features
nessi license start-trial
```

### Checking License Status

```bash
# View your current license information
nessi license info
```

Example output:
```
License Status: Trial
Plan: Enterprise
Is Trial: Yes
Expires: 2025-06-22
Days Remaining: 30
Message: Free trial active. 30 days remaining.
```

### Activating a License

```bash
# Activate a purchased license
nessi license activate --key YOUR_LICENSE_KEY
```

## Technical Implementation

The license management system uses several security measures to ensure the integrity of license validation:

1. **Machine ID**: Each machine is uniquely identified using hardware characteristics
2. **Signature Verification**: All license and trial data is signed using HMAC-SHA256
3. **Anti-Tampering**: Multiple validation checks prevent binary modification
4. **Secure Storage**: License information is stored securely with access controls

## Troubleshooting

### Trial Expiration

When a trial expires, premium features will no longer be accessible. You'll need to purchase a license to continue using these features.

### Maximum Trials Reached

If you've already used 2 trials on your machine, you won't be able to start another trial. Consider purchasing a license to continue using premium features.

### Invalid License

If your license is reported as invalid, check the following:
- License key is entered correctly
- License hasn't expired
- License is for the correct machine ID

## Support

For licensing questions or issues, contact support at licensing@nessi-dev.com.
