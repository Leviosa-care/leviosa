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

func TestGetPendingUser(t *testing.T) {
	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
	user := factories.NewBasicUser(nil)
	commonFields := []string{
		"ID",
		"Email",
		"LastName",
		"FirstName",
		"Gender",
		"EncryptedBirthDate",
		"Telephone",
		"PostalCode",
		"City",
		"Address1",
		"Address2",
	}
	googleUserFields := append(commonFields, "GoogleID")
	appleUserFields := append(commonFields, "AppleID")
	tests := []struct {
		name          string
		version       int64
		emailHash     string
		provider      models.ProviderType
		fields        []string
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:          "invalid provider",
			version:       20250206154643,
			emailHash:     user.EmailHash,
			provider:      "instagram",
			fields:        commonFields,
			expectedUser:  nil,
			expectedError: rp.ErrValidation,
		},
		{
			name:          "no user in database",
			version:       20250206154643,
			emailHash:     user.EmailHash,
			provider:      models.Mail,
			fields:        commonFields,
			expectedUser:  nil,
			expectedError: rp.ErrNotFound,
		},
		{
			name:          "user not in database",
			version:       20250206154724,
			emailHash:     "user_not_in_database@test.com",
			provider:      models.Mail,
			fields:        commonFields,
			expectedUser:  nil,
			expectedError: rp.ErrNotFound,
		},
		{
			name:          "user with 'Mail' provider in database",
			version:       20250206154724,
			emailHash:     user.EmailHash,
			provider:      models.Mail,
			fields:        commonFields,
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name:          "user with 'Google' provider in database",
			version:       20250206154724,
			emailHash:     user.EmailHash,
			provider:      models.Google,
			fields:        googleUserFields,
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name:          "user with 'Apple' provider in database",
			version:       20250206154724,
			emailHash:     user.EmailHash,
			provider:      models.Apple,
			fields:        appleUserFields,
			expectedUser:  user,
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
		defer teardown()
		user, err := repo.GetPendingUser(ctx, tt.emailHash, tt.provider)
		assert.EqualError(t, err, tt.expectedError)
		assert.FieldsEqual(t, user, tt.expectedUser, tt.fields)
	}
}
