package userRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetUserByAppleID TEST_PATH=internal/authuser/infrastructure/postgres/user/get_user_by_apple_id_test.go

func TestGetUserByAppleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return user when Apple ID exists", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		testUser := td.NewTestUserEncx(t)
		err := td.InsertUserEncx(t, ctx, testUser, testPool)
		require.NoError(t, err)

		// Test retrieval by Apple ID
		retrievedUser, err := repo.GetUserByAppleID(ctx, string(testUser.AppleIDEncrypted))
		assert.NoError(t, err)

		assert.Equal(t, testUser.ID, retrievedUser.ID)
		assert.Equal(t, testUser.AppleIDEncrypted, retrievedUser.AppleIDEncrypted)
		assert.Equal(t, testUser.State, retrievedUser.State)
	})

	t.Run("should return ErrRepositoryNotFound when Apple ID does not exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Try to get user with non-existent Apple ID
		_, err := repo.GetUserByAppleID(ctx, "non_existent_apple_id")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should return error when Apple ID is empty", func(t *testing.T) {
		_, err := repo.GetUserByAppleID(ctx, "")
		assert.Error(t, err)
	})
}
