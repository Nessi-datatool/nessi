# Progress Visualization in Nessi

This document describes the progress visualization features in Nessi, which provide real-time feedback during long-running operations.

## Overview

Nessi includes a robust progress visualization system to enhance the user experience by showing what's happening during command execution. This is especially important for long-running operations where users need feedback on the status of their tasks.

The progress visualization system includes:

1. **Progress Indicators** - Spinner-style indicators for operations without measurable progress
2. **Progress Bars** - Visual bars for operations with measurable percentage completion
3. **Success/Warning/Error Messages** - Clear status messages when operations complete

## Progress Indicators

Progress indicators are used for operations that don't have a measurable percentage completion, such as API calls or operations where the total work is unknown.

### Example Usage

```go
// In your command implementation
import "github.com/nessi-dev/nessi/pkg/common"

func RunCommand() {
    // Show a progress indicator
    common.ShowProgress("Processing data...", 2*time.Second)
    
    // Perform the operation
    // ...
    
    // Show success message
    common.ShowSuccess("Data processed successfully")
}
```

### User Experience

When a command uses progress indicators, users will see a spinner with a message:

```
⠋ Processing data...
```

The spinner will continue to animate until the operation completes, at which point a success message is shown:

```
✅ Data processed successfully
```

## Progress Bars

Progress bars are used for operations where the percentage completion can be measured, such as processing a known number of files or records.

### Example Usage

```go
// In your command implementation
import "github.com/nessi-dev/nessi/pkg/common"

func RunCommand() {
    // Create a progress bar with 100 total steps
    bar := common.NewProgressBar(100, "Processing files")
    
    // Update the progress as the operation proceeds
    for i := 0; i < 100; i++ {
        // Perform work on one item
        // ...
        
        // Update the progress bar
        bar.Increment()
        time.Sleep(50 * time.Millisecond) // Simulate work
    }
    
    // Complete the progress bar
    bar.Complete()
    
    // Show success message
    common.ShowSuccess("All files processed successfully")
}
```

### User Experience

When a command uses progress bars, users will see a visual bar with percentage and ETA:

```
[████████████████████░░░░░░░░░░] 67% (ETA: 45s) Processing files
```

The progress bar will update in real-time as the operation proceeds, and when complete, a success message is shown.

## Status Messages

Nessi provides three types of status messages to indicate the outcome of operations:

### Success Messages

Success messages are shown when operations complete successfully:

```go
common.ShowSuccess("Operation completed successfully")
```

This will display:

```
✅ Operation completed successfully
```

### Warning Messages

Warning messages are shown when operations complete with warnings:

```go
common.ShowWarning("Operation completed with warnings")
```

This will display:

```
⚠️ Operation completed with warnings
```

### Error Messages

Error messages are shown when operations fail:

```go
err := fmt.Errorf("Something went wrong: %s", errorDetails)
common.HandleError(err, "operation")
```

This will display:

```
❌ Error: Something went wrong: [error details]

🔍 Suggested solutions:
...
```

## Progress Demo

Nessi includes a progress demo command that showcases all the progress visualization features:

```bash
# Run the progress visualization demo
nessi progress-demo
```

This will demonstrate:

1. Spinner-style progress indicators
2. Progress bars with different speeds and styles
3. Success messages
4. Warning messages
5. Error messages

## Integration with Viral Growth Features

All viral growth features in Nessi use progress visualization:

### Share Command

```bash
nessi viral share my_table --output report.html
```

Shows a progress indicator during report generation:

```
⠋ Generating shareable report...
```

And a success message when complete:

```
✅ Successfully generated shareable report
```

### Badge Command

```bash
nessi viral badge --format markdown
```

Shows a progress indicator during badge generation:

```
⠋ Generating 'Powered by Nessi' badge...
```

And a success message when complete:

```
✅ Badge generated successfully
```

### Community Feedback Command

```bash
nessi viral community feedback --text "Great tool!"
```

Shows a progress indicator during feedback submission:

```
⠋ Submitting feedback...
```

And a success message when complete:

```
✅ Feedback submitted successfully
```

## Error Handling Integration

The progress visualization system is integrated with Nessi's error handling system, providing clear, actionable error messages when operations fail.

For example, if a user tries to generate a report with an invalid output path:

```bash
nessi viral share my_table --output /nonexistent/path/report.html
```

They will see an error message with a code and suggestions:

```
❌ Error: Invalid path: /nonexistent/path does not exist [Error Code: V101]

🔍 Suggested solutions:

📝 The specified output path does not exist
🔧 Solution: Check the path and ensure the directory exists
```

## Implementation Details

The progress visualization system is implemented in the following files:

- `pkg/common/progress.go`: Main implementation of progress indicators and bars
- `cmd/nessi/cli/viral_error_handler.go`: Integration with viral growth features
- `cmd/progress_demo/main.go`: Standalone demo program

## Best Practices

When implementing new commands in Nessi, follow these best practices for progress visualization:

1. **Use progress indicators for operations without measurable progress**
2. **Use progress bars for operations with measurable percentage completion**
3. **Always show success messages when operations complete successfully**
4. **Show warning messages when operations complete with non-critical issues**
5. **Integrate with the error handling system for failures**
6. **Provide clear, actionable error messages with suggestions**

## Conclusion

The progress visualization system in Nessi enhances the user experience by providing real-time feedback during command execution. By following the patterns and best practices described in this document, you can ensure that your commands provide a consistent and user-friendly experience.
