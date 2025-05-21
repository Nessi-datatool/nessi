package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// ProgressIndicator represents a progress indicator for long-running operations
type ProgressIndicator struct {
	spinner    *spinner.Spinner
	message    string
	startTime  time.Time
	isRunning  bool
	showTiming bool
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(message string, showTiming bool) *ProgressIndicator {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Color("cyan")

	return &ProgressIndicator{
		spinner:    s,
		message:    message,
		showTiming: showTiming,
	}
}

// Start starts the progress indicator
func (p *ProgressIndicator) Start() {
	p.startTime = time.Now()
	p.spinner.Start()
	p.isRunning = true
}

// Stop stops the progress indicator
func (p *ProgressIndicator) Stop() {
	if !p.isRunning {
		return
	}
	p.spinner.Stop()
	p.isRunning = false
}

// Success shows a success message after stopping the progress indicator
func (p *ProgressIndicator) Success(message string) {
	p.Stop()
	success := color.New(color.FgGreen, color.Bold)
	success.Printf("\n✅ %s\n", message)

	if p.showTiming {
		elapsed := time.Since(p.startTime)
		fmt.Printf("   Completed in %.2f seconds\n", elapsed.Seconds())
	}
}

// Warning shows a warning message after stopping the progress indicator
func (p *ProgressIndicator) Warning(message string) {
	p.Stop()
	warning := color.New(color.FgYellow, color.Bold)
	warning.Printf("\n⚠️ %s\n", message)

	if p.showTiming {
		elapsed := time.Since(p.startTime)
		fmt.Printf("   Completed with warnings in %.2f seconds\n", elapsed.Seconds())
	}
}

// Error shows an error message after stopping the progress indicator
func (p *ProgressIndicator) Error(message string) {
	p.Stop()
	errorColor := color.New(color.FgRed, color.Bold)
	errorColor.Printf("\n❌ %s\n", message)

	if p.showTiming {
		elapsed := time.Since(p.startTime)
		fmt.Printf("   Failed after %.2f seconds\n", elapsed.Seconds())
	}
}

// UpdateMessage updates the message shown by the progress indicator
func (p *ProgressIndicator) UpdateMessage(message string) {
	p.message = message
	p.spinner.Suffix = " " + message
}

// ProgressBar represents a text-based progress bar
type ProgressBar struct {
	total       int
	current     int
	width       int
	description string
	startTime   time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, description string, width int) *ProgressBar {
	if width <= 0 {
		width = 50 // Default width
	}

	return &ProgressBar{
		total:       total,
		current:     0,
		width:       width,
		description: description,
		startTime:   time.Now(),
	}
}

// Update updates the progress bar and redraws it
func (p *ProgressBar) Update(current int) {
	p.current = current
	p.draw()
}

// Increment increments the progress bar by 1 and redraws it
func (p *ProgressBar) Increment() {
	p.current++
	p.draw()
}

// draw draws the progress bar
func (p *ProgressBar) draw() {
	percent := float64(p.current) / float64(p.total)
	filled := int(percent * float64(p.width))

	// Ensure filled doesn't exceed width
	if filled > p.width {
		filled = p.width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)

	// Calculate ETA
	elapsed := time.Since(p.startTime)
	var eta time.Duration
	if p.current > 0 {
		eta = time.Duration(float64(elapsed) * float64(p.total-p.current) / float64(p.current))
	}

	// Format ETA
	etaStr := ""
	if p.current > 0 && p.current < p.total {
		etaStr = fmt.Sprintf(" ETA: %s", formatDuration(eta))
	}

	// Print progress bar
	fmt.Printf("\r%s: [%s] %d/%d (%d%%)%s",
		p.description,
		bar,
		p.current,
		p.total,
		int(percent*100),
		etaStr,
	)

	// Print newline if complete
	if p.current >= p.total {
		elapsed := time.Since(p.startTime)
		fmt.Printf(" Completed in %s\n", formatDuration(elapsed))
	}
}

// formatDuration formats a duration to a human-readable string
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
