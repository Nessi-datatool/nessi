package rules

import (
	"fmt"
	"sync"
)

// CustomRuleFunc is a function that implements a custom validation rule
type CustomRuleFunc func(data interface{}, params map[string]interface{}) (bool, string, error)

// CustomRuleMetadata contains metadata about a custom rule
type CustomRuleMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Author      string                 `json:"author"`
	Version     string                 `json:"version"`
	Parameters  map[string]interface{} `json:"parameters"`
	Examples    []string               `json:"examples"`
}

// CustomRule represents a custom validation rule
type CustomRule struct {
	Metadata CustomRuleMetadata `json:"metadata"`
	Func     CustomRuleFunc     `json:"-"`
}

// CustomRuleRegistry manages custom validation rules
type CustomRuleRegistry struct {
	rules map[string]*CustomRule
	mu    sync.RWMutex
}

// NewCustomRuleRegistry creates a new custom rule registry
func NewCustomRuleRegistry() *CustomRuleRegistry {
	return &CustomRuleRegistry{
		rules: make(map[string]*CustomRule),
	}
}

// RegisterRule registers a custom rule
func (r *CustomRuleRegistry) RegisterRule(name string, rule *CustomRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[name]; exists {
		return fmt.Errorf("rule %s is already registered", name)
	}

	r.rules[name] = rule
	return nil
}

// UnregisterRule unregisters a custom rule
func (r *CustomRuleRegistry) UnregisterRule(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.rules[name]; !exists {
		return fmt.Errorf("rule %s is not registered", name)
	}

	delete(r.rules, name)
	return nil
}

// GetRule gets a custom rule by name
func (r *CustomRuleRegistry) GetRule(name string) (*CustomRule, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rule, exists := r.rules[name]
	if !exists {
		return nil, fmt.Errorf("rule %s is not registered", name)
	}

	return rule, nil
}

// ListRules lists all registered custom rules
func (r *CustomRuleRegistry) ListRules() []*CustomRule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rules := make([]*CustomRule, 0, len(r.rules))
	for _, rule := range r.rules {
		rules = append(rules, rule)
	}

	return rules
}

// ExecuteRule executes a custom rule
func (r *CustomRuleRegistry) ExecuteRule(name string, data interface{}, params map[string]interface{}) (bool, string, error) {
	rule, err := r.GetRule(name)
	if err != nil {
		return false, "", err
	}

	return rule.Func(data, params)
}

// RuleLibrary represents a shared library of validation rules
type RuleLibrary struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Author      string                  `json:"author"`
	Version     string                  `json:"version"`
	Rules       map[string]*CustomRule  `json:"rules"`
	Tags        []string                `json:"tags"`
	Metadata    map[string]interface{}  `json:"metadata"`
}

// RuleLibraryManager manages shared rule libraries
type RuleLibraryManager struct {
	libraries map[string]*RuleLibrary
	registry  *CustomRuleRegistry
	mu        sync.RWMutex
}

// NewRuleLibraryManager creates a new rule library manager
func NewRuleLibraryManager(registry *CustomRuleRegistry) *RuleLibraryManager {
	return &RuleLibraryManager{
		libraries: make(map[string]*RuleLibrary),
		registry:  registry,
	}
}

// RegisterLibrary registers a rule library
func (m *RuleLibraryManager) RegisterLibrary(library *RuleLibrary) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.libraries[library.Name]; exists {
		return fmt.Errorf("library %s is already registered", library.Name)
	}

	// Register all rules in the library
	for name, rule := range library.Rules {
		if err := m.registry.RegisterRule(name, rule); err != nil {
			// Rollback registrations on error
			for n := range library.Rules {
				if n == name {
					break
				}
				m.registry.UnregisterRule(n)
			}
			return fmt.Errorf("failed to register rule %s: %v", name, err)
		}
	}

	m.libraries[library.Name] = library
	return nil
}

// UnregisterLibrary unregisters a rule library
func (m *RuleLibraryManager) UnregisterLibrary(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	library, exists := m.libraries[name]
	if !exists {
		return fmt.Errorf("library %s is not registered", name)
	}

	// Unregister all rules in the library
	for name := range library.Rules {
		if err := m.registry.UnregisterRule(name); err != nil {
			return fmt.Errorf("failed to unregister rule %s: %v", name, err)
		}
	}

	delete(m.libraries, name)
	return nil
}

// GetLibrary gets a rule library by name
func (m *RuleLibraryManager) GetLibrary(name string) (*RuleLibrary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	library, exists := m.libraries[name]
	if !exists {
		return nil, fmt.Errorf("library %s is not registered", name)
	}

	return library, nil
}

// ListLibraries lists all registered rule libraries
func (m *RuleLibraryManager) ListLibraries() []*RuleLibrary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	libraries := make([]*RuleLibrary, 0, len(m.libraries))
	for _, library := range m.libraries {
		libraries = append(libraries, library)
	}

	return libraries
}

// FindLibrariesByTag finds rule libraries by tag
func (m *RuleLibraryManager) FindLibrariesByTag(tag string) []*RuleLibrary {
	m.mu.RLock()
	defer m.mu.RUnlock()

	libraries := make([]*RuleLibrary, 0)
	for _, library := range m.libraries {
		for _, t := range library.Tags {
			if t == tag {
				libraries = append(libraries, library)
				break
			}
		}
	}

	return libraries
}
