package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MockViralCommand creates a mock viral command executable for testing
func MockViralCommand() (string, error) {
	// Create a temporary directory for the mock executable
	tmpDir, err := os.MkdirTemp("", "nessi-test-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Create a mock nessi executable script
	mockPath := filepath.Join(tmpDir, "nessi")
	mockScript := `#!/bin/bash
set -e
if [[ "$1" == "viral" && "$2" == "share" ]]; then
    echo "Generating shareable report"
    echo "Report generated successfully"
    
    # Parse arguments
    output_file=""
    hashtags="nessi,dataquality"
    
    for i in $(seq 1 $#); do
        arg=${!i}
        next_i=$((i+1))
        next_arg=""
        if [[ $next_i -le $# ]]; then
            next_arg=${!next_i}
        fi
        
        if [[ "$arg" == "--output" && -n "$next_arg" ]]; then
            output_file="$next_arg"
        elif [[ "$arg" == "--hashtags" && -n "$next_arg" ]]; then
            hashtags="$next_arg"
        fi
    done
    
    # Create output file
    if [[ -n "$output_file" ]]; then
        mkdir -p "$(dirname "$output_file")"
        echo "<html>Report content with $hashtags</html>" > "$output_file"
    fi
elif [[ "$1" == "viral" && "$2" == "badge" ]]; then
    echo "Generating 'Powered by Nessi' badge"
    echo "Badge generated successfully"
    
    # Parse arguments
    output_file=""
    format="markdown"
    label="powered by"
    message="nessi"
    color="blue"
    quality_score="0"
    
    for i in $(seq 1 $#); do
        arg=${!i}
        next_i=$((i+1))
        next_arg=""
        if [[ $next_i -le $# ]]; then
            next_arg=${!next_i}
        fi
        
        if [[ "$arg" == "--output" && -n "$next_arg" ]]; then
            output_file="$next_arg"
        elif [[ "$arg" == "--format" && -n "$next_arg" ]]; then
            format="$next_arg"
        elif [[ "$arg" == "--label" && -n "$next_arg" ]]; then
            label="$next_arg"
        elif [[ "$arg" == "--message" && -n "$next_arg" ]]; then
            message="$next_arg"
        elif [[ "$arg" == "--color" && -n "$next_arg" ]]; then
            color="$next_arg"
        elif [[ "$arg" == "--quality-score" && -n "$next_arg" ]]; then
            quality_score="$next_arg"
        fi
    done
    
    # Create output file
    if [[ -n "$output_file" ]]; then
        mkdir -p "$(dirname "$output_file")"
        if [[ "$format" == "html" ]]; then
            echo "<img src='https://img.shields.io/badge/$label-$message-$color'>" > "$output_file"
        else
            badge_text="$label $message"
            if [[ "$quality_score" != "0" ]]; then
                badge_text="$badge_text $quality_score%"
            fi
            echo "[![$badge_text](https://img.shields.io/badge/$label-$message-$color)](https://nessi-dev.github.io)" > "$output_file"
        fi
    fi
elif [[ "$1" == "viral" && "$2" == "community" && "$3" == "feedback" ]]; then
    echo "Submitting feedback"
    echo "Feedback submitted successfully"
    echo "GitHub issue URL: https://github.com/nessi-dev/nessi/issues/new"
elif [[ "$1" == "viral" && "$2" == "community" && "$3" == "contribute" ]]; then
    experience=""
    for arg in "$@"; do
        if [[ "$arg" == "--experience" ]]; then
            # Get the next argument
            get_next=true
        elif [[ "$get_next" == "true" ]]; then
            experience="$arg"
            break
        fi
    done
    
    echo "Finding contribution suggestions"
    if [[ "$experience" == "beginner" ]]; then
        echo "Contribution suggestions for beginners"
        echo "Fix a Documentation Typo (Very Easy)"
        echo "documentation"
    else
        echo "Contribution suggestions for advanced"
        echo "Add a New Core Feature (Advanced)"
        echo "core"
    fi
elif [[ "$1" == "plugins" && "$2" == "list" ]]; then
    echo "[{\"name\":\"social_sharing\",\"version\":\"1.0.0\"},{\"name\":\"badge\",\"version\":\"1.0.0\"},{\"name\":\"community\",\"version\":\"1.0.0\"}]"
elif [[ "$1" == "help" && "$2" == "viral" ]]; then
    echo "Nessi Viral Growth Tools"
    echo "Available commands:"
    echo "  viral share     - Generate shareable reports"
    echo "  viral badge     - Create 'Powered by Nessi' badges"
    echo "  viral community - Engage with the Nessi community"
else
    echo "Error: unknown command \"$1\" for \"nessi\""
    echo "Run 'nessi --help' for usage."
    echo "unknown command \"$1\" for \"nessi\""
    exit 1
fi
`

	// Write the mock script
	err = os.WriteFile(mockPath, []byte(mockScript), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to write mock script: %w", err)
	}

	// Update PATH to include our mock executable
	oldPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s:%s", tmpDir, oldPath)
	os.Setenv("PATH", newPath)

	// Verify the mock works
	cmd := exec.Command("nessi", "viral", "share", "test_table")
	output, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "Report generated successfully") {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("mock verification failed: %v, output: %s", err, output)
	}

	return tmpDir, nil
}

// CleanupMockViralCommand removes the mock viral command
func CleanupMockViralCommand(tmpDir string) {
	if tmpDir != "" {
		os.RemoveAll(tmpDir)
	}
}
