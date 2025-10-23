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

// TEST=TestGetPendingUsers make test-unit-user-test

func TestGetPendingUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve pending users ordered by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create test users with different states
		// Note: helpers functions currently return wrong type, but we'll work around it for now
		user1 := td.NewTestUser("pending1@example.com", "Alice", "Smith")
		user1.State = domain.Pending // Use domain.Pending instead of auth.Pending
		err := crypto.ProcessStruct(ctx, user1)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user1, testPool)

		user2 := td.NewTestUser("active@example.com", "Bob", "Jones")
		user2.State = domain.Active // Use domain.Active instead of auth.Active
		err = crypto.ProcessStruct(ctx, user2)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user2, testPool)

		user3 := td.NewTestUser("pending2@example.com", "Carol", "Brown")
		user3.State = domain.Pending // Use domain.Pending instead of auth.Pending
		err = crypto.ProcessStruct(ctx, user3)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user3, testPool)

		user4 := td.NewTestUser("unverified@example.com", "Dave", "Wilson")
		user4.State = domain.Unverified // Use domain.Unverified instead of auth.Unverified
		err = crypto.ProcessStruct(ctx, user4)
		require.NoError(t, err)
		td.InsertUser(t, ctx, user4, testPool)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.Len(t, pendingUsers, 2, "Should return exactly 2 pending users")

		// Verify users are ordered by created_at DESC (newest first)
		// Since user3 was inserted after user1, it should come first
		assert.Equal(t, user3.EmailHash, pendingUsers[0].EmailHash, "First user should be the most recently created")
		assert.Equal(t, user1.EmailHash, pendingUsers[1].EmailHash, "Second user should be the first created")

		// Verify all users have pending state
		for _, user := range pendingUsers {
			assert.Equal(t, domain.Pending, user.State)
		}

		// Verify encrypted fields are populated (not decrypted at repository layer)
		for _, user := range pendingUsers {
			assert.NotEmpty(t, user.EmailEncrypted)
			assert.NotEmpty(t, user.FirstNameEncrypted)
			assert.NotEmpty(t, user.LastNameEncrypted)
			assert.NotEmpty(t, user.DEKEncrypted)
			assert.Greater(t, user.KeyVersion, 0)
		}
	})

	t.Run("should return empty slice when no pending users exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create users with non-pending states
		activeUser := td.NewTestUser("active@example.com", "Active", "User")
		activeUser.State = domain.Active
		err := crypto.ProcessStruct(ctx, activeUser)
		require.NoError(t, err)
		td.InsertUser(t, ctx, activeUser, testPool)

		unverifiedUser := td.NewTestUser("unverified@example.com", "Unverified", "User")
		unverifiedUser.State = domain.Unverified
		err = crypto.ProcessStruct(ctx, unverifiedUser)
		require.NoError(t, err)
		td.InsertUser(t, ctx, unverifiedUser, testPool)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, pendingUsers, "Should return non-nil slice")
		assert.Empty(t, pendingUsers, "Should return empty slice when no pending users")
	})

	t.Run("should return empty slice when no users exist at all", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, pendingUsers, "Should return non-nil slice")
		assert.Empty(t, pendingUsers, "Should return empty slice when no users exist")
	})

	t.Run("should handle users with and without telephone hash correctly", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// User with telephone
		userWithPhone := td.NewTestUser("withphone@example.com", "With", "Phone")
		userWithPhone.State = domain.Pending
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
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		require.NoError(t, err)
		require.Len(t, pendingUsers, 2)

		// Find users in results
		var withPhoneUser, withoutPhoneUser *domain.User
		for _, user := range pendingUsers {
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

	t.Run("should handle large number of pending users", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const numUsers = 50
		expectedUsers := make([]*domain.User, numUsers)

		// Create many pending users
		for i := 0; i < numUsers; i++ {
			user := td.NewTestUser(
				fmt.Sprintf("pending%d@example.com", i),
				fmt.Sprintf("User%d", i),
				"Test",
			)
			user.State = domain.Pending
			err := crypto.ProcessStruct(ctx, user)
			require.NoError(t, err)
			td.InsertUser(t, ctx, user, testPool)
			expectedUsers[i] = user
		}

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		require.NoError(t, err)
		assert.Len(t, pendingUsers, numUsers, "Should return all pending users")

		// Verify all have pending state
		for _, user := range pendingUsers {
			assert.Equal(t, domain.Pending, user.State)
		}

		// Verify order (newest first - reverse order of insertion)
		assert.Equal(t, expectedUsers[numUsers-1].EmailHash, pendingUsers[0].EmailHash, "First should be last inserted")
		assert.Equal(t, expectedUsers[0].EmailHash, pendingUsers[numUsers-1].EmailHash, "Last should be first inserted")
	})

	t.Run("should handle database connection errors gracefully", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For comprehensive testing, we'd need to simulate connection failures
		t.Skip("Database connection error testing requires mocking or network disruption")
	})
}
