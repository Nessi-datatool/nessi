package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi-dev/pkg/security"
)

func main() {
	fmt.Println("Security Features Verification")
	fmt.Println("=============================")

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "security-verify-")
	if err != nil {
		fmt.Printf("Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Create paths for test files
	certFile := filepath.Join(tempDir, "cert.pem")
	keyFile := filepath.Join(tempDir, "key.pem")

	// Test SSL functionality
	fmt.Println("\n1. Testing SSL Certificate Generation")
	sslConfig := security.SSLConfig{
		Enabled:      true,
		CertFile:     certFile,
		KeyFile:      keyFile,
		AutoGenerate: true,
	}
	
	certManager := security.NewCertManager(sslConfig)
	tlsConfig, err := certManager.GetTLSConfig()
	if err != nil {
		fmt.Printf("❌ Failed to generate TLS config: %v\n", err)
	} else if tlsConfig != nil {
		fmt.Println("✅ Successfully generated TLS config")
	}
	
	// Check if certificate files were created
	if _, err := os.Stat(certFile); err == nil {
		fmt.Println("✅ Certificate file created successfully")
	} else {
		fmt.Printf("❌ Certificate file not created: %v\n", err)
	}
	
	if _, err := os.Stat(keyFile); err == nil {
		fmt.Println("✅ Key file created successfully")
	} else {
		fmt.Printf("❌ Key file not created: %v\n", err)
	}

	fmt.Println("\nSecurity verification completed")
}
