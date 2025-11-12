package domain

import (
	"context"
	"fmt"
	"strings"

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
	IsActive    bool      `json:"is_active"`
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
	var errs errsx.Map

	// Required field validation
	if strings.TrimSpace(r.Name) == "" {
		errs.Set("name", fmt.Errorf("name is required"))
	}
	if strings.TrimSpace(r.Address) == "" {
		errs.Set("address", fmt.Errorf("address is required"))
	}
	if strings.TrimSpace(r.City) == "" {
		errs.Set("city", fmt.Errorf("city is required"))
	}
	if strings.TrimSpace(r.PostalCode) == "" {
		errs.Set("postal_code", fmt.Errorf("postal code is required"))
	}
	if strings.TrimSpace(r.Country) == "" {
		errs.Set("country", fmt.Errorf("country is required"))
	}

	// Optional field validation
	if r.Email != "" {
		if err := validation.ValidateEmail(r.Email); err != nil {
			errs.Set("email", err)
		}
	}
	if r.Phone != "" {
		if err := validation.ValidatePhone(r.Phone); err != nil {
			errs.Set("phone", err)
		}
	}

	return errs.AsError()
}

type UpdateBuildingRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	Address     *string   `json:"address"`
	City        *string   `json:"city"`
	PostalCode  *string   `json:"postal_code"`
	Country     *string   `json:"country"`
	Description *string   `json:"description,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	Email       *string   `json:"email,omitempty"`
	IsActive    *bool     `json:"is_active"`
}

func (r *UpdateBuildingRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// ID validation (required)
	if r.ID == uuid.Nil {
		errs.Set("id", fmt.Errorf("building ID is required"))
	}

	// Name validation (only if provided)
	if r.Name != nil {
		if strings.TrimSpace(*r.Name) == "" {
			errs.Set("name", fmt.Errorf("name cannot be empty"))
		}
	}

	// Address validation (only if provided)
	if r.Address != nil {
		if strings.TrimSpace(*r.Address) == "" {
			errs.Set("address", fmt.Errorf("address cannot be empty"))
		}
	}

	// City validation (only if provided)
	if r.City != nil {
		if strings.TrimSpace(*r.City) == "" {
			errs.Set("city", fmt.Errorf("city cannot be empty"))
		}
	}

	// PostalCode validation (only if provided)
	if r.PostalCode != nil {
		if strings.TrimSpace(*r.PostalCode) == "" {
			errs.Set("postal_code", fmt.Errorf("postal code cannot be empty"))
		}
	}

	// Country validation (only if provided)
	if r.Country != nil {
		if strings.TrimSpace(*r.Country) == "" {
			errs.Set("country", fmt.Errorf("country cannot be empty"))
		}
	}

	// Email validation (only if provided and not empty)
	if r.Email != nil && *r.Email != "" {
		if err := validation.ValidateEmail(*r.Email); err != nil {
			errs.Set("email", err)
		}
	}

	// Phone validation (only if provided and not empty)
	if r.Phone != nil && *r.Phone != "" {
		if err := validation.ValidatePhone(*r.Phone); err != nil {
			errs.Set("phone", err)
		}
	}

	return errs.AsError()
}
