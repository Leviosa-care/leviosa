package models

import (
	"context"
	"reflect"

	"github.com/hengadev/errsx"
)

type UserSignIn struct {
	Email    string `json:"email" validate:"required,email"` // Stored hash for searching
	Password string `json:"password" validate:"required,min=6"`
}

func (u UserSignIn) Valid(ctx context.Context) error {
	var errs = make(errsx.Map)
	vf := reflect.VisibleFields(reflect.TypeOf(u))
	for _, f := range vf {
		switch f.Name {
		case "Email":
			if err := ValidateEmail(u.Email); err != nil {
				errs.Set("email", err)
			}
		case "Password":
			if err := ValidatePassword(u.Password); err != nil {
				errs.Set("password", err)
			}
		default:
			continue
		}
	}
	return errs.AsError()
}
