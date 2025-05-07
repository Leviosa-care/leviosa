package factories

import (
	"time"

	"github.com/hengadev/leviosa/internal/domain/session"
	"github.com/hengadev/leviosa/internal/domain/user/models"
)

func NewBasicSession(overrides map[string]interface{}) *sessionService.Session {
	session := sessionService.Session{
		ID:        "a36bde82-55ae-4fa9-9289-cb82bde68014",
		UserID:    "3d20640f-2df6-4e76-81fe-9f6587fd5980",
		Role:      models.STANDARD,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	for key, value := range overrides {
		switch key {
		case "ID":
			session.ID = value.(string)
		case "UserID":
			session.UserID = value.(string)
		case "Role":
			session.Role = value.(models.Role)
		case "ExpiresAt":
			session.ExpiresAt = value.(time.Time)
		}
	}
	return &session
}

func NewBasicInitSession(overrides map[string]any) map[string]*sessionService.Values {
	session := NewBasicSession(overrides)
	return map[string]*sessionService.Values{
		session.ID: session.Values(),
	}
}
