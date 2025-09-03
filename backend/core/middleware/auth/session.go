package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/Leviosa-care/core/errs"

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
	ID                   uuid.UUID     `json:"-"`
	UserID               uuid.UUID     `json:"-" encx:"encrypt"`
	UserIDEncrypted      []byte        `json:"user_id_encrypted"`
	Role                 identity.Role `json:"-" encx:"encrypt"`
	RoleEncrypted        []byte        `json:"role_encrypted"`
	State                SessionState  `json:"-" encx:"encrypt"`
	StateEncrypted       []byte        `json:"state_encrypted"`
	CreatedAt            time.Time     `json:"-" encx:"encrypt"`
	CreatedAtEncrypted   []byte        `json:"created_at_encrypted"`
	ExpiresAt            time.Time     `json:"-" encx:"encrypt"`
	ExpiresAtEncrypted   []byte        `json:"expires_at_encrypted"`
	CompletedAt          *time.Time    `json:"-" encx:"encrypt"`
	CompletedAtEncrypted []byte        `json:"completed_at_encrypted,omitempty"`
	AccessToken          string        `json:"-" encx:"hash_basic"`
	AccessTokenHash      string        `json:"access_token_hash"`
	RefreshToken         string        `json:"-" encx:"hash_basic"`
	RefreshTokenHash     string        `json:"refresh_token_hash"`
	DEK                  []byte        `json:"-" encx:"encrypt"`
	DEKEncrypted         []byte        `json:"dek_encrypted"`
	KeyVersion           int           `json:"key_version"`
}

func (s *Session) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}

// TokenPair represents access and refresh tokens with their hashed values
type TokenPair struct {
	AccessToken      string `json:"-" encx:"hash_basic"`
	AccessTokenHash  string `json:"access_token_hash"`
	RefreshToken     string `json:"-" encx:"hash_basic"`
	RefreshTokenHash string `json:"refresh_token_hash"`
	DEK              []byte `json:"-" encx:"encrypt"`
	DEKEncrypted     []byte `json:"dek_encrypted"`
	KeyVersion       int    `json:"key_version"`
}

// DecodeSession unmarshals JSON bytes back to a session
// The JSON data should contain plaintext fields for validation
func DecodeSession(data []byte) (*Session, error) {
	// Temporary struct for unmarshaling with correct JSON tags
	var temp struct {
		ID                   uuid.UUID     `json:"id"`
		UserID               uuid.UUID     `json:"user_id"`
		UserIDEncrypted      []byte        `json:"user_id_encrypted"`
		Role                 identity.Role `json:"role"`
		RoleEncrypted        []byte        `json:"role_encrypted"`
		State                SessionState  `json:"state"`
		StateEncrypted       []byte        `json:"state_encrypted"`
		CreatedAt            time.Time     `json:"created_at"`
		CreatedAtEncrypted   []byte        `json:"created_at_encrypted"`
		ExpiresAt            time.Time     `json:"expires_at"`
		ExpiresAtEncrypted   []byte        `json:"expires_at_encrypted"`
		CompletedAt          *time.Time    `json:"completed_at"`
		CompletedAtEncrypted []byte        `json:"completed_at_encrypted,omitempty"`
		AccessTokenHash      string        `json:"access_token_hash"`
		RefreshTokenHash     string        `json:"refresh_token_hash"`
		DEKEncrypted         []byte        `json:"dek_encrypted"`
		KeyVersion           int           `json:"key_version"`
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return nil, errs.NewJSONUnmarshalErr(err)
	}

	// Convert to Session struct
	session := &Session{
		ID:                   temp.ID,
		UserID:               temp.UserID,
		UserIDEncrypted:      temp.UserIDEncrypted,
		Role:                 temp.Role,
		RoleEncrypted:        temp.RoleEncrypted,
		State:                temp.State,
		StateEncrypted:       temp.StateEncrypted,
		CreatedAt:            temp.CreatedAt,
		CreatedAtEncrypted:   temp.CreatedAtEncrypted,
		ExpiresAt:            temp.ExpiresAt,
		ExpiresAtEncrypted:   temp.ExpiresAtEncrypted,
		CompletedAt:          temp.CompletedAt,
		CompletedAtEncrypted: temp.CompletedAtEncrypted,
		AccessTokenHash:      temp.AccessTokenHash,
		RefreshTokenHash:     temp.RefreshTokenHash,
		DEKEncrypted:         temp.DEKEncrypted,
		KeyVersion:           temp.KeyVersion,
	}

	return session, nil
}
