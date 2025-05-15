package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi-dev/pkg/security"
)

// TestResult represents the result of a test
type TestResult struct {
	Name        string
	Passed      bool
	Description string
	Error       error
}

func main() {
	fmt.Println("Running Security Features Test Script")
	fmt.Println("=====================================")

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "security-test-")
	if err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create users file
	usersFile := filepath.Join(tempDir, "users.json")
	if err := os.WriteFile(usersFile, []byte("[]"), 0644); err != nil {
		log.Fatalf("Failed to create users file: %v", err)
	}

	// Create certificate files
	certFile := filepath.Join(tempDir, "cert.pem")
	keyFile := filepath.Join(tempDir, "key.pem")

	// Create auth config
	authConfig := security.AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-secret",
		UsersFile:    usersFile,
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create SSL config
	sslConfig := security.SSLConfig{
		Enabled:      true,
		CertFile:     certFile,
		KeyFile:      keyFile,
		AutoGenerate: true,
	}

	// Run tests
	results := []TestResult{}

	// Test 1: Create auth manager
	authManager, err := security.NewAuthManager(authConfig)
	results = append(results, TestResult{
		Name:        "Create Auth Manager",
		Passed:      err == nil,
		Description: "Creating authentication manager",
		Error:       err,
	})

	if err != nil {
		printResults(results)
		return
	}

	// Test 2: Create cert manager
	certManager := security.NewCertManager(sslConfig)
	tlsConfig, err := certManager.GetTLSConfig()
	results = append(results, TestResult{
		Name:        "Create Cert Manager",
		Passed:      err == nil && tlsConfig != nil,
		Description: "Creating certificate manager and getting TLS config",
		Error:       err,
	})

	// Test 3: Create user
	testUser := security.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     security.RoleUser,
	}
	err = authManager.CreateUser(testUser, "password123")
	results = append(results, TestResult{
		Name:        "Create User",
		Passed:      err == nil,
		Description: "Creating a test user",
		Error:       err,
	})

	// Test 4: Authenticate user
	token, err := authManager.Authenticate("testuser", "password123")
	results = append(results, TestResult{
		Name:        "Authenticate User",
		Passed:      err == nil && token != "",
		Description: "Authenticating with username and password",
		Error:       err,
	})

	// Test 5: Validate token
	claims, err := authManager.ValidateToken(token)
	results = append(results, TestResult{
		Name:        "Validate Token",
		Passed:      err == nil && claims["username"] == "testuser",
		Description: "Validating JWT token",
		Error:       err,
	})

	// Test 6: Generate API key
	apiKey, err := authManager.RegenerateAPIKey("testuser")
	results = append(results, TestResult{
		Name:        "Generate API Key",
		Passed:      err == nil && apiKey != "",
		Description: "Generating API key for user",
		Error:       err,
	})

	// Test 7: Authenticate with API key
	user, err := authManager.AuthenticateWithAPIKey(apiKey)
	results = append(results, TestResult{
		Name:        "Authenticate with API Key",
		Passed:      err == nil && user.Username == "testuser",
		Description: "Authenticating with API key",
		Error:       err,
	})

	// Test 8: Create admin user
	adminUser := security.User{
		Username: "adminuser",
		Email:    "admin@example.com",
		Role:     security.RoleAdmin,
	}
	err = authManager.CreateUser(adminUser, "admin123")
	results = append(results, TestResult{
		Name:        "Create Admin User",
		Passed:      err == nil,
		Description: "Creating an admin user",
		Error:       err,
	})

	// Test 9: Test role middleware
	adminMiddleware := authManager.RoleMiddleware(security.RoleAdmin)
	results = append(results, TestResult{
		Name:        "Create Role Middleware",
		Passed:      adminMiddleware != nil,
		Description: "Creating role-based middleware",
		Error:       nil,
	})

	// Test 10: Test SSL certificate generation
	_, err = os.Stat(certFile)
	certExists := err == nil
	_, err = os.Stat(keyFile)
	keyExists := err == nil
	results = append(results, TestResult{
		Name:        "SSL Certificate Generation",
		Passed:      certExists && keyExists,
		Description: "Generating self-signed SSL certificates",
		Error:       err,
	})

	// Print results
	printResults(results)
}

func printResults(results []TestResult) {
	passed := 0
	failed := 0

	fmt.Println("\nTest Results:")
	fmt.Println("============")

	for i, result := range results {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
			failed++
		} else {
			passed++
		}

		fmt.Printf("%d. %s: %s - %s\n", i+1, status, result.Name, result.Description)
		if result.Error != nil {
			fmt.Printf("   Error: %v\n", result.Error)
		}
	}

	fmt.Println("\nSummary:")
	fmt.Printf("Total: %d, Passed: %d, Failed: %d\n", len(results), passed, failed)

	if failed > 0 {
		fmt.Println("\nSome tests failed!")
		os.Exit(1)
	} else {
		fmt.Println("\nAll tests passed!")
	}
}
