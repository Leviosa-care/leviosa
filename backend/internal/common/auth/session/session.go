package session

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

const (
	PendingSessionDuration = 30 * time.Minute // Shorter duration for registration workflow
	ActiveSessionDuration  = 24 * time.Hour   // Standard duration for authenticated sessions
	SessionDuration        = 24 * time.Hour   // Deprecated: use ActiveSessionDuration
)

type SessionState string

const (
	SessionPending SessionState = "pending"
	SessionActive  SessionState = "active"
)

// SessionInfo contains only the session data needed in request context
// This is a lightweight version of Session for passing through middleware
type SessionInfo struct {
	ID     uuid.UUID     `json:"id"`
	UserID uuid.UUID     `json:"user_id"`
	Role   identity.Role `json:"role"`
	State  SessionState  `json:"state"`
}

type Session struct {
	ID          uuid.UUID     `json:"-"`
	UserID      uuid.UUID     `json:"-" encx:"encrypt,hash_basic"`
	Role        identity.Role `json:"-" encx:"encrypt"`
	State       SessionState  `json:"-" encx:"encrypt"`
	CreatedAt   time.Time     `json:"-" encx:"encrypt"`
	ExpiresAt   time.Time     `json:"-" encx:"encrypt"`
	CompletedAt *time.Time    `json:"-" encx:"encrypt"`
	AccessToken string        `json:"-" encx:"hash_basic"`
	RefreshToken string       `json:"-" encx:"hash_basic"`
}

func (s *Session) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}

// TokenPair represents access and refresh tokens with their hashed values
type TokenPair struct {
	AccessToken  string `json:"-" encx:"hash_basic"`
	RefreshToken string `json:"-" encx:"hash_basic"`
}
