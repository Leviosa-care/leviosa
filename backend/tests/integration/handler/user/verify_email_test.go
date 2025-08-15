package userHandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hengadev/leviosa/internal/domain/otp"
	"github.com/hengadev/leviosa/internal/domain/user"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/internal/repository/postgres/settings"
	"github.com/hengadev/leviosa/internal/repository/postgres/user"
	"github.com/hengadev/leviosa/internal/repository/redis/otp"
	"github.com/hengadev/leviosa/internal/server/app"
	"github.com/hengadev/leviosa/internal/server/handler/user"
	"github.com/hengadev/leviosa/pkg/testdatabase"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

type Test struct {
	name    string
	version int64
	email   models.Email
}

func TestVerifyEmail(t *testing.T) {
	tests := []Test{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := setup(t, tt)

			w := httptest.NewRecorder()
			body := map[string]string{
				"email": string(tt.email),
			}
			jsonBody, _ := json.Marshal(body)
			r := httptest.NewRequest(http.MethodPost, "/auth/email", bytes.NewReader(jsonBody))
			// NOTE: If you want to add other thing to the request
			// req.Header.Set("Content-Type", "application/json")
			handler.VerifyEmail(w, r)
			// TODO: do the asserts here
		})
	}
}

func setup(t *testing.T, tt Test) *userHandler.AppInstance {
	t.Helper()
	ctx := context.Background()
	// TODO: move that to some helper function
	redisTDB, err := testdatabase.NewRedis(ctx)
	require.NoError(t, err)
	postgresTDB, err := testdatabase.NewPostgres(ctx)
	require.NoError(t, err)
	postgresTDB.PostgresUp(ctx, migrations, tt.version)
	require.NoError(t, err)
	userRepo, err := userRepository.New(ctx, postgresTDB.DB)
	require.NoError(t, err)
	cryptoMock := new(CryptoServiceMock)
	userSvc := userService.New(userRepo, cryptoMock)
	otpRepo := otpRepository.New(ctx, redisTDB.DB)
	settingsRepo, err := settingsRepository.New(ctx, postgresTDB.DB)
	require.NoError(t, err)
	// TODO: change that for the real implementation of the rabbitmq thing
	otpSvc, err := otpService.New(ctx, cryptoMock, otpRepo, settingsRepo, &amqp.Connection{})
	require.NoError(t, err)
	Svcs := app.Services{
		User: userSvc,
		OTP:  otpSvc,
	}
	Repos := app.Repos{
		User:     userRepo,
		Settings: settingsRepo,
		OTP:      otpRepo,
	}
	App := app.New(&Svcs, &Repos)
	return userHandler.New(App)
}
