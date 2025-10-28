package userRepository_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetUserByEmailHash TEST_PATH=internal/authuser/infrastructure/postgres/user/get_user_by_email_hash_test.go

func TestGetUserByEmailHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve user by email hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		expectedUserEncx := td.NewTestUserEncx(t)
		err := td.InsertUserEncx(t, ctx, expectedUserEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByEmailHash(ctx, expectedUserEncx.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)
		assert.Equal(t, expectedUserEncx.EmailHash, retrievedUserEncx.EmailHash)
		assert.NotEmpty(t, retrievedUserEncx.ID)
		assert.Equal(t, domain.Unverified, retrievedUserEncx.State) // Default state from helpers

		// Verify encrypted fields are populated
		assert.NotEmpty(t, retrievedUserEncx.EmailEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.FirstNameEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.LastNameEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.DEKEncrypted)
		assert.Greater(t, retrievedUserEncx.KeyVersion, 0)
	})

	t.Run("should retrieve user with all fields populated", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		fullUserEncx := td.NewTestUserEncx(t)
		err := td.InsertUserEncx(t, ctx, fullUserEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByEmailHash(ctx, fullUserEncx.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)

		// Verify basic fields
		assert.Equal(t, fullUserEncx.ID, retrievedUserEncx.ID)
		assert.Equal(t, fullUserEncx.State, retrievedUserEncx.State)
		assert.Equal(t, fullUserEncx.EmailHash, retrievedUserEncx.EmailHash)

		// Verify all encrypted fields are populated (non-empty byte arrays)
		assert.NotEmpty(t, retrievedUserEncx.EmailEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.PictureEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.FirstNameEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.LastNameEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.BirthDateEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.GenderEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.RoleEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.TelephoneEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.PostalCodeEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.CityEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.Address1Encrypted)
		assert.NotEmpty(t, retrievedUserEncx.Address2Encrypted)
		assert.NotEmpty(t, retrievedUserEncx.GoogleIDEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.AppleIDEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.CreatedAtEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.LoggedInAtEncrypted)
		assert.NotEmpty(t, retrievedUserEncx.DEKEncrypted)

		// Verify hashed fields
		assert.Equal(t, fullUserEncx.TelephoneHash, retrievedUserEncx.TelephoneHash)
		assert.Equal(t, fullUserEncx.PasswordHashSecure, retrievedUserEncx.PasswordHashSecure)
		assert.Equal(t, fullUserEncx.KeyVersion, retrievedUserEncx.KeyVersion)
	})

	t.Run("should return not found error when user does not exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		nonExistentHash := "nonexistent_hash_12345"

		// Act
		user, err := repo.GetUserByEmailHash(ctx, nonExistentHash)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should be a not found error")
	})

	t.Run("should return not found error for empty hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		user, err := repo.GetUserByEmailHash(ctx, "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should be a not found error")
	})

	t.Run("should handle case sensitivity in email hashes", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		email := "casetest@example.com"

		originalUserEncx := td.NewTestUserEncx(t)
		originalUserEncx.EmailHash = email
		originalUserEncx.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, originalUserEncx, testPool)
		require.NoError(t, err)

		// Get hash for different case email
		upcaseEmail := strings.ToUpper(email)

		// Act - should find with original hash
		foundUser, err := repo.GetUserByEmailHash(ctx, originalUserEncx.EmailHash)
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)

		// Act - should NOT find with different case hash
		notFoundUser, err := repo.GetUserByEmailHash(ctx, upcaseEmail)
		assert.Error(t, err)
		assert.Nil(t, notFoundUser)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should retrieve correct user when multiple users exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const count = 3

		users := [3]struct {
			email     string
			firstName string
			lastName  string
		}{
			{"multi1@example.com", "User", "One"},
			{"multi2@example.com", "User", "Two"},
			{"multi3@example.com", "User", "Three"},
		}

		expectedHashes := make([]string, count)

		// Insert all users and collect their hashes
		for i, u := range users {
			userEncx := td.NewTestUserEncx(t)
			userEncx.EmailHash = u.email
			userEncx.EmailEncrypted = []byte(u.email)
			userEncx.FirstNameEncrypted = []byte(u.firstName)
			userEncx.LastNameEncrypted = []byte(u.lastName)

			err := td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)

			expectedHashes[i] = userEncx.EmailHash
		}

		// Act & Assert - retrieve each user by their hash
		for i, expectedHash := range expectedHashes {
			retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedHash)
			assert.NoError(t, err, "Should find user %d", i+1)
			assert.NotNil(t, retrievedUser)
			assert.Equal(t, expectedHash, retrievedUser.EmailHash)
		}
	})

	t.Run("should handle special characters in email hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		specialEmail := "test+tag@example-auth.co.uk"

		// Get expected hash
		expectedUserEncx := td.NewTestUserEncx(t)
		expectedUserEncx.EmailHash = specialEmail
		expectedUserEncx.EmailEncrypted = []byte(specialEmail)

		err := td.InsertUserEncx(t, ctx, expectedUserEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByEmailHash(ctx, expectedUserEncx.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)
		assert.Equal(t, expectedUserEncx.EmailHash, retrievedUserEncx.EmailHash)
	})

	t.Run("should handle very long email addresses", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		longEmail := "very.long.email.address@very.long.auth.name.com"

		expectedUserEncx := td.NewTestUserEncx(t)
		expectedUserEncx.EmailEncrypted = []byte(longEmail)
		expectedUserEncx.EmailHash = longEmail

		err := td.InsertUserEncx(t, ctx, expectedUserEncx, testPool)
		require.NoError(t, err)

		// Act
		retrievedUserEncx, err := repo.GetUserByEmailHash(ctx, expectedUserEncx.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUserEncx)
		assert.Equal(t, expectedUserEncx.EmailHash, retrievedUserEncx.EmailHash)
	})

	t.Run("should fail when context is cancelled", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		email := "cancelled@example.com"

		testUserEncx := td.NewTestUserEncx(t)
		testUserEncx.EmailHash = email
		testUserEncx.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Act
		user, err := repo.GetUserByEmailHash(cancelledCtx, testUserEncx.EmailHash)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		// Should be classified as a context-related error by ClassifyPgError
	})

	t.Run("should retrieve users with different states correctly", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const count = 3

		states := [count]domain.UserState{
			domain.Unverified,
			domain.Pending,
			domain.Active,
		}

		expectedHashes := make([]string, len(states))

		for i, state := range states {
			email := fmt.Sprintf("state%d@example.com", i)

			// Create user with specific state
			user := td.NewTestUserEncx(t)
			user.EmailEncrypted = []byte(email)
			user.EmailHash = email
			user.FirstNameEncrypted = []byte(fmt.Sprintf("User%d", i))
			user.LastNameEncrypted = []byte("State")

			user.ID = uuid.New()
			user.State = state

			err := td.InsertUserEncx(t, ctx, user, testPool)
			require.NoError(t, err)

			expectedHashes[i] = user.EmailHash
		}

		// Act & Assert - retrieve each user and verify state
		for i, expectedHash := range expectedHashes {
			retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedHash)
			assert.NoError(t, err, "Should find user with state %s", states[i])
			assert.NotNil(t, retrievedUser)
			assert.Equal(t, expectedHash, retrievedUser.EmailHash)
			assert.Equal(t, states[i], retrievedUser.State)
		}
	})

	t.Run("should handle concurrent retrievals of same user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "concurrent@example.com"

		expectedUserEncx := td.NewTestUserEncx(t)
		expectedUserEncx.EmailHash = email
		expectedUserEncx.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, expectedUserEncx, testPool)
		require.NoError(t, err)

		// Act - perform concurrent retrievals
		numGoroutines := 5
		results := make(chan *domain.UserEncx, numGoroutines)
		errors := make(chan error, numGoroutines)

		for range numGoroutines {
			go func() {
				user, err := repo.GetUserByEmailHash(ctx, expectedUserEncx.EmailHash)
				results <- user
				errors <- err
			}()
		}

		// Assert - collect results
		successCount := 0
		for range numGoroutines {
			user := <-results
			err := <-errors

			if err == nil && user != nil {
				successCount++
				assert.Equal(t, expectedUserEncx.EmailHash, user.EmailHash)
			}
		}

		assert.Equal(t, numGoroutines, successCount, "All concurrent retrievals should succeed")
	})
}
