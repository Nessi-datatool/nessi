package rbac

// BasicRBACManager implements a simple RBAC system for Nessi

import (
	"sync"
)

// BasicRole represents a simple role in the system
type BasicRole string

const (
	// AdminRole has full access to all resources
	AdminRole BasicRole = "admin"
	
	// UserRole has limited access to resources
	UserRole BasicRole = "user"
	
	// ReadOnlyRole has read-only access to resources
	ReadOnlyRole BasicRole = "readonly"
)

// BasicRBACManager implements a simplified version of the RBACManager interface
type BasicRBACManager struct {
	mu       sync.RWMutex
	userRoles map[string]BasicRole
}

// NewBasicRBACManager creates a new BasicRBACManager
func NewBasicRBACManager() *BasicRBACManager {
	return &BasicRBACManager{
		userRoles: make(map[string]BasicRole),
	}
}

// SetUserRole sets the role for a user
func (m *BasicRBACManager) SetUserRole(username string, role BasicRole) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userRoles[username] = role
}

// GetUserRole gets the role for a user
func (m *BasicRBACManager) GetUserRole(username string) (BasicRole, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	role, exists := m.userRoles[username]
	return role, exists
}

// CheckAccess checks if a user has access to a resource
func (m *BasicRBACManager) CheckAccess(username, resource, action string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	role, exists := m.userRoles[username]
	if !exists {
		return false
	}
	
	// Admin has access to everything
	if role == AdminRole {
		return true
	}
	
	// User has access to most things except admin actions
	if role == UserRole {
		return action != "admin" && action != "delete"
	}
	
	// ReadOnly has access only to read actions
	if role == ReadOnlyRole {
		return action == "read" || action == "view"
	}
	
	return false
}
