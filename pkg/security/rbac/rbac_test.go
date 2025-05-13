package rbac

import (
	"testing"
)

func TestRBACManager(t *testing.T) {
	// Create a new RBAC manager
	manager := NewRBACManager()

	// Test setting and getting user roles
	t.Run("UserRoles", func(t *testing.T) {
		// Set a user role
		manager.SetUserRole("user1", AdminRole)

		// Get the user role
		role, exists := manager.GetUserRole("user1")
		if !exists {
			t.Errorf("Expected user role to exist, but it didn't")
		}
		if role != AdminRole {
			t.Errorf("Expected user role to be AdminRole, but it was %s", role)
		}

		// Set another user role
		manager.SetUserRole("user2", UserRole)

		// Get the user role
		role, exists = manager.GetUserRole("user2")
		if !exists {
			t.Errorf("Expected user role to exist, but it didn't")
		}
		if role != UserRole {
			t.Errorf("Expected user role to be UserRole, but it was %s", role)
		}

		// Try to get a non-existent user role
		role, exists = manager.GetUserRole("user3")
		if exists {
			t.Errorf("Expected user role to not exist, but it did")
		}
	})

	// Test setting and getting team roles
	t.Run("TeamRoles", func(t *testing.T) {
		// Set a team role
		manager.SetTeamRole("team1", AdminRole)

		// Get the team role
		role, exists := manager.GetTeamRole("team1")
		if !exists {
			t.Errorf("Expected team role to exist, but it didn't")
		}
		if role != AdminRole {
			t.Errorf("Expected team role to be AdminRole, but it was %s", role)
		}

		// Set another team role
		manager.SetTeamRole("team2", UserRole)

		// Get the team role
		role, exists = manager.GetTeamRole("team2")
		if !exists {
			t.Errorf("Expected team role to exist, but it didn't")
		}
		if role != UserRole {
			t.Errorf("Expected team role to be UserRole, but it was %s", role)
		}

		// Try to get a non-existent team role
		role, exists = manager.GetTeamRole("team3")
		if exists {
			t.Errorf("Expected team role to not exist, but it did")
		}
	})

	// Test creating and deleting custom roles
	t.Run("CustomRoles", func(t *testing.T) {
		// Create a custom role
		accessControls := []AccessControl{
			{
				Role:       CustomRole,
				Resource:   TableResource,
				Permission: ReadPermission,
			},
			{
				Role:       CustomRole,
				Resource:   TableResource,
				Permission: WritePermission,
			},
		}
		err := manager.CreateCustomRole("custom1", accessControls)
		if err != nil {
			t.Errorf("Failed to create custom role: %v", err)
		}

		// Try to create a duplicate custom role
		err = manager.CreateCustomRole("custom1", accessControls)
		if err == nil {
			t.Errorf("Expected error when creating duplicate custom role, but got nil")
		}

		// Delete the custom role
		err = manager.DeleteCustomRole("custom1")
		if err != nil {
			t.Errorf("Failed to delete custom role: %v", err)
		}

		// Try to delete a non-existent custom role
		err = manager.DeleteCustomRole("custom1")
		if err == nil {
			t.Errorf("Expected error when deleting non-existent custom role, but got nil")
		}
	})

	// Test checking access
	t.Run("CheckAccess", func(t *testing.T) {
		// Set up some user roles
		manager.SetUserRole("admin", AdminRole)
		manager.SetUserRole("user", UserRole)
		manager.SetUserRole("viewer", ViewerRole)

		// Test admin access
		hasAccess, err := manager.CheckAccess("admin", TableResource, AdminPermission)
		if err != nil {
			t.Errorf("Failed to check admin access: %v", err)
		}
		if !hasAccess {
			t.Errorf("Expected admin to have admin permission on table resource, but they didn't")
		}

		// Test user access
		hasAccess, err = manager.CheckAccess("user", TableResource, WritePermission)
		if err != nil {
			t.Errorf("Failed to check user access: %v", err)
		}
		if !hasAccess {
			t.Errorf("Expected user to have write permission on table resource, but they didn't")
		}

		// Test user access to admin-only resource
		hasAccess, err = manager.CheckAccess("user", UserResource, WritePermission)
		if err != nil {
			t.Errorf("Failed to check user access: %v", err)
		}
		if hasAccess {
			t.Errorf("Expected user to not have write permission on user resource, but they did")
		}

		// Test viewer access
		hasAccess, err = manager.CheckAccess("viewer", TableResource, ReadPermission)
		if err != nil {
			t.Errorf("Failed to check viewer access: %v", err)
		}
		if !hasAccess {
			t.Errorf("Expected viewer to have read permission on table resource, but they didn't")
		}

		// Test viewer access to write permission
		hasAccess, err = manager.CheckAccess("viewer", TableResource, WritePermission)
		if err != nil {
			t.Errorf("Failed to check viewer access: %v", err)
		}
		if hasAccess {
			t.Errorf("Expected viewer to not have write permission on table resource, but they did")
		}

		// Test non-existent user
		_, err = manager.CheckAccess("nonexistent", TableResource, ReadPermission)
		if err == nil {
			t.Errorf("Expected error when checking access for non-existent user, but got nil")
		}
	})

	// Test checking team access
	t.Run("CheckTeamAccess", func(t *testing.T) {
		// Set up some team roles
		manager.SetTeamRole("admin-team", AdminRole)
		manager.SetTeamRole("user-team", UserRole)
		manager.SetTeamRole("viewer-team", ViewerRole)

		// Test admin team access
		hasAccess, err := manager.CheckTeamAccess("admin-team", TableResource, AdminPermission)
		if err != nil {
			t.Errorf("Failed to check admin team access: %v", err)
		}
		if !hasAccess {
			t.Errorf("Expected admin team to have admin permission on table resource, but they didn't")
		}

		// Test user team access
		hasAccess, err = manager.CheckTeamAccess("user-team", TableResource, WritePermission)
		if err != nil {
			t.Errorf("Failed to check user team access: %v", err)
		}
		if !hasAccess {
			t.Errorf("Expected user team to have write permission on table resource, but they didn't")
		}

		// Test viewer team access
		hasAccess, err = manager.CheckTeamAccess("viewer-team", TableResource, ReadPermission)
		if err != nil {
			t.Errorf("Failed to check viewer team access: %v", err)
		}
		if !hasAccess {
			t.Errorf("Expected viewer team to have read permission on table resource, but they didn't")
		}

		// Test non-existent team
		_, err = manager.CheckTeamAccess("nonexistent", TableResource, ReadPermission)
		if err == nil {
			t.Errorf("Expected error when checking access for non-existent team, but got nil")
		}
	})

	// Test listing users and teams
	t.Run("ListUsersAndTeams", func(t *testing.T) {
		// Set up some user and team roles
		manager.SetUserRole("user1", AdminRole)
		manager.SetUserRole("user2", UserRole)
		manager.SetTeamRole("team1", AdminRole)
		manager.SetTeamRole("team2", UserRole)

		// List users
		users := manager.ListUsers()
		if len(users) < 2 {
			t.Errorf("Expected at least 2 users, but got %d", len(users))
		}
		if users["user1"] != AdminRole {
			t.Errorf("Expected user1 to have AdminRole, but got %s", users["user1"])
		}
		if users["user2"] != UserRole {
			t.Errorf("Expected user2 to have UserRole, but got %s", users["user2"])
		}

		// List teams
		teams := manager.ListTeams()
		if len(teams) < 2 {
			t.Errorf("Expected at least 2 teams, but got %d", len(teams))
		}
		if teams["team1"] != AdminRole {
			t.Errorf("Expected team1 to have AdminRole, but got %s", teams["team1"])
		}
		if teams["team2"] != UserRole {
			t.Errorf("Expected team2 to have UserRole, but got %s", teams["team2"])
		}
	})

	// Test listing custom roles
	t.Run("ListCustomRoles", func(t *testing.T) {
		// Create some custom roles
		accessControls1 := []AccessControl{
			{
				Role:       CustomRole,
				Resource:   TableResource,
				Permission: ReadPermission,
			},
		}
		accessControls2 := []AccessControl{
			{
				Role:       CustomRole,
				Resource:   RuleResource,
				Permission: WritePermission,
			},
		}
		manager.CreateCustomRole("custom1", accessControls1)
		manager.CreateCustomRole("custom2", accessControls2)

		// List custom roles
		customRoles := manager.ListCustomRoles()
		if len(customRoles) < 2 {
			t.Errorf("Expected at least 2 custom roles, but got %d", len(customRoles))
		}
		if len(customRoles["custom1"]) != 1 {
			t.Errorf("Expected custom1 to have 1 access control, but got %d", len(customRoles["custom1"]))
		}
		if len(customRoles["custom2"]) != 1 {
			t.Errorf("Expected custom2 to have 1 access control, but got %d", len(customRoles["custom2"]))
		}
	})
}
