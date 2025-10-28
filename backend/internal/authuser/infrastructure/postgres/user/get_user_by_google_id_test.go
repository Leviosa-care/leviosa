package userRepository_test

import (
	"context"
	"testing"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetUserByGoogleID TEST_PATH=internal/authuser/infrastructure/postgres/user/get_user_by_google_id_test.go

func TestGetUserByGoogleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return user when Google ID exists", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Create test user with Google ID
		testUserEncx := td.NewTestUserEncx(t)

		err := td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Test retrieval by Google ID
		retrievedUser, err := repo.GetUserByGoogleID(ctx, string(testUserEncx.GoogleIDEncrypted))
		assert.NoError(t, err)
		assert.Equal(t, testUserEncx.ID, retrievedUser.ID)
		assert.Equal(t, testUserEncx.GoogleIDEncrypted, retrievedUser.GoogleIDEncrypted)
		assert.Equal(t, testUserEncx.State, retrievedUser.State)
	})

	t.Run("should return ErrRepositoryNotFound when Google ID does not exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Try to get user with non-existent Google ID
		_, err := repo.GetUserByGoogleID(ctx, "non_existent_google_id")
		assert.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})

	t.Run("should return error when Google ID is empty", func(t *testing.T) {
		_, err := repo.GetUserByGoogleID(ctx, "")
		assert.Error(t, err)
	})
}
