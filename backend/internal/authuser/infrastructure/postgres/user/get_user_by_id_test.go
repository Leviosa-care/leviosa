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

// TEST=TestGetUserByID make test-unit-user-test

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve user by ID", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "getuser@example.com"
		td.InsertTestUser(t, ctx, email, "John", "Doe", testPool, crypto)

		// Get the expected user to retrieve its ID
		expectedUser := td.NewTestUser(email, "John", "Doe")
		err := crypto.ProcessStruct(ctx, expectedUser)
		require.NoError(t, err)

		// Get the actual user from DB to get the real ID (since we insert with helpers)
		userByEmail, err := repo.GetUserByEmailHash(ctx, expectedUser.EmailHash)
		require.NoError(t, err)
		actualUserID := userByEmail.ID

		// Act
		retrievedUser, err := repo.GetUserByID(ctx, actualUserID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)
		assert.Equal(t, actualUserID, retrievedUser.ID)
		assert.Equal(t, expectedUser.EmailHash, retrievedUser.EmailHash)
		assert.Equal(t, domain.Unverified, retrievedUser.State) // Default state from helpers

		// Verify encrypted fields are populated
		assert.NotEmpty(t, retrievedUser.EmailEncrypted)
		assert.NotEmpty(t, retrievedUser.FirstNameEncrypted)
		assert.NotEmpty(t, retrievedUser.LastNameEncrypted)
		assert.NotEmpty(t, retrievedUser.DEKEncrypted)
		assert.Greater(t, retrievedUser.KeyVersion, 0)
	})

	t.Run("should return not found error when user does not exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		nonExistentID := uuid.New()

		// Act
		user, err := repo.GetUserByID(ctx, nonExistentID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})

	t.Run("should correctly handle user with telephone hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "userphone@example.com"

		// Create user with phone
		user := td.NewTestUser(email, "Jane", "Smith")
		user.Telephone = "+33123456789"
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Insert user manually to ensure phone is included
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Act
		retrievedUser, err := repo.GetUserByID(ctx, user.ID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.NotEmpty(t, retrievedUser.TelephoneHash)
		assert.NotEmpty(t, retrievedUser.TelephoneEncrypted)
	})

	t.Run("should handle user without telephone hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "nophone@example.com"

		// Create user without phone
		user := td.NewTestUser(email, "Bob", "Jones")
		user.Telephone = "" // No telephone
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Insert user manually
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Act
		retrievedUser, err := repo.GetUserByID(ctx, user.ID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t,
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			retrievedUser.TelephoneHash,
		)
	})
}
