package datalake

import (
	"fmt"
	"regexp"
	"time"
)

// ValidationFormatType defines the type of format validation to perform
type ValidationFormatType string

const (
	// DateFormat validates date formats
	DateFormat ValidationFormatType = "date"

	// EmailFormat validates email formats
	EmailFormat ValidationFormatType = "email"

	// URLFormat validates URL formats
	URLFormat ValidationFormatType = "url"

	// DigitOnlyFormat validates digit-only formats
	DigitOnlyFormat ValidationFormatType = "digit_only"

	// UUIDFormat validates UUID formats
	UUIDFormat ValidationFormatType = "uuid"

	// IPAddressFormat validates IP address formats
	IPAddressFormat ValidationFormatType = "ip_address"

	// PhoneNumberFormat validates phone number formats
	PhoneNumberFormat ValidationFormatType = "phone_number"

	// CustomRegexFormat validates custom regex formats
	CustomRegexFormat ValidationFormatType = "custom_regex"
)

// FormatValidationOptions defines the options for format validation
type FormatValidationOptions struct {
	// The format type to validate
	FormatType ValidationFormatType

	// The field to validate
	Field string

	// The custom regex pattern to use (only for CustomRegexFormat)
	CustomRegex string

	// The date format to use (only for DateFormat)
	DateFormat string

	// The maximum number of invalid values to return
	MaxInvalidValues int
}

// FormatValidationResult represents the result of a format validation
type FormatValidationResult struct {
	// The format type that was validated
	FormatType ValidationFormatType

	// The field that was validated
	Field string

	// Whether the validation passed
	Passed bool

	// The total number of values checked
	TotalValues int

	// The number of valid values
	ValidValues int

	// The number of invalid values
	InvalidValues int

	// Examples of invalid values (up to MaxInvalidValues)
	InvalidExamples []string

	// The timestamp when the validation was performed
	Timestamp time.Time
}

// ValidateFormat validates the format of a field in a Delta table
func (m *MetadataManager) ValidateFormat(options FormatValidationOptions) (*FormatValidationResult, error) {
	// Read the latest table metadata
	table, err := m.ReadTableMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to read table metadata: %w", err)
	}

	// Check if the field exists in the schema
	fieldExists := false
	// Get schema fields
	schemaFields, err := m.GetSchemaFieldsAtVersion(table.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema fields: %w", err)
	}

	for _, field := range schemaFields {
		if field.Name == options.Field {
			fieldExists = true
			break
		}
	}

	if !fieldExists {
		return nil, fmt.Errorf("field '%s' does not exist in the table schema", options.Field)
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
	validator, err := GetValidator(options)
	if err != nil {
		return nil, err
	}

	// Read all parquet files and validate the field values
	for _, file := range table.Files {
		// Read the parquet file
		records, err := m.ReadParquetFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read parquet file '%s': %w", file, err)
		}

		// Validate each record
		for _, record := range records {
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
	}

	// Check if validation passed
	if result.InvalidValues > 0 {
		result.Passed = false
	}

	return result, nil
}

// GetValidator returns a validator function for the specified format type
func GetValidator(options FormatValidationOptions) (func(string) bool, error) {
	switch options.FormatType {
	case DateFormat:
		return getDateValidator(options.DateFormat), nil
	case EmailFormat:
		return getEmailValidator(), nil
	case URLFormat:
		return getURLValidator(), nil
	case DigitOnlyFormat:
		return getDigitOnlyValidator(), nil
	case UUIDFormat:
		return getUUIDValidator(), nil
	case IPAddressFormat:
		return getIPAddressValidator(), nil
	case PhoneNumberFormat:
		return getPhoneNumberValidator(), nil
	case CustomRegexFormat:
		return getCustomRegexValidator(options.CustomRegex)
	default:
		return nil, fmt.Errorf("unsupported format type: %s", options.FormatType)
	}
}

// getDateValidator returns a validator function for date formats
func getDateValidator(format string) func(string) bool {
	if format == "" {
		// Default to RFC3339 format
		format = time.RFC3339
	}

	return func(value string) bool {
		_, err := time.Parse(format, value)
		return err == nil
	}
}

// getEmailValidator returns a validator function for email formats
func getEmailValidator() func(string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return func(value string) bool {
		return emailRegex.MatchString(value)
	}
}

// getURLValidator returns a validator function for URL formats
func getURLValidator() func(string) bool {
	urlRegex := regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
	return func(value string) bool {
		return urlRegex.MatchString(value)
	}
}

// getDigitOnlyValidator returns a validator function for digit-only formats
func getDigitOnlyValidator() func(string) bool {
	digitRegex := regexp.MustCompile(`^\d+$`)
	return func(value string) bool {
		return digitRegex.MatchString(value)
	}
}

// getUUIDValidator returns a validator function for UUID formats
func getUUIDValidator() func(string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return func(value string) bool {
		return uuidRegex.MatchString(value)
	}
}

// getIPAddressValidator returns a validator function for IP address formats
func getIPAddressValidator() func(string) bool {
	ipv4Regex := regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
	ipv6Regex := regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)

	return func(value string) bool {
		return ipv4Regex.MatchString(value) || ipv6Regex.MatchString(value)
	}
}

// getPhoneNumberValidator returns a validator function for phone number formats
func getPhoneNumberValidator() func(string) bool {
	// This is a simplified phone number regex that matches various formats
	phoneRegex := regexp.MustCompile(`^(\+\d{1,3}[- ]?)?\(?(\d{3})\)?[- ]?(\d{3})[- ]?(\d{4})$`)
	return func(value string) bool {
		return phoneRegex.MatchString(value)
	}
}

// getCustomRegexValidator returns a validator function for custom regex formats
func getCustomRegexValidator(pattern string) (func(string) bool, error) {
	if pattern == "" {
		return nil, fmt.Errorf("custom regex pattern is required")
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return func(value string) bool {
		return regex.MatchString(value)
	}, nil
}
