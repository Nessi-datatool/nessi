# Nessi Monitoring System License Guide

## License Overview

The Nessi Monitoring System operates under a trial license system:
- Initial trial period: 14 days
- Requires activation after trial period
- License keys are tied to specific instances
- License management is handled externally

## Trial Period

### Features During Trial
- Full access to all monitoring features
- Complete dashboard functionality
- All metrics collection and visualization
- Alert configuration and management

### Trial Limitations
- 14-day duration from first activation
- System will continue to operate after trial
- Metrics collection continues
- Dashboard remains accessible
- Alerts remain active

## License Activation

### Obtaining a License Key
1. Contact Nessi support at support@nessi.com
2. Provide your instance details:
   - Instance ID
   - Organization name
   - Contact information
   - Intended usage

### Activating Your License
1. Access the license management page:
   ```
   http://localhost:3000/license
   ```

2. Enter your license key:
   - 32-character alphanumeric key
   - Case-sensitive
   - No spaces or special characters

3. Verify activation:
   - Check license status
   - Confirm expiration date
   - Verify all features

### License Status Check
```bash
# Check license status
curl http://localhost:8000/license/status

# Expected response
{
    "status": "active",
    "expiration_date": "2024-04-15",
    "features": ["metrics", "dashboard", "alerts"]
}
```

## License Management

### Viewing License Information
1. Access the license dashboard
2. View current status
3. Check expiration date
4. Monitor usage statistics

### Renewing Your License
1. Contact support before expiration
2. Request new license key
3. Follow activation process
4. Verify renewal status

### License Transfer
1. Contact support for transfer
2. Provide new instance details
3. Deactivate old license
4. Activate on new instance

## Troubleshooting

### Common Issues

1. **Invalid License Key**
   - Verify key format
   - Check for typos
   - Contact support

2. **Activation Failed**
   - Check internet connection
   - Verify instance details
   - Contact support

3. **License Expired**
   - Contact support for renewal
   - Request new license key
   - Follow activation process

### Support Contact
- Email: support@nessi.com
- Phone: +1 (555) 123-4567
- Hours: 9 AM - 5 PM EST

## Best Practices

1. **License Management**
   - Keep license key secure
   - Monitor expiration date
   - Plan for renewal
   - Document activation process

2. **Backup**
   - Store license key securely
   - Document activation details
   - Keep support contact information

3. **Security**
   - Don't share license keys
   - Use secure channels for communication
   - Report suspicious activity

## FAQ

### Q: What happens after the trial period?
A: The system continues to operate, but requires a valid license key for continued use.

### Q: Can I extend my trial period?
A: Trial extensions are available on a case-by-case basis. Contact support for details.

### Q: How do I transfer my license?
A: Contact support for license transfer assistance. Transfers require proper documentation.

### Q: What if I lose my license key?
A: Contact support with your instance details to recover your license information.

### Q: Can I use multiple licenses?
A: Each instance requires its own license. Contact support for multi-instance licensing. 