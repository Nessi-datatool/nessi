package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorSuggestions(t *testing.T) {
	// Initialize default suggestions
	InitDefaultSuggestions()

	// Test getting suggestions for a NessiError
	nessiErr := NewError(ErrInvalidPath, "test path is invalid")
	suggestions := GetSuggestionsForError(nessiErr)

	// Should have at least one suggestion for invalid path
	assert.GreaterOrEqual(t, len(suggestions), 1)

	// Check suggestion fields
	for _, suggestion := range suggestions {
		assert.Equal(t, ErrInvalidPath, suggestion.ErrorCode)
		assert.NotEmpty(t, suggestion.Description)
		assert.NotEmpty(t, suggestion.Solution)
	}

	// Test getting suggestions for a standard error
	stdErr := fmt.Errorf("standard error")
	suggestions = GetSuggestionsForError(stdErr)

	// Should have no suggestions for standard error
	assert.Empty(t, suggestions)

	// Test formatting suggestions
	formattedSuggestions := FormatSuggestions(nessiErr)
	assert.Contains(t, formattedSuggestions, "Suggested solutions")
	assert.Contains(t, formattedSuggestions, "Solution:")

	// Test formatting when no suggestions are available
	formattedSuggestions = FormatSuggestions(stdErr)
	assert.Empty(t, formattedSuggestions)
}

func TestRegisterErrorSuggestion(t *testing.T) {
	// Clear existing suggestions
	errorSuggestions = map[ErrorCode][]ErrorSuggestion{}

	// Register a new suggestion
	newSuggestion := ErrorSuggestion{
		ErrorCode:   ErrInvalidConfig,
		Description: "Test description",
		Solution:    "Test solution",
		DocumentURL: "https://example.com",
	}
	RegisterErrorSuggestion(newSuggestion)

	// Check if suggestion was registered
	suggestions := errorSuggestions[ErrInvalidConfig]
	assert.Equal(t, 1, len(suggestions))
	assert.Equal(t, newSuggestion, suggestions[0])

	// Register another suggestion for the same error code
	anotherSuggestion := ErrorSuggestion{
		ErrorCode:   ErrInvalidConfig,
		Description: "Another description",
		Solution:    "Another solution",
		DocumentURL: "https://example.org",
	}
	RegisterErrorSuggestion(anotherSuggestion)

	// Check if both suggestions are registered
	suggestions = errorSuggestions[ErrInvalidConfig]
	assert.Equal(t, 2, len(suggestions))
	assert.Equal(t, newSuggestion, suggestions[0])
	assert.Equal(t, anotherSuggestion, suggestions[1])
}

func TestFormatSuggestion(t *testing.T) {
	// Test formatting a suggestion with documentation URL
	suggestion := ErrorSuggestion{
		ErrorCode:   ErrInvalidPath,
		Description: "Test description",
		Solution:    "Test solution",
		DocumentURL: "https://example.com",
	}

	formatted := FormatSuggestion(suggestion)
	assert.Contains(t, formatted, "Test description")
	assert.Contains(t, formatted, "Test solution")
	assert.Contains(t, formatted, "https://example.com")

	// Test formatting a suggestion without documentation URL
	suggestion = ErrorSuggestion{
		ErrorCode:   ErrInvalidPath,
		Description: "Test description",
		Solution:    "Test solution",
		DocumentURL: "",
	}

	formatted = FormatSuggestion(suggestion)
	assert.Contains(t, formatted, "Test description")
	assert.Contains(t, formatted, "Test solution")
	assert.NotContains(t, formatted, "Documentation:")
}
