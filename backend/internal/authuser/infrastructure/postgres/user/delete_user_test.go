package userRepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestDeleteUser make test-unit-user-test

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete existing user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "deleteuser@example.com"

		// Create and insert user
		user := td.NewTestUser(email, "Delete", "User")
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Verify user exists before deletion
		existingUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, existingUser)

		// Act
		err = repo.DeleteUser(ctx, user.ID)

		// Assert
		require.NoError(t, err)

		// Verify user no longer exists
		deletedUser, err := repo.GetUserByID(ctx, user.ID)
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
		assert.Nil(t, deletedUser)
	})

	t.Run("should return not found error when deleting non-existent user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		nonExistentID := uuid.New()

		// Act
		err := repo.DeleteUser(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should successfully delete user and allow creating user with same email", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "reuseemail@example.com"

		// Create and insert first user
		firstUser := td.NewTestUser(email, "First", "User")
		err := crypto.ProcessStruct(ctx, firstUser)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, firstUser)
		require.NoError(t, err)

		// Delete first user
		err = repo.DeleteUser(ctx, firstUser.ID)
		require.NoError(t, err)

		// Act - create second user with same email
		secondUser := td.NewTestUser(email, "Second", "User")
		err = crypto.ProcessStruct(ctx, secondUser)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, secondUser)

		// Assert
		require.NoError(t, err)

		// Verify second user exists and is different from first
		retrievedUser, err := repo.GetUserByID(ctx, secondUser.ID)
		require.NoError(t, err)
		err = crypto.DecryptStruct(ctx, retrievedUser)
		require.NoError(t, err)

		assert.Equal(t, secondUser.ID, retrievedUser.ID)
		assert.NotEqual(t, firstUser.ID, retrievedUser.ID)
		assert.Equal(t, "Second", retrievedUser.FirstName)
	})

	t.Run("should handle deletion of user in different states", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create users in different states
		unverifiedUser := td.NewTestUser("unverified@example.com", "Unverified", "User")
		pendingUser := td.NewTestUser("pending@example.com", "Pending", "User")
		activeUser := td.NewTestUser("active@example.com", "Active", "User")

		// Set different states
		unverifiedUser.State = "unverified"
		pendingUser.State = "pending"
		activeUser.State = "active"

		// Process and create users
		for _, user := range []*domain.User{unverifiedUser, pendingUser, activeUser} {
			err := crypto.ProcessStruct(ctx, user)
			require.NoError(t, err)
			err = repo.CreateUser(ctx, user)
			require.NoError(t, err)
		}

		// Act & Assert - delete each user
		for _, user := range []*domain.User{unverifiedUser, pendingUser, activeUser} {
			err := repo.DeleteUser(ctx, user.ID)
			require.NoError(t, err)

			// Verify deletion
			deletedUser, err := repo.GetUserByID(ctx, user.ID)
			require.Error(t, err)
			assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
			assert.Nil(t, deletedUser)
		}
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})
}

