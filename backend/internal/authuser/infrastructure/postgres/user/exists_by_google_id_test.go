package userRepository_test

import (
	"context"
	"testing"

	td "github.com/Leviosa-care/leviosa/backend/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestExistsByGoogleID TEST_PATH=internal/authuser/infrastructure/postgres/user/exists_by_google_id_test.go

func TestExistsByGoogleID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return true when Google ID exists", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Create test user with Google ID
		testUserEncx := td.NewTestUserEncx(t)

		err := td.InsertUserEncx(t, ctx, testUserEncx, testPool)
		require.NoError(t, err)

		// Test existence check
		exists, err := repo.ExistsByGoogleID(ctx, string(testUserEncx.GoogleIDEncrypted))
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return false when Google ID does not exist", func(t *testing.T) {
		// Clean state
		td.ClearUsersTable(t, ctx, testPool)

		// Check for non-existent Google ID
		exists, err := repo.ExistsByGoogleID(ctx, "non_existent_google_id")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
