package domain

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"
	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type CreateSessionRequest struct {
	UserID string        `json:"user_id"`
	Role   identity.Role `json:"role"`
	State  session.SessionState
}

func (r *CreateSessionRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := uuid.Validate(r.UserID); err != nil {
		errs.Set("user ID", err)
	}
	if !r.Role.IsValid() {
		errs.Set("role", "invalid role, role must be 'visitor', 'standard', 'premium', 'guest', 'partner', 'administrator'.")
	}
	if r.State == session.SessionPending && r.Role != identity.Visitor {
		errs.Set("incompatible state and role", "pending user can only be 'visitor'.")
	}

	if r.State == session.SessionActive && r.Role == identity.Visitor {
		errs.Set("incompatible state and role", "active user must have a role different that 'visitor'.")
	}
	return errs.AsError()

}

type CreateSessionResponse struct {
	AccessToken        string    `json:"access_token"`
	RefreshToken       string    `json:"refresh_token"`
	AccessTokenExpiry  time.Time `json:"access_token_expiry"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}

type GetSessionRequest struct {
	Token string `json:"token"`
}

func (r *GetSessionRequest) Valid(ctx context.Context) error {
	return session.ValidateToken(r.Token)
}

type ValidateSessionRequest struct {
	Token string `json:"token"`
}

func (r *ValidateSessionRequest) Valid(ctx context.Context) error {
	return session.ValidateToken(r.Token)
}

type RemoveSessionRequest struct {
	Token string `json:"token"`
}

func (r *RemoveSessionRequest) Valid(ctx context.Context) error {
	return session.ValidateToken(r.Token)
}

type RefreshSessionResponse struct {
	AccessToken        string    `json:"access_token"`
	RefreshToken       string    `json:"refresh_token"`
	AccessTokenExpiry  time.Time `json:"access_token_expiry"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
}
