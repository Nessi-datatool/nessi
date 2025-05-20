package rules

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// YAMLRuleConfig represents the structure of a rule configuration in YAML
type YAMLRuleConfig struct {
	ID          string                 `yaml:"id"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Severity    string                 `yaml:"severity"`
	RuleType    string                 `yaml:"rule_type"`
	Field       string                 `yaml:"field,omitempty"`
	Tags        []string               `yaml:"tags"`
	Config      map[string]interface{} `yaml:"config"`
}

// YAMLRulesConfig represents the structure of a rules configuration file
type YAMLRulesConfig struct {
	Rules []YAMLRuleConfig `yaml:"rules"`
}

// RuleLoader loads rules from various sources
type RuleLoader struct {
	validator *RuleValidator
}

// NewRuleLoader creates a new rule loader
func NewRuleLoader() *RuleLoader {
	return &RuleLoader{
		validator: NewRuleValidator(),
	}
}

// LoadFromYAML loads rules from a YAML file
func (l *RuleLoader) LoadFromYAML(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	var config YAMLRulesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	for _, ruleConfig := range config.Rules {
		rule, err := l.createRuleFromConfig(ruleConfig)
		if err != nil {
			return fmt.Errorf("failed to create rule '%s': %w", ruleConfig.ID, err)
		}

		l.validator.AddRule(rule)
	}

	return nil
}

// GetValidator returns the rule validator with loaded rules
func (l *RuleLoader) GetValidator() *RuleValidator {
	return l.validator
}

// createRuleFromConfig creates a rule from a YAML configuration
func (l *RuleLoader) createRuleFromConfig(config YAMLRuleConfig) (Rule, error) {
	metadata := RuleMetadata{
		ID:          config.ID,
		Name:        config.Name,
		Description: config.Description,
		Severity:    config.Severity,
		Tags:        config.Tags,
		Config:      config.Config,
	}

	switch config.RuleType {
	case "null_check":
		fields, ok := config.Config["fields"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid fields configuration for null check rule")
		}

		fieldNames := make([]string, len(fields))
		for i, field := range fields {
			fieldStr, ok := field.(string)
			if !ok {
				return nil, fmt.Errorf("field name must be a string")
			}
			fieldNames[i] = fieldStr
		}

		return NewNullCheckRule(fieldNames, metadata), nil

	case "range":
		fieldsConfig, ok := config.Config["fields"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid fields configuration for range check rule")
		}

		rangeFields := make(map[string]RangeConfig)
		for field, rangeVal := range fieldsConfig {
			rangeMap, ok := rangeVal.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid range configuration for field '%s'", field)
			}

			min, ok := rangeMap["min"].(float64)
			if !ok {
				return nil, fmt.Errorf("min value must be a number for field '%s'", field)
			}

			max, ok := rangeMap["max"].(float64)
			if !ok {
				return nil, fmt.Errorf("max value must be a number for field '%s'", field)
			}

			rangeFields[field] = RangeConfig{
				Min: min,
				Max: max,
			}
		}

		return NewRangeCheckRule(rangeFields, metadata), nil

	case "regex":
		pattern, ok := config.Config["pattern"].(string)
		if !ok {
			return nil, fmt.Errorf("pattern must be a string for regex rule")
		}

		return NewRegexRule(config.Field, pattern, metadata)

	case "enum":
		values, ok := config.Config["values"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("values must be an array for enum rule")
		}

		strValues := make([]string, len(values))
		for i, v := range values {
			strVal, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("enum values must be strings")
			}
			strValues[i] = strVal
		}

		caseSensitive := true
		if cs, ok := config.Config["case_sensitive"].(bool); ok {
			caseSensitive = cs
		}

		return NewEnumRule(config.Field, strValues, caseSensitive, metadata), nil

	case "length":
		// Handle both int and float64 values for min_length
		var minLength int
		if minLengthFloat, ok := config.Config["min_length"].(float64); ok {
			minLength = int(minLengthFloat)
		} else if minLengthInt, ok := config.Config["min_length"].(int); ok {
			minLength = minLengthInt
		} else {
			return nil, fmt.Errorf("min_length must be a number for length rule")
		}

		// Handle both int and float64 values for max_length
		var maxLength int
		if maxLengthFloat, ok := config.Config["max_length"].(float64); ok {
			maxLength = int(maxLengthFloat)
		} else if maxLengthInt, ok := config.Config["max_length"].(int); ok {
			maxLength = maxLengthInt
		} else {
			return nil, fmt.Errorf("max_length must be a number for length rule")
		}

		return NewLengthRule(config.Field, minLength, maxLength, metadata), nil

	case "date_format":
		format, ok := config.Config["format"].(string)
		if !ok {
			return nil, fmt.Errorf("format must be a string for date format rule")
		}

		return NewDateFormatRule(config.Field, format, metadata), nil

	default:
		return nil, fmt.Errorf("unsupported rule type: %s", config.RuleType)
	}
}

// SaveRulesToYAML saves rules to a YAML file
func SaveRulesToYAML(rules []Rule, filePath string) error {
	config := YAMLRulesConfig{
		Rules: make([]YAMLRuleConfig, len(rules)),
	}

	for i, rule := range rules {
		metadata := rule.GetMetadata()

		yamlConfig := YAMLRuleConfig{
			ID:          metadata.ID,
			Name:        metadata.Name,
			Description: metadata.Description,
			Severity:    metadata.Severity,
			Tags:        metadata.Tags,
			Config:      metadata.Config,
		}

		// Determine rule type and field
		switch rule.(type) {
		case *NullCheckRule:
			yamlConfig.RuleType = "null_check"
		case *RangeCheckRule:
			yamlConfig.RuleType = "range"
		case *RegexRule:
			yamlConfig.RuleType = "regex"
			if r, ok := rule.(*RegexRule); ok {
				yamlConfig.Field = r.field
			}
		case *EnumRule:
			yamlConfig.RuleType = "enum"
			if r, ok := rule.(*EnumRule); ok {
				yamlConfig.Field = r.field
			}
		case *LengthRule:
			yamlConfig.RuleType = "length"
			if r, ok := rule.(*LengthRule); ok {
				yamlConfig.Field = r.field
			}
		case *DateFormatRule:
			yamlConfig.RuleType = "date_format"
			if r, ok := rule.(*DateFormatRule); ok {
				yamlConfig.Field = r.field
			}
		}

		config.Rules[i] = yamlConfig
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal rules to YAML: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}
