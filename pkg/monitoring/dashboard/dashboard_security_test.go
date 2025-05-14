package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/testutil"
	"github.com/nessi-dev/nessi-dev/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDashboardSecurity tests the security features of the dashboard
func TestDashboardSecurity(t *testing.T) {
	// Create a temporary file for auth config
	tempFile, err := os.CreateTemp("", "dashboard-security-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Initialize the file with an empty JSON array
	_, err = tempFile.Write([]byte("[]"))
	require.NoError(t, err)
	tempFile.Close()
	
	// Create auth config
	authConfig := security.AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-security-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
		InMemoryOnly: true,
	}
	
	// Create auth manager
	am, err := security.NewAuthManager(authConfig)
	require.NoError(t, err)
	
	// Create dashboard with security enabled
	options := DashboardOptions{
		ListenAddr:  ":8080",
		AuthManager: am,
	}

	// Create mock monitor
	alertManager, err := alerts.NewAlertManager("/tmp/alerts-test")
	require.NoError(t, err)
	mock := testutil.CreateTestMonitorWithAlertManager(9090, alertManager)

	// Create dashboard
	dash, err := New(mock, options)
	require.NoError(t, err)
	
	// Test that security is enabled
	assert.NotNil(t, dash.authManager)
	
	// Test accessing a protected route without authentication
	req := httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	w := httptest.NewRecorder()
	
	// Create a handler that requires authentication
	handler := dash.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Protected content"))
	}))
	
	// Call the handler
	handler.ServeHTTP(w, req)
	
	// Check that we get an unauthorized response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Create a test user
	testUser := security.User{
		Username: "securityuser",
		Email:    "security@example.com",
		Role:     security.RoleAdmin,
	}
	err = am.CreateUser(testUser, "securitypass")
	require.NoError(t, err)
	
	// Get a token
	token, err := am.Authenticate("securityuser", "securitypass")
	require.NoError(t, err)
	
	// Test accessing a protected route with authentication
	req = httptest.NewRequest(http.MethodGet, "/api/metrics", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	
	// Call the handler
	handler.ServeHTTP(w, req)
	
	// Check that we get a success response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Protected content", w.Body.String())
}
