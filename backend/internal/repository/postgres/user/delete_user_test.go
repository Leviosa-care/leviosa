package userRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/internal/repository/sqlite/user"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestDeleteUser(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	user := factories.NewBasicUser(nil)
	tests := []struct {
		name          string
		version       int64
		userID        string
		expectedError error
	}{
		{
			name:          "no users in database",
			version:       20240811085134,
			userID:        user.ID,
			expectedError: rp.ErrNotDeleted,
		},
		{
			name:          "ID not found in database",
			version:       20240811140841,
			userID:        uuid.NewString(),
			expectedError: rp.ErrNotDeleted,
		},
		{
			name:          "ID not found in database",
			version:       20240811140841,
			userID:        user.ID,
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
			defer teardown()
			err := repo.DeleteUser(ctx, tt.userID)
			assert.EqualError(t, err, tt.expectedError)
		})
	}
}
