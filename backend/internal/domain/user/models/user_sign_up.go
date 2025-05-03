package models

import (
	"context"
	"reflect"
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
	PostalCode string    `json:"postal_code" validate:"required"`
	City       string    `json:"city" validate:"required"`
	Address1   string    `json:"address1" validate:"required"`
	Address2   string    `json:"address2" validate:"required"`
}

// TODO: complete that function with all the remaining fields
func (u UserSignUp) Valid(ctx context.Context) error {
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
		case "Telephone":
			if err := ValidateTelephone(u.Telephone); err != nil {
				errs.Set("telephone", "telephne number should have at leat 10 digits")
			}
		case "Birthday":
			// parsedDate, err := time.Parse(BirthdayLayout, u.BirthDate)
			// nonValidDate, _ := time.Parse(BirthdayLayout, "01-01-01")
			// if err != nil && parsedDate != nonValidDate {
			// 	pbms.Set("birthday", err)
			// }
		case "Gender":
			if err := ValidateGender(u.Gender); err != nil {
				errs.Set("gender", err)
			}
		default:
			continue
		}
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
