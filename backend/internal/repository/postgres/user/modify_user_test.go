package userRepository_test

//
// import (
// 	"context"
// 	"testing"
//
// 	"github.com/hengadev/leviosa/internal/domain/user/models"
// 	rp "github.com/hengadev/leviosa/internal/repository"
// 	"github.com/hengadev/leviosa/internal/repository/postgres"
// 	"github.com/hengadev/leviosa/internal/repository/postgres/user"
// 	"github.com/hengadev/leviosa/tests/utils"
// 	"github.com/hengadev/leviosa/tests/utils/factories"
//
// 	"github.com/hengadev/test-assert"
// )
//
// func TestModifyAccount(t *testing.T) {
// 	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
// 	user := factories.NewBasicUser(nil)
// 	whereMap := map[string]any{"id": user.ID}
// 	changes := map[string]any{
// 		"ID":        "",
// 		"Email":     "",
// 		"Password":  "",
// 		"FirstName": "Jane",
// 		"Gender":    models.GenderInput{Gender: models.GenderWoman},
// 		"GoogleID":  "",
// 		"AppleID":   "",
// 	}
// 	modifiedUser := factories.NewBasicUser(changes)
// 	tests := []struct {
// 		name         string
// 		modifiedUser *models.User
// 		expectedErr  error
// 		version      int64
// 	}{
// 		{
// 			name:         "udpate with no user in database",
// 			modifiedUser: modifiedUser,
// 			expectedErr:  rp.ErrNotUpdated,
// 			version:      20240811085134,
// 		},
// 		{
// 			name:         "update 'user' prohibited field(s)",
// 			modifiedUser: user,
// 			expectedErr:  rp.ErrValidation,
// 			version:      20240811140841,
// 		},
// 		{
// 			name:         "nominal case with valid updatable user",
// 			modifiedUser: modifiedUser,
// 			expectedErr:  nil,
// 			version:      20240811140841,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			ctx := context.Background()
// 			repo, teardown := sqlite.SetupRepository(t, ctx, tt.version, userRepository.New)
// 			defer teardown()
// 			err := repo.ModifyAccount(ctx, tt.modifiedUser, whereMap)
// 			assert.EqualError(t, err, tt.expectedErr)
// 		})
// 	}
// }
