package security_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nessi-dev/nessi/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserContext(t *testing.T) {
	// Create a test user
	testUser := security.User{
		Username: "contextuser",
		Email:    "context@example.com",
		Role:     security.RoleAdmin,
	}

	t.Run("Context Management", func(t *testing.T) {
		// Create context with user
		ctx := context.Background()
		ctx = security.WithUser(ctx, testUser)

		// Get user from context
		user, ok := security.UserFromContext(ctx)
		require.True(t, ok)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.Role, user.Role)

		// Get user from empty context should fail
		emptyCtx := context.Background()
		_, ok = security.UserFromContext(emptyCtx)
		assert.False(t, ok)
	})

	t.Run("Request Context", func(t *testing.T) {
		// Create request with user in context
		req, err := http.NewRequest("GET", "/test", nil)
		require.NoError(t, err)

		ctx := req.Context()
		ctx = security.WithUser(ctx, testUser)
		req = req.WithContext(ctx)

		// Get user from request context
		user, ok := security.UserFromContext(req.Context())
		require.True(t, ok)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.Role, user.Role)
	})

	t.Run("Handler Chain", func(t *testing.T) {
		// Create middleware that adds user to context
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := security.WithUser(r.Context(), testUser)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}

		// Create handler that checks for user in context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := security.UserFromContext(r.Context())
			if !ok {
				http.Error(w, "User not found in context", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(user.Username))
		})

		// Create request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Execute middleware and handler
		middleware(handler).ServeHTTP(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "contextuser", w.Body.String())
	})

	t.Run("Role Checking", func(t *testing.T) {
		// Create context with user
		ctx := context.Background()
		ctx = security.WithUser(ctx, testUser)

		// Check roles
		assert.True(t, security.HasRole(ctx, security.RoleAdmin))
		assert.False(t, security.HasRole(ctx, security.RoleUser))

		// Check with empty context
		emptyCtx := context.Background()
		assert.False(t, security.HasRole(emptyCtx, security.RoleAdmin))
	})
}
