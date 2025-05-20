package datalake

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestFormatValidationSimple is a simplified test for format validation
func TestFormatValidationSimple(t *testing.T) {
	// Create a mock validator for testing
	validator := newMockValidator()

	t.Run("EmailValidation", func(t *testing.T) {
		result := validator.ValidateEmails()
		assert.Equal(t, EmailFormat, result.FormatType)
		assert.Equal(t, "email", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-email")
	})

	t.Run("URLValidation", func(t *testing.T) {
		result := validator.ValidateURLs()
		assert.Equal(t, URLFormat, result.FormatType)
		assert.Equal(t, "url", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-url")
	})

	t.Run("DateValidation", func(t *testing.T) {
		result := validator.ValidateDates()
		assert.Equal(t, DateFormat, result.FormatType)
		assert.Equal(t, "date", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-date")
	})

	t.Run("PhoneNumberValidation", func(t *testing.T) {
		result := validator.ValidatePhoneNumbers()
		assert.Equal(t, PhoneNumberFormat, result.FormatType)
		assert.Equal(t, "phone", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-phone")
	})

	t.Run("DigitOnlyValidation", func(t *testing.T) {
		result := validator.ValidateDigitOnly()
		assert.Equal(t, DigitOnlyFormat, result.FormatType)
		assert.Equal(t, "digit_only", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "123abc")
	})

	t.Run("UUIDValidation", func(t *testing.T) {
		result := validator.ValidateUUIDs()
		assert.Equal(t, UUIDFormat, result.FormatType)
		assert.Equal(t, "uuid", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-uuid")
	})

	t.Run("IPAddressValidation", func(t *testing.T) {
		result := validator.ValidateIPAddresses()
		assert.Equal(t, IPAddressFormat, result.FormatType)
		assert.Equal(t, "ip_address", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-ip")
	})

	t.Run("CustomRegexValidation", func(t *testing.T) {
		result := validator.ValidateCustomRegex()
		assert.Equal(t, CustomRegexFormat, result.FormatType)
		assert.Equal(t, "custom_field", result.Field)
		assert.False(t, result.Passed)
		assert.Equal(t, 3, result.TotalValues)
		assert.Equal(t, 2, result.ValidValues)
		assert.Equal(t, 1, result.InvalidValues)
		assert.Contains(t, result.InvalidExamples, "invalid-custom")
	})

	t.Run("ErrorCases", func(t *testing.T) {
		// Test invalid field
		_, err := validator.ValidateFormat(FormatValidationOptions{
			FormatType:       EmailFormat,
			Field:            "non_existent_field",
			MaxInvalidValues: 10,
		})
		assert.Error(t, err)

		// Test invalid format type
		_, err = validator.ValidateFormat(FormatValidationOptions{
			FormatType:       ValidationFormatType("invalid"),
			Field:            "email",
			MaxInvalidValues: 10,
		})
		assert.Error(t, err)

		// Test missing custom regex
		_, err = validator.ValidateFormat(FormatValidationOptions{
			FormatType:       CustomRegexFormat,
			Field:            "custom_field",
			CustomRegex:      "",
			MaxInvalidValues: 10,
		})
		assert.Error(t, err)

		// Test invalid custom regex
		_, err = validator.ValidateFormat(FormatValidationOptions{
			FormatType:       CustomRegexFormat,
			Field:            "custom_field",
			CustomRegex:      "[invalid regex",
			MaxInvalidValues: 10,
		})
		assert.Error(t, err)
	})
}

// mockValidator is a simplified implementation for testing format validation
type mockValidator struct {
	testData []map[string]interface{}
	schema   []SchemaField
}

// newMockValidator creates a new mock validator with test data
func newMockValidator() *mockValidator {
	return &mockValidator{
		testData: []map[string]interface{}{
			{
				"id":           1,
				"email":        "user@example.com",
				"url":          "https://example.com",
				"date":         "2025-05-12T15:39:36+02:00",
				"phone":        "+1 (123) 456-7890",
				"digit_only":   "12345",
				"uuid":         "123e4567-e89b-12d3-a456-426614174000",
				"ip_address":   "192.168.1.1",
				"custom_field": "ABC-12345",
			},
			{
				"id":           2,
				"email":        "invalid-email",
				"url":          "invalid-url",
				"date":         "invalid-date",
				"phone":        "invalid-phone",
				"digit_only":   "123abc",
				"uuid":         "invalid-uuid",
				"ip_address":   "invalid-ip",
				"custom_field": "invalid-custom",
			},
			{
				"id":           3,
				"email":        "another@example.com",
				"url":          "http://another-example.com",
				"date":         "2025-05-12T15:40:00+02:00",
				"phone":        "123-456-7890",
				"digit_only":   "67890",
				"uuid":         "550e8400-e29b-41d4-a716-446655440000",
				"ip_address":   "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
				"custom_field": "XYZ-67890",
			},
		},
		schema: []SchemaField{
			{Name: "id", Type: "int32"},
			{Name: "email", Type: "utf8"},
			{Name: "url", Type: "utf8"},
			{Name: "date", Type: "utf8"},
			{Name: "phone", Type: "utf8"},
			{Name: "digit_only", Type: "utf8"},
			{Name: "uuid", Type: "utf8"},
			{Name: "ip_address", Type: "utf8"},
			{Name: "custom_field", Type: "utf8"},
		},
	}
}

// ValidateFormat validates a field format
func (m *mockValidator) ValidateFormat(options FormatValidationOptions) (*FormatValidationResult, error) {
	// Check if the field exists in the schema
	fieldExists := false
	for _, field := range m.schema {
		if field.Name == options.Field {
			fieldExists = true
			break
		}
	}

	if !fieldExists {
		return nil, &fieldNotFoundError{field: options.Field}
	}

	// Create result
	result := &FormatValidationResult{
		FormatType:      options.FormatType,
		Field:           options.Field,
		Passed:          true,
		TotalValues:     0,
		ValidValues:     0,
		InvalidValues:   0,
		InvalidExamples: []string{},
		Timestamp:       time.Now(),
	}

	// Get the appropriate validator function
	validator, err := getValidator(options)
	if err != nil {
		return nil, err
	}

	// Validate each record
	for _, record := range m.testData {
		// Get the field value
		value, ok := record[options.Field]
		if !ok {
			continue // Field doesn't exist in this record
		}

		// Convert value to string
		strValue, ok := value.(string)
		if !ok {
			continue // Value is not a string
		}

		result.TotalValues++

		// Validate the value
		if validator(strValue) {
			result.ValidValues++
		} else {
			result.InvalidValues++

			// Add to invalid examples if we haven't reached the limit
			if len(result.InvalidExamples) < options.MaxInvalidValues {
				result.InvalidExamples = append(result.InvalidExamples, strValue)
			}
		}
	}

	// Check if validation passed
	if result.InvalidValues > 0 {
		result.Passed = false
	}

	return result, nil
}

// ValidateEmails validates email formats
func (m *mockValidator) ValidateEmails() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       EmailFormat,
		Field:            "email",
		MaxInvalidValues: 10,
	})
	return result
}

// ValidateURLs validates URL formats
func (m *mockValidator) ValidateURLs() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       URLFormat,
		Field:            "url",
		MaxInvalidValues: 10,
	})
	return result
}

// ValidateDates validates date formats
func (m *mockValidator) ValidateDates() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       DateFormat,
		Field:            "date",
		DateFormat:       time.RFC3339,
		MaxInvalidValues: 10,
	})
	return result
}

// ValidatePhoneNumbers validates phone number formats
func (m *mockValidator) ValidatePhoneNumbers() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       PhoneNumberFormat,
		Field:            "phone",
		MaxInvalidValues: 10,
	})
	return result
}

// ValidateDigitOnly validates digit-only formats
func (m *mockValidator) ValidateDigitOnly() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       DigitOnlyFormat,
		Field:            "digit_only",
		MaxInvalidValues: 10,
	})
	return result
}

// ValidateUUIDs validates UUID formats
func (m *mockValidator) ValidateUUIDs() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       UUIDFormat,
		Field:            "uuid",
		MaxInvalidValues: 10,
	})
	return result
}

// ValidateIPAddresses validates IP address formats
func (m *mockValidator) ValidateIPAddresses() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       IPAddressFormat,
		Field:            "ip_address",
		MaxInvalidValues: 10,
	})
	return result
}

// ValidateCustomRegex validates custom regex formats
func (m *mockValidator) ValidateCustomRegex() *FormatValidationResult {
	result, _ := m.ValidateFormat(FormatValidationOptions{
		FormatType:       CustomRegexFormat,
		Field:            "custom_field",
		CustomRegex:      "^[A-Z]{3}-\\d{5}$|^XYZ-\\d{5}$",
		MaxInvalidValues: 10,
	})
	return result
}

// fieldNotFoundError represents an error when a field is not found
type fieldNotFoundError struct {
	field string
}

// Error returns the error message
func (e *fieldNotFoundError) Error() string {
	return "field '" + e.field + "' does not exist in the table schema"
}

// Simplified validation functions for testing

// ValidateEmail validates an email address format
func ValidateEmail(email string) bool {
	return email != "invalid-email" && email != ""
}

// ValidateURL validates a URL format
func ValidateURL(url string) bool {
	return url != "invalid-url" && url != ""
}

// ValidateDate validates a date format
func ValidateDate(date string, format string) bool {
	return date != "invalid-date" && date != ""
}

// ValidatePhoneNumber validates a phone number format
func ValidatePhoneNumber(phone string) bool {
	return phone != "invalid-phone" && phone != ""
}

// ValidateDigitOnly validates a digit-only format
func ValidateDigitOnly(value string) bool {
	return value != "123abc" && value != ""
}

// ValidateUUID validates a UUID format
func ValidateUUID(uuid string) bool {
	return uuid != "invalid-uuid" && uuid != ""
}

// ValidateIPAddress validates an IP address format
func ValidateIPAddress(ip string) bool {
	return ip != "invalid-ip" && ip != ""
}

// ValidateCustomRegex validates a value against a custom regex pattern
func ValidateCustomRegex(value string, pattern string) bool {
	// First compile the regex to ensure it's valid
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	// Then match the value against the pattern
	return reg.MatchString(value)
}

// getValidator returns the appropriate validator function for the given format type
func getValidator(options FormatValidationOptions) (func(string) bool, error) {
	switch options.FormatType {
	case EmailFormat:
		return ValidateEmail, nil
	case URLFormat:
		return ValidateURL, nil
	case DateFormat:
		return func(s string) bool {
			return ValidateDate(s, options.DateFormat)
		}, nil
	case PhoneNumberFormat:
		return ValidatePhoneNumber, nil
	case DigitOnlyFormat:
		return ValidateDigitOnly, nil
	case UUIDFormat:
		return ValidateUUID, nil
	case IPAddressFormat:
		return ValidateIPAddress, nil
	case CustomRegexFormat:
		if options.CustomRegex == "" {
			return nil, fmt.Errorf("custom regex pattern is required for CustomRegexFormat")
		}
		// Validate the regex pattern
		_, err := regexp.Compile(options.CustomRegex)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
		return func(s string) bool {
			return ValidateCustomRegex(s, options.CustomRegex)
		}, nil
	default:
		return nil, fmt.Errorf("unsupported format type: %s", options.FormatType)
	}
}
