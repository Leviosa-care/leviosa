package userHandler_test

// import (
// 	"context"
// 	"encoding/json"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
// 	"github.com/hengadev/leviosa/internal/domain/user/models"
// 	"github.com/hengadev/leviosa/internal/server/app"
// 	userHandler "github.com/hengadev/leviosa/internal/server/handler/user"
// 	"github.com/hengadev/leviosa/pkg/ctxutil"
// 	test "github.com/hengadev/leviosa/tests/utils"
// 	"github.com/hengadev/leviosa/tests/utils/factories"
// 	"github.com/google/uuid"
//
// 	assert "github.com/hengadev/test-assert"
// )
//
// func TestGetUser(t *testing.T) {
// 	// TEST:
// 	// - no userID in context
// 	// - no logger in context
// 	// - no user in database
// 	// - userID invalid / or not in database
// 	// - nominal case
// 	t.Setenv("TEST_MIGRATION_PATH", test.GetSQLiteMigrationPath())
// 	user := factories.NewBasicUser(nil)
// 	userID := user.ID
// 	tests := []struct {
// 		name               string
// 		version            int64
// 		setupContext       func() context.Context
// 		userID             string
// 		expectedStatusCode int
// 		expectedUser       *models.User
// 	}{
// 		// {
// 		// 	name:    "userID missing in context",
// 		// 	version: 20240819182030,
// 		// 	setupContext: func() context.Context {
// 		// 		return context.WithValue(context.Background(), ctxutil.LoggerKey, *slog.Default())
// 		// 	},
// 		// 	userID:             userID,
// 		// 	expectedStatusCode: 500,
// 		// 	expectedUser:       nil,
// 		// },
// 		// {
// 		// 	name:    "logger missing in context",
// 		// 	version: 20240819182030,
// 		// 	setupContext: func() context.Context {
// 		// 		return context.WithValue(context.Background(), ctxutil.UserIDKey, userID)
// 		// 	},
// 		// 	userID:             userID,
// 		// 	expectedStatusCode: 500,
// 		// 	expectedUser:       nil,
// 		// },
// 		// {
// 		// 	name:    "no user in database",
// 		// 	version: 20240811085134,
// 		// 	// NOTE: no logger in the database
// 		// 	setupContext: func() context.Context {
// 		// 		ctx := context.WithValue(context.Background(), ctxutil.LoggerKey, slog.Default())
// 		// 		return context.WithValue(ctx, ctxutil.UserIDKey, userID)
// 		// 	},
// 		// 	userID:             userID,
// 		// 	expectedStatusCode: 404,
// 		// 	expectedUser:       nil,
// 		// },
// 		{
// 			name:    "ID not in database",
// 			version: 20240819182030,
// 			setupContext: func() context.Context {
// 				ctx := context.WithValue(context.Background(), ctxutil.LoggerKey, slog.Default())
// 				return context.WithValue(ctx, ctxutil.UserIDKey, userID)
// 			},
// 			userID:             uuid.NewString(),
// 			expectedStatusCode: 404,
// 			expectedUser:       nil,
// 		},
// 		// {
// 		// 	name:    "nominal case",
// 		// 	version: 20240811140841,
// 		// 	setupContext: func() context.Context {
// 		// 		ctx := context.WithValue(context.Background(), ctxutil.LoggerKey, slog.Default())
// 		// 		return context.WithValue(ctx, ctxutil.UserIDKey, userID)
// 		// 	},
// 		// 	userID:             userID,
// 		// 	expectedStatusCode: 200,
// 		// 	expectedUser:       user,
// 		// },
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
//
// 			r, _ := http.NewRequest("GET", "/api/v1/user/me", nil)
// 			w := httptest.NewRecorder()
//
// 			// pass userID to context
// 			ctx := tt.setupContext()
// 			r = r.WithContext(ctx)
//
// 			userSvc, userRepo := factories.SetupUser(t, ctx, tt.version)
//
// 			appSvc := &app.Services{User: userSvc}
// 			appRepo := &app.Repos{User: userRepo}
// 			h := app.New(appSvc, appRepo)
// 			handler := userHandler.New(h)
//
// 			handler.GetUser(w, r)
//
// 			// parse body for user
// 			var user *models.User
// 			json.NewDecoder(w.Body).Decode(user)
//
// 			assert.Equal(t, w.Code, tt.expectedStatusCode)
// 			assert.Equal(t, user, tt.expectedUser)
// 		})
// 	}
// }
