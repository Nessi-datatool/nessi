package alerts

// GetRulesByLabels returns all rules that match the given labels
func (am *AlertManager) GetRulesByLabels(labels map[string]string) ([]*AlertRule, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var matchingRules []*AlertRule

	for _, rule := range am.rules {
		if matchesLabels(rule.Labels, labels) {
			matchingRules = append(matchingRules, rule)
		}
	}

	return matchingRules, nil
}

// GetAlertsByLabels returns all alerts that match the given labels
func (am *AlertManager) GetAlertsByLabels(labels map[string]string) ([]*Alert, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var matchingAlerts []*Alert

	for _, alert := range am.alerts {
		if matchesLabels(alert.Labels, labels) {
			matchingAlerts = append(matchingAlerts, alert)
		}
	}

	return matchingAlerts, nil
}

// GetNotifiers returns all registered notifiers
func (am *AlertManager) GetNotifiers() []string {
	am.mu.RLock()
	defer am.mu.RUnlock()

	// Create a list of notifier names
	notifiers := make([]string, 0, len(am.notifiers))
	for name := range am.notifiers {
		notifiers = append(notifiers, name)
	}

	return notifiers
}

// Helper function to check if a map contains all key-value pairs from another map
func matchesLabels(itemLabels, filterLabels map[string]string) bool {
	if len(filterLabels) == 0 {
		return true
	}

	if itemLabels == nil {
		return false
	}

	for k, v := range filterLabels {
		itemValue, ok := itemLabels[k]
		if !ok || itemValue != v {
			return false
		}
	}

	return true
}
