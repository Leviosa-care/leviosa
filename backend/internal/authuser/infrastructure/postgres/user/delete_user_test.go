package userRepository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeleteUser TEST_PATH=internal/authuser/infrastructure/postgres/user/delete_user_test.go

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete existing user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "deleteuser@example.com"

		// Create and insert userEncx
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = email
		userEncx.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Verify user exists before deletion
		existingUser, err := td.GetUserEnxByID(t, ctx, userEncx.ID, testPool)
		require.NoError(t, err)
		require.NotNil(t, existingUser)

		// Act
		err = repo.DeleteUser(ctx, userEncx.ID)

		// Assert
		assert.NoError(t, err)

		// Verify user no longer exists
		deletedUser, err := repo.GetUserByID(ctx, userEncx.ID)
		assert.Error(t, err)
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
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should successfully delete user and allow creating user with same email", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "reuseemail@example.com"

		// Create and insert first user
		// firstUser := td.NewTestUser(email, "First", "User")
		firstUser := td.NewTestUserEncx(t)
		firstUser.EmailEncrypted = []byte(email)
		firstUser.EmailHash = email

		err := td.InsertUserEncx(t, ctx, firstUser, testPool)
		require.NoError(t, err)

		// Delete first user
		err = repo.DeleteUser(ctx, firstUser.ID)
		assert.NoError(t, err)

		// Act - create second user with same email
		secondUser := td.NewTestUserEncx(t)
		secondUser.EmailEncrypted = []byte(email)
		secondUser.EmailHash = email
		secondUserFirstnameEncrypted := []byte("Second")
		secondUser.FirstNameEncrypted = secondUserFirstnameEncrypted

		err = td.InsertUserEncx(t, ctx, secondUser, testPool)

		// Assert
		assert.NoError(t, err)

		// Verify second user exists and is different from first
		retrievedUser, err := td.GetUserEnxByID(t, ctx, secondUser.ID, testPool)
		assert.NoError(t, err)

		assert.Equal(t, secondUser.ID, retrievedUser.ID)
		assert.NotEqual(t, firstUser.ID, retrievedUser.ID)
		assert.Equal(t, secondUserFirstnameEncrypted, retrievedUser.FirstNameEncrypted)
	})

	t.Run("should handle deletion of user in different states", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		createUser := func(email, firstname, lastname string, state domain.UserState) *domain.UserEncx {
			user := td.NewTestUserEncx(t)
			user.EmailEncrypted = []byte(email)
			user.EmailHash = email
			user.FirstNameEncrypted = []byte(firstname)
			user.LastNameEncrypted = []byte(lastname)
			user.State = state
			return user
		}

		// Create users in different states
		unverifiedUser := createUser("unverified@example.com", "Unverified", "User", domain.Unverified)
		pendingUser := createUser("pending@example.com", "Pending", "User", domain.Pending)
		activeUser := createUser("active@example.com", "Active", "User", domain.Active)

		// Process and create users
		for _, user := range []*domain.UserEncx{unverifiedUser, pendingUser, activeUser} {
			err := td.InsertUserEncx(t, ctx, user, testPool)
			require.NoError(t, err)
		}

		// Act & Assert - delete each user
		for _, user := range []*domain.UserEncx{unverifiedUser, pendingUser, activeUser} {
			err := repo.DeleteUser(ctx, user.ID)
			assert.NoError(t, err)

			// Verify deletion
			deletedUser, err := repo.GetUserByID(ctx, user.ID)
			assert.Error(t, err)
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
