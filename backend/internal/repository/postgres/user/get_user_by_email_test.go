package userRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/user/models"
	rp "github.com/hengadev/leviosa/internal/repository"
	"github.com/hengadev/leviosa/internal/repository/sqlite"
	"github.com/hengadev/leviosa/internal/repository/sqlite/user"
	"github.com/hengadev/leviosa/tests/utils"
	"github.com/hengadev/leviosa/tests/utils/factories"

	"github.com/hengadev/test-assert"
)

func TestGetUserByEmail(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	fields := []string{
		"EmailHash",
		"Picture",
		"EncryptedCreatedAt",
		"EncryptedLoggedInAt",
		"Role",
		"EncryptedBirthDate",
		"LastName",
		"FirstName",
		"Gender",
		"Telephone",
		"PostalCode",
		"City",
		"Address1",
		"Address2",
		"GoogleID",
		"AppleID",
	}
	ctx := context.Background()
	tests := []struct {
		name          string
		version       int64
		emailHash     string
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:          "no user in database",
			version:       20240811085134,
			emailHash:     "john.doe@example.com",
			expectedUser:  nil,
			expectedError: rp.ErrNotFound,
		},
		{
			name:          "user not in database",
			version:       20240811140841,
			emailHash:     "test@example.com",
			expectedUser:  nil,
			expectedError: rp.ErrNotFound,
		},
		{
			name:          "user in database",
			version:       20240811140841,
			emailHash:     "john.doe@example.com",
			expectedUser:  factories.NewBasicUser(nil),
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
		defer teardown()
		user, err := repo.GetUserByEmail(ctx, tt.emailHash)
		assert.EqualError(t, err, tt.expectedError)
		assert.FieldsEqual(t, user, tt.expectedUser, fields)
	}
}
