package domain

import (
	"context"
	"time"

	"github.com/Leviosa-care/core/validation"
	"github.com/hengadev/errsx"
)

type CompleteUserRequest struct {
	Password   string      `json:"password" validate:"required"`
	FirstName  string      `json:"first_name" validate:"required"`
	LastName   string      `json:"last_name" validate:"required"`
	BirthDate  time.Time   `json:"birth_date" validate:"required"`
	Gender     GenderInput `json:"gender" validate:"required"`
	Telephone  string      `json:"telephone" validate:"required"`
	PostalCode string      `json:"postal_code" validate:"required"`
	City       string      `json:"city" validate:"required"`
	Address1   string      `json:"address1" validate:"required"`
	Address2   string      `json:"address2"`
}

func (r *CompleteUserRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate password
	if err := ValidatePassword(r.Password); err != nil {
		errs.Set("password", err)
	}

	// Validate names (basic sanitation)
	if err := validateName(r.FirstName, "first_name"); err != nil {
		errs.Set("first_name", err)
	}

	if err := validateName(r.LastName, "last_name"); err != nil {
		errs.Set("last_name", err)
	}

	// Validate birth date
	if err := validateBirthDate(r.BirthDate); err != nil {
		errs.Set("birth_date", err)
	}

	// Validate gender
	if err := r.Gender.ValidateGender(); err != nil {
		errs.Set("gender", err)
	}

	// Validate telephone
	if err := validation.ValidatePhone(r.Telephone); err != nil {
		errs.Set("telephone", err)
	}

	// Validate postal code
	if err := validatePostalCode(r.PostalCode); err != nil {
		errs.Set("postal_code", err)
	}

	// Validate city
	if err := validateCity(r.City); err != nil {
		errs.Set("city", err)
	}

	// Validate address1
	if err := validateAddress1(r.Address1); err != nil {
		errs.Set("address1", err)
	}

	if err := validateAddress2(r.Address2); err != nil {
		errs.Set("address2", err)
	}

	if err := validateFullAddress(r.City, r.PostalCode, r.Address1, r.Address2); err != nil {
		errs.Set("full_address", err)
	}

	return errs.AsError()
}
