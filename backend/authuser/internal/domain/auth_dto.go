package domain

import (
	"context"

	"github.com/Leviosa-care/core/auth/session"
	"github.com/Leviosa-care/core/validation"
	"github.com/hengadev/errsx"
)

type CheckEmailAvailabilityRequest struct {
	Email string `json:"email"`
}

func (r *CheckEmailAvailabilityRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := validation.ValidateEmail(r.Email); err != nil {
		errs.Set("email", err)
	}
	return errs.AsError()
}

type ValidateOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (r ValidateOTPRequest) Valid(ctx context.Context, expectedLength int) error {
	var errs errsx.Map

	if err := validation.ValidateEmail(r.Email); err != nil {
		errs.Set("otp email", err)
	}

	if err := ValidateOTP(ctx, r.Code, expectedLength); err != nil {
		errs.Set("otp code", err)
	}

	return errs.AsError()
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SignInRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := validation.ValidateEmail(r.Email); err != nil {
		errs.Set("email", err)
	}
	if err := ValidatePassword(r.Password); err != nil {
		errs.Set("password", err)
	}
	return errs.AsError()
}

type SignOutRequest struct {
	Token string `json:"token"`
}

func (r *SignOutRequest) Valid(ctx context.Context) error {
	return session.ValidateToken(r.Token)
}

type RequestPasswordResetRequest struct {
	Email string `json:"email"`
}

func (r *RequestPasswordResetRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := validation.ValidateEmail(r.Email); err != nil {
		errs.Set("email", err)
	}
	return errs.AsError()
}
