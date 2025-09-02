package userRepository_test

import (
	"context"
	"fmt"
	"testing"

	td "github.com/Leviosa-care/authuser/test/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TEST=TestExistsByEmailHash make test-unit-user-test

func TestExistsByEmailHash(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when user exists with hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "existing@example.com"
		td.InsertTestUser(t, ctx, email, "John", "Doe", testPool, crypto)

		// Get the actual hash from the crypto service
		testUser := td.NewTestUser(email, "John", "Doe")

		fmt.Println("the user email hash after the NewTestUser call is:", testUser.EmailHash)
		err := crypto.ProcessStruct(ctx, testUser)
		require.NoError(t, err)

		// Act
		exists, err := repo.ExistsByEmailHash(ctx, testUser.EmailHash)

		// Assert
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when hash does not exist", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act - test with non-existent hash
		nonExistentHash := "nonexistent_hash_12345"
		exists, err := repo.ExistsByEmailHash(ctx, nonExistentHash)

		// Assert
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should return false for empty hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)

		// Act
		exists, err := repo.ExistsByEmailHash(ctx, "")

		// Assert
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should handle case sensitivity in hashes", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "test@example.com"
		td.InsertTestUser(t, ctx, email, "Test", "User", testPool, crypto)

		// Act - test with different casing email (should produce different hash)
		upcaseUser := td.NewTestUser("TEST@EXAMPLE.COM", "Test", "User")
		err := crypto.ProcessStruct(ctx, upcaseUser)
		require.NoError(t, err)
		exists, err := repo.ExistsByEmailHash(ctx, upcaseUser.EmailHash)

		// Assert - should not match
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should handle multiple users with different hashes", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		td.InsertTestUser(t, ctx, "user1@example.com", "User", "One", testPool, crypto)
		td.InsertTestUser(t, ctx, "user2@example.com", "User", "Two", testPool, crypto)
		td.InsertTestUser(t, ctx, "user3@example.com", "User", "Three", testPool, crypto)

		// Act & Assert - get real hashes
		user1 := td.NewTestUser("user1@example.com", "User", "One")
		err := crypto.ProcessStruct(ctx, user1)
		require.NoError(t, err)
		exists1, err := repo.ExistsByEmailHash(ctx, user1.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists1)

		user2 := td.NewTestUser("user2@example.com", "User", "Two")
		err = crypto.ProcessStruct(ctx, user2)
		require.NoError(t, err)
		exists2, err := repo.ExistsByEmailHash(ctx, user2.EmailHash)
		require.NoError(t, err)
		assert.True(t, exists2)

		nonExistentHash := "nonexistent_hash_99999"
		existsNon, err := repo.ExistsByEmailHash(ctx, nonExistentHash)
		require.NoError(t, err)
		assert.False(t, existsNon)
	})

	t.Run("should handle very long hash strings", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		longEmail := "very.long.email.address.with.many.dots.and.subdomains@very.long.auth.name.with.many.subdomains.example.com"
		td.InsertTestUser(t, ctx, longEmail, "Long", "Email", testPool, crypto)

		// Act - get real hash for long email
		longEmailUser := td.NewTestUser(longEmail, "Long", "Email")
		err := crypto.ProcessStruct(ctx, longEmailUser)
		require.NoError(t, err)
		exists, err := repo.ExistsByEmailHash(ctx, longEmailUser.EmailHash)

		// Assert
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should handle special characters in hash", func(t *testing.T) {
		// Arrange
		td.ClearUsersTable(t, ctx, testPool)
		email := "test+tag@example-auth.co.uk"
		td.InsertTestUser(t, ctx, email, "Special", "Email", testPool, crypto)

		// Act - get real hash for special chars email
		specialUser := td.NewTestUser(email, "Special", "Email")
		err := crypto.ProcessStruct(ctx, specialUser)
		require.NoError(t, err)
		exists, err := repo.ExistsByEmailHash(ctx, specialUser.EmailHash)

		// Assert
		require.NoError(t, err)
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
		require.Error(t, err)
		assert.False(t, exists)
		// Should be classified as a context-related error by ClassifyPgError
	})
}
