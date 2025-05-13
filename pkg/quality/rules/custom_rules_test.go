package rules

import (
	"testing"
)

func TestCustomRuleRegistry(t *testing.T) {
	// Create a new custom rule registry
	registry := NewCustomRuleRegistry()

	// Create a test rule
	testRule := &CustomRule{
		Metadata: CustomRuleMetadata{
			Name:        "Test Rule",
			Description: "A test rule for unit testing",
			Author:      "Test Author",
			Version:     "1.0.0",
			Parameters: map[string]interface{}{
				"threshold": 0.5,
			},
			Examples: []string{
				"Example 1",
				"Example 2",
			},
		},
		Func: func(data interface{}, params map[string]interface{}) (bool, string, error) {
			// Simple rule that checks if a number is greater than a threshold
			value, ok := data.(float64)
			if !ok {
				return false, "Data is not a number", nil
			}

			threshold, ok := params["threshold"].(float64)
			if !ok {
				threshold = 0.5 // Default threshold
			}

			if value > threshold {
				return true, "Value is greater than threshold", nil
			} else {
				return false, "Value is not greater than threshold", nil
			}
		},
	}

	// Test registering a rule
	t.Run("RegisterRule", func(t *testing.T) {
		err := registry.RegisterRule("test-rule", testRule)
		if err != nil {
			t.Errorf("Failed to register rule: %v", err)
		}

		// Try to register the same rule again
		err = registry.RegisterRule("test-rule", testRule)
		if err == nil {
			t.Errorf("Expected error when registering duplicate rule, but got nil")
		}
	})

	// Test getting a rule
	t.Run("GetRule", func(t *testing.T) {
		rule, err := registry.GetRule("test-rule")
		if err != nil {
			t.Errorf("Failed to get rule: %v", err)
		}
		if rule.Metadata.Name != "Test Rule" {
			t.Errorf("Expected rule name to be 'Test Rule', but it was %s", rule.Metadata.Name)
		}

		// Try to get a non-existent rule
		_, err = registry.GetRule("nonexistent")
		if err == nil {
			t.Errorf("Expected error when getting non-existent rule, but got nil")
		}
	})

	// Test listing rules
	t.Run("ListRules", func(t *testing.T) {
		// Register another rule
		anotherRule := &CustomRule{
			Metadata: CustomRuleMetadata{
				Name:        "Another Rule",
				Description: "Another test rule",
				Author:      "Test Author",
				Version:     "1.0.0",
			},
			Func: func(data interface{}, params map[string]interface{}) (bool, string, error) {
				return true, "Always passes", nil
			},
		}
		err := registry.RegisterRule("another-rule", anotherRule)
		if err != nil {
			t.Errorf("Failed to register rule: %v", err)
		}

		// List all rules
		rules := registry.ListRules()
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules, but got %d", len(rules))
		}
	})

	// Test executing a rule
	t.Run("ExecuteRule", func(t *testing.T) {
		// Execute the rule with a value above the threshold
		passed, message, err := registry.ExecuteRule("test-rule", 0.7, map[string]interface{}{
			"threshold": 0.5,
		})
		if err != nil {
			t.Errorf("Failed to execute rule: %v", err)
		}
		if !passed {
			t.Errorf("Expected rule to pass, but it failed with message: %s", message)
		}

		// Execute the rule with a value below the threshold
		passed, message, err = registry.ExecuteRule("test-rule", 0.3, map[string]interface{}{
			"threshold": 0.5,
		})
		if err != nil {
			t.Errorf("Failed to execute rule: %v", err)
		}
		if passed {
			t.Errorf("Expected rule to fail, but it passed with message: %s", message)
		}

		// Try to execute a non-existent rule
		_, _, err = registry.ExecuteRule("nonexistent", 0.5, nil)
		if err == nil {
			t.Errorf("Expected error when executing non-existent rule, but got nil")
		}
	})

	// Test unregistering a rule
	t.Run("UnregisterRule", func(t *testing.T) {
		err := registry.UnregisterRule("test-rule")
		if err != nil {
			t.Errorf("Failed to unregister rule: %v", err)
		}

		// Try to get the unregistered rule
		_, err = registry.GetRule("test-rule")
		if err == nil {
			t.Errorf("Expected error when getting unregistered rule, but got nil")
		}

		// Try to unregister a non-existent rule
		err = registry.UnregisterRule("nonexistent")
		if err == nil {
			t.Errorf("Expected error when unregistering non-existent rule, but got nil")
		}
	})
}

func TestRuleLibraryManager(t *testing.T) {
	// Create a new custom rule registry
	registry := NewCustomRuleRegistry()

	// Create a new rule library manager
	manager := NewRuleLibraryManager(registry)

	// Create a test rule
	testRule := &CustomRule{
		Metadata: CustomRuleMetadata{
			Name:        "Test Rule",
			Description: "A test rule for unit testing",
			Author:      "Test Author",
			Version:     "1.0.0",
		},
		Func: func(data interface{}, params map[string]interface{}) (bool, string, error) {
			return true, "Always passes", nil
		},
	}

	// Create a test library
	testLibrary := &RuleLibrary{
		Name:        "Test Library",
		Description: "A test rule library",
		Author:      "Test Author",
		Version:     "1.0.0",
		Rules: map[string]*CustomRule{
			"test-rule": testRule,
		},
		Tags: []string{"test", "validation"},
		Metadata: map[string]interface{}{
			"created_at": "2023-01-01",
		},
	}

	// Test registering a library
	t.Run("RegisterLibrary", func(t *testing.T) {
		err := manager.RegisterLibrary(testLibrary)
		if err != nil {
			t.Errorf("Failed to register library: %v", err)
		}

		// Try to register the same library again
		err = manager.RegisterLibrary(testLibrary)
		if err == nil {
			t.Errorf("Expected error when registering duplicate library, but got nil")
		}

		// Verify that the rule was registered
		rule, err := registry.GetRule("test-rule")
		if err != nil {
			t.Errorf("Failed to get rule: %v", err)
		}
		if rule.Metadata.Name != "Test Rule" {
			t.Errorf("Expected rule name to be 'Test Rule', but it was %s", rule.Metadata.Name)
		}
	})

	// Test getting a library
	t.Run("GetLibrary", func(t *testing.T) {
		library, err := manager.GetLibrary("Test Library")
		if err != nil {
			t.Errorf("Failed to get library: %v", err)
		}
		if library.Name != "Test Library" {
			t.Errorf("Expected library name to be 'Test Library', but it was %s", library.Name)
		}

		// Try to get a non-existent library
		_, err = manager.GetLibrary("Nonexistent Library")
		if err == nil {
			t.Errorf("Expected error when getting non-existent library, but got nil")
		}
	})

	// Test listing libraries
	t.Run("ListLibraries", func(t *testing.T) {
		// Register another library
		anotherRule := &CustomRule{
			Metadata: CustomRuleMetadata{
				Name:        "Another Rule",
				Description: "Another test rule",
				Author:      "Test Author",
				Version:     "1.0.0",
			},
			Func: func(data interface{}, params map[string]interface{}) (bool, string, error) {
				return true, "Always passes", nil
			},
		}
		anotherLibrary := &RuleLibrary{
			Name:        "Another Library",
			Description: "Another test rule library",
			Author:      "Test Author",
			Version:     "1.0.0",
			Rules: map[string]*CustomRule{
				"another-rule": anotherRule,
			},
			Tags: []string{"test", "quality"},
		}
		err := manager.RegisterLibrary(anotherLibrary)
		if err != nil {
			t.Errorf("Failed to register library: %v", err)
		}

		// List all libraries
		libraries := manager.ListLibraries()
		if len(libraries) != 2 {
			t.Errorf("Expected 2 libraries, but got %d", len(libraries))
		}
	})

	// Test finding libraries by tag
	t.Run("FindLibrariesByTag", func(t *testing.T) {
		// Find libraries by tag
		libraries := manager.FindLibrariesByTag("test")
		if len(libraries) != 2 {
			t.Errorf("Expected 2 libraries with tag 'test', but got %d", len(libraries))
		}

		libraries = manager.FindLibrariesByTag("validation")
		if len(libraries) != 1 {
			t.Errorf("Expected 1 library with tag 'validation', but got %d", len(libraries))
		}
		if libraries[0].Name != "Test Library" {
			t.Errorf("Expected library name to be 'Test Library', but it was %s", libraries[0].Name)
		}

		libraries = manager.FindLibrariesByTag("quality")
		if len(libraries) != 1 {
			t.Errorf("Expected 1 library with tag 'quality', but got %d", len(libraries))
		}
		if libraries[0].Name != "Another Library" {
			t.Errorf("Expected library name to be 'Another Library', but it was %s", libraries[0].Name)
		}

		libraries = manager.FindLibrariesByTag("nonexistent")
		if len(libraries) != 0 {
			t.Errorf("Expected 0 libraries with tag 'nonexistent', but got %d", len(libraries))
		}
	})

	// Test unregistering a library
	t.Run("UnregisterLibrary", func(t *testing.T) {
		err := manager.UnregisterLibrary("Test Library")
		if err != nil {
			t.Errorf("Failed to unregister library: %v", err)
		}

		// Try to get the unregistered library
		_, err = manager.GetLibrary("Test Library")
		if err == nil {
			t.Errorf("Expected error when getting unregistered library, but got nil")
		}

		// Verify that the rule was unregistered
		_, err = registry.GetRule("test-rule")
		if err == nil {
			t.Errorf("Expected error when getting unregistered rule, but got nil")
		}

		// Try to unregister a non-existent library
		err = manager.UnregisterLibrary("Nonexistent Library")
		if err == nil {
			t.Errorf("Expected error when unregistering non-existent library, but got nil")
		}
	})
}
