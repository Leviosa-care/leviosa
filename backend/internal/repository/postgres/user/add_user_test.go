package userRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/internal/repository/sqlite/user"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestAddAccount(t *testing.T) {
	// TEST:
	// - add mail account no user in database
	// - add google account no user in database
	// - add apple account no user in database
	// - add google with google account already there
	// - add mail with mail account already there
	// - add apple with apple account already there
	// - add oauth with mail account already there
	// - add mail with oauth account already there
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	user := factories.NewBasicUser(nil)
	tests := []struct {
		name        string
		user        *models.User
		version     int64
		expectedErr error
	}{
		// TODO: I should get an error on this one
		{
			name:        "user already exists",
			user:        user,
			version:     20240811140841,
			expectedErr: nil,
		},
		// {
		// 	name:        "nominal case, user does not exist",
		// 	user:        user,
		// 	version:     20240811085134,
		// 	expectedErr: nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
			defer teardown()
			err := repo.AddUser(ctx, tt.user, models.Mail)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
