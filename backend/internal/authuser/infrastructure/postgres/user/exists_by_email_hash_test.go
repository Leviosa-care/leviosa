package userRepository_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestExistsByEmailHash TEST_PATH=internal/authuser/infrastructure/postgres/user/exists_by_email_hash_test.go

func TestExistsByEmailHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when user exists with hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		userEncx := td.NewTestUserEncx(t)
		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Act
		exists, err := repo.ExistsByEmailHash(ctx, userEncx.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when hash does not exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act - test with non-existent hash
		nonExistentHash := "nonexistent_hash_12345"
		exists, err := repo.ExistsByEmailHash(ctx, nonExistentHash)

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should return false for empty hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		exists, err := repo.ExistsByEmailHash(ctx, "")

		// Assert
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should handle case sensitivity in hashes", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "test@example.com"
		userEncx := td.NewTestUserEncx(t)
		userEncx.EmailHash = email
		userEncx.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, userEncx, testPool)
		require.NoError(t, err)

		// Act - test with different casing email (should produce different hash)
		upcaseEmail := strings.ToUpper(email)

		exists, err := repo.ExistsByEmailHash(ctx, upcaseEmail)

		// Assert - should not match
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should handle multiple users with different hashes", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		const count = 2

		email := "user%d@example.com"

		for i := range count {
			userEncx := td.NewTestUserEncx(t)
			email := fmt.Sprintf(email, i)
			userEncx.EmailEncrypted = []byte(email)
			userEncx.EmailHash = email
			userEncx.FirstNameEncrypted = []byte(fmt.Sprintf("User%d", i))
			err := td.InsertUserEncx(t, ctx, userEncx, testPool)
			require.NoError(t, err)
		}

		for i := range count {
			exists, err := repo.ExistsByEmailHash(ctx, fmt.Sprintf(email, i))
			assert.NoError(t, err)
			assert.True(t, exists)
		}

		nonExistentHash := "nonexistent_hash_99999"
		existsNon, err := repo.ExistsByEmailHash(ctx, nonExistentHash)
		assert.NoError(t, err)
		assert.False(t, existsNon)
	})

	t.Run("should handle very long hash strings", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		longEmail := "very.long.email.address@very.long.subdomains.example.com"
		longEmailUserEncx := td.NewTestUserEncx(t)
		longEmailUserEncx.EmailHash = longEmail
		longEmailUserEncx.EmailEncrypted = []byte(longEmail)

		err := td.InsertUserEncx(t, ctx, longEmailUserEncx, testPool)
		require.NoError(t, err)

		// Act - get real hash for long email
		exists, err := repo.ExistsByEmailHash(ctx, longEmailUserEncx.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle special characters in hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		email := "test+tag@example-auth.co.uk"

		specialUser := td.NewTestUserEncx(t)
		specialUser.EmailHash = email
		specialUser.EmailEncrypted = []byte(email)

		err := td.InsertUserEncx(t, ctx, specialUser, testPool)
		require.NoError(t, err)

		// Act - get real hash for special chars email
		exists, err := repo.ExistsByEmailHash(ctx, specialUser.EmailHash)

		// Assert
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return error when context is cancelled", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		// Act
		testHash := "cancelled_context_hash"
		exists, err := repo.ExistsByEmailHash(cancelledCtx, testHash)

		// Assert
		assert.Error(t, err)
		assert.False(t, exists)
		// Should be classified as a context-related error by ClassifyPgError
	})
}
