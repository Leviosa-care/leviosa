package domain

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/validation"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type BuildingResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
	Description string    `json:"description,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
}

type CreateBuildingRequest struct {
	Name        string `json:"name" encx:"encrypt"`
	Address     string `json:"address" encx:"encrypt"`
	City        string `json:"city" encx:"encrypt"`
	PostalCode  string `json:"postal_code" encx:"encrypt"`
	Country     string `json:"country" encx:"encrypt"`
	Description string `json:"description,omitempty" encx:"encrypt"`
	Phone       string `json:"phone,omitempty" encx:"encrypt"`
	Email       string `json:"email,omitempty" encx:"encrypt"`
	IsActive    bool   `json:"is_active"`
}

func (r *CreateBuildingRequest) Valid(ctx context.Context) error {
	// TODO: complete that validation function
	var errs errsx.Map

	if err := validation.ValidateEmail(r.Email); err != nil {
		errs.Set("email", err)
	}
	if err := validation.ValidatePhone(r.Phone); err != nil {
		errs.Set("phone", err)
	}

	return errs.AsError()
}

type UpdateBuildingRequest struct {
	Name        string `json:"name" encx:"encrypt"`
	Address     string `json:"address" encx:"encrypt"`
	City        string `json:"city" encx:"encrypt"`
	PostalCode  string `json:"postal_code" encx:"encrypt"`
	Country     string `json:"country" encx:"encrypt"`
	Description string `json:"description,omitempty" encx:"encrypt"`
	Phone       string `json:"phone,omitempty" encx:"encrypt"`
	Email       string `json:"email,omitempty" encx:"encrypt"`
	IsActive    bool   `json:"is_active"`
}

func (r *UpdateBuildingRequest) Valid(ctx context.Context) error {
	// TODO: complete that validation function
	var errs errsx.Map
	return errs.AsError()
}
