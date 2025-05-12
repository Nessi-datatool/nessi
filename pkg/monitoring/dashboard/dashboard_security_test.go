package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestDashboardSecurity tests the security features of the dashboard
func TestDashboardSecurity(t *testing.T) {
	// Create a test dashboard
	dashboard := &Dashboard{
		config: &Config{
			Address: ":8080",
			Title:   "Test Dashboard",
		},
	}

	// Test authentication middleware
	t.Run("Authentication Required", func(t *testing.T) {
		// Create a request to a protected endpoint
		req := httptest.NewRequest("GET", "/api/metrics", nil)
		w := httptest.NewRecorder()

		// Call the handler with auth middleware
		handler := dashboard.authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		// Should return unauthorized without credentials
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Test with invalid token
		req = httptest.NewRequest("GET", "/api/metrics", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test login handler
	t.Run("Login Handler", func(t *testing.T) {
		// Create a request to the login page (GET)
		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		// Call the login handler
		dashboard.handleLogin(w, req)

		// Should return OK and contain login form
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "login")

		// Test POST to login with invalid credentials
		loginData := map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		// Create a mock auth manager that rejects authentication
		mockAuth := new(mockAuthManagerExtended)
		mockAuth.On("Authenticate", "testuser", "wrongpassword").Return("", security.ErrInvalidCredentials)
		dashboard.authManager = mockAuth

		req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		dashboard.handleLogin(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Test POST to login with valid credentials
		loginData = map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}
		jsonData, err = json.Marshal(loginData)
		require.NoError(t, err)

		// Update mock to accept authentication
		mockAuth = new(mockAuthManagerExtended)
		mockAuth.On("Authenticate", "testuser", "correctpassword").Return("valid-token", nil)
		dashboard.authManager = mockAuth

		req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		dashboard.handleLogin(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	// Test API login handler
	t.Run("API Login Handler", func(t *testing.T) {
		// Test with valid credentials
		loginData := map[string]string{
			"username": "testuser",
			"password": "correctpassword",
		}
		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		// Create a mock auth manager
		mockAuth := new(mockAuthManagerExtended)
		mockAuth.On("Authenticate", "testuser", "correctpassword").Return("valid-token", nil)
		dashboard.authManager = mockAuth

		req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		dashboard.handleAPILogin(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response contains token
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "valid-token", response["token"])

		// Test with invalid credentials
		loginData = map[string]string{
			"username": "testuser",
			"password": "wrongpassword",
		}
		jsonData, err = json.Marshal(loginData)
		require.NoError(t, err)

		// Update mock to reject authentication
		mockAuth = new(mockAuthManagerExtended)
		mockAuth.On("Authenticate", "testuser", "wrongpassword").Return("", security.ErrInvalidCredentials)
		dashboard.authManager = mockAuth

		req = httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		dashboard.handleAPILogin(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test protected routes
	t.Run("Protected Routes", func(t *testing.T) {
		// Create a request to a protected API endpoint
		req := httptest.NewRequest("GET", "/api/users", nil)
		w := httptest.NewRecorder()

		// Create a mock auth manager
		mockAuthManager := &mockAuthManager{}
		dashboard.authManager = mockAuthManager

		// Call the users handler
		dashboard.handleUsers(w, req)

		// Should return unauthorized without credentials
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test role-based access control
	t.Run("Role-Based Access Control", func(t *testing.T) {
		// Create a mock auth manager with role middleware
		mockAuth := new(mockAuthManagerExtended)
		
		// Configure the mock to implement role middleware
		mockAuth.On("RoleMiddleware", security.RoleAdmin).Return(
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Check if request has admin role in context
					if strings.Contains(r.Header.Get("Authorization"), "admin-token") {
						next.ServeHTTP(w, r)
					} else {
						w.WriteHeader(http.StatusForbidden)
					}
				})
			},
		)

		dashboard.authManager = mockAuth

		// Create a test handler protected by admin role
		handler := dashboard.authManager.RoleMiddleware(security.RoleAdmin)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Admin access granted"))
			}),
		)

		// Test with admin token
		req := httptest.NewRequest("GET", "/admin/users", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Admin access granted", w.Body.String())

		// Test with non-admin token
		req = httptest.NewRequest("GET", "/admin/users", nil)
		req.Header.Set("Authorization", "Bearer user-token")
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test API key authentication
	t.Run("API Key Authentication", func(t *testing.T) {
		// Create a mock auth manager
		mockAuth := new(mockAuthManagerExtended)
		
		// Configure the mock to implement API key authentication
		mockAuth.On("AuthenticateWithAPIKey", "valid-api-key").Return(&security.User{
			Username: "apiuser",
			Role:     security.RoleUser,
		}, nil)
		mockAuth.On("AuthenticateWithAPIKey", "invalid-api-key").Return(nil, security.ErrInvalidAPIKey)

		dashboard.authManager = mockAuth

		// Create a test handler protected by auth middleware
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check API key header
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Authenticate with API key
			user, err := dashboard.authManager.AuthenticateWithAPIKey(apiKey)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(user.Username))
		})

		// Test with valid API key
		req := httptest.NewRequest("GET", "/api/metrics", nil)
		req.Header.Set("X-API-Key", "valid-api-key")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "apiuser", w.Body.String())

		// Test with invalid API key
		req = httptest.NewRequest("GET", "/api/metrics", nil)
		req.Header.Set("X-API-Key", "invalid-api-key")
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// mockAuthManager is a simple mock for testing
type mockAuthManager struct{}

func (m *mockAuthManager) Authenticate(username, password string) (string, error) {
	return "mock-token", nil
}

func (m *mockAuthManager) ValidateToken(token string) (map[string]interface{}, error) {
	return map[string]interface{}{"username": "mockuser", "role": "user"}, nil
}

func (m *mockAuthManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func (m *mockAuthManager) RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}

// mockAuthManagerExtended is a more comprehensive mock using testify/mock
type mockAuthManagerExtended struct {
	mock.Mock
}

func (m *mockAuthManagerExtended) Authenticate(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *mockAuthManagerExtended) ValidateToken(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *mockAuthManagerExtended) AuthenticateWithAPIKey(apiKey string) (*security.User, error) {
	args := m.Called(apiKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*security.User), args.Error(1)
}

func (m *mockAuthManagerExtended) AuthMiddleware(next http.Handler) http.Handler {
	args := m.Called(next)
	return args.Get(0).(http.Handler)
}

func (m *mockAuthManagerExtended) RoleMiddleware(role string) func(http.Handler) http.Handler {
	args := m.Called(role)
	return args.Get(0).(func(http.Handler) http.Handler)
}
