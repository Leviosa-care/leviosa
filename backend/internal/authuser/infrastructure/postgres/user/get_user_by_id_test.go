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

// make test-func TEST_NAME=TestGetUserByID TEST_PATH=internal/authuser/infrastructure/postgres/user/get_user_by_id_test.go

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve user by ID", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "getuser@example.com"

		// Get the expected user to retrieve its ID
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailEncrypted = []byte(email)
		userEncx.EmailHash = email

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByID(ctx, userEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)
		assert.Equal(t, userEncx.ID, retrievedUserEncx.ID)
		assert.Equal(t, userEncx.EmailHash, retrievedUserEncx.EmailHash)
		assert.Equal(t, domain.Unverified, retrievedUserEncx.State) // Default state from helpers

		// Verify encrypted fields are populated
		assert.NotEmpty(t, retrievedUserEncx.EmailEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.FirstNameEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.LastNameEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.DEKEncrypted)
		assert.Greater(t, retrievedUserEncx.KeyVersion, 0)
	})

	t.Run("should return not found error when user does not exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		nonExistentID := uuid.New()

		// Act
		user, err := repo.GetUserByID(ctx, nonExistentID)

		// Assert
		assert.Error(t, err)
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

		// Create userEncx with phone
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = email
		userEncx.EmailEncrypted = []byte(email)
		userEncx.TelephoneHash = "+33123456789"
		userEncx.TelephoneEncrypted = []byte("+33123456789")

		// Insert user manually to ensure phone is included
		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByID(ctx, userEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)
		assert.Equal(t, userEncx.ID, retrievedUserEncx.ID)
		assert.NotEmpty(t, retrievedUserEncx.TelephoneHash)
		assert.NotEmpty(t, retrievedUserEncx.TelephoneEncrypted)
	})

	t.Run("should handle user without telephone hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		email := "nophone@example.com"

		// Create userEncx without phone
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = email
		userEncx.EmailEncrypted = []byte(email)
		userEncx.TelephoneEncrypted = []byte("") // No telephone

		// Insert user manually
		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByID(ctx, userEncx.ID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)
		assert.Equal(t, userEncx.ID, retrievedUserEncx.ID)
	})
}
