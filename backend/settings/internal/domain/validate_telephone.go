package domain

import (
	"errors"
	"regexp"
)

var frenchPhoneStrict = regexp.MustCompile(`^(0[1-5]|06|07)\d{8}$`)

func ValidateTelephone(telephone string) error {
	if !frenchPhoneStrict.MatchString(telephone) {
		return errors.New("invalid telephone number")
	}
	return nil
}
