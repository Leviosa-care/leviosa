package factories

//
// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"testing"
// 	"time"
//
// 	sessionService "github.com/hengadev/leviosa/internal/domain/session"
// 	userService "github.com/hengadev/leviosa/internal/domain/user"
// 	userRepository "github.com/hengadev/leviosa/internal/repository/postgres/user"
// 	miniredis "github.com/hengadev/leviosa/internal/repository/redis"
// 	sessionRepository "github.com/hengadev/leviosa/internal/repository/redis/session"
// 	testdb "github.com/hengadev/leviosa/pkg/testdatabase"
// 	test "github.com/hengadev/leviosa/tests/utils"
// )
//
// func SetupUser(t testing.TB, ctx context.Context, version int64) (*userService.Service, *userRepository.Repository) {
// 	t.Helper()
// 	db, err := testdb.NewSQLite(ctx)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if err := testdb.MigrateSQLite(ctx, db, version); err != nil {
// 		t.Error(err)
// 	}
// 	userRepo, _ := userRepository.New(ctx, db)
//
// 	conf := test.PrepareEncryptionConfig()
// 	// TODO: I need a mock here
// 	userService := userService.New(userRepo, conf)
// 	return userService, userRepo
// }
//
// type redisTeardownFunc func()
//
// func SetupSession(t testing.TB, ctx context.Context, initMap miniredis.InitMap[*sessionService.Values]) (*sessionService.Service, *sessionRepository.Repository, redisTeardownFunc) {
// 	t.Helper()
// 	client, err := miniredis.Setup(t, ctx)
// 	if err != nil {
// 		t.Errorf("setup miniredis: %s", err)
// 	}
// 	if err := miniredis.Init(t, ctx, client, sessionRepository.SESSIONPREFIX, initMap); err != nil {
// 		t.Errorf("init miniredis: %s", err)
// 	}
// 	sessionRepo := sessionRepository.New(ctx, client)
// 	sessionService := sessionService.New(sessionRepo)
// 	teardown := func() {
// 		if err := client.Close(); err != nil {
// 			t.Errorf("closing miniredis: %s", err)
// 		}
// 	}
// 	return sessionService, sessionRepo, teardown
// }
//
// func EncodeForBody(t *testing.T, v any) *bytes.Buffer {
// 	t.Helper()
// 	encodedValue, err := json.Marshal(v)
// 	if err != nil {
// 		t.Errorf("encode creds: %s", err)
// 	}
// 	body := bytes.NewBuffer(encodedValue)
// 	return body
// }
//
// func NewBasicCookie(overrides map[string]any) *http.Cookie {
// 	cookie := &http.Cookie{
// 		Name:     sessionService.SessionName,
// 		Value:    NewBasicSession(nil).ID,
// 		Expires:  time.Now().Add(sessionService.SessionDuration),
// 		HttpOnly: true,
// 	}
// 	for key, value := range overrides {
// 		switch key {
// 		case "Name":
// 			cookie.Name = value.(string)
// 		case "Value":
// 			cookie.Value = value.(string)
// 		case "Expires":
// 			cookie.Expires = value.(time.Time)
// 		case "HttpOnly":
// 			cookie.HttpOnly = value.(bool)
// 		}
// 	}
// 	return cookie
// }
