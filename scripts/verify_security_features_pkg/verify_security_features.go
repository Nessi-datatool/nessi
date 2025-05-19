package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi/pkg/security"
)

// This script provides a quick verification of the security features
// without running the full test suite.

func main() {
	fmt.Println("Nessi Security Features Verification")
	fmt.Println("====================================")

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "security-verify-")
	if err != nil {
		fmt.Printf("Error creating temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Create users file path
	usersFilePath := filepath.Join(tempDir, "users.json")

	// Create empty users file with valid JSON array
	usersFile, err := os.Create(usersFilePath)
	if err != nil {
		fmt.Printf("Error creating users file: %v\n", err)
		os.Exit(1)
	}
	_, err = usersFile.WriteString("[]")
	if err != nil {
		fmt.Printf("Error writing to users file: %v\n", err)
		os.Exit(1)
	}
	usersFile.Close()

	// Create certificate paths
	certFile := filepath.Join(tempDir, "cert.pem")
	keyFile := filepath.Join(tempDir, "key.pem")

	// Create auth config
	authConfig := security.AuthConfig{
		Enabled:   true,
		UsersFile: usersFilePath,
	}

	// Create SSL config
	sslConfig := security.SSLConfig{
		Enabled:  true,
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	fmt.Println("\n1. Creating Auth Manager...")
	authManager, err := security.NewAuthManager(authConfig)
	if err != nil {
		fmt.Printf("Error creating auth manager: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Auth Manager created successfully")

	fmt.Println("\n2. Creating Cert Manager...")
	certManager := security.NewCertManager(sslConfig)
	fmt.Println("✓ Cert Manager created successfully")

	fmt.Println("\n3. Testing User Management...")
	// Create test user
	testUser := security.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     security.RoleUser,
	}
	err = authManager.CreateUser(testUser, "testpassword")
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ User created successfully")

	// Get users
	users, err := authManager.GetUsers()
	if err != nil {
		fmt.Printf("Error getting users: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Found %d users\n", len(users))

	fmt.Println("\n4. Testing Authentication...")
	// Authenticate with password
	authenticated, err := authManager.Authenticate("testuser", "testpassword")
	if err != nil {
		fmt.Printf("Error authenticating: %v\n", err)
		os.Exit(1)
	}
	if !authenticated {
		fmt.Println("Error: Authentication failed but no error was returned")
		os.Exit(1)
	}
	fmt.Println("✓ Password authentication successful")

	fmt.Println("\n5. Testing API Key Management...")
	// Generate API key
	apiKey, err := authManager.RegenerateAPIKey("testuser")
	if err != nil {
		fmt.Printf("Error generating API key: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ API key generated successfully")

	// Authenticate with API key
	user, err := authManager.AuthenticateWithAPIKey(apiKey)
	if err != nil {
		fmt.Printf("Error authenticating with API key: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ API key authentication successful for user: %s\n", user.Username)

	fmt.Println("\n6. Testing SSL Certificate Generation...")
	// Get TLS config
	_, err = certManager.GetTLSConfig()
	if err != nil {
		fmt.Printf("Error getting TLS config: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ TLS configuration generated successfully")

	// Verify certificate files were created
	_, err = os.Stat(certFile)
	if err != nil {
		fmt.Printf("Error verifying cert file: %v\n", err)
		os.Exit(1)
	}
	_, err = os.Stat(keyFile)
	if err != nil {
		fmt.Printf("Error verifying key file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Certificate files created successfully")

	fmt.Println("\nAll security features verified successfully!")
}
