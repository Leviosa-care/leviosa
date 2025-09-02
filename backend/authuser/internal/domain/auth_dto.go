package domain

import (
	"context"

	"github.com/Leviosa-care/core/middleware/auth"
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
