package userRepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	td "github.com/Leviosa-care/authuser/test/helpers"

	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestCreateUser make test-unit-user-test

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully create a new user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := td.NewTestUser("newuser@example.com", "John", "Doe")
		user.ID = uuid.New()
		user.State = domain.Unverified

		// Process encryption
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify user was created
		exists, err := repo.ExistsByEmailHash(ctx, user.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should successfully create user with all fields populated", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := &domain.User{
			ID:         uuid.New(),
			State:      domain.Pending,
			Email:      "complete@example.com",
			Password:   "securepassword123",
			Picture:    "https://example.com/avatar.jpg",
			FirstName:  "Jane",
			LastName:   "Smith",
			BirthDate:  time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			Gender:     "female",
			Role:       "user",
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

		// Process encryption
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify user was created
		exists, err := repo.ExistsByEmailHash(ctx, user.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail to create user with duplicate email hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "duplicate@example.com"

		// Create first user
		td.InsertTestUser(t, ctx, email, "First", "User", testPool, crypto)

		// Try to create second user with same email
		duplicateUser := td.NewTestUser(email, "Second", "User")
		duplicateUser.ID = uuid.New()
		duplicateUser.State = domain.Unverified

		err := crypto.ProcessStruct(ctx, duplicateUser)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(ctx, duplicateUser)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrUniqueViolation, "Should be a unique constraint violation")
	})

	t.Run("should fail to create user without required fields", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// User missing required encrypted fields (email, password, etc.)
		invalidUser := &domain.User{
			ID:    uuid.New(),
			State: domain.Unverified,
			// Missing required fields
		}

		// Act
		err := repo.CreateUser(ctx, invalidUser)

		// Assert
		require.Error(t, err)
		// Should fail due to NOT NULL constraints
	})

	t.Run("should handle user with minimal required fields", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := &domain.User{
			ID:       uuid.New(),
			State:    domain.Unverified,
			Email:    "minimal@example.com",
			Password: "password123",
		}

		// Process encryption - this will populate required encrypted fields
		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify user was created
		exists, err := repo.ExistsByEmailHash(ctx, user.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle special characters in email", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := td.NewTestUser("test+tag@example-domain.co.uk", "Special", "Email")
		user.ID = uuid.New()
		user.State = domain.Unverified

		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify user was created
		exists, err := repo.ExistsByEmailHash(ctx, user.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle very long email addresses", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		longEmail := "very.long.email.address.with.many.dots.and.subdomains@very.long.domain.name.with.many.subdomains.example.com"
		user := td.NewTestUser(longEmail, "Long", "Email")
		user.ID = uuid.New()
		user.State = domain.Unverified

		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		// Act
		err = repo.CreateUser(ctx, user)

		// Assert
		require.NoError(t, err)

		// Verify user was created
		exists, err := repo.ExistsByEmailHash(ctx, user.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail when context is cancelled", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		user := td.NewTestUser("cancelled@example.com", "Test", "User")
		user.ID = uuid.New()
		user.State = domain.Unverified

		err := crypto.ProcessStruct(ctx, user)
		require.NoError(t, err)

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Act
		err = repo.CreateUser(cancelledCtx, user)

		// Assert
		require.Error(t, err)
		// Should be classified as a context-related error by ClassifyPgError
	})

	t.Run("should handle concurrent user creations with different emails", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		emails := []string{
			"concurrent1@example.com",
			"concurrent2@example.com",
			"concurrent3@example.com",
		}

		results := make(chan error, len(emails))

		// Act - create users concurrently
		for i, email := range emails {
			go func(idx int, e string) {
				user := td.NewTestUser(e, "Concurrent", "User")
				user.ID = uuid.New()
				user.State = domain.Unverified

				err := crypto.ProcessStruct(ctx, user)
				if err != nil {
					results <- err
					return
				}

				err = repo.CreateUser(ctx, user)
				results <- err
			}(i, email)
		}

		// Assert - collect results
		successCount := 0
		for range len(emails) {
			err := <-results
			if err == nil {
				successCount++
			}
		}

		assert.Equal(t, len(emails), successCount, "All concurrent creations should succeed")

		// Verify all users were created
		for _, email := range emails {
			testUser := td.NewTestUser(email, "Concurrent", "User")
			err := crypto.ProcessStruct(ctx, testUser)
			require.NoError(t, err)

			exists, err := repo.ExistsByEmailHash(ctx, testUser.EmailHash)
			require.NoError(t, err)
			assert.True(t, exists, "User with email %s should exist", email)
		}
	})

	t.Run("should handle different user states", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		states := []domain.UserState{
			domain.Unverified,
			domain.Pending,
			domain.Active,
		}

		for i, state := range states {
			email := fmt.Sprintf("state%d@example.com", i)
			user := td.NewTestUser(email, "State", "User")
			user.ID = uuid.New()
			user.State = state

			err := crypto.ProcessStruct(ctx, user)
			require.NoError(t, err)

			// Act
			err = repo.CreateUser(ctx, user)

			// Assert
			require.NoError(t, err, "Should create user with state %s", state)

			// Verify user was created
			exists, err := repo.ExistsByEmailHash(ctx, user.EmailHash)
			require.NoError(t, err)
			assert.True(t, exists)
		}
	})
}
