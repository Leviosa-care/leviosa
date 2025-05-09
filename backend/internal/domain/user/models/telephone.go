package models

import (
	"errors"
	"regexp"
)

var frenchPhoneStrict = regexp.MustCompile(`^(0[1-5]|06|07)\d{8}$`)

type Telephone string

func ValidateTelephone(telephone string) error {
	if !frenchPhoneStrict.MatchString(telephone) {
		return errors.New("invalid telephone number")
	}
	return nil
}

func NewTelephone(telephone string) (Telephone, error) {
	if err := ValidateTelephone(telephone); err != nil {
		return "", err
	}
	return Telephone(telephone), nil
}

func (t Telephone) String() string {
	return string(t)
}
