package userRepository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestUpdateUser make test-unit-user-test

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully update user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "updateuser@example.com"

		// Create and insert initial user
		user := td.NewTestUser(email, "John", "Doe")
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Update user data
		user.FirstName = "Jane"
		user.LastName = "Smith"
		user.State = domain.Pending
		user.Password = "newpassword123"
		user.Gender = "woman"
		user.Telephone = "+33987654321"
		user.PostalCode = "75001"
		user.City = "Paris"
		user.Address1 = "123 Updated Street"
		user.Address2 = "Apt 456"
		user.BirthDate = time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)

		// Re-encrypt with new data
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.UpdateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify the update by retrieving the user
		updatedUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedUser)

		// Decrypt to verify data was updated
		err = crypto.DecryptStruct(ctx, updatedUser)
		require.NoError(t, err)

		assert.Equal(t, user.ID, updatedUser.ID)
		assert.Equal(t, domain.Pending, updatedUser.State)
		assert.Equal(t, "Jane", updatedUser.FirstName)
		assert.Equal(t, "Smith", updatedUser.LastName)
		assert.Equal(t, "woman", updatedUser.Gender)
		assert.Equal(t, "+33987654321", updatedUser.Telephone)
		assert.Equal(t, "75001", updatedUser.PostalCode)
		assert.Equal(t, "Paris", updatedUser.City)
		assert.Equal(t, "123 Updated Street", updatedUser.Address1)
		assert.Equal(t, "Apt 456", updatedUser.Address2)
		assert.Equal(t, time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC), updatedUser.BirthDate)

		// Verify encrypted fields were updated
		assert.NotEmpty(t, updatedUser.FirstNameEncrypted)
		assert.NotEmpty(t, updatedUser.LastNameEncrypted)
		assert.NotEmpty(t, updatedUser.TelephoneEncrypted)
		assert.NotEmpty(t, updatedUser.TelephoneHash)
	})

	t.Run("should return not found error when updating non-existent user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := td.NewTestUser("nonexistent@example.com", "Ghost", "User")
		user.ID = uuid.New() // Ensure this ID doesn't exist
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.UpdateUser(ctx, user)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should successfully update user state from unverified to pending", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "statechange@example.com"

		// Create user with unverified state
		user := td.NewTestUser(email, "State", "User")
		user.State = domain.Unverified
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Update user state to pending
		user.State = domain.Pending
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.UpdateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify the state change
		updatedUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.Pending, updatedUser.State)
	})

	t.Run("should successfully update user with empty optional fields", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "emptyfields@example.com"

		// Create user with some fields populated
		user := td.NewTestUser(email, "Full", "User")
		user.Address2 = "Suite 100"
		user.Telephone = "+33123456789"
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Update user to remove optional fields
		user.Address2 = ""
		user.Telephone = ""
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.UpdateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify the update
		updatedUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		err = crypto.DecryptStruct(ctx, updatedUser)
		require.NoError(t, err)

		assert.Equal(t, "", updatedUser.Address2)
		assert.Equal(t, "", updatedUser.Telephone)
		assert.Equal(t,
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			updatedUser.TelephoneHash,
		)
	})

	t.Run("should handle database connection errors", func(t *testing.T) {
		// This test would typically involve mocking the database connection
		// For now, we'll skip it since we're using real testcontainers
		t.Skip("Database connection error testing requires mocking")
	})

	t.Run("should successfully update all user fields in complete user flow", func(t *testing.T) {
		// Arrange - simulates the complete user registration flow
		td.ClearUsersTable(t, ctx, testPool)
		email := "complete@example.com"

		// Start with minimal unverified user
		user := &domain.User{
			ID:    uuid.New(),
			State: domain.Unverified,
			Email: email,
		}
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, user)
		require.NoError(t, err)

		// Update with complete user information
		user.State = domain.Pending
		user.FirstName = "Complete"
		user.LastName = "Registration"
		user.Password = "securepassword123"
		user.Gender = "non_binary"
		user.Telephone = "+33456789012"
		user.PostalCode = "69000"
		user.City = "Lyon"
		user.Address1 = "456 Complete Avenue"
		user.Address2 = "Building C"
		user.BirthDate = time.Date(1985, 12, 25, 0, 0, 0, 0, time.UTC)

		// Re-encrypt with complete data
		err = crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.UpdateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify all fields were updated correctly
		completeUser, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		err = crypto.DecryptStruct(ctx, completeUser)
		require.NoError(t, err)

		assert.Equal(t, domain.Pending, completeUser.State)
		assert.Equal(t, "Complete", completeUser.FirstName)
		assert.Equal(t, "Registration", completeUser.LastName)
		assert.Equal(t, "non_binary", completeUser.Gender)
		assert.Equal(t, "+33456789012", completeUser.Telephone)
		assert.Equal(t, "69000", completeUser.PostalCode)
		assert.Equal(t, "Lyon", completeUser.City)
		assert.Equal(t, "456 Complete Avenue", completeUser.Address1)
		assert.Equal(t, "Building C", completeUser.Address2)
		assert.NotEmpty(t, completeUser.PasswordHash) // Password should be hashed

		// Verify all encrypted fields are populated
		assert.NotEmpty(t, completeUser.FirstNameEncrypted)
		assert.NotEmpty(t, completeUser.LastNameEncrypted)
		assert.NotEmpty(t, completeUser.GenderEncrypted)
		assert.NotEmpty(t, completeUser.TelephoneEncrypted)
		assert.NotEmpty(t, completeUser.PostalCodeEncrypted)
		assert.NotEmpty(t, completeUser.CityEncrypted)
		assert.NotEmpty(t, completeUser.Address1Encrypted)
		assert.NotEmpty(t, completeUser.Address2Encrypted)
		assert.NotEmpty(t, completeUser.BirthDateEncrypted)
	})
}
