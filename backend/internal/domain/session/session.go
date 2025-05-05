package sessionService

import (
	"context"
	"os"
	"time"

	"github.com/hengadev/leviosa/internal/domain/user/models"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

const SessionDuration = 30 * 24 * time.Hour

var SessionName = os.Getenv("SESSION_COOKIE_NAME")

type Session struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Role       models.Role `json:"role"`
	LoggedInAt time.Time   `json:"logged_in_at"`
	CreatedAt  time.Time   `json:"created_at"`
	ExpiresAt  time.Time   `json:"expires_at"`
}

// TODO: change that name for session stored
type Values struct {
	UserID     string      `json:"user_id"`
	Role       models.Role `json:"role"`
	LoggedInAt time.Time   `json:"logged_in_at"`
	CreatedAt  time.Time   `json:"created_at"`
	ExpiresAt  time.Time   `json:"expires_at"`
}

func (s *Session) Values() *Values {
	return &Values{
		UserID:     s.UserID,
		Role:       s.Role,
		LoggedInAt: s.LoggedInAt,
		CreatedAt:  s.CreatedAt,
		ExpiresAt:  s.ExpiresAt,
	}
}

func NewSession(userID string, role models.Role) (*Session, error) {
	if err := uuid.Validate(userID); err != nil {
		return nil, err
	}
	return &Session{
		ID:         uuid.NewString(),
		UserID:     userID,
		Role:       role,
		LoggedInAt: time.Now(),
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(SessionDuration),
	}, nil
}

func (s *Session) Valid(ctx context.Context) error {
	var errs = make(errsx.Map)
	if err := uuid.Validate(s.ID); err != nil {
		errs.Set("id", "session ID is not of type UUID")
	}
	if err := uuid.Validate(s.UserID); err != nil {
		errs.Set("userId", "User ID is not of type UUID")
	}
	if time.Now().Add(SessionDuration).Before(s.ExpiresAt) {
		errs.Set("expiredat", "session expired")
	}
	if s.Role != models.VISITOR {
		errs.Set("role", "got UNKNOWN role, expect one of 'BASIC', 'GUEST', 'FREELANCE', 'ADMINISTRATOR'")
	}
	return errs.AsError()
}

func (s Session) AssertComparable() {}
