package models

import (
	"context"

	"github.com/hengadev/errsx"
)

// the user that is send to admin for validation
type UserPending struct {
	EmailHash string `json:"email_hash"`
	LastName  string `json:"lastname"`
	FirstName string `json:"firstname"`
	GoogleID  string `json:"google_id"`
	AppleID   string `json:"apple_id"`
}

// the admin receive this when validating the user
type UserPendingResponse struct {
	Email    string       `json:"email"`
	Role     string       `json:"role"`
	Provider ProviderType `json:"provider"`
}

func (u UserPending) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}

func (u UserPendingResponse) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}
