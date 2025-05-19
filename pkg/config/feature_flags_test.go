package config

import (
	"testing"
)

func TestFeatureFlagManager(t *testing.T) {
	// Create a new feature flag manager
	manager := NewFeatureFlagManager()

	// Test setting and getting flags
	t.Run("SetAndGetFlag", func(t *testing.T) {
		// Set a flag
		manager.SetFlag(FeatureFlagWebhooks, true)

		// Get the flag
		enabled := manager.GetFlag(FeatureFlagWebhooks)
		if !enabled {
			t.Errorf("Expected flag to be enabled, but it was disabled")
		}

		// Set the flag to disabled
		manager.SetFlag(FeatureFlagWebhooks, false)

		// Get the flag again
		enabled = manager.GetFlag(FeatureFlagWebhooks)
		if enabled {
			t.Errorf("Expected flag to be disabled, but it was enabled")
		}
	})

	// Test getting a non-existent flag
	t.Run("GetNonExistentFlag", func(t *testing.T) {
		// Get a flag that hasn't been set
		enabled := manager.GetFlag(FeatureFlagPlugins)

		// It should be disabled by default
		if enabled {
			t.Errorf("Expected non-existent flag to be disabled by default, but it was enabled")
		}
	})

	// Test listing flags
	t.Run("ListFlags", func(t *testing.T) {
		// Set some flags
		manager.SetFlag(FeatureFlagWebhooks, true)
		manager.SetFlag(FeatureFlagPlugins, false)
		manager.SetFlag(FeatureFlagAdvancedMetrics, true)

		// List the flags
		flags := manager.ListFlags()

		// Check that the flags are correct
		if len(flags) != 3 {
			t.Errorf("Expected 3 flags, but got %d", len(flags))
		}

		if !flags[FeatureFlagWebhooks] {
			t.Errorf("Expected webhooks flag to be enabled, but it was disabled")
		}

		if flags[FeatureFlagPlugins] {
			t.Errorf("Expected plugins flag to be disabled, but it was enabled")
		}

		if !flags[FeatureFlagAdvancedMetrics] {
			t.Errorf("Expected advanced metrics flag to be enabled, but it was disabled")
		}
	})

	// Test IsEnabled
	t.Run("IsEnabled", func(t *testing.T) {
		// Set a flag
		manager.SetFlag(FeatureFlagMultiFormat, true)

		// Check if it's enabled
		if !manager.IsEnabled(FeatureFlagMultiFormat) {
			t.Errorf("Expected flag to be enabled, but it was disabled")
		}

		// Set it to disabled
		manager.SetFlag(FeatureFlagMultiFormat, false)

		// Check if it's disabled
		if manager.IsEnabled(FeatureFlagMultiFormat) {
			t.Errorf("Expected flag to be disabled, but it was enabled")
		}
	})

	// Test loading from and saving to config
	t.Run("LoadAndSaveConfig", func(t *testing.T) {
		// Create a config
		config := map[string]interface{}{
			"feature_flags": map[string]interface{}{
				"webhooks":             true,
				"plugins":              false,
				"advanced_metrics":     true,
				"intelligent_alerting": false,
			},
		}

		// Create a new manager
		manager := NewFeatureFlagManager()

		// Load from config
		err := manager.LoadFromConfig(config)
		if err != nil {
			t.Errorf("Failed to load from config: %v", err)
		}

		// Check that the flags were loaded correctly
		if !manager.GetFlag(FeatureFlagWebhooks) {
			t.Errorf("Expected webhooks flag to be enabled, but it was disabled")
		}

		if manager.GetFlag(FeatureFlagPlugins) {
			t.Errorf("Expected plugins flag to be disabled, but it was enabled")
		}

		if !manager.GetFlag(FeatureFlagAdvancedMetrics) {
			t.Errorf("Expected advanced metrics flag to be enabled, but it was disabled")
		}

		if manager.GetFlag(FeatureFlagIntelligentAlerting) {
			t.Errorf("Expected intelligent alerting flag to be disabled, but it was enabled")
		}

		// Save to config
		savedConfig := manager.SaveToConfig()

		// Check that the config was saved correctly
		featureFlags, ok := savedConfig["feature_flags"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected feature_flags to be a map, but it wasn't")
		}

		if len(featureFlags) != 4 {
			t.Errorf("Expected 4 feature flags, but got %d", len(featureFlags))
		}

		if webhooks, ok := featureFlags["webhooks"].(bool); !ok || !webhooks {
			t.Errorf("Expected webhooks flag to be true, but it was %v", featureFlags["webhooks"])
		}

		if plugins, ok := featureFlags["plugins"].(bool); !ok || plugins {
			t.Errorf("Expected plugins flag to be false, but it was %v", featureFlags["plugins"])
		}

		if advancedMetrics, ok := featureFlags["advanced_metrics"].(bool); !ok || !advancedMetrics {
			t.Errorf("Expected advanced_metrics flag to be true, but it was %v", featureFlags["advanced_metrics"])
		}

		if intelligentAlerting, ok := featureFlags["intelligent_alerting"].(bool); !ok || intelligentAlerting {
			t.Errorf("Expected intelligent_alerting flag to be false, but it was %v", featureFlags["intelligent_alerting"])
		}
	})
}
