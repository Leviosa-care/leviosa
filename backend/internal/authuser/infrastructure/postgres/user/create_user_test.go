package userRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateUser TEST_PATH=internal/authuser/infrastructure/postgres/user/create_user_test.go

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	const existsQuery = `SELECT EXISTS(SELECT 1 FROM auth.users WHERE id = $1)`
	const existsByEmailHashQuery = `SELECT EXISTS(SELECT 1 FROM auth.users WHERE email_hash = $1)`

	t.Run("should successfully create a new user", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		userEncx := td.NewTestUserEncx(t)

		// Act
		err := repo.CreateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify user was created
		// TODO: write a sql statement for that
		var exists bool
		err = testPool.QueryRow(ctx, existsQuery, userEncx.ID).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail to create user with duplicate email hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "duplicate@example.com"
		emailEncryped := []byte(email)

		// Create first userEncx
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = email
		userEncx.EmailEncrypted = emailEncryped
		userEncx.FirstNameEncrypted = []byte("First")

		err := repo.CreateUser(ctx, userEncx)
		require.NoError(t, err)

		duplicateUserEncx := td.NewTestUserEncx(t)
		duplicateUserEncx.EmailHash = email
		duplicateUserEncx.EmailEncrypted = emailEncryped
		duplicateUserEncx.FirstNameEncrypted = []byte("Second")

		// Act
		err = repo.CreateUser(ctx, duplicateUserEncx)

		// Assert
		assert.ErrorIs(t, err, errs.ErrUniqueViolation, "Should be a unique constraint violation")
	})

	t.Run("should fail to create user without required fields", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// User missing required encrypted fields (email, password, etc.)
		invalidUser := &domain.UserEncx{
			ID:    uuid.New(),
			State: domain.Unverified,
			// Missing required fields
		}

		// Act
		err := repo.CreateUser(ctx, invalidUser)

		// Assert
		assert.Error(t, err)
		// TODO: check for the specific error returned
		// Should fail due to NOT NULL constraints
	})

	t.Run("should handle user with minimal required fields", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		userEncx := &domain.UserEncx{
			ID:                 uuid.New(),
			State:              domain.Unverified,
			EmailEncrypted:     []byte("minimal@example.com"),
			PasswordHashSecure: "password123",
			CreatedAtEncrypted: []byte("created_at_encrypted"),
			DEKEncrypted:       []byte("dek_encrypted"),
			KeyVersion:         1,
		}

		// Act
		err := repo.CreateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify user was created
		var exists bool
		err = testPool.QueryRow(ctx, existsQuery, userEncx.ID).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle very long email addresses", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		longEmail := "very.long.email.address.with.many.dots@domain.name.com"
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = longEmail
		userEncx.EmailEncrypted = []byte(longEmail)

		// Act
		err := repo.CreateUser(ctx, userEncx)

		// Assert
		assert.NoError(t, err)

		// Verify user was created
		exists, err := repo.ExistsByEmailHash(ctx, userEncx.EmailHash)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should fail when context is cancelled", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// userEncx := td.NewTestUser("cancelled@example.com", "Test", "User")
		userEncx := td.NewTestUserEncx(t)

		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Act
		err := repo.CreateUser(cancelledCtx, userEncx)

		// Assert
		assert.Error(t, err)
		// Should be classified as a context-related error by ClassifyPgError
	})
	t.Run("should handle concurrent user creations with different emails", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const count = 3
		emails := [count]string{
			"concurrent1@example.com",
			"concurrent2@example.com",
			"concurrent3@example.com",
		}

		results := make(chan error, count)

		// Act - create users concurrently
		for i, email := range emails {
			go func(idx int, e string) {
				userEncx := td.NewTestUserEncx(t)
				userEncx.EmailHash = email
				userEncx.EmailEncrypted = []byte(email)

				err := repo.CreateUser(ctx, userEncx)
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

		assert.Equal(t, count, successCount, "All concurrent creations should succeed")

		// Verify all users were created
		for _, email := range emails {
			var exists bool
			err := testPool.QueryRow(ctx, existsByEmailHashQuery, email).Scan(&exists)

			assert.NoError(t, err)
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
			// userEncx := td.NewTestUser(email, "State", "User")
			userEncx := td.NewTestUserEncx(t)
			userEncx.EmailHash = email
			userEncx.EmailEncrypted = []byte(email)
			userEncx.State = state

			// Act
			err := repo.CreateUser(ctx, userEncx)

			// Assert
			assert.NoError(t, err, "Should create user with state %s", state)

			// Verify user was created
			var exists bool
			err = testPool.QueryRow(ctx, existsByEmailHashQuery, email).Scan(&exists)
			assert.NoError(t, err)
			assert.True(t, exists)
		}
	})
}
