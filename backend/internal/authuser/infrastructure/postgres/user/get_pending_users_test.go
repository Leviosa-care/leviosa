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

// make test-func TEST_NAME=TestGetPendingUsers TEST_PATH=internal/authuser/infrastructure/postgres/user/get_pending_users_test.go

func TestGetPendingUsers(t *testing.T) {
	ctx := context.Background()

	// Create test users with different states
	createUser := func(email, firstname, lastname string, state domain.UserState) *domain.UserEncx {
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailEncrypted = []byte(email)
		userEncx.EmailHash = email
		userEncx.FirstNameEncrypted = []byte(firstname)
		userEncx.LastNameEncrypted = []byte(lastname)
		userEncx.State = state

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)
		return userEncx
	}

	t.Run("should successfully retrieve pending users ordered by created_at DESC", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user1 := createUser("pending1@example.com", "Alice", "Smith", domain.Pending)
		_ = createUser("active@example.com", "Bob", "Jones", domain.Active)
		user3 := createUser("pending2@example.com", "Carol", "Jones", domain.Pending)

		_ = createUser("unverified@example.com", "Dave", "Wilson", domain.Unverified)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, pendingUsers, 2, "Should return exactly 2 pending users")

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
		_ = createUser("active@example.com", "Active", "User", domain.Active)
		_ = createUser("unverified@example.com", "Unverified", "User", domain.Unverified)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pendingUsers, "Should return non-nil slice")
		assert.Empty(t, pendingUsers, "Should return empty slice when no pending users")
	})

	t.Run("should return empty slice when no users exist at all", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, pendingUsers, "Should return non-nil slice")
		assert.Empty(t, pendingUsers, "Should return empty slice when no users exist")
	})

	t.Run("should handle users with and without telephone hash correctly", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// TODO: I get some not null violation for postgres

		createUserWithTelephone := func(email, firstname, lastname, telephone string, state domain.UserState) *domain.UserEncx {
			userEncx := td.NewTestUserEncx(t)
			userEncx.EmailHash = email
			userEncx.EmailEncrypted = []byte(email)
			userEncx.FirstNameEncrypted = []byte(firstname)
			userEncx.LastNameEncrypted = []byte(lastname)
			userEncx.TelephoneHash = telephone
			userEncx.TelephoneEncrypted = []byte(telephone)
			userEncx.State = state

			err := td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)

			return userEncx
		}

		// User with telephone
		userWithPhone := createUserWithTelephone("withphone@example.com", "With", "Phone", "+33123456789", domain.Pending)
		// User without telephone
		userWithoutPhone := createUserWithTelephone("nophone@example.com", "No", "Phone", "", domain.Pending)

		// Act
		pendingUsers, err := repo.GetPendingUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, pendingUsers, 2)

		// Find users in results
		var withPhoneUser, withoutPhoneUser *domain.UserEncx
		for _, user := range pendingUsers {
			if user.EmailHash == userWithPhone.EmailHash {
				println("with phone")
				withPhoneUser = user
			} else if user.EmailHash == userWithoutPhone.EmailHash {
				println("without phone")
				withoutPhoneUser = user
			}
		}

		assert.NotNil(t, withPhoneUser, "User with phone should be found")
		assert.NotNil(t, withoutPhoneUser, "User without phone should be found")

		// Verify telephone hash handling
		assert.NotZero(t, withPhoneUser.TelephoneHash, "User with phone should have telephone hash")
		assert.NotZero(t, withPhoneUser.TelephoneEncrypted, "User with phone should have encrypted telephone")
	})

	t.Run("should handle large number of pending users", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const numUsers = 50
		expectedUsers := make([]*domain.UserEncx, numUsers)

		// Create many pending users
		for i := 0; i < numUsers; i++ {
			user := createUser(
				fmt.Sprintf("pending%d@example.com", i),
				fmt.Sprintf("User%d", i),
				"Test",
				domain.Pending,
			)
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
