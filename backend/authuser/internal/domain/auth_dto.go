package domain

import (
	"context"

	"github.com/Leviosa-care/core/middleware/auth"
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
