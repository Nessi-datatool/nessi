package security

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// ShouldSkipIntegrationTests returns true if integration tests should be skipped
// based on environment variables or test flags
func ShouldSkipIntegrationTests(t *testing.T) bool {
	// Skip in short mode (go test -short)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
		return true
	}
	
	// Check for environment variable to skip security integration tests
	if os.Getenv("SKIP_SECURITY_INTEGRATION") == "1" {
		t.Skip("Skipping security integration test due to SKIP_SECURITY_INTEGRATION=1")
		return true
	}
	
	return false
}

// GetTestTimeout returns an appropriate timeout duration for tests
// Uses a consistent timeout that's fast enough for CI but allows tests to complete
func GetTestTimeout() time.Duration {
	// 2 seconds is a good balance - fast enough for CI but allows tests to complete
	return 2 * time.Second
}

// RunWithTimeout runs a test function with a timeout
// If no timeout is provided, it uses the default from GetTestTimeout
func RunWithTimeout(t *testing.T, testFunc func(), timeout ...time.Duration) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		testFunc()
	}()

	// Use provided timeout or default
	actualTimeout := GetTestTimeout()
	if len(timeout) > 0 {
		actualTimeout = timeout[0]
	}

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(actualTimeout):
		t.Fatal("Test timed out after", actualTimeout)
	}
}

// CreateTestAuthManager creates an optimized AuthManager for testing with in-memory storage
// This is the recommended way to create an AuthManager for all tests
func CreateTestAuthManager() (*AuthManager, error) {
	// Use a random suffix to ensure unique usernames across tests
	random := fmt.Sprintf("%d", time.Now().UnixNano()%10000)

	// Create a configuration that doesn't use file I/O for maximum performance
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-secret-" + random,
		UsersFile:    "/tmp/nonexistent-users-file.json", // Won't be used
		TokenExpiry:  1,                                  // 1 hour - minimum value for faster tests
		RequireHTTPS: false,
		InMemoryOnly: true, // Skip all file I/O for better performance
	}

	// Create the auth manager
	am, err := NewAuthManager(config)
	if err != nil {
		return nil, err
	}

	// Pre-create a test user for convenience with unique username
	testUser := User{
		Username: "testuser-" + random,
		Email:    "test" + random + "@example.com",
		Role:     RoleUser,
	}
	
	// Create user (should not error since username is unique)
	err = am.CreateUser(testUser, "password123")
	if err != nil {
		return nil, fmt.Errorf("failed to create test user: %v", err)
	}

	// Pre-create an admin user for convenience with unique username
	adminUser := User{
		Username: "adminuser-" + random,
		Email:    "admin" + random + "@example.com",
		Role:     RoleAdmin,
	}
	
	// Create admin user (should not error since username is unique)
	err = am.CreateUser(adminUser, "admin123")
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %v", err)
	}

	return am, nil
}

// CreateTestSecurityManager creates a SecurityManager for testing
// This wraps CreateTestAuthManager for legacy tests
func CreateTestSecurityManager() *SecurityManager {
	am, err := CreateTestAuthManager()
	if err != nil {
		panic(fmt.Sprintf("Failed to create test auth manager: %v", err))
	}

	return &SecurityManager{
		AuthManager: am,
	}
}
