package models

import (
	"fmt"

	"github.com/hengadev/errsx"
)

type Password string

const passwordMinLength = 8

func ValidatePassword(p string) error {
	var errs errsx.Map
	if len(p) < passwordMinLength {
		errs.Set("password length", fmt.Sprintf("expect at least %d caracter", passwordMinLength))
	}
	return errs.AsError()
}
func NewPassword(p string) (Password, error) {
	var errs errsx.Map
	if err := ValidatePassword(p); err != nil {
		errs.Set("validate password", err)
	}
	return Password(p), errs.AsError()
}

func (p Password) String() string {
	return string(p)
}
