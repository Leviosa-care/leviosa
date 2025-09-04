package userRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetAllUsers make test-unit-user-test

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve all users ordered by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create test users with different states
		user1 := td.NewTestUser("pending@example.com", "Alice", "Smith")
		user1.State = domain.Pending
		err := crypto.ProcessStruct(ctx, user1)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user1, testPool)

		user2 := td.NewTestUser("active@example.com", "Bob", "Jones")
		user2.State = domain.Active
		err = crypto.ProcessStruct(ctx, user2)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user2, testPool)

		user3 := td.NewTestUser("unverified@example.com", "Carol", "Brown")
		user3.State = domain.Unverified
		err = crypto.ProcessStruct(ctx, user3)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user3, testPool)

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.Len(t, allUsers, 3, "Should return all users regardless of state")

		// Verify users are ordered by created_at DESC (newest first)
		// Since user3 was inserted last, it should come first
		assert.Equal(t, user3.EmailHash, allUsers[0].EmailHash, "First user should be the most recently created")
		assert.Equal(t, user2.EmailHash, allUsers[1].EmailHash, "Second user should be the second created")
		assert.Equal(t, user1.EmailHash, allUsers[2].EmailHash, "Third user should be the first created")

		// Verify all user states are preserved
		userStateMap := make(map[string]domain.UserState)
		for _, user := range allUsers {
			userStateMap[user.EmailHash] = user.State
		}
		assert.Equal(t, domain.Pending, userStateMap[user1.EmailHash])
		assert.Equal(t, domain.Active, userStateMap[user2.EmailHash])
		assert.Equal(t, domain.Unverified, userStateMap[user3.EmailHash])

		// Verify encrypted fields are populated (not decrypted at repository layer)
		for _, user := range allUsers {
			assert.NotEmpty(t, user.EmailEncrypted)
			assert.NotEmpty(t, user.FirstNameEncrypted)
			assert.NotEmpty(t, user.LastNameEncrypted)
			assert.NotEmpty(t, user.DEKEncrypted)
			assert.Greater(t, user.KeyVersion, 0)
		}
	})

	t.Run("should return empty slice when no users exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, allUsers, "Should return non-nil slice")
		assert.Empty(t, allUsers, "Should return empty slice when no users exist")
	})

	t.Run("should handle users with and without telephone hash correctly", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// User with telephone
		userWithPhone := td.NewTestUser("withphone@example.com", "With", "Phone")
		userWithPhone.State = domain.Active
		userWithPhone.Telephone = "+33123456789"
		err := crypto.ProcessStruct(ctx, userWithPhone)
		require.NoError(t, err)
		td.InsertUser(t, ctx, userWithPhone, testPool)

		// User without telephone
		userWithoutPhone := td.NewTestUser("nophone@example.com", "No", "Phone")
		userWithoutPhone.State = domain.Pending
		userWithoutPhone.Telephone = ""
		err = crypto.ProcessStruct(ctx, userWithoutPhone)
		require.NoError(t, err)
		td.InsertUser(t, ctx, userWithoutPhone, testPool)

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.Len(t, allUsers, 2)

		// Find users in results
		var withPhoneUser, withoutPhoneUser *domain.User
		for _, user := range allUsers {
			if user.EmailHash == userWithPhone.EmailHash {
				withPhoneUser = user
			} else if user.EmailHash == userWithoutPhone.EmailHash {
				withoutPhoneUser = user
			}
		}

		require.NotNil(t, withPhoneUser, "User with phone should be found")
		require.NotNil(t, withoutPhoneUser, "User without phone should be found")

		// Verify telephone hash handling
		assert.NotEmpty(t, withPhoneUser.TelephoneHash, "User with phone should have telephone hash")
		assert.NotEmpty(t, withPhoneUser.TelephoneEncrypted, "User with phone should have encrypted telephone")

		// User without phone should have empty hash (empty string hashes to a specific value)
		assert.Equal(t,
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			withoutPhoneUser.TelephoneHash,
			"User without phone should have empty string hash")
	})

	t.Run("should return users with all possible states", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create users with all possible states
		states := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}
		expectedUsers := make(map[domain.UserState]*domain.User)

		for i, state := range states {
			user := td.NewTestUser(
				fmt.Sprintf("user%d@example.com", i),
				fmt.Sprintf("User%d", i),
				"Test",
			)
			user.State = state
			err := crypto.ProcessStruct(ctx, user)
			require.NoError(t, err)
			td.InsertUser(t, ctx, user, testPool)
			expectedUsers[state] = user
		}

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, allUsers, len(states), "Should return users with all states")

		// Verify all states are represented
		foundStates := make(map[domain.UserState]bool)
		for _, user := range allUsers {
			foundStates[user.State] = true
		}

		for _, expectedState := range states {
			assert.True(t, foundStates[expectedState], "Should find user with state: %v", expectedState)
		}
	})

	t.Run("should handle large number of users", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const numUsers = 100
		expectedUsers := make([]*domain.User, numUsers)
		states := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}

		// Create many users with various states
		for i := 0; i < numUsers; i++ {
			user := td.NewTestUser(
				fmt.Sprintf("user%d@example.com", i),
				fmt.Sprintf("User%d", i),
				"Test",
			)
			user.State = states[i%len(states)] // Cycle through states
			err := crypto.ProcessStruct(ctx, user)
			require.NoError(t, err)
			td.InsertUser(t, ctx, user, testPool)
			expectedUsers[i] = user
		}

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, allUsers, numUsers, "Should return all users")

		// Verify order (newest first - reverse order of insertion)
		assert.Equal(t, expectedUsers[numUsers-1].EmailHash, allUsers[0].EmailHash, "First should be last inserted")
		assert.Equal(t, expectedUsers[0].EmailHash, allUsers[numUsers-1].EmailHash, "Last should be first inserted")

		// Verify state distribution
		stateCount := make(map[domain.UserState]int)
		for _, user := range allUsers {
			stateCount[user.State]++
		}

		// Each state should appear roughly numUsers/len(states) times
		expectedCount := numUsers / len(states)
		for _, state := range states {
			assert.GreaterOrEqual(t, stateCount[state], expectedCount-1, "State %v should appear at least %d times", state, expectedCount-1)
			assert.LessOrEqual(t, stateCount[state], expectedCount+1, "State %v should appear at most %d times", state, expectedCount+1)
		}
	})

	t.Run("should handle database connection errors gracefully", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For comprehensive testing, we'd need to simulate connection failures
		t.Skip("Database connection error testing requires mocking or network disruption")
	})
}

