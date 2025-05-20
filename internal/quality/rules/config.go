package rules

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// RuleConfig represents the configuration for a rule
type RuleConfig struct {
	Name         string   `yaml:"name"`
	Description  string   `yaml:"description"`
	Type         string   `yaml:"type"`
	Column       string   `yaml:"column"`
	Condition    string   `yaml:"condition,omitempty"`
	Threshold    float64  `yaml:"threshold,omitempty"`
	Severity     string   `yaml:"severity"`
	Pattern      string   `yaml:"pattern,omitempty"`
	Min          float64  `yaml:"min,omitempty"`
	Max          float64  `yaml:"max,omitempty"`
	MinLength    int      `yaml:"min_length,omitempty"`
	MaxLength    int      `yaml:"max_length,omitempty"`
	EnumValues   []string `yaml:"enum_values,omitempty"`
	RegexPattern string   `yaml:"regex_pattern,omitempty"`
	DateFormat   string   `yaml:"date_format,omitempty"`
}

// RulesConfigFile represents a YAML file containing rule configurations
type RulesConfigFile struct {
	Rules []RuleConfig `yaml:"rules"`
}

// LoadRulesFromYAML loads rule configurations from a YAML file
func LoadRulesFromYAML(filePath string) ([]ExtendedRule, error) {
	// Read the YAML file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules config file: %v", err)
	}

	// Parse the YAML data
	var config RulesConfigFile
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse rules config file: %v", err)
	}

	// Convert RuleConfig to ExtendedRule
	rules := make([]ExtendedRule, 0, len(config.Rules))
	for _, ruleConfig := range config.Rules {
		rule := ExtendedRule{
			Rule: Rule{
				Name:        ruleConfig.Name,
				Description: ruleConfig.Description,
				Type:        RuleType(ruleConfig.Type),
				Column:      ruleConfig.Column,
				Condition:   ruleConfig.Condition,
				Threshold:   ruleConfig.Threshold,
				Severity:    Severity(ruleConfig.Severity),
				Pattern:     ruleConfig.Pattern,
				Min:         ruleConfig.Min,
				Max:         ruleConfig.Max,
			},
			MinLength:    ruleConfig.MinLength,
			MaxLength:    ruleConfig.MaxLength,
			EnumValues:   ruleConfig.EnumValues,
			RegexPattern: ruleConfig.RegexPattern,
			DateFormat:   ruleConfig.DateFormat,
		}

		// Validate the rule
		if err := rule.Validate(); err != nil {
			return nil, fmt.Errorf("invalid rule '%s': %v", rule.Name, err)
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// SaveRulesToYAML saves rule configurations to a YAML file
func SaveRulesToYAML(rules []ExtendedRule, filePath string) error {
	// Convert ExtendedRule to RuleConfig
	configs := make([]RuleConfig, 0, len(rules))
	for _, rule := range rules {
		config := RuleConfig{
			Name:         rule.Name,
			Description:  rule.Description,
			Type:         string(rule.Type),
			Column:       rule.Column,
			Condition:    rule.Condition,
			Threshold:    rule.Threshold,
			Severity:     string(rule.Severity),
			Pattern:      rule.Pattern,
			Min:          rule.Min,
			Max:          rule.Max,
			MinLength:    rule.MinLength,
			MaxLength:    rule.MaxLength,
			EnumValues:   rule.EnumValues,
			RegexPattern: rule.RegexPattern,
			DateFormat:   rule.DateFormat,
		}
		configs = append(configs, config)
	}

	configFile := RulesConfigFile{
		Rules: configs,
	}

	// Marshal to YAML
	data, err := yaml.Marshal(configFile)
	if err != nil {
		return fmt.Errorf("failed to marshal rules config: %v", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write rules config file: %v", err)
	}

	return nil
}

// LoadDefaultRules loads a set of default rules
func LoadDefaultRules() []ExtendedRule {
	return []ExtendedRule{
		{
			Rule: Rule{
				Name:        "Not Null Check",
				Description: "Checks that values are not null",
				Type:        RuleTypeCompleteness,
				Column:      "id",
				Threshold:   100.0,
				Severity:    SeverityCritical,
			},
		},
		{
			Rule: Rule{
				Name:        "Unique ID Check",
				Description: "Checks that IDs are unique",
				Type:        RuleTypeUniqueness,
				Column:      "id",
				Threshold:   100.0,
				Severity:    SeverityHigh,
			},
		},
		{
			Rule: Rule{
				Name:        "Age Range Check",
				Description: "Checks that age values are within a valid range",
				Type:        RuleTypeRange,
				Column:      "age",
				Min:         0,
				Max:         120,
				Severity:    SeverityMedium,
			},
		},
		{
			Rule: Rule{
				Name:        "Email Pattern Check",
				Description: "Checks that email values contain @ symbol",
				Type:        RuleTypePattern,
				Column:      "email",
				Pattern:     "@",
				Severity:    SeverityMedium,
			},
		},
		{
			Rule: Rule{
				Name:        "Email Format Check",
				Description: "Checks that email values match a valid email format",
				Type:        RuleTypeRegex,
				Column:      "email",
				Severity:    SeverityHigh,
			},
			RegexPattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		},
		{
			Rule: Rule{
				Name:        "Username Length Check",
				Description: "Checks that username length is between 3 and 20 characters",
				Type:        RuleTypeLength,
				Column:      "username",
				Severity:    SeverityMedium,
			},
			MinLength: 3,
			MaxLength: 20,
		},
		{
			Rule: Rule{
				Name:        "Status Enum Check",
				Description: "Checks that status values are valid",
				Type:        RuleTypeEnum,
				Column:      "status",
				Severity:    SeverityMedium,
			},
			EnumValues: []string{"active", "inactive", "pending"},
		},
		{
			Rule: Rule{
				Name:        "Date Format Check",
				Description: "Checks that date values are in the correct format",
				Type:        RuleTypeDateTime,
				Column:      "created_at",
				Severity:    SeverityMedium,
			},
			DateFormat: "2006-01-02",
		},
	}
}
