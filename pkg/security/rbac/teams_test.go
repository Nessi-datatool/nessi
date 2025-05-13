package rbac

import (
	"testing"
)

func TestTeamManager(t *testing.T) {
	// Create a new team manager
	manager := NewTeamManager()

	// Test creating and getting teams
	t.Run("CreateAndGetTeam", func(t *testing.T) {
		// Create a team
		team, err := manager.CreateTeam("team1", "Team 1", "Test team 1", AdminRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}
		if team.ID != "team1" {
			t.Errorf("Expected team ID to be 'team1', but it was %s", team.ID)
		}
		if team.Name != "Team 1" {
			t.Errorf("Expected team name to be 'Team 1', but it was %s", team.Name)
		}
		if team.Description != "Test team 1" {
			t.Errorf("Expected team description to be 'Test team 1', but it was %s", team.Description)
		}
		if team.Role != AdminRole {
			t.Errorf("Expected team role to be AdminRole, but it was %s", team.Role)
		}
		if len(team.Members) != 0 {
			t.Errorf("Expected team to have 0 members, but it had %d", len(team.Members))
		}

		// Try to create a duplicate team
		_, err = manager.CreateTeam("team1", "Duplicate Team", "This should fail", UserRole)
		if err == nil {
			t.Errorf("Expected error when creating duplicate team, but got nil")
		}

		// Get the team
		team, err = manager.GetTeam("team1")
		if err != nil {
			t.Errorf("Failed to get team: %v", err)
		}
		if team.ID != "team1" {
			t.Errorf("Expected team ID to be 'team1', but it was %s", team.ID)
		}

		// Try to get a non-existent team
		_, err = manager.GetTeam("nonexistent")
		if err == nil {
			t.Errorf("Expected error when getting non-existent team, but got nil")
		}
	})

	// Test updating teams
	t.Run("UpdateTeam", func(t *testing.T) {
		// Create a team
		_, err := manager.CreateTeam("team2", "Team 2", "Test team 2", UserRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}

		// Update the team
		team, err := manager.UpdateTeam("team2", "Updated Team 2", "Updated description", ViewerRole)
		if err != nil {
			t.Errorf("Failed to update team: %v", err)
		}
		if team.Name != "Updated Team 2" {
			t.Errorf("Expected team name to be 'Updated Team 2', but it was %s", team.Name)
		}
		if team.Description != "Updated description" {
			t.Errorf("Expected team description to be 'Updated description', but it was %s", team.Description)
		}
		if team.Role != ViewerRole {
			t.Errorf("Expected team role to be ViewerRole, but it was %s", team.Role)
		}

		// Try to update a non-existent team
		_, err = manager.UpdateTeam("nonexistent", "Nonexistent", "This should fail", AdminRole)
		if err == nil {
			t.Errorf("Expected error when updating non-existent team, but got nil")
		}
	})

	// Test deleting teams
	t.Run("DeleteTeam", func(t *testing.T) {
		// Create a team
		_, err := manager.CreateTeam("team3", "Team 3", "Test team 3", UserRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}

		// Delete the team
		err = manager.DeleteTeam("team3")
		if err != nil {
			t.Errorf("Failed to delete team: %v", err)
		}

		// Try to get the deleted team
		_, err = manager.GetTeam("team3")
		if err == nil {
			t.Errorf("Expected error when getting deleted team, but got nil")
		}

		// Try to delete a non-existent team
		err = manager.DeleteTeam("nonexistent")
		if err == nil {
			t.Errorf("Expected error when deleting non-existent team, but got nil")
		}
	})

	// Test adding and removing members
	t.Run("AddAndRemoveMembers", func(t *testing.T) {
		// Create a team
		_, err := manager.CreateTeam("team4", "Team 4", "Test team 4", UserRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}

		// Add a member
		err = manager.AddMember("team4", "user1")
		if err != nil {
			t.Errorf("Failed to add member: %v", err)
		}

		// Get the team and check the member was added
		team, err := manager.GetTeam("team4")
		if err != nil {
			t.Errorf("Failed to get team: %v", err)
		}
		if len(team.Members) != 1 {
			t.Errorf("Expected team to have 1 member, but it had %d", len(team.Members))
		}
		if team.Members[0] != "user1" {
			t.Errorf("Expected team member to be 'user1', but it was %s", team.Members[0])
		}

		// Add another member
		err = manager.AddMember("team4", "user2")
		if err != nil {
			t.Errorf("Failed to add member: %v", err)
		}

		// Get the team and check the member was added
		team, err = manager.GetTeam("team4")
		if err != nil {
			t.Errorf("Failed to get team: %v", err)
		}
		if len(team.Members) != 2 {
			t.Errorf("Expected team to have 2 members, but it had %d", len(team.Members))
		}

		// Try to add a duplicate member
		err = manager.AddMember("team4", "user1")
		if err == nil {
			t.Errorf("Expected error when adding duplicate member, but got nil")
		}

		// Try to add a member to a non-existent team
		err = manager.AddMember("nonexistent", "user3")
		if err == nil {
			t.Errorf("Expected error when adding member to non-existent team, but got nil")
		}

		// Remove a member
		err = manager.RemoveMember("team4", "user1")
		if err != nil {
			t.Errorf("Failed to remove member: %v", err)
		}

		// Get the team and check the member was removed
		team, err = manager.GetTeam("team4")
		if err != nil {
			t.Errorf("Failed to get team: %v", err)
		}
		if len(team.Members) != 1 {
			t.Errorf("Expected team to have 1 member, but it had %d", len(team.Members))
		}
		if team.Members[0] != "user2" {
			t.Errorf("Expected team member to be 'user2', but it was %s", team.Members[0])
		}

		// Try to remove a non-existent member
		err = manager.RemoveMember("team4", "nonexistent")
		if err == nil {
			t.Errorf("Expected error when removing non-existent member, but got nil")
		}

		// Try to remove a member from a non-existent team
		err = manager.RemoveMember("nonexistent", "user2")
		if err == nil {
			t.Errorf("Expected error when removing member from non-existent team, but got nil")
		}
	})

	// Test getting user teams
	t.Run("GetUserTeams", func(t *testing.T) {
		// Create some teams
		_, err := manager.CreateTeam("team5", "Team 5", "Test team 5", UserRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}
		_, err = manager.CreateTeam("team6", "Team 6", "Test team 6", ViewerRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}

		// Add a user to both teams
		err = manager.AddMember("team5", "user3")
		if err != nil {
			t.Errorf("Failed to add member: %v", err)
		}
		err = manager.AddMember("team6", "user3")
		if err != nil {
			t.Errorf("Failed to add member: %v", err)
		}

		// Get the user's teams
		teams := manager.GetUserTeams("user3")
		if len(teams) != 2 {
			t.Errorf("Expected user to be in 2 teams, but they were in %d", len(teams))
		}

		// Get teams for a user who isn't in any teams
		teams = manager.GetUserTeams("nonexistent")
		if len(teams) != 0 {
			t.Errorf("Expected user to be in 0 teams, but they were in %d", len(teams))
		}
	})

	// Test listing teams
	t.Run("ListTeams", func(t *testing.T) {
		// List all teams
		teams := manager.ListTeams()
		if len(teams) < 5 {
			t.Errorf("Expected at least 5 teams, but got %d", len(teams))
		}
	})

	// Test checking membership
	t.Run("IsMember", func(t *testing.T) {
		// Create a team
		_, err := manager.CreateTeam("team7", "Team 7", "Test team 7", UserRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}

		// Add a member
		err = manager.AddMember("team7", "user4")
		if err != nil {
			t.Errorf("Failed to add member: %v", err)
		}

		// Check if the user is a member
		isMember, err := manager.IsMember("team7", "user4")
		if err != nil {
			t.Errorf("Failed to check membership: %v", err)
		}
		if !isMember {
			t.Errorf("Expected user to be a member, but they weren't")
		}

		// Check if a non-member is a member
		isMember, err = manager.IsMember("team7", "nonexistent")
		if err != nil {
			t.Errorf("Failed to check membership: %v", err)
		}
		if isMember {
			t.Errorf("Expected user to not be a member, but they were")
		}

		// Try to check membership in a non-existent team
		_, err = manager.IsMember("nonexistent", "user4")
		if err == nil {
			t.Errorf("Expected error when checking membership in non-existent team, but got nil")
		}
	})

	// Test getting and setting team roles
	t.Run("TeamRoles", func(t *testing.T) {
		// Create a team
		_, err := manager.CreateTeam("team8", "Team 8", "Test team 8", UserRole)
		if err != nil {
			t.Errorf("Failed to create team: %v", err)
		}

		// Get the team role
		role, err := manager.GetTeamRole("team8")
		if err != nil {
			t.Errorf("Failed to get team role: %v", err)
		}
		if role != UserRole {
			t.Errorf("Expected team role to be UserRole, but it was %s", role)
		}

		// Set the team role
		err = manager.SetTeamRole("team8", AdminRole)
		if err != nil {
			t.Errorf("Failed to set team role: %v", err)
		}

		// Get the team role again
		role, err = manager.GetTeamRole("team8")
		if err != nil {
			t.Errorf("Failed to get team role: %v", err)
		}
		if role != AdminRole {
			t.Errorf("Expected team role to be AdminRole, but it was %s", role)
		}

		// Try to get the role of a non-existent team
		_, err = manager.GetTeamRole("nonexistent")
		if err == nil {
			t.Errorf("Expected error when getting role of non-existent team, but got nil")
		}

		// Try to set the role of a non-existent team
		err = manager.SetTeamRole("nonexistent", ViewerRole)
		if err == nil {
			t.Errorf("Expected error when setting role of non-existent team, but got nil")
		}
	})
}
