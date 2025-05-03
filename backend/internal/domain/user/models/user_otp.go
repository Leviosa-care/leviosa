package models

import (
	"context"

	"github.com/hengadev/errsx"
)

type UserOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (u UserOTP) Valid(ctx context.Context) error {
	var errs errsx.Map
	return errs.AsError()
}
