package security

import (
	"context"
)

// contextKey is a private type for context keys
type contextKey int

const (
	// userContextKey is the key for user context value
	userContextKey contextKey = iota
)

// WithUser returns a new context with the given user
func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext returns the user from the context
func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

// IsAuthenticated checks if the context has an authenticated user
func IsAuthenticated(ctx context.Context) bool {
	_, ok := UserFromContext(ctx)
	return ok
}

// HasRole checks if the authenticated user has the specified role
func HasRole(ctx context.Context, role string) bool {
	user, ok := UserFromContext(ctx)
	if !ok {
		return false
	}
	return user.Role == role
}

// IsAdmin checks if the authenticated user is an admin
func IsAdmin(ctx context.Context) bool {
	return HasRole(ctx, RoleAdmin)
}
