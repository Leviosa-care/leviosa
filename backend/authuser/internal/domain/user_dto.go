package domain

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/contracts/identity"
	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	State      UserState `json:"state"`
	Email      string    `json:"email"`
	Picture    string    `json:"picture,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	LoggedInAt time.Time `json:"logged_in_at"`
	Role       string    `json:"role,omitempty"`
	BirthDate  time.Time `json:"birthdate"`
	LastName   string    `json:"last_name,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	Telephone  string    `json:"telephone,omitempty"`
	PostalCode string    `json:"postal_code,omitempty"`
	City       string    `json:"city,omitempty"`
	Address1   string    `json:"address1,omitempty"`
	Address2   string    `json:"address2,omitempty"`
}

type ApproveUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

func (r *ApproveUserRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if _, err := identity.ParseRole(r.Role); err != nil {
		errs.Set("user role", err)
	}
	return errs.AsError()
}

type UpdateUserRoleRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

func (r *UpdateUserRoleRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if _, err := identity.ParseRole(r.Role); err != nil {
		errs.Set("user role", err)
	}
	return errs.AsError()
}

type UpdateUserRequest struct {
	Picture    *string    `json:"picture,omitempty"`
	FirstName  *string    `json:"first_name,omitempty"`
	LastName   *string    `json:"last_name,omitempty"`
	BirthDate  *time.Time `json:"birthdate,omitempty"`
	Gender     *string    `json:"gender,omitempty"`
	Email      *string    `json:"email,omitempty"`
	Telephone  *string    `json:"telephone,omitempty"`
	PostalCode *string    `json:"postal_code,omitempty"`
	City       *string    `json:"city,omitempty"`
	Address1   *string    `json:"address1,omitempty"`
	Address2   *string    `json:"address2,omitempty"`
}
