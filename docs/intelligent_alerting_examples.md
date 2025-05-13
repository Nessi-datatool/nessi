# Intelligent Alerting System Examples

This document provides examples of how to use the Nessi Intelligent Alerting System in various scenarios.

## Basic Setup

Here's how to set up the Intelligent Alerting System with default settings:

```go
package main

import (
	"log"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

func main() {
	// Create alert manager
	alertManager, err := monitoring.CreateAlertManager()
	if err != nil {
		log.Fatalf("Failed to create alert manager: %v", err)
	}

	// Create monitor with intelligent alerting enabled
	options := monitoring.MonitorOptions{
		MetricsPort: 9090,
		AlertManager: alertManager,
		EnableIntelligentAlerting: true,
		// Use default configuration
		IntelligentAlertingConfig: nil,
	}

	monitor, err := monitoring.New(options)
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}

	// The intelligent alerting system is now running in the background
	// and will automatically create alert rules based on metrics patterns
	
	// Start your application...
}
```

## Custom Configuration

You can customize the Intelligent Alerting System behavior:

```go
// Create custom intelligent alerting configuration
intelligentConfig := &alerts.IntelligentAlertingConfig{
	MinimumDataPoints:      50,           // Require fewer data points
	AnalysisPeriod:         24 * time.Hour, // Look back 24 hours
	UpdateFrequency:        time.Hour,    // Update rules hourly
	Sensitivity:            0.8,          // Higher sensitivity (more alerts)
	EnableOutlierDetection: true,
	EnableTrendDeviation:   true,
	EnableSeasonalPatterns: false,        // Disable seasonal pattern detection
	AutoDisableUnusedRules: true,
	DisableThreshold:       15,           // Disable rules after 15 days of no triggers
}

// Create monitor with custom intelligent alerting config
options := monitoring.MonitorOptions{
	MetricsPort: 9090,
	AlertManager: alertManager,
	EnableIntelligentAlerting: true,
	IntelligentAlertingConfig: intelligentConfig,
}

monitor, err := monitoring.New(options)
```

## Accessing Intelligent Alert Manager

You can access the Intelligent Alert Manager to perform operations:

```go
// Get the intelligent alert manager
iam := monitor.GetIntelligentAlertManager()

// Get the current configuration
config := iam.GetConfig()
log.Printf("Current sensitivity: %f", config.Sensitivity)

// Manually trigger analysis
iam.StartAnalysis()
```

## Handling Different Types of Anomalies

The Intelligent Alerting System can detect different types of anomalies:

### 1. Outlier Detection Example

Outlier detection identifies metrics that fall outside their normal range:

```go
// Example metric data with an outlier
// Normal values around 100, with an outlier at 500
dataPoints := []alerts.MetricDataPoint{
	{Timestamp: time.Now().Add(-5 * time.Hour), Value: 98.5},
	{Timestamp: time.Now().Add(-4 * time.Hour), Value: 102.1},
	{Timestamp: time.Now().Add(-3 * time.Hour), Value: 97.8},
	{Timestamp: time.Now().Add(-2 * time.Hour), Value: 500.0}, // Outlier
	{Timestamp: time.Now().Add(-1 * time.Hour), Value: 101.2},
}

// The system will automatically create rules like:
// - Alert when metric > mean + 3*stdDev (adjusted by sensitivity)
// - Alert when metric < mean - 3*stdDev (adjusted by sensitivity)
```

### 2. Trend Deviation Example

Trend deviation detects when metrics deviate from their expected trend:

```go
// Example metric data with an upward trend
// Values increasing by ~10 per hour, with a deviation
dataPoints := []alerts.MetricDataPoint{
	{Timestamp: time.Now().Add(-5 * time.Hour), Value: 100},
	{Timestamp: time.Now().Add(-4 * time.Hour), Value: 110},
	{Timestamp: time.Now().Add(-3 * time.Hour), Value: 120},
	{Timestamp: time.Now().Add(-2 * time.Hour), Value: 130},
	{Timestamp: time.Now().Add(-1 * time.Hour), Value: 90}, // Deviation from trend
}

// The system will automatically create rules to detect deviations from the trend
```

### 3. Seasonal Pattern Example

Seasonal pattern detection identifies recurring patterns and alerts on deviations:

```go
// Example metric data with daily patterns (higher during business hours)
// Values around 100 during business hours, 50 otherwise
dataPoints := []alerts.MetricDataPoint{
	// Day 1
	{Timestamp: time.Date(2025, 5, 12, 3, 0, 0, 0, time.UTC), Value: 48},  // 3 AM
	{Timestamp: time.Date(2025, 5, 12, 10, 0, 0, 0, time.UTC), Value: 102}, // 10 AM
	{Timestamp: time.Date(2025, 5, 12, 15, 0, 0, 0, time.UTC), Value: 105}, // 3 PM
	{Timestamp: time.Date(2025, 5, 12, 22, 0, 0, 0, time.UTC), Value: 51},  // 10 PM
	
	// Day 2
	{Timestamp: time.Date(2025, 5, 13, 3, 0, 0, 0, time.UTC), Value: 52},  // 3 AM
	{Timestamp: time.Date(2025, 5, 13, 10, 0, 0, 0, time.UTC), Value: 98},  // 10 AM
	{Timestamp: time.Date(2025, 5, 13, 15, 0, 0, 0, time.UTC), Value: 30},  // 3 PM - Anomaly!
	{Timestamp: time.Date(2025, 5, 13, 22, 0, 0, 0, time.UTC), Value: 49},  // 10 PM
}

// The system will automatically create rules for each hour of the day
// based on the historical patterns observed
```

## Integration with Sudden Change Detection

The Intelligent Alerting System integrates with the existing Sudden Change Detection feature to provide comprehensive anomaly detection:

```go
// The system can detect:
// 1. Constant changes (gradual increases/decreases over time)
// 2. Sudden spikes and dips (temporary anomalies)
// 3. Oscillating patterns (alternating values)

// Example of detecting a sudden spike
dataPoints := []alerts.MetricDataPoint{
	{Timestamp: time.Now().Add(-5 * time.Hour), Value: 100},
	{Timestamp: time.Now().Add(-4 * time.Hour), Value: 102},
	{Timestamp: time.Now().Add(-3 * time.Hour), Value: 98},
	{Timestamp: time.Now().Add(-2 * time.Hour), Value: 350}, // Sudden spike
	{Timestamp: time.Now().Add(-1 * time.Hour), Value: 103},
}

// The system will create appropriate alert rules to detect these patterns
```

## Viewing Intelligent Alerts in Grafana

The Intelligent Alerting System integrates with Grafana dashboards:

1. Open the Nessi Intelligent Alerting Dashboard in Grafana
2. View panels showing:
   - Intelligent Alerts by Strategy
   - Intelligent Rules Created
   - Active Intelligent Rules
   - Alert Effectiveness by Strategy
   - Top Metrics with Intelligent Rules

## Programmatic Access via Python API

You can also access the Intelligent Alerting System via the Python API:

```python
from nessi.client import NessiClient
from nessi.models import AlertRule

# Connect to Nessi
client = NessiClient(url="http://localhost:8080", api_key="your-api-key")

# Get all intelligent alert rules
intelligent_rules = client.alerts.get_rules(labels={"type": "intelligent"})

# Get rules by strategy
outlier_rules = client.alerts.get_rules(labels={
    "type": "intelligent", 
    "strategy": "outlier"
})

# Get active alerts from intelligent rules
alerts = client.alerts.get_alerts(labels={"type": "intelligent"})

# Print alert details
for alert in alerts:
    print(f"Alert: {alert.name}")
    print(f"  Severity: {alert.severity}")
    print(f"  Strategy: {alert.labels.get('strategy')}")
    print(f"  Metric: {alert.metric}")
    print(f"  Value: {alert.value}")
    print(f"  Threshold: {alert.threshold}")
```

## Best Practices

1. **Start with Default Configuration**: Begin with the default configuration and adjust as needed based on your specific requirements.

2. **Adjust Sensitivity**: If you're getting too many false positives, lower the sensitivity. If you're missing important alerts, increase it.

3. **Combine with Manual Rules**: Use intelligent alerting alongside manually created rules for critical systems.

4. **Monitor Rule Effectiveness**: Use the Grafana dashboard to monitor how effective your intelligent alert rules are and adjust as needed.

5. **Regular Maintenance**: Periodically review the rules generated by the system to ensure they align with your monitoring needs.

6. **Integrate with Notification Channels**: Configure appropriate notification channels (email, Slack, webhook) to receive alerts when intelligent rules are triggered.
