package userRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestExistsByAppleID TEST_PATH=internal/authuser/infrastructure/postgres/user/exists_by_apple_id_test.go

func TestExistsByAppleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when Apple ID exists", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		testUser := td.NewTestUserEncx(t)
		err := td.InsertUserEncx(t, ctx, testUser, testPool)
		require.NoError(t, err)

		// Test existence check
		exists, err := repo.ExistsByAppleID(ctx, string(testUser.AppleIDEncrypted))
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when Apple ID does not exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Check for non-existent Apple ID
		exists, err := repo.ExistsByAppleID(ctx, "non_existent_apple_id")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
