package specialization

import (
	"context"

	"github.com/Leviosa-care/authuser/internal/domain"
	"github.com/google/uuid"
)

func (s *SpecializationService) GetSpecializationByID(ctx context.Context, specializationID uuid.UUID) (*domain.SpecializationResponse, error) {
	// Get from database
	specializationEncx, err := s.repo.GetSpecializationByID(ctx, specializationID)
	if err != nil {
		return nil, err
	}

	// Decrypt sensitive fields using the new generated function
	specialization, err := domain.DecryptSpecializationEncx(ctx, s.crypto, specializationEncx)
	if err != nil {
		return nil, err
	}

	return specialization.ToResponse(), nil
}