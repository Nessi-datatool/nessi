package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Error codes
type ErrorCode string

const (
	ErrInvalidPath      ErrorCode = "N101"
	ErrInvalidConfig    ErrorCode = "N201"
	ErrAuthFailed       ErrorCode = "N301"
	ErrConnectionFailed ErrorCode = "N401"
	ErrInvalidDeltaTable ErrorCode = "N501"
)

// NessiError represents a structured error
type NessiError struct {
	Code        ErrorCode
	Message     string
	Details     string
	Suggestions []string
}

// Error implements the error interface
func (e *NessiError) Error() string {
	return e.Message
}

// NewError creates a new NessiError
func NewError(code ErrorCode, message string) *NessiError {
	return &NessiError{
		Code:    code,
		Message: message,
	}
}

// PathValidator validates and normalizes file paths
type PathValidator struct {
	AllowNonExistent bool
	RequireDirectory bool
	RequireFile      bool
}

// NewPathValidator creates a new path validator
func NewPathValidator() *PathValidator {
	return &PathValidator{
		AllowNonExistent: false,
		RequireDirectory: false,
		RequireFile:      false,
	}
}

// ValidatePath validates and normalizes a file path
func (v *PathValidator) ValidatePath(path string) (string, error) {
	// Check if path is empty
	if path == "" {
		return "", NewError(ErrInvalidPath, "Path cannot be empty")
	}

	// Expand home directory if path starts with ~
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[1:])
	}

	// Convert to absolute path if relative
	if !filepath.IsAbs(path) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		path = filepath.Join(cwd, path)
	}

	// Clean the path
	path = filepath.Clean(path)

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) && !v.AllowNonExistent {
			return "", NewError(ErrInvalidPath, fmt.Sprintf("Path %s does not exist", path))
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to access path: %w", err)
		}
	} else {
		// Path exists, check if it's a directory or file as required
		if v.RequireDirectory && !info.IsDir() {
			return "", NewError(ErrInvalidPath, fmt.Sprintf("%s is not a directory", path))
		}
		if v.RequireFile && info.IsDir() {
			return "", NewError(ErrInvalidPath, fmt.Sprintf("%s is not a file", path))
		}
	}

	return path, nil
}

// DryRunManager handles dry run operations
type DryRunManager struct {
	Enabled bool
	Actions []string
}

// NewDryRunManager creates a new dry run manager
func NewDryRunManager(enabled bool) *DryRunManager {
	return &DryRunManager{
		Enabled: enabled,
		Actions: []string{},
	}
}

// AddAction adds an action to the dry run manager
func (d *DryRunManager) AddAction(format string, args ...interface{}) {
	action := fmt.Sprintf(format, args...)
	d.Actions = append(d.Actions, action)
}

// ShouldExecute returns whether the operation should be executed
func (d *DryRunManager) ShouldExecute() bool {
	return !d.Enabled
}

// PrintActions prints all actions that would be performed
func (d *DryRunManager) PrintActions() {
	if !d.Enabled || len(d.Actions) == 0 {
		return
	}

	// Create color printers
	header := color.New(color.FgCyan, color.Bold)
	action := color.New(color.FgYellow)

	// Print header
	header.Println("\n=== Dry Run: Actions that would be performed ===")

	// Print actions
	for i, a := range d.Actions {
		fmt.Printf("%d. ", i+1)
		action.Println(a)
	}

	// Print footer
	header.Println("\n=== End of Dry Run ===")
	fmt.Println("To execute these actions, run the command without the --dry-run flag.")
}

// Test functions
func testPathValidation() {
	fmt.Println("\n=== Testing Path Validation ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-test-*")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %s\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp(tempDir, "test-file-*")
	if err != nil {
		fmt.Printf("Failed to create temp file: %s\n", err)
		return
	}
	tempFile.Close()

	// Test cases
	tests := []struct {
		name      string
		validator *PathValidator
		path      string
		wantErr   bool
	}{
		{
			name:      "Empty path",
			validator: NewPathValidator(),
			path:      "",
			wantErr:   true,
		},
		{
			name:      "Non-existent path not allowed",
			validator: NewPathValidator(),
			path:      filepath.Join(tempDir, "non-existent"),
			wantErr:   true,
		},
		{
			name: "Non-existent path allowed",
			validator: &PathValidator{
				AllowNonExistent: true,
			},
			path:    filepath.Join(tempDir, "non-existent"),
			wantErr: false,
		},
		{
			name: "Directory required but file provided",
			validator: &PathValidator{
				RequireDirectory: true,
			},
			path:    tempFile.Name(),
			wantErr: true,
		},
		{
			name: "File required but directory provided",
			validator: &PathValidator{
				RequireFile: true,
			},
			path:    tempDir,
			wantErr: true,
		},
		{
			name:      "Valid directory",
			validator: NewPathValidator(),
			path:      tempDir,
			wantErr:   false,
		},
		{
			name:      "Valid file",
			validator: NewPathValidator(),
			path:      tempFile.Name(),
			wantErr:   false,
		},
	}

	// Run tests
	for _, tt := range tests {
		fmt.Printf("Test: %s\n", tt.name)
		gotPath, err := tt.validator.ValidatePath(tt.path)
		if tt.wantErr {
			if err != nil {
				color.Green("✓ Got expected error: %s\n", err)
			} else {
				color.Red("✗ Expected error but got none\n")
			}
		} else {
			if err != nil {
				color.Red("✗ Unexpected error: %s\n", err)
			} else {
				color.Green("✓ Path validated successfully: %s\n", gotPath)
			}
		}
	}
}

func testDryRun() {
	fmt.Println("\n=== Testing Dry Run Mode ===")

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-test-*")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %s\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Test with dry run enabled
	fmt.Println("\nTest: Dry run enabled")
	manager := NewDryRunManager(true)
	manager.AddAction("Create directory: %s", filepath.Join(tempDir, "test-dir"))
	manager.AddAction("Write file: %s", filepath.Join(tempDir, "test-file.txt"))

	if manager.ShouldExecute() {
		color.Red("✗ ShouldExecute() returned true when dry run is enabled\n")
	} else {
		color.Green("✓ ShouldExecute() correctly returned false when dry run is enabled\n")
	}

	// Print actions
	manager.PrintActions()

	// Test with dry run disabled
	fmt.Println("\nTest: Dry run disabled")
	manager = NewDryRunManager(false)
	manager.AddAction("Create directory: %s", filepath.Join(tempDir, "test-dir"))

	if manager.ShouldExecute() {
		color.Green("✓ ShouldExecute() correctly returned true when dry run is disabled\n")

		// Actually create the directory
		testDir := filepath.Join(tempDir, "test-dir")
		if err := os.Mkdir(testDir, 0755); err != nil {
			color.Red("✗ Failed to create directory: %s\n", err)
		} else {
			color.Green("✓ Directory created successfully: %s\n", testDir)
		}
	} else {
		color.Red("✗ ShouldExecute() returned false when dry run is disabled\n")
	}

	// Print actions (should not print anything since dry run is disabled)
	manager.PrintActions()
}

func main() {
	fmt.Println("=== Nessi Usability Features Test ===\n")

	// Test path validation
	testPathValidation()

	// Test dry run mode
	testDryRun()

	fmt.Println("\n=== Tests Completed ===")
}
