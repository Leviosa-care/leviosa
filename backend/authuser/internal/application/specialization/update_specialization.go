package specialization

import (
	"context"
	"strings"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/Leviosa-care/core/errs"
	"github.com/google/uuid"
)

func (s *SpecializationService) UpdateSpecialization(ctx context.Context, specializationID uuid.UUID, request *domain.UpdateSpecializationRequest) (*domain.SpecializationResponse, error) {
	// Validate request
	if err := request.Valid(ctx); err != nil {
		return nil, errs.ErrInvalidInput
	}

	// Get existing specialization
	specialization, err := s.repo.GetSpecializationByID(ctx, specializationID)
	if err != nil {
		return nil, err
	}

	// Decrypt current data
	if err := s.crypto.Decrypt(ctx, specialization); err != nil {
		return nil, err
	}

	// Update fields if provided
	updated := false
	if request.DisplayName != nil {
		displayName := strings.TrimSpace(*request.DisplayName)
		if displayName != specialization.DisplayName {
			specialization.DisplayName = displayName
			updated = true
		}
	}

	if request.Description != nil {
		description := strings.TrimSpace(*request.Description)
		if description != specialization.Description {
			specialization.Description = description
			updated = true
		}
	}

	if request.IsActive != nil {
		if *request.IsActive != specialization.IsActive {
			specialization.IsActive = *request.IsActive
			updated = true
		}
	}

	// If nothing changed, return current state
	if !updated {
		return specialization.ToResponse(), nil
	}

	// Encrypt updated fields
	if err := s.crypto.Encrypt(ctx, specialization); err != nil {
		return nil, errs.ErrInvalidValue
	}

	// Update in database
	if err := s.repo.UpdateSpecialization(ctx, specialization); err != nil {
		return nil, err
	}

	// Decrypt for response
	if err := s.crypto.Decrypt(ctx, specialization); err != nil {
		return nil, errs.ErrInvalidValue
	}

	return specialization.ToResponse(), nil
}