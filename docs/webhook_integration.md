# Webhook Integration

## Overview

Nessi.dev provides comprehensive webhook integration capabilities, allowing external systems to receive notifications when specific events occur in the system. This document outlines the implementation details and usage of the webhook integration feature.

## Features

- **Event-Based Notifications**: Receive notifications for various events such as data validation results, alerts, rule changes, and more
- **Configurable Endpoints**: Configure webhook endpoints with custom headers and authentication
- **Flexible Event Filtering**: Subscribe to specific events or use wildcards to receive all notifications
- **Retry Mechanism**: Automatic retry with configurable attempts and delay for failed webhook deliveries
- **CLI Management**: Comprehensive CLI commands for managing webhooks

## Supported Events

The following events are supported by the webhook integration:

- `data.profile.created`: Triggered when a new data profile is created
- `data.profile.updated`: Triggered when a data profile is updated
- `data.validation.failed`: Triggered when data validation fails
- `data.validation.passed`: Triggered when data validation passes
- `alert.triggered`: Triggered when an alert is triggered
- `alert.resolved`: Triggered when an alert is resolved
- `rule.created`: Triggered when a new rule is created
- `rule.updated`: Triggered when a rule is updated
- `rule.deleted`: Triggered when a rule is deleted
- `metric.threshold.exceeded`: Triggered when a metric exceeds a threshold
- `system.startup`: Triggered when the system starts up
- `system.shutdown`: Triggered when the system shuts down

## Webhook Payload

When an event occurs, Nessi.dev sends a webhook payload to the configured endpoint. The payload includes:

```json
{
  "id": "event-id",
  "event_type": "event.type",
  "timestamp": "2025-05-13T12:45:27+02:00",
  "payload": {
    // Event-specific data
  }
}
```

## CLI Commands

### List Webhooks

List all registered webhooks:

```bash
nessi webhook list
```

Options:
- `--format`: Output format (table, json)

### Create Webhook

Create a new webhook:

```bash
nessi webhook create --name "My Webhook" --url "https://example.com/webhook" --events "data.validation.failed,alert.triggered" --description "Webhook for validation failures and alerts"
```

Options:
- `--name`: Name of the webhook
- `--url`: URL to send webhook events to
- `--events`: Comma-separated list of event types to subscribe to (use `*` for all events)
- `--headers`: Comma-separated list of headers in key:value format
- `--description`: Description of the webhook

### Get Webhook Details

Get details of a specific webhook:

```bash
nessi webhook get webhook-id
```

Options:
- `--format`: Output format (detailed, json)

### Update Webhook

Update an existing webhook:

```bash
nessi webhook update webhook-id --name "Updated Webhook" --events "data.validation.failed,alert.triggered,rule.created"
```

Options:
- `--name`: Name of the webhook
- `--url`: URL to send webhook events to
- `--events`: Comma-separated list of event types to subscribe to
- `--headers`: Comma-separated list of headers in key:value format
- `--description`: Description of the webhook
- `--retry-count`: Number of retry attempts
- `--retry-delay`: Delay between retry attempts in seconds
- `--enabled`: Whether the webhook is enabled

### Enable/Disable Webhook

Enable or disable a webhook:

```bash
nessi webhook enable webhook-id
nessi webhook disable webhook-id
```

### Delete Webhook

Delete a webhook:

```bash
nessi webhook delete webhook-id
```

### List Supported Events

List all supported event types:

```bash
nessi webhook events
```

### Test Webhook

Test a webhook by sending a test event:

```bash
nessi webhook test webhook-id event.type
```

## Implementation Details

### Webhook Manager

The `WebhookManager` is responsible for managing webhooks and their delivery:

- Registers and unregisters webhooks
- Manages webhook configurations
- Sends events to subscribed webhooks
- Handles retries for failed webhook deliveries

### Webhook Storage

Webhook configurations are stored using one of the following storage implementations:

- **FileStorage**: Stores webhook configurations in JSON files
- **MemoryStorage**: Stores webhook configurations in memory (for testing)

### Webhook Service

The `WebhookService` provides a high-level interface for working with webhooks:

- Creates, updates, and deletes webhooks
- Lists webhooks and their details
- Enables and disables webhooks
- Triggers events and sends them to subscribed webhooks

## Integration with Other Components

The webhook integration is designed to work seamlessly with other Nessi.dev components:

- **Data Quality Engine**: Sends validation results to webhooks
- **Alerting System**: Notifies webhooks when alerts are triggered or resolved
- **Rule Management**: Informs webhooks about rule changes
- **Metrics System**: Sends metric threshold events to webhooks

## Security Considerations

When using webhooks, consider the following security best practices:

- Use HTTPS endpoints to ensure secure communication
- Implement authentication for webhook endpoints (e.g., using API keys in headers)
- Validate the source of webhook requests in your receiving system
- Limit the information included in webhook payloads to what is necessary

## Example: Setting Up a Webhook for Failed Validations

```bash
# Create a webhook for failed validations
nessi webhook create \
  --name "Validation Failures" \
  --url "https://example.com/webhook" \
  --events "data.validation.failed" \
  --headers "Authorization:Bearer token123,Content-Type:application/json" \
  --description "Webhook for validation failures"

# Test the webhook
nessi webhook test webhook-id data.validation.failed
```

## Future Enhancements

Planned enhancements for webhook integration include:

1. Webhook signature verification for enhanced security
2. Event payload customization
3. Webhook delivery logs and history
4. Batch delivery for high-volume events
5. Integration with popular notification services (Slack, Microsoft Teams, etc.)
