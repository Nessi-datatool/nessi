package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Define the directory to serve
	reportsDir := filepath.Join("reports")

	// Create a file server handler
	fs := http.FileServer(http.Dir(reportsDir))

	// Handle requests to the root URL
	http.Handle("/", http.StripPrefix("/", fs))

	// Print available reports
	fmt.Println("Available reports:")
	fmt.Println("- Quality Report: http://localhost:6666/direct-quality-report-fixed.html")
	fmt.Println("- Schema Report: http://localhost:6666/direct-schema-report-fixed.html")
	fmt.Println("- Freshness Report: http://localhost:6666/direct-freshness-report-fixed.html")
	fmt.Println("- Performance Report: http://localhost:6666/direct-performance-report-fixed.html")

	// Start the server
	fmt.Println("Starting server on port 6666...")
	log.Fatal(http.ListenAndServe(":6666", nil))
}
