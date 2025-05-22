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

func TestGetAllUsers(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	usersList := []*models.User{factories.Johndoe, factories.Janedoe, factories.Jeandoe}
	tests := []struct {
		name          string
		version       int64
		expectedUsers []*models.User
		expectedError error
	}{
		{
			name:          "No users in database",
			version:       20240811085134,
			expectedUsers: []*models.User{},
			expectedError: nil,
		},
		{
			name:          "Multiple users in database to retrieve",
			version:       20240819182030,
			expectedUsers: usersList,
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
			defer teardown()
			users, err := repo.GetAllUsers(ctx)
			assert.EqualError(t, err, tt.expectedError)
			fields := []string{}
			for i := range len(users) {
				assert.FieldsEqual(t, users[i], tt.expectedUsers[i], fields)
			}
		})
	}
}
