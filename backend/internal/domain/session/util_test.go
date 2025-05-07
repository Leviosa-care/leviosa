package sessionService_test

import (
	"time"

	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/user/models"
	"github.com/hengadev/leviosa/tests/utils"
)

var baseSession = &sessionService.Session{
	ID:        test.GenerateRandomString(16),
	UserID:    "1",
	Role:      models.STANDARD,
	ExpiresAt: time.Now().Add(sessionService.SessionDuration),
}

var initMap = map[string]*sessionService.Values{
	baseSession.ID: baseSession.Values(),
}
