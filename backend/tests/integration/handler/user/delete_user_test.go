package userHandler_test

// import (
// 	"context"
// 	"encoding/json"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
// 	sessionService "github.com/hengadev/leviosa/internal/domain/session"
// 	"github.com/hengadev/leviosa/internal/domain/user/models"
// 	miniredis "github.com/hengadev/leviosa/internal/repository/redis"
// 	"github.com/hengadev/leviosa/internal/server/app"
// 	userHandler "github.com/hengadev/leviosa/internal/server/handler/user"
// 	"github.com/hengadev/leviosa/pkg/ctxutil"
// 	test "github.com/hengadev/leviosa/tests/utils"
// 	"github.com/hengadev/leviosa/tests/utils/factories"
// 	assert "github.com/hengadev/test-assert"
// )
//
// // TEST:
// // - no logger in context
// // - no userID in context
// // - deleteUser error (I should be able to mock the behaviour no ?)
// // - removeSession error (I should be able to mock the behaviour no ?)
//
// func TestDeleteUser(t *testing.T) {
// 	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
// 	user := factories.NewBasicUser(nil)
// 	tests := []struct {
// 		name               string
// 		version            int64
// 		userID             string
// 		initMap            miniredis.InitMap[*sessionService.Values]
// 		expectedStatusCode int
// 	}{
// 		{
// 			name:               "empty database",
// 			version:            20240811085134,
// 			userID:             user.ID,
// 			initMap:            factories.NewBasicInitSession(nil),
// 			expectedStatusCode: 500,
// 		},
// 		{
// 			name:               "user not in database",
// 			version:            20240811140841,
// 			userID:             test.GenerateRandomString(36),
// 			initMap:            factories.NewBasicInitSession(nil),
// 			expectedStatusCode: 500,
// 		},
// 		{
// 			name:               "session not found",
// 			version:            20240811140841,
// 			userID:             user.ID,
// 			initMap:            nil,
// 			expectedStatusCode: 500,
// 		},
// 		{
// 			name:               "nominal case",
// 			version:            20240811140841,
// 			userID:             user.ID,
// 			initMap:            factories.NewBasicInitSession(nil),
// 			expectedStatusCode: 200,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			r, _ := http.NewRequest("GET", "/api/v1/me", nil)
// 			w := httptest.NewRecorder()
//
// 			r.AddCookie(factories.NewBasicCookie(nil))
//
// 			// pass userID to context
// 			ctx := context.WithValue(context.Background(), ctxutil.UserIDKey, tt.userID)
// 			ctx = context.WithValue(ctx, ctxutil.LoggerKey, &slog.Logger{})
// 			r = r.WithContext(ctx)
//
// 			userSvc, userRepo := factories.SetupUser(t, ctx, tt.version)
//
// 			sessionSvc, sessionRepo, sessionTeardown := factories.SetupSession(t, ctx, tt.initMap)
// 			defer sessionTeardown()
//
// 			svcs := &app.Services{User: userSvc, Session: sessionSvc}
// 			repos := &app.Repos{User: userRepo, Session: sessionRepo}
//
// 			h := app.New(svcs, repos)
// 			userhandler := userHandler.New(h)
//
// 			userhandler.DeleteUser(w, r)
//
// 			// parse the body for the user
// 			var user *models.User
// 			json.NewDecoder(w.Body).Decode(user)
//
// 			assert.Equal(t, w.Code, tt.expectedStatusCode)
// 		})
// 	}
// }
