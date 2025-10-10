package domain

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type Specialization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" encx:"encrypt"`
	DisplayName string    `json:"display_name" encx:"encrypt"`
	Description string    `json:"description" encx:"encrypt"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Specialization) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Required fields validation
	if strings.TrimSpace(s.Name) == "" {
		errs.Set("name", "specialization name is required")
	}

	if strings.TrimSpace(s.DisplayName) == "" {
		errs.Set("display_name", "specialization display name is required")
	}

	// Name format validation (should be lowercase, alphanumeric with underscores)
	name := strings.TrimSpace(s.Name)
	if name != "" {
		if !isValidSpecializationName(name) {
			errs.Set("name", "specialization name must be lowercase alphanumeric with underscores only")
		}
	}

	// Display name length validation
	displayName := strings.TrimSpace(s.DisplayName)
	if len(displayName) > 100 {
		errs.Set("display_name", "display name must be 100 characters or less")
	}

	// Description length validation
	description := strings.TrimSpace(s.Description)
	if len(description) > 500 {
		errs.Set("description", "description must be 500 characters or less")
	}

	return errs.AsError()
}

func (s *Specialization) ToResponse() *SpecializationResponse {
	return &SpecializationResponse{
		ID:          s.ID,
		Name:        s.Name,
		DisplayName: s.DisplayName,
		Description: s.Description,
		IsActive:    s.IsActive,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// isValidSpecializationName checks if the name follows the required format
func isValidSpecializationName(name string) bool {
	if len(name) == 0 {
		return false
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}

	return true
}

