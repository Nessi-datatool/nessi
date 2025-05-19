// Implementation of monitoring and alerting features.

package alerts

// createRule is a wrapper for the AlertManager.CreateRule method
func (iam *IntelligentAlertManager) createRule(rule *AlertRule) error {
	return iam.alertManager.CreateRule(rule)
}

// updateRule is a wrapper for the AlertManager.UpdateRule method
func (iam *IntelligentAlertManager) updateRule(rule *AlertRule) error {
	return iam.alertManager.UpdateRule(rule)
}
