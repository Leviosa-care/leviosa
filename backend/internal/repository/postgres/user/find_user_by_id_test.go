package userRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/internal/repository/sqlite/user"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestFindAccountByID(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	user := factories.NewBasicUser(nil)
	fields := []string{"Email", "Picture", "Role", "LastName", "FirstName", "Gender", "EncryptedBirthDate", "Telephone", "PostalCode", "City", "Address1", "Address2"}
	tests := []struct {
		name          string
		version       int64
		userID        string
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:          "No user in database",
			version:       20240811085134,
			userID:        user.ID,
			expectedUser:  nil,
			expectedError: rp.ErrNotFound,
		},
		{
			name:          "ID not in database",
			version:       20240811140841,
			userID:        uuid.NewString(),
			expectedUser:  nil,
			expectedError: rp.ErrNotFound,
		},
		{
			name:          "nominal case with user in database",
			version:       20240811140841,
			userID:        user.ID,
			expectedUser:  user,
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
			defer teardown()
			user, err := repo.FindAccountByID(ctx, tt.userID)
			assert.EqualError(t, err, tt.expectedError)
			assert.FieldsEqual(t, user, tt.expectedUser, fields)
		})
	}
}
