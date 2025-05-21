package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

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

// Example usage in a command:
//
// dryRun, _ := cmd.Flags().GetBool("dry-run")
// manager := common.NewDryRunManager(dryRun)
//
// // Add actions that would be performed
// manager.AddAction("Create directory: %s", path)
// manager.AddAction("Write file: %s", filePath)
//
// // Check if we should execute
// if manager.ShouldExecute() {
//     // Perform the actual operations
//     os.MkdirAll(path, 0755)
//     os.WriteFile(filePath, data, 0644)
// }
//
// // Print actions at the end
// manager.PrintActions()

// DryRunnable is an interface for operations that support dry run mode
type DryRunnable interface {
	// DryRun performs a dry run of the operation
	DryRun() error

	// Execute performs the actual operation
	Execute() error
}

// FileOperation represents a file operation that supports dry run
type FileOperation struct {
	Type        string // "create", "update", "delete", etc.
	Path        string // File path
	Description string // Description of the operation
	DryRun      bool   // Whether this is a dry run
}

// NewCreateFileOperation creates a new file creation operation
func NewCreateFileOperation(path string, dryRun bool) *FileOperation {
	return &FileOperation{
		Type:        "create",
		Path:        path,
		Description: fmt.Sprintf("Create file: %s", path),
		DryRun:      dryRun,
	}
}

// NewUpdateFileOperation creates a new file update operation
func NewUpdateFileOperation(path string, dryRun bool) *FileOperation {
	return &FileOperation{
		Type:        "update",
		Path:        path,
		Description: fmt.Sprintf("Update file: %s", path),
		DryRun:      dryRun,
	}
}

// NewDeleteFileOperation creates a new file deletion operation
func NewDeleteFileOperation(path string, dryRun bool) *FileOperation {
	return &FileOperation{
		Type:        "delete",
		Path:        path,
		Description: fmt.Sprintf("Delete file: %s", path),
		DryRun:      dryRun,
	}
}

// DryRun performs a dry run of the file operation
func (f *FileOperation) DryRun() error {
	// Print the operation description
	info := color.New(color.FgYellow)
	info.Printf("[DRY RUN] %s\n", f.Description)

	// For create/update operations, check if the file exists
	if f.Type == "create" || f.Type == "update" {
		if _, err := os.Stat(f.Path); err == nil && f.Type == "create" {
			return fmt.Errorf("file already exists: %s", f.Path)
		} else if os.IsNotExist(err) && f.Type == "update" {
			return fmt.Errorf("file does not exist: %s", f.Path)
		}
	}

	// For delete operations, check if the file exists
	if f.Type == "delete" {
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", f.Path)
		}
	}

	return nil
}

// Execute performs the actual file operation
func (f *FileOperation) Execute() error {
	// If this is a dry run, just return
	if f.DryRun {
		return f.DryRun()
	}

	// Print the operation description
	info := color.New(color.FgGreen)
	info.Printf("%s\n", f.Description)

	// Perform the operation based on the type
	switch f.Type {
	case "create":
		// For create, we just check if the file exists
		// The actual creation is done by the caller
		if _, err := os.Stat(f.Path); err == nil {
			return fmt.Errorf("file already exists: %s", f.Path)
		}
	case "update":
		// For update, we just check if the file exists
		// The actual update is done by the caller
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", f.Path)
		}
	case "delete":
		// For delete, we remove the file
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", f.Path)
		}
		return os.Remove(f.Path)
	default:
		return fmt.Errorf("unknown operation type: %s", f.Type)
	}

	return nil
}
