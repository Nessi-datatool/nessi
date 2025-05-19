package main

import (
	"fmt"
	"time"

	"github.com/nessi-dev/nessi/pkg/datalake"
)

// mockValidator is a simplified implementation for testing format validation
type mockValidator struct {
	testData []map[string]interface{}
	schema   []datalake.SchemaField
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
		schema: []datalake.SchemaField{
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
func (m *mockValidator) ValidateFormat(options datalake.FormatValidationOptions) (*datalake.FormatValidationResult, error) {
	// Check if the field exists in the schema
	fieldExists := false
	for _, field := range m.schema {
		if field.Name == options.Field {
			fieldExists = true
			break
		}
	}

	if !fieldExists {
		return nil, fmt.Errorf("field '%s' does not exist in the table schema", options.Field)
	}

	// Create result
	result := &datalake.FormatValidationResult{
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
	validator, err := datalake.GetValidator(options)
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

func main() {
	fmt.Println("Format Validation Demo")
	fmt.Println("======================")
	fmt.Println()

	// Create a mock validator
	validator := newMockValidator()

	// Run email format validation
	fmt.Println("Running email format validation...")
	emailResult, err := validator.ValidateFormat(datalake.FormatValidationOptions{
		FormatType:       datalake.EmailFormat,
		Field:            "email",
		MaxInvalidValues: 10,
	})
	if err != nil {
		fmt.Printf("Error running email format validation: %v\n", err)
		return
	}
	printFormatValidationResult(emailResult)
	fmt.Println()

	// Run URL format validation
	fmt.Println("Running URL format validation...")
	urlResult, err := validator.ValidateFormat(datalake.FormatValidationOptions{
		FormatType:       datalake.URLFormat,
		Field:            "url",
		MaxInvalidValues: 10,
	})
	if err != nil {
		fmt.Printf("Error running URL format validation: %v\n", err)
		return
	}
	printFormatValidationResult(urlResult)
	fmt.Println()

	// Run date format validation
	fmt.Println("Running date format validation...")
	dateResult, err := validator.ValidateFormat(datalake.FormatValidationOptions{
		FormatType:       datalake.DateFormat,
		Field:            "date",
		DateFormat:       time.RFC3339,
		MaxInvalidValues: 10,
	})
	if err != nil {
		fmt.Printf("Error running date format validation: %v\n", err)
		return
	}
	printFormatValidationResult(dateResult)
	fmt.Println()

	// Run phone number format validation
	fmt.Println("Running phone number format validation...")
	phoneResult, err := validator.ValidateFormat(datalake.FormatValidationOptions{
		FormatType:       datalake.PhoneNumberFormat,
		Field:            "phone",
		MaxInvalidValues: 10,
	})
	if err != nil {
		fmt.Printf("Error running phone number format validation: %v\n", err)
		return
	}
	printFormatValidationResult(phoneResult)
	fmt.Println()

	// Run UUID format validation
	fmt.Println("Running UUID format validation...")
	uuidResult, err := validator.ValidateFormat(datalake.FormatValidationOptions{
		FormatType:       datalake.UUIDFormat,
		Field:            "uuid",
		MaxInvalidValues: 10,
	})
	if err != nil {
		fmt.Printf("Error running UUID format validation: %v\n", err)
		return
	}
	printFormatValidationResult(uuidResult)
	fmt.Println()

	// Run custom regex format validation
	fmt.Println("Running custom regex format validation...")
	regexResult, err := validator.ValidateFormat(datalake.FormatValidationOptions{
		FormatType:       datalake.CustomRegexFormat,
		Field:            "custom_field",
		CustomRegex:      "^[A-Z]{3}-\\d{5}$",
		MaxInvalidValues: 10,
	})
	if err != nil {
		fmt.Printf("Error running custom regex format validation: %v\n", err)
		return
	}
	printFormatValidationResult(regexResult)
	fmt.Println()

	fmt.Println("Demo completed!")
}

// printFormatValidationResult prints the format validation result
func printFormatValidationResult(result *datalake.FormatValidationResult) {
	fmt.Printf("Format Validation Result: %s\n", result.FormatType)
	fmt.Printf("Field: %s\n", result.Field)
	fmt.Printf("Passed: %t\n", result.Passed)
	fmt.Printf("Total Values: %d\n", result.TotalValues)
	fmt.Printf("Valid Values: %d\n", result.ValidValues)
	fmt.Printf("Invalid Values: %d\n", result.InvalidValues)
	
	if len(result.InvalidExamples) > 0 {
		fmt.Printf("Invalid Examples (%d):\n", len(result.InvalidExamples))
		for i, example := range result.InvalidExamples {
			fmt.Printf("  %d. %s\n", i+1, example)
		}
	}
}
