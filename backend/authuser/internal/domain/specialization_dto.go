package domain

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type SpecializationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSpecializationRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

func (r *CreateSpecializationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Required fields validation
	if strings.TrimSpace(r.Name) == "" {
		errs.Set("name", "specialization name is required")
	}

	if strings.TrimSpace(r.DisplayName) == "" {
		errs.Set("display_name", "specialization display name is required")
	}

	// Name format validation (should be lowercase, alphanumeric with underscores)
	name := strings.TrimSpace(r.Name)
	if name != "" {
		if !isValidSpecializationName(name) {
			errs.Set("name", "specialization name must be lowercase alphanumeric with underscores only")
		}
	}

	// Display name length validation
	displayName := strings.TrimSpace(r.DisplayName)
	if len(displayName) > 100 {
		errs.Set("display_name", "display name must be 100 characters or less")
	}

	// Description length validation
	description := strings.TrimSpace(r.Description)
	if len(description) > 500 {
		errs.Set("description", "description must be 500 characters or less")
	}

	return errs.AsError()
}

func (r *CreateSpecializationRequest) ToSpecialization() *Specialization {
	return &Specialization{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(r.Name),
		DisplayName: strings.TrimSpace(r.DisplayName),
		Description: strings.TrimSpace(r.Description),
		IsActive:    true, // New specializations are active by default
	}
}

type UpdateSpecializationRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

func (r *UpdateSpecializationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Display name validation if provided
	if r.DisplayName != nil {
		displayName := strings.TrimSpace(*r.DisplayName)
		if displayName == "" {
			errs.Set("display_name", "display name cannot be empty")
		} else if len(displayName) > 100 {
			errs.Set("display_name", "display name must be 100 characters or less")
		}
	}

	// Description validation if provided
	if r.Description != nil {
		description := strings.TrimSpace(*r.Description)
		if len(description) > 500 {
			errs.Set("description", "description must be 500 characters or less")
		}
	}

	return errs.AsError()
}

type GetSpecializationsResponse struct {
	Specializations []SpecializationResponse `json:"specializations"`
	Total           int                      `json:"total"`
}

type AddPartnerSpecializationRequest struct {
	PartnerID        uuid.UUID `json:"partner_id"`
	SpecializationID uuid.UUID `json:"specialization_id"`
}

func (r *AddPartnerSpecializationRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.PartnerID == uuid.Nil {
		errs.Set("partner_id", "partner ID is required")
	}

	if r.SpecializationID == uuid.Nil {
		errs.Set("specialization_id", "specialization ID is required")
	}

	return errs.AsError()
}