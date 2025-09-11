package userRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUserByAppleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return user when Apple ID exists", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Create test user with Apple ID
		testUser := &domain.User{
			ID:             uuid.New(),
			State:          domain.Active,
			EmailHash:      "test@example.com",
			EmailEncrypted: []byte("encrypted_email"),
			PasswordHash:   "hashed_password",
			AppleIDEncrypted: []byte("encrypted_apple_id_123"),
			CreatedAtEncrypted: []byte("encrypted_created_at"),
			DEKEncrypted:   []byte("encrypted_dek"),
			KeyVersion:     1,
		}

		err := testRepo.CreateUser(ctx, testUser)
		require.NoError(t, err)

		// Test retrieval by Apple ID
		retrievedUser, err := testRepo.GetUserByAppleID(ctx, string(testUser.AppleIDEncrypted))
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, retrievedUser.ID)
		assert.Equal(t, testUser.AppleIDEncrypted, retrievedUser.AppleIDEncrypted)
		assert.Equal(t, testUser.State, retrievedUser.State)
	})

	t.Run("should return ErrRepositoryNotFound when Apple ID does not exist", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Try to get user with non-existent Apple ID
		_, err := testRepo.GetUserByAppleID(ctx, "non_existent_apple_id")
		require.Error(t, err)
		assert.True(t, errs.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should return error when Apple ID is empty", func(t *testing.T) {
		_, err := testRepo.GetUserByAppleID(ctx, "")
		require.Error(t, err)
	})
}

func TestRepository_ExistsByAppleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when Apple ID exists", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Create test user with Apple ID
		testUser := &domain.User{
			ID:             uuid.New(),
			State:          domain.Active,
			EmailHash:      "test2@example.com",
			EmailEncrypted: []byte("encrypted_email"),
			PasswordHash:   "hashed_password",
			AppleIDEncrypted: []byte("encrypted_apple_id_456"),
			CreatedAtEncrypted: []byte("encrypted_created_at"),
			DEKEncrypted:   []byte("encrypted_dek"),
			KeyVersion:     1,
		}

		err := testRepo.CreateUser(ctx, testUser)
		require.NoError(t, err)

		// Test existence check
		exists, err := testRepo.ExistsByAppleID(ctx, string(testUser.AppleIDEncrypted))
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when Apple ID does not exist", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Check for non-existent Apple ID
		exists, err := testRepo.ExistsByAppleID(ctx, "non_existent_apple_id")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}