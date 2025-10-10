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
	specializationEncx, err := s.repo.GetSpecializationByID(ctx, specializationID)
	if err != nil {
		return nil, err
	}

	// Decrypt current data using the new generated function
	specialization, err := domain.DecryptSpecializationEncx(ctx, s.crypto, specializationEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("specialization for update", err)
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

	// Encrypt updated fields using the new generated function
	updatedSpecializationEncx, err := domain.ProcessSpecializationEncx(ctx, s.crypto, specialization)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("specialization for update", err)
	}

	// Update in database
	if err := s.repo.UpdateSpecialization(ctx, updatedSpecializationEncx); err != nil {
		return nil, err
	}

	// Decrypt for response using the new generated function
	responseSpecialization, err := domain.DecryptSpecializationEncx(ctx, s.crypto, updatedSpecializationEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("specialization for response", err)
	}

	return responseSpecialization.ToResponse(), nil
}