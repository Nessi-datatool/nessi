package config

import (
	"fmt"
	"sync"
)

// FeatureFlag represents a feature flag that can be toggled
type FeatureFlag string

const (
	// FeatureFlagWebhooks enables webhook functionality
	FeatureFlagWebhooks FeatureFlag = "webhooks"
	// FeatureFlagPlugins enables plugin functionality
	FeatureFlagPlugins FeatureFlag = "plugins"
	// FeatureFlagAdvancedMetrics enables advanced metrics collection
	FeatureFlagAdvancedMetrics FeatureFlag = "advanced_metrics"
	// FeatureFlagIntelligentAlerting enables intelligent alerting
	FeatureFlagIntelligentAlerting FeatureFlag = "intelligent_alerting"
	// FeatureFlagMultiFormat enables multi-format support
	FeatureFlagMultiFormat FeatureFlag = "multi_format"
)

// FeatureFlagManager manages feature flags
type FeatureFlagManager struct {
	flags map[FeatureFlag]bool
	mu    sync.RWMutex
}

// NewFeatureFlagManager creates a new feature flag manager
func NewFeatureFlagManager() *FeatureFlagManager {
	return &FeatureFlagManager{
		flags: make(map[FeatureFlag]bool),
	}
}

// SetFlag sets a feature flag
func (m *FeatureFlagManager) SetFlag(flag FeatureFlag, enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flags[flag] = enabled
}

// GetFlag gets a feature flag
func (m *FeatureFlagManager) GetFlag(flag FeatureFlag) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	enabled, exists := m.flags[flag]
	if !exists {
		return false // Default to disabled if not set
	}
	return enabled
}

// ListFlags lists all feature flags
func (m *FeatureFlagManager) ListFlags() map[FeatureFlag]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Create a copy to avoid race conditions
	flags := make(map[FeatureFlag]bool)
	for k, v := range m.flags {
		flags[k] = v
	}
	
	return flags
}

// IsEnabled checks if a feature flag is enabled
func (m *FeatureFlagManager) IsEnabled(flag FeatureFlag) bool {
	return m.GetFlag(flag)
}

// LoadFromConfig loads feature flags from configuration
func (m *FeatureFlagManager) LoadFromConfig(config map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	featureFlags, ok := config["feature_flags"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("feature_flags configuration not found or invalid")
	}
	
	for flagName, enabled := range featureFlags {
		enabledBool, ok := enabled.(bool)
		if !ok {
			return fmt.Errorf("feature flag %s has invalid value", flagName)
		}
		m.flags[FeatureFlag(flagName)] = enabledBool
	}
	
	return nil
}

// SaveToConfig saves feature flags to configuration
func (m *FeatureFlagManager) SaveToConfig() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	featureFlags := make(map[string]interface{})
	for flag, enabled := range m.flags {
		featureFlags[string(flag)] = enabled
	}
	
	return map[string]interface{}{
		"feature_flags": featureFlags,
	}
}
