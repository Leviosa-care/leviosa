package userRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllUsers TEST_PATH=internal/authuser/infrastructure/postgres/user/get_all_users_test.go

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve all users ordered by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		createUser := func(email, firstname, lastname string, state domain.UserState) *domain.UserEncx {
			user := td.NewTestUserEncx(t)
			user.EmailEncrypted = []byte(email)
			user.EmailHash = email
			user.FirstNameEncrypted = []byte(firstname)
			user.LastNameEncrypted = []byte(lastname)
			user.State = state

			err := td.InsertUserEncx(t, ctx, user, testPool)
			require.NoError(t, err)
			return user
		}

		// Create test users with different states
		user1 := createUser("pending@example.com", "Alice", "SMITH", domain.Pending)
		user2 := createUser("active@example.com", "Bob", "JONES", domain.Active)
		user3 := createUser("unverified@example.com", "Carol", "Brown", domain.Unverified)

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allUsers, 3, "Should return all users regardless of state")

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
		assert.NoError(t, err)
		assert.NotNil(t, allUsers, "Should return non-nil slice")
		assert.Empty(t, allUsers, "Should return empty slice when no users exist")
	})

	t.Run("should handle users with and without telephone encrypted correctly", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// User with telephone
		// userWithPhone := td.NewTestUser("withphone@example.com", "With", "Phone")
		userWithPhone := td.NewTestUserEncx(t)
		userWithPhone.EmailEncrypted = []byte("withphone@example.com")
		userWithPhone.EmailHash = "withphone@example.com"
		userWithPhone.FirstNameEncrypted = []byte("With")
		userWithPhone.LastNameEncrypted = []byte("Phone")
		userWithPhone.State = domain.Active
		userWithPhone.TelephoneHash = "+33123456789"
		userWithPhone.TelephoneEncrypted = []byte("+33123456789")

		err := td.InsertUserEncx(t, ctx, userWithPhone, testPool)
		require.NoError(t, err)

		// User without telephone
		// userWithoutPhone := td.NewTestUser("nophone@example.com", "No", "Phone")
		userWithoutPhone := td.NewTestUserEncx(t)
		userWithoutPhone.State = domain.Pending
		userWithoutPhone.EmailEncrypted = []byte("nophone@example.com")
		userWithoutPhone.EmailHash = "nophone@example.com"
		userWithoutPhone.FirstNameEncrypted = []byte("No")
		userWithoutPhone.LastNameEncrypted = []byte("Phone")
		userWithoutPhone.TelephoneHash = ""
		userWithoutPhone.TelephoneEncrypted = []byte("")

		err = td.InsertUserEncx(t, ctx, userWithoutPhone, testPool)
		require.NoError(t, err)

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allUsers, 2)

		// Find users in results
		var withPhoneUser, withoutPhoneUser *domain.UserEncx
		for _, user := range allUsers {
			if user.EmailHash == userWithPhone.EmailHash {
				withPhoneUser = user
			} else if user.EmailHash == userWithoutPhone.EmailHash {
				withoutPhoneUser = user
			}
		}

		assert.NotNil(t, withPhoneUser, "User with phone should be found")
		assert.NotNil(t, withoutPhoneUser, "User without phone should be found")

		// Verify telephone hash handling
		assert.NotZero(t, withPhoneUser.TelephoneHash, "User with phone should have telephone hash")
		assert.NotZero(t, withPhoneUser.TelephoneEncrypted, "User with phone should have encrypted telephone")
	})

	t.Run("should return users with all possible states", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create users with all possible states
		states := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}
		expectedUsers := make(map[domain.UserState]*domain.UserEncx)

		for i, state := range states {
			userEncx := td.NewTestUserEncx(t)
			userEncx.EmailHash = fmt.Sprintf("user%d@example.com", i)
			userEncx.EmailEncrypted = []byte(fmt.Sprintf("user%d@example.com", i))
			userEncx.FirstNameEncrypted = []byte(fmt.Sprintf("User%d", i))
			userEncx.LastNameEncrypted = []byte("Test")
			userEncx.State = state

			err := td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)

			expectedUsers[state] = userEncx
		}

		// Act
		allUsers, err := repo.GetAllUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, allUsers, len(states), "Should return users with all states")

		// Verify all states are represented
		foundStates := make(map[domain.UserState]bool)
		for _, userEncx := range allUsers {
			foundStates[userEncx.State] = true
		}

		for _, expectedState := range states {
			assert.True(t, foundStates[expectedState], "Should find user with state: %v", expectedState)
		}
	})

	t.Run("should handle large number of users", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const numUsers = 100
		expectedUsers := make([]*domain.UserEncx, numUsers)
		states := []domain.UserState{domain.Pending, domain.Active, domain.Unverified}

		// Create many users with various states
		for i := 0; i < numUsers; i++ {
			userEncx := td.NewTestUserEncx(t)
			userEncx.EmailHash = fmt.Sprintf("user%d@example.com", i)
			userEncx.EmailEncrypted = []byte(fmt.Sprintf("user%d@example.com", i))
			userEncx.FirstNameEncrypted = []byte(fmt.Sprintf("User%d", i))
			userEncx.LastNameEncrypted = []byte("Test")
			userEncx.State = states[i%len(states)] // Cycle through states

			err := td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)

			expectedUsers[i] = userEncx
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
