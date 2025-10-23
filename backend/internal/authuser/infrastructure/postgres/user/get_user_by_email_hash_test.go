package userRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestGetUserByEmailHash make test-unit-user-test

func TestGetUserByEmailHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully retrieve user by email hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "getuser@example.com"
		td.InsertTestUser(t, ctx, email, "John", "Doe", testPool, crypto)

		// Get the expected hash
		expectedUser := td.NewTestUser(email, "John", "Doe")
		err := crypto.ProcessStruct(ctx, expectedUser)
		require.NoError(t, err)

		// Act
		retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedUser.EmailHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)
		assert.Equal(t, expectedUser.EmailHash, retrievedUser.EmailHash)
		assert.NotEmpty(t, retrievedUser.ID)
		assert.Equal(t, domain.Unverified, string(retrievedUser.State)) // Default state from helpers

		// Verify encrypted fields are populated
		assert.NotEmpty(t, retrievedUser.EmailEncrypted)
		assert.NotEmpty(t, retrievedUser.FirstNameEncrypted)
		assert.NotEmpty(t, retrievedUser.LastNameEncrypted)
		assert.NotEmpty(t, retrievedUser.DEKEncrypted)
		assert.Greater(t, retrievedUser.KeyVersion, 0)
	})

	t.Run("should retrieve user with all fields populated", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Create a user with all fields populated
		fullUser := &domain.User{
			ID:         uuid.New(),
			State:      domain.Active,
			Email:      "fulluser@example.com",
			Password:   "securepassword123",
			Picture:    "https://example.com/avatar.jpg",
			FirstName:  "Jane",
			LastName:   "Smith",
			BirthDate:  time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:     "female",
			Role:       "admin",
			Telephone:  "+1234567890",
			PostalCode: "12345",
			City:       "New York",
			Address1:   "123 Main St",
			Address2:   "Apt 4B",
			GoogleID:   "google_12345",
			AppleID:    "apple_67890",
			CreatedAt:  time.Now(),
			LoggedInAt: time.Now().Add(-1 * time.Hour),
		}

		// Process encryption and insert
		err := crypto.ProcessStruct(ctx, fullUser)
		require.NoError(t, err)
		err = repo.CreateUser(ctx, fullUser)
		require.NoError(t, err)

		// Act
		retrievedUser, err := repo.GetUserByEmailHash(ctx, fullUser.EmailHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)

		// Verify basic fields
		assert.Equal(t, fullUser.ID, retrievedUser.ID)
		assert.Equal(t, fullUser.State, retrievedUser.State)
		assert.Equal(t, fullUser.EmailHash, retrievedUser.EmailHash)

		// Verify all encrypted fields are populated (non-empty byte arrays)
		assert.NotEmpty(t, retrievedUser.EmailEncrypted)
		assert.NotEmpty(t, retrievedUser.PictureEncrypted)
		assert.NotEmpty(t, retrievedUser.FirstNameEncrypted)
		assert.NotEmpty(t, retrievedUser.LastNameEncrypted)
		assert.NotEmpty(t, retrievedUser.BirthDateEncrypted)
		assert.NotEmpty(t, retrievedUser.GenderEncrypted)
		assert.NotEmpty(t, retrievedUser.RoleEncrypted)
		assert.NotEmpty(t, retrievedUser.TelephoneEncrypted)
		assert.NotEmpty(t, retrievedUser.PostalCodeEncrypted)
		assert.NotEmpty(t, retrievedUser.CityEncrypted)
		assert.NotEmpty(t, retrievedUser.Address1Encrypted)
		assert.NotEmpty(t, retrievedUser.Address2Encrypted)
		assert.NotEmpty(t, retrievedUser.GoogleIDEncrypted)
		assert.NotEmpty(t, retrievedUser.AppleIDEncrypted)
		assert.NotEmpty(t, retrievedUser.CreatedAtEncrypted)
		assert.NotEmpty(t, retrievedUser.LoggedInAtEncrypted)
		assert.NotEmpty(t, retrievedUser.DEKEncrypted)

		// Verify hashed fields
		assert.Equal(t, fullUser.TelephoneHash, retrievedUser.TelephoneHash)
		assert.Equal(t, fullUser.PasswordHash, retrievedUser.PasswordHash)
		assert.Equal(t, fullUser.KeyVersion, retrievedUser.KeyVersion)
	})

	t.Run("should return not found error when user does not exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		nonExistentHash := "nonexistent_hash_12345"

		// Act
		user, err := repo.GetUserByEmailHash(ctx, nonExistentHash)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should be a not found error")
	})

	t.Run("should return not found error for empty hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		user, err := repo.GetUserByEmailHash(ctx, "")

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound, "Should be a not found error")
	})

	t.Run("should handle case sensitivity in email hashes", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "casetest@example.com"
		td.InsertTestUser(t, ctx, email, "Case", "Test", testPool, crypto)

		// Get hash for original email
		originalUser := td.NewTestUser(email, "Case", "Test")
		err := crypto.ProcessStruct(ctx, originalUser)
		require.NoError(t, err)

		// Get hash for different case email
		upcaseUser := td.NewTestUser("CASETEST@EXAMPLE.COM", "Case", "Test")
		err = crypto.ProcessStruct(ctx, upcaseUser)
		require.NoError(t, err)

		// Act - should find with original hash
		foundUser, err := repo.GetUserByEmailHash(ctx, originalUser.EmailHash)
		require.NoError(t, err)
		assert.NotNil(t, foundUser)

		// Act - should NOT find with different case hash
		notFoundUser, err := repo.GetUserByEmailHash(ctx, upcaseUser.EmailHash)
		require.Error(t, err)
		assert.Nil(t, notFoundUser)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should retrieve correct user when multiple users exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		users := []struct {
			email     string
			firstName string
			lastName  string
		}{
			{"multi1@example.com", "User", "One"},
			{"multi2@example.com", "User", "Two"},
			{"multi3@example.com", "User", "Three"},
		}

		expectedHashes := make([]string, len(users))

		// Insert all users and collect their hashes
		for i, u := range users {
			td.InsertTestUser(t, ctx, u.email, u.firstName, u.lastName, testPool, crypto)

			testUser := td.NewTestUser(u.email, u.firstName, u.lastName)
			err := crypto.ProcessStruct(ctx, testUser)
			require.NoError(t, err)
			expectedHashes[i] = testUser.EmailHash
		}

		// Act & Assert - retrieve each user by their hash
		for i, expectedHash := range expectedHashes {
			retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedHash)
			require.NoError(t, err, "Should find user %d", i+1)
			require.NotNil(t, retrievedUser)
			assert.Equal(t, expectedHash, retrievedUser.EmailHash)
		}
	})

	t.Run("should handle special characters in email hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		specialEmail := "test+tag@example-auth.co.uk"
		td.InsertTestUser(t, ctx, specialEmail, "Special", "Email", testPool, crypto)

		// Get expected hash
		expectedUser := td.NewTestUser(specialEmail, "Special", "Email")
		err := crypto.ProcessStruct(ctx, expectedUser)
		require.NoError(t, err)

		// Act
		retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedUser.EmailHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)
		assert.Equal(t, expectedUser.EmailHash, retrievedUser.EmailHash)
	})

	t.Run("should handle very long email addresses", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		longEmail := "very.long.email.address.with.many.dots.and.subdomains@very.long.auth.name.with.many.subdomains.example.com"
		td.InsertTestUser(t, ctx, longEmail, "Long", "Email", testPool, crypto)

		// Get expected hash
		expectedUser := td.NewTestUser(longEmail, "Long", "Email")
		err := crypto.ProcessStruct(ctx, expectedUser)
		require.NoError(t, err)

		// Act
		retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedUser.EmailHash)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, retrievedUser)
		assert.Equal(t, expectedUser.EmailHash, retrievedUser.EmailHash)
	})

	t.Run("should fail when context is cancelled", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "cancelled@example.com"
		td.InsertTestUser(t, ctx, email, "Test", "User", testPool, crypto)

		testUser := td.NewTestUser(email, "Test", "User")
		err := crypto.ProcessStruct(ctx, testUser)
		require.NoError(t, err)

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Act
		user, err := repo.GetUserByEmailHash(cancelledCtx, testUser.EmailHash)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		// Should be classified as a context-related error by ClassifyPgError
	})

	t.Run("should retrieve users with different states correctly", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		states := []domain.UserState{
			domain.Unverified,
			domain.Pending,
			domain.Active,
		}

		expectedHashes := make([]string, len(states))

		for i, state := range states {
			email := fmt.Sprintf("state%d@example.com", i)

			// Create user with specific state
			user := td.NewTestUser(email, "State", "User")
			user.ID = uuid.New()
			user.State = state

			err := crypto.ProcessStruct(ctx, user)
			require.NoError(t, err)

			err = repo.CreateUser(ctx, user)
			require.NoError(t, err)

			expectedHashes[i] = user.EmailHash
		}

		// Act & Assert - retrieve each user and verify state
		for i, expectedHash := range expectedHashes {
			retrievedUser, err := repo.GetUserByEmailHash(ctx, expectedHash)
			require.NoError(t, err, "Should find user with state %s", states[i])
			require.NotNil(t, retrievedUser)
			assert.Equal(t, expectedHash, retrievedUser.EmailHash)
			assert.Equal(t, states[i], retrievedUser.State)
		}
	})

	t.Run("should handle concurrent retrievals of same user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "concurrent@example.com"
		td.InsertTestUser(t, ctx, email, "Concurrent", "User", testPool, crypto)

		expectedUser := td.NewTestUser(email, "Concurrent", "User")
		err := crypto.ProcessStruct(ctx, expectedUser)
		require.NoError(t, err)

		// Act - perform concurrent retrievals
		numGoroutines := 5
		results := make(chan *domain.User, numGoroutines)
		errors := make(chan error, numGoroutines)

		for range numGoroutines {
			go func() {
				user, err := repo.GetUserByEmailHash(ctx, expectedUser.EmailHash)
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
				assert.Equal(t, expectedUser.EmailHash, user.EmailHash)
			}
		}

		assert.Equal(t, numGoroutines, successCount, "All concurrent retrievals should succeed")
	})
}
