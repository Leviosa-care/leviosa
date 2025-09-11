package userRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUserByGoogleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return user when Google ID exists", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Create test user with Google ID
		testUser := &domain.User{
			ID:             uuid.New(),
			State:          domain.Active,
			EmailHash:      "test@example.com",
			EmailEncrypted: []byte("encrypted_email"),
			PasswordHash:   "hashed_password",
			GoogleIDEncrypted: []byte("encrypted_google_id_123"),
			CreatedAtEncrypted: []byte("encrypted_created_at"),
			DEKEncrypted:   []byte("encrypted_dek"),
			KeyVersion:     1,
		}

		err := testRepo.CreateUser(ctx, testUser)
		require.NoError(t, err)

		// Test retrieval by Google ID
		retrievedUser, err := testRepo.GetUserByGoogleID(ctx, string(testUser.GoogleIDEncrypted))
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, retrievedUser.ID)
		assert.Equal(t, testUser.GoogleIDEncrypted, retrievedUser.GoogleIDEncrypted)
		assert.Equal(t, testUser.State, retrievedUser.State)
	})

	t.Run("should return ErrRepositoryNotFound when Google ID does not exist", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Try to get user with non-existent Google ID
		_, err := testRepo.GetUserByGoogleID(ctx, "non_existent_google_id")
		require.Error(t, err)
		assert.True(t, errs.Is(err, errs.ErrRepositoryNotFound))
	})

	t.Run("should return error when Google ID is empty", func(t *testing.T) {
		_, err := testRepo.GetUserByGoogleID(ctx, "")
		require.Error(t, err)
	})
}

func TestRepository_ExistsByGoogleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when Google ID exists", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Create test user with Google ID
		testUser := &domain.User{
			ID:             uuid.New(),
			State:          domain.Active,
			EmailHash:      "test2@example.com",
			EmailEncrypted: []byte("encrypted_email"),
			PasswordHash:   "hashed_password",
			GoogleIDEncrypted: []byte("encrypted_google_id_456"),
			CreatedAtEncrypted: []byte("encrypted_created_at"),
			DEKEncrypted:   []byte("encrypted_dek"),
			KeyVersion:     1,
		}

		err := testRepo.CreateUser(ctx, testUser)
		require.NoError(t, err)

		// Test existence check
		exists, err := testRepo.ExistsByGoogleID(ctx, string(testUser.GoogleIDEncrypted))
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when Google ID does not exist", func(t *testing.T) {
		// Clean state
		clearUsersTable(t, ctx)

		// Check for non-existent Google ID
		exists, err := testRepo.ExistsByGoogleID(ctx, "non_existent_google_id")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}