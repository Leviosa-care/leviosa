package models

import (
	"context"
	"time"

	"github.com/hengadev/errsx"
)

type UserSignUp struct {
	Email      string    `json:"email" validate:"require"` // Stored hash for searching
	Password   string    `json:"password" validate:"required,min=6"`
	BirthDate  time.Time `json:"birthdate" validate:"require"`
	LastName   string    `json:"lastname" validate:"required"`
	FirstName  string    `json:"firstname" validate:"required"`
	Gender     string    `json:"gender" validate:"required"`
	Telephone  string    `json:"telephone" validate:"required"`
	PostalCode string    `json:"postalCode" validate:"required"`
	City       string    `json:"city" validate:"required"`
	Address1   string    `json:"address1" validate:"required"`
	Address2   string    `json:"address2" validate:"required"`
}

func (u UserSignUp) Valid(ctx context.Context) error {
	var errs = make(errsx.Map)
	if err := ValidateEmail(u.Email); err != nil {
		errs.Set("email", err)
	}
	if err := ValidatePassword(u.Password); err != nil {
		errs.Set("password", err)
	}
	if err := ValidateTelephone(u.Telephone); err != nil {
		errs.Set("telephone", "telephne number should have at leat 10 digits")
	}
	if err := ValidateGender(u.Gender); err != nil {
		errs.Set("gender", err)
	}
	if err := ValidateAddress(u.City, u.PostalCode, u.Address1, u.Address2); err != nil {
		errs.Set("address", err)
	}
	if err := ValidateBirthDate(u.BirthDate); err != nil {
		errs.Set("birthdate", err)
	}
	return errs.AsError()
}

func (user *UserSignUp) ToUser() *User {
	return &User{
		Email:      user.Email,
		Password:   user.Password,
		BirthDate:  user.BirthDate,
		LastName:   user.LastName,
		FirstName:  user.FirstName,
		Gender:     user.Gender,
		Telephone:  user.Telephone,
		PostalCode: user.PostalCode,
		City:       user.City,
		Address1:   user.Address1,
		Address2:   user.Address2,
	}
}
