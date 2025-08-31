package auth

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// TODO: find the right value for this
const (
	SessionDuration   = 24 * time.Hour
	SessionCookieName = "leviosa_session_token"
)

type SessionState string

const (
	SessionPending SessionState = "pending"
	SessionActive  SessionState = "active"
)

type Session struct {
	ID                 uuid.UUID     `json:"-"`
	UserID             uuid.UUID     `json:"-" encx:"encrypt"`
	UserIDEncrypted    []byte        `json:"user_id_encrypted"`
	Role               identity.Role `json:"-" encx:"encrypt"`
	RoleEncrypted      []byte        `json:"role_encrypted"`
	State              SessionState  `json:"-" encx:"encrypt"`
	StateEncrypted     []byte        `json:"state_encrypted"`
	CreatedAt          time.Time     `json:"-" encx:"encrypt"`
	CreatedAtEncrypted []byte        `json:"created_at_encrypted"`
	ExpiresAt          time.Time     `json:"-" encx:"encrypt"`
	ExpiresAtEncrypted []byte        `json:"expires_at_encrypted"`
	Token              string        `json:"-" encx:"hash_basic"`
	TokenHash          string        `json:"token_hash"`
	DEK                []byte        `json:"-" encx:"encrypt"`
	DEKEncrypted       []byte        `json:"dek_encrypted"`
	KeyVersion         int           `json:"key_version"`
}

func (s *Session) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}