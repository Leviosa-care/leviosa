package userHandler_test

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"
//
// 	"github.com/hengadev/leviosa/internal/domain/user/models"
// 	"github.com/hengadev/leviosa/internal/server/app"
// 	userHandler "github.com/hengadev/leviosa/internal/server/handler/user"
// 	"github.com/hengadev/leviosa/pkg/ctxutil"
// 	test "github.com/hengadev/leviosa/tests/utils"
// 	"github.com/hengadev/leviosa/tests/utils/factories"
//
// 	assert "github.com/hengadev/test-assert"
// )
//
// func TestUpdateUser(t *testing.T) {
// 	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
// 	userToModify := *factories.Johndoe
// 	userToModify.Telephone = "0234567890"
// 	var zerotime time.Time
// 	tests := []struct {
// 		user               models.User
// 		expectedStatusCode int
// 		version            int64
// 		name               string
// 	}{
// 		{user: *factories.Johndoe, expectedStatusCode: 500, version: 20240811085134, name: "no user in database"},
// 		{user: *factories.Janedoe, expectedStatusCode: 500, version: 20240811140841, name: "user not in database"},
// 		{user: *factories.Johndoe, expectedStatusCode: 200, version: 20240811140841, name: "user with no modification to do"},
// 		{user: userToModify, expectedStatusCode: 200, version: 20240811140841, name: "nominal case"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
//
// 			userID := tt.user.ID
// 			// remove forbidden field
// 			tt.user.ID = "0"
// 			tt.user.Email = ""
// 			tt.user.Password = ""
// 			tt.user.CreatedAt = zerotime
// 			tt.user.LoggedInAt = zerotime
// 			tt.user.Role = ""
//
// 			// encode credentials for request
// 			body := factories.EncodeForBody(t, tt.user)
//
// 			// create request and responseRecorder
// 			r, _ := http.NewRequest("PUT", "/api/v1/me", body)
// 			w := httptest.NewRecorder()
//
// 			ctx := r.Context()
// 			ctx = context.WithValue(ctx, ctxutil.UserIDKey, userID)
// 			r = r.WithContext(ctx)
//
// 			// setup session service and repo
// 			usersvc, userrepo := factories.SetupUser(t, ctx, tt.version)
//
// 			appsvc := &app.Services{User: usersvc}
// 			apprepo := &app.Repos{User: userrepo}
//
// 			h := app.New(appsvc, apprepo)
// 			handler := userHandler.New(h)
//
// 			handler.UpdateUser(w, r)
//
// 			// status code assertions
// 			assert.Equal(t, w.Code, tt.expectedStatusCode)
// 		})
// 	}
// }
